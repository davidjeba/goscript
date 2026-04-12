package gscompiler

import "fmt"

// ParseError represents a parsing error with its source position.
type ParseError struct {
	Pos Position
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Pos.String(), e.Msg)
}

// Parser is a recursive descent parser for Go-syntax source code.
// It produces a full AST from a stream of tokens.
type Parser struct {
	tokens []Token
	pos    int
	errors []ParseError
}

// NewParser creates a new Parser for the given token stream.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: make([]ParseError, 0),
	}
}

// Parse parses the entire token stream and returns a Program AST.
// Returns an error if any parse errors were encountered.
func (p *Parser) Parse() (*Program, error) {
	prog := &Program{}

	// Parse package declaration
	prog.Package = p.parsePackageDecl()

	// Parse imports
	prog.Imports = p.parseImports()

	// Parse top-level declarations
	for !p.isAtEnd() {
		decl := p.parseTopLevelDecl()
		if decl != nil {
			prog.Decls = append(prog.Decls, decl)
		}
	}

	if len(p.errors) > 0 {
		return prog, p.errors[0]
	}
	return prog, nil
}

// --- Token inspection helpers ---

// peek returns the current token without advancing.
func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF, Line: 0, Column: 0}
	}
	return p.tokens[p.pos]
}

// peekAt returns the token at offset from current position.
func (p *Parser) peekAt(offset int) Token {
	idx := p.pos + offset
	if idx >= len(p.tokens) {
		return Token{Type: TOKEN_EOF, Line: 0, Column: 0}
	}
	return p.tokens[idx]
}

// advance consumes the current token and returns it.
func (p *Parser) advance() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF, Line: 0, Column: 0}
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

// isAtEnd returns true if the parser has consumed all tokens.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TOKEN_EOF
}

// check returns true if the current token matches the given type.
func (p *Parser) check(tt TokenType) bool {
	return p.peek().Type == tt
}

// match advances and returns true if the current token matches the given type.
func (p *Parser) match(tt TokenType) bool {
	if p.check(tt) {
		p.advance()
		return true
	}
	return false
}

// expect consumes the current token if it matches, or records an error.
func (p *Parser) expect(tt TokenType, msg string) Token {
	if p.check(tt) {
		return p.advance()
	}
	tok := p.peek()
	p.error(tok, msg)
	return Token{Type: tt, Line: tok.Line, Column: tok.Column}
}

// expectSemicolon consumes a semicolon token if present, or records an error.
// It's lenient: if the next token is on a new line, the semicolon is optional.
func (p *Parser) expectSemicolon() {
	if !p.match(TOKEN_SEMICOLON) {
		// Semicolon was auto-inserted, that's fine
	}
}

// pos returns the Position of the current token.
func (p *Parser) pos() Position {
	tok := p.peek()
	return Position{Line: tok.Line, Column: tok.Column}
}

// error records a parse error.
func (p *Parser) error(tok Token, msg string) {
	p.errors = append(p.errors, ParseError{
		Pos: Position{Line: tok.Line, Column: tok.Column},
		Msg: msg,
	})
}

// isTypeStart returns true if the current token could start a type expression.
func (p *Parser) isTypeStart() bool {
	switch p.peek().Type {
	case TOKEN_IDENT:
		return true
	case TOKEN_STAR:
		return true
	case TOKEN_LPAREN:
		return true
	case TOKEN_LBRACK:
		return true
	case TOKEN_STRUCT:
		return true
	case TOKEN_INTERFACE:
		return true
	case TOKEN_FUNC:
		return true
	case TOKEN_MAP:
		return true
	case TOKEN_CHAN:
		return true
	case TOKEN_ARROW:
		return true
	default:
		return false
	}
}

// --- Package and imports ---

// parsePackageDecl parses: package <name>
func (p *Parser) parsePackageDecl() *Ident {
	p.expect(TOKEN_PACKAGE, "expected 'package'")
	tok := p.expect(TOKEN_IDENT, "expected package name")
	p.expectSemicolon()
	return &Ident{
		PosNode: Position{Line: tok.Line, Column: tok.Column},
		Name:    tok.Literal,
	}
}

// parseImports parses one or more import declarations.
func (p *Parser) parseImports() []*ImportSpec {
	var imports []*ImportSpec
	for p.check(TOKEN_IMPORT) {
		p.advance() // consume 'import'
		if p.match(TOKEN_LPAREN) {
			// Grouped imports: import ( ... )
			for !p.check(TOKEN_RPAREN) && !p.isAtEnd() {
				spec := p.parseImportSpec()
				if spec != nil {
					imports = append(imports, spec)
				}
				p.expectSemicolon()
			}
			p.expect(TOKEN_RPAREN, "expected ')' in import group")
			p.expectSemicolon()
		} else {
			// Single import: import "path"
			spec := p.parseImportSpec()
			if spec != nil {
				imports = append(imports, spec)
			}
			p.expectSemicolon()
		}
	}
	return imports
}

// parseImportSpec parses: [alias] "path"
func (p *Parser) parseImportSpec() *ImportSpec {
	pos := p.pos()
	var name *Ident
	if p.check(TOKEN_IDENT) && p.peekAt(1).Type == TOKEN_STRING {
		tok := p.advance()
		name = &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: tok.Literal}
	} else if p.check(TOKEN_DOT) {
		tok := p.advance()
		name = &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: "."}
	} else if p.check(TOKEN_NOT) {
		p.advance() // consume _
		name = &Ident{PosNode: pos, Name: "_"}
	}
	tok := p.expect(TOKEN_STRING, "expected import path")
	return &ImportSpec{
		PosNode: pos,
		Name:    name,
		Path:    &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_STRING, Value: tok.Literal},
	}
}

