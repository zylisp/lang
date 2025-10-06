package interpreter

import (
	"fmt"

	"zylisp/lang/sexpr"
)

// Env represents a lexical environment for variable bindings
type Env struct {
	bindings map[string]sexpr.SExpr
	parent   *Env
}

// NewEnv creates a new environment with an optional parent
func NewEnv(parent *Env) *Env {
	return &Env{
		bindings: make(map[string]sexpr.SExpr),
		parent:   parent,
	}
}

// Define binds a value to a name in this environment
func (e *Env) Define(name string, value sexpr.SExpr) {
	e.bindings[name] = value
}

// Set updates an existing binding, searching parent environments
func (e *Env) Set(name string, value sexpr.SExpr) error {
	if _, ok := e.bindings[name]; ok {
		e.bindings[name] = value
		return nil
	}

	if e.parent != nil {
		return e.parent.Set(name, value)
	}

	return fmt.Errorf("undefined variable: %s", name)
}

// Lookup finds a value by name, searching parent environments
func (e *Env) Lookup(name string) (sexpr.SExpr, error) {
	if value, ok := e.bindings[name]; ok {
		return value, nil
	}

	if e.parent != nil {
		return e.parent.Lookup(name)
	}

	return nil, fmt.Errorf("undefined variable: %s", name)
}

// Extend creates a child environment
func (e *Env) Extend() *Env {
	return NewEnv(e)
}
