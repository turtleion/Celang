package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"celang/pkg/compiler"
	"celang/pkg/empathy"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/repl"
	"celang/pkg/syntax"
	"celang/pkg/tester"
	"celang/pkg/uast"
	"celang/pkg/vm"
)

const (
	AppVersion       = "1.0.0"
	bundleMagicBytes = "CELANG_BUNDLE_OK"
	trailerSize      = 24
)

func main() {
	// 1. Periksa apakah biner ini adalah aplikasi mandiri tersemat (standalone bundle)
	if tryRunEmbedded() {
		return
	}

	if len(os.Args) < 2 {
		// Default: launch interactive REPL
		if err := repl.Start(repl.Options{}); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Kesalahan REPL: %v\n", err)
			os.Exit(1)
		}
		return
	}

	command := os.Args[1]
	switch command {
	case "bangun", "build":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang bangun <file.ce> [-o nama_biner] [--target=windows|linux|darwin] [--runner=path] [--style=nama|path]")
			fmt.Println("   Alias:     celang build  <file.ce> [-o nama_biner] [--target=windows|linux|darwin] [--runner=path] [--style=nama|path]")
			os.Exit(1)
		}
		inputFile := ""
		outputFile := ""
		targetOS := ""
		runnerPath := ""
		styleName := ""
		useAST := false
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if (arg == "-o" || arg == "--output") && i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			} else if strings.HasPrefix(arg, "--target=") {
				targetOS = strings.ToLower(strings.TrimPrefix(arg, "--target="))
			} else if (arg == "-t" || arg == "--target") && i+1 < len(os.Args) {
				targetOS = strings.ToLower(os.Args[i+1])
				i++
			} else if strings.HasPrefix(arg, "--runner=") {
				runnerPath = strings.TrimPrefix(arg, "--runner=")
			} else if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if arg == "--ast" {
				useAST = true
			} else if inputFile == "" && !strings.HasPrefix(arg, "-") {
				inputFile = arg
			}
		}
		if inputFile == "" {
			fmt.Println("❌ Mohon tentukan berkas kode (.ce) yang ingin dibangun.")
			os.Exit(1)
		}
		buildFile(inputFile, outputFile, targetOS, runnerPath, useAST, styleName)

	case "repl":
		styleName := "indonesian"
		for _, arg := range os.Args[2:] {
			if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			}
		}
		if err := repl.Start(repl.Options{Style: styleName}); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Kesalahan REPL: %v\n", err)
			os.Exit(1)
		}

	case "run":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang run <file.ce | file.ast> [--style=nama|path] [--ast] [-- argumen1 argumen2 ...]")
			os.Exit(1)
		}
		filePath := ""
		styleName := ""
		useAST := false
		var scriptArgs []string

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "--" {
				scriptArgs = os.Args[i+1:]
				break
			}
			if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if arg == "--ast" || arg == "--legacy" || arg == "--tree-walk" {
				useAST = true
			} else if arg == "--vm" {
				useAST = false // explicit VM confirmation
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}

		if filePath == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin dijalankan.")
			os.Exit(1)
		}
		runFile(filePath, styleName, scriptArgs, useAST)

	case "disasm", "disassemble":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang disasm <file.ce> [--style=nama|path]")
			os.Exit(1)
		}
		filePath := ""
		styleName := ""
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}
		if filePath == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin dibongkar.")
			os.Exit(1)
		}
		disasmFile(filePath, styleName)

	case "compile", "parse":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang compile <file.ce> [-o output.ast] [--style=nama|path]")
			os.Exit(1)
		}
		inputFile := ""
		outputFile := ""
		styleName := ""
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if (arg == "-o" || arg == "--output") && i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			} else if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if inputFile == "" && !strings.HasPrefix(arg, "-") {
				inputFile = arg
			}
		}
		if inputFile == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin dikompilasi.")
			os.Exit(1)
		}
		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, ".ce") + ".ast"
		}
		compileFile(inputFile, outputFile, styleName)

	case "format":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang format [-w | --write] <file.ce> [--style=nama_style | path/ke/style.style] [--from-style=...]")
			os.Exit(1)
		}
		isWrite := false
		styleName := "indonesian"
		fromStyle := ""
		filePath := ""

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "-w" || arg == "--write" {
				isWrite = true
			} else if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if strings.HasPrefix(arg, "--from-style=") {
				fromStyle = strings.TrimPrefix(arg, "--from-style=")
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}

		if filePath == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin diformat.")
			os.Exit(1)
		}
		formatFile(filePath, styleName, isWrite, fromStyle)

	case "init":
		projectName := "projek_celang"
		if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "-") {
			projectName = os.Args[2]
		}
		initProject(projectName)

	case "ast", "uast":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang ast <file.ce> [--style=nama|path]")
			os.Exit(1)
		}
		filePath := ""
		styleName := ""
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if (arg == "-s" || arg == "--style") && i+1 < len(os.Args) {
				styleName = os.Args[i+1]
				i++
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}
		if filePath == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin dicetak AST-nya.")
			os.Exit(1)
		}
		dumpAST(filePath, styleName)

	case "uji", "test":
		var targets []string
		filter := ""
		verbose := false
		noColor := false
		useAST := false

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if (arg == "-f" || arg == "--filter" || arg == "--saring") && i+1 < len(os.Args) {
				filter = os.Args[i+1]
				i++
			} else if strings.HasPrefix(arg, "--filter=") {
				filter = strings.TrimPrefix(arg, "--filter=")
			} else if strings.HasPrefix(arg, "--saring=") {
				filter = strings.TrimPrefix(arg, "--saring=")
			} else if arg == "-v" || arg == "--verbose" {
				verbose = true
			} else if arg == "--no-color" {
				noColor = true
			} else if arg == "--ast" {
				useAST = true
			} else if arg == "--vm" {
				useAST = false
			} else if !strings.HasPrefix(arg, "-") {
				targets = append(targets, arg)
			}
		}

		summary, err := tester.Run(targets, tester.Options{
			Filter:  filter,
			Verbose: verbose,
			NoColor: noColor,
			UseAST:  useAST,
			Out:     os.Stdout,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Kesalahan menjalankan pengujian: %v\n", err)
			os.Exit(1)
		}
		if summary.Failed > 0 {
			os.Exit(1)
		}

	case "versi", "version", "--version", "-v":
		printVersion()

	case "help", "bantuan", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Perintah tidak dikenal: '%s'\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("🌟 Celang Runtime v%s (%s/%s)\n", AppVersion, runtime.GOOS, runtime.GOARCH)
	fmt.Println("   Sistem Bahasa Pemrograman Human-Centric & Syntax-Styling")
	fmt.Println("   Bytecode Compiler & Virtual Machine (VM) Terintegrasi")
	fmt.Println("   Lisensi: MIT")
}

