package evaluator_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/uast"
)

func TestEvaluator_TimeTravelVariables(t *testing.T) {
	// Program:
	// x = 10
	// x = 25
	// x = 50
	// prev = x.sebelumnya   (should be 25)
	// initial = x.awal       (should be 10)
	// current = x            (should be 50)
	prog := uast.NewProgram(
		uast.NewVarDecl("x", uast.NewIntLiteral(10), true, "int"),
		uast.NewAssign(uast.NewIdentifier("x"), "=", uast.NewIntLiteral(25)),
		uast.NewAssign(uast.NewIdentifier("x"), "=", uast.NewIntLiteral(50)),
		uast.NewVarDecl("prev", uast.NewMemberExpr(uast.NewIdentifier("x"), "sebelumnya", false), false, "int"),
		uast.NewVarDecl("initial", uast.NewMemberExpr(uast.NewIdentifier("x"), "awal", false), false, "int"),
		uast.NewVarDecl("hist", uast.NewMemberExpr(uast.NewIdentifier("x"), "histori", false), false, "array"),
	)

	env := evaluator.NewEnvironment()
	evaluator.Eval(prog, env)

	// Test current value of x
	xVal, ok := env.Get("x")
	if !ok || xVal.(*evaluator.Integer).Value != 50 {
		t.Fatalf("Expected x=50, got %+v", xVal)
	}

	// Test x.sebelumnya
	prevVal, ok := env.Get("prev")
	if !ok || prevVal.(*evaluator.Integer).Value != 25 {
		t.Fatalf("Expected prev=25, got %+v", prevVal)
	}

	// Test x.awal
	initVal, ok := env.Get("initial")
	if !ok || initVal.(*evaluator.Integer).Value != 10 {
		t.Fatalf("Expected initial=10, got %+v", initVal)
	}

	// Test x.histori
	histVal, ok := env.Get("hist")
	if !ok {
		t.Fatalf("Expected hist to exist")
	}
	arr, ok := histVal.(*evaluator.Array)
	if !ok || len(arr.Elements) != 3 {
		t.Fatalf("Expected 3 elements in history, got %d", len(arr.Elements))
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 10 ||
		arr.Elements[1].(*evaluator.Integer).Value != 25 ||
		arr.Elements[2].(*evaluator.Integer).Value != 50 {
		t.Fatalf("Unexpected history sequence: %+v", arr.Elements)
	}
}

func TestEvaluator_SmartFallbackCoercion(t *testing.T) {
	tests := []struct {
		name     string
		node     uast.Node
		expected string
	}{
		{
			name: "Pure String Number + Int: '100' + 25 -> 125",
			node: uast.NewBinaryExpr(uast.NewStringLiteral("100"), "+", uast.NewIntLiteral(25)),
			expected: "125",
		},
		{
			name: "String - Int with extraction: '100' - 20 -> 80",
			node: uast.NewBinaryExpr(uast.NewStringLiteral("100"), "-", uast.NewIntLiteral(20)),
			expected: "80",
		},
		{
			name: "RegEx Extraction: 'Rp 50000' / 2 -> 25000",
			node: uast.NewBinaryExpr(uast.NewStringLiteral("Rp 50000"), "/", uast.NewIntLiteral(2)),
			expected: "25000",
		},
		{
			name: "String + Number Human Description: 'Total: ' + 50 -> 'Total: 50'",
			node: uast.NewBinaryExpr(uast.NewStringLiteral("Total: "), "+", uast.NewIntLiteral(50)),
			expected: "Total: 50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := evaluator.NewEnvironment()
			res := evaluator.Eval(tt.node, env)
			if res == nil || res.Inspect() != tt.expected {
				t.Fatalf("Expected %s, got %v (%+v)", tt.expected, res.Inspect(), res)
			}
		})
	}
}

