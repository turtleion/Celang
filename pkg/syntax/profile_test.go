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
			"name": "jawa",
			"display_name": "Basa Jawa",
			"version": "1.0.0"
		},
		"deklarasi": {
			"var": "gawe",
			"const": "tetep",
			"type_assign": "yaiku",
			"struct": "wujud"
		},
		"percabangan": {
			"if": "yen",
			"else_if": "utawa yen",
			"else": "yen ora",
			"switch": "cocogna",
			"case": "nalika",
			"default": "liyane"
		},
		"perulangan": {
			"while": "sasuwene",
			"for": "kanggo",
			"in": "ing"
		},
		"fungsi": {
			"def": "fungsi",
			"return": "balekna",
			"return_type": "ngasilakeun",
			"defer": "pas_bubar"
		},
		"penanganan_galat": {
			"try": "jajal",
			"catch": "cekel",
			"finally": "pungkasan",
			"throw": "uncalna"
		},
		"literal": {
			"true": "bener",
			"false": "salah",
			"null": "suwung"
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

	if prof.Name != "jawa" {
		t.Errorf("Expected Name='jawa', got %q", prof.Name)
	}
	if prof.VarDecl != "gawe" {
		t.Errorf("Expected VarDecl='gawe', got %q", prof.VarDecl)
	}
	if prof.If != "yen" {
		t.Errorf("Expected If='yen', got %q", prof.If)
	}
	if prof.While != "sasuwene" {
		t.Errorf("Expected While='sasuwene', got %q", prof.While)
	}
	if prof.Return != "balekna" {
		t.Errorf("Expected Return='balekna', got %q", prof.Return)
	}
	if prof.Try != "jajal" {
		t.Errorf("Expected Try='jajal', got %q", prof.Try)
	}
	if prof.Catch != "cekel" {
		t.Errorf("Expected Catch='cekel', got %q", prof.Catch)
	}
	if prof.Defer != "pas_bubar" {
		t.Errorf("Expected Defer='pas_bubar', got %q", prof.Defer)
	}
	if prof.Struct != "wujud" {
		t.Errorf("Expected Struct='wujud', got %q", prof.Struct)
	}
}

func TestProfile_BuiltinJawaProfile(t *testing.T) {
	prof := syntax.GetProfile("jawa")
	if prof == nil {
		t.Fatalf("Expected non-nil profile for 'jawa'")
	}
	if prof.Name != "jawa" || prof.DisplayName != "Basa Jawa" {
		t.Errorf("Expected Name='jawa', DisplayName='Basa Jawa', got name=%q, display=%q", prof.Name, prof.DisplayName)
	}
	if prof.VarDecl != "gawe" || prof.If != "yen" {
		t.Errorf("Expected VarDecl='gawe', If='yen', got var=%q, if=%q", prof.VarDecl, prof.If)
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
