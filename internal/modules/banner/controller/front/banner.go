package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/banner/v1"
	"github.com/JarvanDante/my_service/internal/modules/banner/service"
)

type Controller struct{ svc service.IBanner }

func New(svc service.IBanner) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := c.svc.FrontList(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Position: r.Position, Title: r.Title, CoverUrl: r.CoverUrl, Link: r.Link,
		})
	}
	return res, nil
}
