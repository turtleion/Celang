package evaluator

import (
	"fmt"
	"time"
)

func initStdlibWaktu(env *Environment) {
	// sekarang / now()
	sekarangFn := &Builtin{
		Fn: func(args ...Object) Object {
			ms := time.Now().UnixNano() / int64(time.Millisecond)
			return &Integer{Value: ms}
		},
	}
	definePublicFunction(env, "sekarang", sekarangFn)
	definePublicFunction(env, "now", sekarangFn)

	// sekarang_detik / now_seconds()
	sekarangDetikFn := &Builtin{
		Fn: func(args ...Object) Object {
			sec := time.Now().Unix()
			return &Integer{Value: sec}
		},
	}
	definePublicFunction(env, "sekarang_detik", sekarangDetikFn)
	definePublicFunction(env, "now_seconds", sekarangDetikFn)

	// format / format(ms, layout)
	formatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return newError(nil, "waktu.format mengharapkan 1 atau 2 argumen: format(ms, layout?)", "Contoh: waktu.format(waktu.sekarang(), \"2006-01-02 15:04:05\")")
			}
			msVal, _, isInt, ok := CoerceToNumeric(args[0])
			if !ok || !isInt {
				return newError(nil, fmt.Sprintf("argumen timestamp waktu.format harus bilangan bulat milidetik, menerima %s", args[0].Type()), "")
			}
			t := time.Unix(0, msVal*int64(time.Millisecond))
			layout := "2006-01-02 15:04:05"
			if len(args) == 2 {
				layoutStr, ok := args[1].(*String)
				if !ok {
					return newError(nil, "argumen layout format waktu harus teks", "")
				}
				layout = layoutStr.Value
			}
			return &String{Value: t.Format(layout)}
		},
	}
	definePublicFunction(env, "format", formatFn)

	// tidur / sleep(ms)
	tidurFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "waktu.tidur mengharapkan 1 argumen: milidetik", "Contoh: waktu.tidur(500)")
			}
			msVal, _, isInt, ok := CoerceToNumeric(args[0])
			if !ok || !isInt {
				return newError(nil, fmt.Sprintf("argumen milidetik waktu.tidur harus bilangan bulat, menerima %s", args[0].Type()), "")
			}
			if msVal < 0 {
				return newError(nil, "durasi waktu.tidur tidak boleh negatif", "")
			}
			time.Sleep(time.Duration(msVal) * time.Millisecond)
			return &Null{}
		},
	}
	definePublicFunction(env, "tidur", tidurFn)
	definePublicFunction(env, "sleep", tidurFn)
}
