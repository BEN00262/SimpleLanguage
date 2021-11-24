package typechecker

import (
	parser "github.com/BEN00262/simpleLang/parser"
	symTable "github.com/BEN00262/simpleLang/symbolstable"
)

type Types int

const (
	NULL Types = iota
	NUMBER
	PARAMETER
	STRING
	FUNCTION
	ARRAY
	BOOLEAN
	NIL
)

// [{ name: {Type: STRING }}]

// we need a signature for functions and anon functions or iife
type FunctionSignature struct {
	Args    []Types
	Returns Types
}

type TypeInference struct {
	Type      Types
	Signature FunctionSignature
}

type TypeChecker struct {
	program      *parser.ProgramNode
	symbolsTable *symTable.SymbolsTable
}

// get error report cleanly
func NewTypeChecker(program *parser.ProgramNode) *TypeChecker {
	return &TypeChecker{
		program:      program,
		symbolsTable: symTable.InitSymbolsTable(),
	}
}

// have a reference to the actual tokens and get the errors
func (typecheck *TypeChecker) _typecheck(node interface{}) Types {
	switch _node := node.(type) {
	case parser.NumberNode:
		{
			return NUMBER
		}
	case parser.StringNode:
		{
			return STRING
		}
	case parser.ExpressionNode:
		{
			return typecheck._typecheck(_node.Expression)
		}
	case parser.Import:
		{
			return NULL
		}
	case parser.ArrayNode:
		{
			return ARRAY
		}
	case parser.NilNode:
		{
			return NIL
		}
	case parser.ConditionNode:
		{
			// check the type
			return BOOLEAN
		}
	case parser.IFNode:
		{
			// type check the nodes and stuff
			conditionType := typecheck._typecheck(_node.Condition)

			if conditionType != BOOLEAN {
				panic("Requires a condition")
			}

		}
	case parser.ForNode:
		{

		}
	case parser.FunctionDecl:
		{
			typecheck.symbolsTable.PushContext()
			defer typecheck.symbolsTable.PopContext()

			// inject the parameters into the scope
			// for _, param := range _node.Params {
			// 	param.Key
			// }
			// inject the code in then walk the body get the types for the params
			// then return the typings

			// walk through the function and fetch values from the walking
			// insert the node and try to get the actual node type
			// generate the function signature
			// but whats the function signature
			// we insert this into

			// for _, arg := range _node.Params {

			// }

			// we have parameters ---> we can generate type signatures for the parameters by evaluating the function
			// since we will be having the parameters still in the scope after all the evaluation we can just get their types i think

		}
	case parser.Assignment:
		{
			rhsType := typecheck._typecheck(_node.Rvalue)

			if _node.Type == parser.REASSIGNMENT {
				value, err := typecheck.symbolsTable.GetFromContext(_node.Lvalue)

				if err != nil {
					panic(err.Error())
				}

				// get the type
				if _type, ok := (*value).(TypeInference); ok {
					if _type.Type != rhsType {
						panic("Invalid type on assignment")
					}

					goto insertType
				}

				panic("[UNREACHABLE] Failed to get value from symbols table")
			}

		insertType:
			typecheck.symbolsTable.PushToContext(_node.Lvalue, TypeInference{
				Type: rhsType,
			})
		}
	case parser.BinaryNode:
		{
			lhsType := typecheck._typecheck(_node.Lhs)
			rhsType := typecheck._typecheck(_node.Rhs)

			switch _node.Operator {
			case "+":
				{
					switch lhsType {
					case NUMBER:
						{
							switch rhsType {
							case NUMBER:
								return NUMBER
							case STRING:
								return STRING
							default:
								panic("Invalid types")
							}
						}
					case STRING:
						{
							if rhsType != NUMBER && rhsType != STRING {
								panic("Invalid types")
							}

							return STRING
						}
					}

					panic("Unsupported operation")
				}
			}
		}
	}

	// this should be unreacheable
	return NULL
}

func (typecheck *TypeChecker) Walk() {
	// context maintainance
	typecheck.symbolsTable.PushContext()
	defer typecheck.symbolsTable.PopContext()

	for _, node := range typecheck.program.Nodes {
		typecheck._typecheck(node)
	}
}
