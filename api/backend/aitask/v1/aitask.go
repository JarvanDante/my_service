// Package v1 后台 AI 模板/任务契约 + 供应商回调契约。
package v1

import "github.com/gogf/gf/v2/frame/g"

type TemplateItem struct {
	Id        int64          `json:"id"`
	Name      string         `json:"name"`
	BizType   int            `json:"biz_type"`
	Cover     string         `json:"cover"`
	Preview   string         `json:"preview"`
	Params    map[string]any `json:"params"`
	CostGold    float64        `json:"cost_gold"`
	Sort        int            `json:"sort"`
	Status      int            `json:"status"`
	UsageCount  int            `json:"usage_count"`
	CreatedAt   string         `json:"created_at"`
}

// TemplateListReq 模板列表。biz_type/status 都用 string 接收, 空=全部。
type TemplateListReq struct {
	g.Meta  `path:"/ai/templates" method:"get" tags:"Backend/AiTask" summary:"AI模板列表"`
	BizType string `json:"biz_type"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type TemplateListRes struct {
	List  []TemplateItem `json:"list"`
	Total int            `json:"total"`
}

type TemplateCreateReq struct {
	g.Meta   `path:"/ai/templates" method:"post" tags:"Backend/AiTask" summary:"新增AI模板"`
	Name     string         `json:"name" v:"required#模板名必填"`
	BizType  int            `json:"biz_type" v:"required|in:1,2,3,4,5,6#玩法类型必填|玩法类型非法"`
	Cover    string         `json:"cover"`
	Preview  string         `json:"preview"`
	Params   map[string]any `json:"params"`
	CostGold float64        `json:"cost_gold"`
	Sort     int            `json:"sort"`
	Status   int            `json:"status" v:"in:0,1#状态非法"`
}
type TemplateCreateRes struct {
	Id int64 `json:"id"`
}

type TemplateUpdateReq struct {
	g.Meta   `path:"/ai/templates/{id}" method:"put" tags:"Backend/AiTask" summary:"编辑AI模板"`
	Id       int64          `json:"id" in:"path" v:"required|min:1#ID必填"`
	Name     string         `json:"name"`
	BizType  int            `json:"biz_type" v:"in:0,1,2,3,4,5,6#玩法类型非法"`
	Cover    string         `json:"cover"`
	Preview  string         `json:"preview"`
	Params   map[string]any `json:"params"`
	CostGold float64        `json:"cost_gold"`
	Sort     int            `json:"sort"`
	Status   int            `json:"status" v:"in:0,1#状态非法"`
}
type TemplateUpdateRes struct{}

type TemplateDeleteReq struct {
	g.Meta `path:"/ai/templates/{id}" method:"delete" tags:"Backend/AiTask" summary:"删除AI模板"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type TemplateDeleteRes struct{}

// TaskItem 后台任务视图: 比前台多出 provider/provider_task_id, 用于人工排查与对账。
type TaskItem struct {
	Id             int64          `json:"id"`
	UserId         int64          `json:"user_id"`
	TaskNo         string         `json:"task_no"`
	ClientToken    string         `json:"client_token"`
	BizType        int            `json:"biz_type"`
	BizTypeText    string         `json:"biz_type_text"`
	TemplateId     int64          `json:"template_id"`
	Params         map[string]any `json:"params"`
	InputUrl       string         `json:"input_url"`
	CostGold       float64        `json:"cost_gold"`
	Status         int            `json:"status"`
	StatusText     string         `json:"status_text"`
	Provider       string         `json:"provider"`
	ProviderTaskId string         `json:"provider_task_id"`
	Result         map[string]any `json:"result"`
	ErrMsg         string         `json:"err_msg"`
	RetryCount     int            `json:"retry_count"`
	SubmittedAt    string         `json:"submitted_at"`
	FinishedAt     string         `json:"finished_at"`
	CreatedAt      string         `json:"created_at"`
	Nickname       string         `json:"nickname"`
	Phone          string         `json:"phone"`
	Avatar         string         `json:"avatar"`
	GroupName      string         `json:"group_name"`
	ChannelName    string         `json:"channel_name"`
	DeviceType     string         `json:"device_type"`
	Sets           int            `json:"sets"`
}

type TaskStats struct {
	Total        int     `json:"total"`
	Success      int     `json:"success"`
	Refund       int     `json:"refund"`
	Abnormal     int     `json:"abnormal"`
	TotalGold    float64 `json:"total_gold"`
	SuccessGold  float64 `json:"success_gold"`
	RefundGold   float64 `json:"refund_gold"`
	AbnormalGold float64 `json:"abnormal_gold"`
}

// TaskListReq 任务列表。user_id/biz_type/status 全部 string 接收, 空=全部。
type TaskListReq struct {
	g.Meta        `path:"/ai/tasks" method:"get" tags:"Backend/AiTask" summary:"AI任务列表"`
	UserId        string `json:"user_id"`
	BizType       string `json:"biz_type"`
	Status        string `json:"status"`
	TaskNo        string `json:"task_no"`
	Nickname      string `json:"nickname"`
	ChannelName   string `json:"channel_name"`
	DeviceType    string `json:"device_type"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	RegisterStart string `json:"register_start"`
	RegisterEnd   string `json:"register_end"`
	Page          int    `json:"page"`
	Size          int    `json:"size"`
}
type TaskListRes struct {
	List  []TaskItem `json:"list"`
	Total int        `json:"total"`
	Stats TaskStats  `json:"stats"`
}

// TaskRetryReq 重新提交(失败/已退款的任务)。重新扣费 + retry_count+1。
type TaskRetryReq struct {
	g.Meta `path:"/ai/tasks/{id}/retry" method:"post" tags:"Backend/AiTask" summary:"重新提交AI任务"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type TaskRetryRes struct {
	Task TaskItem `json:"task"`
}

// TaskRefundReq 人工退款(客服兜底)。只能退非终态或失败态, 已退款/已取消/已成功一律拒绝。
type TaskRefundReq struct {
	g.Meta `path:"/ai/tasks/{id}/refund" method:"post" tags:"Backend/AiTask" summary:"AI任务人工退款"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Remark string `json:"remark"`
}
type TaskRefundRes struct {
	Refund float64 `json:"refund"`
}

type TaskDeleteReq struct {
	g.Meta `path:"/ai/tasks/{id}" method:"delete" tags:"Backend/AiTask" summary:"删除AI订单"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type TaskDeleteRes struct{}

// CallbackReq 供应商回调(公开接口, 不走后台鉴权中间件, 靠签名自证身份)。
//
// sign = md5(provider_task_id + status + secret), secret 取 app_config 的 ai_callback_secret。
// 这是骨架阶段的最简验签形式; 真实供应商多用 HMAC-SHA256(raw body) + 时间戳防重放,
// 接入时按对方协议整体替换本字段与校验逻辑(见 shared/aiprovider 包注释)。
type CallbackReq struct {
	g.Meta         `path:"/ai/callback" method:"post" tags:"Backend/AiTask" summary:"AI供应商回调"`
	Provider       string         `json:"provider"`
	ProviderTaskId string         `json:"provider_task_id" v:"required#外部任务号必填"`
	Status         int            `json:"status" v:"required|in:3,4#回调状态必填|回调状态非法"`
	Result         map[string]any `json:"result"`
	ErrMsg         string         `json:"err_msg"`
	Sign           string         `json:"sign"`
}
type CallbackRes struct{}
