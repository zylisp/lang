package interpreter

import (
	"testing"

	"github.com/zylisp/lang/parser"
	"github.com/zylisp/lang/sexpr"
)

func testEvalWithPrimitives(t *testing.T, input string, expected sexpr.SExpr) {
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
	LoadPrimitives(env)

	result, err := Eval(expr, env)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}

	if result.String() != expected.String() {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestPrimAdd(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(+ 1 2)", 3},
		{"(+ 1 2 3 4)", 10},
		{"(+)", 0},
		{"(+ 42)", 42},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Number{Value: tt.expected})
		})
	}
}

func TestPrimSub(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(- 5 3)", 2},
		{"(- 10 3 2)", 5},
		{"(- 42)", -42},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Number{Value: tt.expected})
		})
	}
}

func TestPrimMul(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(* 2 3)", 6},
		{"(* 2 3 4)", 24},
		{"(*)", 1},
		{"(* 42)", 42},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Number{Value: tt.expected})
		})
	}
}

func TestPrimDiv(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(/ 6 2)", 3},
		{"(/ 24 3 2)", 4},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Number{Value: tt.expected})
		})
	}
}

func TestPrimComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"(= 1 1)", true},
		{"(= 1 2)", false},
		{"(< 1 2)", true},
		{"(< 2 1)", false},
		{"(> 2 1)", true},
		{"(> 1 2)", false},
		{"(<= 1 1)", true},
		{"(<= 1 2)", true},
		{"(<= 2 1)", false},
		{"(>= 1 1)", true},
		{"(>= 2 1)", true},
		{"(>= 1 2)", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Bool{Value: tt.expected})
		})
	}
}

func TestPrimList(t *testing.T) {
	input := "(list 1 2 3)"
	expected := sexpr.List{
		Elements: []sexpr.SExpr{
			sexpr.Number{Value: 1},
			sexpr.Number{Value: 2},
			sexpr.Number{Value: 3},
		},
	}

	testEvalWithPrimitives(t, input, expected)
}

func TestPrimCar(t *testing.T) {
	input := "(car (list 1 2 3))"
	expected := sexpr.Number{Value: 1}

	testEvalWithPrimitives(t, input, expected)
}

func TestPrimCdr(t *testing.T) {
	input := "(cdr (list 1 2 3))"
	expected := sexpr.List{
		Elements: []sexpr.SExpr{
			sexpr.Number{Value: 2},
			sexpr.Number{Value: 3},
		},
	}

	testEvalWithPrimitives(t, input, expected)
}

func TestPrimCons(t *testing.T) {
	input := "(cons 0 (list 1 2 3))"
	expected := sexpr.List{
		Elements: []sexpr.SExpr{
			sexpr.Number{Value: 0},
			sexpr.Number{Value: 1},
			sexpr.Number{Value: 2},
			sexpr.Number{Value: 3},
		},
	}

	testEvalWithPrimitives(t, input, expected)
}

func TestPrimTypePredicates(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"(number? 42)", true},
		{"(number? (quote x))", false},
		{"(symbol? (quote x))", true},
		{"(symbol? 42)", false},
		{"(list? (list 1 2))", true},
		{"(list? 42)", false},
		{"(null? (list))", true},
		{"(null? (list 1))", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Bool{Value: tt.expected})
		})
	}
}

func TestNestedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(+ (* 2 3) 4)", 10},
		{"(* (+ 1 2) (- 5 2))", 9},
		{"(/ (+ 10 6) (- 6 2))", 4},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testEvalWithPrimitives(t, tt.input, sexpr.Number{Value: tt.expected})
		})
	}
}

func TestUserDefinedFunctions(t *testing.T) {
	env := NewEnv(nil)
	LoadPrimitives(env)

	// Define the square function
	tokens1, _ := parser.Tokenize("(define square (lambda (x) (* x x)))")
	expr1, _ := parser.Read(tokens1)
	_, err := Eval(expr1, env)
	if err != nil {
		t.Fatalf("eval define error: %v", err)
	}

	// Call square function
	tokens2, _ := parser.Tokenize("(square 5)")
	expr2, _ := parser.Read(tokens2)
	result, err := Eval(expr2, env)
	if err != nil {
		t.Fatalf("eval call error: %v", err)
	}

	expected := sexpr.Number{Value: 25}
	if result.String() != expected.String() {
		t.Errorf("got %v, want %v", result, expected)
	}
}
