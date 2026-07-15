package lox

import "time"

type nativeClock struct{}

func (nativeClock) Arity() int { return 0}

func (nativeClock) Call(_ *Interpreter, _ []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}

func (nativeClock) String() string { return "<native fn>"}
