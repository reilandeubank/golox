package parser

import (
	"fmt"
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
	return p.assignment()
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		return p.block()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expr()
	if err != nil {
		return PrintStmt{Expression: Literal{Value: nil}}, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	return PrintStmt{Expression: value}, err
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return VarStmt{}, err
	}
	var initializer Expression
	if p.match(scanner.EQUAL) {
		initializer, err = p.expr()
		if err != nil {
			return VarStmt{}, err
		}
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	return VarStmt{Name: name, Initializer: initializer}, err
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expr()
	if err != nil {
		return ExprStmt{Expression: Literal{Value: nil}}, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	return ExprStmt{Expression: value}, err
}

func (p *Parser) block() (Stmt, error) {
	var statements []Stmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}
	_, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	return BlockStmt{Statements: statements}, err
}

func (p *Parser) assignment() (Expression, error) {
	expr, err := p.equality()
	if err != nil {
		return Literal{Value: nil}, err
	}
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return Literal{Value: nil}, err
		}
		if variable, ok := expr.(Variable); ok {
			name := variable.Name
			return Assign{Name: name, Value: value}, nil
		}
		message := "Invalid assignment target"
		ParseError(equals, message)
		return Literal{Value: nil}, errors.New(message)
	}
	return expr, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.VAR) {
		declaration, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return VarStmt{}, err
		}
		return declaration, nil
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return ExprStmt{}, err
	}
	return stmt, nil
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
		return Literal{Value: false, Type: scanner.FALSE}, nil
	}
	if p.match(scanner.TRUE) {
		return Literal{Value: true, Type: scanner.TRUE}, nil
	}
	if p.match(scanner.NIL) {
		return Literal{Value: nil, Type: scanner.NIL}, nil
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		var prevValue interface{} = p.previous().Literal
		var err error
		switch prevValue.(type) {
		case string:
			return Literal{Value: prevValue, Type: scanner.STRING}, err
		case float64:
			return Literal{Value: prevValue, Type: scanner.NUMBER}, err
		default:
			// Handle other types or error
			message := "unexpected literal type: " + fmt.Sprintf("%T", prevValue)
			ParseError(p.peek(), message)
			err = errors.New(message)
		}
		return Literal{Value: nil, Type: scanner.NIL}, err
	}
	if p.match(scanner.IDENTIFIER) {
		return Variable{Name: p.previous()}, nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expr()
		if err != nil {
			return Literal{Value: nil}, err
		}
		_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after expression.")
		return Grouping{Expression: expr}, err
	}
	message := "expect expression"
	ParseError(p.peek(), message)
	return Literal{Value: nil}, errors.New(message)
}

func (p *Parser) consume(t scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	ParseError(p.peek(), message)
	return scanner.NewToken(scanner.OTHER, "", nil, 0), errors.New(message)
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
	// expr, err := p.expr()
	// if err != nil {
	// 	return Literal{Value: nil}, err
	// }
	// // _, err = p.consume(scanner.EOF, "Expect end of expression")
	// return expr, err
}