package parser

import (
	"fmt"
	"math/big"
	"strings"

	. "github.com/BEN00262/simpleLang/exceptions"
)

type ExternalFunction = func(value ...*interface{}) (interface{}, ExceptionNode)

type ProgramNode struct {
	Nodes []interface{}
}

type Comparison interface {
	IsEqualTo(value interface{}) BoolNode
	IsGreaterThan(value interface{}) BoolNode
	IsGreaterThanOrEqualsTo(value interface{}) BoolNode
	IsLessThanOrEqualsTo(value interface{}) BoolNode
	IsLessThan(value interface{}) BoolNode
}

// just dump the shiets in the current scope
type Import struct {
	FileName string
	Alias    string
}

// module thing
// how the hell will the functions in a module refer
// a module should have its own context ()

// shelve this for later
// register assignments, function decls --> i think
type ModuleValue struct {
	// type
	// interface
}

type Module struct {
	Name                 string
	MethodsAndProperties map[string]interface{}
}

// returns an interface of the stored value --> evaluate the code then on reaching
//
func (module *Module) Get(key string) interface{} {
	return nil
}

type ExpressionNode struct {
	// this can be anything
	Expression interface{}
}

type LogicalComparisonType int

const (
	AND_COMPARATOR LogicalComparisonType = iota + 1
	OR_COMPARATOR
)

type LogicalComparison struct {
	Type LogicalComparisonType
	Rhs  interface{} // Expression
	Lhs  interface{} // Expression
}

/*
	type Getter interface {
		Get(index int64) interface{}
		Range(start int64, end int64) interface{}
	}
*/

type ExceptionNode struct {
	Type    string
	Message string
}

func (exception ExceptionNode) Get(index int64) interface{} {
	if index == 0 {
		return StringNode{
			Value: exception.Type,
		}
	} else if index == 1 {
		return StringNode{
			Value: exception.Message,
		}
	}

	// actually return errors here
	return NilNode{}
}

func (exception ExceptionNode) Range(start int64, end int64) interface{} {
	return ExceptionNode{
		Type: INVALID_OPERATION_EXCEPTION,
	}
}

// generate jump the state to somewhere
// convert the evaluator to state machine
// on getting something react to it immediately
type RaiseExceptionNode struct {
	Exception interface{} // an expression node for now
}

type CatchBlock struct {
	// this is the exception var we actually need it
	Exception string
	Body      []interface{}
}

// a try catch node
// we need a way to inject struff into the catch block
type TryCatchNode struct {
	Try     []interface{}
	Catch   CatchBlock
	Finally []interface{}
}

type IExportables interface {
	IsExported() bool
}

// export visibilty node
type ExportVisibilityNode struct {
	// when wrapped with this u are exported my nigga
	// this only works for functions declarations and constants --> check in the parser for adherence
	Exported interface{}
}

// create a simple . properties thing
// module.first.sample(78)
type ObjectAccessor struct {
	Parent string
	Child  string // improve it later to support further lookups into the structure
}

// implements a len thing
// and also can be looped i dont know how but i just know men

// every method that needs to support the length shit should do this
// implement the Get interface too for string manipulation

type Countable interface {
	Length() NumberNode
}

type Getter interface {
	Get(index int64) interface{}
	Range(start int64, end int64) interface{}
}

type AccessorType = int

const (
	NORMAL AccessorType = iota + 1
	RANGE
)

// create an array accessor node
// it has two things
type ArrayAccessorNode struct {
	// the array name --> should be an expression that resolves to an array stuff else throw an error
	// the index into the array --> the expression should evaluate to a number node
	Type     AccessorType
	Array    ObjectAccessor
	Index    interface{} // ExpressionNode
	EndIndex interface{} // also an expression
}

type ArrayNode struct {
	Elements []interface{}
}

func (array ArrayNode) Length() NumberNode {
	return NumberNode{
		Value: *big.NewInt(int64(len(array.Elements))),
	}
}

// this works for getting at a given index
func (array ArrayNode) Get(index int64) interface{} {
	if (len(array.Elements) == 0) || (index > int64(len(array.Elements)-1) || index < 0) {
		return ExceptionNode{
			Type:    INVALID_INDEX_EXCEPTION,
			Message: "Invalid index access",
		}
	}

	return array.Elements[index]
}

// get a range
func (array ArrayNode) Range(start int64, end int64) interface{} {
	arrayLength := array.Length().Value
	_startCmp := arrayLength.Cmp(big.NewInt(start))
	_endCmp := arrayLength.Cmp(big.NewInt(end))

	// check for the constraints
	if start < 0 || _startCmp == 1 || end < 0 || _endCmp == 1 {
		return ExceptionNode{
			Type:    INVALID_INDEX_EXCEPTION,
			Message: "Invalid range access limits; Ensure the limits are within the array",
		}
	}

	return ArrayNode{
		Elements: array.Elements[start:end],
	}
}

