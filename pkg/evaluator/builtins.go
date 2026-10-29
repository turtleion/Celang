package evaluator

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"celang/pkg/uast"
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
	env.Define("masukkan", masukanFn, false, "function")
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
			case *Bytes:
				return &Integer{Value: int64(len(arg.Elements))}
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
			if (start <= end && (end-start > 1_000_000 || end-start < 0)) ||
				(start > end && (start-end > 1_000_000 || start-end < 0)) {
				return &Error{Message: "ukuran rentang terlalu besar (maksimal 1.000.000 elemen)"}
			}
			var elems []Object
			if start <= end {
				count := end - start + 1
				for idx := int64(0); idx < count; idx++ {
					elems = append(elems, &Integer{Value: start + idx})
				}
			} else {
				count := start - end + 1
				for idx := int64(0); idx < count; idx++ {
					elems = append(elems, &Integer{Value: start - idx})
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
				return &Error{Message: fmt.Sprintf("Kebenaran mengharapkan 1 argumen, menerima %d", len(args))}
			}
			return castValueToType(args[0], "Kebenaran", nil, env)
		},
	}
	env.Define("Kebenaran", logikaFn, false, "function")
	env.Define("kebenaran", logikaFn, false, "function")
	env.Define("Logika", logikaFn, false, "function")
	env.Define("logika", logikaFn, false, "function")
	env.Define("Boolean", logikaFn, false, "function")
	env.Define("boolean", logikaFn, false, "function")

	bitaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) == 0 {
				return &Bytes{Elements: []byte{}}
			}
			if len(args) == 2 {
				strObj, ok1 := args[0].(*String)
				encObj, ok2 := args[1].(*String)
				if !ok1 || !ok2 {
					return &Error{Message: fmt.Sprintf("konstruktor bita dengan encoding mengharapkan teks dan format encoding, menerima %s dan %s", args[0].Type(), args[1].Type())}
				}
				enc := strings.ToLower(strings.TrimSpace(encObj.Value))
				switch enc {
				case "base64":
					data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(strObj.Value))
					if err != nil {
						return &Error{Message: fmt.Sprintf("Gagal mendekode data Base64: %v", err)}
					}
					return &Bytes{Elements: data}
				case "heks", "hex":
					cleanStr := strings.TrimSpace(strObj.Value)
					cleanStr = strings.TrimPrefix(cleanStr, "0x")
					cleanStr = strings.TrimPrefix(cleanStr, "0X")
					data, err := hex.DecodeString(cleanStr)
					if err != nil {
						return &Error{Message: fmt.Sprintf("Gagal mendekode teks heksadesimal: %v", err)}
					}
					return &Bytes{Elements: data}
				default:
					return &Error{Message: fmt.Sprintf("encoding '%s' tidak dikenal, gunakan 'base64' atau 'heks'/'hex'", enc)}
				}
			}
			switch a := args[0].(type) {
			case *Bytes:
				cp := make([]byte, len(a.Elements))
				copy(cp, a.Elements)
				return &Bytes{Elements: cp}
			case *String:
				return &Bytes{Elements: []byte(a.Value)}
			case *Integer:
				if a.Value < 0 {
					return &Error{Message: "ukuran alokasi bita tidak boleh negatif"}
				}
				if a.Value > 100_000_000 {
					return &Error{Message: "ukuran alokasi bita melebihi batas aman (100MB)"}
				}
				return &Bytes{Elements: make([]byte, a.Value)}
			case *Array:
				elems := make([]byte, len(a.Elements))
				for i, el := range a.Elements {
					if iVal, _, isInt, ok := CoerceToNumeric(el); ok && isInt {
						if iVal < 0 || iVal > 255 {
							return &Error{Message: fmt.Sprintf("nilai elemen bita harus di antara 0 dan 255, menerima %d", iVal)}
						}
						elems[i] = byte(iVal)
					} else {
						return &Error{Message: fmt.Sprintf("elemen array untuk bita harus berupa angka (0..255), menerima %s", el.Type())}
					}
				}
				return &Bytes{Elements: elems}
			default:
				return &Error{Message: fmt.Sprintf("tidak dapat membuat bita dari tipe %s", a.Type())}
			}
		},
	}

	dariBase64Fn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				msg := "bita.dari_base64 mengharapkan 1 argumen teks Base64"
				hint := "Contoh: bita.dari_base64(\"SGFsbw==\")"
				if env != nil && env.IsUniversal {
					msg = "bytes.from_base64 expects 1 Base64 string argument"
					hint = "Example: bytes.from_base64(\"SGFsbw==\")"
				}
				return newError(nil, msg, hint)
			}
			strObj, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen dari_base64 harus berupa teks, menerima %s", args[0].Type()), "")
			}
			data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(strObj.Value))
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mendekode data Base64: %v", err), "Pastikan teks merupakan format Base64 yang valid.")
			}
			return &Bytes{Elements: data}
		},
	}

	dariHeksFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				msg := "bita.dari_heks mengharapkan 1 argumen teks heksadesimal"
				hint := "Contoh: bita.dari_heks(\"48656c6c6f\")"
				if env != nil && env.IsUniversal {
					msg = "bytes.from_hex expects 1 hex string argument"
					hint = "Example: bytes.from_hex(\"48656c6c6f\")"
				}
				return newError(nil, msg, hint)
			}
			strObj, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("argumen dari_heks harus berupa teks, menerima %s", args[0].Type()), "")
			}
			cleanStr := strings.TrimSpace(strObj.Value)
			cleanStr = strings.TrimPrefix(cleanStr, "0x")
			cleanStr = strings.TrimPrefix(cleanStr, "0X")
			data, err := hex.DecodeString(cleanStr)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal mendekode teks heksadesimal: %v", err), "Pastikan teks hanya berisi karakter heksadesimal (0-9, a-f, A-F).")
			}
			return &Bytes{Elements: data}
		},
	}

	bitaFn.Properties = map[string]Object{
		"dari_base64": dariBase64Fn,
		"from_base64": dariBase64Fn,
		"dari_heks":   dariHeksFn,
		"dari_hex":    dariHeksFn,
		"from_hex":    dariHeksFn,
		"from_heks":   dariHeksFn,
	}

	env.Define("Bita", bitaFn, false, "function")
	env.Define("bita", bitaFn, false, "function")
	env.Define("Bytes", bitaFn, false, "function")
	env.Define("bytes", bitaFn, false, "function")

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
			tickets := make([]*Ticket, len(arr.Elements))
			for i, el := range arr.Elements {
				ticket, ok := el.(*Ticket)
				if !ok {
					return &Error{Message: fmt.Sprintf("Elemen ke-%d bukan merupakan Tiket", i)}
				}
				tickets[i] = ticket
			}

			results := make([]Object, len(tickets))
			var firstErr Object
			var mu sync.Mutex
			var wg sync.WaitGroup

			for i, ticket := range tickets {
				wg.Add(1)
				idx := i
				t := ticket
				go func() {
					defer wg.Done()
					res := t.Await(0, nil)
					mu.Lock()
					defer mu.Unlock()
					if isError(res) && firstErr == nil {
						firstErr = res
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

			tickets := make([]*Ticket, len(arr.Elements))
			for i, el := range arr.Elements {
				ticket, ok := el.(*Ticket)
				if !ok {
					return &Error{Message: "Semua elemen harus berupa Tiket"}
				}
				tickets[i] = ticket
			}

			resChan := make(chan Object, len(tickets))
			for _, ticket := range tickets {
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

	// Dynamic Import: impor_dinamis(path) / import_dynamic(path)
	imporDinamisFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				if env.IsUniversal {
					return &Error{Message: "import_dynamic() requires 1 string argument (module path or package name)"}
				}
				return &Error{Message: "impor_dinamis() membutuhkan 1 argumen string (path modul atau nama paket)"}
			}

			pathStr, ok := args[0].(*String)
			if !ok {
				if env.IsUniversal {
					return &Error{Message: fmt.Sprintf("import_dynamic() argument must be a string, got %s", args[0].Type())}
				}
				return &Error{Message: fmt.Sprintf("Argumen impor_dinamis() harus berupa teks (string), menerima %s", args[0].Type())}
			}

			importSource := pathStr.Value

			if env.IsUniversal {
				switch strings.ToLower(strings.TrimSpace(importSource)) {
				case "berkas":
					return &Error{Message: "Indonesian module 'berkas' is not allowed in Universal Mode. Import 'fs' instead."}
				case "matematika":
					return &Error{Message: "Indonesian module 'matematika' is not allowed in Universal Mode. Import 'math' instead."}
				case "teks":
					return &Error{Message: "Indonesian module 'teks' is not allowed in Universal Mode. Import 'strings' instead."}
				case "waktu":
					return &Error{Message: "Indonesian module 'waktu' is not allowed in Universal Mode. Import 'time' instead."}
				case "sistem":
					return &Error{Message: "Indonesian module 'sistem' is not allowed in Universal Mode. Import 'os' instead."}
				case "kripto":
					return &Error{Message: "Indonesian module 'kripto' is not allowed in Universal Mode. Import 'crypto' instead."}
				case "arsip", "kompresi":
					return &Error{Message: fmt.Sprintf("Indonesian module '%s' is not allowed in Universal Mode. Import 'archive' or 'compression' instead.", importSource)}
				}
			}

			targetAbsPath, err := resolveModulePath(env.FilePath, importSource)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Gagal menemukan modul '%s': %s", importSource, err.Error())}
			}

			targetEnv, err := loadModule(targetAbsPath, env)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Gagal memuat modul '%s': %s", importSource, err.Error())}
			}

			callerFile := env.FilePath
			callerRoot := env.ModuleRoot
			targetFile := targetEnv.FilePath
			targetRoot := targetEnv.ModuleRoot

			mapPairs := make(map[string]Object)
			var mapOrder []string
			for name, entry := range targetEnv.AllEntries() {
				if err := checkVisibilityAccess(callerFile, callerRoot, targetFile, targetRoot, name, entry.Visibility); err == nil {
					mapPairs[name] = entry.Current()
					mapOrder = append(mapOrder, name)
				}
			}
			sort.Strings(mapOrder)

			return &Map{Pairs: mapPairs, Order: mapOrder}
		},
	}
	env.Define("impor_dinamis", imporDinamisFn, false, "function")
	env.Define("ambil_dinamis", imporDinamisFn, false, "function")
	env.Define("import_dynamic", imporDinamisFn, false, "function")
	env.Define("import_module", imporDinamisFn, false, "function")

	// Register testing assertions & callbacks
	registerTestingBuiltins(env)

	// Register Standard Library global objects (berkas, fs, matematika, math, teks, strings, waktu, time, json)
	RegisterStdlibGlobals(env)
}

