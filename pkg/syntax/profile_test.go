package syntax_test

import (
	"encoding/json"
	"strings"
	"testing"

	"celang/pkg/parser"
	"celang/pkg/syntax"
)

func TestProfile_UnmarshalCategorizedJSON(t *testing.T) {
	categorizedJSON := `{
		"metadata": {
			"name": "sunda",
			"display_name": "Basa Sunda",
			"version": "1.0.0"
		},
		"deklarasi": {
			"var": "jieun",
			"const": "tetep",
			"type_assign": "nyaeta",
			"struct": "wangun"
		},
		"percabangan": {
			"if": "lamun",
			"else_if": "atawa lamun",
			"else": "lamun henteu",
			"switch": "cocogkeun",
			"case": "bisi",
			"default": "lianna"
		},
		"perulangan": {
			"while": "salami",
			"for": "pikeun",
			"in": "dina"
		},
		"fungsi": {
			"def": "fungsi",
			"return": "balikkeun",
			"return_type": "ngahasilkeun",
			"defer": "pas_rengse"
		},
		"penanganan_galat": {
			"try": "cobi",
			"catch": "tewak",
			"finally": "pamungkas",
			"throw": "alungkeun"
		},
		"literal": {
			"true": "bener",
			"false": "salah",
			"null": "kosong"
		},
		"format": {
			"block_style": "brace",
			"indent_width": 4
		}
	}`

	var prof syntax.StyleProfile
	if err := json.Unmarshal([]byte(categorizedJSON), &prof); err != nil {
		t.Fatalf("Failed to unmarshal categorized profile: %v", err)
	}

	if prof.Name != "sunda" {
		t.Errorf("Expected Name='sunda', got %q", prof.Name)
	}
	if prof.VarDecl != "jieun" {
		t.Errorf("Expected VarDecl='jieun', got %q", prof.VarDecl)
	}
	if prof.If != "lamun" {
		t.Errorf("Expected If='lamun', got %q", prof.If)
	}
	if prof.While != "salami" {
		t.Errorf("Expected While='salami', got %q", prof.While)
	}
	if prof.Return != "balikkeun" {
		t.Errorf("Expected Return='balikkeun', got %q", prof.Return)
	}
	if prof.Try != "cobi" {
		t.Errorf("Expected Try='cobi', got %q", prof.Try)
	}
	if prof.Catch != "tewak" {
		t.Errorf("Expected Catch='tewak', got %q", prof.Catch)
	}
	if prof.Defer != "pas_rengse" {
		t.Errorf("Expected Defer='pas_rengse', got %q", prof.Defer)
	}
	if prof.Struct != "wangun" {
		t.Errorf("Expected Struct='wangun', got %q", prof.Struct)
	}
}

func TestProfile_UnmarshalFlatJSON(t *testing.T) {
	flatJSON := `{
		"name": "flat_test",
		"display_name": "Flat Test Profile",
		"var_decl": "let",
		"if": "if",
		"while": "while",
		"func": "fn"
	}`

	var prof syntax.StyleProfile
	if err := json.Unmarshal([]byte(flatJSON), &prof); err != nil {
		t.Fatalf("Failed to unmarshal flat profile: %v", err)
	}

	if prof.Name != "flat_test" {
		t.Errorf("Expected Name='flat_test', got %q", prof.Name)
	}
	if prof.VarDecl != "let" {
		t.Errorf("Expected VarDecl='let', got %q", prof.VarDecl)
	}
	if prof.If != "if" {
		t.Errorf("Expected If='if', got %q", prof.If)
	}
	if prof.Func != "fn" {
		t.Errorf("Expected Func='fn', got %q", prof.Func)
	}
}

func TestProfile_MarshalCategorizedJSON(t *testing.T) {
	indo := syntax.IndonesianProfile()
	data, err := json.Marshal(indo)
	if err != nil {
		t.Fatalf("Failed to marshal Indonesian profile: %v", err)
	}

	str := string(data)
	if !strings.Contains(str, `"metadata"`) ||
		!strings.Contains(str, `"deklarasi"`) ||
		!strings.Contains(str, `"percabangan"`) ||
		!strings.Contains(str, `"fungsi"`) {
		t.Fatalf("Marshaled JSON missing expected categories:\n%s", str)
	}
}

func TestProfile_ApplyToLexer(t *testing.T) {
	prof := &syntax.StyleProfile{
		VarDecl: "gawe, bikin",
		If:      "lamun_bener, lamun-bener",
	}
	l := parser.NewLexer("gawe x = 10\nlamun_bener x == 10 {\n}")
	prof.ApplyToLexer(l)

	tok1 := l.NextToken()
	if tok1.Type != parser.VAR_DECL || tok1.Literal != "gawe" {
		t.Fatalf("Expected VAR_DECL 'gawe', got %v (%s)", tok1.Type, tok1.Literal)
	}

	l.NextToken() // 'x'
	l.NextToken() // '='
	l.NextToken() // '10'

	tok2 := l.NextToken()
	if tok2.Type != parser.IF || tok2.Literal != "lamun_bener" {
		t.Fatalf("Expected IF 'lamun_bener', got %v (%s)", tok2.Type, tok2.Literal)
	}
}
