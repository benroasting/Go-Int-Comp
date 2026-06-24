package lox

import (
	"fmt"
	"os"
)

type Lox struct {
	HadError        bool
	HadRuntimeError bool
	interpreter     *Interpreter
}

func (l *Lox) Run(source string) {
	scanner := NewScanner(source, l)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens, l)
	stmts := parser.Parse()

	if l.HadError {
		return
	}

	if l.interpreter == nil {
		l.interpreter = NewInterpreter(l)
	}
	l.interpreter.Interpret(stmts)
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) TokenError(t Token, message string) {
	if t.Type == EOF {
		l.report(t.Line, " at end", message)
	} else {
		l.report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func (l *Lox) RuntimeError(err runtimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", err.message, err.token.Line)
	l.HadRuntimeError = true
}

func (l *Lox) report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.HadError = true
}
