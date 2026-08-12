// Package backend 后台运营配置控制器(公告/跳转位/敏感词)。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/ops/v1"
	"github.com/JarvanDante/my_service/internal/modules/ops/service"
)

type Controller struct{ svc service.IOps }

func New(svc service.IOps) *Controller { return &Controller{svc: svc} }

// statusOf 空字符串 = 不过滤(-1); "0"/"1" = 精确过滤。
func statusOf(s string) int {
	if s == "" {
		return -1
	}
	if v, e := strconv.Atoi(s); e == nil {
		return v
	}
	return -1
}

// ---------- 公告 ----------

func (c *Controller) AnnList(ctx context.Context, req *v1.AnnListReq) (res *v1.AnnListRes, err error) {
	list, total, err := c.svc.AnnList(ctx, service.PageFilter{
		Status: statusOf(req.Status), Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AnnListRes{Total: total, List: make([]v1.AnnItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.AnnItem{
			Id: r.Id, Title: r.Title, Content: r.Content, TextNode: r.TextNode,
			Cover: r.Cover, JumpUrl: r.JumpUrl, SysType: r.SysType,
			StartAt: r.StartAt, EndAt: r.EndAt, Status: r.Status, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) AnnCreate(ctx context.Context, req *v1.AnnCreateReq) (res *v1.AnnCreateRes, err error) {
	id, err := c.svc.AnnCreate(ctx, service.AnnSaveInput{
		Title: req.Title, Content: req.Content, TextNode: req.TextNode,
		Cover: req.Cover, JumpUrl: req.JumpUrl, SysType: req.SysType,
		StartAt: req.StartAt, EndAt: req.EndAt, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AnnCreateRes{Id: id}, nil
}

func (c *Controller) AnnUpdate(ctx context.Context, req *v1.AnnUpdateReq) (res *v1.AnnUpdateRes, err error) {
	if err = c.svc.AnnUpdate(ctx, service.AnnSaveInput{
		Id: req.Id, Title: req.Title, Content: req.Content, TextNode: req.TextNode,
		Cover: req.Cover, JumpUrl: req.JumpUrl, SysType: req.SysType,
		StartAt: req.StartAt, EndAt: req.EndAt, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AnnUpdateRes{}, nil
}

func (c *Controller) AnnDelete(ctx context.Context, req *v1.AnnDeleteReq) (res *v1.AnnDeleteRes, err error) {
	if err = c.svc.AnnDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.AnnDeleteRes{}, nil
}

// ---------- 跳转位 ----------

func (c *Controller) JtList(ctx context.Context, req *v1.JtListReq) (res *v1.JtListRes, err error) {
	list, total, err := c.svc.JtList(ctx, service.PageFilter{
		Status: statusOf(req.Status), Location: req.Location, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.JtListRes{Total: total, List: make([]v1.JtItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.JtItem{
			Id: r.Id, CnName: r.CnName, EnName: r.EnName, Avatar: r.Avatar,
			Link: r.Link, PicJumpLink: r.PicJumpLink, Location: r.Location,
			Rank: r.Rank, Status: r.Status, CreatedAt: r.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) JtCreate(ctx context.Context, req *v1.JtCreateReq) (res *v1.JtCreateRes, err error) {
	id, err := c.svc.JtCreate(ctx, service.JtSaveInput{
		CnName: req.CnName, EnName: req.EnName, Avatar: req.Avatar,
		Link: req.Link, PicJumpLink: req.PicJumpLink, Location: req.Location,
		Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.JtCreateRes{Id: id}, nil
}

func (c *Controller) JtUpdate(ctx context.Context, req *v1.JtUpdateReq) (res *v1.JtUpdateRes, err error) {
	if err = c.svc.JtUpdate(ctx, service.JtSaveInput{
		Id: req.Id, CnName: req.CnName, EnName: req.EnName, Avatar: req.Avatar,
		Link: req.Link, PicJumpLink: req.PicJumpLink, Location: req.Location,
		Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.JtUpdateRes{}, nil
}

func (c *Controller) JtDelete(ctx context.Context, req *v1.JtDeleteReq) (res *v1.JtDeleteRes, err error) {
	if err = c.svc.JtDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.JtDeleteRes{}, nil
}

// ---------- 敏感词 ----------

func (c *Controller) FwList(ctx context.Context, req *v1.FwListReq) (res *v1.FwListRes, err error) {
	list, total, err := c.svc.FwList(ctx, service.PageFilter{
		Status: -1, Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.FwListRes{Total: total, List: make([]v1.FwItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.FwItem{Id: r.Id, Word: r.Word, CreatedAt: r.CreatedAt})
	}
	return res, nil
}

func (c *Controller) FwAdd(ctx context.Context, req *v1.FwAddReq) (res *v1.FwAddRes, err error) {
	added, err := c.svc.FwAdd(ctx, req.Words)
	if err != nil {
		return nil, err
	}
	return &v1.FwAddRes{Added: added}, nil
}

func (c *Controller) FwDelete(ctx context.Context, req *v1.FwDeleteReq) (res *v1.FwDeleteRes, err error) {
	if err = c.svc.FwDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.FwDeleteRes{}, nil
}
