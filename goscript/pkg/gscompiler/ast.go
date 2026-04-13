package gscompiler

import "fmt"

// Position represents a source position in a .gs file.
type Position struct {
	Line   int
	Column int
}

// Pos returns the position itself, satisfying the Node interface.
func (p Position) Pos() Position { return p }

// String returns a human-readable representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// --- Node interfaces ---

// Node is the base interface for all AST nodes.
type Node interface {
	node()
	Pos() Position
}

// Expr is the interface for all expression nodes.
type Expr interface {
	Node
	exprNode()
}

// Stmt is the interface for all statement nodes.
type Stmt interface {
	Node
	stmtNode()
}

// Decl is the interface for all declaration nodes.
type Decl interface {
	Node
	declNode()
}

// Spec is the interface for specification nodes (import, value, type).
type Spec interface {
	Node
	specNode()
}

// --- Expression nodes ---

// Ident represents an identifier (variable, function, or type name).
type Ident struct {
	PosNode Position
	Name    string
}

func (i *Ident) node()     {}
func (i *Ident) exprNode() {}
func (i *Ident) Pos() Position { return i.PosNode }
func (i *Ident) String() string { return i.Name }

// BasicLit represents a basic literal value (integer, float, string, rune).
type BasicLit struct {
	PosNode Position
	Kind    TokenType
	Value   string
}

func (b *BasicLit) node()     {}
func (b *BasicLit) exprNode() {}
func (b *BasicLit) Pos() Position { return b.PosNode }
func (b *BasicLit) String() string { return b.Value }

// CompositeLit represents a composite literal: []int{1,2,3}, map[string]int{}, T{X: 42}.
type CompositeLit struct {
	PosNode Position
	Type    Expr
	Elts    []Expr // list of *KeyValueExpr or Expr
}

func (c *CompositeLit) node()     {}
func (c *CompositeLit) exprNode() {}
func (c *CompositeLit) Pos() Position { return c.PosNode }
func (c *CompositeLit) String() string {
	return fmt.Sprintf("CompositeLit{Type: %v, Elts: %v}", c.Type, c.Elts)
}

// FuncLit represents a function literal: func(params) { body }.
type FuncLit struct {
	PosNode Position
	Type    *FuncType
	Body    *BlockStmt
}

func (f *FuncLit) node()     {}
func (f *FuncLit) exprNode() {}
func (f *FuncLit) Pos() Position { return f.PosNode }
func (f *FuncLit) String() string { return "func" + f.Type.String() + " " + f.Body.String() }

// CallExpr represents a function call expression: f(args).
type CallExpr struct {
	PosNode  Position
	Func     Expr
	Args     []Expr
	Ellipsis bool // true if last argument uses ...
}

func (c *CallExpr) node()     {}
func (c *CallExpr) exprNode() {}
func (c *CallExpr) Pos() Position { return c.PosNode }
func (c *CallExpr) String() string {
	ellipsis := ""
	if c.Ellipsis {
		ellipsis = "..."
	}
	return fmt.Sprintf("%v(%v%s)", c.Func, exprListString(c.Args), ellipsis)
}

// SelectorExpr represents a field selector: x.field.
type SelectorExpr struct {
	PosNode Position
	X       Expr
	Sel     *Ident
}

func (s *SelectorExpr) node()     {}
func (s *SelectorExpr) exprNode() {}
func (s *SelectorExpr) Pos() Position { return s.PosNode }
func (s *SelectorExpr) String() string {
	return fmt.Sprintf("%v.%v", s.X, s.Sel)
}

// IndexExpr represents an index expression: x[i].
type IndexExpr struct {
	PosNode Position
	X       Expr
	Index   Expr
}

func (i *IndexExpr) node()     {}
func (i *IndexExpr) exprNode() {}
func (i *IndexExpr) Pos() Position { return i.PosNode }
func (i *IndexExpr) String() string {
	return fmt.Sprintf("%v[%v]", i.X, i.Index)
}

