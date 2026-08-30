package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/system/v1"
	"github.com/JarvanDante/my_service/internal/modules/system/service"
)

type Controller struct{ svc service.ISystem }

func New(svc service.ISystem) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := c.svc.FrontNotices(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.Item, 0, len(list))}
	for _, n := range list {
		res.List = append(res.List, v1.Item{Id: n.Id, Title: n.Title, Content: n.Content})
	}
	return res, nil
}
