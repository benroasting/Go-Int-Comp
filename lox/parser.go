package lox

type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
}

type parseError struct {}

func (p *Parser) Parse() (expr Expr) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(parseError); ok {
				expr = nil
			} else {
				panic(r)
			}
		}
	}()
	return p.expression()
}

func NewParser(tokens []Token, l *Lox) *Parser {
	return &Parser{tokens: tokens, lox: l}
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(Greater, GreaterEqual, Less, LessEqual){
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(Minus, Plus){
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(Slash, Star){
		operator := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus){
		operator := p.previous()
		right := p.unary()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(False) {
		return &Literal{Value: false}
	}
	if p.match(True) {
		return &Literal{Value: true}
	}
	if p.match(Nil){
		return &Literal{Value: nil}
	}
	if p.match(Number, String) {
		return &Literal{Value: p.previous().Literal}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return &Grouping{Expression: expr}
	}

	panic(p.parseError(p.peek(), "Expect expression."))
}

func (p *Parser) parseError(t Token, msg string) parseError {
	p.lox.TokenError(t, msg)
	return parseError{}
}

func (p *Parser) consume(t TokenType, msg string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic(p.parseError(p.peek(), msg))
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}
		p.advance()
	}
}
