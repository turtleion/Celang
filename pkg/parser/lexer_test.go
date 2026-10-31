package parser_test

import (
	"testing"

	"celang/pkg/parser"
)

func TestLexer_SubtractionWithoutSpaces(t *testing.T) {
	input := "x-y"
	l := parser.NewLexer(input)

	tok1 := l.NextToken()
	if tok1.Type != parser.IDENT || tok1.Literal != "x" {
		t.Fatalf("Expected IDENT 'x', got %v (%s)", tok1.Type, tok1.Literal)
	}

	tok2 := l.NextToken()
	if tok2.Type != parser.MINUS || tok2.Literal != "-" {
		t.Fatalf("Expected MINUS '-', got %v (%s)", tok2.Type, tok2.Literal)
	}

	tok3 := l.NextToken()
	if tok3.Type != parser.IDENT || tok3.Literal != "y" {
		t.Fatalf("Expected IDENT 'y', got %v (%s)", tok3.Type, tok3.Literal)
	}
}

func TestLexer_ThreeWordCompoundKeyword(t *testing.T) {
	input := "x tidak sama dengan y"
	l := parser.NewLexer(input)
	l.RegisterKeyword("tidak sama dengan", parser.NOT_EQ)

	tok1 := l.NextToken()
	if tok1.Type != parser.IDENT || tok1.Literal != "x" {
		t.Fatalf("Expected IDENT 'x', got %v (%s)", tok1.Type, tok1.Literal)
	}

	tok2 := l.NextToken()
	if tok2.Type != parser.NOT_EQ || tok2.Literal != "tidak sama dengan" {
		t.Fatalf("Expected NOT_EQ 'tidak sama dengan', got %v (%s)", tok2.Type, tok2.Literal)
	}

	tok3 := l.NextToken()
	if tok3.Type != parser.IDENT || tok3.Literal != "y" {
		t.Fatalf("Expected IDENT 'y', got %v (%s)", tok3.Type, tok3.Literal)
	}
}

func TestLexer_HyphenatedKeywords(t *testing.T) {
	input := "laba-laba-sunda x == 1"
	l := parser.NewLexer(input)
	l.RegisterKeyword("laba-laba-sunda", parser.ELSE)

	tok1 := l.NextToken()
	if tok1.Type != parser.ELSE || tok1.Literal != "laba-laba-sunda" {
		t.Fatalf("Expected ELSE 'laba-laba-sunda', got %v (%s)", tok1.Type, tok1.Literal)
	}
}

func TestLexer_EmojiIdentifiersAndKeywords(t *testing.T) {
	input := "buat🚀 x_✨ = 10"
	l := parser.NewLexer(input)
	l.RegisterKeyword("buat🚀", parser.VAR_DECL)

	tok1 := l.NextToken()
	if tok1.Type != parser.VAR_DECL || tok1.Literal != "buat🚀" {
		t.Fatalf("Expected VAR_DECL 'buat🚀', got %v (%s)", tok1.Type, tok1.Literal)
	}

	tok2 := l.NextToken()
	if tok2.Type != parser.IDENT || tok2.Literal != "x_✨" {
		t.Fatalf("Expected IDENT 'x_✨', got %v (%s)", tok2.Type, tok2.Literal)
	}
}

func TestLexer_BuiltinHyphenatedKeywords(t *testing.T) {
	input := "privat-modul! fungsi foo() {}\nprivat-[auth] fungsi bar() {}\nx-y"
	l := parser.NewLexer(input)

	tok1 := l.NextToken()
	if tok1.Type != parser.VIS_MODULE || tok1.Literal != "privat-modul!" {
		t.Fatalf("Expected VIS_MODULE 'privat-modul!', got %v (%s)", tok1.Type, tok1.Literal)
	}

	for l.NextToken().Type != parser.RBRACE {
	}

	tok2 := l.NextToken()
	if tok2.Type != parser.VIS_ALLOWLIST || tok2.Literal != "privat-" {
		t.Fatalf("Expected VIS_ALLOWLIST 'privat-', got %v (%s)", tok2.Type, tok2.Literal)
	}

	for l.NextToken().Type != parser.RBRACE {
	}

	tokX := l.NextToken()
	if tokX.Type != parser.IDENT || tokX.Literal != "x" {
		t.Fatalf("Expected IDENT 'x', got %v (%s)", tokX.Type, tokX.Literal)
	}
	tokMinus := l.NextToken()
	if tokMinus.Type != parser.MINUS || tokMinus.Literal != "-" {
		t.Fatalf("Expected MINUS '-', got %v (%s)", tokMinus.Type, tokMinus.Literal)
	}
	tokY := l.NextToken()
	if tokY.Type != parser.IDENT || tokY.Literal != "y" {
		t.Fatalf("Expected IDENT 'y', got %v (%s)", tokY.Type, tokY.Literal)
	}
}



