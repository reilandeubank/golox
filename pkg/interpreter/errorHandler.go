package interpreter

import (
	"fmt"

	"github.com/reilandeubank/golox/pkg/scanner"
)

var hadErrorFlag bool = false

func HadError() bool {
	return hadErrorFlag
}

func SetErrorFlag(val bool) {
	hadErrorFlag = val
}

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func NewRuntimeError(token scanner.Token, message string) RuntimeError {
	r := RuntimeError{Token: token, Message: message}
	r.Message = r.Error()
	return r
}

func (r *RuntimeError) Error() string {
	msg := fmt.Sprintf("[line %d] Error: %s, offending code:%s", r.Token.Line, r.Message, r.Token.Literal)
	return msg
}
