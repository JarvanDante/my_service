// Package logic 抽奖业务(移植自 tianbi lottery, 并修掉原版的一堆并发/一致性坑)。
//
// tianbi 原版的问题与本实现的修法:
//  1. 原版抽奖全程没有事务, 而且是「先发奖再扣费」: 并发下两个请求都能读到足够余额,
//     奖发出去了扣费才失败, 相当于绕过余额校验白嫖。
//     → 本实现把「校验 → 扣费 → 选奖 → 库存递减 → 发奖 → 写记录」全部放进一个事务,
//     任何一步返回 error 都整体回滚, 不存在「奖发了钱没扣」的中间态。
//  2. 原版扣免费次数是无条件 `$inc: -1`, 没有 `>= 0` 约束, 猛点能把次数刷成负数。
//     → 本实现不存"剩余次数"这种可变计数器, 改为「今天在 lottery_history 里的记录数」
//     反推剩余次数(见 dailyUsed): 记录是抽奖的唯一凭证, 数它天然不会为负, 也不用补偿。
//  3. 原版中奖历史用 `go func(){}` 异步写, 写失败奖已经发出去了, 对不上账。
//     → 本实现 history 与发奖同事务, 写不进去就把奖一起回滚。
//  4. 原版没有行锁也没有幂等, 双击会重复扣费。
//     → 本实现在事务最开头对 users 行 `SELECT ... FOR UPDATE`(见 Draw), 把同一用户的
//     并发抽奖串行化, 后到的请求要等前一次提交后才读次数, 免费次数/每日上限不会被击穿。
//     (不用"幂等键"是因为抽奖每次结果本就不同, 客户端重试应当是一次新抽奖; 需要拦的是
//     "同一瞬间的重复提交", 行锁 + 每日上限已经够了。)
package logic

