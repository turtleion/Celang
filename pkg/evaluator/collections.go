package evaluator

import (
	"fmt"
	"sort"
	"strings"

	"celang/pkg/uast"
)

// objectsEqual checks if two Celang objects are considered equal.
func objectsEqual(a, b Object) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Type() != b.Type() {
		aInt, aFloat, aIsInt, aOk := CoerceToNumeric(a)
		bInt, bFloat, bIsInt, bOk := CoerceToNumeric(b)
		if aOk && bOk {
			if aIsInt && bIsInt {
				return aInt == bInt
			}
			return aFloat == bFloat
		}
		return false
	}
	switch oA := a.(type) {
	case *Integer:
		return oA.Value == b.(*Integer).Value
	case *Float:
		return oA.Value == b.(*Float).Value
	case *String:
		return oA.Value == b.(*String).Value
	case *Boolean:
		return oA.Value == b.(*Boolean).Value
	case *Null:
		return true
	default:
		return a.Inspect() == b.Inspect()
	}
}

// EvalArrayMethod dispatches member method invocations on an Array.
func EvalArrayMethod(arr *Array, method string, env *Environment, meta *uast.NodeMetadata) (Object, bool) {
	isUniversal := env != nil && env.IsUniversal

	// Universal mode guards
	if isUniversal {
		switch method {
		case "tambah":
			return newError(meta, "Indonesian method 'tambah' is not allowed in Universal Mode.", "Use '.push(...)' instead."), true
		case "ambil_terakhir", "hapus_terakhir":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.pop()' instead."), true
		case "geser", "hapus_awal":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.shift()' instead."), true
		case "sisip", "sisipkan":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.insert(index, item)' instead."), true
		case "sisip_depan", "sisipkan_depan", "sisipkan_awal":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.unshift(...)' instead."), true
		case "hapus_indeks":
			return newError(meta, "Indonesian method 'hapus_indeks' is not allowed in Universal Mode.", "Use '.remove_at(index)' instead."), true
		case "balik":
			return newError(meta, "Indonesian method 'balik' is not allowed in Universal Mode.", "Use '.reverse()' instead."), true
		case "urut":
			return newError(meta, "Indonesian method 'urut' is not allowed in Universal Mode.", "Use '.sort()' instead."), true
		case "gabung":
			return newError(meta, "Indonesian method 'gabung' is not allowed in Universal Mode.", "Use '.join(separator)' instead."), true
		case "berisi":
			return newError(meta, "Indonesian method 'berisi' is not allowed in Universal Mode.", "Use '.contains(item)' or '.includes(item)' instead."), true
		case "indeks_dari":
			return newError(meta, "Indonesian method 'indeks_dari' is not allowed in Universal Mode.", "Use '.index_of(item)' instead."), true
		case "bersihkan":
			return newError(meta, "Indonesian method 'bersihkan' is not allowed in Universal Mode.", "Use '.clear()' instead."), true
		case "potong":
			return newError(meta, "Indonesian method 'potong' is not allowed in Universal Mode.", "Use '.slice(start, [end])' instead."), true
		case "gabungkan":
			return newError(meta, "Indonesian method 'gabungkan' is not allowed in Universal Mode.", "Use '.concat(...)' instead."), true
		case "panjang":
			return newError(meta, "Indonesian property 'panjang' is not allowed in Universal Mode.", "Use '.length' or '.len' instead."), true
		case "ukuran":
			return newError(meta, "Indonesian method 'ukuran' is not allowed in Universal Mode.", "Use '.size()' instead."), true
		}
	}

	switch method {
	case "tambah", "push":
		return &Builtin{
			Fn: func(args ...Object) Object {
				arr.Elements = append(arr.Elements, args...)
				return arr
			},
		}, true

	case "hapus_terakhir", "ambil_terakhir", "pop":
		return &Builtin{
			Fn: func(args ...Object) Object {
				n := len(arr.Elements)
				if n == 0 {
					return NULL
				}
				last := arr.Elements[n-1]
				arr.Elements = arr.Elements[:n-1]
				return last
			},
		}, true

	case "hapus_awal", "geser", "shift":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(arr.Elements) == 0 {
					return NULL
				}
				first := arr.Elements[0]
				arr.Elements = arr.Elements[1:]
				return first
			},
		}, true

	case "sisipkan_depan", "sisipkan_awal", "sisip_depan", "unshift":
		return &Builtin{
			Fn: func(args ...Object) Object {
				newElems := make([]Object, len(args)+len(arr.Elements))
				copy(newElems, args)
				copy(newElems[len(args):], arr.Elements)
				arr.Elements = newElems
				return arr
			},
		}, true

	case "sisipkan", "sisip", "insert":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) < 2 {
					return newError(meta, "Metode insert/sisipkan membutuhkan 2 argumen: (indeks, elemen)", "Contoh: arr.sisipkan(0, 'pertama')")
				}
				idxVal, _, isInt, ok := CoerceToNumeric(args[0])
				if !ok || !isInt {
					return newError(meta, "Indeks sisipkan/insert harus berupa bilangan bulat", "")
				}
				idx := int(idxVal)
				if idx < 0 {
					idx = len(arr.Elements) + idx
				}
				if idx < 0 || idx > len(arr.Elements) {
					return newError(meta, fmt.Sprintf("Indeks %d di luar jangkauan array (panjang: %d)", idxVal, len(arr.Elements)), "")
				}
				item := args[1]
				arr.Elements = append(arr.Elements[:idx], append([]Object{item}, arr.Elements[idx:]...)...)
				return arr
			},
		}, true

	case "hapus_indeks", "remove_at", "delete_at":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) < 1 {
					return newError(meta, "Metode remove_at/hapus_indeks membutuhkan 1 argumen: (indeks)", "")
				}
				idxVal, _, isInt, ok := CoerceToNumeric(args[0])
				if !ok || !isInt {
					return newError(meta, "Indeks harus berupa bilangan bulat", "")
				}
				idx := int(idxVal)
				if idx < 0 {
					idx = len(arr.Elements) + idx
				}
				if idx < 0 || idx >= len(arr.Elements) {
					return newError(meta, fmt.Sprintf("Indeks %d di luar jangkauan array (panjang: %d)", idxVal, len(arr.Elements)), "")
				}
				removed := arr.Elements[idx]
				arr.Elements = append(arr.Elements[:idx], arr.Elements[idx+1:]...)
				return removed
			},
		}, true

	case "balik", "reverse":
		return &Builtin{
			Fn: func(args ...Object) Object {
				for i, j := 0, len(arr.Elements)-1; i < j; i, j = i+1, j-1 {
					arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
				}
				return arr
			},
		}, true

	case "urut", "sort":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) > 0 && (args[0].Type() == FUNCTION_OBJ || args[0].Type() == BUILTIN_OBJ) {
					compareFn := args[0]
					var sortErr Object
					sort.SliceStable(arr.Elements, func(i, j int) bool {
						if sortErr != nil {
							return false
						}
						res := applyFunction(compareFn, []Object{arr.Elements[i], arr.Elements[j]}, meta)
						if isError(res) {
							sortErr = res
							return false
						}
						if b, ok := res.(*Boolean); ok {
							return b.Value
						}
						if n, _, isInt, ok := CoerceToNumeric(res); ok && isInt {
							return n < 0
						}
						return CoerceToBoolean(res)
					})
					if sortErr != nil {
						return sortErr
					}
					return arr
				}

				// Natural order sorting (numbers ascending, then strings ascending)
				sort.SliceStable(arr.Elements, func(i, j int) bool {
					elA := arr.Elements[i]
					elB := arr.Elements[j]
					numA, _, _, okA := CoerceToNumeric(elA)
					numB, _, _, okB := CoerceToNumeric(elB)
					if okA && okB {
						return numA < numB
					}
					if sA, ok := elA.(*String); ok {
						if sB, ok := elB.(*String); ok {
							return sA.Value < sB.Value
						}
					}
					return elA.Inspect() < elB.Inspect()
				})
				return arr
			},
		}, true

	case "gabung", "join":
		return &Builtin{
			Fn: func(args ...Object) Object {
				sep := ", "
				if len(args) > 0 {
					if s, ok := args[0].(*String); ok {
						sep = s.Value
					} else {
						sep = args[0].Inspect()
					}
				}
				parts := make([]string, len(arr.Elements))
				for i, el := range arr.Elements {
					if s, ok := el.(*String); ok {
						parts[i] = s.Value
					} else {
						parts[i] = el.Inspect()
					}
				}
				return &String{Value: strings.Join(parts, sep)}
			},
		}, true

	case "berisi", "contains", "includes":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return FALSE
				}
				target := args[0]
				for _, el := range arr.Elements {
					if objectsEqual(el, target) {
						return TRUE
					}
				}
				return FALSE
			},
		}, true

	case "indeks_dari", "index_of":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return &Integer{Value: -1}
				}
				target := args[0]
				for i, el := range arr.Elements {
					if objectsEqual(el, target) {
						return &Integer{Value: int64(i)}
					}
				}
				return &Integer{Value: -1}
			},
		}, true

	case "bersihkan", "clear":
		return &Builtin{
			Fn: func(args ...Object) Object {
				arr.Elements = make([]Object, 0)
				return arr
			},
		}, true

	case "potong", "slice":
		return &Builtin{
			Fn: func(args ...Object) Object {
				n := len(arr.Elements)
				start := 0
				end := n
				if len(args) > 0 {
					sVal, _, isInt, ok := CoerceToNumeric(args[0])
					if ok && isInt {
						start = int(sVal)
					}
				}
				if len(args) > 1 {
					eVal, _, isInt, ok := CoerceToNumeric(args[1])
					if ok && isInt {
						end = int(eVal)
					}
				}
				if start < 0 {
					start = n + start
					if start < 0 {
						start = 0
					}
				}
				if end < 0 {
					end = n + end
					if end < 0 {
						end = 0
					}
				}
				if start > n {
					start = n
				}
				if end > n {
					end = n
				}
				if start >= end {
					return &Array{Elements: []Object{}}
				}
				res := make([]Object, end-start)
				copy(res, arr.Elements[start:end])
				return &Array{Elements: res}
			},
		}, true

	case "gabungkan", "concat":
		return &Builtin{
			Fn: func(args ...Object) Object {
				var res []Object
				res = append(res, arr.Elements...)
				for _, arg := range args {
					if otherArr, ok := arg.(*Array); ok {
						res = append(res, otherArr.Elements...)
					} else {
						res = append(res, arg)
					}
				}
				return &Array{Elements: res}
			},
		}, true

	case "panjang", "length", "len":
		return &Integer{Value: int64(len(arr.Elements))}, true

	case "ukuran", "size":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &Integer{Value: int64(len(arr.Elements))}
			},
		}, true
	}

	return nil, false
}

