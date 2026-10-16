package evaluator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func initStdlibJSON(env *Environment) {
	// -------------------------------------------------------------
	// 1. Parsing & Serialisasi Dasar
	// -------------------------------------------------------------

	// uraikan / parse(str)
	uraikanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.uraikan mengharapkan 1 argumen teks JSON", "Contoh: json.uraikan(\"{\\\"nama\\\": \\\"Celang\\\"}\")")
			}
			jsonStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen json.uraikan harus teks, menerima %s", args[0].Type()), "")
			}

			decoder := json.NewDecoder(strings.NewReader(jsonStr.Value))
			decoder.UseNumber()

			var raw any
			if err := decoder.Decode(&raw); err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengurai teks JSON: %v", err), "Periksa kembali sintaks JSON Anda (tanda kutip ganda, koma, kurung kurawal).")
			}

			return convertGoJSONToCelang(raw)
		},
	}
	definePublicFunction(env, "uraikan", uraikanFn)
	definePublicFunction(env, "parse", uraikanFn)

	// format / stringify(obj, [spasi])
	formatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "json.format mengharapkan 1 atau 2 argumen: format(obj, [spasi])", "Contoh: json.format(data) atau json.format(data, 2)")
			}

			nativeVal, err := convertCelangToJSONNative(args[0])
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mengubah objek ke JSON: %v", err), "")
			}

			indentStr := ""
			if len(args) == 2 {
				switch ind := args[1].(type) {
				case *Integer:
					if ind.Value > 0 {
						indentStr = strings.Repeat(" ", int(ind.Value))
					}
				case *String:
					indentStr = ind.Value
				case *Boolean:
					if ind.Value {
						indentStr = "  "
					}
				}
			}

			var out []byte
			var mErr error
			if indentStr != "" {
				out, mErr = json.MarshalIndent(nativeVal, "", indentStr)
			} else {
				out, mErr = json.Marshal(nativeVal)
			}

			if mErr != nil {
				return newError(nil, fmt.Sprintf("Gagal memformat JSON: %v", mErr), "")
			}

			return &String{Value: string(out)}
		},
	}
	definePublicFunction(env, "format", formatFn)
	definePublicFunction(env, "stringify", formatFn)

	// -------------------------------------------------------------
	// 2. Pemeriksaan & Kerapihan
	// -------------------------------------------------------------

	// apakah_valid / is_valid(str) -> boolean
	apakahValidFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return FALSE
			}
			jsonStr, ok := args[0].(*String)
			if !ok {
				return FALSE
			}
			raw := strings.TrimSpace(jsonStr.Value)
			if raw == "" {
				return FALSE
			}
			var js any
			err := json.Unmarshal([]byte(raw), &js)
			return nativeBoolToBoolean(err == nil)
		},
	}
	definePublicFunction(env, "apakah_valid", apakahValidFn)
	definePublicFunction(env, "is_valid", apakahValidFn)

	// rapi / pretty(obj/str, [spasi=2])
	rapiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "json.rapi mengharapkan 1 atau 2 argumen: json.rapi(obj/str, [spasi=2])", "")
			}
			spasiObj := Object(&Integer{Value: 2})
			if len(args) == 2 {
				spasiObj = args[1]
			}

			target := args[0]
			if str, ok := target.(*String); ok {
				decoder := json.NewDecoder(strings.NewReader(str.Value))
				decoder.UseNumber()
				var raw any
				if err := decoder.Decode(&raw); err == nil {
					target = convertGoJSONToCelang(raw)
				}
			}
			return formatFn.Fn(target, spasiObj)
		},
	}
	definePublicFunction(env, "rapi", rapiFn)
	definePublicFunction(env, "pretty", rapiFn)

	// padatkan / minify(obj/str)
	padatkanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.padatkan mengharapkan 1 argumen: json.padatkan(obj/str)", "")
			}
			target := args[0]
			if str, ok := target.(*String); ok {
				decoder := json.NewDecoder(strings.NewReader(str.Value))
				decoder.UseNumber()
				var raw any
				if err := decoder.Decode(&raw); err == nil {
					target = convertGoJSONToCelang(raw)
				}
			}
			return formatFn.Fn(target)
		},
	}
	definePublicFunction(env, "padatkan", padatkanFn)
	definePublicFunction(env, "minify", padatkanFn)

	// -------------------------------------------------------------
	// 3. Akses & Modifikasi Bersarang via Path (ambil, setel, hapus)
	// -------------------------------------------------------------

	// ambil / get(obj, path, [default])
	ambilFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 || len(args) > 3 {
				return newError(nil, "json.dapatkan mengharapkan 2 atau 3 argumen: json.dapatkan(obj, path, [nilai_bawaan])", "Contoh: json.dapatkan(data, \"pengguna.alamat.kota\", \"Jakarta\")")
			}
			pathStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "argumen path pada json.dapatkan harus berupa teks", "")
			}
			var defVal Object = NULL
			if len(args) == 3 {
				defVal = args[2]
			}
			return getNestedValue(args[0], pathStr.Value, defVal)
		},
	}
	definePublicFunction(env, "dapatkan", ambilFn)
	definePublicFunction(env, "petik", ambilFn)
	definePublicFunction(env, "get", ambilFn)

	// setel / set(obj, path, nilai)
	setelFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "json.setel mengharapkan 3 argumen: json.setel(obj, path, nilai)", "Contoh: json.setel(data, \"pengguna.alamat.kota\", \"Bandung\")")
			}
			pathStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "argumen path pada json.setel harus berupa teks", "")
			}
			return setNestedValue(args[0], pathStr.Value, args[2])
		},
	}
	definePublicFunction(env, "setel", setelFn)
	definePublicFunction(env, "set", setelFn)

	// hapus / delete(obj, path)
	hapusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "json.hapus mengharapkan 2 argumen: json.hapus(obj, path)", "Contoh: json.hapus(data, \"pengguna.kata_sandi\")")
			}
			pathStr, ok := args[1].(*String)
			if !ok {
				return newError(nil, "argumen path pada json.hapus harus berupa teks", "")
			}
			return deleteNestedValue(args[0], pathStr.Value)
		},
	}
	definePublicFunction(env, "hapus", hapusFn)
	definePublicFunction(env, "delete", hapusFn)

	// -------------------------------------------------------------
	// 4. Penggabungan Mendalam (Deep Merge)
	// -------------------------------------------------------------

	// gabung / merge(...objs)
	gabungFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return newError(nil, "json.gabung mengharapkan minimal 2 objek Kamus untuk digabungkan", "Contoh: json.gabung(defaultConfig, userConfig)")
			}
			firstMap, ok := args[0].(*Map)
			if !ok {
				return newError(nil, "seluruh argumen json.gabung harus berupa Kamus (Map)", "")
			}
			result := firstMap
			for i := 1; i < len(args); i++ {
				nextMap, ok2 := args[i].(*Map)
				if !ok2 {
					return newError(nil, "seluruh argumen json.gabung harus berupa Kamus (Map)", "")
				}
				result = deepMergeMaps(result, nextMap)
			}
			return result
		},
	}
	definePublicFunction(env, "gabung", gabungFn)
	definePublicFunction(env, "merge", gabungFn)

	// -------------------------------------------------------------
	// 5. Sensor Data Sensitif (Masking untuk Logging)
	// -------------------------------------------------------------

	// samarkan / mask(obj, [kunciSensitif])
	samarkanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "json.samarkan mengharapkan 1 atau 2 argumen: json.samarkan(obj, [daftar_kunci])", "")
			}
			var customKeys []string
			if len(args) == 2 {
				if arr, ok := args[1].(*Array); ok {
					for _, el := range arr.Elements {
						if s, okS := el.(*String); okS {
							customKeys = append(customKeys, strings.ToLower(s.Value))
						}
					}
				}
			}
			return maskSensitiveJSON(args[0], customKeys)
		},
	}
	definePublicFunction(env, "samarkan", samarkanFn)
	definePublicFunction(env, "mask", samarkanFn)

	// -------------------------------------------------------------
	// 6. JSON Lines / NDJSON (Streaming & Log)
	// -------------------------------------------------------------

	// uraikan_baris / parse_lines(teksJSONL) -> Array of Objects
	uraikanBarisFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.uraikan_baris mengharapkan 1 argumen teks JSONL", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, "argumen json.uraikan_baris harus berupa teks", "")
			}
			scanner := bufio.NewScanner(strings.NewReader(str.Value))
			var elements []Object
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				var raw any
				dec := json.NewDecoder(strings.NewReader(line))
				dec.UseNumber()
				if err := dec.Decode(&raw); err != nil {
					return newError(nil, fmt.Sprintf("Gagal mengurai JSONL pada baris %d: %v", lineNum, err), "")
				}
				elements = append(elements, convertGoJSONToCelang(raw))
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "uraikan_baris", uraikanBarisFn)
	definePublicFunction(env, "parse_lines", uraikanBarisFn)

	// format_baris / stringify_lines(daftar) -> String JSONL
	formatBarisFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.format_baris mengharapkan 1 argumen Daftar (Array)", "")
			}
			arr, ok := args[0].(*Array)
			if !ok {
				return newError(nil, "argumen json.format_baris harus berupa Daftar (Array)", "")
			}
			var b strings.Builder
			for i, el := range arr.Elements {
				nativeVal, err := convertCelangToJSONNative(el)
				if err != nil {
					return newError(nil, fmt.Sprintf("Gagal serialisasi elemen ke-%d: %v", i, err), "")
				}
				lineBytes, err := json.Marshal(nativeVal)
				if err != nil {
					return newError(nil, fmt.Sprintf("Gagal marshal JSON elemen ke-%d: %v", i, err), "")
				}
				b.Write(lineBytes)
				b.WriteByte('\n')
			}
			return &String{Value: b.String()}
		},
	}
	definePublicFunction(env, "format_baris", formatBarisFn)
	definePublicFunction(env, "stringify_lines", formatBarisFn)

	// -------------------------------------------------------------
	// 7. Flatten & Unflatten (Ratakan & Susun)
	// -------------------------------------------------------------

	// ratakan / flatten(obj)
	ratakanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.ratakan mengharapkan 1 argumen objek", "")
			}
			resultPairs := make(map[string]Object)
			var resultOrder []string
			flattenObject("", args[0], resultPairs, &resultOrder)
			return &Map{Pairs: resultPairs, Order: resultOrder}
		},
	}
	definePublicFunction(env, "ratakan", ratakanFn)
	definePublicFunction(env, "flatten", ratakanFn)

	// susun / unflatten(objDatar)
	susunFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "json.susun mengharapkan 1 argumen Kamus datar", "")
			}
			mp, ok := args[0].(*Map)
			if !ok {
				return newError(nil, "argumen json.susun harus berupa Kamus (Map)", "")
			}
			root := &Map{Pairs: make(map[string]Object), Order: []string{}}
			for _, k := range mp.Order {
				setNestedValue(root, k, mp.Pairs[k])
			}
			for k, v := range mp.Pairs {
				if _, okP := root.Pairs[k]; !okP {
					setNestedValue(root, k, v)
				}
			}
			return root
		},
	}
	definePublicFunction(env, "susun", susunFn)
	definePublicFunction(env, "unflatten", susunFn)
}

