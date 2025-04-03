package ast

// Node represents a node in the Abstract Syntax Tree (AST).
// The AST is a tree representation of the source code where each node represents
// a construct occurring in the source code. This is the foundation of how JavaScript
// code is structured and processed by the engine.
type Node interface {
	TokenLiteral() string // Returns the literal value of the token that created this node
	String() string       // Returns a string representation of the node for debugging
}

// Statement represents a statement node in the AST.
// Statements are the building blocks of JavaScript programs - they are instructions
// that perform actions. Examples include variable declarations, function declarations,
// and control flow statements.
type Statement interface {
	Node
	statementNode() // Marker method to distinguish statements from expressions
}

// Expression represents an expression node in the AST.
// Expressions are pieces of code that evaluate to a value. They can be as simple
// as a literal number or as complex as a function call with multiple arguments.
type Expression interface {
	Node
	expressionNode() // Marker method to distinguish expressions from statements
}

// Program represents the root node of every AST.
// It contains a list of statements that make up the entire JavaScript program.
// Think of it as the top-level container for all code in a JavaScript file.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String() + "\n"
	}
	return out
}

// Identifier represents an identifier (variable name, function name, etc.)
// Identifiers are names used to identify variables, functions, or other user-defined items.
// They must start with a letter, underscore, or dollar sign and can contain letters, numbers, underscores, or dollar signs.
type Identifier struct {
	Token Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Literal represents a literal value in the source code.
// Literals are values written directly in the code, such as numbers (42), strings ("hello"),
// booleans (true), null, or undefined.
type Literal struct {
	Token Token
	Value interface{}
}

func (l *Literal) expressionNode()      {}
func (l *Literal) TokenLiteral() string { return l.Token.Literal }
func (l *Literal) String() string       { return l.Token.Literal }

// BinaryExpression represents binary operations like addition, subtraction, etc.
// Binary expressions have a left side, an operator, and a right side.
// Example: 5 + 3 is a binary expression where 5 is the left side, + is the operator, and 3 is the right side.
type BinaryExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) expressionNode()      {}
func (b *BinaryExpression) TokenLiteral() string { return b.Token.Literal }
func (b *BinaryExpression) String() string {
	return "(" + b.Left.String() + " " + b.Operator + " " + b.Right.String() + ")"
}

// VariableDeclaration represents variable declarations using var, let, or const.
// This node captures how variables are declared in JavaScript, including their name
// and optional initial value. The declaration type (var/let/const) is stored in the token.
type VariableDeclaration struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (v *VariableDeclaration) statementNode()       {}
func (v *VariableDeclaration) TokenLiteral() string { return v.Token.Literal }
func (v *VariableDeclaration) String() string {
	var out string
	out += v.TokenLiteral() + " "
	out += v.Name.String()
	if v.Value != nil {
		out += " = " + v.Value.String()
	}
	return out + ";"
}

// FunctionDeclaration represents function declarations in the code.
// Functions are reusable blocks of code that can be called with different arguments.
// This node captures the function's name, parameters, and body.
type FunctionDeclaration struct {
	Token      Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionDeclaration) statementNode()       {}
func (f *FunctionDeclaration) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionDeclaration) String() string {
	var out string
	out += "function " + f.Name.String() + "("

	for i, p := range f.Parameters {
		if i > 0 {
			out += ", "
		}
		out += p.String()
	}

	out += ") " + f.Body.String()
	return out
}

// CallExpression represents function calls in the code.
// When a function is called, it's represented as a call expression with the function
// to be called and the arguments being passed to it.
type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var out string
	out += c.Function.String() + "("

	for i, a := range c.Arguments {
		if i > 0 {
			out += ", "
		}
		out += a.String()
	}

	return out + ")"
}

// BlockStatement represents a block of code enclosed in curly braces.
// Blocks create a new scope and can contain multiple statements.
// They are used in function bodies, if statements, while loops, etc.
type BlockStatement struct {
	Token      Token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var out string
	out += "{\n"
	for _, s := range b.Statements {
		out += "  " + s.String() + "\n"
	}
	return out + "}"
}

// IfStatement represents if/else conditional statements.
// If statements allow code to be executed conditionally based on a boolean expression.
// The Alternative field represents the else clause, which is optional.
type IfStatement struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // can be nil for if without else
}

func (i *IfStatement) statementNode()       {}
func (i *IfStatement) TokenLiteral() string { return i.Token.Literal }
func (i *IfStatement) String() string {
	var out string
	out += "if (" + i.Condition.String() + ") " + i.Consequence.String()
	if i.Alternative != nil {
		out += " else " + i.Alternative.String()
	}
	return out
}

// WhileStatement represents while loops.
// While loops repeatedly execute a block of code as long as a condition is true.
type WhileStatement struct {
	Token     Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode()       {}
func (w *WhileStatement) TokenLiteral() string { return w.Token.Literal }
func (w *WhileStatement) String() string {
	return "while (" + w.Condition.String() + ") " + w.Body.String()
}

// ReturnStatement represents return statements in functions.
// Return statements specify the value to be returned from a function.
// The ReturnValue field can be nil for functions that don't return a value.
type ReturnStatement struct {
	Token       Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	if r.ReturnValue != nil {
		return "return " + r.ReturnValue.String() + ";"
	}
	return "return;"
}

// Helper functions for type checking
func IsExpression(node Node) bool {
	_, ok := node.(Expression)
	return ok
}

func IsStatement(node Node) bool {
	_, ok := node.(Statement)
	return ok
}

// GetNodeType returns a string representation of the node's type
func GetNodeType(node Node) string {
	switch node.(type) {
	case *Program:
		return "Program"
	case *Identifier:
		return "Identifier"
	case *Literal:
		return "Literal"
	case *BinaryExpression:
		return "BinaryExpression"
	case *VariableDeclaration:
		return "VariableDeclaration"
	case *FunctionDeclaration:
		return "FunctionDeclaration"
	case *CallExpression:
		return "CallExpression"
	case *BlockStatement:
		return "BlockStatement"
	case *IfStatement:
		return "IfStatement"
	case *WhileStatement:
		return "WhileStatement"
	case *ReturnStatement:
		return "ReturnStatement"
	default:
		return "Unknown"
	}
}
