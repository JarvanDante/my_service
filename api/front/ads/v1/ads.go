package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	CampaignId string `json:"campaign_id"`
	CreativeId string `json:"creative_id"`
	Title      string `json:"title"`
	MediaURL   string `json:"media_url"`
	LinkURL    string `json:"link_url"`
	SlotCode   string `json:"slot_code"`
}

type ListReq struct {
	g.Meta   `path:"/ads" method:"get" tags:"Front/Ads" summary:"按广告位拉取"`
	SlotCode string `json:"slot_code" v:"required#slot_code必填"`
	Limit    int    `json:"limit"`
}
type ListRes struct {
	List []Item `json:"list"`
}

type EventReq struct {
	g.Meta     `path:"/ads/event" method:"post" tags:"Front/Ads" summary:"广告曝光点击"`
	EventType  string `json:"event_type" v:"required|in:impression,click"`
	CampaignId string `json:"campaign_id"`
	CreativeId string `json:"creative_id"`
	SlotCode   string `json:"slot_code"`
}
type EventRes struct {
	Ok bool `json:"ok"`
}
