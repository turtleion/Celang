package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"

	"celang/pkg/evaluator"
)

// Chunk represents a contiguous sequence of bytecode instructions,
// along with its associated constant pool and debug line information.
type Chunk struct {
	Instructions []byte
	Constants    []evaluator.Object
	Lines        []int
}

func NewChunk() *Chunk {
	return &Chunk{
		Instructions: make([]byte, 0, 128),
		Constants:    make([]evaluator.Object, 0, 32),
		Lines:        make([]int, 0, 128),
	}
}

// Write appends a single byte instruction to the chunk.
func (c *Chunk) Write(b byte, line int) {
	c.Instructions = append(c.Instructions, b)
	c.Lines = append(c.Lines, line)
}

// WriteOp appends an Opcode.
func (c *Chunk) WriteOp(op Opcode, line int) {
	c.Write(byte(op), line)
}

// WriteUint16 appends a 16-bit big-endian unsigned integer (operand).
func (c *Chunk) WriteUint16(val uint16, line int) {
	c.Write(byte(val>>8), line)
	c.Write(byte(val&0xFF), line)
}

// AddConstant adds an object to the constant pool and returns its 16-bit index.
func (c *Chunk) AddConstant(obj evaluator.Object) uint16 {
	c.Constants = append(c.Constants, obj)
	return uint16(len(c.Constants) - 1)
}

// ReadUint16 reads a 16-bit operand starting at offset.
func (c *Chunk) ReadUint16(offset int) uint16 {
	return binary.BigEndian.Uint16(c.Instructions[offset:])
}

// Disassemble prints all instructions in human-readable assembly format.
func (c *Chunk) Disassemble(name string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("== Chunk: %s (Bytes: %d, Constants: %d) ==\n", name, len(c.Instructions), len(c.Constants)))

	offset := 0
	for offset < len(c.Instructions) {
		lineStr, nextOffset := c.DisassembleInstruction(offset)
		sb.WriteString(lineStr)
		sb.WriteString("\n")
		offset = nextOffset
	}

	return sb.String()
}

// DisassembleInstruction formats a single instruction at offset.
func (c *Chunk) DisassembleInstruction(offset int) (string, int) {
	line := c.Lines[offset]
	op := Opcode(c.Instructions[offset])

	prefix := fmt.Sprintf("%04d  line %3d  %-16s", offset, line, op.String())

	switch op {
	case OpConstant, OpDefineGlobal, OpGetGlobal, OpSetGlobal:
		if offset+2 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		constIdx := c.ReadUint16(offset + 1)
		valStr := "<null>"
		if int(constIdx) < len(c.Constants) {
			valStr = c.Constants[constIdx].Inspect()
		}
		return fmt.Sprintf("%s  #%d (%s)", prefix, constIdx, valStr), offset + 3

	case OpGetLocal, OpSetLocal:
		if offset+2 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		slot := c.ReadUint16(offset + 1)
		return fmt.Sprintf("%s  slot %d", prefix, slot), offset + 3

	case OpJumpIfFalse, OpJump:
		if offset+2 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		target := c.ReadUint16(offset + 1)
		return fmt.Sprintf("%s  -> %04d", prefix, target), offset + 3

	case OpLoop:
		if offset+2 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		target := c.ReadUint16(offset + 1)
		return fmt.Sprintf("%s  -> %04d (loop)", prefix, target), offset + 3

	case OpArray, OpMap:
		if offset+2 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		count := c.ReadUint16(offset + 1)
		return fmt.Sprintf("%s  count %d", prefix, count), offset + 3

	case OpCall:
		if offset+1 >= len(c.Instructions) {
			return fmt.Sprintf("%s [truncated]", prefix), offset + 1
		}
		args := c.Instructions[offset+1]
		return fmt.Sprintf("%s  args %d", prefix, args), offset + 2

	default:
		return prefix, offset + 1
	}
}
