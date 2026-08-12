// Package v1 后台兑换码契约(建码/管理/使用记录)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	CardType   int    `json:"card_type"`
	Value      int    `json:"value"`
	TotalTimes int    `json:"total_times"`
	UsedTimes  int    `json:"used_times"`
	Status     int    `json:"status"`
	ExpiredAt  string `json:"expired_at"`
	CreatedAt  string `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/redeemcode" method:"get" tags:"Backend/RedeemCode" summary:"兑换码列表"`
	Status  string `json:"status"`  // 空=全部  0=禁用  1=启用
	Keyword string `json:"keyword"` // 码/名称模糊搜索
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta     `path:"/redeemcode" method:"post" tags:"Backend/RedeemCode" summary:"新增兑换码(code 留空自动生成)"`
	Name       string `json:"name" v:"required#名称必填"`
	Code       string `json:"code"`                                       // 留空自动生成 12 位大写码
	Value      int    `json:"value" v:"required|min:1#金币数必填|金币数需大于0"`     // 金币数
	TotalTimes int    `json:"total_times" v:"required|min:1#次数必填|次数至少为1"` // 总可兑换次数
	ExpiredAt  string `json:"expired_at" v:"required#过期时间必填"`             // 如 2027-12-31 23:59:59
	Status     int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
}

type UpdateReq struct {
	g.Meta     `path:"/redeemcode/{id}" method:"put" tags:"Backend/RedeemCode" summary:"更新兑换码"`
	Id         int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	TotalTimes int    `json:"total_times"`
	ExpiredAt  string `json:"expired_at"`
	Status     int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/redeemcode/{id}" method:"delete" tags:"Backend/RedeemCode" summary:"删除兑换码"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

type RecordItem struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	CodeId    int64  `json:"code_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	CardType  int    `json:"card_type"`
	Value     int    `json:"value"`
	CreatedAt string `json:"created_at"`
}

type RecordListReq struct {
	g.Meta `path:"/redeemcode/records" method:"get" tags:"Backend/RedeemCode" summary:"兑换使用记录"`
	UserId int64  `json:"user_id"` // 0=全部
	Code   string `json:"code"`    // 空=全部
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type RecordListRes struct {
	List  []RecordItem `json:"list"`
	Total int          `json:"total"`
}
