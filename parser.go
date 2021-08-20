package main

type Parser struct {
	Tokens          []Token
	CurrentPosition int
}

func initParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:          tokens,
		CurrentPosition: 0,
	}
}

func (parser *Parser) CurrentToken() Token {
	return parser.Tokens[parser.CurrentPosition]
}

func (parser *Parser) Parse() *ProgramNode {
	program := &ProgramNode{}

	for ; parser.CurrentPosition < len(parser.Tokens); parser.CurrentPosition++ {
		token := parser.CurrentToken()

		switch token.Type {
		case VARIABLE:
			{
				program.Nodes = append(program.Nodes, VariableNode{
					Value: token.Value.(string),
				})
			}

		case NUMBER:
			{
				program.Nodes = append(program.Nodes, NumberNode{
					Value: token.Value.(int),
				})
			}
		}

		parser.CurrentPosition += 1
	}

	return program
}
