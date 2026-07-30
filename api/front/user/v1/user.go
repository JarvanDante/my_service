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

// ---- P2 个人资料 ----

// PublicUser 对外公开信息(看他人)。
type PublicUser struct {
	Id        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Img       string `json:"img"`
	BgImg     string `json:"bg_img"`
	Signature string `json:"signature"`
	Sex       int    `json:"sex"`
	Level     int    `json:"level"`
	Fans      int    `json:"fans"`
	Follow    int    `json:"follow"`
	ShareNum  int    `json:"share_num"`
}

// 他人主页(需登录)
type HomeReq struct {
	g.Meta `path:"/user/home/{id}" method:"get" tags:"Front/User" summary:"他人主页"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填"`
}
type HomeRes struct {
	User       PublicUser `json:"user"`
	IsFollowed bool       `json:"is_followed"`
}

// 修改资料(需登录, 空字段表示不改)
type UpdateReq struct {
	g.Meta    `path:"/user/update" method:"post" tags:"Front/User" summary:"修改资料"`
	Nickname  string `json:"nickname"`
	Img       string `json:"img"`
	BgImg     string `json:"bg_img"`
	Signature string `json:"signature"`
	Sex       int    `json:"sex" v:"in:0,1,2#性别不合法"`
}
type UpdateRes struct{}

// 用户图片(需登录)
type ImagesReq struct {
	g.Meta `path:"/user/images" method:"get" tags:"Front/User" summary:"用户图片"`
}
type ImagesRes struct {
	Images []string `json:"images"`
}

// 按账号查找(需登录)
type FindReq struct {
	g.Meta  `path:"/user/find" method:"get" tags:"Front/User" summary:"按账号查找用户"`
	Account string `json:"account" v:"required#账号必填"`
}
type FindRes struct {
	Found bool        `json:"found"`
	User  *PublicUser `json:"user"`
}
