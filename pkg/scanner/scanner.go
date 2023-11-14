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
		Line:   1,
	}
}

func (s *Scanner) isAtEnd() bool {
	length := utf8.RuneCountInString(s.Source)
	return s.Curr < length
}

func (s *Scanner) advance() rune {
	ch := rune(s.Source[s.Curr])
	s.Curr++
	return ch
}

func (s *Scanner) addToken(thisType TokenType) {
	s.addTokenWithTypeAndLiteral(thisType, nil)
}

func (s *Scanner) addTokenWithTypeAndLiteral(thisType TokenType, literal interface{}) {
	text := s.Source[s.Start:s.Curr]
	s.Tokens = append(s.Tokens, Token{Type: thisType, Lexeme: text, Literal: literal, Line: s.Line})
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// we are at the beginning of the next lexeme
		s.Start = s.Curr
		s.ScanToken()
	}

	t := NewToken(EOF, "", nil, s.Line)
	s.Tokens = append(s.Tokens, t)
	return s.Tokens
}

func (s *Scanner) ScanToken() (TokenType, error) {
	ch := s.advance()
	switch ch {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	default:
		LoxError(s.Line, "Unexpected character.")
	}
	return EOF, nil
}
