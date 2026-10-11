package parser

import (
	"strings"
	"unicode"
)

var keywords = map[string]TokenType{
	"buat":         VAR_DECL,
	"let":          VAR_DECL,
	"var":          VAR_DECL,
	"konstanta":    CONST_DECL,
	"const":        CONST_DECL,
	"itu":          TYPE_ASSIGN,
	"jika":         IF,
	"if":           IF,
	"elif":         ELSE_IF,
	"selain":       ELSE,
	"else":         ELSE,
	"cocok":        SWITCH,
	"switch":       SWITCH,
	"match":        SWITCH,
	"kasus":        CASE,
	"case":         CASE,
	"lainnya":      DEFAULT,
	"default":      DEFAULT,
	"selama":       WHILE,
	"while":        WHILE,
	"untuk":        FOR,
	"for":          FOR,
	"dalam":        IN,
	"in":           IN,
	"fungsi":       FUNC,
	"def":          FUNC,
	"fn":           FUNC,
	"function":     FUNC,
	"kembalikan":   RETURN,
	"return":       RETURN,
	"sambil":       ASYNC,
	"async":        ASYNC,
	"bersamaan":    PARALLEL,
	"parallel":     PARALLEL,
	"benar":        TRUE,
	"true":         TRUE,
	"True":         TRUE,
	"salah":        FALSE,
	"false":        FALSE,
	"False":        FALSE,
	"kosong":       NULL,
	"nihil":        NULL,
	"null":         NULL,
	"None":         NULL,
	"dan":          AND,
	"and":          AND,
	"atau":         OR,
	"or":           OR,
	"bukan":        NOT,
	"tidak":        NOT,
	"not":          NOT,
	"mengembalikan": RETURN_TYPE,
	"maka":         MAKA,
	"selesai":      SELESAI,
	"ambil":        IMPORT,
	"impor":        IMPORT,
	"import":       IMPORT,
	"dari":         FROM,
	"from":         FROM,
	"sebagai":      AS,
	"as":           AS,
	"export":       VIS_PUBLIC,
	"coba":         TRY,
	"try":          TRY,
	"tangkap":      CATCH,
	"catch":        CATCH,
	"except":       CATCH,
	"akhirnya":     FINALLY,
	"finally":      FINALLY,
	"lempar":       THROW,
	"throw":        THROW,
	"raise":        THROW,
	"gaya":         STYLE_BLOCK,
	"style":        STYLE_BLOCK,
	"dialek":       STYLE_BLOCK,
	"bentuk":       STRUCT_DECL,
	"struct":       STRUCT_DECL,
	"privat":       PRIVATE,
	"rujuk":        REF_DECL,
	"ref":          REF_DECL,
	"tipe":         TYPE_DECL,
}

