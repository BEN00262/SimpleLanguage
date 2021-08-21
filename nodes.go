package main

type ProgramNode struct {
	Nodes []interface{}
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
type ComparisonNode struct {
}

type NilNode struct{}

// use an interger
type BoolNode struct {
	Value int
}

type ExternalFunctionNode struct {
	Name       string
	ParamCount int
	Function   ExternalFunction
}

// we need a map of the args
type FunctionCall struct {
	Name     string
	ArgCount int
	Args     map[string]interface{}
}

type FunctionDecl struct {
	Name       string
	ParamCount int
	Code       []interface{}
}

type CommentNode struct {
	comment string
}

type VariableNode struct {
	Value string
}

type NumberNode struct {
	Value int
}
