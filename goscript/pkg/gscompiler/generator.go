package gscompiler

import (
        "fmt"
        "strings"
)

// stdlibMapping maps Go import paths to their JavaScript equivalents.
var stdlibMapping = map[string]string{
        "goscript/dom":     "__gs.dom",
        "goscript/state":   "__gs.state",
        "goscript/api":     "__gs.api",
        "goscript/fmt":     "__gs.fmt",
        "goscript/router":  "__gs.router",
        "goscript/realtime": "__gs.realtime",
        "goscript/ui":      "__gs.ui",
        "fmt":              "__gs.fmt",
        "strings":          "__gs.strings",
        "strconv":          "__gs.strconv",
        "errors":           "__gs.errors",
}

// builtinMapping maps Go built-in function names to JavaScript equivalents.
var builtinMapping = map[string]string{
        "len":    "__gs_len",
        "cap":    "__gs_cap",
        "append": "__gs_append",
        "make":   "__gs_make",
        "new":    "__gs_new",
        "panic":  "__gs_panic",
        "recover": "__gs_recover",
        "delete": "__gs_delete",
        "close":  "__gs_close",
}

// Generator converts a Go AST into JavaScript source code.
// It handles type erasure, struct-to-class conversion, range loop transformation,
// multi-return value destructuring, and other Go-to-JS idioms.
type Generator struct {
        indent     int
        imports    map[string]string // Go import path → JS variable name
        importUsed map[string]bool   // Track which imports are actually used
        output     strings.Builder
        // Track deferred statements in the current function scope
        defers   []string
        inMethod string // Name of the struct for method receiver context
}

// NewGenerator creates a new JavaScript code generator.
func NewGenerator() *Generator {
        return &Generator{
                imports:    make(map[string]string),
                importUsed: make(map[string]bool),
        }
}

// Generate converts a complete Go program AST into JavaScript source code.
func (g *Generator) Generate(program *Program) (string, error) {
        // Collect imports
        for _, imp := range program.Imports {
                path := strings.Trim(imp.Path.Value, `"`)
                if _, ok := stdlibMapping[path]; ok {
                        var alias string
                        if imp.Name != nil {
                                alias = imp.Name.Name
                        } else {
                                // Derive a short name from the last segment
                                parts := strings.Split(path, "/")
                                alias = parts[len(parts)-1]
                        }
                        g.imports[path] = alias
                }
        }

        // Generate declarations
        for _, decl := range program.Decls {
                g.generateDecl(decl)
                g.emit("\n")
        }

        return g.output.String(), nil
}

// --- Output helpers ---

// emit writes a string to the output buffer.
func (g *Generator) emit(s string) {
        g.output.WriteString(s)
}

// emitLine writes a string followed by a newline to the output buffer.
func (g *Generator) emitLine(s string) {
        g.output.WriteString(s)
        g.output.WriteString("\n")
}

// emitIndent writes the current indentation level.
func (g *Generator) emitIndent() {
        for i := 0; i < g.indent; i++ {
                g.output.WriteString("  ")
        }
}

// indentLevel increases the indentation level by one.
func (g *Generator) indentLevel() {
        g.indent++
}

// dedentLevel decreases the indentation level by one.
func (g *Generator) dedentLevel() {
        g.indent--
        if g.indent < 0 {
                g.indent = 0
        }
}

// --- Declaration generation ---

// generateDecl generates JavaScript for a declaration node.
func (g *Generator) generateDecl(decl Decl) {
        switch d := decl.(type) {
        case *FuncDecl:
                g.generateFuncDecl(d)
        case *GenDecl:
                g.generateGenDecl(d)
        }
}

// generateFuncDecl generates JavaScript for a function declaration.
func (g *Generator) generateFuncDecl(fn *FuncDecl) {
        // Save and reset defers
        oldDefers := g.defers
        g.defers = nil

        if fn.Recv != nil && len(fn.Recv.List) > 0 {
                // Method: func (r *T) Method() { ... }
                g.generateMethodDecl(fn)
        } else if fn.Name.Name == "main" {
                // main function: execute automatically
                g.generateMainFunc(fn)
        } else {
                // Regular function: function Name(params) { ... }
                g.generateRegularFunc(fn)
        }

        // Restore defers
        g.defers = oldDefers
}

// generateMethodDecl generates a JavaScript method declaration on a prototype.
func (g *Generator) generateMethodDecl(fn *FuncDecl) {
        // Extract receiver type name
        recvType := g.extractReceiverType(fn.Recv)
        oldMethod := g.inMethod
        g.inMethod = recvType

        paramNames := g.generateParamNames(fn.Type.Params, true)
        g.emitIndent()
        g.emit(fmt.Sprintf("%s.prototype.%s = function(%s) {\n", recvType, fn.Name.Name, paramNames))

        g.indentLevel()
        g.generateBlockStmtBody(fn.Body)
        if len(g.defers) > 0 {
                g.generateDefers()
        }
        g.dedentLevel()

        g.emitIndent()
        g.emit("};\n")
        g.inMethod = oldMethod
}

// generateMainFunc generates the main function body (auto-executed).
func (g *Generator) generateMainFunc(fn *FuncDecl) {
        g.emitLine("// Auto-execute main function")
        g.generateBlockStmt(fn.Body)
        g.emit("\n")
}

// generateRegularFunc generates a regular JavaScript function declaration.
func (g *Generator) generateRegularFunc(fn *FuncDecl) {
        paramNames := g.generateParamNames(fn.Type.Params, false)

        g.emitIndent()
        g.emit(fmt.Sprintf("function %s(%s) {\n", fn.Name.Name, paramNames))

        g.indentLevel()
        g.generateBlockStmtBody(fn.Body)
        if len(g.defers) > 0 {
                g.generateDefers()
        }
        g.dedentLevel()

        g.emitIndent()
        g.emit("}\n")
}

// generateDefers emits deferred statements after function body.
func (g *Generator) generateDefers() {
        if len(g.defers) == 0 {
                return
        }
        // Emit simple defer calls at the end of the function
        for _, d := range g.defers {
                g.emitIndent()
                g.emitLine(d + ";")
        }
}

