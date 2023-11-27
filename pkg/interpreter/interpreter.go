package interpreter

import (
	//"fmt"
	//"reflect"

	"github.com/reilandeubank/golox/pkg/parser"
	//"github.com/reilandeubank/golox/pkg/scanner"
)

type Interpreter struct{
	//globals *environment
	environment *environment
}

func NewInterpreter() Interpreter {
	env := NewEnvironment()
	return Interpreter{environment: &env}
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
