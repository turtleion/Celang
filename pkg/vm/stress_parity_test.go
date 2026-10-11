package vm_test

import (
	"bytes"
	"sync"
	"testing"

	"celang/pkg/compiler"
	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/vm"
)

// runBoth executes Celang code in both AST Evaluator and Bytecode VM,
// returning the stdout outputs and final evaluated objects for comparison.
func runBoth(t *testing.T, input string) (astOut, vmOut string, astRes, vmRes evaluator.Object) {
	t.Helper()

	// 1. Parse program
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// 2. Run AST Evaluator
	var astBuf bytes.Buffer
	astEnv := evaluator.NewEnvironment()
	evaluator.RegisterBuiltinsWithIO(astEnv, nil, &astBuf)
	astRes = evaluator.Eval(prog, astEnv)
	astOut = astBuf.String()

	// 3. Run Bytecode VM
	comp := compiler.NewCompiler()
	chunk, err := comp.Compile(prog)
	if err != nil {
		t.Fatalf("Compiler error: %v", err)
	}

	var vmBuf bytes.Buffer
	vmEnv := evaluator.NewEnvironment()
	evaluator.RegisterBuiltinsWithIO(vmEnv, nil, &vmBuf)
	virtualMachine := vm.NewWithEnv(chunk, vmEnv)
	vmRes, err = virtualMachine.Run()
	if err != nil {
		t.Fatalf("VM runtime error: %v", err)
	}
	vmOut = vmBuf.String()

	return astOut, vmOut, astRes, vmRes
}

func assertParity(t *testing.T, testName, input string) {
	t.Helper()
	astOut, vmOut, astRes, vmRes := runBoth(t, input)

	if astOut != vmOut {
		t.Errorf("[%s] STDOUT Mismatch:\nAST Output:\n%s\nVM Output:\n%s", testName, astOut, vmOut)
	}

	astStr := "kosong"
	if astRes != nil {
		astStr = astRes.Inspect()
	}
	vmStr := "kosong"
	if vmRes != nil {
		vmStr = vmRes.Inspect()
	}

	if astStr != vmStr {
		t.Errorf("[%s] RESULT Mismatch:\nAST: %s\nVM:  %s", testName, astStr, vmStr)
	}
}

// ---------------------------------------------------------------------
// 1. STRESS TEST: Deep Recursion (Fibonacci & Factorial)
// ---------------------------------------------------------------------

func TestStress_DeepRecursion(t *testing.T) {
	input := `
fungsi fib(n) {
    jika n <= 1 {
        kembalikan n
    }
    kembalikan fib(n - 1) + fib(n - 2)
}

buat hasil = fib(18)
cetak("Fib(18):", hasil)
hasil
`
	assertParity(t, "Stress_DeepRecursion", input)
}

// ---------------------------------------------------------------------
// 2. STRESS TEST: Tight High-Iteration Loop (50,000 iterations)
// ---------------------------------------------------------------------

func TestStress_HighIterationLoop(t *testing.T) {
	input := `
buat total = 0
buat i = 1
selama i <= 50000 {
    jika i % 2 == 0 {
        total += i
    } selain {
        total -= 1
    }
    i += 1
}
cetak("Total akhir 50.000 iterasi:", total)
total
`
	assertParity(t, "Stress_HighIterationLoop", input)
}

// ---------------------------------------------------------------------
// 3. STRESS TEST: Large Array & Nested Structure Allocations
// ---------------------------------------------------------------------

func TestStress_LargeArrayAndMapOperations(t *testing.T) {
	input := `
buat data = []
buat n = 1000
buat j = 1

selama j <= n {
    data.tambah({
        "id": j,
        "kuadrat": j * j,
        "genap": j % 2 == 0
    })
    j += 1
}

cetak("Jumlah data dibuat:", panjang(data))
buat elemenTerakhir = data[panjang(data) - 1]
cetak("Elemen ke-1000 kuadrat:", elemenTerakhir.kuadrat)
elemenTerakhir.kuadrat
`
	assertParity(t, "Stress_LargeArrayAndMapOperations", input)
}

// ---------------------------------------------------------------------
// 4. STRESS TEST: Struct Instantiation & Method Invocation Under Load
// ---------------------------------------------------------------------