// generateGenDecl generates JavaScript for a general declaration (var, const, type).
func (g *Generator) generateGenDecl(decl *GenDecl) {
        switch decl.Tok {
        case TOKEN_TYPE:
                for _, spec := range decl.Specs {
                        if ts, ok := spec.(*TypeSpec); ok {
                                g.generateTypeSpec(ts)
                        }
                }
        case TOKEN_VAR, TOKEN_CONST:
                for _, spec := range decl.Specs {
                        if vs, ok := spec.(*ValueSpec); ok {
                                g.generateValueSpec(vs, decl.Tok == TOKEN_CONST)
                        }
                }
        }
}

// generateTypeSpec generates JavaScript for a type declaration.
func (g *Generator) generateTypeSpec(ts *TypeSpec) {
        switch t := ts.Type.(type) {
        case *StructType:
                g.generateStructClass(ts.Name.Name, t)
        case *InterfaceType:
                // Interfaces are erased in JavaScript
                g.emitIndent()
                g.emitLine(fmt.Sprintf("// interface %s (erased)", ts.Name.Name))
        default:
                // Type alias: type Name = OtherType (erased)
                g.emitIndent()
                g.emitLine(fmt.Sprintf("// type %s (erased)", ts.Name.Name))
        }
}

// generateStructClass generates a JavaScript class from a Go struct type.
func (g *Generator) generateStructClass(name string, st *StructType) {
        g.emitIndent()
        g.emit(fmt.Sprintf("class %s {\n", name))

        g.indentLevel()

        if st.Fields != nil {
                // Generate constructor
                g.emitIndent()
                g.emit("constructor(")
                params := make([]string, 0)
                for _, field := range st.Fields.List {
                        if len(field.Names) > 0 {
                                jsName := toCamelCase(field.Names[0].Name)
                                params = append(params, jsName)
                        }
                }
                g.emit(strings.Join(params, ", "))
                g.emit(") {\n")

                g.indentLevel()
                // Initialize all fields (even those not in constructor params)
                for _, field := range st.Fields.List {
                        if len(field.Names) > 0 {
                                jsName := toCamelCase(field.Names[0].Name)
                                g.emitIndent()
                                g.emitLine(fmt.Sprintf("this.%s = %s;", jsName, jsName))
                        }
                }
                g.dedentLevel()

                g.emitIndent()
                g.emit("}\n\n")
        }

        g.dedentLevel()
        g.emitIndent()
        g.emit("}\n")
}

// generateValueSpec generates JavaScript for a variable or constant declaration.
func (g *Generator) generateValueSpec(vs *ValueSpec, isConst bool) {
        keyword := "let"
        if isConst {
                keyword = "const"
        }

        for i, name := range vs.Names {
                g.emitIndent()
                jsName := name.Name
                if i < len(vs.Values) {
                        value := g.generateExpr(vs.Values[i])
                        g.emitLine(fmt.Sprintf("%s %s = %s;", keyword, jsName, value))
                } else if vs.Type != nil {
                        // var x T → let x = defaultValue(T)
                        defaultVal := g.defaultValue(vs.Type)
                        g.emitLine(fmt.Sprintf("%s %s = %s;", keyword, jsName, defaultVal))
                } else {
                        g.emitLine(fmt.Sprintf("%s %s;", keyword, jsName))
                }
        }
}

// --- Statement generation ---

// generateBlockStmt generates JavaScript for a block statement.
func (g *Generator) generateBlockStmt(block *BlockStmt) {
        g.emitLine("{")
        g.indentLevel()
        g.generateBlockStmtBody(block)
        g.dedentLevel()
        g.emitIndent()
        g.emit("}")
}

// generateBlockStmtBody generates the statements inside a block.
func (g *Generator) generateBlockStmtBody(block *BlockStmt) {
        for _, stmt := range block.List {
                g.generateStmt(stmt)
        }
}

// generateStmt generates JavaScript for a single statement.
func (g *Generator) generateStmt(stmt Stmt) {
        switch s := stmt.(type) {
        case *ExprStmt:
                g.generateExprStmt(s)
        case *AssignStmt:
                g.generateAssignStmt(s)
        case *ReturnStmt:
                g.generateReturnStmt(s)
        case *IfStmt:
                g.generateIfStmt(s)
        case *ForStmt:
                g.generateForStmt(s)
        case *RangeStmt:
                g.generateRangeStmt(s)
        case *BlockStmt:
                g.generateBlockStmt(s)
        case *SwitchStmt:
                g.generateSwitchStmt(s)
        case *IncDecStmt:
                g.generateIncDecStmt(s)
        case *DeferStmt:
                g.generateDeferStmt(s)
        case *GoStmt:
                g.generateGoStmt(s)
        case *BranchStmt:
                g.generateBranchStmt(s)
        case *DeclStmt:
                g.generateDecl(s.Decl)
        case *SendStmt:
                g.generateSendStmt(s)
        }
}

// generateExprStmt generates JavaScript for an expression statement.
func (g *Generator) generateExprStmt(s *ExprStmt) {
        g.emitIndent()
        expr := g.generateExpr(s.X)
        // Add semicolon
        if !strings.HasSuffix(expr, ";") && !strings.HasSuffix(expr, "}") {
                expr += ";"
        }
        g.emitLine(expr)
}

// generateAssignStmt generates JavaScript for assignment and short variable declarations.
func (g *Generator) generateAssignStmt(s *AssignStmt) {
        if s.Tok == TOKEN_SHORT_DECL && len(s.Lhs) > 1 {
                // Multi-return short decl: x, err := f()
                // Check if RHS is a single call expression
                if len(s.Rhs) == 1 {
                        if call, ok := s.Rhs[0].(*CallExpr); ok {
                                g.generateMultiReturnAssign(s.Lhs, call)
                                return
                        }
                }
                // Multiple RHS values
                g.generateMultiValueAssign(s.Lhs, s.Rhs, true)
                return
        }

        if s.Tok == TOKEN_SHORT_DECL {
                // x := value → let x = value;
                g.emitIndent()
                name := g.generateExpr(s.Lhs[0])
                value := g.generateExpr(s.Rhs[0])
                g.emitLine(fmt.Sprintf("let %s = %s;", name, value))
                return
        }

        // Regular assignment: x = value, x += value, etc.
        g.emitIndent()
        lhs := g.generateExpr(s.Lhs[0])
        op := g.assignOp(s.Tok)
        rhs := g.generateExpr(s.Rhs[0])
        g.emitLine(fmt.Sprintf("%s %s %s;", lhs, op, rhs))
}

