package lox

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}
	panic(runtimeError{
		token:   name,
		message: "Undefined variable '" + name.Lexeme + "'.",
	})
}

func (e *Environment) Assign(name Token, value any) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}
	panic(runtimeError{
		token:   name,
		message: "Undefined variable '" + name.Lexeme + "'.",
	})
}
