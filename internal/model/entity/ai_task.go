// Code maintained manually (AI 生成模板 + AI 生成任务).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// AI 玩法类型(biz_type): 一张表承载 tianbi 原来的 6 套 AI 玩法, 差异只在 params/result。
const (
	AiBizFaceSwap     = 1 // 换脸
	AiBizUndress      = 2 // 脱衣
	AiBizTextToImage  = 3 // 文生图
	AiBizImageToVideo = 4 // 图生视频
	AiBizTextToNovel  = 5 // 文生小说
	AiBizChat         = 6 // AI对话(女友/客服)
)

// AI 任务状态。失败(4)与已退款(5)必须分开: 退款是一次独立的资金动作,
// 状态本身要能回答"这笔钱退了没有", 否则重复回调时无法判断是否已经补偿过。
const (
	AiStatusQueued    = 1 // 排队中(已扣费, 尚未提交给供应商)
	AiStatusRunning   = 2 // 处理中(供应商已受理)
	AiStatusSucceed   = 3 // 成功(终态)
	AiStatusFailed    = 4 // 失败(终态; 仅在不需要退款的场景出现, 正常失败会直接流转到 5)
	AiStatusRefunded  = 5 // 已退款(终态)
	AiStatusCancelled = 6 // 已取消(终态, 取消时已退款)
)

// AiTemplate AI 生成模板: 玩法预设参数 + 定价。价格只认这张表, 绝不接受客户端传值。
type AiTemplate struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	Name      string      `json:"name"      orm:"name"`
	BizType   int         `json:"bizType"   orm:"biz_type"`
	Cover     string      `json:"cover"     orm:"cover"`
	Preview   string      `json:"preview"   orm:"preview"`
	Params    string      `json:"params"    orm:"params"` // jsonb 原文
	CostGold  float64     `json:"costGold"  orm:"cost_gold"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

// AiTask AI 生成任务。CostGold 是下单时实扣的金额, 退款一律按它退,
// 模板改价不影响历史任务的退款金额。
type AiTask struct {
	Id             int64       `json:"id"             orm:"id"`
	SiteId         int64       `json:"siteId"         orm:"site_id"`
	UserId         int64       `json:"userId"         orm:"user_id"`
	BizType        int         `json:"bizType"        orm:"biz_type"`
	TemplateId     int64       `json:"templateId"     orm:"template_id"`
	TaskNo         string      `json:"taskNo"         orm:"task_no"`
	ClientToken    string      `json:"clientToken"    orm:"client_token"`
	Params         string      `json:"params"         orm:"params"` // jsonb 原文
	InputUrl       string      `json:"inputUrl"       orm:"input_url"`
	CostGold       float64     `json:"costGold"       orm:"cost_gold"`
	Status         int         `json:"status"         orm:"status"`
	Provider       string      `json:"provider"       orm:"provider"`
	ProviderTaskId string      `json:"providerTaskId" orm:"provider_task_id"`
	Result         string      `json:"result"         orm:"result"` // jsonb 原文
	ErrMsg         string      `json:"errMsg"         orm:"err_msg"`
	RetryCount     int         `json:"retryCount"     orm:"retry_count"`
	SubmittedAt    *gtime.Time `json:"submittedAt"    orm:"submitted_at"`
	FinishedAt     *gtime.Time `json:"finishedAt"     orm:"finished_at"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"`
}