// BinaryExpr represents a binary expression: x + y.
type BinaryExpr struct {
	PosNode Position
	X       Expr
	Op      TokenType
	Y       Expr
}

func (b *BinaryExpr) node()     {}
func (b *BinaryExpr) exprNode() {}
func (b *BinaryExpr) Pos() Position { return b.PosNode }
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %s %v)", b.X, TokenString(b.Op), b.Y)
}

// UnaryExpr represents a unary expression: -x, !x, *x, &x, <-x.
type UnaryExpr struct {
	PosNode Position
	Op      TokenType
	X       Expr
}

func (u *UnaryExpr) node()     {}
func (u *UnaryExpr) exprNode() {}
func (u *UnaryExpr) Pos() Position { return u.PosNode }
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%v)", TokenString(u.Op), u.X)
}

// ParenExpr represents a parenthesized expression: (expr).
type ParenExpr struct {
	PosNode Position
	X       Expr
}

func (p *ParenExpr) node()     {}
func (p *ParenExpr) exprNode() {}
func (p *ParenExpr) Pos() Position { return p.PosNode }
func (p *ParenExpr) String() string {
	return fmt.Sprintf("(%v)", p.X)
}

// TypeAssertExpr represents a type assertion: x.(T).
type TypeAssertExpr struct {
	PosNode Position
	X       Expr
	Type    Expr // nil for a type switch
}

func (t *TypeAssertExpr) node()     {}
func (t *TypeAssertExpr) exprNode() {}
func (t *TypeAssertExpr) Pos() Position { return t.PosNode }
func (t *TypeAssertExpr) String() string {
	return fmt.Sprintf("%v.(%v)", t.X, t.Type)
}

// SliceExpr represents a slice expression: arr[low:high] or arr[low:high:max].
type SliceExpr struct {
	PosNode  Position
	X        Expr
	Low      Expr
	High     Expr
	Max      Expr
	Slice3   bool // true if this is a 3-index slice (arr[low:high:max])
}

func (s *SliceExpr) node()     {}
func (s *SliceExpr) exprNode() {}
func (s *SliceExpr) Pos() Position { return s.PosNode }
func (s *SliceExpr) String() string {
	return fmt.Sprintf("%v[%v:%v:%v]", s.X, s.Low, s.High, s.Max)
}

// StarExpr represents a pointer expression: *T.
type StarExpr struct {
	PosNode Position
	X       Expr
}

func (s *StarExpr) node()     {}
func (s *StarExpr) exprNode() {}
func (s *StarExpr) Pos() Position { return s.PosNode }
func (s *StarExpr) String() string {
	return fmt.Sprintf("*%v", s.X)
}

// KeyValueExpr represents a key-value pair in a composite literal: "key": value.
type KeyValueExpr struct {
	PosNode Position
	Key     Expr
	Value   Expr
}

func (k *KeyValueExpr) node()     {}
func (k *KeyValueExpr) exprNode() {}
func (k *KeyValueExpr) Pos() Position { return k.PosNode }
func (k *KeyValueExpr) String() string {
	return fmt.Sprintf("%v: %v", k.Key, k.Value)
}

// --- Statement nodes ---

// ExprStmt represents an expression statement: x = 5.
type ExprStmt struct {
	PosNode Position
	X       Expr
}

func (e *ExprStmt) node()     {}
func (e *ExprStmt) stmtNode() {}
func (e *ExprStmt) Pos() Position { return e.PosNode }
func (e *ExprStmt) String() string { return fmt.Sprintf("%v", e.X) }

// AssignStmt represents an assignment or short variable declaration: x = y, x := y.
type AssignStmt struct {
	PosNode Position
	Lhs     []Expr
	Tok     TokenType // ASSIGN or SHORT_DECL
	Rhs     []Expr
}

func (a *AssignStmt) node()     {}
func (a *AssignStmt) stmtNode() {}
func (a *AssignStmt) Pos() Position { return a.PosNode }
func (a *AssignStmt) String() string {
	return fmt.Sprintf("%v %s %v", exprListString(a.Lhs), TokenString(a.Tok), exprListString(a.Rhs))
}

