// Package gscompiler implements a compiler for the GS (goscript) language.
// It takes .gs files with Go-syntax and compiles them to JavaScript for the browser.
package gscompiler

// TokenType represents the type of a lexical token in Go syntax.
type TokenType int

const (
	// Literals
	TOKEN_IDENT TokenType = iota
	TOKEN_INT
	TOKEN_FLOAT
	TOKEN_STRING
	TOKEN_RUNE
	TOKEN_RAW_STRING

	// Operators
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_STAR     // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %
	TOKEN_AMP      // &
	TOKEN_PIPE     // |
	TOKEN_CARET    // ^
	TOKEN_LSHIFT   // <<
	TOKEN_RSHIFT   // >>
	TOKEN_AMP_CARET // &^

	TOKEN_ASSIGN       // =
	TOKEN_SHORT_DECL   // :=
	TOKEN_PLUS_ASSIGN  // +=
	TOKEN_MINUS_ASSIGN // -=
	TOKEN_STAR_ASSIGN  // *=
	TOKEN_SLASH_ASSIGN // /=
	TOKEN_PERCENT_ASSIGN // %=
	TOKEN_AMP_ASSIGN   // &=
	TOKEN_PIPE_ASSIGN  // |=
	TOKEN_CARET_ASSIGN // ^=
	TOKEN_LSHIFT_ASSIGN // <<=
	TOKEN_RSHIFT_ASSIGN // >>=

	TOKEN_EQ  // ==
	TOKEN_NEQ // !=
	TOKEN_LT  // <
	TOKEN_GT  // >
	TOKEN_LEQ // <=
	TOKEN_GEQ // >=

	TOKEN_LAND // &&
	TOKEN_LOR  // ||
	TOKEN_NOT  // !

	TOKEN_ARROW     // <-
	TOKEN_DOT       // .
	TOKEN_COMMA     // ,
	TOKEN_SEMICOLON // ;
	TOKEN_COLON     // :
	TOKEN_ELLIPSIS  // ...

	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_LBRACE // {
	TOKEN_RBRACE // }
	TOKEN_LBRACK // [
	TOKEN_RBRACK // ]

	// Keywords (all Go keywords)
	TOKEN_BREAK
	TOKEN_CASE
	TOKEN_CHAN
	TOKEN_CONST
	TOKEN_CONTINUE
	TOKEN_DEFAULT
	TOKEN_DEFER
	TOKEN_ELSE
	TOKEN_FALLTHROUGH
	TOKEN_FOR
	TOKEN_FUNC
	TOKEN_GO
	TOKEN_GOTO
	TOKEN_IF
	TOKEN_IMPORT
	TOKEN_INTERFACE
	TOKEN_MAP
	TOKEN_PACKAGE
	TOKEN_RANGE
	TOKEN_RETURN
	TOKEN_SELECT
	TOKEN_STRUCT
	TOKEN_SWITCH
	TOKEN_TYPE
	TOKEN_VAR

	// Special
	TOKEN_NIL   // nil
	TOKEN_TRUE  // true
	TOKEN_FALSE // false

	// Control
	TOKEN_EOF
	TOKEN_NEWLINE
	TOKEN_COMMENT
	TOKEN_ILLEGAL
)

// Token represents a single lexical token with its type, literal value, and position.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// keywords maps Go keyword strings to their corresponding token types.
var keywords = map[string]TokenType{
	"break":       TOKEN_BREAK,
	"case":        TOKEN_CASE,
	"chan":        TOKEN_CHAN,
	"const":       TOKEN_CONST,
	"continue":    TOKEN_CONTINUE,
	"default":     TOKEN_DEFAULT,
	"defer":       TOKEN_DEFER,
	"else":        TOKEN_ELSE,
	"fallthrough": TOKEN_FALLTHROUGH,
	"for":         TOKEN_FOR,
	"func":        TOKEN_FUNC,
	"go":          TOKEN_GO,
	"goto":        TOKEN_GOTO,
	"if":          TOKEN_IF,
	"import":      TOKEN_IMPORT,
	"interface":   TOKEN_INTERFACE,
	"map":         TOKEN_MAP,
	"package":     TOKEN_PACKAGE,
	"range":       TOKEN_RANGE,
	"return":      TOKEN_RETURN,
	"select":      TOKEN_SELECT,
	"struct":      TOKEN_STRUCT,
	"switch":      TOKEN_SWITCH,
	"type":        TOKEN_TYPE,
	"var":         TOKEN_VAR,
	"nil":         TOKEN_NIL,
	"true":        TOKEN_TRUE,
	"false":       TOKEN_FALSE,
}

