package vm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"celang/pkg/bytecode"
	"celang/pkg/evaluator"
	"celang/pkg/uast"
)

var (
	activeVMEnv   *evaluator.Environment
	activeVMEnvMu sync.RWMutex
)

func getActiveEnv() *evaluator.Environment {
	activeVMEnvMu.RLock()
	defer activeVMEnvMu.RUnlock()
	return activeVMEnv
}

func setActiveEnv(env *evaluator.Environment) {
	activeVMEnvMu.Lock()
	defer activeVMEnvMu.Unlock()
	activeVMEnv = env
}

func init() {
	bytecode.SetVMRunner(func(fn *bytecode.CompiledFunction, args []evaluator.Object) evaluator.Object {
		env := getActiveEnv()
		if env == nil {
			env = evaluator.NewEnvironment()
			evaluator.RegisterBuiltins(env, os.Stdout)
		}
		virtualMachine := NewWithEnv(fn.Chunk, env)
		virtualMachine.frames[0].fn = fn
		virtualMachine.frames[0].basePointer = 0
		virtualMachine.sp = 0
		for _, arg := range args {
			_ = virtualMachine.push(arg)
		}
		res, err := virtualMachine.Run()
		if err != nil {
			return &evaluator.Error{Message: err.Error()}
		}
		return res
	})
}

const (
	StackMax  = 2048
	FramesMax = 256
)

// CallFrame tracks the call stack and instruction pointer for active function calls.
type CallFrame struct {
	fn          *bytecode.CompiledFunction
	ip          int
	basePointer int
}

// VM executes compiled bytecode chunks.
type VM struct {
	frames     [FramesMax]CallFrame
	frameCount int
	stack      [StackMax]evaluator.Object
	sp         int
	globals    map[string]evaluator.Object
	constants  map[string]bool
	env        *evaluator.Environment
}

// New creates a new VM initialized with the given root chunk.
func New(chunk *bytecode.Chunk) *VM {
	return NewWithEnv(chunk, nil)
}

// NewWithEnv creates a VM connected to an existing Celang Environment (for standard builtins & globals).
func NewWithEnv(chunk *bytecode.Chunk, env *evaluator.Environment) *VM {
	rootFn := &bytecode.CompiledFunction{
		Name:      "<main>",
		NumParams: 0,
		Chunk:     chunk,
	}

	vm := &VM{
		globals:   make(map[string]evaluator.Object),
		constants: make(map[string]bool),
		env:       env,
		sp:        0,
	}

	vm.frames[0] = CallFrame{
		fn:          rootFn,
		ip:          0,
		basePointer: 0,
	}
	vm.frameCount = 1

	return vm
}

