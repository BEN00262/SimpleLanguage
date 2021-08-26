package main

import (
	"strings"
)

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

type ExpressionNode struct {
	// this can be anything
	expression interface{}
}

type BinaryNode struct {
	Lhs      interface{}
	Operator string
	Rhs      interface{}
}

type Assignment struct {
	Lvalue string
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

// this implements the Equals interface
type NumberNode struct {
	Value int
}

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
