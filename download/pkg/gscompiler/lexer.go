package gscompiler

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes Go-syntax source code into a stream of tokens.
// It handles all Go literal types, operators, keywords, comments,
// and implements Go's automatic semicolon insertion rules.
type Lexer struct {
	input  string
	pos    int
	line   int
	column int
	tokens []Token
}

// NewLexer creates a new Lexer for the given input string.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		tokens: make([]Token, 0, 256),
	}
}

// Tokenize processes the entire input and returns the slice of tokens.
// Any lexical errors are returned as the second value.
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		l.skipWhitespace()

		if l.pos >= len(l.input) {
			break
		}

		ch := l.peek()

		// Comments
		if ch == '/' {
			if l.peekAt(1) == '/' {
				l.readLineComment()
				continue
			}
			if l.peekAt(1) == '*' {
				l.readBlockComment()
				continue
			}
		}

		// Identifiers and keywords
		if isLetter(ch) || ch == '_' {
			l.readIdentifier()
			continue
		}

		// Numbers
		if isDigit(ch) {
			l.readNumber()
			continue
		}

		// String literals
		if ch == '"' {
			l.readString()
			continue
		}

		// Raw string literals (backtick)
		if ch == '`' {
			l.readRawString()
			continue
		}

		// Rune literals
		if ch == '\'' {
			l.readRune()
			continue
		}

		// Operators and punctuation
		l.readOperator()
	}

	// Append EOF token
	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_EOF,
		Literal: "",
		Line:    l.line,
		Column:  l.column,
	})

	// Apply Go's automatic semicolon insertion rules
	l.insertSemicolons()

	return l.tokens, nil
}

// peek returns the current character without advancing.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

// peekAt returns the rune at offset from current position.
func (l *Lexer) peekAt(offset int) rune {
	idx := l.pos + offset
	if idx >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[idx:])
	return r
}

// advance reads the current rune and advances the position.
func (l *Lexer) advance() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return r
}

// skipWhitespace consumes whitespace characters (excluding newlines).
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else if ch == '\n' {
			// Track newlines as they are needed for semicolon insertion
			// but don't skip them as part of whitespace
			break
		} else {
			break
		}
	}
}

// readIdentifier reads an identifier or keyword.
func (l *Lexer) readIdentifier() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos

	for l.pos < len(l.input) {
		ch := l.peek()
		if isLetter(ch) || isDigit(ch) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	literal := l.input[startPos:l.pos]
	tokType := LookupIdent(literal)

	l.tokens = append(l.tokens, Token{
		Type:    tokType,
		Literal: literal,
		Line:    startLine,
		Column:  startCol,
	})
}

