package main

import "fmt"

func main() {
	lexer := initLexer("3")

	tokens := lexer.Lex()

	parser := initParser(tokens)

	program := parser.Parse()

	for _, node := range program.Nodes {
		switch parsedNode := node.(type) {
		case VariableNode:
			{
				fmt.Println(parsedNode.Value)
			}

		case NumberNode:
			{
				fmt.Println(parsedNode.Value)
			}
		}
	}
}
