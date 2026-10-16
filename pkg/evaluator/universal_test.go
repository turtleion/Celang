package evaluator_test

import (
	"bytes"
	"strings"
	"testing"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
)

func runScript(input string) (evaluator.Object, *evaluator.Environment, error) {
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return nil, nil, err
	}

	env := evaluator.NewEnvironment()
	var out bytes.Buffer
	evaluator.RegisterBuiltinsWithIO(env, nil, &out)
	res := evaluator.Eval(prog, env)
	return res, env, nil
}

func TestUniversalMode_DirectiveDetection(t *testing.T) {
	directives := []string{
		`"use universal"`,
		`'use universal'`,
		`"use english"`,
		`'use english'`,
		`"use standard"`,
		`// initial comment
"use universal"`,
	}

	for _, d := range directives {
		if !parser.DetectUniversalDirective(d) {
			t.Errorf("Expected DetectUniversalDirective to be true for:\n%s", d)
		}
	}

	nonDirectives := []string{
		`let x = 10`,
		`buat x = 10`,
		`// "use universal" in comment`,
	}

	for _, nd := range nonDirectives {
		if parser.DetectUniversalDirective(nd) {
			t.Errorf("Expected DetectUniversalDirective to be false for:\n%s", nd)
		}
	}
}

func TestUniversalMode_RejectIndonesianKeywords(t *testing.T) {
	cases := []struct {
		name    string
		code    string
		wantErr string
	}{
		{
			name:    "buat declaration",
			code:    "\"use universal\"\nbuat x = 10",
			wantErr: "Indonesian keyword 'buat' is not allowed in Universal Mode",
		},
		{
			name:    "tetap declaration",
			code:    "\"use universal\"\ntetap x = 10",
			wantErr: "Indonesian keyword 'tetap' is not allowed in Universal Mode",
		},
		{
			name:    "jika condition",
			code:    "\"use universal\"\nif true { jika false {} }",
			wantErr: "Indonesian keyword 'jika' is not allowed in Universal Mode",
		},
		{
			name:    "selain block",
			code:    "\"use universal\"\nif true {} selain {}",
			wantErr: "Indonesian keyword 'selain' is not allowed in Universal Mode",
		},
		{
			name:    "fungsi declaration",
			code:    "\"use universal\"\nfungsi hello() { return 1 }",
			wantErr: "Indonesian keyword 'fungsi' is not allowed in Universal Mode",
		},
		{
			name:    "kembalikan return",
			code:    "\"use universal\"\ndef hello() { kembalikan 1 }",
			wantErr: "Indonesian keyword 'kembalikan' is not allowed in Universal Mode",
		},
		{
			name:    "transaksi block",
			code:    "\"use universal\"\ntransaksi {}",
			wantErr: "Indonesian keyword 'transaksi' is not allowed in Universal Mode",
		},
		{
			name:    "batalkan rollback",
			code:    "\"use universal\"\ntransaction { batalkan \"err\" }",
			wantErr: "Indonesian keyword 'batalkan' is not allowed in Universal Mode",
		},
		{
			name:    "pulihkan recover",
			code:    "\"use universal\"\ntransaction {} pulihkan (e) {}",
			wantErr: "Indonesian keyword 'pulihkan' is not allowed in Universal Mode",
		},
		{
			name:    "pantau watch",
			code:    "\"use universal\"\nlet x = 1\npantau x {}",
			wantErr: "Indonesian keyword 'pantau' is not allowed in Universal Mode",
		},
		{
			name:    "reaktif reactive",
			code:    "\"use universal\"\nlet x = 1\nreaktif y = x * 2",
			wantErr: "Indonesian keyword 'reaktif' is not allowed in Universal Mode",
		},
		{
			name:    "benar boolean",
			code:    "\"use universal\"\nlet x = benar",
			wantErr: "Indonesian literal 'benar' is not allowed in Universal Mode",
		},
		{
			name:    "salah boolean",
			code:    "\"use universal\"\nlet x = salah",
			wantErr: "Indonesian literal 'salah' is not allowed in Universal Mode",
		},
		{
			name:    "nihil null",
			code:    "\"use universal\"\nlet x = nihil",
			wantErr: "Indonesian literal 'nihil' is not allowed in Universal Mode",
		},
		{
			name:    "dan operator",
			code:    "\"use universal\"\nlet x = true dan false",
			wantErr: "Indonesian operator 'dan' is not allowed in Universal Mode",
		},
		{
			name:    "atau operator",
			code:    "\"use universal\"\nlet x = true atau false",
			wantErr: "Indonesian operator 'atau' is not allowed in Universal Mode",
		},
		{
			name:    "bukan prefix operator",
			code:    "\"use universal\"\nlet x = bukan true",
			wantErr: "Indonesian operator 'bukan' is not allowed in Universal Mode",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := runScript(tc.code)
			if err == nil {
				t.Fatalf("Expected parser error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Error %q does not contain expected substring %q", err.Error(), tc.wantErr)
			}
		})
	}
}

