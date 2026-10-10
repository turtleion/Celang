package bytecode_test

import (
	"strings"
	"testing"

	"celang/pkg/bytecode"
	"celang/pkg/evaluator"
)

func TestChunk_WriteAndDisassemble(t *testing.T) {
	chunk := bytecode.NewChunk()
	constIdx := chunk.AddConstant(&evaluator.Integer{Value: 42})

	chunk.WriteOp(bytecode.OpConstant, 1)
	chunk.WriteUint16(constIdx, 1)
	chunk.WriteOp(bytecode.OpAdd, 1)
	chunk.WriteOp(bytecode.OpReturn, 2)

	disasm := chunk.Disassemble("test_chunk")
	if !strings.Contains(disasm, "OP_CONSTANT") {
		t.Errorf("Expected disasm to contain OP_CONSTANT, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "42") {
		t.Errorf("Expected disasm to contain '42', got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_ADD") {
		t.Errorf("Expected disasm to contain OP_ADD, got:\n%s", disasm)
	}
	if !strings.Contains(disasm, "OP_RETURN") {
		t.Errorf("Expected disasm to contain OP_RETURN, got:\n%s", disasm)
	}
}
