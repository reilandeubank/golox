package scanner

import (
	"fmt"
	// "log"
	// "os"
	"strconv"
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

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.Source[s.Curr])
}

func (s *Scanner) addToken(thisType TokenType) {
	//fmt.Println("Adding token: ", thisType)
	s.addTokenWithTypeAndLiteral(thisType, nil)
}

func (s *Scanner) addTokenWithTypeAndLiteral(thisType TokenType, literal interface{}) {
	text := s.Source[s.Start:s.Curr]
	s.Tokens = append(s.Tokens, Token{Type: thisType, Lexeme: text, Literal: literal, Line: s.Line})
}

func (s *Scanner) ScanTokens() []Token {
	// Driving loop
	// Note that s.curr is not incremented in Number, String, or Identifier readers since they handle their own iteration
	for !s.isAtEnd() {
		s.Start = s.Curr
		s.ScanToken()
	}

	// Add EOF token
	s.Tokens = append(s.Tokens, Token{EOF, "EOF", nil, s.Line})
	return s.Tokens
}

func (s *Scanner) ScanToken() {
	ch := s.advance()
	switch ch {
	case '(': s.addToken(LEFT_PAREN)
	case ')': s.addToken(RIGHT_PAREN)
	case '{': s.addToken(LEFT_BRACE)
	case '}': s.addToken(RIGHT_BRACE)
	case ',': s.addToken(COMMA)
	case '.': s.addToken(DOT)
	case '-': s.addToken(MINUS)
	case '+': s.addToken(PLUS)
	case ';': s.addToken(SEMICOLON)
	case '*': s.addToken(STAR)
	// Dual-character tokens
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
				//fmt.Println(s.peek())
			}
		} else {
			s.addToken(SLASH)
		}
	// Handle whitespace
	case ' ': 
	case '\r': 
	case '\t': 
	case '\n':
		s.Line++
	// Handle strings
	case '"': s.tokenizeString()
	default:
		if unicode.IsDigit(rune(ch)) {
			s.tokenizeNumber()
		} else if unicode.IsLetter(rune(ch)) || ch == '_' {
			s.tokenizeIdentifier()
		} else {
			errorStr := fmt.Sprintf("Unexpected character: %c at line %d", ch, s.Line)
			LoxError(s.Line, errorStr)
		}
	}
	
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.Source[s.Curr]) != expected {
		return false
	}
	s.Curr++
	return true;
}

func (s *Scanner) tokenizeString() {
	// Track initial position
	unterminated := true

	//s.Curr++

	// Iterate until end of string or end of file
	for s.Curr < len(s.Source) {

		// Break at closing quote
		if s.Source[s.Curr] == '"' {
			s.Curr++
			unterminated = false
			break
		}

		// Handle newlines
		if s.Source[s.Curr] == '\n' {
			s.Line++
		}

		s.Curr++
	}

	// Check for unterminated string
	if unterminated {
		// Set current position to end of file to prevent further iteration
		s.Curr = len(s.Source)
		errorStr := fmt.Sprintf("Unterminated string at line %d", s.Line)
		LoxError(s.Line, errorStr)
	} else {
		// Return token using substring created from initial and current positions
		s.addTokenWithTypeAndLiteral(STRING, s.Source[s.Start+1 : s.Curr-1])
	}

	// Return token using substring created from initial and current positions
}

// Number reader for Scanner
func (s *Scanner) tokenizeNumber() {
	// Track initial position and whether or not a dot has been found
	foundDot := false

	// Iterate until end of number or end of file
	for s.Curr < len(s.Source) && (unicode.IsDigit(rune(s.Source[s.Curr])) || s.Source[s.Curr] == '.') {

		// Check for dot
		if s.Source[s.Curr] == '.' {
			// Return error if dot has already been found
			if foundDot {
				errorStr := fmt.Sprintf("Invalid number at line %d", s.Line)
				LoxError(s.Line, errorStr)
			}
			// Otherwise, set foundDot to true and skip to next character
			foundDot = true
		}
		// Iterate to next character
		s.Curr++
	}

	floatVal, err := strconv.ParseFloat(s.Source[s.Start:s.Curr], 64)

	if err != nil {
        errorStr := fmt.Sprintf("Invalid number at line %d", s.Line)
		LoxError(s.Line, errorStr)
    }
	// Return token using substring created from initial and current positions
	s.addTokenWithTypeAndLiteral(NUMBER, floatVal)
}

// Identifier reader for Scanner
// Note that although an error is never returned, it is good practice to provide support for it
func (s *Scanner) tokenizeIdentifier() {
	// Iterate until end of identifier or end of file
	for s.Curr < len(s.Source) && (unicode.IsLetter(rune(s.Source[s.Curr])) || unicode.IsDigit(rune(s.Source[s.Curr])) || s.Source[s.Curr] == '_') {
		s.Curr++
	}

	// Check for existing keyword
	identifier := s.Source[s.Start:s.Curr]
	if tokentype, exists := keywords[identifier]; exists {
        s.addToken(tokentype)
    } else {
        // Set to default value if the key is not found
        s.addToken(IDENTIFIER)
    }
}
