package main

import (
	"toybrowser/internal/html"
	"toybrowser/internal/render"
)

func main() {
	// Create a simple HTML document
	doc := html.NewDocument()
	
	// Create the root HTML element
	htmlNode := html.NewNode(html.ElementNode, "html")
	doc.Root.AddChild(htmlNode)

	// Create the head element
	head := html.NewNode(html.ElementNode, "head")
	htmlNode.AddChild(head)

	// Add a title
	title := html.NewNode(html.ElementNode, "title")
	titleText := html.NewNode(html.TextNode, "")
	titleText.Text = "Toy Browser Example"
	title.AddChild(titleText)
	head.AddChild(title)

	// Create the body element
	body := html.NewNode(html.ElementNode, "body")
	htmlNode.AddChild(body)

	// Add a heading
	h1 := html.NewNode(html.ElementNode, "h1")
	h1Text := html.NewNode(html.TextNode, "")
	h1Text.Text = "Welcome to Toy Browser!"
	h1.AddChild(h1Text)
	body.AddChild(h1)

	// Add a paragraph
	p := html.NewNode(html.ElementNode, "p")
	pText := html.NewNode(html.TextNode, "")
	pText.Text = "This is a simple example of our toy browser rendering HTML content."
	p.AddChild(pText)
	body.AddChild(p)

	// Create and configure the renderer
	renderer := render.NewWebViewRenderer("Toy Browser")
	
	// Render the document
	if err := renderer.Render(doc); err != nil {
		panic(err)
	}

	// Start the webview event loop
	renderer.Run()
} 