// --- Top-level declarations ---

// parseTopLevelDecl parses a top-level declaration (func, type, var, const).
func (p *Parser) parseTopLevelDecl() Decl {
	switch p.peek().Type {
	case TOKEN_FUNC:
		return p.parseFuncDecl()
	case TOKEN_TYPE:
		return p.parseGenDecl(TOKEN_TYPE)
	case TOKEN_VAR:
		return p.parseGenDecl(TOKEN_VAR)
	case TOKEN_CONST:
		return p.parseGenDecl(TOKEN_CONST)
	default:
		p.error(p.peek(), fmt.Sprintf("unexpected token %s, expected declaration", TokenString(p.peek().Type)))
		p.advance()
		return nil
	}
}

// --- Function declarations ---

// parseFuncDecl parses: func [receiver] Name(params) [results] [body]
func (p *Parser) parseFuncDecl() *FuncDecl {
	pos := p.pos()
	p.expect(TOKEN_FUNC, "expected 'func'")

	decl := &FuncDecl{PosNode: pos}

	// Check for receiver: func (x *T) Name(...)
	if p.check(TOKEN_LPAREN) {
		// This could be a receiver or could be the start of the name with a type.
		// A receiver is followed by ')' and then an identifier.
		// Save position and try to determine if this is a receiver.
		if p.isReceiver() {
			decl.Recv = p.parseFieldList()
		}
	}

	// Function name
	tok := p.expect(TOKEN_IDENT, "expected function name")
	decl.Name = &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: tok.Literal}

	// Function type (params and results)
	decl.Type = p.parseFuncType()

	// Body
	if p.match(TOKEN_LBRACE) {
		p.pos--
		decl.Body = p.parseBlockStmt()
	}

	return decl
}

// isReceiver checks if the current position indicates a method receiver.
func (p *Parser) isReceiver() bool {
	// Save position
	savePos := p.pos

	// Try to match: ( [ident] [*] Type )
	if !p.match(TOKEN_LPAREN) {
		p.pos = savePos
		return false
	}

	// Skip optional identifier name
	if p.check(TOKEN_IDENT) {
		p.advance()
	}

	// Skip optional *
	if p.check(TOKEN_STAR) {
		p.advance()
	}

	// Expect a type name
	if !p.check(TOKEN_IDENT) {
		p.pos = savePos
		return false
	}
	p.advance()

	// Expect closing paren
	hasRParen := p.check(TOKEN_RPAREN)
	if !hasRParen {
		// Maybe there's more type like pkg.Type
		if p.check(TOKEN_DOT) {
			p.advance()
			if p.check(TOKEN_IDENT) {
				p.advance()
				hasRParen = p.check(TOKEN_RPAREN)
			}
		}
	}

	p.pos = savePos
	return hasRParen
}

// parseFuncType parses: (params) [(results)]
func (p *Parser) parseFuncType() *FuncType {
	pos := p.pos()
	ft := &FuncType{
		PosNode: pos,
		Params:  p.parseFieldList(),
	}

	// Check for results
	if p.check(TOKEN_LPAREN) {
		ft.Results = p.parseFieldList()
	} else if p.isTypeStart() {
		// Single unnamed result
		resultType := p.parseType()
		ft.Results = &FieldList{
			PosNode: pos,
			List:    []*Field{{Type: resultType}},
		}
	}

	return ft
}

// parseFieldList parses: ( field1, field2, ... )
// Also handles empty parameter lists.
func (p *Parser) parseFieldList() *FieldList {
	pos := p.pos()
	p.expect(TOKEN_LPAREN, "expected '('")

	fl := &FieldList{PosNode: pos}

	if p.check(TOKEN_RPAREN) {
		p.advance()
		return fl
	}

	fl.List = p.parseFieldSequence()
	p.expect(TOKEN_RPAREN, "expected ')'")
	return fl
}

// parseFieldSequence parses a comma-separated list of fields until ')' or other terminator.
func (p *Parser) parseFieldSequence() []*Field {
	var fields []*Field

	for !p.check(TOKEN_RPAREN) && !p.isAtEnd() {
		field := p.parseField()
		fields = append(fields, field)
		if !p.match(TOKEN_COMMA) {
			break
		}
	}

	return fields
}

// parseField parses a single field in a parameter list, struct, or interface.
// Field can be: Name Type, or just Type (unnamed).
func (p *Parser) parseField() *Field {
	pos := p.pos()

	// Handle variadic: ...Type
	if p.check(TOKEN_ELLIPSIS) {
		p.advance()
		elt := p.parseType()
		return &Field{PosNode: pos, Type: &Ellipsis{PosNode: pos, Elt: elt}}
	}

	// Try to parse: Name Type, Name1, Name2 Type, or just Type
	if p.check(TOKEN_IDENT) && !p.isFuncParamEnd() {
		// Collect names
		names := []*Ident{}
		for p.check(TOKEN_IDENT) {
			tok := p.advance()
			names = append(names, &Ident{
				PosNode: Position{Line: tok.Line, Column: tok.Column},
				Name:    tok.Literal,
			})
			if !p.match(TOKEN_COMMA) {
				break
			}
		}

		if p.isTypeStart() {
			typ := p.parseType()
			return &Field{PosNode: pos, Names: names, Type: typ}
		}

		// If no type follows, this was just a type name
		if len(names) == 1 {
			return &Field{PosNode: pos, Type: names[0]}
		}
	}

	// Just a type (unnamed parameter or embedded field)
	typ := p.parseType()
	return &Field{PosNode: pos, Type: typ}
}

