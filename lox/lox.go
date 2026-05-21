package lox

import (
	"fmt"
	"os"
)

type Lox struct {
	HadError bool
}

func (l *Lox) Run(source string){
	scanner := NewScanner(source, l)
	tokens := scanner.ScanTokens()
	for _, t := range tokens { fmt.Println(t) }
}

func (l *Lox) Error(line int, message string){
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string){
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.HadError = true
}
