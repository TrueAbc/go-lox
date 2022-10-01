package Syntax

type Stmt interface {
	Accept(visitor VisitorStmt) interface{}
}
type VisitorStmt interface {
	VisitExpressionStmt(expression Stmt) interface{}
	VisitPrintStmt(print Stmt) interface{}
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
