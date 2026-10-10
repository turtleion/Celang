package compiler_test

import (
	"strings"
	"testing"

	"celang/pkg/compiler"
	"celang/pkg/parser"
)

func TestCompiler_BasicExpressions(t *testing.T) {
	input := "10 + 20 * 3"
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	disasm := chunk.Disassemble("arithmetic")
	if !strings.Contains(disasm, "OP_ADD") || !strings.Contains(disasm, "OP_MUL") {
		t.Errorf("Expected OP_ADD and OP_MUL, got:\n%s", disasm)
	}
}

func TestCompiler_VariablesAndControlFlow(t *testing.T) {
	input := `
	buat x = 10
	buat hasil = 0
	jika x > 5 {
		hasil = 100
	} selain {
		hasil = 50
	}
	selama x > 0 {
		x = x - 1
	}
	`
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	disasm := chunk.Disassemble("control_flow")
	if !strings.Contains(disasm, "OP_DEFINE_GLOBAL") {
		t.Errorf("Expected OP_DEFINE_GLOBAL, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_JUMP_IF_FALSE") {
		t.Errorf("Expected OP_JUMP_IF_FALSE, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_LOOP") {
		t.Errorf("Expected OP_LOOP, got:\n%s", disasm)
	}
}

func TestCompiler_FunctionsAndDataStructures(t *testing.T) {
	input := `
	fungsi tambah(a, b) {
		kembalikan a + b
	}
	buat daftar = [1, 2, 3]
	buat kamus = {"kunci": "nilai"}
	buat hasil = tambah(daftar[0], 10)
	`
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	disasm := chunk.Disassemble("functions_and_data")
	if !strings.Contains(disasm, "OP_ARRAY") {
		t.Errorf("Expected OP_ARRAY, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_MAP") {
		t.Errorf("Expected OP_MAP, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_CALL") {
		t.Errorf("Expected OP_CALL, got:\n%s", disasm)
	}
}
