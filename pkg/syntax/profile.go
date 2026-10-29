package syntax

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"celang/pkg/parser"
)

// StyleProfile defines how keywords and block structures are mapped to U-AST.
// It acts as the "CSS for Syntax" in Celang, supporting both categorized (per-bab) and flat JSON.
type StyleProfile struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Version     string `json:"version,omitempty"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`

	// Deklarasi
	VarDecl    string `json:"var_decl"`
	ConstDecl  string `json:"const_decl"`
	TypeAssign string `json:"type_assign"`
	Struct     string `json:"struct,omitempty"`
	Ref        string `json:"ref,omitempty"`
	As         string `json:"as,omitempty"`

	// Percabangan
	If      string `json:"if"`
	ElseIf  string `json:"else_if"`
	Else    string `json:"else"`
	Then    string `json:"then,omitempty"`
	Switch  string `json:"switch"`
	Case    string `json:"case"`
	Default string `json:"default"`

	// Perulangan
	While string `json:"while"`
	For   string `json:"for"`
	In    string `json:"in"`

	// Fungsi
	Func        string `json:"func"`
	Return      string `json:"return"`
	ReturnArrow string `json:"return_arrow"`
	Defer       string `json:"defer,omitempty"`

	// Penanganan Galat
	Try     string `json:"try,omitempty"`
	Catch   string `json:"catch,omitempty"`
	Finally string `json:"finally,omitempty"`
	Throw   string `json:"throw,omitempty"`

	// Operator & Logika
	And          string `json:"and"`
	Or           string `json:"or"`
	Not          string `json:"not"`
	NullCoalesce string `json:"null_coalesce,omitempty"`
	Maybe        string `json:"maybe,omitempty"`

	// Literal
	True  string `json:"true"`
	False string `json:"false"`
	Null  string `json:"null"`

	// Konkurensi
	Async    string `json:"async"`
	Parallel string `json:"parallel"`

	// Format & Tata Letak
	BlockStyle  string `json:"block_style"`
	IndentWidth int    `json:"indent_width"`
	QuoteStyle  string `json:"quote_style,omitempty"`
}

