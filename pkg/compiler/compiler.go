package compiler

import (
	"fmt"
	"strconv"

	"celang/pkg/bytecode"
	"celang/pkg/evaluator"
	"celang/pkg/uast"
)

// Local represents a variable stored on the VM stack within a lexical scope.
type Local struct {
	Name  string
	Depth int
}

// Compiler translates U-AST nodes into bytecode chunks.
type Compiler struct {
	parent     *Compiler
	chunk      *bytecode.Chunk
	locals     []Local
	scopeDepth int
	fnName     string
	numParams  int
}

// NewCompiler creates a new top-level Compiler.
func NewCompiler() *Compiler {
	return &Compiler{
		chunk:      bytecode.NewChunk(),
		locals:     make([]Local, 0, 32),
		scopeDepth: 0,
	}
}

func newFunctionCompiler(parent *Compiler, name string, numParams int) *Compiler {
	return &Compiler{
		parent:     parent,
		chunk:      bytecode.NewChunk(),
		locals:     make([]Local, 0, 16),
		scopeDepth: 0,
		fnName:     name,
		numParams:  numParams,
	}
}

// Compile processes a U-AST Program and returns the compiled bytecode chunk.
func (c *Compiler) Compile(program *uast.Program) (*bytecode.Chunk, error) {
	for i, stmt := range program.Body {
		isLast := (i == len(program.Body)-1)
		if isLast {
			if exprStmt, ok := stmt.(*uast.ExpressionStatement); ok {
				line := 1
				if exprStmt.GetMetadata() != nil && exprStmt.GetMetadata().Pos != nil {
					line = exprStmt.GetMetadata().Pos.Line
				}
				if err := c.compileExpression(exprStmt.Expression); err != nil {
					return nil, err
				}
				c.emitOp(bytecode.OpReturn, line)
				return c.chunk, nil
			}
		}
		if err := c.compileStatement(stmt); err != nil {
			return nil, err
		}
	}
	c.emitOp(bytecode.OpNull, 0)
	c.emitOp(bytecode.OpReturn, 0)
	return c.chunk, nil
}

// CompileExpression compiles a single expression for REPL evaluation.
func (c *Compiler) CompileExpression(expr uast.Expression) (*bytecode.Chunk, error) {
	if err := c.compileExpression(expr); err != nil {
		return nil, err
	}
	c.emitOp(bytecode.OpReturn, 0)
	return c.chunk, nil
}

