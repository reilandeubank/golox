package interpreter

import (
	"reflect"
	"fmt"
	"github.com/reilandeubank/golox/pkg/parser"
	"github.com/reilandeubank/golox/pkg/scanner"
)

func (i *Interpreter) VisitLiteralExpr(literal parser.Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.Grouping) (interface{}, error) {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) VisitUnaryExpr(unary parser.Unary) (interface{}, error) {
	right, err := i.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Type {
	case scanner.MINUS:
		err = checkNumberOperand(unary.Operator, right)
		if err != nil {
			return nil, err
		}
		return -(right.(float64)), nil
	case scanner.BANG:
		return !isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(binary parser.Binary) (interface{}, error) {
	left, err := i.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.Type {
	case scanner.MINUS:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case scanner.SLASH:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case scanner.STAR:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case scanner.PLUS:
		if reflect.TypeOf(left) == reflect.TypeOf("") && reflect.TypeOf(right) == reflect.TypeOf("") {
			return left.(string) + right.(string), nil
		}
		if reflect.TypeOf(left) == reflect.TypeOf(0.0) && reflect.TypeOf(right) == reflect.TypeOf(0.0) {
			return left.(float64) + right.(float64), nil
		}
		return nil, err
	case scanner.GREATER:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case scanner.GREATER_EQUAL:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case scanner.LESS:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case scanner.LESS_EQUAL:
		err = checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case scanner.BANG_EQUAL:
		return !isEqual(left, right), nil
	case scanner.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil // unreachable
}

func (i *Interpreter) VisitVariableExpr(variable parser.Variable) (interface{}, error) {
	return i.environment.get(variable.Name)
}

func (i *Interpreter) VisitAssignExpr(expr parser.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	i.environment.assign(expr.Name, value)
	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr parser.Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == scanner.OR {
		if isTruthy(left) { // short-circuiting
			return left, nil
		}
	} else {
		if !isTruthy(left) { // short-circuiting
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr parser.Call) (interface{}, error) {
	callee, err := expr.Callee.Accept(i)
	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, len(expr.Arguments))
	for j, argument := range expr.Arguments {
		arguments[j], err = i.evaluate(argument)
		if err != nil {
			return nil, err
		}
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, &RuntimeError{Token: expr.Paren, Message: "Can only call functions."}
	}

	if len(arguments) != function.Arity() {
		return nil, &RuntimeError{Token: expr.Paren, Message: "Expected " + fmt.Sprint(function.Arity()) + " arguments but got " + fmt.Sprint(len(arguments)) + "."}
	}

	return function.Call(i, arguments)
}