// Package v1 后台用户管理接口契约(B1)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// AdminUserItem 用户列表项(对齐公司后台用户列表核心字段)。
type AdminUserItem struct {
	Id             int64   `json:"id"`
	Username       string  `json:"username"`
	Nickname       string  `json:"nickname"`
	Phone          string  `json:"phone"`
	Sex            int     `json:"sex"` // 0未知 1男 2女
	Tag            string  `json:"tag"` // JSON 数组字符串
	Img            string  `json:"img"`
	AccountSlat    string  `json:"account_slat"` // 登录二维码凭证 username==>md5
	Balance        float64 `json:"balance"`
	GiftCount      float64 `json:"gift_count"` // 收益
	Credit         float64 `json:"credit"`     // 积分
	MoneyCount     float64 `json:"money_count"`
	IsUp           int     `json:"is_up"`
	IsValid        int     `json:"is_valid"`
	HasBuy         int     `json:"has_buy"`
	Level          int     `json:"level"`
	GroupId        int64   `json:"group_id"`
	GroupName      string  `json:"group_name"`
	GroupRate      int     `json:"group_rate"`
	GroupStartTime int64   `json:"group_start_time"`
	GroupEndTime   int64   `json:"group_end_time"`
	ParentId       int64   `json:"parent_id"`
	ParentName     string  `json:"parent_name"`
	Channel        string  `json:"channel"`
	DeviceType     string  `json:"device_type"`
	DeviceExt      string  `json:"device_ext"`
	DeviceVersion  string  `json:"device_version"`
	MovieFeeRate   int     `json:"movie_fee_rate"`
	PostFeeRate    int     `json:"post_fee_rate"`
	ShareNum       int     `json:"share_num"`
	IsDisabled     int     `json:"is_disabled"`
	RegisterAt     string  `json:"register_at"`
	RegisterIp     string  `json:"register_ip"`
	RegisterArea   string  `json:"register_area"`
	LastLoginAt    string  `json:"last_login_at"`
	LastIp         string  `json:"last_ip"`
	LoginNum       int     `json:"login_num"`
}

// 用户列表(筛选+分页)
type ListReq struct {
	g.Meta      `path:"/users" method:"get" tags:"Backend/User" summary:"用户列表(后台)"`
	Keyword     string `json:"keyword"`                                // 用户名/手机/昵称 模糊
	UserId      int64  `json:"user_id"`                                // 精确用户ID
	Username    string `json:"username"`                               // 精确/模糊用户名
	Phone       string `json:"phone"`                                  // 手机号
	ParentId    int64  `json:"parent_id"`                              // 上级ID
	Channel     string `json:"channel"`                                // 渠道
	GroupId     int64  `json:"group_id"`                               // 用户组/等级
	IsUp        int    `json:"is_up"        v:"in:0,1,2#是否UP仅支持0/1/2"` // 0全部 1是 2否
	IsValid     int    `json:"is_valid"     v:"in:0,1,2#有效用户仅支持0/1/2"` // 0全部 1是 2否
	HasBuy      int    `json:"has_buy"      v:"in:0,1,2#是否购买仅支持0/1/2"` // 0全部 1是 2否
	Status      int    `json:"status"       v:"in:0,1,2#状态仅支持0/1/2"`   // 0全部 1正常 2禁用
	DeviceType  string `json:"device_type"`                            // 设备类型
	StartDate   int    `json:"start_date"`                             // 注册日起 YYYYMMDD
	EndDate     int    `json:"end_date"`                               // 注册日止 YYYYMMDD
	MinLoginNum int    `json:"min_login_num"`                          // 登录次数下限
	MaxLoginNum int    `json:"max_login_num"`                          // 登录次数上限
	Page        int    `json:"page"`
	Size        int    `json:"size"`
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
	Signature string `json:"signature"`
	Fans      int    `json:"fans"`
	Follow    int    `json:"follow"`
	ErrorMsg  string `json:"error_msg"`
}

// 禁用 / 解禁
type DisableReq struct {
	g.Meta `path:"/users/{id}/disable" method:"post" tags:"Backend/User" summary:"禁用/解禁用户(后台)"`
	Id     int64  `json:"id" v:"required|min:1#用户ID必填|用户ID必须大于0"`
	Op     string `json:"op" v:"required|in:disable,enable#操作必填|op 仅支持 disable/enable"`
	Reason string `json:"reason"` // 禁用原因(op=disable 时可填)
}
type DisableRes struct{}

// 批量冻结 / 解冻(对齐公司后台用户列表的批量操作)
type BatchDisableReq struct {
	g.Meta `path:"/users/batch-disable" method:"post" tags:"Backend/User" summary:"批量冻结/解冻用户(后台)"`
	Ids    []int64 `json:"ids" v:"required#用户ID列表必填"`
	Op     string  `json:"op"  v:"required|in:disable,enable#操作必填|op 仅支持 disable/enable"`
	Reason string  `json:"reason"` // 冻结原因(op=disable 时可填, 前台登录时透出)
}
type BatchDisableRes struct {
	Affected int `json:"affected"` // 实际变更行数
}

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