func (c *Compiler) compileStatement(stmt uast.Statement) error {
	line := 1
	if stmt.GetMetadata() != nil && stmt.GetMetadata().Pos != nil {
		line = stmt.GetMetadata().Pos.Line
	}

	switch s := stmt.(type) {
	case *uast.VariableDeclaration:
		if s.Value != nil {
			if err := c.compileExpression(s.Value); err != nil {
				return err
			}
		} else {
			c.emitOp(bytecode.OpNull, line)
		}

		if c.scopeDepth > 0 || c.parent != nil {
			// Local variable on stack
			c.addLocal(s.Name)
		} else {
			// Global variable
			constIdx := c.chunk.AddConstant(&evaluator.String{Value: s.Name})
			c.emitOp(bytecode.OpDefineGlobal, line)
			c.chunk.WriteUint16(constIdx, line)
		}
		return nil

	case *uast.AssignmentStatement:
		if err := c.compileExpression(s.Value); err != nil {
			return err
		}

		switch target := s.Target.(type) {
		case *uast.Identifier:
			slot, isLocal := c.resolveLocal(target.Name)
			if isLocal {
				c.emitOp(bytecode.OpSetLocal, line)
				c.chunk.WriteUint16(slot, line)
			} else {
				constIdx := c.chunk.AddConstant(&evaluator.String{Value: target.Name})
				c.emitOp(bytecode.OpSetGlobal, line)
				c.chunk.WriteUint16(constIdx, line)
			}
			c.emitOp(bytecode.OpPop, line)
			return nil

		case *uast.MemberExpression:
			// Stack currently has [value]
			if err := c.compileExpression(target.Object); err != nil {
				return err
			}
			if target.Computed {
				if intVal, err := strconv.ParseInt(target.Property, 10, 64); err == nil {
					constIdx := c.chunk.AddConstant(&evaluator.Integer{Value: intVal})
					c.emitOp(bytecode.OpConstant, line)
					c.chunk.WriteUint16(constIdx, line)
				} else if slot, isLocal := c.resolveLocal(target.Property); isLocal {
					c.emitOp(bytecode.OpGetLocal, line)
					c.chunk.WriteUint16(slot, line)
				} else {
					constIdx := c.chunk.AddConstant(&evaluator.String{Value: target.Property})
					c.emitOp(bytecode.OpConstant, line)
					c.chunk.WriteUint16(constIdx, line)
				}
			} else {
				constIdx := c.chunk.AddConstant(&evaluator.String{Value: target.Property})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
			}
			c.emitOp(bytecode.OpSetIndex, line)
			c.emitOp(bytecode.OpPop, line)
			return nil

		default:
			return fmt.Errorf("assignment target '%T' belum didukung di bytecode compiler", target)
		}

	case *uast.ExpressionStatement:
		if err := c.compileExpression(s.Expression); err != nil {
			return err
		}
		c.emitOp(bytecode.OpPop, line)
		return nil

	case *uast.IfStatement:
		if err := c.compileExpression(s.Condition); err != nil {
			return err
		}

		// Jump over consequent if false
		jumpFalse := c.emitJump(bytecode.OpJumpIfFalse, line)
		c.emitOp(bytecode.OpPop, line) // pop condition in true branch

		if err := c.compileStatement(s.Consequent); err != nil {
			return err
		}

		jumpEnd := c.emitJump(bytecode.OpJump, line)
		c.patchJump(jumpFalse)
		c.emitOp(bytecode.OpPop, line) // pop condition in false branch

		if s.Alternate != nil {
			if err := c.compileStatement(s.Alternate); err != nil {
				return err
			}
		}

		c.patchJump(jumpEnd)
		return nil

	case *uast.WhileStatement:
		loopStart := len(c.chunk.Instructions)

		if err := c.compileExpression(s.Condition); err != nil {
			return err
		}

		exitJump := c.emitJump(bytecode.OpJumpIfFalse, line)
		c.emitOp(bytecode.OpPop, line)

		if err := c.compileStatement(s.Body); err != nil {
			return err
		}

		c.emitLoop(loopStart, line)
		c.patchJump(exitJump)
		c.emitOp(bytecode.OpPop, line)
		return nil

	case *uast.ForStatement:
		c.beginScope()
		if s.Init != nil {
			if err := c.compileStatement(s.Init); err != nil {
				return err
			}
		}

		loopStart := len(c.chunk.Instructions)
		var exitJump int = -1

		if s.Condition != nil {
			if err := c.compileExpression(s.Condition); err != nil {
				return err
			}
			exitJump = c.emitJump(bytecode.OpJumpIfFalse, line)
			c.emitOp(bytecode.OpPop, line)
		}

		if err := c.compileStatement(s.Body); err != nil {
			return err
		}

		if s.Update != nil {
			if err := c.compileStatement(s.Update); err != nil {
				return err
			}
		}

		c.emitLoop(loopStart, line)

		if exitJump != -1 {
			c.patchJump(exitJump)
			c.emitOp(bytecode.OpPop, line)
		}

		c.endScope(line)
		return nil

	case *uast.ForEachStatement:
		c.beginScope()

		// 1. Compile iterable and place in local slot
		if err := c.compileExpression(s.Iterable); err != nil {
			return err
		}
		iterSlot := len(c.locals)
		c.addLocal("__iter")

		// 2. Call panjang(__iter) and store length in local slot
		constPanjang := c.chunk.AddConstant(&evaluator.String{Value: "panjang"})
		c.emitOp(bytecode.OpGetGlobal, line)
		c.chunk.WriteUint16(constPanjang, line)
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(iterSlot), line)
		c.emitOp(bytecode.OpCall, line)
		c.chunk.Write(1, line)
		lenSlot := len(c.locals)
		c.addLocal("__len")

		// 3. Initialize __idx = 0 in local slot
		constZero := c.chunk.AddConstant(&evaluator.Integer{Value: 0})
		c.emitOp(bytecode.OpConstant, line)
		c.chunk.WriteUint16(constZero, line)
		idxSlot := len(c.locals)
		c.addLocal("__idx")

		// 4. Reserve slot for user item variable
		c.emitOp(bytecode.OpNull, line)
		itemSlot := len(c.locals)
		c.addLocal(s.Item)

		// 5. Loop condition: __idx < __len
		loopStart := len(c.chunk.Instructions)
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(idxSlot), line)
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(lenSlot), line)
		c.emitOp(bytecode.OpLessThan, line)

		exitJump := c.emitJump(bytecode.OpJumpIfFalse, line)
		c.emitOp(bytecode.OpPop, line)

		// 6. Set s.Item = __iter[__idx]
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(iterSlot), line)
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(idxSlot), line)
		c.emitOp(bytecode.OpIndex, line)
		c.emitOp(bytecode.OpSetLocal, line)
		c.chunk.WriteUint16(uint16(itemSlot), line)
		c.emitOp(bytecode.OpPop, line)

		// 7. Loop body
		if s.Body != nil {
			for _, child := range s.Body.Statements {
				if err := c.compileStatement(child); err != nil {
					return err
				}
			}
		}

		// 8. Advance index: __idx = __idx + 1
		c.emitOp(bytecode.OpGetLocal, line)
		c.chunk.WriteUint16(uint16(idxSlot), line)
		constOne := c.chunk.AddConstant(&evaluator.Integer{Value: 1})
		c.emitOp(bytecode.OpConstant, line)
		c.chunk.WriteUint16(constOne, line)
		c.emitOp(bytecode.OpAdd, line)
		c.emitOp(bytecode.OpSetLocal, line)
		c.chunk.WriteUint16(uint16(idxSlot), line)
		c.emitOp(bytecode.OpPop, line)

		// 9. Jump back
		c.emitLoop(loopStart, line)

		// 10. Exit loop
		c.patchJump(exitJump)
		c.emitOp(bytecode.OpPop, line)

		c.endScope(line)
		return nil

	case *uast.BlockStatement:
		c.beginScope()
		for _, child := range s.Statements {
			if err := c.compileStatement(child); err != nil {
				return err
			}
		}
		c.endScope(line)
		return nil

	case *uast.ReturnStatement:
		if s.Value != nil {
			if err := c.compileExpression(s.Value); err != nil {
				return err
			}
		} else {
			c.emitOp(bytecode.OpNull, line)
		}
		c.emitOp(bytecode.OpReturn, line)
		return nil

	case *uast.FunctionDeclaration:
		fnComp := newFunctionCompiler(c, s.Name, len(s.Parameters))
		for _, p := range s.Parameters {
			fnComp.addLocal(p.Name)
		}
		if s.Body != nil {
			for _, child := range s.Body.Statements {
				if err := fnComp.compileStatement(child); err != nil {
					return err
				}
			}
		}
		fnComp.emitOp(bytecode.OpNull, line)
		fnComp.emitOp(bytecode.OpReturn, line)

		compiledFn := &bytecode.CompiledFunction{
			Name:      s.Name,
			NumParams: len(s.Parameters),
			Chunk:     fnComp.chunk,
		}
		constIdx := c.chunk.AddConstant(compiledFn)
		c.emitOp(bytecode.OpConstant, line)
		c.chunk.WriteUint16(constIdx, line)

		if c.scopeDepth > 0 || c.parent != nil {
			c.addLocal(s.Name)
		} else {
			nameIdx := c.chunk.AddConstant(&evaluator.String{Value: s.Name})
			c.emitOp(bytecode.OpDefineGlobal, line)
			c.chunk.WriteUint16(nameIdx, line)
		}
		return nil

	default:
		return fmt.Errorf("statement '%T' belum didukung oleh bytecode compiler", stmt)
	}
}

