package evaluator

import (
	"fmt"
	"io"
	"os"
)

func RegisterBuiltins(env *Environment, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	// 'cetak' prints values to the specified output
	cetakFn := &Builtin{
		Fn: func(args ...Object) Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, arg.Inspect())
			}
			fmt.Fprintln(out)
			return &Null{}
		},
	}
	env.Define("cetak", cetakFn, false, "function")
	env.Define("print", cetakFn, false, "function")

	// 'panjang' / 'len' calculates length of strings or arrays
	panjangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("panjang mengharapkan 1 argumen, menerima %d", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *Map:
				return &Integer{Value: int64(len(arg.Pairs))}
			default:
				return &Error{Message: fmt.Sprintf("tipe %s tidak memiliki panjang", arg.Type())}
			}
		},
	}
	env.Define("panjang", panjangFn, false, "function")
	env.Define("len", panjangFn, false, "function")

	// 'tipe' returns the type name of an object
	tipeFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: "tipe mengharapkan 1 argumen"}
			}
			return &String{Value: string(args[0].Type())}
		},
	}
	env.Define("tipe", tipeFn, false, "function")
	env.Define("type", tipeFn, false, "function")

	// 'rentang' / 'range' generates array of integers
	rentangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 || len(args) > 2 {
				return &Error{Message: "rentang mengharapkan 1 atau 2 argumen: rentang(akhir) atau rentang(awal, akhir)"}
			}
			start := int64(0)
			end := int64(0)
			if len(args) == 1 {
				s, _, isInt, ok := CoerceToNumeric(args[0])
				if !ok || !isInt {
					return &Error{Message: "argumen rentang harus bilangan bulat"}
				}
				end = s
			} else {
				s, _, _, ok1 := CoerceToNumeric(args[0])
				e, _, _, ok2 := CoerceToNumeric(args[1])
				if !ok1 || !ok2 {
					return &Error{Message: "argumen rentang harus bilangan bulat"}
				}
				start = s
				end = e
			}
			var elems []Object
			if start <= end {
				for i := start; i <= end; i++ {
					elems = append(elems, &Integer{Value: i})
				}
			} else {
				for i := start; i >= end; i-- {
					elems = append(elems, &Integer{Value: i})
				}
			}
			return &Array{Elements: elems}
		},
	}
	env.Define("rentang", rentangFn, false, "function")
	env.Define("range", rentangFn, false, "function")
}
