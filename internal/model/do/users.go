// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta         `orm:"table:users, do:true"`
	Id             any         //
	Username       any         //
	Nickname       any         //
	Phone          any         //
	CountryCode    any         //
	Password       any         //
	Slat           any         //
	DeviceId       any         //
	DeviceType     any         //
	DeviceVersion  any         //
	DeviceExt      any         //
	Img            any         //
	BgImg          any         //
	Signature      any         //
	Sex            any         //
	Level          any         //
	Tag            any         //
	Balance        any         //
	Credit         any         //
	GiftCount      any         //
	MoneyCount     any         //
	SendCount      any         //
	HasBuy         any         //
	Rights         any         //
	WithdrawInfo   any         //
	GroupId        any         //
	GroupName      any         //
	GroupRate      any         //
	GroupStartTime any         //
	GroupEndTime   any         //
	ParentId       any         //
	ParentName     any         //
	ChannelName    any         //
	ShareNum       any         //
	Fans           any         //
	Follow         any         //
	Country        any         //
	Province       any         //
	City           any         //
	Location       any         //
	RegisterArea   any         //
	IsChina        any         //
	IsDisabled     any         //
	ErrorMsg       any         //
	IsSystem       any         //
	RegisterAt     *gtime.Time //
	RegisterDate   any         //
	RegisterIp     any         //
	LoginNum       any         //
	LastLoginAt    *gtime.Time //
	LastDate       any         //
	LastIp         any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
