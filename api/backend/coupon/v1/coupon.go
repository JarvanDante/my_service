// Package v1 后台优惠券契约(模板 CRUD + 发放 + 用户券查询)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`
	Scene     int     `json:"scene"`
	FaceValue float64 `json:"face_value"`
	Discount  int     `json:"discount"`
	Threshold float64 `json:"threshold"`
	MaxDeduct float64 `json:"max_deduct"`
	Total     int     `json:"total"`
	Issued    int     `json:"issued"`
	PerLimit  int     `json:"per_limit"`
	ExpireDay int     `json:"expire_day"`
	Status    int     `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type ListReq struct {
	g.Meta  `path:"/coupons" method:"get" tags:"Backend/Coupon" summary:"券模板列表"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta    `path:"/coupons" method:"post" tags:"Backend/Coupon" summary:"新增券模板"`
	Name      string  `json:"name" v:"required#券名必填"`
	Type      int     `json:"type" v:"in:1,2#券类型非法"`
	Scene     int     `json:"scene" v:"in:0,1,2,3#场景非法"`
	FaceValue float64 `json:"face_value"`
	Discount  int     `json:"discount"`
	Threshold float64 `json:"threshold"`
	MaxDeduct float64 `json:"max_deduct"`
	Total     int     `json:"total"`
	PerLimit  int     `json:"per_limit"`
	ExpireDay int     `json:"expire_day"`
	Status    int     `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta    `path:"/coupons/{id}" method:"put" tags:"Backend/Coupon" summary:"更新券模板"`
	Id        int64   `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`
	Scene     int     `json:"scene"`
	FaceValue float64 `json:"face_value"`
	Discount  int     `json:"discount"`
	Threshold float64 `json:"threshold"`
	MaxDeduct float64 `json:"max_deduct"`
	Total     int     `json:"total"`
	PerLimit  int     `json:"per_limit"`
	ExpireDay int     `json:"expire_day"`
	Status    int     `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/coupons/{id}" method:"delete" tags:"Backend/Coupon" summary:"删除券模板"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}

// GrantReq 定向发券(运营补偿/活动发放), 幂等由每人限领约束兜底。
type GrantReq struct {
	g.Meta  `path:"/coupons/{id}/grant" method:"post" tags:"Backend/Coupon" summary:"发放优惠券"`
	Id      int64   `json:"id" in:"path" v:"required|min:1#券模板ID必填"`
	UserIds []int64 `json:"user_ids" v:"required#用户ID列表必填"`
}
type GrantRes struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

type UserItem struct {
	Id         int64   `json:"id"`
	UserId     int64   `json:"user_id"`
	TplId      int64   `json:"tpl_id"`
	Name       string  `json:"name"`
	Type       int     `json:"type"`
	FaceValue  float64 `json:"face_value"`
	Discount   int     `json:"discount"`
	Status     int     `json:"status"`
	StatusText string  `json:"status_text"`
	RefId      string  `json:"ref_id"`
	ExpireAt   string  `json:"expire_at"`
	UsedAt     string  `json:"used_at"`
	CreatedAt  string  `json:"created_at"`
}

type UsersReq struct {
	g.Meta `path:"/coupons/users" method:"get" tags:"Backend/Coupon" summary:"用户券记录"`
	TplId  string `json:"tpl_id"`
	UserId string `json:"user_id"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}
type UsersRes struct {
	List  []UserItem `json:"list"`
	Total int        `json:"total"`
}
