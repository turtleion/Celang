package evaluator

import (
	"fmt"
	"sync"

	"celang/pkg/uast"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

// Eval evaluates a U-AST node within the given Environment.
func Eval(node uast.Node, env *Environment) Object {
	if node == nil {
		return NULL
	}

	switch n := node.(type) {
	// Root / Structural
	case *uast.Program:
		return evalProgram(n.Body, env)

	case *uast.BlockStatement:
		return evalBlockStatement(n.Statements, env)

	case *uast.ExpressionStatement:
		return Eval(n.Expression, env)

	// Statements
	case *uast.VariableDeclaration:
		var val Object = NULL
		if n.Value != nil {
			val = Eval(n.Value, env)
			if isError(val) {
				return val
			}
		}
		if _, err := env.Define(n.Name, val, n.IsMutable, n.TypeHint); err != nil {
			return newError(n.GetMetadata(), err.Error(), "Gunakan nama variabel yang berbeda.")
		}
		return val

	case *uast.AssignmentStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		switch target := n.Target.(type) {
		case *uast.Identifier:
			// Handle "+=", "-=", etc.
			if n.Operator != "=" && len(n.Operator) > 1 && n.Operator[len(n.Operator)-1] == '=' {
				baseOp := n.Operator[:len(n.Operator)-1]
				curVal, ok := env.Get(target.Name)
				if !ok {
					return newError(target.GetMetadata(), fmt.Sprintf("Variabel '%s' belum didefinisikan", target.Name), "Deklarasikan variabel terlebih dahulu sebelum melakukan assignment.")
				}
				val = evalBinaryOp(baseOp, curVal, val, n.GetMetadata())
				if isError(val) {
					return val
				}
			}
			if err := env.Assign(target.Name, val); err != nil {
				return newError(n.GetMetadata(), err.Error(), "Pastikan variabel ini dideklarasikan dengan 'mutable: true'.")
			}
			return val
		default:
			return newError(n.GetMetadata(), "Target assignment tidak valid", "Assignment hanya dapat dilakukan ke variabel atau elemen objek.")
		}

	case *uast.IfStatement:
		condVal := Eval(n.Condition, env)
		if isError(condVal) {
			return condVal
		}
		if CoerceToBoolean(condVal) {
			return Eval(n.Consequent, env)
		} else if n.Alternate != nil {
			return Eval(n.Alternate, env)
		}
		return NULL

	case *uast.WhileStatement:
		var result Object = NULL
		for {
			condVal := Eval(n.Condition, env)
			if isError(condVal) {
				return condVal
			}
			if !CoerceToBoolean(condVal) {
				break
			}
			result = Eval(n.Body, env)
			if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
				return result
			}
		}
		return result

	case *uast.ForStatement:
		forEnv := NewEnclosedEnvironment(env)
		if n.Init != nil {
			initRes := Eval(n.Init, forEnv)
			if isError(initRes) {
				return initRes
			}
		}
		var result Object = NULL
		for {
			if n.Condition != nil {
				condVal := Eval(n.Condition, forEnv)
				if isError(condVal) {
					return condVal
				}
				if !CoerceToBoolean(condVal) {
					break
				}
			}
			result = Eval(n.Body, forEnv)
			if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
				return result
			}
			if n.Update != nil {
				updRes := Eval(n.Update, forEnv)
				if isError(updRes) {
					return updRes
				}
			}
		}
		return result

	case *uast.FunctionDeclaration:
		fn := &Function{
			Name:       n.Name,
			Parameters: n.Parameters,
			Body:       n.Body,
			Env:        env,
		}
		env.Define(n.Name, fn, false, "function")
		return fn

	case *uast.ReturnStatement:
		var val Object = NULL
		if n.Value != nil {
			val = Eval(n.Value, env)
			if isError(val) {
				return val
			}
		}
		return &ReturnValue{Value: val}

	case *uast.ConcurrencyStatement:
		// Natural Concurrency:
		// "sambil" -> asynchronous background Goroutine
		// "bersamaan" -> parallel fan-out synchronized with sync.WaitGroup
		if n.Mode == "sambil" || n.Mode == "async" {
			go func() {
				subEnv := NewEnclosedEnvironment(env)
				for _, task := range n.Tasks {
					Eval(task, subEnv)
				}
			}()
			return NULL
		} else if n.Mode == "bersamaan" || n.Mode == "parallel" {
			var wg sync.WaitGroup
			for _, task := range n.Tasks {
				wg.Add(1)
				t := task
				go func() {
					defer wg.Done()
					subEnv := NewEnclosedEnvironment(env)
					Eval(t, subEnv)
				}()
			}
			wg.Wait()
			return NULL
		}
		return newError(n.GetMetadata(), fmt.Sprintf("Mode konkurensi tidak dikenal: '%s'", n.Mode), "Gunakan mode 'sambil' atau 'bersamaan'.")

	// Expressions
	case *uast.Literal:
		switch n.Kind {
		case uast.LiteralKindInt:
			if v, ok := n.Value.(int64); ok {
				return &Integer{Value: v}
			}
			if f, ok := n.Value.(float64); ok {
				return &Integer{Value: int64(f)}
			}
		case uast.LiteralKindFloat:
			if v, ok := n.Value.(float64); ok {
				return &Float{Value: v}
			}
		case uast.LiteralKindString:
			if s, ok := n.Value.(string); ok {
				return &String{Value: s}
			}
		case uast.LiteralKindBool:
			if b, ok := n.Value.(bool); ok {
				if b {
					return TRUE
				}
				return FALSE
			}
		case uast.LiteralKindNull:
			return NULL
		}
		return NULL

	case *uast.Identifier:
		val, ok := env.Get(n.Name)
		if !ok {
			return newError(n.GetMetadata(), fmt.Sprintf("Variabel '%s' tidak ditemukan", n.Name), "Periksa ejaan nama variabel atau pastikan sudah dideklarasikan sebelumnya.")
		}
		return val

	case *uast.MemberExpression:
		// Check for Time-Travel Variable access (e.g. x.sebelumnya, x.awal, x.histori)
		if id, ok := n.Object.(*uast.Identifier); ok {
			entry, exists := env.GetEntry(id.Name)
			if exists {
				switch n.Property {
				case "sebelumnya", "prev", "previous":
					prev, hasPrev := entry.Previous()
					if !hasPrev {
						return NULL
					}
					return prev
				case "awal", "initial", "first":
					return entry.Initial()
				case "histori", "history":
					arr := &Array{Elements: make([]Object, len(entry.History))}
					copy(arr.Elements, entry.History)
					return arr
				}
			}
		}

		// Fallback to evaluating object
		obj := Eval(n.Object, env)
		if isError(obj) {
			return obj
		}
		if arr, ok := obj.(*Array); ok {
			idxObj := Eval(uast.NewIdentifier(n.Property), env)
			if intIdx, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
				if intIdx >= 0 && int(intIdx) < len(arr.Elements) {
					return arr.Elements[intIdx]
				}
				return newError(n.GetMetadata(), fmt.Sprintf("Indeks array %d di luar jangkauan (panjang: %d)", intIdx, len(arr.Elements)), "Gunakan indeks yang valid antara 0 hingga panjang array - 1.")
			}
		}
		return newError(n.GetMetadata(), fmt.Sprintf("Properti '%s' tidak ditemukan pada objek", n.Property), "Pastikan properti yang diakses valid.")

	case *uast.BinaryExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryOp(n.Operator, left, right, n.GetMetadata())

	case *uast.UnaryExpression:
		right := Eval(n.Argument, env)
		if isError(right) {
			return right
		}
		return evalUnaryOp(n.Operator, right, n.GetMetadata())

	case *uast.CallExpression:
		callee := Eval(n.Callee, env)
		if isError(callee) {
			return callee
		}
		args := make([]Object, len(n.Arguments))
		for i, argNode := range n.Arguments {
			argVal := Eval(argNode, env)
			if isError(argVal) {
				return argVal
			}
			args[i] = argVal
		}
		return applyFunction(callee, args, n.GetMetadata())

	case *uast.ArrayLiteral:
		elems := make([]Object, len(n.Elements))
		for i, elNode := range n.Elements {
			elVal := Eval(elNode, env)
			if isError(elVal) {
				return elVal
			}
			elems[i] = elVal
		}
		return &Array{Elements: elems}
	}

	return newError(node.GetMetadata(), fmt.Sprintf("Node AST tipe '%s' belum didukung oleh evaluator", node.Type()), "Pastikan node U-AST valid.")
}

