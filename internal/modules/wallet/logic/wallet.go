// Package logic 钱包业务(余额/流水/人工调账)。
// 余额本身沉在 users.balance, 变动统一走 shared/balance, 这里只做读聚合与调账入口。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/wallet/service"
	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const walletSiteId = 1

type sWallet struct{}

func New() service.IWallet { return &sWallet{} }

func (s *sWallet) Summary(ctx context.Context, userId int64) (*service.SummaryDTO, error) {
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return nil, err
	}
	out := &service.SummaryDTO{}
	if bal != nil {
		out.Balance = bal.Float64()
	}
	// 在途冻结: 申请中(1) + 审核通过待打款(2)
	frozen, err := g.Model("withdrawal").Ctx(ctx).
		Where("site_id", walletSiteId).Where("user_id", userId).WhereIn("status", g.Slice{1, 2}).
		Sum("amount")
	if err != nil {
		return nil, err
	}
	out.Frozen = frozen
	// 累计提现成功
	paid, err := g.Model("withdrawal").Ctx(ctx).
		Where("site_id", walletSiteId).Where("user_id", userId).Where("status", 4).Sum("amount")
	if err != nil {
		return nil, err
	}
	out.Withdrawn = paid
	in, err := g.Model("user_balance_log").Ctx(ctx).
		Where("site_id", walletSiteId).Where("user_id", userId).Where("direction", 1).Sum("amount")
	if err != nil {
		return nil, err
	}
	out.TotalIn = in
	outSum, err := g.Model("user_balance_log").Ctx(ctx).
		Where("site_id", walletSiteId).Where("user_id", userId).Where("direction", 2).Sum("amount")
	if err != nil {
		return nil, err
	}
	out.TotalOut = outSum
	return out, nil
}

func (s *sWallet) Waters(ctx context.Context, f service.WaterFilter) ([]*service.WaterDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("user_balance_log").Ctx(ctx).Where("site_id", walletSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Direction == 1 || f.Direction == 2 {
		m = m.Where("direction", f.Direction)
	}
	if f.Scene != "" {
		m = m.Where("scene", f.Scene)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserBalanceLog
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.WaterDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.WaterDTO{
			Id: r.Id, UserId: r.UserId, Direction: r.Direction, Scene: r.Scene,
			Amount: r.Amount, BalanceBefore: r.BalanceBefore, BalanceAfter: r.BalanceAfter,
			RefId: r.RefId, Remark: r.Remark, CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sWallet) Adjust(ctx context.Context, adminId, userId int64, amount float64, remark string) (float64, error) {
	if amount == 0 {
		return 0, gerror.New("调账金额不能为0")
	}
	exist, err := g.Model("users").Ctx(ctx).Where("id", userId).Count()
	if err != nil {
		return 0, err
	}
	if exist == 0 {
		return 0, gerror.New("用户不存在")
	}
	refId := "admin:" + gconv.String(adminId)
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if amount > 0 {
			return balance.Add(ctx, tx, userId, amount, balance.SceneAdminAdjust, refId, remark)
		}
		return balance.Deduct(ctx, tx, userId, -amount, balance.SceneAdminAdjust, refId, remark)
	})
	if err != nil {
		return 0, err
	}
	v, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return 0, err
	}
	return v.Float64(), nil
}
