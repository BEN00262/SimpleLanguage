package parser

import (
	"fmt"
	"strings"

	. "github.com/BEN00262/simpleLang/lexer"
	"github.com/gookit/color"
)

// do error reporting then exit
// get the error then exit the program without doing anything
func (parser *Parser) reportError(token Token) {
	red := color.FgRed.Render
	green := color.FgGreen.Render

	_raw_error_string := parser.ActualCode[token.Line][token.ColumnStart : token.ColumnStart+token.Span+1]

	fmt.Printf(
		"%s %s\n\n",
		green(token.Line+1),
		strings.Replace(parser.ActualCode[token.Line], _raw_error_string, red(_raw_error_string), 1),
	)
}
