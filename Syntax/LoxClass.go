package Syntax

type LoxClass struct {
	name string
}

func (lc *LoxClass) Arity() int {
	return 0
}

func (lc *LoxClass) Call(interpreter *Interpreter, args []interface{}) interface{} {
	instance := &LoxInstance{lc}
	return instance
}

func (lc *LoxClass) String() string {
	return lc.name
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

type LoxInstance struct {
	kClass *LoxClass
}

func (li *LoxInstance) String() string {
	return li.kClass.String() + " instance."
}
