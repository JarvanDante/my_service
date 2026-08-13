// Package front 前台图集控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/photo/v1"
	"github.com/JarvanDante/my_service/internal/modules/photo/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IPhoto }

func New(svc service.IPhoto) *Controller { return &Controller{svc: svc} }

// optionalUid 公开接口: 带 token 就取到用户, 不带返回 0(用于标记 is_buy/解锁态)。
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

func toItem(d *service.PhotoDTO) v1.Item {
	return v1.Item{
		Id: d.Id, Title: d.Title, Cover: d.Cover, Intro: d.Intro, Category: d.Category,
		Tags: d.Tags, IsVip: d.IsVip, Price: d.Price, FreeCount: d.FreeCount,
		PicCount: d.PicCount, ViewCount: d.ViewCount, LikeCount: d.LikeCount,
		IsBuy: d.IsBuy, CreatedAt: d.CreatedAt,
	}
}

func toPicV1(in []service.PicDTO) []v1.Pic {
	out := make([]v1.Pic, 0, len(in))
	for _, p := range in {
		out = append(out, v1.Pic{Url: p.Url, Width: p.Width, Height: p.Height})
	}
	return out
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.FrontList(ctx, optionalUid(ctx), service.ListFilter{
		Category: req.Category, Tag: req.Tag, Keyword: req.Keyword,
		Sort: req.Sort, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toItem(d))
	}
	return res, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	d, err := c.svc.Detail(ctx, optionalUid(ctx), req.Id)
	if err != nil {
		return nil, err
	}
	// d.Photo.Pics 已由 logic 按解锁态截断, 这里原样透出即可。
	return &v1.DetailRes{
		Item: toItem(d.Photo), Pics: toPicV1(d.Photo.Pics),
		Playable: d.Playable, NeedPay: d.NeedPay, NeedVip: d.NeedVip,
		Enough: d.Enough, Reason: d.Reason,
		PreviewCount: d.PreviewCount, TotalCount: d.TotalCount,
	}, nil
}

func (c *Controller) Buy(ctx context.Context, req *v1.BuyReq) (res *v1.BuyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	price, bal, err := c.svc.Buy(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.BuyRes{Price: price, Balance: bal}, nil
}