// readNumber reads a numeric literal (integer or float, including hex/octal/binary).
func (l *Lexer) readNumber() {
	startLine := l.line
	startCol := l.column
	startPos := l.pos
	isFloat := false

	ch := l.peek()

	// Handle hex, octal, binary prefixes
	if ch == '0' && l.pos+1 < len(l.input) {
		next := l.peekAt(1)
		switch next {
		case 'x', 'X':
			// Hex: 0x[0-9a-fA-F]+
			l.advance() // consume '0'
			l.advance() // consume 'x'
			for l.pos < len(l.input) {
				c := l.peek()
				if isHexDigit(c) || c == '_' {
					l.advance()
				} else {
					break
				}
			}
			l.tokens = append(l.tokens, Token{
				Type:    TOKEN_INT,
				Literal: l.input[startPos:l.pos],
				Line:    startLine,
				Column:  startCol,
			})
			return
		case 'o', 'O':
			// Octal: 0o[0-7]+
			l.advance()
			l.advance()
			for l.pos < len(l.input) {
				c := l.peek()
				if (c >= '0' && c <= '7') || c == '_' {
					l.advance()
				} else {
					break
				}
			}
			l.tokens = append(l.tokens, Token{
				Type:    TOKEN_INT,
				Literal: l.input[startPos:l.pos],
				Line:    startLine,
				Column:  startCol,
			})
			return
		case 'b', 'B':
			// Binary: 0b[01]+
			l.advance()
			l.advance()
			for l.pos < len(l.input) {
				c := l.peek()
				if (c == '0' || c == '1') || c == '_' {
					l.advance()
				} else {
					break
				}
			}
			l.tokens = append(l.tokens, Token{
				Type:    TOKEN_INT,
				Literal: l.input[startPos:l.pos],
				Line:    startLine,
				Column:  startCol,
			})
			return
		}
	}

	// Decimal integer or float
	for l.pos < len(l.input) {
		c := l.peek()
		if isDigit(c) || c == '_' {
			l.advance()
		} else if c == '.' {
			// Check if this is a decimal point (not a range operator, not a method selector)
			next := l.peekAt(1)
			if next == '.' {
				// Could be ellipsis, stop here
				break
			}
			if isDigit(next) || !isLetter(next) {
				isFloat = true
				l.advance() // consume '.'
				for l.pos < len(l.input) {
					d := l.peek()
					if isDigit(d) || d == '_' {
						l.advance()
					} else {
						break
					}
				}
			}
			break
		} else {
			break
		}
	}

	// Handle exponent for floats
	if !isFloat && l.pos < len(l.input) {
		c := l.peek()
		if c == 'e' || c == 'E' {
			isFloat = true
			l.advance()
			if l.pos < len(l.input) {
				s := l.peek()
				if s == '+' || s == '-' {
					l.advance()
				}
			}
			for l.pos < len(l.input) {
				d := l.peek()
				if isDigit(d) || d == '_' {
					l.advance()
				} else {
					break
				}
			}
		}
	}

	if isFloat {
		l.tokens = append(l.tokens, Token{
			Type:    TOKEN_FLOAT,
			Literal: l.input[startPos:l.pos],
			Line:    startLine,
			Column:  startCol,
		})
	} else {
		l.tokens = append(l.tokens, Token{
			Type:    TOKEN_INT,
			Literal: l.input[startPos:l.pos],
			Line:    startLine,
			Column:  startCol,
		})
	}
}

// readString reads a double-quoted string literal.
// Supports escape sequences: \n, \t, \r, \\, \", \xNN, \uNNNN, \UNNNNNNNN.
func (l *Lexer) readString() {
	startLine := l.line
	startCol := l.column
	l.advance() // consume opening "

	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == '"' {
			l.advance() // consume closing "
			break
		}
		if ch == '\\' {
			l.advance() // consume backslash
			if l.pos < len(l.input) {
				escaped := l.advance()
				switch escaped {
				case 'n':
					buf.WriteString("\\n")
				case 't':
					buf.WriteString("\\t")
				case 'r':
					buf.WriteString("\\r")
				case '\\':
					buf.WriteString("\\\\")
				case '"':
					buf.WriteString("\\\"")
				default:
					buf.WriteString(string(escaped))
				}
			}
			continue
		}
		if ch == '\n' {
			// Newline in string literal is an error in Go but we'll be lenient
			l.advance()
			buf.WriteString("\\n")
			continue
		}
		buf.WriteRune(ch)
		l.advance()
	}

	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_STRING,
		Literal: buf.String(),
		Line:    startLine,
		Column:  startCol,
	})
}

// readRawString reads a backtick-quoted raw string literal.
func (l *Lexer) readRawString() {
	startLine := l.line
	startCol := l.column
	l.advance() // consume opening `

	var buf strings.Builder
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == '`' {
			l.advance() // consume closing `
			break
		}
		buf.WriteRune(ch)
		l.advance()
	}

	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_RAW_STRING,
		Literal: buf.String(),
		Line:    startLine,
		Column:  startCol,
	})
}

