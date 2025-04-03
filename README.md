# Project Golem

A toy browser and JavaScript engine implementation for learning how web browsers work.

## Project Structure

```
golemjs/
├── internal/
│   ├── ast/          # Abstract Syntax Tree implementation
│   ├── lexer/        # JavaScript tokenizer
│   ├── parser/       # JavaScript parser
│   └── interpreter/  # JavaScript interpreter
└── tests/            # Test files

toybrowser/
├── internal/
│   ├── html/         # HTML parser and DOM implementation
│   ├── css/          # CSS parser and styling
│   ├── layout/       # Layout engine
│   └── render/       # Rendering engine
└── examples/         # Example applications
```

## What We've Built

### 1. JavaScript Engine (GolemJS)

We've implemented a JavaScript engine that can:
- Parse JavaScript code into an AST
- Evaluate basic JavaScript expressions
- Handle variables and scope
- Support functions and function calls
- Process basic data types (numbers, strings, booleans)
- Handle control flow (if/else, while loops)

Key components:
- **Lexer**: Breaks JavaScript code into tokens
- **Parser**: Builds an Abstract Syntax Tree (AST)
- **Interpreter**: Evaluates the AST and produces results
- **Environment**: Manages variable scope and closures

### 2. HTML Parser and DOM

We've implemented a basic HTML parser that:
- Parses HTML into a DOM tree
- Handles different node types (elements, text, comments)
- Supports attributes and nested elements
- Maintains parent-child relationships

The DOM implementation includes:
- **Node Types**: Element, Text, Comment, Doctype
- **Node Properties**: Tag name, attributes, children
- **Tree Structure**: Parent-child relationships
- **Document Object**: Root of the DOM tree

### 3. Integration

The project demonstrates how:
- JavaScript can interact with the DOM
- HTML is parsed into a tree structure
- JavaScript code is executed in the context of a web page
- The browser's rendering pipeline works

## Learning Focus

This project is designed to help understand:
1. How JavaScript works at a fundamental level
2. How the DOM represents web pages
3. How JavaScript and the DOM interact
4. How web browsers process and display content

## Next Steps

1. Implement CSS parsing and styling
2. Add layout engine for positioning elements
3. Implement rendering engine
4. Add more JavaScript features (objects, arrays, etc.)
5. Implement event handling
6. Add network capabilities for loading resources

## Running the Project

1. Build the JavaScript engine:
```bash
cd golemjs
go test ./...
```

2. Build the toy browser:
```bash
cd toybrowser
go run examples/simple/main.go
```

## Contributing

This is a learning project. Feel free to:
1. Fork the repository
2. Add new features
3. Fix bugs
4. Improve documentation
5. Share what you've learned

## License

MIT License - feel free to use this code for learning and experimentation.
