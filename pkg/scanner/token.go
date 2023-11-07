package scanner

import (
	"fmt"
)

// Token struct represents a token with its type, lexeme, literal value, and the line it appears on.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// NewToken is a constructor function for creating a new Token instance.
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

// String method provides a string representation of the Token.
func (t Token) String() string {
	return fmt.Sprintf("%d %s %v", t.Type, t.Lexeme, t.Literal)
}
