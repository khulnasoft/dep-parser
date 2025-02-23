package main

import (
	"fmt"
)

type Parser struct {
	Name string
}

// NewParser creates a new Parser instance with the specified name.
// It returns an error if the name is empty.
func NewParser(name string) (*Parser, error) {
	if name == "" {
		return nil, fmt.Errorf("parser name cannot be empty")
	}
	return &Parser{
		Name:           name,
		SupportedFiles: []string{},
		Version:        "1.0.0",
	}, nil
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