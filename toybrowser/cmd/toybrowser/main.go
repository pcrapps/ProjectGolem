package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"toybrowser/internal/html"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: toybrowser <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse HTML
	doc, err := html.ParseHTML(string(body))
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed document tree (for debugging)
	printNode(doc.Root, 0)
}

func printNode(node *html.Node, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	if node.Type == html.TextNode {
		if text := node.Text; text != "" {
			fmt.Printf("%stext: %s\n", indent, text)
		}
		return
	}

	// Print tag name and attributes
	fmt.Printf("%s<%s", indent, node.TagName)
	for name, value := range node.Attrs {
		fmt.Printf(" %s=\"%s\"", name, value)
	}
	fmt.Println(">")

	// Print children
	for _, child := range node.Children {
		printNode(child, depth+1)
	}
}
