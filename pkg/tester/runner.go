package tester

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"celang/pkg/compiler"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	_ "celang/pkg/syntax"
	"celang/pkg/vm"
)

// FileResult holds test outcomes for a single test file.
type FileResult struct {
	FilePath string
	Duration time.Duration
	Cases    []evaluator.TestCaseResult
	Error    string
}

// Summary holds the aggregated metrics of an entire test run.
type Summary struct {
	Files      []FileResult
	TotalFiles int
	TotalTests int
	Passed     int
	Failed     int
	Duration   time.Duration
}

// Options configures the test runner execution and reporter.
type Options struct {
	Filter  string
	Verbose bool
	NoColor bool
	Out     io.Writer
	UseAST  bool
}

// IsTestFile determines whether a file path qualifies as a Celang test file.
func IsTestFile(path string) bool {
	clean := filepath.Clean(path)
	base := filepath.Base(clean)
	if !strings.HasSuffix(base, ".ce") {
		return false
	}
	// Direct filename patterns: foo_uji.ce, foo_test.ce, foo.uji.ce, foo.test.ce
	if strings.HasSuffix(base, "_uji.ce") || strings.HasSuffix(base, "_test.ce") ||
		strings.HasSuffix(base, ".uji.ce") || strings.HasSuffix(base, ".test.ce") {
		return true
	}
	// Files inside a dedicated test directory
	dir := filepath.Base(filepath.Dir(clean))
	if dir == "uji" || dir == "tests" || dir == "test" {
		return true
	}
	return false
}

// DiscoverTestFiles scans targets (or the current directory if empty) for Celang test files.
func DiscoverTestFiles(targets []string) ([]string, error) {
	if len(targets) == 0 {
		return scanDir(".")
	}

	var files []string
	seen := make(map[string]bool)

	for _, target := range targets {
		info, err := os.Stat(target)
		if err != nil {
			return nil, fmt.Errorf("target tidak ditemukan: %s", target)
		}

		if info.IsDir() {
			dirFiles, err := scanDir(target)
			if err != nil {
				return nil, err
			}
			for _, f := range dirFiles {
				if !seen[f] {
					seen[f] = true
					files = append(files, f)
				}
			}
		} else {
			clean := filepath.Clean(target)
			if !seen[clean] {
				seen[clean] = true
				files = append(files, clean)
			}
		}
	}

	sort.Strings(files)
	return files, nil
}

