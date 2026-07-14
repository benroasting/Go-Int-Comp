package lox

type Expr interface {
	exprNode()
}

type Variable struct {
	Name Token
}

func (*Variable) exprNode() {}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (*Binary) exprNode() {}

type Literal struct {
	Value any
}

func (*Literal) exprNode() {}

type Grouping struct {
	Expression Expr
}

func (*Grouping) exprNode() {}

type Unary struct {
	Operator Token
	Right    Expr
}

func (*Unary) exprNode() {}

type Assign struct {
	Name  Token
	Value Expr
}

func (*Assign) exprNode() {}

type Logical struct {
	Left 		Expr
	Operator	Token
	Right 		Expr
}

func (*Logical) exprNode() {}
