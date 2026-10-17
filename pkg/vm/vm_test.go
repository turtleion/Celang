package vm_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestVM_ForLoop_String(t *testing.T) {
	input := `
	buat hasil = ""
	untuk ch dalam "celang" {
		hasil = hasil + ch + "-"
	}
	hasil
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "\"c-e-l-a-n-g-\"" && res.Inspect() != "c-e-l-a-n-g-" {
		t.Errorf("expected c-e-l-a-n-g-, got %s", res.Inspect())
	}
}

func TestVM_ForLoop_Map(t *testing.T) {
	input := `
	buat keys = ""
	buat m = {"a": 1, "b": 2}
	untuk k dalam m {
		keys = keys + k
	}
	keys
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "\"ab\"" && res.Inspect() != "ab" {
		t.Errorf("expected ab, got %s", res.Inspect())
	}
}

func TestVM_ForLoop_Channel(t *testing.T) {
	input := `
	buat ch = Saluran(3)
	ch.kirim(10)
	ch.kirim(20)
	ch.kirim(30)
	ch.tutup()

	buat sum = 0
	untuk val dalam ch {
		sum = sum + val
	}
	sum
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "60" {
		t.Errorf("expected 60, got %s", res.Inspect())
	}
}

func TestVM_Channel_Panjang(t *testing.T) {
	input := `
	buat ch = Saluran(5)
	ch.kirim("a")
	ch.kirim("b")
	panjang(ch) + " - " + ch.panjang
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "\"2 - 2\"" && res.Inspect() != "2 - 2" {
		t.Errorf("expected '2 - 2', got %s", res.Inspect())
	}
}

func TestVM_Berkas_StreamAndJSON(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "vm_stream.log")
	content := "data 1\ndata 2\n"
	_ = os.WriteFile(logFile, []byte(content), 0644)

	input := `
	buat list = []
	untuk baris dalam berkas.alir("` + logFile + `") {
		list.tambah(baris)
	}
	list[0] + " & " + list[1]
	`
	res, _, _ := runVM(t, input)
	if res.Inspect() != "\"data 1 & data 2\"" && res.Inspect() != "data 1 & data 2" {
		t.Errorf("expected 'data 1 & data 2', got %s", res.Inspect())
	}
}

func TestVM_Teks_ModernFeatures(t *testing.T) {
	input := `
	buat s = teks.slug("Belajar Celang Modern!")
	buat m = teks.samarkan("user@example.com")
	buat cl = teks.terdekat("kucing", ["anjing", "kucing_lucu", "kucin"])
	s + " | " + m + " | " + cl
	`
	res, _, _ := runVM(t, input)
	expected := "\"belajar-celang-modern | u***r@example.com | kucin\""
	if res.Inspect() != expected && res.Inspect() != "belajar-celang-modern | u***r@example.com | kucin" {
		t.Errorf("expected %s, got %s", expected, res.Inspect())
	}
}

func TestVM_Matematika_ModernFeatures(t *testing.T) {
	input := `
	buat rp = matematika.rupiah(250000)
	buat words = matematika.terbilang(1500)
	buat avg = matematika.rata_rata([10, 20, 30])
	buat dist = matematika.jarak(0, 0, 3, 4)
	rp + " | " + words + " | " + avg + " | " + dist
	`
	res, _, _ := runVM(t, input)
	expected := "\"Rp 250.000 | seribu lima ratus | 20 | 5\""
	if res.Inspect() != expected && res.Inspect() != "Rp 250.000 | seribu lima ratus | 20 | 5" {
		t.Errorf("expected %s, got %s", expected, res.Inspect())
	}
}

