package evaluator

import (
	. "github.com/BEN00262/simpleLang/parser"
)

// evaluate a block --> check for specific states
func (evaluator *Evaluator) _evaluateBlock(block []interface{}, implicitSymTable bool) (interface{}, ExceptionNode) {
	if implicitSymTable {
		evaluator.symbolsTable.PushContext()
		defer evaluator.symbolsTable.PopContext()
	}

	// evaluate every single piece of code and use it
	for _, _code := range block {
		_return, _error := evaluator.walkTree(_code)

		// find a way to propagate errors down the chain
		if _error.Type != NO_EXCEPTION {
			return nil, _error
		}

		// handle the return node
		// just break the loop and return it
		switch __return := _return.(type) {
		case ReturnNode:
			{
				return __return.Expression, ExceptionNode{Type: NO_EXCEPTION}
			}
		}
	}

	return nil, ExceptionNode{Type: NO_EXCEPTION}
}

func (evaluator *Evaluator) evaluateTryCatchFinally(_tryCatchNode TryCatchNode) (interface{}, ExceptionNode) {
	// evaluate this
	_tryEvaluation, _exceptionThrown := evaluator._evaluateBlock(_tryCatchNode.Try, true)
	_result := _tryEvaluation
	_exception := _exceptionThrown

	if _exceptionThrown.Type != NO_EXCEPTION {
		evaluator.symbolsTable.PushContext()
		defer evaluator.symbolsTable.PopContext()

		evaluator.symbolsTable.PushToContext(_tryCatchNode.Catch.Exception, SymbolTableValue{
			Type:  VALUE,
			Value: _exceptionThrown,
		})

		_result, _exception = evaluator._evaluateBlock(_tryCatchNode.Catch.Body, false)
	}

	if len(_tryCatchNode.Finally) > 0 {
		// we have a finally block execute it
		_result, _exception = evaluator._evaluateBlock(_tryCatchNode.Finally, true)
	}

	return _result, _exception
}
