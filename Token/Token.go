package Token

import (
	"fmt"
)

type Token struct {
	tType   TokenType
	Lexeme  string      // 词位
	literal interface{} // 字面量
	line    int         // token 所在的行
}

func (t *Token) String() string {
	return fmt.Sprintf("type is: %v, lexeme value: %s,  literal is: %v, line number: %d",
		t.tType, t.Lexeme, t.literal, t.line)
}

func NewToken(tType TokenType, lexeme string, literal interface{}, line int) *Token {
	t := &Token{
		tType:   tType,
		Lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
	return t
}
