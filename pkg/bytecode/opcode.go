package bytecode

import "fmt"

// Opcode represents a single bytecode operation (1 byte).
type Opcode byte

const (
	// Literals & Constants
	OpConstant Opcode = iota // [OpConstant, idx_high, idx_low]
	OpNull                   // pushes null
	OpTrue                   // pushes true
	OpFalse                  // pushes false

	// Arithmetic Operations
	OpAdd // pop b, pop a, push (a + b)
	OpSub // pop b, pop a, push (a - b)
	OpMul // pop b, pop a, push (a * b)
	OpDiv // pop b, pop a, push (a / b)
	OpMod // pop b, pop a, push (a % b)
	OpNeg // pop a, push (-a)

	// Comparison & Logical Operations
	OpEqual        // pop b, pop a, push (a == b)
	OpNotEqual     // pop b, pop a, push (a != b)
	OpGreaterThan  // pop b, pop a, push (a > b)
	OpGreaterEqual // pop b, pop a, push (a >= b)
	OpLessThan     // pop b, pop a, push (a < b)
	OpLessEqual    // pop b, pop a, push (a <= b)
	OpNot          // pop a, push (!a)

	// Variables & Scope
	OpPop          // discard top of stack
	OpDefineGlobal // [OpDefineGlobal, idx_high, idx_low]
	OpGetGlobal    // [OpGetGlobal, idx_high, idx_low]
	OpSetGlobal    // [OpSetGlobal, idx_high, idx_low]
	OpGetLocal     // [OpGetLocal, slot_high, slot_low]
	OpSetLocal     // [OpSetLocal, slot_high, slot_low]

	// Jumps & Control Flow
	OpJumpIfFalse // [OpJumpIfFalse, offset_high, offset_low]
	OpJump        // [OpJump, offset_high, offset_low]
	OpLoop        // [OpLoop, offset_high, offset_low] (backward jump)

	// Data Structures
	OpArray // [OpArray, count_high, count_low]
	OpMap   // [OpMap, count_high, count_low]
	OpIndex // pop index, pop target, push target[index]
	OpSetIndex // pop val, pop index, pop target, target[index] = val

	// Functions & Calls
	OpCall   // [OpCall, argCount]
	OpReturn // returns from current frame / halts program
)

// OpcodeNames maps Opcodes to human-readable assembly instructions.
var OpcodeNames = map[Opcode]string{
	OpConstant:     "OP_CONSTANT",
	OpNull:         "OP_NULL",
	OpTrue:         "OP_TRUE",
	OpFalse:        "OP_FALSE",
	OpAdd:          "OP_ADD",
	OpSub:          "OP_SUB",
	OpMul:          "OP_MUL",
	OpDiv:          "OP_DIV",
	OpMod:          "OP_MOD",
	OpNeg:          "OP_NEG",
	OpEqual:        "OP_EQUAL",
	OpNotEqual:     "OP_NOT_EQUAL",
	OpGreaterThan:  "OP_GREATER_THAN",
	OpGreaterEqual: "OP_GREATER_EQUAL",
	OpLessThan:     "OP_LESS_THAN",
	OpLessEqual:    "OP_LESS_EQUAL",
	OpNot:          "OP_NOT",
	OpPop:          "OP_POP",
	OpDefineGlobal: "OP_DEFINE_GLOBAL",
	OpGetGlobal:    "OP_GET_GLOBAL",
	OpSetGlobal:    "OP_SET_GLOBAL",
	OpGetLocal:     "OP_GET_LOCAL",
	OpSetLocal:     "OP_SET_LOCAL",
	OpJumpIfFalse:  "OP_JUMP_IF_FALSE",
	OpJump:         "OP_JUMP",
	OpLoop:         "OP_LOOP",
	OpArray:        "OP_ARRAY",
	OpMap:          "OP_MAP",
	OpIndex:        "OP_INDEX",
	OpSetIndex:     "OP_SET_INDEX",
	OpCall:         "OP_CALL",
	OpReturn:       "OP_RETURN",
}

func (op Opcode) String() string {
	if name, ok := OpcodeNames[op]; ok {
		return name
	}
	return fmt.Sprintf("OP_UNKNOWN(%d)", byte(op))
}
