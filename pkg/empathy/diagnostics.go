package empathy

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"celang/pkg/evaluator"
	"celang/pkg/parser"
	"celang/pkg/uast"
)

// Diagnostic provides a human-centered, empathetic error display.
type Diagnostic struct {
	Title            string
	Message          string
	Line             int
	Column           int
	FriendlyHint     string
	SourceCode       string
	OtherErrorsCount int
	IsUniversal      bool
}

// NewRuntimeDiagnostic constructs a friendly diagnostic from an evaluator error.
func NewRuntimeDiagnostic(err *evaluator.Error, sourceCode string) *Diagnostic {
	line := 1
	col := 1
	if err.Pos != nil {
		line = err.Pos.Line
		col = err.Pos.Column
	}

	isUniversal := parser.DetectUniversalDirective(sourceCode)

	hint := err.FriendlyHint
	if hint == "" {
		if isUniversal {
			hint = suggestEnglishHint(err.Message)
		} else {
			hint = suggestHint(err.Message)
		}
	}

	title := "⚠️ Kesalahan Saat Program Berjalan (Runtime)"
	if isUniversal {
		title = "⚠️ Runtime Error"
	}

	return &Diagnostic{
		Title:        title,
		Message:      err.Message,
		Line:         line,
		Column:       col,
		FriendlyHint: hint,
		SourceCode:   sourceCode,
		IsUniversal:  isUniversal,
	}
}

// parseRawSyntaxError attempts to extract the primary error location and message from formatted strings.
func parseRawSyntaxError(raw string) (line, col int, cleanMsg string, otherCount int) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	re := regexp.MustCompile(`\[Baris\s+(\d+),\s+Kolom\s+(\d+)\]\s+(.*)`)

	var errorLines []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if re.MatchString(l) {
			errorLines = append(errorLines, l)
		}
	}

	if len(errorLines) > 0 {
		m := re.FindStringSubmatch(errorLines[0])
		if len(m) >= 4 {
			lNum, _ := strconv.Atoi(m[1])
			cNum, _ := strconv.Atoi(m[2])
			msg := m[3]
			return lNum, cNum, msg, len(errorLines) - 1
		}
	}

	return 1, 1, raw, 0
}

// NewSyntaxDiagnostic constructs a friendly diagnostic from a parser error string and position.
func NewSyntaxDiagnostic(rawErr string, line, col int, sourceCode string) *Diagnostic {
	actualLine := line
	actualCol := col
	cleanMsg := rawErr
	otherCount := 0

	// If position was unsupplied or default (1, 1), extract real position from rawErr
	if actualLine <= 1 && actualCol <= 1 {
		pLine, pCol, pMsg, count := parseRawSyntaxError(rawErr)
		if pLine > 1 || pCol > 1 {
			actualLine = pLine
			actualCol = pCol
			cleanMsg = pMsg
			otherCount = count
		} else if strings.Contains(rawErr, "[Baris ") {
			cleanMsg = pMsg
			otherCount = count
		}
	} else if strings.Contains(rawErr, "[Baris ") {
		_, _, pMsg, count := parseRawSyntaxError(rawErr)
		cleanMsg = pMsg
		otherCount = count
	}

	isUniversal := parser.DetectUniversalDirective(sourceCode)
	hint := suggestSyntaxHint(cleanMsg, actualLine, actualCol, sourceCode)

	title := "🛑 Kesalahan Penulisan Kode (Syntax Error)"
	if isUniversal {
		title = "🛑 Syntax Error"
	}

	return &Diagnostic{
		Title:            title,
		Message:          cleanMsg,
		Line:             actualLine,
		Column:           actualCol,
		FriendlyHint:     hint,
		SourceCode:       sourceCode,
		OtherErrorsCount: otherCount,
		IsUniversal:      isUniversal,
	}
}

// NewSyntaxDiagnosticFromError constructs a friendly diagnostic directly from any parser error.
func NewSyntaxDiagnosticFromError(err error, sourceCode string) *Diagnostic {
	if err == nil {
		return nil
	}

	isUniversal := parser.DetectUniversalDirective(sourceCode)
	title := "🛑 Kesalahan Penulisan Kode (Syntax Error)"
	if isUniversal {
		title = "🛑 Syntax Error"
	}

	if pel, ok := err.(parser.ParseErrorList); ok && len(pel) > 0 {
		first := pel[0]
		hint := suggestSyntaxHint(first.Message, first.Line, first.Column, sourceCode)
		return &Diagnostic{
			Title:            title,
			Message:          first.Message,
			Line:             first.Line,
			Column:           first.Column,
			FriendlyHint:     hint,
			SourceCode:       sourceCode,
			OtherErrorsCount: len(pel) - 1,
			IsUniversal:      isUniversal,
		}
	}

	return NewSyntaxDiagnostic(err.Error(), 1, 1, sourceCode)
}

