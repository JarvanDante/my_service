package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/video/v1"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
)

func toModuleItem(d *service.ModuleDTO) v1.ModuleItem {
	return v1.ModuleItem{
		Id: d.Id, Name: d.Name, Position: d.Position, Style: d.Style, Icon: d.Icon,
		TagIds: d.TagIds, TagNames: d.TagNames, Size: d.Size, Rank: d.Rank, Status: d.Status,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}

func moduleList(ctx context.Context, mod service.IModule, name, position, status string, page, size int) ([]v1.ModuleItem, int, error) {
	list, total, err := mod.List(ctx, service.ModuleFilter{
		Name: name, Position: position, Status: parseOptionalInt(status),
		Page: page, Size: size,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]v1.ModuleItem, 0, len(list))
	for _, d := range list {
		out = append(out, toModuleItem(d))
	}
	return out, total, nil
}

func (c *Controller) VideoModuleList(ctx context.Context, req *v1.VideoModuleListReq) (res *v1.VideoModuleListRes, err error) {
	list, total, err := moduleList(ctx, c.videoMod, req.Name, req.Position, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.VideoModuleListRes{List: list, Total: total}, nil
}

func (c *Controller) VideoModuleCreate(ctx context.Context, req *v1.VideoModuleCreateReq) (res *v1.VideoModuleCreateRes, err error) {
	id, err := c.videoMod.Create(ctx, service.ModuleInput{
		Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.VideoModuleCreateRes{Id: id}, nil
}

func (c *Controller) VideoModuleUpdate(ctx context.Context, req *v1.VideoModuleUpdateReq) (res *v1.VideoModuleUpdateRes, err error) {
	if err = c.videoMod.Update(ctx, service.ModuleInput{
		Id: req.Id, Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.VideoModuleUpdateRes{}, nil
}

func (c *Controller) VideoModuleDelete(ctx context.Context, req *v1.VideoModuleDeleteReq) (res *v1.VideoModuleDeleteRes, err error) {
	if err = c.videoMod.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.VideoModuleDeleteRes{}, nil
}

func (c *Controller) CartoonModuleList(ctx context.Context, req *v1.CartoonModuleListReq) (res *v1.CartoonModuleListRes, err error) {
	list, total, err := moduleList(ctx, c.cartoonMod, req.Name, req.Position, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.CartoonModuleListRes{List: list, Total: total}, nil
}

func (c *Controller) CartoonModuleCreate(ctx context.Context, req *v1.CartoonModuleCreateReq) (res *v1.CartoonModuleCreateRes, err error) {
	id, err := c.cartoonMod.Create(ctx, service.ModuleInput{
		Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CartoonModuleCreateRes{Id: id}, nil
}

func (c *Controller) CartoonModuleUpdate(ctx context.Context, req *v1.CartoonModuleUpdateReq) (res *v1.CartoonModuleUpdateRes, err error) {
	if err = c.cartoonMod.Update(ctx, service.ModuleInput{
		Id: req.Id, Name: req.Name, Position: req.Position, Style: req.Style, Icon: req.Icon,
		TagIds: req.TagIds, Size: req.Size, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.CartoonModuleUpdateRes{}, nil
}

func (c *Controller) CartoonModuleDelete(ctx context.Context, req *v1.CartoonModuleDeleteReq) (res *v1.CartoonModuleDeleteRes, err error) {
	if err = c.cartoonMod.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CartoonModuleDeleteRes{}, nil
}