// ReturnStmt represents a return statement: return expr1, expr2.
type ReturnStmt struct {
	PosNode Position
	Results []Expr
}

func (r *ReturnStmt) node()     {}
func (r *ReturnStmt) stmtNode() {}
func (r *ReturnStmt) Pos() Position { return r.PosNode }
func (r *ReturnStmt) String() string {
	if len(r.Results) == 0 {
		return "return"
	}
	return fmt.Sprintf("return %v", exprListString(r.Results))
}

// IfStmt represents an if statement: if init; cond { body } else { body }.
type IfStmt struct {
	PosNode Position
	Init    Stmt
	Cond    Expr
	Body    *BlockStmt
	Else    Stmt
}

func (i *IfStmt) node()     {}
func (i *IfStmt) stmtNode() {}
func (i *IfStmt) Pos() Position { return i.PosNode }
func (i *IfStmt) String() string {
	s := "if "
	if i.Init != nil {
		s += i.Init.String() + "; "
	}
	s += i.Cond.String() + " " + i.Body.String()
	if i.Else != nil {
		s += " else " + i.Else.String()
	}
	return s
}

// ForStmt represents a for statement: for init; cond; post { body }.
type ForStmt struct {
	PosNode Position
	Init    Stmt
	Cond    Expr
	Post    Stmt
	Body    *BlockStmt
}

func (f *ForStmt) node()     {}
func (f *ForStmt) stmtNode() {}
func (f *ForStmt) Pos() Position { return f.PosNode }
func (f *ForStmt) String() string {
	s := "for "
	if f.Init != nil {
		s += f.Init.String() + "; "
	}
	if f.Cond != nil {
		s += f.Cond.String() + "; "
	}
	if f.Post != nil {
		s += f.Post.String()
	}
	s += " " + f.Body.String()
	return s
}

// RangeStmt represents a for-range statement: for key, val := range x { body }.
type RangeStmt struct {
	PosNode Position
	Key     *Ident
	Value   *Ident
	Tok     TokenType // ASSIGN or SHORT_DECL
	X       Expr
	Body    *BlockStmt
}

func (r *RangeStmt) node()     {}
func (r *RangeStmt) stmtNode() {}
func (r *RangeStmt) Pos() Position { return r.PosNode }
func (r *RangeStmt) String() string {
	s := "for "
	if r.Key != nil {
		s += r.Key.Name
		if r.Value != nil {
			s += ", " + r.Value.Name
		}
		s += " " + TokenString(r.Tok) + " "
	}
	s += "range " + r.X.String() + " " + r.Body.String()
	return s
}

// BlockStmt represents a block of statements: { stmt1; stmt2; }.
type BlockStmt struct {
	PosNode Position
	List    []Stmt
}

func (b *BlockStmt) node()     {}
func (b *BlockStmt) stmtNode() {}
func (b *BlockStmt) Pos() Position { return b.PosNode }
func (b *BlockStmt) String() string {
	return fmt.Sprintf("{ %v }", stmtListString(b.List))
}

// SwitchStmt represents a switch statement: switch expr { cases }.
type SwitchStmt struct {
	PosNode Position
	Init    Stmt
	Tag     Expr
	Body    *BlockStmt
}

func (s *SwitchStmt) node()     {}
func (s *SwitchStmt) stmtNode() {}
func (s *SwitchStmt) Pos() Position { return s.PosNode }
func (s *SwitchStmt) String() string {
	str := "switch "
	if s.Init != nil {
		str += s.Init.String() + "; "
	}
	if s.Tag != nil {
		str += s.Tag.String() + " "
	}
	str += s.Body.String()
	return str
}

// CaseClause represents a case clause in a switch statement: case expr: body.
type CaseClause struct {
	PosNode Position
	List    []Expr // list of expressions; nil for default case
	Body    []Stmt
}

