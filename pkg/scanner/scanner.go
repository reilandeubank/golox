package scanner

import (
	"fmt"
	// "log"
	// "os"
	"unicode"
	"unicode/utf8"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

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
	return s.Curr >= length
}

func (s *Scanner) advance() rune {
	ch := rune(s.Source[s.Curr])
	s.Curr++
	return ch
}

func (s *Scanner) ScanTokens() []Token {
	// Driving loop
	// Note that s.curr is not incremented in Number, String, or Identifier readers since they handle their own iteration
	for !s.isAtEnd() {
		nextToken, err := s.ScanToken()

		// Handle string literals
		if nextToken == STRING {
			token, err := s.tokenizeString()
			if err == nil {
				s.Tokens = append(s.Tokens, token)
			}

			// Handle number literals
		} else if nextToken == NUMBER {
			token, err := s.tokenizeNumber()
			if err == nil {
				s.Tokens = append(s.Tokens, token)
			}

			// Handle Identifiers
		} else if nextToken == IDENTIFIER {
			token, err := s.tokenizeIdentifier()
			if err == nil {
				s.Tokens = append(s.Tokens, token)
			}

			// Handle Whitespace
		} else if nextToken == WHITESPACE {
			// Whitespace doesn't constitute a token, so just increment
			s.Curr++

			// Handle all other tokens
		} else {
			if err == nil {
				s.Tokens = append(s.Tokens, Token{nextToken, string(s.Source[s.Curr]), nil, s.Line})
			}
			s.Curr++
		}

		if err != nil {
			LoxError(s.Line, "Error scanning token.")
		}
	}

	// Add EOF token
	s.Tokens = append(s.Tokens, Token{EOF, "EOF", nil, s.Line})
	return s.Tokens
}

func (s *Scanner) ScanToken() (TokenType, error) {
	ch := s.advance()
	switch ch {
	case '(':
		return LEFT_PAREN, nil
	case ')':
		return RIGHT_PAREN, nil
	case '{':
		return LEFT_BRACE, nil
	case '}':
		return RIGHT_BRACE, nil
	case ',':
		return COMMA, nil
	case '.':
		return DOT, nil
	case '-':
		return MINUS, nil
	case '+':
		return PLUS, nil
	case ';':
		return SEMICOLON, nil
	case '*':
		return STAR, nil
		// Dual-character tokens
	case '!':
		if s.Source[s.Curr+1] == '=' {
			s.Curr++
			return BANG_EQUAL, nil
		} else {
			return BANG, nil
		}
	case '=':
		if s.Source[s.Curr+1] == '=' {
			s.Curr++
			return EQUAL_EQUAL, nil
		} else {
			return EQUAL, nil
		}
	case '<':
		if s.Source[s.Curr+1] == '=' {
			s.Curr++
			return LESS_EQUAL, nil
		} else {
			return LESS, nil
		}
	case '>':
		if s.Source[s.Curr+1] == '=' {
			s.Curr++
			return GREATER_EQUAL, nil
		} else {
			return GREATER, nil
		}
	case '/':
		if s.Source[s.Curr+1] == '/' {
			for s.Source[s.Curr] != '\n' && s.Curr < len(s.Source) {
				s.Curr++
			}
		} else {
			return SLASH, nil
		}
	// Handle whitespace
	case ' ':
		return WHITESPACE, nil
	case '\r':
		return WHITESPACE, nil
	case '\t':
		return WHITESPACE, nil
	case '\n':
		s.Line++
		return WHITESPACE, nil
	// Handle strings
	case '"':
		return STRING, nil
	}
	// Handle numbers
	if unicode.IsDigit(rune(ch)) {
		return NUMBER, nil
	} else if unicode.IsLetter(rune(ch)) {
		return IDENTIFIER, nil
	}
	LoxError(s.Line, "Unexpected character.")
	return OTHER, fmt.Errorf("Unexpected character: %c at line %d", ch, s.Line)
}

func (s *Scanner) tokenizeString() (Token, error) {
	// Track initial position
	initial := s.Curr
	foundQuote := false

	s.Curr++

	// Iterate until end of string or end of file
	for s.Curr < len(s.Source) {

		// Break at closing quote
		if s.Source[s.Curr] == '"' {
			s.Curr++
			foundQuote = true
			break
		}

		// Handle newlines
		if s.Source[s.Curr] == '\n' {
			s.Line++
		}

		s.Curr++
	}

	// Check for unterminated string
	if !foundQuote {
		// Set current position to end of file to prevent further iteration
		s.Curr = len(s.Source)
		return Token{}, fmt.Errorf("Unterminated string at line %d", s.Line)
	}

	// Return token using substring created from initial and current positions
	return Token{STRING, s.Source[initial:s.Curr], s.Source[initial+1 : s.Curr-1], s.Line}, nil
}

// Number reader for Scanner
func (s *Scanner) tokenizeNumber() (Token, error) {
	// Track initial position and whether or not a dot has been found
	initial := s.Curr
	foundDot := false

	// Iterate until end of number or end of file
	for s.Curr < len(s.Source) && (unicode.IsDigit(rune(s.Source[s.Curr])) || s.Source[s.Curr] == '.') {

		// Check for dot
		if s.Source[s.Curr] == '.' {
			// Return error if dot has already been found
			if foundDot {
				return Token{}, fmt.Errorf("Invalid number at line %d", s.Line)
			}
			// Otherwise, set foundDot to true and skip to next character
			foundDot = true
		}
		// Iterate to next character
		s.Curr++
	}

	// Return token using substring created from initial and current positions
	return Token{NUMBER, s.Source[initial:s.Curr], s.Source[initial:s.Curr], s.Line}, nil
}

// Identifier reader for Scanner
// Note that although an error is never returned, it is good practice to provide support for it
func (s *Scanner) tokenizeIdentifier() (Token, error) {
	// Track initial position
	initial := s.Curr

	// Iterate until end of identifier or end of file
	for s.Curr < len(s.Source) && (unicode.IsLetter(rune(s.Source[s.Curr])) || unicode.IsDigit(rune(s.Source[s.Curr]))) {
		s.Curr++
	}

	// Check for existing keyword
	identifier := s.Source[initial:s.Curr]
	if tokentype, found := keywords[identifier]; found {
		return Token{tokentype, identifier, nil, s.Line}, nil
	}

	// Return token using substring created from initial and current positions
	return Token{IDENTIFIER, identifier, nil, s.Line}, nil
}
