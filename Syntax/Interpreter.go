package Syntax

import (
	"fmt"
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Logger"
	"github.com/trueabc/lox/Token"
)

type Interpreter struct {
	env *Environment

	// 持有最外层的引用
	global *Environment
}

func (i *Interpreter) VisitFunctionStmt(functionstmt Stmt) interface{} {
	class := functionstmt.(*FunctionStmt)
	function := NewLoxFunction(class)
	i.env.Define(class.name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitCallExpr(callexpr Expr) interface{} {
	class := callexpr.(*CallExpr)
	// 将函数名称转为对象
	callee := i.evaluate(class.callee)
	args := make([]interface{}, 0)
	for _, item := range class.arguments {
		args = append(args, i.evaluate(item))
	}
	funCall, ok := callee.(LoxCallable)
	if !ok {
		panic(NewRuntimeError(class.paren, ""+
			"can only call functions and classes."))
	}
	if funCall.Arity() != len(args) {
		panic(NewRuntimeError(class.paren,
			fmt.Sprintf("Expected %d arguments. Got %d arguments", funCall.Arity(), len(args))))
	}
	return funCall.Call(i, args)
}

func (i *Interpreter) VisitWhileStmt(whilestmt Stmt) interface{} {
	class := whilestmt.(*WhileStmt)
	for i.isTruthy(i.evaluate(class.condition)) {
		i.execute(class.body)
	}
	return nil
}

func (i *Interpreter) VisitLogicExpr(logicexpr Expr) interface{} {
	class := logicexpr.(*LogicExpr)
	left := i.evaluate(class.left)
	operator := class.operator
	if operator.TType == Token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(class.right)
}

func (i *Interpreter) VisitIfStmt(ifstmt Stmt) interface{} {
	class := ifstmt.(*IfStmt)
	if i.isTruthy(i.evaluate(class.condition)) {
		i.execute(class.thenBranch)
	} else if class.elseBranch != nil {
		i.execute(class.elseBranch)
	}
	return nil
}

func (i *Interpreter) VisitBlockStmt(block Stmt) interface{} {
	class := block.(*BlockStmt)
	i.executeBlock(class.statements, NewLocalEnvironment(i.env))
	return nil
}

func (i *Interpreter) executeBlock(stmt []Stmt, environment *Environment) {
	previous := i.env
	defer func() {
		i.env = previous
	}()
	i.env = environment
	// 执行内部的子句, 结束后恢复全局作用域
	for _, s := range stmt {
		i.execute(s)
	}
}

func (i *Interpreter) VisitAssignmentExpr(assignment Expr) interface{} {
	class := assignment.(*AssignmentExpr)
	value := i.evaluate(class.value)
	i.env.Assign(class.name, value)
	return value
}

func (i *Interpreter) VisitVariableStmt(variable Stmt) interface{} {
	var value interface{}
	class := variable.(*VariableStmt)
	if class.initializer != nil {
		value = i.evaluate(class.initializer)
	}
	i.env.Define(class.name.Lexeme, value)
	return nil
}

// statements 的返回值都是nil

func (i *Interpreter) VisitExpressionStmt(expression Stmt) interface{} {
	class := expression.(*ExpressionStmt)
	i.evaluate(class.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(print Stmt) interface{} {
	class := print.(*PrintStmt)
	value := i.evaluate(class.Expression)
	fmt.Println(value)
	return nil
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func NewInterpreter() *Interpreter {
	global := NewEnvironment()
	global.Define("clock", ClockFunc{})
	return &Interpreter{env: global, global: global}
}

func (i *Interpreter) Interpret(statements []Stmt) interface{} {
	for _, s := range statements {
		i.executeSingle(s)
	}

	return nil
}

func (i *Interpreter) executeSingle(stmt Stmt) {
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

func (i *Interpreter) VisitBinaryExpr(binary Expr) interface{} {
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

func (i *Interpreter) VisitGroupingExpr(grouping Expr) interface{} {
	class := grouping.(*GroupingExpr)
	return i.evaluate(class.expression)
}

func (i *Interpreter) VisitLiteralExpr(literal Expr) interface{} {
	class := literal.(*LiteralExpr)
	return class.value
}

func (i *Interpreter) VisitUnaryExpr(unary Expr) interface{} {
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

func (i *Interpreter) VisitVariableExpr(variable Expr) interface{} {
	class := variable.(*VariableExpr)
	return i.env.Get(class.name)
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

// 非空对象和bool的true
func (i *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case bool:
		return value.(bool)
	}
	return true
}

func (i *Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}

func (i *Interpreter) checkNumberOperand(operator *Token.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	default:
		panic(NewRuntimeError(operator, "operand must be number. "))
	}
}

func (i *Interpreter) checkNumberOperands(operator *Token.Token, left, right interface{}) {
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