func (c *CaseClause) node()     {}
func (c *CaseClause) stmtNode() {}
func (c *CaseClause) Pos() Position { return c.PosNode }
func (c *CaseClause) String() string {
	if c.List == nil {
		return fmt.Sprintf("default: %v", stmtListString(c.Body))
	}
	return fmt.Sprintf("case %v: %v", exprListString(c.List), stmtListString(c.Body))
}

// IncDecStmt represents an increment or decrement statement: i++, j--.
type IncDecStmt struct {
	PosNode Position
	X       Expr
	Tok     TokenType // PLUS or MINUS (the operator before the =)
}

func (i *IncDecStmt) node()     {}
func (i *IncDecStmt) stmtNode() {}
func (i *IncDecStmt) Pos() Position { return i.PosNode }
func (i *IncDecStmt) String() string {
	return fmt.Sprintf("%v%s", i.X, TokenString(i.Tok)+TokenString(i.Tok))
}

// DeferStmt represents a defer statement: defer f(args).
type DeferStmt struct {
	PosNode Position
	Call    *CallExpr
}

func (d *DeferStmt) node()     {}
func (d *DeferStmt) stmtNode() {}
func (d *DeferStmt) Pos() Position { return d.PosNode }
func (d *DeferStmt) String() string { return "defer " + d.Call.String() }

// GoStmt represents a go statement: go f(args).
type GoStmt struct {
	PosNode Position
	Call    *CallExpr
}

func (g *GoStmt) node()     {}
func (g *GoStmt) stmtNode() {}
func (g *GoStmt) Pos() Position { return g.PosNode }
func (g *GoStmt) String() string { return "go " + g.Call.String() }

// SendStmt represents a channel send statement: ch <- value.
type SendStmt struct {
	PosNode Position
	Chan    Expr
	Value   Expr
}

func (s *SendStmt) node()     {}
func (s *SendStmt) stmtNode() {}
func (s *SendStmt) Pos() Position { return s.PosNode }
func (s *SendStmt) String() string {
	return fmt.Sprintf("%v <- %v", s.Chan, s.Value)
}

// SelectStmt represents a select statement: select { cases }.
type SelectStmt struct {
	PosNode Position
	Body    *BlockStmt
}

func (s *SelectStmt) node()     {}
func (s *SelectStmt) stmtNode() {}
func (s *SelectStmt) Pos() Position { return s.PosNode }
func (s *SelectStmt) String() string { return "select " + s.Body.String() }

// CommClause represents a communication clause in a select statement: case ch <- x: body.
type CommClause struct {
	PosNode Position
	Comm    Stmt // nil for default case
	Body    []Stmt
}

func (c *CommClause) node()     {}
func (c *CommClause) stmtNode() {}
func (c *CommClause) Pos() Position { return c.PosNode }
func (c *CommClause) String() string {
	if c.Comm == nil {
		return fmt.Sprintf("default: %v", stmtListString(c.Body))
	}
	return fmt.Sprintf("case %v: %v", c.Comm, stmtListString(c.Body))
}

// BranchStmt represents a branch statement: break, continue, goto, fallthrough.
type BranchStmt struct {
	PosNode Position
	Tok     TokenType // BREAK, CONTINUE, GOTO, FALLTHROUGH
	Label   *Ident
}

func (b *BranchStmt) node()     {}
func (b *BranchStmt) stmtNode() {}
func (b *BranchStmt) Pos() Position { return b.PosNode }
func (b *BranchStmt) String() string {
	s := TokenString(b.Tok)
	if b.Label != nil {
		s += " " + b.Label.Name
	}
	return s
}

// DeclStmt represents a declaration statement.
type DeclStmt struct {
	PosNode Position
	Decl    Decl
}

func (d *DeclStmt) node()     {}
func (d *DeclStmt) stmtNode() {}
func (d *DeclStmt) Pos() Position { return d.PosNode }
func (d *DeclStmt) String() string { return d.Decl.String() }

// --- Declaration nodes ---

// GenDecl represents a general declaration: var, const, type, import.
type GenDecl struct {
	PosNode Position
	Tok     TokenType // VAR, CONST, TYPE, IMPORT
	Specs   []Spec
}

