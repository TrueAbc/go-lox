package main

import (
	"bufio"
	"fmt"
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
		fmt.Print("> ")
		// scan a line and text get result
		text := reader.Text()

		if len(text) != 0 {
			run(text)
			hadError = false
		} else {
			break
		}
	}
}

// 读取source内容并执行
func run(source string) {
	scanner := Token.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func errorHandle(line int, message string) {

}

func report(line int, where, message string) {
	hadError = true
	err := fmt.Errorf("[line %d ] Error %v : %v", line, where, message)
	fmt.Println(err)
}
