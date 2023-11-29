package interpreter

import (
	"time"
	"github.com/reilandeubank/golox/pkg/parser"
)

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, arguments []interface{}) (interface{}, error)
}

type clock struct{}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()) / 1000, nil
}

func (c clock) String() string {
	return "<native fn>"
}

type LoxFunction struct {
	Declaration parser.FunctionStmt
	Closure     *environment
	IsInitializer bool
}

func (l LoxFunction) String() string {
	return "<fn " + l.Declaration.Name.Lexeme + ">"
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Params)
}

func (l LoxFunction) Call(i *Interpreter, arguments []interface{}) (retVal interface{}, errVal error) {
	env := NewEnvironmentWithEnclosing(*l.Closure)

	for j, param := range l.Declaration.Params {
		env.define(param.Lexeme, arguments[j])
	}

	defer func() {
		if r := recover(); r != nil {
			if returnErr, ok := r.(*ReturnError); ok {
                retVal = returnErr.Value
            } else {
                errVal = &RuntimeError{Message: "Return value error"}
            }
		} 
	}()

	_, err := i.executeBlock(l.Declaration.Body, env)
	if err != nil {
		return nil, err
	}

	// if l.IsInitializer {
	// 	return l.Closure.getAt(0, "this")
	// }

	return retVal, errVal
}