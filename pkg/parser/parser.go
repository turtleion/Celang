package parser

import (
	"fmt"
	"strconv"
	"strings"

	"celang/pkg/uast"
)

const (
	_ int = iota
	LOWEST
	PIPELINE_PREC // |>
	COALESCE_PREC // kalau kosong, ??
	OR_PREC       // atau, ||
	AND_PREC      // dan, &&
	RANGE_PREC    // ..
	EQUALS        // ==, !=
	LESSGREATER   // <, >, <=, >=
	SUM           // +, -
	PRODUCT       // *, /, %
	CAST_PREC     // sebagai, as
	PREFIX        // -x, !x, bukan x, mungkin x
	CALL          // hitung(x)
	INDEX         // arr[0], map["kunci"], obj.prop, obj?.prop
)

var precedences = map[TokenType]int{
	PIPE_ARROW:    PIPELINE_PREC,
	NULL_COALESCE: COALESCE_PREC,
	OR:            OR_PREC,
	PIPE:          OR_PREC,
	AND:           AND_PREC,
	RANGE:         RANGE_PREC,
	EQ:            EQUALS,
	NOT_EQ:        EQUALS,
	LT:            LESSGREATER,
	GT:            LESSGREATER,
	LT_EQ:         LESSGREATER,
	GT_EQ:         LESSGREATER,
	PLUS:          SUM,
	MINUS:         SUM,
	SLASH:         PRODUCT,
	ASTERISK:      PRODUCT,
	PERCENT:       PRODUCT,
	AS:            CAST_PREC,
	LPAREN:        CALL,
	LBRACKET:      INDEX,
	DOT:           INDEX,
	OPTIONAL_DOT:  INDEX,
}

