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
package middleware

import (
	"net/http"
	"strings"
	"library-management-api/internal/service/bizctx"
	"library-management-api/internal/service/jwt"
	"library-management-api/internal/service/session"
	"library-management-api/internal/service/user"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Service HTTP请求中间件
type Service struct {
	bizCtxSvc  *bizctx.Service
	sessionSvc *session.Service
	userSvc    *user.Service
	jwtSvc     *jwt.Service        // 实验三新增：JWT服务
}

func New() *Service {
	return &Service{
		bizCtxSvc:  bizctx.New(),
		sessionSvc: session.New(),
		userSvc:    user.New(),
		jwtSvc:     jwt.New(),
	}
}

// Ctx injects custom business context into request context.
//Ctx中间件——解析用户身份（每个请求都会经过）
// 1. 先从Session中获取用户信息（旧方式）
// 2. 如果Session中没有，再从JWT Token中解析（新方式，支持前端Bearer Token）
func (s *Service) Ctx(r *ghttp.Request) {
	customCtx := &bizctx.Context{
		Session: r.Session,
	}
	s.bizCtxSvc.Init(r, customCtx)

	// 优先从Session获取用户信息（原有方式）
	if user := s.sessionSvc.GetUser(r.Context()); user != nil {
		customCtx.User = &bizctx.User{
			Id:    user.Id,
			Name:  user.Name,
			Phone: user.Phone,
			Role:  user.Role,
		}
		r.Middleware.Next()
		return
	}

	// 从JWT Token获取用户信息（实验三新增）
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := s.jwtSvc.ParseToken(parts[1])
			if err == nil && claims != nil {
				// 从数据库获取最新用户信息
				user, _ := s.userSvc.GetUserById(r.Context(), claims.UserId)
				if user != nil {
					customCtx.User = &bizctx.User{
						Id:    user.Id,
						Name:  user.Name,
						Phone: user.Phone,
						Role:  user.Role,
					}
				}
			}
		}
	}

	r.Middleware.Next()
}

// Auth checks if user is signed in.
//Auth中间件——检查用户是否已登录
// 如果未登录，返回403错误
func (s *Service) Auth(r *ghttp.Request) {
	if !s.userSvc.IsSignedIn(r.Context()) {
		r.Response.WriteStatus(http.StatusForbidden, "请先登录")
		return
	}
	r.Middleware.Next()
}

// Admin checks if user is an admin.
//Admin中间件——检查用户是否为管理员
// 如果角色不是"admin"，返回403权限不足
func (s *Service) Admin(r *ghttp.Request) {
	customCtx := s.bizCtxSvc.Get(r.Context())
	if customCtx == nil || customCtx.User == nil || customCtx.User.Role != "admin" {
		r.Response.WriteStatus(http.StatusForbidden, "权限不足，需要管理员身份")
		return
	}
	r.Middleware.Next()
}



