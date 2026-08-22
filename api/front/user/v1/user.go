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
	Follow      int  `json:"follow"`
	HasPassword bool `json:"has_password"`
	// Ext 站点差异字段(由 Nacos response.user_info_extra 白名单控制)
	Ext map[string]interface{} `json:"ext,omitempty"`
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

// 找回账号(公开): 用账号凭证(登录二维码内容 username==>md5)把账号换绑到当前设备并登录。
// 场景: 用户卸载重装/换机后, 扫码或粘贴之前保存的凭证恢复原账号。
type RestoreReq struct {
	g.Meta        `path:"/user/restore" method:"post" tags:"Front/User" summary:"凭证找回账号"`
	Credential    string `json:"credential"  v:"required#账号凭证必填"`
	DeviceId      string `json:"device_id"   v:"required#设备号必填"`
	DeviceType    string `json:"device_type"`
	DeviceVersion string `json:"device_version"`
}
type RestoreRes struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// AccountLoginReq 用户名+密码登录(公开)。用户名即「我的」页编号。
type AccountLoginReq struct {
	g.Meta        `path:"/user/account-login" method:"post" tags:"Front/User" summary:"账密登录"`
	Username      string `json:"username"       v:"required#用户名必填"`
	Password      string `json:"password"       v:"required#密码必填"`
	DeviceId      string `json:"device_id"      v:"required#设备号必填"`
	DeviceType    string `json:"device_type"`
	DeviceVersion string `json:"device_version"`
}
type AccountLoginRes struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// SetPasswordReq 设置或修改密码(需登录)。首次设密不用旧密码。
type SetPasswordReq struct {
	g.Meta      `path:"/user/password" method:"post" tags:"Front/User" summary:"设置密码"`
	OldPassword string `json:"old_password"`
	Password    string `json:"password" v:"required|length:6,32#新密码必填|密码需6-32位"`
}
type SetPasswordRes struct{}

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

// 默认头像列表(需登录, 改资料时供用户挑选)
type AvatarsReq struct {
	g.Meta `path:"/user/avatars" method:"get" tags:"Front/User" summary:"默认头像列表"`
}
type AvatarsRes struct {
	List []string `json:"list"`
}

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

// ---- P3 社交 ----

// 关注/取关(切换, 需登录)
type FollowReq struct {
	g.Meta `path:"/user/follow" method:"post" tags:"Front/User" summary:"关注/取关(切换)"`
	HomeId int64 `json:"home_id" v:"required|min:1#目标用户ID必填"`
}
type FollowRes struct {
	Followed bool `json:"followed"` // 操作后是否已关注
}

