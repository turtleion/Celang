//go:build !wasm && !js

package evaluator_test

import (
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"runtime"
	"testing"
)

func testEval(input string) evaluator.Object {
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)
	evaluator.RegisterStdlibGlobals(env)
	return evaluator.Eval(prog, env)
}

func TestFFIMemoryAndPointers(t *testing.T) {
	input := `
	impor ffi

	let ptr = ffi.alloc(128)
	let addr = ffi.address(ptr)
	let bukanNull = !ffi.is_null(ptr)
	let nullVal = ffi.is_null(ffi.NULL)

	ffi.write_string(ptr, "Halo FFI Celang")
	let dibaca = ffi.read_string(ptr)
	ffi.free(ptr)

	[addr > 0, bukanNull, nullVal, dibaca]
	`
	res := testEval(input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("hasil evaluasi bukan Array, dapat %T: %v", res, res)
	}

	if b, ok := arr.Elements[0].(*evaluator.Boolean); !ok || !b.Value {
		t.Errorf("addr > 0 diharapkan benar")
	}
	if b, ok := arr.Elements[1].(*evaluator.Boolean); !ok || !b.Value {
		t.Errorf("bukanNull diharapkan benar")
	}
	if b, ok := arr.Elements[2].(*evaluator.Boolean); !ok || !b.Value {
		t.Errorf("nullVal diharapkan benar")
	}
	if s, ok := arr.Elements[3].(*evaluator.String); !ok || s.Value != "Halo FFI Celang" {
		t.Errorf("dibaca salah: %v", arr.Elements[3])
	}
}

func TestFFITypeConstants(t *testing.T) {
	input := `
	impor ffi

	[ffi.Pointer == "pointer", ffi.Penunjuk == "pointer", ffi.Int == "int", ffi.Float == "float", ffi.String == "string", ffi.Void == "void"]
	`
	res := testEval(input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("hasil evaluasi bukan Array, dapat %T: %v", res, res)
	}

	for i, elem := range arr.Elements {
		b, ok := elem.(*evaluator.Boolean)
		if !ok || !b.Value {
			t.Errorf("elemen ke-%d konstanta tipe FFI tidak cocok: %v", i, elem)
		}
	}
}

func TestFFIBytesReadWrite(t *testing.T) {
	input := `
	impor ffi

	let ptr = ffi.alokasi(16)
	let dataAsli = bita([65, 66, 67, 68])
	ffi.tulis_bita(ptr, dataAsli)
	let dataHasil = ffi.baca_bita(ptr, 4)
	ffi.bebas(ptr)

	[dataHasil[0], dataHasil[1], dataHasil[2], dataHasil[3]]
	`
	res := testEval(input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("hasil evaluasi bukan Array, dapat %T: %v", res, res)
	}

	expected := []int64{65, 66, 67, 68}
	for i, exp := range expected {
		n, ok := arr.Elements[i].(*evaluator.Integer)
		if !ok || n.Value != exp {
			t.Errorf("bita ke-%d salah: diharapkan %d, dapat %v", i, exp, arr.Elements[i])
		}
	}
}

func TestFFILibCCall(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Melewati uji LibC di sistem non-Unix")
	}

	libName := "libc.so.6"
	if runtime.GOOS == "darwin" {
		libName = "libSystem.B.dylib"
	}

	input := `
	impor ffi

	let libc = ffi.open("` + libName + `")
	let strlen = libc.fungsi("strlen", [ffi.Pointer], ffi.Int)
	let hasilLen = strlen("CelangFFI")

	hasilLen
	`
	res := testEval(input)
	intRes, ok := res.(*evaluator.Integer)
	if !ok {
		t.Fatalf("strlen gagal atau mengembalikan bukan Integer, dapat %T: %v", res, res)
	}

	if intRes.Value != 9 {
		t.Errorf("strlen('CelangFFI') diharapkan 9, didapat %d", intRes.Value)
	}
}

func TestFFILibMFloatCall(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("Melewati uji LibM di sistem non-Unix")
	}

	libName := "libm.so.6"
	if runtime.GOOS == "darwin" {
		libName = "libSystem.B.dylib"
	}

	input := `
	impor ffi

	let libm = ffi.open("` + libName + `")
	let sqrt = libm.fungsi("sqrt", [ffi.Float], ffi.Float)
	let hasilSqrt = sqrt(16.0)

	hasilSqrt
	`
	res := testEval(input)
	floatRes, ok := res.(*evaluator.Float)
	if !ok {
		t.Fatalf("sqrt gagal atau mengembalikan bukan Float, dapat %T: %v", res, res)
	}

	if floatRes.Value != 4.0 {
		t.Errorf("sqrt(16.0) diharapkan 4.0, didapat %f", floatRes.Value)
	}
}