func (g *GenDecl) node()     {}
func (g *GenDecl) declNode() {}
func (g *GenDecl) Pos() Position { return g.PosNode }
func (g *GenDecl) String() string {
	return fmt.Sprintf("%s %v", TokenString(g.Tok), g.Specs)
}

// FuncDecl represents a function declaration.
type FuncDecl struct {
	PosNode Position
	Recv    *FieldList // receiver (nil for regular functions)
	Name    *Ident
	Type    *FuncType
	Body    *BlockStmt
}

func (f *FuncDecl) node()     {}
func (f *FuncDecl) declNode() {}
func (f *FuncDecl) Pos() Position { return f.PosNode }
func (f *FuncDecl) String() string {
	s := "func "
	if f.Recv != nil {
		s += f.Recv.String() + " "
	}
	s += f.Name.Name + f.Type.String()
	if f.Body != nil {
		s += " " + f.Body.String()
	}
	return s
}

// --- Specification nodes ---

// ImportSpec represents an import declaration.
type ImportSpec struct {
	PosNode Position
	Name    *Ident // alias or nil
	Path    *BasicLit
}

func (i *ImportSpec) node()     {}
func (i *ImportSpec) specNode() {}
func (i *ImportSpec) Pos() Position { return i.PosNode }
func (i *ImportSpec) String() string {
	if i.Name != nil {
		return fmt.Sprintf("%s %s", i.Name.Name, i.Path.Value)
	}
	return i.Path.Value
}

// ValueSpec represents a variable or constant declaration: var x T = expr.
type ValueSpec struct {
	PosNode Position
	Names   []*Ident
	Type    Expr
	Values  []Expr
}

func (v *ValueSpec) node()     {}
func (v *ValueSpec) specNode() {}
func (v *ValueSpec) Pos() Position { return v.PosNode }
func (v *ValueSpec) String() string {
	s := identListString(v.Names)
	if v.Type != nil {
		s += " " + v.Type.String()
	}
	if len(v.Values) > 0 {
		s += " = " + exprListString(v.Values)
	}
	return s
}

// TypeSpec represents a type declaration: type Name Type.
type TypeSpec struct {
	PosNode Position
	Name    *Ident
	Type    Expr
}

func (t *TypeSpec) node()     {}
func (t *TypeSpec) specNode() {}
func (t *TypeSpec) Pos() Position { return t.PosNode }
func (t *TypeSpec) String() string {
	return fmt.Sprintf("%s %v", t.Name.Name, t.Type)
}

// --- Type nodes ---

// Field represents a struct field, method parameter, or function result.
type Field struct {
	PosNode Position
	Names   []*Ident
	Type    Expr
	Tag     *BasicLit
}

func (f *Field) node() {}
func (f *Field) Pos() Position { return f.PosNode }
func (f *Field) String() string {
	if len(f.Names) > 0 {
		return fmt.Sprintf("%v %v", identListString(f.Names), f.Type)
	}
	return f.Type.String()
}

// FieldList represents a list of fields (parameters, results, struct fields, etc.).
type FieldList struct {
	PosNode Position
	List    []*Field
}

func (f *FieldList) node() {}
func (f *FieldList) Pos() Position { return f.PosNode }
func (f *FieldList) String() string {
	if len(f.List) == 0 {
		return ""
	}
	parts := make([]string, len(f.List))
	for i, f := range f.List {
		parts[i] = f.String()
	}
	return fmt.Sprintf("(%v)", joinStrings(parts, ", "))
}

// FuncType represents a function type: func(params) (results).
type FuncType struct {
	PosNode Position
	Params  *FieldList
	Results *FieldList
}

func (f *FuncType) node()     {}
func (f *FuncType) exprNode() {}
func (f *FuncType) Pos() Position { return f.PosNode }
func (f *FuncType) String() string {
	s := f.Params.String()
	if f.Results != nil && len(f.Results.List) > 0 {
		s += " " + f.Results.String()
	}
	return s
}

// StructType represents a struct type: struct { fields }.
type StructType struct {
	PosNode Position
	Fields  *FieldList
}

