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
		!strings.Contains(pyText, "async {") ||
		!strings.Contains(pyText, "print(harga)") {
		t.Fatalf("Pythonic formatted code missing expected keywords:\n%s", pyText)
	}

	// 4. Reverse Parse into Classic Style
	classicFormatter := syntax.NewFormatter(syntax.ClassicProfile())
	classicText := classicFormatter.Format(prog)
	t.Logf("Formatted Classic:\n%s", classicText)

	if !strings.Contains(classicText, "let harga = 50000") ||
		!strings.Contains(classicText, "fn cek_diskon") ||
		!strings.Contains(classicText, "return \"Diskon\"") {
		t.Fatalf("Classic formatted code missing expected keywords:\n%s", classicText)
	}
}
