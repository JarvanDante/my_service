// Package v1 后台用户管理接口契约(B1)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// AdminUserItem 用户列表项。
type AdminUserItem struct {
	Id          int64   `json:"id"`
	Username    string  `json:"username"`
	Nickname    string  `json:"nickname"`
	Phone       string  `json:"phone"`
	Channel     string  `json:"channel"`
	GroupId     int64   `json:"group_id"`
	GroupName   string  `json:"group_name"`
	Level       int     `json:"level"`
	Balance     float64 `json:"balance"`
	Credit      float64 `json:"credit"`
	MoneyCount  float64 `json:"money_count"`
	IsDisabled  int     `json:"is_disabled"`
	RegisterAt  string  `json:"register_at"`
	LastLoginAt string  `json:"last_login_at"`
}

// 用户列表(筛选+分页)
type ListReq struct {
	g.Meta    `path:"/users" method:"get" tags:"Backend/User" summary:"用户列表(后台)"`
	Keyword   string `json:"keyword"`                            // 用户名/手机/昵称 模糊
	Channel   string `json:"channel"`                            // 渠道
	GroupId   int64  `json:"group_id"`                           // 用户组
	Status    int    `json:"status"     v:"in:0,1,2#状态仅支持0/1/2"` // 0全部 1正常 2禁用
	StartDate int    `json:"start_date"`                         // 注册日起 YYYYMMDD
	EndDate   int    `json:"end_date"`                           // 注册日止 YYYYMMDD
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type ListRes struct {
	List  []AdminUserItem `json:"list"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// 用户详情
type DetailReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"Backend/User" summary:"用户详情(后台)"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填|用户ID必须大于0"`
}
type DetailRes struct {
	AdminUserItem
	Sex          int    `json:"sex"`
	Signature    string `json:"signature"`
	Img          string `json:"img"`
	Fans         int    `json:"fans"`
	Follow       int    `json:"follow"`
	ShareNum     int    `json:"share_num"`
	ParentId     int64  `json:"parent_id"`
	ParentName   string `json:"parent_name"`
	GroupRate    int    `json:"group_rate"`
	GroupEndTime int64  `json:"group_end_time"`
	ErrorMsg     string `json:"error_msg"`
	RegisterIp   string `json:"register_ip"`
	LastIp       string `json:"last_ip"`
	LoginNum     int    `json:"login_num"`
}

// 禁用 / 解禁
type DisableReq struct {
	g.Meta `path:"/users/{id}/disable" method:"post" tags:"Backend/User" summary:"禁用/解禁用户(后台)"`
	Id     int64  `json:"id" v:"required|min:1#用户ID必填|用户ID必须大于0"`
	Op     string `json:"op" v:"required|in:disable,enable#操作必填|op 仅支持 disable/enable"`
	Reason string `json:"reason"` // 禁用原因(op=disable 时可填)
}
type DisableRes struct{}

// 调整用户组 / VIP 到期
type GroupReq struct {
	g.Meta       `path:"/users/{id}/group" method:"post" tags:"Backend/User" summary:"调整用户组(后台)"`
	Id           int64  `json:"id"         v:"required|min:1#用户ID必填|用户ID必须大于0"`
	GroupId      int64  `json:"group_id"   v:"min:0#用户组ID不合法"`
	GroupName    string `json:"group_name"`
	GroupRate    int    `json:"group_rate"     v:"min:0#折扣不合法"`
	GroupEndTime int64  `json:"group_end_time" v:"min:0#到期时间不合法"` // epoch 秒, 0 不限
}
type GroupRes struct{}

// 调整金币 / 积分
type BalanceReq struct {
	g.Meta `path:"/users/{id}/balance" method:"post" tags:"Backend/User" summary:"调整金币/积分(后台)"`
	Id     int64   `json:"id"     v:"required|min:1#用户ID必填|用户ID必须大于0"`
	Target string  `json:"target" v:"required|in:balance,credit#target必填|target 仅支持 balance/credit"`
	Amount float64 `json:"amount" v:"required#数额必填(正加负减)"`
	Remark string  `json:"remark"`
}
type BalanceRes struct{}

// 余额流水
type BalanceLogItem struct {
	Id            int64   `json:"id"`
	Direction     int     `json:"direction"` // 1收入 2支出
	Scene         string  `json:"scene"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	RefId         string  `json:"ref_id"`
	Remark        string  `json:"remark"`
	CreatedAt     string  `json:"created_at"`
}
type BalanceLogsReq struct {
	g.Meta `path:"/users/{id}/balance-logs" method:"get" tags:"Backend/User" summary:"用户余额流水(后台)"`
	Id     int64 `json:"id" v:"required|min:1#用户ID必填|用户ID必须大于0"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type BalanceLogsRes struct {
	List  []BalanceLogItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}
