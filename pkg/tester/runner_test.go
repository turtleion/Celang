package tester

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"calc_uji.ce", true},
		{"calc_test.ce", true},
		{"calc.uji.ce", true},
		{"calc.test.ce", true},
		{"uji/aritmatika.ce", true},
		{"tests/math.ce", true},
		{"test/math.ce", true},
		{"main.ce", false},
		{"pkg/calc.ce", false},
		{"calc.go", false},
		{"test.txt", false},
	}

	for _, tt := range tests {
		got := IsTestFile(tt.path)
		if got != tt.expected {
			t.Errorf("IsTestFile(%q) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}

func TestRunFile_BlockSyntax(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "contoh_uji.ce")

	code := `
uji "penjumlahan dasar" {
	buat a = 10
	buat b = 20
	harus_sama(a + b, 30)
}

uji "perkalian desimal" {
	buat x = 2.5
	buat y = 4.0
	harus_sama(x * y, 10.0)
}
`
	if err := os.WriteFile(testPath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	res := RunFile(testPath, Options{})
	if res.Error != "" {
		t.Fatalf("unexpected file error: %s", res.Error)
	}

	if len(res.Cases) != 2 {
		t.Fatalf("expected 2 test cases, got %d", len(res.Cases))
	}

	if !res.Cases[0].Passed || res.Cases[0].Name != "penjumlahan dasar" {
		t.Errorf("case 0 failed: %+v", res.Cases[0])
	}
	if !res.Cases[1].Passed || res.Cases[1].Name != "perkalian desimal" {
		t.Errorf("case 1 failed: %+v", res.Cases[1])
	}
}

func TestRunFile_CallbackSyntax(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "callback_uji.ce")

	code := `
uji("callback lolos", () => {
	buat daftar = [1, 2, 3]
	harus_berisi(daftar, 2)
})

test("english callback", fn() {
	assert_true(true)
	assert_false(false)
	assert_null(null)
	assert_not_null("celang")
})
`
	if err := os.WriteFile(testPath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	res := RunFile(testPath, Options{})
	if res.Error != "" {
		t.Fatalf("unexpected file error: %s", res.Error)
	}

	if len(res.Cases) != 2 {
		t.Fatalf("expected 2 test cases, got %d", len(res.Cases))
	}
	if !res.Cases[0].Passed {
		t.Errorf("case 0 failed: %+v", res.Cases[0])
	}
	if !res.Cases[1].Passed {
		t.Errorf("case 1 failed: %+v", res.Cases[1])
	}
}

func TestRunFile_AssertionsComprehensive(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "semua_assertion_uji.ce")

	code := `
uji "semua variasi assertion" {
	// harus
	harus(10 > 5, "10 harus lebih besar dari 5")

	// harus_sama & harus_beda
	harus_sama([1, 2, 3], [1, 2, 3])
	harus_beda([1, 2], [1, 3])

	// harus_benar & harus_salah
	harus_benar(benar)
	harus_salah(salah)

	// harus_kosong & harus_tidak_kosong
	harus_kosong(kosong)
	harus_tidak_kosong("ada isinya")

	// harus_berisi (string, array, map)
	harus_berisi("halo celang", "celang")
	harus_berisi(["kopi", "teh"], "kopi")
	harus_berisi({"nama": "Antigravity", "umur": 1}, "nama")

	// harus_galat
	harus_galat(() => {
		lempar "terjadi kesalahan fatal"
	}, "kesalahan fatal")
}
`
	if err := os.WriteFile(testPath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	res := RunFile(testPath, Options{})
	if res.Error != "" {
		t.Fatalf("unexpected file error: %s", res.Error)
	}

	if len(res.Cases) != 1 {
		t.Fatalf("expected 1 test case, got %d", len(res.Cases))
	}
	if !res.Cases[0].Passed {
		t.Errorf("case failed unexpectedly: %+v", res.Cases[0])
	}
}

func TestRunFile_FailingAssertion(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "gagal_uji.ce")

	code := `
uji "tes yang sengaja gagal" {
	buat nilai = 42
	harus_sama(nilai, 100, "Nilai matematika salah")
}

uji "tes kedua yang lulus" {
	harus(benar)
}
`
	if err := os.WriteFile(testPath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	res := RunFile(testPath, Options{})
	if res.Error != "" {
		t.Fatalf("unexpected file error: %s", res.Error)
	}

	if len(res.Cases) != 2 {
		t.Fatalf("expected 2 test cases, got %d", len(res.Cases))
	}

	// Case 0 should fail
	if res.Cases[0].Passed {
		t.Errorf("expected case 0 to fail, but it passed")
	}
	if res.Cases[0].Expected != "100" || res.Cases[0].Actual != "42" {
		t.Errorf("unexpected diff: Expected=%q, Actual=%q", res.Cases[0].Expected, res.Cases[0].Actual)
	}

	// Case 1 should pass despite case 0 failing!
	if !res.Cases[1].Passed {
		t.Errorf("expected case 1 to pass, but it failed: %+v", res.Cases[1])
	}
}

func TestRun_WithFilterAndReporter(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "filter_uji.ce")

	code := `
uji "aritmatika tambah" {
	harus_sama(1 + 1, 2)
}

uji "string gabung" {
	harus_sama("a" + "b", "ab")
}
`
	if err := os.WriteFile(testPath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	summary, err := Run([]string{testPath}, Options{
		Filter:  "tambah",
		NoColor: true,
		Out:     &buf,
	})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if summary.TotalTests != 1 {
		t.Fatalf("expected 1 test to match filter 'tambah', got %d", summary.TotalTests)
	}
	if summary.Passed != 1 {
		t.Errorf("expected 1 passed, got %d", summary.Passed)
	}

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("aritmatika tambah")) {
		t.Errorf("output missing filtered test name: %s", output)
	}
	if bytes.Contains(buf.Bytes(), []byte("string gabung")) {
		t.Errorf("output should NOT contain filtered-out test: %s", output)
	}
}
