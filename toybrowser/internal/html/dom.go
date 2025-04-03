package html

// NodeType represents the different types of nodes in the DOM tree.
// In a real browser, there are many more node types, but for our toy browser
// we'll focus on the most common ones:
// - ElementNode: Represents HTML elements like <div>, <p>, etc.
// - TextNode: Contains the actual text content
// - CommentNode: Represents HTML comments
// - DoctypeNode: Represents the DOCTYPE declaration
type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
	CommentNode
	DoctypeNode
)

// Node represents a single node in the DOM tree.
// This is the fundamental building block of the DOM - every piece of content
// in an HTML document becomes a node in the tree.
//
// Key concepts:
// - Type: Determines what kind of node this is (element, text, etc.)
// - TagName: For element nodes, stores the HTML tag name (e.g., "div", "p")
// - Text: For text nodes, stores the actual content
// - Attrs: Stores HTML attributes as key-value pairs
// - Children: Contains all child nodes, creating the tree structure
// - Parent: Reference to the parent node (except for the root)
type Node struct {
	Type     NodeType
	TagName  string
	Text     string
	Attrs    map[string]string
	Children []*Node
	Parent   *Node
}

// NewNode creates a new node with the given type and tag name.
// This is our factory function for creating nodes - it ensures all nodes
// are properly initialized with their required fields.
func NewNode(nodeType NodeType, tagName string) *Node {
	return &Node{
		Type:     nodeType,
		TagName:  tagName,
		Attrs:    make(map[string]string),
		Children: make([]*Node, 0),
	}
}

// AddChild adds a child node to this node and sets up the parent-child relationship.
// This is how we build the DOM tree - by connecting nodes together.
// The parent reference allows us to traverse up the tree, while the children
// array allows us to traverse down.
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// SetAttribute adds or updates an HTML attribute on this node.
// Attributes are key-value pairs that provide additional information
// about elements (like class names, IDs, styles, etc.)
func (n *Node) SetAttribute(name, value string) {
	n.Attrs[name] = value
}

// GetAttribute retrieves the value of an HTML attribute.
// Returns an empty string if the attribute doesn't exist.
func (n *Node) GetAttribute(name string) string {
	return n.Attrs[name]
}

// Document represents the root of the DOM tree.
// In a real browser, the Document object has many more properties and methods,
// but for our toy browser, we'll focus on the essential ones:
// - Root: The root element (usually <html>)
// - Title: The page title from the <title> tag
type Document struct {
	Root  *Node
	Title string
}

// NewDocument creates a new document with an empty root node.
// This is our factory function for creating new documents.
func NewDocument() *Document {
	return &Document{
		Root: NewNode(ElementNode, "document"),
	}
}
