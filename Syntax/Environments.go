package Syntax

import (
	"github.com/trueabc/lox/Token"
)

// 存储变量信息

type Environment struct {
	Enclosing *Environment
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
	} else if ev.Enclosing != nil {
		return ev.Enclosing.Get(token)
	}
	panic(NewRuntimeError(token, "Undefined Var "+token.Lexeme+"."))
}

func (ev *Environment) GetAt(dis int, token *Token.Token) interface{} {
	ancestor := ev.Ancestor(dis)
	return ancestor.VarValues[token.Lexeme]
}

func (ev *Environment) Ancestor(dis int) *Environment {
	env := ev
	for i := 0; i < dis; i++ {
		env = env.Enclosing
	}
	return env
}

func (ev *Environment) Assign(token *Token.Token, value interface{}) interface{} {
	if _, ok := ev.VarValues[token.Lexeme]; ok {
		ev.VarValues[token.Lexeme] = value
		return ev.VarValues[token.Lexeme]
	} else if ev.Enclosing != nil {
		return ev.Enclosing.Assign(token, value)
	}
	panic(NewRuntimeError(token, "Undefined var name "+token.Lexeme+"."))
}

func (ev *Environment) AssignAt(dis int, token *Token.Token, value interface{}) interface{} {
	return ev.Ancestor(dis).Assign(token, value)
}

// NewEnvironment 全局作用域和局部作用域
func NewEnvironment() *Environment {
	return &Environment{nil, make(map[string]interface{})}
}

func NewLocalEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing, make(map[string]interface{})}
}