type (
	prefixParseFn func() uast.Expression
	infixParseFn  func(uast.Expression) uast.Expression
)

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn

	parenDepth int
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(NULL, p.parseNullLiteral)
	p.registerPrefix(BANG, p.parsePrefixExpression)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(AMPERSAND, p.parsePrefixExpression)
	p.registerPrefix(ASTERISK, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(LBRACE, p.parseMapLiteral)
	p.registerPrefix(FUNC, p.parseFunctionExpression)
	p.registerPrefix(AT, p.parseAtExpression)
	p.registerPrefix(MAYBE, p.parseSafeNavigationExpression)
	p.registerPrefix(RANGE, p.parsePrefixRangeExpression)
	p.registerPrefix(TITIP, p.parseTitipExpression)
	p.registerPrefix(JEMPUT, p.parseJemputExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PIPE_ARROW, p.parsePipelineExpression)
	p.registerInfix(NULL_COALESCE, p.parseNullCoalesceExpression)
	p.registerInfix(RANGE, p.parseInfixRangeExpression)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(PERCENT, p.parseInfixExpression)
	p.registerInfix(AS, p.parseCastExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LT_EQ, p.parseInfixExpression)
	p.registerInfix(GT_EQ, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)
	p.registerInfix(DOT, p.parseMemberExpression)
	p.registerInfix(OPTIONAL_DOT, p.parseOptionalMemberExpression)

	// Read two tokens to populate curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	errStr := fmt.Sprintf("[Baris %d, Kolom %d] %s", p.curToken.Line, p.curToken.Column, msg)
	p.errors = append(p.errors, errStr)
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("Diharapkan token '%s', tetapi menemukan '%s' (%s)", t, p.peekToken.Type, p.peekToken.Literal)
	p.addError(msg)
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the full source code into a U-AST Program.
func (p *Parser) ParseProgram() (*uast.Program, error) {
	program := uast.NewProgram()

	for !p.curTokenIs(EOF) {
		// Skip semicolons if present
		if p.curTokenIs(SEMICOLON) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Body = append(program.Body, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parsing gagal dengan %d kesalahan:\n%s", len(p.errors), strings.Join(p.errors, "\n"))
	}

	return program, nil
}

func (p *Parser) parseStatement() uast.Statement {
	switch p.curToken.Type {
	case VAR_DECL:
		return p.parseVariableDeclaration(true)
	case CONST_DECL:
		return p.parseVariableDeclaration(false)
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForEachStatement()
	case FUNC:
		return p.parseFunctionDeclaration()
	case RETURN:
		return p.parseReturnStatement()
	case ASYNC:
		return p.parseConcurrencyStatement("sambil")
	case PARALLEL:
		return p.parseConcurrencyStatement("bersamaan")
	case SWITCH:
		return p.parseSwitchStatement()
	case VIS_PUBLIC, VIS_MODULE, VIS_PRIVATE, VIS_ALLOWLIST:
		return p.parseVisibilityOrExportStatement()
	case IMPORT:
		return p.parseImportStatement()
	case FROM:
		return p.parsePythonicImportStatement()
	case TRY:
		return p.parseTryCatchStatement()
	case THROW:
		return p.parseThrowStatement()
	case STYLE_BLOCK:
		return p.parseStyleBlockStatement()
	case STRUCT_DECL:
		return p.parseStructDeclaration()
	case TYPE_DECL:
		return p.parseTypeDeclaration()
	case REF_DECL:
		return p.parseReferenceDeclaration()
	case DEFER:
		return p.parseDeferStatement()
	case TRANSACTION:
		return p.parseTransactionStatement()
	case ROLLBACK:
		return p.parseRollbackStatement()
	case REACTIVE_DECL:
		return p.parseReactiveDeclaration()
	case WATCH:
		return p.parseWatchStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVisibilityOrExportStatement() uast.Statement {
	meta := p.currentMetadata()
	var vis *uast.Visibility

	switch p.curToken.Type {
	case VIS_PUBLIC:
		vis = &uast.Visibility{Kind: uast.VisibilityPublic}
		p.nextToken()
	case VIS_MODULE:
		vis = &uast.Visibility{Kind: uast.VisibilityModule}
		p.nextToken()
	case VIS_PRIVATE:
		vis = &uast.Visibility{Kind: uast.VisibilityPrivate}
		p.nextToken()
	case VIS_ALLOWLIST:
		p.nextToken() // move to '['
		if !p.curTokenIs(LBRACKET) {
			p.addError("Diharapkan '[' setelah 'privat-'")
			return nil
		}
		p.nextToken() // move inside bracket
		targets := []string{}
		for !p.curTokenIs(RBRACKET) && !p.curTokenIs(EOF) {
			if p.curTokenIs(STRING) || p.curTokenIs(IDENT) {
				targets = append(targets, p.curToken.Literal)
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		if p.peekTokenIs(BANG) {
			p.nextToken() // eat '!'
		}
		p.nextToken() // move to next token after ']' or '!'
		vis = &uast.Visibility{Kind: uast.VisibilityAllowlist, Targets: targets}
	}

	// Standalone export list: publik! { kali, bagi }
	if p.curTokenIs(LBRACE) {
		p.nextToken() // enter brace
		idents := []string{}
		for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				idents = append(idents, p.curToken.Literal)
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		return &uast.VisibilityStatement{
			NodeType:    uast.NodeTypeVisibilityStatement,
			Metadata:    meta,
			Visibility:  *vis,
			Identifiers: idents,
		}
	}

	// Declaration with visibility
	switch p.curToken.Type {
	case VAR_DECL:
		decl := p.parseVariableDeclaration(true)
		if d, ok := decl.(*uast.DestructureDeclaration); ok {
			d.Visibility = vis
		} else if v, ok := decl.(*uast.VariableDeclaration); ok {
			v.Visibility = vis
		}
		return decl
	case CONST_DECL:
		decl := p.parseVariableDeclaration(false)
		if d, ok := decl.(*uast.DestructureDeclaration); ok {
			d.Visibility = vis
		} else if v, ok := decl.(*uast.VariableDeclaration); ok {
			v.Visibility = vis
		}
		return decl
	case FUNC:
		fn := p.parseFunctionDeclaration()
		if fn != nil {
			fn.Visibility = vis
		}
		return fn
	case STRUCT_DECL:
		s := p.parseStructDeclaration()
		if s != nil {
			s.Visibility = vis
		}
		return s
	case TYPE_DECL:
		t := p.parseTypeDeclaration()
		if s, ok := t.(*uast.StructDeclaration); ok && s != nil {
			s.Visibility = vis
		}
		return t
	default:
		p.addError(fmt.Sprintf("Modifier visibilitas tidak dapat diterapkan pada '%s'", p.curToken.Literal))
		return nil
	}
}

func (p *Parser) parseImportStatement() *uast.ImportDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'ambil' / 'impor' / 'import'

	specifiers := []uast.ImportSpecifier{}

	// Case 1: Namespace import: ambil * sebagai math dari "./mod.ce"
	if p.curTokenIs(ASTERISK) {
		p.nextToken() // eat '*'
		if !p.curTokenIs(AS) {
			p.addError("Diharapkan 'sebagai' atau 'as' setelah '*' dalam impor")
			return nil
		}
		p.nextToken() // move to alias identifier
		if !p.curTokenIs(IDENT) {
			p.addError("Diharapkan nama alias identifier setelah 'sebagai'")
			return nil
		}
		local := p.curToken.Literal
		specifiers = append(specifiers, uast.ImportSpecifier{Imported: "*", Local: local, IsAll: true})
		p.nextToken() // move to 'dari' / 'from'
	} else if p.curTokenIs(LBRACE) {
		// Case 2: Named import: ambil { tambah, kali sebagai perkalian } dari "./mod.ce"
		p.nextToken() // enter '{'
		for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				imported := p.curToken.Literal
				local := imported
				if p.peekTokenIs(AS) {
					p.nextToken() // eat 'sebagai' / 'as'
					p.nextToken() // move to alias
					if p.curTokenIs(IDENT) {
						local = p.curToken.Literal
					}
				}
				specifiers = append(specifiers, uast.ImportSpecifier{Imported: imported, Local: local})
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		p.nextToken() // move to 'dari' / 'from'
	} else if p.curTokenIs(STRING) {
		// Case 3: Quoted path import: import "./mod.ce" as math
		source := p.curToken.Literal
		if p.peekTokenIs(AS) && p.peekToken.Line == p.curToken.Line {
			p.nextToken() // eat 'as'
			p.nextToken() // move to alias
			local := p.curToken.Literal
			specifiers = append(specifiers, uast.ImportSpecifier{Imported: "*", Local: local, IsAll: true})
		}
		return &uast.ImportDeclaration{
			NodeType:   uast.NodeTypeImportDeclaration,
			Metadata:   meta,
			Source:     source,
			Specifiers: specifiers,
		}
	} else if p.curTokenIs(IDENT) {
		firstIdent := p.curToken.Literal
		firstLocal := firstIdent
		if p.peekTokenIs(AS) && p.peekToken.Line == p.curToken.Line {
			p.nextToken() // eat 'as' / 'sebagai'
			p.nextToken() // move to alias
			if p.curTokenIs(IDENT) {
				firstLocal = p.curToken.Literal
			}
		}

		if (!p.peekTokenIs(FROM) && !p.peekTokenIs(COMMA)) || p.peekToken.Line > p.curToken.Line {
			// Direct module import: impor matematika [sebagai math]
			specifiers = append(specifiers, uast.ImportSpecifier{Imported: "*", Local: firstLocal, IsAll: true})
			return &uast.ImportDeclaration{
				NodeType:   uast.NodeTypeImportDeclaration,
				Metadata:   meta,
				Source:     firstIdent,
				Specifiers: specifiers,
			}
		}

		// Otherwise: ambil foo [sebagai bar], baz dari ...
		specifiers = append(specifiers, uast.ImportSpecifier{Imported: firstIdent, Local: firstLocal})
		for p.peekTokenIs(COMMA) && p.peekToken.Line == p.curToken.Line {
			p.nextToken() // eat ','
			p.nextToken() // move to next ident
			if !p.curTokenIs(IDENT) {
				p.addError("Diharapkan nama simbol setelah koma dalam impor")
				return nil
			}
			imp := p.curToken.Literal
			loc := imp
			if p.peekTokenIs(AS) && p.peekToken.Line == p.curToken.Line {
				p.nextToken() // eat 'as' / 'sebagai'
				p.nextToken() // move to alias
				if p.curTokenIs(IDENT) {
					loc = p.curToken.Literal
				}
			}
			specifiers = append(specifiers, uast.ImportSpecifier{Imported: imp, Local: loc})
		}
		p.nextToken() // move to 'dari' / 'from'
	} else {
		p.addError(fmt.Sprintf("Sintaks impor tidak valid pada token '%s'", p.curToken.Literal))
		return nil
	}

	if !p.curTokenIs(FROM) {
		p.addError("Diharapkan keyword 'dari' atau 'from' setelah daftar impor")
		return nil
	}
	p.nextToken() // move to source path or module identifier

	if !p.curTokenIs(STRING) && !p.curTokenIs(IDENT) {
		p.addError("Diharapkan nama modul (contoh: matematika) atau path berkas (contoh: \"./modul.ce\")")
		return nil
	}
	source := p.curToken.Literal

	return &uast.ImportDeclaration{
		NodeType:   uast.NodeTypeImportDeclaration,
		Metadata:   meta,
		Source:     source,
		Specifiers: specifiers,
	}
}

func (p *Parser) parsePythonicImportStatement() *uast.ImportDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'dari' / 'from'

	if !p.curTokenIs(STRING) && !p.curTokenIs(IDENT) {
		p.addError("Diharapkan nama modul (contoh: matematika) atau path berkas setelah 'from' / 'dari'")
		return nil
	}
	source := p.curToken.Literal
	p.nextToken() // move to 'import' / 'ambil'

	if !p.curTokenIs(IMPORT) {
		p.addError("Diharapkan 'import' atau 'ambil' setelah path modul")
		return nil
	}
	p.nextToken() // move to first specifier

	specifiers := []uast.ImportSpecifier{}

	if p.curTokenIs(ASTERISK) {
		specifiers = append(specifiers, uast.ImportSpecifier{Imported: "*", Local: "*", IsAll: true})
	} else if p.curTokenIs(LBRACE) {
		// Case: dari matematika ambil { akar, PI as pi }
		p.nextToken() // enter '{'
		for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				imported := p.curToken.Literal
				local := imported
				if p.peekTokenIs(AS) {
					p.nextToken() // eat 'as' / 'sebagai'
					p.nextToken() // move to alias
					if p.curTokenIs(IDENT) {
						local = p.curToken.Literal
					}
				}
				specifiers = append(specifiers, uast.ImportSpecifier{Imported: imported, Local: local})
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
	} else {
		for !p.curTokenIs(SEMICOLON) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				imported := p.curToken.Literal
				local := imported
				if p.peekTokenIs(AS) && p.peekToken.Line == p.curToken.Line {
					p.nextToken() // eat 'as'
					p.nextToken() // move to alias
					if p.curTokenIs(IDENT) {
						local = p.curToken.Literal
					}
				}
				specifiers = append(specifiers, uast.ImportSpecifier{Imported: imported, Local: local})
			}
			if p.peekTokenIs(COMMA) && p.peekToken.Line == p.curToken.Line {
				p.nextToken()
				p.nextToken()
			} else {
				break
			}
		}
	}

	return &uast.ImportDeclaration{
		NodeType:   uast.NodeTypeImportDeclaration,
		Metadata:   meta,
		Source:     source,
		Specifiers: specifiers,
	}
}

func (p *Parser) parseTryCatchStatement() *uast.TryCatchStatement {
	meta := p.currentMetadata()
	p.nextToken() // move to '{'

	if !p.curTokenIs(LBRACE) {
		p.addError("Diharapkan '{' setelah keyword 'coba' / 'try'")
		return nil
	}

	block := p.parseBlockStatement()

	var handler *uast.CatchClause
	if p.peekTokenIs(CATCH) {
		p.nextToken() // move to 'tangkap' / 'catch' / 'except'
		param := ""
		if p.peekTokenIs(LPAREN) {
			p.nextToken() // eat '('
			if p.peekTokenIs(IDENT) {
				p.nextToken()
				param = p.curToken.Literal
			}
			if !p.expectPeek(RPAREN) {
				return nil
			}
		} else if p.peekTokenIs(IDENT) {
			p.nextToken()
			param = p.curToken.Literal
			if p.peekTokenIs(COLON) {
				p.nextToken() // optional ':' for pythonic
			}
		}

		if !p.expectPeek(LBRACE) {
			return nil
		}
		catchBody := p.parseBlockStatement()
		handler = &uast.CatchClause{
			NodeType: uast.NodeTypeCatchClause,
			Param:    param,
			Body:     catchBody,
		}
	}

	var finalizer *uast.BlockStatement
	if p.peekTokenIs(FINALLY) {
		p.nextToken() // move to 'akhirnya' / 'finally'
		if !p.expectPeek(LBRACE) {
			return nil
		}
		finalizer = p.parseBlockStatement()
	}

	return &uast.TryCatchStatement{
		NodeType:  uast.NodeTypeTryCatchStatement,
		Metadata:  meta,
		Block:     block,
		Handler:   handler,
		Finalizer: finalizer,
	}
}

func (p *Parser) parseThrowStatement() *uast.ThrowStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'lempar' / 'throw' / 'raise'

	arg := p.parseExpression(LOWEST)

	return &uast.ThrowStatement{
		NodeType: uast.NodeTypeThrowStatement,
		Metadata: meta,
		Argument: arg,
	}
}

func (p *Parser) parseTransactionStatement() *uast.TransactionStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'transaksi' / 'transaction'

	if !p.curTokenIs(LBRACE) {
		p.addError(fmt.Sprintf("Diharapkan '{' setelah keyword 'transaksi', menerima '%s'", p.curToken.Literal))
		return nil
	}
	body := p.parseBlockStatement()

	stmt := &uast.TransactionStatement{
		NodeType: uast.NodeTypeTransactionStatement,
		Metadata: meta,
		Body:     body,
	}

	if p.peekTokenIs(RECOVER) {
		p.nextToken() // move to 'pulihkan' / 'recover'

		recoverVar := "err"
		if p.peekTokenIs(LPAREN) {
			p.nextToken() // move to '('
			if p.peekTokenIs(IDENT) {
				p.nextToken()
				recoverVar = p.curToken.Literal
			}
			if !p.expectPeek(RPAREN) {
				return nil
			}
		}

		if !p.expectPeek(LBRACE) {
			return nil
		}
		stmt.RecoverVar = recoverVar
		stmt.Recover = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseRollbackStatement() *uast.RollbackStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'batalkan' / 'rollback'

	var arg uast.Expression
	if !p.curTokenIs(SEMICOLON) && !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		arg = p.parseExpression(LOWEST)
	}

	return &uast.RollbackStatement{
		NodeType: uast.NodeTypeRollbackStatement,
		Metadata: meta,
		Argument: arg,
	}
}

