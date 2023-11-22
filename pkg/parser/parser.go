package parser

import (
	//"fmt"
	"github.com/reilandeubank/golox/pkg/scanner"
)

type Parser struct {
	Tokens []scanner.Token
	Curr int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{
		Tokens: tokens,
		Curr: 0,
	}
}