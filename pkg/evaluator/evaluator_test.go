package evaluator_test

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"celang/pkg/evaluator"
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

