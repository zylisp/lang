package parser

import (
	"reflect"
	"testing"

	"zylisp/lang/sexpr"
)

func TestReaderNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected sexpr.SExpr
	}{
		{"42", sexpr.Number{Value: 42}},
		{"-17", sexpr.Number{Value: -17}},
		{"0", sexpr.Number{Value: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize error: %v", err)
			}

			result, err := Read(tokens)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReaderSymbols(t *testing.T) {
	tests := []struct {
		input    string
		expected sexpr.SExpr
	}{
		{"x", sexpr.Symbol{Name: "x"}},
		{"+", sexpr.Symbol{Name: "+"}},
		{"hello-world", sexpr.Symbol{Name: "hello-world"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize error: %v", err)
			}

			result, err := Read(tokens)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReaderLists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected sexpr.SExpr
	}{
		{
			"empty list",
			"()",
			sexpr.List{Elements: []sexpr.SExpr{}},
		},
		{
			"single element",
			"(42)",
			sexpr.List{Elements: []sexpr.SExpr{
				sexpr.Number{Value: 42},
			}},
		},
		{
			"simple list",
			"(+ 1 2)",
			sexpr.List{Elements: []sexpr.SExpr{
				sexpr.Symbol{Name: "+"},
				sexpr.Number{Value: 1},
				sexpr.Number{Value: 2},
			}},
		},
		{
			"nested list",
			"(+ (* 2 3) 4)",
			sexpr.List{Elements: []sexpr.SExpr{
				sexpr.Symbol{Name: "+"},
				sexpr.List{Elements: []sexpr.SExpr{
					sexpr.Symbol{Name: "*"},
					sexpr.Number{Value: 2},
					sexpr.Number{Value: 3},
				}},
				sexpr.Number{Value: 4},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize error: %v", err)
			}

			result, err := Read(tokens)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReaderBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected sexpr.SExpr
	}{
		{"true", sexpr.Bool{Value: true}},
		{"false", sexpr.Bool{Value: false}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize error: %v", err)
			}

			result, err := Read(tokens)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReaderStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected sexpr.SExpr
	}{
		{`"hello"`, sexpr.String{Value: "hello"}},
		{`"hello world"`, sexpr.String{Value: "hello world"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("tokenize error: %v", err)
			}

			result, err := Read(tokens)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReaderErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unclosed list", "(+ 1 2"},
		{"extra closing paren", "(+ 1 2))"},
		{"just closing paren", ")"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize(tt.input)
			if err != nil {
				// Lexer error is fine for some test cases
				return
			}

			_, err = Read(tokens)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
