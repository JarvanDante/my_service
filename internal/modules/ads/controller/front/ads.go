package front

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/front/ads/v1"
	"github.com/JarvanDante/my_service/internal/shared/paas"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct{}

func New() *Controller { return &Controller{} }

func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	list, err := paas.FetchAds(ctx, req.SlotCode, req.Limit)
	if err != nil {
		g.Log().Warningf(ctx, "拉取广告失败 slot=%s: %v", req.SlotCode, err)
		return &v1.ListRes{List: []v1.Item{}}, nil
	}
	res = &v1.ListRes{List: make([]v1.Item, 0, len(list))}
	for _, a := range list {
		res.List = append(res.List, v1.Item{
			CampaignId: a.CampaignId, CreativeId: a.CreativeId, Title: a.Title,
			MediaURL: a.MediaURL, LinkURL: a.LinkURL, SlotCode: a.SlotCode,
		})
	}
	return res, nil
}

func (c *Controller) Event(ctx context.Context, req *v1.EventReq) (res *v1.EventRes, err error) {
	if err := paas.ReportAdEvent(ctx, req.EventType, req.CampaignId, req.CreativeId, req.SlotCode); err != nil {
		g.Log().Warningf(ctx, "广告上报失败: %v", err)
	}
	return &v1.EventRes{Ok: true}, nil
}
