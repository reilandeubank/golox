package main

import (
	"bufio"
	"fmt"
	"os"
	//"strings"
	"github.com/reilandeubank/golox/pkg/scanner"
	"github.com/reilandeubank/golox/pkg/expression"
)

//var hadError bool = false

func main() {
	// args := os.Args[1:]

	// if len(args) > 1 {
	// 	fmt.Println("Usage: golox [script]")
	// 	os.Exit(64)
	// } else if len(args) == 1 {
	// 	runFile(args[0])
	// } else {
	// 	runPrompt()
	// }
	expr := expression.Binary{
		expression.Unary{
			scanner.NewToken(scanner.MINUS, "-", nil, 1),
			expression.Literal{"123", scanner.NUMBER}},
		scanner.NewToken(scanner.STAR, "*", nil, 1),
		expression.Grouping{expression.Literal{"45.67", scanner.NUMBER}},
		}
	fmt.Println(expr.String())
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
	bufscanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>>")
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

func run(source string) {
	thisScanner := scanner.NewScanner(source)
	tokens := thisScanner.ScanTokens()

	// For now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
