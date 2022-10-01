package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Stmt interface {
	Accept(visitor VisitorStmt) interface{}
}
type VisitorStmt interface {
	VisitExpressionStmt(expression Stmt) interface{}
	VisitPrintStmt(print Stmt) interface{}
	VisitVariableStmt(variable Stmt) interface{}
}
type ExpressionStmt struct {
	Expression Expr
}

func (expression *ExpressionStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitExpressionStmt(expression)
}

type PrintStmt struct {
	Expression Expr
}

func (print *PrintStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitPrintStmt(print)
}

type VariableStmt struct {
	name        *Token.Token
	initializer Expr
}

func (variable *VariableStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitVariableStmt(variable)
}
