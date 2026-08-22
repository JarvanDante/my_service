// Package balance 金币余额的统一出入口(事务内使用)。
//
// 约定(与 checkin/redeemcode/redeemgoods 既有实现一致):
//   - 余额在 users.balance, numeric(14,2);
//   - 任何变动都必须写 user_balance_log 流水(direction 1收/2支, before/after 快照);
//   - 扣款一律走「WHERE balance >= ?」条件更新 + RowsAffected 判定, 从 SQL 层防透支,
//     不做「先查再扣」(那是 tianbi 的 TOCTOU 竞态老问题)。
//
// 所有函数都要求调用方已经开启事务并传入 tx, 保证扣款与业务写入同生共死。
package balance

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const siteId = 1 // 单站点样板

// 常用流水场景(scene), 新增场景时在此登记, 便于后台按场景筛选对账。
const (
	SceneWithdrawFreeze = "withdraw_freeze" // 提现冻结(申请即扣)
	SceneWithdrawRefund = "withdraw_refund" // 提现退款(拒绝/撤回)
	SceneAdminAdjust    = "admin_adjust"    // 后台人工调账
	SceneContentBuy     = "content_buy"     // 内容购买(漫画/小说/图集/视频)
	SceneLotteryCost    = "lottery_cost"    // 抽奖消耗
	SceneLotteryPrize   = "lottery_prize"   // 抽奖奖励
	SceneAiCost         = "ai_cost"         // AI生成任务扣费(提交/重试)
	SceneAiRefund       = "ai_refund"       // AI生成任务退款(提交失败/生成失败/取消/人工退)
	SceneRedeemCode     = "redeem_code"     // 兑换码到账
)

// ErrInsufficient 余额不足(供调用方判定并转换成业务提示)。
var ErrInsufficient = gerror.New("金币余额不足")

// Deduct 扣款: 条件更新防透支 + 写支出流水。amount 必须 > 0。
func Deduct(ctx context.Context, tx gdb.TX, userId int64, amount float64, scene, refId, remark string) error {
	if amount <= 0 {
		return gerror.New("扣款金额需大于0")
	}
	res, err := tx.Model("users").Ctx(ctx).
		Where("id", userId).Where("balance >= ?", amount).
		Data(g.Map{"balance": &gdb.Counter{Field: "balance", Value: -amount}}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrInsufficient
	}
	after, err := current(ctx, tx, userId)
	if err != nil {
		return err
	}
	return writeLog(ctx, tx, userId, 2, amount, after+amount, after, scene, refId, remark)
}

// Add 加款: 加余额 + 写收入流水。amount 必须 > 0。
func Add(ctx context.Context, tx gdb.TX, userId int64, amount float64, scene, refId, remark string) error {
	if amount <= 0 {
		return gerror.New("入账金额需大于0")
	}
	if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
		Data(g.Map{"balance": &gdb.Counter{Field: "balance", Value: amount}}).Update(); err != nil {
		return err
	}
	after, err := current(ctx, tx, userId)
	if err != nil {
		return err
	}
	return writeLog(ctx, tx, userId, 1, amount, after-amount, after, scene, refId, remark)
}

// current 取扣/加之后的余额, 作为流水的 balance_after。
func current(ctx context.Context, tx gdb.TX, userId int64) (float64, error) {
	one, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").One()
	if err != nil {
		return 0, err
	}
	if one.IsEmpty() {
		return 0, gerror.New("用户不存在")
	}
	return one["balance"].Float64(), nil
}

func writeLog(ctx context.Context, tx gdb.TX, userId int64, direction int,
	amount, before, after float64, scene, refId, remark string) error {
	_, err := tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
		"site_id": siteId, "user_id": userId, "direction": direction, "scene": scene,
		"amount": amount, "balance_before": before, "balance_after": after,
		"ref_id": refId, "remark": remark,
	}).Insert()
	return err
}
