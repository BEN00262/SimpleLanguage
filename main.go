package main

import "fmt"

func main() {
	lexer := initLexer(`
	fun first() {
		4 + 5
	}

	fun second() {
		6 + 8
	}

	first() + second() + 8 * 6
	`)

	tokens := lexer.Lex()

	// for index, token := range tokens {
	// 	fmt.Printf("%d ==> %#v\n", index, token)
	// }

	parser := initParser(tokens)

	program := parser.Parse()

	// fmt.Println(program.Nodes...)

	evaluator := initEvaluator(program)

	fmt.Println(evaluator.Evaluate())

	// for _, node := range program.Nodes {
	// 	switch parsedNode := node.(type) {
	// 	case VariableNode:
	// 		{
	// 			fmt.Println(parsedNode.Value)
	// 		}

	// 	case NumberNode:
	// 		{
	// 			fmt.Println(parsedNode.Value)
	// 		}
	// 	case Assignment:
	// 		{
	// 			fmt.Println(parsedNode.Lvalue)
	// 		}
	// 	case CommentNode:
	// 		{
	// 			fmt.Println(parsedNode.comment)
	// 		}

	// 	case FunctionDecl:
	// 		{
	// 			fmt.Println(parsedNode.Name)
	// 		}

	// 	case FunctionCall:
	// 		{
	// 			fmt.Println("Function call node")
	// 			fmt.Println(parsedNode.Name)
	// 		}
	// 	case ExpressionNode:
	// 		{
	// 			fmt.Println(parsedNode.expression)
	// 		}
	// 	}
	// }
}
