// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id             int64       `json:"id"             orm:"id"               description:""` //
	Username       string      `json:"username"       orm:"username"         description:""` //
	Nickname       string      `json:"nickname"       orm:"nickname"         description:""` //
	Phone          string      `json:"phone"          orm:"phone"            description:""` //
	CountryCode    string      `json:"countryCode"    orm:"country_code"     description:""` //
	Password       string      `json:"password"       orm:"password"         description:""` //
	Slat           string      `json:"slat"           orm:"slat"             description:""` //
	DeviceId       string      `json:"deviceId"       orm:"device_id"        description:""` //
	DeviceType     string      `json:"deviceType"     orm:"device_type"      description:""` //
	DeviceVersion  string      `json:"deviceVersion"  orm:"device_version"   description:""` //
	DeviceExt      string      `json:"deviceExt"      orm:"device_ext"       description:""` //
	Img            string      `json:"img"            orm:"img"              description:""` //
	BgImg          string      `json:"bgImg"          orm:"bg_img"           description:""` //
	Signature      string      `json:"signature"      orm:"signature"        description:""` //
	Sex            int         `json:"sex"            orm:"sex"              description:""` //
	Level          int         `json:"level"          orm:"level"            description:""` //
	Tag            string      `json:"tag"            orm:"tag"              description:""` //
	Balance        float64     `json:"balance"        orm:"balance"          description:""` //
	Credit         float64     `json:"credit"         orm:"credit"           description:""` //
	GiftCount      float64     `json:"giftCount"      orm:"gift_count"       description:""` //
	MoneyCount     float64     `json:"moneyCount"     orm:"money_count"      description:""` //
	SendCount      float64     `json:"sendCount"      orm:"send_count"       description:""` //
	HasBuy         int         `json:"hasBuy"         orm:"has_buy"          description:""` //
	Rights         string      `json:"rights"         orm:"rights"           description:""` //
	WithdrawInfo   string      `json:"withdrawInfo"   orm:"withdraw_info"    description:""` //
	GroupId        int64       `json:"groupId"        orm:"group_id"         description:""` //
	GroupName      string      `json:"groupName"      orm:"group_name"       description:""` //
	GroupRate      int         `json:"groupRate"      orm:"group_rate"       description:""` //
	GroupStartTime int64       `json:"groupStartTime" orm:"group_start_time" description:""` //
	GroupEndTime   int64       `json:"groupEndTime"   orm:"group_end_time"   description:""` //
	ParentId       int64       `json:"parentId"       orm:"parent_id"        description:""` //
	ParentName     string      `json:"parentName"     orm:"parent_name"      description:""` //
	ChannelName    string      `json:"channelName"    orm:"channel_name"     description:""` //
	ShareNum       int         `json:"shareNum"       orm:"share_num"        description:""` //
	Fans           int         `json:"fans"           orm:"fans"             description:""` //
	Follow         int         `json:"follow"         orm:"follow"           description:""` //
	Country        string      `json:"country"        orm:"country"          description:""` //
	Province       string      `json:"province"       orm:"province"         description:""` //
	City           string      `json:"city"           orm:"city"             description:""` //
	Location       string      `json:"location"       orm:"location"         description:""` //
	RegisterArea   string      `json:"registerArea"   orm:"register_area"    description:""` //
	IsChina        int         `json:"isChina"        orm:"is_china"         description:""` //
	IsDisabled     int         `json:"isDisabled"     orm:"is_disabled"      description:""` //
	ErrorMsg       string      `json:"errorMsg"       orm:"error_msg"        description:""` //
	IsSystem       int         `json:"isSystem"       orm:"is_system"        description:""` //
	RegisterAt     *gtime.Time `json:"registerAt"     orm:"register_at"      description:""` //
	RegisterDate   int         `json:"registerDate"   orm:"register_date"    description:""` //
	RegisterIp     string      `json:"registerIp"     orm:"register_ip"      description:""` //
	LoginNum       int         `json:"loginNum"       orm:"login_num"        description:""` //
	LastLoginAt    *gtime.Time `json:"lastLoginAt"    orm:"last_login_at"    description:""` //
	LastDate       int         `json:"lastDate"       orm:"last_date"        description:""` //
	LastIp         string      `json:"lastIp"         orm:"last_ip"          description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""` //
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"       description:""` //
}