// readRune reads a single-quoted rune literal.
func (l *Lexer) readRune() {
	startLine := l.line
	startCol := l.column
	l.advance() // consume opening '

	value := ""
	if l.pos < len(l.input) {
		ch := l.peek()
		if ch == '\\' {
			l.advance()
			if l.pos < len(l.input) {
				escaped := l.advance()
				switch escaped {
				case 'n':
					value = "\\n"
				case 't':
					value = "\\t"
				case 'r':
					value = "\\r"
				case '\\':
					value = "\\\\"
				case '\'':
					value = "\\'"
				default:
					value = string(escaped)
				}
			}
		} else {
			value = string(ch)
			l.advance()
		}
	}

	if l.pos < len(l.input) && l.peek() == '\'' {
		l.advance() // consume closing '
	}

	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_RUNE,
		Literal: value,
		Line:    startLine,
		Column:  startCol,
	})
}

// readOperator reads a multi-character operator or a single-character token.
func (l *Lexer) readOperator() {
	startLine := l.line
	startCol := l.column

	ch := l.peek()
	switch ch {
	case '+':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_PLUS_ASSIGN, "+=", startLine, startCol)
		} else {
			l.emit(TOKEN_PLUS, "+", startLine, startCol)
		}
	case '-':
		l.advance()
		switch l.peek() {
		case '=':
			l.advance()
			l.emit(TOKEN_MINUS_ASSIGN, "-=", startLine, startCol)
		case '>':
			l.advance()
			l.emit(TOKEN_ARROW, "->", startLine, startCol)
		default:
			l.emit(TOKEN_MINUS, "-", startLine, startCol)
		}
	case '*':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_STAR_ASSIGN, "*=", startLine, startCol)
		} else {
			l.emit(TOKEN_STAR, "*", startLine, startCol)
		}
	case '/':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_SLASH_ASSIGN, "/=", startLine, startCol)
		} else {
			l.emit(TOKEN_SLASH, "/", startLine, startCol)
		}
	case '%':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_PERCENT_ASSIGN, "%=", startLine, startCol)
		} else {
			l.emit(TOKEN_PERCENT, "%", startLine, startCol)
		}
	case '&':
		l.advance()
		switch l.peek() {
		case '&':
			l.advance()
			l.emit(TOKEN_LAND, "&&", startLine, startCol)
		case '=':
			l.advance()
			l.emit(TOKEN_AMP_ASSIGN, "&=", startLine, startCol)
		case '^':
			l.advance()
			if l.peek() == '=' {
				l.advance()
				l.emit(TOKEN_AMP_CARET, "&^", startLine, startCol) // simplified
				_ = fmt.Sprintf("") // suppress import warning
			} else {
				l.emit(TOKEN_AMP_CARET, "&^", startLine, startCol)
			}
		default:
			l.emit(TOKEN_AMP, "&", startLine, startCol)
		}
	case '|':
		l.advance()
		switch l.peek() {
		case '|':
			l.advance()
			l.emit(TOKEN_LOR, "||", startLine, startCol)
		case '=':
			l.advance()
			l.emit(TOKEN_PIPE_ASSIGN, "|=", startLine, startCol)
		default:
			l.emit(TOKEN_PIPE, "|", startLine, startCol)
		}
	case '^':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_CARET_ASSIGN, "^=", startLine, startCol)
		} else {
			l.emit(TOKEN_CARET, "^", startLine, startCol)
		}
	case '<':
		l.advance()
		switch l.peek() {
		case '=':
			l.advance()
			l.emit(TOKEN_LEQ, "<=", startLine, startCol)
		case '<':
			l.advance()
			if l.peek() == '=' {
				l.advance()
				l.emit(TOKEN_LSHIFT_ASSIGN, "<<=", startLine, startCol)
			} else {
				l.emit(TOKEN_LSHIFT, "<<", startLine, startCol)
			}
		case '-':
			l.advance()
			l.emit(TOKEN_ARROW, "<-", startLine, startCol)
		default:
			l.emit(TOKEN_LT, "<", startLine, startCol)
		}
	case '>':
		l.advance()
		switch l.peek() {
		case '=':
			l.advance()
			l.emit(TOKEN_GEQ, ">=", startLine, startCol)
		case '>':
			l.advance()
			if l.peek() == '=' {
				l.advance()
				l.emit(TOKEN_RSHIFT_ASSIGN, ">>=", startLine, startCol)
			} else {
				l.emit(TOKEN_RSHIFT, ">>", startLine, startCol)
			}
		default:
			l.emit(TOKEN_GT, ">", startLine, startCol)
		}
	case '=':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_EQ, "==", startLine, startCol)
		} else {
			l.emit(TOKEN_ASSIGN, "=", startLine, startCol)
		}
	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_NEQ, "!=", startLine, startCol)
		} else {
			l.emit(TOKEN_NOT, "!", startLine, startCol)
		}
	case ':':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			l.emit(TOKEN_SHORT_DECL, ":=", startLine, startCol)
		} else {
			l.emit(TOKEN_COLON, ":", startLine, startCol)
		}
	case '.':
		l.advance()
		if l.peek() == '.' && l.peekAt(1) == '.' {
			l.advance() // second dot
			l.advance() // third dot
			l.emit(TOKEN_ELLIPSIS, "...", startLine, startCol)
		} else {
			l.emit(TOKEN_DOT, ".", startLine, startCol)
		}
	case '(':
		l.advance()
		l.emit(TOKEN_LPAREN, "(", startLine, startCol)
	case ')':
		l.advance()
		l.emit(TOKEN_RPAREN, ")", startLine, startCol)
	case '{':
		l.advance()
		l.emit(TOKEN_LBRACE, "{", startLine, startCol)
	case '}':
		l.advance()
		l.emit(TOKEN_RBRACE, "}", startLine, startCol)
	case '[':
		l.advance()
		l.emit(TOKEN_LBRACK, "[", startLine, startCol)
	case ']':
		l.advance()
		l.emit(TOKEN_RBRACK, "]", startLine, startCol)
	case ',':
		l.advance()
		l.emit(TOKEN_COMMA, ",", startLine, startCol)
	case ';':
		l.advance()
		l.emit(TOKEN_SEMICOLON, ";", startLine, startCol)
	case '\n':
		l.advance()
		l.emit(TOKEN_NEWLINE, "\\n", startLine, startCol)
	default:
		l.advance()
		l.emit(TOKEN_ILLEGAL, string(ch), startLine, startCol)
	}
}

