// Package aiprovider 是 AI 生成任务的「供应商占位层」。
//
// ============================ 这个包为什么存在 ============================
//
// 业务上 6 种 AI 玩法(换脸/脱衣/文生图/图生视频/文生小说/AI对话)形态完全一致:
// 提交任务 → 第三方异步跑 → 回调或轮询拿结果 → 计费与失败退款。真正会变的只有
// 「谁来跑这个任务」。而供应商此刻还没定 —— 所以本包只定义与供应商无关的接口,
// **不包含任何真实第三方的域名、密钥、协议字段、错误码**。
//
// 接入真实供应商时的工作量被限制在一个文件里:
//  1. 新写 internal/shared/aiprovider/xxx.go, 实现 Provider 接口;
//  2. 在它的 init() 里 Register("xxx", &xxxProvider{});
//  3. 把 app_config 的 ai_provider 改成 "xxx"。
//
// 自建工人用 adapter 名 "local": Submit 发 Kafka(换脸/去衣), 结果由 my_ai_worker 回写。
//
// 业务层(internal/modules/aitask)零改动: 它只认 Provider 接口和 ai_task 的状态码,
// 不知道也不需要知道对面是谁。这是把"供应商未定"这件事变成可交付代码的关键。
//
// ======================== 接真实供应商时必须补的东西 ========================
//
// 下面这些在 mock 里全都不需要, 但真实供应商一个都不能少, 否则线上一定出事:
//
//   - 鉴权: API Key / 签名 / OAuth token 的获取与刷新。密钥走配置或密钥管理服务读取,
//     绝不硬编码进代码, 也不要写进本仓库的任何文件。
//   - 超时: 每个 HTTP 调用必须带 context 超时(建议 Submit 5~10s, Query 3~5s)。
//     没有超时的外部调用会把连接池和 goroutine 一起拖死。
//   - 重试与退避: 只对幂等操作(Query, 以及带 TaskNo 幂等键的 Submit)重试,
//     指数退避 + 抖动 + 最大次数上限; 4xx 一般不重试, 5xx/超时才重试。
//   - 限流: 供应商基本都有 QPS/并发上限, 本地要先自限(令牌桶 + 全局并发信号量),
//     被对面 429 之后再退避已经晚了 —— 那时任务已经失败, 要走退款。
//   - 幂等: Submit 必须把我们的 TaskNo 带过去作为幂等键, 网络超时后重试才不会
//     在对面重复建单(重复建单 = 重复计费, 对面的钱和我们用户的钱都会错)。
//   - 回调验签: 按供应商自己的协议做(多数是 HMAC-SHA256 over raw body + 时间戳防重放)。
//     本项目骨架里用的是 md5(provider_task_id + status + secret), 那是给 mock/自测用的
//     最简形式, 接入时按对方协议整体替换, 不要拿它去对接真实供应商。
//   - 回调来源校验: IP 白名单 / mTLS, 以及"回调里带的金额与状态必须与本地任务对得上"。
//   - 对账: 供应商侧的任务清单与本地 ai_task 定期比对, 捞出"本地还在处理中但对面早就
//     结束"的单子(回调丢了)以及"对面成功但本地已退款"的单子(赔钱漏洞)。
//   - 结果落地: 供应商返回的产物 URL 通常是临时链接(几小时过期), 要转存到自己的对象存储
//     再写进 ai_task.result, 否则用户过两天回来看就是一堆死链。
//   - 内容安全: 生成结果入库前过一遍审核(这类玩法尤其需要), 不合规的直接置失败并退款。
package aiprovider

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
)

// SubmitInput 提交给供应商的任务描述。
// TaskNo 是我们自己的业务单号, 必须原样透传给供应商作幂等键(见包注释)。
type SubmitInput struct {
	TaskNo   string
	BizType  int
	Params   map[string]any
	InputURL string
}

// SubmitOutput 供应商受理后返回的外部任务号, 后续回调/轮询都靠它反查本地任务。
type SubmitOutput struct {
	ProviderTaskId string
}

// Provider 供应商适配器接口。
//
// Query 返回的 status 直接复用 ai_task 的状态码(2处理中 3成功 4失败),
// 由 adapter 负责把供应商自己的状态枚举翻译过来 —— 业务层不该认识任何供应商的状态字符串。
type Provider interface {
	Name() string
	Submit(ctx context.Context, in SubmitInput) (*SubmitOutput, error)
	Query(ctx context.Context, providerTaskId string) (status int, result map[string]any, errMsg string, err error)
}

// 与 ai_task 对齐的状态码(adapter 翻译目标), 避免 adapter 里出现魔法数字。
const (
	StatusRunning = 2
	StatusSucceed = 3
	StatusFailed  = 4
)

// DefaultName 未配置 ai_provider 时使用的默认供应商: 本地 mock。
// 生产上线前必须把 app_config 的 ai_provider 改成真实 adapter 名。
const DefaultName = "mock"

var (
	mu       sync.RWMutex
	registry = map[string]Provider{}
)

// Register 注册适配器(在各 adapter 的 init 中调用)。重名覆盖, 方便测试替身。
func Register(name string, p Provider) {
	if name == "" || p == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	registry[name] = p
}

// Get 取适配器。取不到返回 error 而不是 nil Provider —— 让调用方在提交阶段就失败并退款,
// 而不是等到一个 nil 指针在深处 panic。
func Get(name string) (Provider, error) {
	if name == "" {
		name = DefaultName
	}
	mu.RLock()
	defer mu.RUnlock()
	p, ok := registry[name]
	if !ok {
		return nil, gerror.Newf("AI供应商[%s]未接入", name)
	}
	return p, nil
}

// Names 已注册的适配器名(后台排查配置用)。
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}

// ---------------------------------------------------------------- mock

// mockProvider 本地联调 / e2e 用的假供应商: 不发任何网络请求, 立即"受理", 查询即成功。
//
// 它存在的意义是让整条链路(扣费 → 提交 → 回调/轮询 → 结果落库 → 失败退款)在没有供应商的
// 情况下也能被完整测试。外部任务号做成 "mock-<task_no>" 这样的确定性值, 便于自测时
// 由业务单号直接推出外部单号, 也便于日志里一眼看出这是假单。
type mockProvider struct{}

func init() { Register(DefaultName, &mockProvider{}) }

func (m *mockProvider) Name() string { return DefaultName }

func (m *mockProvider) Submit(_ context.Context, in SubmitInput) (*SubmitOutput, error) {
	return &SubmitOutput{ProviderTaskId: "mock-" + in.TaskNo}, nil
}

// Query 固定返回成功 + 占位产物。真实 adapter 这里要区分"还在跑"和"已失败",
// 并把供应商的错误信息翻译成人话写进 errMsg(会原样展示给用户与客服)。
func (m *mockProvider) Query(_ context.Context, providerTaskId string) (int, map[string]any, string, error) {
	return StatusSucceed, map[string]any{
		"mock":     true,
		"task_id":  providerTaskId,
		"url":      "https://example.com/mock/" + providerTaskId + ".png",
		"text":     "mock result",
		"finished": true,
	}, "", nil
}
