package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/comics/v1"
	"github.com/JarvanDante/my_service/internal/modules/comics/service"
)

func toModuleItem(d *service.ModuleDTO) v1.ModuleItem {
	return v1.ModuleItem{
		Id: d.Id, Name: d.Name, Position: d.Position, Style: d.Style, Icon: d.Icon,
		TagIds: d.TagIds, TagNames: d.TagNames, Size: d.Size, Rank: d.Rank, Status: d.Status,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func (c *Controller) ModuleList(ctx context.Context, req *v1.ModuleListReq) (res *v1.ModuleListRes, err error) {
	list, total, err := c.mod.List(ctx, service.ModuleFilter{
		Name: req.Name, Position: req.Position, Status: parseOptionalInt(req.Status),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ModuleListRes{Total: total, List: make([]v1.ModuleItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toModuleItem(d))
	}
	return res, nil
}

func (c *Controller) ModuleCreate(ctx context.Context, req *v1.ModuleCreateReq) (res *v1.ModuleCreateRes, err error) {
	id, err := c.mod.Create(ctx, service.ModuleInput{
		Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ModuleCreateRes{Id: id}, nil
}

func (c *Controller) ModuleUpdate(ctx context.Context, req *v1.ModuleUpdateReq) (res *v1.ModuleUpdateRes, err error) {
	if err = c.mod.Update(ctx, service.ModuleInput{
		Id: req.Id, Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.ModuleUpdateRes{}, nil
}

func (c *Controller) ModuleDelete(ctx context.Context, req *v1.ModuleDeleteReq) (res *v1.ModuleDeleteRes, err error) {
	if err = c.mod.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ModuleDeleteRes{}, nil
}
