package golox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var hadError bool = false

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
	
	if hadError {
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
		hadError = false
	}

	if scanner.Err() != nil {
		fmt.Println("An error occurred:", scanner.Err())
	}
}

func run(source string) {
	reader := strings.NewReader(source)
	scanner := bufio.NewScanner(reader)		//Golang scanners cannot read strings so must pass a string reader
	tokens := scanner.ScanTokens()

	// For now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func loxError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}

