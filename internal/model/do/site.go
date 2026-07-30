// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Site is the golang structure of table site for DAO operations like Where/Data.
type Site struct {
	g.Meta    `orm:"table:site, do:true"`
	Id        any         //
	Name      any         //
	Domain    any         //
	Appid     any         //
	Status    any         //
	Config    any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
