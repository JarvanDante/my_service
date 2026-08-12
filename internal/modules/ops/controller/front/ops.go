// Package front 前台运营配置控制器(公开)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/ops/v1"
	"github.com/JarvanDante/my_service/internal/modules/ops/service"
)

type Controller struct{ svc service.IOps }

func New(svc service.IOps) *Controller { return &Controller{svc: svc} }

func (c *Controller) Announcement(ctx context.Context, req *v1.AnnouncementReq) (res *v1.AnnouncementRes, err error) {
	list, err := c.svc.LiveAnnouncements(ctx, req.SysType)
	if err != nil {
		return nil, err
	}
	res = &v1.AnnouncementRes{List: make([]v1.AnnItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.AnnItem{
			Id: r.Id, Title: r.Title, Content: r.Content, TextNode: r.TextNode,
			Cover: r.Cover, JumpUrl: r.JumpUrl, SysType: r.SysType,
			StartAt: r.StartAt, EndAt: r.EndAt,
		})
	}
	return res, nil
}

func (c *Controller) Jumptab(ctx context.Context, req *v1.JumptabReq) (res *v1.JumptabRes, err error) {
	list, err := c.svc.FrontJumptabs(ctx, req.Location)
	if err != nil {
		return nil, err
	}
	res = &v1.JumptabRes{List: make([]v1.JumptabItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.JumptabItem{
			Id: r.Id, CnName: r.CnName, EnName: r.EnName, Avatar: r.Avatar,
			Link: r.Link, PicJumpLink: r.PicJumpLink, Location: r.Location, Rank: r.Rank,
		})
	}
	return res, nil
}
