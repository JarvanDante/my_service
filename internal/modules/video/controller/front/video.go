// Package front 前台视频控制器(公开浏览)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/video/v1"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
)

type Controller struct{ svc service.IVideo }

func New(svc service.IVideo) *Controller { return &Controller{svc: svc} }

// toItem 只挑前台需要的字段, cover_key/source_key/created_by/status 不外发。
func toItem(v *service.VideoDTO) v1.Item {
	if v == nil {
		return v1.Item{}
	}
	return v1.Item{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, SourceUrl: v.SourceUrl,
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
