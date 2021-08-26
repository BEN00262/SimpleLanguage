package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

			fmt.Println(
				Interpreter(getFileData(*fileName)),
			)
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

			literalParser := initLiteralParsing(getFileData(*fileName))

			switch *operation {
			case "e":
				{
					fmt.Println(literalParser.ExecuteLiteralCode())
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
