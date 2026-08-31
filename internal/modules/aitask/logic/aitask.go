// Package logic AI 生成任务业务(合并 tianbi 的 ai / aimate / aitexttoimage /
// aiimagetovideo / aitexttonovel / AI客服 六套玩法)。
//
// 这一层刻意做成「与供应商无关」: 它只认 shared/aiprovider 的 Provider 接口和 ai_task 的
// 状态码, 不知道对面是谁。供应商定下来后新增一个 adapter 即可, 本文件一行不用改。
//
// 三条必须守住的正确性红线(tianbi 原版这三条全破了):
//
//  1. 钱只能扣一次。
//     客户端双击/断网重试会重复发同一个提交请求。这里用 client_token 做幂等:
//     事务内"先查后插" + DB 层部分唯一索引兜底, 并发下第二条会被唯一约束顶掉;
//     再加 users 行锁把同一用户的并发提交串行化。
//
//  2. 钱只能退一次。
//     退款一律走「条件更新 + RowsAffected 判定」: 先在事务里 UPDATE ... WHERE status IN (...),
//     只有真的改动了一行(说明这一刻它还不是终态)才继续调 balance.Add。
//     供应商反复重推回调、客服手抖点两次退款、取消与回调同时到达 —— 全部只会成功一次。
//     判断"该不该退"绝不用先 SELECT 再 if 的写法, 那是典型的 TOCTOU, 并发下会退两次。
//
//  3. 外部调用不能待在事务里。
//     provider.Submit 是网络请求, 可能几秒才超时。放在事务里就等于用一个数据库事务
//     (以及一条连接、一堆行锁)去等第三方, 高峰期能把连接池打满、把别人的扣费全拖死。
//     所以这里是「事务提交(钱扣了, 单建了, status=排队中) → 再调供应商 → 按结果二次落库」。
//     中间挂掉也不会丢钱: 任务停在"排队中", 后台可以重试或人工退款。
package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/aitask/service"
	"github.com/JarvanDante/my_service/internal/shared/aiprovider"
	"github.com/JarvanDante/my_service/internal/shared/appcfg"
	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const aiSiteId = 1 // 单站点样板, 与其它模块的 xxSiteId 保持一致

// app_config 里的运营参数 key(集中登记, 免得散落在各处拼错)。
const (
	cfgOpen           = "ai_open"             // 总开关
	cfgDefaultCost    = "ai_default_cost"     // 无模板时的默认单价
	cfgDailyLimit     = "ai_daily_limit"      // 每人每日任务数上限, <=0 不限
	cfgCallbackSecret = "ai_callback_secret"  // 回调验签密钥
	cfgProvider       = "ai_provider"         // 当前供应商名, 默认 mock
	cfgTaskTimeout    = "ai_task_timeout_sec" // 排队/处理超时秒数, 超时失败退款
)

type sAiTask struct{}

func New() service.IAiTask { return &sAiTask{} }

// ---------------------------------------------------------------- 小工具

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

