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

	// ast := &Ast.AST{
	// 	Node: program,
	// }

	// fmt.Println(ast.Walk())

	// return nil

	evaluator := initEvaluator(program)

	evaluator.InitGlobalScope()
	LoadGlobalsToContext(evaluator)

	return evaluator.Evaluate(true)

	// ast := &Ast.AST{
	// 	Node: program,
	// }

	// fmt.Println(ast.Walk())

	// Typecheck

	// typeChecker := &Ast.AST{
	// 	Node: program,
	// }

	// fmt.Println(ast.Walk())

	// typechecker := Typecheck.NewTypeChecker(program)
	// typechecker.Walk()

	return nil
}