// push(array, value)
// pop(array)

func (array *ArrayNode) Push(value interface{}) {
	array.Elements = append(array.Elements, value)
}

// poping a value return it
func (array *ArrayNode) Pop() interface{} {
	_last_item_ := array.Get(int64(len(array.Elements) - 1))
	array.Elements = array.Elements[:len(array.Elements)-1]
	return _last_item_
}

func (array *ArrayNode) InsertAt(index int64, value interface{}) {
	if len(array.Elements) == 0 {
		return
	}

	if index >= 0 && index < int64(len(array.Elements)) {
		array.Elements[index] = value
	}
}

type BinaryNode struct {
	Lhs      interface{}
	Operator string
	Rhs      interface{}
}

// we need the type of the Assignment is it a reassignment or a fresh assignment :)
type AssignmentType = int

const (
	ASSIGNMENT AssignmentType = iota + 1
	CONST_ASSIGNMENT
	REASSIGNMENT
)

type Assignment struct {
	Type   AssignmentType
	Lvalue string // will change this to an interface (something that evaluates to something in the symbols table)
	Rvalue interface{}
}

func (assignment Assignment) IsExported() bool {
	if assignment.Type != REASSIGNMENT {
		return true
	}
	return false
}

// expression ( which returns a True or False )
type ConditionNode struct {
	Lhs      interface{}
	Operator string
	Rhs      interface{}
}

type AnonymousFunction struct {
	ParamCount int
	Params     []Param
	Code       []interface{}
}

type IIFENode struct {
	Function AnonymousFunction
	Args     []interface{}
	ArgCount int
}

type BreakNode struct {
}

// we have different type of this
type ForNodeType = int

const (
	WHILE_FOREVER ForNodeType = iota + 1
	WHILE_CONDITIONAL
	FOR_NODE
)

type ForNode struct {
	Type           ForNodeType
	Initialization interface{}
	Condition      interface{}
	Increment      interface{}
	ForBody        BlockNode
}

// we have a block node hapa buana
type IFNode struct {
	Condition interface{}
	ThenBody  BlockNode
	ElseBody  []interface{}
}

type BlockNode struct {
	Code []interface{}
}

type NilNode struct{}

// implement for the nil node here
/*
IsEqualTo(value interface{}) BoolNode
IsGreaterThan(value interface{}) BoolNode
IsGreaterThanOrEqualsTo(value interface{}) BoolNode
IsLessThanOrEqualsTo(value interface{}) BoolNode
IsLessThan(value interface{}) BoolNode
*/

func (null *NilNode) IsEqualTo(value interface{}) BoolNode {
	switch value.(type) {
	case NilNode:
		{
			// print something here

			return BoolNode{
				Value: 1,
			}
		}
	}

	// throw an error here
	return BoolNode{
		Value: 0,
	}
}

