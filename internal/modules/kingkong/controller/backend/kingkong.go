package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/kingkong/v1"
	"github.com/JarvanDante/my_service/internal/modules/kingkong/service"
)

type Controller struct{ svc service.IKingkong }

func New(svc service.IKingkong) *Controller { return &Controller{svc: svc} }

func toItem(r *service.ItemDTO) v1.Item {
	return v1.Item{
		Id: r.Id, Name: r.Name, IconUrl: r.IconUrl,
		OpenMode: r.OpenMode, OpenModeName: r.OpenModeName,
		Link: r.Link, AppLink: r.AppLink, LinkLabel: r.LinkLabel,
		Position: r.Position, PositionName: r.PositionName,
		Sort: r.Sort, Status: r.Status, StatusText: r.StatusText,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	statusFilter := -1
	if req.Status != "" {
		if v, e := strconv.Atoi(req.Status); e == nil {
			statusFilter = v
		}
	}
	list, total, err := c.svc.List(ctx, service.ListFilter{
		Name: req.Name, Position: req.Position, Status: statusFilter,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{Total: total, List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, toItem(r))
	}
	return res, nil
}

func (c *Controller) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.svc.Create(ctx, service.SaveInput{
		Name: req.Name, IconUrl: req.IconUrl, OpenMode: req.OpenMode,
		Link: req.Link, AppLink: req.AppLink, Position: req.Position,
		Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err = c.svc.Update(ctx, service.SaveInput{
		Id: req.Id, Name: req.Name, IconUrl: req.IconUrl, OpenMode: req.OpenMode,
		Link: req.Link, AppLink: req.AppLink, Position: req.Position,
		Sort: req.Sort, Status: req.Status,
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
