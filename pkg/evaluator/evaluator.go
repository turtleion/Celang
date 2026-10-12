package evaluator

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
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
		if _, err := env.DefineWithVisibility(n.Name, val, n.IsMutable, n.TypeHint, n.Visibility); err != nil {
			return newError(n.GetMetadata(), err.Error(), "Gunakan nama variabel yang berbeda.")
		}
		return val

	case *uast.DestructureDeclaration:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		switch n.Kind {
		case uast.DestructureArray:
			arr, isArr := val.(*Array)
			for i, b := range n.Bindings {
				var item Object = NULL
				if isArr && i < len(arr.Elements) {
					item = arr.Elements[i]
				}
				if _, err := env.DefineWithVisibility(b.Target, item, n.IsMutable, "", n.Visibility); err != nil {
					return newError(n.GetMetadata(), err.Error(), "")
				}
			}
		case uast.DestructureMap:
			mp, isMap := val.(*Map)
			inst, isInst := val.(*StructInstance)
			for _, b := range n.Bindings {
				var item Object = NULL
				if isMap {
					if v, ok := mp.Pairs[b.SourceKey]; ok {
						item = v
					}
				} else if isInst {
					if v, ok := inst.Fields[b.SourceKey]; ok {
						item = v
					}
				}
				if _, err := env.DefineWithVisibility(b.Target, item, n.IsMutable, "", n.Visibility); err != nil {
					return newError(n.GetMetadata(), err.Error(), "")
				}
			}
		}
		return val

	case *uast.DestructureAssignment:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		switch n.Kind {
		case uast.DestructureArray:
			arr, isArr := val.(*Array)
			for i, b := range n.Bindings {
				var item Object = NULL
				if isArr && i < len(arr.Elements) {
					item = arr.Elements[i]
				}
				if _, exists := env.Get(b.Target); !exists {
					env.Define(b.Target, item, true, "")
				} else {
					if err := env.Assign(b.Target, item); err != nil {
						return newError(n.GetMetadata(), err.Error(), "")
					}
				}
			}
		case uast.DestructureMap:
			mp, isMap := val.(*Map)
			inst, isInst := val.(*StructInstance)
			for _, b := range n.Bindings {
				var item Object = NULL
				if isMap {
					if v, ok := mp.Pairs[b.SourceKey]; ok {
						item = v
					}
				} else if isInst {
					if v, ok := inst.Fields[b.SourceKey]; ok {
						item = v
					}
				}
				if _, exists := env.Get(b.Target); !exists {
					env.Define(b.Target, item, true, "")
				} else {
					if err := env.Assign(b.Target, item); err != nil {
						return newError(n.GetMetadata(), err.Error(), "")
					}
				}
			}
		}
		return val

	case *uast.StructDeclaration:
		methods := make(map[string]*Function)
		privateFields := make(map[string]bool)
		privateMethods := make(map[string]bool)
		var allFields []uast.StructField
		fieldNames := make(map[string]bool)

		// 1. Resolve embedded structs
		for _, embedName := range n.Embeds {
			embedObj, ok := env.Get(embedName)
			if !ok {
				return newError(n.GetMetadata(), fmt.Sprintf("Bentuk embed '%s' tidak ditemukan", embedName), "Pastikan bentuk yang di-embed sudah didefinisikan sebelumnya.")
			}
			embedDef, ok := embedObj.(*StructDefinition)
			if !ok {
				return newError(n.GetMetadata(), fmt.Sprintf("Simbol embed '%s' bukan sebuah bentuk (struct)", embedName), "")
			}
			for _, ef := range embedDef.Fields {
				if !fieldNames[ef.Name] {
					allFields = append(allFields, ef)
					fieldNames[ef.Name] = true
					if embedDef.PrivateFields[ef.Name] {
						privateFields[ef.Name] = true
					}
				}
			}
			for mName, em := range embedDef.Methods {
				methods[mName] = em
				if embedDef.PrivateMethods[mName] {
					privateMethods[mName] = true
				}
			}
		}

		// 2. Add outer struct fields
		for _, f := range n.Fields {
			if f.IsPrivate {
				privateFields[f.Name] = true
			}
			if !fieldNames[f.Name] {
				allFields = append(allFields, f)
				fieldNames[f.Name] = true
			} else {
				for idx, ef := range allFields {
					if ef.Name == f.Name {
						allFields[idx] = f
						break
					}
				}
			}
		}

		// 3. Add outer struct methods
		for _, m := range n.Methods {
			if m.IsPrivate {
				privateMethods[m.Name] = true
			}
			methods[m.Name] = &Function{
				Name:       m.Name,
				Parameters: m.Parameters,
				Body:       m.Body,
				Env:        env,
			}
		}

		structDef := &StructDefinition{
			Name:           n.Name,
			Fields:         allFields,
			Methods:        methods,
			PrivateFields:  privateFields,
			PrivateMethods: privateMethods,
			Embeds:         n.Embeds,
			Env:            env,
		}
		if _, err := env.DefineWithVisibility(n.Name, structDef, false, "bentuk", n.Visibility); err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}
		return structDef

	case *uast.ReferenceDeclaration:
		id, ok := n.Source.(*uast.Identifier)
		if !ok {
			return newError(n.GetMetadata(), "Sumber rujuk harus berupa variabel (identifier)", "Contoh penggunaan: rujuk alias = target")
		}
		targetEntry, ok := env.GetEntry(id.Name)
		if !ok {
			return newError(n.GetMetadata(), fmt.Sprintf("Variabel '%s' yang dirujuk tidak ditemukan", id.Name), "Deklarasikan variabel target terlebih dahulu sebelum membuat rujukan.")
		}
		entry, err := env.DefineReference(n.Target, targetEntry)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}
		return entry.Current()

	case *uast.TypeAliasDeclaration:
		env.DefineTypeAlias(n.Name, n.BaseType)
		aliasName := n.Name
		baseType := n.BaseType
		aliasCtor := &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return &Error{Message: fmt.Sprintf("%s() mengharapkan 1 argumen, menerima %d", aliasName, len(args))}
				}
				return castValueToType(args[0], baseType, nil, env)
			},
		}
		if _, err := env.DefineWithVisibility(n.Name, aliasCtor, false, "tipe", nil); err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}
		return &String{Value: n.Name}

	case *uast.CastExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		return castValueToType(val, n.TargetType, n.GetMetadata(), env)

	case *uast.DeferStatement:
		env.PushDefer(n.Body)
		return NULL

	case *uast.SafeNavigationExpression:
		return evalSafeNavigation(n.Value, env)

	case *uast.NullCoalesceExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		if left != NULL {
			return left
		}
		return Eval(n.Right, env)

	case *uast.RangeExpression:
		return evalRangeExpression(n, env)

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
			if _, exists := env.Get(target.Name); !exists && n.Operator == "=" {
				env.Define(target.Name, val, true, "")
				return val
			}
			if err := env.Assign(target.Name, val); err != nil {
				return newError(n.GetMetadata(), err.Error(), "Pastikan variabel ini dideklarasikan dengan 'mutable: true'.")
			}
			return val

		case *uast.UnaryExpression:
			if target.Operator == "*" {
				ptrObj := Eval(target.Argument, env)
				if isError(ptrObj) {
					return ptrObj
				}
				ptr, ok := ptrObj.(*Pointer)
				if !ok {
					return newError(n.GetMetadata(), "Dereferencing target bukan pointer", "Pastikan target dereference adalah variabel pointer.")
				}
				if err := ptr.Env.Assign(ptr.VarName, val); err != nil {
					return newError(n.GetMetadata(), err.Error(), "")
				}
				return val
			}
			return newError(n.GetMetadata(), "Target assignment tidak valid", "Assignment hanya dapat dilakukan ke variabel atau elemen pointer.")

		case *uast.MemberExpression:
			obj := Eval(target.Object, env)
			if isError(obj) {
				return obj
			}
			prop := target.Property
			if inst, ok := obj.(*StructInstance); ok {
				if inst.Def.PrivateFields[prop] && !isInternalAccess(inst, env) {
					return newError(n.GetMetadata(),
						fmt.Sprintf("Properti '%s' bersifat privat pada bentuk '%s' dan tidak dapat diubah dari luar", prop, inst.Def.Name),
						"Modifikasi properti privat hanya diperbolehkan dari dalam metode bentuk menggunakan 'diri' atau '@'.")
				}
				if n.Operator != "=" && len(n.Operator) > 1 && n.Operator[len(n.Operator)-1] == '=' {
					baseOp := n.Operator[:len(n.Operator)-1]
					curVal := inst.Fields[prop]
					if curVal == nil {
						curVal = NULL
					}
					val = evalBinaryOp(baseOp, curVal, val, n.GetMetadata())
					if isError(val) {
						return val
					}
				}
				inst.Fields[prop] = val
				return val
			}
			if mp, ok := obj.(*Map); ok {
				if n.Operator != "=" && len(n.Operator) > 1 && n.Operator[len(n.Operator)-1] == '=' {
					baseOp := n.Operator[:len(n.Operator)-1]
					curVal := mp.Pairs[prop]
					if curVal == nil {
						curVal = NULL
					}
					val = evalBinaryOp(baseOp, curVal, val, n.GetMetadata())
					if isError(val) {
						return val
					}
				}
				mp.Pairs[prop] = val
				found := false
				for _, k := range mp.Order {
					if k == prop {
						found = true
						break
					}
				}
				if !found {
					mp.Order = append(mp.Order, prop)
				}
				return val
			}
			if arr, ok := obj.(*Array); ok {
				var intIdx int64
				var hasInt bool
				if parsed, err := strconv.ParseInt(prop, 10, 64); err == nil {
					intIdx = parsed
					hasInt = true
				} else {
					idxObj := Eval(uast.NewIdentifier(prop), env)
					if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
						intIdx = iVal
						hasInt = true
					}
				}
				if hasInt {
					if intIdx >= 0 && int(intIdx) < len(arr.Elements) {
						if n.Operator != "=" && len(n.Operator) > 1 && n.Operator[len(n.Operator)-1] == '=' {
							baseOp := n.Operator[:len(n.Operator)-1]
							curVal := arr.Elements[intIdx]
							val = evalBinaryOp(baseOp, curVal, val, n.GetMetadata())
							if isError(val) {
								return val
							}
						}
						arr.Elements[intIdx] = val
						return val
					}
					return newError(n.GetMetadata(), fmt.Sprintf("Indeks array %d di luar jangkauan (panjang: %d)", intIdx, len(arr.Elements)), "")
				}
			}
			return newError(n.GetMetadata(), fmt.Sprintf("Target properti '%s' tidak dapat ditugaskan", prop), "")

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

	case *uast.SwitchStatement:
		discVal := Eval(n.Discriminant, env)
		if isError(discVal) {
			return discVal
		}
		var result Object = NULL
		matched := false
		var defaultCase *uast.CaseClause

		for i := range n.Cases {
			c := &n.Cases[i]
			if c.Value == nil {
				defaultCase = c
				continue
			}
			caseVal := Eval(c.Value, env)
			if isError(caseVal) {
				return caseVal
			}
			if evalBinaryOp("==", discVal, caseVal, n.GetMetadata()) == TRUE {
				matched = true
				result = evalBlockStatement(c.Body, env)
				break
			}
		}

		if !matched && defaultCase != nil {
			result = evalBlockStatement(defaultCase.Body, env)
		}
		return result

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

	case *uast.ForEachStatement:
		iterVal := Eval(n.Iterable, env)
		if isError(iterVal) {
			return iterVal
		}
		var elements []Object
		if arr, ok := iterVal.(*Array); ok {
			elements = arr.Elements
		} else if str, ok := iterVal.(*String); ok {
			for _, r := range str.Value {
				elements = append(elements, &String{Value: string(r)})
			}
		} else if mp, ok := iterVal.(*Map); ok {
			for _, k := range mp.Order {
				elements = append(elements, &String{Value: k})
			}
		} else {
			return newError(n.GetMetadata(), fmt.Sprintf("Objek tipe %s tidak dapat diiterasi", iterVal.Type()), "Pastikan menggunakan Array, String, atau Map.")
		}

		forEnv := NewEnclosedEnvironment(env)
		var result Object = NULL
		for _, el := range elements {
			forEnv.Define(n.Item, el, true, "")
			result = Eval(n.Body, forEnv)
			if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
				return result
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
		env.DefineWithVisibility(n.Name, fn, false, "function", n.Visibility)
		return fn

	case *uast.VisibilityStatement:
		for _, ident := range n.Identifiers {
			if err := env.SetVisibility(ident, &n.Visibility); err != nil {
				return newError(n.GetMetadata(), err.Error(), fmt.Sprintf("Pastikan simbol '%s' didefinisikan sebelum mengatur visibilitasnya.", ident))
			}
		}
		return NULL

	case *uast.ImportDeclaration:
		targetAbsPath, err := resolveModulePath(env.FilePath, n.Source)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "Pastikan path modul benar dan file tersebut ada.")
		}
		targetEnv, err := loadModule(targetAbsPath, env)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}

		callerFile := env.FilePath
		callerRoot := env.ModuleRoot
		targetFile := targetEnv.FilePath
		targetRoot := targetEnv.ModuleRoot

		for _, spec := range n.Specifiers {
			if spec.IsAll {
				// Namespace import: ambil * sebagai math dari "..."
				mapPairs := make(map[string]Object)
				mapOrder := []string{}
				for name, entry := range targetEnv.AllEntries() {
					if err := checkVisibilityAccess(callerFile, callerRoot, targetFile, targetRoot, name, entry.Visibility); err == nil {
						mapPairs[name] = entry.Current()
						mapOrder = append(mapOrder, name)
					}
				}
				obj := &Map{Pairs: mapPairs, Order: mapOrder}
				if _, err := env.Define(spec.Local, obj, true, "Map"); err != nil {
					return newError(n.GetMetadata(), err.Error(), "")
				}
			} else {
				entry, ok := targetEnv.GetEntry(spec.Imported)
				if !ok {
					return newError(n.GetMetadata(), fmt.Sprintf("Simbol '%s' tidak ditemukan di berkas '%s'", spec.Imported, filepath.Base(targetFile)), "Pastikan nama simbol dieja dengan benar.")
				}
				if err := checkVisibilityAccess(callerFile, callerRoot, targetFile, targetRoot, spec.Imported, entry.Visibility); err != nil {
					return newError(n.GetMetadata(), err.Error(), "Pastikan aturan visibilitas pada file target mengizinkan akses dari file ini.")
				}
				if _, err := env.Define(spec.Local, entry.Current(), true, entry.TypeHint); err != nil {
					return newError(n.GetMetadata(), err.Error(), "")
				}
			}
		}
		return NULL

	case *uast.ThrowStatement:
		val := Eval(n.Argument, env)
		if isError(val) {
			return val
		}
		msg := val.Inspect()
		if s, ok := val.(*String); ok {
			msg = s.Value
		}
		return newError(n.GetMetadata(), msg, "")

	case *uast.TryCatchStatement:
		var result Object = NULL
		blockRes := Eval(n.Block, env)

		if isError(blockRes) {
			if n.Handler != nil {
				catchEnv := NewEnclosedEnvironment(env)
				if n.Handler.Param != "" {
					errMsg := blockRes.Inspect()
					if errObj, ok := blockRes.(*Error); ok {
						errMsg = errObj.Message
					}
					catchEnv.Define(n.Handler.Param, &String{Value: errMsg}, true, "String")
				}
				handlerRes := Eval(n.Handler.Body, catchEnv)
				result = handlerRes
			} else {
				result = blockRes
			}
		} else {
			result = blockRes
		}

		if n.Finalizer != nil {
			finRes := Eval(n.Finalizer, env)
			if isError(finRes) || (finRes != nil && finRes.Type() == RETURN_VALUE_OBJ) {
				return finRes
			}
		}

		return result

	case *uast.StyleBlockStatement:
		var result Object = NULL
		if n.Body != nil {
			for _, stmt := range n.Body.Statements {
				result = Eval(stmt, env)
				if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
					return result
				}
			}
		}
		return result

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

	case *uast.TemplateLiteral:
		var sb strings.Builder
		for _, part := range n.Parts {
			res := Eval(part, env)
			if isError(res) {
				return res
			}
			if s, ok := res.(*String); ok {
				sb.WriteString(s.Value)
			} else {
				sb.WriteString(res.Inspect())
			}
		}
		return &String{Value: sb.String()}

	case *uast.Identifier:
		val, ok := env.Get(n.Name)
		if !ok {
			return newError(n.GetMetadata(), fmt.Sprintf("Variabel '%s' tidak ditemukan", n.Name), "Periksa ejaan nama variabel atau pastikan sudah dideklarasikan sebelumnya.")
		}
		return val

	case *uast.MemberExpression:
		obj := Eval(n.Object, env)
		if isError(obj) {
			return obj
		}
		if n.Safe && (obj == NULL || obj == nil) {
			return NULL
		}

		// 1. Array index access, slicing, or built-in methods (arr.tambah, arr.push)
		if arr, ok := obj.(*Array); ok {
			if n.Property == "tambah" || n.Property == "push" {
				return &Builtin{
					Fn: func(args ...Object) Object {
						arr.Elements = append(arr.Elements, args...)
						return arr
					},
				}
			}
			if rng, isRng := n.PropertyExpr.(*uast.RangeExpression); isRng {
				return evalArraySlice(arr, rng, env, n.GetMetadata())
			}
			var intIdx int64
			var hasInt bool
			if n.Computed && n.PropertyExpr != nil {
				idxObj := Eval(n.PropertyExpr, env)
				if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
					intIdx = iVal
					hasInt = true
				}
			} else if parsed, err := strconv.ParseInt(n.Property, 10, 64); err == nil {
				intIdx = parsed
				hasInt = true
			} else {
				idxObj := Eval(uast.NewIdentifier(n.Property), env)
				if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
					intIdx = iVal
					hasInt = true
				}
			}
			if hasInt {
				if intIdx >= 0 && int(intIdx) < len(arr.Elements) {
					return arr.Elements[intIdx]
				}
				return newError(n.GetMetadata(), fmt.Sprintf("Indeks array %d di luar jangkauan (panjang: %d)", intIdx, len(arr.Elements)), "Gunakan indeks yang valid antara 0 hingga panjang array - 1.")
			}
		}

		// 1b. String index access and slicing (e.g. str[0] or str[0..2])
		if str, ok := obj.(*String); ok {
			if rng, isRng := n.PropertyExpr.(*uast.RangeExpression); isRng {
				return evalStringSlice(str, rng, env, n.GetMetadata())
			}
			var intIdx int64
			var hasInt bool
			if n.Computed && n.PropertyExpr != nil {
				idxObj := Eval(n.PropertyExpr, env)
				if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
					intIdx = iVal
					hasInt = true
				}
			} else if parsed, err := strconv.ParseInt(n.Property, 10, 64); err == nil {
				intIdx = parsed
				hasInt = true
			}
			if hasInt {
				runes := []rune(str.Value)
				if intIdx >= 0 && int(intIdx) < len(runes) {
					return &String{Value: string(runes[intIdx])}
				}
				return newError(n.GetMetadata(), fmt.Sprintf("Indeks teks %d di luar jangkauan (panjang: %d)", intIdx, len(runes)), "")
			}
		}

		// 2. Struct instance field or method access
		if inst, ok := obj.(*StructInstance); ok {
			prop := n.Property
			if val, exists := inst.Fields[prop]; exists {
				if inst.Def.PrivateFields[prop] && !isInternalAccess(inst, env) {
					return newError(n.GetMetadata(),
						fmt.Sprintf("Properti '%s' bersifat privat pada bentuk '%s' dan tidak dapat diakses dari luar", prop, inst.Def.Name),
						"Akses hanya diperbolehkan dari dalam metode bentuk menggunakan 'diri' atau '@'.")
				}
				return val
			}
			if method, exists := inst.Def.Methods[prop]; exists {
				if inst.Def.PrivateMethods[prop] && !isInternalAccess(inst, env) {
					return newError(n.GetMetadata(),
						fmt.Sprintf("Metode '%s' bersifat privat pada bentuk '%s' dan tidak dapat dipanggil dari luar", prop, inst.Def.Name),
						"Metode privat hanya dapat dipanggil secara internal dari dalam metode bentuk.")
				}
				boundEnv := NewEnclosedEnvironment(method.Env)
				boundEnv.Define("diri", inst, false, inst.Def.Name)
				boundEnv.Define("this", inst, false, inst.Def.Name)
				boundEnv.Define("self", inst, false, inst.Def.Name)
				return &Function{
					Name:       method.Name,
					Parameters: method.Parameters,
					Body:       method.Body,
					Env:        boundEnv,
				}
			}
		}

		// 3. Map property access
		if mp, ok := obj.(*Map); ok {
			key := n.Property
			if val, exists := mp.Pairs[key]; exists {
				return val
			}
		}

		// 4. Time-Travel Variable access fallback (e.g. x.sebelumnya, x.awal, x.histori)
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

		// 5. Missing property on Map returns NULL
		if _, ok := obj.(*Map); ok {
			return NULL
		}

		// 6. Missing property on Struct reports friendly error
		if inst, ok := obj.(*StructInstance); ok {
			return newError(n.GetMetadata(), fmt.Sprintf("Properti atau fungsi '%s' tidak ditemukan pada bentuk '%s'", n.Property, inst.Def.Name), "Pastikan nama field atau fungsi dieja dengan benar.")
		}

		return newError(n.GetMetadata(), fmt.Sprintf("Properti '%s' tidak ditemukan pada objek", n.Property), "Pastikan properti yang diakses valid.")

	case *uast.BinaryExpression:
		if n.Operator == "|>" {
			leftVal := Eval(n.Left, env)
			if isError(leftVal) {
				return leftVal
			}
			if call, ok := n.Right.(*uast.CallExpression); ok {
				callee := Eval(call.Callee, env)
				if isError(callee) {
					return callee
				}
				args := make([]Object, 1+len(call.Arguments))
				args[0] = leftVal
				for i, argNode := range call.Arguments {
					argVal := Eval(argNode, env)
					if isError(argVal) {
						return argVal
					}
					args[1+i] = argVal
				}
				return applyFunction(callee, args, n.GetMetadata())
			}
			callee := Eval(n.Right, env)
			if isError(callee) {
				return callee
			}
			return applyFunction(callee, []Object{leftVal}, n.GetMetadata())
		}
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryOp(n.Operator, left, right, n.GetMetadata())

	case *uast.PipelineExpression:
		leftVal := Eval(n.Left, env)
		if isError(leftVal) {
			return leftVal
		}
		if call, ok := n.Right.(*uast.CallExpression); ok {
			callee := Eval(call.Callee, env)
			if isError(callee) {
				return callee
			}
			args := make([]Object, 1+len(call.Arguments))
			args[0] = leftVal
			for i, argNode := range call.Arguments {
				argVal := Eval(argNode, env)
				if isError(argVal) {
					return argVal
				}
				args[1+i] = argVal
			}
			return applyFunction(callee, args, n.GetMetadata())
		}
		callee := Eval(n.Right, env)
		if isError(callee) {
			return callee
		}
		return applyFunction(callee, []Object{leftVal}, n.GetMetadata())

	case *uast.FunctionExpression:
		return &Function{
			Name:       "<anon>",
			Parameters: n.Parameters,
			Body:       n.Body,
			Env:        env,
		}

	case *uast.UnaryExpression:
		if n.Operator == "&" {
			if id, ok := n.Argument.(*uast.Identifier); ok {
				return &Pointer{Env: env, VarName: id.Name}
			}
			return newError(n.GetMetadata(), "Operator '&' hanya dapat mengambil alamat dari variabel", "Gunakan '&nama_variabel'.")
		}
		if n.Operator == "*" {
			ptrObj := Eval(n.Argument, env)
			if isError(ptrObj) {
				return ptrObj
			}
			ptr, ok := ptrObj.(*Pointer)
			if !ok {
				return newError(n.GetMetadata(), fmt.Sprintf("Operator '*' tidak dapat mendereference tipe %s", ptrObj.Type()), "Pastikan objek yang didereference adalah pointer.")
			}
			val, ok := ptr.Env.Get(ptr.VarName)
			if !ok {
				return newError(n.GetMetadata(), fmt.Sprintf("Variabel target pointer '%s' tidak ditemukan", ptr.VarName), "")
			}
			return val
		}
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
		if fn, ok := callee.(*Function); ok {
			return applyFunctionWithASTArgs(fn, n.Arguments, env, n.GetMetadata())
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

	case *uast.MapLiteral:
		m := &Map{
			Pairs: make(map[string]Object),
			Order: make([]string, 0, len(n.Properties)),
		}
		for _, prop := range n.Properties {
			keyVal := Eval(prop.Key, env)
			if isError(keyVal) {
				return keyVal
			}
			valVal := Eval(prop.Value, env)
			if isError(valVal) {
				return valVal
			}
			keyStr := keyVal.Inspect()
			if strObj, ok := keyVal.(*String); ok {
				keyStr = strObj.Value
			}
			if _, exists := m.Pairs[keyStr]; !exists {
				m.Order = append(m.Order, keyStr)
			}
			m.Pairs[keyStr] = valVal
		}
		return m
	}

	return newError(node.GetMetadata(), fmt.Sprintf("Node AST tipe '%s' belum didukung oleh evaluator", node.Type()), "Pastikan node U-AST valid.")
}