// Push pushes an object onto the evaluation stack.
func (vm *VM) push(obj evaluator.Object) error {
	if vm.sp >= StackMax {
		return fmt.Errorf("stack overflow: tumpukan eksekusi melebihi batas maksimum (%d)", StackMax)
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

// Pop removes and returns the top object from the stack.
func (vm *VM) pop() evaluator.Object {
	if vm.sp == 0 {
		return &evaluator.Null{}
	}
	vm.sp--
	val := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil // allow GC
	return val
}

// Peek returns the object at distance from the top of the stack without popping it.
func (vm *VM) peek(distance int) evaluator.Object {
	idx := vm.sp - 1 - distance
	if idx < 0 || idx >= StackMax {
		return &evaluator.Null{}
	}
	return vm.stack[idx]
}

// Run starts bytecode execution.
func (vm *VM) Run() (evaluator.Object, error) {
	if vm.env != nil {
		setActiveEnv(vm.env)
	}
	for vm.frameCount > 0 {
		frame := &vm.frames[vm.frameCount-1]
		chunk := frame.fn.Chunk

		if frame.ip >= len(chunk.Instructions) {
			break
		}

		op := bytecode.Opcode(chunk.Instructions[frame.ip])
		frame.ip++

		switch op {
		case bytecode.OpConstant:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			vm.push(chunk.Constants[constIdx])

		case bytecode.OpNull:
			vm.push(&evaluator.Null{})

		case bytecode.OpTrue:
			vm.push(&evaluator.Boolean{Value: true})

		case bytecode.OpFalse:
			vm.push(&evaluator.Boolean{Value: false})

		case bytecode.OpPop:
			vm.pop()

		case bytecode.OpAdd:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalAdd(a, b)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpSub:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalSub(a, b)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpMul:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalMul(a, b)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpDiv:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalDiv(a, b)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpMod:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalMod(a, b)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpNeg:
			a := vm.pop()
			res, err := vm.evalNeg(a)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(&evaluator.Boolean{Value: a.Inspect() == b.Inspect()})

		case bytecode.OpNotEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(&evaluator.Boolean{Value: a.Inspect() != b.Inspect()})

		case bytecode.OpGreaterThan:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalCompare(a, b, ">")
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpGreaterEqual:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalCompare(a, b, ">=")
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpLessThan:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalCompare(a, b, "<")
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpLessEqual:
			b := vm.pop()
			a := vm.pop()
			res, err := vm.evalCompare(a, b, "<=")
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpNot:
			a := vm.pop()
			vm.push(&evaluator.Boolean{Value: !isTruthy(a)})

		case bytecode.OpDefineGlobal:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			name := chunk.Constants[constIdx].(*evaluator.String).Value
			val := vm.pop()
			vm.globals[name] = val
			if vm.constants != nil {
				delete(vm.constants, name)
			}
			if vm.env != nil {
				vm.env.Define(name, val, true, "")
			}

		case bytecode.OpDefineGlobalConst:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			name := chunk.Constants[constIdx].(*evaluator.String).Value
			val := vm.pop()
			vm.globals[name] = val
			if vm.constants != nil {
				vm.constants[name] = true
			}
			if vm.env != nil {
				vm.env.Define(name, val, false, "")
			}

		case bytecode.OpGetGlobal:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			name := chunk.Constants[constIdx].(*evaluator.String).Value
			var val evaluator.Object
			var exists bool
			if vm.env != nil {
				val, exists = vm.env.Get(name)
			}
			if !exists {
				val, exists = vm.globals[name]
			}
			if !exists {
				return nil, fmt.Errorf("variabel global '%s' belum didefinisikan", name)
			}
			vm.push(val)

		case bytecode.OpSetGlobal:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			name := chunk.Constants[constIdx].(*evaluator.String).Value
			val := vm.peek(0)
			if vm.constants != nil && vm.constants[name] {
				return nil, fmt.Errorf("variabel '%s' bersifat tetap (immutable) dan tidak dapat diubah", name)
			}
			if vm.env != nil {
				if entry, ok := vm.env.GetEntry(name); ok {
					root := entry.Resolve()
					if !root.IsMutable {
						return nil, fmt.Errorf("variabel '%s' bersifat tetap (immutable) dan tidak dapat diubah", name)
					}
				}
				if err := vm.env.Assign(name, val); err != nil {
					return nil, err
				}
				if entry, ok := vm.env.GetEntry(name); ok {
					root := entry.Resolve()
					vm.globals[root.Name] = root.Current()
				}
			} else {
				vm.globals[name] = val
			}

		case bytecode.OpGetLocal:
			slot := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			vm.push(vm.stack[frame.basePointer+int(slot)])

		case bytecode.OpSetLocal:
			slot := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			val := vm.peek(0)
			vm.stack[frame.basePointer+int(slot)] = val

		case bytecode.OpJumpIfFalse:
			target := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			cond := vm.peek(0)
			if !isTruthy(cond) {
				frame.ip = int(target)
			}

		case bytecode.OpJump:
			target := chunk.ReadUint16(frame.ip)
			frame.ip = int(target)

		case bytecode.OpLoop:
			target := chunk.ReadUint16(frame.ip)
			frame.ip = int(target)

		case bytecode.OpArray:
			count := int(chunk.ReadUint16(frame.ip))
			frame.ip += 2
			elements := make([]evaluator.Object, count)
			for i := count - 1; i >= 0; i-- {
				elements[i] = vm.pop()
			}
			vm.push(&evaluator.Array{Elements: elements})

		case bytecode.OpMap:
			count := int(chunk.ReadUint16(frame.ip))
			frame.ip += 2
			pairs := make(map[string]evaluator.Object, count)
			order := make([]string, 0, count)
			type kv struct {
				k string
				v evaluator.Object
			}
			items := make([]kv, count)
			for i := count - 1; i >= 0; i-- {
				v := vm.pop()
				k := vm.pop().Inspect()
				items[i] = kv{k: k, v: v}
			}
			for _, item := range items {
				pairs[item.k] = item.v
				order = append(order, item.k)
			}
			vm.push(&evaluator.Map{Pairs: pairs, Order: order})

		case bytecode.OpIndex:
			indexObj := vm.pop()
			targetObj := vm.pop()
			res, err := vm.evalIndex(targetObj, indexObj)
			if err != nil {
				return nil, err
			}
			vm.push(res)

		case bytecode.OpSetIndex:
			indexObj := vm.pop()
			targetObj := vm.pop()
			val := vm.peek(0)
			if err := vm.evalSetIndex(targetObj, indexObj, val); err != nil {
				return nil, err
			}

		case bytecode.OpCall:
			argCount := int(chunk.Instructions[frame.ip])
			frame.ip++
			if err := vm.callValue(argCount); err != nil {
				return nil, err
			}

		case bytecode.OpReturn:
			var result evaluator.Object = &evaluator.Null{}
			if vm.sp > frame.basePointer {
				result = vm.pop()
			}
			vm.frameCount--
			if vm.frameCount == 0 {
				return result, nil
			}
			// Discard callee and locals
			vm.sp = frame.basePointer - 1
			vm.push(result)

		case bytecode.OpEvalAST:
			constIdx := chunk.ReadUint16(frame.ip)
			frame.ip += 2
			wrapper := chunk.Constants[constIdx].(*evaluator.ASTNodeWrapper)
			env := vm.env
			if env == nil {
				env = evaluator.NewEnvironment()
				evaluator.RegisterBuiltins(env, os.Stdout)
				vm.env = env
			}
			res := evaluator.Eval(wrapper.Node, env)
			if errObj, ok := res.(*evaluator.Error); ok {
				return nil, fmt.Errorf("%s", errObj.Inspect())
			}
			if _, isExpr := wrapper.Node.(uast.Expression); isExpr {
				if res == nil {
					res = &evaluator.Null{}
				}
				vm.push(res)
			}

		case bytecode.OpTimeTravel:
			nameIdx := chunk.ReadUint16(frame.ip)
			propIdx := chunk.ReadUint16(frame.ip + 2)
			frame.ip += 4
			name := chunk.Constants[nameIdx].(*evaluator.String).Value
			prop := chunk.Constants[propIdx].(*evaluator.String).Value
			if vm.env == nil {
				vm.push(&evaluator.Null{})
				break
			}
			entry, exists := vm.env.GetEntry(name)
			if !exists {
				return nil, fmt.Errorf("variabel '%s' tidak ditemukan dalam pencarian histori", name)
			}
			switch prop {
			case "sebelumnya", "prev", "previous":
				prev, hasPrev := entry.Previous()
				if !hasPrev {
					vm.push(&evaluator.Null{})
				} else {
					vm.push(prev)
				}
			case "awal", "initial", "first":
				vm.push(entry.Initial())
			case "histori", "history":
				arr := &evaluator.Array{Elements: make([]evaluator.Object, len(entry.History))}
				copy(arr.Elements, entry.History)
				vm.push(arr)
			default:
				vm.push(&evaluator.Null{})
			}

		case bytecode.OpGetIter:
			iterable := vm.pop()
			if errObj, ok := iterable.(*evaluator.Error); ok {
				return nil, fmt.Errorf("%s", errObj.Inspect())
			}
			iter, err := evaluator.NewIterator(iterable)
			if err != nil {
				return nil, err
			}
			vm.push(iter)

		case bytecode.OpForIter:
			iterSlot := chunk.ReadUint16(frame.ip)
			target := chunk.ReadUint16(frame.ip + 2)
			frame.ip += 4

			iterObj := vm.stack[frame.basePointer+int(iterSlot)]
			iterator, ok := iterObj.(*evaluator.IteratorObject)
			if !ok {
				return nil, fmt.Errorf("objek di slot %d bukan iterator yang valid", iterSlot)
			}
			val, hasNext := iterator.Next()
			if !hasNext {
				frame.ip = int(target)
			} else {
				vm.push(val)
			}

		default:
			return nil, fmt.Errorf("opcode tidak dikenal: 0x%02X (%s)", byte(op), op.String())
		}
	}

	if vm.sp > 0 {
		return vm.pop(), nil
	}
	return &evaluator.Null{}, nil
}

func (vm *VM) callValue(argCount int) error {
	calleeIdx := vm.sp - 1 - argCount
	if calleeIdx < 0 {
		return fmt.Errorf("tumpukan korup: tidak dapat menemukan fungsi callee")
	}
	callee := vm.stack[calleeIdx]

	switch fn := callee.(type) {
	case *bytecode.CompiledFunction:
		if argCount != fn.NumParams {
			return fmt.Errorf("fungsi '%s' membutuhkan %d argumen, tetapi menerima %d", fn.Name, fn.NumParams, argCount)
		}
		if vm.frameCount >= FramesMax {
			return fmt.Errorf("stack overflow: kedalaman pemanggilan fungsi melebihi batas maksimum (%d)", FramesMax)
		}
		frame := &vm.frames[vm.frameCount]
		frame.fn = fn
		frame.ip = 0
		frame.basePointer = vm.sp - argCount
		vm.frameCount++
		return nil

	case *evaluator.Builtin:
		args := make([]evaluator.Object, argCount)
		for i := 0; i < argCount; i++ {
			args[i] = vm.stack[vm.sp-argCount+i]
		}
		vm.sp -= argCount + 1
		res := fn.Fn(args...)
		if errObj, ok := res.(*evaluator.Error); ok {
			return fmt.Errorf("%s", errObj.Inspect())
		}
		vm.push(res)
		return nil

	case *evaluator.Function:
		args := make([]evaluator.Object, argCount)
		for i := 0; i < argCount; i++ {
			args[i] = vm.stack[vm.sp-argCount+i]
		}
		vm.sp -= argCount + 1
		res := evaluator.ApplyFunction(fn, args, nil)
		if errObj, ok := res.(*evaluator.Error); ok {
			return fmt.Errorf("%s", errObj.Inspect())
		}
		vm.push(res)
		return nil

	case *evaluator.StructDefinition:
		args := make([]evaluator.Object, argCount)
		for i := 0; i < argCount; i++ {
			args[i] = vm.stack[vm.sp-argCount+i]
		}
		vm.sp -= argCount + 1
		res := evaluator.ApplyFunction(fn, args, nil)
		if errObj, ok := res.(*evaluator.Error); ok {
			return fmt.Errorf("%s", errObj.Inspect())
		}
		vm.push(res)
		return nil

	default:
		return fmt.Errorf("objek bertipe %s bukan merupakan fungsi yang dapat dipanggil", callee.Type())
	}
}

func (vm *VM) evalIndex(targetObj, indexObj evaluator.Object) (evaluator.Object, error) {
	switch obj := targetObj.(type) {
	case *evaluator.Array:
		if sVal, ok := indexObj.(*evaluator.String); ok {
			if sVal.Value == "tambah" || sVal.Value == "push" {
				return &evaluator.Builtin{
					Fn: func(args ...evaluator.Object) evaluator.Object {
						obj.Elements = append(obj.Elements, args...)
						return obj
					},
				}, nil
			}
		}
		var idx int = -1
		if iVal, ok := indexObj.(*evaluator.Integer); ok {
			idx = int(iVal.Value)
		} else if sVal, ok := indexObj.(*evaluator.String); ok {
			if parsed, err := strconv.Atoi(sVal.Value); err == nil {
				idx = parsed
			} else {
				// Fallback to variable lookup in globals or env
				varObj, exists := vm.globals[sVal.Value]
				if !exists && vm.env != nil {
					varObj, exists = vm.env.Get(sVal.Value)
				}
				if exists {
					if nVal, _, isInt, ok := evaluator.CoerceToNumeric(varObj); ok && isInt {
						idx = int(nVal)
					}
				}
			}
		}

		if idx >= 0 && idx < len(obj.Elements) {
			return obj.Elements[idx], nil
		}
		return nil, fmt.Errorf("indeks array %d di luar jangkauan (panjang: %d)", idx, len(obj.Elements))

	case *evaluator.String:
		var idx int = -1
		if iVal, ok := indexObj.(*evaluator.Integer); ok {
			idx = int(iVal.Value)
		} else if sVal, ok := indexObj.(*evaluator.String); ok {
			if parsed, err := strconv.Atoi(sVal.Value); err == nil {
				idx = parsed
			} else {
				varObj, exists := vm.globals[sVal.Value]
				if !exists && vm.env != nil {
					varObj, exists = vm.env.Get(sVal.Value)
				}
				if exists {
					if nVal, _, isInt, ok := evaluator.CoerceToNumeric(varObj); ok && isInt {
						idx = int(nVal)
					}
				}
			}
		}

		runes := []rune(obj.Value)
		if idx >= 0 && idx < len(runes) {
			return &evaluator.String{Value: string(runes[idx])}, nil
		}
		return nil, fmt.Errorf("indeks teks %d di luar jangkauan (panjang: %d)", idx, len(runes))

	case *evaluator.Map:
		key := indexObj.Inspect()
		if val, exists := obj.Pairs[key]; exists {
			return val, nil
		}
		if sVal, ok := indexObj.(*evaluator.String); ok {
			varObj, exists := vm.globals[sVal.Value]
			if !exists && vm.env != nil {
				varObj, exists = vm.env.Get(sVal.Value)
			}
			if exists {
				if val, found := obj.Pairs[varObj.Inspect()]; found {
					return val, nil
				}
			}
		}
		return &evaluator.Null{}, nil

	case *evaluator.StructInstance:
		key := indexObj.Inspect()
		if val, exists := obj.Fields[key]; exists {
			return val, nil
		}
		if method, exists := obj.Def.Methods[key]; exists {
			boundEnv := evaluator.NewEnclosedEnvironment(method.Env)
			boundEnv.Define("diri", obj, false, obj.Def.Name)
			boundEnv.Define("this", obj, false, obj.Def.Name)
			boundEnv.Define("self", obj, false, obj.Def.Name)
			return &evaluator.Function{
				Name:       method.Name,
				Parameters: method.Parameters,
				Body:       method.Body,
				Env:        boundEnv,
			}, nil
		}
		return &evaluator.Null{}, nil

	case *evaluator.Channel:
		key := indexObj.Inspect()
		if s, ok := indexObj.(*evaluator.String); ok {
			key = s.Value
		}
		switch key {
		case "kirim":
			return &evaluator.Builtin{
				Fn: func(args ...evaluator.Object) evaluator.Object {
					if len(args) == 0 {
						return &evaluator.Error{Message: "kirim() membutuhkan 1 argumen data"}
					}
					return obj.Send(args[0], nil)
				},
			}, nil
		case "terima":
			return &evaluator.Builtin{
				Fn: func(args ...evaluator.Object) evaluator.Object {
					return obj.Receive(0)
				},
			}, nil
		case "tutup":
			return &evaluator.Builtin{
				Fn: func(args ...evaluator.Object) evaluator.Object {
					return obj.Close()
				},
			}, nil
		case "panjang":
			return &evaluator.Integer{Value: int64(len(obj.Ch))}, nil
		case "kapasitas":
			return &evaluator.Integer{Value: int64(obj.Capacity)}, nil
		}
		return &evaluator.Null{}, nil

	case *evaluator.Ticket:
		key := indexObj.Inspect()
		if s, ok := indexObj.(*evaluator.String); ok {
			key = s.Value
		}
		switch key {
		case "hasil", "jemput":
			return &evaluator.Builtin{
				Fn: func(args ...evaluator.Object) evaluator.Object {
					return obj.Await(0, nil)
				},
			}, nil
		}
		return &evaluator.Null{}, nil

	default:
		return nil, fmt.Errorf("tipe %s tidak mendukung operasi indeks", targetObj.Type())
	}
}

func (vm *VM) evalSetIndex(targetObj, indexObj, val evaluator.Object) error {
	switch obj := targetObj.(type) {
	case *evaluator.Array:
		var idx int = -1
		if iVal, ok := indexObj.(*evaluator.Integer); ok {
			idx = int(iVal.Value)
		} else if sVal, ok := indexObj.(*evaluator.String); ok {
			if parsed, err := strconv.Atoi(sVal.Value); err == nil {
				idx = parsed
			} else {
				varObj, exists := vm.globals[sVal.Value]
				if !exists && vm.env != nil {
					varObj, exists = vm.env.Get(sVal.Value)
				}
				if exists {
					if nVal, _, isInt, ok := evaluator.CoerceToNumeric(varObj); ok && isInt {
						idx = int(nVal)
					}
				}
			}
		}

		if idx >= 0 && idx < len(obj.Elements) {
			obj.Elements[idx] = val
			return nil
		}
		return fmt.Errorf("indeks array %d di luar jangkauan (panjang: %d)", idx, len(obj.Elements))

	case *evaluator.Map:
		key := indexObj.Inspect()
		obj.Pairs[key] = val
		found := false
		for _, k := range obj.Order {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			obj.Order = append(obj.Order, key)
		}
		return nil

	case *evaluator.StructInstance:
		key := indexObj.Inspect()
		obj.Fields[key] = val
		return nil

	default:
		return fmt.Errorf("tipe %s tidak mendukung assignment indeks", targetObj.Type())
	}
}

func (vm *VM) evalAdd(a, b evaluator.Object) (evaluator.Object, error) {
	if aStr, ok := a.(*evaluator.String); ok {
		return &evaluator.String{Value: aStr.Value + b.Inspect()}, nil
	}
	if bStr, ok := b.(*evaluator.String); ok {
		return &evaluator.String{Value: a.Inspect() + bStr.Value}, nil
	}

	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, bFloat, bIsInt, bOk := evaluator.CoerceToNumeric(b)
	if !aOk || !bOk {
		return nil, fmt.Errorf("operasi penambahan '+' tidak valid untuk tipe %s dan %s", a.Type(), b.Type())
	}

	if aIsInt && bIsInt {
		return &evaluator.Integer{Value: aInt + bInt}, nil
	}
	return &evaluator.Float{Value: aFloat + bFloat}, nil
}

func (vm *VM) evalSub(a, b evaluator.Object) (evaluator.Object, error) {
	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, bFloat, bIsInt, bOk := evaluator.CoerceToNumeric(b)
	if !aOk || !bOk {
		return nil, fmt.Errorf("operasi pengurangan '-' tidak valid untuk tipe %s dan %s", a.Type(), b.Type())
	}

	if aIsInt && bIsInt {
		return &evaluator.Integer{Value: aInt - bInt}, nil
	}
	return &evaluator.Float{Value: aFloat - bFloat}, nil
}

func (vm *VM) evalMul(a, b evaluator.Object) (evaluator.Object, error) {
	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, bFloat, bIsInt, bOk := evaluator.CoerceToNumeric(b)
	if !aOk || !bOk {
		return nil, fmt.Errorf("operasi perkalian '*' tidak valid untuk tipe %s dan %s", a.Type(), b.Type())
	}

	if aIsInt && bIsInt {
		return &evaluator.Integer{Value: aInt * bInt}, nil
	}
	return &evaluator.Float{Value: aFloat * bFloat}, nil
}

func (vm *VM) evalDiv(a, b evaluator.Object) (evaluator.Object, error) {
	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, bFloat, bIsInt, bOk := evaluator.CoerceToNumeric(b)
	if !aOk || !bOk {
		return nil, fmt.Errorf("operasi pembagian '/' tidak valid untuk tipe %s dan %s", a.Type(), b.Type())
	}

	if aIsInt && bIsInt {
		if bInt == 0 {
			return nil, fmt.Errorf("pembagian dengan nol (zero division)")
		}
		return &evaluator.Integer{Value: aInt / bInt}, nil
	}

	if bFloat == 0.0 {
		return nil, fmt.Errorf("pembagian dengan nol (zero division)")
	}
	return &evaluator.Float{Value: aFloat / bFloat}, nil
}

func (vm *VM) evalMod(a, b evaluator.Object) (evaluator.Object, error) {
	aInt, _, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, _, bIsInt, bOk := evaluator.CoerceToNumeric(b)
	if !aOk || !bOk || !aIsInt || !bIsInt {
		return nil, fmt.Errorf("operasi modulo '%%' hanya valid untuk integer")
	}
	if bInt == 0 {
		return nil, fmt.Errorf("operasi modulo dengan nol")
	}
	return &evaluator.Integer{Value: aInt % bInt}, nil
}

func (vm *VM) evalNeg(a evaluator.Object) (evaluator.Object, error) {
	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	if !aOk {
		return nil, fmt.Errorf("operasi negasi '-' tidak valid untuk tipe %s", a.Type())
	}
	if aIsInt {
		return &evaluator.Integer{Value: -aInt}, nil
	}
	return &evaluator.Float{Value: -aFloat}, nil
}

func (vm *VM) evalCompare(a, b evaluator.Object, op string) (evaluator.Object, error) {
	aInt, aFloat, aIsInt, aOk := evaluator.CoerceToNumeric(a)
	bInt, bFloat, bIsInt, bOk := evaluator.CoerceToNumeric(b)

	if aOk && bOk {
		if aIsInt && bIsInt {
			switch op {
			case ">":
				return &evaluator.Boolean{Value: aInt > bInt}, nil
			case ">=":
				return &evaluator.Boolean{Value: aInt >= bInt}, nil
			case "<":
				return &evaluator.Boolean{Value: aInt < bInt}, nil
			case "<=":
				return &evaluator.Boolean{Value: aInt <= bInt}, nil
			}
		}
		switch op {
		case ">":
			return &evaluator.Boolean{Value: aFloat > bFloat}, nil
		case ">=":
			return &evaluator.Boolean{Value: aFloat >= bFloat}, nil
		case "<":
			return &evaluator.Boolean{Value: aFloat < bFloat}, nil
		case "<=":
			return &evaluator.Boolean{Value: aFloat <= bFloat}, nil
		}
	}

	// Fallback to string comparison
	switch op {
	case ">":
		return &evaluator.Boolean{Value: a.Inspect() > b.Inspect()}, nil
	case ">=":
		return &evaluator.Boolean{Value: a.Inspect() >= b.Inspect()}, nil
	case "<":
		return &evaluator.Boolean{Value: a.Inspect() < b.Inspect()}, nil
	case "<=":
		return &evaluator.Boolean{Value: a.Inspect() <= b.Inspect()}, nil
	default:
		return nil, fmt.Errorf("operator perbandingan '%s' tidak valid", op)
	}
}

func isTruthy(obj evaluator.Object) bool {
	switch val := obj.(type) {
	case *evaluator.Boolean:
		return val.Value
	case *evaluator.Null:
		return false
	case *evaluator.Integer:
		return val.Value != 0
	case *evaluator.Float:
		return val.Value != 0.0
	case *evaluator.String:
		return strings.TrimSpace(val.Value) != ""
	default:
		return true
	}
}