// -------------------------------------------------------------
// Helper Functions (Path Parser, Nested Value, Deep Merge, Masking)
// -------------------------------------------------------------

type pathToken struct {
	isIndex bool
	key     string
	index   int
}

func parseJSONPath(path string) []pathToken {
	var tokens []pathToken
	var cur strings.Builder
	runes := []rune(path)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '.' {
			if cur.Len() > 0 {
				tokens = append(tokens, pathToken{key: cur.String()})
				cur.Reset()
			}
		} else if r == '[' {
			if cur.Len() > 0 {
				tokens = append(tokens, pathToken{key: cur.String()})
				cur.Reset()
			}
			j := i + 1
			var idxStr strings.Builder
			for j < len(runes) && runes[j] != ']' {
				idxStr.WriteRune(runes[j])
				j++
			}
			if idx, err := strconv.Atoi(idxStr.String()); err == nil {
				tokens = append(tokens, pathToken{isIndex: true, index: idx})
			}
			i = j
		} else {
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, pathToken{key: cur.String()})
	}
	return tokens
}

func getNestedValue(root Object, path string, defaultVal Object) Object {
	tokens := parseJSONPath(path)
	cur := root
	for _, tok := range tokens {
		if cur == nil || cur.Type() == NULL_OBJ {
			return defaultVal
		}
		if tok.isIndex {
			arr, ok := cur.(*Array)
			if !ok || tok.index < 0 || tok.index >= len(arr.Elements) {
				return defaultVal
			}
			cur = arr.Elements[tok.index]
		} else {
			switch o := cur.(type) {
			case *Map:
				val, exists := o.Pairs[tok.key]
				if !exists {
					return defaultVal
				}
				cur = val
			case *StructInstance:
				val, exists := o.Fields[tok.key]
				if !exists {
					return defaultVal
				}
				cur = val
			default:
				return defaultVal
			}
		}
	}
	if cur == nil {
		return defaultVal
	}
	return cur
}

