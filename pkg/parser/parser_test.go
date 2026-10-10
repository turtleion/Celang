package parser_test

import (
	"testing"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/uast"
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

func TestParser_VisibilityAndImports(t *testing.T) {
	input := `
	publik! buat nama = "Kayasa"
	privat-modul! fungsi komunikasikan(ip: Teks) mengembalikan Kebenaran {
		kembalikan benar
	}
	privat-["./program5.ce"]! buat umur = 25
	publik! { nama, umur }

	ambil { nama } dari "./program1.ce"
	ambil * sebagai math dari "./matematika.ce"
	ambil { tambah sebagai penjumlahan } dari "./matematika.ce"
	import { perkalian } from "./kalkulator.ce"
	from "./util.ce" import cetak_baris, cetak_kolom as kolom
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser returned error: %v", err)
	}

	if len(prog.Body) != 9 {
		t.Fatalf("Expected 9 statements, got %d", len(prog.Body))
	}
}

func TestParser_TryCatchAndThrow(t *testing.T) {
	input := `
	coba {
		buat hasil = 100 / 0
		lempar "Terjadi masalah"
	} tangkap (kesalahan) {
		cetak(kesalahan)
	} akhirnya {
		cetak("Selesai")
	}

	try {
		throw "Error in JS"
	} catch (err) {
		cetak(err)
	} finally {
		cetak("Done")
	}
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 2 {
		t.Fatalf("Expected 2 try-catch statements, got %d", len(prog.Body))
	}
}

func TestParser_PipelineAndFunctionExpression(t *testing.T) {
	input := `
	buat hasil = "  halo dunia  "
		|> potong_spasi()
		|> kapitalisasi()

	buat angka = [1, 2, 3]
		|> saring(fungsi(x) { kembalikan x > 1 })
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(prog.Body))
	}

	decl1, ok := prog.Body[0].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", prog.Body[0])
	}
	pipe1, ok := decl1.Value.(*uast.PipelineExpression)
	if !ok {
		t.Fatalf("Expected PipelineExpression, got %T", decl1.Value)
	}
	if _, ok := pipe1.Left.(*uast.PipelineExpression); !ok {
		t.Fatalf("Expected nested PipelineExpression on Left, got %T", pipe1.Left)
	}

	decl2, ok := prog.Body[1].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Expected VariableDeclaration, got %T", prog.Body[1])
	}
	pipe2, ok := decl2.Value.(*uast.PipelineExpression)
	if !ok {
		t.Fatalf("Expected PipelineExpression, got %T", decl2.Value)
	}
	call2, ok := pipe2.Right.(*uast.CallExpression)
	if !ok {
		t.Fatalf("Expected CallExpression on Right, got %T", pipe2.Right)
	}
	if len(call2.Arguments) != 1 {
		t.Fatalf("Expected 1 argument in saring call, got %d", len(call2.Arguments))
	}
	if _, ok := call2.Arguments[0].(*uast.FunctionExpression); !ok {
		t.Fatalf("Expected FunctionExpression in argument, got %T", call2.Arguments[0])
	}
}

func TestParser_StyleBlockStatement(t *testing.T) {
	input := `
	buat nama_pengguna = "Budi"

	gaya "javascript" {
		let status_login = true
		function cek_akses(user) {
			return user === "Budi" && status_login
		}
		let diizinkan = cek_akses(nama_pengguna)
	}

	gaya "pythonic" {
		def hitung_bonus(skor): {
			return skor * 2
		}
		total_skor = hitung_bonus(50)
	}
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(prog.Body))
	}

	block1, ok := prog.Body[1].(*uast.StyleBlockStatement)
	if !ok {
		t.Fatalf("Expected StyleBlockStatement at index 1, got %T", prog.Body[1])
	}
	if block1.StyleName != "javascript" {
		t.Fatalf("Expected StyleName 'javascript', got %q", block1.StyleName)
	}
	if len(block1.Body.Statements) != 3 {
		t.Fatalf("Expected 3 statements inside javascript block, got %d", len(block1.Body.Statements))
	}

	block2, ok := prog.Body[2].(*uast.StyleBlockStatement)
	if !ok {
		t.Fatalf("Expected StyleBlockStatement at index 2, got %T", prog.Body[2])
	}
	if block2.StyleName != "pythonic" {
		t.Fatalf("Expected StyleName 'pythonic', got %q", block2.StyleName)
	}
}

