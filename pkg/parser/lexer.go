package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"
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
	"lain":         ELSE,
	"else":         ELSE,
	"pilih":        SWITCH,
	"cocok":        SWITCH,
	"switch":       SWITCH,
	"match":        SWITCH,
	"bila":         CASE,
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
	"berhenti":     BREAK,
	"break":        BREAK,
	"lanjut":       CONTINUE,
	"continue":     CONTINUE,
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
	"ekspor":       VIS_PUBLIC,
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
	"bentuk":         STRUCT_DECL,
	"struct":         STRUCT_DECL,
	"privat":         PRIVATE,
	"privat-modul!":  VIS_MODULE,
	"rujuk":          REF_DECL,
	"ref":            REF_DECL,
	"tipe":           TYPE_DECL,
	"mungkin":        MAYBE,
	"saat_selesai":   DEFER,
	"defer":          DEFER,
	"ketika_selesai": DEFER,
	"kalau_kosong":   NULL_COALESCE,
	"bila_kosong":    NULL_COALESCE,

	// Sunda Dialect keywords
	"jieun":        VAR_DECL,
	"tetep":        CONST_DECL,
	"nyaeta":       TYPE_ASSIGN,
	"wangun":       STRUCT_DECL,
	"lamun":        IF,
	"cocogkeun":    SWITCH,
	"bisi":         CASE,
	"lianna":       DEFAULT,
	"salami":       WHILE,
	"pikeun":       FOR,
	"dina":         IN,
	"balikkeun":    RETURN,
	"eureun":       BREAK,
	"teraskeun":    CONTINUE,
	"ngahasilkeun": RETURN_TYPE,
	"pas_rengse":   DEFER,
	"cobi":         TRY,
	"tewak":        CATCH,
	"pamungkas":    FINALLY,
	"alungkeun":    THROW,
	"sareng":       AND,
	"sanes":        NOT,
	"mun_kosong":   NULL_COALESCE,
	"meureun":      MAYBE,
	"bener":        TRUE,
	"bari":         ASYNC,
	"babarengan":   PARALLEL,
	"salaku":       AS,

	// Modern Concurrency & Generics
	"titip":     TITIP,
	"spawn":     TITIP,
	"jemput":    JEMPUT,
	"await":     JEMPUT,
	"harus":     HARUS,
	"must":      HARUS,
	"sembarang": SEMBARANG,
	"any":       SEMBARANG,
	"saluran":   SALURAN,
	"channel":   SALURAN,

	// Exclusive Killer Features: Transactions & Reactivity
	"transaksi":   TRANSACTION,
	"transaction": TRANSACTION,
	"batalkan":    ROLLBACK,
	"rollback":    ROLLBACK,
	"batalkeun":   ROLLBACK,
	"pulihkan":    RECOVER,
	"recover":     RECOVER,
	"pulihkeun":   RECOVER,
	"reaktif":     REACTIVE_DECL,
	"reactive":    REACTIVE_DECL,
	"pantau":      WATCH,
	"watch":       WATCH,
	"awaskeun":    WATCH,

	// Testing Framework
	"uji":  TEST,
	"test": TEST,
}

// GlobalStyleResolver resolves a style name or file path to its token mappings.
// Initialized by pkg/syntax to avoid circular imports.
var GlobalStyleResolver func(nameOrPath string) map[string]TokenType

type Lexer struct {
	input           string
	position        int  // current character position
	readPosition    int  // next character position
	ch              byte // current character
	line            int
	column          int
	IsUniversal     bool
	StyleName       string
	customKeywords  map[string]TokenType
	customCompounds map[string]map[string]TokenType
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:           input,
		line:            1,
		column:          0,
		IsUniversal:     DetectUniversalDirective(input),
		StyleName:       DetectStyleDirective(input),
		customKeywords:  make(map[string]TokenType),
		customCompounds: make(map[string]map[string]TokenType),
	}
	if l.StyleName != "" {
		l.ApplyStyle(l.StyleName)
	}
	l.readChar()
	return l
}

