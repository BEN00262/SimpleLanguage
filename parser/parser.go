package parser

import (
	"fmt"
	"math/big"
	"strings"

	. "github.com/BEN00262/simpleLang/lexer"
)

type ParsingState = int

const (
	FUNCTION_STATE ParsingState = iota
	LOOP_STATE
	TRY_CATCH_STATE
	INVALID_STATE // used when we are not in a function or a loop
)

// how the fuck will we do this
// we should get a reference to the code for error reporting and stuff
type Parser struct {
	Tokens          []Token
	ActualCode      []string
	State           []ParsingState // push into this state
	CurrentPosition int
	TokensLength    int
	ParsingError    string // this is showing the error and crashing this should a chan
}

func InitParser(tokens []Token, actualCode []string) *Parser {
	return &Parser{
		Tokens:          tokens,
		CurrentPosition: 0,
		State:           []ParsingState{INVALID_STATE},
		TokensLength:    len(tokens),
		ActualCode:      actualCode,
	}
}

func (parser *Parser) pushToParsingState(state ParsingState) {
	parser.State = append(parser.State, state)
}

func (parser *Parser) popFromParsingState() {
	parser.State = parser.State[0 : len(parser.State)-1]
}

// we have to think we cannot have a break in a function alone
// so in break the top state must be loo_state
// for return the parent state must be the function state
func (parser *Parser) EnsureCurrentParsingStateIs(state ParsingState, dig bool) bool {
	if dig {
		// start from the end going to the front trying to check if we are in a function or what
		// we need to make sure that the CurrentCounter is greater than 1 else just return the one
		for CurrentStateCounter := len(parser.State) - 1; CurrentStateCounter > -1; CurrentStateCounter-- {
			if parser.State[CurrentStateCounter] == state {
				return true
			}
		}

		return false
	}

	return parser.State[len(parser.State)-1] == state
}

func (parser *Parser) CurrentToken() Token {
	if parser.CurrentPosition < parser.TokensLength {
		return parser.Tokens[parser.CurrentPosition]
	}

	return Token{}
}

func (parser *Parser) peekAhead() Token {
	futurePosition := parser.CurrentPosition + 1

	if futurePosition < parser.TokensLength {
		return parser.Tokens[futurePosition]
	}

	return Token{}
}

func (parser *Parser) eatToken() {
	parser.CurrentPosition++
}

func (parser *Parser) _parseFactorTail() interface{} {
	_currentToken := parser.CurrentToken()

	// assume that % has the same affinity as these operators
	if _currentToken.Type == OPERATOR && strings.Contains("*/%", _currentToken.Value.(string)) {
		parser.eatToken() // eat the operator * or /

		factor := parser._parseFactor()
		factorTail := parser._parseFactorTail()

		if factorTail != nil {
			_factorTail := factorTail.(BinaryNode)
			_factorTail.Lhs = factor
			factor = _factorTail
		}

		return BinaryNode{
			Rhs:      factor,
			Operator: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseFactor() interface{} {
	_currentToken := parser.CurrentToken()

	if IsTypeAndValue(_currentToken, HALF_CIRCLE_BRACKET, "(") {
		parser.eatToken() // eat the (

		if IsTypeAndValue(parser.CurrentToken(), KEYWORD, FUNC) {
			_iife_func_ := parser.ParseFunction()

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, ")",
				"Expected a ')'",
			)

			// we check if there is anything like a ( if there is parse it as an iife else parse it
			// as an anonymous function
			if !IsTypeAndValue(parser.CurrentToken(), HALF_CIRCLE_BRACKET, "(") {
				// return the func declaration as an anonymous thing
				return _iife_func_
			}

			// just incase its an iife
			_iife_function_ := _iife_func_.(AnonymousFunction)

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, "(",
				"Expected a '('",
			)

			_args := parser._parseFunctionArgs()

			parser.IsExpectedEatElsePanic(
				parser.CurrentToken(),
				HALF_CIRCLE_BRACKET, ")",
				"Expected a ')'",
			)

			return IIFENode{
				Function: _iife_function_,
				Args:     _args,
				ArgCount: len(_args),
			}
		}

		_expression := parser._parseExpression()

		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			HALF_CIRCLE_BRACKET, ")",
			"Expected a ')'",
		)

		return _expression
	} else if _currentToken.Type == NUMBER {
		parser.eatToken()

		// tuko hapa buana
		// fmt.Println(_currentToken)
		_raw_number_ := _currentToken.Value.(string)

		// check if there are any . in the number if there are we have a float
		// improve this later
		if strings.Contains(_raw_number_, ".") {
			_number, isSuccess := new(big.Float).SetString(_raw_number_)

			if !isSuccess {
				parser.reportError(_currentToken)
			}

			return NumberNode{
				Type:   FLOAT,
				FValue: *_number,
			}

		}

		_number, isSuccess := new(big.Int).SetString(_raw_number_, 0)

		if !isSuccess {
			parser.reportError(_currentToken)
		}

		return NumberNode{
			Type:  INTEGER,
			Value: *_number,
		}
	} else if IsTypeAndValue(_currentToken, SQUARE_BRACKET, "[") {
		// start parsing the array here
		// consume the first square bracke
		parser.eatToken() // eat the [

		var _elements_ []interface{}

		for parser.CurrentPosition < parser.TokensLength && !IsTypeAndValue(parser.CurrentToken(), SQUARE_BRACKET, "]") {
			_elements_ = append(_elements_, parser._parseExpression())

			if parser.CurrentToken().Type == COMMA {
				parser.eatToken()
			}
		}

		// final we expect a closing bracket
		// eat ]
		parser.IsExpectedEatElsePanic(
			parser.CurrentToken(),
			SQUARE_BRACKET, "]",
			"Expected ']'",
		)

		return ArrayNode{
			Elements: _elements_,
		}
	} else if _currentToken.Type == STRING {
		parser.eatToken()

		return StringNode{
			Value: _currentToken.Value.(string),
		}
	} else if _currentToken.Type == KEYWORD {

		if strings.Compare(_currentToken.Value.(string), TRUE) == 0 {
			parser.eatToken()
			return BoolNode{
				Value: 1,
			}
		}

		if strings.Compare(_currentToken.Value.(string), FALSE) == 0 {
			parser.eatToken()
			return BoolNode{
				Value: 0,
			}
		}

		// check if we get a fun keyword if so ensure that the next value is not a name and return
		if _currentToken.Value.(string) == FUNC {
			// peer ahead to see
			// report the error well --> do it kesho
			if !IsTypeAndValue(parser.peekAhead(), HALF_CIRCLE_BRACKET, "(") {
				parser.reportError(parser.peekAhead(), "Function expression expects no name")
			}

			return parser.ParseFunction()

		}

		if _currentToken.Value == NIL {
			parser.eatToken()
			return NilNode{}
		}

	} else if _currentToken.Type == VARIABLE {
		return parser.parseVariableExpression()
	}

	return nil
}

