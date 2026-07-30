// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// VipPackage is the golang structure of table vip_package for DAO operations like Where/Data.
type VipPackage struct {
	g.Meta    `orm:"table:vip_package, do:true"`
	Id        any         //
	SiteId    any         //
	Name      any         //
	Days      any         //
	Price     any         //
	GroupId   any         //
	Sort      any         //
	Status    any         //
	CreatedAt *gtime.Time //
}
