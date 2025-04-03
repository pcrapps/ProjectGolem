package ast

import (
	"testing"

	"github.com/biosbuddha/golemjs/internal/ast"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected string
	}{
		{
			name: "Program with multiple statements",
			node: &ast.Program{
				Statements: []ast.Statement{
					&ast.VariableDeclaration{
						Token: ast.Token{Type: "LET", Literal: "let"},
						Name:  &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
						Value: &ast.Literal{Token: ast.Token{Type: "INT", Literal: "5"}, Value: 5},
					},
					&ast.ReturnStatement{
						Token:       ast.Token{Type: "RETURN", Literal: "return"},
						ReturnValue: &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
					},
				},
			},
			expected: "let x = 5;\nreturn x;\n",
		},
		{
			name: "Binary Expression",
			node: &ast.BinaryExpression{
				Token:    ast.Token{Type: "PLUS", Literal: "+"},
				Left:     &ast.Literal{Token: ast.Token{Type: "INT", Literal: "5"}, Value: 5},
				Operator: "+",
				Right:    &ast.Literal{Token: ast.Token{Type: "INT", Literal: "3"}, Value: 3},
			},
			expected: "(5 + 3)",
		},
		{
			name: "Function Declaration",
			node: &ast.FunctionDeclaration{
				Token: ast.Token{Type: "FUNCTION", Literal: "function"},
				Name:  &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "add"}, Value: "add"},
				Parameters: []*ast.Identifier{
					{Token: ast.Token{Type: "IDENT", Literal: "a"}, Value: "a"},
					{Token: ast.Token{Type: "IDENT", Literal: "b"}, Value: "b"},
				},
				Body: &ast.BlockStatement{
					Token: ast.Token{Type: "LBRACE", Literal: "{"},
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Token: ast.Token{Type: "RETURN", Literal: "return"},
							ReturnValue: &ast.BinaryExpression{
								Token:    ast.Token{Type: "PLUS", Literal: "+"},
								Left:     &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "a"}, Value: "a"},
								Operator: "+",
								Right:    &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "b"}, Value: "b"},
							},
						},
					},
				},
			},
			expected: "function add(a, b) {\n  return (a + b);\n}",
		},
		{
			name: "If Statement",
			node: &ast.IfStatement{
				Token: ast.Token{Type: "IF", Literal: "if"},
				Condition: &ast.BinaryExpression{
					Token:    ast.Token{Type: "LT", Literal: "<"},
					Left:     &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
					Operator: "<",
					Right:    &ast.Literal{Token: ast.Token{Type: "INT", Literal: "10"}, Value: 10},
				},
				Consequence: &ast.BlockStatement{
					Token: ast.Token{Type: "LBRACE", Literal: "{"},
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Token:       ast.Token{Type: "RETURN", Literal: "return"},
							ReturnValue: &ast.Literal{Token: ast.Token{Type: "TRUE", Literal: "true"}, Value: true},
						},
					},
				},
				Alternative: &ast.BlockStatement{
					Token: ast.Token{Type: "LBRACE", Literal: "{"},
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Token:       ast.Token{Type: "RETURN", Literal: "return"},
							ReturnValue: &ast.Literal{Token: ast.Token{Type: "FALSE", Literal: "false"}, Value: false},
						},
					},
				},
			},
			expected: "if ((x < 10)) {\n  return true;\n} else {\n  return false;\n}",
		},
		{
			name: "While Statement",
			node: &ast.WhileStatement{
				Token: ast.Token{Type: "WHILE", Literal: "while"},
				Condition: &ast.BinaryExpression{
					Token:    ast.Token{Type: "GT", Literal: ">"},
					Left:     &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
					Operator: ">",
					Right:    &ast.Literal{Token: ast.Token{Type: "INT", Literal: "0"}, Value: 0},
				},
				Body: &ast.BlockStatement{
					Token: ast.Token{Type: "LBRACE", Literal: "{"},
					Statements: []ast.Statement{
						&ast.VariableDeclaration{
							Token: ast.Token{Type: "LET", Literal: "let"},
							Name:  &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
							Value: &ast.BinaryExpression{
								Token:    ast.Token{Type: "MINUS", Literal: "-"},
								Left:     &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
								Operator: "-",
								Right:    &ast.Literal{Token: ast.Token{Type: "INT", Literal: "1"}, Value: 1},
							},
						},
					},
				},
			},
			expected: "while ((x > 0)) {\n  let x = (x - 1);\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeChecking(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		isExpr   bool
		isStmt   bool
		nodeType string
	}{
		{
			name:     "Identifier",
			node:     &ast.Identifier{Token: ast.Token{Type: "IDENT", Literal: "x"}, Value: "x"},
			isExpr:   true,
			isStmt:   false,
			nodeType: "Identifier",
		},
		{
			name:     "Literal",
			node:     &ast.Literal{Token: ast.Token{Type: "INT", Literal: "5"}, Value: 5},
			isExpr:   true,
			isStmt:   false,
			nodeType: "Literal",
		},
		{
			name:     "BinaryExpression",
			node:     &ast.BinaryExpression{Token: ast.Token{Type: "PLUS", Literal: "+"}},
			isExpr:   true,
			isStmt:   false,
			nodeType: "BinaryExpression",
		},
		{
			name:     "VariableDeclaration",
			node:     &ast.VariableDeclaration{Token: ast.Token{Type: "LET", Literal: "let"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "VariableDeclaration",
		},
		{
			name:     "FunctionDeclaration",
			node:     &ast.FunctionDeclaration{Token: ast.Token{Type: "FUNCTION", Literal: "function"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "FunctionDeclaration",
		},
		{
			name:     "CallExpression",
			node:     &ast.CallExpression{Token: ast.Token{Type: "LPAREN", Literal: "("}},
			isExpr:   true,
			isStmt:   false,
			nodeType: "CallExpression",
		},
		{
			name:     "BlockStatement",
			node:     &ast.BlockStatement{Token: ast.Token{Type: "LBRACE", Literal: "{"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "BlockStatement",
		},
		{
			name:     "IfStatement",
			node:     &ast.IfStatement{Token: ast.Token{Type: "IF", Literal: "if"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "IfStatement",
		},
		{
			name:     "WhileStatement",
			node:     &ast.WhileStatement{Token: ast.Token{Type: "WHILE", Literal: "while"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "WhileStatement",
		},
		{
			name:     "ReturnStatement",
			node:     &ast.ReturnStatement{Token: ast.Token{Type: "RETURN", Literal: "return"}},
			isExpr:   false,
			isStmt:   true,
			nodeType: "ReturnStatement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ast.IsExpression(tt.node); got != tt.isExpr {
				t.Errorf("IsExpression() = %v, want %v", got, tt.isExpr)
			}
			if got := ast.IsStatement(tt.node); got != tt.isStmt {
				t.Errorf("IsStatement() = %v, want %v", got, tt.isStmt)
			}
			if got := ast.GetNodeType(tt.node); got != tt.nodeType {
				t.Errorf("GetNodeType() = %v, want %v", got, tt.nodeType)
			}
		})
	}
}
