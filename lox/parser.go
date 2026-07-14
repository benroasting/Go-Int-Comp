package lox

type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
}

type parseError struct{}

func NewParser(tokens []Token, l *Lox) *Parser {
	return &Parser{tokens: tokens, lox: l}
}

func (p *Parser) Parse() (stmts []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(parseError); ok {
				stmts = nil
			} else {
				panic(r)
			}
		}
	}()
	for !p.isAtEnd() {
		s := p.declaration()
		if s != nil {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

// declaration
func (p *Parser) declaration() (stmt Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(parseError); ok {
				p.synchronize()
				stmt = nil
			} else {
				panic(r)
			}
		}
	}()
	if p.match(Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(Identifier, "Expect variable name.")
	var initializer Expr
	if p.match(Equal) {
		initializer = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after variable declaration.")
	return &VarStmt{Name: name, Initializer: initializer}
}

// statements
func (p *Parser) statement() Stmt {
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(LeftBrace) {
		return &BlockStmt{Statements: p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		s := p.declaration()
		if s != nil {
			statements = append(statements, s)
		}
	}
	p.consume(RightBrace , "Expect '}' after block.")
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(Semicolon, "Expect ';' after value.")
	return &PrintStmt{Expression: value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(Semicolon, "Expect ';' after expression.")
	return &ExpressionStmt{Expression: expr}
}

// expressions
func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(Equal) {
		equals := p.previous()
		value := p.assignment()

		if v, ok := expr.(*Variable); ok {
			return &Assign{Name: v.Name, Value: value}
		}

		p.parseError(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(Or) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(And) {
		operator := p.previous()
		right := p.equality()
		expr = &Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr
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

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus) {
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
	if p.match(Nil) {
		return &Literal{Value: nil}
	}
	if p.match(Number, String) {
		return &Literal{Value: p.previous().Literal}
	}
	if p.match(Identifier) {
		return &Variable{Name: p.previous()}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return &Grouping{Expression: expr}
	}

	panic(p.parseError(p.peek(), "Expect expression."))
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

func (p *Parser) consume(t TokenType, msg string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic(p.parseError(p.peek(), msg))
}

func (p *Parser) parseError(t Token, msg string) parseError {
	p.lox.TokenError(t, msg)
	return parseError{}
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

func (p *Parser) ifStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(Else) {
		elseBranch =p.statement()
	}

	return &IfStmt{
		Condition: condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(RightParen) {
		condition = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RightParen) {
		increment = p.expression()
	}
	p.consume(RightParen, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &BlockStmt{Statements: []Stmt{
			body,
			&ExpressionStmt{Expression: increment},
		}}
	}

	if condition == nil {
		condition = &Literal{Value: true}
	}
	body = &WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = &BlockStmt{Statements: []Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after 'while'.")
	body := p.statement()
	return &WhileStmt{Condition: condition, Body: body}
}
