package parser

import (
	"github.com/reilandeubank/golox/pkg/scanner"
)

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

type ExprStmt struct {
	Expression Expression
}

func (e ExprStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExprStmt(e)
}

type PrintStmt struct {
	Expression Expression
}

func (p PrintStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expression
}

func (v VarStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitVarStmt(v)
}