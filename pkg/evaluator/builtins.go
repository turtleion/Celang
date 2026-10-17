package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func RegisterBuiltins(env *Environment, out io.Writer) {
	RegisterBuiltinsWithIO(env, os.Stdin, out)
}

func RegisterBuiltinsWithIO(env *Environment, in io.Reader, out io.Writer) {
	if in == nil {
		in = os.Stdin
	}
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

	// 'masukan' / 'input' prompts and reads a line from input
	reader := bufio.NewReader(in)
	masukanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				if s, ok := args[0].(*String); ok {
					fmt.Fprint(out, s.Value)
				} else {
					fmt.Fprint(out, args[0].Inspect())
				}
			}
			line, err := reader.ReadString('\n')
			if err != nil && len(line) == 0 {
				return NULL
			}
			line = strings.TrimRight(line, "\r\n")
			return &String{Value: line}
		},
	}
	env.Define("masukan", masukanFn, false, "function")
	env.Define("input", masukanFn, false, "function")

	// 'panjang' / 'len' calculates length of strings or arrays
	panjangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("panjang mengharapkan 1 argumen, menerima %d", len(args))}
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *Map:
				return &Integer{Value: int64(len(arg.Pairs))}
			case *Channel:
				return &Integer{Value: int64(len(arg.Ch))}
			default:
				return &Error{Message: fmt.Sprintf("tipe %s tidak memiliki panjang", arg.Type())}
			}
		},
	}
	env.Define("panjang", panjangFn, false, "function")
	env.Define("len", panjangFn, false, "function")
	env.Define("length", panjangFn, false, "function")

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
				s, _, isInt1, ok1 := CoerceToNumeric(args[0])
				e, _, isInt2, ok2 := CoerceToNumeric(args[1])
				if !ok1 || !ok2 || !isInt1 || !isInt2 {
					return &Error{Message: "argumen rentang harus bilangan bulat"}
				}
				start = s
				end = e
			}
			diff := end - start
			if diff < 0 {
				diff = -diff
			}
			if diff > 1_000_000 {
				return &Error{Message: "ukuran rentang terlalu besar (maksimal 1.000.000 elemen)"}
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

	// 'saring' / 'filter' filters elements of an array using a predicate function
	saringFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return &Error{Message: "saring mengharapkan 2 argumen: saring(array, fungsi)"}
			}
			arr, ok := args[0].(*Array)
			if !ok {
				return &Error{Message: fmt.Sprintf("argumen pertama saring harus array, menerima %s", args[0].Type())}
			}
			fn := args[1]
			var filtered []Object
			for _, el := range arr.Elements {
				res := applyFunction(fn, []Object{el}, nil)
				if isError(res) {
					return res
				}
				if CoerceToBoolean(res) {
					filtered = append(filtered, el)
				}
			}
			return &Array{Elements: filtered}
		},
	}
	env.Define("saring", saringFn, false, "function")
	env.Define("filter", saringFn, false, "function")

	// 'petakan' / 'map' transforms elements of an array using a mapping function
	petakanFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return &Error{Message: "petakan mengharapkan 2 argumen: petakan(array, fungsi)"}
			}
			arr, ok := args[0].(*Array)
			if !ok {
				return &Error{Message: fmt.Sprintf("argumen pertama petakan harus array, menerima %s", args[0].Type())}
			}
			fn := args[1]
			mapped := make([]Object, len(arr.Elements))
			for i, el := range arr.Elements {
				res := applyFunction(fn, []Object{el}, nil)
				if isError(res) {
					return res
				}
				mapped[i] = res
			}
			return &Array{Elements: mapped}
		},
	}
	env.Define("petakan", petakanFn, false, "function")
	env.Define("map", petakanFn, false, "function")

	// 'potong_spasi' / 'trim' trims leading and trailing whitespace
	potongSpasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: "potong_spasi mengharapkan 1 argumen"}
			}
			str, ok := args[0].(*String)
			if !ok {
				return &Error{Message: fmt.Sprintf("argumen potong_spasi harus teks, menerima %s", args[0].Type())}
			}
			return &String{Value: strings.TrimSpace(str.Value)}
		},
	}
	env.Define("potong_spasi", potongSpasiFn, false, "function")
	env.Define("trim", potongSpasiFn, false, "function")

	// 'kapitalisasi' / 'capitalize' / 'uppercase' converts text to uppercase
	kapitalisasiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: "kapitalisasi mengharapkan 1 argumen"}
			}
			str, ok := args[0].(*String)
			if !ok {
				return &Error{Message: fmt.Sprintf("argumen kapitalisasi harus teks, menerima %s", args[0].Type())}
			}
			return &String{Value: strings.ToUpper(str.Value)}
		},
	}
	env.Define("kapitalisasi", kapitalisasiFn, false, "function")
	env.Define("capitalize", kapitalisasiFn, false, "function")
	env.Define("uppercase", kapitalisasiFn, false, "function")

	// Type casting constructors: Angka(), Teks(), Logika()
	angkaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("Angka mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return castValueToType(args[0], "Angka", nil, env)
		},
	}
	env.Define("Angka", angkaFn, false, "function")
	env.Define("angka", angkaFn, false, "function")
	env.Define("Number", angkaFn, false, "function")
	env.Define("number", angkaFn, false, "function")

	teksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("Teks mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return castValueToType(args[0], "Teks", nil, env)
		},
	}
	env.Define("Teks", teksFn, false, "function")
	env.Define("teks", teksFn, false, "function")
	env.Define("String", teksFn, false, "function")
	env.Define("string", teksFn, false, "function")

	logikaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("Logika mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return castValueToType(args[0], "Logika", nil, env)
		},
	}
	env.Define("Logika", logikaFn, false, "function")
	env.Define("logika", logikaFn, false, "function")
	env.Define("Boolean", logikaFn, false, "function")
	env.Define("boolean", logikaFn, false, "function")

	// Indonesian Culture Builtins: terbilang(angka), rupiah(angka) / format_rupiah(angka)
	terbilangFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("terbilang mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return &String{Value: Terbilang(args[0])}
		},
	}
	env.Define("terbilang", terbilangFn, false, "function")

	rupiahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("format_rupiah mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return &String{Value: FormatRupiah(args[0])}
		},
	}
	env.Define("rupiah", rupiahFn, false, "function")
	env.Define("format_rupiah", rupiahFn, false, "function")

	// International Currency & Word Builtins: format_usd(n), in_words(n), dollars_in_words(n)
	usdFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("format_usd expects 1 argument, received %d", len(args))}
			}
			return &String{Value: FormatUSD(args[0])}
		},
	}
	env.Define("format_usd", usdFn, false, "function")
	env.Define("format_currency", usdFn, false, "function")

	inWordsFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("in_words expects 1 argument, received %d", len(args))}
			}
			return &String{Value: WordsEnglishCoerced(args[0])}
		},
	}
	env.Define("in_words", inWordsFn, false, "function")
	env.Define("spelled_out", inWordsFn, false, "function")

	dollarsInWordsFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("dollars_in_words expects 1 argument, received %d", len(args))}
			}
			return &String{Value: DollarsInWords(args[0])}
		},
	}
	env.Define("dollars_in_words", dollarsInWordsFn, false, "function")

	// Deep Memory Clone Builtin: salin(obj) / clone(obj)
	salinFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("salin mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return deepClone(args[0])
		},
	}
	env.Define("salin", salinFn, false, "function")
	env.Define("clone", salinFn, false, "function")

	// Modern Concurrency Builtins (Saluran, jemput, jemput_semua, jemput_pertama)
	saluranFn := &Builtin{
		Fn: func(args ...Object) Object {
			cap := 0
			if len(args) > 0 {
				if iVal, _, isInt, ok := CoerceToNumeric(args[0]); ok && isInt {
					cap = int(iVal)
				}
			}
			return NewChannel(cap, "")
		},
	}
	env.Define("Saluran", saluranFn, false, "function")
	env.Define("saluran", saluranFn, false, "function")
	env.Define("Channel", saluranFn, false, "function")
	env.Define("channel", saluranFn, false, "function")

	jemputFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &Error{Message: "jemput mengharapkan setidaknya 1 argumen (tiket)"}
			}
			ticket, ok := args[0].(*Ticket)
			if !ok {
				return &Error{Message: fmt.Sprintf("Target jemput harus berupa Tiket, bukan %s", args[0].Type())}
			}
			var timeout time.Duration = 0
			if len(args) > 1 {
				timeout = parseTimeoutValue(args[1])
			}
			return ticket.Await(timeout, nil)
		},
	}
	env.Define("jemput", jemputFn, false, "function")
	env.Define("await", jemputFn, false, "function")

	jemputSemuaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &Error{Message: "jemput_semua mengharapkan sebuah array berisi tiket"}
			}
			arr, ok := args[0].(*Array)
			if !ok {
				return &Error{Message: "argumen jemput_semua harus berupa Daftar (Array)"}
			}
			results := make([]Object, len(arr.Elements))
			var firstErr *Error
			var mu sync.Mutex
			var wg sync.WaitGroup

			for i, el := range arr.Elements {
				ticket, ok := el.(*Ticket)
				if !ok {
					return &Error{Message: fmt.Sprintf("Elemen ke-%d bukan merupakan Tiket", i)}
				}
				wg.Add(1)
				idx := i
				t := ticket
				go func() {
					defer wg.Done()
					res := t.Await(0, nil)
					mu.Lock()
					defer mu.Unlock()
					if isError(res) && firstErr == nil {
						firstErr = res.(*Error)
					}
					results[idx] = res
				}()
			}
			wg.Wait()

			if firstErr != nil {
				return firstErr
			}
			return &Array{Elements: results}
		},
	}
	env.Define("jemput_semua", jemputSemuaFn, false, "function")
	env.Define("await_all", jemputSemuaFn, false, "function")

	jemputPertamaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &Error{Message: "jemput_pertama mengharapkan sebuah array berisi tiket"}
			}
			arr, ok := args[0].(*Array)
			if !ok || len(arr.Elements) == 0 {
				return &Error{Message: "argumen jemput_pertama harus berupa Daftar tiket yang tidak kosong"}
			}
			resChan := make(chan Object, len(arr.Elements))
			for _, el := range arr.Elements {
				ticket, ok := el.(*Ticket)
				if !ok {
					return &Error{Message: "Semua elemen harus berupa Tiket"}
				}
				t := ticket
				go func() {
					resChan <- t.Await(0, nil)
				}()
			}
			return <-resChan
		},
	}
	env.Define("jemput_pertama", jemputPertamaFn, false, "function")
	env.Define("await_race", jemputPertamaFn, false, "function")
	env.Define("await_first", jemputPertamaFn, false, "function")

	// Register Standard Library global objects (berkas, fs, matematika, math, teks, strings, waktu, time, json)
	RegisterStdlibGlobals(env)
}

