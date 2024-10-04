package comp

type TokenType int

const (
	EOF     = -1
	NEWLINE = 0
	NUMBER  = 1
	IDENT   = 2
	STRING  = 3
	// Keywords.
	LABEL    = 101
	GOTO     = 102
	PRINT    = 103
	INPUT    = 104
	LET      = 105
	IF       = 106
	THEN     = 107
	ENDIF    = 108
	WHILE    = 109
	REPEAT   = 110
	ENDWHILE = 111
	// Operators.
	EQ       = 201
	PLUS     = 202
	MINUS    = 203
	ASTERISK = 204
	SLASH    = 205
	EQEQ     = 206
	NOTEQ    = 207
	LT       = 208
	LTEQ     = 209
	GT       = 210
	GTEQ     = 211
)

var typeName = map[TokenType]string{
	EOF:     "EOF",
	NEWLINE: "NEWLINE",
	NUMBER:  "Number",
	IDENT:   "Ident",
	STRING:  "String",
	// Keywords.
	LABEL:    "LABEL",
	GOTO:     "GOTO",
	PRINT:    "PRINT",
	INPUT:    "INPUT",
	LET:      "LET",
	IF:       "IF",
	THEN:     "THEN",
	ENDIF:    "ENDIF",
	WHILE:    "WHILE",
	REPEAT:   "REPEAT",
	ENDWHILE: "ENDWHILE",
	// Operators.
	EQ:       "Equal",
	PLUS:     "Plus",
	MINUS:    "Minus",
	ASTERISK: "Asterik",
	SLASH:    "Slash",
	EQEQ:     "EqEq",
	NOTEQ:    "NotEq",
	LT:       "LargerThan",
	LTEQ:     "LargerThanOrEq",
	GT:       "GraterThan",
	GTEQ:     "GraterThanOrEq",
}

func (tt TokenType) String() string {
	return typeName[tt]
}
