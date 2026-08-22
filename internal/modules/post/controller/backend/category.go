package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/post/v1"
	"github.com/JarvanDante/my_service/internal/modules/post/service"
)

func parseOptionalInt(raw string) int {
	if raw == "" {
		return -1
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return -1
	}
	return v
}

func (c *Controller) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	if c.cat == nil {
		return &v1.CategoryListRes{List: []v1.CategoryItem{}}, nil
	}
	list, total, err := c.cat.List(ctx, service.CategoryFilter{
		Kind: parseOptionalInt(req.Kind), Status: parseOptionalInt(req.Status),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.CategoryListRes{Total: total, List: make([]v1.CategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.CategoryItem{
			Id: d.Id, Name: d.Name, Kind: d.Kind, Rank: d.Rank, Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) CategoryCreate(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error) {
	id, err := c.cat.Create(ctx, service.CategoryInput{
		Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CategoryCreateRes{Id: id}, nil
}

func (c *Controller) CategoryUpdate(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error) {
	if err = c.cat.Update(ctx, service.CategoryInput{
		Id: req.Id, Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.CategoryUpdateRes{}, nil
}

func (c *Controller) CategoryDelete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error) {
	if err = c.cat.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CategoryDeleteRes{}, nil
}
