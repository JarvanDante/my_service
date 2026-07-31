// Package backend 后台系统模块控制器(B7)。
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/system/v1"
	"github.com/JarvanDante/my_service/internal/modules/system/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ sys service.ISystem }

func New(svc service.ISystem) *Controller { return &Controller{sys: svc} }

func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

// Push 发布公告/推送。
func (c *Controller) Push(ctx context.Context, req *v1.PushReq) (res *v1.PushRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.sys.Push(ctx, service.PushInput{
		Title: req.Title, Content: req.Content, Type: req.Type, OperatorId: op,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PushRes{Id: id}, nil
}

// NoticeList 公告列表。
func (c *Controller) NoticeList(ctx context.Context, req *v1.NoticeListReq) (res *v1.NoticeListRes, err error) {
	dto, err := c.sys.Notices(ctx, service.NoticeListInput{
		Type: req.Type, Status: req.Status, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.NoticeItem, 0, len(dto.List))
	for _, n := range dto.List {
		items = append(items, v1.NoticeItem{
			Id: n.Id, Title: n.Title, Content: n.Content, Type: n.Type,
			Status: n.Status, CreatedBy: n.CreatedBy, CreatedAt: n.CreatedAt,
		})
	}
	return &v1.NoticeListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

// NoticeStatus 上/下线公告。
func (c *Controller) NoticeStatus(ctx context.Context, req *v1.NoticeStatusReq) (res *v1.NoticeStatusRes, err error) {
	if err = c.sys.SetNoticeStatus(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &v1.NoticeStatusRes{}, nil
}

// CustomerUrlGet 查看客服链接。
func (c *Controller) CustomerUrlGet(ctx context.Context, req *v1.CustomerUrlGetReq) (res *v1.CustomerUrlGetRes, err error) {
	url, err := c.sys.GetCustomerUrl(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CustomerUrlGetRes{Url: url}, nil
}

// CustomerUrlPut 配置客服链接。
func (c *Controller) CustomerUrlPut(ctx context.Context, req *v1.CustomerUrlPutReq) (res *v1.CustomerUrlPutRes, err error) {
	if err = c.sys.SetCustomerUrl(ctx, req.Url); err != nil {
		return nil, err
	}
	return &v1.CustomerUrlPutRes{}, nil
}