// generateMultiReturnAssign handles x, err := f() by destructuring.
func (g *Generator) generateMultiReturnAssign(lhs []Expr, call *CallExpr) {
        // const [_gs_r0, _gs_r1, ...] = f();
        // let x = _gs_r0; let err = _gs_r1;
        tmpVars := make([]string, len(lhs))
        for i := range lhs {
                tmpVars[i] = fmt.Sprintf("_gs_r%d", i)
        }

        callExpr := g.generateExpr(call)
        g.emitIndent()
        g.emitLine(fmt.Sprintf("const [%s] = %s;", strings.Join(tmpVars, ", "), callExpr))

        for i, lv := range lhs {
                name := g.generateExpr(lv)
                g.emitIndent()
                g.emitLine(fmt.Sprintf("let %s = %s;", name, tmpVars[i]))
        }
}

// generateMultiValueAssign handles x, y = a, b or x, y := a, b (multi-value).
func (g *Generator) generateMultiValueAssign(lhs []Expr, rhs []Expr, isShortDecl bool) {
        keyword := "let"
        if !isShortDecl {
                keyword = ""
        }

        g.emitIndent()
        lhsStrs := make([]string, len(lhs))
        for i, l := range lhs {
                lhsStrs[i] = g.generateExpr(l)
        }
        rhsStrs := make([]string, len(rhs))
        for i, r := range rhs {
                rhsStrs[i] = g.generateExpr(r)
        }

        if keyword != "" {
                g.emitLine(fmt.Sprintf("%s %s = [%s];", keyword, strings.Join(lhsStrs, ", "), strings.Join(rhsStrs, ", ")))
        } else {
                g.emitLine(fmt.Sprintf("[%s] = [%s];", strings.Join(lhsStrs, ", "), strings.Join(rhsStrs, ", ")))
        }
}

// assignOp converts a Go assignment operator token to its JavaScript string.
func (g *Generator) assignOp(tok TokenType) string {
        switch tok {
        case TOKEN_ASSIGN:
                return "="
        case TOKEN_PLUS_ASSIGN:
                return "+="
        case TOKEN_MINUS_ASSIGN:
                return "-="
        case TOKEN_STAR_ASSIGN:
                return "*="
        case TOKEN_SLASH_ASSIGN:
                return "/="
        case TOKEN_PERCENT_ASSIGN:
                return "%="
        case TOKEN_AMP_ASSIGN:
                return "&="
        case TOKEN_PIPE_ASSIGN:
                return "|="
        case TOKEN_CARET_ASSIGN:
                return "^="
        case TOKEN_LSHIFT_ASSIGN:
                return "<<="
        case TOKEN_RSHIFT_ASSIGN:
                return ">>="
        default:
                return "="
        }
}

// generateReturnStmt generates JavaScript for a return statement.
func (g *Generator) generateReturnStmt(s *ReturnStmt) {
        g.emitIndent()
        if len(s.Results) == 0 {
                g.emitLine("return;")
        } else if len(s.Results) == 1 {
                g.emitLine(fmt.Sprintf("return %s;", g.generateExpr(s.Results[0])))
        } else {
                // Multiple return values: return [a, b, c]
                vals := make([]string, len(s.Results))
                for i, r := range s.Results {
                        vals[i] = g.generateExpr(r)
                }
                g.emitLine(fmt.Sprintf("return [%s];", strings.Join(vals, ", ")))
        }
}

// generateIfStmt generates JavaScript for an if statement.
func (g *Generator) generateIfStmt(s *IfStmt) {
        // If there's an init statement, wrap in a block
        if s.Init != nil {
                g.emitIndent()
                g.emitLine("{")
                g.indentLevel()
                g.generateStmt(s.Init)
        }

        g.emitIndent()
        cond := g.generateExpr(s.Cond)
        g.emitLine(fmt.Sprintf("if (%s) {", cond))

        g.indentLevel()
        g.generateBlockStmtBody(s.Body)
        g.dedentLevel()

        if s.Else != nil {
                g.emitIndent()
                switch elseStmt := s.Else.(type) {
                case *IfStmt:
                        g.emit("} else ")
                        // Generate the if statement inline without extra indent
                        g.dedentLevel()
                        g.generateIfStmt(elseStmt)
                        g.indentLevel()
                case *BlockStmt:
                        g.emitLine("} else {")
                        g.indentLevel()
                        g.generateBlockStmtBody(elseStmt)
                        g.dedentLevel()
                        g.emitIndent()
                        g.emitLine("}")
                }
        } else {
                g.emitIndent()
                g.emitLine("}")
        }

        if s.Init != nil {
                g.dedentLevel()
                g.emitIndent()
                g.emitLine("}")
        }
}

// generateForStmt generates JavaScript for a for statement.
func (g *Generator) generateForStmt(s *ForStmt) {
        g.emitIndent()

        if s.Cond == nil && s.Init == nil && s.Post == nil {
                // Infinite loop: for { body }
                g.emitLine("while (true) {")
        } else if s.Init == nil && s.Post == nil {
                // while loop: for cond { body }
                g.emitLine(fmt.Sprintf("while (%s) {", g.generateExpr(s.Cond)))
        } else {
                // Classic for loop: for init; cond; post { body }
                var initStr, condStr, postStr string
                if s.Init != nil {
                        initStr = g.generateStmtAsString(s.Init)
                }
                if s.Cond != nil {
                        condStr = g.generateExpr(s.Cond)
                } else {
                        condStr = "true"
                }
                if s.Post != nil {
                        postStr = g.generateStmtAsString(s.Post)
                }
                g.emitLine(fmt.Sprintf("for (%s; %s; %s) {", initStr, condStr, postStr))
        }

        g.indentLevel()
        g.generateBlockStmtBody(s.Body)
        g.dedentLevel()

        g.emitIndent()
        g.emitLine("}")
}