// RegisterKeyword registers a custom keyword or comma/pipe-separated alias list for a given token type.
func (l *Lexer) RegisterKeyword(rawWord string, tok TokenType) {
	if l.customKeywords == nil {
		l.customKeywords = make(map[string]TokenType)
	}
	if l.customCompounds == nil {
		l.customCompounds = make(map[string]map[string]TokenType)
	}

	for _, item := range strings.FieldsFunc(rawWord, func(r rune) bool {
		return r == ',' || r == '|' || r == ';'
	}) {
		word := strings.TrimSpace(item)
		if word == "" {
			continue
		}

		parts := strings.Fields(word)
		if len(parts) >= 2 {
			first := parts[0]
			second := strings.Join(parts[1:], " ")
			if _, ok := l.customCompounds[first]; !ok {
				l.customCompounds[first] = make(map[string]TokenType)
			}
			l.customCompounds[first][second] = tok
			continue
		}

		l.customKeywords[word] = tok

		// If word contains hyphens ("lato-lato"), also register underscore version ("lato_lato")
		if strings.Contains(word, "-") {
			l.customKeywords[strings.ReplaceAll(word, "-", "_")] = tok
		}
		// If word contains underscores ("lato_lato"), also register hyphen version ("lato-lato")
		if strings.Contains(word, "_") {
			l.customKeywords[strings.ReplaceAll(word, "_", "-")] = tok
		}
	}
}