func (parser *Parser) _parseTermTail() interface{} {
	_currentToken := parser.CurrentToken()

	// also check for conditionals here
	if _currentToken.Type == OPERATOR && strings.Contains("+-", _currentToken.Value.(string)) {
		parser.eatToken() // eat the operator e.g + or -

		term := parser._parseTerm()
		termTail := parser._parseTermTail()

		if termTail != nil {
			_termTail := termTail.(BinaryNode)
			_termTail.Lhs = term
			term = _termTail
		}

		return BinaryNode{
			Rhs:      term,
			Operator: _currentToken.Value.(string),
		}
	}

	return nil
}

func (parser *Parser) _parseTerm() interface{} {
	factor := parser._parseFactor()
	factorTail := parser._parseFactorTail()

	if factorTail != nil {
		_factorTail := factorTail.(BinaryNode)
		_factorTail.Lhs = factor
		factor = _factorTail
	}

	return factor
}

func (parser *Parser) _parseInnerComparison() interface{} {
	term := parser._parseTerm()
	termTail := parser._parseTermTail()

	if termTail != nil {
		_termTail := termTail.(BinaryNode)
		_termTail.Lhs = term
		term = _termTail
	}

	return ExpressionNode{
		Expression: term,
	}
}

func (parser *Parser) _parseComparison() interface{} {
	_lhs := parser._parseInnerComparison()
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == CONDITION && Contains([]string{">", ">=", "<", "<="}, _currentToken.Value.(string)) {
		parser.eatToken()
		_rhs := parser._parseInnerComparison()

		return ConditionNode{
			Lhs:      _lhs,
			Rhs:      _rhs,
			Operator: _currentToken.Value.(string),
		}
	}

	return _lhs
}

func (parser *Parser) _parseEqualsNotEquals() interface{} {
	_lhs := parser._parseComparison()

	_currentToken := parser.CurrentToken()

	// this fuckers should eat their shit
	if _currentToken.Type == CONDITION && Contains([]string{"!=", "=="}, _currentToken.Value.(string)) {
		parser.eatToken()

		_expression := parser._parseComparison()

		return ConditionNode{
			Operator: _currentToken.Value.(string),
			Rhs:      _expression,
			Lhs:      _lhs,
		}
	}

	return _lhs
}

// we need to parse 'and' and 'or' too buana
func (parser *Parser) _parseExpression() interface{} {
	_lhs := parser._parseEqualsNotEquals()
	_currentToken := parser.CurrentToken()

	if _currentToken.Type == KEYWORD && (_currentToken.Value.(string) == OR || _currentToken.Value.(string) == AND) {
		// we eat the token
		parser.eatToken()
		_rhs := parser._parseExpression()

		_type := AND_COMPARATOR

		if _currentToken.Value.(string) == OR {
			_type = OR_COMPARATOR
		}

		return LogicalComparison{
			Type: _type,
			Lhs:  _lhs,
			Rhs:  _rhs,
		}
	}

	return _lhs
}

