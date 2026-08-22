package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/group/v1"
	"github.com/JarvanDante/my_service/internal/modules/group/service"
)

type Controller struct{ svc service.IGroup }

func New(svc service.IGroup) *Controller { return &Controller{svc: svc} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := c.svc.FrontList(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.ListRes{List: make([]v1.Item, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.Item{
			Id: r.Id, Name: r.Name, Intro: r.Intro, Avatar: r.Avatar,
			Link: r.Link, Platform: r.Platform,
		})
	}
	return res, nil
}
