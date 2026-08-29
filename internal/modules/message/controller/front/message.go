// Package front 前台系统消息控制器(需登录)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/message/v1"
	"github.com/JarvanDante/my_service/internal/modules/message/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IMessage }

func New(svc service.IMessage) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.MyList(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.MsgItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.MsgItem{
			Id: r.Id, Type: r.Type, Title: r.Title, Content: r.Content,
			IsRead: r.IsRead, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Unread(ctx context.Context, req *v1.UnreadReq) (res *v1.UnreadRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	u, err := c.svc.UnreadAll(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.UnreadRes{
		Count: u.Sys + u.Comment + u.Like, Sys: u.Sys, Comment: u.Comment, Like: u.Like,
	}, nil
}

func (c *Controller) InteractList(ctx context.Context, req *v1.InteractListReq) (res *v1.InteractListRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.InteractList(ctx, userId, req.Channel, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.InteractListRes{Total: total, List: make([]v1.InteractItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.InteractItem{
			Id: r.Id, Channel: r.Channel, SubType: r.SubType, IsRead: r.IsRead,
			CreatedAt: r.CreatedAt, ActorId: r.ActorId, ActorName: r.ActorName,
			ActorAvatar: r.ActorAvatar, ActorSex: r.ActorSex, ActorCount: r.ActorCount,
			MediaType: r.MediaType, ContentId: r.ContentId, ObjectTitle: r.ObjectTitle,
			TargetType: r.TargetType, CommentId: r.CommentId, RootCommentId: r.RootCommentId,
			Snippet: r.Snippet,
		})
	}
	return res, nil
}

func (c *Controller) InteractRead(ctx context.Context, req *v1.InteractReadReq) (res *v1.InteractReadRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.MarkInteractRead(ctx, userId, req.Id, req.All, req.Channel); err != nil {
		return nil, err
	}
	return &v1.InteractReadRes{}, nil
}

func (c *Controller) Read(ctx context.Context, req *v1.ReadReq) (res *v1.ReadRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.MarkRead(ctx, userId, req.Id, req.All); err != nil {
		return nil, err
	}
	return &v1.ReadRes{}, nil
}