// 我的关注列表(需登录)
type FollowsReq struct {
	g.Meta `path:"/user/follows" method:"get" tags:"Front/User" summary:"我的关注列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type FollowsRes struct {
	List  []PublicUser `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// 我的粉丝列表(需登录)
type FansReq struct {
	g.Meta `path:"/user/fans" method:"get" tags:"Front/User" summary:"我的粉丝列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type FansRes struct {
	List  []PublicUser `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// ---- P4 推广/兑换码 ----

// 绑定推荐人(by 账号)
type BindParentReq struct {
	g.Meta  `path:"/user/bind-parent" method:"post" tags:"Front/User" summary:"绑定推荐人"`
	Account string `json:"account" v:"required#推荐人账号必填"`
}
type BindParentRes struct{}

// 绑定邀请码(by 邀请码=分享码)
type BindCodeReq struct {
	g.Meta `path:"/user/bind-code" method:"post" tags:"Front/User" summary:"绑定邀请码"`
	Code   string `json:"code" v:"required#邀请码必填"`
}
type BindCodeRes struct{}

// 使用兑换码
type RedeemReq struct {
	g.Meta `path:"/user/code/redeem" method:"post" tags:"Front/User" summary:"使用兑换码"`
	Code   string `json:"code" v:"required#兑换码必填"`
}
type RedeemRes struct {
	Type   string `json:"type"`
	AddNum int    `json:"add_num"`
	Name   string `json:"name"`
}

// 兑换记录
type CodeLog struct {
	Id        int64  `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	AddNum    int    `json:"add_num"`
	CreatedAt string `json:"created_at"`
}
type CodeLogsReq struct {
	g.Meta `path:"/user/code/logs" method:"get" tags:"Front/User" summary:"兑换记录"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type CodeLogsRes struct {
	List  []CodeLog `json:"list"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

// 分享信息
type ShareReq struct {
	g.Meta `path:"/user/share" method:"get" tags:"Front/User" summary:"分享信息"`
}
type ShareRes struct {
	ShareCode string `json:"share_code"`
	ShareUrl  string `json:"share_url"`
	ShareNum  int    `json:"share_num"`
}

// 分享记录(暂空)
type ShareLog struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	TargetId  int64  `json:"target_id"`
	Channel   string `json:"channel"`
	CreatedAt string `json:"created_at"`
}

// 上报分享
type ShareReportReq struct {
	g.Meta   `path:"/user/share/report" method:"post" tags:"Front/User" summary:"上报分享"`
	Type     string `json:"type"`      // app/poster/link/invite...
	TargetId int64  `json:"target_id"` // 分享对象(可选)
	Channel  string `json:"channel"`   // 渠道(可选)
}
type ShareReportRes struct{}
type ShareLogsReq struct {
	g.Meta `path:"/user/share/logs" method:"get" tags:"Front/User" summary:"分享记录"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type ShareLogsRes struct {
	List  []ShareLog `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type Invitee struct {
	Nickname   string `json:"nickname"`
	InviteCode string `json:"invite_code"`
	CreatedAt  string `json:"created_at"`
}

type InviteesReq struct {
	g.Meta `path:"/user/share/invitees" method:"get" tags:"Front/User" summary:"我邀请的用户"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type InviteesRes struct {
	List  []Invitee `json:"list"`
	Total int       `json:"total"`
}

// ---- P5 成长(签到/任务) ----

// 每日签到
type SignReq struct {
	g.Meta `path:"/user/sign" method:"post" tags:"Front/User" summary:"每日签到"`
}
type SignRes struct {
	Today  int     `json:"today"`
	Days   []int   `json:"days"`
	Count  int     `json:"count"`
	Reward float64 `json:"reward"`
}

// 任务列表
type Task struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	MaxNum      int    `json:"max_num"`
	DoneToday   int    `json:"done_today"`
}
type TasksReq struct {
	g.Meta `path:"/user/tasks" method:"get" tags:"Front/User" summary:"任务列表"`
}
type TasksRes struct {
	List []Task `json:"list"`
}

// 完成任务
type TaskDoReq struct {
	g.Meta `path:"/user/task/do" method:"post" tags:"Front/User" summary:"完成任务领奖"`
	TaskId int64 `json:"task_id" v:"required|min:1#任务ID必填"`
}
type TaskDoRes struct {
	DoneToday int     `json:"done_today"`
	MaxNum    int     `json:"max_num"`
	Reward    float64 `json:"reward"`
}

// 任务记录
type TaskLog struct {
	Id        int64  `json:"id"`
	TaskId    int64  `json:"task_id"`
	Type      string `json:"type"`
	Num       int    `json:"num"`
	LogDate   int    `json:"log_date"`
	CreatedAt string `json:"created_at"`
}
type TaskLogsReq struct {
	g.Meta `path:"/user/task/logs" method:"get" tags:"Front/User" summary:"任务记录"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type TaskLogsRes struct {
	List  []TaskLog `json:"list"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

// 成长信息
type UpReq struct {
	g.Meta `path:"/user/up" method:"get" tags:"Front/User" summary:"成长信息"`
}
type UpRes struct {
	Level             int     `json:"level"`
	Credit            float64 `json:"credit"`
	Balance           float64 `json:"balance"`
	SignDaysThisMonth int     `json:"sign_days_this_month"`
	Fans              int     `json:"fans"`
	Follow            int     `json:"follow"`
	ShareNum          int     `json:"share_num"`
}

// ---- P6 资产(充值/VIP/兑换) ----

type RechargePackage struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Coin   float64 `json:"coin"`
	Bonus  float64 `json:"bonus"`
}
type RechargeReq struct {
	g.Meta `path:"/user/recharge" method:"get" tags:"Front/User" summary:"充值套餐"`
}
type RechargeRes struct {
	List []RechargePackage `json:"list"`
}

type RechargeDoReq struct {
	g.Meta    `path:"/user/recharge/do" method:"post" tags:"Front/User" summary:"发起充值"`
	PackageId int64 `json:"package_id" v:"required|min:1#套餐ID必填"`
}
type RechargeDoRes struct {
	OrderNo string  `json:"order_no"`
	Amount  float64 `json:"amount"`
	Coin    float64 `json:"coin"`
}

type RechargeMockPayReq struct {
	g.Meta  `path:"/user/recharge/mock-pay" method:"post" tags:"Front/User" summary:"开发环境 mock 支付到账"`
	OrderNo string `json:"order_no" v:"required#订单号必填"`
}
type RechargeMockPayRes struct{}

type VipPackage struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Days  int     `json:"days"`
	Price float64 `json:"price"`
}
type VipReq struct {
	g.Meta `path:"/user/vip" method:"get" tags:"Front/User" summary:"VIP套餐"`
}
type VipRes struct {
	List []VipPackage `json:"list"`
}

type VipDoReq struct {
	g.Meta    `path:"/user/vip/do" method:"post" tags:"Front/User" summary:"开通/续费VIP"`
	PackageId int64 `json:"package_id" v:"required|min:1#套餐ID必填"`
}
type VipDoRes struct{}

type VipLog struct {
	Id        int64   `json:"id"`
	PackageId int64   `json:"package_id"`
	Days      int     `json:"days"`
	Price     float64 `json:"price"`
	StartAt   int64   `json:"start_at"`
	EndAt     int64   `json:"end_at"`
	CreatedAt string  `json:"created_at"`
}
type VipLogsReq struct {
	g.Meta `path:"/user/vip/logs" method:"get" tags:"Front/User" summary:"VIP记录"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type VipLogsRes struct {
	List  []VipLog `json:"list"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Size  int      `json:"size"`
}

type ExchangeReq struct {
	g.Meta `path:"/user/exchange" method:"get" tags:"Front/User" summary:"兑换信息"`
}
type ExchangeRes struct {
	Rate    int     `json:"rate"` // 多少积分兑 1 金币
	Credit  float64 `json:"credit"`
	Balance float64 `json:"balance"`
}

type ExchangeDoReq struct {
	g.Meta `path:"/user/exchange/do" method:"post" tags:"Front/User" summary:"积分兑换金币"`
	Coin   int `json:"coin" v:"required|min:1#兑换金币数必填"`
}
type ExchangeDoRes struct{}

// ---- P7 私信 ----

type Chat struct {
	PeerId      int64  `json:"peer_id"`
	Nickname    string `json:"nickname"`
	Img         string `json:"img"`
	LastContent string `json:"last_content"`
	LastAt      string `json:"last_at"`
	Unread      int    `json:"unread"`
}
type ChatsReq struct {
	g.Meta `path:"/user/chats" method:"get" tags:"Front/User" summary:"会话列表"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type ChatsRes struct {
	List  []Chat `json:"list"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type Message struct {
	Id        int64  `json:"id"`
	FromId    int64  `json:"from_id"`
	ToId      int64  `json:"to_id"`
	Content   string `json:"content"`
	Mine      bool   `json:"mine"`
	CreatedAt string `json:"created_at"`
}
type ChatMessagesReq struct {
	g.Meta `path:"/user/chat/messages" method:"get" tags:"Front/User" summary:"会话消息"`
	PeerId int64 `json:"peer_id" v:"required|min:1#对方ID必填"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}
type ChatMessagesRes struct {
	List  []Message `json:"list"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

type ChatSendReq struct {
	g.Meta  `path:"/user/chat/send" method:"post" tags:"Front/User" summary:"发消息"`
	ToId    int64  `json:"to_id" v:"required|min:1#接收人ID必填"`
	Content string `json:"content" v:"required#消息内容必填"`
}
type ChatSendRes struct{}

type ChatDelReq struct {
	g.Meta `path:"/user/chat/del" method:"post" tags:"Front/User" summary:"删除会话"`
	PeerId int64 `json:"peer_id" v:"required|min:1#对方ID必填"`
}
type ChatDelRes struct{}

type CustomerUrlReq struct {
	g.Meta `path:"/user/customer-url" method:"get" tags:"Front/User" summary:"客服链接"`
}
type CustomerUrlRes struct {
	Url string `json:"url"`
}
