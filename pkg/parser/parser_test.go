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

func TestParser_ErgonomicFeatures(t *testing.T) {
	input := `
	buat kota = mungkin pengguna.profil.alamat.kota
	buat nama = input_nama kalau kosong "Anonim"
	buat port = cfg.port ?? 8080
	saat_selesai tutup_berkas(berkas)
	buat rentang_angka = 1..10
	buat sub_arr = arr[1..3]
	buat open_slice = arr[2..]
	buat prefix_slice = arr[..4]
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(prog.Body) != 8 {
		t.Fatalf("Expected 8 statements, got %d", len(prog.Body))
	}

	// 1. mungkin
	stmt0, ok := prog.Body[0].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 0 expected VariableDeclaration, got %T", prog.Body[0])
	}
	safeNav, ok := stmt0.Value.(*uast.SafeNavigationExpression)
	if !ok {
		t.Fatalf("Expected SafeNavigationExpression, got %T", stmt0.Value)
	}
	if safeNav.Value == nil {
		t.Fatal("SafeNavigationExpression.Value is nil")
	}

	// 2. kalau kosong
	stmt1, ok := prog.Body[1].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 1 expected VariableDeclaration, got %T", prog.Body[1])
	}
	coalesce1, ok := stmt1.Value.(*uast.NullCoalesceExpression)
	if !ok {
		t.Fatalf("Expected NullCoalesceExpression, got %T", stmt1.Value)
	}
	if coalesce1.Left == nil || coalesce1.Right == nil {
		t.Fatal("NullCoalesceExpression has nil Left or Right")
	}

	// 3. ??
	stmt2, ok := prog.Body[2].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 2 expected VariableDeclaration, got %T", prog.Body[2])
	}
	if _, ok := stmt2.Value.(*uast.NullCoalesceExpression); !ok {
		t.Fatalf("Expected NullCoalesceExpression for ??, got %T", stmt2.Value)
	}

	// 4. saat_selesai
	stmt3, ok := prog.Body[3].(*uast.DeferStatement)
	if !ok {
		t.Fatalf("Stmt 3 expected DeferStatement, got %T", prog.Body[3])
	}
	if stmt3.Body == nil {
		t.Fatal("DeferStatement.Body is nil")
	}

	// 5. 1..10
	stmt4, ok := prog.Body[4].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 4 expected VariableDeclaration, got %T", prog.Body[4])
	}
	rng, ok := stmt4.Value.(*uast.RangeExpression)
	if !ok {
		t.Fatalf("Expected RangeExpression, got %T", stmt4.Value)
	}
	if rng.Start == nil || rng.End == nil {
		t.Fatal("Expected non-nil Start and End for 1..10")
	}

	// 6. arr[1..3]
	stmt5, ok := prog.Body[5].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 5 expected VariableDeclaration, got %T", prog.Body[5])
	}
	mem5, ok := stmt5.Value.(*uast.MemberExpression)
	if !ok {
		t.Fatalf("Expected MemberExpression for slice, got %T", stmt5.Value)
	}
	rng5, ok := mem5.PropertyExpr.(*uast.RangeExpression)
	if !ok || rng5.Start == nil || rng5.End == nil {
		t.Fatalf("Expected RangeExpression with start and end, got %#v", mem5.PropertyExpr)
	}

	// 7. arr[2..]
	stmt6, ok := prog.Body[6].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 6 expected VariableDeclaration, got %T", prog.Body[6])
	}
	mem6, ok := stmt6.Value.(*uast.MemberExpression)
	if !ok {
		t.Fatalf("Expected MemberExpression for open slice, got %T", stmt6.Value)
	}
	rng6, ok := mem6.PropertyExpr.(*uast.RangeExpression)
	if !ok || rng6.Start == nil || rng6.End != nil {
		t.Fatalf("Expected RangeExpression with start and nil end, got %#v", mem6.PropertyExpr)
	}

	// 8. arr[..4]
	stmt7, ok := prog.Body[7].(*uast.VariableDeclaration)
	if !ok {
		t.Fatalf("Stmt 7 expected VariableDeclaration, got %T", prog.Body[7])
	}
	mem7, ok := stmt7.Value.(*uast.MemberExpression)
	if !ok {
		t.Fatalf("Expected MemberExpression for prefix slice, got %T", stmt7.Value)
	}
	rng7, ok := mem7.PropertyExpr.(*uast.RangeExpression)
	if !ok || rng7.Start != nil || rng7.End == nil {
		t.Fatalf("Expected RangeExpression with nil start and non-nil end, got %#v", mem7.PropertyExpr)
	}
}

func TestParser_ArrowFunctions(t *testing.T) {
	input := `
	buat f1 = x => x * 2
	buat f2 = (a, b) => a + b
	buat f3 = () => "halo"
	buat f4 = (x: Angka) => x + 1
	buat f5 = x => { cetak(x) }
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	if len(p.Errors()) > 0 {
		t.Fatalf("Parser had errors: %v", p.Errors())
	}
	if len(prog.Body) != 5 {
		t.Fatalf("Expected 5 statements, got %d", len(prog.Body))
	}

	// f1: x => x * 2
	decl1 := prog.Body[0].(*uast.VariableDeclaration)
	fn1, ok := decl1.Value.(*uast.FunctionExpression)
	if !ok || len(fn1.Parameters) != 1 || fn1.Parameters[0].Name != "x" {
		t.Fatalf("Expected 1 param 'x' for f1, got %#v", fn1)
	}

	// f2: (a, b) => a + b
	decl2 := prog.Body[1].(*uast.VariableDeclaration)
	fn2, ok := decl2.Value.(*uast.FunctionExpression)
	if !ok || len(fn2.Parameters) != 2 || fn2.Parameters[0].Name != "a" || fn2.Parameters[1].Name != "b" {
		t.Fatalf("Expected 2 params 'a','b' for f2, got %#v", fn2)
	}

	// f3: () => "halo"
	decl3 := prog.Body[2].(*uast.VariableDeclaration)
	fn3, ok := decl3.Value.(*uast.FunctionExpression)
	if !ok || len(fn3.Parameters) != 0 {
		t.Fatalf("Expected 0 params for f3, got %#v", fn3)
	}

	// f4: (x: Angka) => x + 1
	decl4 := prog.Body[3].(*uast.VariableDeclaration)
	fn4, ok := decl4.Value.(*uast.FunctionExpression)
	if !ok || len(fn4.Parameters) != 1 || fn4.Parameters[0].TypeHint != "Angka" {
		t.Fatalf("Expected 1 typed param for f4, got %#v", fn4)
	}

	// f5: x => { cetak(x) }
	decl5 := prog.Body[4].(*uast.VariableDeclaration)
	fn5, ok := decl5.Value.(*uast.FunctionExpression)
	if !ok || len(fn5.Body.Statements) != 1 {
		t.Fatalf("Expected block with 1 stmt for f5, got %#v", fn5)
	}
}

