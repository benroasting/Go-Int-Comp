package lox

type Environment struct {
	enclosing 	*Environment
	values 		map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func NewChildEnvironment(parent *Environment) *Environment {
	return &Environment{
		enclosing: 	parent,
		values: 	make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
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
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
	}
	panic(runtimeError{
		token:   name,
		message: "Undefined variable '" + name.Lexeme + "'.",
	})
}