// EvalMapMethod dispatches member method invocations on a Map.
func EvalMapMethod(mp *Map, method string, env *Environment, meta *uast.NodeMetadata) (Object, bool) {
	isUniversal := env != nil && env.IsUniversal

	// Universal mode guards
	if isUniversal {
		switch method {
		case "kunci":
			return newError(meta, "Indonesian method 'kunci' is not allowed in Universal Mode.", "Use '.keys()' instead."), true
		case "nilai":
			return newError(meta, "Indonesian method 'nilai' is not allowed in Universal Mode.", "Use '.values()' instead."), true
		case "entri":
			return newError(meta, "Indonesian method 'entri' is not allowed in Universal Mode.", "Use '.entries()' instead."), true
		case "memiliki":
			return newError(meta, "Indonesian method 'memiliki' is not allowed in Universal Mode.", "Use '.has(key)' or '.contains(key)' instead."), true
		case "hapus":
			return newError(meta, "Indonesian method 'hapus' is not allowed in Universal Mode.", "Use '.delete(key)' or '.remove(key)' instead."), true
		case "gabung":
			return newError(meta, "Indonesian method 'gabung' is not allowed in Universal Mode.", "Use '.merge(otherMap)' instead."), true
		case "ambil":
			return newError(meta, "Indonesian method 'ambil' is not allowed in Universal Mode.", "Use '.get(key, [default])' instead."), true
		case "atur":
			return newError(meta, "Indonesian method 'atur' is not allowed in Universal Mode.", "Use '.set(key, value)' instead."), true
		case "bersihkan":
			return newError(meta, "Indonesian method 'bersihkan' is not allowed in Universal Mode.", "Use '.clear()' instead."), true
		case "panjang":
			return newError(meta, "Indonesian property 'panjang' is not allowed in Universal Mode.", "Use '.length' or '.len' instead."), true
		case "ukuran":
			return newError(meta, "Indonesian method 'ukuran' is not allowed in Universal Mode.", "Use '.size()' instead."), true
		}
	}

	switch method {
	case "kunci", "keys":
		return &Builtin{
			Fn: func(args ...Object) Object {
				keys := make([]Object, len(mp.Order))
				for i, k := range mp.Order {
					keys[i] = &String{Value: k}
				}
				return &Array{Elements: keys}
			},
		}, true

	case "nilai", "values":
		return &Builtin{
			Fn: func(args ...Object) Object {
				vals := make([]Object, len(mp.Order))
				for i, k := range mp.Order {
					vals[i] = mp.Pairs[k]
				}
				return &Array{Elements: vals}
			},
		}, true

	case "entri", "entries":
		return &Builtin{
			Fn: func(args ...Object) Object {
				entries := make([]Object, len(mp.Order))
				for i, k := range mp.Order {
					entries[i] = &Array{
						Elements: []Object{
							&String{Value: k},
							mp.Pairs[k],
						},
					}
				}
				return &Array{Elements: entries}
			},
		}, true

	case "memiliki", "has", "contains":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return FALSE
				}
				key := args[0].Inspect()
				if s, ok := args[0].(*String); ok {
					key = s.Value
				}
				_, exists := mp.Pairs[key]
				return nativeBool(exists)
			},
		}, true

	case "ambil", "get":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return NULL
				}
				key := args[0].Inspect()
				if s, ok := args[0].(*String); ok {
					key = s.Value
				}
				if val, exists := mp.Pairs[key]; exists {
					return val
				}
				if len(args) > 1 {
					return args[1] // default value
				}
				return NULL
			},
		}, true

	case "atur", "set":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) < 2 {
					return newError(meta, "Metode set/atur membutuhkan 2 argumen: (kunci, nilai)", "")
				}
				key := args[0].Inspect()
				if s, ok := args[0].(*String); ok {
					key = s.Value
				}
				val := args[1]
				if _, exists := mp.Pairs[key]; !exists {
					mp.Order = append(mp.Order, key)
				}
				mp.Pairs[key] = val
				return mp
			},
		}, true

	case "hapus", "delete", "remove":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return FALSE
				}
				key := args[0].Inspect()
				if s, ok := args[0].(*String); ok {
					key = s.Value
				}
				if _, exists := mp.Pairs[key]; !exists {
					return FALSE
				}
				delete(mp.Pairs, key)
				for i, k := range mp.Order {
					if k == key {
						mp.Order = append(mp.Order[:i], mp.Order[i+1:]...)
						break
					}
				}
				return TRUE
			},
		}, true

	case "gabung", "merge":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return mp
				}
				other, ok := args[0].(*Map)
				if !ok {
					return newError(meta, fmt.Sprintf("Metode merge/gabung mengharapkan Map, menerima %s", args[0].Type()), "")
				}
				for _, k := range other.Order {
					if _, exists := mp.Pairs[k]; !exists {
						mp.Order = append(mp.Order, k)
					}
					mp.Pairs[k] = other.Pairs[k]
				}
				return mp
			},
		}, true

	case "bersihkan", "clear":
		return &Builtin{
			Fn: func(args ...Object) Object {
				mp.Pairs = make(map[string]Object)
				mp.Order = make([]string, 0)
				return mp
			},
		}, true

	case "panjang", "length", "len":
		return &Integer{Value: int64(len(mp.Order))}, true

	case "ukuran", "size":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &Integer{Value: int64(len(mp.Order))}
			},
		}, true
	}

	return nil, false
}
