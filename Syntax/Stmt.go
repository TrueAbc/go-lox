package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Stmt interface {
	Accept(visitor VisitorStmt) interface{}
}
type VisitorStmt interface {
	VisitExpressionStmt(expressionstmt Stmt) interface{}
	VisitPrintStmt(printstmt Stmt) interface{}
	VisitVariableStmt(variablestmt Stmt) interface{}
	VisitBlockStmt(blockstmt Stmt) interface{}
	VisitWhileStmt(whilestmt Stmt) interface{}
	VisitIfStmt(ifstmt Stmt) interface{}
	VisitReturnStmt(returnstmt Stmt) interface{}
	VisitFunctionStmt(functionstmt Stmt) interface{}
}
type ExpressionStmt struct {
	Expression Expr
}

func (expressionstmt *ExpressionStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitExpressionStmt(expressionstmt)
}

type PrintStmt struct {
	Expression Expr
}

func (printstmt *PrintStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitPrintStmt(printstmt)
}

type VariableStmt struct {
	name        *Token.Token
	initializer Expr
}

func (variablestmt *VariableStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitVariableStmt(variablestmt)
}

type BlockStmt struct {
	statements []Stmt
}

func (blockstmt *BlockStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitBlockStmt(blockstmt)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (whilestmt *WhileStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitWhileStmt(whilestmt)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (ifstmt *IfStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitIfStmt(ifstmt)
}

type ReturnStmt struct {
	keyword *Token.Token
	value   Expr
}

func (returnstmt *ReturnStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitReturnStmt(returnstmt)
}

type FunctionStmt struct {
	name   *Token.Token
	params []*Token.Token
	body   []Stmt
}

func (functionstmt *FunctionStmt) Accept(visitor VisitorStmt) interface{} {
	return visitor.VisitFunctionStmt(functionstmt)
}