// generateRangeStmt generates JavaScript for a for-range statement.
func (g *Generator) generateRangeStmt(s *RangeStmt) {
        g.emitIndent()
        xExpr := g.generateExpr(s.X)

        if s.Key == nil && s.Value == nil {
                // for range x → for (const _ of x)
                g.emitLine(fmt.Sprintf("for (const _ of %s) {", xExpr))
        } else if s.Key != nil && s.Value == nil {
                keyName := s.Key.Name
                // for key := range x → for (const key of x)
                // or for _, key := range x (key is actually value)
                if keyName == "_" {
                        g.emitLine(fmt.Sprintf("for (const _ of %s) {", xExpr))
                } else {
                        g.emitLine(fmt.Sprintf("for (const %s of %s) {", keyName, xExpr))
                }
        } else {
                keyName := "_"
                valName := "_"
                if s.Key != nil {
                        keyName = s.Key.Name
                }
                if s.Value != nil {
                        valName = s.Value.Name
                }

                if keyName == "_" {
                        // for _, val := range items → for (const val of items)
                        g.emitLine(fmt.Sprintf("for (const %s of %s) {", valName, xExpr))
                } else {
                        // for key, val := range x → for (const [key, val] of x.entries())
                        // Use Object.entries for maps, .entries() for arrays
                        g.emitLine(fmt.Sprintf("for (const [%s, %s] of __gs_entries(%s)) {", keyName, valName, xExpr))
                }
        }

        g.indentLevel()
        g.generateBlockStmtBody(s.Body)
        g.dedentLevel()

        g.emitIndent()
        g.emitLine("}")
}

// generateSwitchStmt generates JavaScript for a switch statement.
func (g *Generator) generateSwitchStmt(s *SwitchStmt) {
        g.emitIndent()

        if s.Init != nil {
                g.emitLine("{")
                g.indentLevel()
                g.generateStmt(s.Init)
        }

        if s.Tag != nil {
                tagExpr := g.generateExpr(s.Tag)
                g.emitLine(fmt.Sprintf("switch (%s) {", tagExpr))
        } else {
                g.emitLine("switch (true) {")
        }

        g.indentLevel()
        if s.Body != nil {
                for _, stmt := range s.Body.List {
                        if cc, ok := stmt.(*CaseClause); ok {
                                g.generateCaseClause(cc)
                        }
                }
        }
        g.dedentLevel()

        g.emitIndent()
        g.emitLine("}")

        if s.Init != nil {
                g.dedentLevel()
                g.emitIndent()
                g.emitLine("}")
        }
}

// generateCaseClause generates JavaScript for a case clause.
func (g *Generator) generateCaseClause(cc *CaseClause) {
        if cc.List == nil {
                // default case
                g.emitIndent()
                g.emitLine("default:")
        } else {
                for _, expr := range cc.List {
                        g.emitIndent()
                        g.emitLine(fmt.Sprintf("case %s:", g.generateExpr(expr)))
                }
        }

        g.indentLevel()
        for _, stmt := range cc.Body {
                g.generateStmt(stmt)
        }
        // Go fallthrough is implicit, JS needs break
        g.emitIndent()
        g.emitLine("break;")
        g.dedentLevel()
}

// generateIncDecStmt generates JavaScript for increment/decrement.
func (g *Generator) generateIncDecStmt(s *IncDecStmt) {
        g.emitIndent()
        expr := g.generateExpr(s.X)
        op := "++"
        if s.Tok == TOKEN_MINUS {
                op = "--"
        }
        g.emitLine(fmt.Sprintf("%s%s;", expr, op))
}

// generateDeferStmt generates JavaScript for a defer statement.
// Simplified: collect deferred calls and emit them at function end.
func (g *Generator) generateDeferStmt(s *DeferStmt) {
        callExpr := g.generateExpr(s.Call)
        g.defers = append(g.defers, callExpr)
}

// generateGoStmt generates JavaScript for a go statement.
// go f() → (async () => { f(); })()
func (g *Generator) generateGoStmt(s *GoStmt) {
        g.emitIndent()
        callExpr := g.generateExpr(s.Call)
        g.emitLine(fmt.Sprintf("(async () => { %s; })();", callExpr))
}

// generateBranchStmt generates JavaScript for break/continue/goto.
func (g *Generator) generateBranchStmt(s *BranchStmt) {
        g.emitIndent()
        switch s.Tok {
        case TOKEN_BREAK:
                g.emitLine("break;")
        case TOKEN_CONTINUE:
                g.emitLine("continue;")
        case TOKEN_FALLTHROUGH:
                // fallthrough is handled by removing the break in case clause
        case TOKEN_GOTO:
                if s.Label != nil {
                        g.emitLine(fmt.Sprintf("// goto %s (not supported)", s.Label.Name))
                }
        }
}

// generateSendStmt generates JavaScript for a channel send.
func (g *Generator) generateSendStmt(s *SendStmt) {
        g.emitIndent()
        chExpr := g.generateExpr(s.Chan)
        valExpr := g.generateExpr(s.Value)
        g.emitLine(fmt.Sprintf("__gs_send(%s, %s);", chExpr, valExpr))
}

// generateStmtAsString generates a statement as an inline string (for use in for loop headers).
func (g *Generator) generateStmtAsString(stmt Stmt) string {
        switch s := stmt.(type) {
        case *AssignStmt:
                if s.Tok == TOKEN_SHORT_DECL {
                        name := g.generateExpr(s.Lhs[0])
                        value := g.generateExpr(s.Rhs[0])
                        return fmt.Sprintf("let %s = %s", name, value)
                }
                lhs := g.generateExpr(s.Lhs[0])
                op := g.assignOp(s.Tok)
                rhs := g.generateExpr(s.Rhs[0])
                return fmt.Sprintf("%s %s %s", lhs, op, rhs)
        case *ExprStmt:
                return g.generateExpr(s.X)
        case *IncDecStmt:
                expr := g.generateExpr(s.X)
                op := "++"
                if s.Tok == TOKEN_MINUS {
                        op = "--"
                }
                return fmt.Sprintf("%s%s", expr, op)
        }
        return ""
}

// --- Expression generation ---