func TestParser_Generics(t *testing.T) {
	input := `
	bentuk Kotak<T> {
		isi: T
	}

	bentuk Pasangan<K, V> {
		kunci: K,
		nilai: V
	}

	bentuk Terikat<Nilai: Angka> {
		angka: Nilai
	}

	bentuk Harus<Nilai harus Angka> {
		angka: Nilai
	}

	fungsi tukar<T>(a: T, b: T) {
		kembalikan [b, a]
	}

	buat k = Kotak<Angka>(42)
	buat p = Pasangan<Teks, Angka>("umur", 25)
	buat sal = Saluran<Angka>(5)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	if len(p.Errors()) > 0 {
		t.Fatalf("Parser had errors: %v", p.Errors())
	}

	// 1. Kotak<T>
	st1 := prog.Body[0].(*uast.StructDeclaration)
	if len(st1.TypeParameters) != 1 || st1.TypeParameters[0].Name != "T" {
		t.Fatalf("Expected 1 type parameter 'T', got %#v", st1.TypeParameters)
	}

	// 2. Pasangan<K, V>
	st2 := prog.Body[1].(*uast.StructDeclaration)
	if len(st2.TypeParameters) != 2 || st2.TypeParameters[0].Name != "K" || st2.TypeParameters[1].Name != "V" {
		t.Fatalf("Expected 2 type parameters 'K', 'V', got %#v", st2.TypeParameters)
	}

	// 3. Terikat<Nilai: Angka>
	st3 := prog.Body[2].(*uast.StructDeclaration)
	if len(st3.TypeParameters) != 1 || st3.TypeParameters[0].Constraint != "Angka" {
		t.Fatalf("Expected constraint 'Angka', got %#v", st3.TypeParameters)
	}

	// 4. Harus<Nilai harus Angka>
	st4 := prog.Body[3].(*uast.StructDeclaration)
	if len(st4.TypeParameters) != 1 || st4.TypeParameters[0].Constraint != "Angka" {
		t.Fatalf("Expected constraint 'Angka' for Harus, got %#v", st4.TypeParameters)
	}

	// 5. fungsi tukar<T>
	fn1 := prog.Body[4].(*uast.FunctionDeclaration)
	if len(fn1.TypeParameters) != 1 || fn1.TypeParameters[0].Name != "T" {
		t.Fatalf("Expected type parameter 'T' on function, got %#v", fn1.TypeParameters)
	}

	// 6. Kotak<Angka>(42)
	call1 := prog.Body[5].(*uast.VariableDeclaration).Value.(*uast.CallExpression)
	if len(call1.TypeArguments) != 1 || call1.TypeArguments[0] != "Angka" {
		t.Fatalf("Expected 1 type argument 'Angka', got %#v", call1.TypeArguments)
	}

	// 7. Pasangan<Teks, Angka>("umur", 25)
	call2 := prog.Body[6].(*uast.VariableDeclaration).Value.(*uast.CallExpression)
	if len(call2.TypeArguments) != 2 || call2.TypeArguments[0] != "Teks" || call2.TypeArguments[1] != "Angka" {
		t.Fatalf("Expected type arguments ['Teks', 'Angka'], got %#v", call2.TypeArguments)
	}

	// 8. Saluran<Angka>(5)
	call3 := prog.Body[7].(*uast.VariableDeclaration).Value.(*uast.CallExpression)
	if len(call3.TypeArguments) != 1 || call3.TypeArguments[0] != "Angka" {
		t.Fatalf("Expected type argument 'Angka' on Saluran, got %#v", call3.TypeArguments)
	}
}

func TestParser_TitipJemput(t *testing.T) {
	input := `
	buat tiket = titip hitung(10)
	buat hasil = jemput tiket
	buat hasil2 = jemput tiket dalam 3
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	if len(p.Errors()) > 0 {
		t.Fatalf("Parser had errors: %v", p.Errors())
	}

	// 1. titip hitung(10)
	v1 := prog.Body[0].(*uast.VariableDeclaration)
	t1, ok := v1.Value.(*uast.TitipExpression)
	if !ok {
		t.Fatalf("Expected TitipExpression, got %T", v1.Value)
	}
	if _, isCall := t1.Argument.(*uast.CallExpression); !isCall {
		t.Fatalf("Expected CallExpression as titip argument, got %T", t1.Argument)
	}

	// 2. jemput tiket
	v2 := prog.Body[1].(*uast.VariableDeclaration)
	j1, ok := v2.Value.(*uast.JemputExpression)
	if !ok {
		t.Fatalf("Expected JemputExpression, got %T", v2.Value)
	}
	if j1.Timeout != nil {
		t.Fatalf("Expected nil timeout, got %#v", j1.Timeout)
	}

	// 3. jemput tiket dalam 3
	v3 := prog.Body[2].(*uast.VariableDeclaration)
	j2, ok := v3.Value.(*uast.JemputExpression)
	if !ok {
		t.Fatalf("Expected JemputExpression, got %T", v3.Value)
	}
	if j2.Timeout == nil {
		t.Fatalf("Expected non-nil timeout")
	}
}