// UnmarshalJSON deserializes a StyleProfile from either categorized (per-bab) or flat JSON.
func (p *StyleProfile) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	isCategorized := false
	for _, k := range []string{
		"metadata", "meta",
		"deklarasi", "declaration", "declarations",
		"percabangan", "control_flow", "branching",
		"perulangan", "loops", "loop",
		"fungsi", "functions", "function",
		"penanganan_galat", "error_handling", "errors",
		"operator_logika", "operators", "operator",
		"literal", "literals",
		"konkurensi", "concurrency",
		"format", "formatting",
	} {
		if _, ok := raw[k]; ok {
			isCategorized = true
			break
		}
	}

	if !isCategorized {
		type FlatProfile StyleProfile
		var flat FlatProfile
		if err := json.Unmarshal(data, &flat); err != nil {
			return err
		}
		*p = StyleProfile(flat)
		return nil
	}

	getSection := func(keys ...string) map[string]interface{} {
		for _, k := range keys {
			if msg, ok := raw[k]; ok {
				var sec map[string]interface{}
				if err := json.Unmarshal(msg, &sec); err == nil {
					return sec
				}
			}
		}
		return nil
	}

	getString := func(sec map[string]interface{}, keys ...string) string {
		if sec == nil {
			return ""
		}
		for _, k := range keys {
			if v, ok := sec[k]; ok {
				if s, ok := v.(string); ok {
					return s
				}
				if arr, ok := v.([]interface{}); ok {
					var parts []string
					for _, item := range arr {
						if s, ok := item.(string); ok {
							parts = append(parts, s)
						}
					}
					return strings.Join(parts, ", ")
				}
			}
		}
		return ""
	}

	getInt := func(sec map[string]interface{}, keys ...string) int {
		if sec == nil {
			return 0
		}
		for _, k := range keys {
			if v, ok := sec[k]; ok {
				if f, ok := v.(float64); ok {
					return int(f)
				}
			}
		}
		return 0
	}

	metaSec := getSection("metadata", "meta")
	p.Name = getString(metaSec, "name")
	p.DisplayName = getString(metaSec, "display_name", "title")
	p.Version = getString(metaSec, "version")
	p.Author = getString(metaSec, "author")
	p.Description = getString(metaSec, "description")

	dekSec := getSection("deklarasi", "declaration", "declarations")
	p.VarDecl = getString(dekSec, "var", "var_decl")
	p.ConstDecl = getString(dekSec, "const", "const_decl")
	p.TypeAssign = getString(dekSec, "type_assign", "type")
	p.Struct = getString(dekSec, "struct")
	p.Ref = getString(dekSec, "ref", "rujuk")
	p.As = getString(dekSec, "as")

	cabSec := getSection("percabangan", "control_flow", "branching")
	p.If = getString(cabSec, "if")
	p.ElseIf = getString(cabSec, "else_if", "elif")
	p.Else = getString(cabSec, "else")
	p.Then = getString(cabSec, "then", "maka")
	p.Switch = getString(cabSec, "switch", "match")
	p.Case = getString(cabSec, "case")
	p.Default = getString(cabSec, "default", "other")

	loopSec := getSection("perulangan", "loops", "loop")
	p.While = getString(loopSec, "while")
	p.For = getString(loopSec, "for")
	p.In = getString(loopSec, "in")

	fnSec := getSection("fungsi", "functions", "function")
	p.Func = getString(fnSec, "def", "func", "fn")
	p.Return = getString(fnSec, "return")
	p.ReturnArrow = getString(fnSec, "return_type", "return_arrow", "arrow")
	p.Defer = getString(fnSec, "defer", "saat_selesai")

	errSec := getSection("penanganan_galat", "error_handling", "errors")
	p.Try = getString(errSec, "try", "coba")
	p.Catch = getString(errSec, "catch", "tangkap")
	p.Finally = getString(errSec, "finally", "akhirnya")
	p.Throw = getString(errSec, "throw", "lempar")

	opSec := getSection("operator_logika", "operators", "operator")
	p.And = getString(opSec, "and")
	p.Or = getString(opSec, "or")
	p.Not = getString(opSec, "not")
	p.NullCoalesce = getString(opSec, "null_coalesce")
	p.Maybe = getString(opSec, "maybe", "mungkin")

	litSec := getSection("literal", "literals")
	p.True = getString(litSec, "true")
	p.False = getString(litSec, "false")
	p.Null = getString(litSec, "null")

	concSec := getSection("konkurensi", "concurrency")
	p.Async = getString(concSec, "async")
	p.Parallel = getString(concSec, "parallel")

	fmtSec := getSection("format", "formatting")
	p.BlockStyle = getString(fmtSec, "block_style")
	p.IndentWidth = getInt(fmtSec, "indent_width")
	p.QuoteStyle = getString(fmtSec, "quote_style")

	return nil
}

// MarshalJSON serializes StyleProfile into the clean categorized format.
func (p StyleProfile) MarshalJSON() ([]byte, error) {
	categorized := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":         p.Name,
			"display_name": p.DisplayName,
			"version":      p.Version,
			"author":       p.Author,
			"description":  p.Description,
		},
		"deklarasi": map[string]interface{}{
			"var":         p.VarDecl,
			"const":       p.ConstDecl,
			"type_assign": p.TypeAssign,
			"struct":      p.Struct,
			"ref":         p.Ref,
			"as":          p.As,
		},
		"percabangan": map[string]interface{}{
			"if":      p.If,
			"else_if": p.ElseIf,
			"else":    p.Else,
			"then":    p.Then,
			"switch":  p.Switch,
			"case":    p.Case,
			"default": p.Default,
		},
		"perulangan": map[string]interface{}{
			"while": p.While,
			"for":   p.For,
			"in":    p.In,
		},
		"fungsi": map[string]interface{}{
			"def":         p.Func,
			"return":      p.Return,
			"return_type": p.ReturnArrow,
			"defer":       p.Defer,
		},
		"penanganan_galat": map[string]interface{}{
			"try":     p.Try,
			"catch":   p.Catch,
			"finally": p.Finally,
			"throw":   p.Throw,
		},
		"operator_logika": map[string]interface{}{
			"and":           p.And,
			"or":            p.Or,
			"not":           p.Not,
			"null_coalesce": p.NullCoalesce,
			"maybe":         p.Maybe,
		},
		"literal": map[string]interface{}{
			"true":  p.True,
			"false": p.False,
			"null":  p.Null,
		},
		"konkurensi": map[string]interface{}{
			"async":    p.Async,
			"parallel": p.Parallel,
		},
		"format": map[string]interface{}{
			"block_style":  p.BlockStyle,
			"indent_width": p.IndentWidth,
			"quote_style":  p.QuoteStyle,
		},
	}
	return json.MarshalIndent(categorized, "", "  ")
}