func (c *Compiler) compileExpression(expr uast.Expression) error {
	line := 1
	if expr.GetMetadata() != nil && expr.GetMetadata().Pos != nil {
		line = expr.GetMetadata().Pos.Line
	}

	switch e := expr.(type) {
	case *uast.Literal:
		switch e.Kind {
		case uast.LiteralKindInt:
			if iVal, ok := e.Value.(int64); ok {
				constIdx := c.chunk.AddConstant(&evaluator.Integer{Value: iVal})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
				return nil
			} else if iVal, ok := e.Value.(int); ok {
				constIdx := c.chunk.AddConstant(&evaluator.Integer{Value: int64(iVal)})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
				return nil
			}
		case uast.LiteralKindFloat:
			if fVal, ok := e.Value.(float64); ok {
				constIdx := c.chunk.AddConstant(&evaluator.Float{Value: fVal})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
				return nil
			}
		case uast.LiteralKindString:
			sVal := e.Value.(string)
			constIdx := c.chunk.AddConstant(&evaluator.String{Value: sVal})
			c.emitOp(bytecode.OpConstant, line)
			c.chunk.WriteUint16(constIdx, line)
			return nil
		case uast.LiteralKindBool:
			bVal := e.Value.(bool)
			if bVal {
				c.emitOp(bytecode.OpTrue, line)
			} else {
				c.emitOp(bytecode.OpFalse, line)
			}
			return nil
		case uast.LiteralKindNull:
			c.emitOp(bytecode.OpNull, line)
			return nil
		}
		return fmt.Errorf("literal kind '%s' belum didukung di bytecode compiler", e.Kind)

	case *uast.Identifier:
		slot, isLocal := c.resolveLocal(e.Name)
		if isLocal {
			c.emitOp(bytecode.OpGetLocal, line)
			c.chunk.WriteUint16(slot, line)
		} else {
			constIdx := c.chunk.AddConstant(&evaluator.String{Value: e.Name})
			c.emitOp(bytecode.OpGetGlobal, line)
			c.chunk.WriteUint16(constIdx, line)
		}
		return nil

	case *uast.BinaryExpression:
		if e.Operator == "&&" || e.Operator == "dan" {
			if err := c.compileExpression(e.Left); err != nil {
				return err
			}
			endJump := c.emitJump(bytecode.OpJumpIfFalse, line)
			c.emitOp(bytecode.OpPop, line)
			if err := c.compileExpression(e.Right); err != nil {
				return err
			}
			c.patchJump(endJump)
			return nil
		}
		if e.Operator == "||" || e.Operator == "atau" {
			if err := c.compileExpression(e.Left); err != nil {
				return err
			}
			elseJump := c.emitJump(bytecode.OpJumpIfFalse, line)
			endJump := c.emitJump(bytecode.OpJump, line)
			c.patchJump(elseJump)
			c.emitOp(bytecode.OpPop, line)
			if err := c.compileExpression(e.Right); err != nil {
				return err
			}
			c.patchJump(endJump)
			return nil
		}
		if e.Operator == "|>" {
			return c.compilePipeline(e.Left, e.Right, line)
		}

		if err := c.compileExpression(e.Left); err != nil {
			return err
		}
		if err := c.compileExpression(e.Right); err != nil {
			return err
		}

		switch e.Operator {
		case "+":
			c.emitOp(bytecode.OpAdd, line)
		case "-":
			c.emitOp(bytecode.OpSub, line)
		case "*":
			c.emitOp(bytecode.OpMul, line)
		case "/":
			c.emitOp(bytecode.OpDiv, line)
		case "%":
			c.emitOp(bytecode.OpMod, line)
		case "==", "===":
			c.emitOp(bytecode.OpEqual, line)
		case "!=", "!==":
			c.emitOp(bytecode.OpNotEqual, line)
		case ">":
			c.emitOp(bytecode.OpGreaterThan, line)
		case ">=":
			c.emitOp(bytecode.OpGreaterEqual, line)
		case "<":
			c.emitOp(bytecode.OpLessThan, line)
		case "<=":
			c.emitOp(bytecode.OpLessEqual, line)
		default:
			return fmt.Errorf("operator biner '%s' belum didukung di bytecode compiler", e.Operator)
		}
		return nil

	case *uast.UnaryExpression:
		if err := c.compileExpression(e.Argument); err != nil {
			return err
		}
		switch e.Operator {
		case "-":
			c.emitOp(bytecode.OpNeg, line)
		case "!", "bukan":
			c.emitOp(bytecode.OpNot, line)
		default:
			return fmt.Errorf("operator unary '%s' belum didukung di bytecode compiler", e.Operator)
		}
		return nil

	case *uast.ArrayLiteral:
		for _, el := range e.Elements {
			if err := c.compileExpression(el); err != nil {
				return err
			}
		}
		c.emitOp(bytecode.OpArray, line)
		c.chunk.WriteUint16(uint16(len(e.Elements)), line)
		return nil

	case *uast.MapLiteral:
		for _, prop := range e.Properties {
			if err := c.compileExpression(prop.Key); err != nil {
				return err
			}
			if err := c.compileExpression(prop.Value); err != nil {
				return err
			}
		}
		c.emitOp(bytecode.OpMap, line)
		c.chunk.WriteUint16(uint16(len(e.Properties)), line)
		return nil

	case *uast.MemberExpression:
		if err := c.compileExpression(e.Object); err != nil {
			return err
		}
		if e.Computed {
			if intVal, err := strconv.ParseInt(e.Property, 10, 64); err == nil {
				constIdx := c.chunk.AddConstant(&evaluator.Integer{Value: intVal})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
			} else if slot, isLocal := c.resolveLocal(e.Property); isLocal {
				c.emitOp(bytecode.OpGetLocal, line)
				c.chunk.WriteUint16(slot, line)
			} else {
				constIdx := c.chunk.AddConstant(&evaluator.String{Value: e.Property})
				c.emitOp(bytecode.OpConstant, line)
				c.chunk.WriteUint16(constIdx, line)
			}
		} else {
			constIdx := c.chunk.AddConstant(&evaluator.String{Value: e.Property})
			c.emitOp(bytecode.OpConstant, line)
			c.chunk.WriteUint16(constIdx, line)
		}
		c.emitOp(bytecode.OpIndex, line)
		return nil

	case *uast.PipelineExpression:
		return c.compilePipeline(e.Left, e.Right, line)

	case *uast.CallExpression:
		if err := c.compileExpression(e.Callee); err != nil {
			return err
		}
		for _, arg := range e.Arguments {
			if err := c.compileExpression(arg); err != nil {
				return err
			}
		}
		c.emitOp(bytecode.OpCall, line)
		c.chunk.Write(byte(len(e.Arguments)), line)
		return nil

	case *uast.FunctionExpression:
		fnComp := newFunctionCompiler(c, "", len(e.Parameters))
		for _, p := range e.Parameters {
			fnComp.addLocal(p.Name)
		}
		if e.Body != nil {
			for _, child := range e.Body.Statements {
				if err := fnComp.compileStatement(child); err != nil {
					return err
				}
			}
		}
		fnComp.emitOp(bytecode.OpNull, line)
		fnComp.emitOp(bytecode.OpReturn, line)

		compiledFn := &bytecode.CompiledFunction{
			Name:      "",
			NumParams: len(e.Parameters),
			Chunk:     fnComp.chunk,
		}
		constIdx := c.chunk.AddConstant(compiledFn)
		c.emitOp(bytecode.OpConstant, line)
		c.chunk.WriteUint16(constIdx, line)
		return nil

	case *uast.TemplateLiteral:
		if len(e.Parts) == 0 {
			c.emitOp(bytecode.OpConstant, line)
			constIdx := c.chunk.AddConstant(&evaluator.String{Value: ""})
			c.chunk.WriteUint16(constIdx, line)
			return nil
		}
		if err := c.compileExpression(e.Parts[0]); err != nil {
			return err
		}
		for i := 1; i < len(e.Parts); i++ {
			if err := c.compileExpression(e.Parts[i]); err != nil {
				return err
			}
			c.emitOp(bytecode.OpAdd, line)
		}
		return nil

	default:
		return fmt.Errorf("ekspresi '%T' belum didukung di bytecode compiler", expr)
	}
}