func printHelp() {
	fmt.Println("🌟 Celang - Bahasa Pemrograman 'Syntax Styling' & Human-Centric")
	fmt.Println("\nPerintah yang tersedia:")
	fmt.Println("  celang                                  Masuk ke Interactive REPL (otomatis)")
	fmt.Println("  celang repl [--style=nama|path]         Membuka sesi Interactive REPL eksplisit")
	fmt.Println("  celang run <file> [--ast] [-- args...]  Menjalankan kode (.ce/.ast) via Bytecode VM (default, tercepat) atau AST Evaluator (--ast)")
	fmt.Println("  celang bangun <file> [-o out] [--ast]   Membangun aplikasi mandiri (standalone executable)")
	fmt.Println("  celang build  <file> [-o out] [--ast]   Alias internasional untuk 'bangun'")
	fmt.Println("  celang disasm <file.ce>                 Membongkar kode menjadi instruksi rakitan Bytecode")
	fmt.Println("  celang format [-w] <file> [--style=...] Memformat kode (-w untuk langsung menimpa file)")
	fmt.Println("  celang uji [target...] [-f filter] [--ast] Menjalankan pengujian otomatis via VM (default) atau AST")
	fmt.Println("  celang test [target...] [-f filter] [--ast] Alias internasional untuk 'uji'")
	fmt.Println("  celang compile <file.ce> [-o out.ast]   Mengompilasi kode teks menjadi berkas U-AST JSON")
	fmt.Println("  celang ast <file.ce>                    Mencetak pohon struktur U-AST JSON ke terminal")
	fmt.Println("  celang init [nama_proyek]               Menginisialisasi proyek baru (celang.mod & utama.ce)")
	fmt.Println("  celang versi, --version, -v             Melihat nomor versi Celang")
	fmt.Println("  celang bantuan, --help, -h              Menampilkan bantuan ini")
}

