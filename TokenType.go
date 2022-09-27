package main

// 字符的类型, 词法解析的时候识别的类型

type TokenType int

const (
	LEFT_PAREN  = iota + 1 // (
	RIGHT_PAREN            // )
	LEFT_BRACE             // {
	RIGHT_BRACE            // }
	COMMA                  // ,
	DOT                    // .
	MINUS                  // -
	PLUS                   // +
	SEMICOLON              // ;
	SLASH                  // 反斜线
	STAR                   // *

	// todo 部分关键字含义不清楚

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	/*
	  literals 字面量
	*/
	IDENTIFIER
	STRING
	NUMBER

	/*
		keywords for lox language
	*/

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

// 单纯用于打印
var TokenTypeMap = map[TokenType]string{
	LEFT_PAREN:  "(", // (
	RIGHT_PAREN: ")", // )
	LEFT_BRACE:  "{", // {
	RIGHT_BRACE: "}", // }
	COMMA:       ",", // ,
	DOT:         ".", // .
	MINUS:       "-", // -
	PLUS:        "+", // +
	SEMICOLON:   ";", // ;
	SLASH:       "/", // 反斜线
	STAR:        "*", // *
	// todo 部分关键字含义不清楚

	BANG:          "!",
	BANG_EQUAL:    "!=",
	EQUAL:         "=",
	EQUAL_EQUAL:   "==",
	GREATER:       ">",
	GREATER_EQUAL: ">=",
	LESS:          "<",
	LESS_EQUAL:    "<=",

	/*
	  literals 字面量
	*/
	IDENTIFIER: "identifier",
	STRING:     "str",
	NUMBER:     "num",

	/*
		keywords for lox language
	*/

	AND:    "and",
	CLASS:  "class",
	ELSE:   "else",
	FALSE:  "false",
	FUN:    "fun",
	FOR:    "for",
	IF:     "if",
	NIL:    "nil",
	OR:     "or",
	PRINT:  "print",
	RETURN: "return",
	SUPER:  "super",
	THIS:   "this",
	TRUE:   "true",
	VAR:    "var",
	WHILE:  "while",

	EOF: "eof",
}

func (t TokenType) String() string {
	return TokenTypeMap[t]
}