// LookupIdent checks whether the given identifier is a Go keyword.
// If it is, it returns the keyword token type; otherwise it returns TOKEN_IDENT.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENT
}

// TokenString returns a human-readable string representation of a token type.
func TokenString(tt TokenType) string {
	switch tt {
	case TOKEN_IDENT:
		return "IDENT"
	case TOKEN_INT:
		return "INT"
	case TOKEN_FLOAT:
		return "FLOAT"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_RUNE:
		return "RUNE"
	case TOKEN_RAW_STRING:
		return "RAW_STRING"
	case TOKEN_PLUS:
		return "+"
	case TOKEN_MINUS:
		return "-"
	case TOKEN_STAR:
		return "*"
	case TOKEN_SLASH:
		return "/"
	case TOKEN_PERCENT:
		return "%"
	case TOKEN_AMP:
		return "&"
	case TOKEN_PIPE:
		return "|"
	case TOKEN_CARET:
		return "^"
	case TOKEN_LSHIFT:
		return "<<"
	case TOKEN_RSHIFT:
		return ">>"
	case TOKEN_AMP_CARET:
		return "&^"
	case TOKEN_ASSIGN:
		return "="
	case TOKEN_SHORT_DECL:
		return ":="
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
	case TOKEN_EQ:
		return "=="
	case TOKEN_NEQ:
		return "!="
	case TOKEN_LT:
		return "<"
	case TOKEN_GT:
		return ">"
	case TOKEN_LEQ:
		return "<="
	case TOKEN_GEQ:
		return ">="
	case TOKEN_LAND:
		return "&&"
	case TOKEN_LOR:
		return "||"
	case TOKEN_NOT:
		return "!"
	case TOKEN_ARROW:
		return "<-"
	case TOKEN_DOT:
		return "."
	case TOKEN_COMMA:
		return ","
	case TOKEN_SEMICOLON:
		return ";"
	case TOKEN_COLON:
		return ":"
	case TOKEN_ELLIPSIS:
		return "..."
	case TOKEN_LPAREN:
		return "("
	case TOKEN_RPAREN:
		return ")"
	case TOKEN_LBRACE:
		return "{"
	case TOKEN_RBRACE:
		return "}"
	case TOKEN_LBRACK:
		return "["
	case TOKEN_RBRACK:
		return "]"
	case TOKEN_BREAK:
		return "break"
	case TOKEN_CASE:
		return "case"
	case TOKEN_CHAN:
		return "chan"
	case TOKEN_CONST:
		return "const"
	case TOKEN_CONTINUE:
		return "continue"
	case TOKEN_DEFAULT:
		return "default"
	case TOKEN_DEFER:
		return "defer"
	case TOKEN_ELSE:
		return "else"
	case TOKEN_FALLTHROUGH:
		return "fallthrough"
	case TOKEN_FOR:
		return "for"
	case TOKEN_FUNC:
		return "func"
	case TOKEN_GO:
		return "go"
	case TOKEN_GOTO:
		return "goto"
	case TOKEN_IF:
		return "if"
	case TOKEN_IMPORT:
		return "import"
	case TOKEN_INTERFACE:
		return "interface"
	case TOKEN_MAP:
		return "map"
	case TOKEN_PACKAGE:
		return "package"
	case TOKEN_RANGE:
		return "range"
	case TOKEN_RETURN:
		return "return"
	case TOKEN_SELECT:
		return "select"
	case TOKEN_SWITCH:
		return "switch"
	case TOKEN_TYPE:
		return "type"
	case TOKEN_VAR:
		return "var"
	case TOKEN_NIL:
		return "nil"
	case TOKEN_TRUE:
		return "true"
	case TOKEN_FALSE:
		return "false"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_NEWLINE:
		return "NEWLINE"
	case TOKEN_COMMENT:
		return "COMMENT"
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	default:
		return "UNKNOWN"
	}
}
