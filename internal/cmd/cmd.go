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
package cmd

import (
	"context"
	// 导入各个控制器（自己写的）
	"library-management-api/internal/controller/user"
	"library-management-api/internal/controller/book"
	"library-management-api/internal/controller/borrow"
	// 导入中间件（自己写的）
	"library-management-api/internal/service/middleware"

	"github.com/gogf/gf/v2/frame/g"     // GoFrame框架核心
	"github.com/gogf/gf/v2/net/ghttp"   // HTTP服务器
	"github.com/gogf/gf/v2/os/gcmd"     // 命令行工具
)

var (
	// Main 定义一个命令行命令，运行 "main" 时执行 mainFunc 函数
	Main = &gcmd.Command{
		Name:  "main",
		Brief: "start http server of library management system",
		Func:  mainFunc,
	}
)

// mainFunc 主函数：启动HTTP服务器并注册所有API路由
// 【重点】为什么路由要集中写在这里
// 答：集中管理所有API路径，方便查看和维护，避免路由分散在各处
func mainFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 第1步：创建HTTP服务器实例（默认监听8000端口）
	s := g.Server()

	// 第2步：创建中间件实例，用于身份验证和权限控制
	middlewareSvc := middleware.New()   //中间件

	// 第3步：注册全局中间件（所有请求都会经过）
	s.Use(ghttp.MiddlewareHandlerResponse)   // GoFrame内置：统一响应格式
	s.Use(middlewareSvc.Ctx)                 //Ctx中间件：每个请求都解析用户身份
	s.Use(ghttp.MiddlewareCORS)              // CORS：允许跨域（前端3000端口调后端8000端口）

	// 第4步：创建控制器实例
	userCtrl := user.NewV1()     //用户控制器
	bookCtrl := book.NewV1()     //图书控制器
	borrowCtrl := borrow.NewV1() //借阅控制器

	// ============================================================
	// 【核心】RESTful API路由设计
	// 设计原则：
	//   - 名词复数表示资源：/books、/users、/borrows
	//   - HTTP方法表示操作：POST创建、GET查询、PUT更新、DELETE删除
	//   - 按功能模块分组：认证(/auth)、用户(/users)、图书(/books)、借阅(/borrows)
	//   - 分层权限：公开 > 登录 > 管理员
	// ============================================================
	s.Group("/api", func(group *ghttp.RouterGroup) {

		// ------ 认证服务（无需登录）------
		// POST   /api/auth/verification-codes  发送验证码
		// POST   /api/auth/register            邮箱注册
		// POST   /api/auth/login               登录
		// DELETE /api/auth/logout              退出登录
		// PUT    /api/auth/password            重置密码
		group.Group("/auth", func(group *ghttp.RouterGroup) {
			group.POST("/verification-codes", userCtrl.SendVerificationCode)
			group.POST("/register", userCtrl.Register)
			group.POST("/login", userCtrl.Login)
			group.DELETE("/logout", userCtrl.SignOut)
			group.PUT("/password", userCtrl.ResetPassword)
		})

		// ------ 用户服务 ------
		// POST   /api/users                    手机号注册
		// GET    /api/users/me                 个人信息
		group.Group("/users", func(group *ghttp.RouterGroup) {
			group.POST("/", userCtrl.SignUp)
			group.GET("/me", userCtrl.Profile)
		})

		// ------ 图书查询（公开，无需登录）------
		// GET    /api/books                    图书列表
		// GET    /api/books/{id}               图书详情
		group.GET("/books", bookCtrl.List)
		group.GET("/books/{id}", bookCtrl.Get)

		// ------ 以下接口需要登录 ------
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(middlewareSvc.Auth)    //检查登录状态

			// ------ 图书管理（需管理员权限）------
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middlewareSvc.Admin)  //检查管理员权限
				group.POST("/books", bookCtrl.Create)
				group.PUT("/books/{id}", bookCtrl.Update)
				group.DELETE("/books/{id}", bookCtrl.Delete)
			})

			// ------ 借阅管理（登录即可）------
			// POST   /api/borrows               借书
			// GET    /api/borrows               借阅列表
			// PUT    /api/borrows/{id}/return   还书
			group.POST("/borrows", borrowCtrl.Borrow)
			group.GET("/borrows", borrowCtrl.List)
			group.PUT("/borrows/{id}/return", borrowCtrl.Return)
		})
	})

	// 启动HTTP服务器
	s.Run()
	return nil
}

