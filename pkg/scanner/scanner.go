package scanner

// "fmt"
// "log"
// "os"
// "unicode"

type Scanner struct {
	source string
	tokens []Token
	curr   int
	Line   int
}

func (s *Scanner) scanTokens() {
}
