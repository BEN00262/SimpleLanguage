package ast

import (
	"fmt"
	"strings"

	. "github.com/BEN00262/simpleLang/parser"
)

// this beauty prints the AST
// get the ast walk it then print the shit out of it
type AST struct {
	Depth   int
	Program *ProgramNode
	String  string
}

// we start a depth of 0
func initAST(program *ProgramNode) *AST {
	return &AST{
		Depth:   0,
		Program: program,
	}
}

// implement an interface to interact with the AST and then expose the shit to the language itself
// so that one can modify stuff

// increment the depth of the ast
func (ast *AST) IncrementDepth() {
	ast.Depth += 1
}

// decrement the depth of the ast
func (ast *AST) DecrementDepth() {
	ast.Depth -= 1
}

// a print at the given spacing :)
func (ast *AST) DisplayAtCurrentSpacing(whatToPrint string) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", ast.Depth), whatToPrint)
}

func (ast *AST) AppendToFinalString(content string) {
	ast.String += content
}

func (ast *AST) _walk(child interface{}) string {
	// start the walking of the AST and printing it
	switch _node := child.(type) {
	case VariableNode:
		{
			// just print it at the required depth
			return ast.DisplayAtCurrentSpacing("<Variable Node>")
		}
	case NumberNode:
		{
			// just print it also at the required depth
			return ast.DisplayAtCurrentSpacing("<Number Node>")
		}

	case Assignment:
		{
			// display the left and the right
			ast.AppendToFinalString(
				ast.DisplayAtCurrentSpacing(
					fmt.Sprintf("%s = %s", _node.Lvalue, ast._walk(_node.Rvalue)),
				),
			)
		}
	case BinaryNode:
		{
			return ast.DisplayAtCurrentSpacing(
				fmt.Sprintf("%s %s %s", ast._walk(_node.Lhs), _node.Operator, ast._walk(_node.Rhs)),
			)
		}
	case ExpressionNode:
		{
			ast.AppendToFinalString(
				ast.DisplayAtCurrentSpacing("Start of an expression\n"),
			)

			// add the stuff in
			// ast.AppendToFinalString(
			// 	ast._walk(_node) + "\n",
			// )

			ast.AppendToFinalString(
				ast.DisplayAtCurrentSpacing("end of an expression\n"),
			)
		}
	case FunctionDecl:
		{
			ast.AppendToFinalString(
				ast.DisplayAtCurrentSpacing("<Start of function>\n"),
			)

			ast.IncrementDepth()

			// we can walk the code and print shit here by the way
			for _, _code := range _node.Code {
				ast.AppendToFinalString(
					ast._walk(_code) + "\n",
				)
			}

			ast.DecrementDepth()

			ast.AppendToFinalString(
				ast.DisplayAtCurrentSpacing("<End of function>\n"),
			)
		}
	}

	return ast.String
}

func (ast *AST) walk() {
	// start the walking of the AST and printing it
	for _, child := range ast.Program.Nodes {
		ast._walk(child)
	}

	fmt.Println(ast.String)
}
