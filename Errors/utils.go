package Errors

import (
	"fmt"
	"github.com/trueabc/lox/Logger"
	"github.com/trueabc/lox/Token"
)

var HadError = false
var HadRunTimeError = false

func report(line int, where, message string) string {
	HadError = true
	data := fmt.Sprintf("[line %d ] Error %v : %v", line, where, message)
	Logger.Errorf(data)
	return data
}

func LoxError(token *Token.Token, mess string) string {
	if token.TType == Token.EOF {
		return report(token.Line, " at end", mess)
	} else {
		return report(token.Line, "at '"+token.Lexeme+"'", mess)
	}
}

func LoxRuntimeError(token *Token.Token, mess string) {
	HadRunTimeError = true
	// todo need to be panic
	report(token.Line, "at '"+token.Lexeme+"'",
		mess)
}
