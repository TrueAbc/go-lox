package Syntax

import (
	"github.com/trueabc/lox/Token"
)

// 存储变量信息

type Environment struct {
	VarValues map[string]interface{}
}

// Define 定义一个变量
func (ev *Environment) Define(name string, value interface{}) {
	// todo 暂时允许变量的重复定义
	ev.VarValues[name] = value
}

func (ev *Environment) Get(token *Token.Token) interface{} {
	if v, ok := ev.VarValues[token.Lexeme]; ok {
		return v
	}
	panic(NewRuntimeError(token, "Undefined Var "+token.Lexeme+"."))
}

func (ev *Environment) Assign(token *Token.Token, value interface{}) interface{} {
	if _, ok := ev.VarValues[token.Lexeme]; ok {
		ev.VarValues[token.Lexeme] = value
	}
	panic(NewRuntimeError(token, "Undefined var name "+token.Lexeme+"."))
}

func NewEnvironment() *Environment {
	return &Environment{make(map[string]interface{})}
}
