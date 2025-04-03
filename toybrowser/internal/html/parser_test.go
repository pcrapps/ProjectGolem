package html

import (ing"
)

// TestParseHTML tests our HTML parser implementation. This test suite demonstrates
// how HTML documents are structured as a tree of nodes, where each node can be
// either an element (like <div> or <p>) or a text node containing the actual content.
//
// The test cases show different aspects of HTML parsing:
// 1. Basic element structure
// 2. Nested elements (parent-child relationships)
// 3. Attributes on elements
func TestParseHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple paragraph",
			// This test demonstrates the basic structure of an HTML element:
			// - An opening tag (<p>)
			// - Text content
			// - A closing tag (</p>)
			// The resulting DOM tree will have:
			// document (root)
			//   └── p (element node)
			//        └── text: Hello, World! (text node)
			input: "<p>Hello, World!</p>",
			expected: "document\n  p\n    text: Hello, World!",
		},
		{
			name: "nested elements",
			// This test shows how HTML elements can be nested inside each other,
			// creating a tree structure. The DOM tree will be:
			// document (root)
			//   └── div (element node)
			//        ├── p (element node)
			//        │    └── text: Hello (text node)
			//        └── p (element node)
			//             └── text: World (text node)
			// This demonstrates parent-child relationships in the DOM.
			input: "<div><p>Hello</p><p>World</p></div>",
			expected: "document\n  div\n    p\n      text: Hello\n    p\n      text: World",
		},
		{
			name: "with attributes",
			// This test demonstrates how HTML attributes are handled:
			// - Attributes are key-value pairs (e.g., class="container")
			// - They provide additional information about elements
			// The DOM tree will include these attributes:
			// document (root)
			//   └── div (element node with class="container")
			//        └── p (element node with id="greeting")
			//             └── text: Hello (text node)
			input: `<div class="container"><p id="greeting">Hello</p></div>`,
			expected: "document\n  div class=\"container\"\n    p id=\"greeting\"\n      text: Hello",
		},
		{
			name: "self-closing tags",
			input: `<img src="test.jpg"/><br/>`,
			expected: "document\n  img src=\"test.jpg\"\n  br\n",
		},
		{
			name: "comments",
			input: `<!-- Header --><h1>Title</h1><!-- Footer -->`,
			expected: "document\n  comment: Header\n  h1\n    text: Title\n  comment: Footer\n",
		},
		{
			name: "doctype",
			input: `<!DOCTYPE html><html><head></head><body></body></html>`,
			expected: "document\n  doctype: DOCTYPE html\n  html\n    head\n    body\n",
		},
		{
			name: "mixed content",
			input: `<!DOCTYPE html>
<!-- Page Start -->
<html>
  <head>
    <meta charset="utf-8"/>
    <title>Test</title>
  </head>
  <body>
    <h1>Hello</h1>
    <img src="test.jpg"/>
    <!-- Section -->
    <p>Text</p>
  </body>
</html>`,
			expected: "document\n  doctype: DOCTYPE html\n  comment: Page Start\n  html\n    head\n      meta charset=\"utf-8\"\n      title\n        text: Test\n    body\n      h1\n        text: Hello\n      img src=\"test.jpg\"\n      comment: Section\n      p\n        text: Text\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the HTML input into a DOM tree
			doc, err := ParseHTML(tt.input)
			if err != nil {
				t.Fatalf("ParseHTML() error = %v", err)
			}
			// Convert the DOM tree to a string representation for comparison
			got := docToString(doc.Root)
			if got != tt.expected {
				t.Errorf("ParseHTML() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// docToString converts a document tree to a string representation for testing.
// This function helps us visualize the structure of the DOM tree by:
// 1. Using indentation to show parent-child relationships
// 2. Including attributes in the output
// 3. Clearly marking text nodes
func docToString(node *Node) string {
	var result string
	docToStringHelper(node, 0, &result)
	return result
}

// docToStringHelper recursively builds the string representation of the DOM tree.
// The depth parameter controls indentation to show the tree structure.
func docToStringHelper(node *Node, depth int, result *string) {
	// Create indentation based on the node's depth in the tree
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	// Handle text nodes differently from element nodes
	if node.Type == TextNode {
		*result += indent + "text: " + node.Text + "\n"
		return
	}

	// For element nodes, include their attributes in the output
	attrStr := ""
	for name, value := range node.Attrs {
		attrStr += " " + name + "=\"" + value + "\""
	}
	*result += indent + node.TagName + attrStr + "\n"
	// Recursively process all child nodes
	for _, child := range node.Children {
		docToStringHelper(child, depth+1, result)
	}
} 