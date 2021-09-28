package lexer

import (
	"math/big"
	"strings"
	"unicode"
)

func Contains(array []string, str string) bool {
	for _, a := range array {
		if strings.Compare(a, str) == 0 {
			return true
		}
	}

	return false
}

type Lexer struct {
	code            string
	tokens          []Token
	currentPosition int
	line            int // get the line and start it
	column          int // the code column --> get the start column
}

func InitLexer(code string) *Lexer {
	return &Lexer{
		code:            code,
		currentPosition: 0,
	}
}

func (lexer *Lexer) CurrentLexeme() rune {
	return rune(lexer.code[lexer.currentPosition])
}

func (lexer *Lexer) peakAhead() rune {
	return rune(lexer.code[lexer.currentPosition+1])
}

func (lexer *Lexer) eatLexeme() {
	lexer.currentPosition++
}

// this appends the token to the system
func (lexer *Lexer) AppendToken(tokenType TokenType, value interface{}, colStart int, colEnd int) {

	// use a column span
	lexer.tokens = append(lexer.tokens, Token{
		Type:        tokenType,
		Line:        lexer.line,
		Value:       value,
		ColumnStart: colStart,
		ColumnEnd:   colEnd,
	})
}

func isValidVariableName(lexeme rune, continuous bool) bool {
	_start := unicode.IsLetter(lexeme) || lexeme == '_'

	if continuous {
		return _start || unicode.IsDigit(lexeme)
	}

	return _start
}

