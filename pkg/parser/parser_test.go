package parser_test

import (
	"testing"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
)

func TestParser_BasicSyntax(t *testing.T) {
	input := `
	// Deklarasi variabel
	buat nama = "Budi"
	buat umur itu Angka = 25
	buat skor = 50
	konstanta MAX_LOGIN = 3

	// Penugasan ulang
	umur = 26
	skor += 10

	// Map dan Array
	buat profil = {"nama": "Budi", "skor" itu 100}
	buat angka = [10, 20, 30]

	// Time-travel akses
	buat prev = umur.sebelumnya

	// Percabangan
	jika (umur >= 18) {
		cetak("Dewasa")
	} selain jika (umur >= 13) {
		cetak("Remaja")
	} selain {
		cetak("Anak-anak")
	}

	// Perulangan
	buat hitung = 1
	selama hitung <= 3 {
		hitung += 1
	}

	untuk item dalam angka {
		cetak(item)
	}

	// Fungsi dengan type hint
	fungsi hitung_luas(panjang itu Angka, lebar: Angka) mengembalikan Angka {
		kembalikan panjang * lebar
	}

	// Pointer
	buat nilai_awal = 100
	buat ptr = &nilai_awal
	*ptr = 200

	// Switch
	buat peran = "admin"
	cocok peran {
		kasus "admin":
			cetak("Akses penuh")
		lainnya:
			cetak("Akses terbatas")
	}

	// Konkurensi
	sambil {
		cetak("Latar belakang")
	}

	bersamaan {
		cetak("Paralel 1")
		cetak("Paralel 2")
	}
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("Parser returned error: %v", err)
	}

	if prog == nil {
		t.Fatalf("Expected non-nil program")
	}

	// Evaluate parsed U-AST to verify end-to-end execution!
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)
	result := evaluator.Eval(prog, env)

	if result != nil && result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Evaluator error on parsed AST: %s", result.Inspect())
	}

	// Verify environment state
	namaVal, ok := env.Get("nama")
	if !ok || namaVal.(*evaluator.String).Value != "Budi" {
		t.Fatalf("Expected nama='Budi', got %+v", namaVal)
	}

	umurVal, ok := env.Get("umur")
	if !ok || umurVal.(*evaluator.Integer).Value != 26 {
		t.Fatalf("Expected umur=26, got %+v", umurVal)
	}

	maxVal, ok := env.Get("MAX_LOGIN")
	if !ok || maxVal.(*evaluator.Integer).Value != 3 {
		t.Fatalf("Expected MAX_LOGIN=3, got %+v", maxVal)
	}

	profilVal, ok := env.Get("profil")
	if !ok {
		t.Fatalf("Expected profil map to exist")
	}
	profilMap, ok := profilVal.(*evaluator.Map)
	if !ok || profilMap.Pairs["nama"].Inspect() != "Budi" || profilMap.Pairs["skor"].(*evaluator.Integer).Value != 100 {
		t.Fatalf("Unexpected profil map content: %+v", profilMap)
	}

	nilaiVal, ok := env.Get("nilai_awal")
	if !ok || nilaiVal.(*evaluator.Integer).Value != 200 {
		t.Fatalf("Expected nilai_awal=200 after pointer mutation, got %+v", nilaiVal)
	}
}

func TestParser_FunctionExecution(t *testing.T) {
	input := `
	fungsi faktorial(n itu Angka) mengembalikan Angka {
		jika (n <= 1) {
			kembalikan 1
		}
		kembalikan n * faktorial(n - 1)
	}

	buat hasil = faktorial(5)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	evaluator.Eval(prog, env)

	hasilVal, ok := env.Get("hasil")
	if !ok || hasilVal.(*evaluator.Integer).Value != 120 {
		t.Fatalf("Expected hasil=120, got %+v", hasilVal)
	}
}