func (null *NilNode) IsGreaterThan(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

func (null *NilNode) IsGreaterThanOrEqualsTo(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

func (null *NilNode) IsLessThanOrEqualsTo(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

// IsLessThan
func (null *NilNode) IsLessThan(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

// use an interger
type BoolNode struct {
	Value int
}

// this has something to do with true of false
func (boolNode *BoolNode) True() bool {
	return boolNode.Value == 1
}

func (boolNode *BoolNode) IsEqualTo(value interface{}) BoolNode {
	switch _lhs := value.(type) {
	case BoolNode:
		{
			// print something here
			if boolNode.Value == _lhs.Value {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	// throw an error here
	return BoolNode{
		Value: 0,
	}
}

func (boolNode *BoolNode) IsGreaterThan(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

func (boolNode *BoolNode) IsGreaterThanOrEqualsTo(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

func (boolNode *BoolNode) IsLessThan(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

func (boolNode *BoolNode) IsLessThanOrEqualsTo(value interface{}) BoolNode {
	return BoolNode{
		Value: 0,
	}
}

// implement one of the interfaces and throw

type ExternalFunctionNode struct {
	Name       string
	ParamCount int
	Function   ExternalFunction
}

// the name should not be a string it should be an ObjectAccessor that resolves to a string
// handle accessors by all means or throw errors
type FunctionCall struct {
	Name     ObjectAccessor
	ArgCount int
	Args     []interface{}
}

// support variadic functions ... what happens is we create an array of the rest and pass them as such
type FunctionDecl struct {
	Name       string
	ParamCount int
	Params     []Param
	Code       []interface{}
}

func (function FunctionDecl) IsExported() bool {
	return true
}

type CommentNode struct {
	comment string
}

type ReturnNode struct {
	Expression interface{}
}

type VariableNode struct {
	Value string
}

// Arithmetic operations
type ArthOp interface {
	Add(right interface{}) interface{}
	Sub(right interface{}) interface{}
	Mod(right interface{}) interface{}
	Div(right interface{}) interface{}
	Mul(right interface{}) interface{}
}

// this implements the Equals interface
// this should not be an int ( use float64 laters )
// this should be a big integer thing
// we should place also a float value hapa
type NumberType int

const (
	INTEGER = iota + 1
	FLOAT
)

type NumberNode struct {
	Type   NumberType
	Value  big.Int
	FValue big.Float
}

// implementing arithmetic operations
func (number *NumberNode) Add(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Add(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Add(
									left_float_number_copy,
									right_float_number_copy,
								),
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							// left_float_number_copy.

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Add(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							return NumberNode{
								Type: INTEGER,
								Value: *left_float_number_copy.Add(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					}
				}
			}
		}
	case StringNode:
		{
			return StringNode{
				Value: fmt.Sprintf("%s%s", Print(number), _right.Value),
			}
		}
	}

	// we should return an error code here
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Unsupported operation on type 'number'",
	}
}

// implementing arithmetic operations
func (number *NumberNode) Sub(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Sub(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Sub(
									left_float_number_copy,
									right_float_number_copy,
								),
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Sub(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							return NumberNode{
								Type: INTEGER,
								Value: *left_float_number_copy.Sub(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					}
				}
			}
		}
	}

	// we should return an error code here
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Unsupported operation on type 'number'",
	}
}

// implementing arithmetic operations
func (number *NumberNode) Mod(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT, INTEGER:
						goto error_exit
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							goto error_exit
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							if right_float_number_copy.Cmp(big.NewInt(0)) == 0 {
								return ExceptionNode{
									Type:    DIVISION_BY_ZERO_EXCEPTION,
									Message: "Division or Modulo by zero error",
								}
							}

							return NumberNode{
								Type: INTEGER,
								Value: *left_float_number_copy.Mod(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					}
				}
			}
		}
	}

	// we should return an error code here
error_exit:
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Unsupported operation on type 'number'",
	}
}

// implementing arithmetic operations
func (number *NumberNode) Div(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT, INTEGER:
						{
							goto error_exit
						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							goto error_exit
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							if right_float_number_copy.Cmp(big.NewInt(0)) == 0 {
								return ExceptionNode{
									Type:    DIVISION_BY_ZERO_EXCEPTION,
									Message: "Division or Modulo by zero error",
								}
							}

							return NumberNode{
								Type: INTEGER,
								Value: *left_float_number_copy.Div(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					}
				}
			}
		}
	}

	// we should return an error code here
error_exit:
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Unsupported operation on type 'number'",
	}
}

// implementing arithmetic operations
func (number *NumberNode) Mul(right interface{}) interface{} {
	switch _right := right.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Mul(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Mul(
									left_float_number_copy,
									right_float_number_copy,
								),
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							return NumberNode{
								Type: FLOAT,
								FValue: *left_float_number_copy.Mul(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							return NumberNode{
								Type: INTEGER,
								Value: *left_float_number_copy.Mul(
									left_float_number_copy,
									right_float_number_copy,
								),
							}
						}
					}
				}
			}
		}
	}

	// we should return an error code here
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Unsupported operation on type 'number'",
	}
}

// comperison interface implementation
func (number *NumberNode) IsEqualTo(value interface{}) BoolNode {
	switch _right := value.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == 0 {
								return BoolNode{
									Value: 1,
								}
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					}
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

func (number *NumberNode) IsGreaterThan(value interface{}) BoolNode {
	switch _right := value.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == 1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == 1 {
								return BoolNode{
									Value: 1,
								}
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == 1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == 1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					}
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

// IsGreaterThanOrEqualsTo
func (number *NumberNode) IsGreaterThanOrEqualsTo(value interface{}) BoolNode {
	switch _right := value.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == 1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == 1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == 1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == 1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					}
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

func (number *NumberNode) IsLessThan(value interface{}) BoolNode {
	switch _right := value.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == -1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == -1 {
								return BoolNode{
									Value: 1,
								}
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							if left_float_number_copy.Cmp(right_float_number_copy) == -1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							if left_float_number_copy.Cmp(right_float_number_copy) == -1 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					}
				}
			}

		}
	}

	return BoolNode{
		Value: 0,
	}
}

