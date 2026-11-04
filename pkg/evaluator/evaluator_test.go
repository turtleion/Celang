package evaluator_test

import (
	"bytes"
	"math"
	"strings"
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
			name:     "Pure String Number + Int: '100' + 25 -> 125",
			node:     uast.NewBinaryExpr(uast.NewStringLiteral("100"), "+", uast.NewIntLiteral(25)),
			expected: "125",
		},
		{
			name:     "String - Int with extraction: '100' - 20 -> 80",
			node:     uast.NewBinaryExpr(uast.NewStringLiteral("100"), "-", uast.NewIntLiteral(20)),
			expected: "80",
		},
		{
			name:     "RegEx Extraction: 'Rp 50000' / 2 -> 25000",
			node:     uast.NewBinaryExpr(uast.NewStringLiteral("Rp 50000"), "/", uast.NewIntLiteral(2)),
			expected: "25000",
		},
		{
			name:     "String + Number Human Description: 'Total: ' + 50 -> 'Total: 50'",
			node:     uast.NewBinaryExpr(uast.NewStringLiteral("Total: "), "+", uast.NewIntLiteral(50)),
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
		{
			input: `
			buat nama = "Rani"
			buat tabungan = 50000
			f"Halo, ${nama}! Tabunganmu adalah Rp${tabungan}."
			`,
			expected: "Halo, Rani! Tabunganmu adalah Rp50000.",
		},
		{
			input: `
			buat tabungan = 50000
			f"Jika ditambah 10000: Rp${tabungan + 10000}"
			`,
			expected: "Jika ditambah 10000: Rp60000",
		},
		{
			input: `
			buat tabungan = 50000
			f"Jika ditambah 10000: Rp{tabungan + 10000}"
			`,
			expected: "Jika ditambah 10000: Rp60000",
		},
		{
			input: `
			buat nama = "Dedi"
			f'Nama: {nama}'
			`,
			expected: "Nama: Dedi",
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

func TestAudit_EdgeCasesAndCollisions(t *testing.T) {
	// 1. Test Map with key 'awal' or 'sebelumnya'
	inputMap := `
	buat data = { "awal": 99, "sebelumnya": 88 }
	[data.awal, data.sebelumnya]
	`
	l := parser.NewLexer(inputMap)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)
	res := evaluator.Eval(prog, env)
	t.Logf("Result of map with 'awal'/'sebelumnya': %v", res.Inspect())

	// 2. Test Struct destructuring
	inputStructDestruct := `
	bentuk Point {
		x itu Angka,
		y itu Angka
	}
	buat pt = Point(10, 20)
	buat { x, y } = pt
	[x, y]
	`
	l2 := parser.NewLexer(inputStructDestruct)
	p2 := parser.NewParser(l2)
	prog2, err2 := p2.ParseProgram()
	if err2 != nil {
		t.Fatalf("parser error: %v", err2)
	}
	res2 := evaluator.Eval(prog2, env)
	t.Logf("Result of struct destructuring: %v", res2.Inspect())

	jsonProg, _ := parser.NewParser(parser.NewLexer(`json.format(pt)`)).ParseProgram()
	jsonRes := evaluator.Eval(jsonProg, env)
	t.Logf("Result of json.format(pt): %v", jsonRes.Inspect())

	// 3. Test Destructuring Assignment without buat: { a, b } = obj
	inputDestructAssign := `
	buat a = 1
	buat b = 2
	{ a, b } = { "a": 100, "b": 200 }
	[a, b]
	`
	l3 := parser.NewLexer(inputDestructAssign)
	p3 := parser.NewParser(l3)
	prog3, err3 := p3.ParseProgram()
	if err3 != nil {
		t.Logf("Parser error on { a, b } = obj: %v", err3)
	} else {
		res3 := evaluator.Eval(prog3, env)
		t.Logf("Result of { a, b } = obj: %v", res3.Inspect())
	}

	// 4. Test Circular Struct Inspect
	inputCircular := `
	bentuk Node {
		buat val,
		buat next
	}
	buat n = Node(1, null)
	n.next = n
	`
	l4 := parser.NewLexer(inputCircular)
	p4 := parser.NewParser(l4)
	prog4, err4 := p4.ParseProgram()
	if err4 != nil {
		t.Fatalf("parser error: %v", err4)
	}
	res4 := evaluator.Eval(prog4, env)
	nObj, ok := env.Get("n")
	if ok {
		t.Logf("n found, type: %s", nObj.Type())
		inspectStr := nObj.Inspect()
		t.Logf("Circular struct inspect result: %s", inspectStr)
		if !strings.Contains(inspectStr, "Node{...}") {
			t.Errorf("Expected cycle detection in Inspect, got: %s", inspectStr)
		}

		// Also verify json.format handles cycle cleanly with error
		circJSONProg, _ := parser.NewParser(parser.NewLexer(`json.format(n)`)).ParseProgram()
		circJSONRes := evaluator.Eval(circJSONProg, env)
		t.Logf("Result of json.format on circular struct: %s", circJSONRes.Inspect())
		if circErr, ok := circJSONRes.(*evaluator.Error); !ok || !strings.Contains(circErr.Message, "sirkular") {
			t.Errorf("Expected circular error from json.format, got: %s", circJSONRes.Inspect())
		}
	}
	_ = res4
}

func TestEvaluator_InputBuiltin(t *testing.T) {
	inBuf := strings.NewReader("Kayasa\n")
	var outBuf bytes.Buffer

	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltinsWithIO(env, inBuf, &outBuf)

	inputCode := `
	buat nama = masukan("Nama kamu? ")
	buat salam = "Halo " + nama
	salam
	`
	l := parser.NewLexer(inputCode)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	res := evaluator.Eval(prog, env)
	str, ok := res.(*evaluator.String)
	if !ok || str.Value != "Halo Kayasa" {
		t.Fatalf("Expected 'Halo Kayasa', got %v", res.Inspect())
	}

	if outBuf.String() != "Nama kamu? " {
		t.Fatalf("Expected prompt 'Nama kamu? ', got %q", outBuf.String())
	}
}

func TestEvaluator_TimeTravelMemoryBounded(t *testing.T) {
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)

	// Mutate variable 1,000 times in a loop
	code := `
	buat hitungan = 0
	selama hitungan < 1000 {
		hitungan = hitungan + 1
	}
	[hitungan.awal, hitungan.sebelumnya, hitungan, hitungan.histori]
	`
	l := parser.NewLexer(code)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	res := evaluator.Eval(prog, env)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 4 {
		t.Fatalf("Expected 4 elements, got %v", res.Inspect())
	}

	// 1. Initial value (hitungan.awal) must stay 0
	initVal, ok := arr.Elements[0].(*evaluator.Integer)
	if !ok || initVal.Value != 0 {
		t.Errorf("Expected hitungan.awal = 0, got %v", arr.Elements[0].Inspect())
	}

	// 2. Previous value (hitungan.sebelumnya) must be 999
	prevVal, ok := arr.Elements[1].(*evaluator.Integer)
	if !ok || prevVal.Value != 999 {
		t.Errorf("Expected hitungan.sebelumnya = 999, got %v", arr.Elements[1].Inspect())
	}

	// 3. Current value must be 1000
	currVal, ok := arr.Elements[2].(*evaluator.Integer)
	if !ok || currVal.Value != 1000 {
		t.Errorf("Expected hitungan = 1000, got %v", arr.Elements[2].Inspect())
	}

	// 4. History size must be bounded by DefaultHistoryCapacity (50)
	histArr, ok := arr.Elements[3].(*evaluator.Array)
	if !ok {
		t.Fatalf("Expected array for hitungan.histori, got %v", arr.Elements[3].Inspect())
	}
	if len(histArr.Elements) != evaluator.DefaultHistoryCapacity {
		t.Errorf("Expected history bounded to %d elements, got %d", evaluator.DefaultHistoryCapacity, len(histArr.Elements))
	}

	// 5. Verify the sliding window content: last element must be 1000, first element must be 1000 - 50 + 1 = 951
	firstHistVal := histArr.Elements[0].(*evaluator.Integer).Value
	lastHistVal := histArr.Elements[len(histArr.Elements)-1].(*evaluator.Integer).Value
	if lastHistVal != 1000 || firstHistVal != 951 {
		t.Errorf("Expected history window from 951 to 1000, got %d to %d", firstHistVal, lastHistVal)
	}
}

func testEvalSource(t *testing.T, input string) evaluator.Object {
	t.Helper()
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}
	env := evaluator.NewEnvironment()
	evaluator.RegisterBuiltins(env, nil)
	return evaluator.Eval(prog, env)
}