func setNestedValue(root Object, path string, val Object) Object {
	tokens := parseJSONPath(path)
	if len(tokens) == 0 {
		return root
	}
	cur := root
	for i := 0; i < len(tokens)-1; i++ {
		tok := tokens[i]
		nextTok := tokens[i+1]

		if tok.isIndex {
			arr, ok := cur.(*Array)
			if !ok {
				return root
			}
			for len(arr.Elements) <= tok.index {
				arr.Elements = append(arr.Elements, NULL)
			}
			target := arr.Elements[tok.index]
			if target == nil || target.Type() == NULL_OBJ {
				if nextTok.isIndex {
					target = &Array{Elements: []Object{}}
				} else {
					target = &Map{Pairs: make(map[string]Object), Order: []string{}}
				}
				arr.Elements[tok.index] = target
			}
			cur = target
		} else {
			mp, ok := cur.(*Map)
			if !ok {
				if si, isSi := cur.(*StructInstance); isSi {
					target, exists := si.Fields[tok.key]
					if !exists || target == nil || target.Type() == NULL_OBJ {
						if nextTok.isIndex {
							target = &Array{Elements: []Object{}}
						} else {
							target = &Map{Pairs: make(map[string]Object), Order: []string{}}
						}
						si.Fields[tok.key] = target
					}
					cur = target
					continue
				}
				return root
			}
			target, exists := mp.Pairs[tok.key]
			if !exists || target == nil || target.Type() == NULL_OBJ {
				if nextTok.isIndex {
					target = &Array{Elements: []Object{}}
				} else {
					target = &Map{Pairs: make(map[string]Object), Order: []string{}}
				}
				mp.Pairs[tok.key] = target
				if !exists {
					mp.Order = append(mp.Order, tok.key)
				}
			}
			cur = target
		}
	}

	lastTok := tokens[len(tokens)-1]
	if lastTok.isIndex {
		if arr, ok := cur.(*Array); ok {
			for len(arr.Elements) <= lastTok.index {
				arr.Elements = append(arr.Elements, NULL)
			}
			arr.Elements[lastTok.index] = val
		}
	} else {
		if mp, ok := cur.(*Map); ok {
			if _, exists := mp.Pairs[lastTok.key]; !exists {
				mp.Order = append(mp.Order, lastTok.key)
			}
			mp.Pairs[lastTok.key] = val
		} else if si, ok := cur.(*StructInstance); ok {
			si.Fields[lastTok.key] = val
		}
	}
	return root
}

