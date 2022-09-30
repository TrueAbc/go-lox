package Token

import (
	"fmt"
)

type Token struct {
	TType   TokenType
	Lexeme  string      // 词位
	Literal interface{} // 字面量
	Line    int         // token 所在的行
}

func (t *Token) String() string {
	return fmt.Sprintf("type is: %v, lexeme value: %s,  Literal is: %v, Line number: %d",
		t.TType, t.Lexeme, t.Literal, t.Line)
}

func NewToken(tType TokenType, lexeme string, literal interface{}, line int) *Token {
	t := &Token{
		TType:   tType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
	return t
}
