package Syntax

import (
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Token"
)

// Resolver 用于变量的解析

type Resolver struct {
	*Interpreter

	// stack for scopes
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
}

func (r *Resolver) VisitSuperExpr(superexpr Expr) interface{} {
	class := superexpr.(*SuperExpr)
	r.resolveLocal(class, class.keyword)
	return nil
}

func (r *Resolver) VisitThisExpr(thisexpr Expr) interface{} {
	class := thisexpr.(*ThisExpr)
	if r.currentClass == NoneClass {
		Errors.LoxError(class.keyword, ""+
			"Can't use 'this' outside of a class.")
		return nil
	}
	r.resolveLocal(class, class.keyword)
	return nil
}

func (r *Resolver) VisitGetExpr(getexpr Expr) interface{} {
	class := getexpr.(*GetExpr)
	r.resolveExpr(class.object)
	return nil
}

func (r *Resolver) VisitSetExpr(setexpr Expr) interface{} {
	class := setexpr.(*SetExpr)
	r.resolveExpr(class.value)
	r.resolveExpr(class.object)
	return nil
}

func (r *Resolver) VisitClassStmt(classstmt Stmt) interface{} {
	class := classstmt.(*ClassStmt)
	enclosingClass := r.currentClass
	r.currentClass = Class

	r.declare(class.name)
	r.define(class.name)

	// 循环依赖可以最后添加图检测环的算法
	if class.superClass != nil &&
		class.name.Lexeme == class.superClass.name.Lexeme {
		Errors.LoxError(class.superClass.name,
			"A class can't inherit from itself.")
	}

	if class.superClass != nil {
		r.resolveExpr(class.superClass)

		r.beginScope()
		r.peek()["super"] = true
	}
	// for this pointer and methods
	r.beginScope()
	r.peek()["this"] = true

	for _, item := range class.methods {
		declaration := METHOD
		if item.(*FunctionStmt).name.Lexeme == "init" {
			declaration = ISINITIALIZER
		}
		r.resolveFunction(item, declaration)
	}
	r.endScope()
	r.currentClass = enclosingClass
	if class.superClass != nil {
		r.endScope()
	}
	return nil
}

func (r *Resolver) peek() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) VisitFunctionStmt(stmt Stmt) interface{} {
	class := stmt.(*FunctionStmt)
	r.declare(class.name)
	r.define(class.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr Expr) interface{} {
	class := expr.(*VariableExpr)
	if len(r.scopes) != 0 {
		if v, ok := r.peek()[class.name.Lexeme]; ok && !v {
			// var is initialized in its own initializer
			Errors.LoxError(class.name, ""+
				"Can't read local var in its own initializer.")
		}
	}
	r.resolveLocal(class, class.name)

	return nil
}

func (r *Resolver) VisitAssignmentExpr(expr Expr) interface{} {
	class := expr.(*AssignmentExpr)
	r.resolveExpr(class.value)
	r.resolveLocal(class, class.name)
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt Stmt) interface{} {
	class := stmt.(*BlockStmt)
	r.beginScope()
	r.ResolveStmts(class.statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVariableStmt(stmt Stmt) interface{} {
	class := stmt.(*VariableStmt)
	r.declare(class.name)
	if class.initializer != nil {
		r.resolveExpr(class.initializer)
	}
	r.define(class.name)
	// 将声明和定义分开的原因
	//var a = "outer";
	//{
	//	var a = a;
	//}
	return nil
}

func (r *Resolver) declare(name *Token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.peek()
	// 只是声明, 没有赋值
	if _, ok := scope[name.Lexeme]; ok {
		Errors.LoxError(name, ""+
			"Already a variable with this name in this scope.")
	} else {
		scope[name.Lexeme] = false
	}
}

func (r *Resolver) define(name *Token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) ResolveStmts(stmts []Stmt) {
	for _, item := range stmts {
		r.resolveStmt(item)
	}
}

// 解析statement
func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

// 处理expr
func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

// 获取变量的深度
func (r *Resolver) resolveDeep(expr Expr, deep int) {
	// todo 记录变量访问到的作用域
	r.resolve(expr, deep)
}

func (r *Resolver) resolveLocal(expr Expr, token *Token.Token) {
	// [0, 1, 2, ..., n-1]
	// find in i
	// n - 1 - i
	// if resolve all but not found, it's in global.
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][token.Lexeme]; ok {
			r.resolveDeep(expr, len(r.scopes)-1-i)
		}
	}
}

func (r *Resolver) resolveFunction(stmt Stmt, functionType FunctionType) {
	class := stmt.(*FunctionStmt)
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, token := range class.params {
		r.declare(token)
		r.define(token)
	}
	r.ResolveStmts(class.body)

	r.endScope()

	r.currentFunction = enclosingFunction
}

// 下面不涉及变量和作用域的操作, 但是需要重写进行遍历

func (r *Resolver) VisitExpressionStmt(stmt Stmt) interface{} {
	class := stmt.(*ExpressionStmt)
	r.resolveExpr(class.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt Stmt) interface{} {
	class := stmt.(*IfStmt)
	r.resolveExpr(class.condition)
	r.resolveStmt(class.thenBranch)
	if class.elseBranch != nil {
		r.resolveStmt(class.elseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt Stmt) interface{} {
	class := stmt.(*PrintStmt)
	r.resolveExpr(class.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt Stmt) interface{} {
	class := stmt.(*ReturnStmt)
	if r.currentFunction == None {
		Errors.LoxError(class.keyword, "Can't return from top-level code.")
	}

	if class.value != nil {
		if r.currentFunction == ISINITIALIZER {
			Errors.LoxError(class.keyword, "Can't return a value"+
				" from an initializer")
		}
		r.resolveExpr(class.value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt Stmt) interface{} {
	class := stmt.(*WhileStmt)
	r.resolveExpr(class.condition)
	r.resolveStmt(class.body)
	return nil
}

func (r *Resolver) VisitCallExpr(expr Expr) interface{} {
	class := expr.(*CallExpr)
	r.resolveExpr(class.callee)
	for _, item := range class.arguments {
		r.resolveExpr(item)
	}

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr Expr) interface{} {
	class := expr.(*GroupingExpr)
	r.resolveExpr(class)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr Expr) interface{} {
	return nil
}

func (r *Resolver) VisitLogicExpr(expr Expr) interface{} {
	class := expr.(*LogicExpr)
	r.resolveExpr(class.left)
	r.resolveExpr(class.right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr Expr) interface{} {
	class := expr.(*UnaryExpr)
	r.resolveExpr(class.right)
	return nil
}

func NewResolver(i *Interpreter) *Resolver {
	scopes := make([]map[string]bool, 0)
	//scopes[0] = make(map[string]bool) // 代表全局?
	// 用于检测return语句在当前情况是否可行
	return &Resolver{i, scopes,
		None, NoneClass}
}

type FunctionType int32

// 标识return语句的可用范围
const (
	None          FunctionType = iota + 1
	FUNCTION                   // normal function
	ISINITIALIZER              // constructor no return value
	METHOD                     // function bind to a class
)

type ClassType int32

const (
	NoneClass ClassType = iota + 1
	Class
)