// ApplyStyle resolves and applies a style by name or file path if GlobalStyleResolver is set.
func (l *Lexer) ApplyStyle(nameOrPath string) {
	if GlobalStyleResolver != nil && nameOrPath != "" {
		kws := GlobalStyleResolver(nameOrPath)
		for word, tok := range kws {
			l.RegisterKeyword(word, tok)
		}
	}
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
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: FAT_ARROW, Literal: "=>", Line: tok.Line, Column: tok.Column}
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
		if l.peekChar() == '.' {
			l.readChar()
			tok = Token{Type: RANGE, Literal: "..", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: DOT, Literal: ".", Line: tok.Line, Column: tok.Column}
		}
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			tok = Token{Type: NULL_COALESCE, Literal: "??", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '.' {
			l.readChar()
			tok = Token{Type: OPTIONAL_DOT, Literal: "?.", Line: tok.Line, Column: tok.Column}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '"', '\'', '`':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		return tok
	case 'f', 'F':
		if l.peekChar() == '"' || l.peekChar() == '\'' || l.peekChar() == '`' {
			quote := l.peekChar()
			l.readChar() // Skip 'f' or 'F', l.ch becomes opening quote
			tok.Type = STRING
			tok.Literal = l.readString(quote)
			return tok
		}
		fallthrough
	default:
		isStart := false
		if l.ch < 0x80 {
			isStart = isLetter(l.ch)
		} else {
			r, _ := utf8.DecodeRuneInString(l.input[l.position:])
			isStart = isIdentifierStartRune(r)
		}

		if isStart {
			tok.Literal = l.readIdentifier()

			// Check custom or built-in hyphenated keywords (e.g. "laba-laba-sunda", "privat-modul!", "lato-lato")
			if matchedTok, ok := l.matchHyphenatedKeyword(tok.Literal); ok {
				return matchedTok
			}

			// Check custom multi-word compound keywords with space (e.g. "tidak sama dengan", "selain jika")
			if matchedTok, ok := l.matchCompoundKeyword(tok.Literal); ok {
				return matchedTok
			}

			tok.Type = l.lookupIdent(tok.Literal)

			// Handle compound keywords: "selain jika", "lain jika", "else if", "atawa lamun"
			if ((tok.Literal == "selain" || tok.Literal == "lain" || tok.Literal == "else") && l.peekNextWord() == "jika") ||
				(tok.Literal == "else" && l.peekNextWord() == "if") ||
				(tok.Literal == "atawa" && l.peekNextWord() == "lamun") {
				l.skipWhitespaceAndComments()
				secondWord := l.readIdentifier()
				tok.Type = ELSE_IF
				tok.Literal = tok.Literal + " " + secondWord
			}

			// Handle compound keywords: "lamun henteu" (ELSE)
			if tok.Literal == "lamun" && l.peekNextWord() == "henteu" {
				l.skipWhitespaceAndComments()
				secondWord := l.readIdentifier()
				tok.Type = ELSE
				tok.Literal = tok.Literal + " " + secondWord
			}

			// Handle compound keywords: "kalau kosong", "bila kosong", "mun kosong"
			if (tok.Literal == "kalau" || tok.Literal == "bila" || tok.Literal == "mun") && (l.peekNextWord() == "kosong" || l.peekNextWord() == "nihil" || l.peekNextWord() == "null") {
				l.skipWhitespaceAndComments()
				secondWord := l.readIdentifier()
				tok.Type = NULL_COALESCE
				tok.Literal = tok.Literal + " " + secondWord
			}

			// Handle compound keywords: "saat selesai", "pas rengse"
			if (tok.Literal == "saat" && l.peekNextWord() == "selesai") || (tok.Literal == "pas" && l.peekNextWord() == "rengse") {
				l.skipWhitespaceAndComments()
				secondWord := l.readIdentifier()
				tok.Type = DEFER
				tok.Literal = tok.Literal + " " + secondWord
			}

			// Handle visibility modifiers: publik!, privat!, and privat-[...]
			if tok.Literal == "publik" && l.ch == '!' {
				l.readChar()
				tok.Literal = "publik!"
				tok.Type = VIS_PUBLIC
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

func (l *Lexer) matchHyphenatedKeyword(firstWord string) (Token, bool) {
	if l.ch != '-' {
		return Token{}, false
	}

	savedPos := l.position
	savedReadPos := l.readPosition
	savedCh := l.ch
	savedLine := l.line
	savedCol := l.column

	accumulated := firstWord
	bestMatchLit := ""
	var bestMatchTok TokenType
	var bestPos, bestReadPos int
	var bestCh byte
	var bestLine, bestCol int

	for l.ch == '-' && (isLetter(l.peekChar()) || l.peekChar() >= 0x80) {
		l.readChar() // eat '-'
		part := l.readIdentifier()
		if part == "" {
			break
		}
		accumulated += "-" + part

		// Check if accumulated matches a custom keyword
		if l.customKeywords != nil {
			if targetTok, ok := l.customKeywords[accumulated]; ok {
				bestMatchLit = accumulated
				bestMatchTok = targetTok
				bestPos = l.position
				bestReadPos = l.readPosition
				bestCh = l.ch
				bestLine = l.line
				bestCol = l.column
			}
		}
		// Also check built-in keywords
		if targetTok, ok := keywords[accumulated]; ok {
			bestMatchLit = accumulated
			bestMatchTok = targetTok
			bestPos = l.position
			bestReadPos = l.readPosition
			bestCh = l.ch
			bestLine = l.line
			bestCol = l.column
		}

		// Also check if followed by '!' (e.g. "privat-modul!")
		if l.ch == '!' {
			withBang := accumulated + "!"
			if targetTok, ok := keywords[withBang]; ok {
				l.readChar() // eat '!'
				bestMatchLit = withBang
				bestMatchTok = targetTok
				bestPos = l.position
				bestReadPos = l.readPosition
				bestCh = l.ch
				bestLine = l.line
				bestCol = l.column
				break
			}
			if l.customKeywords != nil {
				if targetTok, ok := l.customKeywords[withBang]; ok {
					l.readChar() // eat '!'
					bestMatchLit = withBang
					bestMatchTok = targetTok
					bestPos = l.position
					bestReadPos = l.readPosition
					bestCh = l.ch
					bestLine = l.line
					bestCol = l.column
					break
				}
			}
		}
	}

	if bestMatchLit != "" {
		l.position = bestPos
		l.readPosition = bestReadPos
		l.ch = bestCh
		l.line = bestLine
		l.column = bestCol
		return Token{Type: bestMatchTok, Literal: bestMatchLit, Line: savedLine, Column: savedCol}, true
	}

	// Restore state if not matched
	l.position = savedPos
	l.readPosition = savedReadPos
	l.ch = savedCh
	l.line = savedLine
	l.column = savedCol
	return Token{}, false
}

func (l *Lexer) matchCompoundKeyword(firstWord string) (Token, bool) {
	if l.customCompounds == nil {
		return Token{}, false
	}
	remainderMap, ok := l.customCompounds[firstWord]
	if !ok {
		return Token{}, false
	}

	savedPos := l.position
	savedReadPos := l.readPosition
	savedCh := l.ch
	savedLine := l.line
	savedCol := l.column

	accumulated := ""
	bestMatchLit := ""
	var bestMatchTok TokenType
	var bestPos, bestReadPos int
	var bestCh byte
	var bestLine, bestCol int

	for {
		if l.ch != ' ' && l.ch != '\t' {
			break
		}
		// skip spaces and tabs only (not newlines)
		for l.ch == ' ' || l.ch == '\t' {
			l.readChar()
		}
		if !isLetter(l.ch) && l.ch < 0x80 {
			break
		}

		word := l.readIdentifier()
		if word == "" {
			break
		}
		if accumulated == "" {
			accumulated = word
		} else {
			accumulated += " " + word
		}

		if targetTok, exists := remainderMap[accumulated]; exists {
			bestMatchLit = firstWord + " " + accumulated
			bestMatchTok = targetTok
			bestPos = l.position
			bestReadPos = l.readPosition
			bestCh = l.ch
			bestLine = l.line
			bestCol = l.column
		}
	}

	if bestMatchLit != "" {
		l.position = bestPos
		l.readPosition = bestReadPos
		l.ch = bestCh
		l.line = bestLine
		l.column = bestCol
		return Token{Type: bestMatchTok, Literal: bestMatchLit, Line: savedLine, Column: savedCol}, true
	}

	l.position = savedPos
	l.readPosition = savedReadPos
	l.ch = savedCh
	l.line = savedLine
	l.column = savedCol
	return Token{}, false
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for {
		if l.ch == 0 {
			break
		}
		if l.ch < 0x80 {
			if isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
				l.readChar()
				continue
			}
			break
		}
		// Multi-byte UTF-8 rune
		r, width := utf8.DecodeRuneInString(l.input[l.position:])
		if r == utf8.RuneError || !isIdentifierPartRune(r) {
			break
		}
		for i := 0; i < width; i++ {
			l.readChar()
		}
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
	isStart := false
	if l.ch < 0x80 {
		isStart = isLetter(l.ch)
	} else {
		r, _ := utf8.DecodeRuneInString(l.input[l.position:])
		isStart = isIdentifierStartRune(r)
	}
	if isStart {
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

func (l *Lexer) lookupIdent(ident string) TokenType {
	if l != nil && l.customKeywords != nil {
		if tok, ok := l.customKeywords[ident]; ok {
			return tok
		}
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isEmojiRune(r rune) bool {
	return (r >= 0x1F000 && r <= 0x1FAFF) || // Miscellaneous Symbols and Pictographs, Emoticons, Transport, etc.
		(r >= 0x2600 && r <= 0x27BF) ||  // Dingbats, Miscellaneous Symbols
		(r >= 0x2B00 && r <= 0x2BFF) ||  // Miscellaneous Symbols and Arrows
		(r >= 0x2300 && r <= 0x23FF) ||  // Miscellaneous Technical
		(r >= 0xFE00 && r <= 0xFE0F) ||  // Variation Selectors
		(r >= 0x200D && r <= 0x200D)     // Zero Width Joiner
}

func isIdentifierStartRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || isEmojiRune(r)
}

func isIdentifierPartRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || isEmojiRune(r)
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

