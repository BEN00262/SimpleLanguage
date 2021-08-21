package main

import "fmt"

func main() {
	lexer := initLexer(`
	# void => void
	fun sample() {
	}
	`)

	tokens := lexer.Lex()

	// for index, token := range tokens {
	// 	fmt.Printf("%d ==> %#v\n", index, token)
	// }

	parser := initParser(tokens)
	program := parser.Parse()

	// fmt.Println(program.Nodes...)
	// use comments as docstrings to document the functions when running tools over the AST

	evaluator := initEvaluator(program)

	evaluator.InitGlobalScope()

	evaluator.InjectIntoGlobalScope("print", SymbolTableValue{
		Type: EXTERNALFUNC,
		Value: ExternalFunctionNode{
			Name:       "print",
			ParamCount: 0,
			Function: func(value ...interface{}) {
				// you can use ur function hapa sasa :)
				fmt.Println("Hurray we are in an external function here")
			},
		},
	})

	fmt.Println(evaluator.Evaluate())
}
