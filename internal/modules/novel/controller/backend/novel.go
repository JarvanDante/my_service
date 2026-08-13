// Package backend 后台小说控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/novel/v1"
	"github.com/JarvanDante/my_service/internal/modules/novel/service"
)

type Controller struct{ svc service.INovel }

func New(svc service.INovel) *Controller { return &Controller{svc: svc} }

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
			Category: d.Category, Tags: d.Tags, IsVip: d.IsVip, Price: d.Price,
			FreeChapter: d.FreeChapter, ChapterCount: d.ChapterCount, WordCount: d.WordCount,
			IsAudio: d.IsAudio, ViewCount: d.ViewCount, BuyCount: d.BuyCount,
			LikeCount: d.LikeCount, UpdateStatus: d.UpdateStatus, Rank: d.Rank,
			Status: d.Status, PublishId: d.PublishId, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Title: req.Title, Author: req.Author, Cover: req.Cover, Intro: req.Intro,
		Category: req.Category, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price,
		FreeChapter: req.FreeChapter, IsAudio: req.IsAudio, UpdateStatus: req.UpdateStatus,
		Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Title: req.Title, Author: req.Author, Cover: req.Cover, Intro: req.Intro,
		Category: req.Category, Tags: req.Tags, IsVip: req.IsVip, Price: req.Price,
		FreeChapter: req.FreeChapter, IsAudio: req.IsAudio, UpdateStatus: req.UpdateStatus,
		Rank: req.Rank, Status: req.Status,
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
			Id: d.Id, NovelId: d.NovelId, Seq: d.Seq, Title: d.Title,
			Content: d.Content, WordCount: d.WordCount, AudioUrl: d.AudioUrl,
			Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) ChapterCreate(ctx context.Context, req *v1.ChapterCreateReq) (res *v1.ChapterCreateRes, err error) {
	id, err := c.svc.ChapterCreate(ctx, service.ChapterInput{
		NovelId: req.Id, Seq: req.Seq, Title: req.Title,
		Content: req.Content, AudioUrl: req.AudioUrl, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ChapterCreateRes{Id: id}, nil
}

func (c *Controller) ChapterUpdate(ctx context.Context, req *v1.ChapterUpdateReq) (res *v1.ChapterUpdateRes, err error) {
	if err = c.svc.ChapterUpdate(ctx, service.ChapterInput{
		Id: req.Id, Seq: req.Seq, Title: req.Title,
		Content: req.Content, AudioUrl: req.AudioUrl, Status: req.Status,
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
