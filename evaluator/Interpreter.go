package evaluator

import (
	Lexer "github.com/BEN00262/simpleLang/lexer"
	Parser "github.com/BEN00262/simpleLang/parser"
)

func Interpreter(codeString string) interface{} {
	lexer := Lexer.InitLexer(codeString)
	tokens := lexer.Lex()

	parser := Parser.InitParser(tokens, lexer.SplitCode)
	program := parser.Parse()

	// fmt.Println(program)
	// fmt.Println()

	evaluator := initEvaluator(program)

	evaluator.InitGlobalScope()
	LoadGlobalsToContext(evaluator)

	return evaluator.Evaluate(true)
}
