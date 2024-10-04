package comp

import (
	"os"
	"unicode"
)

type Lexer struct {
	source  string
	curChar rune
	curPos  int
}

func (l *Lexer) GetCurChar() rune {
	return l.curChar
}

func (l *Lexer) GetCurPos() int {
	return l.curPos
}

func NewLexer(s string) *Lexer {
	return &Lexer{
		source:  s + "\n",
		curChar: ' ',
		curPos:  -1,
	}
}

func (l *Lexer) NextChar() {
	l.curPos += 1
	if l.curPos >= len(l.source) {
		l.curChar = '\000' // EOF
	} else {
		l.curChar = []rune(l.source)[l.curPos]
	}
}

func (l *Lexer) Peek() rune {
	if l.curPos+1 >= len(l.source) {
		return '\000'
	}
	return []rune(l.source)[l.curPos+1]
}

func (l *Lexer) abort(msg string) {
	println("Lexing error. " + msg)
	os.Exit(0)
}

func (l *Lexer) skipWhitespace() {
	for l.curChar == ' ' ||
		l.curChar == '\t' ||
		l.curChar == '\r' {
		l.NextChar()
	}
}

func (l *Lexer) skipComment() {
	if l.curChar == '#' {
		for l.curChar != '\n' {
			l.NextChar()
		}
	}
}

func (l *Lexer) GetToken() *Token {
	l.skipWhitespace()
	l.skipComment()
	var token *Token

	//switch l.curChar {
	if l.curChar == '+' {
		token = NewToken(string(l.curChar), PLUS)
	} else if l.curChar == '-' {
		token = NewToken(string(l.curChar), MINUS)
	} else if l.curChar == '*' {
		token = NewToken(string(l.curChar), ASTERISK)
	} else if l.curChar == '/' {

		token = NewToken(string(l.curChar), SLASH)
	} else if l.curChar == '=' {
		if l.Peek() == '=' {
			lastChar := l.curChar
			l.NextChar()
			token = NewToken(string(lastChar+l.curChar), EQEQ)
		} else {
			token = NewToken(string(l.curChar), EQ)
		}
	} else if l.curChar == '<' {
		if l.Peek() == '=' {
			lastChar := l.curChar
			l.NextChar()
			token = NewToken(string(lastChar+l.curChar), LTEQ)
		} else {
			token = NewToken(string(l.curChar), LT)
		}
	} else if l.curChar == '>' {
		if l.Peek() == '=' {
			lastChar := l.curChar
			l.NextChar()
			token = NewToken(string(lastChar+l.curChar), GTEQ)
		} else {
			token = NewToken(string(l.curChar), GT)
		}
	} else if l.curChar == '!' {
		if l.Peek() == '=' {
			lastChar := l.curChar
			l.NextChar()
			token = NewToken(string(lastChar+l.curChar), NOTEQ)
		} else {
			l.abort("Expected !=, got !" + string(l.Peek()))
		}
	} else if l.curChar == '"' {
		l.NextChar()
		startPos := l.curPos
		for l.curChar != '"' {
			if l.curChar == '\r' || l.curChar == '\n' || l.curChar == '\t' ||
				l.curChar == '\\' || l.curChar == '%' {
				l.abort("Illegal character in string.")
			}
			l.NextChar()
		}
		tokText := l.source[startPos:l.curPos]
		token = NewToken(tokText, STRING)
	} else if unicode.IsDigit(l.curChar) {
		startPos := l.curPos
		for unicode.IsDigit(l.Peek()) {
			l.NextChar()
		}
		if l.Peek() == '.' {
			l.NextChar()
			if !unicode.IsDigit(l.Peek()) {
				// Error
				l.abort("Illegal character in number.")
			}
			for unicode.IsDigit(l.Peek()) {
				l.NextChar()
			}
		}
		tokText := l.source[startPos : l.curPos+1]
		token = NewToken(tokText, NUMBER)
	} else if unicode.IsLetter(l.curChar) {
		startPos := l.curPos
		for unicode.IsDigit(l.Peek()) || unicode.IsLetter(l.Peek()) {
			l.NextChar()
		}

		tokText := l.source[startPos : l.curPos+1]
		keyword, err := CheckIfKeyword(tokText)
		if err != nil {
			token = NewToken(tokText, IDENT)
		} else {
			token = NewToken(tokText, keyword)
		}

	} else if l.curChar == '\n' {
		token = NewToken(string(l.curChar), NEWLINE)
	} else if l.curChar == '\000' {
		token = NewToken("", EOF)
	} else {
		// unknown token
		l.abort("Unknown token: " + string(l.curChar))
	}

	l.NextChar()
	return token
}
