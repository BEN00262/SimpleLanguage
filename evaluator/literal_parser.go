package evaluator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	CODE_BLOCK   = regexp.MustCompile(`#!START_CODE((\s*?.*?)*?)#!END_CODE\n*`)
	HEADER_BLOCK = regexp.MustCompile(`#!H(\d)((\s*?.*?)*?)#!H(\d)`)
)

type LiteralParsing struct {
	KnuthCode string
	filePath  string
}

func InitLiteralParsing(filePath string, knuthCode string) *LiteralParsing {
	return &LiteralParsing{
		KnuthCode: knuthCode,
		filePath:  filePath,
	}
}

func (literal *LiteralParsing) ExecuteLiteralCode() interface{} {
	var codeStrings []string

	for _, codeBlock := range CODE_BLOCK.FindAllStringSubmatch(literal.KnuthCode, -1) {
		if codeBlock != nil {
			codeStrings = append(codeStrings, codeBlock[1])
		}
	}

	return Interpreter(literal.filePath, strings.Join(codeStrings, "\n"))
}

func (literal *LiteralParsing) GenerateDocumentation() string {
	for _, headerBlock := range HEADER_BLOCK.FindAllStringSubmatch(literal.KnuthCode, -1) {
		if headerBlock != nil {
			lhs := headerBlock[1]
			rhs := headerBlock[4]

			if strings.Compare(lhs, rhs) != 0 {
				continue
			}

			tag := "h1"

			switch lhs {
			case "2":
				tag = "h2"
			case "3":
				tag = "h3"
			case "4":
				tag = "h4"
			case "5":
				tag = "h5"
			case "6":
				tag = "h6"
			}

			literal.KnuthCode = strings.ReplaceAll(literal.KnuthCode, headerBlock[0], fmt.Sprintf(
				"<%s>%s</%s>\n", tag, headerBlock[2], tag),
			)
		}
	}

	for _, codeBlock := range CODE_BLOCK.FindAllStringSubmatch(literal.KnuthCode, -1) {
		if codeBlock != nil {
			literal.KnuthCode = strings.ReplaceAll(literal.KnuthCode, codeBlock[0], fmt.Sprintf(
				"<pre>%s</pre>\n", codeBlock[1]),
			)
		}
	}

	return literal.KnuthCode
}
