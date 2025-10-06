package parser

import (
	"fmt"
	"strconv"

	"github.com/zylisp/lang/sexpr"
)

// Reader parses tokens into S-expressions
type Reader struct {
	tokens []Token
	pos    int
}

// NewReader creates a new reader for the given tokens
func NewReader(tokens []Token) *Reader {
	return &Reader{tokens: tokens, pos: 0}
}

// Read parses tokens into an S-expression
func Read(tokens []Token) (sexpr.SExpr, error) {
	reader := NewReader(tokens)
	expr, err := reader.readExpr()
	if err != nil {
		return nil, err
	}

	// Check for extra tokens after the expression (excluding EOF)
	if !reader.isAtEnd() && reader.peek().Type != EOF {
		tok := reader.peek()
		return nil, fmt.Errorf("unexpected token after expression at line %d, col %d: %v",
			tok.Line, tok.Col, tok.Type)
	}

	return expr, nil
}

// readExpr reads a single expression
func (r *Reader) readExpr() (sexpr.SExpr, error) {
	if r.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of input")
	}

	tok := r.peek()

	switch tok.Type {
	case LPAREN:
		return r.readList()
	case NUMBER:
		return r.readNumber()
	case SYMBOL:
		return r.readSymbol()
	case STRING:
		return r.readString()
	case BOOL:
		return r.readBool()
	case RPAREN:
		return nil, fmt.Errorf("unexpected closing paren at line %d, col %d",
			tok.Line, tok.Col)
	case EOF:
		return nil, fmt.Errorf("unexpected end of file")
	default:
		return nil, fmt.Errorf("unexpected token %v at line %d, col %d",
			tok.Type, tok.Line, tok.Col)
	}
}

// readList reads a list expression
func (r *Reader) readList() (sexpr.SExpr, error) {
	r.advance() // consume LPAREN

	elements := []sexpr.SExpr{}

	for !r.isAtEnd() && r.peek().Type != RPAREN {
		expr, err := r.readExpr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, expr)
	}

	if r.isAtEnd() {
		return nil, fmt.Errorf("unclosed list")
	}

	r.advance() // consume RPAREN

	return sexpr.List{Elements: elements}, nil
}

// readNumber reads a number expression
func (r *Reader) readNumber() (sexpr.SExpr, error) {
	tok := r.advance()

	value, err := strconv.ParseInt(tok.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %q at line %d, col %d: %v",
			tok.Value, tok.Line, tok.Col, err)
	}

	return sexpr.Number{Value: value}, nil
}

// readSymbol reads a symbol expression
func (r *Reader) readSymbol() (sexpr.SExpr, error) {
	tok := r.advance()
	return sexpr.Symbol{Name: tok.Value}, nil
}

// readString reads a string expression
func (r *Reader) readString() (sexpr.SExpr, error) {
	tok := r.advance()
	return sexpr.String{Value: tok.Value}, nil
}

// readBool reads a boolean expression
func (r *Reader) readBool() (sexpr.SExpr, error) {
	tok := r.advance()
	value := tok.Value == "true"
	return sexpr.Bool{Value: value}, nil
}

// Helper functions

func (r *Reader) peek() Token {
	if r.isAtEnd() {
		return Token{Type: EOF}
	}
	return r.tokens[r.pos]
}

func (r *Reader) advance() Token {
	if r.isAtEnd() {
		return Token{Type: EOF}
	}
	tok := r.tokens[r.pos]
	r.pos++
	return tok
}

func (r *Reader) isAtEnd() bool {
	return r.pos >= len(r.tokens) || r.tokens[r.pos].Type == EOF
}