// isFuncParamEnd returns true if the current token ends a parameter list item.
func (p *Parser) isFuncParamEnd() bool {
	tt := p.peek().Type
	if tt == TOKEN_COMMA || tt == TOKEN_RPAREN || tt == TOKEN_ELLIPSIS {
		return true
	}
	// Check if next-next token suggests this is just a type
	// e.g., "string" followed by "," or ")" is a type-only field
	return false
}

// --- General declarations (type, var, const) ---

// parseGenDecl parses: var/const/type spec1, spec2;
func (p *Parser) parseGenDecl(tokType TokenType) *GenDecl {
	pos := p.pos()
	p.advance() // consume var/const/type

	decl := &GenDecl{PosNode: pos, Tok: tokType}

	if p.match(TOKEN_LPAREN) {
		// Grouped declaration: var ( ... )
		for !p.check(TOKEN_RPAREN) && !p.isAtEnd() {
			spec := p.parseSpec(tokType)
			decl.Specs = append(decl.Specs, spec)
			p.expectSemicolon()
		}
		p.expect(TOKEN_RPAREN, "expected ')'")
		p.expectSemicolon()
	} else {
		spec := p.parseSpec(tokType)
		decl.Specs = append(decl.Specs, spec)
		p.expectSemicolon()
	}

	return decl
}

// parseSpec parses a single specification based on the declaration token type.
func (p *Parser) parseSpec(tokType TokenType) Spec {
	switch tokType {
	case TOKEN_TYPE:
		return p.parseTypeSpec()
	case TOKEN_VAR, TOKEN_CONST:
		return p.parseValueSpec()
	default:
		return nil
	}
}

// parseTypeSpec parses: Name Type
func (p *Parser) parseTypeSpec() *TypeSpec {
	pos := p.pos()
	nameTok := p.expect(TOKEN_IDENT, "expected type name")
	name := &Ident{PosNode: Position{Line: nameTok.Line, Column: nameTok.Column}, Name: nameTok.Literal}

	typ := p.parseType()

	return &TypeSpec{PosNode: pos, Name: name, Type: typ}
}

// parseValueSpec parses: Name1, Name2 [Type] [= expr1, expr2]
func (p *Parser) parseValueSpec() *ValueSpec {
	pos := p.pos()
	spec := &ValueSpec{PosNode: pos}

	// Parse names
	for {
		tok := p.expect(TOKEN_IDENT, "expected identifier")
		spec.Names = append(spec.Names, &Ident{
			PosNode: Position{Line: tok.Line, Column: tok.Column},
			Name:    tok.Literal,
		})
		if !p.match(TOKEN_COMMA) {
			break
		}
	}

	// Optional type
	if p.isTypeStart() && !p.check(TOKEN_ASSIGN) {
		spec.Type = p.parseType()
	}

	// Optional value
	if p.match(TOKEN_ASSIGN) {
		for {
			expr := p.parseExpr()
			spec.Values = append(spec.Values, expr)
			if !p.match(TOKEN_COMMA) {
				break
			}
		}
	}

	return spec
}

// --- Type parsing ---

// parseType parses a Go type expression.
func (p *Parser) parseType() Expr {
	pos := p.pos()

	switch p.peek().Type {
	case TOKEN_STAR:
		// Pointer type: *T
		p.advance()
		x := p.parseType()
		return &StarExpr{PosNode: pos, X: x}

	case TOKEN_LBRACK:
		// Array, slice, or map type: [N]T, []T, map[K]V
		if p.peekAt(1).Type == TOKEN_RBRACK {
			// Slice type: []T
			p.advance() // consume [
			p.advance() // consume ]
			elt := p.parseType()
			return &SliceType{PosNode: pos, Elt: elt}
		}
		// Array type: [N]T
		p.advance() // consume [
		var length Expr
		if p.check(TOKEN_ELLIPSIS) {
			p.advance()
			length = &Ellipsis{PosNode: p.pos()}
		} else {
			length = p.parseExpr()
		}
		p.expect(TOKEN_RBRACK, "expected ']' in array type")
		elt := p.parseType()
		return &ArrayType{PosNode: pos, Len: length, Elt: elt}

	case TOKEN_MAP:
		// Map type: map[K]V
		p.advance() // consume map
		p.expect(TOKEN_LBRACK, "expected '[' in map type")
		key := p.parseType()
		p.expect(TOKEN_RBRACK, "expected ']' in map type")
		value := p.parseType()
		return &MapType{PosNode: pos, Key: key, Value: value}

	case TOKEN_CHAN:
		// Channel type: chan T
		p.advance()
		value := p.parseType()
		return &ChanType{PosNode: pos, Dir: TOKEN_CHAN, Value: value}

	case TOKEN_ARROW:
		// Receive-only channel: <-chan T
		p.advance()
		if p.check(TOKEN_CHAN) {
			p.advance()
			value := p.parseType()
			return &ChanType{PosNode: pos, Dir: TOKEN_ARROW, Value: value}
		}
		p.error(p.peek(), "expected 'chan' after '<-'")
		return &Ident{PosNode: pos, Name: "unknown"}

	case TOKEN_FUNC:
		// Function type: func(params) results
		p.advance()
		return p.parseFuncType()

	case TOKEN_INTERFACE:
		// Interface type: interface { methods }
		return p.parseInterfaceType()

	case TOKEN_STRUCT:
		// Struct type: struct { fields }
		return p.parseStructType()

	case TOKEN_IDENT:
		// Named type or pkg.Type
		tok := p.advance()
		ident := &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: tok.Literal}
		if p.check(TOKEN_DOT) {
			p.advance() // consume .
			sel := p.expect(TOKEN_IDENT, "expected identifier after '.'")
			return &SelectorExpr{
				PosNode: pos,
				X:       ident,
				Sel:     &Ident{PosNode: Position{Line: sel.Line, Column: sel.Column}, Name: sel.Literal},
			}
		}
		return ident

	case TOKEN_LPAREN:
		// Parenthesized type: (T)
		p.advance()
		typ := p.parseType()
		p.expect(TOKEN_RPAREN, "expected ')' in type expression")
		return typ

	default:
		p.error(p.peek(), fmt.Sprintf("expected type, got %s", TokenString(p.peek().Type)))
		return &Ident{PosNode: pos, Name: "unknown"}
	}
}

