package parser

import (
	"fmt"
	"strconv"
	"github.com/reilandeubank/golox/pkg/scanner"
	"errors"
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

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.Curr++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.Tokens[p.Curr]
}

func (p *Parser) previous() scanner.Token {
	return p.Tokens[p.Curr-1]
}

func (p *Parser) expr() (Expression, error) {
	return p.equality()
}

func (p *Parser) equality() (Expression, error) {
	expr, err := p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		var right Expression
		right, err= p.comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) comparison() (Expression, error) {
	expr, err := p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		var right Expression
		right, err= p.term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) term() (Expression, error) {
	expr, err := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		var right Expression
		right, err = p.factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) factor() (Expression, error) {
	expr, err := p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		var right Expression
		right, err = p.unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) unary() (Expression, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return Unary{Operator: operator, Right: right}, err
	}
	return p.primary()
}

func (p *Parser) primary() (Expression, error) {
	if p.match(scanner.FALSE) {
		return Literal{Value: "false"}, nil
	}
	if p.match(scanner.TRUE) {
		return Literal{Value: "true"}, nil
	}
	if p.match(scanner.NIL) {
		return Literal{Value: "nil"}, nil
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		var value string
		var err error
		switch v := p.previous().Literal.(type) {
		case string:
			value = v
		case float64:
			value = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			// Handle other types or error
			message := "unexpected literal type: " + fmt.Sprintf("%T", v)
			ParseError(p.peek(), message)
			err = errors.New(message)
		}
		return Literal{Value: value}, err
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expr()
		if err != nil {
			return Literal{Value: "nil"}, err
		}
		_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after expression.")
		return Grouping{Expression: expr}, err
	}
	message := "expect expression"
	ParseError(p.peek(), message)
	return Literal{Value: "nil"}, errors.New(message)
}

func (p *Parser) consume(t scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	ParseError(p.peek(), message)
	return scanner.NewToken(scanner.OTHER, "", nil, 0), errors.New(message)
}

func (p *Parser) Parse() (Expression, error) {
	expr, err := p.expr()
	if err != nil {
		return Literal{Value: "nil"}, err
	}
	//_, err = p.consume(scanner.EOF, "Expect end of expression")
	return expr, err
}