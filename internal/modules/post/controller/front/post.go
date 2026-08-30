// Package front 前台帖子控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/post/v1"
	"github.com/JarvanDante/my_service/internal/modules/post/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct {
	svc service.IPost
	cat service.ICategory
}

func New(svc service.IPost, cat service.ICategory) *Controller {
	return &Controller{svc: svc, cat: cat}
}

func (c *Controller) CategoryList(ctx context.Context, _ *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	if c.cat == nil {
		return &v1.CategoryListRes{List: []v1.FrontCategoryItem{}}, nil
	}
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

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func optionalUid(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
}

func itemOf(d *service.PostDTO, withAudit bool) v1.Item {
	it := v1.Item{
		Id: d.Id, UserId: d.UserId, Nickname: d.Nickname, Img: d.Img, Sex: d.Sex, IsVip: d.IsVip,
		Title: d.Title, Content: d.Content, Pics: d.Pics, Topics: d.Topics,
		VideoUrl: d.VideoUrl, MediaId: d.MediaId, ViewCount: d.ViewCount,
		LikeCount: d.LikeCount, CommentCount: d.CommentCount, CreatedAt: d.CreatedAt,
	}
	if withAudit {
		it.Status, it.RejectReason = d.Status, d.RejectReason
	}
	return it
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.svc.Create(ctx, service.CreateInput{
		UserId: userId, Title: req.Title, Content: req.Content,
		Pics: req.Pics, Topics: req.Topics, VideoUrl: req.VideoUrl, MediaId: req.MediaId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.FrontList(ctx, service.ListFilter{
		Sort: req.Sort, Keyword: req.Keyword, UserId: req.UserId,
		Category: req.Category, FollowOnly: req.Follow == 1, ViewerId: optionalUid(ctx),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, itemOf(d, false))
	}
	return res, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	d, err := c.svc.Detail(ctx, req.Id, optionalUid(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{Post: itemOf(d, false)}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.My(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, itemOf(d, true))
	}
	return res, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.DeleteOwn(ctx, userId, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