func deleteNestedValue(root Object, path string) Object {
	tokens := parseJSONPath(path)
	if len(tokens) == 0 {
		return root
	}
	cur := root
	for i := 0; i < len(tokens)-1; i++ {
		tok := tokens[i]
		if tok.isIndex {
			arr, ok := cur.(*Array)
			if !ok || tok.index < 0 || tok.index >= len(arr.Elements) {
				return root
			}
			cur = arr.Elements[tok.index]
		} else {
			if mp, ok := cur.(*Map); ok {
				target, exists := mp.Pairs[tok.key]
				if !exists {
					return root
				}
				cur = target
			} else if si, ok := cur.(*StructInstance); ok {
				target, exists := si.Fields[tok.key]
				if !exists {
					return root
				}
				cur = target
			} else {
				return root
			}
		}
	}

	lastTok := tokens[len(tokens)-1]
	if lastTok.isIndex {
		if arr, ok := cur.(*Array); ok {
			if lastTok.index >= 0 && lastTok.index < len(arr.Elements) {
				arr.Elements = append(arr.Elements[:lastTok.index], arr.Elements[lastTok.index+1:]...)
			}
		}
	} else {
		if mp, ok := cur.(*Map); ok {
			delete(mp.Pairs, lastTok.key)
			newOrder := make([]string, 0, len(mp.Order))
			for _, k := range mp.Order {
				if k != lastTok.key {
					newOrder = append(newOrder, k)
				}
			}
			mp.Order = newOrder
		} else if si, ok := cur.(*StructInstance); ok {
			delete(si.Fields, lastTok.key)
		}
	}
	return root
}

