package main

import (
	"fmt"
	"os"

	"celang/pkg/evaluator"
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
			fmt.Println("❌ Penggunaan: celang run <path-ke-file.ast>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Perintah tidak dikenal: '%s'\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Celang - Bahasa Pemrograman 'Syntax Styling' & Human-Centric")
	fmt.Println("\nPerintah yang tersedia:")
	fmt.Println("  celang run <file.ast>    Menjalankan file Universal AST (U-AST JSON)")
	fmt.Println("  celang help              Menampilkan bantuan ini")
}

func runFile(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Gagal membaca file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	prog, err := uast.FromJSON(data)
	if err != nil {
		fmt.Printf("❌ Format U-AST tidak valid: %v\n", err)
		os.Exit(1)
	}

	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, os.Stdout)

	result := evaluator.Eval(prog, env)
	if result != nil && result.Type() == evaluator.ERROR_OBJ {
		fmt.Println("\n❌ " + result.Inspect())
		os.Exit(1)
	}
}