func evalProgram(stmts []uast.Statement, env *Environment) Object {
	var result Object = NULL
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == RETURN_VALUE_OBJ {
				return result.(*ReturnValue).Value
			}
			if result.Type() == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalBlockStatement(stmts []uast.Statement, env *Environment) Object {
	var result Object = NULL
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalBinaryOp(op string, left, right Object, meta *uast.NodeMetadata) Object {
	// String Concatenation vs Smart Numeric Operations
	if op == "+" {
		// If both are pure strings, concatenate
		if leftStr, ok := left.(*String); ok {
			if rightStr, ok := right.(*String); ok {
				return &String{Value: leftStr.Value + rightStr.Value}
			}
			// String + Number: check if string is purely numeric
			if IsPurelyNumericString(leftStr.Value) {
				if _, rightF, _, ok := CoerceToNumeric(right); ok {
					_, leftF, isLeftInt, _ := CoerceToNumeric(left)
					if isLeftInt && right.Type() == INTEGER_OBJ {
						return &Integer{Value: int64(leftF) + right.(*Integer).Value}
					}
					return &Float{Value: leftF + rightF}
				}
			}
			// Non-numeric string concatenation: "Total: " + 50 -> "Total: 50"
			return &String{Value: leftStr.Value + right.Inspect()}
		}
		if _, ok := right.(*String); ok {
			rightStr := right.(*String)
			if IsPurelyNumericString(rightStr.Value) {
				if _, leftF, _, ok := CoerceToNumeric(left); ok {
					_, rightF, isRightInt, _ := CoerceToNumeric(right)
					if isRightInt && left.Type() == INTEGER_OBJ {
						return &Integer{Value: left.(*Integer).Value + int64(rightF)}
					}
					return &Float{Value: leftF + rightF}
				}
			}
			return &String{Value: left.Inspect() + rightStr.Value}
		}
	}

	// Smart Fallback Arithmetic (-, *, /, %) with adaptive coercion
	switch op {
	case "+", "-", "*", "/", "%":
		leftInt, leftFloat, leftIsInt, leftOk := CoerceToNumeric(left)
		rightInt, rightFloat, rightIsInt, rightOk := CoerceToNumeric(right)
		if !leftOk || !rightOk {
			return newError(meta, fmt.Sprintf("Operasi '%s' tidak dapat mengekstrak angka dari tipe %s dan %s", op, left.Type(), right.Type()), "Pastikan nilai dapat dikonversi ke numerik.")
		}

		// Integer arithmetic
		if leftIsInt && rightIsInt {
			switch op {
			case "+":
				return &Integer{Value: leftInt + rightInt}
			case "-":
				return &Integer{Value: leftInt - rightInt}
			case "*":
				return &Integer{Value: leftInt * rightInt}
			case "/":
				if rightInt == 0 {
					return newError(meta, "Pembagian dengan nol (zero division)", "Pastikan penyebut bernilai selain 0.")
				}
				return &Integer{Value: leftInt / rightInt}
			case "%":
				if rightInt == 0 {
					return newError(meta, "Operasi modulo dengan nol", "Pastikan penyebut bernilai selain 0.")
				}
				return &Integer{Value: leftInt % rightInt}
			}
		}

		// Floating point arithmetic
		switch op {
		case "+":
			return &Float{Value: leftFloat + rightFloat}
		case "-":
			return &Float{Value: leftFloat - rightFloat}
		case "*":
			return &Float{Value: leftFloat * rightFloat}
		case "/":
			if rightFloat == 0.0 {
				return newError(meta, "Pembagian dengan nol (zero division)", "Pastikan penyebut bernilai selain 0.")
			}
			return &Float{Value: leftFloat / rightFloat}
		}

	case "==", "!=":
		eq := left.Inspect() == right.Inspect()
		if op == "==" {
			return nativeBool(eq)
		}
		return nativeBool(!eq)

	case "<", "<=", ">", ">=":
		leftInt, leftFloat, leftIsInt, leftOk := CoerceToNumeric(left)
		rightInt, rightFloat, rightIsInt, rightOk := CoerceToNumeric(right)
		if leftOk && rightOk {
			if leftIsInt && rightIsInt {
				switch op {
				case "<":
					return nativeBool(leftInt < rightInt)
				case "<=":
					return nativeBool(leftInt <= rightInt)
				case ">":
					return nativeBool(leftInt > rightInt)
				case ">=":
					return nativeBool(leftInt >= rightInt)
				}
			}
			switch op {
			case "<":
				return nativeBool(leftFloat < rightFloat)
			case "<=":
				return nativeBool(leftFloat <= rightFloat)
			case ">":
				return nativeBool(leftFloat > rightFloat)
			case ">=":
				return nativeBool(leftFloat >= rightFloat)
			}
		}
		// Fallback to string comparison
		switch op {
		case "<":
			return nativeBool(left.Inspect() < right.Inspect())
		case "<=":
			return nativeBool(left.Inspect() <= right.Inspect())
		case ">":
			return nativeBool(left.Inspect() > right.Inspect())
		case ">=":
			return nativeBool(left.Inspect() >= right.Inspect())
		}

	case "&&", "and", "dan":
		return nativeBool(CoerceToBoolean(left) && CoerceToBoolean(right))

	case "||", "or", "atau":
		return nativeBool(CoerceToBoolean(left) || CoerceToBoolean(right))
	}

	return newError(meta, fmt.Sprintf("Operator binary tidak dikenal: '%s'", op), "Periksa kembali operator yang digunakan.")
}

func evalUnaryOp(op string, right Object, meta *uast.NodeMetadata) Object {
	switch op {
	case "!", "bukan", "not":
		return nativeBool(!CoerceToBoolean(right))
	case "-":
		if intVal, floatVal, isInt, ok := CoerceToNumeric(right); ok {
			if isInt {
				return &Integer{Value: -intVal}
			}
			return &Float{Value: -floatVal}
		}
		return newError(meta, fmt.Sprintf("Operator '-' tidak dapat diterapkan pada %s", right.Type()), "Gunakan angka untuk negasi numerik.")
	}
	return newError(meta, fmt.Sprintf("Operator unary '%s' tidak dikenal", op), "Operator unary yang didukung: '!', 'bukan', '-'.")
}

func applyFunction(callee Object, args []Object, meta *uast.NodeMetadata) Object {
	switch fn := callee.(type) {
	case *Builtin:
		return fn.Fn(args...)
	case *Function:
		fnEnv := NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			var argVal Object = NULL
			if i < len(args) {
				argVal = args[i]
			} else if param.DefaultValue != nil {
				argVal = Eval(param.DefaultValue, fnEnv)
			}
			fnEnv.Define(param.Name, argVal, true, param.TypeHint)
		}
		res := Eval(fn.Body, fnEnv)
		if ret, ok := res.(*ReturnValue); ok {
			return ret.Value
		}
		return res
	default:
		return newError(meta, fmt.Sprintf("Objek bukan fungsi: %s", callee.Type()), "Pastikan Anda memanggil fungsi yang sah.")
	}
}

func nativeBool(b bool) *Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isError(obj Object) bool {
	return obj != nil && obj.Type() == ERROR_OBJ
}

func newError(meta *uast.NodeMetadata, msg, hint string) *Error {
	var pos *uast.Position
	if meta != nil {
		pos = meta.Pos
	}
	return &Error{
		Message:      msg,
		Pos:          pos,
		FriendlyHint: hint,
	}
}
