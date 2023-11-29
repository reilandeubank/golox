package parser

import (
	//"fmt"
	"github.com/reilandeubank/golox/pkg/scanner"
)

// Expression interface

// Interface in go is similar to an abstract class in Java
// Expression is an interface that all expressions will implement
type Expression interface {
	Accept(v ExprVisitor) (interface{}, error)
	// String() string
}

// Literal

// Literal is a struct that implements the Expression interface
type Literal struct {
	Value interface{}
	Type  scanner.TokenType
}

// Accept() is a method that returns a string representation of the expression
func (l Literal) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLiteralExpr(l)
}

// Grouping

// Grouping is a struct that implements the Expression interface
type Grouping struct {
	Expression Expression
}

// Accept() is a method that returns a string representation of the expression
func (g Grouping) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGroupingExpr(g)
}

// Unary

// Unary is a struct that implements the Expression interface
type Unary struct {
	Operator scanner.Token
	Right    Expression
}

// Accept() is a method that returns a string representation of the expression
func (u Unary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitUnaryExpr(u)
}
// Binary

// Binary is a struct that implements the Expression interface
type Binary struct {
	Left     Expression
	Operator scanner.Token
	Right    Expression
}

// Accept() is a method that returns a string representation of the expression
func (b Binary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitBinaryExpr(b)
}

// Variable

// Variable is a struct that implements the Expression interface
type Variable struct {
	Name scanner.Token
}

// Accept() is a method that returns a string representation of the expression
func (v Variable) Accept(vi ExprVisitor) (interface{}, error) {
	return vi.VisitVariableExpr(v)
}

// Assignment

// Assignment is a struct that implements the Expression interface
type Assign struct {
	Name  scanner.Token
	Value Expression
}

// Accept() is a method that returns a string representation of the expression
func (a Assign) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitAssignExpr(a)
}

// Logical

// Logical is a struct that implements the Expression interface
type Logical struct {
	Left     Expression
	Operator scanner.Token
	Right    Expression
}

// Accept() is a method that returns a string representation of the expression
func (l Logical) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLogicalExpr(l)
}

// Call

// Call is a struct that implements the Expression interface
type Call struct {
	Callee    Expression
	Paren     scanner.Token
	Arguments []Expression
}

// Accept() is a method that returns a string representation of the expression
func (c Call) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitCallExpr(c)
}