// use a stack to monitor the current state we are --> are we parsing a function or a loop
// push the current parsing state to a stack then check on special tokens like return and break
func (parser *Parser) _parse(token Token) interface{} {
	switch token.Type {
	case KEYWORD:
		{
			switch token.Value {
			case FUNC:
				{
					return parser.ParseFunction()
				}
			case TRY:
				{
					// work with a try catch block
					return parser.parseTryCatchBlock()
				}
			case FOR:
				{
					parser.pushToParsingState(LOOP_STATE)
					_for_node_ := parser.ParseForLoop()
					parser.popFromParsingState()

					return _for_node_
				}
			case IF:
				{
					return parser.ParseIfStatement()
				}
			case RAISE:
				{
					parser.IsExpectedEatElsePanic(
						parser.CurrentToken(),
						KEYWORD, RAISE,
						fmt.Sprintf("Expected a raise statement got %#v", parser.CurrentToken().Value),
					)

					return RaiseExceptionNode{Exception: parser._parseExpression()}
				}
			case BREAK:
				{
					// check the current state we are in
					if !parser.EnsureCurrentParsingStateIs(LOOP_STATE, false) {
						parser.reportError(
							token,
							"Break cannot be used outside a loop",
						)
					}

					parser.eatToken()

					return BreakNode{}
				}
			case RETURN:
				{
					// work with this and check
					if !parser.EnsureCurrentParsingStateIs(FUNCTION_STATE, true) {
						// panic here
						// we should not panic rather jump to an error state --> use gotos

						parser.reportError(
							token,
							"Return cannot be used outside a function definition",
						)

						// parser.ParsingError = "Return cannot be used outside a function definition"
						// goto parsingError
					}

					parser.eatToken()

					_expression := parser._parseExpression()

					if _expression == nil {
						_expression = NilNode{}
					}

					return ReturnNode{
						Expression: _expression,
					}
				}
			case EXPOSE:
				{
					parser.eatToken()

					// the node returned should be of type Assignment(not a reassignment) FunctionDecl
					_toBeExported := parser._parse(parser.CurrentToken())

					if _isExportable, ok := _toBeExported.(IExportables); ok {
						if !_isExportable.IsExported() {
							goto failedToExport
						}

						// create a an expose node
						return ExportVisibilityNode{
							Exported: _toBeExported,
						}
					}

				failedToExport:
					panic("Can't place an expose statement before that")
				}
			case DEF, CONST:
				{
					_currentTokenType := parser.CurrentToken().Value
					parser.eatToken() // eat the def or const keyword or const key

					lvalue := parser.CurrentToken().Value.(string)

					parser.eatToken() // eat the variable name

					parser.IsExpectedEatElsePanic(
						parser.CurrentToken(),
						ASSIGN, "=",
						"Expected '='",
					)

					rvalue := parser.ParseAssignment()
					_type := ASSIGNMENT

					if _currentTokenType == CONST {
						// change the type to const type
						_type = CONST_ASSIGNMENT // used to check later when the person tries to reassign
					}

					return Assignment{
						Type:   _type,
						Lvalue: lvalue,
						Rvalue: rvalue,
					}
				}
			case IMPORT:
				{
					parser.eatToken() // eat the import keyword
					fileName := parser.CurrentToken()
					parser.eatToken() // eat the filename
					var alias string

					// we check if we have an alias thing if so get it and ensure its a * or a variable
					if IsTypeAndValue(parser.CurrentToken(), KEYWORD, AS) {
						// we then need to get the alias
						parser.eatToken()

						// check the alias value
						_currentToken := parser.CurrentToken()

						if _currentToken.Type == VARIABLE || (_currentToken.Type == OPERATOR && _currentToken.Value.(string) == "*") {
							alias = _currentToken.Value.(string)
							parser.eatToken()
							goto end
						}

						panic("Expected a variable or * for an alias")
					}

				end:
					return Import{
						FileName: fileName.Value.(string),
						Alias:    alias,
					}
				}
			default:
				{
					return parser._parseExpression()
				}
			}
		}
	case VARIABLE:
		{
			if IsTypeAndValue(parser.peekAhead(), ASSIGN, "=") {
				parser.eatToken() // eat the variable
				parser.eatToken() // eat the assignment operator

				lvalue := parser.ParseAssignment()

				return Assignment{
					Type:   REASSIGNMENT,
					Lvalue: token.Value.(string),
					Rvalue: lvalue,
				}

			}

			// why the fuck does this not execute
			return parser._parseExpression()
		}

	case COMMENT:
		{
			// increment the fuck out of this
			parser.eatToken()

			// we should return nothing
			return CommentNode{
				comment: token.Value.(string),
			}
		}
	case CURLY_BRACES:
		{
			return parser.parseBlockScope()
		}
	default:
		{
			// we try to pass the expression
			return parser._parseExpression()
		}
	}

	// return something hapa
	// just return the error
	// set an error state
	// parsingError:
	// 	return nil
}

func (parser *Parser) Parse() *ProgramNode {
	program := &ProgramNode{}

	// start parsing the tokens here make everything to do its own complete shit
	for parser.CurrentPosition < parser.TokensLength {
		program.Nodes = append(program.Nodes, parser._parse(parser.CurrentToken()))
		// parser.CurrentPosition += 1
	}

	return program
}
