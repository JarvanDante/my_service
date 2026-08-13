// Package v1 前台 AI 生成任务契约(换脸/脱衣/文生图/图生视频/文生小说/AI对话统一入口)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// TemplateItem 模板(含价格): 前台按这里的 cost_gold 展示"本次消耗", 提交时价格仍由服务端重算。
type TemplateItem struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	BizType  int            `json:"biz_type"`
	Cover    string         `json:"cover"`
	Preview  string         `json:"preview"`
	Params   map[string]any `json:"params"`
	CostGold float64        `json:"cost_gold"`
	Sort     int            `json:"sort"`
}

// TemplatesReq 模板列表(公开)。biz_type=0 表示全部玩法。
type TemplatesReq struct {
	g.Meta  `path:"/ai/templates" method:"get" tags:"Front/AiTask" summary:"AI模板列表"`
	BizType int `json:"biz_type"`
}
type TemplatesRes struct {
	List []TemplateItem `json:"list"`
}

// TaskItem 任务对外视图。不下发 provider/provider_task_id: 那是我们与供应商之间的内部标识,
// 泄漏出去等于把回调伪造的一半材料白送给攻击者。
type TaskItem struct {
	Id          int64          `json:"id"`
	TaskNo      string         `json:"task_no"`
	BizType     int            `json:"biz_type"`
	TemplateId  int64          `json:"template_id"`
	Params      map[string]any `json:"params"`
	InputUrl    string         `json:"input_url"`
	CostGold    float64        `json:"cost_gold"`
	Status      int            `json:"status"` // 1排队中 2处理中 3成功 4失败 5已退款 6已取消
	Result      map[string]any `json:"result"`
	ErrMsg      string         `json:"err_msg"`
	RetryCount  int            `json:"retry_count"`
	SubmittedAt string         `json:"submitted_at"`
	FinishedAt  string         `json:"finished_at"`
	CreatedAt   string         `json:"created_at"`
}

// SubmitReq 提交生成任务(需登录)。
//
// 注意这里**没有金额字段**: 价格一律取模板 cost_gold, 无模板时取配置 ai_default_cost,
// 客户端传什么都不作数(tianbi 原版把价格放在请求里, 改个包就能白嫖)。
//
// client_token 由客户端生成(建议 uuid), 同一个 token 重复提交只会有一个任务、只扣一次费;
// 断网重试时请沿用同一个 token, 不要每次重发都换一个。
type SubmitReq struct {
	g.Meta      `path:"/ai/submit" method:"post" tags:"Front/AiTask" summary:"提交AI生成任务"`
	BizType     int            `json:"biz_type" v:"required|in:1,2,3,4,5,6#玩法类型必填|玩法类型非法"`
	TemplateId  int64          `json:"template_id"`
	Params      map[string]any `json:"params"`
	InputUrl    string         `json:"input_url"`
	ClientToken string         `json:"client_token"`
}
type SubmitRes struct {
	Task     TaskItem `json:"task"`
	Balance  float64  `json:"balance"`  // 扣费后余额
	Repeated bool     `json:"repeated"` // true=命中 client_token 幂等, 返回的是已有任务, 本次没有扣费
}

// TaskReq 查单个任务(需登录, 只能查自己的)。
type TaskReq struct {
	g.Meta `path:"/ai/task" method:"get" tags:"Front/AiTask" summary:"AI任务详情"`
	Id     int64 `json:"id" v:"required|min:1#任务ID必填"`
}
type TaskRes struct {
	Task TaskItem `json:"task"`
}

// TasksReq 我的任务列表(需登录)。status 用 string 接收: 空串=全部, "0" 也是合法筛选值,
// 用 int 的话零值与"没传"分不开。
type TasksReq struct {
	g.Meta  `path:"/ai/tasks" method:"get" tags:"Front/AiTask" summary:"我的AI任务列表"`
	BizType string `json:"biz_type"`
	Status  string `json:"status"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type TasksRes struct {
	List  []TaskItem `json:"list"`
	Total int        `json:"total"`
}

// CancelReq 取消任务(需登录)。只有还在排队中(未提交给供应商)的任务能取消, 取消即退款。
type CancelReq struct {
	g.Meta `path:"/ai/cancel" method:"post" tags:"Front/AiTask" summary:"取消AI任务并退款"`
	Id     int64 `json:"id" v:"required|min:1#任务ID必填"`
}
type CancelRes struct {
	Refund  float64 `json:"refund"`
	Balance float64 `json:"balance"`
}
