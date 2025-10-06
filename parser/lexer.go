package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	LPAREN TokenType = iota
	RPAREN
	NUMBER
	SYMBOL
	STRING
	BOOL
	EOF
	ILLEGAL
)

func (tt TokenType) String() string {
	switch tt {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case NUMBER:
		return "NUMBER"
	case SYMBOL:
		return "SYMBOL"
	case STRING:
		return "STRING"
	case BOOL:
		return "BOOL"
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	default:
		return "UNKNOWN"
	}
}

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", t.Type, t.Value)
}

// Lexer tokenizes Zylisp source code
type Lexer struct {
	input  string
	pos    int // current position
	line   int // current line
	col    int // current column
	tokens []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

// Tokenize returns all tokens from the input
func Tokenize(input string) ([]Token, error) {
	lexer := NewLexer(input)
	return lexer.Tokenize()
}

// Tokenize produces all tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	for {
		tok := l.nextToken()
		l.tokens = append(l.tokens, tok)

		if tok.Type == EOF {
			break
		}

		if tok.Type == ILLEGAL {
			return nil, fmt.Errorf("illegal token at line %d, col %d: %q",
				tok.Line, tok.Col, tok.Value)
		}
	}

	return l.tokens, nil
}

// nextToken returns the next token
func (l *Lexer) nextToken() Token {
	l.skipWhitespaceAndComments()

	if l.isAtEnd() {
		return l.makeToken(EOF, "")
	}

	ch := l.peek()

	switch ch {
	case '(':
		return l.makeSingleCharToken(LPAREN)
	case ')':
		return l.makeSingleCharToken(RPAREN)
	case '"':
		return l.scanString()
	}

	if isDigit(ch) || (ch == '-' && l.peekNext() != 0 && isDigit(l.peekNext())) {
		return l.scanNumber()
	}

	if isSymbolStart(ch) {
		return l.scanSymbol()
	}

	return l.makeToken(ILLEGAL, string(ch))
}

// skipWhitespaceAndComments skips whitespace and comments
func (l *Lexer) skipWhitespaceAndComments() {
	for !l.isAtEnd() {
		ch := l.peek()

		if ch == ';' {
			// Skip comment to end of line
			for !l.isAtEnd() && l.peek() != '\n' {
				l.advance()
			}
			continue
		}

		if isWhitespace(ch) {
			l.advance()
			continue
		}

		break
	}
}

// scanNumber scans a number token
func (l *Lexer) scanNumber() Token {
	start := l.pos
	startCol := l.col

	if l.peek() == '-' {
		l.advance()
	}

	for !l.isAtEnd() && isDigit(l.peek()) {
		l.advance()
	}

	value := l.input[start:l.pos]
	return Token{Type: NUMBER, Value: value, Line: l.line, Col: startCol}
}

// scanSymbol scans a symbol token
func (l *Lexer) scanSymbol() Token {
	start := l.pos
	startCol := l.col

	for !l.isAtEnd() && isSymbolChar(l.peek()) {
		l.advance()
	}

	value := l.input[start:l.pos]

	// Check for boolean literals
	if value == "true" || value == "false" {
		return Token{Type: BOOL, Value: value, Line: l.line, Col: startCol}
	}

	return Token{Type: SYMBOL, Value: value, Line: l.line, Col: startCol}
}

// scanString scans a string token
func (l *Lexer) scanString() Token {
	startCol := l.col
	l.advance() // consume opening quote

	var value strings.Builder

	for !l.isAtEnd() && l.peek() != '"' {
		ch := l.peek()

		if ch == '\\' {
			l.advance()
			if l.isAtEnd() {
				return l.makeToken(ILLEGAL, "unterminated string")
			}

			// Handle escape sequences
			escaped := l.peek()
			switch escaped {
			case 'n':
				value.WriteByte('\n')
			case 't':
				value.WriteByte('\t')
			case 'r':
				value.WriteByte('\r')
			case '"':
				value.WriteByte('"')
			case '\\':
				value.WriteByte('\\')
			default:
				value.WriteByte(escaped)
			}
			l.advance()
		} else {
			value.WriteByte(ch)
			l.advance()
		}
	}

	if l.isAtEnd() {
		return l.makeToken(ILLEGAL, "unterminated string")
	}

	l.advance() // consume closing quote

	return Token{Type: STRING, Value: value.String(), Line: l.line, Col: startCol}
}

// Helper functions

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) advance() byte {
	if l.isAtEnd() {
		return 0
	}

	ch := l.input[l.pos]
	l.pos++

	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return ch
}

func (l *Lexer) isAtEnd() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer) makeToken(typ TokenType, value string) Token {
	return Token{Type: typ, Value: value, Line: l.line, Col: l.col}
}

func (l *Lexer) makeSingleCharToken(typ TokenType) Token {
	ch := l.peek()
	l.advance()
	return Token{Type: typ, Value: string(ch), Line: l.line, Col: l.col - 1}
}

// Character classification

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSymbolStart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || isSymbolSpecial(ch)
}

func isSymbolChar(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || isDigit(ch) || isSymbolSpecial(ch)
}

func isSymbolSpecial(ch byte) bool {
	return strings.ContainsRune("+-*/<>=!?&|%$_", rune(ch))
}
