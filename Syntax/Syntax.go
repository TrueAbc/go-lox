package Syntax

/*
 需要通过优先级和关联性规则消除文法的歧义

优先级：determines which operator is evaluated first in an expression containing a mixture of different operators
关联性：determines which operator is evaluated first in a series of the same operator.
6 - 3 - 2 从左到右匹配
a = ( b = c ) 从右到左

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

// 符号的优先级定义: 低到高

// Name	Operators	Associates
// Equality	== !=	Left
// Comparison	> >= < <=	Left
// Term	- +	Left
// Factor	/ *	Left
// Unary	! -	Right

// 修改后低优先级的production只能匹配同级别或者高级别
// todo 暂时不是关注的重点
// 这里使用递归下降进行解析
//expression     → equality ;
//equality       → comparison ( ( "!=" | "==" ) comparison )* ;
//comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
//term           → factor ( ( "-" | "+" ) factor )* ;
//factor         → unary ( ( "/" | "*" ) unary )* ;
//unary       Token equals = previous();
//      Expr value = assignment();
//
//      if (expr instanceof Expr.Variable) {
//        Token name = ((Expr.Variable)expr).name;
//        return new Expr.Assign(name, value);
//      }
//
//      error(equals, "Invalid assignment target.");    → ( "!" | "-" ) unary
//| primary ;
//primary        → NUMBER | STRING | "true" | "false" | "nil"
//| "(" expression ")" ;

// 转换语法的大致逻辑
// Grammar notation	Code representation
// Terminal	Code to match and consume a token
//Nonterminal	Call to that rule’s function
//|	if or switch statement
//* or +	while or for loop
//?	if statement

//  添加 statemenet 支持语法
//program        → statement* EOF ;
//
//statement      → exprStmt
//| printStmt ;
//
//exprStmt       → expression ";" ;
//printStmt      → "print" expression ";" ;

// 定义变量声明, 情况2 不允许, if内部的声明
//if (monday) print "Ugh, already?";
//if (monday) var beverage = "espresso";

//program        → declaration* EOF ;
//declaration    → varDecl | statement ;
//statement      → exprStmt| printStmt ;
//varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
// 变量声明之后, ast树的根结点也需要可以获取变量
//primary        → "true" | "false" | "nil"
//| NUMBER | STRING
//| "(" expression ")"
//| IDENTIFIER ;

// 赋值语句的特点
//expression     → assignment ;
//assignment     → IDENTIFIER "=" assignment
//| equality ;

// 作用域支持
//statement      → exprStmt
//| printStmt
//| block ;
//
//block          → "{" declaration* "}" ;
