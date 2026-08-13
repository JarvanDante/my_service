// Package front 前台提现控制器。
package front

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/withdrawal/v1"
	"github.com/JarvanDante/my_service/internal/modules/withdrawal/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IWithdrawal }

func New(svc service.IWithdrawal) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Config(ctx context.Context, req *v1.ConfigReq) (res *v1.ConfigRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Config(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigRes{
		Open: d.Open, MinAmount: d.MinAmount, MaxAmount: d.MaxAmount, Multiple: d.Multiple,
		FeeRate: d.FeeRate, DailyLimit: d.DailyLimit, DailyUsed: d.DailyUsed,
		Balance: d.Balance, Frozen: d.Frozen,
	}, nil
}

func (c *Controller) Apply(ctx context.Context, req *v1.ApplyReq) (res *v1.ApplyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	r, err := c.svc.Apply(ctx, userId, req.CardId, req.Amount)
	if err != nil {
		return nil, err
	}
	return &v1.ApplyRes{
		Id: r.Id, TradeNo: r.TradeNo, Amount: r.Amount, Fee: r.Fee, RealAmount: r.RealAmount,
	}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	status := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.My(ctx, service.ListFilter{
		UserId: userId, Status: status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, TradeNo: d.TradeNo, Amount: d.Amount, Fee: d.Fee,
			RealAmount: d.RealAmount, Status: d.Status, StatusText: service.StatusText(d.Status),
			AccountNo: d.AccountNo, AccountName: d.AccountName, BankName: d.BankName,
			Remark: d.Remark, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Cancel(ctx context.Context, req *v1.CancelReq) (res *v1.CancelRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.Cancel(ctx, userId, req.Id); err != nil {
		return nil, err
	}
	return &v1.CancelRes{}, nil
}

func (c *Controller) CardList(ctx context.Context, req *v1.CardListReq) (res *v1.CardListRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, err := c.svc.CardList(ctx, userId)
	if err != nil {
		return nil, err
	}
	res = &v1.CardListRes{List: make([]v1.CardItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.CardItem{
			Id: d.Id, AccountType: d.AccountType, AccountName: d.AccountName,
			AccountNo: d.AccountNo, BankName: d.BankName, IsDefault: d.IsDefault,
			CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) CardAdd(ctx context.Context, req *v1.CardAddReq) (res *v1.CardAddRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.svc.CardAdd(ctx, service.CardInput{
		UserId: userId, AccountType: req.AccountType, AccountName: req.AccountName,
		AccountNo: req.AccountNo, BankName: req.BankName, IsDefault: req.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CardAddRes{Id: id}, nil
}

func (c *Controller) CardUpdate(ctx context.Context, req *v1.CardUpdateReq) (res *v1.CardUpdateRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.CardUpdate(ctx, service.CardInput{
		Id: req.Id, UserId: userId, AccountType: req.AccountType, AccountName: req.AccountName,
		AccountNo: req.AccountNo, BankName: req.BankName, IsDefault: req.IsDefault,
	}); err != nil {
		return nil, err
	}
	return &v1.CardUpdateRes{}, nil
}

func (c *Controller) CardDel(ctx context.Context, req *v1.CardDelReq) (res *v1.CardDelRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.CardDel(ctx, userId, req.Ids); err != nil {
		return nil, err
	}
	return &v1.CardDelRes{}, nil
}
