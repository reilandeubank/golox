package interpreter

import (
	"fmt"
	"reflect"

	"github.com/reilandeubank/golox/pkg/parser"
	"github.com/reilandeubank/golox/pkg/scanner"
)

func isTruthy(object interface{}) bool {
	if object == nil || object == 0.0 { // only false, nil, and 0.0 are falsey
		return false
	}

	if reflect.TypeOf(object) == reflect.TypeOf(false) {
		return object.(bool)
	}

	return true
}

func isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil { // No implicit type conversion for equality, like Go
		return true
	} else if a == nil {
		return false
	}
	return a == b
}

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) VisitLiteralExpr(literal parser.Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.Grouping) (interface{}, error) {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) evaluate(expr parser.Expression) (interface{}, error) {
	return expr.Accept(i)
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
	right, _ := i.evaluate(binary.Right)
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
		return nil, &RuntimeError{Token: binary.Operator, Message: "Operands must be two numbers or two strings"}
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

func checkNumberOperand(operator scanner.Token, operand interface{}) error {
	if reflect.TypeOf(operand) == reflect.TypeOf(0.0) {
		return nil
	}
	return &RuntimeError{Token: operator, Message: "Operator must be a number"}
}

func checkNumberOperands(operator scanner.Token, left interface{}, right interface{}) error {
	if reflect.TypeOf(left) == reflect.TypeOf(0.0) && reflect.TypeOf(right) == reflect.TypeOf(0.0) {
		return nil
	}
	return &RuntimeError{Token: operator, Message: "Operators must be numbers"}
}

func (i *Interpreter) Interpret(expression parser.Expression) error {
	value, err := i.evaluate(expression)
	if value != nil {
		println(stringify(value))
	}
	return err
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	// Type assertion for float64
	if val, ok := object.(float64); ok {
		return fmt.Sprintf("%g", val) // %g removes trailing zeros
	}

	// Default to using fmt.Sprintf for other types
	return fmt.Sprintf("%v", object)
}