func TestVM_Waktu_ModernFeatures(t *testing.T) {
	input := `
	buat t0 = waktu.urai("2026-10-15 12:00:00")
	buat t_indo = waktu.format_indonesia(t0, "tanggal")
	buat rel = waktu.lalu(t0 - 120000, t0)
	buat dur = waktu.format_durasi(waktu.menit(5) + waktu.detik(30))
	t_indo + " | " + rel + " | " + dur
	`
	res, _, _ := runVM(t, input)
	expected := "\"15 Oktober 2026 | 2 menit yang lalu | 5 menit 30 detik\""
	if res.Inspect() != expected && res.Inspect() != "15 Oktober 2026 | 2 menit yang lalu | 5 menit 30 detik" {
		t.Errorf("expected %s, got %s", expected, res.Inspect())
	}
}

func TestVM_JSON_ModernFeatures(t *testing.T) {
	input := `
	buat data = json.uraikan("{\"app\": {\"nama\": \"Celang\", \"versi\": 2}}")
	json.setel(data, "app.fitur.modern", benar)
	buat nama = json.dapatkan(data, "app.nama")
	buat modern = json.dapatkan(data, "app.fitur.modern")
	buat masked = json.samarkan({"user": "Andi", "token": "secret123"})
	nama + " | " + modern + " | " + masked.token
	`
	res, _, _ := runVM(t, input)
	expected := "\"Celang | true | ***\""
	if res.Inspect() != expected && res.Inspect() != "Celang | true | ***" {
		t.Errorf("expected %s, got %s", expected, res.Inspect())
	}
}

func TestVM_Sistem_ModernFeatures(t *testing.T) {
	evaluator.SetScriptArgs([]string{"halo_vm", "arg_kedua"})
	input := `
	buat plat = sistem.os()
	sistem.setel_env("CELANG_VM_ENV", "nilai_123")
	buat val = sistem.dapatkan_env("CELANG_VM_ENV")
	sistem.hapus_env("CELANG_VM_ENV")
	buat arg0 = sistem.argumen(0)
	buat res = os.exec(["echo", "vm_eksekusi"])
	buat out = teks.trim(res.keluaran)
	plat + " | " + val + " | " + arg0 + " | " + out
	`
	res, _, _ := runVM(t, input)
	expected := "\"" + runtime.GOOS + " | nilai_123 | halo_vm | vm_eksekusi\""
	if res.Inspect() != expected && res.Inspect() != runtime.GOOS+" | nilai_123 | halo_vm | vm_eksekusi" {
		t.Errorf("expected %s, got %s", expected, res.Inspect())
	}
}

func TestVM_ConstImmutability(t *testing.T) {
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, os.Stdout)

	// Step 1: Define const t = 3
	prog1, _ := parser.NewParser(parser.NewLexer("const t = 3")).ParseProgram()
	chunk1, _ := compiler.NewCompilerWithEnv(env).Compile(prog1)
	vm1 := vm.NewWithEnv(chunk1, env)
	if _, err := vm1.Run(); err != nil {
		t.Fatalf("failed to define const: %v", err)
	}

	// Step 2: Try to reassign t = 4 in the same environment
	prog2, _ := parser.NewParser(parser.NewLexer("t = 4")).ParseProgram()
	comp2 := compiler.NewCompilerWithEnv(env)
	chunk2, errComp := comp2.Compile(prog2)
	if errComp == nil {
		// If compile succeeded (e.g. dynamic), VM must catch it at runtime
		vm2 := vm.NewWithEnv(chunk2, env)
		_, errRun := vm2.Run()
		if errRun == nil || !strings.Contains(errRun.Error(), "bersifat tetap") {
			t.Fatalf("expected VM runtime error for const reassignment, got: %v", errRun)
		}
	} else if !strings.Contains(errComp.Error(), "bersifat tetap") {
		t.Fatalf("expected compile error mentioning 'bersifat tetap', got: %v", errComp)
	}
}

