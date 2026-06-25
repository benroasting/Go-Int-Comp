package lox

type Stmt interface {
	stmtNode()
}

type ExpressionStmt struct {
	Expression Expr
}

func (*ExpressionStmt) stmtNode() {}

type PrintStmt struct {
	Expression Expr
}

func (*PrintStmt) stmtNode() {}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (*VarStmt) stmtNode() {}

type BlockStmt struct {
	Statements 	[]Stmt
}

func (*BlockStmt) stmtNode() {}
