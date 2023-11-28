package parser

import (
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

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt

	for !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return []Stmt{}, err
		}
		statements = append(statements, dec)
	}

	return statements, nil
}