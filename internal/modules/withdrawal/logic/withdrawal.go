// Package logic 提现业务(移植自 tianbi withdrawal + bankcard)。
//
// 与 tianbi 的关键差异(都是修 bug, 不是换写法):
//  1. tianbi 的"查余额 → 扣余额"分两步, 存在 TOCTOU 竞态; 这里全程「事务 + WHERE balance>=? 条件更新」;
//  2. tianbi 的状态迁移直接 $set, 并发下可能重复退款; 这里每次迁移都带 `WHERE status=旧值`,
//     用 RowsAffected==0 判定"已被别人处理过", 保证退款只发生一次;
//  3. 手续费、最低/最高、倍数、日限全部服务端从 app_config 取, 不信客户端传值。
package logic

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/withdrawal/service"
	"github.com/JarvanDante/my_service/internal/shared/appcfg"
	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const wdSiteId = 1

type sWithdrawal struct{}

func New() service.IWithdrawal { return &sWithdrawal{} }

// round2 金额一律保留两位, 与 numeric(14,2) 对齐, 避免浮点尾差导致对账不平。
func round2(f float64) float64 { return math.Round(f*100) / 100 }

type accountSnap struct {
	AccountType int    `json:"account_type"`
	AccountName string `json:"account_name"`
	AccountNo   string `json:"account_no"`
	BankName    string `json:"bank_name"`
}