// Render formats the diagnostic into a beautiful, empathetic terminal view.
func (d *Diagnostic) Render() string {
	var sb strings.Builder

	border := strings.Repeat("─", 50)
	sb.WriteString(fmt.Sprintf("\n╭─ %s %s\n", d.Title, border[:max(0, 50-len(d.Title))]))
	if d.IsUniversal {
		sb.WriteString(fmt.Sprintf("│ 📍 Location: Line %d, Column %d\n", d.Line, d.Column))
		sb.WriteString(fmt.Sprintf("│ 💬 Explanation: %s\n", d.Message))
		if d.OtherErrorsCount > 0 {
			sb.WriteString(fmt.Sprintf("│    (Found %d additional secondary error(s))\n", d.OtherErrorsCount))
		}
	} else {
		sb.WriteString(fmt.Sprintf("│ 📍 Lokasi: Baris %d, Kolom %d\n", d.Line, d.Column))
		sb.WriteString(fmt.Sprintf("│ 💬 Penjelasan: %s\n", d.Message))
		if d.OtherErrorsCount > 0 {
			sb.WriteString(fmt.Sprintf("│    (Ditemukan juga %d kemungkinan kesalahan turunan lainnya)\n", d.OtherErrorsCount))
		}
	}
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
					pointerText := "^ di sini"
					if d.IsUniversal {
						pointerText = "^ here"
					}
					sb.WriteString(fmt.Sprintf("│%s%s\n", strings.Repeat(" ", indentLen)+pointerPad, pointerText))
				}
			}
		}
	}

	if d.FriendlyHint != "" {
		sb.WriteString("│\n")
		header := "💡 Saran Solusi:"
		if d.IsUniversal {
			header = "💡 Suggested Solution:"
		}
		sb.WriteString(fmt.Sprintf("│ %s\n│    %s\n", header, d.FriendlyHint))
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

func suggestEnglishHint(msg string) string {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "not defined") || strings.Contains(lower, "not found") {
		return "Ensure this variable is declared with 'let name = ...' or 'const name = ...' beforehand."
	}
	if strings.Contains(lower, "immutable") || strings.Contains(lower, "constant") {
		return "This variable was declared with 'const'. Change it to 'let' if you need to reassign it."
	}
	if strings.Contains(lower, "division") || strings.Contains(lower, "zero") {
		return "Division by zero is undefined in mathematics. Make sure the denominator is not zero."
	}
	return "Check your variable types and logic involved in this operation."
}