func loadProgram(filePath string, styleOpt ...string) (*uast.Program, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("gagal membaca file '%s': %w", filePath, err)
	}

	sourceStr := string(data)

	// If file ends with .ast or .json, attempt JSON deserialize
	if strings.HasSuffix(filePath, ".ast") || strings.HasSuffix(filePath, ".json") {
		prog, err := uast.FromJSON(data)
		if err == nil {
			return prog, sourceStr, nil
		}
	}

	// Otherwise, parse as source text
	l := parser.NewLexer(sourceStr)
	styleName := ""
	if len(styleOpt) > 0 && styleOpt[0] != "" {
		styleName = styleOpt[0]
	}
	if styleName != "" {
		if prof := syntax.GetProfile(styleName); prof != nil {
			prof.ApplyToLexer(l)
		}
	}
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return nil, sourceStr, err
	}

	return prog, sourceStr, nil
}

func handleSyntaxError(err error, sourceStr string) {
	diag := empathy.NewSyntaxDiagnosticFromError(err, sourceStr)
	fmt.Print(diag.Render())
	os.Exit(1)
}

func runFile(filePath string, styleName string, scriptArgs []string, useAST bool) {
	prog, sourceStr, err := loadProgram(filePath, styleName)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	env := evaluator.NewEnvironment()
	if prog.IsUniversal {
		env.IsUniversal = true
	}
	absPath, _ := filepath.Abs(filePath)
	env.FilePath = absPath
	evaluator.RegisterBuiltins(env, os.Stdout)

	// Inject CLI arguments into global 'argumen' and 'args'
	evaluator.SetScriptArgs(scriptArgs)
	argElements := make([]evaluator.Object, len(scriptArgs))
	for i, arg := range scriptArgs {
		argElements[i] = &evaluator.String{Value: arg}
	}
	argsArray := &evaluator.Array{Elements: argElements}
	env.Define("argumen", argsArray, false, "Array")
	env.Define("args", argsArray, false, "Array")

	if useAST {
		result := evaluator.Eval(prog, env)
		if result != nil && result.Type() == evaluator.ERROR_OBJ {
			diag := empathy.NewRuntimeDiagnostic(result.(*evaluator.Error), sourceStr)
			fmt.Print(diag.Render())
			os.Exit(1)
		}
		return
	}

	// Default: Run via Bytecode Compiler & Virtual Machine (VM)
	comp := compiler.NewCompilerWithEnv(env)
	chunk, err := comp.Compile(prog)
	if err != nil {
		diag := empathy.NewSyntaxDiagnostic(err.Error(), 1, 1, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}
	virtualMachine := vm.NewWithEnv(chunk, env)
	result, err := virtualMachine.Run()
	if err != nil {
		var evalErr *evaluator.Error
		if ee, ok := err.(*evaluator.Error); ok {
			evalErr = ee
		} else {
			evalErr = &evaluator.Error{Message: err.Error()}
		}
		diag := empathy.NewRuntimeDiagnostic(evalErr, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}
	_ = result
}

func disasmFile(filePath string, styleOpt ...string) {
	prog, sourceStr, err := loadProgram(filePath, styleOpt...)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		fmt.Printf("❌ Gagal mengompilasi ke bytecode: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(chunk.Disassemble(filepath.Base(filePath)))
}

func compileFile(inputFile, outputFile string, styleOpt ...string) {
	prog, sourceStr, err := loadProgram(inputFile, styleOpt...)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	jsonData, err := uast.ToJSON(prog, true)
	if err != nil {
		fmt.Printf("❌ Gagal serialisasi ke U-AST: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		fmt.Printf("❌ Gagal menulis file output '%s': %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Berhasil mengompilasi '%s' -> '%s' (U-AST)\n", inputFile, outputFile)
}

func formatFile(filePath, styleName string, isWrite bool, fromStyleOpt ...string) {
	fromStyle := ""
	if len(fromStyleOpt) > 0 && fromStyleOpt[0] != "" {
		fromStyle = fromStyleOpt[0]
	} else if styleName != "indonesian" {
		fromStyle = styleName
	}
	prog, sourceStr, err := loadProgram(filePath, fromStyle)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	if styleName == "indonesian" {
		if pragma := detectPragmaStyle(sourceStr); pragma != "" {
			styleName = pragma
		}
	}

	profile := syntax.GetProfile(styleName)
	formatter := syntax.NewFormatter(profile)
	formatted := formatter.Format(prog)

	if isWrite {
		if err := os.WriteFile(filePath, []byte(formatted), 0644); err != nil {
			fmt.Printf("❌ Gagal menimpa berkas '%s': %v\n", filePath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Berkas '%s' berhasil diformat ulang secara langsung (gaya: %s)\n", filePath, profile.DisplayName)
	} else {
		fmt.Printf("--- Menampilkan Kode dengan Gaya: %s ---\n", profile.DisplayName)
		fmt.Print(formatted)
	}
}

func initProject(name string) {
	targetDir := name
	if targetDir != "." {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Printf("❌ Gagal membuat direktori proyek '%s': %v\n", targetDir, err)
			os.Exit(1)
		}
	}

	modContent := fmt.Sprintf("modul %s\nversi 1.0.0\n", filepath.Base(name))
	modPath := filepath.Join(targetDir, "celang.mod")
	if err := os.WriteFile(modPath, []byte(modContent), 0644); err != nil {
		fmt.Printf("❌ Gagal membuat file 'celang.mod': %v\n", err)
		os.Exit(1)
	}

	mainContent := fmt.Sprintf(`// Entrypoint utama aplikasi Celang
fungsi utama() {
    cetak("Halo dari Celang! Selamat datang di %s 🚀")
}

utama()
`, filepath.Base(name))
	mainPath := filepath.Join(targetDir, "utama.ce")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		fmt.Printf("❌ Gagal membuat file 'utama.ce': %v\n", err)
		os.Exit(1)
	}

	readmeContent := fmt.Sprintf(`# %s

Aplikasi berbasis bahasa pemrograman [Celang](https://github.com/celang).

## Cara Menjalankan
`+"```bash\ncelang run utama.ce\n```\n", filepath.Base(name))
	readmePath := filepath.Join(targetDir, "README.md")
	_ = os.WriteFile(readmePath, []byte(readmeContent), 0644)

	fmt.Printf("🎉 Proyek Celang '%s' berhasil diinisialisasi!\n", name)
	fmt.Printf("   📁 Direktori : %s\n", targetDir)
	fmt.Printf("   📄 Berkas    : celang.mod, utama.ce, README.md\n")
	fmt.Println("\nJalankan proyekmu dengan:")
	if targetDir != "." {
		fmt.Printf("   cd %s && celang run utama.ce\n", targetDir)
	} else {
		fmt.Println("   celang run utama.ce")
	}
}

func dumpAST(filePath string, styleOpt ...string) {
	prog, sourceStr, err := loadProgram(filePath, styleOpt...)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	jsonData, err := uast.ToJSON(prog, true)
	if err != nil {
		fmt.Printf("❌ Gagal mengonversi AST ke JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func detectPragmaStyle(sourceStr string) string {
	lines := strings.Split(sourceStr, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "// @gaya:") || strings.HasPrefix(trimmed, "// @style:") ||
			strings.HasPrefix(trimmed, "# @gaya:") || strings.HasPrefix(trimmed, "# @style:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			break
		}
	}
	return ""
}

// tryRunEmbedded memeriksa apakah biner eksekusi saat ini memiliki muatan kode Celang yang tersemat di ekornya.
func tryRunEmbedded() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return false
	}

	f, err := os.Open(exePath)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil || stat.Size() < trailerSize {
		return false
	}

	trailer := make([]byte, trailerSize)
	_, err = f.ReadAt(trailer, stat.Size()-trailerSize)
	if err != nil {
		return false
	}

	magic := string(trailer[8:])
	if magic != bundleMagicBytes {
		return false
	}

	payloadLen := binary.BigEndian.Uint64(trailer[:8])
	if payloadLen == 0 || int64(payloadLen) > stat.Size()-trailerSize {
		return false
	}

	payloadOffset := stat.Size() - trailerSize - int64(payloadLen)
	payloadBytes := make([]byte, payloadLen)
	_, err = f.ReadAt(payloadBytes, payloadOffset)
	if err != nil {
		return false
	}

	// Dekompresi muatan gzip
	gz, err := gzip.NewReader(bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Gagal mendekompresi berkas aplikasi tersemat: %v\n", err)
		os.Exit(1)
	}
	defer gz.Close()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Gagal membaca muatan aplikasi tersemat: %v\n", err)
		os.Exit(1)
	}

	// Format payload: "<nama_berkas_asli>\n<kode_sumber>"
	newlineIdx := bytes.IndexByte(decompressed, '\n')
	virtualFilename := "program.ce"
	var sourceBytes []byte
	if newlineIdx != -1 {
		virtualFilename = string(decompressed[:newlineIdx])
		sourceBytes = decompressed[newlineIdx+1:]
	} else {
		sourceBytes = decompressed
	}

	sourceStr := string(sourceBytes)
	l := parser.NewLexer(sourceStr)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	env := evaluator.NewEnvironment()
	if prog.IsUniversal {
		env.IsUniversal = true
	}
	env.FilePath = filepath.Join(filepath.Dir(exePath), virtualFilename)
	evaluator.RegisterBuiltins(env, os.Stdout)

	// Teruskan seluruh argumen CLI ke skrip
	scriptArgs := os.Args[1:]
	evaluator.SetScriptArgs(scriptArgs)
	argElements := make([]evaluator.Object, len(scriptArgs))
	for i, arg := range scriptArgs {
		argElements[i] = &evaluator.String{Value: arg}
	}
	argsArray := &evaluator.Array{Elements: argElements}
	env.Define("argumen", argsArray, false, "Array")
	env.Define("args", argsArray, false, "Array")

	// Jalankan via Bytecode VM (dengan fallback ke AST jika ada fitur spesifik)
	comp := compiler.NewCompilerWithEnv(env)
	chunk, err := comp.Compile(prog)
	if err != nil {
		result := evaluator.Eval(prog, env)
		if result != nil && result.Type() == evaluator.ERROR_OBJ {
			diag := empathy.NewRuntimeDiagnostic(result.(*evaluator.Error), sourceStr)
			fmt.Print(diag.Render())
			os.Exit(1)
		}
		return true
	}

	virtualMachine := vm.NewWithEnv(chunk, env)
	result, err := virtualMachine.Run()
	if err != nil {
		var evalErr *evaluator.Error
		if ee, ok := err.(*evaluator.Error); ok {
			evalErr = ee
		} else {
			evalErr = &evaluator.Error{Message: err.Error()}
		}
		diag := empathy.NewRuntimeDiagnostic(evalErr, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}
	_ = result

	return true
}

// resolveRunnerPath mencari atau menyiapkan biner pelari untuk sistem target (misal: windows/amd64).
func resolveRunnerPath(explicitRunner, targetOS string) (string, error) {
	// 1. Jika pengguna menentukan runner eksplisit
	if explicitRunner != "" {
		if _, err := os.Stat(explicitRunner); err != nil {
			return "", fmt.Errorf("berkas pelari kustom '--runner=%s' tidak ditemukan: %w", explicitRunner, err)
		}
		return explicitRunner, nil
	}

	// 2. Jika target sama dengan OS saat ini atau target kosong, gunakan diri sendiri
	if targetOS == "" || targetOS == runtime.GOOS {
		exePath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("gagal mendeteksi lokasi biner Celang: %w", err)
		}
		return filepath.EvalSymlinks(exePath)
	}

	// 3. Jika target adalah Windows (dari Linux atau Mac)
	if targetOS == "windows" {
		candidates := []string{
			"celang.exe",
			"./celang.exe",
		}

		if selfExe, err := os.Executable(); err == nil {
			candidates = append(candidates, filepath.Join(filepath.Dir(selfExe), "celang.exe"))
		}

		if homeDir, err := os.UserHomeDir(); err == nil {
			candidates = append(candidates, filepath.Join(homeDir, ".celang", "runners", "windows-amd64.exe"))
		}

		for _, cand := range candidates {
			if _, err := os.Stat(cand); err == nil {
				return cand, nil
			}
		}

		// Jika Go terpasang, otomatis kompilasi runner dan simpan ke cache
		if goPath, err := exec.LookPath("go"); err == nil {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				runnersDir := filepath.Join(homeDir, ".celang", "runners")
				_ = os.MkdirAll(runnersDir, 0755)
				cachePath := filepath.Join(runnersDir, "windows-amd64.exe")

				fmt.Println("⏳ Menyiapkan template pelari Windows (windows/amd64) ke cache...")
				cmd := exec.Command(goPath, "build", "-ldflags=-s -w", "-o", cachePath, "celang/cmd/celang")
				cmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", "CGO_ENABLED=0")
				if _, err := cmd.CombinedOutput(); err == nil {
					fmt.Println("✅ Template pelari Windows siap digunakan!")
					return cachePath, nil
				}
			}
		}

		return "", fmt.Errorf("biner pelari Windows ('celang.exe') belum ditemukan.\n" +
			"   💡 Solusi: Jalankan kompilasi dengan menyertakan runner:\n" +
			"      celang bangun <file.ce> -o <app.exe> --runner=path/ke/celang.exe\n" +
			"   Atau letakkan 'celang.exe' di folder ini atau di ~/.celang/runners/windows-amd64.exe")
	}

	return "", fmt.Errorf("target sistem '%s' belum didukung otomatis", targetOS)
}

// buildFile mengompilasi berkas Celang menjadi biner aplikasi mandiri (standalone executable).
func buildFile(inputFile, outputFile, targetOS, runnerPath string, useAST bool, styleOpt ...string) {
	styleName := ""
	if len(styleOpt) > 0 && styleOpt[0] != "" {
		styleName = styleOpt[0]
	}

	// 1. Validasi sintaks berkas sumber terlebih dahulu
	prog, sourceStr, err := loadProgram(inputFile, styleName)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}
	_ = prog

	// Otomatis deteksi targetOS jika nama outputFile berakhiran .exe
	if targetOS == "" && strings.HasSuffix(strings.ToLower(outputFile), ".exe") {
		targetOS = "windows"
	}

	// 2. Tentukan nama berkas luaran
	if outputFile == "" {
		baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		if targetOS == "windows" || (targetOS == "" && runtime.GOOS == "windows") {
			outputFile = baseName + ".exe"
		} else {
			outputFile = baseName
		}
	} else if targetOS == "windows" && !strings.HasSuffix(strings.ToLower(outputFile), ".exe") {
		outputFile += ".exe"
	}

	// 3. Tentukan biner pelari (runner) dasar
	exePath, err := resolveRunnerPath(runnerPath, targetOS)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	selfFile, err := os.Open(exePath)
	if err != nil {
		fmt.Printf("❌ Gagal membuka biner induk/runner '%s': %v\n", exePath, err)
		os.Exit(1)
	}
	defer selfFile.Close()

	stat, err := selfFile.Stat()
	if err != nil {
		fmt.Printf("❌ Gagal membaca ukuran biner induk: %v\n", err)
		os.Exit(1)
	}

	// Cari ukuran biner murni (jika biner induk sudah pernah dibundle sebelumnya)
	baseBinarySize := stat.Size()
	if baseBinarySize > trailerSize {
		trailer := make([]byte, trailerSize)
		_, err := selfFile.ReadAt(trailer, baseBinarySize-trailerSize)
		if err == nil && string(trailer[8:]) == bundleMagicBytes {
			payloadLen := binary.BigEndian.Uint64(trailer[:8])
			if int64(payloadLen) < baseBinarySize-trailerSize {
				baseBinarySize = baseBinarySize - trailerSize - int64(payloadLen)
			}
		}
	}

	// 4. Baca biner dasar ke memori
	baseBytes := make([]byte, baseBinarySize)
	_, err = io.ReadFull(io.NewSectionReader(selfFile, 0, baseBinarySize), baseBytes)
	if err != nil {
		fmt.Printf("❌ Gagal membaca konten biner dasar: %v\n", err)
		os.Exit(1)
	}

	// 5. Kompresi muatan (nama berkas + kode sumber)
	var payloadBuf bytes.Buffer
	gw := gzip.NewWriter(&payloadBuf)
	payloadHeader := filepath.Base(inputFile) + "\n"
	bundledSource := sourceStr
	if styleName != "" && detectPragmaStyle(sourceStr) == "" {
		bundledSource = "// @style: " + styleName + "\n" + sourceStr
	}
	if _, err := gw.Write([]byte(payloadHeader + bundledSource)); err != nil {
		fmt.Printf("❌ Gagal mengompresi kode program: %v\n", err)
		os.Exit(1)
	}
	if err := gw.Close(); err != nil {
		fmt.Printf("❌ Gagal menyelesaikan kompresi kode: %v\n", err)
		os.Exit(1)
	}
	compressedPayload := payloadBuf.Bytes()

	// 6. Tulis berkas biner mandiri baru
	outFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("❌ Gagal membuat berkas luaran '%s': %v\n", outputFile, err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Tulis biner dasar
	if _, err := outFile.Write(baseBytes); err != nil {
		fmt.Printf("❌ Gagal menulis biner dasar ke '%s': %v\n", outputFile, err)
		os.Exit(1)
	}

	// Tulis payload terkompresi
	if _, err := outFile.Write(compressedPayload); err != nil {
		fmt.Printf("❌ Gagal menulis payload ke '%s': %v\n", outputFile, err)
		os.Exit(1)
	}

	// Tulis trailer (8 byte panjang + 16 byte magic)
	var trailer [trailerSize]byte
	binary.BigEndian.PutUint64(trailer[:8], uint64(len(compressedPayload)))
	copy(trailer[8:], []byte(bundleMagicBytes))
	if _, err := outFile.Write(trailer[:]); err != nil {
		fmt.Printf("❌ Gagal menulis trailer ke '%s': %v\n", outputFile, err)
		os.Exit(1)
	}

	// Pastikan izin berkas dapat dieksekusi (executable) pada Linux/macOS
	_ = os.Chmod(outputFile, 0755)

	outStat, _ := outFile.Stat()
	outSizeMB := float64(outStat.Size()) / (1024 * 1024)

	targetDisplay := runtime.GOOS + " (" + runtime.GOARCH + ")"
	if targetOS != "" {
		targetDisplay = targetOS
		if targetOS == "windows" {
			targetDisplay = "windows (amd64)"
		}
	}

	fmt.Println("✨ Berhasil Membangun Aplikasi Mandiri! ✨")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("📄 Berkas Sumber  : %s\n", inputFile)
	fmt.Printf("📦 Berkas Biner   : %s\n", outputFile)
	fmt.Printf("⚖️  Ukuran Berkas : %.2f MB\n", outSizeMB)
	fmt.Printf("🖥️  Target Sistem : %s\n", targetDisplay)
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("💡 Berkas '%s' ini adalah aplikasi mandiri (standalone executable).\n", outputFile)
	fmt.Println("   Dapat langsung dijalankan atau dibagikan ke teman/klien tanpa")
	fmt.Println("   perlu menginstal Celang ataupun Go di komputer mereka!")
	if targetOS == "windows" || (targetOS == "" && runtime.GOOS == "windows") {
		fmt.Printf("   Jalankan dengan klik dua kali atau ketik: %s\n", outputFile)
	} else {
		fmt.Printf("   Jalankan dengan: ./%s\n", outputFile)
	}
}