// generateExpr generates JavaScript for an expression node.
func (g *Generator) generateExpr(expr Expr) string {
        switch e := expr.(type) {
        case *Ident:
                return g.generateIdent(e)
        case *BasicLit:
                return g.generateBasicLit(e)
        case *CompositeLit:
                return g.generateCompositeLit(e)
        case *FuncLit:
                return g.generateFuncLit(e)
        case *CallExpr:
                return g.generateCallExpr(e)
        case *SelectorExpr:
                return g.generateSelectorExpr(e)
        case *IndexExpr:
                return g.generateIndexExpr(e)
        case *BinaryExpr:
                return g.generateBinaryExpr(e)
        case *UnaryExpr:
                return g.generateUnaryExpr(e)
        case *ParenExpr:
                return fmt.Sprintf("(%s)", g.generateExpr(e.X))
        case *TypeAssertExpr:
                return g.generateExpr(e.X)
        case *SliceExpr:
                return g.generateSliceExpr(e)
        case *StarExpr:
                return g.generateExpr(e.X)
        case *KeyValueExpr:
                return fmt.Sprintf("%s: %s", g.generateExpr(e.Key), g.generateExpr(e.Value))
        case *SliceType:
                return "[]" // type expression, erased
        case *MapType:
                return "{}" // type expression, erased
        case *StructType:
                return "Object" // type expression, erased
        case *ArrayType:
                return "Array" // type expression, erased
        case *InterfaceType:
                return "Object" // type expression, erased
        case *ChanType:
                return "__gs_chan()" // type expression
        case *FuncType:
                return "function" // type expression
        case *Ellipsis:
                return g.generateExpr(e.Elt)
        default:
                return "undefined"
        }
}

// generateIdent generates JavaScript for an identifier.
func (g *Generator) generateIdent(id *Ident) string {
        switch id.Name {
        case "nil":
                return "null"
        case "true":
                return "true"
        case "false":
                return "false"
        default:
                return id.Name
        }
}

// generateBasicLit generates JavaScript for a literal value.
func (g *Generator) generateBasicLit(lit *BasicLit) string {
        switch lit.Kind {
        case TOKEN_STRING:
                // Convert Go string to JS string (preserve escape sequences)
                return fmt.Sprintf(`"%s"`, lit.Value)
        case TOKEN_RAW_STRING:
                // Raw string: use template literal
                return fmt.Sprintf("`%s`", lit.Value)
        case TOKEN_INT:
                return lit.Value
        case TOKEN_FLOAT:
                return lit.Value
        case TOKEN_RUNE:
                return fmt.Sprintf(`"%s"`, lit.Value)
        default:
                return lit.Value
        }
}

// generateCompositeLit generates JavaScript for a composite literal.
func (g *Generator) generateCompositeLit(lit *CompositeLit) string {
        if lit.Type == nil {
                return g.generateCompositeLitBody(lit)
        }

        switch t := lit.Type.(type) {
        case *SliceType:
                // []int{1, 2, 3} → [1, 2, 3]
                if len(lit.Elts) == 0 {
                        return "[]"
                }
                return g.generateCompositeLitBody(lit)
        case *MapType:
                // map[string]int{"a": 1} → {"a": 1}
                if len(lit.Elts) == 0 {
                        return "{}"
                }
                return g.generateCompositeLitBody(lit)
        case *StructType:
                // Struct literal inline: T{X: 1} → new T() or Object.assign
                // We can't easily construct named structs inline, so use a helper
                return g.generateCompositeLitBody(lit)
        case *Ident:
                // T{X: 1, Y: "hi"} → new T(1, "hi") or Object.assign(new T(), {x: 1, y: "hi"})
                return g.generateStructCompositeLit(t.Name, lit.Elts)
        case *SelectorExpr:
                // pkg.T{...} → similar handling
                typeName := t.Sel.Name
                return g.generateStructCompositeLit(typeName, lit.Elts)
        default:
                return g.generateCompositeLitBody(lit)
        }
}

// generateStructCompositeLit generates JS for T{Field: value, ...}.
func (g *Generator) generateStructCompositeLit(typeName string, elts []Expr) string {
        // Check if all elements are KeyValueExpr with field names
        allKeyValues := true
        for _, elt := range elts {
                if _, ok := elt.(*KeyValueExpr); !ok {
                        allKeyValues = false
                        break
                }
        }

        if allKeyValues && len(elts) > 0 {
                // Use Object.assign(new T(), { key: value, ... })
                parts := make([]string, len(elts))
                for i, elt := range elts {
                        kv := elt.(*KeyValueExpr)
                        key := ""
                        if ident, ok := kv.Key.(*Ident); ok {
                                key = toCamelCase(ident.Name)
                        } else {
                                key = g.generateExpr(kv.Key)
                        }
                        parts[i] = fmt.Sprintf("%s: %s", key, g.generateExpr(kv.Value))
                }
                return fmt.Sprintf("Object.assign(new %s(), { %s })", typeName, strings.Join(parts, ", "))
        }

        // Positional: new T(val1, val2)
        vals := make([]string, len(elts))
        for i, elt := range elts {
                vals[i] = g.generateExpr(elt)
        }
        return fmt.Sprintf("new %s(%s)", typeName, strings.Join(vals, ", "))
}

