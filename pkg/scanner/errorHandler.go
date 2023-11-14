package scanner

import (
	"fmt"
	// "log"
	"os"
	//"unicode/utf8"
)

var hadErrorFlag bool = false

func HadError() bool {
	return hadErrorFlag
}

func SetErrorFlag(val bool) {
	hadErrorFlag = val
}

func LoxError(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadErrorFlag = true
}