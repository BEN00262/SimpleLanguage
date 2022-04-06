package evaluator

import (
	Lexer "github.com/BEN00262/simpleLang/lexer"
	Parser "github.com/BEN00262/simpleLang/parser"
	TypeCheck "github.com/BEN00262/simpleLang/typechecker"
)

func Interpreter(filePath string, codeString string) interface{} {
	lexer := Lexer.InitLexer(codeString)
	tokens := lexer.Lex()

	parser := Parser.InitParser(tokens, lexer.SplitCode)
	program := parser.Parse()

	typechecker := TypeCheck.NewTypeChecker(program)
	evaluator := initEvaluator(typechecker.Walk(), filePath)

	evaluator.InitGlobalScope()
	LoadGlobalsToContext(evaluator)

	return evaluator.Evaluate(true)
}
