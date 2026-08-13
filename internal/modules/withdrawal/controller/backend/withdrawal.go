// Package backend 后台提现控制器。
package backend

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/withdrawal/v1"
	"github.com/JarvanDante/my_service/internal/modules/withdrawal/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IWithdrawal }

func New(svc service.IWithdrawal) *Controller { return &Controller{svc: svc} }

func adminId(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	status := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	var userId int64
	if req.UserId != "" {
		if v, e := strconv.ParseInt(req.UserId, 10, 64); e == nil {
			userId = v
		}
	}
	list, total, sum, pending, err := c.svc.List(ctx, service.ListFilter{
		UserId: userId, Status: status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{
		Total: total, SumAmount: sum, PendingNum: pending,
		List: make([]v1.Item, 0, len(list)),
	}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, TradeNo: d.TradeNo, UserId: d.UserId, Nickname: d.Nickname,
			Amount: d.Amount, Fee: d.Fee, RealAmount: d.RealAmount, FeeRate: d.FeeRate,
			BalanceAfter: d.BalanceAfter, AccountName: d.AccountName, AccountNo: d.AccountNo,
			BankName: d.BankName, Status: d.Status, StatusText: service.StatusText(d.Status),
			AuditBy: d.AuditBy, AuditAt: d.AuditAt, PaidAt: d.PaidAt,
			Remark: d.Remark, PayVoucher: d.PayVoucher, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, adminId(ctx), req.Id, req.Pass, req.Remark); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}

func (c *Controller) MarkPaid(ctx context.Context, req *v1.MarkPaidReq) (res *v1.MarkPaidRes, err error) {
	if err = c.svc.MarkPaid(ctx, adminId(ctx), req.Id, req.Voucher, req.Remark); err != nil {
		return nil, err
	}
	return &v1.MarkPaidRes{}, nil
}

func (c *Controller) Refund(ctx context.Context, req *v1.RejectPaidReq) (res *v1.RejectPaidRes, err error) {
	if err = c.svc.RefundPaid(ctx, adminId(ctx), req.Id, req.Remark); err != nil {
		return nil, err
	}
	return &v1.RejectPaidRes{}, nil
}
