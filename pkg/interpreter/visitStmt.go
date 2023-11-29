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

func (i *Interpreter) VisitIfStmt(ifStmt parser.IfStmt) (interface{}, error) {
	condition, err := i.evaluate(ifStmt.Condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(condition) {
		return i.execute(ifStmt.ThenBranch)
	} else if ifStmt.ElseBranch != nil {
		return i.execute(ifStmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(whileStmt parser.WhileStmt) (interface{}, error) {
	for {
		condition, err := i.evaluate(whileStmt.Condition)	// replaces normal for loop condition to allow for error handling
		if err != nil {
			return nil, err
		}

		if !isTruthy(condition) {
			break
		}

		_, err = i.execute(whileStmt.Body)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(functionStmt parser.FunctionStmt) (interface{}, error) {
	function := LoxFunction{Declaration: functionStmt, Closure: i.environment, IsInitializer: false}
	i.environment.define(functionStmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(returnStmt parser.ReturnStmt) (interface{}, error) {
	var value interface{}
	var err error
	if returnStmt.Value != nil {
		value, err = i.evaluate(returnStmt.Value)
		if err != nil {
			return nil, err
		}
	}
	panic(&ReturnError{Value: value})
	//return returnStmt, nil
}