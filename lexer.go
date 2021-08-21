package main

import (
	"strconv"
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
	currentPosition int
}

func initLexer(code string) *Lexer {
	return &Lexer{
		code:            code,
		currentPosition: 0,
	}
}

func (lexer *Lexer) CurrentLexeme() rune {
	return rune(lexer.code[lexer.currentPosition])
}

func (lexer *Lexer) eatLexeme() {
	lexer.currentPosition++
}

func (lexer *Lexer) Lex() []Token {
	var tokens []Token

	for lexer.currentPosition < len(lexer.code) {
		lexeme := lexer.CurrentLexeme()

		if unicode.IsLetter(lexeme) {
			_currentPosition := lexer.currentPosition
			_currentLexeme := lexer.CurrentLexeme()

			for ; unicode.IsLetter(_currentLexeme); lexer.currentPosition++ {
				_currentLexeme = lexer.CurrentLexeme()
			}

			lexer.currentPosition -= 1

			_variable := lexer.code[_currentPosition:lexer.currentPosition]
			tokenType := VARIABLE

			if Contains(KEYWORDS, _variable) {
				tokenType = KEYWORD
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: _variable,
			})

			continue
		} else if unicode.IsDigit(lexeme) {
			intValue, err := strconv.Atoi(string(lexeme))

			if err != nil {
				panic(err.Error())
			}

			tokens = append(tokens, Token{
				Type:  NUMBER,
				Value: intValue,
			})
		} else if lexeme == ',' {
			tokens = append(tokens, Token{
				Type:  COMMA,
				Value: lexeme,
			})
		} else if lexeme == '=' {
			tokens = append(tokens, Token{
				Type:  ASSIGN,
				Value: lexeme,
			})
		} else if strings.Contains("+-*/", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  OPERATOR,
				Value: string(lexeme),
			})
		} else if lexeme == '#' {
			lexer.eatLexeme()

			_currentLexeme := lexer.CurrentLexeme()
			_currentPosition := lexer.currentPosition

			for ; _currentLexeme != '\n'; lexer.currentPosition++ {
				_currentLexeme = lexer.CurrentLexeme()
			}

			lexer.currentPosition -= 1

			tokens = append(tokens, Token{
				Type:  COMMENT,
				Value: lexer.code[_currentPosition:lexer.currentPosition],
			})
			continue
		} else if strings.Contains("()", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  HALF_CIRCLE_BRACKET,
				Value: lexeme,
			})
		} else if strings.Contains("{}", string(lexeme)) {
			tokens = append(tokens, Token{
				Type:  CURLY_BRACES,
				Value: lexeme,
			})
		}

		lexer.currentPosition += 1
	}

	return tokens
}
