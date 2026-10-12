package vm_test

import (
	"bytes"
	"testing"

	"celang/pkg/compiler"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/vm"
)

func runVM(t *testing.T, input string) (evaluator.Object, *evaluator.Environment, *bytes.Buffer) {
	t.Helper()
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	var out bytes.Buffer
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltinsWithIO(env, nil, &out)

	virtualMachine := vm.NewWithEnv(chunk, env)
	res, err := virtualMachine.Run()
	if err != nil {
		t.Fatalf("vm runtime error: %v", err)
	}

	return res, env, &out
}

func TestVM_Arithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"10 + 20 * 3", "70"},
		{"(10 + 20) * 3", "90"},
		{"100 / 4", "25"},
		{"17 % 5", "2"},
		{"-15 + 20", "5"},
		{`"halo " + "dunia"`, "halo dunia"},
	}

	for _, tt := range tests {
		res, _, _ := runVM(t, tt.input)
		if res.Inspect() != tt.expected {
			t.Errorf("input '%s': expected %s, got %s", tt.input, tt.expected, res.Inspect())
		}
	}
}

func TestVM_ConditionalsAndBooleans(t *testing.T) {
	input := `
	buat x = 15
	buat pesan = ""
	jika x > 20 {
		pesan = "lebih besar"
	} selain jika x == 15 {
		pesan = "sama dengan 15"
	} selain {
		pesan = "lebih kecil"
	}
	pesan
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "sama dengan 15" {
		t.Errorf("expected 'sama dengan 15', got %s", res.Inspect())
	}
}

func TestVM_LoopsAndAccumulator(t *testing.T) {
	input := `
	buat sum = 0
	buat i = 1
	selama i <= 100 {
		sum = sum + i
		i = i + 1
	}
	sum
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "5050" {
		t.Errorf("expected 5050, got %s", res.Inspect())
	}
}

func TestVM_ForLoop(t *testing.T) {
	input := `
	buat total = 0
	untuk x dalam [1, 2, 3, 4] {
		total = total + x
	}
	total
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "10" {
		t.Errorf("expected 10, got %s", res.Inspect())
	}
}

func TestVM_DataStructures(t *testing.T) {
	input := `
	buat arr = [10, 20, 30]
	arr[1] = 99
	buat kamus = {"nama": "Celang", "versi": 2}
	kamus["status"] = "cepat"
	arr[1] + " - " + kamus["status"]
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "99 - cepat" {
		t.Errorf("expected '99 - cepat', got %s", res.Inspect())
	}
}

