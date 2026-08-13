// Package backend 后台UGC投稿控制器(列表/审核)。
package backend

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/publish/v1"
	"github.com/JarvanDante/my_service/internal/modules/publish/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IPublish }

func New(svc service.IPublish) *Controller { return &Controller{svc: svc} }

func adminId(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
}

// atoiOr 后台筛选参数一律 string 接收: 空串/非法值都当"不筛选", 返回 def 哨兵。
func atoiOr(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status:  atoiOr(req.Status, -1), // -1=全部(0 是合法的"待审")
		UserId:  int64(atoiOr(req.UserId, 0)),
		Type:    atoiOr(req.Type, 0),
		Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, UserId: d.UserId, Type: d.Type, Title: d.Title, Intro: d.Intro,
			Cover: d.Cover, Resource: d.Resource, Tags: d.Tags, Status: d.Status,
			RejectReason: d.RejectReason, AuditBy: d.AuditBy, AuditAt: d.AuditAt,
			CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, req.Id, adminId(ctx), req.Pass, req.RejectReason); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}