func (p *Parser) parseReactiveDeclaration() *uast.ReactiveDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'reaktif' / 'reactive'

	if !p.curTokenIs(IDENT) {
		p.addError(fmt.Sprintf("Diharapkan nama variabel setelah 'reaktif', menerima '%s'", p.curToken.Literal))
		return nil
	}
	name := p.curToken.Literal

	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken() // move past '='

	val := p.parseExpression(LOWEST)

	return &uast.ReactiveDeclaration{
		NodeType: uast.NodeTypeReactiveDeclaration,
		Metadata: meta,
		Name:     name,
		Value:    val,
	}
}

func (p *Parser) parseWatchStatement() *uast.WatchStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'pantau' / 'watch'

	if !p.curTokenIs(IDENT) {
		p.addError(fmt.Sprintf("Diharapkan nama target variabel setelah 'pantau', menerima '%s'", p.curToken.Literal))
		return nil
	}
	target := p.curToken.Literal

	newVar := "baru"
	oldVar := "lama"

	if p.peekTokenIs(LPAREN) {
		p.nextToken() // move to '('
		if p.peekTokenIs(IDENT) {
			p.nextToken()
			newVar = p.curToken.Literal
			if p.peekTokenIs(COMMA) {
				p.nextToken() // eat ','
				if p.expectPeek(IDENT) {
					oldVar = p.curToken.Literal
				}
			}
		}
		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}
	body := p.parseBlockStatement()

	return &uast.WatchStatement{
		NodeType:   uast.NodeTypeWatchStatement,
		Metadata:   meta,
		Target:     target,
		NewValVar:  newVar,
		OldValVar:  oldVar,
		Body:       body,
	}
}

func (p *Parser) parseStyleBlockStatement() *uast.StyleBlockStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'gaya' / 'style' / 'dialek'

	if !p.curTokenIs(STRING) {
		p.addError(fmt.Sprintf("Diharapkan nama gaya berupa teks (string), menerima '%s'", p.curToken.Literal))
		return nil
	}
	styleName := p.curToken.Literal

	if p.peekTokenIs(COLON) {
		p.nextToken() // optional colon
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &uast.StyleBlockStatement{
		NodeType:  uast.NodeTypeStyleBlockStatement,
		Metadata:  meta,
		StyleName: styleName,
		Body:      body,
	}
}


func (p *Parser) parseVariableDeclaration(isMutable bool) uast.Statement {
	meta := p.currentMetadata()
	p.nextToken() // eat VAR_DECL or CONST_DECL

	// 1. Array destructuring: buat [x, y] = koordinat
	if p.curTokenIs(LBRACKET) {
		p.nextToken() // eat '['
		var bindings []uast.DestructureBinding
		for !p.curTokenIs(RBRACKET) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				bindings = append(bindings, uast.DestructureBinding{Target: p.curToken.Literal})
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		if !p.expectPeek(ASSIGN) {
			return nil
		}
		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		return &uast.DestructureDeclaration{
			NodeType:  uast.NodeTypeDestructureDeclaration,
			Metadata:  meta,
			Kind:      uast.DestructureArray,
			Bindings:  bindings,
			Value:     value,
			IsMutable: isMutable,
		}
	}

	// 2. Map destructuring: buat { nama, peran } = profil
	if p.curTokenIs(LBRACE) {
		p.nextToken() // eat '{'
		var bindings []uast.DestructureBinding
		for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				srcKey := p.curToken.Literal
				target := srcKey
				if p.peekTokenIs(COLON) || p.peekTokenIs(AS) {
					p.nextToken() // eat ':' or 'sebagai'
					if p.expectPeek(IDENT) {
						target = p.curToken.Literal
					}
				}
				bindings = append(bindings, uast.DestructureBinding{SourceKey: srcKey, Target: target})
			}
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		if !p.expectPeek(ASSIGN) {
			return nil
		}
		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		return &uast.DestructureDeclaration{
			NodeType:  uast.NodeTypeDestructureDeclaration,
			Metadata:  meta,
			Kind:      uast.DestructureMap,
			Bindings:  bindings,
			Value:     value,
			IsMutable: isMutable,
		}
	}

	if !p.curTokenIs(IDENT) {
		p.addError(fmt.Sprintf("Diharapkan nama variabel (identifier), menerima '%s'", p.curToken.Literal))
		return nil
	}

	name := p.curToken.Literal
	typeHint := ""

	// Check for type annotation: "buat umur itu Angka" or "buat umur: Angka"
	if p.peekTokenIs(TYPE_ASSIGN) || p.peekTokenIs(COLON) {
		p.nextToken() // eat 'itu' or ':'
		p.nextToken()
		typeHint = p.parseTypeHint()
	}

	var value uast.Expression
	if p.peekTokenIs(ASSIGN) {
		p.nextToken() // eat '='
		p.nextToken() // move to value
		value = p.parseExpression(LOWEST)
	}

	return &uast.VariableDeclaration{
		NodeType:  uast.NodeTypeVariableDeclaration,
		Metadata:  meta,
		Name:      name,
		Value:     value,
		IsMutable: isMutable,
		TypeHint:  typeHint,
	}
}

func (p *Parser) parseStructDeclaration() *uast.StructDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'bentuk' / 'struct'

	if !p.curTokenIs(IDENT) {
		p.addError(fmt.Sprintf("Diharapkan nama bentuk (struct), menerima '%s'", p.curToken.Literal))
		return nil
	}

	name := p.curToken.Literal
	var typeParams []uast.TypeParameter
	if p.peekTokenIs(LT) {
		typeParams = p.parseTypeParameters()
	}
	decl := p.parseStructBody(meta, name)
	if decl != nil {
		decl.TypeParameters = typeParams
	}
	return decl
}

func (p *Parser) parseTypeDeclaration() uast.Statement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'tipe'
	if !p.curTokenIs(IDENT) {
		p.addError("Diharapkan nama tipe setelah 'tipe'")
		return nil
	}
	name := p.curToken.Literal

	var typeParams []uast.TypeParameter
	if p.peekTokenIs(LT) {
		typeParams = p.parseTypeParameters()
	}

	if p.peekTokenIs(ASSIGN) {
		p.nextToken() // eat '='
		p.nextToken() // move to base type
		baseType := p.parseTypeHint()
		return &uast.TypeAliasDeclaration{
			NodeType: uast.NodeTypeTypeAliasDeclaration,
			Metadata: meta,
			Name:     name,
			BaseType: baseType,
		}
	}

	// Legacy or Go-style: tipe Point struct { ... } or tipe Point { ... }
	if p.peekTokenIs(STRUCT_DECL) {
		p.nextToken() // eat 'struct'
	}
	decl := p.parseStructBody(meta, name)
	if decl != nil {
		decl.TypeParameters = typeParams
	}
	return decl
}

