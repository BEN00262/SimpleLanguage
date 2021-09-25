package evaluator

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/BEN00262/simpleLang/parser"
	symTable "github.com/BEN00262/simpleLang/symbolstable"
)

type SymbolTableValueType = int

const (
	FUNCTION SymbolTableValueType = iota + 1
	VALUE
	ARRAY
	EXTERNALFUNC // called to the external runtime
	EXTVALUE
)

type SymbolTableValue struct {
	Type  SymbolTableValueType
	Value interface{}
}

// create a runtime (used for other things, like creating standalone binaries :))
// use the language to mask away malware
// actually write my first ransomware using this language

// file access ( file_open ) --> returns a string node --> then we can call all the other shit on this
// what are we doing here we need to work with pointers to the values

type Evaluator struct {
	program      *ProgramNode
	symbolsTable *symTable.SymbolsTable
}

func initEvaluator(program *ProgramNode) *Evaluator {
	return &Evaluator{
		program:      program,
		symbolsTable: symTable.InitSymbolsTable(),
	}
}

// create a method to be used by the REPL
func NewEvaluatorContext() *Evaluator {
	eval := &Evaluator{
		symbolsTable: symTable.InitSymbolsTable(),
	}

	eval.symbolsTable.PushContext()

	return eval
}

func (eval *Evaluator) ReplExecute(program *ProgramNode) interface{} {
	eval.program = program
	return eval.replEvaluate()
}

func (eval *Evaluator) TearDownRepl() {
	eval.symbolsTable.PopContext()
}

func (eval *Evaluator) executeFunctionCode(code []interface{}) interface{} {
	var returnValue interface{}
	var err error

	for _, _code := range code {
		returnValue, err = eval.walkTree(_code)

		if err != nil {
			panic(err.Error())
		}

		switch _val := returnValue.(type) {
		case ReturnNode:
			{
				return _val.Expression
			}
		}
	}

	return returnValue
}

var (
	INTERPOLATION = regexp.MustCompile(`{((\s*?.*?)*?)}`)
)

func Compare(comp Comparison, op string, rhs interface{}) BoolNode {
	switch op {
	case "==":
		{
			// call the comparison stuff and return the value
			return comp.IsEqualTo(rhs)
		}
	case "!=":
		{
			_comp_ := comp.IsEqualTo(rhs)

			if _comp_.Value == 1 {
				_comp_.Value = 0
			} else {
				_comp_.Value = 1
			}

			return _comp_
		}
	case "<=":
		{
			return comp.IsLessThanOrEqualsTo(rhs)
		}
	case ">=":
		{
			return comp.IsGreaterThanOrEqualsTo(rhs)
		}
	case ">":
		{
			return comp.IsGreaterThan(rhs)
		}
	case "<":
		{
			return comp.IsLessThan(rhs)
		}
	}

	// panic here the operation is unsupported
	return BoolNode{
		Value: 0,
	}
}

// a function to perform string interpolation and return the string node
func (eval *Evaluator) _stringInterpolate(stringNode StringNode) StringNode {
	for _, stringBlock := range INTERPOLATION.FindAllStringSubmatch(stringNode.Value, -1) {
		if stringBlock != nil {
			_interpolated_string_ := ""
			// fetch the interpolator from the current context
			// we should actually evaluate it as an expression --> its gonna be slow AF
			// if u use it in a loop fuck u

			// evaluate the value and get the results
			// value, _ := eval.symbolsTable.GetFromContext(stringBlock[1])

			_value_ := eval._eval(stringBlock[1])

			switch _value := _value_.(type) {
			case NumberNode:
				{
					// do the work and change the values
					_interpolated_string_ = fmt.Sprintf("%d", _value.Value)
				}
			case StringNode:
				{
					_interpolated_string_ = fmt.Sprintf("%s", _value.Value)
				}
			}

			stringNode.Value = strings.ReplaceAll(stringNode.Value, stringBlock[0], _interpolated_string_)
		}
	}

	return stringNode
}

