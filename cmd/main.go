package main

import (
	"fmt"
)

type Parser struct {
	Name string
}

func NewParser(name string) *Parser {
	return &Parser{
		Name: name,
	}
}

func main() {
	parser, err := NewParser("ExampleParser")
	if err != nil {
		fmt.Printf("Failed to create parser: %v\n", err)
		return
	}
	
	// Example of parsing a dependency file
	fmt.Printf("Created %s parser that supports: %v\n", 
		parser.Name, parser.SupportedFiles)
	
	// TODO: Add actual parsing logic
}