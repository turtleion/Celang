package syntax

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
	While        string   `json:"while"`          // e.g. "selama", "while"
	For          string   `json:"for"`            // e.g. "untuk", "for"
	In           string   `json:"in"`             // e.g. "dalam", "in"
	Func         string   `json:"func"`           // e.g. "fungsi", "def", "fn"
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
		BlockStyle:  "brace", // or colon/indent
		IndentWidth: 4,
	}
}

// ClassicProfile returns C/JavaScript-inspired syntax styling.
func ClassicProfile() *StyleProfile {
	return &StyleProfile{
		Name:        "classic",
		DisplayName: "C / JavaScript Style",
		VarDecl:     "let",
		ConstDecl:   "const",
		TypeAssign:  ":",
		If:          "if",
		ElseIf:      "else if",
		Else:        "else",
		While:       "while",
		For:         "for",
		In:          "in",
		Func:        "fn",
		Return:      "return",
		Async:       "async",
		Parallel:    "parallel",
		True:        "true",
		False:       "false",
		Null:        "null",
		And:         "&&",
		Or:          "||",
		Not:         "!",
		ReturnArrow: "->",
		BlockStyle:  "brace",
		IndentWidth: 4,
	}
}

// GetProfile returns a profile by name, falling back to Indonesian profile.
func GetProfile(name string) *StyleProfile {
	switch name {
	case "python", "pythonic", "py":
		return PythonicProfile()
	case "classic", "c", "js", "ts":
		return ClassicProfile()
	case "id", "indo", "indonesia", "indonesian":
		return IndonesianProfile()
	default:
		return IndonesianProfile()
	}
}
