package empathy

import (
	"fmt"
	"strings"

	"celang/pkg/evaluator"
	"celang/pkg/uast"
)

// Diagnostic provides a human-centered, empathetic error display.
type Diagnostic struct {
	Title        string
	Message      string
	Line         int
	Column       int
	FriendlyHint string
	SourceCode   string
}

// NewRuntimeDiagnostic constructs a friendly diagnostic from an evaluator error.
func NewRuntimeDiagnostic(err *evaluator.Error, sourceCode string) *Diagnostic {
	line := 1
	col := 1
	if err.Pos != nil {
		line = err.Pos.Line
		col = err.Pos.Column
	}

	hint := err.FriendlyHint
	if hint == "" {
		hint = suggestHint(err.Message)
	}

	return &Diagnostic{
		Title:        "⚠️ Kesalahan Saat Program Berjalan (Runtime)",
		Message:      err.Message,
		Line:         line,
		Column:       col,
		FriendlyHint: hint,
		SourceCode:   sourceCode,
	}
}

// NewSyntaxDiagnostic constructs a friendly diagnostic from a parser error.
func NewSyntaxDiagnostic(rawErr string, line, col int, sourceCode string) *Diagnostic {
	hint := suggestSyntaxHint(rawErr)
	return &Diagnostic{
		Title:        "🛑 Kesalahan Penulisan Kode (Syntax Error)",
		Message:      rawErr,
		Line:         line,
		Column:       col,
		FriendlyHint: hint,
		SourceCode:   sourceCode,
	}
}

// Render formats the diagnostic into a beautiful, empathetic terminal view.
func (d *Diagnostic) Render() string {
	var sb strings.Builder

	border := strings.Repeat("─", 50)
	sb.WriteString(fmt.Sprintf("\n╭─ %s %s\n", d.Title, border[:max(0, 50-len(d.Title))]))
	sb.WriteString(fmt.Sprintf("│ 📍 Lokasi: Baris %d, Kolom %d\n", d.Line, d.Column))
	sb.WriteString(fmt.Sprintf("│ 💬 Penjelasan: %s\n", d.Message))
	sb.WriteString("│\n")

	// Render code snippet if source code is available
	if d.SourceCode != "" && d.Line > 0 {
		lines := strings.Split(d.SourceCode, "\n")
		startLine := max(1, d.Line-1)
		endLine := min(len(lines), d.Line+1)

		for i := startLine; i <= endLine; i++ {
			lineIdx := i - 1
			if lineIdx >= 0 && lineIdx < len(lines) {
				lineContent := lines[lineIdx]
				gutter := fmt.Sprintf("│ %3d | ", i)
				sb.WriteString(gutter)
				sb.WriteString(lineContent)
				sb.WriteString("\n")

				if i == d.Line {
					indentLen := len(gutter) - 2
					pointerPad := strings.Repeat(" ", max(0, d.Column-1))
					sb.WriteString(fmt.Sprintf("│%s^ %s\n", strings.Repeat(" ", indentLen)+pointerPad, d.Message))
				}
			}
		}
	}

	if d.FriendlyHint != "" {
		sb.WriteString("│\n")
		sb.WriteString(fmt.Sprintf("│ 💡 Saran Solusi:\n│    %s\n", d.FriendlyHint))
	}

	sb.WriteString("╰" + strings.Repeat("─", 60) + "\n")
	return sb.String()
}

func suggestHint(msg string) string {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "belum didefinisikan") || strings.Contains(lower, "tidak ditemukan") {
		return "Apakah kamu lupa mendeklarasikan variabel ini dengan 'buat nama = ...' terlebih dahulu?"
	}
	if strings.Contains(lower, "immutable") || strings.Contains(lower, "tetap") {
		return "Variabel ini dibuat dengan 'konstanta'. Jika nilainya perlu diubah, ganti deklarasinya menjadi 'buat'."
	}
	if strings.Contains(lower, "pembagian dengan nol") {
		return "Dalam matematika, pembagian dengan 0 tidak terdefinisi. Pastikan penyebut bernilai selain 0."
	}
	if strings.Contains(lower, "pointer") {
		return "Pastikan kamu menggunakan '&variabel' untuk membuat pointer, dan '*pointer' untuk mengakses/mengubah nilainya."
	}
	return "Periksa kembali logika dan tipe data variabel yang terlibat dalam operasi ini."
}

func suggestSyntaxHint(raw string) string {
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "kurung") || strings.Contains(lower, "{") || strings.Contains(lower, "}") {
		return "Pastikan setiap kurung buka '{' memiliki kurung tutup '}' pasangannya."
	}
	if strings.Contains(lower, ":") || strings.Contains(lower, "itu") {
		return "Gunakan ':' atau 'itu' untuk memisahkan kunci dan nilai dalam map, atau untuk tipe data."
	}
	return "Periksa ejaan kata kunci atau tanda baca pada baris ini."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// EnhancePos annotates AST nodes with position for empathetic reporting.
func EnhancePos(node uast.Node, file string, line, col int) {
	meta := node.GetMetadata()
	if meta != nil && meta.Pos == nil {
		meta.Pos = &uast.Position{
			File:   file,
			Line:   line,
			Column: col,
		}
	}
}
