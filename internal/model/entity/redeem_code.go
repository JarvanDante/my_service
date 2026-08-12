// Code maintained manually (兑换码 + 使用记录).
package entity

import "github.com/gogf/gf/v2/os/gtime"

type RedeemCode struct {
	Id         int64       `json:"id"         orm:"id"`
	SiteId     int64       `json:"siteId"     orm:"site_id"`
	Name       string      `json:"name"       orm:"name"`
	Code       string      `json:"code"       orm:"code"`
	CardType   int         `json:"cardType"   orm:"card_type"` // 1金币(当前仅支持)
	Value      int         `json:"value"      orm:"value"`
	TotalTimes int         `json:"totalTimes" orm:"total_times"`
	UsedTimes  int         `json:"usedTimes"  orm:"used_times"`
	Status     int         `json:"status"     orm:"status"` // 1启用 0禁用
	ExpiredAt  *gtime.Time `json:"expiredAt"  orm:"expired_at"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}

type RedeemCodeRecord struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	CodeId    int64       `json:"codeId"    orm:"code_id"`
	Code      string      `json:"code"      orm:"code"`
	Name      string      `json:"name"      orm:"name"`
	CardType  int         `json:"cardType"  orm:"card_type"`
	Value     int         `json:"value"     orm:"value"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
