package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/IXnamI/interpreter_in_go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is an abstract programming language!\n", user.Username)
	fmt.Println("Enter your commands below: ")
	repl.StartEval(os.Stdin, os.Stdout)
}
