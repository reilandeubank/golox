package scanner

import (
	"fmt"
	// "log"
	"os"
	//"unicode/utf8"
)

var hadErrorFlag bool = false

func hadError() bool {
	return hadErrorFlag
}

func setErrorFlag(val bool) {
	hadErrorFlag = val
}

func loxError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadErrorFlag = true
}