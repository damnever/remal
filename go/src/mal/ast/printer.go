package ast

import "fmt"

func PrintAST(ast *AST) {
	ast.Walk(func(node Node) bool {
		if _, ok := node.(*Comment); ok {
			return true
		}
		fmt.Printf("%s\n", node)
		return true
	})
}
