package main

import (
	"strconv"
	"unicode"
)

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

// return a list of Tokens
func (lexer *Lexer) Lex() []Token {
	var tokens []Token

	for ; lexer.currentPosition < len(lexer.code); lexer.currentPosition++ {
		lexeme := lexer.CurrentLexeme()

		if unicode.IsLetter(lexeme) {
			tokens = append(tokens, Token{
				Type:  VARIABLE,
				Value: string(lexeme),
			})
		} else if unicode.IsDigit(lexeme) {

			// convert the string to a float 64
			intValue, err := strconv.Atoi(string(lexeme))

			if err != nil {
				panic(err.Error())
			}

			tokens = append(tokens, Token{
				Type:  NUMBER,
				Value: intValue,
			})
		}

		lexer.currentPosition += 1
	}

	return tokens
}
