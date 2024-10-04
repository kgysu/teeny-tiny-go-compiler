package comp

import "errors"

type Token struct {
	text string
	kind TokenType
}

func (t *Token) Kind() TokenType {
	return t.kind
}
func (t *Token) Text() string {
	return t.text
}

func NewToken(tokenText string, tokenKind TokenType) *Token {
	return &Token{
		text: tokenText,
		kind: tokenKind,
	}
}

func CheckIfKeyword(tokText string) (TokenType, error) {
	for k, v := range typeName {
		if v == tokText && k >= 100 && k < 200 {
			return k, nil
		}
	}
	return 0, errors.New("")
}