func deepMergeMaps(a, b *Map) *Map {
	resPairs := make(map[string]Object)
	var resOrder []string
	seen := make(map[string]bool)

	for _, k := range a.Order {
		resPairs[k] = a.Pairs[k]
		resOrder = append(resOrder, k)
		seen[k] = true
	}
	for k, v := range a.Pairs {
		if !seen[k] {
			resPairs[k] = v
			resOrder = append(resOrder, k)
			seen[k] = true
		}
	}

	for _, k := range b.Order {
		bVal := b.Pairs[k]
		if aVal, exists := resPairs[k]; exists {
			aMap, aOk := aVal.(*Map)
			bMap, bOk := bVal.(*Map)
			if aOk && bOk {
				resPairs[k] = deepMergeMaps(aMap, bMap)
			} else {
				resPairs[k] = bVal
			}
		} else {
			resPairs[k] = bVal
			resOrder = append(resOrder, k)
			seen[k] = true
		}
	}
	for k, bVal := range b.Pairs {
		if !seen[k] {
			resPairs[k] = bVal
			resOrder = append(resOrder, k)
			seen[k] = true
		}
	}
	return &Map{Pairs: resPairs, Order: resOrder}
}

var defaultSensitiveKeywords = []string{
	"password", "kata_sandi", "token", "pin", "kartu",
	"secret", "rahasia", "cvv", "nik", "api_key", "auth",
}

func maskSensitiveJSON(obj Object, customKeys []string) Object {
	keysToMask := defaultSensitiveKeywords
	if len(customKeys) > 0 {
		keysToMask = append(keysToMask, customKeys...)
	}

	isSensitive := func(k string) bool {
		lowerK := strings.ToLower(k)
		for _, sens := range keysToMask {
			if strings.Contains(lowerK, sens) {
				return true
			}
		}
		return false
	}

	switch o := obj.(type) {
	case *Map:
		newPairs := make(map[string]Object)
		for k, v := range o.Pairs {
			if isSensitive(k) {
				newPairs[k] = &String{Value: "***"}
			} else {
				newPairs[k] = maskSensitiveJSON(v, customKeys)
			}
		}
		return &Map{Pairs: newPairs, Order: o.Order}
	case *Array:
		newElems := make([]Object, len(o.Elements))
		for i, el := range o.Elements {
			newElems[i] = maskSensitiveJSON(el, customKeys)
		}
		return &Array{Elements: newElems}
	default:
		return obj
	}
}

