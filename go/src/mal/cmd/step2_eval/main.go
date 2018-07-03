package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"mal"
	"mal/ast"
	"mal/types"
)

func READ(line string) (*ast.AST, error) {
	a := new(ast.AST)
	err := a.Parse(line)
	return a, err
}

func EVAL(a *ast.AST, evaler *mal.Evaler) []types.Valuer {
	return evaler.EvalAST(a)
}

func PRINT(vs []types.Valuer) {
	for _, v := range vs {
		fmt.Printf("%s\n", v)
	}
}

func REP(line string) error {
	a, err := READ(line)
	if err != nil {
		return err
	}
	PRINT(EVAL(a, new(mal.Evaler)))
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
