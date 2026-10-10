package syntax

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// StyleProfile defines how keywords and block structures are mapped to U-AST.
// It acts as the "CSS for Syntax" in Celang.
type StyleProfile struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`

	// Keywords
	VarDecl      string   `json:"var_decl"`       // e.g. "buat", "let", ""
	ConstDecl    string   `json:"const_decl"`     // e.g. "konstanta", "const"
	TypeAssign   string   `json:"type_assign"`    // e.g. "itu", ":"
	If           string   `json:"if"`             // e.g. "jika", "if"
	ElseIf       string   `json:"else_if"`        // e.g. "selain jika", "elif", "else if"
	Else         string   `json:"else"`           // e.g. "selain", "else"
	Switch       string   `json:"switch"`         // e.g. "cocok", "switch"
	Case         string   `json:"case"`           // e.g. "kasus", "case"
	Default      string   `json:"default"`        // e.g. "lainnya", "default"
	While        string   `json:"while"`          // e.g. "selama", "while"
	For          string   `json:"for"`            // e.g. "untuk", "for"
	In           string   `json:"in"`             // e.g. "dalam", "of", "in"
	Func         string   `json:"func"`           // e.g. "fungsi", "def", "function"
	Return       string   `json:"return"`         // e.g. "kembalikan", "return"
	Async        string   `json:"async"`          // e.g. "sambil", "async"
	Parallel     string   `json:"parallel"`       // e.g. "bersamaan", "parallel"
	True         string   `json:"true"`           // e.g. "benar", "True", "true"
	False        string   `json:"false"`          // e.g. "salah", "False", "false"
	Null         string   `json:"null"`           // e.g. "kosong", "None", "null"
	And          string   `json:"and"`            // e.g. "dan", "and", "&&"
	Or           string   `json:"or"`             // e.g. "atau", "or", "||"
	Not          string   `json:"not"`            // e.g. "bukan", "not", "!"
	ReturnArrow  string   `json:"return_arrow"`   // e.g. "mengembalikan", "->"

	// Structural styling
	BlockStyle   string `json:"block_style"`   // "brace" ({ }), "keyword" (maka ... selesai)
	IndentWidth  int    `json:"indent_width"`  // default indentation spaces (e.g. 4)
}

// IndonesianProfile returns the default human-centric Indonesian style profile.
func IndonesianProfile() *StyleProfile {
	return &StyleProfile{
		Name:        "indonesian",
		DisplayName: "Bahasa Indonesia (Human-Centric Default)",
		VarDecl:     "buat",
		ConstDecl:   "konstanta",
		TypeAssign:  "itu",
		If:          "jika",
		ElseIf:      "selain jika",
		Else:        "selain",
		Switch:      "cocok",
		Case:        "kasus",
		Default:     "lainnya",
		While:       "selama",
		For:         "untuk",
		In:          "dalam",
		Func:        "fungsi",
		Return:      "kembalikan",
		Async:       "sambil",
		Parallel:    "bersamaan",
		True:        "benar",
		False:       "salah",
		Null:        "kosong",
		And:         "dan",
		Or:          "atau",
		Not:         "bukan",
		ReturnArrow: "mengembalikan",
		BlockStyle:  "brace",
		IndentWidth: 4,
	}
}

// PythonicProfile returns Python-inspired syntax styling.
func PythonicProfile() *StyleProfile {
	return &StyleProfile{
		Name:        "pythonic",
		DisplayName: "Pythonic Style",
		VarDecl:     "",
		ConstDecl:   "const",
		TypeAssign:  ":",
		If:          "if",
		ElseIf:      "elif",
		Else:        "else",
		Switch:      "match",
		Case:        "case",
		Default:     "case _",
		While:       "while",
		For:         "for",
		In:          "in",
		Func:        "def",
		Return:      "return",
		Async:       "async",
		Parallel:    "parallel",
		True:        "True",
		False:       "False",
		Null:        "None",
		And:         "and",
		Or:          "or",
		Not:         "not",
		ReturnArrow: "->",
		BlockStyle:  "brace",
		IndentWidth: 4,
	}
}

// JavaScriptProfile returns pure modern JavaScript syntax styling.
func JavaScriptProfile() *StyleProfile {
	return &StyleProfile{
		Name:        "javascript",
		DisplayName: "JavaScript Style (100% JS)",
		VarDecl:     "let",
		ConstDecl:   "const",
		TypeAssign:  ":",
		If:          "if",
		ElseIf:      "else if",
		Else:        "else",
		Switch:      "switch",
		Case:        "case",
		Default:     "default",
		While:       "while",
		For:         "for",
		In:          "of",
		Func:        "function",
		Return:      "return",
		Async:       "async",
		Parallel:    "parallel",
		True:        "true",
		False:       "false",
		Null:        "null",
		And:         "&&",
		Or:          "||",
		Not:         "!",
		ReturnArrow: "",
		BlockStyle:  "brace",
		IndentWidth: 4,
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
	// 1. Check if nameOrPath is an existing file on disk
	if strings.HasSuffix(nameOrPath, ".style") || strings.HasSuffix(nameOrPath, ".json") {
		if prof, err := LoadProfileFromFile(nameOrPath); err == nil {
			return prof
		}
	}

	// 2. Predefined keyword profiles
	switch strings.ToLower(nameOrPath) {
	case "python", "pythonic", "py":
		return PythonicProfile()
	case "javascript", "js", "classic", "ts":
		return JavaScriptProfile()
	case "id", "indo", "indonesia", "indonesian":
		return IndonesianProfile()
	default:
		// Attempt file loading if it exists
		if _, err := os.Stat(nameOrPath); err == nil {
			if prof, err := LoadProfileFromFile(nameOrPath); err == nil {
				return prof
			}
		}
		return IndonesianProfile()
	}
}
