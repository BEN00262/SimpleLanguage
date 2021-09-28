package parser

import (
	"fmt"
	"os"

	. "github.com/BEN00262/simpleLang/lexer"
	"github.com/gookit/color"
)

// do error reporting then exit
// get the error then exit the program without doing anything
func (parser *Parser) reportError(token Token, errorMessage ...string) {
	red := color.FgRed.Render
	green := color.FgGreen.Render

	for _, message := range errorMessage {
		color.BgYellow.Println(red(message))
	}

	// input[:index] + string(replacement) + input[index+1:]
	_actual_line_ := parser.ActualCode[token.Line]
	_end_index_ := token.ColumnStart + token.Span + 1

	fmt.Printf(
		"%s %s\n\n",
		green(token.Line+1),
		_actual_line_[:token.ColumnStart]+red(_actual_line_[token.ColumnStart:_end_index_])+_actual_line_[_end_index_:],
	)

	os.Exit(-1)
}
