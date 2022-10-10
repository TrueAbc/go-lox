package Syntax

import (
	"time"
)

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, args []interface{}) interface{}
}

// ClockFunc Native Function
type ClockFunc struct {
}

func (c ClockFunc) Arity() int {
	return 0
}

func (c ClockFunc) Call(interpreter *Interpreter, args []interface{}) interface{} {
	return time.Now().Second()
}

func (c ClockFunc) String() string {
	return "<native fn>"
}

type LoxFunction struct {
	funcStmt *FunctionStmt
	Closure  *Environment
}

func (l *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	env := NewLocalEnvironment(l.Closure)
	env.Define("this", instance)
	return NewLoxFunction(l.funcStmt, env)
}

func (l *LoxFunction) Arity() int {
	return len(l.funcStmt.params)
}

func (l *LoxFunction) Call(interpreter *Interpreter, args []interface{}) (result interface{}) {
	// 默认是使用全局变量, 闭包在这里需要考虑其他
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(*ReturnObj); ok {
				result = v.Value
			}
		}
	}()
	env := NewLocalEnvironment(l.Closure)
	for id, item := range l.funcStmt.params {
		env.Define(item.Lexeme, args[id])
	}
	interpreter.executeBlock(l.funcStmt.body, env)
	return result
}

func (l *LoxFunction) String() string {
	return "<fn " + l.funcStmt.name.Lexeme + " >"
}

func NewLoxFunction(declaration *FunctionStmt, closure *Environment) *LoxFunction {
	f := &LoxFunction{funcStmt: declaration, Closure: closure}
	return f
}
