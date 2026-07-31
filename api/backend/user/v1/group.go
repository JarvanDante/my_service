// Package v1 后台用户组配置接口契约(B4)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type UserGroupItem struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Rate   int    `json:"rate"`   // 折扣(%), 100=无折扣
	Rights string `json:"rights"` // 权益 JSON
	Remark string `json:"remark"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

// 用户组列表
type GroupListReq struct {
	g.Meta `path:"/user-groups" method:"get" tags:"Backend/UserGroup" summary:"用户组列表"`
}
type GroupListRes struct {
	List []UserGroupItem `json:"list"`
}

// 创建用户组
type GroupCreateReq struct {
	g.Meta `path:"/user-groups" method:"post" tags:"Backend/UserGroup" summary:"创建用户组"`
	Name   string `json:"name"   v:"required#组名必填"`
	Rate   int    `json:"rate"   v:"between:0,100#折扣须在0~100"`
	Rights string `json:"rights"` // JSON 字符串, 默认 {}
	Remark string `json:"remark"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type GroupCreateRes struct {
	Id int64 `json:"id"`
}

// 更新用户组(同步组内用户快照)
type GroupUpdateReq struct {
	g.Meta `path:"/user-groups/{id}" method:"put" tags:"Backend/UserGroup" summary:"更新用户组"`
	Id     int64  `json:"id"     v:"required|min:1#组ID必填|组ID必须大于0"`
	Name   string `json:"name"   v:"required#组名必填"`
	Rate   int    `json:"rate"   v:"between:0,100#折扣须在0~100"`
	Rights string `json:"rights"`
	Remark string `json:"remark"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type GroupUpdateRes struct{}

// 删除用户组
type GroupDeleteReq struct {
	g.Meta `path:"/user-groups/{id}" method:"delete" tags:"Backend/UserGroup" summary:"删除用户组"`
	Id     int64 `json:"id" v:"required|min:1#组ID必填|组ID必须大于0"`
}
type GroupDeleteRes struct{}
