package evaluator

import (
	"fmt"
	"strings"
)

func initStdlibTeks(env *Environment) {
	// pisah / split(str, sep)
	pisahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.pisah mengharapkan 2 argumen: string dan pemisah", "Contoh: teks.pisah(\"a,b,c\", \",\")")
			}
			str, ok1 := args[0].(*String)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen untuk teks.pisah harus berupa teks", "")
			}
			parts := strings.Split(str.Value, sep.Value)
			elements := make([]Object, len(parts))
			for i, p := range parts {
				elements[i] = &String{Value: p}
			}
			return &Array{Elements: elements}
		},
	}
	definePublicFunction(env, "pisah", pisahFn)
	definePublicFunction(env, "split", pisahFn)

	// gabung / join(arr, sep)
	gabungFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.gabung mengharapkan 2 argumen: array dan pemisah", "Contoh: teks.gabung([\"a\", \"b\"], \", \")")
			}
			arr, ok1 := args[0].(*Array)
			sep, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.gabung harus berupa (array, pemisah teks)", "")
			}
			strParts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				if s, ok := el.(*String); ok {
					strParts[i] = s.Value
				} else {
					strParts[i] = el.Inspect()
				}
			}
			return &String{Value: strings.Join(strParts, sep.Value)}
		},
	}
	definePublicFunction(env, "gabung", gabungFn)
	definePublicFunction(env, "join", gabungFn)

	// kapital / upper(str)
	kapitalFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kapital mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen teks.kapital harus berupa teks, menerima %s", args[0].Type()), "")
			}
			return &String{Value: strings.ToUpper(str.Value)}
		},
	}
	definePublicFunction(env, "kapital", kapitalFn)
	definePublicFunction(env, "upper", kapitalFn)

	// kecil / lower(str)
	kecilFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.kecil mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen teks.kecil harus berupa teks, menerima %s", args[0].Type()), "")
			}
			return &String{Value: strings.ToLower(str.Value)}
		},
	}
	definePublicFunction(env, "kecil", kecilFn)
	definePublicFunction(env, "lower", kecilFn)

	// potong / trim(str)
	potongFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "teks.potong mengharapkan 1 argumen teks", "")
			}
			str, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen teks.potong harus berupa teks, menerima %s", args[0].Type()), "")
			}
			return &String{Value: strings.TrimSpace(str.Value)}
		},
	}
	definePublicFunction(env, "potong", potongFn)
	definePublicFunction(env, "trim", potongFn)

	// ganti / replace(str, old, new)
	gantiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return newError(nil, "teks.ganti mengharapkan 3 argumen: teks, teks_lama, teks_baru", "Contoh: teks.ganti(\"halo dunia\", \"dunia\", \"celang\")")
			}
			str, ok1 := args[0].(*String)
			oldS, ok2 := args[1].(*String)
			newS, ok3 := args[2].(*String)
			if !ok1 || !ok2 || !ok3 {
				return newError(nil, "seluruh argumen teks.ganti harus berupa teks", "")
			}
			return &String{Value: strings.ReplaceAll(str.Value, oldS.Value, newS.Value)}
		},
	}
	definePublicFunction(env, "ganti", gantiFn)
	definePublicFunction(env, "replace", gantiFn)

	// berisi / contains(str, substr)
	berisiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.berisi mengharapkan 2 argumen: teks dan teks_cari", "")
			}
			str, ok1 := args[0].(*String)
			sub, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "seluruh argumen teks.berisi harus berupa teks", "")
			}
			if strings.Contains(str.Value, sub.Value) {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "berisi", berisiFn)
	definePublicFunction(env, "contains", berisiFn)

	// dimulai_dengan / starts_with(str, prefix)
	dimulaiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.dimulai_dengan mengharapkan 2 argumen: teks dan awalan", "")
			}
			str, ok1 := args[0].(*String)
			pref, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.dimulai_dengan harus berupa teks", "")
			}
			if strings.HasPrefix(str.Value, pref.Value) {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "dimulai_dengan", dimulaiFn)
	definePublicFunction(env, "starts_with", dimulaiFn)

	// diakhiri_dengan / ends_with(str, suffix)
	diakhiriFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "teks.diakhiri_dengan mengharapkan 2 argumen: teks dan akhiran", "")
			}
			str, ok1 := args[0].(*String)
			suff, ok2 := args[1].(*String)
			if !ok1 || !ok2 {
				return newError(nil, "argumen teks.diakhiri_dengan harus berupa teks", "")
			}
			if strings.HasSuffix(str.Value, suff.Value) {
				return TRUE
			}
			return FALSE
		},
	}
	definePublicFunction(env, "diakhiri_dengan", diakhiriFn)
	definePublicFunction(env, "ends_with", diakhiriFn)
}
