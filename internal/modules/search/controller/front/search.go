// Package front 前台搜索控制器(公开)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/search/v1"
	"github.com/JarvanDante/my_service/internal/modules/search/service"
)

type Controller struct{ svc service.ISearch }

func New(svc service.ISearch) *Controller { return &Controller{svc: svc} }

// toItems nil 切片保持 nil, 让 v1 的 omitempty 生效(用不到的分组直接不下发);
// 非 nil 但为空时返回空数组, 前端拿到 [] 而不是 null。
func toItems(list []*service.ItemDTO) []v1.Item {
	if list == nil {
		return nil
	}
	out := make([]v1.Item, 0, len(list))
	for _, d := range list {
		out = append(out, v1.Item{
			Id: d.Id, MediaType: d.MediaType, Title: d.Title, Cover: d.Cover,
			Author: d.Author, Price: d.Price, IsVip: d.IsVip,
			ViewCount: d.ViewCount, CreatedAt: d.CreatedAt,
		})
	}
	return out
}

func (c *Controller) Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error) {
	r, err := c.svc.Search(ctx, service.SearchInput{
		Keyword: req.Keyword, Type: req.Type, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SearchRes{
		Videos: toItems(r.Videos), Posts: toItems(r.Posts), Comics: toItems(r.Comics),
		Novels: toItems(r.Novels), Photos: toItems(r.Photos),
		List: toItems(r.List), TotalHit: r.TotalHit, Total: r.Total,
	}, nil
}

func (c *Controller) Suggest(ctx context.Context, req *v1.SuggestReq) (res *v1.SuggestRes, err error) {
	list, err := c.svc.Suggest(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}
	return &v1.SuggestRes{List: list}, nil
}
