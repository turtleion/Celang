package evaluator_test

import (
	"os"
	"path/filepath"
	"testing"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
)

func TestModule_SingleFileDefaultPrivateAndPublic(t *testing.T) {
	tempDir := t.TempDir()

	p1Path := filepath.Join(tempDir, "program1.ce")
	p2Path := filepath.Join(tempDir, "program2.ce")

	// 1. Without publik! (default privat!)
	p1Content := `buat nama = "Kayasa"`
	if err := os.WriteFile(p1Path, []byte(p1Content), 0644); err != nil {
		t.Fatalf("Failed to write p1: %v", err)
	}

	p2Content := `ambil { nama } dari "./program1.ce"`
	l := parser.NewLexer(p2Content)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = p2Path
	result := evaluator.Eval(prog, env)

	if result.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected error when importing private variable, got %s", result.Inspect())
	}

	// 2. With publik!
	p1ContentPub := `publik! buat nama = "Kayasa"`
	if err := os.WriteFile(p1Path, []byte(p1ContentPub), 0644); err != nil {
		t.Fatalf("Failed to write p1 with publik: %v", err)
	}

	env2 := evaluator.NewEnvironment()
	env2.FilePath = p2Path
	result2 := evaluator.Eval(prog, env2)

	if result2 != nil && result2.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Unexpected error with publik!: %s", result2.Inspect())
	}

	val, ok := env2.Get("nama")
	if !ok || val.Inspect() != "Kayasa" {
		t.Fatalf("Expected 'Kayasa', got %+v", val)
	}
}

func TestModule_ModuleScopedVisibility(t *testing.T) {
	tempDir := t.TempDir()

	// Structure:
	// tempDir/
	//   Transport/
	//     celang.mod
	//     core.ce
	//     lib/
	//       avt.ce
	//   Code/
	//     program3.ce

	transportDir := filepath.Join(tempDir, "Transport")
	libDir := filepath.Join(transportDir, "lib")
	codeDir := filepath.Join(tempDir, "Code")

	if err := os.MkdirAll(libDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create celang.mod marker
	if err := os.WriteFile(filepath.Join(transportDir, "celang.mod"), []byte("modul Transport"), 0644); err != nil {
		t.Fatal(err)
	}

	// Transport/lib/avt.ce
	avtPath := filepath.Join(libDir, "avt.ce")
	avtContent := `
	privat-modul! fungsi komunikasikan(ip: Teks) {
		kembalikan "OK:" + ip
	}
	`
	if err := os.WriteFile(avtPath, []byte(avtContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Transport/core.ce (same module -> allowed)
	coreContent := `
	ambil { komunikasikan } dari "./lib/avt.ce"
	buat status = komunikasikan("127.0.0.1")
	`
	corePath := filepath.Join(transportDir, "core.ce")
	lCore := parser.NewLexer(coreContent)
	pCore := parser.NewParser(lCore)
	progCore, err := pCore.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	envCore := evaluator.NewEnvironment()
	envCore.FilePath = corePath
	resCore := evaluator.Eval(progCore, envCore)
	if resCore != nil && resCore.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("core.ce in same module failed: %s", resCore.Inspect())
	}

	statusVal, ok := envCore.Get("status")
	if !ok || statusVal.Inspect() != "OK:127.0.0.1" {
		t.Fatalf("Expected status='OK:127.0.0.1', got %+v", statusVal)
	}

	// Code/program3.ce (outside module -> denied)
	p3Content := `
	ambil { komunikasikan } dari "../Transport/lib/avt.ce"
	`
	p3Path := filepath.Join(codeDir, "program3.ce")
	lP3 := parser.NewLexer(p3Content)
	pP3 := parser.NewParser(lP3)
	progP3, err := pP3.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	envP3 := evaluator.NewEnvironment()
	envP3.FilePath = p3Path
	resP3 := evaluator.Eval(progP3, envP3)

	if resP3.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected error for program3 outside module, got %s", resP3.Inspect())
	}
}

func TestModule_AllowlistVisibility(t *testing.T) {
	tempDir := t.TempDir()

	p4Path := filepath.Join(tempDir, "program4.ce")
	p5Path := filepath.Join(tempDir, "program5.ce")
	p6Path := filepath.Join(tempDir, "program6.ce")

	// program4.ce grants allowlist only to program5.ce
	p4Content := `
	privat-["./program5.ce"]! buat nama = "Abdi"
	`
	if err := os.WriteFile(p4Path, []byte(p4Content), 0644); err != nil {
		t.Fatal(err)
	}

	// program5.ce: allowed!
	p5Content := `
	ambil { nama } dari "./program4.ce"
	`
	lP5 := parser.NewLexer(p5Content)
	pP5 := parser.NewParser(lP5)
	progP5, err := pP5.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	envP5 := evaluator.NewEnvironment()
	envP5.FilePath = p5Path
	resP5 := evaluator.Eval(progP5, envP5)
	if resP5 != nil && resP5.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("program5 (allowlisted) failed: %s", resP5.Inspect())
	}

	namaVal, ok := envP5.Get("nama")
	if !ok || namaVal.Inspect() != "Abdi" {
		t.Fatalf("Expected nama='Abdi', got %+v", namaVal)
	}

	// program6.ce: not on allowlist -> denied!
	p6Content := `
	ambil { nama } dari "./program4.ce"
	`
	lP6 := parser.NewLexer(p6Content)
	pP6 := parser.NewParser(lP6)
	progP6, err := pP6.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	envP6 := evaluator.NewEnvironment()
	envP6.FilePath = p6Path
	resP6 := evaluator.Eval(progP6, envP6)

	if resP6.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected error for program6 not in allowlist, got %s", resP6.Inspect())
	}
}

