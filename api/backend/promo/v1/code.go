// Package v1 后台兑换码管理接口契约(B3)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type CodeItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	CodeKey   string `json:"code_key"`
	Type      string `json:"type"` // point / group
	ObjectId  int64  `json:"object_id"`
	AddNum    int    `json:"add_num"`
	CanUseNum int    `json:"can_use_num"`
	UsedNum   int    `json:"used_num"`
	Status    int    `json:"status"` // 0未使用 1已使用 -1作废
	ExpiredAt int64  `json:"expired_at"`
	CreatedAt string `json:"created_at"`
}

// 兑换码列表
type CodeListReq struct {
	g.Meta  `path:"/codes" method:"get" tags:"Backend/Promo" summary:"兑换码列表"`
	Keyword string `json:"keyword"`  // 码/名称模糊
	CodeKey string `json:"code_key"` // 批次
	Type    string `json:"type"     v:"in:,point,group#type 仅支持 point/group"`
	Status  int    `json:"status"   v:"in:0,1,2,3#状态仅支持0/1/2/3"` // 0全部 1可用 2已使用 3作废
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type CodeListRes struct {
	List  []CodeItem `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// 批量生成
type CodeGenReq struct {
	g.Meta    `path:"/codes" method:"post" tags:"Backend/Promo" summary:"批量生成兑换码"`
	Name      string `json:"name"        v:"required#名称必填"`
	Type      string `json:"type"        v:"required|in:point,group#type必填|type 仅支持 point/group"`
	ObjectId  int64  `json:"object_id"   v:"min:0#object_id 不合法"`         // type=group 时为用户组ID
	AddNum    int    `json:"add_num"     v:"required|min:1#面额必填|面额必须大于0"` // 金币数/天数
	CanUseNum int    `json:"can_use_num" v:"min:0#可用次数不合法"`               // 默认1
	Count     int    `json:"count"       v:"required|between:1,1000#数量必填|数量须在1~1000"`
	ExpiredAt int64  `json:"expired_at"  v:"min:0#过期时间不合法"` // epoch秒, 0不过期
}
type CodeGenRes struct {
	CodeKey string   `json:"code_key"`
	Count   int      `json:"count"`
	Codes   []string `json:"codes"`
}

// 作废
type CodeVoidReq struct {
	g.Meta `path:"/codes/{id}/void" method:"post" tags:"Backend/Promo" summary:"作废兑换码"`
	Id     int64 `json:"id" v:"required|min:1#兑换码ID必填|兑换码ID必须大于0"`
}
type CodeVoidRes struct{}

// 兑换记录
type CodeLogItem struct {
	Id        int64  `json:"id"`
	CodeId    int64  `json:"code_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	AddNum    int    `json:"add_num"`
	CreatedAt string `json:"created_at"`
}
type CodeLogListReq struct {
	g.Meta `path:"/code-logs" method:"get" tags:"Backend/Promo" summary:"兑换记录"`
	CodeId int64  `json:"code_id" v:"min:0#code_id 不合法"`
	UserId int64  `json:"user_id" v:"min:0#user_id 不合法"`
	Code   string `json:"code"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type CodeLogListRes struct {
	List  []CodeLogItem `json:"list"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}