type Lexer struct {
	input        string
	position     int  // current character position
	readPosition int  // next character position
	ch           byte // current character
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	var tok Token
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
			}
			tok = Token{Type: EQ, Literal: "==", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: ASSIGN, Literal: "=", Line: tok.Line, Column: tok.Column}
		}
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: PLUS_ASSIGN, Literal: "+=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: PLUS, Literal: "+", Line: tok.Line, Column: tok.Column}
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: MINUS_ASSIGN, Literal: "-=", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: ARROW, Literal: "->", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: MINUS, Literal: "-", Line: tok.Line, Column: tok.Column}
		}
	case '*':
		tok = Token{Type: ASTERISK, Literal: "*", Line: tok.Line, Column: tok.Column}
	case '/':
		tok = Token{Type: SLASH, Literal: "/", Line: tok.Line, Column: tok.Column}
	case '%':
		tok = Token{Type: PERCENT, Literal: "%", Line: tok.Line, Column: tok.Column}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
			}
			tok = Token{Type: NOT_EQ, Literal: "!=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: BANG, Literal: "!", Line: tok.Line, Column: tok.Column}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: LT_EQ, Literal: "<=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: LT, Literal: "<", Line: tok.Line, Column: tok.Column}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: GT_EQ, Literal: ">=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: GT, Literal: ">", Line: tok.Line, Column: tok.Column}
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: AND, Literal: "&&", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: AMPERSAND, Literal: "&", Line: tok.Line, Column: tok.Column}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: OR, Literal: "||", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: PIPE_ARROW, Literal: "|>", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: PIPE, Literal: "|", Line: tok.Line, Column: tok.Column}
		}
	case ',':
		tok = Token{Type: COMMA, Literal: ",", Line: tok.Line, Column: tok.Column}
	case '@':
		tok = Token{Type: AT, Literal: "@", Line: tok.Line, Column: tok.Column}
	case ':':
		tok = Token{Type: COLON, Literal: ":", Line: tok.Line, Column: tok.Column}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: ";", Line: tok.Line, Column: tok.Column}
	case '(':
		tok = Token{Type: LPAREN, Literal: "(", Line: tok.Line, Column: tok.Column}
	case ')':
		tok = Token{Type: RPAREN, Literal: ")", Line: tok.Line, Column: tok.Column}
	case '{':
		tok = Token{Type: LBRACE, Literal: "{", Line: tok.Line, Column: tok.Column}
	case '}':
		tok = Token{Type: RBRACE, Literal: "}", Line: tok.Line, Column: tok.Column}
	case '[':
		tok = Token{Type: LBRACKET, Literal: "[", Line: tok.Line, Column: tok.Column}
	case ']':
		tok = Token{Type: RBRACKET, Literal: "]", Line: tok.Line, Column: tok.Column}
	case '.':
		tok = Token{Type: DOT, Literal: ".", Line: tok.Line, Column: tok.Column}
	case '"', '\'', '`':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		return tok
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)

			// Handle compound keywords: "selain jika" and "else if"
			if (tok.Literal == "selain" || tok.Literal == "else") && l.peekNextWord() == "jika" || (tok.Literal == "else" && l.peekNextWord() == "if") {
				l.skipWhitespaceAndComments()
				secondWord := l.readIdentifier()
				tok.Type = ELSE_IF
				tok.Literal = tok.Literal + " " + secondWord
			}

			// Handle visibility modifiers: publik!, privat-modul!, privat!, and privat-[...]
			if tok.Literal == "publik" && l.ch == '!' {
				l.readChar()
				tok.Literal = "publik!"
				tok.Type = VIS_PUBLIC
				return tok
			}
			if tok.Literal == "privat-modul" && l.ch == '!' {
				l.readChar()
				tok.Literal = "privat-modul!"
				tok.Type = VIS_MODULE
				return tok
			}
			if tok.Literal == "privat" {
				if l.ch == '!' {
					l.readChar()
					tok.Literal = "privat!"
					tok.Type = VIS_PRIVATE
					return tok
				}
				if l.ch == '-' && l.peekChar() == '[' {
					l.readChar() // eat '-'
					tok.Literal = "privat-"
					tok.Type = VIS_ALLOWLIST
					return tok
				}
			}

			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		// Whitespace
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}

		// Single line comments: // or #
		if l.ch == '#' || (l.ch == '/' && l.peekChar() == '/') {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Multi-line comment: /* ... */
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar() // eat '/'
			l.readChar() // eat '*'
			for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
				if l.ch == '\n' {
					l.line++
					l.column = 0
				}
				l.readChar()
			}
			if l.ch != 0 {
				l.readChar() // eat '*'
				l.readChar() // eat '/'
			}
			continue
		}

		break
	}
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || (l.ch == '-' && isLetter(l.peekChar())) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) readString(delimiter byte) string {
	l.readChar() // Skip opening quote
	var sb strings.Builder
	for l.ch != delimiter && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '"':
				sb.WriteByte('"')
			case '\'':
				sb.WriteByte('\'')
			case '`':
				sb.WriteByte('`')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			sb.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == delimiter {
		l.readChar() // Skip closing quote
	}
	return sb.String()
}

func (l *Lexer) readNumber() (TokenType, string) {
	startPos := l.position
	isFloat := false

	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // eat '.'
		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	}

	lit := l.input[startPos:l.position]
	// Remove underscores from number literal string
	cleanLit := strings.ReplaceAll(lit, "_", "")
	if isFloat {
		return FLOAT, cleanLit
	}
	return INT, cleanLit
}

func (l *Lexer) peekNextWord() string {
	// Lookahead without permanently mutating lexer position
	savedPos := l.position
	savedReadPos := l.readPosition
	savedCh := l.ch
	savedLine := l.line
	savedCol := l.column

	l.skipWhitespaceAndComments()
	word := ""
	if isLetter(l.ch) {
		word = l.readIdentifier()
	}

	// Restore lexer state
	l.position = savedPos
	l.readPosition = savedReadPos
	l.ch = savedCh
	l.line = savedLine
	l.column = savedCol

	return word
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