// IndonesianProfile returns the default human-centric Indonesian style profile.
func IndonesianProfile() *StyleProfile {
	return &StyleProfile{
		Name:         "indonesian",
		DisplayName:  "Bahasa Indonesia (Human-Centric Default)",
		VarDecl:      "buat",
		ConstDecl:    "konstanta",
		TypeAssign:   "itu",
		Struct:       "bentuk",
		Ref:          "rujuk",
		As:           "sebagai",
		If:           "jika",
		ElseIf:       "selain jika",
		Else:         "selain",
		Then:         "maka",
		Switch:       "pilih",
		Case:         "bila",
		Default:      "lainnya",
		While:        "selama",
		For:          "untuk",
		In:           "dalam",
		Func:         "fungsi",
		Return:       "kembalikan",
		ReturnArrow:  "mengembalikan",
		Defer:        "saat_selesai",
		Try:          "coba",
		Catch:        "tangkap",
		Finally:      "akhirnya",
		Throw:        "lempar",
		And:          "dan",
		Or:           "atau",
		Not:          "bukan",
		NullCoalesce: "kalau kosong",
		Maybe:        "mungkin",
		True:         "benar",
		False:        "salah",
		Null:         "kosong",
		Async:        "sambil",
		Parallel:     "bersamaan",
		BlockStyle:   "brace",
		IndentWidth:  4,
		QuoteStyle:   "double",
	}
}

// PythonicProfile returns Python-inspired syntax styling.
func PythonicProfile() *StyleProfile {
	return &StyleProfile{
		Name:         "pythonic",
		DisplayName:  "Pythonic Style",
		VarDecl:      "",
		ConstDecl:    "const",
		TypeAssign:   ":",
		Struct:       "class",
		Ref:          "ref",
		As:           "as",
		If:           "if",
		ElseIf:       "elif",
		Else:         "else",
		Then:         "",
		Switch:       "match",
		Case:         "case",
		Default:      "case _",
		While:        "while",
		For:          "for",
		In:           "in",
		Func:         "def",
		Return:       "return",
		ReturnArrow:  "->",
		Defer:        "defer",
		Try:          "try",
		Catch:        "except",
		Finally:      "finally",
		Throw:        "raise",
		And:          "and",
		Or:           "or",
		Not:          "not",
		NullCoalesce: "??",
		Maybe:        "?",
		True:         "True",
		False:        "False",
		Null:         "None",
		Async:        "async",
		Parallel:     "parallel",
		BlockStyle:   "brace",
		IndentWidth:  4,
		QuoteStyle:   "double",
	}
}

// JavaScriptProfile returns pure modern JavaScript syntax styling.
func JavaScriptProfile() *StyleProfile {
	return &StyleProfile{
		Name:         "javascript",
		DisplayName:  "JavaScript Style (100% JS)",
		VarDecl:      "let",
		ConstDecl:    "const",
		TypeAssign:   ":",
		Struct:       "class",
		Ref:          "ref",
		As:           "as",
		If:           "if",
		ElseIf:       "else if",
		Else:         "else",
		Then:         "",
		Switch:       "switch",
		Case:         "case",
		Default:      "default",
		While:        "while",
		For:          "for",
		In:           "of",
		Func:         "function",
		Return:      "return",
		ReturnArrow:  "",
		Defer:        "defer",
		Try:          "try",
		Catch:        "catch",
		Finally:      "finally",
		Throw:        "throw",
		And:          "&&",
		Or:           "||",
		Not:          "!",
		NullCoalesce: "??",
		Maybe:        "?",
		True:         "true",
		False:        "false",
		Null:         "null",
		Async:        "async",
		Parallel:     "parallel",
		BlockStyle:   "brace",
		IndentWidth:  4,
		QuoteStyle:   "double",
	}
}

// ClassicProfile is an alias to JavaScriptProfile for convenience.
func ClassicProfile() *StyleProfile {
	return JavaScriptProfile()
}