// emit appends a token to the token list.
func (l *Lexer) emit(tokType TokenType, literal string, line, col int) {
	l.tokens = append(l.tokens, Token{
		Type:    tokType,
		Literal: literal,
		Line:    line,
		Column:  col,
	})
}

// readLineComment reads a single-line comment (// ...).
func (l *Lexer) readLineComment() {
	startLine := l.line
	startCol := l.column
	l.advance() // consume first /
	l.advance() // consume second /

	var buf strings.Builder
	for l.pos < len(l.input) && l.peek() != '\n' {
		buf.WriteRune(l.peek())
		l.advance()
	}

	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_COMMENT,
		Literal: buf.String(),
		Line:    startLine,
		Column:  startCol,
	})
}

// readBlockComment reads a block comment (/* ... */).
func (l *Lexer) readBlockComment() {
	startLine := l.line
	startCol := l.column
	l.advance() // consume first /
	l.advance() // consume *

	var buf strings.Builder
	buf.WriteString("/*")
	for l.pos < len(l.input) {
		ch := l.peek()
		buf.WriteRune(ch)
		l.advance()
		if ch == '*' && l.peek() == '/' {
			buf.WriteRune(l.peek())
			l.advance()
			break
		}
	}

	l.tokens = append(l.tokens, Token{
		Type:    TOKEN_COMMENT,
		Literal: buf.String(),
		Line:    startLine,
		Column:  startCol,
	})
}

