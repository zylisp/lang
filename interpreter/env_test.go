package interpreter

import (
	"testing"

	"github.com/zylisp/lang/sexpr"
)

func TestEnvDefineAndLookup(t *testing.T) {
	env := NewEnv(nil)

	// Define a variable
	env.Define("x", sexpr.Number{Value: 42})

	// Look it up
	value, err := env.Lookup("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	num, ok := value.(sexpr.Number)
	if !ok {
		t.Fatalf("expected Number, got %T", value)
	}

	if num.Value != 42 {
		t.Errorf("got %d, want 42", num.Value)
	}
}

func TestEnvUndefinedVariable(t *testing.T) {
	env := NewEnv(nil)

	_, err := env.Lookup("undefined")
	if err == nil {
		t.Error("expected error for undefined variable")
	}
}

func TestEnvParentLookup(t *testing.T) {
	parent := NewEnv(nil)
	parent.Define("x", sexpr.Number{Value: 42})

	child := NewEnv(parent)
	child.Define("y", sexpr.Number{Value: 17})

	// Child can see its own bindings
	value, err := child.Lookup("y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(sexpr.Number).Value != 17 {
		t.Errorf("got %v, want 17", value)
	}

	// Child can see parent bindings
	value, err = child.Lookup("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", value)
	}

	// Parent cannot see child bindings
	_, err = parent.Lookup("y")
	if err == nil {
		t.Error("parent should not see child bindings")
	}
}

func TestEnvShadowing(t *testing.T) {
	parent := NewEnv(nil)
	parent.Define("x", sexpr.Number{Value: 42})

	child := NewEnv(parent)
	child.Define("x", sexpr.Number{Value: 17})

	// Child sees its own binding
	value, err := child.Lookup("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(sexpr.Number).Value != 17 {
		t.Errorf("got %v, want 17", value)
	}

	// Parent still has original binding
	value, err = parent.Lookup("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", value)
	}
}

func TestEnvExtend(t *testing.T) {
	parent := NewEnv(nil)
	parent.Define("x", sexpr.Number{Value: 42})

	child := parent.Extend()

	// Child can see parent binding
	value, err := child.Lookup("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(sexpr.Number).Value != 42 {
		t.Errorf("got %v, want 42", value)
	}
}