func deepClone(obj Object) Object {
	if obj == nil {
		return NULL
	}
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: o.Value}
	case *Float:
		return &Float{Value: o.Value}
	case *String:
		return &String{Value: o.Value}
	case *Boolean:
		return &Boolean{Value: o.Value}
	case *Null:
		return &Null{}
	case *Array:
		newElems := make([]Object, len(o.Elements))
		for i, el := range o.Elements {
			newElems[i] = deepClone(el)
		}
		return &Array{Elements: newElems}
	case *Map:
		newPairs := make(map[string]Object, len(o.Pairs))
		for k, v := range o.Pairs {
			newPairs[k] = deepClone(v)
		}
		var newOrder []string
		if o.Order != nil {
			newOrder = make([]string, len(o.Order))
			copy(newOrder, o.Order)
		}
		return &Map{Pairs: newPairs, Order: newOrder}
	case *StructInstance:
		newFields := make(map[string]Object, len(o.Fields))
		for k, v := range o.Fields {
			newFields[k] = deepClone(v)
		}
		cloneInst := &StructInstance{
			Def:    o.Def,
			Fields: newFields,
			Env:    NewEnclosedEnvironment(o.Def.Env),
		}
		cloneInst.Env.Define("diri", cloneInst, false, o.Def.Name)
		cloneInst.Env.Define("this", cloneInst, false, o.Def.Name)
		cloneInst.Env.Define("self", cloneInst, false, o.Def.Name)
		return cloneInst
	default:
		return o
	}
}

