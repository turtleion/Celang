package repl_test

import (
	"bytes"
	"strings"
	"testing"

	"celang/pkg/repl"
)

func TestREPL_BasicEvaluation(t *testing.T) {
	input := "10 + 20\n:keluar\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	err := repl.Start(repl.Options{
		In:       in,
		Out:      &out,
		NoBanner: true,
	})
	if err != nil {
		t.Fatalf("REPL error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "=> 30") {
		t.Errorf("Expected '=> 30' in output, got: %s", output)
	}
	if !strings.Contains(output, "Sampai jumpa lagi!") {
		t.Errorf("Expected exit message, got: %s", output)
	}
}

func TestREPL_MultilineBlock(t *testing.T) {
	input := `fungsi kali(a, b) {
    kembalikan a * b
}
kali(6, 7)
:keluar
`
	in := strings.NewReader(input)
	var out bytes.Buffer

	err := repl.Start(repl.Options{
		In:       in,
		Out:      &out,
		NoBanner: true,
	})
	if err != nil {
		t.Fatalf("REPL error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "=> 42") {
		t.Errorf("Expected '=> 42' in output, got: %s", output)
	}
}

func TestREPL_MetaCommands(t *testing.T) {
	input := `buat nama = "Celang"
:lingkungan
:gaya pythonic
:ast 1 + 2
:bantuan
:keluar
`
	in := strings.NewReader(input)
	var out bytes.Buffer

	err := repl.Start(repl.Options{
		In:       in,
		Out:      &out,
		NoBanner: true,
	})
	if err != nil {
		t.Fatalf("REPL error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "nama") || !strings.Contains(output, "Celang") {
		t.Errorf("Expected variable in :lingkungan, got: %s", output)
	}
	if !strings.Contains(output, "Gaya berhasil diganti") {
		t.Errorf("Expected style change confirmation, got: %s", output)
	}
	if !strings.Contains(output, "BinaryExpression") {
		t.Errorf("Expected AST output, got: %s", output)
	}
	if !strings.Contains(output, "Daftar Perintah Khusus") {
		t.Errorf("Expected help output, got: %s", output)
	}
}