func TestStress_StructAndMethods(t *testing.T) {
	input := `
bentuk AkunDompet {
    nama: Teks,
    saldo itu Angka = 0,

    fungsi setor(nominal: Angka) {
        diri.saldo += nominal
        kembalikan diri.saldo
    }

    fungsi tarik(nominal: Angka) {
        diri.saldo -= nominal
        kembalikan diri.saldo
    }
}

buat dompet = AkunDompet("Yosatsu", 1000)
buat putaran = 1

selama putaran <= 500 {
    dompet.setor(20)
    dompet.tarik(10)
    putaran += 1
}

cetak("Saldo akhir dompet:", dompet.saldo)
dompet.saldo
`
	assertParity(t, "Stress_StructAndMethods", input)
}

// ---------------------------------------------------------------------
// 5. STRESS TEST: Closures & Higher-Order Functions Factory
// ---------------------------------------------------------------------

func TestStress_ClosuresFactory(t *testing.T) {
	input := `
fungsi buatPengali(faktor) {
    kembalikan fungsi(x) {
        kembalikan x * faktor
    }
}

buat kaliDua = buatPengali(2)
buat kaliTiga = buatPengali(3)
buat kaliSepuluh = buatPengali(10)

buat kumulatif = 0
buat idx = 1

selama idx <= 200 {
    kumulatif += kaliDua(idx) + kaliTiga(idx) + kaliSepuluh(idx)
    idx += 1
}

cetak("Akumulasi closure:", kumulatif)
kumulatif
`
	assertParity(t, "Stress_ClosuresFactory", input)
}

// ---------------------------------------------------------------------
// 6. STRESS TEST: Destructuring Array & Map Reassignment
// ---------------------------------------------------------------------

func TestStress_DestructuringStress(t *testing.T) {
	input := `
buat point = [100, 200, 300]
buat [x, y, z] = point
buat swapTotal = 0
buat k = 1

selama k <= 200 {
    [x, y] = [y + 1, x + 2]
    swapTotal += (x + y)
    k += 1
}

cetak("Hasil swap stress:", x, y, swapTotal)
swapTotal
`
	assertParity(t, "Stress_DestructuringStress", input)
}

// ---------------------------------------------------------------------
// 7. STRESS TEST: Pipeline Operators Chaining
// ---------------------------------------------------------------------

func TestStress_PipelineChaining(t *testing.T) {
	input := `
fungsi tambahSatu(x) { kembalikan x + 1 }
fungsi kaliDua(x) { kembalikan x * 2 }
fungsi kuadrat(x) { kembalikan x * x }

buat hasil = 5 |> tambahSatu |> kaliDua |> kuadrat
cetak("Hasil pipeline berantai:", hasil)
hasil
`
	assertParity(t, "Stress_PipelineChaining", input)
}

// ---------------------------------------------------------------------
// 8. STRESS TEST: Concurrent Fan-Out (bersamaan & goroutines)
// ---------------------------------------------------------------------

func TestStress_ConcurrencyFanOut(t *testing.T) {
	input := `
buat counter = 0
bersamaan {
    buat a = 1
    selama a <= 100 {
        a += 1
    }
    buat b = 1
    selama b <= 100 {
        b += 1
    }
}
cetak("Blok bersamaan selesai dieksekusi dengan aman")
`
	assertParity(t, "Stress_ConcurrencyFanOut", input)
}

// ---------------------------------------------------------------------
// 9. STRESS TEST: Parallel Goroutine Race Detector Test
// ---------------------------------------------------------------------

func TestStress_MultiThreadedStress(t *testing.T) {
	const workerCount = 10
	var wg sync.WaitGroup
	wg.Add(workerCount)

	progStr := `
fungsi hitung(n) {
    buat sum = 0
    buat i = 1
    selama i <= n {
        sum += i
        i += 1
    }
    kembalikan sum
}
hitung(1000)
`

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			l := parser.NewLexer(progStr)
			p := parser.NewParser(l)
			prog, err := p.ParseProgram()
			if err != nil {
				t.Errorf("worker %d parser err: %v", id, err)
				return
			}
			comp := compiler.NewCompiler()
			chunk, err := comp.Compile(prog)
			if err != nil {
				t.Errorf("worker %d compiler err: %v", id, err)
				return
			}
			vmInst := vm.New(chunk)
			res, err := vmInst.Run()
			if err != nil {
				t.Errorf("worker %d vm err: %v", id, err)
				return
			}
			if res.Inspect() != "500500" {
				t.Errorf("worker %d expected 500500, got %s", id, res.Inspect())
			}
		}(i)
	}

	wg.Wait()
}