import (
	"context"
	crand "crypto/rand"
	"math/big"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/lottery/service"
	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const ltSiteId = 1 // 单站点样板, 与其它模块的 xxSiteId 保持一致

// errCouponOut 券模板缺失/已发完的内部信号: 不让整次抽奖失败, 降级成"谢谢参与"。
var errCouponOut = gerror.New("优惠券已发完")

type sLottery struct{}

func New() service.ILottery { return &sLottery{} }

// ---------------------------------------------------------------- 转换

func actDTO(a *entity.LotteryActivity) *service.ActivityDTO {
	created := ""
	if a.CreatedAt != nil {
		created = a.CreatedAt.String()
	}
	return &service.ActivityDTO{
		Id: a.Id, Name: a.Name, LotteryType: a.LotteryType, PayType: a.PayType,
		CostGold: a.CostGold, DailyFree: a.DailyFree, DailyLimit: a.DailyLimit,
		Notice: a.Notice, Status: a.Status, CreatedAt: created,
	}
}

func prizeDTO(p *entity.LotteryPrize) *service.PrizeDTO {
	created := ""
	if p.CreatedAt != nil {
		created = p.CreatedAt.String()
	}
	return &service.PrizeDTO{
		Id: p.Id, ActivityId: p.ActivityId, Name: p.Name, Cover: p.Cover, Desc: p.Desc,
		Type: p.Type, Amount: p.Amount, CouponTplId: p.CouponTplId, Odds: p.Odds,
		Stock: p.Stock, Awarded: p.Awarded, Rank: p.Rank, Status: p.Status, CreatedAt: created,
	}
}

func histDTO(h *entity.LotteryHistory) *service.HistoryDTO {
	created := ""
	if h.CreatedAt != nil {
		created = h.CreatedAt.String()
	}
	return &service.HistoryDTO{
		Id: h.Id, UserId: h.UserId, ActivityId: h.ActivityId, LotteryType: h.LotteryType,
		PayType: h.PayType, CostGold: h.CostGold, PrizeId: h.PrizeId, PrizeName: h.PrizeName,
		PrizeType: h.PrizeType, PrizeAmount: h.PrizeAmount, Status: h.Status,
		Remark: h.Remark, CreatedAt: created,
	}
}

// maskName 跑马灯用的昵称脱敏: 只留首尾字符, 中间打码(别把用户昵称原样广播出去)。
func maskName(s string) string {
	r := []rune(strings.TrimSpace(s))
	switch {
	case len(r) == 0:
		return "神秘用户"
	case len(r) == 1:
		return string(r) + "**"
	case len(r) == 2:
		return string(r[0]) + "**"
	default:
		return string(r[0]) + "**" + string(r[len(r)-1])
	}
}

// nicknames 批量取昵称(避免跑马灯 N+1 查询)。
func nicknames(ctx context.Context, ids []int64) map[int64]string {
	out := map[int64]string{}
	if len(ids) == 0 {
		return out
	}
	all, err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Fields("id,nickname").All()
	if err != nil {
		return out
	}
	for _, r := range all {
		out[r["id"].Int64()] = r["nickname"].String()
	}
	return out
}

// ---------------------------------------------------------------- 公共查询

// activityByType 按玩法取启用中的活动。
func (s *sLottery) activityByType(ctx context.Context, lotteryType int) (*entity.LotteryActivity, error) {
	var a *entity.LotteryActivity
	if err := g.Model("lottery_activity").Ctx(ctx).
		Where("site_id", ltSiteId).Where("lottery_type", lotteryType).Where("status", 1).
		Scan(&a); err != nil {
		return nil, err
	}
	if a == nil {
		return nil, gerror.New("抽奖活动不存在或已停用")
	}
	return a, nil
}

// dailyUsed 统计"今天该用户在本活动下的抽奖记录数"(总次数, 以及其中免费次数)。
// 以记录数反推剩余次数, 不维护可变计数器 —— 这是修 tianbi「免费次数能刷成负数」的关键。
// 传 tx 时在事务内统计(配合 users 行锁, 结果对并发是稳定的); 传 nil 走普通连接(只读展示)。
func dailyUsed(ctx context.Context, tx gdb.TX, userId, activityId int64) (total, free int, err error) {
	m := g.Model("lottery_history").Ctx(ctx)
	if tx != nil {
		m = tx.Model("lottery_history").Ctx(ctx)
	}
	base := m.Where("site_id", ltSiteId).Where("user_id", userId).
		Where("activity_id", activityId).
		Where("created_at >= ?", gtime.Now().StartOfDay())
	total, err = base.Clone().Count()
	if err != nil {
		return 0, 0, err
	}
	free, err = base.Clone().Where("pay_type", service.PayFree).Count()
	if err != nil {
		return 0, 0, err
	}
	return total, free, nil
}

// leftOf 由"已用次数"算出剩余次数; daily_limit=0 表示不限, 用 -1 表达。
func leftOf(act *entity.LotteryActivity, total, free int) (freeLeft, drawLeft int) {
	freeLeft = act.DailyFree - free
	if freeLeft < 0 {
		freeLeft = 0
	}
	drawLeft = -1
	if act.DailyLimit > 0 {
		drawLeft = act.DailyLimit - total
		if drawLeft < 0 {
			drawLeft = 0
		}
	}
	return
}

func userBalance(ctx context.Context, tx gdb.TX, userId int64) (float64, error) {
	if userId <= 0 {
		return 0, nil
	}
	m := g.Model("users").Ctx(ctx)
	if tx != nil {
		m = tx.Model("users").Ctx(ctx)
	}
	v, err := m.Where("id", userId).Fields("balance").Value()
	if err != nil || v == nil {
		return 0, err
	}
	return v.Float64(), nil
}

// ---------------------------------------------------------------- 前台

func (s *sLottery) Info(ctx context.Context, userId int64, lotteryType int) (*service.InfoDTO, error) {
	act, err := s.activityByType(ctx, lotteryType)
	if err != nil {
		return nil, err
	}
	var prizes []*entity.LotteryPrize
	if err := g.Model("lottery_prize").Ctx(ctx).
		Where("site_id", ltSiteId).Where("activity_id", act.Id).Where("status", 1).
		OrderAsc("rank").OrderAsc("id").Scan(&prizes); err != nil {
		return nil, err
	}
	out := &service.InfoDTO{
		Activity: actDTO(act), Prizes: make([]*service.PrizeDTO, 0, len(prizes)),
		LoggedIn: userId > 0,
	}
	for _, p := range prizes {
		out.Prizes = append(out.Prizes, prizeDTO(p))
	}
	// 未登录时给出"满次数"的展示值, 让前端能直接渲染活动页。
	out.FreeLeft, out.DrawLeft = leftOf(act, 0, 0)
	if userId > 0 {
		total, free, err := dailyUsed(ctx, nil, userId, act.Id)
		if err != nil {
			return nil, err
		}
		out.FreeLeft, out.DrawLeft = leftOf(act, total, free)
		if out.Balance, err = userBalance(ctx, nil, userId); err != nil {
			return nil, err
		}
	}
	// 跑马灯用真实中奖记录(tianbi 是前端写死的假数据, 运营一改奖品就穿帮)。
	// 过滤掉"谢谢参与": 那不是中奖, 播出来只会劝退。
	var feed []*entity.LotteryHistory
	if err := g.Model("lottery_history").Ctx(ctx).
		Where("site_id", ltSiteId).Where("activity_id", act.Id).
		WhereNot("prize_type", service.PrizeThanks).
		OrderDesc("id").Limit(20).Scan(&feed); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(feed))
	for _, h := range feed {
		ids = append(ids, h.UserId)
	}
	nameMap := nicknames(ctx, ids)
	out.Marquee = make([]*service.HistoryDTO, 0, len(feed))
	for _, h := range feed {
		d := histDTO(h)
		d.Nickname = maskName(nameMap[h.UserId])
		d.UserId = 0 // 跑马灯不下发用户ID
		out.Marquee = append(out.Marquee, d)
	}
	return out, nil
}

