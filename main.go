package main

import (
	"bufio"
	"fmt"
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Syntax"
	"github.com/trueabc/lox/Token"
	"os"
	"path/filepath"
)

var hadError = false

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: go-lox [script]")
		// sysexits.h 的一个错误码, 错误使用command
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		// 交互式的运行
		runPrompt()
	}
}

func runFile(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	file, err := os.ReadFile(abs)
	if err != nil {
		return
	}
	run(string(file))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		// scan a line and text get result
		text := reader.Text()

		if len(text) != 0 {
			run(text)
			hadError = false
		} else {
			break
		}
		fmt.Print("> ")
	}
}

// 读取source内容并执行
func run(source string) {
	scanner := Token.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := Syntax.NewParser(tokens)
	// res is an ast
	res := parser.Parse()

	interpreter := Syntax.NewInterpreter()
	interpreter.Interpret(res)

	if Errors.HadError {
		return
	}
	if Errors.HadRunTimeError {
		return
	}
	//fmt.Println(Syntax.AstPrinter{}.Print(res))
}
