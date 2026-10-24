package evaluator

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"celang/pkg/uast"
)

// ObjectsEqual checks if two Celang objects are considered equal.
func ObjectsEqual(a, b Object) bool {
	visited := make(map[[2]any]bool)
	return objectsEqualWithVisited(a, b, visited, 0)
}

func objectsEqualWithVisited(a, b Object, visited map[[2]any]bool, depth int) bool {
	if depth > 100 {
		return true
	}
	if a == nil || b == nil {
		return a == b
	}
	if a.Type() == NULL_OBJ && b.Type() == NULL_OBJ {
		return true
	}
	if a.Type() == NULL_OBJ || b.Type() == NULL_OBJ {
		return false
	}
	if a.Type() != b.Type() {
		// Only numeric types (INTEGER and FLOAT) can be equal across different types
		if (a.Type() == INTEGER_OBJ || a.Type() == FLOAT_OBJ) && (b.Type() == INTEGER_OBJ || b.Type() == FLOAT_OBJ) {
			_, aFloat, _, aOk := CoerceToNumeric(a)
			_, bFloat, _, bOk := CoerceToNumeric(b)
			if aOk && bOk {
				return aFloat == bFloat
			}
		}
		return false
	}

	pair := [2]any{a, b}
	if visited[pair] {
		return true
	}
	visited[pair] = true

	switch oA := a.(type) {
	case *Integer:
		return oA.Value == b.(*Integer).Value
	case *Float:
		return oA.Value == b.(*Float).Value
	case *String:
		return oA.Value == b.(*String).Value
	case *Boolean:
		return oA.Value == b.(*Boolean).Value
	case *Bytes:
		return bytes.Equal(oA.Elements, b.(*Bytes).Elements)
	case *Null:
		return true
	case *Array:
		oB := b.(*Array)
		if len(oA.Elements) != len(oB.Elements) {
			return false
		}
		for i := range oA.Elements {
			if !objectsEqualWithVisited(oA.Elements[i], oB.Elements[i], visited, depth+1) {
				return false
			}
		}
		return true
	case *Map:
		oB := b.(*Map)
		if len(oA.Pairs) != len(oB.Pairs) {
			return false
		}
		for k, valA := range oA.Pairs {
			valB, exists := oB.Pairs[k]
			if !exists || !objectsEqualWithVisited(valA, valB, visited, depth+1) {
				return false
			}
		}
		return true
	case *StructInstance:
		oB := b.(*StructInstance)
		if oA.Def.Name != oB.Def.Name || len(oA.Fields) != len(oB.Fields) {
			return false
		}
		for k, valA := range oA.Fields {
			valB, exists := oB.Fields[k]
			if !exists || !objectsEqualWithVisited(valA, valB, visited, depth+1) {
				return false
			}
		}
		return true
	default:
		return a.Inspect() == b.Inspect()
	}
}

func objectsEqual(a, b Object) bool {
	return ObjectsEqual(a, b)
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
		case "balik", "balikkan":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.reverse()' instead."), true
		case "urut", "urutkan":
			return newError(meta, fmt.Sprintf("Indonesian method '%s' is not allowed in Universal Mode.", method), "Use '.sort()' instead."), true
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

	case "balik", "balikkan", "reverse":
		return &Builtin{
			Fn: func(args ...Object) Object {
				for i, j := 0, len(arr.Elements)-1; i < j; i, j = i+1, j-1 {
					arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
				}
				return arr
			},
		}, true

	case "urut", "urutkan", "sort":
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
					iA, fA, isIntA, okA := CoerceToNumeric(elA)
					iB, fB, isIntB, okB := CoerceToNumeric(elB)
					if okA && okB {
						if isIntA && isIntB {
							return iA < iB
						}
						return fA < fB
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
					if !ok || !isInt {
						msg := "Indeks awal potong/slice harus berupa bilangan bulat"
						if isUniversal {
							msg = "Start index for slice must be an integer"
						}
						return newError(meta, msg, "")
					}
					start = int(sVal)
				}
				if len(args) > 1 {
					eVal, _, isInt, ok := CoerceToNumeric(args[1])
					if !ok || !isInt {
						msg := "Indeks akhir potong/slice harus berupa bilangan bulat"
						if isUniversal {
							msg = "End index for slice must be an integer"
						}
						return newError(meta, msg, "")
					}
					end = int(eVal)
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
				kList := mp.GetKeys()
				keys := make([]Object, len(kList))
				for i, k := range kList {
					keys[i] = &String{Value: k}
				}
				return &Array{Elements: keys}
			},
		}, true

	case "nilai", "values":
		return &Builtin{
			Fn: func(args ...Object) Object {
				kList := mp.GetKeys()
				vals := make([]Object, len(kList))
				for i, k := range kList {
					vals[i] = mp.Pairs[k]
				}
				return &Array{Elements: vals}
			},
		}, true

	case "entri", "entries":
		return &Builtin{
			Fn: func(args ...Object) Object {
				kList := mp.GetKeys()
				entries := make([]Object, len(kList))
				for i, k := range kList {
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
				for _, k := range other.GetKeys() {
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
		return &Integer{Value: int64(len(mp.GetKeys()))}, true

	case "ukuran", "size":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &Integer{Value: int64(len(mp.GetKeys()))}
			},
		}, true
	}

	return nil, false
}

