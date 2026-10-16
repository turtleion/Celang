package empathy_test

import (
	"strings"
	"testing"

	"celang/pkg/empathy"
	"celang/pkg/parser"
)

func TestEmpathy_SyntaxDiagnosticExtraction(t *testing.T) {
	sourceCode := `buat a = 1
buat b = 2
jika a > b {
    cetak("a")
} lain {
    cetak("b")
}
`
	rawErr := `parsing gagal dengan 3 kesalahan:
[Baris 5, Kolom 3] Diharapkan ':' atau 'itu' antara kunci dan nilai map
[Baris 5, Kolom 3] Tidak ada fungsi parse untuk token '}' (})
[Baris 7, Kolom 1] Tidak ada fungsi parse untuk token '}' (})`

	// Even if caller passes line=1, col=1, NewSyntaxDiagnostic should auto-extract real location
	diag := empathy.NewSyntaxDiagnostic(rawErr, 1, 1, sourceCode)

	if diag.Line != 5 {
		t.Errorf("expected line 5, got %d", diag.Line)
	}
	if diag.Column != 3 {
		t.Errorf("expected column 3, got %d", diag.Column)
	}
	if diag.OtherErrorsCount != 2 {
		t.Errorf("expected OtherErrorsCount 2, got %d", diag.OtherErrorsCount)
	}
	if !strings.Contains(diag.Message, "antara kunci dan nilai map") {
		t.Errorf("expected clean message without header, got: %s", diag.Message)
	}

	rendered := diag.Render()
	if !strings.Contains(rendered, "Baris 5, Kolom 3") {
		t.Errorf("expected rendered diagnostic to show Baris 5, Kolom 3, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "5 | } lain {") {
		t.Errorf("expected line 5 in rendered snippet, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "^ di sini") {
		t.Errorf("expected '^ di sini' marker, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Ditemukan juga 2 kemungkinan kesalahan turunan lainnya") {
		t.Errorf("expected other errors note, got:\n%s", rendered)
	}
}

func TestEmpathy_NewSyntaxDiagnosticFromError(t *testing.T) {
	sourceCode := `buat x = 10
batalkan "error"
`
	l := parser.NewLexer(sourceCode)
	p := parser.NewParser(l)
	_, err := p.ParseProgram()
	if err == nil {
		// Note: batalkan as statement outside transaction
	}

	// Test with ParseErrorList
	errList := parser.ParseErrorList{
		{
			Line:    2,
			Column:  1,
			Message: "Perintah 'batalkan' hanya dapat digunakan di dalam transaksi",
		},
		{
			Line:    3,
			Column:  1,
			Message: "Kesalahan kedua",
		},
	}

	diag := empathy.NewSyntaxDiagnosticFromError(errList, sourceCode)
	if diag.Line != 2 || diag.Column != 1 {
		t.Errorf("expected line 2 col 1, got %d:%d", diag.Line, diag.Column)
	}
	if diag.OtherErrorsCount != 1 {
		t.Errorf("expected 1 other error, got %d", diag.OtherErrorsCount)
	}
	if !strings.Contains(diag.FriendlyHint, "transaksi") {
		t.Errorf("expected rollback hint, got: %s", diag.FriendlyHint)
	}
}

func TestEmpathy_SmartHints(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		rawMsg     string
		line       int
		col        int
		expectHint string
	}{
		{
			name:       "Lain jika typo hint",
			source:     "jika a > 1 { \n } lain jika a == 0 { \n }",
			rawMsg:     "Tidak ada fungsi parse untuk token",
			line:       2,
			col:        3,
			expectHint: "selain jika",
		},
		{
			name:       "Lain else typo hint",
			source:     "jika a > 1 { \n } lain { \n }",
			rawMsg:     "Diharapkan ':' atau 'itu' antara kunci dan nilai map",
			line:       2,
			col:        3,
			expectHint: "selain",
		},
		{
			name:       "Function keyword typo",
			source:     "function hitung() { }",
			rawMsg:     "Tidak ada fungsi parse",
			line:       1,
			col:        1,
			expectHint: "fungsi",
		},
		{
			name:       "Print keyword typo",
			source:     "print('halo')",
			rawMsg:     "Variabel print tidak ditemukan",
			line:       1,
			col:        1,
			expectHint: "cetak",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag := empathy.NewSyntaxDiagnostic(tt.rawMsg, tt.line, tt.col, tt.source)
			if !strings.Contains(strings.ToLower(diag.FriendlyHint), strings.ToLower(tt.expectHint)) {
				t.Errorf("expected hint to mention %q, got: %q", tt.expectHint, diag.FriendlyHint)
			}
		})
	}
}
