package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Expr interface {
	Accept(visitor VisitorExpr) interface{}
}
type VisitorExpr interface {
	VisitBinaryExpr(binary Expr) interface{}
	VisitGroupingExpr(grouping Expr) interface{}
	VisitLiteralExpr(literal Expr) interface{}
	VisitUnaryExpr(unary Expr) interface{}
	VisitVariableExpr(variable Expr) interface{}
}
type BinaryExpr struct {
	left     Expr
	operator *Token.Token
	right    Expr
}

func (binary *BinaryExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitBinaryExpr(binary)
}

type GroupingExpr struct {
	expression Expr
}

func (grouping *GroupingExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitGroupingExpr(grouping)
}

type LiteralExpr struct {
	value interface{}
}

func (literal *LiteralExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitLiteralExpr(literal)
}

type UnaryExpr struct {
	operator *Token.Token
	right    Expr
}

func (unary *UnaryExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitUnaryExpr(unary)
}

type VariableExpr struct {
	name *Token.Token
}

func (variable *VariableExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitVariableExpr(variable)
}
