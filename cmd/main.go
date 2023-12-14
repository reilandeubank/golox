package main

import (
	"bufio"
	"fmt"
	"os"
	//"strings"
	"github.com/reilandeubank/golox/pkg/scanner"
	//"github.com/reilandeubank/golox/pkg/expression"
	"github.com/reilandeubank/golox/pkg/interpreter"
	"github.com/reilandeubank/golox/pkg/parser"
)

var i interpreter.Interpreter = interpreter.NewInterpreter()
func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		err := runFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(64)
		}
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = run(string(bytes))

	if scanner.HadError() {
		os.Exit(65)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(70)
	}
	return nil
}

func runPrompt() {
	bufscanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		if !bufscanner.Scan() {
			break
		}

		line := bufscanner.Text()
		run(line)
		scanner.SetErrorFlag(false)
	}

	if bufscanner.Err() != nil {
		fmt.Println("An error occurred:", bufscanner.Err())
	}
}

func run(source string) error {
	thisScanner := scanner.NewScanner(source)
	tokens := thisScanner.ScanTokens()

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()

	if err != nil {
		return err
	}

	err = i.Interpret(statements)
	if err != nil {
		os.Exit(70)
		return err
	}
	return nil
}