func TestModule_StandaloneExportAndAliases(t *testing.T) {
	tempDir := t.TempDir()

	mathPath := filepath.Join(tempDir, "matematika.ce")
	appPath := filepath.Join(tempDir, "app.ce")

	mathContent := `
	fungsi tambah(a, b) {
		kembalikan a + b
	}

	fungsi kali(a, b) {
		kembalikan a * b
	}

	publik! { tambah, kali }
	`
	if err := os.WriteFile(mathPath, []byte(mathContent), 0644); err != nil {
		t.Fatal(err)
	}

	appContent := `
	ambil { tambah sebagai penjumlahan } dari "./matematika.ce"
	ambil * sebagai math dari "./matematika.ce"

	buat hasil_tambah = penjumlahan(5, 7)
	buat hasil_kali = math.kali(3, 4)
	`
	lApp := parser.NewLexer(appContent)
	pApp := parser.NewParser(lApp)
	progApp, err := pApp.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	envApp := evaluator.NewEnvironment()
	envApp.FilePath = appPath
	resApp := evaluator.Eval(progApp, envApp)
	if resApp != nil && resApp.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("app evaluation failed: %s", resApp.Inspect())
	}

	tambahVal, ok := envApp.Get("hasil_tambah")
	if !ok || tambahVal.(*evaluator.Integer).Value != 12 {
		t.Fatalf("Expected hasil_tambah=12, got %+v", tambahVal)
	}

	kaliVal, ok := envApp.Get("hasil_kali")
	if !ok || kaliVal.(*evaluator.Integer).Value != 12 {
		t.Fatalf("Expected hasil_kali=12, got %+v", kaliVal)
	}
}

func TestModule_CleanPackageAndStdlibImports(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Unquoted stdlib imports
	script := `
	ambil { akar, PI } dari matematika
	impor matematika sebagai math
	dari json ambil uraikan
	dari matematika ambil *

	buat a = akar(16)
	buat b = math.akar(25)
	buat c = uraikan("{\"kunci\": 42}")
	buat d = floor(3.9)
	`
	l := parser.NewLexer(script)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filepath.Join(tempDir, "main.ce")
	res := evaluator.Eval(prog, env)
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error: %s", res.Inspect())
	}

	aVal, ok := env.Get("a")
	if !ok || aVal.Inspect() != "4" {
		t.Fatalf("Expected a=4, got %+v", aVal)
	}

	bVal, ok := env.Get("b")
	if !ok || bVal.Inspect() != "5" {
		t.Fatalf("Expected b=5, got %+v", bVal)
	}

	dVal, ok := env.Get("d")
	if !ok || dVal.Inspect() != "3" {
		t.Fatalf("Expected d=3, got %+v", dVal)
	}

	// 2. Unquoted package from celang_modules/<pkg>/utama.ce
	pkgDir := filepath.Join(tempDir, "celang_modules", "kalkulator")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}
	pkgMain := filepath.Join(pkgDir, "utama.ce")
	pkgCode := `
	publik! fungsi kali2(x) {
		kembalikan x * 2
	}
	`
	if err := os.WriteFile(pkgMain, []byte(pkgCode), 0644); err != nil {
		t.Fatal(err)
	}

	pkgScript := `
	ambil { kali2 } dari kalkulator
	buat hasil = kali2(21)
	`
	l2 := parser.NewLexer(pkgScript)
	p2 := parser.NewParser(l2)
	prog2, err := p2.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env2 := evaluator.NewEnvironment()
	env2.FilePath = filepath.Join(tempDir, "app.ce")
	res2 := evaluator.Eval(prog2, env2)
	if res2 != nil && res2.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error for package: %s", res2.Inspect())
	}

	hasilVal, ok := env2.Get("hasil")
	if !ok || hasilVal.(*evaluator.Integer).Value != 42 {
		t.Fatalf("Expected hasil=42, got %+v", hasilVal)
	}
}