func TestVM_Functions(t *testing.T) {
	input := `
	fungsi faktorial(n) {
		jika n <= 1 {
			kembalikan 1
		}
		kembalikan n * faktorial(n - 1)
	}
	faktorial(5)
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "120" {
		t.Errorf("expected 120, got %s", res.Inspect())
	}
}

func TestVM_BuiltinInteroperability(t *testing.T) {
	input := `
	cetak("Halo dari Bytecode VM!")
	buat p = panjang([1, 2, 3, 4, 5])
	buat akar = matematika.akar(16)
	p + akar
	`
	res, _, out := runVM(t, input)
	if out.String() != "Halo dari Bytecode VM!\n" {
		t.Errorf("expected output 'Halo dari Bytecode VM!\\n', got %q", out.String())
	}
	if res.Inspect() != "9" {
		t.Errorf("expected 9, got %s", res.Inspect())
	}
}

func BenchmarkLoop_ASTEvaluator(b *testing.B) {
	code := `
	buat sum = 0
	buat i = 0
	selama i < 1000 {
		sum = sum + 1
		i = i + 1
	}
	`
	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		env := evaluator.NewEnvironment()
		evaluator.Eval(prog, env)
	}
}

func BenchmarkLoop_BytecodeVM(b *testing.B) {
	code := `
	buat sum = 0
	buat i = 0
	selama i < 1000 {
		sum = sum + 1
		i = i + 1
	}
	`
	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()
	comp := compiler.NewCompiler()
	chunk, _ := comp.Compile(prog)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		env := evaluator.NewEnvironment()
		vm := vm.NewWithEnv(chunk, env)
		vm.Run()
	}
}

func BenchmarkFunctionLoop_ASTEvaluator(b *testing.B) {
	code := `
	fungsi hitung() {
		buat sum = 0
		buat i = 0
		selama i < 1000 {
			sum = sum + 1
			i = i + 1
		}
		kembalikan sum
	}
	hitung()
	`
	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		env := evaluator.NewEnvironment()
		evaluator.Eval(prog, env)
	}
}

func BenchmarkFunctionLoop_BytecodeVM(b *testing.B) {
	code := `
	fungsi hitung() {
		buat sum = 0
		buat i = 0
		selama i < 1000 {
			sum = sum + 1
			i = i + 1
		}
		kembalikan sum
	}
	hitung()
	`
	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()
	comp := compiler.NewCompiler()
	chunk, _ := comp.Compile(prog)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		vm := vm.New(chunk)
		vm.Run()
	}
}

func TestVM_SuperBigUpdateFeatures(t *testing.T) {
	// Test struct with @ shorthand and private block
	code1 := `
	bentuk Dompet {
		privat saldo: Angka = 500,
		setor(n: Angka) {
			@saldo += n
			kembalikan @saldo
		}
	}
	buat d = Dompet()
	d.setor(250)
	`
	res1, _, _ := runVM(t, code1)
	if res1.Inspect() != "750" {
		t.Errorf("expected 750, got %s", res1.Inspect())
	}

	// Test casting and type alias
	code2 := `
	tipe ID = Angka
	buat x = "123" sebagai ID
	x + 1
	`
	res2, _, _ := runVM(t, code2)
	if res2.Inspect() != "124" {
		t.Errorf("expected 124, got %s", res2.Inspect())
	}

	// Test transparent reference
	code3 := `
	buat a = 10
	rujuk b = a
	b = 40
	a
	`
	res3, _, _ := runVM(t, code3)
	if res3.Inspect() != "40" {
		t.Errorf("expected 40, got %s", res3.Inspect())
	}
}

func TestVM_ErgonomicFeatures(t *testing.T) {
	// 1. Safe navigation (mungkin) & Null Coalesce (kalau kosong / ??)
	code1 := `
	bentuk Detail {
		kota: Teks
	}
	bentuk Akun {
		detail: Detail
	}
	buat a = Akun(Detail("Surabaya"))
	buat b = kosong

	buat k1 = mungkin a.detail.kota
	buat k2 = mungkin b.detail.kota kalau kosong "Jakarta"
	k1 + " " + k2
	`
	res1, _, _ := runVM(t, code1)
	if res1.Inspect() != "\"Surabaya Jakarta\"" && res1.Inspect() != "Surabaya Jakarta" {
		t.Errorf("expected 'Surabaya Jakarta', got %s", res1.Inspect())
	}

	// 2. Defer (saat_selesai)
	code2 := `
	buat log = []
	fungsi tes() {
		saat_selesai log.tambah("tutup")
		log.tambah("kerja")
		kembalikan 100
	}
	tes()
	log[1]
	`
	res2, _, _ := runVM(t, code2)
	if res2.Inspect() != "\"tutup\"" && res2.Inspect() != "tutup" {
		t.Errorf("expected 'tutup', got %s", res2.Inspect())
	}

	// 3. Range & Slicing
	code3 := `
	buat angka = 1..5
	buat sub = angka[1..3]
	sub[0] + sub[1]
	`
	res3, _, _ := runVM(t, code3)
	if res3.Inspect() != "5" { // sub is [2, 3] -> 2 + 3 = 5
		t.Errorf("expected 5, got %s", res3.Inspect())
	}
}

