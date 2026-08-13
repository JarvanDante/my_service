// Package front 前台排行/热搜控制器(公开)。
package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/ranks/v1"
	"github.com/JarvanDante/my_service/internal/modules/ranks/service"
)

type Controller struct{ svc service.IRank }

func New(svc service.IRank) *Controller { return &Controller{svc: svc} }

func (c *Controller) Rank(ctx context.Context, req *v1.RankReq) (res *v1.RankRes, err error) {
	list, err := c.svc.Rank(ctx, req.MediaType, req.Period)
	if err != nil {
		return nil, err
	}
	res = &v1.RankRes{List: make([]v1.RankItem, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.RankItem{
			ContentId: r.ContentId, MediaType: r.MediaType, Score: r.Score, RankNo: r.RankNo,
		})
	}
	return res, nil
}

func (c *Controller) Hot(ctx context.Context, req *v1.HotReq) (res *v1.HotRes, err error) {
	words, err := c.svc.HotKeywords(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.HotRes{List: make([]v1.HotItem, 0, len(words))}
	for _, w := range words {
		res.List = append(res.List, v1.HotItem{Keyword: w})
	}
	return res, nil
}