func registerTestingBuiltins(env *Environment) {
	// 1. harus / assert / must
	harusFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus() membutuhkan minimal 1 argumen (kondisi)"}
			}
			msg := "Kondisi tidak terpenuhi (diharapkan bernilai benar)"
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					msg = s.Value
				} else {
					msg = args[1].Inspect()
				}
			}
			if !CoerceToBoolean(args[0]) {
				return &AssertionError{
					Message:  msg,
					Expected: "benar (true)",
					Actual:   args[0].Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus", harusFn, false, "function")
	env.Define("assert", harusFn, false, "function")
	env.Define("must", harusFn, false, "function")

	// 2. harus_sama / assert_eq / assert_equal
	harusSamaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return &AssertionError{Message: "harus_sama() membutuhkan minimal 2 argumen: (aktual, ekspektasi)"}
			}
			actual := args[0]
			expected := args[1]
			msg := "Nilai aktual tidak sama dengan ekspektasi"
			if len(args) >= 3 {
				if s, ok := args[2].(*String); ok {
					msg = s.Value
				} else {
					msg = args[2].Inspect()
				}
			}
			if !ObjectsEqual(actual, expected) {
				return &AssertionError{
					Message:  msg,
					Expected: expected.Inspect(),
					Actual:   actual.Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_sama", harusSamaFn, false, "function")
	env.Define("assert_eq", harusSamaFn, false, "function")
	env.Define("assert_equal", harusSamaFn, false, "function")

	// 3. harus_beda / assert_ne / assert_not_equal
	harusBedaFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return &AssertionError{Message: "harus_beda() membutuhkan minimal 2 argumen: (aktual, ekspektasi)"}
			}
			actual := args[0]
			expected := args[1]
			msg := "Nilai diharapkan berbeda, tetapi bernilai sama"
			if len(args) >= 3 {
				if s, ok := args[2].(*String); ok {
					msg = s.Value
				} else {
					msg = args[2].Inspect()
				}
			}
			if ObjectsEqual(actual, expected) {
				return &AssertionError{
					Message:  msg,
					Expected: "berbeda dari " + expected.Inspect(),
					Actual:   actual.Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_beda", harusBedaFn, false, "function")
	env.Define("assert_ne", harusBedaFn, false, "function")
	env.Define("assert_not_equal", harusBedaFn, false, "function")

	// 4. harus_benar / assert_true
	harusBenarFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus_benar() membutuhkan 1 argumen"}
			}
			msg := "Diharapkan bernilai benar (true)"
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					msg = s.Value
				}
			}
			if b, ok := args[0].(*Boolean); !ok || !b.Value {
				return &AssertionError{
					Message:  msg,
					Expected: "benar (true)",
					Actual:   args[0].Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_benar", harusBenarFn, false, "function")
	env.Define("assert_true", harusBenarFn, false, "function")

	// 5. harus_salah / assert_false
	harusSalahFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus_salah() membutuhkan 1 argumen"}
			}
			msg := "Diharapkan bernilai salah (false)"
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					msg = s.Value
				}
			}
			if b, ok := args[0].(*Boolean); !ok || b.Value {
				return &AssertionError{
					Message:  msg,
					Expected: "salah (false)",
					Actual:   args[0].Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_salah", harusSalahFn, false, "function")
	env.Define("assert_false", harusSalahFn, false, "function")

	// 6. harus_kosong / assert_null
	harusKosongFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus_kosong() membutuhkan 1 argumen"}
			}
			msg := "Diharapkan bernilai kosong (null)"
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					msg = s.Value
				}
			}
			if args[0] != nil && args[0].Type() != NULL_OBJ {
				return &AssertionError{
					Message:  msg,
					Expected: "kosong (null)",
					Actual:   args[0].Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_kosong", harusKosongFn, false, "function")
	env.Define("assert_null", harusKosongFn, false, "function")

	// 7. harus_tidak_kosong / assert_not_null
	harusTidakKosongFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus_tidak_kosong() membutuhkan 1 argumen"}
			}
			msg := "Diharapkan tidak bernilai kosong (bukan null)"
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					msg = s.Value
				}
			}
			if args[0] == nil || args[0].Type() == NULL_OBJ {
				return &AssertionError{
					Message:  msg,
					Expected: "bukan kosong",
					Actual:   "kosong (null)",
				}
			}
			return NULL
		},
	}
	env.Define("harus_tidak_kosong", harusTidakKosongFn, false, "function")
	env.Define("assert_not_null", harusTidakKosongFn, false, "function")

	// 8. harus_berisi / assert_contains
	harusBerisiFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 2 {
				return &AssertionError{Message: "harus_berisi() membutuhkan minimal 2 argumen: (kontainer, elemen)"}
			}
			container := args[0]
			item := args[1]
			msg := "Kontainer tidak memuat elemen yang diharapkan"
			if len(args) >= 3 {
				if s, ok := args[2].(*String); ok {
					msg = s.Value
				}
			}
			found := false
			switch c := container.(type) {
			case *String:
				if itemStr, ok := item.(*String); ok {
					found = strings.Contains(c.Value, itemStr.Value)
				} else {
					found = strings.Contains(c.Value, item.Inspect())
				}
			case *Array:
				for _, el := range c.Elements {
					if ObjectsEqual(el, item) {
						found = true
						break
					}
				}
			case *Map:
				if itemStr, ok := item.(*String); ok {
					_, found = c.Pairs[itemStr.Value]
				}
			case *Bytes:
				if itemInt, ok := item.(*Integer); ok {
					targetByte := byte(itemInt.Value)
					for _, b := range c.Elements {
						if b == targetByte {
							found = true
							break
						}
					}
				}
			default:
				return &AssertionError{
					Message:  fmt.Sprintf("harus_berisi tidak mendukung tipe kontainer '%s'", container.Type()),
					Expected: "teks, larik, atau peta",
					Actual:   string(container.Type()),
				}
			}
			if !found {
				return &AssertionError{
					Message:  msg,
					Expected: "memuat " + item.Inspect(),
					Actual:   container.Inspect(),
				}
			}
			return NULL
		},
	}
	env.Define("harus_berisi", harusBerisiFn, false, "function")
	env.Define("assert_contains", harusBerisiFn, false, "function")

	// 9. harus_galat / assert_throws
	harusGalatFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "harus_galat() membutuhkan minimal 1 argumen: fungsi/closure yang diuji"}
			}
			fnObj := args[0]
			var expectedErrMsg string
			if len(args) >= 2 {
				if s, ok := args[1].(*String); ok {
					expectedErrMsg = s.Value
				}
			}

			res := applyFunction(fnObj, []Object{}, nil)
			if !isError(res) {
				return &AssertionError{
					Message:  "Diharapkan menghasilkan galat (error), namun fungsi berjalan sukses tanpa galat",
					Expected: "galat/error",
					Actual:   res.Inspect(),
				}
			}

			if expectedErrMsg != "" {
				actualMsg := res.Inspect()
				if errObj, ok := res.(*Error); ok {
					actualMsg = errObj.Message
				} else if aErr, ok := res.(*AssertionError); ok {
					actualMsg = aErr.Message
				}
				if !strings.Contains(strings.ToLower(actualMsg), strings.ToLower(expectedErrMsg)) {
					return &AssertionError{
						Message:  fmt.Sprintf("Galat terjadi tetapi pesannya tidak memuat '%s'", expectedErrMsg),
						Expected: "pesan memuat: " + expectedErrMsg,
						Actual:   actualMsg,
					}
				}
			}

			return NULL
		},
	}
	env.Define("harus_galat", harusGalatFn, false, "function")
	env.Define("assert_throws", harusGalatFn, false, "function")

	// 10. uji / test callback
	ujiCallbackFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &AssertionError{Message: "uji() membutuhkan minimal 1 argumen (fungsi pengujian)"}
			}
			testName := "unnamed test"
			var fnObj Object
			if len(args) >= 2 {
				if s, ok := args[0].(*String); ok {
					testName = s.Value
				} else {
					testName = args[0].Inspect()
				}
				fnObj = args[1]
			} else {
				fnObj = args[0]
			}

			tc := env.GetTestContext()
			if tc != nil && tc.Filter != "" {
				if !strings.Contains(strings.ToLower(testName), strings.ToLower(tc.Filter)) {
					return NULL
				}
			}

			start := time.Now()
			res := applyFunction(fnObj, []Object{}, nil)
			dur := time.Since(start)

			if isError(res) {
				errMsg := res.Inspect()
				expected := ""
				actual := ""
				var pos *uast.Position
				if aErr, ok := res.(*AssertionError); ok {
					errMsg = aErr.Message
					expected = aErr.Expected
					actual = aErr.Actual
					pos = aErr.Pos
				} else if errObj, ok := res.(*Error); ok {
					errMsg = errObj.Message
					pos = errObj.Pos
				}

				if tc != nil {
					tc.RecordResult(TestCaseResult{
						Name:     testName,
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
					Name:     testName,
					Passed:   true,
					Duration: dur,
				})
			}
			return NULL
		},
	}
	env.Define("uji", ujiCallbackFn, false, "function")
	env.Define("test", ujiCallbackFn, false, "function")
}

