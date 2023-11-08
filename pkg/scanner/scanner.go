package scanner

import (
	// "fmt"
	// "log"
	// "os"
	"unicode/utf8"
)

type Scanner struct {
	Source string
	Tokens []Token
	Start  int
	Curr   int
	Line   int
}

func NewScanner(sourceText string) Scanner {
	var tokens []Token
	return Scanner{
		Source: sourceText,
		Tokens: tokens,
		Start:  0,
		Curr:   0,
		Line:   0,
	}
}

func (s *Scanner) isAtEnd() bool {
	length := utf8.RuneCountInString(s.Source)
	return s.Curr < length
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		// we are at the beginning of the next lexeme
		s.Start = s.Curr
		s.scanToken()
	}

	t := NewToken(EOF, "", nil, s.Line)
	s.Tokens = append(s.Tokens, t)
	return s.Tokens
}

func (s *Scanner) scanToken() (TokenType, error) {
	return EOF, nil
}
