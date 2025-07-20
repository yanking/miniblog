package main

import (
	"github.com/yanking/miniblog/cmd/mb-apiserver/app"
	_ "go.uber.org/automaxprocs"
	"os"
)

func main() {
	command := app.NewMiniBlogCommand()

	// 执行命令并处理错误
	if err := command.Execute(); err != nil {
		// 如果发生错误，则退出程序
		// 返回退出码，可以使其他程序（例如 bash 脚本）根据退出码来判断服务运行状态
		os.Exit(1)
	}
}
