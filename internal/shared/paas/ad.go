package paas

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type AdItem struct {
	CampaignId string `json:"campaign_id"`
	CreativeId string `json:"creative_id"`
	Title      string `json:"title"`
	MediaURL   string `json:"media_url"`
	LinkURL    string `json:"link_url"`
	SlotCode   string `json:"slot_code"`
	Priority   int    `json:"priority"`
	Weight     int    `json:"weight"`
}

func adBase(ctx context.Context) string {
	return strings.TrimRight(g.Cfg().MustGet(ctx, "paas.ad_base").String(), "/")
}

func adJSON(ctx context.Context, method, path string, body any, out any) error {
	base := adBase(ctx)
	if base == "" {
		return gerror.New("未配置 paas.ad_base")
	}
	key := g.Cfg().MustGet(ctx, "paas.app_key").String()
	secret := g.Cfg().MustGet(ctx, "paas.app_secret").String()
	if key == "" || secret == "" {
		return gerror.New("未配置 paas.app_key / app_secret")
	}
	c := g.Client().SetHeader("X-App-Key", key).SetHeader("X-App-Secret", secret)
	var raw []byte
	switch method {
	case "GET":
		r, e := c.Get(ctx, base+path)
		if e != nil {
			return gerror.Wrap(e, "请求广告中台失败")
		}
		defer r.Close()
		raw = r.ReadAll()
	case "POST":
		r, e := c.ContentJson().Post(ctx, base+path, body)
		if e != nil {
			return gerror.Wrap(e, "请求广告中台失败")
		}
		defer r.Close()
		raw = r.ReadAll()
	default:
		return gerror.New("不支持的方法")
	}
	return parseEnvelope(raw, out)
}

func FetchAds(ctx context.Context, slotCode string, limit int) ([]AdItem, error) {
	if adBase(ctx) == "" || strings.TrimSpace(slotCode) == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	q := url.Values{}
	q.Set("slot_code", strings.TrimSpace(slotCode))
	q.Set("limit", fmt.Sprintf("%d", limit))
	var data struct {
		List []AdItem `json:"list"`
	}
	if err := adJSON(ctx, "GET", "/open/ads?"+q.Encode(), nil, &data); err != nil {
		return nil, err
	}
	return data.List, nil
}

func ReportAdEvent(ctx context.Context, eventType, campaignId, creativeId, slotCode string) error {
	if adBase(ctx) == "" {
		return nil
	}
	return adJSON(ctx, "POST", "/open/events", g.Map{
		"event_type":  eventType,
		"campaign_id": campaignId,
		"creative_id": creativeId,
		"slot_code":   slotCode,
	}, nil)
}
