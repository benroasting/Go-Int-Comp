package lox

import "testing"

func TestAstString(t *testing.T) {
	expr := &Binary{
		Left: &Unary{
			Operator: Token{Type: Minus, Lexeme: "-", Line: 1},
			Right:    &Literal{Value: 123.0},
		},
		Operator: Token{Type: Star, Lexeme: "*", Line: 1},
		Right: &Grouping{
			Expression: &Literal{Value: 45.67},
		},
	}

	got := AstString(expr)
	want := "(* (- 123) (group 45.67))"
	if got != want {
		t.Errorf("AstString = %q, want %q", got, want)
	}
}
