package lox

import (
	"fmt"
	"strings"
)

func AstString(e Expr) string {
	switch e := e.(type) {
	case *Binary:
		return parenthesize(e.Operator.Lexeme, e.Left, e.Right)
	case *Unary:
		return parenthesize(e.Operator.Lexeme, e.Right)
	case *Grouping:
		return parenthesize("group", e.Expression)
	case *Literal:
		if e.Value == nil {
			return "nil"
		}
		return fmt.Sprint(e.Value)
	}
	panic(fmt.Sprintf("unknown expr type: %T", e))
}

func parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, e := range exprs {
		sb.WriteString(" ")
		sb.WriteString(AstString(e))
	}
	sb.WriteString(")")
	return sb.String()
}
