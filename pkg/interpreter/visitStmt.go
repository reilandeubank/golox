package interpreter

import (
	"fmt"

	"github.com/reilandeubank/golox/pkg/parser"
	//"github.com/reilandeubank/golox/pkg/scanner"
)

func (i *Interpreter) VisitExprStmt(exprStmt parser.ExprStmt) (interface{}, error) {
	_, err := i.evaluate(exprStmt.Expression)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(printStmt parser.PrintStmt) (interface{}, error) {
	value, err := i.evaluate(printStmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(varStmt parser.VarStmt) (interface{}, error) {
	var value interface{}
	var err error
	if varStmt.Initializer != nil {
		value, err = i.evaluate(varStmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.define(varStmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(blockStmt parser.BlockStmt) (interface{}, error) {
	return i.executeBlock(blockStmt.Statements, NewEnvironmentWithEnclosing(*i.environment))
}