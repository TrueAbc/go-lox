package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Expr interface {
	Accept(visitor VisitorExpr) interface{}
}
type VisitorExpr interface {
	VisitBinaryExpr(binaryexpr Expr) interface{}
	VisitGroupingExpr(groupingexpr Expr) interface{}
	VisitLiteralExpr(literalexpr Expr) interface{}
	VisitUnaryExpr(unaryexpr Expr) interface{}
	VisitVariableExpr(variableexpr Expr) interface{}
	VisitLogicExpr(logicexpr Expr) interface{}
	VisitAssignmentExpr(assignmentexpr Expr) interface{}
	VisitCallExpr(callexpr Expr) interface{}
}
type BinaryExpr struct {
	left     Expr
	operator *Token.Token
	right    Expr
}

func (binaryexpr *BinaryExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitBinaryExpr(binaryexpr)
}

type GroupingExpr struct {
	expression Expr
}

func (groupingexpr *GroupingExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitGroupingExpr(groupingexpr)
}

type LiteralExpr struct {
	value interface{}
}

func (literalexpr *LiteralExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitLiteralExpr(literalexpr)
}

type UnaryExpr struct {
	operator *Token.Token
	right    Expr
}

func (unaryexpr *UnaryExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitUnaryExpr(unaryexpr)
}

type VariableExpr struct {
	name *Token.Token
}

func (variableexpr *VariableExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitVariableExpr(variableexpr)
}

type LogicExpr struct {
	left     Expr
	operator *Token.Token
	right    Expr
}

func (logicexpr *LogicExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitLogicExpr(logicexpr)
}

type AssignmentExpr struct {
	name  *Token.Token
	value Expr
}

func (assignmentexpr *AssignmentExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitAssignmentExpr(assignmentexpr)
}

type CallExpr struct {
	callee    Expr
	paren     *Token.Token
	arguments []Expr
}

func (callexpr *CallExpr) Accept(visitor VisitorExpr) interface{} {
	return visitor.VisitCallExpr(callexpr)
}
