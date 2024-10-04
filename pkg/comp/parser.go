package comp

type Parser struct {
	lexer          *Lexer
	emitter        *Emitter
	curToken       *Token
	peekToken      *Token
	symbols        map[string]struct{}
	labelsDeclared map[string]struct{}
	labelsGotoed   map[string]struct{}
}

func NewParser(l *Lexer, e *Emitter) *Parser {
	p := &Parser{
		lexer:          l,
		emitter:        e,
		symbols:        make(map[string]struct{}),
		labelsDeclared: make(map[string]struct{}),
		labelsGotoed:   make(map[string]struct{}),
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) checkToken(kind TokenType) bool {
	return kind == p.curToken.Kind()
}

func (p *Parser) checkPeek(kind TokenType) bool {
	return kind == p.peekToken.Kind()
}

func (p *Parser) match(kind TokenType) {
	if !p.checkToken(kind) {
		p.abort("Expected " + kind.String() + ", got" + p.curToken.Kind().String())
	}
	p.nextToken()
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.GetToken()
}

func (p *Parser) abort(msg string) {
	panic("Parser failed: " + msg)
}

func (p *Parser) Program() {
	p.emitter.headerLine("#include <stdio.h>")
	p.emitter.headerLine("int main(void) {")

	for p.checkToken(NEWLINE) {
		p.nextToken()
	}

	for !p.checkToken(EOF) {
		p.statement()
	}

	p.emitter.emitLine("return 0;")
	p.emitter.emitLine("}")

	for k := range p.labelsGotoed {
		if _, isPresent := p.labelsDeclared[k]; !isPresent {
			p.abort("Attempting to GOTO undeclared label: " + k)
		}
	}
}

func (p *Parser) statement() {
	// "PRINT" (expression | string)
	if p.checkToken(PRINT) {
		p.nextToken()
		if p.checkToken(STRING) {
			p.emitter.emitLine("printf(\"" + p.curToken.Text() + "\\n\");")
			p.nextToken()
		} else {
			p.emitter.emit("printf(\"%" + ".2f\\n\", (float) (")
			p.expression()
			p.emitter.emitLine("));")
		}
	} else if p.checkToken(IF) {
		// "IF" comparison "THEN" {statement} "ENDIF"
		p.nextToken()
		p.emitter.emit("if (")
		p.comparison()
		p.match(THEN)
		p.nl()
		p.emitter.emitLine(") {")
		for !p.checkToken(ENDIF) {
			p.statement()
		}
		p.match(ENDIF)
		p.emitter.emitLine("}")
	} else if p.checkToken(WHILE) {
		// "WHILE" comparison "REPEAT" {statement} "ENDWHILE"
		p.nextToken()
		p.emitter.emit("while (")
		p.comparison()
		p.match(REPEAT)
		p.nl()
		p.emitter.emitLine(") {")
		for !p.checkToken(ENDWHILE) {
			p.statement()
		}
		p.match(ENDWHILE)
		p.emitter.emitLine("}")
	} else if p.checkToken(LABEL) {
		// "LABEL" ident
		p.nextToken()
		if _, isPresent := p.labelsDeclared[p.curToken.Text()]; isPresent {
			p.abort("Label already exists: " + p.curToken.Text())
		}
		p.labelsDeclared[p.curToken.Text()] = struct{}{}
		p.emitter.emitLine(p.curToken.Text() + ":")
		p.match(IDENT)
	} else if p.checkToken(GOTO) {
		// "GOTO" ident
		p.nextToken()
		p.labelsGotoed[p.curToken.Text()] = struct{}{}
		p.emitter.emitLine("goto " + p.curToken.Text() + ";")
		p.match(IDENT)
	} else if p.checkToken(LET) {
		// "LET" ident "=" expression
		p.nextToken()
		if _, isPresent := p.symbols[p.curToken.Text()]; !isPresent {
			p.symbols[p.curToken.Text()] = struct{}{}
			p.emitter.headerLine("float " + p.curToken.Text() + ";")
		}
		p.emitter.emit(p.curToken.Text() + " = ")
		p.match(IDENT)
		p.match(EQ)
		p.expression()
		p.emitter.emitLine(";")
	} else if p.checkToken(INPUT) {
		// "INPUT" ident
		p.nextToken()
		if _, isPresent := p.symbols[p.curToken.Text()]; !isPresent {
			p.symbols[p.curToken.Text()] = struct{}{}
			p.emitter.headerLine("float " + p.curToken.Text() + ";")
		}
		p.emitter.emitLine("if (0 == scanf(\"%" + "f\", &" + p.curToken.Text() + ")) {")
		p.emitter.emitLine(p.curToken.Text() + " = 0;")
		p.emitter.emit("scanf(\"%")
		p.emitter.emitLine("*s\");")
		p.emitter.emitLine("}")
		p.match(IDENT)
	} else {
		// Invalid statement. Error!
		p.abort("Invalid statement at " + p.curToken.Text() + " (" + p.curToken.Kind().String() + ")")
	}

	// Newline.
	p.nl()
}

func (p *Parser) expression() {
	// expression ::= term {( "-" | "+") term}
	p.term()
	for p.checkToken(PLUS) || p.checkToken(MINUS) {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
		p.term()
	}
}

func (p *Parser) term() {
	// term ::= unary {( "/" | "*") unary}
	p.unary()
	for p.checkToken(ASTERISK) || p.checkToken(SLASH) {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
		p.unary()
	}
}

func (p *Parser) unary() {
	// unary ::= ["+" | "-"] primary
	if p.checkToken(PLUS) || p.checkToken(MINUS) {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
	}
	p.primary()
}

func (p *Parser) primary() {
	// primary ::= number | ident
	if p.checkToken(NUMBER) {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
	} else if p.checkToken(IDENT) {
		// Ensure the variable already exists.
		if _, isPresent := p.symbols[p.curToken.Text()]; !isPresent {
			p.abort("Referencing variable before assignement: " + p.curToken.Text())
		}
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
	} else {
		p.abort("Unexpected token at " + p.curToken.Text())
	}
}

func (p *Parser) comparison() {
	p.expression()
	if p.isComparisonOperator() {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
		p.expression()
	} else {
		p.abort("Expected comparison operator at: " + p.curToken.Text())
	}
	for p.isComparisonOperator() {
		p.emitter.emit(p.curToken.Text())
		p.nextToken()
		p.expression()
	}
}

func (p *Parser) isComparisonOperator() bool {
	return p.checkToken(GT) ||
		p.checkToken(GTEQ) ||
		p.checkToken(LT) ||
		p.checkToken(LTEQ) ||
		p.checkToken(EQEQ) ||
		p.checkToken(NOTEQ)
}

func (p *Parser) nl() {
	// nl ::= '\n'+

	// Require at least one newline.
	p.match(NEWLINE)
	// But we will allow extra newlines too, of course.
	for p.checkToken(NEWLINE) {
		p.nextToken()
	}
}