func evalProgram(stmts []uast.Statement, env *Environment) Object {
	defer func() {
		for _, defStmt := range env.PopAllDefers() {
			Eval(defStmt, env)
		}
	}()
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

// ApplyFunction calls a function or builtin object with the given arguments.
func ApplyFunction(callee Object, args []Object, meta *uast.NodeMetadata) Object {
	return applyFunction(callee, args, meta)
}

func applyFunction(callee Object, args []Object, meta *uast.NodeMetadata) Object {
	switch fn := callee.(type) {
	case *Builtin:
		return fn.Fn(args...)
	case *Function:
		fnEnv := NewEnclosedEnvironment(fn.Env)
		defer func() {
			for _, defStmt := range fnEnv.PopAllDefers() {
				Eval(defStmt, fnEnv)
			}
		}()
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
	case *StructDefinition:
		instEnv := NewEnclosedEnvironment(fn.Env)
		instance := &StructInstance{
			Def:    fn,
			Fields: make(map[string]Object),
			Env:    instEnv,
		}
		instEnv.Define("diri", instance, false, fn.Name)
		instEnv.Define("this", instance, false, fn.Name)
		instEnv.Define("self", instance, false, fn.Name)

		if len(args) == 1 && args[0].Type() == MAP_OBJ {
			mp := args[0].(*Map)
			for _, f := range fn.Fields {
				if v, ok := mp.Pairs[f.Name]; ok {
					instance.Fields[f.Name] = v
				} else if f.DefaultValue != nil {
					defVal := Eval(f.DefaultValue, instEnv)
					if isError(defVal) {
						return defVal
					}
					instance.Fields[f.Name] = defVal
				} else {
					instance.Fields[f.Name] = NULL
				}
			}
		} else {
			for i, f := range fn.Fields {
				if i < len(args) {
					instance.Fields[f.Name] = args[i]
				} else if f.DefaultValue != nil {
					defVal := Eval(f.DefaultValue, instEnv)
					if isError(defVal) {
						return defVal
					}
					instance.Fields[f.Name] = defVal
				} else {
					instance.Fields[f.Name] = NULL
				}
			}
		}
		return instance
	case Callable:
		return fn.Call(args...)
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

func isInternalAccess(inst *StructInstance, env *Environment) bool {
	if env == nil || inst == nil {
		return false
	}
	selfObj, ok := env.Get("diri")
	if !ok {
		selfObj, ok = env.Get("this")
	}
	if !ok {
		selfObj, ok = env.Get("self")
	}
	if !ok {
		return false
	}
	if callerInst, ok := selfObj.(*StructInstance); ok {
		return callerInst == inst
	}
	return false
}

func castValueToType(val Object, targetType string, meta *uast.NodeMetadata, env *Environment) Object {
	if env != nil {
		if alias, ok := env.GetTypeAlias(targetType); ok {
			targetType = alias
		}
	}

	switch targetType {
	case "Angka", "angka", "Int", "int", "Integer":
		switch v := val.(type) {
		case *Integer:
			return v
		case *Float:
			return &Integer{Value: int64(v.Value)}
		case *String:
			if iv, err := strconv.ParseInt(v.Value, 10, 64); err == nil {
				return &Integer{Value: iv}
			}
			if fv, err := strconv.ParseFloat(v.Value, 64); err == nil {
				return &Integer{Value: int64(fv)}
			}
			return newError(meta, fmt.Sprintf("Gagal mengubah teks '%s' menjadi Angka", v.Value), "Pastikan format teks adalah angka yang valid.")
		case *Boolean:
			if v.Value {
				return &Integer{Value: 1}
			}
			return &Integer{Value: 0}
		default:
			return newError(meta, fmt.Sprintf("Tidak dapat mengubah tipe %s menjadi Angka", val.Type()), "")
		}

	case "Pecahan", "pecahan", "Float", "float":
		switch v := val.(type) {
		case *Float:
			return v
		case *Integer:
			return &Float{Value: float64(v.Value)}
		case *String:
			if fv, err := strconv.ParseFloat(v.Value, 64); err == nil {
				return &Float{Value: fv}
			}
			return newError(meta, fmt.Sprintf("Gagal mengubah teks '%s' menjadi Pecahan", v.Value), "Pastikan format teks adalah desimal yang valid.")
		case *Boolean:
			if v.Value {
				return &Float{Value: 1.0}
			}
			return &Float{Value: 0.0}
		default:
			return newError(meta, fmt.Sprintf("Tidak dapat mengubah tipe %s menjadi Pecahan", val.Type()), "")
		}

	case "Teks", "teks", "String", "string":
		switch v := val.(type) {
		case *String:
			return v
		default:
			return &String{Value: v.Inspect()}
		}

	case "Logika", "logika", "Bool", "bool", "Boolean":
		return nativeBool(CoerceToBoolean(val))

	default:
		return newError(meta, fmt.Sprintf("Tipe tujuan casting '%s' tidak didukung", targetType), "Tipe yang didukung: Angka, Pecahan, Teks, Logika, atau alias yang terdaftar.")
	}
}

func applyFunctionWithASTArgs(fn *Function, argNodes []uast.Expression, callerEnv *Environment, meta *uast.NodeMetadata) Object {
	fnEnv := NewEnclosedEnvironment(fn.Env)
	defer func() {
		for _, defStmt := range fnEnv.PopAllDefers() {
			Eval(defStmt, fnEnv)
		}
	}()
	for i, param := range fn.Parameters {
		if param.IsReference {
			if i >= len(argNodes) {
				return newError(meta, fmt.Sprintf("Parameter rujuk '%s' membutuhkan argumen", param.Name), "")
			}
			id, ok := argNodes[i].(*uast.Identifier)
			if !ok {
				return newError(meta,
					fmt.Sprintf("Argumen untuk parameter rujuk '%s' harus berupa variabel (contoh: %s(x)), bukan ekspresi atau literal", param.Name, fn.Name),
					"Buat variabel terlebih dahulu sebelum meneruskannya sebagai referensi.")
			}
			targetEntry, ok := callerEnv.GetEntry(id.Name)
			if !ok {
				return newError(meta, fmt.Sprintf("Variabel '%s' tidak ditemukan di lingkungan pemanggil", id.Name), "")
			}
			if _, err := fnEnv.DefineReference(param.Name, targetEntry); err != nil {
				return newError(meta, err.Error(), "")
			}
		} else {
			var argVal Object = NULL
			if i < len(argNodes) {
				argVal = Eval(argNodes[i], callerEnv)
				if isError(argVal) {
					return argVal
				}
			} else if param.DefaultValue != nil {
				argVal = Eval(param.DefaultValue, fnEnv)
				if isError(argVal) {
					return argVal
				}
			}
			fnEnv.Define(param.Name, argVal, true, param.TypeHint)
		}
	}
	res := Eval(fn.Body, fnEnv)
	if ret, ok := res.(*ReturnValue); ok {
		return ret.Value
	}
	return res
}

func evalSafeNavigation(expr uast.Expression, env *Environment) Object {
	if expr == nil {
		return NULL
	}

	switch e := expr.(type) {
	case *uast.MemberExpression:
		obj := evalSafeNavigation(e.Object, env)
		if obj == NULL || obj == nil || isError(obj) {
			return NULL
		}

		// Safe member access on obj
		if arr, ok := obj.(*Array); ok {
			if rng, isRng := e.PropertyExpr.(*uast.RangeExpression); isRng {
				return evalArraySlice(arr, rng, env, e.GetMetadata())
			}
			var intIdx int64
			var hasInt bool
			if e.Computed && e.PropertyExpr != nil {
				idxObj := evalSafeNavigation(e.PropertyExpr, env)
				if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
					intIdx = iVal
					hasInt = true
				}
			} else if parsed, err := strconv.ParseInt(e.Property, 10, 64); err == nil {
				intIdx = parsed
				hasInt = true
			}
			if hasInt && intIdx >= 0 && int(intIdx) < len(arr.Elements) {
				return arr.Elements[intIdx]
			}
			return NULL
		}

		if str, ok := obj.(*String); ok {
			if rng, isRng := e.PropertyExpr.(*uast.RangeExpression); isRng {
				return evalStringSlice(str, rng, env, e.GetMetadata())
			}
			return NULL
		}

		if inst, ok := obj.(*StructInstance); ok {
			prop := e.Property
			if val, exists := inst.Fields[prop]; exists {
				return val
			}
			if method, exists := inst.Def.Methods[prop]; exists {
				boundEnv := NewEnclosedEnvironment(method.Env)
				boundEnv.Define("diri", inst, false, inst.Def.Name)
				boundEnv.Define("this", inst, false, inst.Def.Name)
				boundEnv.Define("self", inst, false, inst.Def.Name)
				return &Function{
					Name:       method.Name,
					Parameters: method.Parameters,
					Body:       method.Body,
					Env:        boundEnv,
				}
			}
			return NULL
		}

		if mp, ok := obj.(*Map); ok {
			if val, exists := mp.Pairs[e.Property]; exists {
				return val
			}
			return NULL
		}

		return NULL

	case *uast.CallExpression:
		callee := evalSafeNavigation(e.Callee, env)
		if callee == NULL || callee == nil || isError(callee) {
			return NULL
		}
		var args []Object
		for _, arg := range e.Arguments {
			argVal := Eval(arg, env)
			if isError(argVal) {
				return NULL
			}
			args = append(args, argVal)
		}
		res := applyFunction(callee, args, e.GetMetadata())
		if isError(res) {
			return NULL
		}
		return res

	case *uast.Identifier:
		val, ok := env.Get(e.Name)
		if !ok || val == nil {
			return NULL
		}
		return val

	default:
		res := Eval(expr, env)
		if isError(res) {
			return NULL
		}
		return res
	}
}

func evalRangeExpression(rng *uast.RangeExpression, env *Environment) Object {
	var startVal int64 = 0
	var endVal int64 = 0
	if rng.Start != nil {
		s := Eval(rng.Start, env)
		if isError(s) {
			return s
		}
		if iVal, _, isInt, ok := CoerceToNumeric(s); ok && isInt {
			startVal = iVal
		} else {
			return newError(rng.GetMetadata(), "Batas awal rentang harus berupa bilangan bulat", "")
		}
	}
	if rng.End != nil {
		e := Eval(rng.End, env)
		if isError(e) {
			return e
		}
		if iVal, _, isInt, ok := CoerceToNumeric(e); ok && isInt {
			endVal = iVal
		} else {
			return newError(rng.GetMetadata(), "Batas akhir rentang harus berupa bilangan bulat", "")
		}
	}

	var elements []Object
	if startVal <= endVal {
		for i := startVal; i <= endVal; i++ {
			elements = append(elements, &Integer{Value: i})
		}
	} else {
		for i := startVal; i >= endVal; i-- {
			elements = append(elements, &Integer{Value: i})
		}
	}
	return &Array{Elements: elements}
}

func evalArraySlice(arr *Array, rng *uast.RangeExpression, env *Environment, meta *uast.NodeMetadata) Object {
	start := 0
	end := len(arr.Elements)
	if rng.Start != nil {
		sVal := Eval(rng.Start, env)
		if isError(sVal) {
			return sVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(sVal); ok && isInt {
			start = int(iVal)
		}
	}
	if rng.End != nil {
		eVal := Eval(rng.End, env)
		if isError(eVal) {
			return eVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(eVal); ok && isInt {
			end = int(iVal)
		}
	}
	if start < 0 {
		start = 0
	}
	if start > len(arr.Elements) {
		start = len(arr.Elements)
	}
	if end < start {
		end = start
	}
	if end > len(arr.Elements) {
		end = len(arr.Elements)
	}

	sliced := make([]Object, end-start)
	copy(sliced, arr.Elements[start:end])
	return &Array{Elements: sliced}
}

func evalStringSlice(str *String, rng *uast.RangeExpression, env *Environment, meta *uast.NodeMetadata) Object {
	runes := []rune(str.Value)
	start := 0
	end := len(runes)
	if rng.Start != nil {
		sVal := Eval(rng.Start, env)
		if isError(sVal) {
			return sVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(sVal); ok && isInt {
			start = int(iVal)
		}
	}
	if rng.End != nil {
		eVal := Eval(rng.End, env)
		if isError(eVal) {
			return eVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(eVal); ok && isInt {
			end = int(iVal)
		}
	}
	if start < 0 {
		start = 0
	}
	if start > len(runes) {
		start = len(runes)
	}
	if end < start {
		end = start
	}
	if end > len(runes) {
		end = len(runes)
	}

	return &String{Value: string(runes[start:end])}
}