// parseStructType parses: struct { fields }
func (p *Parser) parseStructType() *StructType {
	pos := p.pos()
	p.expect(TOKEN_STRUCT, "expected 'struct'")

	fields := p.parseStructFields()

	return &StructType{PosNode: pos, Fields: fields}
}

// parseStructFields parses the field list inside a struct definition.
func (p *Parser) parseStructFields() *FieldList {
	pos := p.pos()
	p.expect(TOKEN_LBRACE, "expected '{' in struct type")

	fl := &FieldList{PosNode: pos}
	if p.check(TOKEN_RBRACE) {
		p.advance()
		return fl
	}

	for !p.check(TOKEN_RBRACE) && !p.isAtEnd() {
		field := p.parseStructField()
		fl.List = append(fl.List, field)
		p.expectSemicolon()
	}

	p.expect(TOKEN_RBRACE, "expected '}' in struct type")
	return fl
}

// parseStructField parses a struct field: Name Type [Tag], or embedded Type.
func (p *Parser) parseStructField() *Field {
	pos := p.pos()

	// Collect names
	var names []*Ident
	for p.check(TOKEN_IDENT) {
		tok := p.advance()
		names = append(names, &Ident{
			PosNode: Position{Line: tok.Line, Column: tok.Column},
			Name:    tok.Literal,
		})
		if !p.match(TOKEN_COMMA) {
			break
		}
	}

	if p.isTypeStart() {
		typ := p.parseType()
		var tag *BasicLit
		if p.check(TOKEN_STRING) || p.check(TOKEN_RAW_STRING) {
			tagTok := p.advance()
			tag = &BasicLit{
				PosNode: Position{Line: tagTok.Line, Column: tagTok.Column},
				Kind:    tagTok.Type,
				Value:   tagTok.Literal,
			}
		}
		return &Field{PosNode: pos, Names: names, Type: typ, Tag: tag}
	}

	if len(names) == 1 {
		return &Field{PosNode: pos, Type: names[0]}
	}

	p.error(p.peek(), "expected type in struct field")
	return &Field{PosNode: pos}
}

// parseInterfaceType parses: interface { methods }
func (p *Parser) parseInterfaceType() *InterfaceType {
	pos := p.pos()
	p.expect(TOKEN_INTERFACE, "expected 'interface'")

	methods := &FieldList{PosNode: pos}
	p.expect(TOKEN_LBRACE, "expected '{' in interface type")

	for !p.check(TOKEN_RBRACE) && !p.isAtEnd() {
		method := p.parseInterfaceMethod()
		methods.List = append(methods.List, method)
		p.expectSemicolon()
	}

	p.expect(TOKEN_RBRACE, "expected '}' in interface type")
	return &InterfaceType{PosNode: pos, Methods: methods}
}

// parseInterfaceMethod parses a method specification in an interface.
// Can be: Name(params) results, or Type (embedded interface).
func (p *Parser) parseInterfaceMethod() *Field {
	pos := p.pos()

	if p.check(TOKEN_FUNC) {
		p.advance()
		ft := p.parseFuncType()
		return &Field{PosNode: pos, Type: ft}
	}

	// Name(params) results or embedded Type
	if p.check(TOKEN_IDENT) {
		nameTok := p.advance()
		name := &Ident{PosNode: Position{Line: nameTok.Line, Column: nameTok.Column}, Name: nameTok.Literal}

		if p.check(TOKEN_LPAREN) {
			// Method: Name(params) results
			ft := p.parseFuncType()
			return &Field{PosNode: pos, Names: []*Ident{name}, Type: ft}
		}

		// Embedded type
		return &Field{PosNode: pos, Type: name}
	}

	// Embedded interface name
	typ := p.parseType()
	return &Field{PosNode: pos, Type: typ}
}

// --- Statement parsing ---

// parseBlockStmt parses: { stmt1; stmt2; ... }
func (p *Parser) parseBlockStmt() *BlockStmt {
	pos := p.pos()
	p.expect(TOKEN_LBRACE, "expected '{'")

	block := &BlockStmt{PosNode: pos}

	for !p.check(TOKEN_RBRACE) && !p.isAtEnd() {
		stmt := p.parseStmt()
		if stmt != nil {
			block.List = append(block.List, stmt)
		}
		p.expectSemicolon()
	}

	p.expect(TOKEN_RBRACE, "expected '}'")
	return block
}