func deepClone(obj Object) Object {
	visited := make(map[any]Object)
	return deepCloneWithVisited(obj, visited, 0)
}

func deepCloneWithVisited(obj Object, visited map[any]Object, depth int) Object {
	if obj == nil {
		return NULL
	}
	if depth > 100 {
		return obj
	}
	if cloned, exists := visited[obj]; exists {
		return cloned
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
	case *Bytes:
		cloned := make([]byte, len(o.Elements))
		copy(cloned, o.Elements)
		return &Bytes{Elements: cloned}
	case *Null:
		return &Null{}
	case *Array:
		newArr := &Array{Elements: make([]Object, len(o.Elements))}
		visited[o] = newArr
		for i, el := range o.Elements {
			newArr.Elements[i] = deepCloneWithVisited(el, visited, depth+1)
		}
		return newArr
	case *Map:
		newMap := &Map{
			Pairs: make(map[string]Object, len(o.Pairs)),
		}
		if o.Order != nil {
			newMap.Order = make([]string, len(o.Order))
			copy(newMap.Order, o.Order)
		}
		visited[o] = newMap
		for k, v := range o.Pairs {
			newMap.Pairs[k] = deepCloneWithVisited(v, visited, depth+1)
		}
		return newMap
	case *StructInstance:
		newFields := make(map[string]Object, len(o.Fields))
		cloneInst := &StructInstance{
			Def:    o.Def,
			Fields: newFields,
			Env:    NewEnclosedEnvironment(o.Def.Env),
		}
		visited[o] = cloneInst
		for k, v := range o.Fields {
			newFields[k] = deepCloneWithVisited(v, visited, depth+1)
		}
		cloneInst.Env.Define("diri", cloneInst, false, o.Def.Name)
		cloneInst.Env.Define("this", cloneInst, false, o.Def.Name)
		cloneInst.Env.Define("self", cloneInst, false, o.Def.Name)
		return cloneInst
	default:
		return o
	}
}

