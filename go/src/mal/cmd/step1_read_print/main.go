package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"mal/ast"
)

func READ(line string) (*ast.AST, error) {
	a := new(ast.AST)
	err := a.Parse(line)
	return a, err
}

func EVAL(a *ast.AST, envs map[string]string) *ast.AST {
	return a
}

func PRINT(a *ast.AST) {
	ast.PrintAST(a)
}

func REP(line string) error {
	a, err := READ(line)
	if err != nil {
		return err
	}
	PRINT(EVAL(a, nil))
	return nil
}

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("user> ")
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
			return
		}
		if err := REP(line); err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		}
	}
}