// IsLessThanOrEqualsTo
func (number *NumberNode) IsLessThanOrEqualsTo(value interface{}) BoolNode {
	switch _right := value.(type) {
	case NumberNode:
		{
			switch number.Type {
			case FLOAT:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == -1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{

							left_float_number_copy := new(big.Float).Set(&number.FValue)
							right_float_number_copy := new(big.Float).SetInt(&_right.Value)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == -1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}

						}
					}
				}
			case INTEGER:
				{
					switch _right.Type {
					case FLOAT:
						{
							left_float_number_copy := new(big.Float).SetInt(&number.Value)
							right_float_number_copy := new(big.Float).Set(&_right.FValue)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == -1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					case INTEGER:
						{
							left_float_number_copy := new(big.Int).Set(&number.Value)
							right_float_number_copy := new(big.Int).Set(&_right.Value)

							_comp := left_float_number_copy.Cmp(right_float_number_copy)
							if _comp == -1 || _comp == 0 {
								return BoolNode{
									Value: 1,
								}
							}
						}
					}
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

type StringNode struct {
	Value string
}

func (stringNode *StringNode) Add(right interface{}) interface{} {
	switch _right := right.(type) {
	case NumberNode, StringNode:
		return StringNode{
			Value: fmt.Sprintf("%s%s", stringNode.Value, Print(_right)),
		}
	}

	// we have an an invalid operation
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Operation not supported in strings",
	}
}

func (stringNode *StringNode) Mul(right interface{}) interface{} {

	// ensure the number has type integer if not throw an error buana

	switch _right := right.(type) {
	case NumberNode:
		{
			// ensure the number is integer
			if _right.Type == FLOAT {
				goto error_exit
			}

			return StringNode{
				Value: strings.Repeat(stringNode.Value, int(_right.Value.Int64())),
			}
		}
	}

error_exit:
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Operation not supported in strings",
	}
}

// not implemented for the language
func (stringNode *StringNode) Mod(right interface{}) interface{} {
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Operation not supported in strings",
	}
}

func (stringNode *StringNode) Div(right interface{}) interface{} {
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Operation not supported in strings",
	}
}

func (stringNode *StringNode) Sub(right interface{}) interface{} {
	return ExceptionNode{
		Type:    INVALID_OPERATION_EXCEPTION,
		Message: "Operation not supported in strings",
	}
}

// make the string indexeable and countable
func (stringNode StringNode) Length() NumberNode {
	// the number node is a big integer stuff :)
	return NumberNode{
		Type:  INTEGER,
		Value: *big.NewInt(int64(len(stringNode.Value))),
	}
}

// make it indexeable
func (stringNode StringNode) Get(index int64) interface{} {
	result := stringNode.Length().Value
	result = *result.Sub(&result, big.NewInt(1))

	if index < 0 || result.Cmp(big.NewInt(index)) == -1 {
		return ExceptionNode{
			Type:    INVALID_INDEX_EXCEPTION,
			Message: "Invalid index access",
		}
	}

	return StringNode{
		Value: string(stringNode.Value[index]),
	}
}

// get a range
func (stringNode StringNode) Range(start int64, end int64) interface{} {
	arrayLength := stringNode.Length().Value
	_startCmp := arrayLength.Cmp(big.NewInt(start))
	_endCmp := arrayLength.Cmp(big.NewInt(end))

	// check for the constraints
	if start < 0 || _startCmp == 1 || end < 0 || _endCmp == 1 {
		return ExceptionNode{
			Type:    INVALID_INDEX_EXCEPTION,
			Message: "Invalid range access limits; Ensure the limits are within the array",
		}
	}

	return StringNode{
		Value: stringNode.Value[start:end],
	}
}

func (stringNode *StringNode) IsEqualTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case StringNode:
		{
			if strings.Compare(stringNode.Value, _rvalue.Value) == 0 {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

func (stringNode *StringNode) IsGreaterThan(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case StringNode:
		{
			if stringNode.Value > _rvalue.Value {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

// IsGreaterThanOrEqualsTo
func (stringNode *StringNode) IsGreaterThanOrEqualsTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case StringNode:
		{
			if stringNode.Value == _rvalue.Value || stringNode.Value > _rvalue.Value {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

// IsLessThanOrEqualsTo
func (stringNode *StringNode) IsLessThan(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case StringNode:
		{
			if stringNode.Value < _rvalue.Value {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}

// IsLessThanOrEqualsTo
func (stringNode *StringNode) IsLessThanOrEqualsTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case StringNode:
		{
			if stringNode.Value <= _rvalue.Value {
				return BoolNode{
					Value: 1,
				}
			}
		}
	}

	return BoolNode{
		Value: 0,
	}
}
