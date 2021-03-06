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
	baseFilePath string
	program      *ProgramNode
	symbolsTable *symTable.SymbolsTable
	IsExported   bool // used to show the current assignment is exported
	evalCache    CodeCache
}

func initEvaluator(program *ProgramNode, filePath string) *Evaluator {
	return &Evaluator{
		baseFilePath: filePath,
		program:      program,
		symbolsTable: symTable.InitSymbolsTable(),
		evalCache:    CodeCache{},
	}
}

// create a method to be used by the REPL
func NewEvaluatorContext() *Evaluator {
	eval := &Evaluator{
		symbolsTable: symTable.InitSymbolsTable(),
		evalCache:    CodeCache{},
	}

	eval.symbolsTable.PushContext()

	return eval
}

func (eval *Evaluator) ReplExecute(program *ProgramNode) interface{} {
	eval.program = program
	return eval.Evaluate(false)
}

func (eval *Evaluator) TearDownRepl() {
	eval.symbolsTable.PopContext()
}

// we need a way to inform of the return node stuff
// we can use the exceptions i think

func (eval *Evaluator) executeFunctionCode(code []interface{}) (interface{}, ExceptionNode) {
	for _, _code := range code {
		returnValue, exception := eval.walkTree(_code)

		if exception.Type == INTERNAL_RETURN_EXCEPTION {
			return returnValue, ExceptionNode{Type: NO_EXCEPTION}
		}

		if exception.Type != NO_EXCEPTION {
			return nil, exception
		}

	}

	return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
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
			_value_, exception := eval._eval(stringBlock[1])

			if exception.Type != NO_EXCEPTION {
				// TODO: look at this again
				return StringNode{}, exception
			}

			stringNode.Value = strings.ReplaceAll(stringNode.Value, stringBlock[0], Print(_value_))
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

			return ((*_value).(SymbolTableValue)).Value, ExceptionNode{Type: NO_EXCEPTION}
		}
	case TryCatchNode:
		{
			return eval.evaluateTryCatchFinally(_node)
		}
	case RaiseExceptionNode:
		{
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
			var _array_elements_ []interface{}

			for _, _element_ := range _node.Elements {
				_element, exception := eval.walkTree(_element_)

				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}

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

			return eval.walkTree(_node.Exported)
		}
	case IFNode:
		{
			_condition, _ := eval.walkTree(_node.Condition)
			_bool_condition, ok := _condition.(BoolNode) // ensure thats this is a bool node btw

			if !ok {
				return nil, ExceptionNode{
					Type:    INVALID_OPERATION_EXCEPTION,
					Message: "Conditional expression did not evaluate to a boolean expression",
				}
			}

			if _bool_condition.True() {
				_res, _exception := eval.walkTree(_node.ThenBody)

				switch _exception.Type {
				case INTERNAL_RETURN_EXCEPTION, INTERNAL_BREAK_EXCEPTION:
					{
						// follow the execution of this
						return _res, _exception
					}
				case NO_EXCEPTION:
					{
						return nil, ExceptionNode{Type: NO_EXCEPTION}
					}
				default:
					return nil, _exception
				}
			} else {
				for _, _blocks := range _node.ElseBody {
					_res, _exception := eval.walkTree(_blocks)

					switch _exception.Type {
					case INTERNAL_RETURN_EXCEPTION, INTERNAL_BREAK_EXCEPTION:
						{
							return _res, _exception
						}
					case NO_EXCEPTION:
						{
							return nil, ExceptionNode{Type: NO_EXCEPTION}
						}
					default:
						return nil, _exception
					}
				}
			}
		}
	case BlockNode:
		{
			// scope
			eval.symbolsTable.PushContext()
			defer eval.symbolsTable.PopContext()

			for _, _code := range _node.Code {
				ret, exception := eval.walkTree(_code)

				if exception.Type == INTERNAL_RETURN_EXCEPTION {
					return ret, exception
				}

				if exception.Type == INTERNAL_BREAK_EXCEPTION {
					return nil, ExceptionNode{Type: INTERNAL_BREAK_EXCEPTION}
				}

				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}
			}
		}
	case BreakNode:
		{
			return nil, ExceptionNode{Type: INTERNAL_BREAK_EXCEPTION}
		}
	case NilNode:
		{
			return _node, ExceptionNode{Type: NO_EXCEPTION}
		}
	case ReturnNode:
		{
			_ret, exception := eval.walkTree(_node.Expression)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			return _ret, ExceptionNode{Type: INTERNAL_RETURN_EXCEPTION}
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

					for {
						// we start executing the body of the loop
						_ret, _exception := eval.walkTree(_node.ForBody)

						if _exception.Type == INTERNAL_RETURN_EXCEPTION {
							return _ret, _exception
						}

						if _exception.Type == INTERNAL_BREAK_EXCEPTION {
							return nil, ExceptionNode{Type: NO_EXCEPTION}
						}

						if _exception.Type != NO_EXCEPTION {
							return nil, _exception
						}
					}
				}
			case FOR_NODE:
				{

					/*
						for loop prologue
					*/
					_initialization, ok := _node.Initialization.(Assignment)

					if !ok {
						// raise an exception here
						return nil, ExceptionNode{
							Type:    INVALID_OPERATION_EXCEPTION,
							Message: "Invalid initialization for a 'for' node",
						}
					}

					// evaluate the initialization
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
					_condition_bool_, ok := _condition.(BoolNode)

					if !ok {
						return nil, ExceptionNode{
							Type:    INVALID_OPERATION_EXCEPTION,
							Message: "Expected an expression that evaluated to a boolean",
						}
					}

					if !_condition_bool_.True() {
						return nil, ExceptionNode{Type: NO_EXCEPTION}
					}

					/*
						end of for loop prologue
					*/

					for _condition_bool_.True() {
						_ret, _exception := eval.walkTree(_node.ForBody)

						if _exception.Type == INTERNAL_RETURN_EXCEPTION {
							return _ret, _exception
						}

						if _exception.Type == INTERNAL_BREAK_EXCEPTION {
							return nil, ExceptionNode{Type: NO_EXCEPTION}
						}

						if _exception.Type != NO_EXCEPTION {
							return nil, _exception
						}

						_increment_return_value_, exception := eval.walkTree(_node.Increment)

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						if _increment_return_value, ok := _increment_return_value_.(NumberNode); ok {
							// push the value to the symbols table again which is not good buana
							eval.symbolsTable.PushToContext(_initialization.Lvalue, SymbolTableValue{
								Type:  VALUE,
								Value: _increment_return_value,
							})

							// re-evaluate the condition again
							_condition, exception = eval.walkTree(_node.Condition)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							// convert the condition to a BoolNode and check the return value
							if _condition_bool_, ok = _condition.(BoolNode); !ok {
								return nil, ExceptionNode{
									Type:    INVALID_OPERATION_EXCEPTION,
									Message: "Expected an expression that evaluated to a boolean",
								}
							}

							continue
						}

						goto invalid_loop_exit
					}

					// goto to the end of the looping
					goto valid_loop_exit

				invalid_loop_exit:
					return nil, ExceptionNode{
						Type:    INVALID_OPERATION_EXCEPTION,
						Message: "Failed in loop post processing",
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
					// check for this if the node is not a conditional throw an error
					_condition_bool_, ok := _condition.(BoolNode)

					if !ok {
						return nil, ExceptionNode{
							Type:    INVALID_OPERATION_EXCEPTION,
							Message: "Expected an expression that evaluated to a boolean",
						}
					}

					if !_condition_bool_.True() {
						return nil, ExceptionNode{Type: NO_EXCEPTION}
					}

					for _condition_bool_.True() {

						_ret, _exception := eval.walkTree(_node.ForBody)

						if _exception.Type == INTERNAL_RETURN_EXCEPTION {
							return _ret, _exception
						}

						if _exception.Type == INTERNAL_BREAK_EXCEPTION {
							return nil, ExceptionNode{Type: NO_EXCEPTION}
						}

						if _exception.Type != NO_EXCEPTION {
							return nil, _exception
						}

						// re-evaluate the condition again
						_condition, exception = eval.walkTree(_node.Condition)

						if exception.Type != NO_EXCEPTION {
							return nil, exception
						}

						// convert the condition to a BoolNode and check the return value
						if _condition_bool_, ok = _condition.(BoolNode); !ok {
							// raise an exception here
							return nil, ExceptionNode{
								Type:    INVALID_OPERATION_EXCEPTION,
								Message: "Expected an expression that evaluated to a boolean",
							}
						}
					}
				}
			}

		valid_loop_exit:
			return nil, ExceptionNode{Type: NO_EXCEPTION}
		}
	case StringNode:
		{
			// first check if the string is being interpolated if so interpolate it
			return eval._stringInterpolate(_node)
		}
	case LogicalComparison:
		{
			_lhs, _exception := eval.walkTree(_node.Lhs)

			if _exception.Type != NO_EXCEPTION {
				return nil, ExceptionNode{Type: NO_EXCEPTION}
			}

			_lhs_boolean_, ok := _lhs.(BoolNode)

			if !ok {
				return nil, ExceptionNode{
					Type:    INVALID_OPERATION_EXCEPTION,
					Message: "lhs a boolean",
				}
			}

			_rhs, _exception := eval.walkTree(_node.Rhs)

			if _exception.Type != NO_EXCEPTION {
				return nil, ExceptionNode{Type: NO_EXCEPTION}
			}

			_rhs_boolean_, ok := _rhs.(BoolNode)

			if !ok {
				return nil, ExceptionNode{
					Type:    INVALID_OPERATION_EXCEPTION,
					Message: "lhs a boolean",
				}
			}

			// check the type of the comparator and do stuff
			switch _node.Type {
			case AND_COMPARATOR:
				{
					// combine the results and return a result
					if _lhs_boolean_.True() && _rhs_boolean_.True() {
						return BoolNode{
							Value: 1,
						}, ExceptionNode{Type: NO_EXCEPTION}
					}

					return BoolNode{
						Value: 0,
					}, ExceptionNode{Type: NO_EXCEPTION}
				}
			case OR_COMPARATOR:
				{
					// combine the results and return a result
					if _lhs_boolean_.True() || _rhs_boolean_.True() {
						return BoolNode{
							Value: 1,
						}, ExceptionNode{Type: NO_EXCEPTION}
					}

					return BoolNode{
						Value: 0,
					}, ExceptionNode{Type: NO_EXCEPTION}
				}
			}

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

			return eval.executeFunctionCode(_function_decl_.Code)
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

					return eval.executeFunctionCode(_function_decl_.Code)
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

					return eval.executeFunctionCode(_function_decl_.Code)
				}
			}

			return nil, ExceptionNode{Type: NO_EXCEPTION}
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
							_return := _implemented.Get(_index_.Value.Int64())

							if _exception, ok := _return.(ExceptionNode); ok {
								return nil, _exception
							}

							return _return, ExceptionNode{Type: NO_EXCEPTION}
						}
					case RANGE:
						{
							_end_index_, exception := eval.walkTree(_node.EndIndex)

							if exception.Type != NO_EXCEPTION {
								return nil, exception
							}

							if _eIndex_, ok := _end_index_.(NumberNode); ok {
								_return := _implemented.Range(_index_.Value.Int64(), _eIndex_.Value.Int64())

								if _exception, ok := _return.(ExceptionNode); ok {
									return nil, _exception
								}

								return _return, ExceptionNode{Type: NO_EXCEPTION}
							}
						}
					}
				}

				// fmt.Errorf("Failed to fetch element at the given index")
				return nil, ExceptionNode{
					Type:    INVALID_INDEX_EXCEPTION,
					Message: fmt.Sprintf("Failed to fetch element at the given index '%s'", Print(_index_)),
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
			return nil, eval.LoadModule(_node)
		}
	case ConditionNode:
		{
			// evaluate this stuff
			_lhs, exception := eval.walkTree(_node.Lhs)

			if exception.Type != NO_EXCEPTION {
				return nil, exception
			}

			// BUG
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

func (eval *Evaluator) Evaluate(initSymbolsTable bool) interface{} {
	var ret interface{}
	var exception ExceptionNode

	if initSymbolsTable {
		eval.symbolsTable.PushContext()
		defer eval.symbolsTable.PopContext()
	}

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
