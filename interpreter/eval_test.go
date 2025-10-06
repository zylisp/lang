package interpreter

import (
	"testing"

	"zylisp/lang/parser"
	"zylisp/lang/sexpr"
)

func testEval(t *testing.T, input string, expected sexpr.SExpr) {
	t.Helper()

	tokens, err := parser.Tokenize(input)
	if err != nil {
		t.Fatalf("tokenize error: %v", err)
	}

	expr, err := parser.Read(tokens)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	env := NewEnv(nil)
	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	if result.String() != expected.String() {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestEvalSelfEvaluating(t *testing.T) {
	tests := []struct {
		input    string
		expected sexpr.SExpr
	}{
		{"42", sexpr.Number{Value: 42}},
		{`"hello"`, sexpr.String{Value: "hello"}},
		{"true", sexpr.Bool{Value: true}},
		{"false", sexpr.Bool{Value: false}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEval(t, tt.input, tt.expected)
		})
	}
}

func TestEvalDefine(t *testing.T) {
	tokens, _ := parser.Tokenize("(define x 42)")
	expr, _ := parser.Read(tokens)

	env := NewEnv(nil)
	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	// Should return the value
	if result.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", result)
	}

	// Should be defined in environment
	value, err := env.Lookup("x")
	if err != nil {
		t.Fatalf("lookup error: %v", err)
	}
	if value.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", value)
	}
}

func TestEvalSymbolLookup(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", sexpr.Number{Value: 42})

	tokens, _ := parser.Tokenize("x")
	expr, _ := parser.Read(tokens)

	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	if result.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", result)
	}
}

func TestEvalLambda(t *testing.T) {
	tokens, _ := parser.Tokenize("(lambda (x) x)")
	expr, _ := parser.Read(tokens)

	env := NewEnv(nil)
	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	_, ok := result.(sexpr.Func)
	if !ok {
		t.Errorf("expected Func, got %T", result)
	}
}

func TestEvalIf(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(if true 1 2)", 1},
		{"(if false 1 2)", 2},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens, _ := parser.Tokenize(tt.input)
			expr, _ := parser.Read(tokens)

			env := NewEnv(nil)
			result, err := Eval(expr, env)
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}

			if result.(sexpr.Number).Value != tt.expected {
				t.Errorf("got %v, want %d", result, tt.expected)
			}
		})
	}
}

func TestEvalQuote(t *testing.T) {
	tokens, _ := parser.Tokenize("(quote (+ 1 2))")
	expr, _ := parser.Read(tokens)

	env := NewEnv(nil)
	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	list, ok := result.(sexpr.List)
	if !ok {
		t.Fatalf("expected List, got %T", result)
	}

	if len(list.Elements) != 3 {
		t.Errorf("got %d elements, want 3", len(list.Elements))
	}
}