func (c *Compiler) compilePipeline(left, right uast.Expression, line int) error {
	if call, ok := right.(*uast.CallExpression); ok {
		if err := c.compileExpression(call.Callee); err != nil {
			return err
		}
		if err := c.compileExpression(left); err != nil {
			return err
		}
		for _, arg := range call.Arguments {
			if err := c.compileExpression(arg); err != nil {
				return err
			}
		}
		c.emitOp(bytecode.OpCall, line)
		c.chunk.Write(byte(len(call.Arguments)+1), line)
		return nil
	}

	if err := c.compileExpression(right); err != nil {
		return err
	}
	if err := c.compileExpression(left); err != nil {
		return err
	}
	c.emitOp(bytecode.OpCall, line)
	c.chunk.Write(1, line)
	return nil
}

func (c *Compiler) emitOp(op bytecode.Opcode, line int) {
	c.chunk.WriteOp(op, line)
}

func (c *Compiler) emitJump(op bytecode.Opcode, line int) int {
	c.emitOp(op, line)
	c.chunk.Write(0xFF, line)
	c.chunk.Write(0xFF, line)
	return len(c.chunk.Instructions) - 2
}

func (c *Compiler) patchJump(offset int) {
	jumpDist := len(c.chunk.Instructions)
	c.chunk.Instructions[offset] = byte(jumpDist >> 8)
	c.chunk.Instructions[offset+1] = byte(jumpDist & 0xFF)
}

func (c *Compiler) emitLoop(loopStart int, line int) {
	c.emitOp(bytecode.OpLoop, line)
	c.chunk.WriteUint16(uint16(loopStart), line)
}

func (c *Compiler) beginScope() {
	c.scopeDepth++
}

func (c *Compiler) endScope(line int) {
	c.scopeDepth--
	for len(c.locals) > 0 && c.locals[len(c.locals)-1].Depth > c.scopeDepth {
		c.emitOp(bytecode.OpPop, line)
		c.locals = c.locals[:len(c.locals)-1]
	}
}

func (c *Compiler) addLocal(name string) {
	c.locals = append(c.locals, Local{
		Name:  name,
		Depth: c.scopeDepth,
	})
}

func (c *Compiler) resolveLocal(name string) (uint16, bool) {
	for i := len(c.locals) - 1; i >= 0; i-- {
		if c.locals[i].Name == name {
			return uint16(i), true
		}
	}
	return 0, false
}
