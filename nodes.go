package main

type ProgramNode struct {
	Nodes []interface{}
}

type VariableNode struct {
	Value string
}

type NumberNode struct {
	Value int
}
