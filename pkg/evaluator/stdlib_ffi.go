//go:build !wasm && !js

package evaluator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// FFI_POINTER_TYPE adalah tipe objek penunjuk memori mentah untuk FFI.
const FFI_POINTER_TYPE ObjectType = "FFI_POINTER"

// FFIPointer merepresentasikan alamat memori (pointer) mentah C.
type FFIPointer struct {
	Address uintptr
}

func (p *FFIPointer) Type() ObjectType { return FFI_POINTER_TYPE }
func (p *FFIPointer) Inspect() string {
	if p.Address == 0 {
		return "Pointer(NULL)"
	}
	return fmt.Sprintf("Pointer(0x%x)", p.Address)
}

var (
	ffiAllocMu      sync.RWMutex
	ffiAllocBuffers = make(map[uintptr][]byte)
)

func normalizeFFIType(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	switch lower {
	case "int", "integer", "angka":
		return "int"
	case "float", "double", "pecahan":
		return "float"
	case "string", "teks", "str":
		return "string"
	case "pointer", "penunjuk", "ptr":
		return "pointer"
	case "bool", "boolean", "kebenaran", "logika":
		return "bool"
	case "void", "kosong", "null":
		return "void"
	default:
		return lower
	}
}

func readCString(addr uintptr) string {
	if addr == 0 {
		return ""
	}
	var bytes []byte
	const maxLen = 10 * 1024 * 1024 // batas aman 10MB
	for i := 0; i < maxLen; i++ {
		b := *(*byte)(unsafe.Pointer(addr + uintptr(i)))
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
	}
	return string(bytes)
}

func writeCString(addr uintptr, s string) int {
	if addr == 0 {
		return 0
	}
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		*(*byte)(unsafe.Pointer(addr + uintptr(i))) = b[i]
	}
	*(*byte)(unsafe.Pointer(addr + uintptr(len(b)))) = 0
	return len(b)
}

func readCBytes(addr uintptr, length int) []byte {
	if addr == 0 || length <= 0 {
		return []byte{}
	}
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = *(*byte)(unsafe.Pointer(addr + uintptr(i)))
	}
	return res
}

func writeCBytes(addr uintptr, b []byte) int {
	if addr == 0 || len(b) == 0 {
		return 0
	}
	for i := 0; i < len(b); i++ {
		*(*byte)(unsafe.Pointer(addr + uintptr(i))) = b[i]
	}
	return len(b)
}

