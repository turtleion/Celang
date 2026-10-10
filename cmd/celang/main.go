package main

import (
	"fmt"
	"os"
	"strings"

	"celang/pkg/empathy"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/syntax"
	"celang/pkg/uast"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("❌ Penggunaan: celang run <file.ce | file.ast>")
			os.Exit(1)
		}
		runFile(os.Args[2])

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
			fmt.Println("❌ Penggunaan: celang format <file.ce | file.ast> [--style=nama_style | path/ke/style.style]")
			os.Exit(1)
		}
		inputFile := os.Args[2]
		styleName := "indonesian"
		for i := 3; i < len(os.Args); i++ {
			if strings.HasPrefix(os.Args[i], "--style=") {
				styleName = strings.TrimPrefix(os.Args[i], "--style=")
			}
		}
		formatFile(inputFile, styleName)

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Perintah tidak dikenal: '%s'\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("🌟 Celang - Bahasa Pemrograman 'Syntax Styling' & Human-Centric")
	fmt.Println("\nPerintah yang tersedia:")
	fmt.Println("  celang run <file.ce | file.ast>           Menjalankan file kode teks (.ce) atau U-AST (.ast)")
	fmt.Println("  celang compile <file.ce> [-o out.ast]     Mengubah kode teks menjadi file Universal AST (U-AST)")
	fmt.Println("  celang format <file> [--style=nama|path]  Memformat/merender kode ke gaya sintaks tertentu")
	fmt.Println("                                            (pilihan style: indonesian, javascript, pythonic, atau file .style)")
	fmt.Println("  celang help                               Menampilkan bantuan ini")
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

func runFile(filePath string) {
	prog, sourceStr, err := loadProgram(filePath)
	if err != nil {
		diag := empathy.NewSyntaxDiagnostic(err.Error(), 1, 1, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}

	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, os.Stdout)

	result := evaluator.Eval(prog, env)
	if result != nil && result.Type() == evaluator.ERROR_OBJ {
		diag := empathy.NewRuntimeDiagnostic(result.(*evaluator.Error), sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}
}

func compileFile(inputFile, outputFile string) {
	prog, sourceStr, err := loadProgram(inputFile)
	if err != nil {
		diag := empathy.NewSyntaxDiagnostic(err.Error(), 1, 1, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
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

func formatFile(filePath, styleName string) {
	prog, sourceStr, err := loadProgram(filePath)
	if err != nil {
		diag := empathy.NewSyntaxDiagnostic(err.Error(), 1, 1, sourceStr)
		fmt.Print(diag.Render())
		os.Exit(1)
	}

	profile := syntax.GetProfile(styleName)
	formatter := syntax.NewFormatter(profile)
	formatted := formatter.Format(prog)

	fmt.Printf("--- Menampilkan Kode dengan Gaya: %s ---\n", profile.DisplayName)
	fmt.Print(formatted)
}
