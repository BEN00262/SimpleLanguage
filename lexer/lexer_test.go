package lexer

import "testing"

// test the lexing of the system
func TestTokenGeneration(t *testing.T) {

	lexer := InitLexer("def name = 90")

	// we get an array of lexeme lets test them
	tokens := lexer.Lex()
	expectedTokens := []Token{
		{
			Type:  KEYWORD,
			Value: "def",
		},
		{
			Type:  VARIABLE,
			Value: "name",
		},
		{
			Type:  ASSIGN,
			Value: "=",
		},
		{
			Type:  NUMBER,
			Value: "90",
		},
	}

	// match the two arrays
	for index, token := range tokens {
		_token := expectedTokens[index]

		if _token.Type != token.Type {
			t.Fatalf("%v != %v", _token, token)
		}
	}
}
