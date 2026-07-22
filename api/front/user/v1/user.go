// Package v1 前台用户接口契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

// UserInfo 前台用户信息(登录/详情共用)。
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

// 设备登录/游客注册
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

// 当前用户详情(需登录, Header: Authorization)
type InfoReq struct {
	g.Meta `path:"/user/info" method:"get" tags:"Front/User" summary:"个人信息"`
}
type InfoRes struct {
	UserInfo
}
