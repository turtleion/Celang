package evaluator

import (
	"bytes"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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
		if n.IsUniversal {
			env.IsUniversal = true
		}
		return evalProgram(n.Body, env)

	case *uast.BlockStatement:
		return evalBlockStatement(n.Statements, env)

	case *uast.ExpressionStatement:
		res := Eval(n.Expression, env)
		if isError(res) {
			if aErr, ok := res.(*AssertionError); ok && aErr.Pos == nil && n.GetMetadata() != nil {
				aErr.Pos = n.GetMetadata().Pos
			}
		}
		return res

	// Statements
	case *uast.VariableDeclaration:
		var val Object = NULL
		if n.Value != nil {
			val = Eval(n.Value, env)
			if isError(val) || (val != nil && val.Type() == ROLLBACK_SIGNAL_OBJ) {
				return val
			}
		}
		if existing, exists := env.GetLocalEntry(n.Name); exists && existing.IsConstant {
			return newError(n.GetMetadata(),
				fmt.Sprintf("variabel '%s' bersifat tetap (konstan) dan tidak dapat dideklarasikan ulang", n.Name),
				"Gunakan nama variabel yang berbeda.")
		}
		entry, err := env.DefineWithVisibility(n.Name, val, n.IsMutable, n.TypeHint, n.Visibility)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "Gunakan nama variabel yang berbeda.")
		}
		if !n.IsMutable {
			entry.IsConstant = true
		}
		return val

	case *uast.DestructureDeclaration:
		val := Eval(n.Value, env)
		if isError(val) || (val != nil && val.Type() == ROLLBACK_SIGNAL_OBJ) {
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
		if isError(val) || (val != nil && val.Type() == ROLLBACK_SIGNAL_OBJ) {
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
			TypeParameters: n.TypeParameters,
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
		if isError(val) || (val != nil && val.Type() == ROLLBACK_SIGNAL_OBJ) {
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
			if rootName := findRootIdentifier(target.Object); rootName != "" {
				env.RecordSavepointForVar(rootName)
			}
			obj := Eval(target.Object, env)
			if isError(obj) {
				return obj
			}
			prop := target.Property
			var computedIndexObj Object
			if target.Computed && target.PropertyExpr != nil {
				computedIndexObj = Eval(target.PropertyExpr, env)
				if isError(computedIndexObj) {
					return computedIndexObj
				}
				if s, ok := computedIndexObj.(*String); ok {
					prop = s.Value
				} else {
					prop = computedIndexObj.Inspect()
				}
			}

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
				if computedIndexObj != nil {
					if iVal, _, isInt, ok := CoerceToNumeric(computedIndexObj); ok && isInt {
						intIdx = iVal
						hasInt = true
					}
				} else if parsed, err := strconv.ParseInt(prop, 10, 64); err == nil {
					intIdx = parsed
					hasInt = true
				} else if target.Computed && prop != "" {
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
			if b, ok := obj.(*Bytes); ok {
				var intIdx int64
				var hasInt bool
				if computedIndexObj != nil {
					if iVal, _, isInt, ok := CoerceToNumeric(computedIndexObj); ok && isInt {
						intIdx = iVal
						hasInt = true
					}
				} else if parsed, err := strconv.ParseInt(prop, 10, 64); err == nil {
					intIdx = parsed
					hasInt = true
				} else if target.Computed && prop != "" {
					idxObj := Eval(uast.NewIdentifier(prop), env)
					if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
						intIdx = iVal
						hasInt = true
					}
				}
				if hasInt {
					if intIdx < 0 {
						intIdx = int64(len(b.Elements)) + intIdx
					}
					if intIdx >= 0 && int(intIdx) < len(b.Elements) {
						var byteVal byte
						if iVal, _, isInt, ok := CoerceToNumeric(val); ok && isInt {
							if iVal < 0 || iVal > 255 {
								return newError(n.GetMetadata(), fmt.Sprintf("Nilai bita harus di antara 0 dan 255, menerima %d", iVal), "")
							}
							byteVal = byte(iVal)
						} else {
							return newError(n.GetMetadata(), fmt.Sprintf("Nilai bita harus berupa angka (0..255), menerima %s", val.Type()), "")
						}
						b.Elements[intIdx] = byteVal
						return val
					}
					return newError(n.GetMetadata(), fmt.Sprintf("Indeks bita %d di luar jangkauan (panjang: %d)", intIdx, len(b.Elements)), "")
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
		if result != nil && result.Type() == BREAK_SIGNAL_OBJ {
			result = NULL
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
			if result != nil {
				if result.Type() == BREAK_SIGNAL_OBJ {
					break
				}
				if result.Type() == CONTINUE_SIGNAL_OBJ {
					continue
				}
				if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ {
					return result
				}
			}
		}
		if result != nil && (result.Type() == BREAK_SIGNAL_OBJ || result.Type() == CONTINUE_SIGNAL_OBJ) {
			return NULL
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
			if result != nil {
				if result.Type() == BREAK_SIGNAL_OBJ {
					break
				}
				if result.Type() == CONTINUE_SIGNAL_OBJ {
					if n.Update != nil {
						updRes := Eval(n.Update, forEnv)
						if isError(updRes) {
							return updRes
						}
					}
					continue
				}
				if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ {
					return result
				}
			}
			if n.Update != nil {
				updRes := Eval(n.Update, forEnv)
				if isError(updRes) {
					return updRes
				}
			}
		}
		if result != nil && (result.Type() == BREAK_SIGNAL_OBJ || result.Type() == CONTINUE_SIGNAL_OBJ) {
			return NULL
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
		} else if b, ok := iterVal.(*Bytes); ok {
			for _, by := range b.Elements {
				elements = append(elements, &Integer{Value: int64(by)})
			}
		} else if mp, ok := iterVal.(*Map); ok {
			if len(mp.Order) > 0 {
				for _, k := range mp.Order {
					elements = append(elements, &String{Value: k})
				}
			} else {
				for k := range mp.Pairs {
					elements = append(elements, &String{Value: k})
				}
			}
		} else if ch, ok := iterVal.(*Channel); ok {
			forEnv := NewEnclosedEnvironment(env)
			var result Object = NULL
			for {
				el, ok := ch.ReceiveWithStatus(0)
				if !ok {
					break
				}
				forEnv.Define(n.Item, el, true, "")
				result = Eval(n.Body, forEnv)
				if result != nil {
					if result.Type() == BREAK_SIGNAL_OBJ {
						break
					}
					if result.Type() == CONTINUE_SIGNAL_OBJ {
						continue
					}
					if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ {
						return result
					}
				}
			}
			if result != nil && (result.Type() == BREAK_SIGNAL_OBJ || result.Type() == CONTINUE_SIGNAL_OBJ) {
				return NULL
			}
			return result
		} else {
			return newError(n.GetMetadata(), fmt.Sprintf("Objek tipe %s tidak dapat diiterasi", iterVal.Type()), "Pastikan menggunakan Array, String, Map, atau Saluran.")
		}

		forEnv := NewEnclosedEnvironment(env)
		var result Object = NULL
		for _, el := range elements {
			forEnv.Define(n.Item, el, true, "")
			result = Eval(n.Body, forEnv)
			if result != nil {
				if result.Type() == BREAK_SIGNAL_OBJ {
					break
				}
				if result.Type() == CONTINUE_SIGNAL_OBJ {
					continue
				}
				if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ {
					return result
				}
			}
		}
		if result != nil && (result.Type() == BREAK_SIGNAL_OBJ || result.Type() == CONTINUE_SIGNAL_OBJ) {
			return NULL
		}
		return result

	case *uast.FunctionDeclaration:
		fn := &Function{
			Name:           n.Name,
			TypeParameters: n.TypeParameters,
			Parameters:     n.Parameters,
			Body:           n.Body,
			Env:            env,
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
		if env.IsUniversal {
			switch strings.ToLower(strings.TrimSpace(n.Source)) {
			case "berkas":
				return newError(n.GetMetadata(), "Indonesian module 'berkas' is not allowed in Universal Mode.", "Import 'fs' instead.")
			case "matematika":
				return newError(n.GetMetadata(), "Indonesian module 'matematika' is not allowed in Universal Mode.", "Import 'math' instead.")
			case "teks":
				return newError(n.GetMetadata(), "Indonesian module 'teks' is not allowed in Universal Mode.", "Import 'strings' instead.")
			case "waktu":
				return newError(n.GetMetadata(), "Indonesian module 'waktu' is not allowed in Universal Mode.", "Import 'time' instead.")
			case "sistem":
				return newError(n.GetMetadata(), "Indonesian module 'sistem' is not allowed in Universal Mode.", "Import 'os' instead.")
			case "kripto":
				return newError(n.GetMetadata(), "Indonesian module 'kripto' is not allowed in Universal Mode.", "Import 'crypto' instead.")
			case "arsip", "kompresi":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian module '%s' is not allowed in Universal Mode.", n.Source), "Import 'archive' or 'compression' instead.")
			}
		}
		targetAbsPath, err := resolveModulePath(env.FilePath, n.Source)
		if err != nil {
			msg := err.Error()
			hint := "Pastikan path modul benar dan file tersebut ada."
			if env.IsUniversal {
				hint = "Ensure module path is correct and the file exists."
			}
			return newError(n.GetMetadata(), msg, hint)
		}
		targetEnv, err := loadModule(targetAbsPath, env)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}

		callerFile := env.FilePath
		callerRoot := env.ModuleRoot
		targetFile := targetEnv.FilePath
		targetRoot := targetEnv.ModuleRoot

		vis := n.Visibility
		if vis == nil {
			vis = &uast.Visibility{Kind: uast.VisibilityPrivate}
		}

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
				if spec.Local == "*" {
					for _, name := range mapOrder {
						val := mapPairs[name]
						entry, _ := targetEnv.GetEntry(name)
						hint := ""
						if entry != nil {
							hint = entry.TypeHint
						}
						if _, err := env.DefineWithVisibility(name, val, true, hint, vis); err != nil {
							return newError(n.GetMetadata(), err.Error(), "")
						}
					}
				} else {
					obj := &Map{Pairs: mapPairs, Order: mapOrder}
					if _, err := env.DefineWithVisibility(spec.Local, obj, true, "Map", vis); err != nil {
						return newError(n.GetMetadata(), err.Error(), "")
					}
				}
			} else {
				entry, ok := targetEnv.GetEntry(spec.Imported)
				if !ok {
					return newError(n.GetMetadata(), fmt.Sprintf("Simbol '%s' tidak ditemukan di berkas '%s'", spec.Imported, filepath.Base(targetFile)), "Pastikan nama simbol dieja dengan benar.")
				}
				if err := checkVisibilityAccess(callerFile, callerRoot, targetFile, targetRoot, spec.Imported, entry.Visibility); err != nil {
					return newError(n.GetMetadata(), err.Error(), "Pastikan aturan visibilitas pada file target mengizinkan akses dari file ini.")
				}
				if _, err := env.DefineWithVisibility(spec.Local, entry.Current(), true, entry.TypeHint, vis); err != nil {
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
			if isError(finRes) || (finRes != nil && (finRes.Type() == RETURN_VALUE_OBJ || finRes.Type() == BREAK_SIGNAL_OBJ || finRes.Type() == CONTINUE_SIGNAL_OBJ || finRes.Type() == ROLLBACK_SIGNAL_OBJ)) {
				return finRes
			}
		}

		return result

	case *uast.StyleBlockStatement:
		var result Object = NULL
		if n.Body != nil {
			for _, stmt := range n.Body.Statements {
				result = Eval(stmt, env)
				if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ) {
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

	case *uast.BreakStatement:
		return BREAK_SIGNAL

	case *uast.ContinueStatement:
		return CONTINUE_SIGNAL

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

	case *uast.TransactionStatement:
		env.BeginTransaction()
		bodyRes := Eval(n.Body, env)

		if bodyRes != nil && bodyRes.Type() == ROLLBACK_SIGNAL_OBJ {
			env.RollbackTransaction()
			sig := bodyRes.(*RollbackSignal)
			if n.Recover != nil {
				recEnv := NewEnclosedEnvironment(env)
				if n.RecoverVar != "" {
					recVal := sig.Reason
					if recVal == nil || recVal.Type() == NULL_OBJ {
						recVal = &String{Value: "transaksi dibatalkan"}
					}
					recEnv.Define(n.RecoverVar, recVal, false, "")
				}
				return Eval(n.Recover, recEnv)
			}
			return NULL
		} else if isError(bodyRes) {
			env.RollbackTransaction()
			if n.Recover != nil {
				recEnv := NewEnclosedEnvironment(env)
				if n.RecoverVar != "" {
					recVal := bodyRes
					if errObj, ok := bodyRes.(*Error); ok {
						recVal = &String{Value: errObj.Message}
					}
					recEnv.Define(n.RecoverVar, recVal, false, "")
				}
				return Eval(n.Recover, recEnv)
			}
			return bodyRes
		}

		env.CommitTransaction()
		return bodyRes

	case *uast.RollbackStatement:
		var reason Object = NULL
		if n.Argument != nil {
			reason = Eval(n.Argument, env)
			if isError(reason) {
				return reason
			}
		}
		return &RollbackSignal{Reason: reason}

	case *uast.ReactiveDeclaration:
		entry, err := env.DefineReactive(n.Name, n.Value, env)
		if err != nil {
			return newError(n.GetMetadata(), err.Error(), "")
		}
		return entry.Current()

	case *uast.WatchStatement:
		targetEntry, ok := env.GetEntry(n.Target)
		if !ok || targetEntry == nil {
			msg := fmt.Sprintf("Variabel '%s' tidak ditemukan untuk dipantau", n.Target)
			hint := "Pastikan variabel target sudah dideklarasikan sebelum dipantau."
			if env.IsUniversal {
				msg = fmt.Sprintf("Target variable '%s' not found to watch", n.Target)
				hint = "Ensure target variable is declared before watching."
			}
			return newError(n.GetMetadata(), msg, hint)
		}
		newVar := n.NewValVar
		if newVar == "" {
			if env.IsUniversal {
				newVar = "newVal"
			} else {
				newVar = "baru"
			}
		}
		oldVar := n.OldValVar
		if oldVar == "" {
			if env.IsUniversal {
				oldVar = "oldVal"
			} else {
				oldVar = "lama"
			}
		}
		watcher := &ReactiveWatcher{
			TargetName: n.Target,
			NewValVar:  newVar,
			OldValVar:  oldVar,
			Body:       n.Body,
			Env:        env,
		}
		root := targetEntry.Resolve()
		root.Watchers = append(root.Watchers, watcher)
		return NULL

	case *uast.TestStatement:
		return evalTestStatement(n, env)

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
		if env.IsUniversal {
			switch n.Name {
			case "cetak":
				return newError(n.GetMetadata(), "Indonesian built-in 'cetak' is not allowed in Universal Mode.", "Use 'print' instead.")
			case "masukan":
				return newError(n.GetMetadata(), "Indonesian built-in 'masukan' is not allowed in Universal Mode.", "Use 'input' instead.")
			case "panjang":
				return newError(n.GetMetadata(), "Indonesian built-in 'panjang' is not allowed in Universal Mode.", "Use 'len' or 'length' instead.")
			case "rentang":
				return newError(n.GetMetadata(), "Indonesian built-in 'rentang' is not allowed in Universal Mode.", "Use 'range' instead.")
			case "tipe":
				return newError(n.GetMetadata(), "Indonesian built-in 'tipe' is not allowed in Universal Mode.", "Use 'type' instead.")
			case "saring":
				return newError(n.GetMetadata(), "Indonesian built-in 'saring' is not allowed in Universal Mode.", "Use 'filter' instead.")
			case "petakan":
				return newError(n.GetMetadata(), "Indonesian built-in 'petakan' is not allowed in Universal Mode.", "Use 'map' instead.")
			case "potong_spasi":
				return newError(n.GetMetadata(), "Indonesian built-in 'potong_spasi' is not allowed in Universal Mode.", "Use 'trim' instead.")
			case "kapitalisasi":
				return newError(n.GetMetadata(), "Indonesian built-in 'kapitalisasi' is not allowed in Universal Mode.", "Use 'capitalize' or 'uppercase' instead.")
			case "Angka", "angka":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian type '%s' is not allowed in Universal Mode.", n.Name), "Use 'Number' instead.")
			case "Teks", "teks":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian type '%s' is not allowed in Universal Mode.", n.Name), "Use 'String' instead.")
			case "Logika", "logika":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian type '%s' is not allowed in Universal Mode.", n.Name), "Use 'Boolean' instead.")
			case "terbilang":
				return newError(n.GetMetadata(), "Indonesian function 'terbilang' is not allowed in Universal Mode.", "Use 'in_words' or '.words' instead.")
			case "rupiah", "format_rupiah":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian function '%s' is not allowed in Universal Mode.", n.Name), "Use 'format_usd' or '.usd' instead.")
			case "salin":
				return newError(n.GetMetadata(), "Indonesian function 'salin' is not allowed in Universal Mode.", "Use 'clone' instead.")
			case "Saluran", "saluran":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian type '%s' is not allowed in Universal Mode.", n.Name), "Use 'Channel' or 'channel' instead.")
			case "Bita", "bita":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian type '%s' is not allowed in Universal Mode.", n.Name), "Use 'Bytes' or 'bytes' instead.")
			case "jemput":
				return newError(n.GetMetadata(), "Indonesian function 'jemput' is not allowed in Universal Mode.", "Use 'await' instead.")
			case "jemput_semua":
				return newError(n.GetMetadata(), "Indonesian function 'jemput_semua' is not allowed in Universal Mode.", "Use 'await_all' instead.")
			case "jemput_pertama":
				return newError(n.GetMetadata(), "Indonesian function 'jemput_pertama' is not allowed in Universal Mode.", "Use 'await_race' or 'await_first' instead.")
			case "berkas":
				return newError(n.GetMetadata(), "Indonesian module 'berkas' is not allowed in Universal Mode.", "Use 'fs' instead.")
			case "matematika":
				return newError(n.GetMetadata(), "Indonesian module 'matematika' is not allowed in Universal Mode.", "Use 'math' instead.")
			case "waktu":
				return newError(n.GetMetadata(), "Indonesian module 'waktu' is not allowed in Universal Mode.", "Use 'time' instead.")
			case "sistem":
				return newError(n.GetMetadata(), "Indonesian module 'sistem' is not allowed in Universal Mode.", "Use 'os' instead.")
			case "kripto":
				return newError(n.GetMetadata(), "Indonesian module 'kripto' is not allowed in Universal Mode.", "Use 'crypto' instead.")
			case "arsip", "kompresi":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian module '%s' is not allowed in Universal Mode.", n.Name), "Use 'archive' or 'compression' instead.")
			case "impor_dinamis", "ambil_dinamis":
				return newError(n.GetMetadata(), fmt.Sprintf("Indonesian function '%s' is not allowed in Universal Mode.", n.Name), "Use 'import_dynamic' instead.")
			}
		}
		val, ok := env.Get(n.Name)
		if !ok {
			if env.IsUniversal {
				return newError(n.GetMetadata(), fmt.Sprintf("Variable '%s' is not defined", n.Name), "Check spelling or ensure it is declared before use.")
			}
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
		return evalMemberAccess(obj, n, env)


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
		if n.Operator == "&&" || n.Operator == "dan" || n.Operator == "and" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if !CoerceToBoolean(left) {
				return FALSE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBool(CoerceToBoolean(right))
		}
		if n.Operator == "||" || n.Operator == "atau" || n.Operator == "or" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if CoerceToBoolean(left) {
				return TRUE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBool(CoerceToBoolean(right))
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
		if mem, ok := n.Callee.(*uast.MemberExpression); ok && mem.Safe && (callee == NULL || callee == nil) {
			return NULL
		}
		if fn, ok := callee.(*Function); ok {
			return applyFunctionWithASTArgs(fn, n.Arguments, env, n.GetMetadata(), n.TypeArguments)
		}
		args := make([]Object, len(n.Arguments))
		for i, argNode := range n.Arguments {
			argVal := Eval(argNode, env)
			if isError(argVal) {
				return argVal
			}
			args[i] = argVal
		}
		return applyFunctionWithGenerics(callee, args, n.TypeArguments, n.GetMetadata())

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
			var keyStr string
			if id, ok := prop.Key.(*uast.Identifier); ok {
				keyStr = id.Name
			} else {
				keyVal := Eval(prop.Key, env)
				if isError(keyVal) {
					return keyVal
				}
				if strObj, ok := keyVal.(*String); ok {
					keyStr = strObj.Value
				} else {
					keyStr = keyVal.Inspect()
				}
			}
			valVal := Eval(prop.Value, env)
			if isError(valVal) {
				return valVal
			}
			if _, exists := m.Pairs[keyStr]; !exists {
				m.Order = append(m.Order, keyStr)
			}
			m.Pairs[keyStr] = valVal
		}
		return m

	case *uast.TitipExpression:
		ticket := NewTicket()
		subEnv := NewEnclosedEnvironment(env)
		go func() {
			res := Eval(n.Argument, subEnv)
			if isError(res) {
				ticket.Fail(res)
			} else {
				ticket.Resolve(res)
			}
		}()
		return ticket

	case *uast.JemputExpression:
		target := Eval(n.Argument, env)
		if isError(target) {
			return target
		}
		ticket, ok := target.(*Ticket)
		if !ok {
			return newError(n.GetMetadata(),
				fmt.Sprintf("Target jemput harus berupa Tiket, bukan %s", target.Type()),
				"Pastikan menggunakan 'titip' untuk menghasilkan tiket.")
		}
		var timeoutDuration time.Duration = 0
		if n.Timeout != nil {
			tVal := Eval(n.Timeout, env)
			if isError(tVal) {
				return tVal
			}
			timeoutDuration = parseTimeoutValue(tVal)
		}
		return ticket.Await(timeoutDuration, n.GetMetadata())
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
			if result.Type() == ERROR_OBJ || result.Type() == ASSERTION_ERROR_OBJ {
				return result
			}
			if result.Type() == ROLLBACK_SIGNAL_OBJ {
				return newError(nil, "perintah 'batalkan' hanya dapat digunakan di dalam blok transaksi", "Bungkus operasi dengan 'transaksi { ... }' terlebih dahulu.")
			}
			if result.Type() == BREAK_SIGNAL_OBJ {
				return newError(nil, "perintah 'berhenti' (break) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
			}
			if result.Type() == CONTINUE_SIGNAL_OBJ {
				return newError(nil, "perintah 'lanjut' (continue) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
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
			if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ || result.Type() == ASSERTION_ERROR_OBJ || result.Type() == ROLLBACK_SIGNAL_OBJ || result.Type() == BREAK_SIGNAL_OBJ || result.Type() == CONTINUE_SIGNAL_OBJ {
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

		// Array concatenation / append
		if leftArr, ok := left.(*Array); ok {
			if rightArr, ok := right.(*Array); ok {
				newElems := make([]Object, len(leftArr.Elements)+len(rightArr.Elements))
				copy(newElems, leftArr.Elements)
				copy(newElems[len(leftArr.Elements):], rightArr.Elements)
				return &Array{Elements: newElems}
			}
			newElems := make([]Object, len(leftArr.Elements)+1)
			copy(newElems, leftArr.Elements)
			newElems[len(leftArr.Elements)] = right
			return &Array{Elements: newElems}
		}

		// Bytes concatenation
		if leftB, ok := left.(*Bytes); ok {
			if rightB, ok := right.(*Bytes); ok {
				newElems := make([]byte, len(leftB.Elements)+len(rightB.Elements))
				copy(newElems, leftB.Elements)
				copy(newElems[len(leftB.Elements):], rightB.Elements)
				return &Bytes{Elements: newElems}
			}
			if rightStr, ok := right.(*String); ok {
				newElems := make([]byte, len(leftB.Elements)+len(rightStr.Value))
				copy(newElems, leftB.Elements)
				copy(newElems[len(leftB.Elements):], []byte(rightStr.Value))
				return &Bytes{Elements: newElems}
			}
		}
	}

	if op == "*" {
		if leftB, ok := left.(*Bytes); ok {
			if cnt, _, isInt, ok := CoerceToNumeric(right); ok && isInt {
				if cnt <= 0 {
					return &Bytes{Elements: []byte{}}
				}
				if cnt > 1_000_000 || int64(len(leftB.Elements))*cnt > 50_000_000 {
					return newError(meta, "hasil pengulangan bita melebihi batas aman (50MB)", "")
				}
				return &Bytes{Elements: bytes.Repeat(leftB.Elements, int(cnt))}
			}
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
		case "%":
			if rightFloat == 0.0 {
				return newError(meta, "Operasi modulo dengan nol", "Pastikan penyebut bernilai selain 0.")
			}
			return &Float{Value: math.Mod(leftFloat, rightFloat)}
		}

	case "==", "!=", "===", "!==":
		eq := ObjectsEqual(left, right)
		if op == "==" || op == "===" {
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
	return applyFunctionWithGenerics(callee, args, nil, meta)
}

func applyFunctionWithGenerics(callee Object, args []Object, typeArgs []string, meta *uast.NodeMetadata) Object {
	switch fn := callee.(type) {
	case *Builtin:
		res := fn.Fn(args...)
		if len(typeArgs) > 0 {
			if ch, ok := res.(*Channel); ok {
				ch.ElemType = typeArgs[0]
			}
		}
		return res
	case *Function:
		fnEnv := NewFunctionEnvironment(fn.Env)
		defer func() {
			for _, defStmt := range fnEnv.PopAllDefers() {
				Eval(defStmt, fnEnv)
			}
		}()

		concreteTypes := make(map[string]string)
		if len(fn.TypeParameters) > 0 {
			if len(typeArgs) > 0 && len(typeArgs) != len(fn.TypeParameters) {
				return newError(meta, fmt.Sprintf("Fungsi '%s' mengharapkan %d parameter tipe generik, tetapi diberikan %d", fn.Name, len(fn.TypeParameters), len(typeArgs)), "")
			}
			if len(typeArgs) > 0 {
				for i, tp := range fn.TypeParameters {
					argT := typeArgs[i]
					if tp.Constraint != "" && !isTypeCompatible(argT, tp.Constraint) {
						return newError(meta, fmt.Sprintf("Argumen tipe '%s' tidak memenuhi batasan '%s' pada fungsi '%s'", argT, tp.Constraint, fn.Name), fmt.Sprintf("Gunakan tipe yang memenuhi batasan %s.", tp.Constraint))
					}
					concreteTypes[tp.Name] = argT
				}
			}
		}

		for i, param := range fn.Parameters {
			var argVal Object = NULL
			if i < len(args) {
				argVal = args[i]
			} else if param.DefaultValue != nil {
				argVal = Eval(param.DefaultValue, fnEnv)
			}

			if expectedT, ok := concreteTypes[param.TypeHint]; ok {
				if argVal != NULL && !checkValueMatchesType(argVal, expectedT) {
					return newError(meta, fmt.Sprintf("Argumen '%s' bertipe %s, tetapi parameter generik '%s' mengharuskan tipe %s", param.Name, objectTypeName(argVal), param.TypeHint, expectedT), fmt.Sprintf("Gunakan nilai bertipe %s.", expectedT))
				}
			} else if len(fn.TypeParameters) > 0 && len(typeArgs) == 0 {
				for _, tp := range fn.TypeParameters {
					if param.TypeHint == tp.Name {
						if tp.Constraint != "" && !checkValueMatchesType(argVal, tp.Constraint) {
							return newError(meta, fmt.Sprintf("Argumen '%s' bertipe %s tidak memenuhi batasan '%s' pada fungsi '%s'", param.Name, objectTypeName(argVal), tp.Constraint, fn.Name), fmt.Sprintf("Gunakan nilai bertipe %s.", tp.Constraint))
						}
						concreteTypes[tp.Name] = objectTypeName(argVal)
					}
				}
			}

			fnEnv.Define(param.Name, argVal, true, param.TypeHint)
		}
		res := Eval(fn.Body, fnEnv)
		if ret, ok := res.(*ReturnValue); ok {
			return ret.Value
		}
		if res != nil {
			if res.Type() == BREAK_SIGNAL_OBJ {
				return newError(meta, "perintah 'berhenti' (break) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
			}
			if res.Type() == CONTINUE_SIGNAL_OBJ {
				return newError(meta, "perintah 'lanjut' (continue) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
			}
		}
		return res
	case *StructDefinition:
		instEnv := NewEnclosedEnvironment(fn.Env)
		instance := &StructInstance{
			Def:           fn,
			TypeArguments: typeArgs,
			Fields:        make(map[string]Object),
			Env:           instEnv,
		}
		instEnv.Define("diri", instance, false, fn.Name)
		instEnv.Define("this", instance, false, fn.Name)
		instEnv.Define("self", instance, false, fn.Name)

		concreteTypes := make(map[string]string)
		if len(fn.TypeParameters) > 0 {
			if len(typeArgs) > 0 && len(typeArgs) != len(fn.TypeParameters) {
				return newError(meta, fmt.Sprintf("Bentuk '%s' mengharapkan %d parameter tipe generik, tetapi diberikan %d", fn.Name, len(fn.TypeParameters), len(typeArgs)), "")
			}
			if len(typeArgs) > 0 {
				for i, tp := range fn.TypeParameters {
					argT := typeArgs[i]
					if tp.Constraint != "" && !isTypeCompatible(argT, tp.Constraint) {
						return newError(meta, fmt.Sprintf("Argumen tipe '%s' tidak memenuhi batasan '%s' pada bentuk '%s'", argT, tp.Constraint, fn.Name), fmt.Sprintf("Gunakan tipe yang memenuhi batasan %s.", tp.Constraint))
					}
					concreteTypes[tp.Name] = argT
				}
			}
		}

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

		// Validate generic type constraints on fields
		for _, f := range fn.Fields {
			val := instance.Fields[f.Name]
			if expectedT, ok := concreteTypes[f.TypeHint]; ok {
				if val != NULL && val != nil && !checkValueMatchesType(val, expectedT) {
					return newError(meta, fmt.Sprintf("Nilai untuk bidang '%s' bertipe %s, tetapi parameter generik '%s' mengharuskan tipe %s", f.Name, objectTypeName(val), f.TypeHint, expectedT), fmt.Sprintf("Gunakan nilai bertipe %s.", expectedT))
				}
			} else if len(fn.TypeParameters) > 0 && len(typeArgs) == 0 {
				for _, tp := range fn.TypeParameters {
					if f.TypeHint == tp.Name {
						if tp.Constraint != "" && !checkValueMatchesType(val, tp.Constraint) {
							return newError(meta, fmt.Sprintf("Nilai untuk bidang '%s' bertipe %s tidak memenuhi batasan '%s' pada bentuk '%s'", f.Name, objectTypeName(val), tp.Constraint, fn.Name), fmt.Sprintf("Gunakan nilai bertipe %s.", tp.Constraint))
						}
						concreteTypes[tp.Name] = objectTypeName(val)
					}
				}
			}
		}

		if len(typeArgs) == 0 && len(fn.TypeParameters) > 0 {
			var inferred []string
			for _, tp := range fn.TypeParameters {
				if ct, ok := concreteTypes[tp.Name]; ok {
					inferred = append(inferred, ct)
				}
			}
			if len(inferred) == len(fn.TypeParameters) {
				instance.TypeArguments = inferred
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
	return obj != nil && (obj.Type() == ERROR_OBJ || obj.Type() == ASSERTION_ERROR_OBJ)
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
		case *Bytes:
			return &String{Value: string(v.Elements)}
		default:
			return &String{Value: v.Inspect()}
		}

	case "Logika", "logika", "Kebenaran", "kebenaran", "Bool", "bool", "Boolean", "boolean":
		return nativeBool(CoerceToBoolean(val))

	case "Bita", "bita", "Bytes", "bytes":
		switch v := val.(type) {
		case *Bytes:
			return v
		case *String:
			return &Bytes{Elements: []byte(v.Value)}
		case *Integer:
			if v.Value < 0 {
				return newError(meta, "Ukuran alokasi bita tidak boleh negatif", "Gunakan bilangan bulat positif untuk alokasi buffer bita.")
			}
			if v.Value > 100_000_000 {
				return newError(meta, "Ukuran alokasi bita melebihi batas aman (100MB)", "")
			}
			return &Bytes{Elements: make([]byte, v.Value)}
		case *Array:
			elems := make([]byte, len(v.Elements))
			for i, el := range v.Elements {
				if iVal, _, isInt, ok := CoerceToNumeric(el); ok && isInt {
					if iVal < 0 || iVal > 255 {
						return newError(meta, fmt.Sprintf("Nilai elemen bita harus di antara 0 dan 255, menerima %d", iVal), "")
					}
					elems[i] = byte(iVal)
				} else {
					return newError(meta, fmt.Sprintf("Elemen array untuk bita harus berupa angka (0..255), menerima %s", el.Type()), "")
				}
			}
			return &Bytes{Elements: elems}
		default:
			return newError(meta, fmt.Sprintf("Tidak dapat mengubah tipe %s menjadi Bita", val.Type()), "")
		}

	default:
		return newError(meta, fmt.Sprintf("Tipe tujuan casting '%s' tidak didukung", targetType), "Tipe yang didukung: Angka, Pecahan, Teks, Kebenaran, Bita, atau alias yang terdaftar.")
	}
}

func applyFunctionWithASTArgs(fn *Function, argNodes []uast.Expression, callerEnv *Environment, meta *uast.NodeMetadata, typeArgs []string) Object {
	fnEnv := NewFunctionEnvironment(fn.Env)
	defer func() {
		for _, defStmt := range fnEnv.PopAllDefers() {
			Eval(defStmt, fnEnv)
		}
	}()

	concreteTypes := make(map[string]string)
	if len(fn.TypeParameters) > 0 {
		if len(typeArgs) > 0 && len(typeArgs) != len(fn.TypeParameters) {
			return newError(meta, fmt.Sprintf("Fungsi '%s' mengharapkan %d parameter tipe generik, tetapi diberikan %d", fn.Name, len(fn.TypeParameters), len(typeArgs)), "")
		}
		if len(typeArgs) > 0 {
			for i, tp := range fn.TypeParameters {
				argT := typeArgs[i]
				if tp.Constraint != "" && !isTypeCompatible(argT, tp.Constraint) {
					return newError(meta, fmt.Sprintf("Argumen tipe '%s' tidak memenuhi batasan '%s' pada fungsi '%s'", argT, tp.Constraint, fn.Name), fmt.Sprintf("Gunakan tipe yang memenuhi batasan %s.", tp.Constraint))
				}
				concreteTypes[tp.Name] = argT
			}
		}
	}

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

			if expectedT, ok := concreteTypes[param.TypeHint]; ok {
				if argVal != NULL && !checkValueMatchesType(argVal, expectedT) {
					return newError(meta, fmt.Sprintf("Argumen '%s' bertipe %s, tetapi parameter generik '%s' mengharuskan tipe %s", param.Name, objectTypeName(argVal), param.TypeHint, expectedT), fmt.Sprintf("Gunakan nilai bertipe %s.", expectedT))
				}
			} else if len(fn.TypeParameters) > 0 && len(typeArgs) == 0 {
				for _, tp := range fn.TypeParameters {
					if param.TypeHint == tp.Name {
						if tp.Constraint != "" && !checkValueMatchesType(argVal, tp.Constraint) {
							return newError(meta, fmt.Sprintf("Argumen '%s' bertipe %s tidak memenuhi batasan '%s' pada fungsi '%s'", param.Name, objectTypeName(argVal), tp.Constraint, fn.Name), fmt.Sprintf("Gunakan nilai bertipe %s.", tp.Constraint))
						}
						concreteTypes[tp.Name] = objectTypeName(argVal)
					}
				}
			}

			fnEnv.Define(param.Name, argVal, true, param.TypeHint)
		}
	}
	res := Eval(fn.Body, fnEnv)
	if ret, ok := res.(*ReturnValue); ok {
		return ret.Value
	}
	if res != nil {
		if res.Type() == BREAK_SIGNAL_OBJ {
			return newError(meta, "perintah 'berhenti' (break) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
		}
		if res.Type() == CONTINUE_SIGNAL_OBJ {
			return newError(meta, "perintah 'lanjut' (continue) hanya dapat digunakan di dalam perulangan", "Hapus perintah tersebut atau bungkus dengan blok perulangan.")
		}
	}
	return res
}

func evalMemberAccess(obj Object, n *uast.MemberExpression, env *Environment) Object {
	var computedObj Object
	propName := n.Property
	if n.Computed && n.PropertyExpr != nil {
		computedObj = Eval(n.PropertyExpr, env)
		if isError(computedObj) {
			return computedObj
		}
		if s, ok := computedObj.(*String); ok {
			propName = s.Value
		} else {
			propName = computedObj.Inspect()
		}
	}

	// 1. Array index access, slicing, or built-in methods (arr.push, arr.pop, arr.insert, etc.)
	if arr, ok := obj.(*Array); ok {
		if methodObj, ok := EvalArrayMethod(arr, propName, env, n.GetMetadata()); ok {
			return methodObj
		}
		if rng, isRng := n.PropertyExpr.(*uast.RangeExpression); isRng {
			return evalArraySlice(arr, rng, env, n.GetMetadata())
		}
		var intIdx int64
		var hasInt bool
		if computedObj != nil {
			if iVal, _, isInt, ok := CoerceToNumeric(computedObj); ok && isInt {
				intIdx = iVal
				hasInt = true
			}
		} else if parsed, err := strconv.ParseInt(propName, 10, 64); err == nil {
			intIdx = parsed
			hasInt = true
		} else if n.Computed && propName != "" {
			idxObj := Eval(uast.NewIdentifier(propName), env)
			if isError(idxObj) {
				return idxObj
			}
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
		if computedObj != nil {
			if iVal, _, isInt, ok := CoerceToNumeric(computedObj); ok && isInt {
				intIdx = iVal
				hasInt = true
			}
		} else if parsed, err := strconv.ParseInt(propName, 10, 64); err == nil {
			intIdx = parsed
			hasInt = true
		} else if n.Computed && propName != "" {
			idxObj := Eval(uast.NewIdentifier(propName), env)
			if isError(idxObj) {
				return idxObj
			}
			if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
				intIdx = iVal
				hasInt = true
			}
		}
		if hasInt {
			runes := []rune(str.Value)
			if intIdx >= 0 && int(intIdx) < len(runes) {
				return &String{Value: string(runes[intIdx])}
			}
			return newError(n.GetMetadata(), fmt.Sprintf("Indeks teks %d di luar jangkauan (panjang: %d)", intIdx, len(runes)), "")
		}
	}

	// 1c. Bytes index access, slicing, and methods (e.g. b[0], b[0..2], b.ke_teks())
	if b, ok := obj.(*Bytes); ok {
		if methodObj, ok := EvalBytesMethod(b, propName, env, n.GetMetadata()); ok {
			return methodObj
		}
		if rng, isRng := n.PropertyExpr.(*uast.RangeExpression); isRng {
			return evalBytesSlice(b, rng, env, n.GetMetadata())
		}
		var intIdx int64
		var hasInt bool
		if computedObj != nil {
			if iVal, _, isInt, ok := CoerceToNumeric(computedObj); ok && isInt {
				intIdx = iVal
				hasInt = true
			}
		} else if parsed, err := strconv.ParseInt(propName, 10, 64); err == nil {
			intIdx = parsed
			hasInt = true
		} else if n.Computed && propName != "" {
			idxObj := Eval(uast.NewIdentifier(propName), env)
			if isError(idxObj) {
				return idxObj
			}
			if iVal, _, isInt, ok := CoerceToNumeric(idxObj); ok && isInt {
				intIdx = iVal
				hasInt = true
			}
		}
		if hasInt {
			if intIdx < 0 {
				intIdx = int64(len(b.Elements)) + intIdx
			}
			if intIdx >= 0 && int(intIdx) < len(b.Elements) {
				return &Integer{Value: int64(b.Elements[intIdx])}
			}
			return newError(n.GetMetadata(), fmt.Sprintf("Indeks bita %d di luar jangkauan (panjang: %d)", intIdx, len(b.Elements)), "")
		}
	}

	// 2. Struct instance field or method access
	if inst, ok := obj.(*StructInstance); ok {
		prop := propName
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

	// 3. Map property access and built-in methods (mp.keys, mp.values, mp.get, etc.)
	if mp, ok := obj.(*Map); ok {
		key := propName
		if val, exists := mp.Pairs[key]; exists {
			return val
		}
		if methodObj, ok := EvalMapMethod(mp, key, env, n.GetMetadata()); ok {
			return methodObj
		}
	}

	// 3b. Channel member access
	if ch, ok := obj.(*Channel); ok {
		if env.IsUniversal {
			switch n.Property {
			case "kirim":
				return newError(n.GetMetadata(), "Indonesian method 'kirim' is not allowed in Universal Mode.", "Use '.send(data)' instead.")
			case "terima":
				return newError(n.GetMetadata(), "Indonesian method 'terima' is not allowed in Universal Mode.", "Use '.receive()' or '.recv()' instead.")
			case "tutup":
				return newError(n.GetMetadata(), "Indonesian method 'tutup' is not allowed in Universal Mode.", "Use '.close()' instead.")
			case "panjang":
				return newError(n.GetMetadata(), "Indonesian property 'panjang' is not allowed in Universal Mode.", "Use '.length' or '.len' instead.")
			case "kapasitas":
				return newError(n.GetMetadata(), "Indonesian property 'kapasitas' is not allowed in Universal Mode.", "Use '.capacity' or '.cap' instead.")
			}
		}
		switch n.Property {
		case "kirim", "send":
			return &Builtin{
				Fn: func(args ...Object) Object {
					if len(args) == 0 {
						msg := "Metode kirim() membutuhkan 1 argumen data"
						hint := "Contoh: saluran.kirim(data)"
						if env.IsUniversal {
							msg = "Method send() requires 1 data argument"
							hint = "Example: channel.send(data)"
						}
						return newError(n.GetMetadata(), msg, hint)
					}
					return ch.Send(args[0], n.GetMetadata())
				},
			}
		case "terima", "receive", "recv":
			return &Builtin{
				Fn: func(args ...Object) Object {
					var to time.Duration = 0
					if len(args) > 0 {
						to = parseTimeoutValue(args[0])
					}
					return ch.Receive(to)
				},
			}
		case "tutup", "close":
			return &Builtin{
				Fn: func(args ...Object) Object {
					return ch.Close()
				},
			}
		case "panjang", "length", "len":
			return &Integer{Value: int64(len(ch.Ch))}
		case "kapasitas", "capacity", "cap":
			return &Integer{Value: int64(ch.Capacity)}
		}
	}

	// 3c. Ticket member access
	if t, ok := obj.(*Ticket); ok {
		if env.IsUniversal {
			switch n.Property {
			case "selesai":
				return newError(n.GetMetadata(), "Indonesian property 'selesai' is not allowed in Universal Mode.", "Use '.done' instead.")
			case "hasil":
				return newError(n.GetMetadata(), "Indonesian property 'hasil' is not allowed in Universal Mode.", "Use '.result' instead.")
			case "jemput":
				return newError(n.GetMetadata(), "Indonesian method 'jemput' is not allowed in Universal Mode.", "Use '.await()' instead.")
			}
		}
		switch n.Property {
		case "selesai", "done":
			return nativeBool(t.IsDone())
		case "hasil", "result", "jemput", "await":
			return &Builtin{
				Fn: func(args ...Object) Object {
					var to time.Duration = 0
					if len(args) > 0 {
						to = parseTimeoutValue(args[0])
					}
					return t.Await(to, n.GetMetadata())
				},
			}
		}
	}

	// 3d. Numeric cultural, nominal, terbilang, data size & duration helpers (e.g. 50.ribu, 3.detik, 10.mb, 15000.rupiah, 50000.terbilang, 50.k, 25.usd)
	numProp := propName
	if res, ok := EvalNumericMemberWithEnv(obj, numProp, env); ok {
		return res
	}

	// 3e. Builtin static properties / methods access (e.g. bita.dari_base64, bytes.from_hex)
	if b, ok := obj.(*Builtin); ok {
		if env.IsUniversal {
			switch propName {
			case "dari_base64":
				return newError(n.GetMetadata(), "Indonesian method 'dari_base64' is not allowed in Universal Mode.", "Use '.from_base64(...)' instead.")
			case "dari_heks":
				return newError(n.GetMetadata(), "Indonesian method 'dari_heks' is not allowed in Universal Mode.", "Use '.from_hex(...)' instead.")
			}
		}
		if prop, exists := b.GetProperty(propName); exists {
			return prop
		}
	}

	// 4. Time-Travel Variable access fallback (e.g. x.sebelumnya, x.awal, x.histori, x.previous, x.initial, x.history)
	if id, ok := n.Object.(*uast.Identifier); ok {
		entry, exists := env.GetEntry(id.Name)
		if exists {
			if env.IsUniversal {
				switch n.Property {
				case "sebelumnya":
					return newError(n.GetMetadata(), "Indonesian property '.sebelumnya' is not allowed in Universal Mode.", "Use '.previous' or '.prev' instead.")
				case "awal":
					return newError(n.GetMetadata(), "Indonesian property '.awal' is not allowed in Universal Mode.", "Use '.initial' instead.")
				case "histori":
					return newError(n.GetMetadata(), "Indonesian property '.histori' is not allowed in Universal Mode.", "Use '.history' instead.")
				}
			}
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
		msg := fmt.Sprintf("Properti atau fungsi '%s' tidak ditemukan pada bentuk '%s'", n.Property, inst.Def.Name)
		hint := "Pastikan nama field atau fungsi dieja dengan benar."
		if env.IsUniversal {
			msg = fmt.Sprintf("Property or method '%s' not found on struct '%s'", n.Property, inst.Def.Name)
			hint = "Ensure the field or method name is spelled correctly."
		}
		return newError(n.GetMetadata(), msg, hint)
	}

	msg := fmt.Sprintf("Properti '%s' tidak ditemukan pada objek", n.Property)
	hint := "Pastikan properti yang diakses valid."
	if env.IsUniversal {
		msg = fmt.Sprintf("Property '%s' not found on object", n.Property)
		hint = "Ensure the accessed property exists and is valid."
	}
	return newError(n.GetMetadata(), msg, hint)
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
		res := evalMemberAccess(obj, e, env)
		if isError(res) {
			return NULL
		}
		return res

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

	if (endVal >= startVal && (endVal-startVal > 1_000_000 || endVal-startVal < 0)) ||
		(startVal > endVal && (startVal-endVal > 1_000_000 || startVal-endVal < 0)) {
		msg := "Ukuran rentang terlalu besar (maksimal 1.000.000 elemen)"
		if env.IsUniversal {
			msg = "Range size is too large (maximum 1,000,000 elements)"
		}
		return newError(rng.GetMetadata(), msg, "Gunakan perulangan dengan kondisi 'selama' (while) untuk iterasi dalam skala besar.")
	}

	var elements []Object
	if startVal <= endVal {
		count := endVal - startVal + 1
		for idx := int64(0); idx < count; idx++ {
			elements = append(elements, &Integer{Value: startVal + idx})
		}
	} else {
		count := startVal - endVal + 1
		for idx := int64(0); idx < count; idx++ {
			elements = append(elements, &Integer{Value: startVal - idx})
		}
	}
	return &Array{Elements: elements}
}

func evalArraySlice(arr *Array, rng *uast.RangeExpression, env *Environment, meta *uast.NodeMetadata) Object {
	n := len(arr.Elements)
	start := 0
	end := n
	if rng.Start != nil {
		sVal := Eval(rng.Start, env)
		if isError(sVal) {
			return sVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(sVal); ok && isInt {
			start = int(iVal)
		} else {
			msg := "Batas awal rentang slicing harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "Start index for slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if rng.End != nil {
		eVal := Eval(rng.End, env)
		if isError(eVal) {
			return eVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(eVal); ok && isInt {
			end = int(iVal)
		} else {
			msg := "Batas akhir rentang slicing harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "End index for slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if start < 0 {
		start = n + start
		if start < 0 {
			start = 0
		}
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = n + end
		if end < 0 {
			end = 0
		}
	}
	if end < start {
		end = start
	}
	if end > n {
		end = n
	}

	sliced := make([]Object, end-start)
	copy(sliced, arr.Elements[start:end])
	return &Array{Elements: sliced}
}

func evalStringSlice(str *String, rng *uast.RangeExpression, env *Environment, meta *uast.NodeMetadata) Object {
	runes := []rune(str.Value)
	n := len(runes)
	start := 0
	end := n
	if rng.Start != nil {
		sVal := Eval(rng.Start, env)
		if isError(sVal) {
			return sVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(sVal); ok && isInt {
			start = int(iVal)
		} else {
			msg := "Batas awal rentang teks harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "Start index for string slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if rng.End != nil {
		eVal := Eval(rng.End, env)
		if isError(eVal) {
			return eVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(eVal); ok && isInt {
			end = int(iVal)
		} else {
			msg := "Batas akhir rentang teks harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "End index for string slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if start < 0 {
		start = n + start
		if start < 0 {
			start = 0
		}
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = n + end
		if end < 0 {
			end = 0
		}
	}
	if end < start {
		end = start
	}
	if end > n {
		end = n
	}

	return &String{Value: string(runes[start:end])}
}

func evalBytesSlice(b *Bytes, rng *uast.RangeExpression, env *Environment, meta *uast.NodeMetadata) Object {
	n := len(b.Elements)
	start := 0
	end := n
	if rng.Start != nil {
		sVal := Eval(rng.Start, env)
		if isError(sVal) {
			return sVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(sVal); ok && isInt {
			start = int(iVal)
		} else {
			msg := "Batas awal rentang slicing harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "Start index for slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if rng.End != nil {
		eVal := Eval(rng.End, env)
		if isError(eVal) {
			return eVal
		}
		if iVal, _, isInt, ok := CoerceToNumeric(eVal); ok && isInt {
			end = int(iVal)
		} else {
			msg := "Batas akhir rentang slicing harus berupa bilangan bulat"
			if env.IsUniversal {
				msg = "End index for slicing must be an integer"
			}
			return newError(rng.GetMetadata(), msg, "")
		}
	}
	if start < 0 {
		start = n + start
		if start < 0 {
			start = 0
		}
	}
	if start > n {
		start = n
	}
	if end < 0 {
		end = n + end
		if end < 0 {
			end = 0
		}
	}
	if end < start {
		end = start
	}
	if end > n {
		end = n
	}

	sliced := make([]byte, end-start)
	copy(sliced, b.Elements[start:end])
	return &Bytes{Elements: sliced}
}

func findRootIdentifier(expr uast.Expression) string {
	curr := expr
	for curr != nil {
		switch e := curr.(type) {
		case *uast.Identifier:
			return e.Name
		case *uast.MemberExpression:
			curr = e.Object
		default:
			return ""
		}
	}
	return ""
}

func evalTestStatement(stmt *uast.TestStatement, env *Environment) Object {
	tc := env.GetTestContext()
	if tc != nil && tc.Filter != "" {
		if !strings.Contains(strings.ToLower(stmt.Name), strings.ToLower(tc.Filter)) {
			return NULL
		}
	}

	start := time.Now()
	testEnv := NewEnclosedEnvironment(env)
	res := Eval(stmt.Body, testEnv)
	dur := time.Since(start)

	var pos *uast.Position
	if stmt.Metadata != nil {
		pos = stmt.Metadata.Pos
	}

	if isError(res) {
		errMsg := res.Inspect()
		expected := ""
		actual := ""
		if aErr, ok := res.(*AssertionError); ok {
			errMsg = aErr.Message
			expected = aErr.Expected
			actual = aErr.Actual
			if aErr.Pos != nil {
				pos = aErr.Pos
			}
		} else if errObj, ok := res.(*Error); ok {
			errMsg = errObj.Message
			if errObj.Pos != nil {
				pos = errObj.Pos
			}
		}

		if tc != nil {
			tc.RecordResult(TestCaseResult{
				Name:     stmt.Name,
				Passed:   false,
				Duration: dur,
				Error:    errMsg,
				Expected: expected,
				Actual:   actual,
				Pos:      pos,
			})
			return NULL
		}
		return res
	}

	if tc != nil {
		tc.RecordResult(TestCaseResult{
			Name:     stmt.Name,
			Passed:   true,
			Duration: dur,
			Pos:      pos,
		})
	}
	return NULL
}