func createFFIDynamicCaller(proc uintptr, inKinds []string, outKind string, name, libPath string) (*Builtin, error) {
	var inTypes []reflect.Type
	for _, k := range inKinds {
		switch k {
		case "int":
			inTypes = append(inTypes, reflect.TypeOf(int64(0)))
		case "float":
			inTypes = append(inTypes, reflect.TypeOf(float64(0)))
		case "string", "pointer":
			inTypes = append(inTypes, reflect.TypeOf(uintptr(0)))
		case "bool":
			inTypes = append(inTypes, reflect.TypeOf(bool(false)))
		default:
			return nil, fmt.Errorf("tipe parameter C '%s' tidak didukung", k)
		}
	}

	var outTypes []reflect.Type
	switch outKind {
	case "int":
		outTypes = append(outTypes, reflect.TypeOf(int64(0)))
	case "float":
		outTypes = append(outTypes, reflect.TypeOf(float64(0)))
	case "string", "pointer":
		outTypes = append(outTypes, reflect.TypeOf(uintptr(0)))
	case "bool":
		outTypes = append(outTypes, reflect.TypeOf(bool(false)))
	case "void":
		// kosong
	default:
		return nil, fmt.Errorf("tipe kembalian C '%s' tidak didukung", outKind)
	}

	fnType := reflect.FuncOf(inTypes, outTypes, false)
	fnVal := reflect.New(fnType)
	purego.RegisterFunc(fnVal.Interface(), proc)
	callable := fnVal.Elem()

	builtin := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != len(inKinds) {
				return newError(nil, fmt.Sprintf("Fungsi FFI '%s' membutuhkan %d argumen, tetapi diberikan %d", name, len(inKinds), len(args)), "")
			}

			callArgs := make([]reflect.Value, len(args))
			var pinnedBuffers [][]byte

			for i, arg := range args {
				expectedKind := inKinds[i]
				switch expectedKind {
				case "int":
					var val int64
					if n, ok := arg.(*Integer); ok {
						val = n.Value
					} else if f, ok := arg.(*Float); ok {
						val = int64(f.Value)
					} else {
						return newError(nil, fmt.Sprintf("Argumen %d fungsi FFI '%s' harus berupa Angka (int), menerima %s", i+1, name, arg.Type()), "")
					}
					callArgs[i] = reflect.ValueOf(val)

				case "float":
					var val float64
					if f, ok := arg.(*Float); ok {
						val = f.Value
					} else if n, ok := arg.(*Integer); ok {
						val = float64(n.Value)
					} else {
						return newError(nil, fmt.Sprintf("Argumen %d fungsi FFI '%s' harus berupa Pecahan/Angka (float), menerima %s", i+1, name, arg.Type()), "")
					}
					callArgs[i] = reflect.ValueOf(val)

				case "string":
					if s, ok := arg.(*String); ok {
						buf := append([]byte(s.Value), 0)
						pinnedBuffers = append(pinnedBuffers, buf)
						callArgs[i] = reflect.ValueOf(uintptr(unsafe.Pointer(&buf[0])))
					} else if arg == NULL {
						callArgs[i] = reflect.ValueOf(uintptr(0))
					} else if p, ok := arg.(*FFIPointer); ok {
						callArgs[i] = reflect.ValueOf(p.Address)
					} else {
						return newError(nil, fmt.Sprintf("Argumen %d fungsi FFI '%s' harus berupa Teks (string), menerima %s", i+1, name, arg.Type()), "")
					}

				case "pointer":
					if p, ok := arg.(*FFIPointer); ok {
						callArgs[i] = reflect.ValueOf(p.Address)
					} else if arg == NULL {
						callArgs[i] = reflect.ValueOf(uintptr(0))
					} else if n, ok := arg.(*Integer); ok {
						callArgs[i] = reflect.ValueOf(uintptr(n.Value))
					} else if s, ok := arg.(*String); ok {
						buf := append([]byte(s.Value), 0)
						pinnedBuffers = append(pinnedBuffers, buf)
						callArgs[i] = reflect.ValueOf(uintptr(unsafe.Pointer(&buf[0])))
					} else if b, ok := arg.(*Bytes); ok {
						if len(b.Elements) > 0 {
							callArgs[i] = reflect.ValueOf(uintptr(unsafe.Pointer(&b.Elements[0])))
						} else {
							callArgs[i] = reflect.ValueOf(uintptr(0))
						}
					} else {
						return newError(nil, fmt.Sprintf("Argumen %d fungsi FFI '%s' harus berupa ffi.Pointer, menerima %s", i+1, name, arg.Type()), "")
					}

				case "bool":
					callArgs[i] = reflect.ValueOf(CoerceToBoolean(arg))
				}
			}

			// Panggil fungsi C secara native via ABI
			results := callable.Call(callArgs)
			_ = pinnedBuffers // pastikan memori string/buffer tetap hidup selama panggilan

			if len(results) == 0 || outKind == "void" {
				return NULL
			}

			resVal := results[0]
			switch outKind {
			case "int":
				return &Integer{Value: resVal.Int()}
			case "float":
				return &Float{Value: resVal.Float()}
			case "bool":
				return nativeBoolToBoolean(resVal.Bool())
			case "pointer":
				return &FFIPointer{Address: resVal.Interface().(uintptr)}
			case "string":
				addr := resVal.Interface().(uintptr)
				if addr == 0 {
					return NULL
				}
				return &String{Value: readCString(addr)}
			default:
				return NULL
			}
		},
	}

	return builtin, nil
}

