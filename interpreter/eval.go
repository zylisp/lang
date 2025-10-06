package interpreter

import (
	"fmt"

	"zylisp/lang/sexpr"
)

// Eval evaluates an S-expression in an environment
func Eval(expr sexpr.SExpr, env *Env) (sexpr.SExpr, error) {
	switch e := expr.(type) {

	// Self-evaluating types
	case sexpr.Number:
		return e, nil
	case sexpr.String:
		return e, nil
	case sexpr.Bool:
		return e, nil
	case sexpr.Nil:
		return e, nil

	// Symbol lookup
	case sexpr.Symbol:
		return env.Lookup(e.Name)

	// List evaluation
	case sexpr.List:
		return evalList(e, env)

	default:
		return nil, fmt.Errorf("cannot evaluate: %v", expr)
	}
}

// evalList evaluates a list expression
func evalList(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	if len(list.Elements) == 0 {
		return sexpr.Nil{}, nil
	}

	first := list.Elements[0]

	// Check for special forms
	if sym, ok := first.(sexpr.Symbol); ok {
		switch sym.Name {
		case "define":
			return evalDefine(list, env)
		case "lambda":
			return evalLambda(list, env)
		case "if":
			return evalIf(list, env)
		case "quote":
			return evalQuote(list, env)
		}
	}

	// Function application
	return evalApply(list, env)
}

// evalDefine handles (define name value)
func evalDefine(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	if len(list.Elements) != 3 {
		return nil, fmt.Errorf("define requires 2 arguments, got %d",
			len(list.Elements)-1)
	}

	name, ok := list.Elements[1].(sexpr.Symbol)
	if !ok {
		return nil, fmt.Errorf("define: first argument must be a symbol")
	}

	value, err := Eval(list.Elements[2], env)
	if err != nil {
		return nil, err
	}

	env.Define(name.Name, value)
	return value, nil
}

// evalLambda handles (lambda (params...) body)
func evalLambda(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	if len(list.Elements) != 3 {
		return nil, fmt.Errorf("lambda requires 2 arguments, got %d",
			len(list.Elements)-1)
	}

	paramsList, ok := list.Elements[1].(sexpr.List)
	if !ok {
		return nil, fmt.Errorf("lambda: parameters must be a list")
	}

	var params []sexpr.Symbol
	for _, p := range paramsList.Elements {
		sym, ok := p.(sexpr.Symbol)
		if !ok {
			return nil, fmt.Errorf("lambda: parameter must be a symbol, got %v", p)
		}
		params = append(params, sym)
	}

	body := list.Elements[2]

	return sexpr.Func{
		Params: params,
		Body:   body,
		Env:    env,
	}, nil
}

// evalIf handles (if test then else)
func evalIf(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	if len(list.Elements) != 4 {
		return nil, fmt.Errorf("if requires 3 arguments, got %d",
			len(list.Elements)-1)
	}

	test, err := Eval(list.Elements[1], env)
	if err != nil {
		return nil, err
	}

	if isTruthy(test) {
		return Eval(list.Elements[2], env)
	}
	return Eval(list.Elements[3], env)
}

// evalQuote handles (quote expr)
func evalQuote(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	if len(list.Elements) != 2 {
		return nil, fmt.Errorf("quote requires 1 argument, got %d",
			len(list.Elements)-1)
	}

	return list.Elements[1], nil
}

// evalApply handles function application
func evalApply(list sexpr.List, env *Env) (sexpr.SExpr, error) {
	// Evaluate the function
	fn, err := Eval(list.Elements[0], env)
	if err != nil {
		return nil, err
	}

	// Evaluate arguments
	var args []sexpr.SExpr
	for _, arg := range list.Elements[1:] {
		value, err := Eval(arg, env)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}

	// Apply function
	switch f := fn.(type) {
	case sexpr.Primitive:
		return f.Fn(args, env)

	case sexpr.Func:
		return applyFunc(f, args)

	default:
		return nil, fmt.Errorf("not a function: %v", fn)
	}
}

// applyFunc applies a user-defined function
func applyFunc(fn sexpr.Func, args []sexpr.SExpr) (sexpr.SExpr, error) {
	if len(args) != len(fn.Params) {
		return nil, fmt.Errorf("function expects %d arguments, got %d",
			len(fn.Params), len(args))
	}

	// Create new environment extending the function's closure
	funcEnv := fn.Env.(*Env).Extend()

	// Bind parameters to arguments
	for i, param := range fn.Params {
		funcEnv.Define(param.Name, args[i])
	}

	// Evaluate body in new environment
	return Eval(fn.Body, funcEnv)
}

// isTruthy determines if a value is truthy
func isTruthy(value sexpr.SExpr) bool {
	switch v := value.(type) {
	case sexpr.Bool:
		return v.Value
	case sexpr.Nil:
		return false
	default:
		return true
	}
}
