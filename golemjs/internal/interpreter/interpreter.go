package interpreter

import (
	"fmt"
	"strings"

	"github.com/biosbuddha/golemjs/internal/ast"
)

// Object represents a JavaScript object in our interpreter.
// In JavaScript, everything is an object, including:
// - Numbers, strings, booleans (primitive objects)
// - Arrays and objects (compound objects)
// - Functions (callable objects)
// - null and undefined (special objects)
type Object interface {
	Type() ObjectType
	Inspect() string
}

// ObjectType represents the different types of JavaScript objects.
// This helps us distinguish between different kinds of values and
// implement appropriate behavior for each type.
type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ = "BUILTIN"
	ARRAY_OBJ = "ARRAY"
	HASH_OBJ = "HASH"
)

// Null represents JavaScript's null value.
// It's a special value that represents the intentional absence of any object value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Error represents a JavaScript error object.
// Errors can occur during evaluation and need to be handled appropriately.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Integer represents JavaScript numbers.
// In our toy implementation, we only handle integers for simplicity.
// A real JavaScript engine would handle floating-point numbers as well.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// String represents JavaScript strings.
// Strings are immutable sequences of characters.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Boolean represents JavaScript boolean values.
// There are only two possible values: true and false.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// ReturnValue represents a return statement's value.
// This is a special object that helps us implement the return statement
// by allowing us to propagate the return value up the call stack.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Function represents a JavaScript function.
// Functions are objects that can be called with arguments.
// They contain:
// - Parameters: The function's formal parameters
// - Body: The function's body (an AST node)
// - Env: The environment where the function was defined (for closures)
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out strings.Builder
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// BuiltinFunction represents a built-in JavaScript function.
// These are functions implemented in Go that provide core functionality
// like console.log, parseInt, etc.
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Array represents JavaScript arrays.
// Arrays are ordered collections of values that can be of any type.
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// HashKey represents a key in a JavaScript object.
// In JavaScript, object keys are always strings.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashPair represents a key-value pair in a JavaScript object.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a JavaScript object (not to be confused with HashKey).
// Objects are collections of properties (key-value pairs).
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Value.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Hashable represents an object that can be used as a hash key.
// In JavaScript, only strings can be used as object keys.
type Hashable interface {
	HashKey() HashKey
}

// Environment represents a JavaScript scope.
// Environments are used to implement variable scoping and closures.
// They form a chain (like a linked list) where each environment
// has a reference to its outer (parent) environment.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new environment.
// The outer parameter is used to create nested scopes.
func NewEnvironment(outer *Environment) *Environment {
	env := &Environment{store: make(map[string]Object), outer: outer}
	return env
}

// Get retrieves a variable from the environment.
// If the variable isn't found in the current environment,
// it looks in the outer environment (implementing variable shadowing).
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores a variable in the current environment.
// Note that this doesn't modify variables in outer environments.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Interpreter represents our JavaScript interpreter.
// It's responsible for evaluating AST nodes and producing JavaScript values.
type Interpreter struct {
	env *Environment
}

// New creates a new interpreter with a fresh environment.
func New() *Interpreter {
	env := NewEnvironment(nil)
	return &Interpreter{env: env}
}

// Eval evaluates an AST node and returns the resulting JavaScript value.
// This is the main entry point for evaluation.
func (i *Interpreter) Eval(node ast.Node) Object {
	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node)
	case *ast.ExpressionStatement:
		return i.Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &String{Value: node.Value}
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
	case *ast.ArrayLiteral:
		elements := i.evalExpressions(node.Elements)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	case *ast.IndexExpression:
		left := i.Eval(node.Left)
		if isError(left) {
			return left
		}
		index := i.Eval(node.Index)
		if isError(index) {
			return index
		}
		return i.evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return i.evalHashLiteral(node)
	}
	return nil
}

// evalProgram evaluates a program (the root node of the AST).
// It evaluates each statement in sequence and returns the last value.
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

// evalBlockStatement evaluates a block of statements.
// It creates a new environment for the block to implement proper scoping.
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

// evalPrefixExpression evaluates prefix expressions like -5 or !true.
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

// evalInfixExpression evaluates infix expressions like 5 + 5 or true && false.
func (i *Interpreter) evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return i.evalIntegerInfixExpression(operator, left, right)
	case operator == "+":
		return i.evalStringInfixExpression(left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression evaluates arithmetic expressions between integers.
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression evaluates string concatenation.
func (i *Interpreter) evalStringInfixExpression(left, right Object) Object {
	if left.Type() != STRING_OBJ || right.Type() != STRING_OBJ {
		return newError("type mismatch: %s + %s", left.Type(), right.Type())
	}
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
}

// evalIfExpression evaluates if expressions and their else clauses.
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

// evalIdentifier evaluates identifiers (variable names).
func (i *Interpreter) evalIdentifier(node *ast.Identifier) Object {
	if val, ok := i.env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

// evalExpressions evaluates a list of expressions (used for function arguments).
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

// applyFunction applies a function to its arguments.
// This handles both user-defined functions and built-in functions.
func (i *Interpreter) applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := i.extendFunctionEnv(fn, args)
		evaluated := i.Eval(fn.Body)
		return i.unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv creates a new environment for a function call.
// This implements proper scoping for function parameters and local variables.
func (i *Interpreter) extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// unwrapReturnValue handles return values from functions.
func (i *Interpreter) unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// evalIndexExpression evaluates array and object indexing expressions.
func (i *Interpreter) evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ:
		return i.evalArrayIndexExpression(left, index)
	case left.Type() == HASH_OBJ:
		return i.evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalArrayIndexExpression evaluates array indexing expressions.
func (i *Interpreter) evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

// evalHashIndexExpression evaluates object property access expressions.
func (i *Interpreter) evalHashIndexExpression(hash, index Object) Object {
	hashObject := hash.(*Hash)
	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

// evalHashLiteral evaluates object literals.
func (i *Interpreter) evalHashLiteral(node *ast.HashLiteral) Object {
	pairs := make(map[HashKey]HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := i.Eval(keyNode)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := i.Eval(valueNode)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}
	return &Hash{Pairs: pairs}
}

// Helper functions for type conversion and error checking
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

// Built-in functions
var TRUE = &Boolean{Value: true}
var FALSE = &Boolean{Value: false}
var NULL = &Null{}

var builtins = map[string]*Builtin{
	"len": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			newElements := make([]Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &Array{Elements: newElements}
		},
	},
	"puts": &Builtin{
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
} 