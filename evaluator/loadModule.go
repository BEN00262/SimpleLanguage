package evaluator

import (
	"io/ioutil"

	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
)

// create a dependancy graph

func (eval *Evaluator) loadModule(modulePath string) {
	// ensure the filename exists --> also check for errors in the lexer and the parser too
	module, err := ioutil.ReadFile(modulePath)

	if err != nil {
		panic(err.Error())
	}

	lexer := InitLexer(string(module))
	parser := InitParser(lexer.Lex())

	for _, node := range parser.Parse().Nodes {
		_, err = eval.walkTree(node)

		if err != nil {
			panic(err.Error())
		}
	}
}

func (eval *Evaluator) _eval(codeString string) (result interface{}) {
	lexer := InitLexer(codeString)
	parser := InitParser(lexer.Lex())
	var err error

	for _, node := range parser.Parse().Nodes {

		switch node.(type) {
		case ExpressionNode:
			{
				result, err = eval.walkTree(node)

				if err != nil {
					panic(err.Error())
				}
			}
		default:
			panic("Expected an expression")
		}
		break
	}
	return
}
