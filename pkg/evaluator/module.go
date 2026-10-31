package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"celang/pkg/parser"
	_ "celang/pkg/syntax"
	"celang/pkg/uast"
)

// findModuleRoot searches upward from a file's directory for a module marker (celang.mod or modul.ce).
// If no marker is found, it falls back to the file's parent directory.
func findModuleRoot(filePath string) string {
	if filePath == "" {
		wd, _ := os.Getwd()
		return wd
	}

	dir := filepath.Dir(filePath)
	current := dir

	for {
		// Check for marker file (Option B)
		if _, err := os.Stat(filepath.Join(current, "celang.mod")); err == nil {
			return current
		}
		if _, err := os.Stat(filepath.Join(current, "modul.ce")); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current || parent == "." || parent == "/" {
			break
		}
		current = parent
	}

	return dir
}

// resolvePackagePath searches for a package in celang_modules directories ascending from baseDir.
func resolvePackagePath(baseDir, packageName string) (string, error) {
	current := baseDir
	for {
		pkgDir := filepath.Join(current, "celang_modules", packageName)
		candidates := []string{
			filepath.Join(pkgDir, "utama.ce"),
			filepath.Join(pkgDir, "main.ce"),
			filepath.Join(pkgDir, "modul.ce"),
			filepath.Join(pkgDir, packageName+".ce"),
			filepath.Join(pkgDir, "index.ce"),
			filepath.Join(current, "celang_modules", packageName+".ce"),
		}
		for _, c := range candidates {
			if info, err := os.Stat(c); err == nil && !info.IsDir() {
				return filepath.Abs(c)
			}
		}

		parent := filepath.Dir(current)
		if parent == current || parent == "." || parent == "/" {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("paket '%s' tidak ditemukan di 'celang_modules'", packageName)
}

// resolveDirectoryEntrypoint searches a directory for common module entrypoints.
func resolveDirectoryEntrypoint(dirPath string) (string, error) {
	baseName := filepath.Base(dirPath)
	candidates := []string{
		filepath.Join(dirPath, "modul.ce"),
		filepath.Join(dirPath, "utama.ce"),
		filepath.Join(dirPath, "index.ce"),
		filepath.Join(dirPath, baseName+".ce"),
		filepath.Join(dirPath, "modul.ast"),
		filepath.Join(dirPath, "utama.ast"),
		filepath.Join(dirPath, "index.ast"),
		filepath.Join(dirPath, baseName+".ast"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return filepath.Abs(c)
		}
	}
	return "", fmt.Errorf("folder '%s' tidak memiliki berkas entrypoint (modul.ce, utama.ce, index.ce, %s.ce)", dirPath, baseName)
}

// resolveModulePath resolves an import path relative to the caller's file.
func resolveModulePath(callerFilePath, importSource string) (string, error) {
	if isStdlibModule(importSource) {
		return "stdlib:" + normalizeStdlibName(importSource), nil
	}

	baseDir := ""
	if callerFilePath != "" {
		baseDir = filepath.Dir(callerFilePath)
	} else {
		baseDir, _ = os.Getwd()
	}

	// If importSource is a package name (not a relative/absolute path and not ending in .ce/.ast),
	// search celang_modules first
	if !strings.HasPrefix(importSource, ".") && !filepath.IsAbs(importSource) &&
		!strings.HasSuffix(importSource, ".ce") && !strings.HasSuffix(importSource, ".ast") {
		if pkgPath, err := resolvePackagePath(baseDir, importSource); err == nil {
			return pkgPath, nil
		}
	}

	candidate := importSource
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(baseDir, candidate)
	}

	// 1. Check exact path (file or directory)
	if info, err := os.Stat(candidate); err == nil {
		if !info.IsDir() {
			return filepath.Abs(candidate)
		}
		// Candidate is a directory: resolve entrypoint file
		if entrypoint, err := resolveDirectoryEntrypoint(candidate); err == nil {
			return entrypoint, nil
		}
	}

	// 2. Try with .ce extension
	if info, err := os.Stat(candidate + ".ce"); err == nil && !info.IsDir() {
		return filepath.Abs(candidate + ".ce")
	}

	// 3. Try with .ast extension
	if info, err := os.Stat(candidate + ".ast"); err == nil && !info.IsDir() {
		return filepath.Abs(candidate + ".ast")
	}

	return "", fmt.Errorf("modul '%s' tidak ditemukan (dicari di: '%s')", importSource, candidate)
}

// loadModule loads and evaluates a module, caching the resulting environment.
func loadModule(targetAbsPath string, env *Environment) (*Environment, error) {
	if strings.HasPrefix(targetAbsPath, "stdlib:") {
		canonName := strings.TrimPrefix(targetAbsPath, "stdlib:")
		return getStdlibModuleEnvironment(canonName)
	}

	if env.ModuleRegistry == nil {
		env.ModuleRegistry = NewModuleRegistry()
	}

	env.ModuleRegistry.mu.RLock()
	cached, ok := env.ModuleRegistry.modules[targetAbsPath]
	env.ModuleRegistry.mu.RUnlock()
	if ok {
		return cached, nil
	}

	data, err := os.ReadFile(targetAbsPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca berkas modul '%s': %w", targetAbsPath, err)
	}

	var prog *uast.Program
	if filepath.Ext(targetAbsPath) == ".ast" {
		prog, err = uast.FromJSON(data)
		if err != nil {
			return nil, fmt.Errorf("gagal deserialisasi U-AST modul '%s': %w", targetAbsPath, err)
		}
	} else {
		l := parser.NewLexer(string(data))
		p := parser.NewParser(l)
		prog, err = p.ParseProgram()
		if err != nil {
			return nil, fmt.Errorf("gagal parsing modul '%s': %w", targetAbsPath, err)
		}
	}

	targetEnv := NewEnvironment()
	targetEnv.FilePath = targetAbsPath
	targetEnv.ModuleRoot = findModuleRoot(targetAbsPath)
	targetEnv.ModuleRegistry = env.ModuleRegistry
	targetEnv.IsUniversal = prog.IsUniversal

	// Cache early to support recursive or mutual imports
	env.ModuleRegistry.mu.Lock()
	env.ModuleRegistry.modules[targetAbsPath] = targetEnv
	env.ModuleRegistry.mu.Unlock()

	result := Eval(prog, targetEnv)
	if isError(result) {
		env.ModuleRegistry.mu.Lock()
		delete(env.ModuleRegistry.modules, targetAbsPath)
		env.ModuleRegistry.mu.Unlock()
		return nil, fmt.Errorf("kesalahan saat mengevaluasi modul '%s': %s", filepath.Base(targetAbsPath), result.Inspect())
	}

	return targetEnv, nil
}

// checkVisibilityAccess verifies if caller is permitted to access a target symbol.
func checkVisibilityAccess(callerFile, callerModuleRoot, targetFile, targetModuleRoot, symbolName string, vis *uast.Visibility) error {
	if vis == nil || vis.Kind == uast.VisibilityPrivate {
		return fmt.Errorf("Simbol '%s' bersifat privat untuk berkas '%s'. Tambahkan 'publik!' atau 'privat-modul!' jika ingin mengeksposnya.", symbolName, filepath.Base(targetFile))
	}

	if vis.Kind == uast.VisibilityPublic {
		return nil
	}

	if vis.Kind == uast.VisibilityModule {
		callerRoot := callerModuleRoot
		if callerRoot == "" && callerFile != "" {
			callerRoot = findModuleRoot(callerFile)
		}
		targetRoot := targetModuleRoot
		if targetRoot == "" && targetFile != "" {
			targetRoot = findModuleRoot(targetFile)
		}

		cleanCaller, _ := filepath.Abs(callerRoot)
		cleanTarget, _ := filepath.Abs(targetRoot)

		if cleanCaller != cleanTarget {
			return fmt.Errorf("Simbol '%s' di '%s' memiliki visibilitas 'privat-modul!' dan hanya dapat diakses di dalam modul yang sama ('%s'). Berkas pemanggil: '%s'.",
				symbolName, filepath.Base(targetFile), filepath.Base(cleanTarget), filepath.Base(callerFile))
		}
		return nil
	}

	if vis.Kind == uast.VisibilityAllowlist {
		targetDir := filepath.Dir(targetFile)
		callerAbs, _ := filepath.Abs(callerFile)

		for _, t := range vis.Targets {
			resolved := t
			if !filepath.IsAbs(resolved) {
				resolved = filepath.Clean(filepath.Join(targetDir, t))
			}
			resolvedAbs, _ := filepath.Abs(resolved)

			if callerAbs == resolvedAbs || filepath.Dir(callerAbs) == resolvedAbs {
				return nil
			}
		}

		return fmt.Errorf("Simbol '%s' di '%s' memiliki visibilitas privat-%v! dan tidak mengizinkan akses dari '%s'.",
			symbolName, filepath.Base(targetFile), vis.Targets, filepath.Base(callerFile))
	}

	return nil
}
