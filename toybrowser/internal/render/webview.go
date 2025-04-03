package render

import (
	"fmt"
	"html/template"
	"strings"

	"toybrowser/internal/html"

	"github.com/webview/webview"
)

// WebViewRenderer represents our webview-based renderer.
// It creates a native window and renders HTML content into it.
type WebViewRenderer struct {
	webview webview.WebView
	doc     *html.Document
}

// NewWebViewRenderer creates a new webview renderer.
// The title parameter sets the window title.
func NewWebViewRenderer(title string) *WebViewRenderer {
	w := webview.New(true)
	defer w.Init()
	w.SetTitle(title)
	w.SetSize(800, 600, webview.HintNone)
	return &WebViewRenderer{
		webview: w,
	}
}

// Render renders the HTML document in the webview window.
// This converts our DOM tree into HTML and displays it.
func (r *WebViewRenderer) Render(doc *html.Document) error {
	r.doc = doc
	html := r.generateHTML(doc.Root)
	r.webview.SetHTML(html)
	return nil
}

// Run starts the webview event loop.
// This should be called after setting up the content.
func (r *WebViewRenderer) Run() {
	r.webview.Run()
}

// generateHTML converts our DOM tree into HTML.
// This is a simple implementation that handles basic elements and text nodes.
func (r *WebViewRenderer) generateHTML(node *html.Node) string {
	var sb strings.Builder

	// Handle different node types
	switch node.Type {
	case html.ElementNode:
		// Start tag with attributes
		sb.WriteString("<")
		sb.WriteString(node.TagName)
		for name, value := range node.Attrs {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", name, value))
		}
		sb.WriteString(">")

		// Children
		for _, child := range node.Children {
			sb.WriteString(r.generateHTML(child))
		}

		// End tag
		sb.WriteString("</")
		sb.WriteString(node.TagName)
		sb.WriteString(">")

	case html.TextNode:
		// Escape HTML special characters in text
		sb.WriteString(template.HTMLEscapeString(node.Text))

	case html.CommentNode:
		sb.WriteString("<!--")
		sb.WriteString(node.Text)
		sb.WriteString("-->")

	case html.DoctypeNode:
		sb.WriteString("<!DOCTYPE ")
		sb.WriteString(node.Text)
		sb.WriteString(">")
	}

	return sb.String()
}

// InjectJavaScript injects JavaScript code into the webview.
// This allows us to execute JavaScript in the context of the rendered page.
func (r *WebViewRenderer) InjectJavaScript(js string) error {
	r.webview.Eval(js)
	return nil
}

// Bind binds a Go function to JavaScript.
// This allows JavaScript code to call Go functions.
func (r *WebViewRenderer) Bind(name string, fn interface{}) error {
	return r.webview.Bind(name, fn)
} 