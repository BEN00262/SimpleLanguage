package evaluator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
)

// create a dependancy graph

func (eval *Evaluator) loadModule(modulePath string) ExceptionNode {
	// ensure the filename exists --> also check for errors in the lexer and the parser too
	// have a * we dump to the global scope
	// otherwise we namespace

	_basePath, err := os.Executable()

	if err != nil {
		// we get the error code for not working here
		// what happens for now we use the path in the current directory

		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: err.Error(),
		}
	}

	// find a way to pass the values along
	// check if the module has a .happ extension if not add it
	if filepath.Ext(modulePath) != ".happ" {
		// append that
		modulePath += ".happ"
	}

	// should find the actual root folder of the stuff then get the files from there

	// system includes
	module, err := ioutil.ReadFile(filepath.Join(filepath.Dir(_basePath), "includes", modulePath))

	if err != nil {
		// we get the error code for not working here
		// what happens for now we use the path in the current directory

		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: err.Error(),
		}
	}

	// we redo this --> create a graph to prevent alot of shiets
	// make these very fast ---> i think
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
