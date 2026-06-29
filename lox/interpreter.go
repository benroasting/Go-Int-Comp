package lox

import "fmt"

type Interpreter struct {
	lox         *Lox
	environment *Environment
}

type runtimeError struct {
	token   Token
	message string
}

func NewInterpreter(l *Lox) *Interpreter {
	return &Interpreter{
		lox:         l,
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(runtimeError); ok {
				i.lox.RuntimeError(rerr)
			} else {
				panic(r)
			}
		}
	}()
	for _, s := range stmts {
		i.execute(s)
	}
}

func (i *Interpreter) execute(s Stmt) {
	switch s := s.(type) {
	case *ExpressionStmt:
		i.evaluate(s.Expression)
	case *PrintStmt:
		value := i.evaluate(s.Expression)
		fmt.Println(stringify(value))
	case *VarStmt:
		var value any
		if s.Initializer != nil {
			value = i.evaluate(s.Initializer)
		}
		i.environment.Define(s.Name.Lexeme, value)
	case *BlockStmt:
		i.executeBlock(s.Statements, NewChildEnvironment(i.environment))
		
	}
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	previous := i.environment
	defer func() { i.environment = previous }()
	i.environment = env
	for _, s := range stmts {
		i.execute(s)
	}
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
				token:   e.Operator,
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
	case *Variable:
		return i.environment.Get(e.Name)
	case *Assign:
		value := i.evaluate(e.Value)
		i.environment.Assign(e.Name, value)
		return value
	}
	return nil
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

func checkNumberOperands(operator Token, left, right any) {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}
	panic(runtimeError{token: operator, message: "Operands must be numbers."})
}