func (p *Parser) parseReferenceDeclaration() *uast.ReferenceDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'rujuk' / 'ref'
	if !p.curTokenIs(IDENT) {
		p.addError("Diharapkan nama variabel alias setelah 'rujuk'")
		return nil
	}
	target := p.curToken.Literal

	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken() // move to source expression
	source := p.parseExpression(LOWEST)

	return &uast.ReferenceDeclaration{
		NodeType: uast.NodeTypeReferenceDeclaration,
		Metadata: meta,
		Target:   target,
		Source:   source,
	}
}

func (p *Parser) parseStructBody(meta *uast.NodeMetadata, name string) *uast.StructDeclaration {
	if !p.expectPeek(LBRACE) {
		return nil
	}
	p.nextToken() // enter '{'

	var fields []uast.StructField
	var methods []uast.StructMethod
	var embeds []string

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(SEMICOLON) || p.curTokenIs(COMMA) {
			p.nextToken()
			continue
		}

		if p.curTokenIs(PRIVATE) {
			if p.peekTokenIs(LBRACE) {
				p.nextToken() // eat 'privat'
				p.nextToken() // enter '{' of private block
				for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
					if p.curTokenIs(SEMICOLON) || p.curTokenIs(COMMA) {
						p.nextToken()
						continue
					}
					p.parseStructMember(true, &fields, &methods, &embeds)
					p.nextToken()
				}
				p.nextToken() // eat '}' of private block
				continue
			} else {
				p.nextToken() // eat 'privat'
				p.parseStructMember(true, &fields, &methods, &embeds)
				p.nextToken()
				continue
			}
		}

		p.parseStructMember(false, &fields, &methods, &embeds)
		p.nextToken()
	}

	return &uast.StructDeclaration{
		NodeType: uast.NodeTypeStructDeclaration,
		Metadata: meta,
		Name:     name,
		Fields:   fields,
		Methods:  methods,
		Embeds:   embeds,
	}
}

func (p *Parser) parseStructMember(isPrivate bool, fields *[]uast.StructField, methods *[]uast.StructMethod, embeds *[]string) {
	if p.curTokenIs(FUNC) {
		p.nextToken() // eat 'fungsi'
		if !p.curTokenIs(IDENT) {
			p.addError("Diharapkan nama fungsi pada bentuk")
			return
		}
		methodName := p.curToken.Literal
		if !p.expectPeek(LPAREN) {
			return
		}
		params := p.parseFunctionParameters()
		retType := ""
		if p.peekTokenIs(RETURN_TYPE) || p.peekTokenIs(ARROW) || p.peekTokenIs(COLON) {
			p.nextToken()
			p.nextToken()
			retType = p.parseTypeHint()
		}
		if !p.expectPeek(LBRACE) {
			return
		}
		body := p.parseBlockStatement()
		*methods = append(*methods, uast.StructMethod{
			Name:       methodName,
			Parameters: params,
			ReturnType: retType,
			Body:       body,
			IsPrivate:  isPrivate,
		})
		return
	}

	if p.curTokenIs(IDENT) {
		if p.peekTokenIs(LPAREN) {
			// Method without 'fungsi' keyword: setor(jumlah) { ... }
			methodName := p.curToken.Literal
			p.nextToken() // move to LPAREN
			params := p.parseFunctionParameters()
			retType := ""
			if p.peekTokenIs(RETURN_TYPE) || p.peekTokenIs(ARROW) || p.peekTokenIs(COLON) {
				p.nextToken()
				p.nextToken()
				retType = p.parseTypeHint()
			}
			if !p.expectPeek(LBRACE) {
				return
			}
			body := p.parseBlockStatement()
			*methods = append(*methods, uast.StructMethod{
				Name:       methodName,
				Parameters: params,
				ReturnType: retType,
				Body:       body,
				IsPrivate:  isPrivate,
			})
			return
		}

		if p.peekTokenIs(COLON) || p.peekTokenIs(TYPE_ASSIGN) || p.peekTokenIs(ASSIGN) {
			// Struct field: nama: Teks, saldo itu Angka = 0, saldo = 0
			fieldName := p.curToken.Literal
			typeHint := ""
			var defVal uast.Expression

			if p.peekTokenIs(COLON) || p.peekTokenIs(TYPE_ASSIGN) {
				p.nextToken() // eat ':' or 'itu'
				p.nextToken()
				typeHint = p.parseTypeHint()
			}

			if p.peekTokenIs(ASSIGN) {
				p.nextToken() // eat '='
				p.nextToken() // move to default expr
				defVal = p.parseExpression(LOWEST)
			}

			*fields = append(*fields, uast.StructField{
				Name:         fieldName,
				TypeHint:     typeHint,
				DefaultValue: defVal,
				IsPrivate:    isPrivate,
			})
			return
		}

		// Embedded struct: e.g. Identitas, or Identitas
		*embeds = append(*embeds, p.curToken.Literal)
		return
	}
}

func (p *Parser) parseAssignmentStatement() uast.Statement {
	meta := p.currentMetadata()
	var target uast.Expression

	if p.curTokenIs(ASTERISK) {
		p.nextToken() // move to IDENT
		id := p.parseIdentifier()
		target = &uast.UnaryExpression{
			NodeType: uast.NodeTypeUnaryExpression,
			Metadata: meta,
			Operator: "*",
			Argument: id,
			IsPrefix: true,
		}
	} else {
		target = p.parseIdentifier()
	}

	if !p.peekTokenIs(ASSIGN) && !p.peekTokenIs(PLUS_ASSIGN) && !p.peekTokenIs(MINUS_ASSIGN) {
		return &uast.ExpressionStatement{
			NodeType:   uast.NodeTypeExpressionStatement,
			Metadata:   meta,
			Expression: target,
		}
	}

	p.nextToken() // move to operator (=, +=, -=)
	operator := p.curToken.Literal

	p.nextToken() // move to value
	val := p.parseExpression(LOWEST)

	return &uast.AssignmentStatement{
		NodeType: uast.NodeTypeAssignmentStatement,
		Metadata: meta,
		Target:   target,
		Operator: operator,
		Value:    val,
	}
}

func (p *Parser) parseSwitchStatement() *uast.SwitchStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'cocok' or 'switch'

	hasParen := false
	if p.curTokenIs(LPAREN) {
		hasParen = true
		p.nextToken()
	}

	discriminant := p.parseExpression(LOWEST)

	if hasParen {
		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}
	p.nextToken() // enter brace

	cases := []uast.CaseClause{}

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(CASE) {
			p.nextToken() // eat 'kasus' / 'case'
			caseVal := p.parseExpression(LOWEST)
			if p.peekTokenIs(COLON) {
				p.nextToken() // eat ':'
			}
			p.nextToken()

			body := []uast.Statement{}
			for !p.curTokenIs(CASE) && !p.curTokenIs(DEFAULT) && !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
				if p.curTokenIs(SEMICOLON) {
					p.nextToken()
					continue
				}
				stmt := p.parseStatement()
				if stmt != nil {
					body = append(body, stmt)
				}
				p.nextToken()
			}
			cases = append(cases, uast.CaseClause{
				NodeType: uast.NodeTypeCaseClause,
				Value:    caseVal,
				Body:     body,
			})
		} else if p.curTokenIs(DEFAULT) {
			p.nextToken() // eat 'lainnya' / 'default'
			if p.peekTokenIs(COLON) {
				p.nextToken() // eat ':'
			}
			p.nextToken()

			body := []uast.Statement{}
			for !p.curTokenIs(CASE) && !p.curTokenIs(DEFAULT) && !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
				if p.curTokenIs(SEMICOLON) {
					p.nextToken()
					continue
				}
				stmt := p.parseStatement()
				if stmt != nil {
					body = append(body, stmt)
				}
				p.nextToken()
			}
			cases = append(cases, uast.CaseClause{
				NodeType: uast.NodeTypeCaseClause,
				Value:    nil,
				Body:     body,
			})
		} else {
			p.nextToken()
		}
	}

	return &uast.SwitchStatement{
		NodeType:     uast.NodeTypeSwitchStatement,
		Metadata:     meta,
		Discriminant: discriminant,
		Cases:        cases,
	}
}