func (s *StructType) node()     {}
func (s *StructType) exprNode() {}
func (s *StructType) Pos() Position { return s.PosNode }
func (s *StructType) String() string {
	return "struct " + s.Fields.String()
}

// InterfaceType represents an interface type: interface { methods }.
type InterfaceType struct {
	PosNode Position
	Methods *FieldList
}

func (i *InterfaceType) node()     {}
func (i *InterfaceType) exprNode() {}
func (i *InterfaceType) Pos() Position { return i.PosNode }
func (i *InterfaceType) String() string {
	return "interface " + i.Methods.String()
}

// ArrayType represents an array type: [N]T.
type ArrayType struct {
	PosNode Position
	Len     Expr
	Elt     Expr
}

func (a *ArrayType) node()     {}
func (a *ArrayType) exprNode() {}
func (a *ArrayType) Pos() Position { return a.PosNode }
func (a *ArrayType) String() string {
	return fmt.Sprintf("[%v]%v", a.Len, a.Elt)
}

// SliceType represents a slice type: []T.
type SliceType struct {
	PosNode Position
	Elt     Expr
}

func (s *SliceType) node()     {}
func (s *SliceType) exprNode() {}
func (s *SliceType) Pos() Position { return s.PosNode }
func (s *SliceType) String() string {
	return fmt.Sprintf("[]%v", s.Elt)
}

// MapType represents a map type: map[K]V.
type MapType struct {
	PosNode Position
	Key     Expr
	Value   Expr
}

func (m *MapType) node()     {}
func (m *MapType) exprNode() {}
func (m *MapType) Pos() Position { return m.PosNode }
func (m *MapType) String() string {
	return fmt.Sprintf("map[%v]%v", m.Key, m.Value)
}

// ChanType represents a channel type: chan T, <-chan T, chan<- T.
type ChanType struct {
	PosNode Position
	Dir     TokenType // CHAN (bidirectional), ARROW (recv-only), SEND (send-only)
	Value   Expr
}

func (c *ChanType) node()     {}
func (c *ChanType) exprNode() {}
func (c *ChanType) Pos() Position { return c.PosNode }
func (c *ChanType) String() string {
	switch c.Dir {
	case TOKEN_ARROW:
		return fmt.Sprintf("<-chan %v", c.Value)
	case TOKEN_CHAN:
		return fmt.Sprintf("chan %v", c.Value)
	default:
		return fmt.Sprintf("chan<- %v", c.Value)
	}
}

// Ellipsis represents an ellipsis type ...T used for variadic parameters.
type Ellipsis struct {
	PosNode Position
	Elt     Expr
}

func (e *Ellipsis) node()     {}
func (e *Ellipsis) exprNode() {}
func (e *Ellipsis) Pos() Position { return e.PosNode }
func (e *Ellipsis) String() string {
	return fmt.Sprintf("...%v", e.Elt)
}

// --- Program node ---

// Program represents a complete Go source file.
type Program struct {
	PosNode Position
	Package *Ident
	Imports []*ImportSpec
	Decls   []Decl
}

func (p *Program) node() {}
func (p *Program) Pos() Position { return p.PosNode }
func (p *Program) String() string {
	return fmt.Sprintf("package %s", p.Package.Name)
}

// --- Helper functions ---

// exprListString joins a list of expression strings.
func exprListString(exprs []Expr) string {
	parts := make([]string, len(exprs))
	for i, e := range exprs {
		parts[i] = e.String()
	}
	return joinStrings(parts, ", ")
}

// stmtListString joins a list of statement strings.
func stmtListString(stmts []Stmt) string {
	parts := make([]string, len(stmts))
	for i, s := range stmts {
		parts[i] = s.String()
	}
	return joinStrings(parts, "; ")
}

// identListString joins a list of identifier names.
func identListString(idents []*Ident) string {
	parts := make([]string, len(idents))
	for i, id := range idents {
		parts[i] = id.Name
	}
	return joinStrings(parts, ", ")
}

// joinStrings joins a slice of strings with a separator.
func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
