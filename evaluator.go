package main

import "fmt"

type SymbolTableValueType = int

const (
	FUNCTION SymbolTableValueType = iota + 1
	VALUE
	EXTERNALFUNC // called to the external runtime
	EXTVALUE
)

type SymbolTableValue struct {
	Type  SymbolTableValueType
	Value interface{}
}

type ExternalFunction = func(value ...interface{})

type Evaluator struct {
	program      *ProgramNode
	symbolsTable *SymbolsTable
}

func initEvaluator(program *ProgramNode) *Evaluator {
	return &Evaluator{
		program:      program,
		symbolsTable: initSymbolsTable(),
	}
}

func (eval *Evaluator) walkTree(node interface{}) (interface{}, error) {
	switch _node := node.(type) {
	case VariableNode:
		{
			_value, err := eval.symbolsTable.getFromContext(_node.Value)

			if err != nil {
				return nil, err
			}

			_parsedValue := _value.(SymbolTableValue)

			if _parsedValue.Type != VALUE {
				return nil, fmt.Errorf("%s is not a variable", _parsedValue.Value.(string))
			}

			return _parsedValue.Value, nil
		}
	case NumberNode:
		{
			return _node, nil
		}

	case ExpressionNode:
		{
			return eval.walkTree(_node.expression)
		}
	case BinaryNode:
		{
			// we have to check the binary Node to ascertain
			// return the evaluation here
			lhs, err := eval.walkTree(_node.Lhs)

			if err != nil {
				// throw an error here
				panic(err.Error())
			}

			rhs, err := eval.walkTree(_node.Rhs)

			if err != nil {
				panic(err.Error())
			}

			// now we just do the math
			switch _lhs := lhs.(type) {
			case NumberNode:
				{
					// also check the rhs
					switch _rhs := rhs.(type) {
					case NumberNode:
						{
							// check the operator
							// find a way to inject the operators
							switch _node.Operator {
							case "+":
								{
									return NumberNode{
										Value: _lhs.Value + _rhs.Value,
									}, nil
								}
							case "-":
								{
									return NumberNode{
										Value: _lhs.Value - _rhs.Value,
									}, nil
								}
							case "*":
								{
									return NumberNode{
										Value: _lhs.Value * _rhs.Value,
									}, nil
								}
							}
						}
					}
				}
			default:
				{
					fmt.Println("This is the lhs")
					fmt.Println(_lhs)
				}
			}

			panic(fmt.Errorf("Inavlid operation %#v", _node))
		}
	case FunctionDecl:
		{
			eval.symbolsTable.pushToContext(_node.Name, SymbolTableValue{
				Type:  FUNCTION,
				Value: _node,
			})
		}
	case FunctionCall:
		{
			function, err := eval.symbolsTable.getFromContext(_node.Name)

			if err != nil {
				return nil, err
			}

			// check if the value found is a function if not throw an error
			_function := function.(SymbolTableValue)

			if _function.Type != FUNCTION || _function.Type != EXTERNALFUNC {
				// throw an error here

			}

			if _function.Type == EXTERNALFUNC {
				// this is an externa function
				// just call the function

				_function_decl_ := _function.Value.(ExternalFunctionNode)

				if _function_decl_.ParamCount != _node.ArgCount {
					// throw an error here
					return nil, fmt.Errorf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount)
				}
				// call the function with the args
				// pass them as a list of interfaces
				_function_decl_.Function()
				return nil, nil
			}

			_function_decl_ := _function.Value.(FunctionDecl)

			// check if the Args and the Params match
			if _function_decl_.ParamCount != _node.ArgCount {
				// throw an error here
				return nil, fmt.Errorf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount)
			}

			var returnValue interface{}

			// check if the type is an external function

			eval.symbolsTable.pushContext()

			// start the execution here
			// inject the args into the current context first then execute the function

			for _, _code := range _function_decl_.Code {
				returnValue, err = eval.walkTree(_code)

				if err != nil {
					panic(err.Error())
				}
			}

			eval.symbolsTable.popContext()

			return returnValue, nil
		}
	case CommentNode:
		{
			return _node.comment, nil
		}
	case Assignment:
		{
			_value, _ := eval.walkTree(_node.Rvalue)

			eval.symbolsTable.pushToContext(_node.Lvalue, SymbolTableValue{
				Type:  VALUE,
				Value: _value,
			})
		}
	default:
		{
			fmt.Println(_node)
			return nil, fmt.Errorf("Uknown node %#v", _node)
		}
	}

	return nil, nil
}

// think about this very hard
func (eval *Evaluator) InitGlobalScope() {
	eval.symbolsTable.pushContext()
}

func (eval *Evaluator) InjectIntoGlobalScope(key string, value interface{}) {
	eval.symbolsTable.pushToContext(key, value)

}

// we need a way to inject functions into the scope ( print functions and stuff )
// the functions will just take an interface and work with it
func (eval *Evaluator) Evaluate() interface{} {
	var ret interface{}
	var err error

	eval.symbolsTable.pushContext()

	for _, node := range eval.program.Nodes {
		ret, err = eval.walkTree(node)

		if err != nil {
			panic(err.Error())
		}
	}

	eval.symbolsTable.popContext()
	return ret
}