func TestEvaluator_FunctionsAndRecursion(t *testing.T) {
	// Function: faktorial(n)
	// if n <= 1 return 1
	// return n * faktorial(n - 1)
	prog := uast.NewProgram(
		uast.NewFuncDecl(
			"faktorial",
			[]uast.Parameter{{Name: "n"}},
			uast.NewBlock(
				uast.NewIf(
					uast.NewBinaryExpr(uast.NewIdentifier("n"), "<=", uast.NewIntLiteral(1)),
					uast.NewBlock(uast.NewReturn(uast.NewIntLiteral(1))),
					nil,
				),
				uast.NewReturn(
					uast.NewBinaryExpr(
						uast.NewIdentifier("n"),
						"*",
						uast.NewCallExpr(
							uast.NewIdentifier("faktorial"),
							uast.NewBinaryExpr(uast.NewIdentifier("n"), "-", uast.NewIntLiteral(1)),
						),
					),
				),
			),
			"int",
		),
		uast.NewVarDecl(
			"hasil",
			uast.NewCallExpr(uast.NewIdentifier("faktorial"), uast.NewIntLiteral(5)),
			false,
			"int",
		),
	)

	env := evaluator.NewEnvironment()
	evaluator.Eval(prog, env)

	hasil, ok := env.Get("hasil")
	if !ok || hasil.(*evaluator.Integer).Value != 120 {
		t.Fatalf("Expected faktorial(5) = 120, got %+v", hasil)
	}
}

func TestEvaluator_BuiltinsAndOutput(t *testing.T) {
	var buf bytes.Buffer
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, &buf)

	prog := uast.NewProgram(
		uast.NewExprStmt(
			uast.NewCallExpr(
				uast.NewIdentifier("cetak"),
				uast.NewStringLiteral("Halo dari Celang!"),
				uast.NewIntLiteral(2026),
			),
		),
	)

	evaluator.Eval(prog, env)

	expectedOutput := "Halo dari Celang! 2026\n"
	if buf.String() != expectedOutput {
		t.Fatalf("Expected output %q, got %q", expectedOutput, buf.String())
	}
}

func TestEvaluator_Concurrency(t *testing.T) {
	// Test 'bersamaan' (fan-out wait group execution)
	env := evaluator.NewEnvironment()
	env.Define("counter", &evaluator.Integer{Value: 0}, true, "int")

	var mu sync.Mutex
	// Register a thread-safe increment builtin
	env.Define("inc", &evaluator.Builtin{
		Fn: func(args ...evaluator.Object) evaluator.Object {
			mu.Lock()
			defer mu.Unlock()
			cur, _ := env.Get("counter")
			env.Assign("counter", &evaluator.Integer{Value: cur.(*evaluator.Integer).Value + 1})
			return &evaluator.Null{}
		},
	}, false, "function")

	prog := uast.NewProgram(
		uast.NewConcurrencyStmt(
			"bersamaan",
			uast.NewExprStmt(uast.NewCallExpr(uast.NewIdentifier("inc"))),
			uast.NewExprStmt(uast.NewCallExpr(uast.NewIdentifier("inc"))),
			uast.NewExprStmt(uast.NewCallExpr(uast.NewIdentifier("inc"))),
		),
	)

	evaluator.Eval(prog, env)

	cnt, _ := env.Get("counter")
	if cnt.(*evaluator.Integer).Value != 3 {
		t.Fatalf("Expected counter=3 after bersamaan execution, got %v", cnt.Inspect())
	}

	// Test 'sambil' (background async goroutine)
	progAsync := uast.NewProgram(
		uast.NewConcurrencyStmt(
			"sambil",
			uast.NewExprStmt(uast.NewCallExpr(uast.NewIdentifier("inc"))),
		),
	)
	evaluator.Eval(progAsync, env)
	// Give background routine a tiny tick to execute
	time.Sleep(50 * time.Millisecond)

	cntAfter, _ := env.Get("counter")
	if cntAfter.(*evaluator.Integer).Value != 4 {
		t.Fatalf("Expected counter=4 after sambil execution, got %v", cntAfter.Inspect())
	}
}