// parseStmt parses a single statement.
func (p *Parser) parseStmt() Stmt {
	switch p.peek().Type {
	case TOKEN_IF:
		return p.parseIfStmt()
	case TOKEN_FOR:
		return p.parseForStmt()
	case TOKEN_SWITCH:
		return p.parseSwitchStmt()
	case TOKEN_SELECT:
		return p.parseSelectStmt()
	case TOKEN_RETURN:
		return p.parseReturnStmt()
	case TOKEN_BREAK, TOKEN_CONTINUE, TOKEN_GOTO, TOKEN_FALLTHROUGH:
		return p.parseBranchStmt()
	case TOKEN_DEFER:
		return p.parseDeferStmt()
	case TOKEN_GO:
		return p.parseGoStmt()
	case TOKEN_VAR, TOKEN_CONST:
		decl := p.parseGenDecl(p.peek().Type)
		return &DeclStmt{PosNode: decl.PosNode, Decl: decl}
	case TOKEN_TYPE:
		decl := p.parseGenDecl(TOKEN_TYPE)
		return &DeclStmt{PosNode: decl.PosNode, Decl: decl}
	default:
		// Try to parse an expression statement or assignment
		return p.parseSimpleStmt()
	}
}

// parseSimpleStmt parses expression statements, assignments, and short variable declarations.
func (p *Parser) parseSimpleStmt() Stmt {
	pos := p.pos()

	expr := p.parseExpr()

	// Check for assignment operators
	switch p.peek().Type {
	case TOKEN_ASSIGN, TOKEN_SHORT_DECL,
		TOKEN_PLUS_ASSIGN, TOKEN_MINUS_ASSIGN, TOKEN_STAR_ASSIGN, TOKEN_SLASH_ASSIGN, TOKEN_PERCENT_ASSIGN,
		TOKEN_AMP_ASSIGN, TOKEN_PIPE_ASSIGN, TOKEN_CARET_ASSIGN, TOKEN_LSHIFT_ASSIGN, TOKEN_RSHIFT_ASSIGN:
		tok := p.advance()
		tokType := tok.Type

		// Collect LHS (may be multiple for x, y = ...)
		lhs := []Expr{expr}
		for p.check(TOKEN_COMMA) {
			p.advance()
			lhs = append(lhs, p.parseExpr())
		}

		// Parse RHS
		rhs := p.parseExprList()

		return &AssignStmt{
			PosNode: pos,
			Lhs:     lhs,
			Tok:     tokType,
			Rhs:     rhs,
		}
	}

	// Check for inc/dec: x++, x--
	if p.check(TOKEN_PLUS) && p.peekAt(1).Type == TOKEN_PLUS {
		p.advance() // consume first +
		p.advance() // consume second +
		return &IncDecStmt{PosNode: pos, X: expr, Tok: TOKEN_PLUS}
	}
	if p.check(TOKEN_MINUS) && p.peekAt(1).Type == TOKEN_MINUS {
		p.advance()
		p.advance()
		return &IncDecStmt{PosNode: pos, X: expr, Tok: TOKEN_MINUS}
	}

	// Check for send: ch <- value (expr must be a channel)
	if p.check(TOKEN_ARROW) {
		p.advance()
		value := p.parseExpr()
		return &SendStmt{PosNode: pos, Chan: expr, Value: value}
	}

	return &ExprStmt{PosNode: pos, X: expr}
}

// parseExprList parses a comma-separated list of expressions.
func (p *Parser) parseExprList() []Expr {
	var exprs []Expr
	exprs = append(exprs, p.parseExpr())
	for p.match(TOKEN_COMMA) {
		exprs = append(exprs, p.parseExpr())
	}
	return exprs
}

// --- Control flow ---

// parseIfStmt parses: if [init;] cond { body } [else [if|{...}]]
func (p *Parser) parseIfStmt() *IfStmt {
	pos := p.pos()
	p.expect(TOKEN_IF, "expected 'if'")

	stmt := &IfStmt{PosNode: pos}

	// Parse condition; may have init statement before semicolon
	stmt.Cond = p.parseExpr()

	// Check for semicolon (init statement)
	if p.check(TOKEN_SEMICOLON) {
		p.advance()
		// The condition we just parsed is actually the init statement
		stmt.Init = &ExprStmt{PosNode: stmt.Cond.Pos(), X: stmt.Cond}
		stmt.Cond = p.parseExpr()
	}

	stmt.Body = p.parseBlockStmt()

	// else clause
	if p.match(TOKEN_ELSE) {
		if p.check(TOKEN_IF) {
			stmt.Else = p.parseIfStmt()
		} else if p.check(TOKEN_LBRACE) {
			stmt.Else = p.parseBlockStmt()
		} else {
			p.error(p.peek(), "expected 'if' or '{' after 'else'")
		}
	}

	return stmt
}