func TestParser_StringInterpolation(t *testing.T) {
	input := `
	buat pesan1 = "Halo {nama}, saldo: Rp {saldo}."
	buat pesan2 = ` + "`" + `Total: ${harga * 2}` + "`" + `
	buat pesan3 = "Teks tanpa kurung"
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(prog.Body))
	}

	decl1 := prog.Body[0].(*uast.VariableDeclaration)
	templ1, ok := decl1.Value.(*uast.TemplateLiteral)
	if !ok {
		t.Fatalf("Expected TemplateLiteral for pesan1, got %T", decl1.Value)
	}
	if len(templ1.Parts) != 5 {
		t.Fatalf("Expected 5 parts in templ1, got %d", len(templ1.Parts))
	}

	decl2 := prog.Body[1].(*uast.VariableDeclaration)
	templ2, ok := decl2.Value.(*uast.TemplateLiteral)
	if !ok {
		t.Fatalf("Expected TemplateLiteral for pesan2, got %T", decl2.Value)
	}
	if len(templ2.Parts) != 2 {
		t.Fatalf("Expected 2 parts in templ2, got %d", len(templ2.Parts))
	}

	decl3 := prog.Body[2].(*uast.VariableDeclaration)
	if _, ok := decl3.Value.(*uast.Literal); !ok {
		t.Fatalf("Expected plain Literal for pesan3, got %T", decl3.Value)
	}
}

func TestParser_Destructuring(t *testing.T) {
	input := `
	buat [x, y] = [10, 20]
	buat { nama, peran: p } = profil
	[a, b] = koordinat
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(prog.Body))
	}

	d1, ok := prog.Body[0].(*uast.DestructureDeclaration)
	if !ok || d1.Kind != uast.DestructureArray || len(d1.Bindings) != 2 {
		t.Fatalf("Expected Array DestructureDeclaration, got %#v", prog.Body[0])
	}

	d2, ok := prog.Body[1].(*uast.DestructureDeclaration)
	if !ok || d2.Kind != uast.DestructureMap || len(d2.Bindings) != 2 {
		t.Fatalf("Expected Map DestructureDeclaration, got %#v", prog.Body[1])
	}
	if d2.Bindings[1].SourceKey != "peran" || d2.Bindings[1].Target != "p" {
		t.Fatalf("Expected peran: p alias binding, got %#v", d2.Bindings[1])
	}

	d3, ok := prog.Body[2].(*uast.DestructureAssignment)
	if !ok || d3.Kind != uast.DestructureArray || len(d3.Bindings) != 2 {
		t.Fatalf("Expected Array DestructureAssignment, got %#v", prog.Body[2])
	}
}

func TestParser_StructDeclaration(t *testing.T) {
	input := `
	bentuk RekeningBank {
		pemilik: Teks,
		saldo itu Angka = 0,

		fungsi setor(jumlah: Angka) {
			diri.saldo += jumlah
		}
	}
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Body))
	}

	s, ok := prog.Body[0].(*uast.StructDeclaration)
	if !ok {
		t.Fatalf("Expected StructDeclaration, got %T", prog.Body[0])
	}
	if s.Name != "RekeningBank" {
		t.Errorf("Expected struct name RekeningBank, got %s", s.Name)
	}
	if len(s.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(s.Fields))
	}
	if len(s.Methods) != 1 {
		t.Errorf("Expected 1 method, got %d", len(s.Methods))
	}
	if s.Methods[0].Name != "setor" {
		t.Errorf("Expected method name setor, got %s", s.Methods[0].Name)
	}
}



