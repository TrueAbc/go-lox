package Syntax

/*
基本的上下文无关文法
expression     → literal
               | unary
               | binary
               | grouping ;

literal        → NUMBER | STRING | "true" | "false" | "nil" ;
grouping       → "(" expression ")" ;
unary          → ( "-" | "!" ) expression ;
binary         → expression operator expression ;
operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
               | "+"  | "-"  | "*" | "/" ;
*/

// 需要完成元编程, 实现抽象语法树
//type Expr interface {
//}
//
//// 上下文无关文法转为抽象语法树
//
//// Binary 二元运算表达式
//type Binary struct {
//	left     Expr
//	operator Token.Token
//	right    Expr
//}

// 该文件只是标注使用, 产生的代码是Expr.go

// 这里需要使用访问者模式实现AST的结点的操作
// 添加一个新的操作是添加一个新的访问者类