func TestEvaluator_Pointers(t *testing.T) {
	// Program:
	// buat x = 10
	// buat ptr = &x
	// *ptr = 99
	// cetak(x)            // 99
	// cetak(x.sebelumnya) // 10
	env := evaluator.NewEnvironment()
	prog := uast.NewProgram(
		uast.NewVarDecl("x", uast.NewIntLiteral(10), true, "int"),
		uast.NewVarDecl("ptr", uast.NewUnaryExpr("&", uast.NewIdentifier("x"), true), false, "pointer"),
		uast.NewAssign(uast.NewUnaryExpr("*", uast.NewIdentifier("ptr"), true), "=", uast.NewIntLiteral(99)),
	)

	evaluator.Eval(prog, env)

	xVal, ok := env.Get("x")
	if !ok || xVal.(*evaluator.Integer).Value != 99 {
		t.Fatalf("Expected x=99 via pointer mutation, got %+v", xVal)
	}

	entry, _ := env.GetEntry("x")
	prev, ok := entry.Previous()
	if !ok || prev.(*evaluator.Integer).Value != 10 {
		t.Fatalf("Expected x.sebelumnya=10 after pointer mutation, got %+v", prev)
	}
}

func TestEvaluator_SwitchStatement(t *testing.T) {
	// Program:
	// buat status = "admin"
	// cocok status {
	//     kasus "user": buat hasil = "pengguna"
	//     kasus "admin": buat hasil = "pengelola"
	//     lainnya: buat hasil = "tamu"
	// }
	env := evaluator.NewEnvironment()
	prog := uast.NewProgram(
		uast.NewVarDecl("status", uast.NewStringLiteral("admin"), false, "string"),
		uast.NewSwitch(
			uast.NewIdentifier("status"),
			uast.NewCase(uast.NewStringLiteral("user"), uast.NewVarDecl("hasil", uast.NewStringLiteral("pengguna"), true, "string")),
			uast.NewCase(uast.NewStringLiteral("admin"), uast.NewVarDecl("hasil", uast.NewStringLiteral("pengelola"), true, "string")),
			uast.NewCase(nil, uast.NewVarDecl("hasil", uast.NewStringLiteral("tamu"), true, "string")),
		),
	)

	evaluator.Eval(prog, env)

	hasilVal, ok := env.Get("hasil")
	if !ok || hasilVal.(*evaluator.String).Value != "pengelola" {
		t.Fatalf("Expected hasil='pengelola', got %+v", hasilVal)
	}
}

func TestEvaluator_TryCatchAndThrow(t *testing.T) {
	// Program 1: Catch runtime division-by-zero error and run finally
	// buat hasil = ""
	// coba {
	//     buat x = 10 / 0
	// } tangkap (err) {
	//     hasil = "tertangkap: " + err
	// } akhirnya {
	//     hasil = hasil + " [akhir]"
	// }
	env := evaluator.NewEnvironment()
	prog1 := uast.NewProgram(
		uast.NewVarDecl("hasil", uast.NewStringLiteral(""), true, "string"),
		uast.NewTryCatch(
			uast.NewBlock(
				uast.NewVarDecl("x", uast.NewBinaryExpr(uast.NewIntLiteral(10), "/", uast.NewIntLiteral(0)), false, "int"),
			),
			uast.NewCatchClause(
				"err",
				uast.NewBlock(
					uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewBinaryExpr(uast.NewStringLiteral("tertangkap: "), "+", uast.NewIdentifier("err"))),
				),
			),
			uast.NewBlock(
				uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewBinaryExpr(uast.NewIdentifier("hasil"), "+", uast.NewStringLiteral(" [akhir]"))),
			),
		),
	)

	evaluator.Eval(prog1, env)
	hasilVal, ok := env.Get("hasil")
	if !ok {
		t.Fatalf("Expected hasil to exist")
	}
	expected := "tertangkap: Pembagian dengan nol (zero division) [akhir]"
	if str, ok := hasilVal.(*evaluator.String); !ok || str.Value != expected {
		t.Fatalf("Expected %q, got %+v", expected, hasilVal)
	}

	// Program 2: Throw custom error
	// coba {
	//     lempar "Ada galat kustom"
	// } tangkap (e) {
	//     hasil = e
	// }
	prog2 := uast.NewProgram(
		uast.NewTryCatch(
			uast.NewBlock(
				uast.NewThrow(uast.NewStringLiteral("Ada galat kustom")),
			),
			uast.NewCatchClause(
				"e",
				uast.NewBlock(
					uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewIdentifier("e")),
				),
			),
			nil,
		),
	)

	evaluator.Eval(prog2, env)
	hasilVal2, ok := env.Get("hasil")
	if !ok || hasilVal2.(*evaluator.String).Value != "Ada galat kustom" {
		t.Fatalf("Expected hasil='Ada galat kustom', got %+v", hasilVal2)
	}

	// Program 3: No error in try, finally still executes
	// hasil = "awal"
	// coba {
	//     hasil = "sukses"
	// } tangkap (e) {
	//     hasil = "gagal"
	// } akhirnya {
	//     hasil = hasil + " selesai"
	// }
	prog3 := uast.NewProgram(
		uast.NewTryCatch(
			uast.NewBlock(
				uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewStringLiteral("sukses")),
			),
			uast.NewCatchClause(
				"e",
				uast.NewBlock(
					uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewStringLiteral("gagal")),
				),
			),
			uast.NewBlock(
				uast.NewAssign(uast.NewIdentifier("hasil"), "=", uast.NewBinaryExpr(uast.NewIdentifier("hasil"), "+", uast.NewStringLiteral(" selesai"))),
			),
		),
	)

	evaluator.Eval(prog3, env)
	hasilVal3, ok := env.Get("hasil")
	if !ok || hasilVal3.(*evaluator.String).Value != "sukses selesai" {
		t.Fatalf("Expected hasil='sukses selesai', got %+v", hasilVal3)
	}
}