func TestEvaluator_StructVisibilityAndAtShorthand(t *testing.T) {
	// 1. Test internal access with @ and diri, and block external access to privat field
	inputSuccess := `
	bentuk Dompet {
		pemilik: Teks,
		privat saldo: Angka = 1000,
		privat {
			pin: Angka = 1234,
			cekPin(p: Angka) {
				kembalikan p == @pin
			}
		}

		setor(jumlah: Angka) {
			@saldo += jumlah
			kembalikan @saldo
		}

		ambilSaldo(p: Angka) {
			jika @cekPin(p) {
				kembalikan @saldo
			}
			kembalikan -1
		}
	}

	buat d = Dompet("Andi")
	buat s1 = d.setor(500)
	buat s2 = d.ambilSaldo(1234)
	buat s3 = d.ambilSaldo(0)
	[s1, s2, s3, d.pemilik]
	`
	res := testEvalSource(t, inputSuccess)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 1500 {
		t.Errorf("expected s1 = 1500, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 1500 {
		t.Errorf("expected s2 = 1500, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != -1 {
		t.Errorf("expected s3 = -1, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.String).Value != "Andi" {
		t.Errorf("expected d.pemilik = Andi, got %v", arr.Elements[3])
	}

	// 2. Direct external read on private field should fail
	inputFailRead := `
	bentuk Dompet {
		privat saldo: Angka = 1000
	}
	buat d = Dompet()
	d.saldo
	`
	errObj := testEvalSource(t, inputFailRead)
	if errObj.Type() != evaluator.ERROR_OBJ {
		t.Errorf("expected error on external read of privat saldo, got %s", errObj.Inspect())
	} else {
		errVal := errObj.(*evaluator.Error)
		if !strings.Contains(errVal.Message, "privat") {
			t.Errorf("expected error message mentioning 'privat', got %s", errVal.Message)
		}
	}

	// 3. Direct external write on private field should fail
	inputFailWrite := `
	bentuk Dompet {
		privat saldo: Angka = 1000
	}
	buat d = Dompet()
	d.saldo = 9999
	`
	errObjWrite := testEvalSource(t, inputFailWrite)
	if errObjWrite.Type() != evaluator.ERROR_OBJ {
		t.Errorf("expected error on external write of privat saldo, got %s", errObjWrite.Inspect())
	}

	// 4. Direct external call to private method should fail
	inputFailMethod := `
	bentuk Kunci {
		privat rahasia() {
			kembalikan 42
		}
	}
	buat k = Kunci()
	k.rahasia()
	`
	errObjMethod := testEvalSource(t, inputFailMethod)
	if errObjMethod.Type() != evaluator.ERROR_OBJ {
		t.Errorf("expected error on external call of private method, got %s", errObjMethod.Inspect())
	}
}

func TestEvaluator_ReferenceTransparent(t *testing.T) {
	// 1. Primitive reference binding
	input := `
	buat a = 10
	rujuk b = a
	b = 25
	buat c = a
	a = 50
	buat d = b
	[c, d, b.sebelumnya]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 25 {
		t.Errorf("expected c = 25, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 50 {
		t.Errorf("expected d = 50, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 25 {
		t.Errorf("expected b.sebelumnya = 25, got %v", arr.Elements[2])
	}

	// 2. Function parameter reference
	inputFn := `
	fungsi tambahLima(rujuk x: Angka) {
		x += 5
	}
	buat skor = 100
	tambahLima(skor)
	skor
	`
	resFn := testEvalSource(t, inputFn)
	intVal, ok := resFn.(*evaluator.Integer)
	if !ok || intVal.Value != 105 {
		t.Errorf("expected skor = 105, got %s", resFn.Inspect())
	}
}

func TestEvaluator_StructEmbedding(t *testing.T) {
	input := `
	bentuk Identitas {
		nama: Teks,
		umur: Angka,
		sapa() {
			kembalikan "Halo, aku " + @nama
		}
	}

	bentuk Admin {
		Identitas,
		level: Angka = 1,
		sapaAdmin() {
			kembalikan @sapa() + " (Level " + Teks(@level) + ")"
		}
	}

	buat adm = Admin("Budi", 30, 99)
	[adm.nama, adm.umur, adm.level, adm.sapa(), adm.sapaAdmin()]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.String).Value != "Budi" {
		t.Errorf("expected adm.nama = Budi, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 30 {
		t.Errorf("expected adm.umur = 30, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 99 {
		t.Errorf("expected adm.level = 99, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.String).Value != "Halo, aku Budi" {
		t.Errorf("expected adm.sapa() = 'Halo, aku Budi', got %v", arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.String).Value != "Halo, aku Budi (Level 99)" {
		t.Errorf("expected adm.sapaAdmin() = 'Halo, aku Budi (Level 99)', got %v", arr.Elements[4])
	}
}

func TestEvaluator_TypeAliasAndCasting(t *testing.T) {
	input := `
	tipe ID = Angka

	buat a = "42" sebagai Angka
	buat b = 3.99 sebagai Angka
	buat c = 100 sebagai Teks
	buat d = "99" sebagai ID
	buat e = ID("123")
	buat f = Angka("500")
	buat g = Teks(777)
	buat h = Logika(1)

	[a, b, c, d, e, f, g, h]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 42 {
		t.Errorf("expected a = 42, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 3 {
		t.Errorf("expected b = 3, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.String).Value != "100" {
		t.Errorf("expected c = '100', got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.Integer).Value != 99 {
		t.Errorf("expected d = 99, got %v", arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.Integer).Value != 123 {
		t.Errorf("expected e = 123, got %v", arr.Elements[4])
	}
	if arr.Elements[5].(*evaluator.Integer).Value != 500 {
		t.Errorf("expected f = 500, got %v", arr.Elements[5])
	}
	if arr.Elements[6].(*evaluator.String).Value != "777" {
		t.Errorf("expected g = '777', got %v", arr.Elements[6])
	}
	if !arr.Elements[7].(*evaluator.Boolean).Value {
		t.Errorf("expected h = true, got %v", arr.Elements[7])
	}
}

func TestEvaluator_DeepClone(t *testing.T) {
	input := `
	bentuk Profil {
		nama: Teks,
		skor: Angka
	}

	buat p1 = Profil("Dewi", 100)
	buat p2 = salin(p1)
	p2.nama = "Sari"
	p2.skor = 200

	buat arr1 = [1, 2, [3, 4]]
	buat arr2 = clone(arr1)
	arr2[0] = 99

	[p1.nama, p1.skor, p2.nama, p2.skor, arr1[0], arr2[0]]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.String).Value != "Dewi" {
		t.Errorf("expected p1.nama = Dewi, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected p1.skor = 100, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.String).Value != "Sari" {
		t.Errorf("expected p2.nama = Sari, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.Integer).Value != 200 {
		t.Errorf("expected p2.skor = 200, got %v", arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.Integer).Value != 1 {
		t.Errorf("expected arr1[0] = 1, got %v", arr.Elements[4])
	}
	if arr.Elements[5].(*evaluator.Integer).Value != 99 {
		t.Errorf("expected arr2[0] = 99, got %v", arr.Elements[5])
	}
}

func TestEvaluator_ErgonomicFeatures(t *testing.T) {
	// 1. Safe Navigation (mungkin) & Null Coalesce (kalau kosong / ??)
	inputSafeNav := `
	bentuk Profil {
		kota: Teks
	}
	bentuk Pengguna {
		profil: Profil
	}

	buat p1 = kosong
	buat p2 = Pengguna(Profil("Bandung"))
	buat p3 = Pengguna(kosong)

	buat kota1 = mungkin p1.profil.kota
	buat kota2 = mungkin p2.profil.kota
	buat kota3 = mungkin p3.profil.kota

	buat fb1 = kota1 kalau kosong "Default"
	buat fb2 = kota2 kalau kosong "Default"
	buat fb3 = kota3 ?? "Default"
	buat fbFalse = salah kalau kosong benar
	buat fbZero = 0 kalau kosong 99

	[kota1 == kosong, kota2, kota3 == kosong, fb1, fb2, fb3, fbFalse, fbZero]
	`
	res1 := testEvalSource(t, inputSafeNav)
	arr1, ok := res1.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res1.Inspect())
	}
	if !arr1.Elements[0].(*evaluator.Boolean).Value {
		t.Errorf("expected kota1 == kosong, got %v", arr1.Elements[0])
	}
	if arr1.Elements[1].(*evaluator.String).Value != "Bandung" {
		t.Errorf("expected kota2 = Bandung, got %v", arr1.Elements[1])
	}
	if !arr1.Elements[2].(*evaluator.Boolean).Value {
		t.Errorf("expected kota3 == kosong, got %v", arr1.Elements[2])
	}
	if arr1.Elements[3].(*evaluator.String).Value != "Default" {
		t.Errorf("expected fb1 = Default, got %v", arr1.Elements[3])
	}
	if arr1.Elements[4].(*evaluator.String).Value != "Bandung" {
		t.Errorf("expected fb2 = Bandung, got %v", arr1.Elements[4])
	}
	if arr1.Elements[5].(*evaluator.String).Value != "Default" {
		t.Errorf("expected fb3 = Default, got %v", arr1.Elements[5])
	}
	if arr1.Elements[6].(*evaluator.Boolean).Value != false {
		t.Errorf("expected fbFalse = false, got %v", arr1.Elements[6])
	}
	if arr1.Elements[7].(*evaluator.Integer).Value != 0 {
		t.Errorf("expected fbZero = 0, got %v", arr1.Elements[7])
	}

	// 2. Defer / saat_selesai (LIFO order & early return)
	inputDefer := `
	buat log = []

	fungsi tes_defer() {
		saat_selesai log.tambah("akhir_1")
		saat_selesai log.tambah("akhir_2")
		log.tambah("tengah")
		kembalikan 42
	}

	buat hasil = tes_defer()
	[hasil, log]
	`
	res2 := testEvalSource(t, inputDefer)
	arr2, ok := res2.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res2.Inspect())
	}
	if arr2.Elements[0].(*evaluator.Integer).Value != 42 {
		t.Errorf("expected hasil = 42, got %v", arr2.Elements[0])
	}
	logArr := arr2.Elements[1].(*evaluator.Array)
	if len(logArr.Elements) != 3 {
		t.Fatalf("expected 3 log elements, got %d", len(logArr.Elements))
	}
	// LIFO: "tengah", then "akhir_2", then "akhir_1"
	if logArr.Elements[0].(*evaluator.String).Value != "tengah" {
		t.Errorf("expected log[0] = tengah, got %v", logArr.Elements[0])
	}
	if logArr.Elements[1].(*evaluator.String).Value != "akhir_2" {
		t.Errorf("expected log[1] = akhir_2, got %v", logArr.Elements[1])
	}
	if logArr.Elements[2].(*evaluator.String).Value != "akhir_1" {
		t.Errorf("expected log[2] = akhir_1, got %v", logArr.Elements[2])
	}

	// 2b. Defer inside loops and try-catch
	inputDeferLoop := `
	buat logLoop = []
	fungsi tes_defer_loop() {
		untuk i dalam [1] {
			saat_selesai logLoop.tambah("defer_loop")
		}
		coba {
			lempar "err"
		} tangkap e {
			saat_selesai logLoop.tambah("defer_catch")
		}
		logLoop.tambah("fungsi_selesai")
	}
	tes_defer_loop()
	logLoop
	`
	resDeferLoop := testEvalSource(t, inputDeferLoop)
	arrDeferLoop, ok := resDeferLoop.(*evaluator.Array)
	if !ok || len(arrDeferLoop.Elements) != 3 {
		t.Fatalf("expected 3 log elements from defer loop/catch, got: %v", resDeferLoop)
	}
	if arrDeferLoop.Elements[0].(*evaluator.String).Value != "fungsi_selesai" {
		t.Errorf("expected first element 'fungsi_selesai', got %v", arrDeferLoop.Elements[0])
	}
	if arrDeferLoop.Elements[1].(*evaluator.String).Value != "defer_catch" {
		t.Errorf("expected second element 'defer_catch', got %v", arrDeferLoop.Elements[1])
	}
	if arrDeferLoop.Elements[2].(*evaluator.String).Value != "defer_loop" {
		t.Errorf("expected third element 'defer_loop', got %v", arrDeferLoop.Elements[2])
	}

	// 3. Range & Slicing (1..5, arr[1..3], str[0..2])
	inputRange := `
	buat r = 1..5

	buat total = 0
	untuk x dalam 1..5 {
		total += x
	}

	buat data = [10, 20, 30, 40, 50]
	buat s1 = data[1..3]
	buat s2 = data[..2]
	buat s3 = data[3..]

	buat kata = "Celang"
	buat sub_kata = kata[0..2]

	[r, total, s1, s2, s3, sub_kata]
	`
	res3 := testEvalSource(t, inputRange)
	arr3, ok := res3.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res3.Inspect())
	}
	// r = [1, 2, 3, 4, 5]
	rArr := arr3.Elements[0].(*evaluator.Array)
	if len(rArr.Elements) != 5 || rArr.Elements[4].(*evaluator.Integer).Value != 5 {
		t.Errorf("expected r = [1, 2, 3, 4, 5], got %s", rArr.Inspect())
	}
	// total = 15
	if arr3.Elements[1].(*evaluator.Integer).Value != 15 {
		t.Errorf("expected total = 15, got %v", arr3.Elements[1])
	}
	// s1 = [20, 30]
	s1Arr := arr3.Elements[2].(*evaluator.Array)
	if len(s1Arr.Elements) != 2 || s1Arr.Elements[0].(*evaluator.Integer).Value != 20 || s1Arr.Elements[1].(*evaluator.Integer).Value != 30 {
		t.Errorf("expected s1 = [20, 30], got %s", s1Arr.Inspect())
	}
	// s2 = [10, 20]
	s2Arr := arr3.Elements[3].(*evaluator.Array)
	if len(s2Arr.Elements) != 2 || s2Arr.Elements[0].(*evaluator.Integer).Value != 10 || s2Arr.Elements[1].(*evaluator.Integer).Value != 20 {
		t.Errorf("expected s2 = [10, 20], got %s", s2Arr.Inspect())
	}
	// s3 = [40, 50]
	s3Arr := arr3.Elements[4].(*evaluator.Array)
	if len(s3Arr.Elements) != 2 || s3Arr.Elements[0].(*evaluator.Integer).Value != 40 || s3Arr.Elements[1].(*evaluator.Integer).Value != 50 {
		t.Errorf("expected s3 = [40, 50], got %s", s3Arr.Inspect())
	}
	// sub_kata = "Ce"
	if arr3.Elements[5].(*evaluator.String).Value != "Ce" {
		t.Errorf("expected sub_kata = Ce, got %v", arr3.Elements[5])
	}
}

func TestEvaluator_ArrowFunctions(t *testing.T) {
	// 1. Basic arrow functions
	input := `
	buat tambah = (a, b) => a + b
	buat kuadrat = x => x * x
	buat salam = () => "halo celang"
	[tambah(10, 5), kuadrat(6), salam()]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 15 {
		t.Errorf("expected 15, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 36 {
		t.Errorf("expected 36, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.String).Value != "halo celang" {
		t.Errorf("expected 'halo celang', got %v", arr.Elements[2])
	}

	// 2. Higher order functions with arrow and pipeline
	inputPipeline := `
	buat angka = [1, 2, 3, 4]
	buat genap = angka |> saring(x => x % 2 == 0)
	buat kaliDua = genap |> petakan(x => x * 2)
	kaliDua
	`
	res2 := testEvalSource(t, inputPipeline)
	arr2, ok := res2.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res2.Inspect())
	}
	if len(arr2.Elements) != 2 || arr2.Elements[0].(*evaluator.Integer).Value != 4 || arr2.Elements[1].(*evaluator.Integer).Value != 8 {
		t.Errorf("expected [4, 8], got %s", arr2.Inspect())
	}
}

func TestEvaluator_TitipJemput(t *testing.T) {
	// 1. Basic titip & jemput
	input := `
	fungsi hitungLambat(x) {
		kembalikan x * 10
	}
	buat tiket = titip hitungLambat(5)
	buat hasil = jemput tiket
	hasil
	`
	res := testEvalSource(t, input)
	if intVal, ok := res.(*evaluator.Integer); !ok || intVal.Value != 50 {
		t.Fatalf("expected 50, got %s", res.Inspect())
	}

	// 1b. titip block statement
	inputBlock := `
	buat tiket = titip {
		buat a = 15
		buat b = 25
		kembalikan a + b
	}
	jemput tiket
	`
	resBlock := testEvalSource(t, inputBlock)
	if intVal, ok := resBlock.(*evaluator.Integer); !ok || intVal.Value != 40 {
		t.Fatalf("expected 40 from titip block, got %s", resBlock.Inspect())
	}

	// 2. jemput dengan batas waktu (timeout)
	inputTimeout := `
	buat tiket = titip 42
	buat hasil = jemput tiket dalam 1
	hasil
	`
	resTimeout := testEvalSource(t, inputTimeout)
	if intVal, ok := resTimeout.(*evaluator.Integer); !ok || intVal.Value != 42 {
		t.Fatalf("expected 42, got %s", resTimeout.Inspect())
	}

	// 3. jemput_semua
	inputSemua := `
	buat t1 = titip 10
	buat t2 = titip 20
	buat hasil = jemput_semua([t1, t2])
	hasil
	`
	resSemua := testEvalSource(t, inputSemua)
	arrSemua, ok := resSemua.(*evaluator.Array)
	if !ok || len(arrSemua.Elements) != 2 {
		t.Fatalf("expected array of 2 elements, got %s", resSemua.Inspect())
	}
	if arrSemua.Elements[0].(*evaluator.Integer).Value != 10 || arrSemua.Elements[1].(*evaluator.Integer).Value != 20 {
		t.Errorf("expected [10, 20], got %s", arrSemua.Inspect())
	}

	// 4. jemput_pertama
	inputPertama := `
	buat t1 = titip "cepat"
	buat hasil = jemput_pertama([t1])
	hasil
	`
	resPertama := testEvalSource(t, inputPertama)
	if strVal, ok := resPertama.(*evaluator.String); !ok || strVal.Value != "cepat" {
		t.Fatalf("expected 'cepat', got %s", resPertama.Inspect())
	}
}

func TestEvaluator_Saluran(t *testing.T) {
	// 1. Saluran buffered send & receive
	input := `
	buat ch = Saluran(2)
	ch.kirim(100)
	ch.kirim(200)
	buat a = ch.terima()
	buat b = ch.terima()
	[a, b]
	`
	res := testEvalSource(t, input)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("expected array of 2, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 100 || arr.Elements[1].(*evaluator.Integer).Value != 200 {
		t.Errorf("expected [100, 200], got %s", arr.Inspect())
	}

	// 2. Loop streaming: untuk x dalam ch
	inputLoop := `
	buat ch = Saluran(3)
	ch.kirim("satu")
	ch.kirim("dua")
	ch.tutup()
	buat hasil = []
	untuk item dalam ch {
		hasil.tambah(item)
	}
	hasil
	`
	resLoop := testEvalSource(t, inputLoop)
	arrLoop, ok := resLoop.(*evaluator.Array)
	if !ok || len(arrLoop.Elements) != 2 {
		t.Fatalf("expected [satu, dua], got %s", resLoop.Inspect())
	}
	if arrLoop.Elements[0].(*evaluator.String).Value != "satu" || arrLoop.Elements[1].(*evaluator.String).Value != "dua" {
		t.Errorf("expected ['satu', 'dua'], got %s", arrLoop.Inspect())
	}

	// 3. Typed Channel constraint
	inputTyped := `
	buat ch = Saluran<Angka>(1)
	ch.kirim(999)
	ch.terima()
	`
	resTyped := testEvalSource(t, inputTyped)
	if intVal, ok := resTyped.(*evaluator.Integer); !ok || intVal.Value != 999 {
		t.Fatalf("expected 999, got %s", resTyped.Inspect())
	}

	inputTypedErr := `
	buat ch = Saluran<Angka>(1)
	ch.kirim("bukan angka")
	`
	resTypedErr := testEvalSource(t, inputTypedErr)
	if resTypedErr.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("expected error sending string to Saluran<Angka>, got %s", resTypedErr.Inspect())
	}
}

func TestEvaluator_Generics(t *testing.T) {
	// 1. Struct generic
	inputStruct := `
	bentuk Kotak<T> {
		isi: T
	}
	buat k1 = Kotak<Angka>(123)
	buat k2 = Kotak<Teks>("rahasia")
	[k1.isi, k2.isi]
	`
	res := testEvalSource(t, inputStruct)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 123 {
		t.Errorf("expected 123, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.String).Value != "rahasia" {
		t.Errorf("expected 'rahasia', got %v", arr.Elements[1])
	}

	// 2. Struct with type constraints: Nilai: Angka and Nilai harus Angka
	inputConstraint := `
	bentuk Dompet<Nilai: Angka> {
		saldo: Nilai
	}
	bentuk Celengan<Nilai harus Angka> {
		tabungan: Nilai
	}
	buat d = Dompet<Angka>(50000)
	buat c = Celengan<Angka>(100000)
	[d.saldo, c.tabungan]
	`
	resConst := testEvalSource(t, inputConstraint)
	arrConst, ok := resConst.(*evaluator.Array)
	if !ok || len(arrConst.Elements) != 2 {
		t.Fatalf("expected array of 2, got %s", resConst.Inspect())
	}
	if arrConst.Elements[0].(*evaluator.Integer).Value != 50000 || arrConst.Elements[1].(*evaluator.Integer).Value != 100000 {
		t.Errorf("expected [50000, 100000], got %s", arrConst.Inspect())
	}

	// 3. Struct type constraint violation returns error
	inputErr := `
	bentuk Dompet<Nilai: Angka> {
		saldo: Nilai
	}
	buat d = Dompet<Teks>("uang palsu")
	`
	resErr := testEvalSource(t, inputErr)
	if resErr.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("expected error when type constraint violated, got %s", resErr.Inspect())
	}

	// 4. Generic function
	inputFunc := `
	fungsi identitas<T>(x: T) {
		kembalikan x
	}
	buat a = identitas<Angka>(77)
	buat b = identitas<Teks>("selamat")
	[a, b]
	`
	resFunc := testEvalSource(t, inputFunc)
	arrFunc, ok := resFunc.(*evaluator.Array)
	if !ok || len(arrFunc.Elements) != 2 {
		t.Fatalf("expected array of 2, got %s", resFunc.Inspect())
	}
	if arrFunc.Elements[0].(*evaluator.Integer).Value != 77 || arrFunc.Elements[1].(*evaluator.String).Value != "selamat" {
		t.Errorf("expected [77, 'selamat'], got %s", arrFunc.Inspect())
	}
}

func TestBug_ComputedAssignment(t *testing.T) {
	// 1. Map with variable key assignment: data[k] = 123 and data[k] += 10
	inputMap := `
	buat data = {}
	buat k = "kunci_rahasia"
	data[k] = 123
	data[k] += 10
	[data[k], data["kunci_rahasia"]]
	`
	resMap := testEvalSource(t, inputMap)
	if resMap.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("resMap error: %s", resMap.Inspect())
	}
	arrMap, ok := resMap.(*evaluator.Array)
	if !ok || len(arrMap.Elements) != 2 || arrMap.Elements[0].(*evaluator.Integer).Value != 133 || arrMap.Elements[1].(*evaluator.Integer).Value != 133 {
		t.Fatalf("expected [133, 133], got %s", resMap.Inspect())
	}

	// 2. Array with computed index expression: arr[1 + 1] = 99 and arr[i + 1] += 1
	inputArr := `
	buat arr = [10, 20, 30]
	buat i = 1
	arr[i + 1] = 99
	arr[i + 1] += 1
	[arr[2], arr[i + 1]]
	`
	resArr := testEvalSource(t, inputArr)
	if resArr.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("resArr error: %s", resArr.Inspect())
	}
	arrRes, ok := resArr.(*evaluator.Array)
	if !ok || len(arrRes.Elements) != 2 || arrRes.Elements[0].(*evaluator.Integer).Value != 100 || arrRes.Elements[1].(*evaluator.Integer).Value != 100 {
		t.Fatalf("expected [100, 100], got %s", resArr.Inspect())
	}
}

func TestBug_AuditFixes(t *testing.T) {
	// 1. Float modulo
	resMod := testEvalSource(t, `5.5 % 2`)
	if flt, ok := resMod.(*evaluator.Float); !ok || flt.Value != 1.5 {
		t.Fatalf("expected 5.5 %% 2 = 1.5, got %s", resMod.Inspect())
	}

	// 2. Short-circuiting for && and ||
	inputShortCircuit := `
	buat x = null
	buat hasil_and = 0
	jika x != null dan x.properti_pasti_error > 0 {
		hasil_and = 1
	} selain {
		hasil_and = 2
	}

	buat hasil_or = 0
	jika 1 == 1 atau x.properti_pasti_error > 0 {
		hasil_or = 10
	}

	[hasil_and, hasil_or]
	`
	resSC := testEvalSource(t, inputShortCircuit)
	if resSC.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("resSC error: %s", resSC.Inspect())
	}
	arrSC := resSC.(*evaluator.Array)
	if arrSC.Elements[0].(*evaluator.Integer).Value != 2 || arrSC.Elements[1].(*evaluator.Integer).Value != 10 {
		t.Fatalf("expected [2, 10], got %s", resSC.Inspect())
	}

	// 3. Safe navigation with computed key: data?.[k]
	inputSafeNav := `
	buat data = {"pesan": "halo dunia"}
	buat k = "pesan"
	data?.[k]
	`
	resSafe := testEvalSource(t, inputSafeNav)
	if s, ok := resSafe.(*evaluator.String); !ok || s.Value != "halo dunia" {
		t.Fatalf("expected 'halo dunia', got %s", resSafe.Inspect())
	}

	// 4. Unicode panjang()
	resPanjang := testEvalSource(t, `panjang("halo 👋")`)
	if i, ok := resPanjang.(*evaluator.Integer); !ok || i.Value != 6 {
		t.Fatalf("expected panjang('halo 👋') = 6, got %s", resPanjang.Inspect())
	}

	// 5. deepClone preserving Map Order and iteration
	inputClone := `
	buat m = {"z": 1, "a": 2, "m": 3}
	buat c = salin(m)
	buat kunci = []
	untuk k dalam c {
		kunci.tambah(k)
	}
	kunci
	`
	resClone := testEvalSource(t, inputClone)
	if resClone.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("resClone error: %s", resClone.Inspect())
	}
	arrClone := resClone.(*evaluator.Array)
	if len(arrClone.Elements) != 3 || arrClone.Elements[0].(*evaluator.String).Value != "z" ||
		arrClone.Elements[1].(*evaluator.String).Value != "a" ||
		arrClone.Elements[2].(*evaluator.String).Value != "m" {
		t.Fatalf("expected ['z', 'a', 'm'], got %s", resClone.Inspect())
	}
}

func TestFeature_SatuanBudayaDanTerbilang(t *testing.T) {
	input := `
	buat harga = 50.ribu
	buat gaji = 2.5.juta
	buat hutang = 1.miliar
	buat ram = 16.gb

	buat h_rp = harga.rupiah
	buat h_tb = harga.terbilang
	buat g_tb = gaji.terbilang
	buat seribu_tb = 1000.terbilang
	buat angka_kata = 105.kata
	buat fn_tb = terbilang(15)
	buat fn_rp = format_rupiah(1250000)

	[harga, gaji, hutang, ram, h_rp, h_tb, g_tb, seribu_tb, angka_kata, fn_tb, fn_rp]
	`
	res := testEvalSource(t, input)
	if res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("res error: %s", res.Inspect())
	}
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 11 {
		t.Fatalf("expected 11 elements, got %s", res.Inspect())
	}

	if arr.Elements[0].(*evaluator.Integer).Value != 50000 {
		t.Errorf("50.ribu expected 50000, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 2500000 {
		t.Errorf("2.5.juta expected 2500000, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 1000000000 {
		t.Errorf("1.miliar expected 1000000000, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.Integer).Value != 16*1024*1024*1024 {
		t.Errorf("16.gb expected %d, got %v", int64(16*1024*1024*1024), arr.Elements[3])
	}
	if arr.Elements[4].(*evaluator.String).Value != "Rp 50.000" {
		t.Errorf("harga.rupiah expected 'Rp 50.000', got %v", arr.Elements[4])
	}
	if arr.Elements[5].(*evaluator.String).Value != "lima puluh ribu rupiah" {
		t.Errorf("harga.terbilang expected 'lima puluh ribu rupiah', got %v", arr.Elements[5])
	}
	if arr.Elements[6].(*evaluator.String).Value != "dua juta lima ratus ribu rupiah" {
		t.Errorf("gaji.terbilang expected 'dua juta lima ratus ribu rupiah', got %v", arr.Elements[6])
	}
	if arr.Elements[7].(*evaluator.String).Value != "seribu rupiah" {
		t.Errorf("1000.terbilang expected 'seribu rupiah', got %v", arr.Elements[7])
	}
	if arr.Elements[8].(*evaluator.String).Value != "seratus lima" {
		t.Errorf("105.kata expected 'seratus lima', got %v", arr.Elements[8])
	}
	if arr.Elements[9].(*evaluator.String).Value != "lima belas rupiah" {
		t.Errorf("terbilang(15) expected 'lima belas rupiah', got %v", arr.Elements[9])
	}
	if arr.Elements[10].(*evaluator.String).Value != "Rp 1.250.000" {
		t.Errorf("format_rupiah(1250000) expected 'Rp 1.250.000', got %v", arr.Elements[10])
	}
}

func TestFeature_TransaksiDanRollback(t *testing.T) {
	// 1. Sukses commit transaksi biasa
	codeCommit := `
	buat saldo = 100
	transaksi {
		saldo = 80
	}
	saldo
	`
	resCommit := testEvalSource(t, codeCommit)
	if intVal, ok := resCommit.(*evaluator.Integer); !ok || intVal.Value != 80 {
		t.Fatalf("expected saldo 80 after commit, got: %v", resCommit)
	}

	// 2. Rollback via 'batalkan' eksplisit dan tangkap di pulihkan (err)
	codeRollback := `
	buat saldo = 100
	buat status = "aman"
	transaksi {
		saldo = 50
		saldo = 20
		batalkan "dana tidak mencukupi"
		saldo = 0
	} pulihkan (alasan) {
		status = alasan
	}
	[saldo, status]
	`
	resRollback := testEvalSource(t, codeRollback)
	arrRollback, ok := resRollback.(*evaluator.Array)
	if !ok || len(arrRollback.Elements) != 2 {
		t.Fatalf("expected array [saldo, status], got: %v", resRollback)
	}
	if arrRollback.Elements[0].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected saldo reverted to 100, got: %v", arrRollback.Elements[0])
	}
	if arrRollback.Elements[1].(*evaluator.String).Value != "dana tidak mencukupi" {
		t.Errorf("expected status 'dana tidak mencukupi', got: %v", arrRollback.Elements[1])
	}

	// 3. Rollback otomatis akibat runtime error / lempar
	codeError := `
	buat saldo = 500
	buat err_msg = ""
	transaksi {
		saldo = 200
		lempar "koneksi bank terputus"
	} pulihkan (err) {
		err_msg = err
	}
	[saldo, err_msg]
	`
	resError := testEvalSource(t, codeError)
	arrError, ok := resError.(*evaluator.Array)
	if !ok || len(arrError.Elements) != 2 {
		t.Fatalf("expected array [saldo, err_msg], got: %v", resError)
	}
	if arrError.Elements[0].(*evaluator.Integer).Value != 500 {
		t.Errorf("expected saldo 500 after error rollback, got: %v", arrError.Elements[0])
	}
	if arrError.Elements[1].(*evaluator.String).Value != "koneksi bank terputus" {
		t.Errorf("expected err_msg 'koneksi bank terputus', got: %v", arrError.Elements[1])
	}

	// 4. Nested Transactions: inner rolls back, outer succeeds
	codeNestedInner := `
	buat x = 10
	buat y = 20
	transaksi {
		x = 100
		transaksi {
			y = 200
			batalkan "batal dalam"
		}
	}
	[x, y]
	`
	resNestedInner := testEvalSource(t, codeNestedInner)
	arrNestedInner, ok := resNestedInner.(*evaluator.Array)
	if !ok || len(arrNestedInner.Elements) != 2 {
		t.Fatalf("expected array [x, y], got: %v", resNestedInner)
	}
	if arrNestedInner.Elements[0].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected x 100, got: %v", arrNestedInner.Elements[0])
	}
	if arrNestedInner.Elements[1].(*evaluator.Integer).Value != 20 {
		t.Errorf("expected y 20 (rolled back), got: %v", arrNestedInner.Elements[1])
	}

	// 5. Nested Transactions: outer rolls back all
	codeNestedOuter := `
	buat x = 10
	buat y = 20
	transaksi {
		x = 100
		transaksi {
			y = 200
		}
		batalkan "batal luar"
	}
	[x, y]
	`
	resNestedOuter := testEvalSource(t, codeNestedOuter)
	arrNestedOuter, ok := resNestedOuter.(*evaluator.Array)
	if !ok || len(arrNestedOuter.Elements) != 2 {
		t.Fatalf("expected array [x, y], got: %v", resNestedOuter)
	}
	if arrNestedOuter.Elements[0].(*evaluator.Integer).Value != 10 {
		t.Errorf("expected x 10, got: %v", arrNestedOuter.Elements[0])
	}
	if arrNestedOuter.Elements[1].(*evaluator.Integer).Value != 20 {
		t.Errorf("expected y 20, got: %v", arrNestedOuter.Elements[1])
	}

	// 6. Rollback with shadowing in inner scope does not destroy outer variable
	codeShadow := `
	buat x = 100
	transaksi {
		untuk i dalam [1] {
			buat x = 999
			batalkan "batal"
		}
	}
	x
	`
	resShadow := testEvalSource(t, codeShadow)
	iShadow, ok := resShadow.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected outer variable x to survive rollback, got: %v", resShadow)
	}
	if iShadow.Value != 100 {
		t.Errorf("expected outer x = 100, got: %d", iShadow.Value)
	}
}

func TestFeature_ReaktifDanPengamat(t *testing.T) {
	// 1. Variabel Reaktif (auto-update derived dependencies)
	codeReactive := `
	buat harga = 50.ribu
	buat jumlah = 2
	reaktif total = harga * jumlah
	buat t1 = total
	harga = 75.ribu
	buat t2 = total
	jumlah = 3
	buat t3 = total
	[t1, t2, t3]
	`
	resReactive := testEvalSource(t, codeReactive)
	arrReactive, ok := resReactive.(*evaluator.Array)
	if !ok || len(arrReactive.Elements) != 3 {
		t.Fatalf("expected array [t1, t2, t3], got: %v", resReactive)
	}
	if arrReactive.Elements[0].(*evaluator.Integer).Value != 100000 {
		t.Errorf("expected t1 100000, got: %v", arrReactive.Elements[0])
	}
	if arrReactive.Elements[1].(*evaluator.Integer).Value != 150000 {
		t.Errorf("expected t2 150000, got: %v", arrReactive.Elements[1])
	}
	if arrReactive.Elements[2].(*evaluator.Integer).Value != 225000 {
		t.Errorf("expected t3 225000, got: %v", arrReactive.Elements[2])
	}

	// 2. Pengamat (pantau)
	codeWatch := `
	buat saldo = 100
	buat logs = []
	pantau saldo (baru, lama) {
		logs = logs + ["ubah: " + lama + " -> " + baru]
	}
	saldo = 150
	saldo = 200
	logs
	`
	resWatch := testEvalSource(t, codeWatch)
	arrWatch, ok := resWatch.(*evaluator.Array)
	if !ok || len(arrWatch.Elements) != 2 {
		t.Fatalf("expected 2 log entries, got: %v", resWatch)
	}
	if arrWatch.Elements[0].(*evaluator.String).Value != "ubah: 100 -> 150" {
		t.Errorf("expected log 0 'ubah: 100 -> 150', got: %v", arrWatch.Elements[0])
	}
	if arrWatch.Elements[1].(*evaluator.String).Value != "ubah: 150 -> 200" {
		t.Errorf("expected log 1 'ubah: 150 -> 200', got: %v", arrWatch.Elements[1])
	}

	// 3. Sinergi: Pantau variabel reaktif saat transaksi di-rollback
	codeSynergy := `
	buat saldo = 100
	reaktif bonus = saldo * 2
	buat catatan = []
	pantau bonus (baru, lama) {
		catatan = catatan + [baru]
	}
	transaksi {
		saldo = 500
		batalkan "tes batal"
	}
	[saldo, bonus, catatan]
	`
	resSynergy := testEvalSource(t, codeSynergy)
	arrSynergy, ok := resSynergy.(*evaluator.Array)
	if !ok || len(arrSynergy.Elements) != 3 {
		t.Fatalf("expected synergy array [saldo, bonus, catatan], got: %v", resSynergy)
	}
	if arrSynergy.Elements[0].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected saldo reverted to 100, got: %v", arrSynergy.Elements[0])
	}
	if arrSynergy.Elements[1].(*evaluator.Integer).Value != 200 {
		t.Errorf("expected bonus reverted to 200, got: %v", arrSynergy.Elements[1])
	}
	catatanArr := arrSynergy.Elements[2].(*evaluator.Array)
	// catatan was modified inside transaction, so it was rolled back to [] before receiving the reverted 200 event
	if len(catatanArr.Elements) != 1 {
		t.Fatalf("expected 1 note in catatan (rolled back then reverted), got: %v", catatanArr)
	}
	if catatanArr.Elements[0].(*evaluator.Integer).Value != 200 {
		t.Errorf("expected note 200 (reverted), got: %v", catatanArr.Elements[0])
	}
}

func TestFeature_BugFixesAudit(t *testing.T) {
	// 1. Bug 1: Computed properties on numeric members (num[prop])
	codeComp := `
	buat prop = "rupiah"
	buat harga = 50.ribu
	harga[prop]
	`
	resComp := testEvalSource(t, codeComp)
	if s, ok := resComp.(*evaluator.String); !ok || s.Value != "Rp 50.000" {
		t.Errorf("expected 'Rp 50.000' for computed numeric property, got: %v", resComp)
	}

	// 2. Bug 2: math.MinInt64 does not overflow or crash
	minStr := evaluator.TerbilangAngka(math.MinInt64)
	if !strings.HasPrefix(minStr, "minus ") {
		t.Errorf("expected MinInt64 to start with 'minus ', got: %s", minStr)
	}

	// 3. Bug 3: Negative float formatting retains minus sign
	codeNeg := `
	buat f = -0.5
	[f.terbilang, f.rupiah, format_rupiah(f)]
	`
	resNeg := testEvalSource(t, codeNeg)
	arrNeg, ok := resNeg.(*evaluator.Array)
	if !ok || len(arrNeg.Elements) != 3 {
		t.Fatalf("expected 3 elements for negative float, got: %v", resNeg)
	}
	if arrNeg.Elements[0].(*evaluator.String).Value != "minus nol koma lima nol rupiah" {
		t.Errorf("expected 'minus nol koma lima nol rupiah', got: %v", arrNeg.Elements[0])
	}
	if arrNeg.Elements[1].(*evaluator.String).Value != "-Rp 0,50" {
		t.Errorf("expected '-Rp 0,50', got: %v", arrNeg.Elements[1])
	}
	if arrNeg.Elements[2].(*evaluator.String).Value != "-Rp 0,50" {
		t.Errorf("expected '-Rp 0,50', got: %v", arrNeg.Elements[2])
	}

	// 4. Bug 4: RollbackSignal returned from function aborts variable assignment
	codeFnRollback := `
	buat saldo = 100
	fungsi gagal() {
		saldo = 50
		batalkan "batal dari fungsi"
		saldo = 0
	}
	transaksi {
		buat x = gagal()
	} pulihkan (err) {
		saldo
	}
	saldo
	`
	resFnRollback := testEvalSource(t, codeFnRollback)
	if i, ok := resFnRollback.(*evaluator.Integer); !ok || i.Value != 100 {
		t.Errorf("expected saldo reverted to 100 after function batalkan, got: %v", resFnRollback)
	}

	// 5. Bug 5: Phantom variables created inside rolled-back transaction are deleted
	codePhantom := `
	transaksi {
		buat akun_rahasia = 999
		batalkan "batal transaksi"
	}
	coba {
		akun_rahasia
	} tangkap (e) {
		"terhapus"
	}
	`
	resPhantom := testEvalSource(t, codePhantom)
	if s, ok := resPhantom.(*evaluator.String); !ok || s.Value != "terhapus" {
		t.Errorf("expected phantom variable to be deleted after rollback, got: %v", resPhantom)
	}

	// 6. Bug 6: In-place array element, map property, and nested mutations rollback
	codeCompositeRollback := `
	buat data = [10, 20, 30]
	buat user = { "nama": "Budi", "poin": 100 }
	buat matriks = [[1, 2], [3, 4]]
	transaksi {
		data[0] = 999
		user.poin = 500
		matriks[0][1] = 888
		batalkan "batal mutasi komposit"
	}
	[data[0], user.poin, matriks[0][1]]
	`
	resCompRollback := testEvalSource(t, codeCompositeRollback)
	arrCompRollback, ok := resCompRollback.(*evaluator.Array)
	if !ok || len(arrCompRollback.Elements) != 3 {
		t.Fatalf("expected 3 elements for composite rollback, got: %v", resCompRollback)
	}
	if arrCompRollback.Elements[0].(*evaluator.Integer).Value != 10 {
		t.Errorf("expected data[0] reverted to 10, got: %v", arrCompRollback.Elements[0])
	}
	if arrCompRollback.Elements[1].(*evaluator.Integer).Value != 100 {
		t.Errorf("expected user.poin reverted to 100, got: %v", arrCompRollback.Elements[1])
	}
	if arrCompRollback.Elements[2].(*evaluator.Integer).Value != 2 {
		t.Errorf("expected matriks[0][1] reverted to 2, got: %v", arrCompRollback.Elements[2])
	}

	// 7. Bug 7: Infinite recursion guard on self-mutating watcher
	codeWatcherLoop := `
	buat count = 1
	pantau count (baru, lama) {
		count = count + 1
	}
	coba {
		count = 2
	} tangkap (err) {
		"berhasil dicegah"
	}
	`
	resWatcherLoop := testEvalSource(t, codeWatcherLoop)
	if s, ok := resWatcherLoop.(*evaluator.String); !ok || s.Value != "berhasil dicegah" {
		t.Errorf("expected watcher infinite recursion to be guarded, got: %v", resWatcherLoop)
	}

	// 8. Bug 8: Reactive variable history sliding window after 50 updates
	codeHistSliding := `
	buat val = 0
	reaktif derived = val * 2
	untuk i dalam rentang(1, 60) {
		val = i
	}
	[panjang(derived.histori), derived.sebelumnya]
	`
	resHistSliding := testEvalSource(t, codeHistSliding)
	arrHistSliding, ok := resHistSliding.(*evaluator.Array)
	if !ok || len(arrHistSliding.Elements) != 2 {
		t.Fatalf("expected 2 elements for sliding history, got: %v", resHistSliding)
	}
	if arrHistSliding.Elements[0].(*evaluator.Integer).Value != 50 {
		t.Errorf("expected history capped at 50, got: %v", arrHistSliding.Elements[0])
	}
	if arrHistSliding.Elements[1].(*evaluator.Integer).Value != 118 {
		t.Errorf("expected previous derived value 118, got: %v", arrHistSliding.Elements[1])
	}
}

func TestEvaluator_LainBranching(t *testing.T) {
	code := `
	fungsi cek(x) {
		jika x == 1 {
			"satu"
		} lain jika x == 2 {
			"dua"
		} lain {
			"banyak"
		}
	}
	[cek(1), cek(2), cek(3)]
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got: %v", res)
	}
	if arr.Elements[0].(*evaluator.String).Value != "satu" ||
		arr.Elements[1].(*evaluator.String).Value != "dua" ||
		arr.Elements[2].(*evaluator.String).Value != "banyak" {
		t.Errorf("unexpected evaluation: %v", res)
	}
}

func TestEvaluator_BreakAndContinue(t *testing.T) {
	code := `
	buat sum = 0
	buat i = 0
	selama true {
		i += 1
		jika i == 2 {
			lanjut
		}
		jika i > 5 {
			berhenti
		}
		sum += i
	}
	sum
	`
	res := testEvalSource(t, code)
	if intVal, ok := res.(*evaluator.Integer); !ok || intVal.Value != (1+3+4+5) {
		t.Fatalf("expected sum 13, got %s", res.Inspect())
	}

	codeFor := `
	buat res = []
	untuk x dalam rentang(1, 10) {
		jika x % 2 == 0 {
			lanjut
		}
		jika x > 7 {
			berhenti
		}
		res.tambah(x)
	}
	res
	`
	resFor := testEvalSource(t, codeFor)
	if arr, ok := resFor.(*evaluator.Array); !ok || arr.Inspect() != "[1, 3, 5, 7]" {
		t.Fatalf("expected [1, 3, 5, 7], got %s", resFor.Inspect())
	}
}

func TestEvaluator_BreakAndContinue_InvalidScopes(t *testing.T) {
	// 1. Top-level break should error
	resBreak := testEvalSource(t, `berhenti`)
	if errObj, ok := resBreak.(*evaluator.Error); !ok || !strings.Contains(errObj.Message, "hanya dapat digunakan di dalam perulangan") {
		t.Fatalf("expected error for top-level break, got %v", resBreak)
	}

	// 2. Top-level continue should error
	resContinue := testEvalSource(t, `lanjut`)
	if errObj, ok := resContinue.(*evaluator.Error); !ok || !strings.Contains(errObj.Message, "hanya dapat digunakan di dalam perulangan") {
		t.Fatalf("expected error for top-level continue, got %v", resContinue)
	}

	// 3. Break inside function should error and not escape to outer loop
	codeFnBreak := `
	fungsi coba_berhenti() {
		berhenti
	}
	buat loopCount = 0
	selama loopCount < 3 {
		loopCount += 1
		coba_berhenti()
	}
	loopCount
	`
	resFnBreak := testEvalSource(t, codeFnBreak)
	if _, ok := resFnBreak.(*evaluator.Error); !ok {
		t.Fatalf("expected error when function executes break, got %v", resFnBreak)
	}
}

func TestEvaluator_NegativeIndexSlicing(t *testing.T) {
	code := `
	buat arr = [10, 20, 30, 40, 50]
	buat s1 = arr[-2..]
	buat s2 = arr[..-1]
	buat s3 = arr[-3..-1]
	buat s4 = arr.potong(-2)

	buat kata = "Nusantara"
	buat sub1 = kata[-3..]
	buat sub2 = kata[..-4]

	[s1, s2, s3, s4, sub1, sub2]
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].Inspect() != "[40, 50]" {
		t.Errorf("expected arr[-2..] = [40, 50], got %s", arr.Elements[0].Inspect())
	}
	if arr.Elements[1].Inspect() != "[10, 20, 30, 40]" {
		t.Errorf("expected arr[..-1] = [10, 20, 30, 40], got %s", arr.Elements[1].Inspect())
	}
	if arr.Elements[2].Inspect() != "[30, 40]" {
		t.Errorf("expected arr[-3..-1] = [30, 40], got %s", arr.Elements[2].Inspect())
	}
	if arr.Elements[3].Inspect() != "[40, 50]" {
		t.Errorf("expected arr.potong(-2) = [40, 50], got %s", arr.Elements[3].Inspect())
	}
	if arr.Elements[4].Inspect() != "ara" {
		t.Errorf("expected kata[-3..] = 'ara', got %s", arr.Elements[4].Inspect())
	}
	if arr.Elements[5].Inspect() != "Nusan" {
		t.Errorf("expected kata[..-4] = 'Nusan', got %s", arr.Elements[5].Inspect())
	}
}

func TestEvaluator_MapMethods_WithoutInitialOrder(t *testing.T) {
	// Test that map operations (kunci, nilai, entri, panjang, ukuran) work even if Map was created without Order
	code := `
	waktu_info = waktu.rincian()
	kunci_list = waktu_info.kunci()
	p = waktu_info.panjang
	u = waktu_info.ukuran()
	[kunci_list.panjang, p, u]
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	kLen := arr.Elements[0].(*evaluator.Integer).Value
	pLen := arr.Elements[1].(*evaluator.Integer).Value
	uLen := arr.Elements[2].(*evaluator.Integer).Value
	if kLen != 12 || pLen != 12 || uLen != 12 {
		t.Fatalf("expected 12 keys/length/size, got %d, %d, %d", kLen, pLen, uLen)
	}
}

func TestEvaluator_FloatNaturalSort(t *testing.T) {
	code := `
	buat arr = [0.9, 0.1, 0.5, 0.2]
	buat diurut = arr.urutkan()
	[diurut[0], diurut[1], diurut[2], diurut[3]]
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	v0 := arr.Elements[0].(*evaluator.Float).Value
	v1 := arr.Elements[1].(*evaluator.Float).Value
	v2 := arr.Elements[2].(*evaluator.Float).Value
	v3 := arr.Elements[3].(*evaluator.Float).Value
	if v0 != 0.1 || v1 != 0.2 || v2 != 0.5 || v3 != 0.9 {
		t.Errorf("expected sorted floats [0.1, 0.2, 0.5, 0.9], got [%v, %v, %v, %v]", v0, v1, v2, v3)
	}
}

func TestEvaluator_Rentang_StrictIntegers(t *testing.T) {
	code := `rentang(1.5, 3.5)`
	res := testEvalSource(t, code)
	if _, ok := res.(*evaluator.Error); !ok {
		t.Fatalf("expected error passing floats to rentang(), got %s", res.Inspect())
	}
}

func TestEvaluator_SafeNavigationParity(t *testing.T) {
	code := `
	buat arr = [10, 20, 30]
	buat mp = {"nama": "Celang"}
	buat teks = "Halo"
	buat objKosong = null

	buat pArr = arr?.panjang
	arr?.tambah(40)
	buat pArr2 = arr?.panjang
	buat kMp = mp?.kunci()
	buat ch0 = teks?.[0]

	buat pNull = objKosong?.panjang
	buat mNull = objKosong?.kunci()

	[pArr, pArr2, kMp.panjang, ch0, pNull == null, mNull == null]
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok {
		t.Fatalf("expected array, got %s", res.Inspect())
	}
	if arr.Elements[0].(*evaluator.Integer).Value != 3 {
		t.Errorf("expected arr?.panjang == 3, got %v", arr.Elements[0])
	}
	if arr.Elements[1].(*evaluator.Integer).Value != 4 {
		t.Errorf("expected arr?.panjang after push == 4, got %v", arr.Elements[1])
	}
	if arr.Elements[2].(*evaluator.Integer).Value != 1 {
		t.Errorf("expected mp?.kunci().panjang == 1, got %v", arr.Elements[2])
	}
	if arr.Elements[3].(*evaluator.String).Value != "H" {
		t.Errorf("expected teks?.[0] == 'H', got %v", arr.Elements[3])
	}
	if !arr.Elements[4].(*evaluator.Boolean).Value {
		t.Errorf("expected objKosong?.panjang == null, got %v", arr.Elements[4])
	}
	if !arr.Elements[5].(*evaluator.Boolean).Value {
		t.Errorf("expected objKosong?.kunci() == null, got %v", arr.Elements[5])
	}
}

func TestEvaluator_SliceBounds_Validation(t *testing.T) {
	// 1. Array slicing with invalid non-integer bound
	resArr := testEvalSource(t, `buat a = [1, 2, 3]; a["satu"..2]`)
	if _, ok := resArr.(*evaluator.Error); !ok {
		t.Fatalf("expected error for non-integer array slice bound, got %s", resArr.Inspect())
	}

	// 2. String slicing with invalid non-integer bound
	resStr := testEvalSource(t, `buat s = "halo"; s[0.."dua"]`)
	if _, ok := resStr.(*evaluator.Error); !ok {
		t.Fatalf("expected error for non-integer string slice bound, got %s", resStr.Inspect())
	}

	// 3. potong/slice method with non-integer arg
	resPotong := testEvalSource(t, `buat a = [1, 2, 3]; a.potong("bukan_angka")`)
	if _, ok := resPotong.(*evaluator.Error); !ok {
		t.Fatalf("expected error for non-integer potong argument, got %s", resPotong.Inspect())
	}
}

func TestEvaluator_StructuralEquality(t *testing.T) {
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
	buat e11 = ({"a": [1, 2]} == {"a": [1, 2]})

	[e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11]
	`
	res := testEvalSource(t, code)
	arr := res.(*evaluator.Array)
	expectations := []bool{false, false, false, true, true, true, false, false, true, true, true}
	for i, exp := range expectations {
		got := arr.Elements[i].(*evaluator.Boolean).Value
		if got != exp {
			t.Errorf("e%d expected %v, got %v", i+1, exp, got)
		}
	}
}

func TestEvaluator_RangeSizeGuard(t *testing.T) {
	res1 := testEvalSource(t, `rentang(1, 10000000)`)
	if _, ok := res1.(*evaluator.Error); !ok {
		t.Fatalf("expected error for excessive rentang size, got %s", res1.Inspect())
	}

	res2 := testEvalSource(t, `1..10000000`)
	if _, ok := res2.(*evaluator.Error); !ok {
		t.Fatalf("expected error for excessive range expression size, got %s", res2.Inspect())
	}
}

func TestEvaluator_RangeMaxIntOverflow(t *testing.T) {
	// 1. rentang reaching math.MaxInt64 without hanging
	res1 := testEvalSource(t, `rentang(9223372036854775805, 9223372036854775807)`)
	arr1, ok := res1.(*evaluator.Array)
	if !ok || len(arr1.Elements) != 3 {
		t.Fatalf("expected 3 elements for range to MaxInt64, got %v", res1)
	}

	// 2. range operator .. reaching math.MaxInt64 without hanging
	res2 := testEvalSource(t, `9223372036854775805..9223372036854775807`)
	arr2, ok := res2.(*evaluator.Array)
	if !ok || len(arr2.Elements) != 3 {
		t.Fatalf("expected 3 elements for .. to MaxInt64, got %v", res2)
	}
}

func TestEvaluator_MathOverflow(t *testing.T) {
	code := `matematika.pangkat(10, 25)`
	res := testEvalSource(t, code)
	if flt, ok := res.(*evaluator.Float); !ok || flt.Value < 0 {
		t.Fatalf("expected positive Float for 10^25 without int64 overflow, got %s", res.Inspect())
	}
}

func TestAudit_ChannelSendClosedRecover(t *testing.T) {
	ch := evaluator.NewChannel(0, "")
	go func() {
		time.Sleep(20 * time.Millisecond)
		ch.Close()
	}()
	res := ch.Send(&evaluator.String{Value: "test"}, nil)
	if _, ok := res.(*evaluator.Error); !ok {
		t.Fatalf("expected error on send to closed channel, got %v", res)
	}
}

func TestAudit_CircularDataStructures(t *testing.T) {
	code := `
	buat a = [1]
	a.tambah(a)
	buat b = [1]
	b.tambah(b)
	buat eq = (a == b)
	buat c = salin(a)
	eq
	`
	res := testEvalSource(t, code)
	if b, ok := res.(*evaluator.Boolean); !ok || !b.Value {
		t.Fatalf("expected circular equality to be true without stack overflow, got %s", res.Inspect())
	}
}

func TestAudit_ConstantRedeclaration(t *testing.T) {
	code := `
	tetap KUNCI = 100
	tetap KUNCI = 200
	`
	res := testEvalSource(t, code)
	if _, ok := res.(*evaluator.Error); !ok {
		t.Fatalf("expected error for re-declaring constant 'tetap', got %s", res.Inspect())
	}
}

func TestAudit_ArrayDotNotationProperty(t *testing.T) {
	code := `
	buat foo = 0
	buat arr = [10, 20]
	arr.foo
	`
	res := testEvalSource(t, code)
	if _, ok := res.(*evaluator.Error); !ok {
		t.Fatalf("expected error for accessing invalid dot property 'foo' on Array, got %s", res.Inspect())
	}
}

func TestAudit_FinallyBlockSignals(t *testing.T) {
	code := `
	buat iterasi = 0
	selama (iterasi < 5) {
		iterasi = iterasi + 1
		coba {
			// normal
		} akhirnya {
			berhenti
		}
	}
	iterasi
	`
	res := testEvalSource(t, code)
	if iVal, ok := res.(*evaluator.Integer); !ok || iVal.Value != 1 {
		t.Fatalf("expected loop to terminate on iteration 1 via finally break, got %s", res.Inspect())
	}
}

func TestAudit_JemputAllInvalidElement(t *testing.T) {
	code := `
	buat t1 = titip 1
	jemput_semua([t1, "bukan_tiket"])
	`
	res := testEvalSource(t, code)
	if _, ok := res.(*evaluator.Error); !ok {
		t.Fatalf("expected error for jemput_semua with non-ticket, got %s", res.Inspect())
	}
}

func TestAudit_JSONMaxAutoExpand(t *testing.T) {
	code := `
	buat data = []
	json.setel(data, "[10000000]", 42)
	`
	res := testEvalSource(t, code)
	arr, ok := res.(*evaluator.Array)
	if !ok || len(arr.Elements) != 0 {
		t.Fatalf("expected data not to expand to 10M elements, got length %d", len(arr.Elements))
	}
}

func TestAudit_TitipAssertionErrorDoesNotPanic(t *testing.T) {
	code := `
	fungsi gagal() {
		harus_sama(1, 2)
	}
	buat t = titip gagal()
	jemput t
	`
	res := testEvalSource(t, code)
	if res == nil {
		t.Fatalf("expected result, got nil")
	}
	if aErr, ok := res.(*evaluator.AssertionError); !ok {
		t.Fatalf("expected AssertionError, got %T (%s)", res, res.Inspect())
	} else if !strings.Contains(aErr.Message, "Gagal:") && !strings.Contains(aErr.Message, "ekspektasi") {
		t.Errorf("unexpected assertion error message: %s", aErr.Message)
	}

	// Test with jemput_semua
	codeSemua := `
	fungsi gagalSemua() {
		harus_sama(10, 20)
	}
	buat t1 = titip gagalSemua()
	jemput_semua([t1])
	`
	resSemua := testEvalSource(t, codeSemua)
	if _, ok := resSemua.(*evaluator.AssertionError); !ok {
		t.Fatalf("expected jemput_semua to return AssertionError without panic, got %T (%s)", resSemua, resSemua.Inspect())
	}
}

func TestAudit_MemberAccessSideEffectEvaluatedOnce(t *testing.T) {
	codeArray := `
	buat arr = [100, 200, 300]
	buat hitungan = 0
	buat dapatkanIndeks = fungsi() {
		hitungan = hitungan + 1
		kembalikan 1
	}
	buat el = arr[dapatkanIndeks()]
	[el, hitungan]
	`
	resArr := testEvalSource(t, codeArray)
	arr, ok := resArr.(*evaluator.Array)
	if !ok || len(arr.Elements) != 2 {
		t.Fatalf("expected 2-element array, got %v", resArr)
	}
	if elVal, ok := arr.Elements[0].(*evaluator.Integer); !ok || elVal.Value != 200 {
		t.Errorf("expected el = 200, got %v", arr.Elements[0])
	}
	if countVal, ok := arr.Elements[1].(*evaluator.Integer); !ok || countVal.Value != 1 {
		t.Errorf("expected hitungan = 1 (evaluated once), got %v", arr.Elements[1])
	}

	codeBytes := `
	buat b = bita([10, 20, 30])
	buat hitunganB = 0
	buat dapatkanIndeksB = fungsi() {
		hitunganB = hitunganB + 1
		kembalikan 2
	}
	buat bVal = b[dapatkanIndeksB()]
	[bVal, hitunganB]
	`
	resBytes := testEvalSource(t, codeBytes)
	arrB, ok := resBytes.(*evaluator.Array)
	if !ok || len(arrB.Elements) != 2 {
		t.Fatalf("expected 2-element array, got %v", resBytes)
	}
	if bVal, ok := arrB.Elements[0].(*evaluator.Integer); !ok || bVal.Value != 30 {
		t.Errorf("expected bVal = 30, got %v", arrB.Elements[0])
	}
	if countValB, ok := arrB.Elements[1].(*evaluator.Integer); !ok || countValB.Value != 1 {
		t.Errorf("expected hitunganB = 1 (evaluated once), got %v", arrB.Elements[1])
	}
}

func TestAudit_MemberAccessErrorPropagation(t *testing.T) {
	code := `
	buat arr = [1, 2, 3]
	buat lemparGalat = fungsi() {
		lempar "Koneksi database putus!"
	}
	arr[lemparGalat()]
	`
	res := testEvalSource(t, code)
	if res == nil || res.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("expected Error object, got %v", res)
	}
	errObj := res.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "Koneksi database putus!") {
		t.Errorf("expected error message to be propagated, got: %s", errObj.Message)
	}
}
