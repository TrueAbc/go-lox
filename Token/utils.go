package Token

import (
	"fmt"
	"github.com/trueabc/lox/Logger"
)

var HadError = false

func report(line int, where, message string) string {
	HadError = true
	data := fmt.Sprintf("[line %d ] Error %v : %v", line, where, message)
	Logger.Errorf(data)
	return data
}

func LoxError(token *Token, mess string) string {
	if token.TType == EOF {
		return report(token.Line, " at end", mess)
	} else {
		return report(token.Line, "at '"+token.Lexeme+"'", mess)
	}
}
