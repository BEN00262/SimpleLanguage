package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	. "github.com/BEN00262/simpleLang/evaluator"
)

var (
	fileName  = flag.String("filename", "", "< filename > filename to execute")
	mode      = flag.String("mode", "e", "< mode > i or e or l mode to run the interpreter") // l ( generate doc or execute the code )
	operation = flag.String("operation", "d", "< operation > e or d operation to run the literal compiler in")
	outfile   = flag.String("out", "", "< out file > filename to write the documentation to")
)

func getFileData(filename string) string {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func main() {
	flag.Parse()

	switch *mode {
	case "e":
		{
			if *fileName == "" {
				flag.Usage()
				return
			}

			Interpreter(getFileData(*fileName))
		}
	case "i":
		{
			REPL()
		}
	case "l":
		{
			if *fileName == "" {
				flag.Usage()
				return
			}

			literalParser := InitLiteralParsing(getFileData(*fileName))

			switch *operation {
			case "e":
				{
					literalParser.ExecuteLiteralCode()
				}
			case "d":
				{
					htmlDocumentation := literalParser.GenerateDocumentation()

					if *outfile == "" {
						fmt.Println(htmlDocumentation)
						return
					}

					if err := ioutil.WriteFile(*outfile, []byte(htmlDocumentation), 0600); err != nil {
						panic(err.Error())
					}
				}
			}

		}
	}
}

// func main() {
// 	// experiment with calling function from Daisy

// 	evaluator := NewEvaluatorContext()

// 	// inject all global functions
// 	evaluator.InitGlobalScope()

// 	LoadGlobalsToContext(evaluator)

// 	lexer := InitLexer(`
// 	fun printName(first, second) {
// 		return first + second
// 	}
// 	`)
// 	parser := InitParser(lexer.Lex())
// 	evaluator.ReplExecute(parser.Parse())

// 	value := DaisyInvoke(evaluator, "printName", ToDaisy(400), ToDaisy(23))

// 	fmt.Println(FromDaisy(value))

// 	evaluator.TearDownRepl()
// }
