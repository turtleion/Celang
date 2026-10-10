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
	env.Define("cetak", &Builtin{
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
	}, false, "function")

	// 'panjang' calculates length of strings or arrays
	env.Define("panjang", &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("panjang mengharapkan 1 argumen, menerima %d", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return &Error{Message: fmt.Sprintf("tipe %s tidak memiliki panjang", arg.Type())}
			}
		},
	}, false, "function")

	// 'tipe' returns the type name of an object
	env.Define("tipe", &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: "tipe mengharapkan 1 argumen"}
			}
			return &String{Value: string(args[0].Type())}
		},
	}, false, "function")
}
