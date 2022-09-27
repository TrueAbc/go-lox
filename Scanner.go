package main

import "strconv"

var KEY_WORDS = map[string]TokenType{
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
	source  string
	tokens  []*Token
	start   int
	current int
	line    int
}

// NewScanner 读取source分割为token
func NewScanner(source string) *Scanner {
	s := &Scanner{source: source, tokens: make([]*Token, 0)}
	s.start = 0
	s.current = 0
	s.line = 1
	return s
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	n := s.advance()
	switch n {
	case '(':
		s.addTokenDefault(LEFT_PAREN)
	case ')':
		s.addTokenDefault(RIGHT_PAREN)
	case '{':
		s.addTokenDefault(LEFT_BRACE)
	case '}':
		s.addTokenDefault(RIGHT_BRACE)
	case ',':
		s.addTokenDefault(COMMA)
	case '.':
		s.addTokenDefault(DOT)
	case '-':
		s.addTokenDefault(MINUS)
	case '+':
		s.addTokenDefault(PLUS)
	case ';':
		s.addTokenDefault(SEMICOLON)
	case '*':
		s.addTokenDefault(STAR)

	// 两个阶段的关键字
	case '!':
		if s.match('=') {
			s.addTokenDefault(BANG_EQUAL)
		} else {
			s.addTokenDefault(BANG)
		}
	case '=':
		if s.match('=') {
			s.addTokenDefault(EQUAL_EQUAL)
		} else {
			s.addTokenDefault(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addTokenDefault(LESS_EQUAL)
		} else {
			s.addTokenDefault(LESS)
		}
	case '>':
		if s.match('=') {
			s.addTokenDefault(GREATER_EQUAL)
		} else {
			s.addTokenDefault(GREATER)
		}
	case '/':
		// 可能代表注释
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addTokenDefault(SLASH)
		}
	//newline and whitespace
	case ' ', '\r', '\t':
		break
		// Ignore whitespace.
	case '\n':
		s.line++

	case '"':
		// string 字面量匹配
		s.string()

	// 关键字和标识符, or orchid 前者是关键字, 后者不是
	// 正则概念, 是否进行贪婪匹配
	// 默认都是标识符处理, 再看是否是关键字
	default:
		if s.isDigit(n) {
			s.number()
		} else if s.isAlpha(n) {
			s.identifier()
		} else {
			Errorf("%d Unexpected character.", s.line)
		}
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		// cpp的字符串结束 '\0', ascii 为0
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current += 1
	return true
}

// 访问单个字符, 增加current下标
func (s *Scanner) advance() byte {
	defer func() {
		s.current += 1
	}()
	return s.source[s.current]
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]

	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) addTokenDefault(tokenType TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		// 这种模式可以支持多行的string, 但是单行string无法添加分隔符
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		Errorf("%d unterminated string", s.line)
		return
	}

	// 消费最后一个 '"'
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	literal, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(NUMBER, literal)

}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	t, ok := KEY_WORDS[text]
	if !ok {
		t = IDENTIFIER
	}
	s.addTokenDefault(t)
}

func (s *Scanner) isDigit(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func (s *Scanner) isAlpha(c byte) bool {
	if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' {
		return true
	}
	return false
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
