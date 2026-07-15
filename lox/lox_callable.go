package lox

type LoxCallable interface {
		Arity()	int //the number of arguments a function takes
		Call(interpreter *Interpreter, arguments []any) any 
}