// parseForStmt parses all forms of for statements:
// - for init; cond; post { body }
// - for cond { body }
// - for { body } (infinite loop)
// - for range x { body }
// - for key := range x { body }
// - for key, val := range x { body }
func (p *Parser) parseForStmt() Stmt {
	pos := p.pos()
	p.expect(TOKEN_FOR, "expected 'for'")

	// for { body } - infinite loop
	if p.check(TOKEN_LBRACE) {
		body := p.parseBlockStmt()
		return &ForStmt{PosNode: pos, Body: body}
	}

	// Try to detect range statement
	// Save position and try to parse as range
	if p.isRangeStmt(pos) {
		return p.parseRangeStmt(pos)
	}

	// for init; cond; post { body } or for cond { body }
	// We need to determine if this is a 3-clause or simple for loop.
	// Strategy: parse first simple stmt, then check for semicolons.

	initStmt := p.parseSimpleStmt()
	p.expectSemicolon()

	// If we see a semicolon, this is a 3-clause for loop
	condStmt := p.parseSimpleStmt()
	p.expectSemicolon()

	// Check for closing brace (3-clause)
	if p.check(TOKEN_LBRACE) {
		var init, cond Stmt
		var post Stmt
		var body *BlockStmt

		// The first parsed statement could be init or empty
		init = initStmt

		// cond could be empty or an expression
		if exprStmt, ok := condStmt.(*ExprStmt); ok {
			cond = exprStmt
		} else {
			// init was actually cond, cond is actually post
			cond = initStmt
			init = nil
		}

		body = p.parseBlockStmt()
		return &ForStmt{PosNode: pos, Init: init, Cond: cond, Body: body}
	}

	// Handle post clause
	post := condStmt
	p.expectSemicolon()
	body := p.parseBlockStmt()

	return &ForStmt{PosNode: pos, Init: initStmt, Body: body, Post: post}
}

// isRangeStmt checks if the current parse state represents a range statement.
func (p *Parser) isRangeStmt(pos Position) bool {
	savePos := p.pos
	defer func() { p.pos = savePos }()

	// Try to parse: [ident [, ident]] := range expr
	// or: range expr
	if p.check(TOKEN_RANGE) {
		return true
	}

	if p.check(TOKEN_IDENT) {
		p.advance()
		if p.check(TOKEN_SHORT_DECL) || p.check(TOKEN_ASSIGN) {
			p.advance()
			if p.check(TOKEN_RANGE) {
				return true
			}
		} else if p.check(TOKEN_COMMA) {
			p.advance()
			if p.check(TOKEN_IDENT) {
				p.advance()
				if p.check(TOKEN_SHORT_DECL) || p.check(TOKEN_ASSIGN) {
					p.advance()
					if p.check(TOKEN_RANGE) {
						return true
					}
				}
			}
		}
	}

	return false
}

// parseRangeStmt parses: for [key [, val]] := range expr { body }
func (p *Parser) parseRangeStmt(pos Position) *RangeStmt {
	stmt := &RangeStmt{PosNode: pos}

	// Parse key (optional)
	if p.check(TOKEN_IDENT) && !p.check(TOKEN_RANGE) {
		tok := p.advance()
		stmt.Key = &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: tok.Literal}

		// Check for value
		if p.match(TOKEN_COMMA) {
			valTok := p.expect(TOKEN_IDENT, "expected identifier in range")
			stmt.Value = &Ident{PosNode: Position{Line: valTok.Line, Column: valTok.Column}, Name: valTok.Literal}
		}
	}

	if p.match(TOKEN_SHORT_DECL) {
		stmt.Tok = TOKEN_SHORT_DECL
	} else if p.match(TOKEN_ASSIGN) {
		stmt.Tok = TOKEN_ASSIGN
	}

	p.expect(TOKEN_RANGE, "expected 'range'")
	stmt.X = p.parseExpr()
	stmt.Body = p.parseBlockStmt()

	return stmt
}

// parseSwitchStmt parses: switch [init;] tag { cases }
func (p *Parser) parseSwitchStmt() *SwitchStmt {
	pos := p.pos()
	p.expect(TOKEN_SWITCH, "expected 'switch'")

	stmt := &SwitchStmt{PosNode: pos}

	// Check for type switch: switch x.(type) { ... }
	// For now, parse as expression switch

	// Optional init; tag
	if p.check(TOKEN_LBRACE) {
		// No tag
		stmt.Body = p.parseBlockStmt()
		return stmt
	}

	// Try parsing: check if it's init; tag
	// Save position to try
	savePos := p.pos
	expr := p.parseSimpleStmt()
	if p.check(TOKEN_SEMICOLON) {
		p.advance()
		stmt.Init = expr
		stmt.Tag = p.parseExpr()
	} else {
		p.pos = savePos
		stmt.Tag = p.parseExpr()
	}

	stmt.Body = p.parseBlockStmt()
	return stmt
}

// parseSelectStmt parses: select { cases }
func (p *Parser) parseSelectStmt() *SelectStmt {
	pos := p.pos()
	p.expect(TOKEN_SELECT, "expected 'select'")

	body := p.parseBlockStmt()
	return &SelectStmt{PosNode: pos, Body: body}
}

// parseReturnStmt parses: return [expr1, expr2, ...]
func (p *Parser) parseReturnStmt() *ReturnStmt {
	pos := p.pos()
	p.expect(TOKEN_RETURN, "expected 'return'")

	stmt := &ReturnStmt{PosNode: pos}

	if p.check(TOKEN_SEMICOLON) || p.check(TOKEN_RBRACE) || p.isAtEnd() {
		return stmt
	}

	stmt.Results = p.parseExprList()
	return stmt
}

// parseBranchStmt parses: break/continue/goto/fallthrough [label]
func (p *Parser) parseBranchStmt() *BranchStmt {
	pos := p.pos()
	tok := p.advance()

	stmt := &BranchStmt{
		PosNode: pos,
		Tok:     tok.Type,
	}

	if p.check(TOKEN_IDENT) {
		labelTok := p.advance()
		stmt.Label = &Ident{PosNode: Position{Line: labelTok.Line, Column: labelTok.Column}, Name: labelTok.Literal}
	}

	return stmt
}

