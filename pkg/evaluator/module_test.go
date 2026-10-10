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
