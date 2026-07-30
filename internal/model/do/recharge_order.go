// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RechargeOrder is the golang structure of table recharge_order for DAO operations like Where/Data.
type RechargeOrder struct {
	g.Meta    `orm:"table:recharge_order, do:true"`
	Id        any         //
	SiteId    any         //
	OrderNo   any         //
	UserId    any         //
	PackageId any         //
	Amount    any         //
	Coin      any         //
	Status    any         //
	PayAt     *gtime.Time //
	CreatedAt *gtime.Time //
}