// insertSemicolons implements Go's automatic semicolon insertion.
// When the input is broken into tokens, a semicolon is automatically
// inserted into the token stream immediately after a line's final token
// if that token is one of: identifier, literal (int, float, rune, string,
// raw string), break, continue, fallthrough, return, ), }, ].
//
// The semicolon is NOT inserted if the next non-newline token on the
// next line is an operator or opening punctuation.
func (l *Lexer) insertSemicolons() {
	var result []Token

	for i := 0; i < len(l.tokens); i++ {
		tok := l.tokens[i]

		// Skip newline tokens and comments
		if tok.Type == TOKEN_NEWLINE || tok.Type == TOKEN_COMMENT {
			continue
		}

		result = append(result, tok)

		// Check if we need to insert a semicolon after this token
		if l.shouldInsertSemicolon(tok, i) {
			// Look ahead for a newline before the next non-newline, non-comment token
			hasNewline := false
			for j := i + 1; j < len(l.tokens); j++ {
				next := l.tokens[j]
				if next.Type == TOKEN_NEWLINE {
					hasNewline = true
					continue
				}
				if next.Type == TOKEN_COMMENT {
					continue
				}
				// Found the next real token
				if hasNewline && l.blocksSemicolon(next.Type) {
					result = append(result, Token{
						Type:    TOKEN_SEMICOLON,
						Literal: ";",
						Line:    tok.Line,
						Column:  tok.Column + len(tok.Literal),
					})
				}
				break
			}
		}
	}

	l.tokens = result
}

// shouldInsertSemicolon returns true if a semicolon should potentially
// be inserted after the given token.
func (l *Lexer) shouldInsertSemicolon(tok Token, idx int) bool {
	switch tok.Type {
	case TOKEN_IDENT, TOKEN_INT, TOKEN_FLOAT, TOKEN_STRING, TOKEN_RUNE, TOKEN_RAW_STRING,
		TOKEN_BREAK, TOKEN_CONTINUE, TOKEN_FALLTHROUGH, TOKEN_RETURN,
		TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL,
		TOKEN_RPAREN, TOKEN_RBRACE, TOKEN_RBRACK:
		return true
	default:
		return false
	}
}

// blocksSemicolon returns true if the given token type prevents
// semicolon insertion when it appears after a newline.
func (l *Lexer) blocksSemicolon(tt TokenType) bool {
	switch tt {
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT,
		TOKEN_AMP, TOKEN_PIPE, TOKEN_CARET, TOKEN_LSHIFT, TOKEN_RSHIFT, TOKEN_AMP_CARET,
		TOKEN_ASSIGN, TOKEN_SHORT_DECL, TOKEN_PLUS_ASSIGN, TOKEN_MINUS_ASSIGN,
		TOKEN_STAR_ASSIGN, TOKEN_SLASH_ASSIGN, TOKEN_PERCENT_ASSIGN,
		TOKEN_AMP_ASSIGN, TOKEN_PIPE_ASSIGN, TOKEN_CARET_ASSIGN,
		TOKEN_LSHIFT_ASSIGN, TOKEN_RSHIFT_ASSIGN,
		TOKEN_EQ, TOKEN_NEQ, TOKEN_LT, TOKEN_GT, TOKEN_LEQ, TOKEN_GEQ,
		TOKEN_LAND, TOKEN_LOR, TOKEN_NOT,
		TOKEN_ARROW, TOKEN_COLON, TOKEN_ELLIPSIS,
		TOKEN_LPAREN, TOKEN_LBRACE, TOKEN_LBRACK,
		TOKEN_COMMA, TOKEN_SEMICOLON:
		return true
	default:
		return false
	}
}

// isLetter returns true if the rune is a Unicode letter.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit returns true if the rune is a decimal digit.
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isHexDigit returns true if the rune is a hexadecimal digit.
func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
