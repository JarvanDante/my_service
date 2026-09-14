// Package front 前台小说控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/novel/v1"
	"github.com/JarvanDante/my_service/internal/modules/novel/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

type Controller struct {
	svc service.INovel
	cat service.ICategory
	mod service.IModule
}

func New(svc service.INovel, cat service.ICategory, mod service.IModule) *Controller {
	return &Controller{svc: svc, cat: cat, mod: mod}
}

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

func toItem(d *service.NovelDTO) v1.Item {
	return v1.Item{
		Id: d.Id, Title: d.Title, Author: d.Author, Cover: d.Cover, Intro: d.Intro,
		Category: d.Category, Tags: d.Tags, IsVip: d.IsVip, Price: d.Price,
		FreeChapter: d.FreeChapter, ChapterCount: d.ChapterCount, WordCount: d.WordCount,
		IsAudio: d.IsAudio, ViewCount: d.ViewCount, LikeCount: d.LikeCount,
		UpdateStatus: d.UpdateStatus, IsBuy: d.IsBuy, CreatedAt: d.CreatedAt,
	}
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

func (c *Controller) CategoryList(ctx context.Context, _ *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	list, err := c.cat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.CategoryListRes{List: make([]v1.FrontCategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.FrontCategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
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

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	d, err := c.svc.Detail(ctx, optionalUid(ctx), req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{
		Item: toItem(d.Novel), Playable: d.Playable, NeedPay: d.NeedPay,
		NeedVip: d.NeedVip, Enough: d.Enough, Reason: d.Reason,
	}, nil
}

func (c *Controller) Chapters(ctx context.Context, req *v1.ChaptersReq) (res *v1.ChaptersRes, err error) {
	title, list, err := c.svc.Chapters(ctx, optionalUid(ctx), req.Id, req.Desc)
	if err != nil {
		return nil, err
	}
	res = &v1.ChaptersRes{NovelId: req.Id, Title: title, List: make([]v1.ChapterItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.ChapterItem{
			Id: d.Id, Seq: d.Seq, Title: d.Title, WordCount: d.WordCount,
			HasAudio: d.AudioUrl != "", IsFree: d.IsFree, Playable: d.Playable,
		})
	}
	return res, nil
}

func (c *Controller) Read(ctx context.Context, req *v1.ReadReq) (res *v1.ReadRes, err error) {
	d, err := c.svc.Read(ctx, optionalUid(ctx), req.ChapterId)
	if err != nil {
		return nil, err
	}
	return &v1.ReadRes{
		ChapterId: d.ChapterId, NovelId: d.NovelId, Seq: d.Seq, Title: d.Title,
		Content: d.Content, WordCount: d.WordCount, AudioUrl: d.AudioUrl,
		PrevId: d.PrevId, NextId: d.NextId,
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

func (c *Controller) MayLike(ctx context.Context, req *v1.MayLikeReq) (res *v1.MayLikeRes, err error) {
	list, err := c.svc.MayLike(ctx, req.Id, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.MayLikeRes{List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toItem(d))
	}
	return res, nil
}