// parseDeferStmt parses: defer expr(args)
func (p *Parser) parseDeferStmt() *DeferStmt {
	pos := p.pos()
	p.expect(TOKEN_DEFER, "expected 'defer'")

	call := p.parseCallExpr()
	return &DeferStmt{PosNode: pos, Call: call}
}

// parseGoStmt parses: go expr(args)
func (p *Parser) parseGoStmt() *GoStmt {
	pos := p.pos()
	p.expect(TOKEN_GO, "expected 'go'")

	call := p.parseCallExpr()
	return &GoStmt{PosNode: pos, Call: call}
}

// parseCallExpr parses a function call expression.
func (p *Parser) parseCallExpr() *CallExpr {
	expr := p.parseExpr()
	if call, ok := expr.(*CallExpr); ok {
		return call
	}
	// Wrap in a call if the expression is not already one
	p.error(p.peek(), "expected function call")
	return &CallExpr{PosNode: pos, Func: expr}
}

// --- Expression parsing with operator precedence ---

// Operator precedence levels (lowest to highest).
const (
	precLowest  = iota
	precLor     // ||
	precLand    // &&
	precEql     // == != < > <= >=
	precAdd     // + - | ^
	precMul     // * / % << >> & &^
	precUnary   // unary operators
	precPostfix // . [ ( call
)

// precedence returns the precedence of the given operator token.
func precedence(op TokenType) int {
	switch op {
	case TOKEN_LOR:
		return precLor
	case TOKEN_LAND:
		return precLand
	case TOKEN_EQ, TOKEN_NEQ, TOKEN_LT, TOKEN_GT, TOKEN_LEQ, TOKEN_GEQ:
		return precEql
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_PIPE, TOKEN_CARET:
		return precAdd
	case TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT, TOKEN_LSHIFT, TOKEN_RSHIFT, TOKEN_AMP, TOKEN_AMP_CARET:
		return precMul
	default:
		return precLowest
	}
}

// parseExpr parses a full expression using the Pratt parser technique.
func (p *Parser) parseExpr() Expr {
	return p.parseBinaryExpr(precLowest)
}

// parseBinaryExpr parses binary expressions with proper precedence.
func (p *Parser) parseBinaryExpr(minPrec int) Expr {
	left := p.parseUnaryExpr()

	for {
		op := p.peek()
		prec := precedence(op.Type)
		if prec < minPrec {
			break
		}

		p.advance() // consume operator
		right := p.parseBinaryExpr(prec + 1)

		left = &BinaryExpr{
			PosNode: left.Pos(),
			X:       left,
			Op:      op.Type,
			Y:       right,
		}
	}

	return left
}

// parseUnaryExpr parses unary expressions: !x, -x, +x, *x, &x, <-x.
func (p *Parser) parseUnaryExpr() Expr {
	pos := p.pos()

	switch p.peek().Type {
	case TOKEN_NOT:
		p.advance()
		x := p.parseUnaryExpr()
		return &UnaryExpr{PosNode: pos, Op: TOKEN_NOT, X: x}
	case TOKEN_MINUS:
		p.advance()
		x := p.parseUnaryExpr()
		return &UnaryExpr{PosNode: pos, Op: TOKEN_MINUS, X: x}
	case TOKEN_PLUS:
		p.advance()
		x := p.parseUnaryExpr()
		return &UnaryExpr{PosNode: pos, Op: TOKEN_PLUS, X: x}
	case TOKEN_STAR:
		p.advance()
		x := p.parseUnaryExpr()
		return &StarExpr{PosNode: pos, X: x}
	case TOKEN_AMP:
		p.advance()
		x := p.parseUnaryExpr()
		return &UnaryExpr{PosNode: pos, Op: TOKEN_AMP, X: x}
	case TOKEN_ARROW:
		p.advance()
		x := p.parseUnaryExpr()
		return &UnaryExpr{PosNode: pos, Op: TOKEN_ARROW, X: x}
	}

	return p.parsePostfixExpr()
}

// parsePostfixExpr parses postfix expressions: primary followed by ., [], (), type assertion.
func (p *Parser) parsePostfixExpr() Expr {
	expr := p.parsePrimaryExpr()

	for {
		pos := p.pos()
		switch p.peek().Type {
		case TOKEN_DOT:
			p.advance() // consume .
			if p.check(TOKEN_IDENT) {
				sel := p.advance()
				expr = &SelectorExpr{
					PosNode: pos,
					X:       expr,
					Sel:     &Ident{PosNode: Position{Line: sel.Line, Column: sel.Column}, Name: sel.Literal},
				}
			} else if p.check(TOKEN_LPAREN) {
				// Type assertion: x.(type) or x.(T)
				p.advance() // consume (
				var typ Expr
				if p.check(TOKEN_TYPE) {
					typeTok := p.advance()
					typ = &Ident{PosNode: Position{Line: typeTok.Line, Column: typeTok.Column}, Name: "type"}
				} else {
					typ = p.parseType()
				}
				p.expect(TOKEN_RPAREN, "expected ')' in type assertion")
				expr = &TypeAssertExpr{PosNode: pos, X: expr, Type: typ}
			}

		case TOKEN_LBRACK:
			p.advance() // consume [
			index := p.parseExpr()
			// Check for slice expression: [low:high] or [low:high:max]
			if p.match(TOKEN_COLON) {
				sliceExpr := &SliceExpr{
					PosNode: pos,
					X:       expr,
					Low:     index,
				}
				if !p.check(TOKEN_RBRACK) {
					sliceExpr.High = p.parseExpr()
					if p.match(TOKEN_COLON) {
						sliceExpr.Slice3 = true
						sliceExpr.Max = p.parseExpr()
					}
				}
				p.expect(TOKEN_RBRACK, "expected ']' in slice expression")
				expr = sliceExpr
			} else {
				p.expect(TOKEN_RBRACK, "expected ']' in index expression")
				expr = &IndexExpr{PosNode: pos, X: expr, Index: index}
			}

		case TOKEN_LPAREN:
			// Function call
			p.advance() // consume (
			call := &CallExpr{PosNode: pos, Func: expr}

			if !p.check(TOKEN_RPAREN) {
				// Check for ellipsis argument
				call.Args = p.parseCallArgs()
				if p.check(TOKEN_ELLIPSIS) {
					p.advance()
					call.Ellipsis = true
				}
			}

			p.expect(TOKEN_RPAREN, "expected ')' in function call")
			expr = call

		default:
			return expr
		}
	}
}

