package lox

type Expr interface {
	exprNode()
}

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
