package html

import ("
)

// NodeType represents the type of a DOM node. In the DOM, different types of nodes
// have different behaviors and properties. For example, element nodes can have
// children and attributes, while text nodes can only contain text.
type NodeType int

const (
	ElementNode NodeType = iota // Represents HTML elements like <div>, <p>, etc.
	TextNode                   // Represents text content within elements
	DocumentNode               // Represents the root document node
	CommentNode                // Represents HTML comments <!-- comment -->
	DoctypeNode                // Represents the DOCTYPE declaration
)

// List of void elements (self-closing tags). These are HTML elements that cannot
// have content and don't need a closing tag. For example: <img>, <br>, <input>.
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// Parser represents an HTML parser. The parser maintains state about its current
// position in the input and the current node being processed. It uses a stack to
// keep track of parent nodes while building the DOM tree.
type Parser struct {
	pos     int      // Current position in the input string
	input   string   // The HTML text being parsed
	current *Node    // The current node being processed
	stack   []*Node  // Stack of parent nodes for maintaining hierarchy
}

// NewParser creates a new HTML parser with the given input string.
// The parser starts at the beginning of the input with an empty stack.
func NewParser(input string) *Parser {
	return &Parser{
		pos:   0,
		input: input,
		stack: make([]*Node, 0),
	}
}

// Parse parses the HTML input and returns a Document. This is the main parsing
// function that implements a basic HTML parser. It handles:
// 1. Opening and closing tags
// 2. Self-closing tags
// 3. Text nodes
// 4. Comments
// 5. DOCTYPE declarations
// 6. Attributes
func (p *Parser) Parse() (*Document, error) {
	// Create a new document with a document node as root
	doc := NewDocument()
	p.current = doc.Root
	p.stack = append(p.stack, p.current)

	// Process the input character by character
	for p.pos < len(p.input) {
		if p.input[p.pos] == '<' {
			// We've found a tag or special construct
			if p.pos+1 >= len(p.input) {
				break
			}

			switch p.input[p.pos+1] {
			case '!':
				// Handle comments and DOCTYPE declarations
				if p.pos+3 < len(p.input) && p.input[p.pos+2] == '-' && p.input[p.pos+3] == '-' {
					// Parse HTML comment <!-- comment -->
					p.consumeChar() // '<'
					p.consumeChar() // '!'
					p.consumeChar() // '-'
					p.consumeChar() // '-'
					comment := p.consumeUntil('-')
					if p.pos+2 < len(p.input) && p.input[p.pos+1] == '-' && p.input[p.pos+2] == '>' {
						p.pos += 3 // Skip "-->"
						node := NewNode(CommentNode, "")
						node.Text = strings.TrimSpace(comment)
						p.current.AddChild(node)
					}
				} else if strings.HasPrefix(p.input[p.pos:], "<!DOCTYPE") {
					// Parse DOCTYPE declaration
					p.consumeChar() // '<'
					p.consumeChar() // '!'
					doctype := p.consumeUntil('>')
					p.consumeChar() // '>'
					node := NewNode(DoctypeNode, "")
					node.Text = strings.TrimSpace(doctype)
					p.current.AddChild(node)
				}
			case '/':
				// Handle closing tags
				p.consumeChar() // '<'
				p.consumeChar() // '/'
				tagName := p.consumeUntil('>')
				p.consumeChar() // '>'
				tagName = strings.ToLower(strings.TrimSpace(tagName))

				// Pop nodes from the stack until we find the matching opening tag
				for len(p.stack) > 1 {
					last := p.stack[len(p.stack)-1]
					if last.TagName == tagName {
						p.stack = p.stack[:len(p.stack)-1]
						p.current = p.stack[len(p.stack)-1]
						break
					}
					p.stack = p.stack[:len(p.stack)-1]
				}
			default:
				// Handle opening tags
				p.consumeChar() // '<'
				tag := p.consumeUntil('>')
				p.consumeChar() // '>'

				// Check if it's a self-closing tag (ends with '/')
				selfClosing := false
				if strings.HasSuffix(tag, "/") {
					tag = strings.TrimSuffix(tag, "/")
					selfClosing = true
				}

				// Parse tag name and attributes
				parts := strings.Fields(tag)
				if len(parts) == 0 {
					continue
				}
				tagName := strings.ToLower(parts[0])

				// Create the element node
				node := NewNode(ElementNode, tagName)

				// Parse attributes (name="value" pairs)
				for i := 1; i < len(parts); i++ {
					attr := parts[i]
					if strings.Contains(attr, "=") {
						kv := strings.SplitN(attr, "=", 2)
						name := strings.ToLower(kv[0])
						value := strings.Trim(kv[1], "\"'")
						node.SetAttribute(name, value)
					}
				}

				// Add the node to the current parent
				p.current.AddChild(node)

				// For non-void and non-self-closing elements, push onto stack
				if !voidElements[tagName] && !selfClosing {
					p.current = node
					p.stack = append(p.stack, node)
				}
			}
		} else {
			// Handle text content between tags
			text := p.consumeUntil('<')
			if text = strings.TrimSpace(text); text != "" {
				textNode := NewNode(TextNode, "")
				textNode.Text = text
				p.current.AddChild(textNode)
			}
		}
	}

	return doc, nil
}

// consumeChar consumes and returns the current character, advancing the position.
// Returns 0 if we've reached the end of the input.
func (p *Parser) consumeChar() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	char := p.input[p.pos]
	p.pos++
	return char
}

// consumeUntil consumes characters until the given character is found.
// Returns the consumed text, not including the target character.
func (p *Parser) consumeUntil(char byte) string {
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != char {
		p.pos++
	}
	return p.input[start:p.pos]
}

// ParseHTML parses HTML text into a Document. This is the main entry point
// for parsing HTML. It creates a new parser and returns the resulting document.
func ParseHTML(input string) (*Document, error) {
	parser := NewParser(input)
	return parser.Parse()
}
