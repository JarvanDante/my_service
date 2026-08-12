// Package front 前台标签控制器(公开浏览, 无需登录)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/tag/v1"
	"github.com/JarvanDante/my_service/internal/modules/tag/service"
)

type Controller struct{ svc service.ITag }

func New(svc service.ITag) *Controller { return &Controller{svc: svc} }

func (c *Controller) RepoList(ctx context.Context, req *v1.RepoListReq) (res *v1.RepoListRes, err error) {
	list, err := c.svc.Repo(ctx, req.Type, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.RepoListRes{List: make([]v1.RepoTagItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.RepoTagItem{Id: r.Id, Name: r.Name})
	}
	return res, nil
}
