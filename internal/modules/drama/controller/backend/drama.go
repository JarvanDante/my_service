// Package backend 后台短剧控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/drama/v1"
	"github.com/JarvanDante/my_service/internal/modules/drama/service"
)

type Controller struct {
	svc service.IDrama
	cat service.ICategory
	mod service.IModule
}

func New(svc service.IDrama, cat service.ICategory, mod service.IModule) *Controller {
	return &Controller{svc: svc, cat: cat, mod: mod}
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	status := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: status, Category: req.Category, Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, Title: d.Title, Cover: d.Cover, Intro: d.Intro, Category: d.Category,
			Categories: d.Categories, Tags: d.Tags, IsVip: d.IsVip, Price: d.Price, FreeEpisode: d.FreeEpisode,
			EpisodeCount: d.EpisodeCount, ViewCount: d.ViewCount, BuyCount: d.BuyCount, LikeCount: d.LikeCount,
			UpdateStatus: d.UpdateStatus, Rank: d.Rank, IsRecommend: d.IsRecommend, Status: d.Status,
			MediaCode: d.MediaCode, PublishedAt: d.PublishedAt, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Title: req.Title, Cover: req.Cover, Intro: req.Intro, Category: req.Category,
		Categories: req.Categories, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price, FreeEpisode: req.FreeEpisode,
		UpdateStatus: req.UpdateStatus, Rank: req.Rank, IsRecommend: req.IsRecommend, Status: req.Status,
		PublishedAt: req.PublishedAt,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err = c.svc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

func (c *Controller) Audit(ctx context.Context, req *v1.AuditReq) (res *v1.AuditRes, err error) {
	if err = c.svc.Audit(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return &v1.AuditRes{}, nil
}

func (c *Controller) Episodes(ctx context.Context, req *v1.EpisodesReq) (res *v1.EpisodesRes, err error) {
	list, total, err := c.svc.EpisodeList(ctx, req.Id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.EpisodesRes{Total: total, List: make([]v1.EpisodeItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.EpisodeItem{
			Id: d.Id, DramaId: d.DramaId, Seq: d.Seq, Title: d.Title, MediaCode: d.MediaCode,
			Duration: d.Duration, Cover: d.Cover, PayMode: d.PayMode, Price: d.Price,
			Status: d.Status, MediaReady: d.MediaReady, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) EpisodeUpdate(ctx context.Context, req *v1.EpisodeUpdateReq) (res *v1.EpisodeUpdateRes, err error) {
	if err = c.svc.EpisodeUpdate(ctx, service.EpisodeInput{
		Id: req.Id, PayMode: req.PayMode, Price: req.Price, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.EpisodeUpdateRes{}, nil
}

func (c *Controller) Sync(ctx context.Context, req *v1.SyncReq) (res *v1.SyncRes, err error) {
	d, err := c.svc.SyncEpisodes(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.SyncRes{Created: d.Created, Updated: d.Updated, Missing: d.Missing, Total: d.Total}, nil
}

func (c *Controller) MediaDramas(ctx context.Context, req *v1.MediaDramaListReq) (res *v1.MediaDramaListRes, err error) {
	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := c.svc.ListMediaDramas(ctx, page, size, req.Keyword)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MediaDramaItem, 0, len(list))
	for _, a := range list {
		items = append(items, v1.MediaDramaItem{
			Id: a.Id, Title: a.Title, CoverUrl: a.CoverUrl, Intro: a.Intro, ChapterCount: a.ChapterCount,
			EpisodeCount: a.EpisodeCount, UpdateStatus: a.UpdateStatus, Picked: a.Picked, LocalId: a.LocalId,
		})
	}
	return &v1.MediaDramaListRes{List: items, Total: total, Page: page, Size: size}, nil
}

func (c *Controller) MediaPick(ctx context.Context, req *v1.MediaDramaPickReq) (res *v1.MediaDramaPickRes, err error) {
	id, err := c.svc.PickMedia(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.MediaDramaPickRes{Id: id}, nil
}