func TestEvaluator_PipelineOperator(t *testing.T) {
	// Program 1: String pipeline: "  halo dunia  " |> potong_spasi() |> kapitalisasi()
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)

	prog1 := uast.NewProgram(
		uast.NewVarDecl("teks",
			uast.NewPipelineExpr(
				uast.NewPipelineExpr(
					uast.NewStringLiteral("  halo dunia  "),
					uast.NewCallExpr(uast.NewIdentifier("potong_spasi")),
				),
				uast.NewCallExpr(uast.NewIdentifier("kapitalisasi")),
			),
			false, "string",
		),
	)

	evaluator.Eval(prog1, env)
	teksVal, ok := env.Get("teks")
	if !ok || teksVal.(*evaluator.String).Value != "HALO DUNIA" {
		t.Fatalf("Expected teks='HALO DUNIA', got %+v", teksVal)
	}

	// Program 2: Pipeline with extra arguments: "Celang" |> sambung("Hebat")
	prog2 := uast.NewProgram(
		uast.NewFuncDecl("sambung",
			[]uast.Parameter{{Name: "a"}, {Name: "b"}},
			uast.NewBlock(
				uast.NewReturn(uast.NewBinaryExpr(uast.NewBinaryExpr(uast.NewIdentifier("a"), "+", uast.NewStringLiteral(" ")), "+", uast.NewIdentifier("b"))),
			),
			"string",
		),
		uast.NewVarDecl("hasil_sambung",
			uast.NewPipelineExpr(
				uast.NewStringLiteral("Celang"),
				uast.NewCallExpr(uast.NewIdentifier("sambung"), uast.NewStringLiteral("Hebat")),
			),
			false, "string",
		),
	)
	evaluator.Eval(prog2, env)
	sambungVal, ok := env.Get("hasil_sambung")
	if !ok || sambungVal.(*evaluator.String).Value != "Celang Hebat" {
		t.Fatalf("Expected hasil_sambung='Celang Hebat', got %+v", sambungVal)
	}

	// Program 3: Array transformation with anonymous functions
	// [1, 2, 3, 4, 5, 6]
	//     |> saring(fungsi(x) { kembalikan x % 2 == 0 })
	//     |> petakan(fungsi(x) { kembalikan x * 10 })
	prog3 := uast.NewProgram(
		uast.NewVarDecl("angka",
			uast.NewArrayLiteral(
				uast.NewIntLiteral(1),
				uast.NewIntLiteral(2),
				uast.NewIntLiteral(3),
				uast.NewIntLiteral(4),
				uast.NewIntLiteral(5),
				uast.NewIntLiteral(6),
			),
			false, "array",
		),
		uast.NewVarDecl("hasil_angka",
			uast.NewPipelineExpr(
				uast.NewPipelineExpr(
					uast.NewIdentifier("angka"),
					uast.NewCallExpr(
						uast.NewIdentifier("saring"),
						uast.NewFuncExpr(
							[]uast.Parameter{{Name: "x"}},
							uast.NewBlock(
								uast.NewReturn(
									uast.NewBinaryExpr(
										uast.NewBinaryExpr(uast.NewIdentifier("x"), "%", uast.NewIntLiteral(2)),
										"==",
										uast.NewIntLiteral(0),
									),
								),
							),
							"bool",
						),
					),
				),
				uast.NewCallExpr(
					uast.NewIdentifier("petakan"),
					uast.NewFuncExpr(
						[]uast.Parameter{{Name: "x"}},
						uast.NewBlock(
							uast.NewReturn(
								uast.NewBinaryExpr(uast.NewIdentifier("x"), "*", uast.NewIntLiteral(10)),
							),
						),
						"int",
					),
				),
			),
			false, "array",
		),
	)

	evaluator.Eval(prog3, env)
	hasilAngkaVal, ok := env.Get("hasil_angka")
	if !ok {
		t.Fatalf("Expected hasil_angka to exist")
	}
	arr, ok := hasilAngkaVal.(*evaluator.Array)
	if !ok || len(arr.Elements) != 3 {
		t.Fatalf("Expected 3 elements [20, 40, 60], got %v", hasilAngkaVal)
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 20 ||
		arr.Elements[1].(*evaluator.Integer).Value != 40 ||
		arr.Elements[2].(*evaluator.Integer).Value != 60 {
		t.Fatalf("Expected [20, 40, 60], got %v", arr.Inspect())
	}
}

