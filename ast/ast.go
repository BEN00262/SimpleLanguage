package ast

import (
	"fmt"

	parser "github.com/BEN00262/simpleLang/parser"
)

// walk down the AST and generate javascript code
// WASM

type AST struct {
	Node *parser.ProgramNode
}

func (ast *AST) _walk(node interface{}) string {
	switch _node := node.(type) {
	case parser.NumberNode:
		{
			return parser.Print(_node)
		}
	case parser.VariableNode:
		{
			return _node.Value
		}
	case parser.ExpressionNode:
		{
			return ast._walk(_node.Expression)
		}
	case parser.Assignment:
		{
			rvalue := ast._walk(_node.Rvalue)

			switch _node.Type {
			case parser.CONST_ASSIGNMENT:
				return fmt.Sprintf("const %s = %s", _node.Lvalue, rvalue)
			case parser.ASSIGNMENT:
				return fmt.Sprintf("let %s = %s;", _node.Lvalue, rvalue)
			case parser.REASSIGNMENT:
				return fmt.Sprintf("%s = %s;", _node.Lvalue, rvalue)
			}
		}
	}

	return ""
}

func (ast *AST) Walk() string {
	javascript := ""

	for _, node := range ast.Node.Nodes {
		javascript += ast._walk(node) + "\n"
	}

	return javascript
}
