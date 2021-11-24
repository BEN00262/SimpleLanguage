package ast

import (
	"fmt"
	"strings"

	parser "github.com/BEN00262/simpleLang/parser"
)

// walk down the AST and generate javascript code
// WASM

type AST struct {
	Node  *parser.ProgramNode
	depth int
}

func (ast *AST) ident(additional string) string {
	return strings.Repeat(" ", ast.depth) + additional
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
	case parser.ReturnNode:
		{
			return "return " + ast._walk(_node.Expression)
		}
	case parser.BinaryNode:
		{
			lhs := ast._walk(_node.Lhs)
			rhs := ast._walk(_node.Rhs)

			return fmt.Sprintf("%s %s %s", lhs, _node.Operator, rhs)
		}
	case parser.FunctionDecl:
		{
			ast.depth += 1
			defer func() { ast.depth -= 1 }()
			// function parameters
			func_parameters := []string{}

			for _, param := range _node.Params {
				func_parameters = append(func_parameters, param.Key)
			}

			// now generate the function signature
			func_body := ""

			for _, code := range _node.Code {
				func_body += ast.ident(ast._walk(code)) + "\n"
			}

			return ast.ident(fmt.Sprintf(
				"function %s (%s) {\n%s}\n", _node.Name, strings.Join(func_parameters, ","), func_body,
			))
		}
	case parser.Assignment:
		{
			rvalue := ast._walk(_node.Rvalue)

			switch _node.Type {
			case parser.CONST_ASSIGNMENT:
				return ast.ident(fmt.Sprintf("const %s = %s", _node.Lvalue, rvalue))
			case parser.ASSIGNMENT:
				return ast.ident(fmt.Sprintf("let %s = %s;", _node.Lvalue, rvalue))
			case parser.REASSIGNMENT:
				return ast.ident(fmt.Sprintf("%s = %s;", _node.Lvalue, rvalue))
			}
		}
	}

	return ""
}

// generate c for now and then build it to an executable
func (ast *AST) Walk() string {
	javascript := ""

	for _, node := range ast.Node.Nodes {
		javascript += ast.ident(ast._walk(node)) + "\n"
	}

	return javascript
}
