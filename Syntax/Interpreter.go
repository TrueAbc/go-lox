package Syntax

import (
	"github.com/trueabc/lox/Token"
)

type Interpreter struct {
}

func (i Interpreter) VisitBinaryExpr(binary Expr) interface{} {
	class := binary.(*Binary)
	left := i.evaluate(class.left)
	right := i.evaluate(class.right)
	switch class.operator.TType {
	case Token.MINUS:
		return left.(float64) - right.(float64)
	case Token.STAR:
		return left.(float64) * right.(float64)
	case Token.SLASH:
		// 分母为0
		return left.(float64) / right.(float64)
	case Token.PLUS:
		// 字符串拼接和数字相加
		l1, ok1 := left.(float64)
		r1, ok2 := right.(float64)
		if ok1 && ok2 {
			return l1 + r1
		}
		s1 := left.(string)
		s2 := right.(string)
		return s1 + s2
	case Token.GREATER:
		return left.(float64) > right.(float64)
	case Token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case Token.LESS:
		return left.(float64) < right.(float64)
	case Token.LESS_EQUAL:
		return left.(float64) <= right.(float64)

	case Token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case Token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i Interpreter) VisitGroupingExpr(grouping Expr) interface{} {
	class := grouping.(*Grouping)
	return i.evaluate(class.expression)
}

func (i Interpreter) VisitLiteralExpr(literal Expr) interface{} {
	class := literal.(*Literal)
	return class.value
}

func (i Interpreter) VisitUnaryExpr(unary Expr) interface{} {
	class := unary.(*Unary)
	right := i.evaluate(class.right)
	switch class.operator.TType {
	case Token.MINUS:
		return -(right.(float64))
	case Token.BANG:
		return !i.isTruthy(right)
	}
	return nil
}

func (i Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case bool:
		return value.(bool)
	}
	return true
}

func (i Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}
