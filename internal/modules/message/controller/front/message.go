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
	count, err := c.svc.UnreadCount(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.UnreadRes{Count: count}, nil
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
