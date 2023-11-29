package interpreter

import (
	// "fmt"
	//"reflect"

	"github.com/reilandeubank/golox/pkg/parser"
	//"github.com/reilandeubank/golox/pkg/scanner"
)

type Interpreter struct{
	globals *environment
	environment *environment
}

func NewInterpreter() Interpreter {
	global := NewEnvironment()
	global.define("clock", &clock{})
	global.define("toStr", &toStr{})
	return Interpreter{environment: &global, globals: &global}
}

func (i *Interpreter) execute(stmt parser.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr parser.Expression) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) executeBlock(statements []parser.Stmt, environment environment) (interface{}, error) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = &environment
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}