package interpreter

import (
	"github.com/reilandeubank/golox/pkg/scanner"
)

type environment struct {
	enclosing *environment
	values map[string]interface{}
}

func NewEnvironment() environment {
	return environment{enclosing: nil, values: make(map[string]interface{})}
}

func NewEnvironmentWithEnclosing(Enclosing environment) environment {
	return environment{enclosing: &Enclosing, values: make(map[string]interface{})}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value	// this allows for variable redefinition. May be weird in normal code, but is useful for REPL
}

func (e *environment) get(name scanner.Token) (interface{}, error) {
	value, ok := e.values[name.Lexeme]
	if !ok && e.enclosing != nil {
		return e.enclosing.get(name)
	} else if !ok {
		return nil, &RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."}
	}
	return value, nil
}

func (e *environment) assign(name scanner.Token, value interface{}) error {
	_, ok := e.values[name.Lexeme]
	if !ok && e.enclosing != nil {
		return e.enclosing.assign(name, value)
	} else if !ok {
		return &RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme}
	}
	e.values[name.Lexeme] = value
	return nil
}