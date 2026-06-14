// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package v1

import "github.com/gogf/gf/v2/frame/g"

// UpdatePhoneReq 更新手机号请求
type UpdatePhoneReq struct {
	g.Meta `path:"/users/me/phone" method:"put" tags:"用户服务" summary:"更新手机号"`
	Phone  string `v:"required" dc:"手机号"`
}

// UpdatePhoneRes 更新手机号响应
type UpdatePhoneRes struct {
	Message string `json:"message" dc:"提示信息"`
}
