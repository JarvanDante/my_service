// Package backend 后台漫画控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/comics/v1"
	"github.com/JarvanDante/my_service/internal/modules/comics/service"
)

type Controller struct {
	svc service.IComics
	cat service.ICategory
	mod service.IModule
}

func New(svc service.IComics, cat service.ICategory, mod service.IModule) *Controller {
	return &Controller{svc: svc, cat: cat, mod: mod}
}

func toPicDTO(in []v1.Pic) []service.PicDTO {
	if in == nil {
		return nil
	}
	out := make([]service.PicDTO, 0, len(in))
	for _, p := range in {
		out = append(out, service.PicDTO{Url: p.Url, Width: p.Width, Height: p.Height})
	}
	return out
}

func toPicV1(in []service.PicDTO) []v1.Pic {
	out := make([]v1.Pic, 0, len(in))
	for _, p := range in {
		out = append(out, v1.Pic{Url: p.Url, Width: p.Width, Height: p.Height})
	}
	return out
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	status := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			status = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Status: status, Category: req.Category, Keyword: req.Keyword,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.Item{
			Id: d.Id, Title: d.Title, Author: d.Author, Cover: d.Cover, Intro: d.Intro,
			Category: d.Category, Categories: d.Categories, Tags: d.Tags, IsVip: d.IsVip, Price: d.Price,
			FreeChapter: d.FreeChapter, ChapterCount: d.ChapterCount, ViewCount: d.ViewCount,
			BuyCount: d.BuyCount, LikeCount: d.LikeCount, UpdateStatus: d.UpdateStatus,
			Rank: d.Rank, IsRecommend: d.IsRecommend, Status: d.Status, PublishId: d.PublishId, MediaCode: d.MediaCode, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Title: req.Title, Author: req.Author, Cover: req.Cover, Intro: req.Intro,
		Category: req.Category, Categories: req.Categories, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price,
		FreeChapter: req.FreeChapter, UpdateStatus: req.UpdateStatus,
		Rank: req.Rank, IsRecommend: req.IsRecommend, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Title: req.Title, Author: req.Author, Cover: req.Cover, Intro: req.Intro,
		Category: req.Category, Categories: req.Categories, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price,
		FreeChapter: req.FreeChapter, UpdateStatus: req.UpdateStatus,
		Rank: req.Rank, IsRecommend: req.IsRecommend, Status: req.Status,
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

func (c *Controller) Chapters(ctx context.Context, req *v1.ChaptersReq) (res *v1.ChaptersRes, err error) {
	list, total, err := c.svc.ChapterList(ctx, req.Id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.ChaptersRes{Total: total, List: make([]v1.ChapterItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.ChapterItem{
			Id: d.Id, ComicsId: d.ComicsId, Seq: d.Seq, Title: d.Title,
			Pics: toPicV1(d.Pics), PicCount: d.PicCount, Status: d.Status,
			CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) ChapterGet(ctx context.Context, req *v1.ChapterGetReq) (res *v1.ChapterGetRes, err error) {
	d, err := c.svc.ChapterGet(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ChapterGetRes{ChapterItem: v1.ChapterItem{
		Id: d.Id, ComicsId: d.ComicsId, Seq: d.Seq, Title: d.Title,
		Pics: toPicV1(d.Pics), PicCount: d.PicCount, Status: d.Status,
		CreatedAt: d.CreatedAt,
	}}, nil
}

func (c *Controller) ChapterCreate(ctx context.Context, req *v1.ChapterCreateReq) (res *v1.ChapterCreateRes, err error) {
	id, err := c.svc.ChapterCreate(ctx, service.ChapterInput{
		ComicsId: req.Id, Seq: req.Seq, Title: req.Title,
		Pics: toPicDTO(req.Pics), Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ChapterCreateRes{Id: id}, nil
}

func (c *Controller) ChapterUpdate(ctx context.Context, req *v1.ChapterUpdateReq) (res *v1.ChapterUpdateRes, err error) {
	if err = c.svc.ChapterUpdate(ctx, service.ChapterInput{
		Id: req.Id, Seq: req.Seq, Title: req.Title,
		Pics: toPicDTO(req.Pics), Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.ChapterUpdateRes{}, nil
}

func (c *Controller) ChapterDelete(ctx context.Context, req *v1.ChapterDeleteReq) (res *v1.ChapterDeleteRes, err error) {
	if err = c.svc.ChapterDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ChapterDeleteRes{}, nil
}

func (c *Controller) MediaComics(ctx context.Context, req *v1.MediaComicsListReq) (res *v1.MediaComicsListRes, err error) {
	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := c.svc.ListMediaComics(ctx, page, size, req.Keyword)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MediaComicsItem, 0, len(list))
	for _, a := range list {
		items = append(items, v1.MediaComicsItem{
			Id: a.Id, Title: a.Title, CoverUrl: a.CoverUrl, Intro: a.Intro,
			ChapterCount: a.ChapterCount, Picked: a.Picked, LocalId: a.LocalId,
		})
	}
	return &v1.MediaComicsListRes{List: items, Total: total, Page: page, Size: size}, nil
}

func (c *Controller) MediaPick(ctx context.Context, req *v1.MediaComicsPickReq) (res *v1.MediaComicsPickRes, err error) {
	id, err := c.svc.PickMedia(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.MediaComicsPickRes{Id: id}, nil
}
