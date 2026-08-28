package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/video/v1"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
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

func (c *Controller) CartoonCategoryList(ctx context.Context, req *v1.CartoonCategoryListReq) (res *v1.CartoonCategoryListRes, err error) {
	list, total, err := c.cartoonCat.List(ctx, service.CategoryFilter{
		Kind: parseOptionalInt(req.Kind), Status: parseOptionalInt(req.Status),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.CartoonCategoryListRes{Total: total, List: make([]v1.CategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.CategoryItem{
			Id: d.Id, Name: d.Name, Kind: d.Kind, Rank: d.Rank, Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) CartoonCategoryCreate(ctx context.Context, req *v1.CartoonCategoryCreateReq) (res *v1.CartoonCategoryCreateRes, err error) {
	id, err := c.cartoonCat.Create(ctx, service.CategoryInput{
		Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CartoonCategoryCreateRes{Id: id}, nil
}

func (c *Controller) CartoonCategoryUpdate(ctx context.Context, req *v1.CartoonCategoryUpdateReq) (res *v1.CartoonCategoryUpdateRes, err error) {
	if err = c.cartoonCat.Update(ctx, service.CategoryInput{
		Id: req.Id, Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.CartoonCategoryUpdateRes{}, nil
}

func (c *Controller) CartoonCategoryDelete(ctx context.Context, req *v1.CartoonCategoryDeleteReq) (res *v1.CartoonCategoryDeleteRes, err error) {
	if err = c.cartoonCat.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CartoonCategoryDeleteRes{}, nil
}

func (c *Controller) DouyinCategoryList(ctx context.Context, req *v1.DouyinCategoryListReq) (res *v1.DouyinCategoryListRes, err error) {
	list, total, err := c.douyinCat.List(ctx, service.CategoryFilter{
		Kind: parseOptionalInt(req.Kind), Status: parseOptionalInt(req.Status),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.DouyinCategoryListRes{Total: total, List: make([]v1.CategoryItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.CategoryItem{
			Id: d.Id, Name: d.Name, Kind: d.Kind, Rank: d.Rank, Status: d.Status, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) DouyinCategoryCreate(ctx context.Context, req *v1.DouyinCategoryCreateReq) (res *v1.DouyinCategoryCreateRes, err error) {
	id, err := c.douyinCat.Create(ctx, service.CategoryInput{
		Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DouyinCategoryCreateRes{Id: id}, nil
}

func (c *Controller) DouyinCategoryUpdate(ctx context.Context, req *v1.DouyinCategoryUpdateReq) (res *v1.DouyinCategoryUpdateRes, err error) {
	if err = c.douyinCat.Update(ctx, service.CategoryInput{
		Id: req.Id, Name: req.Name, Kind: req.Kind, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.DouyinCategoryUpdateRes{}, nil
}

func (c *Controller) DouyinCategoryDelete(ctx context.Context, req *v1.DouyinCategoryDeleteReq) (res *v1.DouyinCategoryDeleteRes, err error) {
	if err = c.douyinCat.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DouyinCategoryDeleteRes{}, nil
}
