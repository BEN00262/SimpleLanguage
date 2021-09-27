package evaluator

import (
	"io/ioutil"

	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
)

// create a dependancy graph

func (eval *Evaluator) loadModule(modulePath string) ExceptionNode {
	// ensure the filename exists --> also check for errors in the lexer and the parser too
	module, err := ioutil.ReadFile(modulePath)

	if err != nil {
		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: err.Error(),
		}
	}

	lexer := InitLexer(string(module))
	parser := InitParser(lexer.Lex())

	for _, node := range parser.Parse().Nodes {
		_, exception := eval.walkTree(node)

		if exception.Type != NO_EXCEPTION {
			return exception
		}
	}

	return ExceptionNode{Type: NO_EXCEPTION}
}

func (eval *Evaluator) _eval(codeString string) (result interface{}, exception ExceptionNode) {
	lexer := InitLexer(codeString)
	parser := InitParser(lexer.Lex())

	for _, node := range parser.Parse().Nodes {

		switch node.(type) {
		case ExpressionNode:
			{
				result, exception = eval.walkTree(node)

				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}
			}
		default:
			return nil, ExceptionNode{
				Type:    INVALID_OPERATION_EXCEPTION,
				Message: "Expected an expression",
			}
		}
		break
	}
	return
}
