package parser

import (
	"reflect"
	"testing"
)

func TestLexerSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			"empty",
			"",
			[]TokenType{EOF},
		},
		{
			"single number",
			"42",
			[]TokenType{NUMBER, EOF},
		},
		{
			"single symbol",
			"hello",
			[]TokenType{SYMBOL, EOF},
		},
		{
			"empty list",
			"()",
			[]TokenType{LPAREN, RPAREN, EOF},
		},
		{
			"simple list",
			"(+ 1 2)",
			[]TokenType{LPAREN, SYMBOL, NUMBER, NUMBER, RPAREN, EOF},
		},
		{
			"nested list",
			"(+ (* 2 3) 4)",
			[]TokenType{LPAREN, SYMBOL, LPAREN, SYMBOL, NUMBER, NUMBER,
				RPAREN, NUMBER, RPAREN, EOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("got %d tokens, want %d", len(tokens), len(tt.expected))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i] {
					t.Errorf("token %d: got %v, want %v",
						i, tok.Type, tt.expected[i])
				}
			}
		})
	}
}

func TestLexerTokenValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			"numbers",
			"42 -17 0",
			[]Token{
				{Type: NUMBER, Value: "42"},
				{Type: NUMBER, Value: "-17"},
				{Type: NUMBER, Value: "0"},
				{Type: EOF, Value: ""},
			},
		},
		{
			"symbols",
			"+ hello-world foo?",
			[]Token{
				{Type: SYMBOL, Value: "+"},
				{Type: SYMBOL, Value: "hello-world"},
				{Type: SYMBOL, Value: "foo?"},
				{Type: EOF, Value: ""},
			},
		},
		{
			"strings",
			`"hello" "world"`,
			[]Token{
				{Type: STRING, Value: "hello"},
				{Type: STRING, Value: "world"},
				{Type: EOF, Value: ""},
			},
		},
		{
			"booleans",
			"true false",
			[]Token{
				{Type: BOOL, Value: "true"},
				{Type: BOOL, Value: "false"},
				{Type: EOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("got %d tokens, want %d", len(tokens), len(tt.expected))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d type: got %v, want %v",
						i, tok.Type, tt.expected[i].Type)
				}
				if tok.Value != tt.expected[i].Value {
					t.Errorf("token %d value: got %q, want %q",
						i, tok.Value, tt.expected[i].Value)
				}
			}
		})
	}
}

func TestLexerComments(t *testing.T) {
	input := `
; This is a comment
(+ 1 2) ; inline comment
; another comment
42
`
	expected := []TokenType{LPAREN, SYMBOL, NUMBER, NUMBER, RPAREN, NUMBER, EOF}

	tokens, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var types []TokenType
	for _, tok := range tokens {
		types = append(types, tok.Type)
	}

	if !reflect.DeepEqual(types, expected) {
		t.Errorf("got %v, want %v", types, expected)
	}
}

func TestLexerStringEscapes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"quote\"here"`, `quote"here`},
		{`"backslash\\here"`, `backslash\here`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tokens) != 2 { // STRING + EOF
				t.Fatalf("got %d tokens, want 2", len(tokens))
			}

			if tokens[0].Value != tt.expected {
				t.Errorf("got %q, want %q", tokens[0].Value, tt.expected)
			}
		})
	}
}