// LoadProfileFromFile reads a custom .style (or .json) file from disk.
func LoadProfileFromFile(filePath string) (*StyleProfile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca berkas style '%s': %w", filePath, err)
	}

	var profile StyleProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("format berkas style '%s' tidak valid: %w", filePath, err)
	}

	if profile.Name == "" {
		profile.Name = filePath
	}
	if profile.DisplayName == "" {
		profile.DisplayName = profile.Name
	}
	if profile.IndentWidth == 0 {
		profile.IndentWidth = 4
	}

	return &profile, nil
}

// GetProfile resolves a profile by predefined name or by file path.
func GetProfile(nameOrPath string) *StyleProfile {
	if strings.HasSuffix(nameOrPath, ".style") || strings.HasSuffix(nameOrPath, ".json") {
		if prof, err := LoadProfileFromFile(nameOrPath); err == nil {
			return prof
		}
	}

	switch strings.ToLower(nameOrPath) {
	case "python", "pythonic", "py":
		return PythonicProfile()
	case "javascript", "js", "classic", "ts":
		return JavaScriptProfile()
	case "id", "indo", "indonesia", "indonesian":
		return IndonesianProfile()
	default:
		if _, err := os.Stat(nameOrPath); err == nil {
			if prof, err := LoadProfileFromFile(nameOrPath); err == nil {
				return prof
			}
		}
		return IndonesianProfile()
	}
}

func init() {
	// Register global style resolver with the parser package so any Lexer can dynamically resolve styles.
	parser.GlobalStyleResolver = func(nameOrPath string) map[string]parser.TokenType {
		prof := GetProfile(nameOrPath)
		if prof == nil {
			return nil
		}
		return prof.ExportKeywords()
	}
}

// ExportKeywords maps all declared keywords in the StyleProfile to their corresponding parser.TokenType.
func (p *StyleProfile) ExportKeywords() map[string]parser.TokenType {
	m := make(map[string]parser.TokenType)
	if p == nil {
		return m
	}
	reg := func(field string, tok parser.TokenType) {
		if field == "" {
			return
		}
		for _, item := range strings.FieldsFunc(field, func(r rune) bool {
			return r == ',' || r == '|' || r == ';'
		}) {
			word := strings.TrimSpace(item)
			if word != "" {
				m[word] = tok
			}
		}
	}

	reg(p.VarDecl, parser.VAR_DECL)
	reg(p.ConstDecl, parser.CONST_DECL)
	reg(p.TypeAssign, parser.TYPE_ASSIGN)
	reg(p.Struct, parser.STRUCT_DECL)
	reg(p.Ref, parser.REF_DECL)
	reg(p.As, parser.AS)

	reg(p.If, parser.IF)
	reg(p.ElseIf, parser.ELSE_IF)
	reg(p.Else, parser.ELSE)
	reg(p.Then, parser.MAKA)
	reg(p.Switch, parser.SWITCH)
	reg(p.Case, parser.CASE)
	reg(p.Default, parser.DEFAULT)

	reg(p.While, parser.WHILE)
	reg(p.For, parser.FOR)
	reg(p.In, parser.IN)

	reg(p.Func, parser.FUNC)
	reg(p.Return, parser.RETURN)
	reg(p.ReturnArrow, parser.RETURN_TYPE)
	reg(p.Defer, parser.DEFER)

	reg(p.Try, parser.TRY)
	reg(p.Catch, parser.CATCH)
	reg(p.Finally, parser.FINALLY)
	reg(p.Throw, parser.THROW)

	reg(p.And, parser.AND)
	reg(p.Or, parser.OR)
	reg(p.Not, parser.NOT)
	reg(p.NullCoalesce, parser.NULL_COALESCE)
	reg(p.Maybe, parser.MAYBE)

	reg(p.True, parser.TRUE)
	reg(p.False, parser.FALSE)
	reg(p.Null, parser.NULL)

	reg(p.Async, parser.ASYNC)
	reg(p.Parallel, parser.PARALLEL)

	return m
}

// ApplyToLexer registers this style's keywords directly onto a Lexer instance.
func (p *StyleProfile) ApplyToLexer(l *parser.Lexer) {
	if p == nil || l == nil {
		return
	}
	for word, tok := range p.ExportKeywords() {
		l.RegisterKeyword(word, tok)
	}
}
