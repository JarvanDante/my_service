// Package v1 前台用户接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type UserInfo struct {
	Id        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Phone     string  `json:"phone"`
	Img       string  `json:"img"`
	Signature string  `json:"signature"`
	Sex       int     `json:"sex"`
	Level     int     `json:"level"`
	Balance   float64 `json:"balance"`
	Credit    float64 `json:"credit"`
	GroupName string  `json:"group_name"`
	Fans      int     `json:"fans"`
	Follow    int     `json:"follow"`
}

// 设备登录(公开)
type LoginReq struct {
	g.Meta        `path:"/user/login" method:"post" tags:"Front/User" summary:"设备登录"`
	DeviceId      string `json:"device_id"      v:"required#设备号必填"`
	DeviceType    string `json:"device_type"`
	DeviceVersion string `json:"device_version"`
}
type LoginRes struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// 个人信息(需登录)
type InfoReq struct {
	g.Meta `path:"/user/info" method:"get" tags:"Front/User" summary:"个人信息"`
}
type InfoRes struct {
	UserInfo
}

// 退出登录(需登录)
type LogoutReq struct {
	g.Meta `path:"/user/logout" method:"post" tags:"Front/User" summary:"退出登录"`
}
type LogoutRes struct{}

// 刷新 token(需登录)
type RefreshReq struct {
	g.Meta `path:"/user/token/refresh" method:"post" tags:"Front/User" summary:"刷新token"`
}
type RefreshRes struct {
	Token string `json:"token"`
}

// 绑定手机(需登录)
type BindPhoneReq struct {
	g.Meta `path:"/user/bind-phone" method:"post" tags:"Front/User" summary:"绑定手机"`
	Phone  string `json:"phone" v:"required|phone#手机号必填|手机号格式不正确"`
	Code   string `json:"code"` // 短信验证码(接入短信服务后校验)
}
type BindPhoneRes struct{}
