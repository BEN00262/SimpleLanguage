package evaluator

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/BEN00262/simpleLang/exceptions"
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
	IMPORTED_MODULE // has its own symbols table that is copied around
)

// add a type whether its exposed or not
// we load into a new scope --> "VALUE" IsExported Boolean Value
type SymbolTableValue struct {
	Type             SymbolTableValueType
	IsExported       bool
	Value            interface{}
	ReferenceToScope *symTable.ContextValue // this will not be used always --> used esp in modules and stuff
}

// create a runtime (used for other things, like creating standalone binaries :))
// use the language to mask away malware
// actually write my first ransomware using this language

// file access ( file_open ) --> returns a string node --> then we can call all the other shit on this
// what are we doing here we need to work with pointers to the values

type Evaluator struct {
	program      *ProgramNode
	symbolsTable *symTable.SymbolsTable
	IsExported   bool // used to show the current assignment is exported
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

// we need a way to inform of the return node stuff
// we can use the exceptions i think

func (eval *Evaluator) executeFunctionCode(code []interface{}) (interface{}, ExceptionNode) {
	var returnValue interface{}
	var exception ExceptionNode

	for _, _code := range code {
		returnValue, exception = eval.walkTree(_code)

		if exception.Type != NO_EXCEPTION {
			return nil, exception
		}

		switch _val := returnValue.(type) {
		case ReturnNode:
			{
				// we need a way to break the loop of execution
				return _val.Expression, ExceptionNode{Type: INTERNAL_RETURN_EXCEPTION}
			}
		}
	}

	return returnValue, ExceptionNode{Type: NO_EXCEPTION}
}

var (
	INTERPOLATION = regexp.MustCompile(`{((\s*?.*?)*?)}`)
)

func evaluateAndReturn(_response interface{}) (interface{}, ExceptionNode) {
	if _exception, ok := _response.(ExceptionNode); ok {
		return nil, _exception
	}

	return _response, ExceptionNode{Type: NO_EXCEPTION}
}

// return something
func doArithmetic(left ArthOp, operator string, right interface{}) (interface{}, ExceptionNode) {
	switch operator {
	case "+":
		{
			return evaluateAndReturn(left.Add(right))
		}
	case "-":
		{
			return evaluateAndReturn(left.Sub(right))
		}
	case "*":
		{
			return evaluateAndReturn(left.Mul(right))
		}
	case "%":
		{
			return evaluateAndReturn(left.Mod(right))
		}
	case "/":
		{
			return evaluateAndReturn(left.Div(right))
		}
	}

	// return an exception
	return nil, ExceptionNode{
		Type:    INVALID_OPERATOR_EXCEPTION,
		Message: fmt.Sprintf("Unsupported binary operator, '%s'", operator),
	}
}

// simply pass the error down the line
// until we find an error handler that handles it
func Compare(comp Comparison, op string, rhs interface{}) (BoolNode, ExceptionNode) {
	switch op {
	case "==":
		{
			// call the comparison stuff and return the value
			return comp.IsEqualTo(rhs), ExceptionNode{Type: NO_EXCEPTION}
		}
	case "!=":
		{
			_comp_ := comp.IsEqualTo(rhs)

			if _comp_.Value == 1 {
				_comp_.Value = 0
			} else {
				_comp_.Value = 1
			}

			return _comp_, ExceptionNode{Type: NO_EXCEPTION}
		}
	case "<=":
		{
			return comp.IsLessThanOrEqualsTo(rhs), ExceptionNode{Type: NO_EXCEPTION}
		}
	case ">=":
		{
			return comp.IsGreaterThanOrEqualsTo(rhs), ExceptionNode{Type: NO_EXCEPTION}
		}
	case ">":
		{
			return comp.IsGreaterThan(rhs), ExceptionNode{Type: NO_EXCEPTION}
		}
	case "<":
		{
			return comp.IsLessThan(rhs), ExceptionNode{Type: NO_EXCEPTION}
		}
	}

	// panic here the operation is unsupported
	// we return an error code buana i think thats a good way to throw stuff down the line
	return BoolNode{Value: 0}, ExceptionNode{
		Type:    INVALID_OPERATOR_EXCEPTION,
		Message: fmt.Sprintf("Unsupported comparison operator '%s'", op),
	}
}

// a function to perform string interpolation and return the string node
func (eval *Evaluator) _stringInterpolate(stringNode StringNode) (StringNode, ExceptionNode) {
	for _, stringBlock := range INTERPOLATION.FindAllStringSubmatch(stringNode.Value, -1) {
		if stringBlock != nil {
			_interpolated_string_ := ""
			// fetch the interpolator from the current context
			// we should actually evaluate it as an expression --> its gonna be slow AF
			// if u use it in a loop fuck u

			// evaluate the value and get the results
			// value, _ := eval.symbolsTable.GetFromContext(stringBlock[1])

			_value_, exception := eval._eval(stringBlock[1])

			if exception.Type != NO_EXCEPTION {
				return StringNode{}, exception
			}

			switch _value := _value_.(type) {
			case NumberNode:
				{
					// do the work and change the values
					_interpolated_string_ = _value.Value.Text(10)
				}
			case StringNode:
				{
					_interpolated_string_ = fmt.Sprintf("%s", _value.Value)
				}
			}

			stringNode.Value = strings.ReplaceAll(stringNode.Value, stringBlock[0], _interpolated_string_)
		}
	}

	return stringNode, ExceptionNode{Type: NO_EXCEPTION}
}

// do passes over the code inorder to use the documentation strings well for typechecking
func (eval *Evaluator) walkTree(node interface{}) (interface{}, ExceptionNode) {
	switch _node := node.(type) {
	case VariableNode:
		{
			_value, err := eval.symbolsTable.GetFromContext(_node.Value)

			// this one is a none existent value
			if err != nil {
				return nil, ExceptionNode{
					Type:    NAME_EXCEPTION,
					Message: fmt.Sprintf("'%s' is not defined", _node.Value),
				}
			}

			_parsedValue := (*_value).(SymbolTableValue)

			// this will remain that way --> we need a way to actually throw the exception
			return _parsedValue.Value, ExceptionNode{Type: NO_EXCEPTION}
		}
	case TryCatchNode:
		{
			// evaluate the try catch stuff
			// find a way to throw errors
			// this errors will be used then
			// if we get an exception node just pass it down the line
			// in the case of this one tuko poa
			// now check the result and find the exception

			return eval.evaluateTryCatchFinally(_node)
		}
	case RaiseExceptionNode:
		{
			// raise exception("something", "some explanation")
			// we just return the exeption
			// the result should be an exception node

			_result, _exception := eval.walkTree(_node.Exception)

			if _exception.Type != NO_EXCEPTION {
				return nil, _exception
			}

			if _extracted_exception, ok := _result.(ExceptionNode); ok {
				return nil, _extracted_exception
			}

			return nil, ExceptionNode{
				Type:    INVALID_EXCEPTION_EXCEPTION,
				Message: fmt.Sprintf("%#v is not an exception", _result),
			}
		}
	case ArrayNode:
		{
			// handle the array node shit
			// return stuff here
			// also implement a type check for arrays in the symbols table
			var _array_elements_ []interface{}

			for _, _element_ := range _node.Elements {
				_element, exception := eval.walkTree(_element_)

				// we check if the exception returned is not an ok exception if so just exit
				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}

				// we need to check if the element is of type exeception if it is cease the execution and find a catch
				// pass the error back until we reach a handler and if not just throw and exception

				_array_elements_ = append(_array_elements_, _element)
			}

			return ArrayNode{
				Elements: _array_elements_,
			}, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ExportVisibilityNode:
		{
			eval.IsExported = true
			defer (func(eval *Evaluator) {
				eval.IsExported = false
			})(eval)

			_, _exception := eval.walkTree(_node.Exported)
			return nil, _exception
		}
	case IFNode:
		{
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			_condition, _ := eval.walkTree(_node.Condition)

			_bool_condition := _condition.(BoolNode)

			if _bool_condition.Value == 1 {
				for _, _code := range _node.ThenBody {
					res, exception := eval.walkTree(_code)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// check for the return type
					switch _node_ := res.(type) {
					case BreakNode:
						{
							return BreakNode{}, ExceptionNode{Type: NO_EXCEPTION}
						}
					case ReturnNode:
						{
							return _node_, ExceptionNode{Type: NO_EXCEPTION}
						}
					}
				}

				return nil, ExceptionNode{Type: NO_EXCEPTION}
			} else {
				// we could have thrown an error in other languages but we cant here fuck
				for _, _code := range _node.ElseBody {
					res, exception := eval.walkTree(_code)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// check if the
					switch _node_ := res.(type) {
					case ReturnNode:
						{
							return res, ExceptionNode{Type: NO_EXCEPTION}
						}
					case BreakNode:
						{
							// check the state we are in if it allows this
							return _node_, ExceptionNode{Type: NO_EXCEPTION}
						}
					}
				}
			}

			return nil, ExceptionNode{Type: NO_EXCEPTION}
		}
	case BlockNode:
		{
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			for _, _code := range _node.Code {
				// we can throw errors in golang
				ret, exception := eval.walkTree(_code)

				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}

				// ensure the return is not a break node or return node if so just return a nil
				switch _node := ret.(type) {
				case ReturnNode:
					{
						return ReturnNode{
							Expression: _node.Expression,
						}, ExceptionNode{Type: NO_EXCEPTION}
					}
				case BreakNode:
					{
						return BreakNode{}, ExceptionNode{Type: NO_EXCEPTION}
					}
				}
			}
		}
	case BreakNode:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case NilNode:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ReturnNode:
		{
			//
			_ret, exception := eval.walkTree(_node.Expression)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			return ReturnNode{Expression: _ret}, ExceptionNode{Type: NO_EXCEPTION}
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
							retToken, exception := eval.walkTree(_code)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, ExceptionNode{Type: NO_EXCEPTION}
								}
							}
						}
					}
				}
			case FOR_NODE:
				{
					_initialization := _node.Initialization.(Assignment)

					_, exception := eval.walkTree(_initialization)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// get the condition
					_condition, exception := eval.walkTree(_node.Condition)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// convert the condition to a BoolNode and check the return value
					_condition_bool_ := _condition.(BoolNode)

					if _condition_bool_.Value == 0 {
						// this is a false thing
						// do not proceed anywhere
						return nil, ExceptionNode{Type: NO_EXCEPTION}
					}

					isExecuting := true

					for isExecuting && _condition_bool_.Value == 1 {
						for _, _code := range _node.ForBody {
							retToken, exception := eval.walkTree(_code)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, ExceptionNode{Type: NO_EXCEPTION}
								}
							}
						}

						_increment_return_value_, exception := eval.walkTree(_node.Increment)

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						_increment_return_value := _increment_return_value_.(NumberNode)

						eval.symbolsTable.PushToContext(_initialization.Lvalue, SymbolTableValue{
							Type: VALUE,
							Value: NumberNode{
								Value: _increment_return_value.Value,
							},
						})

						// re-evaluate the condition again
						_condition, exception = eval.walkTree(_node.Condition)

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						// convert the condition to a BoolNode and check the return value
						_condition_bool_ = _condition.(BoolNode)
					}

				}
			case WHILE_CONDITIONAL:
				{
					// the condition must evaluate to BoolNode inorder to be used here
					_condition, exception := eval.walkTree(_node.Condition)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// convert the condition to a BoolNode and check the return value
					_condition_bool_ := _condition.(BoolNode)

					if _condition_bool_.Value == 0 {
						// this is a false thing
						// do not proceed anywhere
						return nil, ExceptionNode{Type: NO_EXCEPTION}
					}

					isExecuting := true

					for isExecuting && _condition_bool_.Value == 1 {
						for _, _code := range _node.ForBody {
							retToken, exception := eval.walkTree(_code)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							// if the token is a break statement just exit the execution
							switch _node_ := retToken.(type) {
							case BreakNode:
								{
									isExecuting = false
								}
							case ReturnNode:
								{
									return _node_, ExceptionNode{Type: NO_EXCEPTION}
								}
							}
						}

						// re-evaluate the condition again
						_condition, exception = eval.walkTree(_node.Condition)

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						// convert the condition to a BoolNode and check the return value
						_condition_bool_ = _condition.(BoolNode)
					}
				}
			}

			return nil, ExceptionNode{Type: NO_EXCEPTION}
		}
	case StringNode:
		{
			// first check if the string is being interpolated if so interpolate it
			return eval._stringInterpolate(_node)
		}
	case IIFENode:
		{
			// we just call the anonymous function and parse the args
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			_function_decl_ := _node.Function

			// we get the value then execute the code here
			if _function_decl_.ParamCount != _node.ArgCount {
				return nil, ExceptionNode{
					Type:    ARITY_EXCEPTION,
					Message: fmt.Sprintf("IIFE function expected %d args but only %d args given", _function_decl_.ParamCount, _node.ArgCount),
				}
			}

			_ret, _exception := eval.executeFunctionCode(_function_decl_.Code)

			// check if the exception is of INTERNAL_RETURN_EXCEPTION if so just return the results
			if _exception.Type == INTERNAL_RETURN_EXCEPTION {
				return _ret, ExceptionNode{Type: NO_EXCEPTION}
			}

			return _ret, _exception
		}
	case NumberNode:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ExpressionNode:
		{
			return eval.walkTree(_node.Expression)
		}
	case BinaryNode:
		{
			// we have to check the binary Node to ascertain
			// return the evaluation here
			lhs, exception := eval.walkTree(_node.Lhs)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			rhs, exception := eval.walkTree(_node.Rhs)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			// additions allowed --> string + number / number + string / number + number
			// we just pass them to the interface stuff

			// return doArithmetic(lhs, _node.Operator, rhs)

			switch _lhs := lhs.(type) {
			case NumberNode:
				{
					return doArithmetic(&_lhs, _node.Operator, rhs)
				}
			case StringNode:
				{
					return doArithmetic(&_lhs, _node.Operator, rhs)
				}
			}

			// we should not panic buana in this system
			return nil, ExceptionNode{
				Type:    INVALID_OPERATION_EXCEPTION,
				Message: fmt.Sprintf("Invalid operation %#v", _node),
			}
		}
	case FunctionDecl:
		{
			eval.symbolsTable.PushToContext(_node.Name, SymbolTableValue{
				Type:       FUNCTION,
				IsExported: eval.IsExported,
				Value:      _node,
			})
		}
	case AnonymousFunction:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case FunctionCall:
		{
			// ideally we are using the top level scope but not for namespaces
			// how tf are we going to solve this buana
			// we need a way to inject a context here --> thats what it is
			// how will it work --> set a global pointer to sth
			// node.Name should not be a string it should be an interface i think
			// so that we can call it correctly
			/*
				print(name.juma)
				name.juma(7,8) --> inject the context here and start using them
			*/

			// function, err := eval.symbolsTable.GetFromContext(_node.Name)

			function, _exception := eval.walkTree(_node.Name)

			if _exception.Type != NO_EXCEPTION {
				return nil, _exception
			}

			_function := (function).(SymbolTableValue)

			if _function.Type != FUNCTION && _function.Type != EXTERNALFUNC {
				return nil, ExceptionNode{
					Type:    NAME_EXCEPTION,
					Message: fmt.Sprintf("'%#v' is not a function", _function.Value),
				}
			}

			if _function.Type == EXTERNALFUNC {
				// this is an externa function
				// just call the function

				_function_decl_ := _function.Value.(ExternalFunctionNode)

				if _function_decl_.ParamCount != _node.ArgCount {
					// throw an error here
					return nil, ExceptionNode{
						Type:    ARITY_EXCEPTION,
						Message: fmt.Sprintf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount),
					}
				}

				// evaluate each argument --> i think
				var _args []*interface{}

				// get out the execution of the code when the return occurs
				// we evaluate the args -->

				for _, _myArg := range _node.Args {
					_val, exception := eval.walkTree(_myArg)

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}

					// get the type of the _val
					switch _val_ := _val.(type) {
					case SymbolTableValue:
						{
							_args = append(_args, &_val_.Value)
						}
					default:
						{
							_args = append(_args, &_val)
						}
					}
				}

				return _function_decl_.Function(_args...)
			}

			var returnValue interface{}
			var exception ExceptionNode

			// we are pushing the scopes here
			// push a normal scope --> otherwise we should kind of get a refernce to the main scope
			// i think of the function

			// before we jump into this context please place the other context
			if _function.ReferenceToScope != nil {
				eval.symbolsTable.CopyContextToTop(*_function.ReferenceToScope)
				defer eval.symbolsTable.PopContext()
			}

			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			// TODO: make this code DRY laters
			switch _function_decl_ := _function.Value.(type) {
			case FunctionDecl:
				{
					if _function_decl_.ParamCount != _node.ArgCount {
						return nil, ExceptionNode{
							Type:    ARITY_EXCEPTION,
							Message: fmt.Sprintf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount),
						}
					}

					// push the function args into the current scope
					for _, Param := range _function_decl_.Params {
						// find the _args and push them into the current
						// if we walk we find the values
						res, exception := eval.walkTree(_node.Args[Param.Position])

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						valueType := VALUE

						switch _ret_ := res.(type) {
						case AnonymousFunction:
							valueType = FUNCTION
						case FunctionDecl:
							valueType = FUNCTION
						case ArrayNode:
							valueType = ARRAY
						case SymbolTableValue:
							res = _ret_.Value
						}

						// we push to the context here --> ideally we should have a way to
						eval.symbolsTable.PushToContext(Param.Key, SymbolTableValue{
							Type:  valueType,
							Value: res,
						})
					}

					// push all the args to the new scope here -> we should switch right here
					// i think or am wrong i dunno

					// this is the place we are executing the functions
					returnValue, exception = eval.executeFunctionCode(_function_decl_.Code)

					// check if the exception is of INTERNAL_RETURN_EXCEPTION if so just return the results
					if exception.Type == INTERNAL_RETURN_EXCEPTION {
						return returnValue, ExceptionNode{Type: NO_EXCEPTION}
					}

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}
				}
			case AnonymousFunction:
				{
					if _function_decl_.ParamCount != _node.ArgCount {
						return nil, ExceptionNode{
							Type:    ARITY_EXCEPTION,
							Message: fmt.Sprintf("'%s' expected %d args but only %d args given", _node.Name, _function_decl_.ParamCount, _node.ArgCount),
						}
					}

					// push the function args into the current scope
					for _, Param := range _function_decl_.Params {
						// find the _args and push them into the current
						// if we walk we find the values
						res, exception := eval.walkTree(_node.Args[Param.Position])

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						valueType := VALUE

						switch res.(type) {
						case AnonymousFunction:
							valueType = FUNCTION
						case ArrayNode:
							valueType = ARRAY
						}

						// we push to the context here --> ideally we should have a way to
						eval.symbolsTable.PushToContext(Param.Key, SymbolTableValue{
							Type:  valueType,
							Value: res,
						})
					}

					returnValue, exception = eval.executeFunctionCode(_function_decl_.Code)

					// check if the exception is of INTERNAL_RETURN_EXCEPTION if so just return the results
					if exception.Type == INTERNAL_RETURN_EXCEPTION {
						return returnValue, ExceptionNode{Type: NO_EXCEPTION}
					}

					if exception.Type != NO_EXCEPTION {
						return nil, exception
					}
				}
			}

			return returnValue, ExceptionNode{Type: NO_EXCEPTION}
		}
	case BoolNode:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ObjectAccessor:
		{
			// walk the tree
			// parent and the children
			// first get the parent context check its type

			_parent_, _error_ := eval.symbolsTable.GetFromContext(_node.Parent)

			if _error_ != nil {
				return nil, ExceptionNode{Type: NAME_EXCEPTION, Message: _error_.Error()}
			}

			if _node.Child == "" {
				// we dont have a kid here
				// just return whatever we had
				return (*_parent_).(SymbolTableValue), ExceptionNode{Type: NO_EXCEPTION}
			}

			// we have the parent --> check for now if its a module
			if _parent_converted_, ok := (*_parent_).(SymbolTableValue); ok {
				// we have the parent check the type
				if _parent_converted_.Type != IMPORTED_MODULE {
					return nil, ExceptionNode{
						Type:    INVALID_OPERATION_EXCEPTION,
						Message: "We only support accessing for module imports only as of now",
					}
				}

				// we have the converted type
				// get the child
				if _import_, ok := _parent_converted_.Value.(ImportModule); ok {
					// return the value found in this context
					// after finishing just dump the current state --> find a much better way
					eval.symbolsTable.CopyContextToTop(_import_.context)
					defer eval.symbolsTable.PopContext()
					// get the child value
					_child_, err := eval.symbolsTable.GetFromContext(_node.Child)

					if err != nil {
						return nil, ExceptionNode{Type: NAME_EXCEPTION, Message: _error_.Error()}
					}

					// the *_child is an actual symbols table value
					_child := (*_child_).(SymbolTableValue)

					if !_child.IsExported {
						return nil, ExceptionNode{
							Type:    ACCESS_VIOLATION_EXCEPTION,
							Message: "Access violation",
						}
					}

					_child.ReferenceToScope = &_import_.context
					return _child, ExceptionNode{Type: NO_EXCEPTION}
				}
			}
		}
	case CommentNode:
		{
			// we dont return a comment node for now just assume it
			// we dont care about the comment
			return nil, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ArrayAccessorNode:
		{
			_index_of_element_, exception := eval.walkTree(_node.Index)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			// we should also check the type of the stuff

			if _index_, ok := _index_of_element_.(NumberNode); ok {
				_array, _exception := eval.walkTree(_node.Array)

				if _exception.Type != NO_EXCEPTION {
					return nil, _exception
				}

				_array_ := _array.(SymbolTableValue)

				if _implemented, ok := _array_.Value.(Getter); ok {
					switch _node.Type {
					case NORMAL:
						{
							return _implemented.Get(_index_.Value.Int64()), ExceptionNode{Type: NO_EXCEPTION}
						}
					case RANGE:
						{
							_end_index_, exception := eval.walkTree(_node.EndIndex)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							if _eIndex_, ok := _end_index_.(NumberNode); ok {
								return _implemented.Range(_index_.Value.Int64(), _eIndex_.Value.Int64()), ExceptionNode{Type: NO_EXCEPTION}
							}
						}
					}
				}

				// fmt.Errorf("Failed to fetch element at the given index")
				return nil, ExceptionNode{
					Type:    INVALID_INDEX_EXCEPTION,
					Message: fmt.Sprintf("Failed to fetch element at the given index '%d'", _index_.Value),
				}
			}

			// ensure the _index_of_element is a number node else return an error node
			return nil, ExceptionNode{
				Type:    INVALID_OPERATION_EXCEPTION,
				Message: fmt.Sprint("Given index expression does not evaluate to a number"),
			}
		}
	case Assignment:
		{
			_value, _exception := eval.walkTree(_node.Rvalue)

			// if the _exception is not
			if _exception.Type != NO_EXCEPTION {
				return nil, _exception
			}

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
			case ASSIGNMENT, CONST_ASSIGNMENT:
				{
					// we push it here
					eval.symbolsTable.PushToContext(_node.Lvalue, SymbolTableValue{
						Type:       _type,
						IsExported: eval.IsExported,
						Value:      _value,
					})
				}
			case REASSIGNMENT:
				{
					// check for the constants in the parser
					// check here if the value is a constant just return an exception
					eval.symbolsTable.PushToParentContext(_node.Lvalue, SymbolTableValue{
						Type:  _type,
						Value: _value,
					})
				}
			}
		}
	case Import:
		{
			// we need to pass back the exception
			// we create something else our own stuff
			return nil, eval.loadModule(_node)
		}
	case ConditionNode:
		{
			// evaluate this stuff
			_lhs, exception := eval.walkTree(_node.Lhs)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			_rhs, exception := eval.walkTree(_node.Rhs)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			// start the switching here
			switch _lhs_ := _lhs.(type) {
			case NumberNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs)
				}
			case StringNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs)
				}
			case BoolNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs)
				}
			case NilNode:
				{
					return Compare(&_lhs_, _node.Operator, _rhs)
				}
			default:
				return nil, ExceptionNode{
					Type:    INVALID_OPERATION_EXCEPTION,
					Message: fmt.Sprintf("%#v does not implement the Comparison interface", _lhs_),
				}
			}
		}
	default:
		{
			/*
				ExceptionNode{
					Type:    INVALID_NODE_EXCEPTION,
					Message: fmt.Sprintf("Unknown node %#v", _node),
				}
			*/
			// throw errors
			return nil, ExceptionNode{
				Type: NO_EXCEPTION,
			}
		}
	}

	return nil, ExceptionNode{Type: NO_EXCEPTION}
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
	var exception ExceptionNode

	for _, node := range eval.program.Nodes {
		ret, exception = eval.walkTree(node)

		if exception.Type != NO_EXCEPTION {
			// we just return the exception node
			return exception
		}
	}

	return ret
}

func (eval *Evaluator) Evaluate() interface{} {
	var ret interface{}
	var exception ExceptionNode

	eval.symbolsTable.PushContext()
	defer eval.symbolsTable.PopContext()

	for _, node := range eval.program.Nodes {
		ret, exception = eval.walkTree(node)

		// we should not panic or return an error at all instead use the internal data structures
		// start on this kesho

		if exception.Type != NO_EXCEPTION {
			fmt.Printf("[ %s ] %s\n\n", exception.Type, exception.Message)
			return nil
		}
	}

	return ret
}
