// Code maintained manually (优惠券模板 + 用户券).
package entity

import "github.com/gogf/gf/v2/os/gtime"

// CouponTpl 优惠券模板。type: 1抵用券 2折扣券; scene: 1充值 2内容购买 3通用。
type CouponTpl struct {
	Id        int64       `json:"id"         orm:"id"`
	SiteId    int64       `json:"siteId"     orm:"site_id"`
	Name      string      `json:"name"       orm:"name"`
	Type      int         `json:"type"       orm:"type"`
	Scene     int         `json:"scene"      orm:"scene"`
	FaceValue float64     `json:"faceValue"  orm:"face_value"`
	Discount  int         `json:"discount"   orm:"discount"`
	Threshold float64     `json:"threshold"  orm:"threshold"`
	MaxDeduct float64     `json:"maxDeduct"  orm:"max_deduct"`
	Total     int         `json:"total"      orm:"total"`
	Issued    int         `json:"issued"     orm:"issued"`
	PerLimit  int         `json:"perLimit"   orm:"per_limit"`
	ExpireDay int         `json:"expireDay"  orm:"expire_day"`
	Status    int         `json:"status"     orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt"  orm:"updated_at"`
}

// UserCoupon 用户券(领取时对模板做值快照)。status: 1未使用 2已使用 3已过期。
type UserCoupon struct {
	Id        int64       `json:"id"        orm:"id"`
	SiteId    int64       `json:"siteId"    orm:"site_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	TplId     int64       `json:"tplId"     orm:"tpl_id"`
	Name      string      `json:"name"      orm:"name"`
	Type      int         `json:"type"      orm:"type"`
	Scene     int         `json:"scene"     orm:"scene"`
	FaceValue float64     `json:"faceValue" orm:"face_value"`
	Discount  int         `json:"discount"  orm:"discount"`
	Threshold float64     `json:"threshold" orm:"threshold"`
	MaxDeduct float64     `json:"maxDeduct" orm:"max_deduct"`
	Status    int         `json:"status"    orm:"status"`
	UsedAt    *gtime.Time `json:"usedAt"    orm:"used_at"`
	RefId     string      `json:"refId"     orm:"ref_id"`
	ExpireAt  *gtime.Time `json:"expireAt"  orm:"expire_at"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
