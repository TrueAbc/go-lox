package main

import (
	"bufio"
	"fmt"
	"github.com/trueabc/lox/Logger"
	"io/fs"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		Logger.Debugf("Usaget: generate_ast <out dir>")
		os.Exit(64)
	}
	outDir := args[1]

	defineAst(outDir, "Expr", []string{
		"Binary   : Expr left, *Token.Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : interface{} value",
		"Unary    : *Token.Token operator, Expr right",
		"Variable : *Token.Token name",
		"This : *Token.Token keyword",
		"Super    : *Token.Token keyword, *Token.Token method",
		"Get      : Expr object, *Token.Token name",
		"Set      : Expr object, *Token.Token name, Expr value",
		"Logic : Expr left, *Token.Token operator, Expr right",
		"Assignment : *Token.Token name, Expr value",
		"Call     : Expr callee, *Token.Token paren, []Expr arguments",
	})

	defineAst(outDir, "Stmt", []string{
		"Expression : Expr Expression",
		"Print : Expr Expression",
		"Variable : *Token.Token name, Expr initializer",
		"Block : []Stmt statements",
		"Class      : *Token.Token name, *VariableExpr superClass,  []Stmt methods",
		"While : Expr condition, Stmt body",
		"If : Expr condition, Stmt thenBranch," +
			" Stmt elseBranch",
		"Return     : *Token.Token keyword, Expr value",
		"Function   : *Token.Token name, []*Token.Token params," +
			" []Stmt body",
	})
}

func defineAst(outDir string, inter string, obj []string) {
	file, _ := os.OpenFile(outDir+"/"+inter+".go", os.O_TRUNC|os.O_RDWR|os.O_CREATE, fs.ModePerm)
	pacakage := "package Syntax\n"
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(pacakage)

	write.WriteString("import (\n")
	write.WriteString("\t" + "\"github.com/trueabc/lox/Token\"\n")
	write.WriteString(")\n")

	// 表达式定义
	write.WriteString("type " + inter + " interface { \n")
	write.WriteString(fmt.Sprintf("\tAccept(visitor Visitor%s) interface{}", inter))
	write.WriteString("\n")
	write.WriteString("}\n")

	// Visitor定义
	defineVisitor(write, inter, obj)

	// AST 结点
	for _, item := range obj {
		members := strings.Split(item, ":")
		name := strings.TrimSpace(members[0])
		fields := strings.Split(
			strings.TrimSpace(members[1]), ",")
		defineType(write, inter, name, fields)
	}

	write.Flush()
}

func defineVisitor(out *bufio.Writer, baseName string, fields []string) {
	out.WriteString(fmt.Sprintf("type Visitor%s interface { \n",
		baseName))

	for _, i := range fields {
		temp := strings.Split(i, ":")
		t := strings.TrimSpace(temp[0])
		out.WriteString(fmt.Sprintf("\tVisit%s%s(%s %s) interface{}\n", t, baseName, strings.ToLower(t+baseName), baseName))
	}

	out.WriteString("}\n")
}

func defineType(out *bufio.Writer, baseName, className string, fields []string) {
	out.WriteString("type " + className + baseName + " struct { \n")
	for _, i := range fields {
		i = strings.TrimSpace(i)
		items := strings.Split(i, " ")
		name := strings.TrimSpace(items[1])
		t := strings.TrimSpace(items[0])
		out.WriteString("\t " + name + "\t" + t + "\n")
	}

	out.WriteString(" }")
	out.WriteString("\n")
	// 添加accept方法
	// Accept(visitor Visitor) interface{}
	// func (b *Binary) Accept(visitor Visitor) interface{}  {
	//
	//}
	out.WriteString(fmt.Sprintf("func (%s *%s) Accept(visitor Visitor%s) interface{} {\n",
		strings.ToLower(className+baseName), className+baseName, baseName))
	out.WriteString("return visitor." + "Visit" + className + baseName + "(" + strings.ToLower(className+baseName) + ")")

	out.WriteString("}\n")

}