func toOrderDTO(r *entity.Withdrawal) *service.OrderDTO {
	var snap accountSnap
	if r.AccountInfo != "" {
		_ = json.Unmarshal([]byte(r.AccountInfo), &snap)
	}
	fmtTime := func(t *gtime.Time) string {
		if t == nil {
			return ""
		}
		return t.String()
	}
	return &service.OrderDTO{
		Id: r.Id, TradeNo: r.TradeNo, UserId: r.UserId, Amount: r.Amount, Fee: r.Fee,
		RealAmount: r.RealAmount, FeeRate: r.FeeRate, BalanceAfter: r.BalanceAfter,
		AccountName: snap.AccountName, AccountNo: snap.AccountNo, BankName: snap.BankName,
		Status: r.Status, AuditBy: r.AuditBy, AuditAt: fmtTime(r.AuditAt),
		PaidAt: fmtTime(r.PaidAt), Remark: r.Remark, PayVoucher: r.PayVoucher,
		CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sWithdrawal) Config(ctx context.Context, userId int64) (*service.ConfigDTO, error) {
	out := &service.ConfigDTO{
		Open:       appcfg.Bool(ctx, "withdraw_open", true),
		MinAmount:  appcfg.Float(ctx, "withdraw_min", 100),
		MaxAmount:  appcfg.Float(ctx, "withdraw_max", 50000),
		Multiple:   appcfg.Float(ctx, "withdraw_multiple", 10),
		FeeRate:    appcfg.Float(ctx, "withdraw_fee_rate", 2),
		DailyLimit: appcfg.Int(ctx, "withdraw_daily_limit", 3),
	}
	if userId > 0 {
		bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
		if err != nil {
			return nil, err
		}
		if bal != nil {
			out.Balance = bal.Float64()
		}
		frozen, err := g.Model("withdrawal").Ctx(ctx).
			Where("site_id", wdSiteId).Where("user_id", userId).
			WhereIn("status", g.Slice{service.StatusApplying, service.StatusPassed}).Sum("amount")
		if err != nil {
			return nil, err
		}
		out.Frozen = frozen
		used, err := s.todayCount(ctx, userId)
		if err != nil {
			return nil, err
		}
		out.DailyUsed = used
	}
	return out, nil
}

// todayCount 今日已申请笔数(撤回的也算, 防止刷单绕过日限)。
func (s *sWithdrawal) todayCount(ctx context.Context, userId int64) (int, error) {
	start := gtime.Now().StartOfDay()
	return g.Model("withdrawal").Ctx(ctx).
		Where("site_id", wdSiteId).Where("user_id", userId).
		Where("created_at >= ?", start).Count()
}

func (s *sWithdrawal) Apply(ctx context.Context, userId, cardId int64, amount float64) (*service.ApplyResult, error) {
	cfg, err := s.Config(ctx, userId)
	if err != nil {
		return nil, err
	}
	if !cfg.Open {
		return nil, gerror.New("提现功能维护中, 请稍后再试")
	}
	amount = round2(amount)
	if amount < cfg.MinAmount {
		return nil, gerror.Newf("单笔最低提现 %.2f", cfg.MinAmount)
	}
	if cfg.MaxAmount > 0 && amount > cfg.MaxAmount {
		return nil, gerror.Newf("单笔最高提现 %.2f", cfg.MaxAmount)
	}
	if cfg.Multiple > 0 && math.Mod(amount, cfg.Multiple) != 0 {
		return nil, gerror.Newf("提现金额需为 %.0f 的整数倍", cfg.Multiple)
	}
	if cfg.DailyLimit > 0 && cfg.DailyUsed >= cfg.DailyLimit {
		return nil, gerror.Newf("今日提现次数已达上限(%d次)", cfg.DailyLimit)
	}
	// 频控: 与最近一单的间隔(移植自 tianbi 的 10 秒限制, 改为可配)
	interval := appcfg.Int(ctx, "withdraw_interval_sec", 10)
	if interval > 0 {
		var last *entity.Withdrawal
		if err := g.Model("withdrawal").Ctx(ctx).
			Where("site_id", wdSiteId).Where("user_id", userId).
			OrderDesc("id").Limit(1).Scan(&last); err != nil {
			return nil, err
		}
		if last != nil && last.CreatedAt != nil &&
			gtime.Now().Sub(last.CreatedAt).Seconds() < float64(interval) {
			return nil, gerror.Newf("操作过于频繁, 请 %d 秒后再试", interval)
		}
	}
	// 收款账户归属校验
	var card *entity.BankCard
	if err := g.Model("bank_card").Ctx(ctx).
		Where("site_id", wdSiteId).Where("id", cardId).Where("user_id", userId).
		Scan(&card); err != nil {
		return nil, err
	}
	if card == nil {
		return nil, gerror.New("收款账户不存在")
	}
	snap, _ := json.Marshal(accountSnap{
		AccountType: card.AccountType, AccountName: card.AccountName,
		AccountNo: card.AccountNo, BankName: card.BankName,
	})
	fee := round2(amount * cfg.FeeRate / 100)
	real := round2(amount - fee)
	if real <= 0 {
		return nil, gerror.New("扣除手续费后到账金额为0, 请提高提现金额")
	}
	tradeNo := "W" + gtime.Now().Format("YmdHis") + gconv.String(userId) + grand.Digits(4)

	out := &service.ApplyResult{TradeNo: tradeNo, Amount: amount, Fee: fee, RealAmount: real}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 冻结: 条件扣款(防透支) + 支出流水
		if err := balance.Deduct(ctx, tx, userId, amount,
			balance.SceneWithdrawFreeze, tradeNo, "提现冻结"); err != nil {
			return err
		}
		// 2. 取扣后余额作为快照
		one, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").One()
		if err != nil {
			return err
		}
		after := one["balance"].Float64()
		// 3. 建单
		id, err := tx.Model("withdrawal").Ctx(ctx).Data(g.Map{
			"site_id": wdSiteId, "trade_no": tradeNo, "user_id": userId,
			"amount": amount, "fee": fee, "real_amount": real, "fee_rate": cfg.FeeRate,
			"balance_after": after, "account_info": string(snap),
			"status": service.StatusApplying,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		out.Id = id
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *sWithdrawal) listBy(ctx context.Context, f service.ListFilter) (*gdb.Model, error) {
	m := g.Model("withdrawal").Ctx(ctx).Where("site_id", wdSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Status > 0 {
		m = m.Where("status", f.Status)
	}
	return m, nil
}

func (s *sWithdrawal) My(ctx context.Context, f service.ListFilter) ([]*service.OrderDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m, err := s.listBy(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Withdrawal
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.OrderDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toOrderDTO(r))
	}
	return out, total, nil
}

func (s *sWithdrawal) List(ctx context.Context, f service.ListFilter) ([]*service.OrderDTO, int, float64, int, error) {
	list, total, err := s.My(ctx, f)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	m, _ := s.listBy(ctx, f)
	sum, err := m.Clone().Sum("amount")
	if err != nil {
		return nil, 0, 0, 0, err
	}
	pending, err := g.Model("withdrawal").Ctx(ctx).
		Where("site_id", wdSiteId).Where("status", service.StatusApplying).Count()
	if err != nil {
		return nil, 0, 0, 0, err
	}
	// 补昵称(后台列表要展示)
	ids := make([]int64, 0, len(list))
	for _, d := range list {
		ids = append(ids, d.UserId)
	}
	if len(ids) > 0 {
		var users []*entity.Users
		if err := g.Model("users").Ctx(ctx).WhereIn("id", ids).
			Fields("id,nickname,username").Scan(&users); err == nil {
			nameOf := make(map[int64]string, len(users))
			for _, u := range users {
				n := u.Nickname
				if n == "" {
					n = u.Username
				}
				nameOf[u.Id] = n
			}
			for _, d := range list {
				d.Nickname = nameOf[d.UserId]
			}
		}
	}
	return list, total, sum, pending, nil
}

// transit 状态迁移的统一实现: 条件更新 + 可选退款, 全在一个事务里。
// refund=true 时把 amount 原路退回并写收入流水。
func (s *sWithdrawal) transit(ctx context.Context, id int64, userId int64,
	fromStatus []interface{}, toStatus int, refund bool, data g.Map, scene, remark string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var order *entity.Withdrawal
		if err := tx.Model("withdrawal").Ctx(ctx).
			Where("site_id", wdSiteId).Where("id", id).Scan(&order); err != nil {
			return err
		}
		if order == nil {
			return gerror.New("提现单不存在")
		}
		if userId > 0 && order.UserId != userId {
			return gerror.New("提现单不存在")
		}
		upd := g.Map{"status": toStatus, "updated_at": gtime.Now()}
		for k, v := range data {
			upd[k] = v
		}
		res, err := tx.Model("withdrawal").Ctx(ctx).
			Where("site_id", wdSiteId).Where("id", id).
			WhereIn("status", fromStatus).Data(upd).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("当前状态不允许该操作(可能已被处理)")
		}
		if refund {
			return balance.Add(ctx, tx, order.UserId, order.Amount, scene, order.TradeNo, remark)
		}
		return nil
	})
}