func TestVM_AuditFixes(t *testing.T) {
	// 1. Float modulo in VM
	resMod, _, _ := runVM(t, "5.5 % 2")
	if flt, ok := resMod.(*evaluator.Float); !ok || flt.Value != 1.5 {
		t.Fatalf("expected 5.5 %% 2 = 1.5, got %s", resMod.Inspect())
	}

	// 2. Map with property named 'awal' (should not trigger time travel)
	inputAwal := `
	buat data = {"awal": "nilai_awal_map"}
	[data.awal, data["awal"]]
	`
	resAwal, _, _ := runVM(t, inputAwal)
	arrAwal, ok := resAwal.(*evaluator.Array)
	if !ok || len(arrAwal.Elements) != 2 ||
		arrAwal.Elements[0].(*evaluator.String).Value != "nilai_awal_map" ||
		arrAwal.Elements[1].(*evaluator.String).Value != "nilai_awal_map" {
		t.Fatalf("expected ['nilai_awal_map', 'nilai_awal_map'], got %s", resAwal.Inspect())
	}

	// 3. Truthiness parity with Evaluator: empty array is falsy
	inputTruthy := `
	buat res = 0
	jika [] {
		res = 1
	} selain {
		res = 2
	}
	res
	`
	resTruthy, _, _ := runVM(t, inputTruthy)
	if i, ok := resTruthy.(*evaluator.Integer); !ok || i.Value != 2 {
		t.Fatalf("expected empty array to be falsy (res = 2), got %s", resTruthy.Inspect())
	}
}