func suggestSyntaxHint(msg string, line, col int, sourceCode string) string {
	lowerMsg := strings.ToLower(msg)

	if parser.DetectUniversalDirective(sourceCode) {
		if strings.Contains(lowerMsg, "is not allowed in universal mode") {
			return "Universal Mode requires standard English syntax and keywords."
		}
		if strings.Contains(lowerMsg, "'}'") || strings.Contains(lowerMsg, "'{'") || strings.Contains(lowerMsg, "brace") {
			return "There is an unclosed '{' or an extra '}'. Ensure every opening brace has a matching closing brace."
		}
		if strings.Contains(lowerMsg, "token ')'") || strings.Contains(lowerMsg, "token '('") {
			return "Ensure every opening parenthesis '(' has a matching closing parenthesis ')'."
		}
		if strings.Contains(lowerMsg, "token ']'") || strings.Contains(lowerMsg, "token '['") {
			return "Ensure every opening bracket '[' has a matching closing bracket ']'."
		}
		return "Check syntax, keyword spelling, and punctuation on this line."
	}

	var currentLine, prevLine string
	if sourceCode != "" && line > 0 {
		srcLines := strings.Split(sourceCode, "\n")
		if line-1 < len(srcLines) {
			currentLine = strings.TrimSpace(srcLines[line-1])
		}
		if line-2 >= 0 && line-2 < len(srcLines) {
			prevLine = strings.TrimSpace(srcLines[line-2])
		}
	}

	// 1. Check for alternative/misspelled keywords in current and previous lines
	combinedContext := prevLine + " " + currentLine
	if strings.Contains(combinedContext, "lain jika") {
		return "Kata kunci percabangan 'else if' di Celang adalah 'selain jika' atau 'elif'."
	}
	if strings.Contains(combinedContext, "lain {") || strings.Contains(combinedContext, "} lain") {
		return "Kata kunci percabangan 'else' di Celang adalah 'selain' atau 'else' (bukan 'lain')."
	}
	if strings.Contains(combinedContext, "function ") {
		return "Di Celang, gunakan kata kunci 'fungsi' atau 'def' untuk membuat fungsi."
	}
	if strings.Contains(combinedContext, "kelas ") || strings.Contains(combinedContext, "class ") {
		return "Di Celang, gunakan kata kunci 'bentuk' untuk membuat struktur data atau class (contoh: bentuk Pengguna { ... })."
	}
	if strings.Contains(combinedContext, "print(") || strings.Contains(combinedContext, "console.log(") {
		return "Gunakan fungsi 'cetak(...)' untuk menampilkan teks ke layar."
	}

	// 2. Map vs block ambiguity
	if strings.Contains(lowerMsg, "kunci dan nilai map") || strings.Contains(lowerMsg, "antara kunci dan nilai") {
		if strings.HasSuffix(prevLine, "{") || strings.HasSuffix(currentLine, "{") {
			return "Kurung kurawal '{' di baris ini dianggap sebagai Map/Kamus karena kata sebelumnya bukan pernyataan yang valid. Periksa apakah kamu salah mengetik kata kunci percabangan atau nama fungsi."
		}
		return "Gunakan ':' atau 'itu' untuk memisahkan kunci dan nilai dalam map (contoh: { nama: 'Budi', umur: 25 })."
	}

	// 3. Transactions & Rollback
	if strings.Contains(lowerMsg, "rollback") || strings.Contains(lowerMsg, "batalkan") {
		return "Perintah 'batalkan' hanya dapat digunakan di dalam blok 'transaksi { ... }'. Pastikan blok transaksi dibuka dan ditutup dengan kurung kurawal yang benar."
	}
	if strings.Contains(lowerMsg, "recover") || strings.Contains(lowerMsg, "pulihkan") {
		return "Kata kunci 'pulihkan' harus menempel setelah blok transaksi: 'transaksi { ... } pulihkan (alasan) { ... }'."
	}

	// 4. Reactivity & Watchers
	if strings.Contains(lowerMsg, "watch") || strings.Contains(lowerMsg, "pantau") {
		return "Gunakan format: 'pantau target (baru, lama) { ... }' untuk mengamati perubahan nilai variabel."
	}
	if strings.Contains(lowerMsg, "reactive") || strings.Contains(lowerMsg, "reaktif") {
		return "Gunakan format: 'reaktif nama = ekspresi' untuk membuat variabel reaktif turunan."
	}

	// 5. Braces & Parentheses
	if strings.Contains(lowerMsg, "'}'") || strings.Contains(lowerMsg, "'{'") || strings.Contains(lowerMsg, "kurung kurawal") {
		return "Ada kurung kurawal buka '{' yang belum ditutup, atau kurung tutup '}' yang berlebih. Pastikan setiap blok kurung kurawal memiliki kurung tutup pasangannya."
	}
	if strings.Contains(lowerMsg, "diharapkan token ')'") || strings.Contains(lowerMsg, "diharapkan token '('") {
		return "Pastikan setiap tanda kurung buka '(' memiliki tanda kurung tutup ')' pasangannya."
	}
	if strings.Contains(lowerMsg, "diharapkan token ']'") || strings.Contains(lowerMsg, "diharapkan token '['") {
		return "Pastikan setiap kurung siku buka '[' memiliki kurung siku tutup ']' pasangannya."
	}

	// 6. Generic syntax hints
	if strings.Contains(lowerMsg, "tidak ada fungsi parse") {
		return "Simbol atau kata ini tidak diharapkan muncul di posisi ini. Periksa apakah baris sebelumnya belum selesai atau ada tanda kurung/koma yang kurang."
	}

	return "Periksa kembali ejaan kata kunci atau tanda baca pada baris ini."
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
