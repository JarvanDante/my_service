// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// VipLog is the golang structure of table vip_log for DAO operations like Where/Data.
type VipLog struct {
	g.Meta    `orm:"table:vip_log, do:true"`
	Id        any         //
	SiteId    any         //
	UserId    any         //
	PackageId any         //
	Days      any         //
	Price     any         //
	StartAt   any         //
	EndAt     any         //
	CreatedAt *gtime.Time //
}