func TestModule_DirectoryEntrypointResolution(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Folder with modul.ce
	dir1 := filepath.Join(tempDir, "komp_modul")
	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "modul.ce"), []byte(`publik! buat pesan = "dari modul.ce"`), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. Folder with utama.ce
	dir2 := filepath.Join(tempDir, "komp_utama")
	if err := os.MkdirAll(dir2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "utama.ce"), []byte(`publik! buat pesan = "dari utama.ce"`), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Folder with index.ce
	dir3 := filepath.Join(tempDir, "komp_index")
	if err := os.MkdirAll(dir3, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir3, "index.ce"), []byte(`publik! buat pesan = "dari index.ce"`), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Folder with <dirname>.ce
	dir4 := filepath.Join(tempDir, "komp_nama")
	if err := os.MkdirAll(dir4, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir4, "komp_nama.ce"), []byte(`publik! buat pesan = "dari komp_nama.ce"`), 0644); err != nil {
		t.Fatal(err)
	}

	consumerScript := `
	ambil { pesan as p1 } dari "./komp_modul"
	ambil { pesan as p2 } dari "./komp_utama"
	ambil { pesan as p3 } dari "./komp_index"
	ambil { pesan as p4 } dari "./komp_nama"
	`
	l := parser.NewLexer(consumerScript)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filepath.Join(tempDir, "main.ce")
	res := evaluator.Eval(prog, env)
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error: %s", res.Inspect())
	}

	for _, tc := range []struct {
		name     string
		expected string
	}{
		{"p1", "dari modul.ce"},
		{"p2", "dari utama.ce"},
		{"p3", "dari index.ce"},
		{"p4", "dari komp_nama.ce"},
	} {
		val, ok := env.Get(tc.name)
		if !ok || val.Inspect() != tc.expected {
			t.Errorf("Expected %s=%q, got %+v", tc.name, tc.expected, val)
		}
	}
}

func TestModule_ReExportBarrel(t *testing.T) {
	tempDir := t.TempDir()

	// sub1.ce
	sub1 := filepath.Join(tempDir, "sub1.ce")
	sub1Code := `
	publik! buat angka_a = 100
	buat rahasia = 999
	publik! fungsi kali(x) { kembalikan x * 2 }
	`
	if err := os.WriteFile(sub1, []byte(sub1Code), 0644); err != nil {
		t.Fatal(err)
	}

	// sub2.ce
	sub2 := filepath.Join(tempDir, "sub2.ce")
	sub2Code := `
	publik! buat kata = "celang"
	publik! buat versi = "2.0"
	`
	if err := os.WriteFile(sub2, []byte(sub2Code), 0644); err != nil {
		t.Fatal(err)
	}

	// barrel.ce re-exports sub1 and sub2
	barrel := filepath.Join(tempDir, "barrel.ce")
	barrelCode := `
	ekspor { angka_a, kali sebagai ganda } dari "./sub1.ce"
	ekspor * dari "./sub2.ce"
	`
	if err := os.WriteFile(barrel, []byte(barrelCode), 0644); err != nil {
		t.Fatal(err)
	}

	// consumer.ce imports through barrel
	consumerScript := `
	ambil { angka_a, ganda, kata, versi } dari "./barrel.ce"
	buat hasil_kali = ganda(21)
	`
	l := parser.NewLexer(consumerScript)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filepath.Join(tempDir, "consumer.ce")
	res := evaluator.Eval(prog, env)
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error: %s", res.Inspect())
	}

	if val, ok := env.Get("angka_a"); !ok || val.(*evaluator.Integer).Value != 100 {
		t.Errorf("Expected angka_a=100, got %+v", val)
	}
	if val, ok := env.Get("hasil_kali"); !ok || val.(*evaluator.Integer).Value != 42 {
		t.Errorf("Expected hasil_kali=42, got %+v", val)
	}
	if val, ok := env.Get("kata"); !ok || val.Inspect() != "celang" {
		t.Errorf("Expected kata='celang', got %+v", val)
	}
	if val, ok := env.Get("versi"); !ok || val.Inspect() != "2.0" {
		t.Errorf("Expected versi='2.0', got %+v", val)
	}

	// Test privacy barrier: importing private 'rahasia' through barrel must fail
	leakScript := `ambil { rahasia } dari "./barrel.ce"`
	lLeak := parser.NewLexer(leakScript)
	pLeak := parser.NewParser(lLeak)
	progLeak, _ := pLeak.ParseProgram()
	envLeak := evaluator.NewEnvironment()
	envLeak.FilePath = filepath.Join(tempDir, "leak.ce")
	resLeak := evaluator.Eval(progLeak, envLeak)
	if resLeak == nil || resLeak.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected error when importing private variable through barrel, got %+v", resLeak)
	}
}

