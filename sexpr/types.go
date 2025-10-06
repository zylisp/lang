package sexpr

import "fmt"

// SExpr is the base interface for all S-expression types
type SExpr interface {
	String() string
}

// Number represents an integer
type Number struct {
	Value int64
}

func (n Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// Symbol represents a name/identifier
type Symbol struct {
	Name string
}

func (s Symbol) String() string {
	return s.Name
}

// String represents a string literal
type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("%q", s.Value)
}

// Bool represents a boolean value
type Bool struct {
	Value bool
}

func (b Bool) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// Nil represents the empty value
type Nil struct{}

func (n Nil) String() string {
	return "nil"
}

// List represents a sequence of expressions
type List struct {
	Elements []SExpr
}

func (l List) String() string {
	if len(l.Elements) == 0 {
		return "()"
	}

	result := "("
	for i, elem := range l.Elements {
		if i > 0 {
			result += " "
		}
		result += elem.String()
	}
	result += ")"
	return result
}

// Func represents a user-defined function
type Func struct {
	Params []Symbol
	Body   SExpr
	Env    *Env // Will define in interpreter package
}

func (f Func) String() string {
	return "<function>"
}

// Primitive represents a built-in function
type Primitive struct {
	Name string
	Fn   func([]SExpr, *Env) (SExpr, error)
}

func (p Primitive) String() string {
	return fmt.Sprintf("<primitive:%s>", p.Name)
}

// Env is forward-declared here but implemented in interpreter
type Env interface {
	Define(name string, value SExpr)
	Lookup(name string) (SExpr, error)
}
