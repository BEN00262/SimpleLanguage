package evaluator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/parser"
)

func visit(files *[]interface{}) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		*files = append(*files, StringNode{
			Value: path,
		})

		return nil
	}
}

func _print(value interface{}) string {
	switch _value := value.(type) {
	case StringNode:
		{
			return fmt.Sprintf("%s", _value.Value)
		}
	case NumberNode:
		{
			// find the string stuff and then use it
			return fmt.Sprintf("%v", _value.Value.Text(10))
		}
	case BoolNode:
		{
			if _value.Value == 1 {
				return fmt.Sprint("True")
			} else {
				return fmt.Sprint("False")
			}
		}
	case ArrayNode:
		{
			var _arguments_ []string

			for _, _element := range _value.Elements {
				_arguments_ = append(_arguments_, _print(_element))
			}

			return "[ " + strings.Join(_arguments_, ",") + " ]"
		}
	}
	return ""
}

// work arrays into the system and also dictionaries
// [{Type: "DIR", filePath:"" }] --> we can then easily have ourselves a malware :)
// the encryption bit is easy ( just make the execution portion quite fast just that )
// add functionality to obfuscate the generated binaries
// thinking is we introduce an IR compile down the code to it then ship it
// lets test and see
// introduce global directives e.g #@ obfuscate and #@ compile
type ExternalDependencies = map[string]SymbolTableValue

var (
	GLOBALS = ExternalDependencies{
		// internal functions for the system to use :)
		// get a reference to an element in the array then work on it
		"push": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "push",
				ParamCount: 2,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					_arr_ := *value[0]
					_value_to_push_ := *value[1]

					switch _array_ := _arr_.(type) {
					case ArrayNode:
						{
							// does not change the value bana
							_array_.Push(_value_to_push_)
							return _array_, ExceptionNode{Type: NO_EXCEPTION}
						}
					}

					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		"pop": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "pop",
				ParamCount: 1,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					_arr_ := *value[0]

					switch _array_ := _arr_.(type) {
					case ArrayNode:
						{
							// return the value
							// what we do is just return the array
							_array_.Pop()
							return _array_, ExceptionNode{Type: NO_EXCEPTION}
						}
					}
					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		"insertAt": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "insertAt",
				ParamCount: 3, // array, index, value
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					_arr_ := *value[0]
					_index_ := *value[1]
					_value_ := *value[2]

					switch _array_ := _arr_.(type) {
					case ArrayNode:
						{
							switch _index := _index_.(type) {
							case NumberNode:
								{
									_array_.InsertAt(_index.Value.Int64(), _value_)
									return _array_, ExceptionNode{Type: NO_EXCEPTION}
								}
							}
						}
					}
					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		"len": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "len",
				ParamCount: 1,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					_first_argument_ := *value[0]

					if countable, ok := _first_argument_.(Countable); ok {
						return countable.Length(), ExceptionNode{Type: NO_EXCEPTION}
					}

					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},

		// type introspection
		"type": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "type",
				ParamCount: 1,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					// only take the first value the rest we dont need
					_argument := *value[0]
					_type := "nil"

					switch _val := _argument.(type) {
					case StringNode:
						_type = "string"
					case NumberNode:
						_type = "number"
					case BoolNode:
						_type = "boolean"
					case FunctionDecl, AnonymousFunction:
						_type = "function"
					case ExceptionNode:
						// we should have an exception type
						// should we return the actual type in the execption
						_type = _val.Type
					}

					return StringNode{
						Value: _type,
					}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},

		// Exception creator
		"Exception": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "Exception",
				ParamCount: 2,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					if _type, ok := (*value[0]).(StringNode); ok {
						if _message, ok := (*value[1]).(StringNode); ok {
							return ExceptionNode{
								Type:    _type.Value,
								Message: _message.Value,
							}, ExceptionNode{Type: NO_EXCEPTION}
						}
					}

					return ExceptionNode{
						Type:    "InvalidException",
						Message: "Invalid exception thrown",
					}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		// print
		"print": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "print",
				ParamCount: 1,
				Function: func(values ...*interface{}) (interface{}, ExceptionNode) {

					// read and write ( fd, format, data )
					// pass the data to be printed down here
					// how the tf are we going to do this

					for _, _value_ := range values {
						fmt.Printf("%s", _print(*_value_))
					}

					fmt.Println()
					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		// file system methods
		"openFile": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "openFile",
				ParamCount: 1,
				Function: func(filename ...*interface{}) (interface{}, ExceptionNode) {
					_filename := *filename[0]

					switch _file_ := _filename.(type) {
					case StringNode:
						{
							fileData, err := ioutil.ReadFile(_file_.Value)

							if err != nil {
								return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
							}

							return StringNode{
								Value: string(fileData),
							}, ExceptionNode{Type: NO_EXCEPTION}
						}
					}
					return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		// walks the filesystem's tree and returns the filenames as an array of strings which can be accessed from the language
		"walkFS": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "walkFS",
				ParamCount: 1,
				Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
					// start looping over the files and check them out
					_startDirectory := *value[0]
					var _files_ []interface{}

					// start here for now
					if directory, ok := _startDirectory.(StringNode); ok {
						filepath.Walk(directory.Value, visit(&_files_))
					}

					return ArrayNode{
						Elements: _files_,
					}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
		// write file
		// get the filename and the file data
		"writeFile": SymbolTableValue{
			Type: EXTERNALFUNC,
			Value: ExternalFunctionNode{
				Name:       "writeFile",
				ParamCount: 2,
				Function: func(values ...*interface{}) (interface{}, ExceptionNode) {
					// on success we return True
					_file_ := *values[0]
					_content_ := *values[1]

					// now write the data
					switch _filename_ := _file_.(type) {
					case StringNode:
						{
							// now check if also _content_ is a stringNode
							switch _file_content_ := _content_.(type) {
							case StringNode:
								{
									// write the data to the file quick
									if err := ioutil.WriteFile(_filename_.Value, []byte(_file_content_.Value), 0600); err != nil {
										return BoolNode{
											Value: 0,
										}, ExceptionNode{Type: NO_EXCEPTION}
									}

									// a successful write of the file to the system
									return BoolNode{
										Value: 1,
									}, ExceptionNode{Type: NO_EXCEPTION}
								}
							}

						}
					}

					// if it fails it returns a false
					return BoolNode{
						Value: 0,
					}, ExceptionNode{Type: NO_EXCEPTION}
				},
			},
		},
	}
)

// introduce a byte node to hold byte data -> i think that will be good

func LoadGlobalsToContext(eval *Evaluator) {
	for Key, Value := range GLOBALS {
		eval.InjectIntoGlobalScope(Key, Value)
	}

	// inject the eval functionality
	eval.InjectIntoGlobalScope("eval", SymbolTableValue{
		Type: EXTERNALFUNC,
		Value: ExternalFunctionNode{
			Name:       "eval",
			ParamCount: 1,
			Function: func(value ...*interface{}) (interface{}, ExceptionNode) {
				codeString := *value[0]

				if _code, ok := codeString.(StringNode); ok {
					return eval._eval(_code.Value)
				}

				return NilNode{}, ExceptionNode{Type: NO_EXCEPTION}
			},
		},
	})
}