// EvalBytesMethod dispatches member method invocations and properties on a Bytes object.
func EvalBytesMethod(b *Bytes, prop string, env *Environment, meta *uast.NodeMetadata) (Object, bool) {
	isUniversal := env != nil && env.IsUniversal

	if isUniversal {
		switch prop {
		case "ke_teks":
			return newError(meta, "Indonesian method 'ke_teks' is not allowed in Universal Mode.", "Use '.to_string()' or '.to_text()' instead."), true
		case "ke_heks":
			return newError(meta, "Indonesian method 'ke_heks' is not allowed in Universal Mode.", "Use '.to_hex()' instead."), true
		case "ke_daftar":
			return newError(meta, "Indonesian method 'ke_daftar' is not allowed in Universal Mode.", "Use '.to_array()' instead."), true
		case "ke_base64":
			return newError(meta, "Indonesian method 'ke_base64' is not allowed in Universal Mode.", "Use '.to_base64()' instead."), true
		case "tambah":
			return newError(meta, "Indonesian method 'tambah' is not allowed in Universal Mode.", "Use '.append()' instead."), true
		case "potong":
			return newError(meta, "Indonesian method 'potong' is not allowed in Universal Mode.", "Use '.slice()' instead."), true
		case "salin":
			return newError(meta, "Indonesian method 'salin' is not allowed in Universal Mode.", "Use '.clone()' instead."), true
		case "heks":
			return newError(meta, "Indonesian method 'heks' is not allowed in Universal Mode.", "Use '.hex()' or '.to_hex()' instead."), true
		case "panjang":
			return newError(meta, "Indonesian property 'panjang' is not allowed in Universal Mode.", "Use '.length' or '.len' instead."), true
		}
	}

	switch prop {
	case "ke_teks", "to_string", "to_text":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &String{Value: string(b.Elements)}
			},
		}, true

	case "ke_heks", "to_hex", "heks", "hex":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &String{Value: hex.EncodeToString(b.Elements)}
			},
		}, true

	case "ke_base64", "to_base64":
		return &Builtin{
			Fn: func(args ...Object) Object {
				return &String{Value: base64.StdEncoding.EncodeToString(b.Elements)}
			},
		}, true

	case "ke_daftar", "to_array":
		return &Builtin{
			Fn: func(args ...Object) Object {
				elems := make([]Object, len(b.Elements))
				for i, by := range b.Elements {
					elems[i] = &Integer{Value: int64(by)}
				}
				return &Array{Elements: elems}
			},
		}, true

	case "tambah", "append":
		return &Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					if iVal, _, isInt, ok := CoerceToNumeric(arg); ok && isInt {
						b.Elements = append(b.Elements, byte(iVal))
					} else if otherB, ok := arg.(*Bytes); ok {
						b.Elements = append(b.Elements, otherB.Elements...)
					} else if otherArr, ok := arg.(*Array); ok {
						for _, el := range otherArr.Elements {
							if elInt, _, isI, ok := CoerceToNumeric(el); ok && isI {
								b.Elements = append(b.Elements, byte(elInt))
							}
						}
					} else if s, ok := arg.(*String); ok {
						b.Elements = append(b.Elements, []byte(s.Value)...)
					}
				}
				return b
			},
		}, true

	case "potong", "slice", "sub":
		return &Builtin{
			Fn: func(args ...Object) Object {
				if len(args) == 0 {
					return b
				}
				start := 0
				end := len(b.Elements)
				if iVal, _, isInt, ok := CoerceToNumeric(args[0]); ok && isInt {
					start = int(iVal)
					if start < 0 {
						start = len(b.Elements) + start
					}
					if start < 0 {
						start = 0
					}
					if start > len(b.Elements) {
						start = len(b.Elements)
					}
				}
				if len(args) > 1 {
					if iVal, _, isInt, ok := CoerceToNumeric(args[1]); ok && isInt {
						end = int(iVal)
						if end < 0 {
							end = len(b.Elements) + end
						}
						if end < start {
							end = start
						}
						if end > len(b.Elements) {
							end = len(b.Elements)
						}
					}
				}
				sliced := make([]byte, end-start)
				copy(sliced, b.Elements[start:end])
				return &Bytes{Elements: sliced}
			},
		}, true

	case "panjang", "length", "len":
		return &Integer{Value: int64(len(b.Elements))}, true

	case "salin", "clone":
		return &Builtin{
			Fn: func(args ...Object) Object {
				cp := make([]byte, len(b.Elements))
				copy(cp, b.Elements)
				return &Bytes{Elements: cp}
			},
		}, true
	}

	return nil, false
}