// pickWeighted 加权随机: 按 odds 做前缀和, 取 [0,total) 的随机数落在哪个区间就是哪个奖。
//
// 为什么用 crypto/rand 而不是 math/rand:
//   - math/rand 是确定性的线性同余/PCG 序列, 种子一旦被猜到(默认种子在老版本里是固定的,
//     即便播种也常用时间戳), 攻击者可以预测下一次的随机数, 从而挑时机抽走大奖 —— 抽奖直接
//     关联真金白银的奖品, 属于"对抗性场景", 必须用密码学安全的随机源;
//   - math/rand 的全局源还需要自己保证并发安全, crypto/rand 天生并发安全。
//   - 代价是每次要一次系统调用, 但抽奖 QPS 远达不到成为瓶颈的量级。
//
// odds <= 0 的奖品视为"不参与抽取"(运营想临时下掉某个奖但保留配置), 直接跳过。
func pickWeighted(pool []*entity.LotteryPrize) (*entity.LotteryPrize, error) {
	cand := make([]*entity.LotteryPrize, 0, len(pool))
	prefix := make([]int64, 0, len(pool))
	var total int64
	for _, p := range pool {
		if p.Odds <= 0 {
			continue
		}
		total += int64(p.Odds)
		cand = append(cand, p)
		prefix = append(prefix, total)
	}
	if total <= 0 {
		return nil, gerror.New("奖池为空, 请联系管理员配置奖品")
	}
	n, err := crand.Int(crand.Reader, big.NewInt(total))
	if err != nil {
		return nil, err
	}
	r := n.Int64()
	for i, up := range prefix {
		if r < up {
			return cand[i], nil
		}
	}
	return cand[len(cand)-1], nil // 理论不可达, 兜底防越界
}

// thanksPrize 取本活动配置的"谢谢参与"奖品, 用于库存/券被抢空时的降级。
// 没配的话返回 nil, 调用方用一个虚拟奖品(prize_id=0)顶上, 保证抽奖流程不因配置缺失而失败。
func thanksPrize(ctx context.Context, tx gdb.TX, activityId int64) *entity.LotteryPrize {
	var p *entity.LotteryPrize
	_ = tx.Model("lottery_prize").Ctx(ctx).
		Where("site_id", ltSiteId).Where("activity_id", activityId).
		Where("type", service.PrizeThanks).Where("status", 1).
		OrderAsc("id").Scan(&p)
	return p
}

