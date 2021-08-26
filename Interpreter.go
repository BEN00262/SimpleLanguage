package main

import "fmt"

func Interpreter(codeString string) interface{} {
	lexer := initLexer(codeString)
	tokens := lexer.Lex()

	for _, token := range tokens {
		fmt.Println(token)
	}

	parser := initParser(tokens)
	program := parser.Parse()

	fmt.Println(program)

	// ast := initAST(program)

	// ast.walk()

	// for _, nodes := range program.Nodes {
	// 	fmt.Println(nodes)
	// }

	// evaluator := initEvaluator(program)

	// evaluator.InitGlobalScope()
	// LoadGlobalsToContext(evaluator)

	// return evaluator.Evaluate()
	return nil
}
