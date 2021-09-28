package evaluator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
	. "github.com/BEN00262/simpleLang/symbolstable"
)

func (eval *Evaluator) _evaluateProgramNode(nodes []interface{}) ExceptionNode {
	for _, node := range nodes {
		_, exception := eval.walkTree(node)

		if exception.Type != NO_EXCEPTION {
			return exception
		}
	}

	return ExceptionNode{Type: NO_EXCEPTION}
}

// create a dependancy graph
type ImportModule struct {
	context ContextValue
}

func (eval *Evaluator) loadModule(module Import) ExceptionNode {
	// ensure the filename exists --> also check for errors in the lexer and the parser too
	// have a * we dump to the global scope
	// otherwise we namespace

	// create push our own context then use it later

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
	if filepath.Ext(module.FileName) != ".happ" {
		// append that
		module.FileName += ".happ"
	}

	// should find the actual root folder of the stuff then get the files from there

	// system includes
	importedModule, err := ioutil.ReadFile(filepath.Join(filepath.Dir(_basePath), "includes", module.FileName))

	if err != nil {
		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: err.Error(),
		}
	}

	lexer := InitLexer(string(importedModule))
	parser := InitParser(lexer.Lex())

	if module.Alias != "*" {
		eval.symbolsTable.PushContext()

		_exception := eval._evaluateProgramNode(parser.Parse().Nodes)

		if _exception.Type != NO_EXCEPTION {
			return _exception
		}

		_module_context := eval.symbolsTable.GetTopContext()

		eval.symbolsTable.PushToContext(module.Alias, SymbolTableValue{
			Type: IMPORTED_MODULE,
			Value: ImportModule{
				context: _module_context,
			},
		})

		return _exception
	}

	return eval._evaluateProgramNode(parser.Parse().Nodes)
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
