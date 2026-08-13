// Package front 前台优惠券控制器。
package front

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/coupon/v1"
	"github.com/JarvanDante/my_service/internal/modules/coupon/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.ICoupon }

func New(svc service.ICoupon) *Controller { return &Controller{svc: svc} }

// optionalUid 公开接口用: 未登录返回 0, 不报错。
func optionalUid(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
}

func uid(ctx context.Context) (int64, error) {
	id := optionalUid(ctx)
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Tpls(ctx context.Context, req *v1.TplsReq) (res *v1.TplsRes, err error) {
	list, err := c.svc.Tpls(ctx, optionalUid(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.TplsRes{List: make([]v1.TplItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.TplItem{
			Id: d.Id, Name: d.Name, Type: d.Type, Scene: d.Scene, FaceValue: d.FaceValue,
			Discount: d.Discount, Threshold: d.Threshold, MaxDeduct: d.MaxDeduct,
			ExpireDay: d.ExpireDay, Received: d.Received,
		})
	}
	return res, nil
}

func (c *Controller) Receive(ctx context.Context, req *v1.ReceiveReq) (res *v1.ReceiveRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.svc.Receive(ctx, userId, req.TplId)
	if err != nil {
		return nil, err
	}
	return &v1.ReceiveRes{Id: id}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	status := 0
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.My(ctx, userId, status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.MyItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.MyItem{
			Id: d.Id, Name: d.Name, Type: d.Type, Scene: d.Scene, FaceValue: d.FaceValue,
			Discount: d.Discount, Threshold: d.Threshold, MaxDeduct: d.MaxDeduct,
			Status: d.Status, StatusText: service.StatusText(d.Status),
			ExpireAt: d.ExpireAt, UsedAt: d.UsedAt, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Available(ctx context.Context, req *v1.BestReq) (res *v1.BestRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, bestId, bestDeduct, err := c.svc.Available(ctx, userId, req.Scene, req.Amount)
	if err != nil {
		return nil, err
	}
	res = &v1.BestRes{
		BestId: bestId, BestDeduct: bestDeduct, PayAmount: req.Amount - bestDeduct,
		List: make([]v1.AvailableItem, 0, len(list)),
	}
	for _, d := range list {
		res.List = append(res.List, v1.AvailableItem{
			Id: d.Id, Name: d.Name, Deduct: d.Deduct, Type: d.Type,
			Expire: d.ExpireAt, IsBest: d.Id == bestId,
		})
	}
	return res, nil
}