func scanDir(root string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := info.Name()
			if path != root && (strings.HasPrefix(base, ".") || base == "node_modules" || base == "dist" ||
				base == "bin" || base == "vendor") {
				return filepath.SkipDir
			}
			return nil
		}
		if IsTestFile(path) {
			clean := filepath.Clean(path)
			if !seen[clean] {
				seen[clean] = true
				files = append(files, clean)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// RunFile executes all tests found in a single Celang test file.
func RunFile(filePath string, opt Options) FileResult {
	start := time.Now()

	sourceBytes, err := os.ReadFile(filePath)
	if err != nil {
		return FileResult{
			FilePath: filePath,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Gagal membaca berkas: %v", err),
		}
	}

	source := string(sourceBytes)
	l := parser.NewLexer(source)
	p := parser.NewParser(l)
	p.IsUniversal = parser.DetectUniversalDirective(source)
	prog, parseErr := p.ParseProgram()
	if parseErr != nil || len(p.Errors()) > 0 {
		var errs []string
		if parseErr != nil {
			errs = append(errs, parseErr.Error())
		}
		errs = append(errs, p.Errors()...)
		return FileResult{
			FilePath: filePath,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("Kesalahan sintaksis (%d galat):\n  • %s", len(errs), strings.Join(errs, "\n  • ")),
		}
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filePath
	env.ModuleRoot = filepath.Dir(filePath)
	env.IsUniversal = p.IsUniversal

	var stdout io.Writer = io.Discard
	if opt.Verbose {
		stdout = opt.Out
	}
	evaluator.RegisterBuiltinsWithIO(env, os.Stdin, stdout)

	tc := &evaluator.TestContext{Filter: opt.Filter}
	env.SetTestContext(tc)

	var runErr error
	if opt.UseAST {
		evalRes := evaluator.Eval(prog, env)
		if evalRes != nil && (evalRes.Type() == evaluator.ERROR_OBJ || evalRes.Type() == evaluator.ASSERTION_ERROR_OBJ) {
			runErr = fmt.Errorf("%s", evalRes.Inspect())
		}
	} else {
		comp := compiler.NewCompilerWithEnv(env)
		chunk, err := comp.Compile(prog)
		if err != nil {
			// Fall back to AST evaluator if compilation fails on non-bytecode constructs
			evalRes := evaluator.Eval(prog, env)
			if evalRes != nil && (evalRes.Type() == evaluator.ERROR_OBJ || evalRes.Type() == evaluator.ASSERTION_ERROR_OBJ) {
				runErr = fmt.Errorf("%s", evalRes.Inspect())
			}
		} else {
			virtualMachine := vm.NewWithEnv(chunk, env)
			_, err = virtualMachine.Run()
			if err != nil {
				runErr = err
			}
		}
	}
	dur := time.Since(start)

	if runErr != nil {
		if len(tc.Results) == 0 {
			return FileResult{
				FilePath: filePath,
				Duration: dur,
				Error:    runErr.Error(),
			}
		}
	}

	return FileResult{
		FilePath: filePath,
		Duration: dur,
		Cases:    tc.Results,
	}
}

// Run executes all test files, reports progress to opt.Out, and returns a Summary.
func Run(targets []string, opt Options) (*Summary, error) {
	if opt.Out == nil {
		opt.Out = os.Stdout
	}
	if os.Getenv("NO_COLOR") != "" {
		opt.NoColor = true
	}

	files, err := DiscoverTestFiles(targets)
	if err != nil {
		return nil, err
	}

	cReset := "\033[0m"
	cBold := "\033[1m"
	cGreen := "\033[32m"
	cRed := "\033[31m"
	cGray := "\033[90m"
	cCyan := "\033[36m"
	cYellow := "\033[33m"

	if opt.NoColor {
		cReset = ""
		cBold = ""
		cGreen = ""
		cRed = ""
		cGray = ""
		cCyan = ""
		cYellow = ""
	}

	if len(files) == 0 {
		fmt.Fprintf(opt.Out, "%sℹ Tidak ada berkas pengujian (.ce) yang ditemukan.%s\n", cYellow, cReset)
		return &Summary{}, nil
	}

	fmt.Fprintf(opt.Out, "\n%s%s🔍 Menjalankan Celang Test Runner (%d berkas)%s\n\n", cBold, cCyan, len(files), cReset)

	summary := &Summary{
		TotalFiles: len(files),
	}
	suiteStart := time.Now()

	for _, file := range files {
		res := RunFile(file, opt)
		summary.Files = append(summary.Files, res)

		relPath, relErr := filepath.Rel(".", file)
		if relErr != nil {
			relPath = file
		}

		if res.Error != "" {
			summary.Failed++
			summary.TotalTests++
			fmt.Fprintf(opt.Out, "%s✕ %s%s %s(%s)%s\n", cRed+cBold, relPath, cReset, cGray, formatDuration(res.Duration), cReset)
			fmt.Fprintf(opt.Out, "  %s%s%s\n\n", cRed, indentLines(res.Error, "  "), cReset)
			continue
		}

		filePassed := true
		for _, c := range res.Cases {
			summary.TotalTests++
			if c.Passed {
				summary.Passed++
			} else {
				summary.Failed++
				filePassed = false
			}
		}

		if len(res.Cases) == 0 {
			fmt.Fprintf(opt.Out, "%s• %s%s %s(tidak ada blok pengujian)%s\n", cGray, relPath, cReset, cGray, cReset)
			continue
		}

		if filePassed {
			fmt.Fprintf(opt.Out, "%s✓ %s%s %s(%s)%s\n", cGreen+cBold, relPath, cReset, cGray, formatDuration(res.Duration), cReset)
		} else {
			fmt.Fprintf(opt.Out, "%s✕ %s%s %s(%s)%s\n", cRed+cBold, relPath, cReset, cGray, formatDuration(res.Duration), cReset)
		}

		for _, c := range res.Cases {
			if c.Passed {
				fmt.Fprintf(opt.Out, "  %s✓%s %s %s(%s)%s\n", cGreen, cReset, c.Name, cGray, formatDuration(c.Duration), cReset)
			} else {
				fmt.Fprintf(opt.Out, "  %s✕ %s%s %s(%s)%s\n", cRed+cBold, c.Name, cReset, cGray, formatDuration(c.Duration), cReset)
				if c.Pos != nil {
					fmt.Fprintf(opt.Out, "    %sbaris %d, kolom %d:%s %s\n", cYellow, c.Pos.Line, c.Pos.Column, cReset, c.Error)
				} else {
					fmt.Fprintf(opt.Out, "    %s%s%s\n", cRed, c.Error, cReset)
				}
				if c.Expected != "" || c.Actual != "" {
					fmt.Fprintf(opt.Out, "    %sEkspektasi: %s%s\n", cGreen, c.Expected, cReset)
					fmt.Fprintf(opt.Out, "    %sAktual:     %s%s\n", cRed, c.Actual, cReset)
				}
				fmt.Fprintln(opt.Out)
			}
		}
	}

	summary.Duration = time.Since(suiteStart)

	// Summary Footer
	fmt.Fprintf(opt.Out, "\n%s----------------------------------------------------------------------%s\n", cGray, cReset)
	if summary.Failed == 0 {
		fmt.Fprintf(opt.Out, "%s%s✓ SUKSES:%s %d lulus dari %d pengujian (%d berkas) [%s]\n",
			cGreen, cBold, cReset, summary.Passed, summary.TotalTests, summary.TotalFiles, formatDuration(summary.Duration))
	} else {
		fmt.Fprintf(opt.Out, "%s%s✕ GAGAL:%s %d gagal, %d lulus dari %d pengujian (%d berkas) [%s]\n",
			cRed, cBold, cReset, summary.Failed, summary.Passed, summary.TotalTests, summary.TotalFiles, formatDuration(summary.Duration))
	}
	fmt.Fprintf(opt.Out, "%s----------------------------------------------------------------------%s\n\n", cGray, cReset)

	return summary, nil
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.1fµs", float64(d.Nanoseconds())/1000.0)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func indentLines(str, indent string) string {
	lines := strings.Split(str, "\n")
	for i, l := range lines {
		lines[i] = indent + l
	}
	return strings.Join(lines, "\n")
}
