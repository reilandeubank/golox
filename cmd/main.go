package main

import (
	"bufio"
	"fmt"
	"os"
	//"strings"

	"github.com/reilandeubank/golox/pkg/scanner"
)

//var hadError bool = false

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	run(string(bytes))

	if scanner.HadError() {
		os.Exit(65)
	}
	return nil
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>>")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		run(line)
		scanner.SetErrorFlag(false)
	}

	if scanner.Err() != nil {
		fmt.Println("An error occurred:", scanner.Err())
	}
}

func run(source string) {
	thisScanner := scanner.NewScanner(source)
	tokens := thisScanner.ScanTokens()

	// For now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}

	foo := scanner.Tester() //imports are working

	fmt.Println(foo)
}
