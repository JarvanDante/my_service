// Package service AI 生成任务对外接口(前台提交/查询/取消 + 供应商回调 + 后台管理)。
package service

import "context"

// TemplateDTO AI 模板。
type TemplateDTO struct {
	Id        int64
	Name      string
	BizType   int
	Cover     string
	Preview   string
	Params    map[string]any
	CostGold   float64
	Sort       int
	Status     int
	UsageCount int
	CreatedAt  string
}

// TaskDTO AI 任务。Provider/ProviderTaskId 只给后台用, 前台控制器不下发。
type TaskDTO struct {
	Id             int64
	UserId         int64
	TaskNo         string
	ClientToken    string
	BizType        int
	TemplateId     int64
	Params         map[string]any
	InputUrl       string
	CostGold       float64
	Status         int
	Provider       string
	ProviderTaskId string
	Result         map[string]any
	ErrMsg         string
	RetryCount     int
	SubmittedAt    string
	FinishedAt     string
	CreatedAt      string
	Nickname       string
	Phone          string
	Avatar         string
	GroupName      string
	ChannelName    string
	DeviceType     string
	Sets           int
}

// TaskStats 后台订单汇总(按当前筛选, 不分页)。
type TaskStats struct {
	Total        int
	Success      int
	Refund       int
	Abnormal     int
	TotalGold    float64
	SuccessGold  float64
	RefundGold   float64
	AbnormalGold float64
}

// SubmitInput 提交入参。价格不在这里 —— 服务端按 template_id 自己算。
type SubmitInput struct {
	BizType     int
	TemplateId  int64
	Params      map[string]any
	InputUrl    string
	ClientToken string
}

// SubmitOutput 提交结果。Repeated=true 表示命中 client_token 幂等, 本次没有扣费。
type SubmitOutput struct {
	Task     *TaskDTO
	Balance  float64
	Repeated bool
}

// TemplateFilter 模板筛选。BizType/Status 为 -1 表示不筛。
type TemplateFilter struct {
	BizType int
	Status  int
	Keyword string
	Page    int
	Size    int
}

// TaskFilter 任务筛选。UserId/BizType/Status 为 0/-1 表示不筛。
type TaskFilter struct {
	UserId        int64
	BizType       int
	Status        int
	TaskNo        string
	Nickname      string
	ChannelName   string
	DeviceType    string
	StartTime     string
	EndTime       string
	RegisterStart string
	RegisterEnd   string
	Page          int
	Size          int
}

// TemplateInput 后台模板保存入参。
type TemplateInput struct {
	Id       int64
	Name     string
	BizType  int
	Cover    string
	Preview  string
	Params   map[string]any
	CostGold float64
	Sort     int
	Status   int
}

// CallbackInput 供应商回调入参(已由控制器解出, 验签在 logic 内做)。
type CallbackInput struct {
	Provider       string
	ProviderTaskId string
	Status         int
	Result         map[string]any
	ErrMsg         string
	Sign           string
}

type IAiTask interface {
	// ---- 前台 ----

	// FrontTemplates 启用中的模板(公开)。bizType<=0 表示全部玩法。
	FrontTemplates(ctx context.Context, bizType int) ([]*TemplateDTO, error)
	// Submit 提交任务: 事务内完成"校验+扣费+建单", 事务提交后才调用供应商;
	// 供应商提交失败会自动退款。client_token 命中则直接返回旧任务, 不二次扣费。
	Submit(ctx context.Context, userId int64, in SubmitInput) (*SubmitOutput, error)
	// Task 查单个任务(仅自己的)。处理中的任务会顺带向供应商轮询一次并落库(回调丢失的兜底)。
	Task(ctx context.Context, userId, id int64) (*TaskDTO, error)
	// Tasks 我的任务列表。
	Tasks(ctx context.Context, userId int64, f TaskFilter) ([]*TaskDTO, int, error)
	// Cancel 取消排队中的任务并退款; 已提交给供应商的任务不能取消。
	Cancel(ctx context.Context, userId, id int64) (refund, balance float64, err error)

	// ---- 回调 ----

	// Callback 供应商回调: 验签 + 幂等落终态; 失败状态在同一事务里自动退款。
	// 重复回调返回 nil(成功), 不重复处理也不重复退款。
	Callback(ctx context.Context, in CallbackInput) error
	HandleWorkerResult(ctx context.Context, jobID, status, outputURL, outputKey, errMsg string) error
	ExpireStale(ctx context.Context) (int, error)

	// ---- 后台 ----

	TemplateList(ctx context.Context, f TemplateFilter) ([]*TemplateDTO, int, error)
	TemplateCreate(ctx context.Context, in TemplateInput) (int64, error)
	TemplateUpdate(ctx context.Context, in TemplateInput) error
	TemplateDelete(ctx context.Context, id int64) error

	TaskList(ctx context.Context, f TaskFilter) ([]*TaskDTO, int, *TaskStats, error)
	// TaskRetry 重新提交任务(失败/已退款要重新扣费; 排队中的任务钱没退, 只重投不扣费), retry_count+1。
	TaskRetry(ctx context.Context, id int64) (*TaskDTO, error)
	// TaskRefund 人工退款: 仅非终态(排队中/处理中)或失败态可退, 条件更新防重复退款。
	TaskRefund(ctx context.Context, id int64, remark string) (float64, error)
	// TaskDelete 删除订单记录。处理中拒绝; 排队中先退款再删, 避免吞金币。
	TaskDelete(ctx context.Context, id int64) error
}