func TestUniversalMode_EnglishKeywordsExecution(t *testing.T) {
	code := `"use universal"
let a = 15
const b = 25
let sum = a + b

let branch = ""
if sum > 30 {
    branch = "greater"
} else {
    branch = "smaller"
}

def add(x, y) {
    return x + y
}
let res = add(5, 10)
`
	_, env, err := runScript(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	sumVal, _ := env.Get("sum")
	if sumVal.(*evaluator.Integer).Value != 40 {
		t.Errorf("Expected sum=40, got %v", sumVal)
	}

	branchVal, _ := env.Get("branch")
	if branchVal.(*evaluator.String).Value != "greater" {
		t.Errorf("Expected branch='greater', got %v", branchVal)
	}

	resVal, _ := env.Get("res")
	if resVal.(*evaluator.Integer).Value != 15 {
		t.Errorf("Expected res=15, got %v", resVal)
	}
}

func TestUniversalMode_CurrencyAndNumbers(t *testing.T) {
	code := `"use universal"
let kVal = 50.k
let thousandVal = 10.thousand
let millionVal = 5.million
let billionVal = 2.billion
let trillionVal = 1.trillion

let price1 = 1250.usd
let price2 = 25.50.usd

let words1 = 42.words
let words2 = 105.in_words
let dw1 = 25.50.dollars_in_words
let dw2 = 10.dollars_in_words

let tSec = 5.sec
let tMin = 2.min
let tHour = 1.hour
let tDay = 1.day
`
	_, env, err := runScript(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assertInt := func(name string, expected int64) {
		val, ok := env.Get(name)
		if !ok {
			t.Fatalf("Variable %s not found", name)
		}
		if iv, ok := val.(*evaluator.Integer); !ok || iv.Value != expected {
			t.Errorf("Expected %s = %d, got %v", name, expected, val)
		}
	}

	assertString := func(name string, expected string) {
		val, ok := env.Get(name)
		if !ok {
			t.Fatalf("Variable %s not found", name)
		}
		if sv, ok := val.(*evaluator.String); !ok || sv.Value != expected {
			t.Errorf("Expected %s = %q, got %v", name, expected, val)
		}
	}

	assertInt("kVal", 50000)
	assertInt("thousandVal", 10000)
	assertInt("millionVal", 5000000)
	assertInt("billionVal", 2000000000)
	assertInt("trillionVal", 1000000000000)

	assertString("price1", "$1,250.00")
	assertString("price2", "$25.50")

	assertString("words1", "forty-two")
	assertString("words2", "one hundred five")
	assertString("dw1", "twenty-five dollars and fifty cents")
	assertString("dw2", "ten dollars")

	assertInt("tSec", 5000)
	assertInt("tMin", 120000)
	assertInt("tHour", 3600000)
	assertInt("tDay", 86400000)
}

func TestUniversalMode_RejectIndonesianCulture(t *testing.T) {
	evalCases := []struct {
		name    string
		code    string
		wantErr string
	}{
		{
			name:    "ribu suffix",
			code:    "\"use universal\"\nlet x = 50.ribu",
			wantErr: "Indonesian suffix 'ribu' is not allowed in Universal Mode",
		},
		{
			name:    "juta suffix",
			code:    "\"use universal\"\nlet x = 10.juta",
			wantErr: "Indonesian suffix 'juta' is not allowed in Universal Mode",
		},
		{
			name:    "rupiah suffix",
			code:    "\"use universal\"\nlet x = 15000.rupiah",
			wantErr: "Indonesian suffix 'rupiah' is not allowed in Universal Mode",
		},
		{
			name:    "terbilang suffix",
			code:    "\"use universal\"\nlet x = 50000.terbilang",
			wantErr: "Indonesian suffix 'terbilang' is not allowed in Universal Mode",
		},
		{
			name:    "detik suffix",
			code:    "\"use universal\"\nlet x = 5.detik",
			wantErr: "Indonesian suffix 'detik' is not allowed in Universal Mode",
		},
		{
			name:    "cetak builtin",
			code:    "\"use universal\"\ncetak(\"hello\")",
			wantErr: "Indonesian built-in 'cetak' is not allowed in Universal Mode",
		},
		{
			name:    "terbilang builtin",
			code:    "\"use universal\"\nterbilang(100)",
			wantErr: "Indonesian function 'terbilang' is not allowed in Universal Mode",
		},
		{
			name:    "rupiah builtin",
			code:    "\"use universal\"\nrupiah(100)",
			wantErr: "Indonesian function 'rupiah' is not allowed in Universal Mode",
		},
		{
			name:    "berkas module import",
			code:    "\"use universal\"\nimport \"berkas\"",
			wantErr: "Indonesian module 'berkas' is not allowed in Universal Mode",
		},
		{
			name:    "sebelumnya time travel",
			code:    "\"use universal\"\nlet x = 10\nx = 20\nlet y = x.sebelumnya",
			wantErr: "Indonesian property '.sebelumnya' is not allowed in Universal Mode",
		},
	}

	for _, tc := range evalCases {
		t.Run(tc.name, func(t *testing.T) {
			res, _, err := runScript(tc.code)
			if err != nil {
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("Parser error %q does not contain %q", err.Error(), tc.wantErr)
				}
				return
			}
			if res == nil || res.Type() != evaluator.ERROR_OBJ {
				t.Fatalf("Expected evaluator error containing %q, got %v", tc.wantErr, res)
			}
			errObj := res.(*evaluator.Error)
			if !strings.Contains(errObj.Message, tc.wantErr) {
				t.Errorf("Error %q does not contain expected substring %q", errObj.Message, tc.wantErr)
			}
		})
	}
}

func TestUniversalMode_TimeTravel(t *testing.T) {
	code := `"use universal"
let x = 10
x = 20
x = 30

let p = x.previous
let initVal = x.initial
let h = x.history
`
	_, env, err := runScript(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	pVal, _ := env.Get("p")
	if pVal.(*evaluator.Integer).Value != 20 {
		t.Errorf("Expected p=20, got %v", pVal)
	}

	initVal, _ := env.Get("initVal")
	if initVal.(*evaluator.Integer).Value != 10 {
		t.Errorf("Expected initVal=10, got %v", initVal)
	}

	hVal, _ := env.Get("h")
	arr := hVal.(*evaluator.Array)
	if len(arr.Elements) != 3 {
		t.Errorf("Expected history length 3, got %d", len(arr.Elements))
	}
}

func TestUniversalMode_ModernFeatures(t *testing.T) {
	// Transactions & Reactivity & Watchers with English keywords
	code := `"use universal"
let bankBalance = 100
let userBalance = 50

transaction {
    bankBalance = bankBalance - 200
    userBalance = userBalance + 200
    if bankBalance < 0 {
        rollback "insufficient balance"
    }
} recover (reason) {
    let recoveredReason = reason
}

let count = 10
reactive doubled = count * 2

let watchTriggered = false
watch count (newVal, oldVal) {
    watchTriggered = true
}
count = 25
`
	_, env, err := runScript(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Transaction should have rolled back
	bb, _ := env.Get("bankBalance")
	if bb.(*evaluator.Integer).Value != 100 {
		t.Errorf("Expected bankBalance=100 after rollback, got %v", bb)
	}

	ub, _ := env.Get("userBalance")
	if ub.(*evaluator.Integer).Value != 50 {
		t.Errorf("Expected userBalance=50 after rollback, got %v", ub)
	}

	// Reactive doubled
	d, _ := env.Get("doubled")
	if d.(*evaluator.Integer).Value != 50 { // 25 * 2 = 50
		t.Errorf("Expected reactive doubled=50, got %v", d)
	}

	// Watcher triggered
	wt, _ := env.Get("watchTriggered")
	if wt.(*evaluator.Boolean).Value != true {
		t.Errorf("Expected watchTriggered=true, got %v", wt)
	}
}

func TestUniversalMode_BackwardCompatibilityWithoutDirective(t *testing.T) {
	// Without directive, flexible mode allows Indonesian and mixed code seamlessly
	code := `buat saldo = 50.ribu
jika saldo > 10.ribu {
    buat status = "aman"
} selain {
    buat status = "kurang"
}
buat teksSaldo = 50000.terbilang
buat rp = 50000.rupiah
`
	_, env, err := runScript(code)
	if err != nil {
		t.Fatalf("Backward compatibility broken: %v", err)
	}

	saldo, _ := env.Get("saldo")
	if saldo.(*evaluator.Integer).Value != 50000 {
		t.Errorf("Expected saldo=50000, got %v", saldo)
	}

	status, _ := env.Get("status")
	if status.(*evaluator.String).Value != "aman" {
		t.Errorf("Expected status='aman', got %v", status)
	}

	rp, _ := env.Get("rp")
	if rp.(*evaluator.String).Value != "Rp 50.000" {
		t.Errorf("Expected rp='Rp 50.000', got %v", rp)
	}
}