func initStdlibFFI(env *Environment) {
	// -------------------------------------------------------------
	// 1. Tipe Data FFI (Konstanta Tipe)
	// -------------------------------------------------------------
	definePublicConstant(env, "Int", &String{Value: "int"})
	definePublicConstant(env, "Integer", &String{Value: "int"})
	definePublicConstant(env, "Angka", &String{Value: "int"})

	definePublicConstant(env, "Float", &String{Value: "float"})
	definePublicConstant(env, "Double", &String{Value: "float"})
	definePublicConstant(env, "Pecahan", &String{Value: "float"})

	definePublicConstant(env, "String", &String{Value: "string"})
	definePublicConstant(env, "Teks", &String{Value: "string"})

	// Sesuai permintaan pengguna: ffi.Pointer sebagai nama utama, ffi.Penunjuk sebagai alias
	definePublicConstant(env, "Pointer", &String{Value: "pointer"})
	definePublicConstant(env, "Penunjuk", &String{Value: "pointer"})

	definePublicConstant(env, "Void", &String{Value: "void"})
	definePublicConstant(env, "Kosong", &String{Value: "void"})

	definePublicConstant(env, "Bool", &String{Value: "bool"})
	definePublicConstant(env, "Kebenaran", &String{Value: "bool"})

	definePublicConstant(env, "NULL", &FFIPointer{Address: 0})
	definePublicConstant(env, "NIL", &FFIPointer{Address: 0})

	// -------------------------------------------------------------
	// 2. Membuka Pustaka Dinamis (open / buka)
	// -------------------------------------------------------------
	openFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.open / ffi.buka membutuhkan 1 argumen: path/nama berkas pustaka (.so / .dll / .dylib)", "Contoh: ffi.open(\"libc.so.6\")")
			}
			pathStr, ok := args[0].(*String)
			if !ok {
				return newError(nil, fmt.Sprintf("path pustaka harus teks, menerima %s", args[0].Type()), "")
			}

			handle, err := ffiOpen(pathStr.Value)
			if err != nil {
				return newError(nil, fmt.Sprintf("Gagal memuat pustaka dinamis '%s': %v", pathStr.Value, err), "Pastikan nama berkas pustaka (.so / .dll / .dylib) sudah benar dan terpasang di sistem.")
			}

			libMap := &Map{
				Pairs: make(map[string]Object),
				Order: []string{},
			}

			libMap.Pairs["path"] = pathStr
			libMap.Pairs["jalur"] = pathStr
			libMap.Pairs["handle"] = &FFIPointer{Address: handle}
			libMap.Pairs["pegangan"] = &FFIPointer{Address: handle}

			// lib.func / lib.fungsi(nama, [tipe_parameter...], tipe_kembalian)
			funcFn := &Builtin{
				Fn: func(fnArgs ...Object) Object {
					if len(fnArgs) < 2 || len(fnArgs) > 3 {
						return newError(nil, "metode fungsi() membutuhkan 2 atau 3 argumen: nama, [tipe_param...], tipe_kembalian (opsional, default: void)", "Contoh: lib.fungsi(\"puts\", [ffi.String], ffi.Int)")
					}

					nameStr, ok := fnArgs[0].(*String)
					if !ok {
						return newError(nil, fmt.Sprintf("nama fungsi C harus teks, menerima %s", fnArgs[0].Type()), "")
					}

					var inKinds []string
					if paramList, ok := fnArgs[1].(*Array); ok {
						for _, p := range paramList.Elements {
							if s, ok := p.(*String); ok {
								inKinds = append(inKinds, normalizeFFIType(s.Value))
							} else {
								return newError(nil, "daftar tipe parameter harus berupa teks tipe FFI (misal: ffi.String, ffi.Int, ffi.Pointer)", "")
							}
						}
					} else {
						return newError(nil, "parameter kedua harus berupa Array daftar tipe parameter (misal: [ffi.String, ffi.Int])", "")
					}

					outKind := "void"
					if len(fnArgs) == 3 {
						if retStr, ok := fnArgs[2].(*String); ok {
							outKind = normalizeFFIType(retStr.Value)
						} else {
							return newError(nil, "tipe kembalian harus berupa teks tipe FFI (misal: ffi.Int, ffi.Void, ffi.Pointer)", "")
						}
					}

					proc, err := ffiSymbol(handle, nameStr.Value)
					if err != nil {
						return newError(nil, fmt.Sprintf("Simbol fungsi '%s' tidak ditemukan di pustaka '%s': %v", nameStr.Value, pathStr.Value, err), "")
					}

					builtinCallable, err := createFFIDynamicCaller(proc, inKinds, outKind, nameStr.Value, pathStr.Value)
					if err != nil {
						return newError(nil, fmt.Sprintf("Gagal menyiapkan jembatan ABI untuk '%s': %v", nameStr.Value, err), "")
					}

					return builtinCallable
				},
			}
			libMap.Pairs["func"] = funcFn
			libMap.Pairs["fungsi"] = funcFn

			// lib.symbol / lib.simbol(nama)
			symbolFn := &Builtin{
				Fn: func(symArgs ...Object) Object {
					if len(symArgs) != 1 {
						return newError(nil, "metode simbol() membutuhkan 1 argumen: nama simbol C", "")
					}
					nameStr, ok := symArgs[0].(*String)
					if !ok {
						return newError(nil, "nama simbol harus berupa teks", "")
					}
					proc, err := ffiSymbol(handle, nameStr.Value)
					if err != nil {
						return newError(nil, fmt.Sprintf("Simbol '%s' tidak ditemukan: %v", nameStr.Value, err), "")
					}
					return &FFIPointer{Address: proc}
				},
			}
			libMap.Pairs["symbol"] = symbolFn
			libMap.Pairs["simbol"] = symbolFn

			// lib.close / lib.tutup()
			closeFn := &Builtin{
				Fn: func(closeArgs ...Object) Object {
					_ = ffiClose(handle)
					return NULL
				},
			}
			libMap.Pairs["close"] = closeFn
			libMap.Pairs["tutup"] = closeFn

			return libMap
		},
	}

	definePublicFunction(env, "open", openFn)
	definePublicFunction(env, "buka", openFn)

	// -------------------------------------------------------------
	// 3. Utilitas Memori & Pointer
	// -------------------------------------------------------------

	// ffi.alloc / ffi.alokasi(ukuran)
	allocFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.alloc membutuhkan 1 argumen: ukuran alokasi byte", "Contoh: ffi.alloc(1024)")
			}
			size := int64(0)
			if n, ok := args[0].(*Integer); ok {
				size = n.Value
			} else if f, ok := args[0].(*Float); ok {
				size = int64(f.Value)
			} else {
				return newError(nil, "ukuran alokasi harus berupa Angka", "")
			}
			if size <= 0 {
				return newError(nil, "ukuran alokasi harus lebih dari 0 byte", "")
			}

			buf := make([]byte, size)
			addr := uintptr(unsafe.Pointer(&buf[0]))

			ffiAllocMu.Lock()
			ffiAllocBuffers[addr] = buf
			ffiAllocMu.Unlock()

			return &FFIPointer{Address: addr}
		},
	}
	definePublicFunction(env, "alloc", allocFn)
	definePublicFunction(env, "alokasi", allocFn)

	// ffi.free / ffi.bebas(ptr)
	freeFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.free membutuhkan 1 argumen: ffi.Pointer", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				ffiAllocMu.Lock()
				delete(ffiAllocBuffers, p.Address)
				ffiAllocMu.Unlock()
			}
			return NULL
		},
	}
	definePublicFunction(env, "free", freeFn)
	definePublicFunction(env, "bebas", freeFn)

	// ffi.pointer / ffi.to_pointer(addr)
	toPointerFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.pointer membutuhkan 1 argumen: alamat numerik", "")
			}
			if n, ok := args[0].(*Integer); ok {
				return &FFIPointer{Address: uintptr(n.Value)}
			} else if f, ok := args[0].(*Float); ok {
				return &FFIPointer{Address: uintptr(int64(f.Value))}
			} else if p, ok := args[0].(*FFIPointer); ok {
				return p
			}
			return newError(nil, "argumen harus berupa Angka alamat memori", "")
		},
	}
	definePublicFunction(env, "pointer", toPointerFn)
	definePublicFunction(env, "to_pointer", toPointerFn)

	// ffi.address / ffi.alamat(ptr)
	addressFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.address membutuhkan 1 argumen: ffi.Pointer", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				return &Integer{Value: int64(p.Address)}
			}
			return newError(nil, "argumen harus berupa ffi.Pointer", "")
		},
	}
	definePublicFunction(env, "address", addressFn)
	definePublicFunction(env, "alamat", addressFn)

	// ffi.is_null / ffi.apakah_kosong(ptr)
	isNullFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.is_null membutuhkan 1 argumen: ffi.Pointer", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				return nativeBoolToBoolean(p.Address == 0)
			}
			if args[0] == NULL {
				return nativeBoolToBoolean(true)
			}
			return nativeBoolToBoolean(false)
		},
	}
	definePublicFunction(env, "is_null", isNullFn)
	definePublicFunction(env, "apakah_kosong", isNullFn)

	// ffi.read_string / ffi.baca_teks(ptr)
	readStringFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.read_string membutuhkan 1 argumen: ffi.Pointer", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				return &String{Value: readCString(p.Address)}
			}
			return newError(nil, "argumen harus berupa ffi.Pointer", "")
		},
	}
	definePublicFunction(env, "read_string", readStringFn)
	definePublicFunction(env, "baca_teks", readStringFn)

	// ffi.write_string / ffi.tulis_teks(ptr, teks)
	writeStringFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.write_string membutuhkan 2 argumen: ffi.Pointer dan Teks", "")
			}
			p, okP := args[0].(*FFIPointer)
			s, okS := args[1].(*String)
			if !okP || !okS {
				return newError(nil, "argumen pertama harus ffi.Pointer, kedua harus Teks", "")
			}
			written := writeCString(p.Address, s.Value)
			return &Integer{Value: int64(written)}
		},
	}
	definePublicFunction(env, "write_string", writeStringFn)
	definePublicFunction(env, "tulis_teks", writeStringFn)

	// ffi.read_bytes / ffi.baca_bita(ptr, panjang)
	readBytesFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.read_bytes membutuhkan 2 argumen: ffi.Pointer dan panjang byte", "")
			}
			p, okP := args[0].(*FFIPointer)
			if !okP {
				return newError(nil, "argumen pertama harus berupa ffi.Pointer", "")
			}
			length := 0
			if n, ok := args[1].(*Integer); ok {
				length = int(n.Value)
			} else if f, ok := args[1].(*Float); ok {
				length = int(f.Value)
			} else {
				return newError(nil, "panjang byte harus berupa Angka", "")
			}
			return &Bytes{Elements: readCBytes(p.Address, length)}
		},
	}
	definePublicFunction(env, "read_bytes", readBytesFn)
	definePublicFunction(env, "baca_bita", readBytesFn)

	// ffi.write_bytes / ffi.tulis_bita(ptr, bita)
	writeBytesFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.write_bytes membutuhkan 2 argumen: ffi.Pointer dan Bita", "")
			}
			p, okP := args[0].(*FFIPointer)
			b, okB := args[1].(*Bytes)
			if !okP || !okB {
				return newError(nil, "argumen pertama harus ffi.Pointer, kedua harus Bita", "")
			}
			written := writeCBytes(p.Address, b.Elements)
			return &Integer{Value: int64(written)}
		},
	}
	definePublicFunction(env, "write_bytes", writeBytesFn)
	definePublicFunction(env, "tulis_bita", writeBytesFn)
}
