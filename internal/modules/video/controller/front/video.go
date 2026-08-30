// Package front 前台视频控制器(公开浏览)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/video/v1"
	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

type Controller struct {
	svc        service.IVideo
	videoCat   service.ICategory
	cartoonCat service.ICategory
	douyinCat  service.ICategory
	videoMod   service.IModule
	cartoonMod service.IModule
}

func New(svc service.IVideo, videoCat, cartoonCat, douyinCat service.ICategory, videoMod, cartoonMod service.IModule) *Controller {
	return &Controller{svc: svc, videoCat: videoCat, cartoonCat: cartoonCat, douyinCat: douyinCat, videoMod: videoMod, cartoonMod: cartoonMod}
}

// toItem 只挑前台需要的字段, cover_key/source_key/created_by/status 不外发。
func toItem(v *service.VideoDTO) v1.Item {
	if v == nil {
		return v1.Item{}
	}
	return v1.Item{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, SourceUrl: v.SourceUrl,
		Category: v.Category, Categories: v.Categories,
		Duration: v.Duration, CreatedAt: v.CreatedAt,
		UpUserId: v.UpUserId, UpNickname: v.UpNickname, UpAvatar: v.UpAvatar, Followed: v.Followed,
		CommentCount: v.CommentCount,
		PreviewSec: v.PreviewSec, NeedVip: v.NeedVip,
	}
}

func viewerId(ctx context.Context) int64 {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxUserId).Int64()
}

func uid(ctx context.Context) (int64, error) {
	id := viewerId(ctx)
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) CategoryList(ctx context.Context, _ *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	if c.videoCat == nil {
		return &v1.CategoryListRes{List: []v1.FrontCategoryItem{}}, nil
	}
	list, err := c.videoCat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.CategoryListRes{List: make([]v1.FrontCategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.FrontCategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.svc.FrontList(ctx, service.FrontListInput{
		Keyword: req.Keyword, Category: req.Category, Categories: kit.NamesCSV(req.Categories),
		Tag: req.Tag, Tags: kit.NamesCSV(req.Tags), Kind: entity.VideoKindVideo,
		Sort: req.Sort, Page: req.Page, Size: req.Size, ViewerId: viewerId(ctx),
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.Item, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.ListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	v, err := c.svc.FrontDetail(ctx, req.Id, viewerId(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{Item: toItem(v)}, nil
}

func (c *Controller) CartoonCategoryList(ctx context.Context, _ *v1.CartoonCategoryListReq) (res *v1.CartoonCategoryListRes, err error) {
	if c.cartoonCat == nil {
		return &v1.CartoonCategoryListRes{List: []v1.FrontCategoryItem{}}, nil
	}
	list, err := c.cartoonCat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.CartoonCategoryListRes{List: make([]v1.FrontCategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.FrontCategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
}

func (c *Controller) CartoonList(ctx context.Context, req *v1.CartoonListReq) (res *v1.CartoonListRes, err error) {
	dto, err := c.svc.FrontList(ctx, service.FrontListInput{
		Keyword: req.Keyword, Category: req.Category, Categories: kit.NamesCSV(req.Categories),
		Tag: req.Tag, Tags: kit.NamesCSV(req.Tags), Kind: entity.VideoKindCartoon,
		Sort: req.Sort, Page: req.Page, Size: req.Size, ViewerId: viewerId(ctx),
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.Item, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.CartoonListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) DouyinCategoryList(ctx context.Context, _ *v1.DouyinCategoryListReq) (res *v1.DouyinCategoryListRes, err error) {
	if c.douyinCat == nil {
		return &v1.DouyinCategoryListRes{List: []v1.FrontCategoryItem{}}, nil
	}
	list, err := c.douyinCat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.DouyinCategoryListRes{List: make([]v1.FrontCategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.FrontCategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
}

func (c *Controller) DouyinList(ctx context.Context, req *v1.DouyinListReq) (res *v1.DouyinListRes, err error) {
	dto, err := c.svc.FrontList(ctx, service.FrontListInput{
		Keyword: req.Keyword, Category: req.Category, Tag: req.Tag, Kind: entity.VideoKindDouyin,
		Sort: req.Sort, Page: req.Page, Size: req.Size,
		ViewerId: viewerId(ctx), FollowOnly: req.Follow == 1, UpUserId: req.UpUserId,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.Item, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.DouyinListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) DouyinSubmit(ctx context.Context, req *v1.DouyinSubmitReq) (res *v1.DouyinSubmitRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	vid, err := c.svc.SubmitDouyin(ctx, service.SubmitDouyinInput{
		UserId: id, Title: req.Title, Description: req.Description,
		CoverUrl: req.CoverUrl, CoverKey: req.CoverKey,
		SourceUrl: req.SourceUrl, SourceKey: req.SourceKey,
		Duration: req.Duration, Tags: req.Tags,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DouyinSubmitRes{Id: vid}, nil
}

func (c *Controller) DouyinMy(ctx context.Context, req *v1.DouyinMyReq) (res *v1.DouyinMyRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.svc.MyDouyin(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.DouyinMineItem, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, v1.DouyinMineItem{
			Item: toItem(v), Status: v.Status, RejectReason: v.RejectReason,
		})
	}
	return &v1.DouyinMyRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func toFrontModules(list []*service.ModuleFrontDTO) []v1.FrontModuleItem {
	out := make([]v1.FrontModuleItem, 0, len(list))
	for _, d := range list {
		row := v1.FrontModuleItem{
			Id: d.Id, Name: d.Name, Style: d.Style, Icon: d.Icon, Size: d.Size, Tags: d.Tags, Categories: d.Categories,
			Items: make([]v1.Item, 0, len(d.Items)),
		}
		for _, item := range d.Items {
			row.Items = append(row.Items, toItem(item))
		}
		out = append(out, row)
	}
	return out
}

func (c *Controller) VideoModuleList(ctx context.Context, req *v1.VideoModuleListReq) (res *v1.VideoModuleListRes, err error) {
	list, err := c.videoMod.FrontRepo(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	return &v1.VideoModuleListRes{List: toFrontModules(list)}, nil
}

func (c *Controller) CartoonModuleList(ctx context.Context, req *v1.CartoonModuleListReq) (res *v1.CartoonModuleListRes, err error) {
	list, err := c.cartoonMod.FrontRepo(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	return &v1.CartoonModuleListRes{List: toFrontModules(list)}, nil
}