func (lexer *Lexer) Lex() []Token {
	// var tokens []Token

	for lexer.currentPosition < len(lexer.code) {
		lexeme := lexer.CurrentLexeme()

		// use a regex to match all the letters here
		// or the lexeme is _

		if isValidVariableName(lexeme, false) {
			_currentPosition := lexer.currentPosition

			for ; lexer.currentPosition < len(lexer.code) && isValidVariableName(lexer.CurrentLexeme(), true); lexer.currentPosition++ {

			}

			// what happens is if we have one stuff eg
			_variable := lexer.code[_currentPosition:lexer.currentPosition]
			// lexer.currentPosition -= 1

			tokenType := VARIABLE

			if Contains(KEYWORDS, _variable) {
				tokenType = KEYWORD
			}

			lexer.AppendToken(tokenType, _variable, _currentPosition, lexer.currentPosition-1)

			// tokens = append(tokens, Token{
			// 	Type:  tokenType,
			// 	Value: _variable,
			// })
			continue
		} else if unicode.IsDigit(lexeme) {
			_start_position_ := lexer.currentPosition

			// we want to get all the number

			for ; lexer.currentPosition < len(lexer.code) && unicode.IsDigit(lexer.CurrentLexeme()); lexer.currentPosition++ {

			}

			_number_ := lexer.code[_start_position_:lexer.currentPosition]

			lexer.currentPosition -= 1

			// we should use a different value not an integer
			_number, isSuccess := new(big.Int).SetString(_number_, 0)

			if !isSuccess {
				panic("Failed to convert number")
			}

			lexer.AppendToken(NUMBER, *_number, _start_position_, lexer.currentPosition-1)

			// tokens = append(tokens, Token{
			// 	Type:  NUMBER,
			// 	Value: *_number,
			// })
		} else if lexeme == ';' {
			// tokens = append(tokens, Token{
			// 	Type:  SEMI_COLON,
			// 	Value: ";",
			// })

			lexer.AppendToken(SEMI_COLON, ";", lexer.currentPosition, lexer.currentPosition)
		} else if lexeme == '!' {
			_current_column_ := lexer.currentPosition
			// _peakAheadLexeme := lexer.peakAhead()

			if lexer.peakAhead() == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: "!=",
				// })

				lexer.AppendToken(CONDITION, "!=", _current_column_, lexer.currentPosition-1)
			}

		} else if lexeme == '=' {
			// check the next if its a another =
			_current_column_ := lexer.currentPosition
			// return the condition token
			// _peakAheadLexeme := lexer.peakAhead()

			if lexer.peakAhead() == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: "==",
				// })

				lexer.AppendToken(CONDITION, "==", _current_column_, lexer.currentPosition)
			} else {
				// tokens = append(tokens, Token{
				// 	Type:  ASSIGN,
				// 	Value: "=",
				// })

				lexer.AppendToken(ASSIGN, "=", _current_column_, _current_column_)
			}
		} else if strings.Contains(`"'`, string(lexeme)) {
			// we do have the starting point go on until we reach the end of the strings
			_initial_position := lexer.currentPosition + 1

			// eat the current lexeme
			lexer.eatLexeme()

			for ; lexer.currentPosition < len(lexer.code) && lexer.CurrentLexeme() != lexeme; lexer.currentPosition++ {
			}

			_string_value_ := lexer.code[_initial_position:lexer.currentPosition]

			lexer.AppendToken(STRING, _string_value_, _initial_position-1, lexer.currentPosition-1)
			// tokens = append(tokens, Token{
			// 	Type:  STRING,
			// 	Value: _string_value_,
			// })

		} else if Contains([]string{">", "<"}, string(lexeme)) {
			_peakAheadLexeme := lexer.peakAhead()
			_current_column_ := lexer.currentPosition

			if _peakAheadLexeme == '=' {
				lexer.eatLexeme()

				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: string(lexeme) + string(_peakAheadLexeme),
				// })

				lexer.AppendToken(CONDITION, string(lexeme)+string(_peakAheadLexeme), _current_column_, lexer.currentPosition)
			} else {
				// tokens = append(tokens, Token{
				// 	Type:  CONDITION,
				// 	Value: string(lexeme),
				// })

				lexer.AppendToken(CONDITION, string(lexeme), _current_column_, lexer.currentPosition)
			}
		} else if strings.Contains("+-*/%", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  OPERATOR,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(OPERATOR, string(lexeme), lexer.currentPosition, lexer.currentPosition)
		} else if lexeme == '.' {
			// tokens = append(tokens, Token{
			// 	Type:  DOT,
			// 	Value: ".",
			// })

			lexer.AppendToken(DOT, ".", lexer.currentPosition, lexer.currentPosition)
		} else if lexeme == ':' {
			// return the lexeme
			// tokens = append(tokens, Token{
			// 	Type:  COLON,
			// 	Value: ":",
			// })

			lexer.AppendToken(COLON, ":", lexer.currentPosition, lexer.currentPosition)
		} else if strings.Contains("[]", string(lexeme)) {
			// return the token now
			// tokens = append(tokens, Token{
			// 	Type:  SQUARE_BRACKET,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(SQUARE_BRACKET, string(lexeme), lexer.currentPosition, lexer.currentPosition)
		} else if lexeme == ',' {
			// add the comma tokens for this very reason
			// tokens = append(tokens, Token{
			// 	Type:  COMMA,
			// 	Value: ",",
			// })

			lexer.AppendToken(COMMA, ",", lexer.currentPosition, lexer.currentPosition)
		} else if lexeme == '\n' {
			lexer.line += 1
			lexer.column = 0
		} else if lexeme == '#' {
			lexer.eatLexeme()

			// #[ comments ]# for multiline comments

			// if the nextToken is a [
			// eat all the characters until we find ]#

			_start_position_ := lexer.currentPosition

			for ; lexer.currentPosition < len(lexer.code) && lexer.CurrentLexeme() != '\n'; lexer.currentPosition++ {

			}

			_comment_ := lexer.code[_start_position_:lexer.currentPosition]
			lexer.currentPosition -= 1

			// tokens = append(tokens, Token{
			// 	Type:  COMMENT,
			// 	Value: _comment_,
			// })

			lexer.AppendToken(COMMENT, _comment_, _start_position_, lexer.currentPosition)

		} else if strings.Contains("()", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  HALF_CIRCLE_BRACKET,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(HALF_CIRCLE_BRACKET, string(lexeme), lexer.currentPosition, lexer.currentPosition)
		} else if strings.Contains("{}", string(lexeme)) {
			// tokens = append(tokens, Token{
			// 	Type:  CURLY_BRACES,
			// 	Value: string(lexeme),
			// })

			lexer.AppendToken(CURLY_BRACES, string(lexeme), lexer.currentPosition, lexer.currentPosition)
		}

		lexer.currentPosition += 1
	}
	return lexer.tokens
}
