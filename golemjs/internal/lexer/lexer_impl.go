package lexer

// LexerImpl represents our concrete lexer implementation.
// The lexer is the first step in processing JavaScript code. It takes the raw source code
// and breaks it down into tokens - the smallest meaningful units of the language.
// For example, "let x = 42;" is broken down into tokens: [LET, IDENT("x"), ASSIGN, INT("42"), SEMICOLON]
type LexerImpl struct {
	input        string // The source code to be tokenized
	position     int    // Current position in input (points to current char)
	readPosition int    // Current reading position in input (after current char)
	ch           byte   // Current char under examination
}

// New creates a new Lexer instance.
// It initializes the lexer with the input string and reads the first character.
func New(input string) *LexerImpl {
	l := &LexerImpl{input: input}
	l.readChar() // Initialize first character
	return l
}

// readChar advances the position and reads the next character.
// This is a fundamental operation that moves the lexer through the input string.
// When it reaches the end of input, it sets the current character to 0 (NUL).
func (l *LexerImpl) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken returns the next token from the input.
// This is the main function of the lexer - it examines the current character
// and returns the appropriate token based on what it finds.
// The function handles:
// - Operators (+, -, *, /, etc.)
// - Delimiters (parentheses, braces, etc.)
// - Keywords (let, function, if, etc.)
// - Identifiers (variable names)
// - Numbers
// - Illegal characters
func (l *LexerImpl) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch)}
		}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch)}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = Token{Type: BANG, Literal: string(l.ch)}
		}
	case '/':
		tok = Token{Type: SLASH, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.ch)}
	case '<':
		tok = Token{Type: LT, Literal: string(l.ch)}
	case '>':
		tok = Token{Type: GT, Literal: string(l.ch)}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch)}
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// peekChar looks at the next character without consuming it.
// This is used for handling multi-character operators like == and !=.
// It allows us to look ahead one character to determine if we're dealing
// with a two-character operator or a single-character one.
func (l *LexerImpl) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips over any whitespace characters.
// Whitespace is not significant in JavaScript (except in strings),
// so we can safely skip over spaces, tabs, newlines, and carriage returns.
func (l *LexerImpl) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier and advances the lexer's position.
// Identifiers are used for variable names, function names, etc.
// They can contain letters, numbers, underscores, and dollar signs,
// but must start with a letter, underscore, or dollar sign.
func (l *LexerImpl) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number and advances the lexer's position.
// Currently handles only integer numbers. In a full JavaScript implementation,
// this would need to handle floating-point numbers, scientific notation, etc.
func (l *LexerImpl) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter checks if the character is a letter.
// In JavaScript, identifiers can contain letters (a-z, A-Z),
// underscores (_), and dollar signs ($).
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if the character is a digit.
// Used for parsing numbers in the source code.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent checks if the identifier is a keyword.
// Keywords are special identifiers that have specific meaning in JavaScript.
// Examples include: let, function, if, else, return, etc.
func lookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "if":
		return IF
	case "else":
		return ELSE
	case "return":
		return RETURN
	default:
		return IDENT
	}
}
