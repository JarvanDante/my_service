// Package logic 财务模块业务实现(B2)。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/modules/finance/domain"
	"github.com/JarvanDante/my_service/internal/modules/finance/service"
)

type sFinance struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IFinance { return &sFinance{repo: repo} }

// ---------- 充值套餐 ----------

func (s *sFinance) RechargePackages(ctx context.Context) ([]*service.RechargePackageDTO, error) {
	list, err := s.repo.ListRechargePackages(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.RechargePackageDTO, 0, len(list))
	for _, p := range list {
		out = append(out, &service.RechargePackageDTO{
			Id: p.Id, Name: p.Name, Amount: p.Amount, Coin: p.Coin, Bonus: p.Bonus,
			Sort: p.Sort, Status: p.Status,
		})
	}
	return out, nil
}

func checkRechargePkg(in service.RechargePackageInput) error {
	if in.Name == "" {
		return gerror.New("套餐名必填")
	}
	if in.Amount <= 0 || in.Coin <= 0 {
		return gerror.New("价格与到账金币必须大于0")
	}
	if in.Bonus < 0 || in.Status < 0 || in.Status > 1 {
		return gerror.New("参数不合法")
	}
	return nil
}

func (s *sFinance) CreateRechargePackage(ctx context.Context, in service.RechargePackageInput) (int64, error) {
	if err := checkRechargePkg(in); err != nil {
		return 0, err
	}
	return s.repo.CreateRechargePackage(ctx, in.Name, in.Amount, in.Coin, in.Bonus, in.Sort, in.Status)
}

func (s *sFinance) UpdateRechargePackage(ctx context.Context, in service.RechargePackageInput) error {
	if in.Id <= 0 {
		return gerror.New("套餐ID无效")
	}
	if err := checkRechargePkg(in); err != nil {
		return err
	}
	return s.repo.UpdateRechargePackage(ctx, in.Id, in.Name, in.Amount, in.Coin, in.Bonus, in.Sort, in.Status)
}

func (s *sFinance) DeleteRechargePackage(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("套餐ID无效")
	}
	return s.repo.DeleteRechargePackage(ctx, id)
}

// ---------- VIP 套餐 ----------

func (s *sFinance) VipPackages(ctx context.Context) ([]*service.VipPackageDTO, error) {
	list, err := s.repo.ListVipPackages(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.VipPackageDTO, 0, len(list))
	for _, p := range list {
		out = append(out, &service.VipPackageDTO{
			Id: p.Id, Name: p.Name, Days: p.Days, Price: p.Price,
			GroupId: p.GroupId, Sort: p.Sort, Status: p.Status,
		})
	}
	return out, nil
}

func checkVipPkg(in service.VipPackageInput) error {
	if in.Name == "" {
		return gerror.New("套餐名必填")
	}
	if in.Days <= 0 || in.Price <= 0 {
		return gerror.New("时长与价格必须大于0")
	}
	if in.Status < 0 || in.Status > 1 {
		return gerror.New("参数不合法")
	}
	return nil
}

func (s *sFinance) CreateVipPackage(ctx context.Context, in service.VipPackageInput) (int64, error) {
	if err := checkVipPkg(in); err != nil {
		return 0, err
	}
	return s.repo.CreateVipPackage(ctx, in.Name, in.Days, in.Price, in.GroupId, in.Sort, in.Status)
}

func (s *sFinance) UpdateVipPackage(ctx context.Context, in service.VipPackageInput) error {
	if in.Id <= 0 {
		return gerror.New("套餐ID无效")
	}
	if err := checkVipPkg(in); err != nil {
		return err
	}
	return s.repo.UpdateVipPackage(ctx, in.Id, in.Name, in.Days, in.Price, in.GroupId, in.Sort, in.Status)
}

func (s *sFinance) DeleteVipPackage(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("套餐ID无效")
	}
	return s.repo.DeleteVipPackage(ctx, id)
}

// ---------- 订单 / 流水 ----------

func (s *sFinance) Orders(ctx context.Context, in service.OrderListInput) (*service.PageDTO[service.OrderDTO], error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.ListOrders(ctx, domain.OrderFilter{
		OrderNo: in.OrderNo, UserId: in.UserId, Status: in.Status,
		StartDate: in.StartDate, EndDate: in.EndDate,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.OrderDTO, 0, len(list))
	for _, o := range list {
		out = append(out, &service.OrderDTO{
			Id: o.Id, OrderNo: o.OrderNo, UserId: o.UserId, PackageId: o.PackageId,
			Amount: o.Amount, Coin: o.Coin, Status: o.Status,
			PayAt: fmtTime(o.PayAt), CreatedAt: fmtTime(o.CreatedAt),
		})
	}
	return &service.PageDTO[service.OrderDTO]{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sFinance) BalanceLogs(ctx context.Context, in service.BalanceLogListInput) (*service.PageDTO[service.BalanceLogDTO], error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.ListBalanceLogs(ctx, domain.BalanceLogFilter{
		UserId: in.UserId, Scene: in.Scene, Direction: in.Direction,
		StartDate: in.StartDate, EndDate: in.EndDate,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.BalanceLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.BalanceLogDTO{
			Id: l.Id, UserId: l.UserId, Direction: l.Direction, Scene: l.Scene,
			Amount: l.Amount, BalanceBefore: l.BalanceBefore, BalanceAfter: l.BalanceAfter,
			RefId: l.RefId, Remark: l.Remark, CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return &service.PageDTO[service.BalanceLogDTO]{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

// ---------- 支付回调 ----------

// PayCallback 校验签名后置订单已支付并到账。
// 签名: md5(order_no + secret), secret 取配置 pay.callbackSecret; 未配置(开发环境)跳过校验。
func (s *sFinance) PayCallback(ctx context.Context, orderNo, sign string) error {
	if orderNo == "" {
		return gerror.New("订单号必填")
	}
	secret := g.Cfg().MustGet(ctx, "pay.callbackSecret").String()
	if secret != "" {
		if sign == "" || gmd5.MustEncryptString(orderNo+secret) != sign {
			return gerror.New("签名校验失败")
		}
	}
	return s.repo.MarkOrderPaid(ctx, orderNo)
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}
