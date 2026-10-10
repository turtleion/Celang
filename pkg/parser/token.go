package parser

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers & Literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	BANG     TokenType = "!"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="

	PLUS_ASSIGN  TokenType = "+="
	MINUS_ASSIGN TokenType = "-="
	ARROW        TokenType = "->"
	PIPE         TokenType = "|"

	// Delimiters
	COMMA     TokenType = ","
	COLON     TokenType = ":"
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	DOT       TokenType = "."

	// Keywords
	VAR_DECL    TokenType = "VAR_DECL"    // buat, let, var
	CONST_DECL  TokenType = "CONST_DECL"  // konstanta, const
	TYPE_ASSIGN TokenType = "TYPE_ASSIGN" // itu
	IF          TokenType = "IF"          // jika, if
	ELSE_IF     TokenType = "ELSE_IF"     // selain jika, elif
	ELSE        TokenType = "ELSE"        // selain, else
	WHILE       TokenType = "WHILE"       // selama, while
	FOR         TokenType = "FOR"         // untuk, for
	IN          TokenType = "IN"          // dalam, in
	FUNC        TokenType = "FUNC"        // fungsi, def, fn, function
	RETURN      TokenType = "RETURN"      // kembalikan, return
	ASYNC       TokenType = "ASYNC"       // sambil, async
	PARALLEL    TokenType = "PARALLEL"    // bersamaan, parallel
	TRUE        TokenType = "TRUE"        // benar, true, True
	FALSE       TokenType = "FALSE"       // salah, false, False
	NULL        TokenType = "NULL"        // kosong, nihil, null, None
	AND         TokenType = "AND"         // dan, and, &&
	OR          TokenType = "OR"          // atau, or, ||
	NOT         TokenType = "NOT"         // bukan, tidak, not, !
	RETURN_TYPE TokenType = "RETURN_TYPE" // mengembalikan
	MAKA        TokenType = "MAKA"        // maka
	SELESAI     TokenType = "SELESAI"     // selesai
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