// takeStock 库存条件递减(防超发): stock=-1 不限量只累加已发数; stock>0 才减。
// 返回 false 表示这一份被别人抢走了。
func takeStock(ctx context.Context, tx gdb.TX, prizeId int64) (bool, error) {
	res, err := tx.Model("lottery_prize").Ctx(ctx).
		Where("id", prizeId).Where("stock = -1 OR stock > 0").
		Data("stock = CASE WHEN stock > 0 THEN stock - 1 ELSE stock END," +
			" awarded = awarded + 1, updated_at = now()").Update()
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// returnStock 归还库存(降级时用): 已经条件递减过的那一份没真发出去, 要放回去,
// 否则 awarded 会虚高、限量奖品会凭空少一份。事务内执行, 与降级本身同生共死。
func returnStock(ctx context.Context, tx gdb.TX, prizeId int64) error {
	if prizeId <= 0 {
		return nil
	}
	_, err := tx.Model("lottery_prize").Ctx(ctx).Where("id", prizeId).
		Data("stock = CASE WHEN stock >= 0 THEN stock + 1 ELSE stock END," +
			" awarded = GREATEST(awarded - 1, 0), updated_at = now()").Update()
	return err
}

// grantVip 发 VIP 天数: 已有未过期的 vip_log 则从其 end_at 续期, 否则从当下起算。
// start_at/end_at 是秒级时间戳(与 00004 建表 + paywall.IsVipActive 的约定一致)。
func grantVip(ctx context.Context, tx gdb.TX, userId int64, days int) error {
	if days <= 0 {
		return gerror.New("VIP奖品天数配置非法")
	}
	now := gtime.Now().Unix()
	start := now
	one, err := tx.Model("vip_log").Ctx(ctx).
		Where("site_id", ltSiteId).Where("user_id", userId).Where("end_at > ?", now).
		Fields("end_at").OrderDesc("end_at").One()
	if err != nil {
		return err
	}
	if !one.IsEmpty() {
		start = one["end_at"].Int64()
	}
	_, err = tx.Model("vip_log").Ctx(ctx).Data(g.Map{
		"site_id": ltSiteId, "user_id": userId, "package_id": 0, "days": days,
		"price": 0, "start_at": start, "end_at": start + int64(days)*86400,
	}).Insert()
	return err
}

// issueCoupon 事务内直接发券。
//
// 为什么不复用 coupon 模块的 Receive/issue:
//   - 那两个方法内部自己 `g.DB().Transaction(...)` 开事务。在抽奖事务里调它, GoFrame 会开
//     嵌套事务(SAVEPOINT), 它内部 commit 掉的 savepoint 并不代表真的落库, 而一旦外层因为
//     后续步骤失败而回滚, 券就被一起撤销 —— 更糟的是, 如果它内部走的是新连接, 券会脱离
//     外层事务先行落库, 造成"抽奖回滚了券还在"的对不上账;
//   - coupon.issue 还带"每人限领 per_limit"校验, 那是给用户主动领券用的风控, 中奖发放不该
//     被它拦(中三次奖就该有三张券)。
//
// 所以这里直接复用 coupon 的表结构, 在同一个 tx 内做「coupon_tpl 总量条件递增 + user_coupon
// 快照插入」, 与 coupon.issue 的语义保持一致(面值等字段快照, 模板改价不影响已发放)。
// 券模板缺失或已发完时返回 errCouponOut, 由调用方降级为"谢谢参与", 而不是让整次抽奖失败。
func issueCoupon(ctx context.Context, tx gdb.TX, userId, tplId int64) (int64, error) {
	if tplId <= 0 {
		return 0, errCouponOut
	}
	var tpl *entity.CouponTpl
	if err := tx.Model("coupon_tpl").Ctx(ctx).
		Where("site_id", ltSiteId).Where("id", tplId).Where("status", 1).Scan(&tpl); err != nil {
		return 0, err
	}
	if tpl == nil {
		return 0, errCouponOut
	}
	res, err := tx.Model("coupon_tpl").Ctx(ctx).
		Where("id", tplId).Where("total = -1 OR issued < total").
		Data(g.Map{
			"issued":     &gdb.Counter{Field: "issued", Value: 1},
			"updated_at": gtime.Now(),
		}).Update()
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, errCouponOut
	}
	expireDay := tpl.ExpireDay
	if expireDay <= 0 {
		expireDay = 7
	}
	return tx.Model("user_coupon").Ctx(ctx).Data(g.Map{
		"site_id": ltSiteId, "user_id": userId, "tpl_id": tplId, "name": tpl.Name,
		"type": tpl.Type, "scene": tpl.Scene, "face_value": tpl.FaceValue,
		"discount": tpl.Discount, "threshold": tpl.Threshold, "max_deduct": tpl.MaxDeduct,
		"status": 1, "expire_at": gtime.Now().AddDate(0, 0, expireDay),
	}).InsertAndGetId()
}

