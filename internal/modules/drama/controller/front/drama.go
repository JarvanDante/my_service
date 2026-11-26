// Package front 前台短剧控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/drama/v1"
	"github.com/JarvanDante/my_service/internal/modules/drama/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

type Controller struct {
	svc service.IDrama
	cat service.ICategory
	mod service.IModule
}

func New(svc service.IDrama, cat service.ICategory, mod service.IModule) *Controller {
	return &Controller{svc: svc, cat: cat, mod: mod}
}

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

func toItem(d *service.DramaDTO) v1.Item {
	return v1.Item{
		Id: d.Id, Title: d.Title, Cover: d.Cover, Intro: d.Intro, Category: d.Category,
		Categories: d.Categories, Tags: d.Tags, IsVip: d.IsVip, Price: d.Price, FreeEpisode: d.FreeEpisode,
		EpisodeCount: d.EpisodeCount, ViewCount: d.ViewCount, LikeCount: d.LikeCount,
		UpdateStatus: d.UpdateStatus, CreatedAt: d.CreatedAt,
	}
}

func (c *Controller) CategoryList(ctx context.Context, _ *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	list, err := c.cat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.CategoryListRes{List: make([]v1.CategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.CategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.FrontList(ctx, service.ListFilter{
		Category: req.Category, Categories: kit.NamesCSV(req.Categories),
		Tag: req.Tag, Tags: kit.NamesCSV(req.Tags), Keyword: req.Keyword,
		Sort: req.Sort, PayType: req.PayType, OnlyRecommend: req.Recommend == 1, Page: req.Page, Size: req.Size,
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
	d, err := c.svc.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{Item: toItem(d)}, nil
}

func toEpisode(d *service.EpisodeDTO) v1.EpisodeItem {
	return v1.EpisodeItem{
		Id: d.Id, Seq: d.Seq, Title: d.Title, Duration: d.Duration, Cover: d.Cover, PayMode: d.PayMode,
		IsFree: d.IsFree, IsVip: d.IsVip, Price: d.Price, IsBuy: d.IsBuy,
		Playable: d.Playable, NeedPay: d.NeedPay, NeedVip: d.NeedVip,
	}
}

func (c *Controller) Episodes(ctx context.Context, req *v1.EpisodesReq) (res *v1.EpisodesRes, err error) {
	title, list, err := c.svc.Episodes(ctx, optionalUid(ctx), req.Id, req.Desc)
	if err != nil {
		return nil, err
	}
	res = &v1.EpisodesRes{DramaId: req.Id, Title: title, List: make([]v1.EpisodeItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toEpisode(d))
	}
	return res, nil
}

func (c *Controller) Play(ctx context.Context, req *v1.PlayReq) (res *v1.PlayRes, err error) {
	d, err := c.svc.Play(ctx, optionalUid(ctx), req.EpisodeId)
	if err != nil {
		return nil, err
	}
	e := d.Episode
	return &v1.PlayRes{
		EpisodeId: e.Id, DramaId: e.DramaId, Seq: e.Seq, Title: e.Title, Cover: e.Cover, Duration: e.Duration,
		PlayUrl: d.PlayUrl, PreviewSec: d.PreviewSec, Playable: e.Playable, IsBuy: e.IsBuy,
		NeedPay: e.NeedPay, NeedVip: e.NeedVip, Price: e.Price, Enough: d.Enough, Reason: d.Reason,
		PrevId: d.PrevId, NextId: d.NextId,
	}, nil
}

func (c *Controller) Buy(ctx context.Context, req *v1.BuyReq) (res *v1.BuyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	price, bal, err := c.svc.Buy(ctx, userId, req.EpisodeId)
	if err != nil {
		return nil, err
	}
	return &v1.BuyRes{Price: price, Balance: bal}, nil
}

func toFrontModule(d *service.ModuleFrontDTO) v1.FrontModuleItem {
	row := v1.FrontModuleItem{
		Id: d.Id, Name: d.Name, Style: d.Style, Icon: d.Icon, Size: d.Size, Tags: d.Tags, Categories: d.Categories,
		Items: make([]v1.Item, 0, len(d.Items)),
	}
	for _, item := range d.Items {
		row.Items = append(row.Items, toItem(item))
	}
	return row
}

func (c *Controller) ModuleList(ctx context.Context, req *v1.ModuleListReq) (res *v1.ModuleListRes, err error) {
	list, err := c.mod.FrontRepo(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	res = &v1.ModuleListRes{List: make([]v1.FrontModuleItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toFrontModule(d))
	}
	return res, nil
}

func (c *Controller) ModuleRefresh(ctx context.Context, req *v1.ModuleRefreshReq) (res *v1.ModuleRefreshRes, err error) {
	d, err := c.mod.FrontRefresh(ctx, req.Id, kit.ParseI64CSV(req.Exclude))
	if err != nil {
		return nil, err
	}
	return &v1.ModuleRefreshRes{FrontModuleItem: toFrontModule(d)}, nil
}
