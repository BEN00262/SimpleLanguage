package evaluator

import (
	"fmt"

	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
)

func Interpreter(codeString string) interface{} {
	lexer := InitLexer(codeString)
	tokens := lexer.Lex()

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	parser := InitParser(tokens)
	program := parser.Parse()

	fmt.Println(program)
	fmt.Println()

	// // ast := initAST(program)

	// // ast.walk()

	// // for _, nodes := range program.Nodes {
	// // 	fmt.Println(nodes)
	// // }

	evaluator := initEvaluator(program)

	evaluator.InitGlobalScope()
	LoadGlobalsToContext(evaluator)

	evaluator.Evaluate()
	return nil
}
