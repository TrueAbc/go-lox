package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Expr interface {
	Accept(visitor Visitor) interface{}
}
type Visitor interface {
	VisitBinaryExpr(binary Expr) interface{}
	VisitGroupingExpr(grouping Expr) interface{}
	VisitLiteralExpr(literal Expr) interface{}
	VisitUnaryExpr(unary Expr) interface{}
}
type Binary struct {
	left     Expr
	operator *Token.Token
	right    Expr
}

func (binary *Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(binary)
}

type Grouping struct {
	expression Expr
}

func (grouping *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpr(grouping)
}

type Literal struct {
	value interface{}
}

func (literal *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpr(literal)
}

type Unary struct {
	operator *Token.Token
	right    Expr
}

func (unary *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(unary)
}
