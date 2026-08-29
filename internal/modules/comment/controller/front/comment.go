// Package front 前台评论控制器(发表需登录, 列表公开)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/comment/v1"
	"github.com/JarvanDante/my_service/internal/modules/comment/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IComment }

func New(svc service.IComment) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Add(ctx context.Context, req *v1.AddReq) (res *v1.AddRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	id, status, err := c.svc.Add(ctx, service.AddInput{
		UserId: userId, MediaType: req.MediaType, ContentId: req.ContentId,
		ParentId: req.ParentId, Content: req.Content, Pics: req.Pics,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AddRes{Id: id, Status: status}, nil
}

func itemOf(d service.ItemDTO) v1.Item {
	it := v1.Item{
		Id: d.Id, UserId: d.UserId, Nickname: d.Nickname, Img: d.Img, IsVip: d.IsVip,
		ParentId: d.ParentId, RootId: d.RootId, ReplyUserId: d.ReplyUserId, ReplyNickname: d.ReplyNickname,
		Content: d.Content, Pics: d.Pics, LikeCount: d.LikeCount, ReplyCount: d.ReplyCount,
		Liked: d.Liked, CreatedAt: d.CreatedAt,
	}
	if it.Pics == nil {
		it.Pics = []string{}
	}
	for _, r := range d.Replies {
		it.Replies = append(it.Replies, itemOf(r))
	}
	return it
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	viewer := int64(0)
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		viewer = r.GetCtxVar(consts.CtxUserId).Int64()
	}
	list, total, err := c.svc.List(ctx, req.MediaType, req.ContentId, req.Page, req.Size, req.Sort, viewer)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, itemOf(d))
	}
	return res, nil
}

func (c *Controller) Like(ctx context.Context, req *v1.LikeReq) (res *v1.LikeRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	count, liked, err := c.svc.Like(ctx, userId, req.Id, req.Flag)
	if err != nil {
		return nil, err
	}
	return &v1.LikeRes{Liked: liked, LikeCount: count}, nil
}
