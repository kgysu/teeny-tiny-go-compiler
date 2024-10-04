package main

import (
	"fmt"
	"os"
	"tiny-compiler/pkg/comp"
)

func main() {
	fmt.Println("Teeny tiny compiler!")

	args := os.Args
	if len(args) < 2 {
		panic("Error: Compiler needs source file as argument.")
	}

	source, err := os.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	lexer := comp.NewLexer(string(source))
	emitter := comp.NewEmitter("out.c")
	parser := comp.NewParser(lexer, emitter)

	parser.Program()
	emitter.WriteFile()
	fmt.Println("Compiling completed.")
}
