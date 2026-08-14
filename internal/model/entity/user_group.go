package entity

import "github.com/gogf/gf/v2/os/gtime"

type UserGroup struct {
	Id               int64       `json:"id"               orm:"id"`
	SiteId           int64       `json:"siteId"           orm:"site_id"`
	Name             string      `json:"name"             orm:"name"`
	Rate             int         `json:"rate"             orm:"rate"`
	Rights           string      `json:"rights"           orm:"rights"`
	Remark           string      `json:"remark"           orm:"remark"`
	Sort             int         `json:"sort"             orm:"sort"`
	Status           int         `json:"status"           orm:"status"`
	Img              string      `json:"img"              orm:"img"`
	TitleHeat        string      `json:"titleHeat"        orm:"title_heat"`
	TitleDescription string      `json:"titleDescription" orm:"title_description"`
	TitlePicture     string      `json:"titlePicture"     orm:"title_picture"`
	Level            int         `json:"level"            orm:"level"`
	PromotionType    int         `json:"promotionType"    orm:"promotion_type"`
	Price            float64     `json:"price"            orm:"price"`
	OldPrice         float64     `json:"oldPrice"         orm:"old_price"`
	DayNum           int         `json:"dayNum"           orm:"day_num"`
	GiftNum          int         `json:"giftNum"          orm:"gift_num"`
	DownloadNum      int         `json:"downloadNum"      orm:"download_num"`
	DayTips          string      `json:"dayTips"          orm:"day_tips"`
	PriceTips        string      `json:"priceTips"        orm:"price_tips"`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"`
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"`
}