func TestEvaluator_StyleBlockScopingAndExecution(t *testing.T) {
	// Program:
	// buat nama_pengguna = "Budi"
	// gaya "javascript" {
	//     let status_login = true
	//     function cek_akses(user) {
	//         return user == "Budi" && status_login
	//     }
	//     let diizinkan = cek_akses(nama_pengguna)
	// }
	// buat hasil_akses = diizinkan
	env := evaluator.NewEnvironment()
	prog := uast.NewProgram(
		uast.NewVarDecl("nama_pengguna", uast.NewStringLiteral("Budi"), false, "string"),
		uast.NewStyleBlock(
			"javascript",
			uast.NewBlock(
				uast.NewVarDecl("status_login", uast.NewBoolLiteral(true), true, "bool"),
				uast.NewFuncDecl(
					"cek_akses",
					[]uast.Parameter{{Name: "user"}},
					uast.NewBlock(
						uast.NewReturn(
							uast.NewBinaryExpr(
								uast.NewBinaryExpr(uast.NewIdentifier("user"), "==", uast.NewStringLiteral("Budi")),
								"&&",
								uast.NewIdentifier("status_login"),
							),
						),
					),
					"bool",
				),
				uast.NewVarDecl(
					"diizinkan",
					uast.NewCallExpr(uast.NewIdentifier("cek_akses"), uast.NewIdentifier("nama_pengguna")),
					true,
					"bool",
				),
			),
		),
		uast.NewVarDecl("hasil_akses", uast.NewIdentifier("diizinkan"), false, "bool"),
	)

	evaluator.Eval(prog, env)
	hasilVal, ok := env.Get("hasil_akses")
	if !ok || hasilVal.(*evaluator.Boolean).Value != true {
		t.Fatalf("Expected hasil_akses=true from javascript style block, got %+v", hasilVal)
	}

	fnVal, ok := env.Get("cek_akses")
	if !ok || fnVal.Type() != evaluator.FUNCTION_OBJ {
		t.Fatalf("Expected cek_akses function to exist in global environment, got %+v", fnVal)
	}
}

