package lexer

var (
	KEYWORDS = []string{
		FUNC, TRUE,
		FALSE, FOR, IF, ELSE,
		BREAK, RETURN, NIL, IMPORT,
		AS, EXPORT, DEF, TRY, CATCH,
		FINALLY, RAISE, CONST, EXPOSE, NAMESPACE,
	}
)

type TokenType = int

const (
	NUMBER              TokenType = iota + 1
	STRING                        // 2
	VARIABLE                      // 3
	KEYWORD                       // 4
	ASSIGN                        // 5
	OPERATOR                      // 6
	COMMENT                       // 7
	HALF_CIRCLE_BRACKET           // 8
	CURLY_BRACES                  // 9
	COMMA                         //10
	SEMI_COLON                    // 11
	CONDITION                     //12
	SQUARE_BRACKET                //13
	COLON                         //14
	DOT                           //15
)

type Token struct {
	Type        TokenType
	Line        int
	ColumnStart int // start of the column
	Span        int // end of the column
	Value       interface{}
}