// do passes over the code inorder to use the documentation strings well for typechecking
func (eval *Evaluator) walkTree(node interface{}) (interface{}, error) {
	switch _node := node.(type) {
	case VariableNode:
		{
			_value, err := eval.symbolsTable.GetFromContext(_node.Value)

			if err != nil {
				return nil, err
			}

			_parsedValue := (*_value).(SymbolTableValue)

			return _parsedValue.Value, nil
		}
	case ArrayNode:
		{
			// handle the array node shit
			// return stuff here
			// also implement a type check for arrays in the symbols table
			var _array_elements_ []interface{}

			for _, _element_ := range _node.Elements {
				_element, err := eval.walkTree(_element_)

				if err != nil {
					return nil, err
				}

				_array_elements_ = append(_array_elements_, _element)
			}

			return ArrayNode{
				Elements: _array_elements_,
			}, nil
		}
	case IFNode:
		{
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			_condition, _ := eval.walkTree(_node.Condition)

			_bool_condition := _condition.(BoolNode)

			if _bool_condition.Value == 1 {
				for _, _code := range _node.ThenBody {
					res, err := eval.walkTree(_code)

					if err != nil {
						return nil, err
					}

					// check for the return type
					switch _node_ := res.(type) {
					case BreakNode:
						{
							return BreakNode{}, nil
						}
					case ReturnNode:
						{
							return _node_, nil
						}
					}
				}

				return nil, nil
			} else {
				// we could have thrown an error in other languages but we cant here fuck
				for _, _code := range _node.ElseBody {
					res, err := eval.walkTree(_code)

					if err != nil {
						return nil, err
					}

					// check if the
					switch _node_ := res.(type) {
					case ReturnNode:
						{
							return res, nil
						}
					case BreakNode:
						{
							// check the state we are in if it allows this
							return _node_, nil
						}
					}
				}
			}

			return nil, nil
		}
	case BlockNode:
		{
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			for _, _code := range _node.Code {
				// we can throw errors in golang
				ret, err := eval.walkTree(_code)

				if err != nil {
					return nil, err
				}

				// ensure the return is not a break node or return node if so just return a nil
				switch _node := ret.(type) {
				case ReturnNode:
					{
						return ReturnNode{
							Expression: _node.Expression,
						}, nil
					}
				case BreakNode:
					{
						return BreakNode{}, nil
					}
				default:
					{
						return nil, nil
					}
				}
			}
		}
	case BreakNode:
		{
			return _node, nil
		}
	case NilNode:
		{
			return _node, nil
		}
	case ReturnNode:
		{
			_ret, err := eval.walkTree(_node.Expression)

			if err != nil {
				return nil, err
			}

			return ReturnNode{Expression: _ret}, nil
		}
	case ForNode:
		{
			// evaluate a for node
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			// do our thing
			switch _node.Type {
			case WHILE_FOREVER:
				{
					// we just execute the code forever until we get a break statement and exit
					// execute this over and over again
					isExecuting := true

					for isExecuting {
						for _, _code := range _node.ForBody {
							retToken, err := eval.walkTree(_code)

							if err != nil {
								return nil, err
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, nil
								}
							}
						}
					}
				}
			case FOR_NODE:
				{
					_initialization := _node.Initialization.(Assignment)

					_, err := eval.walkTree(_initialization)

					if err != nil {
						return nil, err
					}

					// get the condition
					_condition, err := eval.walkTree(_node.Condition)

					if err != nil {
						return nil, err
					}

					// convert the condition to a BoolNode and check the return value
					_condition_bool_ := _condition.(BoolNode)

					if _condition_bool_.Value == 0 {
						// this is a false thing
						// do not proceed anywhere
						return nil, nil
					}

					isExecuting := true

					for isExecuting && _condition_bool_.Value == 1 {
						for _, _code := range _node.ForBody {
							retToken, err := eval.walkTree(_code)

							if err != nil {
								return nil, err
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, nil
								}
							}
						}

						_increment_return_value_, err := eval.walkTree(_node.Increment)

						if err != nil {
							return nil, err
						}

						_increment_return_value := _increment_return_value_.(NumberNode)

						eval.symbolsTable.PushToContext(_initialization.Lvalue, SymbolTableValue{
							Type: VALUE,
							Value: NumberNode{
								Value: _increment_return_value.Value,
							},
						})

						// re-evaluate the condition again
						_condition, err = eval.walkTree(_node.Condition)

						if err != nil {
							return nil, err
						}

						// convert the condition to a BoolNode and check the return value
						_condition_bool_ = _condition.(BoolNode)
					}

				}
			case WHILE_CONDITIONAL:
				{
					// the condition must evaluate to BoolNode inorder to be used here
					_condition, err := eval.walkTree(_node.Condition)

					if err != nil {
						return nil, err
					}

					// convert the condition to a BoolNode and check the return value
					_condition_bool_ := _condition.(BoolNode)

					if _condition_bool_.Value == 0 {
						// this is a false thing
						// do not proceed anywhere
						return nil, nil
					}

					isExecuting := true

					for isExecuting && _condition_bool_.Value == 1 {
						for _, _code := range _node.ForBody {
							retToken, err := eval.walkTree(_code)

							if err != nil {
								return nil, err
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, nil
								}
							}
						}

						// re-evaluate the condition again
						_condition, err = eval.walkTree(_node.Condition)

						if err != nil {
							return nil, err
						}

						// convert the condition to a BoolNode and check the return value
						_condition_bool_ = _condition.(BoolNode)
					}
				}
			}

			return nil, nil
		}
	case StringNode:
		{
			// first check if the string is being interpolated if so interpolate it
			return eval._stringInterpolate(_node), nil
		}
	case IIFENode:
		{
			// we just call the anonymous function and parse the args
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			_function_decl_ := _node.Function

			// we get the value then execute the code here
			if _function_decl_.ParamCount != _node.ArgCount {
				return nil, fmt.Errorf("IIFE function expected %d args but only %d args given", _function_decl_.ParamCount, _node.ArgCount)
			}

			return eval.executeFunctionCode(_function_decl_.Code), nil
		}
	case NumberNode:
		{
			return _node, nil
		}
	case ExpressionNode:
		{
			return eval.walkTree(_node.Expression)
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

			// additions allowed --> string + number / number + string / number + number

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

							case "%":
								{
									return NumberNode{
										Value: _lhs.Value % _rhs.Value,
									}, nil
								}

							case "/":
								{
									return NumberNode{
										Value: _lhs.Value / _rhs.Value,
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

			panic(fmt.Errorf("Invalid operation %#v", _node))
		}
	case FunctionDecl:
		{
			eval.symbolsTable.PushToContext(_node.Name, SymbolTableValue{
				Type:  FUNCTION,
				Value: _node,
			})
		}
	case AnonymousFunction:
		{
			return _node, nil
		}
	case FunctionCall:
		{
			function, err := eval.symbolsTable.GetFromContext(_node.Name)

			if err != nil {
				return nil, err
			}

			// check if the value found is a function if not throw an error
			_function := (*function).(SymbolTableValue)

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

				// evaluate each argument --> i think
				var _args []*interface{}

				// get out the execution of the code when the return occurs
				// we evaluate the args -->

				for _, _myArg := range _node.Args {
					_val, err := eval.walkTree(_myArg)

					if err != nil {
						panic(err.Error())
					}

					// get the type of the _val
					switch _val_ := _val.(type) {
					case ReturnNode:
						{
							// we break out of the function execution with the given thing
							// print this value
							// fmt.Println(_val)
							return _val_.Expression, nil
						}
					}

					_args = append(_args, &_val)
				}

				return _function_decl_.Function(_args...), nil
			}

			var returnValue interface{}
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			switch _function_decl_ := _function.Value.(type) {
			case FunctionDecl:
				{
					if _function_decl_.ParamCount != _node.ArgCount {
						return nil, fmt.Errorf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount)
					}

					// push the function args into the current scope
					for _, Param := range _function_decl_.Params {
						// find the _args and push them into the current
						// if we walk we find the values
						res, err := eval.walkTree(_node.Args[Param.Position])

						if err != nil {
							panic(err.Error())
						}

						valueType := VALUE

						switch res.(type) {
						case AnonymousFunction:
							{
								valueType = FUNCTION
							}
						case ArrayNode:
							{
								valueType = ARRAY
							}
						}

						eval.symbolsTable.PushToContext(Param.Key, SymbolTableValue{
							Type:  valueType,
							Value: res,
						})
					}

					// this is the place we are executing the functions
					returnValue = eval.executeFunctionCode(_function_decl_.Code)
				}
			case AnonymousFunction:
				{
					if _function_decl_.ParamCount != _node.ArgCount {
						return nil, fmt.Errorf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount)
					}

					returnValue = eval.executeFunctionCode(_function_decl_.Code)
				}
			}

			return returnValue, nil
		}
	case BoolNode:
		{
			return _node, nil
		}
	case CommentNode:
		{
			return _node, nil
		}
	case ArrayAccessorNode:
		{
			_index_of_element_, err := eval.walkTree(_node.Index)

			if err != nil {
				return nil, err
			}

			// we should also check the type of the stuff

			if _index_, ok := _index_of_element_.(NumberNode); ok {
				_array_, err := eval.symbolsTable.GetFromContext(_node.Array)

				if err != nil {
					return nil, err
				}

				_array_symbols_table_ := (*_array_).(SymbolTableValue)

				if _implemented, ok := _array_symbols_table_.Value.(Getter); ok {

					switch _node.Type {
					case NORMAL:
						{
							return _implemented.Get(_index_.Value), nil
						}
					case RANGE:
						{
							_end_index_, err := eval.walkTree(_node.EndIndex)

							if err != nil {
								return nil, fmt.Errorf("The end expression failed to evaluate")
							}

							if _eIndex_, ok := _end_index_.(NumberNode); ok {
								return _implemented.Range(_index_.Value, _eIndex_.Value), nil
							}
						}
					}
				}

				return nil, fmt.Errorf("Failed to fetch element at the given index")
			}

			// ensure the _index_of_element is a number node else return an error node
			return nil, fmt.Errorf("Given index expression does not evaluate to a number")
		}
	case Assignment:
		{
			_value, _ := eval.walkTree(_node.Rvalue)
			_type := VALUE

			switch _value.(type) {
			case AnonymousFunction:
				{
					_type = FUNCTION
				}
			case ArrayNode:
				{
					_type = ARRAY
				}
			}

			switch _node.Type {
			case ASSIGNMENT:
				{
					eval.symbolsTable.PushToContext(_node.Lvalue, SymbolTableValue{
						Type:  _type,
						Value: _value,
					})
				}
			case REASSIGNMENT:
				{
					eval.symbolsTable.PushToParentContext(_node.Lvalue, SymbolTableValue{
						Type:  _type,
						Value: _value,
					})
				}
			}
		}
	case Import:
		{
			eval.loadModule(_node.FileName)
		}
	case ConditionNode:
		{
			// evaluate this stuff
			_lhs, err := eval.walkTree(_node.Lhs)

			if err != nil {
				return nil, err
			}

			_rhs, err := eval.walkTree(_node.Rhs)

			if err != nil {
				return nil, err
			}

			// start the switching here
			switch _lhs_ := _lhs.(type) {
			case NumberNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs), nil
				}
			case StringNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs), nil
				}
			case BoolNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs), nil
				}
			case NilNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs), nil
				}
			default:
				return nil, fmt.Errorf("%#v does not implement the Comparison interface", _lhs_)
			}
		}
	default:
		{
			// fmt.Println(_node)
			return nil, nil //fmt.Errorf("Uknown node %#v", _node)
		}
	}

	return nil, nil
}

// think about this very hard
func (eval *Evaluator) InitGlobalScope() {
	eval.symbolsTable.PushContext()
}

func (eval *Evaluator) InjectIntoGlobalScope(key string, value interface{}) {
	eval.symbolsTable.PushToContext(key, value)

}

func (eval *Evaluator) replEvaluate() interface{} {
	var ret interface{}
	var err error

	for _, node := range eval.program.Nodes {
		ret, err = eval.walkTree(node)

		if err != nil {
			panic(err.Error())
		}
	}

	return ret
}

func (eval *Evaluator) Evaluate() interface{} {
	var ret interface{}
	var err error

	eval.symbolsTable.PushContext()

	for _, node := range eval.program.Nodes {
		ret, err = eval.walkTree(node)

		if err != nil {
			panic(err.Error())
		}
	}

	eval.symbolsTable.PopContext()
	return ret
}
