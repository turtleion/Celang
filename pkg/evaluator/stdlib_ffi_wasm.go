//go:build js && wasm

package evaluator

import (
	"fmt"
	"sync"
)

// FFI_POINTER_TYPE adalah tipe objek penunjuk memori mentah untuk FFI di WASM.
const FFI_POINTER_TYPE ObjectType = "FFI_POINTER"

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
	wasmAllocMu      sync.RWMutex
	wasmAllocBuffers = make(map[uintptr][]byte)
	wasmNextAddr     = uintptr(0x1000)
)

func initStdlibFFI(env *Environment) {
	// Konstanta tipe data
	definePublicConstant(env, "Int", &String{Value: "int"})
	definePublicConstant(env, "Integer", &String{Value: "int"})
	definePublicConstant(env, "Angka", &String{Value: "int"})

	definePublicConstant(env, "Float", &String{Value: "float"})
	definePublicConstant(env, "Double", &String{Value: "float"})
	definePublicConstant(env, "Pecahan", &String{Value: "float"})

	definePublicConstant(env, "String", &String{Value: "string"})
	definePublicConstant(env, "Teks", &String{Value: "string"})

	definePublicConstant(env, "Pointer", &String{Value: "pointer"})
	definePublicConstant(env, "Penunjuk", &String{Value: "pointer"})

	definePublicConstant(env, "Void", &String{Value: "void"})
	definePublicConstant(env, "Kosong", &String{Value: "void"})

	definePublicConstant(env, "Bool", &String{Value: "bool"})
	definePublicConstant(env, "Kebenaran", &String{Value: "bool"})

	definePublicConstant(env, "NULL", &FFIPointer{Address: 0})
	definePublicConstant(env, "NIL", &FFIPointer{Address: 0})

	// Open / buka pustaka biner
	openStubFn := &Builtin{
		Fn: func(args ...Object) Object {
			return newError(nil, "ffi.open / ffi.buka membutuhkan berkas biner native (.so / .dll / .dylib) yang hanya tersedia pada sistem operasi komputer asli (Linux, Windows, macOS). Jalankan Celang di terminal lokal untuk memuat pustaka biner tersebut.", "Contoh di terminal komputer: ./celang run program_ffi.ce")
		},
	}
	definePublicFunction(env, "open", openStubFn)
	definePublicFunction(env, "buka", openStubFn)

	// alloc / alokasi
	allocFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.alloc membutuhkan 1 argumen: ukuran alokasi byte", "")
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
				return newError(nil, "ukuran alokasi harus lebih dari 0", "")
			}

			wasmAllocMu.Lock()
			addr := wasmNextAddr
			wasmNextAddr += uintptr(size + 16)
			wasmAllocBuffers[addr] = make([]byte, size)
			wasmAllocMu.Unlock()

			return &FFIPointer{Address: addr}
		},
	}
	definePublicFunction(env, "alloc", allocFn)
	definePublicFunction(env, "alokasi", allocFn)

	// free / bebas
	freeFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.free membutuhkan 1 argumen: ffi.Pointer", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				wasmAllocMu.Lock()
				delete(wasmAllocBuffers, p.Address)
				wasmAllocMu.Unlock()
			}
			return NULL
		},
	}
	definePublicFunction(env, "free", freeFn)
	definePublicFunction(env, "bebas", freeFn)

	// pointer / to_pointer
	pointerFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.pointer membutuhkan 1 argumen", "")
			}
			if n, ok := args[0].(*Integer); ok {
				return &FFIPointer{Address: uintptr(n.Value)}
			} else if p, ok := args[0].(*FFIPointer); ok {
				return p
			}
			return newError(nil, "argumen harus berupa alamat numerik atau ffi.Pointer", "")
		},
	}
	definePublicFunction(env, "pointer", pointerFn)
	definePublicFunction(env, "to_pointer", pointerFn)

	// address / alamat
	addrFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.address membutuhkan 1 argumen", "")
			}
			if p, ok := args[0].(*FFIPointer); ok {
				return &Integer{Value: int64(p.Address)}
			}
			return newError(nil, "argumen harus berupa ffi.Pointer", "")
		},
	}
	definePublicFunction(env, "address", addrFn)
	definePublicFunction(env, "alamat", addrFn)

	// is_null / apakah_kosong
	isNullFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.is_null membutuhkan 1 argumen", "")
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

	// write_string / tulis_teks
	writeStrFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.write_string membutuhkan 2 argumen: ffi.Pointer dan Teks", "")
			}
			p, okP := args[0].(*FFIPointer)
			s, okS := args[1].(*String)
			if !okP || !okS {
				return newError(nil, "argumen pertama harus ffi.Pointer, kedua harus Teks", "")
			}
			wasmAllocMu.Lock()
			defer wasmAllocMu.Unlock()
			buf, ok := wasmAllocBuffers[p.Address]
			if !ok {
				return newError(nil, "penunjuk memori tidak valid atau sudah dibebaskan", "")
			}
			bytes := []byte(s.Value)
			n := copy(buf, bytes)
			if n < len(buf) {
				buf[n] = 0
			}
			return &Integer{Value: int64(n)}
		},
	}
	definePublicFunction(env, "write_string", writeStrFn)
	definePublicFunction(env, "tulis_teks", writeStrFn)

	// read_string / baca_teks
	readStrFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError(nil, "ffi.read_string membutuhkan 1 argumen: ffi.Pointer", "")
			}
			p, okP := args[0].(*FFIPointer)
			if !okP {
				return newError(nil, "argumen harus berupa ffi.Pointer", "")
			}
			wasmAllocMu.RLock()
			buf, ok := wasmAllocBuffers[p.Address]
			wasmAllocMu.RUnlock()
			if !ok {
				return &String{Value: ""}
			}
			// Read until null terminator or end
			end := 0
			for end < len(buf) && buf[end] != 0 {
				end++
			}
			return &String{Value: string(buf[:end])}
		},
	}
	definePublicFunction(env, "read_string", readStrFn)
	definePublicFunction(env, "baca_teks", readStrFn)

	// read_bytes / baca_bita
	readBytesFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.read_bytes membutuhkan 2 argumen", "")
			}
			p, okP := args[0].(*FFIPointer)
			if !okP {
				return newError(nil, "argumen pertama harus ffi.Pointer", "")
			}
			length := 0
			if n, ok := args[1].(*Integer); ok {
				length = int(n.Value)
			}
			wasmAllocMu.RLock()
			buf, ok := wasmAllocBuffers[p.Address]
			wasmAllocMu.RUnlock()
			if !ok || length <= 0 {
				return &Bytes{Elements: []byte{}}
			}
			if length > len(buf) {
				length = len(buf)
			}
			res := make([]byte, length)
			copy(res, buf[:length])
			return &Bytes{Elements: res}
		},
	}
	definePublicFunction(env, "read_bytes", readBytesFn)
	definePublicFunction(env, "baca_bita", readBytesFn)

	// write_bytes / tulis_bita
	writeBytesFn := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError(nil, "ffi.write_bytes membutuhkan 2 argumen", "")
			}
			p, okP := args[0].(*FFIPointer)
			b, okB := args[1].(*Bytes)
			if !okP || !okB {
				return newError(nil, "argumen harus ffi.Pointer dan Bita", "")
			}
			wasmAllocMu.Lock()
			defer wasmAllocMu.Unlock()
			buf, ok := wasmAllocBuffers[p.Address]
			if !ok {
				return newError(nil, "penunjuk memori tidak valid atau sudah dibebaskan", "")
			}
			n := copy(buf, b.Elements)
			return &Integer{Value: int64(n)}
		},
	}
	definePublicFunction(env, "write_bytes", writeBytesFn)
	definePublicFunction(env, "tulis_bita", writeBytesFn)
}
