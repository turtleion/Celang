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
	FAT_ARROW    TokenType = "=>"
	PIPE         TokenType = "|"
	PIPE_ARROW   TokenType = "|>"
	AMPERSAND    TokenType = "&"
	AT           TokenType = "@"

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
	DOT          TokenType = "."
	RANGE        TokenType = ".."
	OPTIONAL_DOT TokenType = "?."

	// Keywords
	VAR_DECL    TokenType = "VAR_DECL"    // buat, let, var
	CONST_DECL  TokenType = "CONST_DECL"  // konstanta, const
	TYPE_ASSIGN TokenType = "TYPE_ASSIGN" // itu
	IF          TokenType = "IF"          // jika, if
	ELSE_IF     TokenType = "ELSE_IF"     // selain jika, elif
	ELSE        TokenType = "ELSE"        // selain, else
	SWITCH      TokenType = "SWITCH"      // cocok, switch
	CASE        TokenType = "CASE"        // kasus, case
	DEFAULT     TokenType = "DEFAULT"     // lainnya, default
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

	// Visibility
	VIS_PUBLIC    TokenType = "VIS_PUBLIC"    // publik!, export
	VIS_MODULE    TokenType = "VIS_MODULE"    // privat-modul!
	VIS_PRIVATE   TokenType = "VIS_PRIVATE"   // privat!
	VIS_ALLOWLIST TokenType = "VIS_ALLOWLIST" // privat-

	// Modules & Imports
	IMPORT TokenType = "IMPORT" // ambil, impor, import
	FROM   TokenType = "FROM"   // dari, from
	AS     TokenType = "AS"     // sebagai, as

	// Error Handling
	TRY     TokenType = "TRY"     // coba, try
	CATCH   TokenType = "CATCH"   // tangkap, catch, except
	FINALLY TokenType = "FINALLY" // akhirnya, finally
	THROW   TokenType = "THROW"   // lempar, throw, raise

	// In-Code Styling
	STYLE_BLOCK TokenType = "STYLE_BLOCK" // gaya, style, dialek

	// Data Structures / Structs
	STRUCT_DECL TokenType = "STRUCT_DECL" // bentuk, struct
	PRIVATE     TokenType = "PRIVATE"     // privat
	REF_DECL    TokenType = "REF_DECL"    // rujuk, ref
	TYPE_DECL   TokenType = "TYPE_DECL"   // tipe

	// Human Ergonomics & Null Safety
	NULL_COALESCE TokenType = "NULL_COALESCE" // kalau kosong, bila kosong, ??
	MAYBE         TokenType = "MAYBE"         // mungkin
	DEFER         TokenType = "DEFER"         // saat_selesai, defer

	// Modern Concurrency & Generics
	TITIP     TokenType = "TITIP"     // titip
	JEMPUT    TokenType = "JEMPUT"    // jemput
	HARUS     TokenType = "HARUS"     // harus
	SEMBARANG TokenType = "SEMBARANG" // sembarang
	SALURAN   TokenType = "SALURAN"   // saluran
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
