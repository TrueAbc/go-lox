

# 该makefile 执行的路径应该是当前文件夹
# 生成语法结点
clean:
	-rm ../Expr.go

build:
	go run generateAST.go ../

all: clean build
	gofmt -w ../Expr.go