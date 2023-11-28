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

func (r *RuntimeError) Error() string {
	msg := fmt.Sprintf("[line %d] Runtime Error: %s", r.Token.Line, r.Message)
	return msg
}
