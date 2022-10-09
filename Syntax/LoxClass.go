package Syntax

type LoxClass struct {
	name string
}

func (lc *LoxClass) String() string {
	return lc.name
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}
