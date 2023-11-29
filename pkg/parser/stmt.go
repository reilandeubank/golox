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

type BlockStmt struct {
	Statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitBlockStmt(b)
}

type IfStmt struct {
	Condition  Expression
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i IfStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expression
	Body      Stmt
}

func (w WhileStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitWhileStmt(w)
}

type FunctionStmt struct {
	Name        scanner.Token
	Params      []scanner.Token
	Body        []Stmt
}

func (f FunctionStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitFunctionStmt(f)
}

type ReturnStmt struct {
	Keyword scanner.Token
	Value   Expression
}

func (r ReturnStmt) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitReturnStmt(r)
}