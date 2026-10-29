package parser

import (
	"strings"
)

// DetectUniversalDirective inspects the source code to determine if it starts with the universal mode directive.
// Supported directives on line 1 (ignoring comments/whitespace):
// "use universal", 'use universal', "use english", 'use english', "use standard", 'use standard'
func DetectUniversalDirective(source string) bool {
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "/*") {
			continue
		}
		clean := strings.TrimSuffix(trimmed, ";")
		return clean == `"use universal"` || clean == `'use universal'` ||
			clean == `"use english"` || clean == `'use english'` ||
			clean == `"use standard"` || clean == `'use standard'`
	}
	return false
}

// DetectStyleDirective inspects the beginning lines of source code to detect an in-code style directive.
// Supports:
// // @style: nama_style.style
// # style: nama_style
// "use style nama_style", "gaya nama_style.style"
func DetectStyleDirective(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if i > 10 { // Only look at top of file
			break
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Comment directives: // @style: xxx or # style: xxx
		if strings.HasPrefix(trimmed, "//") {
			commentContent := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
			if strings.HasPrefix(commentContent, "@style:") {
				return strings.TrimSpace(strings.TrimPrefix(commentContent, "@style:"))
			}
			if strings.HasPrefix(commentContent, "style:") {
				return strings.TrimSpace(strings.TrimPrefix(commentContent, "style:"))
			}
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			commentContent := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			if strings.HasPrefix(commentContent, "style:") || strings.HasPrefix(commentContent, "@style:") {
				return strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(commentContent, "@style:"), "style:"))
			}
			continue
		}
		// String literal directives: "use style xxx" or "gaya xxx"
		clean := strings.Trim(strings.TrimSuffix(trimmed, ";"), "\"'`")
		if strings.HasPrefix(clean, "use style ") {
			return strings.TrimSpace(strings.TrimPrefix(clean, "use style "))
		}
		if strings.HasPrefix(clean, "gaya ") {
			return strings.TrimSpace(strings.TrimPrefix(clean, "gaya "))
		}
	}
	return ""
}

// IndonesianToEnglishKeywords maps Indonesian Celang keywords to their English equivalents in Universal Mode.
var IndonesianToEnglishKeywords = map[string]string{
	"buat":          "let",
	"konstanta":     "const",
	"konstan":       "const",
	"tetap":         "const",
	"jika":          "if",
	"kalau":         "if",
	"selain jika":   "else if",
	"lain jika":     "else if",
	"selain":        "else",
	"lain":          "else",
	"maka":          "then (or omit)",
	"pilih":         "switch",
	"cocok":         "switch",
	"bila":          "case",
	"kasus":         "case",
	"lainnya":       "default",
	"selama":        "while",
	"untuk":         "for",
	"dalam":         "in",
	"fungsi":        "function (or fn / def)",
	"kembalikan":    "return",
	"mengembalikan": "return",
	"berhenti":      "break",
	"lanjut":        "continue",
	"transaksi":     "transaction",
	"batalkan":      "rollback",
	"pulihkan":      "recover",
	"pantau":        "watch",
	"reaktif":       "reactive",
	"coba":          "try",
	"tangkap":       "catch",
	"akhirnya":      "finally",
	"lempar":        "throw",
	"benar":         "true",
	"salah":         "false",
	"kosong":        "null",
	"dan":           "and",
	"atau":          "or",
	"bukan":         "not",
	"saat_selesai":  "defer",
	"titip":         "spawn",
	"jemput":        "await",
	"saluran":       "channel",
	"sembarang":     "any",
	"harus":         "must",
	"bentuk":        "struct / class",
	"rujuk":         "ref",
	"sebagai":       "as",
	"uji":           "test",
	"ekspor":        "export",
	"ambil":         "import",
	"impor":         "import",
	"dari":          "from",
}

// IsIndonesianKeyword checks if a word is an Indonesian-specific keyword and returns its English replacement.
func IsIndonesianKeyword(ident string) (string, bool) {
	repl, ok := IndonesianToEnglishKeywords[ident]
	return repl, ok
}
