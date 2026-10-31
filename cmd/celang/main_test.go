package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_SubcommandsWithStyle(t *testing.T) {
	tmpDir := t.TempDir()

	styleFile := filepath.Join(tmpDir, "custom.style")
	styleContent := `{
		"deklarasi": {
			"var": "bikin"
		},
		"percabangan": {
			"if": "lamun",
			"else": "lamun_henteu"
		}
	}`
	if err := os.WriteFile(styleFile, []byte(styleContent), 0644); err != nil {
		t.Fatalf("Failed to write style file: %v", err)
	}

	sourceFile := filepath.Join(tmpDir, "program.ce")
	sourceContent := `bikin x = 42
lamun x == 42 {
    cetak(x)
} lamun_henteu {
    cetak("salah")
}
`
	if err := os.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}

	// 1. Test loadProgram with style
	prog, sourceStr, err := loadProgram(sourceFile, styleFile)
	if err != nil {
		t.Fatalf("loadProgram failed: %v", err)
	}
	if prog == nil || len(prog.Body) != 2 {
		t.Fatalf("Expected 2 statements in program body, got %v", len(prog.Body))
	}
	if !strings.Contains(sourceStr, "bikin x = 42") {
		t.Fatalf("Expected sourceStr to contain 'bikin x = 42', got %s", sourceStr)
	}

	// 2. Test compileFile with style
	astOutput := filepath.Join(tmpDir, "program.ast")
	compileFile(sourceFile, astOutput, styleFile)
	astBytes, err := os.ReadFile(astOutput)
	if err != nil {
		t.Fatalf("Failed to read compiled AST file: %v", err)
	}
	if !strings.Contains(string(astBytes), "VariableDeclaration") {
		t.Fatalf("Compiled AST missing VariableDeclaration: %s", string(astBytes))
	}

	// 3. Test formatFile with style
	formatFile(sourceFile, "indonesian", false, styleFile)
}
