package Syntax

import (
	"fmt"
	"github.com/trueabc/lox/Token"
	"strings"
)

type AstPrinter struct {
}

func (a AstPrinter) VisitBinaryExpr(binary Expr) interface{} {
	class := binary.(*BinaryExpr)
	return a.parenthesize(class.operator.Lexeme,
		class.left, class.right)
}

func (a AstPrinter) VisitGroupingExpr(grouping Expr) interface{} {
	//TODO implement me
	class := grouping.(*GroupingExpr)
	return a.parenthesize("group", class.expression)
}

func (a AstPrinter) VisitLiteralExpr(literal Expr) interface{} {
	//TODO implement me
	class := literal.(*LiteralExpr)
	if class.value == nil {
		return Token.NIL
	}
	return class.value
}

func (a AstPrinter) VisitUnaryExpr(unary Expr) interface{} {
	class := unary.(*UnaryExpr)
	return a.parenthesize(class.operator.Lexeme, class.right)
}

func (a AstPrinter) Print(expr Expr) interface{} {
	return expr.Accept(a)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var sb strings.Builder

	sb.WriteString("(" + name)
	for _, exp := range exprs {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", exp.Accept(a)))

	}
	sb.WriteString(")")

	return sb.String()
}