func (s *sWithdrawal) Cancel(ctx context.Context, userId, id int64) error {
	return s.transit(ctx, id, userId,
		[]interface{}{service.StatusApplying}, service.StatusCanceled, true,
		g.Map{"remark": "用户撤回"}, balance.SceneWithdrawRefund, "提现撤回退款")
}

func (s *sWithdrawal) Audit(ctx context.Context, adminId, id int64, pass bool, remark string) error {
	data := g.Map{"audit_by": adminId, "audit_at": gtime.Now(), "remark": remark}
	if pass {
		return s.transit(ctx, id, 0,
			[]interface{}{service.StatusApplying}, service.StatusPassed, false, data, "", "")
	}
	if strings.TrimSpace(remark) == "" {
		data["remark"] = "审核拒绝"
	}
	return s.transit(ctx, id, 0,
		[]interface{}{service.StatusApplying}, service.StatusRejected, true, data,
		balance.SceneWithdrawRefund, "提现被拒退款")
}

func (s *sWithdrawal) MarkPaid(ctx context.Context, adminId, id int64, voucher, remark string) error {
	data := g.Map{"paid_at": gtime.Now(), "pay_voucher": voucher, "audit_by": adminId}
	if remark != "" {
		data["remark"] = remark
	}
	return s.transit(ctx, id, 0,
		[]interface{}{service.StatusPassed}, service.StatusPaid, false, data, "", "")
}

