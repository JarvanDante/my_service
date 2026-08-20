// Package front 前台视频控制器(公开浏览)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/video/v1"
	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
)

type Controller struct {
	svc service.IVideo
	cat service.ICategory
}

func New(svc service.IVideo, cat service.ICategory) *Controller {
	return &Controller{svc: svc, cat: cat}
}

// toItem 只挑前台需要的字段, cover_key/source_key/created_by/status 不外发。
func toItem(v *service.VideoDTO) v1.Item {
	if v == nil {
		return v1.Item{}
	}
	return v1.Item{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, SourceUrl: v.SourceUrl,
		Category: v.Category, Categories: v.Categories,
		Duration: v.Duration, CreatedAt: v.CreatedAt,
	}
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.svc.FrontList(ctx, service.FrontListInput{
		Keyword: req.Keyword, Sort: req.Sort, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.Item, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.ListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	v, err := c.svc.FrontDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{Item: toItem(v)}, nil
}

func (c *Controller) CartoonCategoryList(ctx context.Context, _ *v1.CartoonCategoryListReq) (res *v1.CartoonCategoryListRes, err error) {
	if c.cat == nil {
		return &v1.CartoonCategoryListRes{List: []v1.FrontCategoryItem{}}, nil
	}
	list, err := c.cat.Repo(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.CartoonCategoryListRes{List: make([]v1.FrontCategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.FrontCategoryItem{Id: d.Id, Name: d.Name, Kind: d.Kind})
	}
	return res, nil
}

func (c *Controller) CartoonList(ctx context.Context, req *v1.CartoonListReq) (res *v1.CartoonListRes, err error) {
	dto, err := c.svc.FrontList(ctx, service.FrontListInput{
		Keyword: req.Keyword, Category: req.Category, Kind: entity.VideoKindCartoon,
		Sort: req.Sort, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.Item, 0, len(dto.List))
	for _, v := range dto.List {
		items = append(items, toItem(v))
	}
	return &v1.CartoonListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}
