package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"celang/pkg/compiler"
	"celang/pkg/empathy"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/repl"
	"celang/pkg/syntax"
	"celang/pkg/uast"
	"celang/pkg/vm"
)

const AppVersion = "1.0.0"

func main() {
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
			fmt.Println("❌ Penggunaan: celang run <file.ce | file.ast> [--ast] [-- argumen1 argumen2 ...]")
			os.Exit(1)
		}
		filePath := ""
		useAST := false
		var scriptArgs []string

		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "--" {
				scriptArgs = os.Args[i+1:]
				break
			}
			if arg == "--ast" || arg == "--legacy" || arg == "--tree-walk" {
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
		runFile(filePath, scriptArgs, useAST)

	case "disasm", "disassemble":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang disasm <file.ce>")
			os.Exit(1)
		}
		disasmFile(os.Args[2])

	case "compile", "parse":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang compile <file.ce> [-o output.ast]")
			os.Exit(1)
		}
		inputFile := os.Args[2]
		outputFile := ""
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "-o" && i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				break
			}
		}
		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, ".ce") + ".ast"
		}
		compileFile(inputFile, outputFile)

	case "format":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang format [-w | --write] <file.ce> [--style=nama_style | path/ke/style.style]")
			os.Exit(1)
		}
		isWrite := false
		styleName := "indonesian"
		filePath := ""

		for _, arg := range os.Args[2:] {
			if arg == "-w" || arg == "--write" {
				isWrite = true
			} else if strings.HasPrefix(arg, "--style=") {
				styleName = strings.TrimPrefix(arg, "--style=")
			} else if filePath == "" && !strings.HasPrefix(arg, "-") {
				filePath = arg
			}
		}

		if filePath == "" {
			fmt.Println("❌ Mohon sertakan berkas kode yang ingin diformat.")
			os.Exit(1)
		}
		formatFile(filePath, styleName, isWrite)

	case "init":
		projectName := "projek_celang"
		if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "-") {
			projectName = os.Args[2]
		}
		initProject(projectName)

	case "ast", "uast":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang ast <file.ce>")
			os.Exit(1)
		}
		dumpAST(os.Args[2])

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
	fmt.Println("  celang disasm <file.ce>                 Membongkar kode menjadi instruksi rakitan Bytecode")
	fmt.Println("  celang format [-w] <file> [--style=...] Memformat kode (-w untuk langsung menimpa file)")
	fmt.Println("  celang compile <file.ce> [-o out.ast]   Mengompilasi kode teks menjadi berkas U-AST JSON")
	fmt.Println("  celang ast <file.ce>                    Mencetak pohon struktur U-AST JSON ke terminal")
	fmt.Println("  celang init [nama_proyek]               Menginisialisasi proyek baru (celang.mod & utama.ce)")
	fmt.Println("  celang versi, --version, -v             Melihat nomor versi Celang")
	fmt.Println("  celang bantuan, --help, -h              Menampilkan bantuan ini")
}

func loadProgram(filePath string) (*uast.Program, string, error) {
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

func runFile(filePath string, scriptArgs []string, useAST bool) {
	prog, sourceStr, err := loadProgram(filePath)
	if err != nil {
		handleSyntaxError(err, sourceStr)
	}

	env := evaluator.NewEnvironment()
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
	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		fmt.Printf("❌ Kesalahan kompilasi Bytecode: %v\n", err)
		os.Exit(1)
	}
	virtualMachine := vm.NewWithEnv(chunk, env)
	result, err := virtualMachine.Run()
	if err != nil {
		fmt.Printf("❌ Kesalahan Runtime VM: %v\n", err)
		os.Exit(1)
	}
	_ = result
}

func disasmFile(filePath string) {
	prog, sourceStr, err := loadProgram(filePath)
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

func compileFile(inputFile, outputFile string) {
	prog, sourceStr, err := loadProgram(inputFile)
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

func formatFile(filePath, styleName string, isWrite bool) {
	prog, sourceStr, err := loadProgram(filePath)
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

func dumpAST(filePath string) {
	prog, sourceStr, err := loadProgram(filePath)
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
