package syntax_test

import (
	"strings"
	"testing"

	"celang/pkg/parser"
	"celang/pkg/syntax"
)

func TestFormatter_StyleTransformations(t *testing.T) {
	inputIndo := `
	buat harga = 50000
	harga = 75000

	fungsi cek_diskon(nilai itu Angka) mengembalikan Teks {
		jika (nilai > 60000) {
			kembalikan "Diskon"
		} selain {
			kembalikan "Normal"
		}
	}

	buat daftar = [1, 2, 3]
	untuk item dalam daftar {
		cetak(item)
	}

	cocok harga {
		kasus 75000:
			cetak("Cocok")
		lainnya:
			cetak("Lainnya")
	}

	cetak(harga)
	sambil {
		cetak("Latar belakang")
	}
	`

	// 1. Forward Parse: Text -> U-AST
	l := parser.NewLexer(inputIndo)
	p := parser.NewParser(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// 2. Reverse Parse into Indonesian Style
	indoFormatter := syntax.NewFormatter(syntax.IndonesianProfile())
	indoText := indoFormatter.Format(prog)
	t.Logf("Formatted Indonesian:\n%s", indoText)

	if !strings.Contains(indoText, "buat harga = 50000") ||
		!strings.Contains(indoText, "fungsi cek_diskon") ||
		!strings.Contains(indoText, "kembalikan \"Diskon\"") ||
		!strings.Contains(indoText, "cocok harga {") ||
		!strings.Contains(indoText, "kasus 75000:") ||
		!strings.Contains(indoText, "sambil {") {
		t.Fatalf("Indonesian formatted code missing expected keywords:\n%s", indoText)
	}

	// 3. Reverse Parse into Pythonic Style
	pyFormatter := syntax.NewFormatter(syntax.PythonicProfile())
	pyText := pyFormatter.Format(prog)
	t.Logf("Formatted Pythonic:\n%s", pyText)

	if !strings.Contains(pyText, "harga = 50000") ||
		!strings.Contains(pyText, "def cek_diskon") ||
		!strings.Contains(pyText, "return \"Diskon\"") ||
		!strings.Contains(pyText, "for item in daftar {") ||
		!strings.Contains(pyText, "match harga {") ||
		!strings.Contains(pyText, "case 75000:") ||
		!strings.Contains(pyText, "async {") ||
		!strings.Contains(pyText, "print(harga)") {
		t.Fatalf("Pythonic formatted code missing expected keywords:\n%s", pyText)
	}

	// 4. Reverse Parse into JavaScript Style (100% JS)
	jsFormatter := syntax.NewFormatter(syntax.JavaScriptProfile())
	jsText := jsFormatter.Format(prog)
	t.Logf("Formatted JavaScript (100%% JS):\n%s", jsText)

	if !strings.Contains(jsText, "let harga = 50000") ||
		!strings.Contains(jsText, "function cek_diskon") ||
		!strings.Contains(jsText, "for (let item of daftar) {") ||
		!strings.Contains(jsText, "switch (harga) {") ||
		!strings.Contains(jsText, "case 75000:") ||
		!strings.Contains(jsText, "return \"Diskon\"") {
		t.Fatalf("JavaScript formatted code missing expected keywords:\n%s", jsText)
	}

	// 5. Reverse Parse into Custom SDL Profile (Bahasa Sunda from file)
	sundaProfile, err := syntax.LoadProfileFromFile("../../examples/sunda.style")
	if err == nil {
		sundaFormatter := syntax.NewFormatter(sundaProfile)
		sundaText := sundaFormatter.Format(prog)
		t.Logf("Formatted Sunda:\n%s", sundaText)

		if !strings.Contains(sundaText, "jieun harga = 50000") ||
			!strings.Contains(sundaText, "balikkeun \"Diskon\"") ||
			!strings.Contains(sundaText, "cocogkeun harga {") ||
			!strings.Contains(sundaText, "bisi 75000:") ||
			!strings.Contains(sundaText, "bari {") {
			t.Fatalf("Sunda formatted code missing expected keywords:\n%s", sundaText)
		}
	}
}