func flattenObject(prefix string, obj Object, result map[string]Object, order *[]string) {
	switch o := obj.(type) {
	case *Map:
		for _, k := range o.Order {
			v := o.Pairs[k]
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			flattenObject(p, v, result, order)
		}
		for k, v := range o.Pairs {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			if _, exists := result[p]; !exists {
				flattenObject(p, v, result, order)
			}
		}
	case *Array:
		for i, el := range o.Elements {
			p := fmt.Sprintf("%s[%d]", prefix, i)
			flattenObject(p, el, result, order)
		}
	default:
		if prefix != "" {
			if _, exists := result[prefix]; !exists {
				*order = append(*order, prefix)
			}
			result[prefix] = obj
		}
	}
}

func convertGoJSONToCelang(val any) Object {
	if val == nil {
		return &Null{}
	}

	switch v := val.(type) {
	case bool:
		if v {
			return TRUE
		}
		return FALSE
	case json.Number:
		s := v.String()
		if !strings.Contains(s, ".") && !strings.Contains(s, "e") && !strings.Contains(s, "E") {
			if i, err := v.Int64(); err == nil {
				return &Integer{Value: i}
			}
		}
		if f, err := v.Float64(); err == nil {
			return &Float{Value: f}
		}
		return &String{Value: s}
	case string:
		return &String{Value: v}
	case []any:
		elements := make([]Object, len(v))
		for i, item := range v {
			elements[i] = convertGoJSONToCelang(item)
		}
		return &Array{Elements: elements}
	case map[string]any:
		pairs := make(map[string]Object, len(v))
		order := make([]string, 0, len(v))
		for k := range v {
			order = append(order, k)
		}
		sort.Strings(order)
		for _, k := range order {
			pairs[k] = convertGoJSONToCelang(v[k])
		}
		return &Map{Pairs: pairs, Order: order}
	default:
		return &String{Value: fmt.Sprintf("%v", v)}
	}
}

func convertCelangToJSONNative(obj Object) (any, error) {
	return convertCelangToJSONNativeWithVisited(obj, make(map[any]bool))
}

func convertCelangToJSONNativeWithVisited(obj Object, visited map[any]bool) (any, error) {
	if obj == nil {
		return nil, nil
	}

	switch o := obj.(type) {
	case *Null:
		return nil, nil
	case *Boolean:
		return o.Value, nil
	case *Integer:
		return o.Value, nil
	case *Float:
		return o.Value, nil
	case *String:
		return o.Value, nil
	case *Array:
		if visited[o] {
			return nil, fmt.Errorf("referensi sirkular terdeteksi pada array saat serialisasi JSON")
		}
		visited[o] = true
		defer delete(visited, o)
		res := make([]any, len(o.Elements))
		for i, el := range o.Elements {
			c, err := convertCelangToJSONNativeWithVisited(el, visited)
			if err != nil {
				return nil, err
			}
			res[i] = c
		}
		return res, nil
	case *Map:
		if visited[o] {
			return nil, fmt.Errorf("referensi sirkular terdeteksi pada map saat serialisasi JSON")
		}
		visited[o] = true
		defer delete(visited, o)
		res := make(map[string]any)
		for k, v := range o.Pairs {
			c, err := convertCelangToJSONNativeWithVisited(v, visited)
			if err != nil {
				return nil, err
			}
			res[k] = c
		}
		return res, nil
	case *StructInstance:
		if visited[o] {
			return nil, fmt.Errorf("referensi sirkular terdeteksi pada bentuk '%s' saat serialisasi JSON", o.Def.Name)
		}
		visited[o] = true
		defer delete(visited, o)
		res := make(map[string]any)
		for k, v := range o.Fields {
			c, err := convertCelangToJSONNativeWithVisited(v, visited)
			if err != nil {
				return nil, err
			}
			res[k] = c
		}
		return res, nil
	case *ReturnValue:
		return convertCelangToJSONNativeWithVisited(o.Value, visited)
	default:
		return nil, fmt.Errorf("tipe %s tidak dapat diserialisasi ke JSON", obj.Type())
	}
}
