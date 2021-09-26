package evaluator

import (
	. "github.com/BEN00262/simpleLang/parser"
)

// states ---> exception --> pass the exception
func (evaluator *Evaluator) evaluateRaiseNode() {

}

// evaluate a block --> check for specific states
func (evaluator *Evaluator) _evaluateBlock(block []interface{}, implicitSymTable bool) interface{} {
	if implicitSymTable {
		evaluator.symbolsTable.PushContext()
		defer evaluator.symbolsTable.PopContext()
	}

	// evaluate every single piece of code and use it
	for _, _code := range block {
		_return, _error := evaluator.walkTree(_code)

		// find a way to propagate errors down the chain
		if _error != nil {
			panic(_error)
		}

		// handle the return node
		// just break the loop and return it
		switch __return := _return.(type) {
		case ReturnNode:
			{
				return __return.Expression
			}
		case ExceptionNode:
			{
				// pass the exception node down
				return _return
			}
		}
	}

	return nil
}

func (evaluator *Evaluator) evaluateTryCatchFinally(_tryCatchNode TryCatchNode) (interface{}, error) {
	// evaluate this
	_tryEvaluation := evaluator._evaluateBlock(_tryCatchNode.Try, true)
	_result := _tryEvaluation

	// check the return value if its an Exception node handle it
	if exceptionThrown, ok := _tryEvaluation.(ExceptionNode); ok {
		evaluator.symbolsTable.PushContext()
		defer evaluator.symbolsTable.PopContext()

		// we have an exception lets evaluate the catch block and then do the finally block
		// we also need to inject the exeception into the symbols table here
		evaluator.symbolsTable.PushToContext(_tryCatchNode.Catch.Exception, SymbolTableValue{
			Type:  VALUE,
			Value: exceptionThrown,
		})

		_result = evaluator._evaluateBlock(_tryCatchNode.Catch.Body, false)
	}

	if len(_tryCatchNode.Finally) > 0 {
		// we have a finally block execute it
		_result = evaluator._evaluateBlock(_tryCatchNode.Finally, true)
	}

	return _result, nil
}
