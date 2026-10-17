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

// IndonesianToEnglishKeywords maps Indonesian Celang keywords to their English equivalents in Universal Mode.
var IndonesianToEnglishKeywords = map[string]string{
	"buat":          "let",
	"konstanta":     "const",
	"konstan":       "const",
	"tetap":         "const",
	"jika":          "if",
	"kalau":         "if",
	"bila":          "if",
	"selain jika":   "else if",
	"lain jika":     "else if",
	"selain":        "else",
	"lain":          "else",
	"maka":          "then (or omit)",
	"cocok":         "switch",
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
	"titip":         "defer",
	"jemput":        "await",
	"bentuk":        "struct / class",
	"rujuk":         "ref",
	"sebagai":       "as",
}

// IsIndonesianKeyword checks if a word is an Indonesian-specific keyword and returns its English replacement.
func IsIndonesianKeyword(ident string) (string, bool) {
	repl, ok := IndonesianToEnglishKeywords[ident]
	return repl, ok
}
