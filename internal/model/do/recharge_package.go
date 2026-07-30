// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RechargePackage is the golang structure of table recharge_package for DAO operations like Where/Data.
type RechargePackage struct {
	g.Meta    `orm:"table:recharge_package, do:true"`
	Id        any         //
	SiteId    any         //
	Name      any         //
	Amount    any         //
	Coin      any         //
	Bonus     any         //
	Sort      any         //
	Status    any         //
	CreatedAt *gtime.Time //
}