func (s *sWithdrawal) RefundPaid(ctx context.Context, adminId, id int64, remark string) error {
	if strings.TrimSpace(remark) == "" {
		remark = "打款失败退款"
	}
	return s.transit(ctx, id, 0,
		[]interface{}{service.StatusPassed}, service.StatusRejected, true,
		g.Map{"audit_by": adminId, "audit_at": gtime.Now(), "remark": remark},
		balance.SceneWithdrawRefund, "提现打款失败退款")
}

// ---------------- 收款账户 ----------------

func toCardDTO(r *entity.BankCard) *service.CardDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.CardDTO{
		Id: r.Id, AccountType: r.AccountType, AccountName: r.AccountName,
		AccountNo: r.AccountNo, BankName: r.BankName, IsDefault: r.IsDefault, CreatedAt: created,
	}
}

func (s *sWithdrawal) CardList(ctx context.Context, userId int64) ([]*service.CardDTO, error) {
	var list []*entity.BankCard
	if err := g.Model("bank_card").Ctx(ctx).
		Where("site_id", wdSiteId).Where("user_id", userId).
		OrderDesc("is_default").OrderDesc("id").Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.CardDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toCardDTO(r))
	}
	return out, nil
}

// checkCardFields 移植自 tianbi: 户名/账号/开户行不得为空且不得含空格(空格是打款失败的高发原因)。
func checkCardFields(name, no, bank string, needAll bool) error {
	if needAll && (strings.TrimSpace(name) == "" || strings.TrimSpace(no) == "") {
		return gerror.New("开户人与账号不能为空")
	}
	for _, v := range []string{name, no, bank} {
		if strings.ContainsAny(v, " \t\n") {
			return gerror.New("开户人/账号/开户行不能包含空格")
		}
	}
	return nil
}

func (s *sWithdrawal) CardAdd(ctx context.Context, in service.CardInput) (int64, error) {
	if err := checkCardFields(in.AccountName, in.AccountNo, in.BankName, true); err != nil {
		return 0, err
	}
	if in.AccountType <= 0 {
		in.AccountType = 1
	}
	dup, err := g.Model("bank_card").Ctx(ctx).
		Where("site_id", wdSiteId).Where("user_id", in.UserId).
		Where("account_no", in.AccountNo).Count()
	if err != nil {
		return 0, err
	}
	if dup > 0 {
		return 0, gerror.New("该账号已添加过")
	}
	return g.DB().Model("bank_card").Ctx(ctx).Data(g.Map{
		"site_id": wdSiteId, "user_id": in.UserId, "account_type": in.AccountType,
		"account_name": in.AccountName, "account_no": in.AccountNo,
		"bank_name": in.BankName, "is_default": in.IsDefault,
	}).InsertAndGetId()
}

func (s *sWithdrawal) CardUpdate(ctx context.Context, in service.CardInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if err := checkCardFields(in.AccountName, in.AccountNo, in.BankName, false); err != nil {
		return err
	}
	data := g.Map{"updated_at": gtime.Now(), "is_default": in.IsDefault}
	if in.AccountType > 0 {
		data["account_type"] = in.AccountType
	}
	if in.AccountName != "" {
		data["account_name"] = in.AccountName
	}
	if in.AccountNo != "" {
		data["account_no"] = in.AccountNo
	}
	if in.BankName != "" {
		data["bank_name"] = in.BankName
	}
	res, err := g.Model("bank_card").Ctx(ctx).
		Where("site_id", wdSiteId).Where("id", in.Id).Where("user_id", in.UserId).
		Data(data).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("收款账户不存在")
	}
	return nil
}

func (s *sWithdrawal) CardDel(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return gerror.New("请选择要删除的账户")
	}
	res, err := g.Model("bank_card").Ctx(ctx).
		Where("site_id", wdSiteId).Where("user_id", userId).WhereIn("id", ids).Delete()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("收款账户不存在")
	}
	return nil
}