// generateCompositeLitBody generates the { ... } portion of a composite literal.
func (g *Generator) generateCompositeLitBody(lit *CompositeLit) string {
        if len(lit.Elts) == 0 {
                return "[]"
        }
        parts := make([]string, len(lit.Elts))
        for i, elt := range lit.Elts {
                parts[i] = g.generateExpr(elt)
        }
        return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// generateFuncLit generates JavaScript for a function literal.
func (g *Generator) generateFuncLit(fl *FuncLit) string {
        paramNames := g.generateParamNames(fl.Type.Params, false)

        var body strings.Builder
        body.WriteString(fmt.Sprintf("function(%s) {", paramNames))
        // For simplicity, just use the generated block content
        for _, stmt := range fl.Body.List {
                // We can't easily use the generator here due to state, so we'll inline
                body.WriteString(fmt.Sprintf("\n    %s", g.generateStmtSimple(stmt)))
        }
        body.WriteString("\n  }")
        return body.String()
}

// generateStmtSimple generates a simple statement string.
func (g *Generator) generateStmtSimple(stmt Stmt) string {
        switch s := stmt.(type) {
        case *ExprStmt:
                return g.generateExpr(s.X) + ";"
        case *ReturnStmt:
                if len(s.Results) == 0 {
                        return "return;"
                }
                return fmt.Sprintf("return %s;", g.generateExpr(s.Results[0]))
        case *AssignStmt:
                if s.Tok == TOKEN_SHORT_DECL {
                        return fmt.Sprintf("let %s = %s;", g.generateExpr(s.Lhs[0]), g.generateExpr(s.Rhs[0]))
                }
                return fmt.Sprintf("%s %s %s;", g.generateExpr(s.Lhs[0]), g.assignOp(s.Tok), g.generateExpr(s.Rhs[0]))
        default:
                return fmt.Sprintf("/* unsupported stmt: %T */", s)
        }
}

// generateCallExpr generates JavaScript for a function call.
func (g *Generator) generateCallExpr(call *CallExpr) string {
        // Handle special built-in functions
        if sel, ok := call.Func.(*SelectorExpr); ok {
                return g.generateSelectorCall(sel, call.Args, call.Ellipsis)
        }

        // Handle direct built-in calls
        if ident, ok := call.Func.(*Ident); ok {
                switch ident.Name {
                case "fmt", "console":
                        return g.generateFmtCall(call.Args)
                case "make":
                        return g.generateMakeCall(call.Args)
                case "len":
                        if len(call.Args) > 0 {
                                return fmt.Sprintf("%s.length", g.generateExpr(call.Args[0]))
                        }
                case "cap":
                        if len(call.Args) > 0 {
                                return fmt.Sprintf("%s.length", g.generateExpr(call.Args[0]))
                        }
                case "append":
                        return g.generateAppendCall(call.Args, call.Ellipsis)
                case "new":
                        if len(call.Args) > 0 {
                                return fmt.Sprintf("new (%s)()", g.generateExpr(call.Args[0]))
                        }
                case "panic":
                        if len(call.Args) > 0 {
                                return fmt.Sprintf("throw %s", g.generateExpr(call.Args[0]))
                        }
                case "String":
                        if len(call.Args) > 0 {
                                return fmt.Sprintf("String(%s)", g.generateExpr(call.Args[0]))
                        }
                case "errors":
                        // errors.New("msg") - handled elsewhere
                        break
                }
        }

        // Handle import-aliased calls like fmt.Println
        // Check if ident matches an import alias
        if ident, ok := call.Func.(*Ident); ok {
                if jsName, isImport := g.importForAlias(ident.Name); isImport {
                        // This is a package call, check method
                        return g.generatePackageCall(jsName, call.Args, call.Ellipsis)
                }
        }

        // Regular function call
        funcExpr := g.generateExpr(call.Func)
        args := g.generateArgs(call.Args, call.Ellipsis)
        return fmt.Sprintf("%s(%s)", funcExpr, args)
}

// generateSelectorCall generates JS for a method call like fmt.Println().
func (g *Generator) generateSelectorCall(sel *SelectorExpr, args []Expr, ellipsis bool) string {
        // Get the object part
        if ident, ok := sel.X.(*Ident); ok {
                switch ident.Name {
                case "fmt":
                        return g.generateFmtMethodCall(sel.Sel.Name, args, ellipsis)
                case "strings":
                        return g.generateStringsCall(sel.Sel.Name, args)
                case "strconv":
                        return g.generateStrconvCall(sel.Sel.Name, args)
                case "errors":
                        return g.generateErrorsCall(sel.Sel.Name, args)
                }

                // Check if the object is an imported package
                if jsPkg, isImport := g.importForAlias(ident.Name); isImport {
                        return fmt.Sprintf("%s.%s(%s)", jsPkg, sel.Sel.Name, g.generateArgs(args, ellipsis))
                }
        }

        // Regular method call
        obj := g.generateExpr(sel.X)
        method := sel.Sel.Name
        argsStr := g.generateArgs(args, ellipsis)

        // Handle special methods
        switch method {
        case "Sprintf":
                return g.generateSprintf(args)
        }

        return fmt.Sprintf("%s.%s(%s)", obj, method, argsStr)
}

// generateFmtCall generates JS for fmt.Println(), fmt.Printf(), etc.
func (g *Generator) generateFmtCall(args []Expr) string {
        if len(args) == 0 {
                return "console.log()"
        }
        return fmt.Sprintf("console.log(%s)", g.generateArgs(args, false))
}

// generateFmtMethodCall generates JS for fmt.Println, fmt.Sprintf, etc.
func (g *Generator) generateFmtMethodCall(method string, args []Expr, ellipsis bool) string {
        switch method {
        case "Println":
                if len(args) == 0 {
                        return "console.log()"
                }
                return fmt.Sprintf("console.log(%s)", g.generateArgs(args, false))
        case "Printf":
                if len(args) >= 1 {
                        // Convert format string to template literal
                        return g.generateSprintf(args)
                }
                return "console.log()"
        case "Sprintf":
                return g.generateSprintf(args)
        case "Fprintf", "Fprintln":
                // Simplified: just log
                if len(args) > 1 {
                        return fmt.Sprintf("console.log(%s)", g.generateArgs(args[1:], false))
                }
                return "console.log()"
        case "Errorf":
                return g.generateSprintf(args)
        default:
                return fmt.Sprintf("__gs.fmt.%s(%s)", method, g.generateArgs(args, ellipsis))
        }
}

// generateSprintf converts fmt.Sprintf("format", args...) to template literal.
func (g *Generator) generateSprintf(args []Expr) string {
        if len(args) == 0 {
                return `""`
        }
        formatExpr := g.generateExpr(args[0])
        rest := args[1:]

        // For simple cases, just use the format string directly
        // In a full implementation we'd parse the format string
        if len(rest) == 0 {
                return formatExpr
        }

        // Use a runtime sprintf helper
        restArgs := make([]string, len(rest))
        for i, a := range rest {
                restArgs[i] = g.generateExpr(a)
        }
        return fmt.Sprintf("__gs_sprintf(%s, %s)", formatExpr, strings.Join(restArgs, ", "))
}

// generateStringsCall generates JS for strings package functions.
func (g *Generator) generateStringsCall(method string, args []Expr) string {
        switch method {
        case "Contains":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.includes(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "HasPrefix":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.startsWith(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "HasSuffix":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.endsWith(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "ToUpper":
                if len(args) >= 1 {
                        return fmt.Sprintf("%s.toUpperCase()", g.generateExpr(args[0]))
                }
        case "ToLower":
                if len(args) >= 1 {
                        return fmt.Sprintf("%s.toLowerCase()", g.generateExpr(args[0]))
                }
        case "TrimSpace":
                if len(args) >= 1 {
                        return fmt.Sprintf("%s.trim()", g.generateExpr(args[0]))
                }
        case "Split":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.split(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "Join":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.join(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "Replace":
                if len(args) >= 4 {
                        return fmt.Sprintf("%s.replaceAll(%s, %s)", g.generateExpr(args[0]), g.generateExpr(args[1]), g.generateExpr(args[2]))
                }
        case "ReplaceAll":
                if len(args) >= 3 {
                        return fmt.Sprintf("%s.replaceAll(%s, %s)", g.generateExpr(args[0]), g.generateExpr(args[1]), g.generateExpr(args[2]))
                }
        case "Index":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.indexOf(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "LastIndex":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.lastIndexOf(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "Repeat":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s.repeat(%s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        }
        return fmt.Sprintf("__gs.strings.%s(%s)", method, g.generateArgs(args, false))
}

// generateStrconvCall generates JS for strconv package functions.
func (g *Generator) generateStrconvCall(method string, args []Expr) string {
        switch method {
        case "Itoa":
                if len(args) >= 1 {
                        return fmt.Sprintf("String(%s)", g.generateExpr(args[0]))
                }
        case "Atoi":
                if len(args) >= 1 {
                        return fmt.Sprintf("parseInt(%s, 10)", g.generateExpr(args[0]))
                }
        case "FormatInt":
                if len(args) >= 1 {
                        return fmt.Sprintf("String(%s)", g.generateExpr(args[0]))
                }
        case "ParseInt":
                if len(args) >= 2 {
                        return fmt.Sprintf("parseInt(%s, %s)", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        case "FormatFloat":
                if len(args) >= 1 {
                        return fmt.Sprintf("String(%s)", g.generateExpr(args[0]))
                }
        case "ParseFloat":
                if len(args) >= 1 {
                        return fmt.Sprintf("parseFloat(%s)", g.generateExpr(args[0]))
                }
        case "Quote":
                if len(args) >= 1 {
                        return fmt.Sprintf("JSON.stringify(%s)", g.generateExpr(args[0]))
                }
        }
        return fmt.Sprintf("__gs.strconv.%s(%s)", method, g.generateArgs(args, false))
}

// generateErrorsCall generates JS for errors package functions.
func (g *Generator) generateErrorsCall(method string, args []Expr) string {
        switch method {
        case "New":
                if len(args) >= 1 {
                        return fmt.Sprintf("new Error(%s)", g.generateExpr(args[0]))
                }
        case "Is":
                if len(args) >= 2 {
                        return fmt.Sprintf("%s instanceof %s", g.generateExpr(args[0]), g.generateExpr(args[1]))
                }
        }
        return fmt.Sprintf("__gs.errors.%s(%s)", method, g.generateArgs(args, false))
}

// generateMakeCall generates JS for make() built-in.
func (g *Generator) generateMakeCall(args []Expr) string {
        if len(args) == 0 {
                return "null"
        }
        typeExpr := g.generateExpr(args[0])

        // Check if it's a slice type
        if _, ok := args[0].(*SliceType); ok {
                if len(args) >= 2 {
                        sizeExpr := g.generateExpr(args[1])
                        return fmt.Sprintf("new Array(%s).fill(null)", sizeExpr)
                }
                return "[]"
        }

        // Check if it's a map type
        if _, ok := args[0].(*MapType); ok {
                return "{}"
        }

        // Check if it's a chan type
        if _, ok := args[0].(*ChanType); ok {
                return "__gs_chan()"
        }

        return fmt.Sprintf("new (%s)()", typeExpr)
}

// generateAppendCall generates JS for append() built-in.
func (g *Generator) generateAppendCall(args []Expr, ellipsis bool) string {
        if len(args) == 0 {
                return "[]"
        }
        sliceExpr := g.generateExpr(args[0])

        if ellipsis && len(args) == 2 {
                // append(slice, other...)
                other := g.generateExpr(args[1])
                return fmt.Sprintf("[...%s, ...%s]", sliceExpr, other)
        }

        // append(slice, elem1, elem2, ...)
        rest := args[1:]
        restExprs := make([]string, len(rest))
        for i, r := range rest {
                restExprs[i] = g.generateExpr(r)
        }
        return fmt.Sprintf("[...%s, %s]", sliceExpr, strings.Join(restExprs, ", "))
}

// generateArgs generates the comma-separated argument list for a function call.
func (g *Generator) generateArgs(args []Expr, ellipsis bool) string {
        parts := make([]string, len(args))
        for i, a := range args {
                parts[i] = g.generateExpr(a)
        }
        if ellipsis {
                return fmt.Sprintf("...%s", strings.Join(parts, ", "))
        }
        return strings.Join(parts, ", ")
}

// generatePackageCall generates JS for a call on an imported package.
func (g *Generator) generatePackageCall(pkgName string, args []Expr, ellipsis bool) string {
        // This is a simplified handler for import-based calls
        return fmt.Sprintf("%s(%s)", pkgName, g.generateArgs(args, ellipsis))
}

// generateSelectorExpr generates JavaScript for a selector expression (x.field).
func (g *Generator) generateSelectorExpr(sel *SelectorExpr) string {
        obj := g.generateExpr(sel.X)
        field := toCamelCase(sel.Sel.Name)

        // Special case: if accessing a Go package constant/function
        // that starts with uppercase, keep it as-is
        if isExported(sel.Sel.Name) {
                field = sel.Sel.Name
        }

        return fmt.Sprintf("%s.%s", obj, field)
}

// generateIndexExpr generates JavaScript for an index expression (x[i]).
func (g *Generator) generateIndexExpr(idx *IndexExpr) string {
        x := g.generateExpr(idx.X)
        index := g.generateExpr(idx.Index)
        return fmt.Sprintf("%s[%s]", x, index)
}

// generateBinaryExpr generates JavaScript for a binary expression.
func (g *Generator) generateBinaryExpr(expr *BinaryExpr) string {
        left := g.generateExpr(expr.X)
        right := g.generateExpr(expr.Y)

        switch expr.Op {
        case TOKEN_LAND:
                return fmt.Sprintf("(%s && %s)", left, right)
        case TOKEN_LOR:
                return fmt.Sprintf("(%s || %s)", left, right)
        case TOKEN_EQ:
                return fmt.Sprintf("(%s === %s)", left, right)
        case TOKEN_NEQ:
                return fmt.Sprintf("(%s !== %s)", left, right)
        default:
                op := TokenString(expr.Op)
                return fmt.Sprintf("(%s %s %s)", left, op, right)
        }
}

// generateUnaryExpr generates JavaScript for a unary expression.
func (g *Generator) generateUnaryExpr(expr *UnaryExpr) string {
        x := g.generateExpr(expr.X)
        switch expr.Op {
        case TOKEN_NOT:
                return fmt.Sprintf("(!%s)", x)
        case TOKEN_MINUS:
                return fmt.Sprintf("(-%s)", x)
        case TOKEN_PLUS:
                return fmt.Sprintf("(+%s)", x)
        case TOKEN_AMP:
                return x // &x → x (reference is implicit in JS)
        case TOKEN_ARROW:
                return fmt.Sprintf("__gs_recv(%s)", x) // <-ch
        case TOKEN_STAR:
                return x // *x → x (dereference is implicit in JS)
        default:
                return x
        }
}

// generateSliceExpr generates JavaScript for a slice expression (arr[low:high]).
func (g *Generator) generateSliceExpr(slice *SliceExpr) string {
        x := g.generateExpr(slice.X)

        if slice.Low == nil && slice.High == nil {
                // arr[:]
                return fmt.Sprintf("%s.slice()", x)
        }
        if slice.Low == nil {
                // arr[:high]
                return fmt.Sprintf("%s.slice(0, %s)", x, g.generateExpr(slice.High))
        }
        if slice.High == nil {
                // arr[low:]
                return fmt.Sprintf("%s.slice(%s)", x, g.generateExpr(slice.Low))
        }
        // arr[low:high]
        return fmt.Sprintf("%s.slice(%s, %s)", x, g.generateExpr(slice.Low), g.generateExpr(slice.High))
}

// --- Helper methods ---

// generateParamNames generates a comma-separated list of parameter names.
// skipReceiver indicates whether to skip the first parameter (for method calls where receiver is implicit).
func (g *Generator) generateParamNames(params *FieldList, skipReceiver bool) string {
        if params == nil {
                return ""
        }
        parts := make([]string, 0)
        startIdx := 0
        if skipReceiver && len(params.List) > 0 && len(params.List[0].Names) == 0 {
                // Skip unnamed first parameter (this happens in some edge cases)
        }

        for i := startIdx; i < len(params.List); i++ {
                field := params.List[i]
                if len(field.Names) > 0 {
                        for _, name := range field.Names {
                                parts = append(parts, name.Name)
                        }
                } else if i > startIdx || !skipReceiver {
                        // Unnamed parameter - generate a placeholder
                        parts = append(parts, fmt.Sprintf("_p%d", i))
                }
        }
        return strings.Join(parts, ", ")
}

// extractReceiverType extracts the type name from a method receiver field list.
func (g *Generator) extractReceiverType(recv *FieldList) string {
        if recv == nil || len(recv.List) == 0 {
                return ""
        }
        field := recv.List[0]
        switch t := field.Type.(type) {
        case *Ident:
                return t.Name
        case *StarExpr:
                if inner, ok := t.X.(*Ident); ok {
                        return inner.Name
                }
                return "unknown"
        case *SelectorExpr:
                return t.Sel.Name
        default:
                return "unknown"
        }
}

// importForAlias checks if the given name is an import alias and returns the JS package name.
func (g *Generator) importForAlias(alias string) (string, bool) {
        for goPath, jsName := range g.imports {
                // Check if the alias matches
                parts := strings.Split(goPath, "/")
                lastPart := parts[len(parts)-1]
                if alias == lastPart || alias == jsName {
                        return jsName, true
                }
        }
        return "", false
}

// defaultValue returns a JavaScript default value for a given Go type.
func (g *Generator) defaultValue(expr Expr) string {
        switch t := expr.(type) {
        case *Ident:
                switch t.Name {
                case "int", "int8", "int16", "int32", "int64",
                        "uint", "uint8", "uint16", "uint32", "uint64":
                        return "0"
                case "float32", "float64":
                        return "0"
                case "string":
                        return `""`
                case "bool":
                        return "false"
                case "error":
                        return "null"
                }
        case *SliceType:
                return "[]"
        case *MapType:
                return "{}"
        case *ArrayType:
                return "[]"
        case *StarExpr:
                return "null"
        }
        return "undefined"
}

// isErrorValue checks if a variable name looks like an error variable.
func (g *Generator) isErrorValue(name string) bool {
        return name == "err" || strings.HasSuffix(name, "Error") || strings.HasSuffix(name, "Err")
}

// generateMultiReturn handles destructuring of multi-return function calls.
func (g *Generator) generateMultiReturn(names []*Ident, call *CallExpr) string {
        // This is handled in generateAssignStmt for x, err := f()
        return g.generateExpr(call)
}

// --- Utility functions ---

// toCamelCase converts a Go exported identifier (PascalCase) to camelCase for JavaScript.
// Unexported identifiers are kept as-is.
func toCamelCase(name string) string {
        if len(name) == 0 {
                return name
        }
        // If the name starts with an uppercase letter (exported), convert to camelCase
        if name[0] >= 'A' && name[0] <= 'Z' {
                if len(name) == 1 {
                        return strings.ToLower(name)
                }
                return strings.ToLower(name[:1]) + name[1:]
        }
        return name
}

// isExported checks if a Go identifier is exported (starts with uppercase).
func isExported(name string) bool {
        if len(name) == 0 {
                return false
        }
        return name[0] >= 'A' && name[0] <= 'Z'
}