// encodeJSON 落 jsonb 用。nil / 编码失败一律退化成 "{}", 保证列上的 NOT NULL 约束不被打破。
func encodeJSON(v map[string]any) string {
	if len(v) == 0 {
		return "{}"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// decodeJSON 读 jsonb。解析失败返回空 map 而不是报错 —— 一条脏数据不该让整个列表挂掉。
func decodeJSON(raw string) map[string]any {
	out := map[string]any{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func tplDTO(r *entity.AiTemplate) *service.TemplateDTO {
	return &service.TemplateDTO{
		Id: r.Id, Name: r.Name, BizType: r.BizType, Cover: r.Cover, Preview: r.Preview,
		Params: decodeJSON(r.Params), CostGold: r.CostGold, Sort: r.Sort,
		Status: r.Status, CreatedAt: fmtTime(r.CreatedAt),
	}
}

func taskDTO(r *entity.AiTask) *service.TaskDTO {
	return &service.TaskDTO{
		Id: r.Id, UserId: r.UserId, TaskNo: r.TaskNo, ClientToken: r.ClientToken,
		BizType: r.BizType, TemplateId: r.TemplateId, Params: decodeJSON(r.Params),
		InputUrl: r.InputUrl, CostGold: r.CostGold, Status: r.Status,
		Provider: r.Provider, ProviderTaskId: r.ProviderTaskId, Result: decodeJSON(r.Result),
		ErrMsg: r.ErrMsg, RetryCount: r.RetryCount, SubmittedAt: fmtTime(r.SubmittedAt),
		FinishedAt: fmtTime(r.FinishedAt), CreatedAt: fmtTime(r.CreatedAt),
	}
}

// genTaskNo 业务单号: A + 时间 + 用户 + 随机, 与 withdrawal 的 tradeNo 同一套风格。
// 唯一性最终由 uk_ai_task_no 保证, 这里只负责"基本不撞"。
func genTaskNo(userId int64) string {
	return "A" + gtime.Now().Format("YmdHis") + gconv.String(userId) + grand.Digits(6)
}

// truncErr 错误信息落库前截断: err_msg 列是 varchar(500), 供应商的报文可能很长,
// 超长会直接插入失败, 反而把"记录失败原因"这件事本身搞失败。
func truncErr(s string) string {
	r := []rune(s)
	if len(r) > 480 {
		return string(r[:480]) + "..."
	}
	return s
}

// currentProvider 取当前供应商适配器。供应商未定/配错时返回 error,
// 调用方据此走退款, 而不是让任务和用户的钱一起卡死。
func currentProvider(ctx context.Context) (aiprovider.Provider, error) {
	return aiprovider.Get(appcfg.String(ctx, cfgProvider, aiprovider.DefaultName))
}

func userBalance(ctx context.Context, userId int64) float64 {
	v, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil || v == nil {
		return 0
	}
	return v.Float64()
}

// findTask 按主键取任务。
func findTask(ctx context.Context, id int64) (*entity.AiTask, error) {
	var t *entity.AiTask
	if err := g.Model("ai_task").Ctx(ctx).
		Where("site_id", aiSiteId).Where("id", id).Scan(&t); err != nil {
		return nil, err
	}
	if t == nil {
		return nil, gerror.New("任务不存在")
	}
	return t, nil
}

// ---------------------------------------------------------------- 前台: 模板

func (s *sAiTask) FrontTemplates(ctx context.Context, bizType int) ([]*service.TemplateDTO, error) {
	m := g.Model("ai_template").Ctx(ctx).Where("site_id", aiSiteId).Where("status", 1)
	if bizType > 0 {
		m = m.Where("biz_type", bizType)
	}
	var list []*entity.AiTemplate
	if err := m.OrderDesc("sort").OrderDesc("id").Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.TemplateDTO, 0, len(list))
	for _, r := range list {
		out = append(out, tplDTO(r))
	}
	return out, nil
}

// ---------------------------------------------------------------- 前台: 提交

// priceOf 定价: 有模板按模板价, 没模板按配置默认价。**永远不看客户端传来的任何金额**。
// 模板必须启用且玩法匹配, 否则等于让用户拿"文生小说"的便宜模板去跑"图生视频"。
func priceOf(ctx context.Context, bizType int, templateId int64) (float64, *entity.AiTemplate, error) {
	if templateId <= 0 {
		return appcfg.Float(ctx, cfgDefaultCost, 20), nil, nil
	}
	var tpl *entity.AiTemplate
	if err := g.Model("ai_template").Ctx(ctx).
		Where("site_id", aiSiteId).Where("id", templateId).Scan(&tpl); err != nil {
		return 0, nil, err
	}
	if tpl == nil || tpl.Status != 1 {
		return 0, nil, gerror.New("模板不存在或已停用")
	}
	if tpl.BizType != bizType {
		return 0, nil, gerror.New("模板与玩法类型不匹配")
	}
	return tpl.CostGold, tpl, nil
}

// findByClientToken 按客户端幂等键查已有任务。
func findByClientToken(ctx context.Context, tx gdb.TX, userId int64, token string) (*entity.AiTask, error) {
	if token == "" {
		return nil, nil
	}
	m := g.Model("ai_task").Ctx(ctx)
	if tx != nil {
		m = tx.Model("ai_task").Ctx(ctx)
	}
	var t *entity.AiTask
	if err := m.Where("site_id", aiSiteId).Where("user_id", userId).
		Where("client_token", token).Scan(&t); err != nil {
		return nil, err
	}
	return t, nil
}

// dailyUsed 今天该用户已提交的任务数。以任务行数统计, 不维护可变计数器:
// 计数器要在失败/退款时回滚, 少回滚一次上限就永久少一次, 数行数天然自洽。
func dailyUsed(ctx context.Context, tx gdb.TX, userId int64) (int, error) {
	m := g.Model("ai_task").Ctx(ctx)
	if tx != nil {
		m = tx.Model("ai_task").Ctx(ctx)
	}
	return m.Where("site_id", aiSiteId).Where("user_id", userId).
		Where("created_at >= ?", gtime.Now().StartOfDay()).Count()
}

func (s *sAiTask) Submit(ctx context.Context, userId int64, in service.SubmitInput) (*service.SubmitOutput, error) {
	if userId <= 0 {
		return nil, gerror.New("未登录")
	}
	if !appcfg.Bool(ctx, cfgOpen, true) {
		return nil, gerror.New("AI功能维护中, 请稍后再试")
	}
	if in.BizType < 1 || in.BizType > 6 {
		return nil, gerror.New("玩法类型非法")
	}
	// 快路径: 事务外先查一次 client_token, 命中直接返回, 连事务都不用开。
	// (真正防并发的那道判断在事务里, 这里只是省开销。)
	if old, err := findByClientToken(ctx, nil, userId, in.ClientToken); err != nil {
		return nil, err
	} else if old != nil {
		return &service.SubmitOutput{Task: taskDTO(old), Balance: userBalance(ctx, userId), Repeated: true}, nil
	}

	cost, _, err := priceOf(ctx, in.BizType, in.TemplateId)
	if err != nil {
		return nil, err
	}
	provName := appcfg.String(ctx, cfgProvider, aiprovider.DefaultName)
	taskNo := genTaskNo(userId)

	var (
		taskId   int64
		repeated bool
		existing *entity.AiTask
	)
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 0. 用户行锁: 把同一用户的并发提交串行化(双击/重试的第二发要等第一发提交完再读)。
		//    只锁这个用户, 不影响别人提交。
		lock, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("id").LockUpdate().One()
		if err != nil {
			return err
		}
		if lock.IsEmpty() {
			return gerror.New("用户不存在")
		}
		// 1. 幂等键复查(拿到行锁之后): 前一发如果已经建单, 这里必然能看到, 直接复用不扣费。
		if old, err := findByClientToken(ctx, tx, userId, in.ClientToken); err != nil {
			return err
		} else if old != nil {
			existing, repeated = old, true
			return nil
		}
		// 2. 每日上限(行锁之后统计, 不会被并发击穿)。<=0 表示不限。
		limit := appcfg.Int(ctx, cfgDailyLimit, 20)
		if limit > 0 {
			used, err := dailyUsed(ctx, tx, userId)
			if err != nil {
				return err
			}
			if used >= limit {
				return gerror.Newf("今日AI任务次数已达上限(%d次)", limit)
			}
		}
		// 3. 扣费: 条件扣款, SQL 层防透支(余额不足直接报错回滚, 任务不会建出来)。
		if cost > 0 {
			if err := balance.Deduct(ctx, tx, userId, cost,
				balance.SceneAiCost, taskNo, "AI生成任务消耗"); err != nil {
				if gerror.Is(err, balance.ErrInsufficient) {
					return gerror.New("金币余额不足")
				}
				return err
			}
		}
		// 4. 建单(排队中)。此刻钱已经扣了、单已经落库, 即使下一步调供应商时进程挂掉,
		//    也只是留下一条"排队中"的任务, 后台能重试或退款, 不会出现"扣了钱没有单"。
		id, err := tx.Model("ai_task").Ctx(ctx).Data(g.Map{
			"site_id": aiSiteId, "user_id": userId, "biz_type": in.BizType,
			"template_id": in.TemplateId, "task_no": taskNo, "client_token": in.ClientToken,
			"params": encodeJSON(in.Params), "input_url": in.InputUrl, "cost_gold": cost,
			"status": entity.AiStatusQueued, "provider": provName,
			"submitted_at": gtime.Now(),
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		taskId = id
		return nil
	})
	if err != nil {
		return nil, err
	}
	if repeated {
		return &service.SubmitOutput{Task: taskDTO(existing), Balance: userBalance(ctx, userId), Repeated: true}, nil
	}

	// 5. 事务已提交, 现在才调供应商(见包注释第 3 条: 网络请求绝不待在事务里)。
	s.dispatch(ctx, taskId, taskNo, userId, in.BizType, in.Params, in.InputUrl, cost)

	t, err := findTask(ctx, taskId)
	if err != nil {
		return nil, err
	}
	return &service.SubmitOutput{Task: taskDTO(t), Balance: userBalance(ctx, userId)}, nil
}

// dispatch 把排队中的任务投给供应商, 并按结果二次落库。三条分支:
//
//   - 供应商受理: status 排队中 → 处理中, 记下 provider_task_id;
//   - 供应商拒单(Submit 返回 error): 这一单确实办不成 → 立刻退款补偿,
//     用户的钱不能因为我们没投递成功而蒸发;
//   - 供应商未接入/配置错(拿不到 adapter): **保持排队中**, 只记 err_msg。
//     这两种失败要区别对待 —— 前者是单笔业务失败, 后者是运营配置问题, 影响的是所有人。
//     配置错时把成千上万单一起退掉, 用户会看到一堆"失败", 运营改好配置后还得逐单安抚;
//     留在"排队中"则改好配置后一键重试即可, 期间用户看到的是"排队中"(事实如此),
//     真要退也有后台人工退款兜底。这也是 e2e 里能造出"排队中"任务的路径。
//
// 两处落库都用 WHERE status=排队中 的条件更新: 万一供应商的回调比我们的返回还快
// (确实会发生), 回调已经把任务推到终态了, 这里就不能再把它拽回"处理中"。
func (s *sAiTask) dispatch(ctx context.Context, taskId int64, taskNo string, userId int64,
	bizType int, params map[string]any, inputUrl string, cost float64) {
	p, err := currentProvider(ctx)
	if err != nil {
		// 供应商未接入: 任务留在排队中, 等运营配好后由后台重试。
		g.Log().Errorf(ctx, "AI任务无法投递(供应商未接入), 任务保持排队中 task_no=%s: %v", taskNo, err)
		_, _ = g.Model("ai_task").Ctx(ctx).
			Where("id", taskId).Where("status", entity.AiStatusQueued).
			Data(g.Map{"err_msg": truncErr(err.Error()), "updated_at": gtime.Now()}).Update()
		return
	}
	out, err := p.Submit(ctx, aiprovider.SubmitInput{
		TaskNo: taskNo, BizType: bizType, Params: params, InputURL: inputUrl,
	})
	if err == nil {
		if _, uerr := g.Model("ai_task").Ctx(ctx).
			Where("id", taskId).Where("status", entity.AiStatusQueued).
			Data(g.Map{
				"status": entity.AiStatusRunning, "provider": p.Name(),
				"provider_task_id": out.ProviderTaskId, "updated_at": gtime.Now(),
			}).Update(); uerr != nil {
			g.Log().Errorf(ctx, "AI任务投递成功但落库失败 task_no=%s: %v", taskNo, uerr)
		}
		return
	}
	// 供应商拒单 → 退款补偿。
	g.Log().Warningf(ctx, "AI任务投递失败, 触发退款 task_no=%s: %v", taskNo, err)
	if rerr := refundQueued(ctx, taskId, userId, cost, taskNo, truncErr("提交供应商失败: "+err.Error())); rerr != nil {
		// 退款也失败: 任务停在"排队中", 由后台人工退款兜底。必须打日志, 这是要人管的钱。
		g.Log().Errorf(ctx, "AI任务投递失败后退款失败, 需人工处理 task_no=%s: %v", taskNo, rerr)
	}
}

// refundQueued 把"排队中"的任务退款(投递失败专用)。
// 条件更新只认 status=排队中: 若期间回调已把它推成终态, 这里 RowsAffected=0, 直接不退,
// 避免与回调的退款撞车退两次。
func refundQueued(ctx context.Context, taskId, userId int64, cost float64, taskNo, errMsg string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("ai_task").Ctx(ctx).
			Where("id", taskId).Where("status", entity.AiStatusQueued).
			Data(g.Map{
				"status": entity.AiStatusRefunded, "err_msg": errMsg,
				"finished_at": gtime.Now(), "updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return nil // 已被回调/取消处理过, 不重复退款
		}
		if cost <= 0 {
			return nil
		}
		return balance.Add(ctx, tx, userId, cost, balance.SceneAiRefund, taskNo, "AI任务提交失败退款")
	})
}

// ---------------------------------------------------------------- 前台: 查询

func (s *sAiTask) Task(ctx context.Context, userId, id int64) (*service.TaskDTO, error) {
	t, err := findTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.UserId != userId {
		return nil, gerror.New("任务不存在")
	}
	// 轮询兜底: 回调是"尽力而为"的 —— 供应商可能没推、推丢了、或者我们当时正好在重启。
	// 只要任务还处在"处理中", 用户每次点开详情就顺手向供应商问一次最新状态并落库。
	// 这条路径与回调走同一套条件更新, 两边同时到达也只会有一方生效。
	if t.Status == entity.AiStatusRunning && t.ProviderTaskId != "" {
		if p, perr := currentProvider(ctx); perr == nil {
			st, result, errMsg, qerr := p.Query(ctx, t.ProviderTaskId)
			if qerr != nil {
				g.Log().Warningf(ctx, "AI任务轮询失败 task_no=%s: %v", t.TaskNo, qerr)
			} else if st == aiprovider.StatusSucceed || st == aiprovider.StatusFailed {
				if aerr := applyTerminal(ctx, t.Id, st, result, errMsg); aerr != nil {
					g.Log().Errorf(ctx, "AI任务轮询结果落库失败 task_no=%s: %v", t.TaskNo, aerr)
				}
				if t, err = findTask(ctx, id); err != nil {
					return nil, err
				}
			}
		}
	}
	return taskDTO(t), nil
}

func (s *sAiTask) Tasks(ctx context.Context, userId int64, f service.TaskFilter) ([]*service.TaskDTO, int, error) {
	if userId <= 0 {
		return nil, 0, gerror.New("未登录")
	}
	f.UserId = userId
	return queryTasks(ctx, f)
}

func taskNeedUserJoin(f service.TaskFilter) bool {
	return f.Nickname != "" || f.ChannelName != "" || f.DeviceType != "" ||
		f.RegisterStart != "" || f.RegisterEnd != ""
}

func applyTaskFilter(m *gdb.Model, f service.TaskFilter) *gdb.Model {
	m = m.Where("ai_task.site_id", aiSiteId)
	if f.UserId > 0 {
		m = m.Where("ai_task.user_id", f.UserId)
	}
	if f.BizType > 0 {
		m = m.Where("ai_task.biz_type", f.BizType)
	}
	if f.Status > 0 {
		m = m.Where("ai_task.status", f.Status)
	}
	if f.TaskNo != "" {
		m = m.Where("ai_task.task_no", f.TaskNo)
	}
	if f.StartTime != "" {
		m = m.Where("ai_task.created_at >= ?", f.StartTime)
	}
	if f.EndTime != "" {
		m = m.Where("ai_task.created_at <= ?", f.EndTime)
	}
	if taskNeedUserJoin(f) {
		m = m.LeftJoin("users", "users.id = ai_task.user_id")
		if f.Nickname != "" {
			like := "%" + f.Nickname + "%"
			m = m.Where("(users.nickname ILIKE ? OR users.username ILIKE ? OR users.phone ILIKE ?)", like, like, like)
		}
		if f.ChannelName != "" {
			m = m.Where("users.channel_name ILIKE ?", "%"+f.ChannelName+"%")
		}
		if f.DeviceType != "" {
			m = m.Where("users.device_type", f.DeviceType)
		}
		if f.RegisterStart != "" {
			m = m.Where("users.register_at >= ?", f.RegisterStart)
		}
		if f.RegisterEnd != "" {
			m = m.Where("users.register_at <= ?", f.RegisterEnd)
		}
	}
	return m
}

func queryTasks(ctx context.Context, f service.TaskFilter) ([]*service.TaskDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 || f.Size > 100 {
		f.Size = 20
	}
	base := applyTaskFilter(g.Model("ai_task").Ctx(ctx), f)
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.AiTask
	if err := base.Clone().Fields("ai_task.*").OrderDesc("ai_task.id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.TaskDTO, 0, len(list))
	for _, r := range list {
		out = append(out, taskDTO(r))
	}
	return out, total, nil
}

func queryTaskStats(ctx context.Context, f service.TaskFilter) *service.TaskStats {
	stats := &service.TaskStats{}
	rows, err := applyTaskFilter(g.Model("ai_task").Ctx(ctx), f).
		Fields("ai_task.status AS status, count(1) AS cnt, coalesce(sum(ai_task.cost_gold),0) AS gold").
		Group("ai_task.status").All()
	if err != nil {
		return stats
	}
	for _, row := range rows {
		n := row["cnt"].Int()
		gold := row["gold"].Float64()
		stats.Total += n
		stats.TotalGold += gold
		switch row["status"].Int() {
		case entity.AiStatusSucceed:
			stats.Success = n
			stats.SuccessGold = gold
		case entity.AiStatusRefunded, entity.AiStatusCancelled:
			stats.Refund += n
			stats.RefundGold += gold
		case entity.AiStatusFailed:
			stats.Abnormal = n
			stats.AbnormalGold = gold
		}
	}
	return stats
}

func fillTaskUsers(ctx context.Context, list []*service.TaskDTO) {
	if len(list) == 0 {
		return
	}
	ids := make([]int64, 0, len(list))
	seen := map[int64]struct{}{}
	for _, r := range list {
		if r.UserId <= 0 {
			continue
		}
		if _, ok := seen[r.UserId]; ok {
			continue
		}
		seen[r.UserId] = struct{}{}
		ids = append(ids, r.UserId)
	}
	if len(ids) == 0 {
		return
	}
	var users []*entity.Users
	if err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Scan(&users); err != nil {
		return
	}
	byID := map[int64]*entity.Users{}
	for _, u := range users {
		byID[u.Id] = u
	}
	for _, r := range list {
		r.Sets = setsOf(r.Params)
		u := byID[r.UserId]
		if u == nil {
			continue
		}
		r.Nickname = u.Nickname
		r.Phone = u.Phone
		r.Avatar = u.Img
		r.GroupName = u.GroupName
		r.ChannelName = u.ChannelName
		r.DeviceType = u.DeviceType
	}
}

func setsOf(params map[string]any) int {
	if params == nil {
		return 1
	}
	switch v := params["sets"].(type) {
	case float64:
		if int(v) > 0 {
			return int(v)
		}
	case int:
		if v > 0 {
			return v
		}
	case json.Number:
		if n, err := v.Int64(); err == nil && n > 0 {
			return int(n)
		}
	}
	return 1
}

// ---------------------------------------------------------------- 前台: 取消

// Cancel 取消并退款。只允许取消"排队中"的任务: 一旦投给供应商, 对方已经在烧算力,
// 这时候的取消要走供应商的取消协议(各家都不一样), 属于接入真实供应商时再补的能力。
//
// 防重复退款同样靠条件更新: 两次取消请求同时进来, 只有一次能把 status 从 1 改成 6。
func (s *sAiTask) Cancel(ctx context.Context, userId, id int64) (float64, float64, error) {
	if userId <= 0 {
		return 0, 0, gerror.New("未登录")
	}
	t, err := findTask(ctx, id)
	if err != nil {
		return 0, 0, err
	}
	if t.UserId != userId {
		return 0, 0, gerror.New("任务不存在")
	}
	var refund float64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("ai_task").Ctx(ctx).
			Where("id", id).Where("user_id", userId).Where("status", entity.AiStatusQueued).
			Data(g.Map{
				"status": entity.AiStatusCancelled, "err_msg": "用户取消",
				"finished_at": gtime.Now(), "updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("任务已在处理或已结束, 无法取消")
		}
		if t.CostGold > 0 {
			if err := balance.Add(ctx, tx, userId, t.CostGold,
				balance.SceneAiRefund, t.TaskNo, "AI任务取消退款"); err != nil {
				return err
			}
			refund = t.CostGold
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return refund, userBalance(ctx, userId), nil
}

// ---------------------------------------------------------------- 回调

// Callback 供应商回调。
//
// 幂等是这个方法的全部重点: 供应商为了"确保送达"普遍会重推同一条回调(有的按分钟级重推
// 十几次), 网络抖动也会造成重复投递。所以:
//   - 落终态用 UPDATE ... WHERE status IN (排队中, 处理中): 只有第一次能改到行;
//   - RowsAffected==0 意味着这单早就是终态了 —— 直接返回成功(不是错误!),
//     否则供应商会认为投递失败, 继续无限重推;
//   - 退款只在 RowsAffected==1 的分支里做, 且与状态更新在同一个事务里, 要么一起成功
//     要么一起回滚, 不存在"状态改了钱没退"或"钱退了状态没改"的中间态。
func (s *sAiTask) Callback(ctx context.Context, in service.CallbackInput) error {
	if in.ProviderTaskId == "" {
		return gerror.New("外部任务号必填")
	}
	if in.Status != entity.AiStatusSucceed && in.Status != entity.AiStatusFailed {
		return gerror.New("回调状态非法")
	}
	// 验签: md5(provider_task_id + status + secret)。
	// 这是骨架期的最简形式(与本项目支付回调保持一致的风格); 真实供应商多用
	// HMAC-SHA256(raw body) + 时间戳/nonce 防重放, 接入时按对方协议整体替换这一段。
	secret := appcfg.String(ctx, cfgCallbackSecret, "")
	if secret != "" {
		want := gmd5.MustEncryptString(in.ProviderTaskId + strconv.Itoa(in.Status) + secret)
		if in.Sign == "" || in.Sign != want {
			return gerror.New("回调签名校验失败")
		}
	}
	var t *entity.AiTask
	m := g.Model("ai_task").Ctx(ctx).
		Where("site_id", aiSiteId).Where("provider_task_id", in.ProviderTaskId)
	if in.Provider != "" {
		m = m.Where("provider", in.Provider)
	}
	if err := m.OrderDesc("id").Scan(&t); err != nil {
		return err
	}
	if t == nil {
		return gerror.New("任务不存在")
	}
	return applyTerminal(ctx, t.Id, in.Status, in.Result, in.ErrMsg)
}

func (s *sAiTask) HandleWorkerResult(ctx context.Context, jobID, status, outputURL, outputKey, errMsg string) error {
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return gerror.New("job_id 为空")
	}
	var t *entity.AiTask
	if err := g.Model("ai_task").Ctx(ctx).
		Where("site_id", aiSiteId).Where("task_no", jobID).Scan(&t); err != nil {
		return err
	}
	if t == nil {
		return gerror.Newf("任务不存在: %s", jobID)
	}
	if status == "ready" {
		url := strings.TrimSpace(outputURL)
		if url == "" && outputKey != "" {
			url = publicObjectURL(ctx, faceswapBucket(ctx), outputKey)
		}
		return applyTerminal(ctx, t.Id, entity.AiStatusSucceed, map[string]any{
			"url":        url,
			"output_key": outputKey,
		}, "")
	}
	if errMsg == "" {
		errMsg = "生成失败"
	}
	return applyTerminal(ctx, t.Id, entity.AiStatusFailed, nil, errMsg)
}

func faceswapBucket(ctx context.Context) string {
	b := strings.TrimSpace(g.Cfg().MustGet(ctx, "faceswap.bucket", "my-storage").String())
	if b == "" {
		return "my-storage"
	}
	return b
}

func publicObjectURL(ctx context.Context, bucket, key string) string {
	base := strings.TrimRight(g.Cfg().MustGet(ctx, "minio.publicURL").String(), "/")
	bucket = strings.Trim(bucket, "/")
	key = strings.TrimLeft(key, "/")
	if base == "" || bucket == "" || key == "" {
		return ""
	}
	return base + "/" + bucket + "/" + key
}

func (s *sAiTask) ExpireStale(ctx context.Context) (int, error) {
	sec := appcfg.Int(ctx, cfgTaskTimeout, 600)
	if sec <= 0 {
		return 0, nil
	}
	cutoff := gtime.Now().Add(-time.Duration(sec) * time.Second)
	var list []*entity.AiTask
	if err := g.Model("ai_task").Ctx(ctx).
		Where("site_id", aiSiteId).
		WhereIn("status", []int{entity.AiStatusQueued, entity.AiStatusRunning}).
		Where("submitted_at IS NOT NULL").
		Where("submitted_at < ?", cutoff).
		OrderAsc("id").Limit(100).Scan(&list); err != nil {
		return 0, err
	}
	n := 0
	for _, t := range list {
		if t == nil {
			continue
		}
		if err := applyTerminal(ctx, t.Id, entity.AiStatusFailed, nil, "任务超时未完成"); err != nil {
			g.Log().Warningf(ctx, "AI任务超时退款失败 task_no=%s: %v", t.TaskNo, err)
			continue
		}
		n++
	}
	if n > 0 {
		g.Log().Infof(ctx, "AI任务超时清理 %d 条", n)
	}
	return n, nil
}

// applyTerminal 把任务落到终态, 失败则退款。回调与轮询共用这一个入口, 两条路径同时到达
// 也只会有一方真正生效(条件更新 + 行锁)。
//
// 参数 status: 3成功 4失败。失败会被落成 5已退款(钱退回去了, 状态必须能体现)。
func applyTerminal(ctx context.Context, taskId int64, status int, result map[string]any, errMsg string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先锁行: 并发的重复回调会在这里排队, 后到的那个拿到锁时已经能看到终态,
		// 下面的条件更新自然影响 0 行。行锁 + 条件更新是双保险, 缺一都可能退两次。
		var t *entity.AiTask
		if err := tx.Model("ai_task").Ctx(ctx).Where("id", taskId).LockUpdate().Scan(&t); err != nil {
			return err
		}
		if t == nil {
			return gerror.New("任务不存在")
		}
		data := g.Map{"finished_at": gtime.Now(), "updated_at": gtime.Now()}
		if status == entity.AiStatusSucceed {
			data["status"] = entity.AiStatusSucceed
			data["result"] = encodeJSON(result)
			data["err_msg"] = ""
		} else {
			// 失败即退款, 所以终态直接写"已退款"而不是"失败":
			// 状态本身就是"这笔钱退过了"的凭证, 后续任何重试都能据此判断。
			data["status"] = entity.AiStatusRefunded
			data["err_msg"] = truncErr(errMsg)
		}
		res, err := tx.Model("ai_task").Ctx(ctx).Where("id", taskId).
			WhereIn("status", []int{entity.AiStatusQueued, entity.AiStatusRunning}).
			Data(data).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			// 已经是终态: 重复回调。返回 nil = 告诉供应商"收到了", 让它别再重推。
			return nil
		}
		if status == entity.AiStatusFailed && t.CostGold > 0 {
			return balance.Add(ctx, tx, t.UserId, t.CostGold,
				balance.SceneAiRefund, t.TaskNo, "AI任务生成失败退款")
		}
		return nil
	})
}

// ---------------------------------------------------------------- 后台: 模板

func (s *sAiTask) TemplateList(ctx context.Context, f service.TemplateFilter) ([]*service.TemplateDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 || f.Size > 100 {
		f.Size = 20
	}
	base := g.Model("ai_template").Ctx(ctx).Where("site_id", aiSiteId)
	if f.BizType > 0 {
		base = base.Where("biz_type", f.BizType)
	}
	if f.Status >= 0 {
		base = base.Where("status", f.Status)
	}
	if f.Keyword != "" {
		base = base.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.AiTemplate
	if err := base.Clone().OrderDesc("sort").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.TemplateDTO, 0, len(list))
	for _, r := range list {
		out = append(out, tplDTO(r))
	}
	fillTemplateUsage(ctx, out)
	return out, total, nil
}

func fillTemplateUsage(ctx context.Context, list []*service.TemplateDTO) {
	if len(list) == 0 {
		return
	}
	ids := make([]int64, 0, len(list))
	for _, r := range list {
		ids = append(ids, r.Id)
	}
	rows, err := g.Model("ai_task").Ctx(ctx).
		Fields("template_id", "count(1) AS cnt").
		Where("site_id", aiSiteId).
		WhereIn("template_id", ids).
		Group("template_id").All()
	if err != nil {
		return
	}
	cnt := map[int64]int{}
	for _, row := range rows {
		cnt[row["template_id"].Int64()] = row["cnt"].Int()
	}
	for _, r := range list {
		r.UsageCount = cnt[r.Id]
	}
}

func (s *sAiTask) TemplateCreate(ctx context.Context, in service.TemplateInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("模板名必填")
	}
	if in.BizType < 1 || in.BizType > 6 {
		return 0, gerror.New("玩法类型非法")
	}
	if in.CostGold < 0 {
		return 0, gerror.New("价格不能为负")
	}
	return g.Model("ai_template").Ctx(ctx).Data(g.Map{
		"site_id": aiSiteId, "name": in.Name, "biz_type": in.BizType, "cover": in.Cover,
		"preview": in.Preview, "params": encodeJSON(in.Params), "cost_gold": in.CostGold,
		"sort": in.Sort, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sAiTask) TemplateUpdate(ctx context.Context, in service.TemplateInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.CostGold < 0 {
		return gerror.New("价格不能为负")
	}
	data := g.Map{
		"cover": in.Cover, "preview": in.Preview, "params": encodeJSON(in.Params),
		"cost_gold": in.CostGold, "sort": in.Sort, "status": in.Status,
		"updated_at": gtime.Now(),
	}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.BizType > 0 {
		data["biz_type"] = in.BizType
	}
	_, err := g.Model("ai_template").Ctx(ctx).
		Where("site_id", aiSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

// TemplateDelete 删除模板。已产生的任务不受影响: 任务上冗余了下单时的 cost_gold 与 params,
// 模板删了也不影响历史任务的展示与退款金额。
func (s *sAiTask) TemplateDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("ai_template").Ctx(ctx).
		Where("site_id", aiSiteId).Where("id", id).Delete()
	return err
}

// ---------------------------------------------------------------- 后台: 任务

func (s *sAiTask) TaskList(ctx context.Context, f service.TaskFilter) ([]*service.TaskDTO, int, *service.TaskStats, error) {
	list, total, err := queryTasks(ctx, f)
	if err != nil {
		return nil, 0, nil, err
	}
	fillTaskUsers(ctx, list)
	return list, total, queryTaskStats(ctx, f), nil
}

// TaskRetry 重新提交任务, retry_count+1, 再投一次供应商。
//
// 扣不扣费取决于上一次的钱在谁手里 —— 这是这个方法唯一容易搞错的地方:
//   - 失败(4)/已退款(5): 钱已经退回给用户了(status=5 本身就是退款凭证),
//     所以重试是一次全新的消耗, 必须重新扣费;
//   - 排队中(1): 钱还扣在我们这儿(通常是供应商未接入导致投递不出去),
//     重试只是把同一单再投一次, **绝不能再扣一次**, 否则用户为一次生成付两次钱。
//
// 并发保护用的是精确 CAS: WHERE id=? AND status=<刚读到的状态>。同一单被点两次重试,
// 第二次的 RowsAffected 必为 0(状态已经被第一次改成排队中), 直接报错, 不会重复扣费。
func (s *sAiTask) TaskRetry(ctx context.Context, id int64) (*service.TaskDTO, error) {
	t, err := findTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.Status != entity.AiStatusFailed && t.Status != entity.AiStatusRefunded &&
		t.Status != entity.AiStatusQueued {
		return nil, gerror.New("只有排队中/失败/已退款的任务可以重新提交")
	}
	recharge := t.Status != entity.AiStatusQueued // 排队中的任务钱没退, 不能再扣
	// 需要重新扣费时, 价格按"当前"模板价重算(重试是新消耗, 该按现价收);
	// 不重新扣费时沿用原价, 保证退款金额与实扣一致。
	cost := t.CostGold
	if recharge {
		if c, _, perr := priceOf(ctx, t.BizType, t.TemplateId); perr == nil {
			cost = c // 模板已删/已停用时保持原价, 不因为运营改配置就卡死重试
		}
	}
	taskNo := t.TaskNo
	if recharge {
		taskNo = genTaskNo(t.UserId) // 新的一笔消耗, 换新单号, 与新的流水一一对应
	}
	provName := appcfg.String(ctx, cfgProvider, aiprovider.DefaultName)

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("ai_task").Ctx(ctx).Where("id", id).Where("status", t.Status).
			Data(g.Map{
				"status": entity.AiStatusQueued, "task_no": taskNo, "cost_gold": cost,
				"provider": provName, "provider_task_id": "", "result": "{}",
				"err_msg": "", "retry_count": &gdb.Counter{Field: "retry_count", Value: 1},
				"submitted_at": gtime.Now(), "finished_at": nil, "updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("任务状态已变更, 请刷新后重试")
		}
		if recharge && cost > 0 {
			if err := balance.Deduct(ctx, tx, t.UserId, cost,
				balance.SceneAiCost, taskNo, "AI任务重新提交消耗"); err != nil {
				if gerror.Is(err, balance.ErrInsufficient) {
					return gerror.New("用户金币余额不足, 无法重新提交")
				}
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 同 Submit: 事务提交后再投递(网络请求不进事务)。
	s.dispatch(ctx, id, taskNo, t.UserId, t.BizType, decodeJSON(t.Params), t.InputUrl, cost)

	nt, err := findTask(ctx, id)
	if err != nil {
		return nil, err
	}
	return taskDTO(nt), nil
}

// TaskRefund 人工退款(客服兜底: 供应商产出质量太差、用户投诉等)。
//
// 允许退的只有"还没退过钱"的状态: 排队中/处理中/失败。已成功(3)不退(东西给了),
// 已退款(5)/已取消(6)不退(退过了)。同样是条件更新 + RowsAffected 判定, 连点两次只退一次。
func (s *sAiTask) TaskRefund(ctx context.Context, id int64, remark string) (float64, error) {
	t, err := findTask(ctx, id)
	if err != nil {
		return 0, err
	}
	if remark == "" {
		remark = "AI任务人工退款"
	}
	var refund float64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("ai_task").Ctx(ctx).Where("id", id).
			WhereIn("status", []int{entity.AiStatusQueued, entity.AiStatusRunning, entity.AiStatusFailed}).
			Data(g.Map{
				"status": entity.AiStatusRefunded, "err_msg": truncErr(remark),
				"finished_at": gtime.Now(), "updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("该任务已成功或已退款, 不能重复退款")
		}
		if t.CostGold > 0 {
			if err := balance.Add(ctx, tx, t.UserId, t.CostGold,
				balance.SceneAiRefund, t.TaskNo, remark); err != nil {
				return err
			}
			refund = t.CostGold
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return refund, nil
}

// TaskDelete 删订单。处理中的任务 worker 还可能回写, 不能删。
// 排队中还扣着用户的钱, 先按取消退款再删, 避免后台一点删除就把金币吞掉。
func (s *sAiTask) TaskDelete(ctx context.Context, id int64) error {
	t, err := findTask(ctx, id)
	if err != nil {
		return err
	}
	if t.Status == entity.AiStatusRunning {
		return gerror.New("处理中的任务不能删除, 请等完成或先人工退款")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if t.Status == entity.AiStatusQueued || t.Status == entity.AiStatusFailed {
			res, err := tx.Model("ai_task").Ctx(ctx).Where("id", id).
				WhereIn("status", []int{entity.AiStatusQueued, entity.AiStatusFailed}).
				Data(g.Map{
					"status": entity.AiStatusRefunded, "err_msg": "后台删除订单退款",
					"finished_at": gtime.Now(), "updated_at": gtime.Now(),
				}).Update()
			if err != nil {
				return err
			}
			if n, _ := res.RowsAffected(); n > 0 && t.CostGold > 0 {
				if err := balance.Add(ctx, tx, t.UserId, t.CostGold,
					balance.SceneAiRefund, t.TaskNo, "后台删除AI订单退款"); err != nil {
					return err
				}
			}
		}
		_, err := tx.Model("ai_task").Ctx(ctx).Where("id", id).Delete()
		return err
	})
}
