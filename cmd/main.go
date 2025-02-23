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
	parser := NewParser("ExampleParser")
	fmt.Println("Parser created:", parser)
}