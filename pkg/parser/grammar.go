package parser

import (
	"fmt"
	"github.com/reilandeubank/golox/pkg/scanner"
	"errors"
)

func (p *Parser) expr() (Expression, error) {
	return p.assignment()
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.RETURN) {
		return p.returnStatement()
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return BlockStmt{}, err
		}
		return BlockStmt{Statements: statements}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return WhileStmt{}, err
	}
	var initializer Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return WhileStmt{}, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return WhileStmt{}, err
		}
	
	}
	var condition Expression
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expr()
		if err != nil {
			return WhileStmt{}, err
		}
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return WhileStmt{}, err
	}
	var increment Expression
	if !p.check(scanner.RIGHT_PAREN) {
		increment, err = p.expr()
		if err != nil {
			return WhileStmt{}, err
		}
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return WhileStmt{}, err
	}
	body, err := p.statement()
	if err != nil {
		return WhileStmt{}, err
	}
	if increment != nil {
		body = BlockStmt{Statements: []Stmt{body, ExprStmt{Expression: increment}}}
	}
	if condition == nil {
		condition = Literal{Value: true}
	}
	body = WhileStmt{Condition: condition, Body: body}
	if initializer != nil {
		body = BlockStmt{Statements: []Stmt{initializer, body}}
	}
	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return IfStmt{}, err
	}
	condition, err := p.expr()
	if err != nil {
		return IfStmt{}, err
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return IfStmt{}, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return IfStmt{}, err
	}
	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return IfStmt{}, err
		}
	}
	return IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
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

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expression
	if !p.check(scanner.SEMICOLON) {
		var err error
		value, err = p.expr()
		if err != nil {
			return ReturnStmt{}, err
		}
	}
	_, err := p.consume(scanner.SEMICOLON, "Expect ';' after return value.")
	return ReturnStmt{Keyword: keyword, Value: value}, err
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return WhileStmt{}, err
	}
	condition, err := p.expr()
	if err != nil {
		return WhileStmt{}, err
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return WhileStmt{}, err
	}
	body, err := p.statement()
	if err != nil {
		return WhileStmt{}, err
	}
	return WhileStmt{Condition: condition, Body: body}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expr()
	if err != nil {
		return ExprStmt{Expression: Literal{Value: nil}}, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	return ExprStmt{Expression: value}, err
}

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return FunctionStmt{}, err
	}
	_, err = p.consume(scanner.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return FunctionStmt{}, err
	}
	var parameters []scanner.Token
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				message := "Cannot have more than 255 parameters."
				ParseError(p.peek(), message)
			}
			param, err := p.consume(scanner.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return FunctionStmt{}, err
			}
			parameters = append(parameters, param)
			if !p.match(scanner.COMMA) {
				break
			}
		}
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return FunctionStmt{}, err
	}
	_, err = p.consume(scanner.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return FunctionStmt{}, err
	}
	body, err := p.block()
	if err != nil {
		return FunctionStmt{}, err
	}
	return FunctionStmt{Name: name, Params: parameters, Body: body}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}
	_, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	return statements, err
}

func (p *Parser) assignment() (Expression, error) {
	expr, err := p.or()
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

func (p *Parser) or() (Expression, error) {
	expr, err := p.and()
	for p.match(scanner.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return Literal{Value: nil}, err
		}
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) and() (Expression, error) {
	expr, err := p.equality()
	for p.match(scanner.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return Literal{Value: nil}, err
		}
		expr = Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.FUN) {
		return p.function("function")
	}
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
	return p.call()
}

func (p *Parser) finishCall(callee Expression) (Expression, error) {
	var arguments []Expression
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				message := "Cannot have more than 255 arguments."
				ParseError(p.peek(), message)
			}
			argument, err := p.expr()
			if err != nil {
				return Literal{Value: nil}, err
			}
			arguments = append(arguments, argument)
			if !p.match(scanner.COMMA) {
				break
			}
		}
	}
	paren, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments.")
	return Call{Callee: callee, Paren: paren, Arguments: arguments}, err
}

func (p *Parser) call() (Expression, error) {
	expr, err := p.primary()
	for {
		if p.match(scanner.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr, err
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