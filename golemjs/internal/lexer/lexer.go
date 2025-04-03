package lexer

// TokenType represents the type of a token in the JavaScript language.
// Each token type corresponds to a specific construct in the language,
// such as keywords, operators, or identifiers.
type TokenType string

const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL" // Represents an illegal/unrecognized character
	EOF     TokenType = "EOF"     // End of file - indicates we've reached the end of input

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // Variable names, function names, etc. (e.g., "x", "add", "foobar")
	INT    TokenType = "INT"    // Integer literals (e.g., "123", "42")
	STRING TokenType = "STRING" // String literals (e.g., "hello", "world")

	// Operators
	ASSIGN   TokenType = "="  // Assignment operator (e.g., x = 42)
	PLUS     TokenType = "+"  // Addition operator
	MINUS    TokenType = "-"  // Subtraction operator
	BANG     TokenType = "!"  // Logical NOT operator
	ASTERISK TokenType = "*"  // Multiplication operator
	SLASH    TokenType = "/"  // Division operator
	LT       TokenType = "<"  // Less than operator
	GT       TokenType = ">"  // Greater than operator
	EQ       TokenType = "==" // Equality operator
	NOT_EQ   TokenType = "!=" // Inequality operator

	// Delimiters
	COMMA     TokenType = ","  // Separates items in lists (e.g., function arguments)
	SEMICOLON TokenType = ";"  // Statement terminator
	LPAREN    TokenType = "("  // Left parenthesis - used for grouping and function calls
	RPAREN    TokenType = ")"  // Right parenthesis
	LBRACE    TokenType = "{"  // Left brace - starts a block of code
	RBRACE    TokenType = "}"  // Right brace - ends a block of code

	// Keywords
	FUNCTION TokenType = "FUNCTION" // "function" keyword for function declarations
	LET      TokenType = "LET"      // "let" keyword for variable declarations
	TRUE     TokenType = "TRUE"     // Boolean literal "true"
	FALSE    TokenType = "FALSE"    // Boolean literal "false"
	IF       TokenType = "IF"       // "if" keyword for conditional statements
	ELSE     TokenType = "ELSE"     // "else" keyword for else clauses
	RETURN   TokenType = "RETURN"   // "return" keyword for returning values from functions
)

// Token represents a single token in the input.
// Each token has a type (what kind of token it is) and a literal value
// (the actual characters that make up the token).
type Token struct {
	Type    TokenType // The type of token (e.g., IDENT, INT, PLUS)
	Literal string    // The actual characters that make up the token
}

// Lexer represents the lexer interface.
// The lexer is responsible for breaking down the input string into tokens.
// It provides a NextToken() method that returns the next token in the input.
type Lexer interface {
	NextToken() Token
}
