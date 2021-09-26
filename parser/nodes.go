package parser

import (
	"fmt"
	"strings"
)

type ExternalFunction = func(value ...*interface{}) interface{}

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

/*
	type Getter interface {
		Get(index int) interface{}
		Range(start int, end int) interface{}
	}
*/

type ExceptionNode struct {
	Type    string
	Message string
}

func (exception ExceptionNode) Get(index int) interface{} {
	// o for the type and 1 for the message
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

func (exception ExceptionNode) Range(start int, end int) interface{} {
	return NilNode{}
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

// implements a len thing
// and also can be looped i dont know how but i just know men

// every method that needs to support the length shit should do this
// implement the Get interface too for string manipulation

type Countable interface {
	Length() NumberNode
}

type Getter interface {
	Get(index int) interface{}
	Range(start int, end int) interface{}
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
	Array    string
	Index    interface{} // ExpressionNode
	EndIndex interface{} // also an expression
}

type ArrayNode struct {
	Elements []interface{}
}

func (array ArrayNode) Length() NumberNode {
	return NumberNode{
		Value: len(array.Elements),
	}
}

// this works for getting at a given index
func (array ArrayNode) Get(index int) interface{} {
	if (len(array.Elements) == 0) || (index > len(array.Elements)-1 || index < 0) {
		return NilNode{}
	}

	return array.Elements[index]
}

// get a range
func (array ArrayNode) Range(start int, end int) interface{} {
	// how tf do we implement this return a
	// check for the constraints
	if (start < 0 || start > array.Length().Value) || (start < 0 || start > array.Length().Value) {
		return NilNode{}
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
	_last_item_ := array.Get(len(array.Elements) - 1)
	array.Elements = array.Elements[:len(array.Elements)-1]
	return _last_item_
}

func (array *ArrayNode) InsertAt(index int, value interface{}) {
	if len(array.Elements) == 0 {
		return
	}

	if index >= 0 && index < len(array.Elements) {
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
	REASSIGNMENT
)

type Assignment struct {
	Type   AssignmentType
	Lvalue string // will change this to an interface (something that evaluates to something in the symbols table)
	Rvalue interface{}
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
	ForBody        []interface{}
}

type IFNode struct {
	Condition interface{}
	ThenBody  []interface{}
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

// we need a map of the args
type FunctionCall struct {
	Name     string
	ArgCount int
	Args     []interface{}
}

type FunctionDecl struct {
	Name       string
	ParamCount int
	Params     []Param
	Code       []interface{}
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
type NumberNode struct {
	Value int
}

// implementing arithmetic operations
func (number *NumberNode) Add(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			// just add the numbers
			return NumberNode{
				Value: number.Value + _right.Value,
			}
		}
	case StringNode:
		{
			// we convert the number to a string then add them together
			return StringNode{
				Value: fmt.Sprintf("%d %s", number.Value, _right.Value),
			}
		}
	}

	// we should return an error code here
	return NilNode{}
}

// implementing arithmetic operations
func (number *NumberNode) Sub(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			// just add the numbers
			return NumberNode{
				Value: number.Value - _right.Value,
			}
		}
	}

	// we should return an error code here
	return NilNode{}
}

// implementing arithmetic operations
func (number *NumberNode) Mod(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			// just add the numbers
			return NumberNode{
				Value: number.Value % _right.Value,
			}
		}
	}

	// we should return an error code here
	return NilNode{}
}

// implementing arithmetic operations
func (number *NumberNode) Div(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			// just add the numbers
			return NumberNode{
				Value: number.Value / _right.Value,
			}
		}
	}

	// we should return an error code here
	return NilNode{}
}

// implementing arithmetic operations
func (number *NumberNode) Mul(right interface{}) interface{} {
	// work with interfaces
	switch _right := right.(type) {
	case NumberNode:
		{
			// just add the numbers
			return NumberNode{
				Value: number.Value * _right.Value,
			}
		}
	}

	// we should return an error code here
	return NilNode{}
}

// comperison interface implementation
func (numberNode *NumberNode) IsEqualTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case NumberNode:
		{
			if numberNode.Value == _rvalue.Value {
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

func (numberNode *NumberNode) IsGreaterThan(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case NumberNode:
		{
			if numberNode.Value > _rvalue.Value {
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
func (numberNode *NumberNode) IsGreaterThanOrEqualsTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case NumberNode:
		{
			if numberNode.Value >= _rvalue.Value {
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

func (numberNode *NumberNode) IsLessThan(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case NumberNode:
		{
			if numberNode.Value < _rvalue.Value {
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
func (numberNode *NumberNode) IsLessThanOrEqualsTo(value interface{}) BoolNode {
	switch _rvalue := value.(type) {
	case NumberNode:
		{
			if numberNode.Value <= _rvalue.Value {
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

type StringNode struct {
	Value string
}

func (stringNode *StringNode) Add(right interface{}) interface{} {
	switch _right := right.(type) {
	case NumberNode:
		{
			return StringNode{
				Value: fmt.Sprintf("%s %d", stringNode.Value, _right.Value),
			}
		}
	case StringNode:
		{
			return StringNode{
				Value: fmt.Sprintf("%s%s", stringNode.Value, _right.Value),
			}
		}
	}

	return NilNode{}
}

func (stringNode *StringNode) Mul(right interface{}) interface{} {
	switch _right := right.(type) {
	case NumberNode:
		{
			return StringNode{
				Value: strings.Repeat(stringNode.Value, _right.Value),
			}
		}
	}

	return NilNode{}
}

// not implemented for the language
func (stringNode *StringNode) Mod(right interface{}) interface{} {
	return NilNode{}
}

func (stringNode *StringNode) Div(right interface{}) interface{} {
	return NilNode{}
}

func (stringNode *StringNode) Sub(right interface{}) interface{} {
	return NilNode{}
}

// make the string indexeable and countable
func (stringNode StringNode) Length() NumberNode {
	return NumberNode{
		Value: len(stringNode.Value),
	}
}

// make it indexeable
func (stringNode StringNode) Get(index int) interface{} {
	if index < 0 || index > stringNode.Length().Value-1 {
		return NilNode{}
	}

	return StringNode{
		Value: string(stringNode.Value[index]),
	}
}

// get a range
func (stringNode StringNode) Range(start int, end int) interface{} {
	// how tf do we implement this return a
	// check for the constraints
	if (start < 0 || start > stringNode.Length().Value) || (start < 0 || start > stringNode.Length().Value) {
		return NilNode{}
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
			if strings.Compare(stringNode.Value, _rvalue.Value) > 1 {
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
			if strings.Compare(stringNode.Value, _rvalue.Value) == 0 || strings.Compare(stringNode.Value, _rvalue.Value) > 1 {
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
