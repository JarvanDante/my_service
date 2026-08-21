// Package backend 后台评论审核。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/comment/v1"
	"github.com/JarvanDante/my_service/internal/modules/comment/service"
)

type Controller struct{ svc service.IComment }

func New(svc service.IComment) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.AdminList(ctx, service.AdminListFilter{
		Status: statusFilter, Kind: req.Kind, Keyword: req.Keyword,
		UserId: req.UserId, MediaType: req.MediaType,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, UserId: d.UserId, Nickname: d.Nickname, Img: d.Img, IsVip: d.IsVip,
			MediaType: d.MediaType, ContentId: d.ContentId, ParentId: d.ParentId, RootId: d.RootId,
			Content: d.Content, LikeCount: d.LikeCount, ReplyCount: d.ReplyCount,
			Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, req.Id, req.Pass); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}
