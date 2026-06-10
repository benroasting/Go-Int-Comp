package lox

import "fmt"

type Interpreter struct {
	lox *Lox
}

type runtimeError struct {
	token 	Token
	message string
}

func NewInterpreter(l *Lox) *Interpreter {
	return &Interpreter{lox: l}
}

func (i *Interpreter) Interpret(e Expr) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(runtimeError); ok {
				i.lox.RuntimeError(rerr)
			} else {
				panic(r)
			}
		}
	}()
	value := i.evaluate(e)
	fmt.Println(stringify(value))
}

func (i *Interpreter) evaluate(e Expr) any {
	switch e := e.(type) {
	case *Literal:
		return e.Value
	case *Grouping:
		return i.evaluate(e.Expression)
	case *Unary:
		right := i.evaluate(e.Right)
		switch e.Operator.Type {
		case Minus:
			checkNumberOperand(e.Operator, right)
			return -right.(float64)
		case Bang:
			return !isTruthy(right)
		}
		panic("unreachable")
	case *Binary:
		left := i.evaluate(e.Left)
		right := i.evaluate(e.Right)

		switch e.Operator.Type {
		case Minus:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) - right.(float64)
		case Slash:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) / right.(float64)
		case Star:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) * right.(float64)
		case Plus:
			if l, ok := left.(float64); ok {
				if r, ok := right.(float64); ok {
					return l + r
				}
			}
			if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r
				}
			}
			panic(runtimeError{
				token: 	e.Operator,
				message: "Operands must be two numbers or two strings.",
			})
		case Greater:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) > right.(float64)
		case GreaterEqual:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) >= right.(float64)
		case Less: 
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) < right.(float64)
		case LessEqual:
			checkNumberOperands(e.Operator, left, right)
			return left.(float64) <= right.(float64)
		case BangEqual:
			return !isEqual(left, right)
		case EqualEqual:
        	return isEqual(left, right)
		}
		panic("unreachable")
	}
	return nil // placeholder - add Unary and Binary later
}

func stringify(v any) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprint(v)
}

func isEqual(a, b any) bool {
    return a == b
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func checkNumberOperand(operator Token, operand any) {
	 if _, ok := operand.(float64); ok {
		return
	 }
	 panic(runtimeError{token: operator, message: "Operand must be a number."})
}

func checkNumberOperands(operator Token, left, right any){
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}
	panic(runtimeError{token: operator, message: "Operands must be numbers."})
}
