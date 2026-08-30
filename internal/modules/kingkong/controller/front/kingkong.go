package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/kingkong/v1"
	"github.com/JarvanDante/my_service/internal/modules/kingkong/service"
)

type Controller struct{ svc service.IKingkong }

func New(svc service.IKingkong) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := c.svc.FrontList(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Name: r.Name, IconUrl: r.IconUrl,
			OpenMode: r.OpenMode, Link: r.Link, AppLink: r.AppLink,
			Position: r.Position, Sort: r.Sort,
		})
	}
	return res, nil
}
