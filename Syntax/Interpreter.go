package Syntax

import (
	"fmt"
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Logger"
	"github.com/trueabc/lox/Token"
)

type Interpreter struct {
}

// statements 的返回值都是nil

func (i Interpreter) VisitExpressionStmt(expression Stmt) interface{} {
	class := expression.(*ExpressionStmt)
	i.evaluate(class.Expression)
	return nil
}

func (i Interpreter) VisitPrintStmt(print Stmt) interface{} {
	class := print.(*PrintStmt)
	value := i.evaluate(class.Expression)
	fmt.Println(value)
	return nil
}

func (i Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i Interpreter) Interpret(statements []Stmt) interface{} {
	for _, s := range statements {
		i.executeSingle(s)
	}

	return nil
}

func (i Interpreter) executeSingle(stmt Stmt) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				err := r.(*RuntimeError)
				Errors.LoxRuntimeError(err.Token, err.Content)
			}
			Logger.Errorf("%v", r)
		}
	}()
	i.execute(stmt)
}

// expression计算结果 四类expression

func (i Interpreter) VisitBinaryExpr(binary Expr) interface{} {
	class := binary.(*BinaryExpr)
	left := i.evaluate(class.left)
	right := i.evaluate(class.right)
	switch class.operator.TType {
	case Token.MINUS:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) - right.(float64)
	case Token.STAR:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) * right.(float64)
	case Token.SLASH:
		// 分母为0
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) / right.(float64)
	case Token.PLUS:
		// 字符串拼接和数字相加
		l1, ok1 := left.(float64)
		r1, ok2 := right.(float64)
		if ok1 && ok2 {
			return l1 + r1
		}
		s1, ok1 := left.(string)
		s2, ok2 := right.(string)
		if ok1 && ok2 {
			return s1 + s2
		}
		panic(NewRuntimeError(class.operator, "Operands must be two numbers or strings."))
	case Token.GREATER:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) > right.(float64)
	case Token.GREATER_EQUAL:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) >= right.(float64)
	case Token.LESS:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) < right.(float64)
	case Token.LESS_EQUAL:
		i.checkNumberOperands(class.operator, left, right)
		return left.(float64) <= right.(float64)

	case Token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case Token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i Interpreter) VisitGroupingExpr(grouping Expr) interface{} {
	class := grouping.(*GroupingExpr)
	return i.evaluate(class.expression)
}

func (i Interpreter) VisitLiteralExpr(literal Expr) interface{} {
	class := literal.(*LiteralExpr)
	return class.value
}

func (i Interpreter) VisitUnaryExpr(unary Expr) interface{} {
	class := unary.(*UnaryExpr)
	right := i.evaluate(class.right)
	switch class.operator.TType {
	case Token.MINUS:
		i.checkNumberOperand(class.operator, right)
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

func (i Interpreter) checkNumberOperand(operator *Token.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	default:
		panic(NewRuntimeError(operator, "operand must be number. "))
	}
}

func (i Interpreter) checkNumberOperands(operator *Token.Token, left, right interface{}) {
	_, ok1 := left.(float64)
	_, ok2 := right.(float64)
	if ok1 && ok2 {
		return
	}
	panic(NewRuntimeError(operator, "operand must be number."))
}

type RuntimeError struct {
	Token   *Token.Token
	Content string
}

func (re *RuntimeError) Error() string {
	return re.Content
}

func NewRuntimeError(token *Token.Token, content string) *RuntimeError {
	return &RuntimeError{Token: token, Content: content}
}
