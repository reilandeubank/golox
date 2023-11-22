package expression

import (
	"github.com/reilandeubank/golox/pkg/scanner"
)

// Expression interface

// Interface in go is similar to an abstract class in Java
// Expression is an interface that all expressions will implement
type Expression interface {
	Accept() string
	String() string
}

// Parenthesize() is a method that returns a string representation of the expression
func Parenthesize(name string, exprs ...Expression) string {

	// Open the string
	output := "(" + name

	// Iterate over expressions and add them to the string
	for _, expr := range exprs {
		output += " " + expr.Accept()
	}

	// Close the string
	output += ")"
	return output
}

// Expression implementations

// Literal

// Literal is a struct that implements the Expression interface
type Literal struct {
	Value string
	Type  scanner.TokenType
}

// Accept() is a method that returns a string representation of the expression
func (l Literal) Accept() string {
	return l.Value
}

// visitLiteralExpr() is a method that returns a string representation of the expression
func visitLiteralExpr(l Literal) string {
	return l.Value
}

// String() is a method that returns a string representation of the expression
func (l Literal) String() string {
	return l.Value
}

// Grouping

// Grouping is a struct that implements the Expression interface
type Grouping struct {
	Expression Expression
}

// Accept() is a method that returns a string representation of the expression
func (g Grouping) Accept() string {
	return ("(" + g.Expression.Accept() + ")")
}

// visitGroupingExpr() is a method that returns a string representation of the expression
func visitGroupingExpr(expr Grouping) string {
	return Parenthesize("group", expr.Expression)
}

// String() is a method that returns a string representation of the expression
func (g Grouping) String() string {
	return ("(" + g.Expression.Accept() + ")")
}

// Unary

// Unary is a struct that implements the Expression interface
type Unary struct {
	Operator scanner.Token
	Right    Expression
}

// Accept() is a method that returns a string representation of the expression
func (u Unary) Accept() string {
	return "(" + u.Operator.Lexeme + u.Right.Accept() + ")"
}

// visitUnaryExpr() is a method that returns a string representation of the expression
func visitUnaryExpr(expr Unary) string {
	return Parenthesize(expr.Operator.Lexeme, expr.Right)
}

// String() is a method that returns a string representation of the expression
func (u Unary) String() string {
	return "(" + u.Operator.Lexeme + u.Right.Accept() + ")"
}

// Binary

// Binary is a struct that implements the Expression interface
type Binary struct {
	Left     Expression
	Operator scanner.Token
	Right    Expression
}

// Accept() is a method that returns a string representation of the expression
func (b Binary) Accept() string {
	return b.Left.Accept() + b.Operator.Lexeme + b.Right.Accept()
}

// visitBinaryExpr() is a method that returns a string representation of the expression
func visitBinaryExpr(expr Binary) string {
	return Parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// String() is a method that returns a string representation of the expression
func (b Binary) String() string {
	return b.Left.Accept() + b.Operator.Lexeme + b.Right.Accept()
}

type Logical struct {
	Left     Expression
	Operator scanner.Token
	Right    Expression
}

func (l Logical) Accept() string {
	return l.Left.Accept() + l.Operator.Lexeme + l.Right.Accept()
}

func visitLogicalExpr(expr Logical) Logical {

	left := expr.Left.Accept()

	if expr.Operator.Type == scanner.OR {
		operatorToken := scanner.NewToken(scanner.OR, "or", nil, expr.Operator.Line)
		if left == "true" {
			return Logical{
				Left:     Literal{"true", scanner.TRUE},
				Operator: operatorToken,
				Right:    Literal{"true", scanner.TRUE},	// lazy evaluation
			}
		} else {
			return Logical{
				Left:     Literal{"false", scanner.FALSE},
				Operator: operatorToken,
				Right:    expr.Right,			// only evaluate the right side if the left side is false
			}
		}
	} else {
		operatorToken := scanner.NewToken(scanner.AND, "and", nil, expr.Operator.Line)
		if left == "true" {
			return Logical{
				Left:     Literal{"true", scanner.TRUE},
				Operator: operatorToken,
				Right:    expr.Right,		// only evaluate the right side if the left side is true	
			}
		} else {
			return Logical{
				Left:     Literal{"false", scanner.FALSE},
				Operator: operatorToken,
				Right:    Literal{"false", scanner.FALSE},	// lazy evaluation
			}
		}
	}
}

func (l Logical) String() string {
	return l.Left.Accept() + l.Operator.Lexeme + l.Right.Accept()
}

type Assign struct {
	Name  scanner.Token
	Value Expression
}

func (a Assign) Accept() string {
	return a.Name.Lexeme + " = " + a.Value.Accept()
}

func (a Assign) String() string {
	return a.Name.Lexeme + " = " + a.Value.Accept()
}

// Variable

type Variable struct {
	Name scanner.Token
}

func (v Variable) Accept() string {
	return v.Name.Lexeme
}

func (v Variable) String() string {
	return v.Name.Lexeme
}

// Call

type Call struct {
	Callee Expression
	Paren  scanner.Token
	Args   []Expression
}

func (c Call) Accept() string {
	return c.Callee.Accept()
}

func (c Call) String() string {
	return c.Callee.Accept()
}

// Note: accept(), visit(), and String() are all very similar but serve different functions within the interpreter
// accept() is used to generate the AST
// visit() is used to generate the string representation of the AST
// String() is used to generate readable output for debugging/users