func TestEvaluator_StringInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
			buat nama = "Budi"
			buat umur = 25
			"Halo {nama}, umur {umur} tahun!"
			`,
			expected: "Halo Budi, umur 25 tahun!",
		},
		{
			input: `
			buat harga = 50000
			"Total belanja: Rp {harga * 2}"
			`,
			expected: "Total belanja: Rp 100000",
		},
		{
			input: `
			buat status = "aktif"
			` + "`" + `User status: ${status}` + "`" + `
			`,
			expected: "User status: aktif",
		},
		{
			input: `
			"Teks biasa tanpa kurung"
			`,
			expected: "Teks biasa tanpa kurung",
		},
	}

	for _, tt := range tests {
		l := parser.NewLexer(tt.input)
		p := parser.NewParser(l)
		prog, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parser error: %v", err)
		}
		env := evaluator.NewEnvironment()
		res := evaluator.Eval(prog, env)
		strVal, ok := res.(*evaluator.String)
		if !ok || strVal.Value != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, res.Inspect())
		}
	}
}

func TestEvaluator_Destructuring(t *testing.T) {
	// 1. Array destructuring declaration
	inputArrayDecl := `
	buat [x, y, z] = [10, 20, 30]
	[x, y, z]
	`
	l := parser.NewLexer(inputArrayDecl)
	p := parser.NewParser(l)
	prog, _ := p.ParseProgram()
	env := evaluator.NewEnvironment()
	res := evaluator.Eval(prog, env)
	arr := res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.Integer).Value != 10 ||
		arr.Elements[1].(*evaluator.Integer).Value != 20 ||
		arr.Elements[2].(*evaluator.Integer).Value != 30 {
		t.Fatalf("array destructuring declaration failed: %s", res.Inspect())
	}

	// 2. Map destructuring declaration with aliases
	inputMapDecl := `
	buat profil = { "nama": "Andi", "peran": "Admin" }
	buat { nama, peran: p } = profil
	[nama, p]
	`
	l = parser.NewLexer(inputMapDecl)
	p = parser.NewParser(l)
	prog, _ = p.ParseProgram()
	env = evaluator.NewEnvironment()
	res = evaluator.Eval(prog, env)
	arr = res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.String).Value != "Andi" ||
		arr.Elements[1].(*evaluator.String).Value != "Admin" {
		t.Fatalf("map destructuring declaration failed: %s", res.Inspect())
	}

	// 3. Destructuring assignments
	inputAssign := `
	buat a = 0
	buat b = 0
	[a, b] = [100, 200]
	[a, b]
	`
	l = parser.NewLexer(inputAssign)
	p = parser.NewParser(l)
	prog, _ = p.ParseProgram()
	env = evaluator.NewEnvironment()
	res = evaluator.Eval(prog, env)
	arr = res.(*evaluator.Array)
	if arr.Elements[0].(*evaluator.Integer).Value != 100 ||
		arr.Elements[1].(*evaluator.Integer).Value != 200 {
		t.Fatalf("destructuring assignment failed: %s", res.Inspect())
	}
}

func TestEvaluator_Structs(t *testing.T) {
	input := `
	bentuk RekeningBank {
		pemilik: Teks,
		saldo itu Angka = 0,

		fungsi setor(jumlah: Angka) {
			diri.saldo += jumlah
			kembalikan diri.saldo
		}

		fungsi tarik(jumlah: Angka) {
			jika jumlah > diri.saldo {
				lempar "Saldo tidak cukup!"
			}
			diri.saldo -= jumlah
			kembalikan diri.saldo
		}
	}

	buat akun = RekeningBank("Budi", 50000)
	buat s1 = akun.saldo
	buat s2 = akun.setor(25000)
	buat s3 = akun.tarik(10000)
	akun.saldo += 5000
	buat s4 = akun.saldo

	[s1, s2, s3, s4, akun.pemilik]
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	env := evaluator.NewEnvironment()
	res := evaluator.Eval(prog, env)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array result, got %s", res.Inspect())
	}

	// s1 = 50000
	if arr.Elements[0].(*evaluator.Integer).Value != 50000 {
		t.Errorf("expected initial saldo 50000, got %s", arr.Elements[0].Inspect())
	}
	// s2 = 75000
	if arr.Elements[1].(*evaluator.Integer).Value != 75000 {
		t.Errorf("expected after setor 75000, got %s", arr.Elements[1].Inspect())
	}
	// s3 = 65000
	if arr.Elements[2].(*evaluator.Integer).Value != 65000 {
		t.Errorf("expected after tarik 65000, got %s", arr.Elements[2].Inspect())
	}
	// s4 = 70000 (after akun.saldo += 5000)
	if arr.Elements[3].(*evaluator.Integer).Value != 70000 {
		t.Errorf("expected after direct mutation 70000, got %s", arr.Elements[3].Inspect())
	}
	// pemilik = "Budi"
	if arr.Elements[4].(*evaluator.String).Value != "Budi" {
		t.Errorf("expected pemilik 'Budi', got %s", arr.Elements[4].Inspect())
	}
}




