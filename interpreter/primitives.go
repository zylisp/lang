package interpreter

import (
	"fmt"

	"github.com/zylisp/lang/sexpr"
)

// LoadPrimitives adds all primitive functions to an environment
func LoadPrimitives(env *Env) {
	// Arithmetic
	env.Define("+", makePrimitive("+", primAdd))
	env.Define("-", makePrimitive("-", primSub))
	env.Define("*", makePrimitive("*", primMul))
	env.Define("/", makePrimitive("/", primDiv))

	// Comparison
	env.Define("=", makePrimitive("=", primEq))
	env.Define("<", makePrimitive("<", primLt))
	env.Define(">", makePrimitive(">", primGt))
	env.Define("<=", makePrimitive("<=", primLte))
	env.Define(">=", makePrimitive(">=", primGte))

	// List operations
	env.Define("list", makePrimitive("list", primList))
	env.Define("car", makePrimitive("car", primCar))
	env.Define("cdr", makePrimitive("cdr", primCdr))
	env.Define("cons", makePrimitive("cons", primCons))

	// Type predicates
	env.Define("number?", makePrimitive("number?", primIsNumber))
	env.Define("symbol?", makePrimitive("symbol?", primIsSymbol))
	env.Define("list?", makePrimitive("list?", primIsList))
	env.Define("null?", makePrimitive("null?", primIsNull))
}

func makePrimitive(name string, fn func([]sexpr.SExpr, *Env) (sexpr.SExpr, error)) sexpr.Primitive {
	return sexpr.Primitive{
		Name: name,
		Fn: func(args []sexpr.SExpr, envInterface interface{}) (sexpr.SExpr, error) {
			env := envInterface.(*Env)
			return fn(args, env)
		},
	}
}

// Arithmetic primitives

func primAdd(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) == 0 {
		return sexpr.Number{Value: 0}, nil
	}

	var sum int64
	for _, arg := range args {
		num, ok := arg.(sexpr.Number)
		if !ok {
			return nil, fmt.Errorf("+: expected number, got %v", arg)
		}
		sum += num.Value
	}

	return sexpr.Number{Value: sum}, nil
}

func primSub(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("-: requires at least 1 argument")
	}

	first, ok := args[0].(sexpr.Number)
	if !ok {
		return nil, fmt.Errorf("-: expected number, got %v", args[0])
	}

	if len(args) == 1 {
		return sexpr.Number{Value: -first.Value}, nil
	}

	result := first.Value
	for _, arg := range args[1:] {
		num, ok := arg.(sexpr.Number)
		if !ok {
			return nil, fmt.Errorf("-: expected number, got %v", arg)
		}
		result -= num.Value
	}

	return sexpr.Number{Value: result}, nil
}

func primMul(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) == 0 {
		return sexpr.Number{Value: 1}, nil
	}

	product := int64(1)
	for _, arg := range args {
		num, ok := arg.(sexpr.Number)
		if !ok {
			return nil, fmt.Errorf("*: expected number, got %v", arg)
		}
		product *= num.Value
	}

	return sexpr.Number{Value: product}, nil
}

func primDiv(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("/: requires at least 1 argument")
	}

	first, ok := args[0].(sexpr.Number)
	if !ok {
		return nil, fmt.Errorf("/: expected number, got %v", args[0])
	}

	if len(args) == 1 {
		if first.Value == 0 {
			return nil, fmt.Errorf("/: division by zero")
		}
		return sexpr.Number{Value: 1 / first.Value}, nil
	}

	result := first.Value
	for _, arg := range args[1:] {
		num, ok := arg.(sexpr.Number)
		if !ok {
			return nil, fmt.Errorf("/: expected number, got %v", arg)
		}
		if num.Value == 0 {
			return nil, fmt.Errorf("/: division by zero")
		}
		result /= num.Value
	}

	return sexpr.Number{Value: result}, nil
}

// Comparison primitives

func primEq(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("=: requires 2 arguments, got %d", len(args))
	}

	a, ok1 := args[0].(sexpr.Number)
	b, ok2 := args[1].(sexpr.Number)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("=: expected numbers")
	}

	return sexpr.Bool{Value: a.Value == b.Value}, nil
}

func primLt(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("<: requires 2 arguments, got %d", len(args))
	}

	a, ok1 := args[0].(sexpr.Number)
	b, ok2 := args[1].(sexpr.Number)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("<: expected numbers")
	}

	return sexpr.Bool{Value: a.Value < b.Value}, nil
}

func primGt(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(">: requires 2 arguments, got %d", len(args))
	}

	a, ok1 := args[0].(sexpr.Number)
	b, ok2 := args[1].(sexpr.Number)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf(">: expected numbers")
	}

	return sexpr.Bool{Value: a.Value > b.Value}, nil
}

func primLte(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("<=: requires 2 arguments, got %d", len(args))
	}

	a, ok1 := args[0].(sexpr.Number)
	b, ok2 := args[1].(sexpr.Number)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("<=: expected numbers")
	}

	return sexpr.Bool{Value: a.Value <= b.Value}, nil
}

func primGte(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(">=: requires 2 arguments, got %d", len(args))
	}

	a, ok1 := args[0].(sexpr.Number)
	b, ok2 := args[1].(sexpr.Number)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf(">=: expected numbers")
	}

	return sexpr.Bool{Value: a.Value >= b.Value}, nil
}

// List primitives

func primList(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	return sexpr.List{Elements: args}, nil
}

func primCar(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("car: requires 1 argument, got %d", len(args))
	}

	list, ok := args[0].(sexpr.List)
	if !ok {
		return nil, fmt.Errorf("car: expected list, got %v", args[0])
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("car: cannot take car of empty list")
	}

	return list.Elements[0], nil
}

func primCdr(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cdr: requires 1 argument, got %d", len(args))
	}

	list, ok := args[0].(sexpr.List)
	if !ok {
		return nil, fmt.Errorf("cdr: expected list, got %v", args[0])
	}

	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("cdr: cannot take cdr of empty list")
	}

	return sexpr.List{Elements: list.Elements[1:]}, nil
}

func primCons(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("cons: requires 2 arguments, got %d", len(args))
	}

	list, ok := args[1].(sexpr.List)
	if !ok {
		return nil, fmt.Errorf("cons: second argument must be a list, got %v", args[1])
	}

	elements := make([]sexpr.SExpr, 0, len(list.Elements)+1)
	elements = append(elements, args[0])
	elements = append(elements, list.Elements...)

	return sexpr.List{Elements: elements}, nil
}

// Type predicates

func primIsNumber(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("number?: requires 1 argument, got %d", len(args))
	}

	_, ok := args[0].(sexpr.Number)
	return sexpr.Bool{Value: ok}, nil
}

func primIsSymbol(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("symbol?: requires 1 argument, got %d", len(args))
	}

	_, ok := args[0].(sexpr.Symbol)
	return sexpr.Bool{Value: ok}, nil
}

func primIsList(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("list?: requires 1 argument, got %d", len(args))
	}

	_, ok := args[0].(sexpr.List)
	return sexpr.Bool{Value: ok}, nil
}

func primIsNull(args []sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("null?: requires 1 argument, got %d", len(args))
	}

	list, ok := args[0].(sexpr.List)
	if !ok {
		return sexpr.Bool{Value: false}, nil
	}

	return sexpr.Bool{Value: len(list.Elements) == 0}, nil
}
