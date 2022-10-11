package Syntax

import (
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Token"
)

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction

	superClass *LoxClass
}

func (lc *LoxClass) FindMethod(name string) *LoxFunction {
	if v, ok := lc.methods[name]; ok {
		return v
	}
	return nil
}
func (lc *LoxClass) Arity() int {
	initializer := lc.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (lc *LoxClass) Call(interpreter *Interpreter, args []interface{}) interface{} {
	instance := NewLoxInstance(lc)
	initializer := lc.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, args)
	}

	return instance
}

func (lc *LoxClass) String() string {
	return lc.name
}

func NewLoxClass(name string, methods map[string]*LoxFunction,
	superClass interface{}) *LoxClass {
	if superClass != nil {
		return &LoxClass{name: name, methods: methods, superClass: superClass.(*LoxClass)}
	} else {
		superClass = nil
		return &LoxClass{name: name, methods: methods, superClass: nil}
	}
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
		return method.Bind(li)
	}
	if li.kClass.superClass != nil {
		return li.kClass.superClass.FindMethod(token.Lexeme)
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
