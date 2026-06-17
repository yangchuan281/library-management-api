// ============================================================
// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package main

import (
	// 导入MySQL驱动（GoFrame框架提供的，用于连接MySQL数据库）
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"library-management-api/internal/cmd"
	"github.com/gogf/gf/v2/os/gctx"
)

// main 函数：程序启动的第一个入口
// 调用 cmd.Main.Run() 启动HTTP服务器（路由注册在 cmd.go 中）
func main() {
		// 【核心】启动HTTP服务器，具体的路由配置请看 internal/cmd/cmd.go
	cmd.Main.Run(gctx.New())
}


