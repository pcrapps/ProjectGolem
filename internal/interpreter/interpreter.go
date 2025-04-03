package interpreter

import (
	"fmt"

	"github.com/biosbuddha/golemjs/internal/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type Interpreter struct {
	env *Environment
}

func New() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

func (i *Interpreter) Eval(node ast.Node) Object {
	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node)
	case *ast.ExpressionStatement:
		return i.Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := i.Eval(node.Right)
		if isError(right) {
			return right
		}
		return i.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := i.Eval(node.Left)
		if isError(left) {
			return left
		}
		right := i.Eval(node.Right)
		if isError(right) {
			return right
		}
		return i.evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return i.evalBlockStatement(node)
	case *ast.IfExpression:
		return i.evalIfExpression(node)
	case *ast.ReturnStatement:
		val := i.Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *ast.LetStatement:
		val := i.Eval(node.Value)
		if isError(val) {
			return val
		}
		i.env.Set(node.Name.Value, val)
		return val
	case *ast.Identifier:
		return i.evalIdentifier(node)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Body: body, Env: i.env}
	case *ast.CallExpression:
		function := i.Eval(node.Function)
		if isError(function) {
			return function
		}
		args := i.evalExpressions(node.Arguments)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return i.applyFunction(function, args)
	}
	return nil
}

func (i *Interpreter) evalProgram(program *ast.Program) Object {
	var result Object
	for _, statement := range program.Statements {
		result = i.Eval(statement)
		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func (i *Interpreter) evalBlockStatement(block *ast.BlockStatement) Object {
	var result Object
	for _, statement := range block.Statements {
		result = i.Eval(statement)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func (i *Interpreter) evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return i.evalBangOperatorExpression(right)
	case "-":
		return i.evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func (i *Interpreter) evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return i.evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (i *Interpreter) evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func (i *Interpreter) evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func (i *Interpreter) evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (i *Interpreter) evalIfExpression(ie *ast.IfExpression) Object {
	condition := i.Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return i.Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return i.Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func (i *Interpreter) evalIdentifier(node *ast.Identifier) Object {
	if val, ok := i.env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: " + node.Value)
}

func (i *Interpreter) evalExpressions(exps []ast.Expression) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := i.Eval(e)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func (i *Interpreter) applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := i.extendFunctionEnv(fn, args)
		interpreter := &Interpreter{env: extendedEnv}
		evaluated := interpreter.Eval(fn.Body)
		return i.unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func (i *Interpreter) extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func (i *Interpreter) unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

var TRUE = &Boolean{Value: true}
var FALSE = &Boolean{Value: false}
var NULL = &Null{}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return "function"
}

type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return FUNCTION_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