func TestVM_SatuanBudayaDanTerbilang(t *testing.T) {
	input := `
	buat harga = 50.ribu
	buat gaji = 2.5.juta
	buat hutang = 1.miliar
	[harga, gaji, hutang, harga.rupiah, gaji.terbilang]
	`
	res, _, _ := runVM(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 5 {
		t.Fatalf("expected 5 elements, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 50000 {
		t.Errorf("expected 50000, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 2500000 {
		t.Errorf("expected 2500000, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 1000000000 {
		t.Errorf("expected 1000000000, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.String).Value != "Rp 50.000" {
		t.Errorf("expected 'Rp 50.000', got %v", arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.String).Value != "dua juta lima ratus ribu rupiah" {
		t.Errorf("expected 'dua juta lima ratus ribu rupiah', got %v", arr.Elements[4])
	}
}

func TestVM_TransaksiDanRollback(t *testing.T) {
	input := `
	buat saldo = 100
	buat status = "aman"
	transaksi {
		saldo = 40
		batalkan "saldo kurang"
	} pulihkan (err) {
		status = err
	}
	[saldo, status]
	`
	res, _, _ := runVM(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected saldo 100, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.String).Value != "saldo kurang" {
		t.Errorf("expected status 'saldo kurang', got %v", arr.Elements[1])
	}
}

func TestVM_ReaktifDanPengamat(t *testing.T) {
	input := `
	buat harga = 50.ribu
	buat jumlah = 2
	reaktif total = harga * jumlah
	buat t1 = total
	harga = 100.ribu
	buat t2 = total
	[t1, t2]
	`
	res, _, _ := runVM(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 100000 {
		t.Errorf("expected t1 100000, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 200000 {
		t.Errorf("expected t2 200000, got %v", arr.Elements[1])
	}
}

func TestVM_Channel_UniversalAndIndonesian(t *testing.T) {
	codeIndo := `
	buat ch = Saluran(2)
	ch.kirim(123)
	buat val = ch.terima()
	buat p = ch.panjang
	buat k = ch.kapasitas
	ch.tutup()
	[val, p, k]
	`
	res, _, _ := runVM(t, codeIndo)
	arr := res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.Integer).Value != 123 ||
		arr.Elements[1].(*evaluator.Integer).Value != 0 ||
		arr.Elements[2].(*evaluator.Integer).Value != 2 {
		t.Fatalf("unexpected channel result in VM: %s", arr.Inspect())
	}

	codeUniv := `
	"use universal"
	let ch = Channel(2)
	ch.send(456)
	let val = ch.recv()
	let p = ch.length
	let k = ch.capacity
	ch.close()
	[val, p, k]
	`
	resUniv, _, _ := runVM(t, codeUniv)
	arrUniv := resUniv.(*evaluator.Array)
	if arrUniv.Elements[0].(*evaluator.Integer).Value != 456 ||
		arrUniv.Elements[1].(*evaluator.Integer).Value != 0 ||
		arrUniv.Elements[2].(*evaluator.Integer).Value != 2 {
		t.Fatalf("unexpected universal channel result in VM: %s", arrUniv.Inspect())
	}
}

func TestVM_Ticket_UniversalAndIndonesian(t *testing.T) {
	codeIndo := `
	buat t1 = titip 789
	buat r = t1.jemput()
	buat s = t1.selesai
	[r, s]
	`
	res, _, _ := runVM(t, codeIndo)
	arr := res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.Integer).Value != 789 ||
		!arr.Elements[1].(*evaluator.Boolean).Value {
		t.Fatalf("unexpected ticket result in VM: %s", arr.Inspect())
	}

	codeUniv := `
	"use universal"
	let t1 = spawn 321
	let r = t1.await()
	let s = t1.done
	[r, s]
	`
	resUniv, _, _ := runVM(t, codeUniv)
	arrUniv := resUniv.(*evaluator.Array)
	if arrUniv.Elements[0].(*evaluator.Integer).Value != 321 ||
		!arrUniv.Elements[1].(*evaluator.Boolean).Value {
		t.Fatalf("unexpected universal ticket result in VM: %s", arrUniv.Inspect())
	}
}

func TestVM_MapAndStruct_StringIndexAssignment(t *testing.T) {
	code := `
	bentuk Pengguna {
		nama: Teks
	}
	buat user = Pengguna("Awal")
	user["nama"] = "Budi"

	buat data = {}
	buat k = "kunci"
	data[k] = 999
	[user["nama"], data["kunci"]]
	`
	res, _, _ := runVM(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.String).Value != "Budi" {
		t.Errorf("expected user['nama'] = 'Budi', got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 999 {
		t.Errorf("expected data['kunci'] = 999, got %v", arr.Elements[1])
	}
}

func TestVM_StructuralEquality(t *testing.T) {
	code := `
	buat e1 = (null == "kosong")
	buat e2 = (true == "benar")
	buat e3 = (false == "salah")
	buat e4 = (null == null)
	buat e5 = (true == true)
	buat e6 = (1 == 1.0)
	buat e7 = ([1, 2] == "[1, 2]")
	buat e8 = ([1, 2] == [1, "2"])
	buat e9 = ([1, 2] == [1, 2])
	buat e10 = ({"a": 1, "b": 2} == {"b": 2, "a": 1})

	[e1, e2, e3, e4, e5, e6, e7, e8, e9, e10]
	`
	res, _, _ := runVM(t, code)
	arr := res.(*evaluator.Array)
	expectations := []bool{false, false, false, true, true, true, false, false, true, true}
	for i, exp := range expectations {
		got := arr.Elements[i].(*evaluator.Boolean).Value
		if got != exp {
			t.Errorf("VM e%d expected %v, got %v", i+1, exp, got)
		}
	}
}

func TestAudit_VM_MapLookupNoGlobalFallback(t *testing.T) {
	code := `
	buat tamu = "admin"
	buat peran = {"admin": "AKSES_PENUH"}
	peran["tamu"]
	`
	res, _, _ := runVM(t, code)
	if res.Type() != evaluator.NULL_OBJ {
		t.Fatalf("expected peran['tamu'] to be null without falling back to global 'tamu', got %s", res.Inspect())
	}
}

func TestAudit_VM_OpEvalAST_FunctionLocals(t *testing.T) {
	code := `
	fungsi periksa(nilai) {
		cocok nilai {
			kasus 1:
				kembalikan "satu"
			kasus 2:
				kembalikan "dua"
			lainnya:
				kembalikan "lainnya"
		}
	}
	[periksa(1), periksa(2), periksa(3)]
	`
	res, _, _ := runVM(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.String).Value != "satu" {
		t.Errorf("expected 'satu', got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].(*evaluator.String).Value != "dua" {
		t.Errorf("expected 'dua', got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].(*evaluator.String).Value != "lainnya" {
		t.Errorf("expected 'lainnya', got %s", arr.Elements[2].Inspect())
	}
}




