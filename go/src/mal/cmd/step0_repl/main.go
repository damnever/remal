package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("user> ")
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "ERR: %v", err)
			}
			break
		}
		fmt.Print(line)
	}
}
