package main

import (
	"bufio"
	"fmt"
	"os"

	"go-int-comp/lox"
)

func main() {
	l := &lox.Lox{}

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(74)
		}

		l.Run(string(data))
		if l.HadError {
			os.Exit(65)
		}
		if l.HadRuntimeError {
			os.Exit(70)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			l.Run(scanner.Text())
			l.HadError = false
			l.HadRuntimeError = false
		}
	}
}