func (p *Parser) parseIfStatement() *uast.IfStatement {
	meta := p.currentMetadata()

	p.nextToken() // eat 'jika' or 'if'

	// Support optional parentheses around condition
	hasParen := false
	if p.curTokenIs(LPAREN) {
		hasParen = true
		p.nextToken() // eat '('
	}

	condition := p.parseExpression(LOWEST)

	if hasParen {
		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	// Optional 'maka' keyword
	if p.peekTokenIs(MAKA) {
		p.nextToken()
	}

	if p.peekTokenIs(COLON) {
		p.nextToken()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	consequent := p.parseBlockStatement()

	var alternate uast.Statement
	if p.peekTokenIs(ELSE_IF) {
		p.nextToken() // move to ELSE_IF ("selain jika" / "elif")
		alternate = p.parseIfStatement()
	} else if p.peekTokenIs(ELSE) {
		p.nextToken() // move to ELSE ("selain")
		if p.peekTokenIs(IF) {
			// In case written as two tokens "selain" "jika"
			p.nextToken()
			alternate = p.parseIfStatement()
		} else {
			if p.peekTokenIs(COLON) {
				p.nextToken()
			}
			if !p.expectPeek(LBRACE) {
				return nil
			}
			alternate = p.parseBlockStatement()
		}
	}

	return &uast.IfStatement{
		NodeType:   uast.NodeTypeIfStatement,
		Metadata:   meta,
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func (p *Parser) parseWhileStatement() *uast.WhileStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'selama' or 'while'

	hasParen := false
	if p.curTokenIs(LPAREN) {
		hasParen = true
		p.nextToken()
	}

	condition := p.parseExpression(LOWEST)

	if hasParen {
		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	if p.peekTokenIs(COLON) {
		p.nextToken()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &uast.WhileStatement{
		NodeType:  uast.NodeTypeWhileStatement,
		Metadata:  meta,
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseForEachStatement() *uast.ForEachStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'untuk' or 'for'

	if !p.curTokenIs(IDENT) {
		p.addError("Diharapkan variabel iterasi setelah 'untuk'")
		return nil
	}
	itemName := p.curToken.Literal

	if !p.expectPeek(IN) {
		p.addError("Diharapkan keyword 'dalam' setelah variabel iterasi")
		return nil
	}

	p.nextToken() // move to iterable expression
	iterable := p.parseExpression(LOWEST)

	if p.peekTokenIs(COLON) {
		p.nextToken()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &uast.ForEachStatement{
		NodeType: uast.NodeTypeForEachStatement,
		Metadata: meta,
		Item:     itemName,
		Iterable: iterable,
		Body:     body,
	}
}

func (p *Parser) parseFunctionDeclaration() *uast.FunctionDeclaration {
	meta := p.currentMetadata()
	p.nextToken() // eat 'fungsi' or 'def'

	if !p.curTokenIs(IDENT) {
		p.addError("Diharapkan nama fungsi setelah keyword 'fungsi'")
		return nil
	}
	funcName := p.curToken.Literal

	var typeParams []uast.TypeParameter
	if p.peekTokenIs(LT) {
		typeParams = p.parseTypeParameters()
	}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	returnType := ""
	if p.peekTokenIs(RETURN_TYPE) || p.peekTokenIs(ARROW) {
		p.nextToken() // eat 'mengembalikan' or '->'
		p.nextToken()
		returnType = p.parseTypeHint()
	}

	if p.peekTokenIs(COLON) {
		p.nextToken()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &uast.FunctionDeclaration{
		NodeType:       uast.NodeTypeFunctionDeclaration,
		Metadata:       meta,
		Name:           funcName,
		TypeParameters: typeParams,
		Parameters:     params,
		Body:           body,
		ReturnType:     returnType,
	}
}

func (p *Parser) parseFunctionParameters() []uast.Parameter {
	params := []uast.Parameter{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := p.parseParameter()
	params = append(params, param)

	for p.peekTokenIs(COMMA) {
		p.nextToken() // eat ','
		p.nextToken() // move to next param
		param = p.parseParameter()
		params = append(params, param)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseParameter() uast.Parameter {
	isRef := false
	if p.curTokenIs(REF_DECL) {
		isRef = true
		p.nextToken() // eat 'rujuk' / 'ref'
	}
	name := p.curToken.Literal
	typeHint := ""
	var defaultVal uast.Expression

	if p.peekTokenIs(TYPE_ASSIGN) || p.peekTokenIs(COLON) {
		p.nextToken() // eat 'itu' or ':'
		p.nextToken()
		typeHint = p.parseTypeHint()
	}

	if p.peekTokenIs(ASSIGN) {
		p.nextToken() // eat '='
		p.nextToken()
		defaultVal = p.parseExpression(LOWEST)
	}

	return uast.Parameter{
		Name:         name,
		TypeHint:     typeHint,
		DefaultValue: defaultVal,
		IsReference:  isRef,
	}
}

func (p *Parser) parseTypeHint() string {
	var parts []string
	parts = append(parts, p.curToken.Literal)

	if p.peekTokenIs(LT) {
		p.nextToken() // eat '<'
		parts = append(parts, "<")
		p.nextToken() // move to inner type
		parts = append(parts, p.curToken.Literal)
		for p.peekTokenIs(COMMA) {
			p.nextToken()
			parts = append(parts, ", ")
			p.nextToken()
			parts = append(parts, p.curToken.Literal)
		}
		if p.expectPeek(GT) {
			parts = append(parts, ">")
		}
	}

	// Support union types: "Angka|Pecahan" or "Angka atau Pecahan"
	for p.peekTokenIs(PIPE) || p.peekTokenIs(OR) {
		p.nextToken() // eat '|' or 'atau'
		parts = append(parts, "|")
		p.nextToken() // move to next type
		parts = append(parts, p.curToken.Literal)
	}

	return strings.Join(parts, "")
}

func (p *Parser) parseReturnStatement() *uast.ReturnStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'kembalikan' or 'return'

	var val uast.Expression
	if !p.curTokenIs(SEMICOLON) && !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		val = p.parseExpression(LOWEST)
	}

	return &uast.ReturnStatement{
		NodeType: uast.NodeTypeReturnStatement,
		Metadata: meta,
		Value:    val,
	}
}

func (p *Parser) parseConcurrencyStatement(mode string) *uast.ConcurrencyStatement {
	meta := p.currentMetadata()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	block := p.parseBlockStatement()

	return &uast.ConcurrencyStatement{
		NodeType: uast.NodeTypeConcurrencyStatement,
		Metadata: meta,
		Mode:     mode,
		Tasks:    block.Statements,
	}
}

func (p *Parser) parseBlockStatement() *uast.BlockStatement {
	meta := p.currentMetadata()
	block := &uast.BlockStatement{
		NodeType:   uast.NodeTypeBlockStatement,
		Metadata:   meta,
		Statements: []uast.Statement{},
	}

	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(SEMICOLON) {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionStatement() uast.Statement {
	meta := p.currentMetadata()
	expr := p.parseExpression(LOWEST)

	if p.peekTokenIs(ASSIGN) || p.peekTokenIs(PLUS_ASSIGN) || p.peekTokenIs(MINUS_ASSIGN) {
		p.nextToken() // move to operator (=, +=, -=)
		operator := p.curToken.Literal

		p.nextToken() // move to right-hand value
		val := p.parseExpression(LOWEST)

		// Check if target is an ArrayLiteral -> convert to DestructureAssignment!
		if arrLit, ok := expr.(*uast.ArrayLiteral); ok && operator == "=" {
			var bindings []uast.DestructureBinding
			for _, el := range arrLit.Elements {
				if id, ok := el.(*uast.Identifier); ok {
					bindings = append(bindings, uast.DestructureBinding{Target: id.Name})
				}
			}
			return &uast.DestructureAssignment{
				NodeType: uast.NodeTypeDestructureAssignment,
				Metadata: meta,
				Kind:     uast.DestructureArray,
				Bindings: bindings,
				Value:    val,
			}
		}

		// Check if target is a MapLiteral -> convert to DestructureAssignment!
		if mapLit, ok := expr.(*uast.MapLiteral); ok && operator == "=" {
			var bindings []uast.DestructureBinding
			for _, prop := range mapLit.Properties {
				keyStr := ""
				if id, ok := prop.Key.(*uast.Identifier); ok {
					keyStr = id.Name
				} else if lit, ok := prop.Key.(*uast.Literal); ok {
					keyStr = fmt.Sprintf("%v", lit.Value)
				}
				target := keyStr
				if id, ok := prop.Value.(*uast.Identifier); ok {
					target = id.Name
				}
				bindings = append(bindings, uast.DestructureBinding{SourceKey: keyStr, Target: target})
			}
			return &uast.DestructureAssignment{
				NodeType: uast.NodeTypeDestructureAssignment,
				Metadata: meta,
				Kind:     uast.DestructureMap,
				Bindings: bindings,
				Value:    val,
			}
		}

		return &uast.AssignmentStatement{
			NodeType: uast.NodeTypeAssignmentStatement,
			Metadata: meta,
			Target:   expr,
			Operator: operator,
			Value:    val,
		}
	}

	return &uast.ExpressionStatement{
		NodeType:   uast.NodeTypeExpressionStatement,
		Metadata:   meta,
		Expression: expr,
	}
}

func (p *Parser) parseExpression(precedence int) uast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("Tidak ada fungsi parse untuk token '%s' (%s)", p.curToken.Type, p.curToken.Literal))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(SEMICOLON) && precedence < p.peekPrecedence() {
		if p.peekToken.Line > p.curToken.Line && p.parenDepth == 0 {
			if !p.peekTokenIs(PIPE_ARROW) && !p.peekTokenIs(DOT) {
				break
			}
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() uast.Expression {
	meta := p.currentMetadata()
	name := p.curToken.Literal

	// 1. Single-param arrow function: x => ...
	if p.peekTokenIs(FAT_ARROW) {
		p.nextToken() // move to '=>'
		p.nextToken() // move past '=>' to body
		body := p.parseArrowFunctionBody()
		return &uast.FunctionExpression{
			NodeType:   uast.NodeTypeFunctionExpression,
			Metadata:   meta,
			Parameters: []uast.Parameter{{Name: name}},
			Body:       body,
		}
	}

	// 2. Generic call: Kotak<Angka>(10)
	if p.isGenericCallAhead() {
		ident := &uast.Identifier{
			NodeType: uast.NodeTypeIdentifier,
			Metadata: meta,
			Name:     name,
		}
		p.nextToken() // move to '<'
		typeArgs := p.parseTypeArguments()
		if !p.expectPeek(LPAREN) {
			return nil
		}
		return p.parseCallArguments(ident, typeArgs)
	}

	return &uast.Identifier{
		NodeType: uast.NodeTypeIdentifier,
		Metadata: meta,
		Name:     name,
	}
}

func (p *Parser) parseIntegerLiteral() uast.Expression {
	meta := p.currentMetadata()
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Gagal mengubah '%s' ke integer", p.curToken.Literal))
		return nil
	}
	return &uast.Literal{
		NodeType: uast.NodeTypeLiteral,
		Metadata: meta,
		Kind:     uast.LiteralKindInt,
		Raw:      p.curToken.Literal,
		Value:    val,
	}
}

func (p *Parser) parseFloatLiteral() uast.Expression {
	meta := p.currentMetadata()
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Gagal mengubah '%s' ke float", p.curToken.Literal))
		return nil
	}
	return &uast.Literal{
		NodeType: uast.NodeTypeLiteral,
		Metadata: meta,
		Kind:     uast.LiteralKindFloat,
		Raw:      p.curToken.Literal,
		Value:    val,
	}
}

func (p *Parser) parseStringLiteral() uast.Expression {
	return p.parseInterpolationOrLiteral(p.curToken.Literal, p.currentMetadata())
}

func (p *Parser) parseInterpolationOrLiteral(rawStr string, meta *uast.NodeMetadata) uast.Expression {
	if !strings.Contains(rawStr, "{") {
		return &uast.Literal{
			NodeType: uast.NodeTypeLiteral,
			Metadata: meta,
			Kind:     uast.LiteralKindString,
			Raw:      strconv.Quote(rawStr),
			Value:    rawStr,
		}
	}

	var parts []uast.Expression
	var currentLit strings.Builder
	i := 0
	hasInterpolation := false

	for i < len(rawStr) {
		// Handle escape sequences: \{ or \${
		if rawStr[i] == '\\' && i+1 < len(rawStr) {
			if rawStr[i+1] == '{' {
				currentLit.WriteByte('{')
				i += 2
				continue
			}
			if rawStr[i+1] == '$' && i+2 < len(rawStr) && rawStr[i+2] == '{' {
				currentLit.WriteString("${")
				i += 3
				continue
			}
			currentLit.WriteByte('\\')
			currentLit.WriteByte(rawStr[i+1])
			i += 2
			continue
		}

		// Check for ${expr} or {expr}
		isJSInterp := rawStr[i] == '$' && i+1 < len(rawStr) && rawStr[i+1] == '{'
		isIDInterp := rawStr[i] == '{'

		if isJSInterp || isIDInterp {
			startExpr := i + 1
			if isJSInterp {
				startExpr = i + 2
			}

			// Find matching '}'
			braceDepth := 1
			endExpr := -1
			inSubQuote := byte(0)

			for j := startExpr; j < len(rawStr); j++ {
				if inSubQuote != 0 {
					if rawStr[j] == '\\' {
						j++
					} else if rawStr[j] == inSubQuote {
						inSubQuote = 0
					}
					continue
				}
				if rawStr[j] == '"' || rawStr[j] == '\'' || rawStr[j] == '`' {
					inSubQuote = rawStr[j]
					continue
				}
				if rawStr[j] == '{' {
					braceDepth++
				} else if rawStr[j] == '}' {
					braceDepth--
					if braceDepth == 0 {
						endExpr = j
						break
					}
				}
			}

			if endExpr != -1 {
				exprStr := strings.TrimSpace(rawStr[startExpr:endExpr])
				subLex := NewLexer(exprStr)
				subParser := NewParser(subLex)
				subExpr := subParser.parseExpression(LOWEST)

				if subExpr != nil && len(subParser.errors) == 0 && subParser.peekToken.Type == EOF {
					hasInterpolation = true
					if currentLit.Len() > 0 {
						parts = append(parts, &uast.Literal{
							NodeType: uast.NodeTypeLiteral,
							Metadata: meta,
							Kind:     uast.LiteralKindString,
							Raw:      strconv.Quote(currentLit.String()),
							Value:    currentLit.String(),
						})
						currentLit.Reset()
					}
					parts = append(parts, subExpr)
					i = endExpr + 1
					continue
				}
			}
		}

		currentLit.WriteByte(rawStr[i])
		i++
	}

	if hasInterpolation {
		if currentLit.Len() > 0 {
			parts = append(parts, &uast.Literal{
				NodeType: uast.NodeTypeLiteral,
				Metadata: meta,
				Kind:     uast.LiteralKindString,
				Raw:      strconv.Quote(currentLit.String()),
				Value:    currentLit.String(),
			})
		}
		return &uast.TemplateLiteral{
			NodeType: uast.NodeTypeTemplateLiteral,
			Metadata: meta,
			Parts:    parts,
		}
	}

	return &uast.Literal{
		NodeType: uast.NodeTypeLiteral,
		Metadata: meta,
		Kind:     uast.LiteralKindString,
		Raw:      strconv.Quote(currentLit.String()),
		Value:    currentLit.String(),
	}
}

func (p *Parser) parseBooleanLiteral() uast.Expression {
	val := p.curToken.Type == TRUE || p.curToken.Literal == "benar" || p.curToken.Literal == "true" || p.curToken.Literal == "True"
	return &uast.Literal{
		NodeType: uast.NodeTypeLiteral,
		Metadata: p.currentMetadata(),
		Kind:     uast.LiteralKindBool,
		Raw:      strconv.FormatBool(val),
		Value:    val,
	}
}

func (p *Parser) parseNullLiteral() uast.Expression {
	return &uast.Literal{
		NodeType: uast.NodeTypeLiteral,
		Metadata: p.currentMetadata(),
		Kind:     uast.LiteralKindNull,
		Raw:      "null",
		Value:    nil,
	}
}

func (p *Parser) parsePrefixExpression() uast.Expression {
	meta := p.currentMetadata()
	operator := p.curToken.Literal
	p.nextToken()
	arg := p.parseExpression(PREFIX)

	return &uast.UnaryExpression{
		NodeType: uast.NodeTypeUnaryExpression,
		Metadata: meta,
		Operator: operator,
		Argument: arg,
		IsPrefix: true,
	}
}

func (p *Parser) parseInfixExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	operator := p.curToken.Literal
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	return &uast.BinaryExpression{
		NodeType: uast.NodeTypeBinaryExpression,
		Metadata: meta,
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func (p *Parser) parsePipelineExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // eat '|>'
	right := p.parseExpression(PIPELINE_PREC)
	if right == nil {
		p.addError("Diharapkan ekspresi fungsi setelah operator pipeline '|>'")
		return nil
	}

	return &uast.PipelineExpression{
		NodeType: uast.NodeTypePipelineExpression,
		Metadata: meta,
		Left:     left,
		Right:    right,
	}
}

func (p *Parser) parseFunctionExpression() uast.Expression {
	meta := p.currentMetadata()
	if p.peekTokenIs(IDENT) {
		p.nextToken() // optional identifier for named function expression
	}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	returnType := ""
	if p.peekTokenIs(RETURN_TYPE) || p.peekTokenIs(ARROW) {
		p.nextToken() // eat 'mengembalikan' or '->'
		p.nextToken()
		returnType = p.parseTypeHint()
	}

	if p.peekTokenIs(COLON) {
		p.nextToken()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &uast.FunctionExpression{
		NodeType:   uast.NodeTypeFunctionExpression,
		Metadata:   meta,
		Parameters: params,
		Body:       body,
		ReturnType: returnType,
	}
}


func (p *Parser) parseGroupedExpression() uast.Expression {
	if p.isArrowAheadFromParen() {
		meta := p.currentMetadata()
		params := p.parseFunctionParameters()
		if !p.expectPeek(FAT_ARROW) {
			return nil
		}
		p.nextToken() // move past '=>' to body
		body := p.parseArrowFunctionBody()
		return &uast.FunctionExpression{
			NodeType:   uast.NodeTypeFunctionExpression,
			Metadata:   meta,
			Parameters: params,
			Body:       body,
		}
	}

	p.parenDepth++
	defer func() { p.parenDepth-- }()
	p.nextToken() // eat '('
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() uast.Expression {
	meta := p.currentMetadata()
	elements := []uast.Expression{}

	if p.peekTokenIs(RBRACKET) {
		p.nextToken()
		return &uast.ArrayLiteral{
			NodeType: uast.NodeTypeArrayLiteral,
			Metadata: meta,
			Elements: elements,
		}
	}

	p.nextToken()
	elements = append(elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken() // eat ','
		p.nextToken()
		elements = append(elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return &uast.ArrayLiteral{
		NodeType: uast.NodeTypeArrayLiteral,
		Metadata: meta,
		Elements: elements,
	}
}

func (p *Parser) parseMapLiteral() uast.Expression {
	meta := p.currentMetadata()
	props := []uast.MapProperty{}

	if p.peekTokenIs(RBRACE) {
		p.nextToken()
		return &uast.MapLiteral{
			NodeType:   uast.NodeTypeMapLiteral,
			Metadata:   meta,
			Properties: props,
		}
	}

	p.nextToken()
	key := p.parseExpression(LOWEST)

	var val uast.Expression
	if p.peekTokenIs(COLON) || p.peekTokenIs(TYPE_ASSIGN) {
		p.nextToken()
		p.nextToken() // move to value
		val = p.parseExpression(LOWEST)
	} else if id, ok := key.(*uast.Identifier); ok && (p.peekTokenIs(COMMA) || p.peekTokenIs(RBRACE)) {
		// Shorthand property: { a, b }
		val = id
	} else {
		p.addError("Diharapkan ':' atau 'itu' antara kunci dan nilai map")
		return nil
	}

	props = append(props, uast.MapProperty{Key: key, Value: val})

	for p.peekTokenIs(COMMA) {
		p.nextToken() // eat ','
		if p.peekTokenIs(RBRACE) {
			break // allow trailing comma
		}
		p.nextToken() // move to next key
		nextKey := p.parseExpression(LOWEST)

		var nextVal uast.Expression
		if p.peekTokenIs(COLON) || p.peekTokenIs(TYPE_ASSIGN) {
			p.nextToken()
			p.nextToken()
			nextVal = p.parseExpression(LOWEST)
		} else if id, ok := nextKey.(*uast.Identifier); ok && (p.peekTokenIs(COMMA) || p.peekTokenIs(RBRACE)) {
			nextVal = id
		} else {
			p.addError("Diharapkan ':' atau 'itu' antara kunci dan nilai map")
			return nil
		}

		props = append(props, uast.MapProperty{Key: nextKey, Value: nextVal})
	}

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return &uast.MapLiteral{
		NodeType:   uast.NodeTypeMapLiteral,
		Metadata:   meta,
		Properties: props,
	}
}

func (p *Parser) parseCallArguments(callee uast.Expression, typeArgs []string) uast.Expression {
	meta := p.currentMetadata()
	args := []uast.Expression{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return &uast.CallExpression{
			NodeType:      uast.NodeTypeCallExpression,
			Metadata:      meta,
			Callee:        callee,
			TypeArguments: typeArgs,
			Arguments:     args,
		}
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken() // eat ','
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return &uast.CallExpression{
		NodeType:      uast.NodeTypeCallExpression,
		Metadata:      meta,
		Callee:        callee,
		TypeArguments: typeArgs,
		Arguments:     args,
	}
}

func (p *Parser) parseCallExpression(callee uast.Expression) uast.Expression {
	return p.parseCallArguments(callee, nil)
}

func (p *Parser) parseIndexExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // move to index expression
	index := p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	propStr := ""
	if lit, ok := index.(*uast.Literal); ok {
		if s, ok := lit.Value.(string); ok {
			propStr = s
		} else {
			propStr = lit.Raw
		}
	} else if id, ok := index.(*uast.Identifier); ok {
		propStr = id.Name
	}

	return &uast.MemberExpression{
		NodeType:     uast.NodeTypeMemberExpression,
		Metadata:     meta,
		Object:       left,
		Property:     propStr,
		PropertyExpr: index,
		Computed:     true,
	}
}

func (p *Parser) isPropertyIdentifierToken(t TokenType) bool {
	if t == IDENT {
		return true
	}
	switch t {
	case VAR_DECL, CONST_DECL, TYPE_ASSIGN, IF, ELSE_IF, ELSE, SWITCH, CASE, DEFAULT,
		WHILE, FOR, IN, FUNC, RETURN, ASYNC, PARALLEL, TRUE, FALSE, NULL, AND, OR, NOT,
		RETURN_TYPE, MAKA, SELESAI, VIS_PUBLIC, VIS_MODULE, VIS_PRIVATE, VIS_ALLOWLIST,
		IMPORT, FROM, AS, TRY, CATCH, FINALLY, THROW, STYLE_BLOCK, STRUCT_DECL, PRIVATE,
		REF_DECL, TYPE_DECL, NULL_COALESCE, MAYBE, DEFER, TITIP, JEMPUT, HARUS, SEMBARANG,
		SALURAN, TRANSACTION, ROLLBACK, RECOVER, REACTIVE_DECL, WATCH:
		return true
	default:
		return false
	}
}

func (p *Parser) expectPeekProperty() bool {
	if p.isPropertyIdentifierToken(p.peekToken.Type) {
		p.nextToken()
		return true
	}
	p.peekError(IDENT)
	return false
}

func (p *Parser) parseMemberExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	if !p.expectPeekProperty() {
		return nil
	}
	prop := p.curToken.Literal

	memExpr := &uast.MemberExpression{
		NodeType: uast.NodeTypeMemberExpression,
		Metadata: meta,
		Object:   left,
		Property: prop,
		Computed: false,
	}

	if p.isGenericCallAhead() {
		p.nextToken() // move to '<'
		typeArgs := p.parseTypeArguments()
		if !p.expectPeek(LPAREN) {
			return nil
		}
		return p.parseCallArguments(memExpr, typeArgs)
	}

	return memExpr
}

func (p *Parser) currentMetadata() *uast.NodeMetadata {
	return &uast.NodeMetadata{
		Pos: &uast.Position{
			Line:   p.curToken.Line,
			Column: p.curToken.Column,
		},
	}
}

func (p *Parser) parseAtExpression() uast.Expression {
	meta := p.currentMetadata()
	if !p.expectPeek(IDENT) {
		p.addError("Diharapkan nama properti atau metode setelah '@' (contoh: @saldo atau @hitung())")
		return nil
	}
	prop := p.curToken.Literal

	return &uast.MemberExpression{
		NodeType: uast.NodeTypeMemberExpression,
		Metadata: meta,
		Object: &uast.Identifier{
			NodeType: uast.NodeTypeIdentifier,
			Metadata: meta,
			Name:     "diri",
		},
		Property: prop,
		Computed: false,
	}
}

func (p *Parser) parseCastExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // move to target type token
	targetType := p.parseTypeHint()

	return &uast.CastExpression{
		NodeType:   uast.NodeTypeCastExpression,
		Metadata:   meta,
		Value:      left,
		TargetType: targetType,
	}
}

func (p *Parser) parseDeferStatement() *uast.DeferStatement {
	meta := p.currentMetadata()
	p.nextToken() // eat 'saat_selesai' / 'defer'

	var body uast.Statement
	if p.curTokenIs(LBRACE) {
		body = p.parseBlockStatement()
	} else {
		body = p.parseStatement()
	}

	return &uast.DeferStatement{
		NodeType: uast.NodeTypeDeferStatement,
		Metadata: meta,
		Body:     body,
	}
}

func (p *Parser) parseSafeNavigationExpression() uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // eat 'mungkin'

	val := p.parseExpression(PREFIX)
	if val == nil {
		p.addError("Diharapkan ekspresi properti atau pemanggilan setelah 'mungkin'")
		return nil
	}

	return &uast.SafeNavigationExpression{
		NodeType: uast.NodeTypeSafeNavigationExpression,
		Metadata: meta,
		Value:    val,
	}
}

func (p *Parser) parseNullCoalesceExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // eat 'kalau kosong' / '??'

	right := p.parseExpression(COALESCE_PREC)
	if right == nil {
		p.addError("Diharapkan ekspresi nilai cadangan setelah 'kalau kosong' atau '??'")
		return nil
	}

	return &uast.NullCoalesceExpression{
		NodeType: uast.NodeTypeNullCoalesceExpression,
		Metadata: meta,
		Left:     left,
		Right:    right,
	}
}

func (p *Parser) parsePrefixRangeExpression() uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // eat '..'

	var end uast.Expression
	if !p.curTokenIs(RBRACKET) && !p.curTokenIs(EOF) && !p.curTokenIs(SEMICOLON) && !p.curTokenIs(RPAREN) {
		end = p.parseExpression(RANGE_PREC)
	}

	return &uast.RangeExpression{
		NodeType:  uast.NodeTypeRangeExpression,
		Metadata:  meta,
		Start:     nil,
		End:       end,
		Inclusive: true,
	}
}

func (p *Parser) parseInfixRangeExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()

	var end uast.Expression
	if !p.peekTokenIs(RBRACKET) && !p.peekTokenIs(EOF) && !p.peekTokenIs(SEMICOLON) && !p.peekTokenIs(RPAREN) {
		p.nextToken() // move to expression after '..'
		end = p.parseExpression(RANGE_PREC)
	}

	return &uast.RangeExpression{
		NodeType:  uast.NodeTypeRangeExpression,
		Metadata:  meta,
		Start:     left,
		End:       end,
		Inclusive: true,
	}
}

func (p *Parser) parseOptionalMemberExpression(left uast.Expression) uast.Expression {
	meta := p.currentMetadata()
	if p.peekTokenIs(LBRACKET) {
		p.nextToken() // move to '['
		p.nextToken() // move to index expression
		index := p.parseExpression(LOWEST)
		if !p.expectPeek(RBRACKET) {
			return nil
		}
		propStr := ""
		if lit, ok := index.(*uast.Literal); ok {
			if s, ok := lit.Value.(string); ok {
				propStr = s
			} else {
				propStr = lit.Raw
			}
		} else if id, ok := index.(*uast.Identifier); ok {
			propStr = id.Name
		}
		return &uast.MemberExpression{
			NodeType:     uast.NodeTypeMemberExpression,
			Metadata:     meta,
			Object:       left,
			Property:     propStr,
			PropertyExpr: index,
			Computed:     true,
			Safe:         true,
		}
	}

	if !p.expectPeekProperty() {
		return nil
	}
	prop := p.curToken.Literal

	return &uast.MemberExpression{
		NodeType: uast.NodeTypeMemberExpression,
		Metadata: meta,
		Object:   left,
		Property: prop,
		Computed: false,
		Safe:     true,
	}
}

func (p *Parser) parseArrowFunctionBody() *uast.BlockStatement {
	if p.curTokenIs(LBRACE) {
		return p.parseBlockStatement()
	}
	meta := p.currentMetadata()
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	return &uast.BlockStatement{
		NodeType: uast.NodeTypeBlockStatement,
		Metadata: meta,
		Statements: []uast.Statement{
			&uast.ReturnStatement{
				NodeType: uast.NodeTypeReturnStatement,
				Metadata: meta,
				Value:    expr,
			},
		},
	}
}

func (p *Parser) isArrowAheadFromParen() bool {
	lClone := *p.l
	tok := p.peekToken // the token after '('
	depth := 1
	for tok.Type != EOF {
		if tok.Type == LPAREN {
			depth++
		} else if tok.Type == RPAREN {
			depth--
			if depth == 0 {
				nextTok := lClone.NextToken()
				return nextTok.Type == FAT_ARROW
			}
		}
		tok = lClone.NextToken()
	}
	return false
}

func (p *Parser) isGenericCallAhead() bool {
	if !p.peekTokenIs(LT) {
		return false
	}
	lClone := *p.l
	tok := lClone.NextToken() // token after '<'
	hasTypeArg := false
	for tok.Type != EOF {
		if tok.Type == IDENT || tok.Type == SEMBARANG || tok.Type == COMMA {
			if tok.Type == IDENT || tok.Type == SEMBARANG {
				hasTypeArg = true
			}
		} else if tok.Type == GT {
			if !hasTypeArg {
				return false
			}
			nextTok := lClone.NextToken()
			return nextTok.Type == LPAREN
		} else {
			return false
		}
		tok = lClone.NextToken()
	}
	return false
}

func (p *Parser) parseTypeArguments() []string {
	// curToken is '<'
	typeArgs := []string{}
	p.nextToken() // move to first type
	for !p.curTokenIs(GT) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) || p.curTokenIs(SEMBARANG) {
			typeArgs = append(typeArgs, p.curToken.Literal)
		}
		if p.peekTokenIs(COMMA) {
			p.nextToken() // eat ','
		}
		p.nextToken()
	}
	return typeArgs
}

func (p *Parser) parseTypeParameters() []uast.TypeParameter {
	p.nextToken() // eat '<'
	typeParams := []uast.TypeParameter{}
	p.nextToken() // move to first type parameter

	for !p.curTokenIs(GT) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) || p.curTokenIs(SEMBARANG) {
			name := p.curToken.Literal
			constraint := ""
			if p.peekTokenIs(COLON) || p.peekTokenIs(HARUS) || p.peekTokenIs(TYPE_ASSIGN) {
				p.nextToken() // eat ':', 'harus', or 'itu'
				if p.peekTokenIs(IDENT) || p.peekTokenIs(SEMBARANG) {
					p.nextToken()
					constraint = p.curToken.Literal
				}
			}
			typeParams = append(typeParams, uast.TypeParameter{
				Name:       name,
				Constraint: constraint,
			})
		}
		if p.peekTokenIs(COMMA) {
			p.nextToken() // eat ','
		}
		p.nextToken()
	}
	return typeParams
}

func (p *Parser) parseTitipExpression() uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // move to argument expression
	arg := p.parseExpression(PREFIX)
	return &uast.TitipExpression{
		NodeType: uast.NodeTypeTitipExpression,
		Metadata: meta,
		Argument: arg,
	}
}

func (p *Parser) parseJemputExpression() uast.Expression {
	meta := p.currentMetadata()
	p.nextToken() // move to argument expression
	arg := p.parseExpression(PREFIX)
	var timeout uast.Expression
	if p.peekTokenIs(IN) {
		p.nextToken() // eat 'dalam'
		p.nextToken() // move to timeout expr
		timeout = p.parseExpression(PREFIX)
	}
	return &uast.JemputExpression{
		NodeType: uast.NodeTypeJemputExpression,
		Metadata: meta,
		Argument: arg,
		Timeout:  timeout,
	}
}

