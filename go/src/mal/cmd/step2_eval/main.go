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

func EVAL(a *ast.AST, evaler *mal.Evaler) ([]types.Valuer, error) {
	return evaler.EvalAST(a)
}

func PRINT(vs []types.Valuer) {
	for _, v := range vs {
		fmt.Println(v.SPrint(true))
	}
}

func REP(line string, evaler *mal.Evaler) error {
	a, err := READ(line)
	if err != nil {
		return err
	}
	vs, err := EVAL(a, evaler)
	if err != nil {
		return err
	}
	PRINT(vs)
	return nil
}

func main() {
	evaler := mal.NewEvaler(mal.NewEnv(nil, nil, nil))

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
		if err := REP(line, evaler); err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		}
	}
}
