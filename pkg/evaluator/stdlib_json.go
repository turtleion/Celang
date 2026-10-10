package evaluator

import (
	"encoding/json"
	"fmt"
	"strings"
)

func initStdlibJSON(env *Environment) {
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

	// format / stringify(obj, spasi?)
	formatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "json.format mengharapkan 1 atau 2 argumen: format(obj, spasi?)", "Contoh: json.format(data) atau json.format(data, 2)")
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
		pairs := make(map[string]Object)
		order := make([]string, 0, len(v))
		for k, item := range v {
			pairs[k] = convertGoJSONToCelang(item)
			order = append(order, k)
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