func TestModule_UniversalReExport(t *testing.T) {
	tempDir := t.TempDir()

	sub := filepath.Join(tempDir, "math_sub.ce")
	subCode := `"use universal"
export const PI = 314
export function double_it(n) {
    return n * 2
}
`
	if err := os.WriteFile(sub, []byte(subCode), 0644); err != nil {
		t.Fatal(err)
	}

	barrel := filepath.Join(tempDir, "math_barrel.ce")
	barrelCode := `"use universal"
export { PI, double_it as times_two } from "./math_sub.ce"
`
	if err := os.WriteFile(barrel, []byte(barrelCode), 0644); err != nil {
		t.Fatal(err)
	}

	consumerScript := `"use universal"
import { PI, times_two } from "./math_barrel.ce"
const res = times_two(PI)
`
	l := parser.NewLexer(consumerScript)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filepath.Join(tempDir, "consumer_univ.ce")
	res := evaluator.Eval(prog, env)
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error: %s", res.Inspect())
	}

	resVal, ok := env.Get("res")
	if !ok || resVal.(*evaluator.Integer).Value != 628 {
		t.Fatalf("Expected res=628, got %+v", resVal)
	}
}

func TestModule_DynamicImport(t *testing.T) {
	tempDir := t.TempDir()

	modulePath := filepath.Join(tempDir, "helper.ce")
	moduleCode := `
	publik! buat greeting = "halo dari modul"
	publik! fungsi jumlah(a, b) {
		kembalikan a + b
	}
	buat rahasia = 1234
	`
	if err := os.WriteFile(modulePath, []byte(moduleCode), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Dynamic import in Indonesian mode
	script := `
	buat m = impor_dinamis("./helper.ce")
	buat salam = m.greeting
	buat total = m.jumlah(10, 20)
	buat ada_rahasia = m.kunci().ada_yang(k => k == "rahasia")
	`
	l := parser.NewLexer(script)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	env := evaluator.NewEnvironment()
	env.FilePath = filepath.Join(tempDir, "main.ce")
	evaluator.RegisterBuiltins(env, nil)
	res := evaluator.Eval(prog, env)
	if res != nil && res.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Eval error in dynamic import: %s", res.Inspect())
	}

	salamVal, ok := env.Get("salam")
	if !ok || salamVal.Inspect() != "halo dari modul" {
		t.Errorf("Expected salam='halo dari modul', got %+v", salamVal)
	}

	totalVal, ok := env.Get("total")
	if !ok || totalVal.(*evaluator.Integer).Value != 30 {
		t.Errorf("Expected total=30, got %+v", totalVal)
	}

	adaRahasiaVal, ok := env.Get("ada_rahasia")
	if !ok || adaRahasiaVal != evaluator.FALSE {
		t.Errorf("Expected ada_rahasia=false, got %+v", adaRahasiaVal)
	}

	// 2. Dynamic import in Universal mode
	univScript := `"use universal"
let mod = import_dynamic("./helper.ce")
let msg = mod.greeting
let sum = mod.jumlah(5, 7)
`
	lUniv := parser.NewLexer(univScript)
	pUniv := parser.NewParser(lUniv)
	progUniv, err := pUniv.ParseProgram()
	if err != nil {
		t.Fatalf("Universal parse error: %v", err)
	}

	envUniv := evaluator.NewEnvironment()
	envUniv.FilePath = filepath.Join(tempDir, "main_univ.ce")
	evaluator.RegisterBuiltins(envUniv, nil)
	resUniv := evaluator.Eval(progUniv, envUniv)
	if resUniv != nil && resUniv.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Universal dynamic import eval error: %s", resUniv.Inspect())
	}

	msgVal, ok := envUniv.Get("msg")
	if !ok || msgVal.Inspect() != "halo dari modul" {
		t.Errorf("Expected msg='halo dari modul', got %+v", msgVal)
	}

	sumVal, ok := envUniv.Get("sum")
	if !ok || sumVal.(*evaluator.Integer).Value != 12 {
		t.Errorf("Expected sum=12, got %+v", sumVal)
	}

	// 3. Universal mode rejects Indonesian 'impor_dinamis'
	rejectScript := `"use universal"
let m = impor_dinamis("./helper.ce")
`
	lRej := parser.NewLexer(rejectScript)
	pRej := parser.NewParser(lRej)
	progRej, _ := pRej.ParseProgram()
	envRej := evaluator.NewEnvironment()
	envRej.FilePath = filepath.Join(tempDir, "reject.ce")
	evaluator.RegisterBuiltins(envRej, nil)
	resRej := evaluator.Eval(progRej, envRej)
	if resRej == nil || resRej.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("Expected universal mode rejection for 'impor_dinamis', got %+v", resRej)
	}
}

