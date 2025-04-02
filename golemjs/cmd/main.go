package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("GolemJS - JavaScript Interpreter")
	fmt.Println("Version 0.1.0")

	if len(os.Args) > 1 {
		fmt.Printf("Input file: %s\n", os.Args[1])
	} else {
		fmt.Println("No input file specified")
		fmt.Println("Usage: golemjs <filename.js>")
	}
}