// parseCallArgs parses comma-separated function call arguments.
func (p *Parser) parseCallArgs() []Expr {
	var args []Expr
	args = append(args, p.parseExpr())
	for p.match(TOKEN_COMMA) {
		args = append(args, p.parseExpr())
	}
	return args
}

// parsePrimaryExpr parses primary expressions: identifiers, literals, composite literals, function literals.
func (p *Parser) parsePrimaryExpr() Expr {
	pos := p.pos()

	switch p.peek().Type {
	case TOKEN_IDENT:
		tok := p.advance()
		return &Ident{PosNode: Position{Line: tok.Line, Column: tok.Column}, Name: tok.Literal}

	case TOKEN_INT:
		tok := p.advance()
		return &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_INT, Value: tok.Literal}

	case TOKEN_FLOAT:
		tok := p.advance()
		return &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_FLOAT, Value: tok.Literal}

	case TOKEN_STRING:
		tok := p.advance()
		return &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_STRING, Value: tok.Literal}

	case TOKEN_RAW_STRING:
		tok := p.advance()
		return &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_RAW_STRING, Value: tok.Literal}

	case TOKEN_RUNE:
		tok := p.advance()
		return &BasicLit{PosNode: Position{Line: tok.Line, Column: tok.Column}, Kind: TOKEN_RUNE, Value: tok.Literal}

	case TOKEN_NIL:
		p.advance()
		return &Ident{PosNode: pos, Name: "nil"}

	case TOKEN_TRUE:
		p.advance()
		return &Ident{PosNode: pos, Name: "true"}

	case TOKEN_FALSE:
		p.advance()
		return &Ident{PosNode: pos, Name: "false"}

	case TOKEN_FUNC:
		// Function literal: func(params) { body }
		p.advance()
		ft := p.parseFuncType()
		body := p.parseBlockStmt()
		return &FuncLit{PosNode: pos, Type: ft, Body: body}

	case TOKEN_LPAREN:
		// Parenthesized expression
		p.advance()
		expr := p.parseExpr()
		p.expect(TOKEN_RPAREN, "expected ')'")
		return &ParenExpr{PosNode: pos, X: expr}

	case TOKEN_LBRACK:
		// Slice literal: []T{...}
		p.advance() // consume [
		if p.check(TOKEN_RBRACK) {
			// []T{...}
			p.advance() // consume ]
			elt := p.parseType()
			lit := p.parseCompositeLitBody(pos, &SliceType{PosNode: pos, Elt: elt})
			return lit
		}
		// Could be [expr] — should not happen here, but handle gracefully
		p.pos = pos.pos()
		return p.parseCompositeLit()

	case TOKEN_MAP:
		// Map literal: map[K]V{...}
		mp := p.parseType().(*MapType)
		lit := p.parseCompositeLitBody(pos, mp)
		return lit

	case TOKEN_STRUCT:
		// Struct literal: struct literal shouldn't appear here normally,
		// but handle composite literals with struct type
		return p.parseCompositeLit()

	default:
		// Try composite literal (type followed by {)
		if p.isTypeStart() {
			typ := p.parseType()
			if p.check(TOKEN_LBRACE) {
				return p.parseCompositeLitBody(pos, typ)
			}
			return typ
		}

		p.error(p.peek(), fmt.Sprintf("unexpected token %s in expression", TokenString(p.peek().Type)))
		return &Ident{PosNode: pos, Name: "unknown"}
	}
}

// parseCompositeLit parses a composite literal: Type{...}.
func (p *Parser) parseCompositeLit() Expr {
	pos := p.pos()
	typ := p.parseType()
	return p.parseCompositeLitBody(pos, typ)
}

// parseCompositeLitBody parses the { ... } portion of a composite literal.
func (p *Parser) parseCompositeLitBody(pos Position, typ Expr) *CompositeLit {
	p.expect(TOKEN_LBRACE, "expected '{' in composite literal")

	lit := &CompositeLit{PosNode: pos, Type: typ}

	for !p.check(TOKEN_RBRACE) && !p.isAtEnd() {
		expr := p.parseExpr()
		if p.match(TOKEN_COLON) {
			// Key: value
			value := p.parseExpr()
			lit.Elts = append(lit.Elts, &KeyValueExpr{PosNode: pos, Key: expr, Value: value})
		} else {
			lit.Elts = append(lit.Elts, expr)
		}
		if !p.match(TOKEN_COMMA) {
			break
		}
	}

	p.expect(TOKEN_RBRACE, "expected '}' in composite literal")
	return lit
}