// Draw 抽奖: 全程一个事务, 详见包注释。
func (s *sLottery) Draw(ctx context.Context, userId int64, lotteryType, payType int) (*service.DrawDTO, error) {
	if userId <= 0 {
		return nil, gerror.New("未登录")
	}
	act, err := s.activityByType(ctx, lotteryType)
	if err != nil {
		return nil, err
	}
	if payType != service.PayFree && payType != service.PayGold {
		return nil, gerror.New("消耗方式非法")
	}
	if payType == service.PayGold && act.PayType != service.PayGold {
		return nil, gerror.New("该活动只支持免费次数抽奖")
	}
	if payType == service.PayFree && act.DailyFree <= 0 {
		return nil, gerror.New("该活动没有免费次数")
	}

	out := &service.DrawDTO{}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 0. 用户行锁: 把同一用户的并发抽奖串行化(修 tianbi 双击重复扣费)。
		//    锁 users 行而不是活动行, 是为了只挡住"同一个人"的并发, 不影响别人抽奖的吞吐。
		lock, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("id").LockUpdate().One()
		if err != nil {
			return err
		}
		if lock.IsEmpty() {
			return gerror.New("用户不存在")
		}
		// 1. 每日次数校验(行锁之后统计, 结果不会被并发击穿)
		total, free, err := dailyUsed(ctx, tx, userId, act.Id)
		if err != nil {
			return err
		}
		if act.DailyLimit > 0 && total >= act.DailyLimit {
			return gerror.New("今日抽奖次数已达上限")
		}
		if payType == service.PayFree && free >= act.DailyFree {
			return gerror.New("今日免费次数已用完")
		}
		// 2. 扣费: 金币走 balance.Deduct(条件扣款, SQL 层防透支); 免费次数由第 1 步的条件判断兜住。
		cost := 0.0
		if payType == service.PayGold {
			cost = act.CostGold
			if cost <= 0 {
				return gerror.New("活动未配置抽奖金币价")
			}
			if err := balance.Deduct(ctx, tx, userId, cost,
				balance.SceneLotteryCost, gconv.String(act.Id), "抽奖消耗:"+act.Name); err != nil {
				return err
			}
		}
		// 3. 加权随机选奖(只取启用中、且还有库存的奖品进池)
		var pool []*entity.LotteryPrize
		if err := tx.Model("lottery_prize").Ctx(ctx).
			Where("site_id", ltSiteId).Where("activity_id", act.Id).Where("status", 1).
			Where("stock = -1 OR stock > 0").OrderAsc("id").Scan(&pool); err != nil {
			return err
		}
		win, err := pickWeighted(pool)
		if err != nil {
			return err
		}
		remark := ""
		// 4. 库存条件递减。抢不到(别人先一步拿走最后一份)时**降级为"谢谢参与"**, 而不是重抽。
		//    取舍: 重抽要在事务里循环 SELECT+UPDATE, 极端情况(奖池只剩一个抢手奖)会长时间
		//    占着行锁和事务, 放大锁冲突; 而且重抽还要重新算权重, 逻辑复杂度陡增。降级只多一次
		//    UPDATE, 用户侧体验是"没中奖", 与真实概率的偏差只在库存耗尽的边界上, 完全可接受。
		ok, err := takeStock(ctx, tx, win.Id)
		if err != nil {
			return err
		}
		if !ok {
			remark = "奖品[" + win.Name + "]库存已被抽完, 降级为谢谢参与"
			win = fallbackThanks(ctx, tx, act.Id)
		}
		// 5. 发奖分派(全部在事务内)
		status := service.HistoryDone
		switch win.Type {
		case service.PrizeGold:
			if win.Amount > 0 {
				if err := balance.Add(ctx, tx, userId, win.Amount,
					balance.SceneLotteryPrize, gconv.String(win.Id), "抽奖奖励:"+win.Name); err != nil {
					return err
				}
			}
		case service.PrizeVip:
			if err := grantVip(ctx, tx, userId, int(win.Amount)); err != nil {
				return err
			}
		case service.PrizeCoupon:
			cid, err := issueCoupon(ctx, tx, userId, win.CouponTplId)
			switch {
			case err == nil:
				remark = "优惠券ID:" + gconv.String(cid)
			case gerror.Is(err, errCouponOut):
				// 与库存降级同一套取舍: 券发完了不让整次抽奖失败, 降级成谢谢参与。
				// 这份券奖的库存在第 4 步已经扣过, 没真发出去就得还回来。
				if e := returnStock(ctx, tx, win.Id); e != nil {
					return e
				}
				remark = "奖品[" + win.Name + "]券已发完, 降级为谢谢参与"
				win = fallbackThanks(ctx, tx, act.Id)
			default:
				return err
			}
		case service.PrizeGoods:
			status = service.HistoryWaitShip // 实物待发货, 等用户填地址
		case service.PrizeThanks:
			// 不发任何东西
		}
		// 6. 写中奖记录(与发奖同事务, 写不进去就整体回滚)
		hid, err := tx.Model("lottery_history").Ctx(ctx).Data(g.Map{
			"site_id": ltSiteId, "user_id": userId, "activity_id": act.Id,
			"lottery_type": act.LotteryType, "pay_type": payType, "cost_gold": cost,
			"prize_id": win.Id, "prize_name": win.Name, "prize_type": win.Type,
			"prize_amount": win.Amount, "status": status, "remark": remark,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		// 7. 实物奖顺手建一张待填写的收货单(唯一索引保证一奖一单)
		if win.Type == service.PrizeGoods {
			if _, err := tx.Model("prize_addr").Ctx(ctx).Data(g.Map{
				"site_id": ltSiteId, "history_id": hid, "user_id": userId,
				"delivery_status": service.AddrWaitFill,
			}).Insert(); err != nil {
				return err
			}
		}
		bal, err := userBalance(ctx, tx, userId)
		if err != nil {
			return err
		}
		freeLeft, drawLeft := leftOf(act, total+1, free+boolInt(payType == service.PayFree))
		out = &service.DrawDTO{
			HistoryId: hid, PrizeId: win.Id, PrizeName: win.Name, PrizeType: win.Type,
			PrizeAmount: win.Amount, PrizeCover: win.Cover, PrizeDesc: win.Desc,
			Status: status, NeedAddr: win.Type == service.PrizeGoods, Remark: remark,
			FreeLeft: freeLeft, DrawLeft: drawLeft, Balance: bal, CostGold: cost,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// fallbackThanks 降级奖品: 优先用活动配置的"谢谢参与", 没配就用一个 prize_id=0 的虚拟奖,
// 保证 history 一定写得出去(奖品配置缺失不该让用户的这次抽奖凭空消失)。
func fallbackThanks(ctx context.Context, tx gdb.TX, activityId int64) *entity.LotteryPrize {
	if p := thanksPrize(ctx, tx, activityId); p != nil {
		p.Amount = 0
		return p
	}
	return &entity.LotteryPrize{Id: 0, Name: "谢谢参与", Type: service.PrizeThanks}
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (s *sLottery) My(ctx context.Context, userId int64, lotteryType, page, size int) ([]*service.HistoryDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("lottery_history").Ctx(ctx).Where("site_id", ltSiteId).Where("user_id", userId)
	if lotteryType > 0 {
		m = m.Where("lottery_type", lotteryType)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.LotteryHistory
	if err := m.Clone().OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.HistoryDTO, 0, len(list))
	for _, h := range list {
		out = append(out, histDTO(h))
	}
	return out, total, nil
}

// FillAddr 用户自助填收货地址(tianbi 没有这个接口, 实物奖只能线下找客服登记)。
// 三重限制全部写进 SQL 条件, 用影响行数判定, 不做"先查再写":
// 只能填自己的(user_id) / 必须是实物奖(先校验 history.prize_type) / 只有待填写才能填(delivery_status=0)。
func (s *sLottery) FillAddr(ctx context.Context, userId, historyId int64, receiver, phone, address string) error {
	if historyId <= 0 {
		return gerror.New("中奖记录ID必填")
	}
	if strings.TrimSpace(receiver) == "" || strings.TrimSpace(phone) == "" || strings.TrimSpace(address) == "" {
		return gerror.New("收货人/手机号/地址均不能为空")
	}
	var h *entity.LotteryHistory
	if err := g.Model("lottery_history").Ctx(ctx).
		Where("site_id", ltSiteId).Where("id", historyId).Where("user_id", userId).
		Scan(&h); err != nil {
		return err
	}
	if h == nil {
		return gerror.New("中奖记录不存在")
	}
	if h.PrizeType != service.PrizeGoods {
		return gerror.New("该奖品无需填写收货地址")
	}
	res, err := g.Model("prize_addr").Ctx(ctx).
		Where("site_id", ltSiteId).Where("history_id", historyId).Where("user_id", userId).
		Where("delivery_status", service.AddrWaitFill).
		Data(g.Map{
			"receiver": receiver, "phone": phone, "address": address,
			"delivery_status": service.AddrWaitShip, "updated_at": gtime.Now(),
		}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("收货地址已填写或该奖品不可填写")
	}
	return nil
}

// addrList 收货单列表(前后台共用), 顺带补上奖品名。
func addrList(ctx context.Context, m *gdb.Model, page, size int) ([]*service.AddrDTO, int, error) {
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.PrizeAddr
	if err := m.Clone().OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	hids := make([]int64, 0, len(list))
	uids := make([]int64, 0, len(list))
	for _, a := range list {
		hids = append(hids, a.HistoryId)
		uids = append(uids, a.UserId)
	}
	prizeName := map[int64]string{}
	if len(hids) > 0 {
		all, err := g.Model("lottery_history").Ctx(ctx).WhereIn("id", hids).
			Fields("id,prize_name").All()
		if err != nil {
			return nil, 0, err
		}
		for _, r := range all {
			prizeName[r["id"].Int64()] = r["prize_name"].String()
		}
	}
	nameMap := nicknames(ctx, uids)
	out := make([]*service.AddrDTO, 0, len(list))
	for _, a := range list {
		created := ""
		if a.CreatedAt != nil {
			created = a.CreatedAt.String()
		}
		out = append(out, &service.AddrDTO{
			Id: a.Id, HistoryId: a.HistoryId, UserId: a.UserId, Nickname: nameMap[a.UserId],
			PrizeName: prizeName[a.HistoryId], Receiver: a.Receiver, Phone: a.Phone,
			Address: a.Address, DeliveryStatus: a.DeliveryStatus, ExpressNo: a.ExpressNo,
			CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sLottery) MyAddrs(ctx context.Context, userId int64, page, size int) ([]*service.AddrDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("prize_addr").Ctx(ctx).Where("site_id", ltSiteId).Where("user_id", userId)
	return addrList(ctx, m, page, size)
}

// ---------------------------------------------------------------- 后台: 活动

func (s *sLottery) Activities(ctx context.Context, status int) ([]*service.ActivityDTO, error) {
	m := g.Model("lottery_activity").Ctx(ctx).Where("site_id", ltSiteId)
	if status >= 0 {
		m = m.Where("status", status)
	}
	var list []*entity.LotteryActivity
	if err := m.OrderAsc("lottery_type").OrderAsc("id").Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.ActivityDTO, 0, len(list))
	for _, a := range list {
		out = append(out, actDTO(a))
	}
	return out, nil
}

func (s *sLottery) ActivityCreate(ctx context.Context, in service.ActivityInput) (int64, error) {
	if strings.TrimSpace(in.Name) == "" {
		return 0, gerror.New("活动名不能为空")
	}
	if in.LotteryType != service.TypeVipDay && in.LotteryType != service.TypeWelfare {
		return 0, gerror.New("玩法类型非法")
	}
	if in.PayType != service.PayFree && in.PayType != service.PayGold {
		in.PayType = service.PayFree
	}
	if in.PayType == service.PayGold && in.CostGold <= 0 {
		return 0, gerror.New("金币抽奖必须配置单次金币价")
	}
	// 唯一索引兜底, 这里先查一次只是为了给运营一句人话而不是数据库报错。
	n, err := g.Model("lottery_activity").Ctx(ctx).
		Where("site_id", ltSiteId).Where("lottery_type", in.LotteryType).Count()
	if err != nil {
		return 0, err
	}
	if n > 0 {
		return 0, gerror.New("该玩法已存在活动, 请直接编辑")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	return g.Model("lottery_activity").Ctx(ctx).Data(g.Map{
		"site_id": ltSiteId, "name": in.Name, "lottery_type": in.LotteryType,
		"pay_type": in.PayType, "cost_gold": in.CostGold, "daily_free": in.DailyFree,
		"daily_limit": in.DailyLimit, "notice": in.Notice, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sLottery) ActivityUpdate(ctx context.Context, in service.ActivityInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{
		"cost_gold": in.CostGold, "daily_free": in.DailyFree,
		"daily_limit": in.DailyLimit, "updated_at": gtime.Now(),
	}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.PayType == service.PayFree || in.PayType == service.PayGold {
		data["pay_type"] = in.PayType
	}
	if in.Notice != "" {
		data["notice"] = in.Notice
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("lottery_activity").Ctx(ctx).
		Where("site_id", ltSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

// ActivityDelete 删活动连带删奖品(奖品脱离活动就是孤儿数据), 但保留 lottery_history:
// 中奖记录是资产发放的凭证, 任何时候都不能因为运营删活动而消失。
func (s *sLottery) ActivityDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("lottery_prize").Ctx(ctx).
			Where("site_id", ltSiteId).Where("activity_id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("lottery_activity").Ctx(ctx).
			Where("site_id", ltSiteId).Where("id", id).Delete()
		return err
	})
}

// ---------------------------------------------------------------- 后台: 奖品

func (s *sLottery) Prizes(ctx context.Context, activityId int64) ([]*service.PrizeDTO, error) {
	m := g.Model("lottery_prize").Ctx(ctx).Where("site_id", ltSiteId)
	if activityId > 0 {
		m = m.Where("activity_id", activityId)
	}
	var list []*entity.LotteryPrize
	if err := m.OrderAsc("activity_id").OrderAsc("rank").OrderAsc("id").Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.PrizeDTO, 0, len(list))
	for _, p := range list {
		out = append(out, prizeDTO(p))
	}
	return out, nil
}

func (s *sLottery) PrizeCreate(ctx context.Context, in service.PrizeInput) (int64, error) {
	if in.ActivityId <= 0 {
		return 0, gerror.New("活动ID必填")
	}
	if strings.TrimSpace(in.Name) == "" {
		return 0, gerror.New("奖品名不能为空")
	}
	if in.Type < service.PrizeGold || in.Type > service.PrizeThanks {
		return 0, gerror.New("奖品类型非法")
	}
	if in.Type == service.PrizeCoupon && in.CouponTplId <= 0 {
		return 0, gerror.New("优惠券奖品必须指定券模板")
	}
	if in.Stock < -1 {
		in.Stock = -1
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	return g.Model("lottery_prize").Ctx(ctx).Data(g.Map{
		"site_id": ltSiteId, "activity_id": in.ActivityId, "name": in.Name,
		"cover": in.Cover, "desc": in.Desc, "type": in.Type, "amount": in.Amount,
		"coupon_tpl_id": in.CouponTplId, "odds": in.Odds, "stock": in.Stock,
		"rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sLottery) PrizeUpdate(ctx context.Context, in service.PrizeInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	// odds/stock/amount/rank 允许被改成 0 或 -1, 所以无条件写入; 其余空值视为"不改"。
	data := g.Map{
		"amount": in.Amount, "odds": in.Odds, "stock": in.Stock,
		"rank": in.Rank, "updated_at": gtime.Now(),
	}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Cover != "" {
		data["cover"] = in.Cover
	}
	if in.Desc != "" {
		data["desc"] = in.Desc
	}
	if in.Type >= service.PrizeGold && in.Type <= service.PrizeThanks {
		data["type"] = in.Type
	}
	if in.CouponTplId > 0 {
		data["coupon_tpl_id"] = in.CouponTplId
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("lottery_prize").Ctx(ctx).
		Where("site_id", ltSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sLottery) PrizeDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("lottery_prize").Ctx(ctx).
		Where("site_id", ltSiteId).Where("id", id).Delete()
	return err
}

// ---------------------------------------------------------------- 后台: 记录 / 发货

func (s *sLottery) Histories(ctx context.Context, f service.HistoryFilter) ([]*service.HistoryDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("lottery_history").Ctx(ctx).Where("site_id", ltSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.LotteryType > 0 {
		m = m.Where("lottery_type", f.LotteryType)
	}
	if f.PrizeType > 0 {
		m = m.Where("prize_type", f.PrizeType)
	}
	if f.Status > 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.LotteryHistory
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	uids := make([]int64, 0, len(list))
	for _, h := range list {
		uids = append(uids, h.UserId)
	}
	nameMap := nicknames(ctx, uids)
	out := make([]*service.HistoryDTO, 0, len(list))
	for _, h := range list {
		d := histDTO(h)
		d.Nickname = nameMap[h.UserId]
		out = append(out, d)
	}
	return out, total, nil
}

func (s *sLottery) Addrs(ctx context.Context, f service.AddrFilter) ([]*service.AddrDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("prize_addr").Ctx(ctx).Where("site_id", ltSiteId)
	if f.DeliveryStatus >= 0 {
		m = m.Where("delivery_status", f.DeliveryStatus)
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	return addrList(ctx, m, f.Page, f.Size)
}

// Ship 发货: 收货单 1→2 用条件更新 + 影响行数判定(重复发货/未填地址就发货都会被拦),
// 同时把对应中奖记录置为已发货, 两处状态在一个事务里保持一致。
func (s *sLottery) Ship(ctx context.Context, addrId int64, expressNo string) error {
	if addrId <= 0 {
		return gerror.New("ID非法")
	}
	if strings.TrimSpace(expressNo) == "" {
		return gerror.New("快递单号必填")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var addr *entity.PrizeAddr
		if err := tx.Model("prize_addr").Ctx(ctx).
			Where("site_id", ltSiteId).Where("id", addrId).Scan(&addr); err != nil {
			return err
		}
		if addr == nil {
			return gerror.New("收货单不存在")
		}
		res, err := tx.Model("prize_addr").Ctx(ctx).
			Where("id", addrId).Where("delivery_status", service.AddrWaitShip).
			Data(g.Map{
				"delivery_status": service.AddrShipped, "express_no": expressNo,
				"updated_at": gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("该收货单不是待发货状态")
		}
		_, err = tx.Model("lottery_history").Ctx(ctx).
			Where("id", addr.HistoryId).Where("status", service.HistoryWaitShip).
			Data(g.Map{"status": service.HistoryShipped, "updated_at": gtime.Now()}).Update()
		return err
	})
}
