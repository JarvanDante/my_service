// Package backend 后台财务控制器(B2)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/finance/v1"
	"github.com/JarvanDante/my_service/internal/modules/finance/service"
)

type Controller struct{ fin service.IFinance }

func New(svc service.IFinance) *Controller { return &Controller{fin: svc} }

// ---------- 充值套餐 ----------

func (c *Controller) RechargePkgList(ctx context.Context, req *v1.RechargePkgListReq) (res *v1.RechargePkgListRes, err error) {
	list, err := c.fin.RechargePackages(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.RechargePackageItem, 0, len(list))
	for _, p := range list {
		items = append(items, v1.RechargePackageItem{
			Id: p.Id, Name: p.Name, Amount: p.Amount, Coin: p.Coin, Bonus: p.Bonus,
			Sort: p.Sort, Status: p.Status,
		})
	}
	return &v1.RechargePkgListRes{List: items}, nil
}

func (c *Controller) RechargePkgCreate(ctx context.Context, req *v1.RechargePkgCreateReq) (res *v1.RechargePkgCreateRes, err error) {
	id, err := c.fin.CreateRechargePackage(ctx, service.RechargePackageInput{
		Name: req.Name, Amount: req.Amount, Coin: req.Coin, Bonus: req.Bonus,
		Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RechargePkgCreateRes{Id: id}, nil
}

func (c *Controller) RechargePkgUpdate(ctx context.Context, req *v1.RechargePkgUpdateReq) (res *v1.RechargePkgUpdateRes, err error) {
	if err = c.fin.UpdateRechargePackage(ctx, service.RechargePackageInput{
		Id: req.Id, Name: req.Name, Amount: req.Amount, Coin: req.Coin, Bonus: req.Bonus,
		Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.RechargePkgUpdateRes{}, nil
}

func (c *Controller) RechargePkgDelete(ctx context.Context, req *v1.RechargePkgDeleteReq) (res *v1.RechargePkgDeleteRes, err error) {
	if err = c.fin.DeleteRechargePackage(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.RechargePkgDeleteRes{}, nil
}

// ---------- VIP 套餐 ----------

func (c *Controller) VipPkgList(ctx context.Context, req *v1.VipPkgListReq) (res *v1.VipPkgListRes, err error) {
	list, err := c.fin.VipPackages(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.VipPackageItem, 0, len(list))
	for _, p := range list {
		items = append(items, v1.VipPackageItem{
			Id: p.Id, Name: p.Name, Days: p.Days, Price: p.Price,
			GroupId: p.GroupId, Sort: p.Sort, Status: p.Status,
		})
	}
	return &v1.VipPkgListRes{List: items}, nil
}

func (c *Controller) VipPkgCreate(ctx context.Context, req *v1.VipPkgCreateReq) (res *v1.VipPkgCreateRes, err error) {
	id, err := c.fin.CreateVipPackage(ctx, service.VipPackageInput{
		Name: req.Name, Days: req.Days, Price: req.Price, GroupId: req.GroupId,
		Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.VipPkgCreateRes{Id: id}, nil
}

func (c *Controller) VipPkgUpdate(ctx context.Context, req *v1.VipPkgUpdateReq) (res *v1.VipPkgUpdateRes, err error) {
	if err = c.fin.UpdateVipPackage(ctx, service.VipPackageInput{
		Id: req.Id, Name: req.Name, Days: req.Days, Price: req.Price, GroupId: req.GroupId,
		Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.VipPkgUpdateRes{}, nil
}

func (c *Controller) VipPkgDelete(ctx context.Context, req *v1.VipPkgDeleteReq) (res *v1.VipPkgDeleteRes, err error) {
	if err = c.fin.DeleteVipPackage(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.VipPkgDeleteRes{}, nil
}

// ---------- 订单 / 流水 ----------

func (c *Controller) Orders(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	dto, err := c.fin.Orders(ctx, service.OrderListInput{
		OrderNo: req.OrderNo, UserId: req.UserId, Status: req.Status,
		StartDate: req.StartDate, EndDate: req.EndDate, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.OrderItem, 0, len(dto.List))
	for _, o := range dto.List {
		items = append(items, v1.OrderItem{
			Id: o.Id, OrderNo: o.OrderNo, UserId: o.UserId, PackageId: o.PackageId,
			Amount: o.Amount, Coin: o.Coin, Status: o.Status, PayAt: o.PayAt, CreatedAt: o.CreatedAt,
		})
	}
	return &v1.OrderListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) BalanceLogs(ctx context.Context, req *v1.BalanceLogListReq) (res *v1.BalanceLogListRes, err error) {
	dto, err := c.fin.BalanceLogs(ctx, service.BalanceLogListInput{
		UserId: req.UserId, Scene: req.Scene, Direction: req.Direction,
		StartDate: req.StartDate, EndDate: req.EndDate, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.BalanceLogItem, 0, len(dto.List))
	for _, l := range dto.List {
		items = append(items, v1.BalanceLogItem{
			Id: l.Id, UserId: l.UserId, Direction: l.Direction, Scene: l.Scene,
			Amount: l.Amount, BalanceBefore: l.BalanceBefore, BalanceAfter: l.BalanceAfter,
			RefId: l.RefId, Remark: l.Remark, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.BalanceLogListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

// ---------- 支付回调(公开) ----------

func (c *Controller) PayCallback(ctx context.Context, req *v1.PayCallbackReq) (res *v1.PayCallbackRes, err error) {
	if err = c.fin.PayCallback(ctx, req.OrderNo, req.Sign); err != nil {
		return nil, err
	}
	return &v1.PayCallbackRes{}, nil
}
