package Syntax

import (
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Token"
)

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (lc *LoxClass) FindMethod(name string) *LoxFunction {
	if v, ok := lc.methods[name]; ok {
		return v
	}
	return nil
}
func (lc *LoxClass) Arity() int {
	return 0
}

func (lc *LoxClass) Call(interpreter *Interpreter, args []interface{}) interface{} {
	instance := NewLoxInstance(lc)
	return instance
}

func (lc *LoxClass) String() string {
	return lc.name
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

type LoxInstance struct {
	kClass *LoxClass
	fields map[string]interface{}
}

func (li *LoxInstance) Get(token *Token.Token) interface{} {
	if v, ok := li.fields[token.Lexeme]; ok {
		return v
	}
	if method := li.kClass.FindMethod(token.Lexeme); method != nil {
		return method
	}
	Errors.LoxRuntimeError(token, "Undefined property '"+token.Lexeme+"'.")
	return nil
}

func (li *LoxInstance) Set(token *Token.Token, value interface{}) interface{} {
	li.fields[token.Lexeme] = value
	return value
}

func (li *LoxInstance) String() string {
	return li.kClass.String() + " instance."
}

func NewLoxInstance(kClass *LoxClass) *LoxInstance {
	return &LoxInstance{kClass: kClass, fields: make(map[string]interface{})}
}
