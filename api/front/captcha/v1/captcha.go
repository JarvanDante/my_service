package v1

import "github.com/gogf/gf/v2/frame/g"

type IssueReq struct {
	g.Meta `path:"/captcha" method:"get" tags:"Front/Captcha" summary:"开屏图形验证码"`
}
type IssueRes struct {
	Id    string `json:"id"`
	Image string `json:"image"`
}

type VerifyReq struct {
	g.Meta `path:"/captcha/verify" method:"post" tags:"Front/Captcha" summary:"校验开屏验证码"`
	Id     string `json:"id"   v:"required#验证码已过期，请刷新"`
	Code   string `json:"code" v:"required#请输入验证码"`
}
type VerifyRes struct {
	Ok bool `json:"ok"`
}
