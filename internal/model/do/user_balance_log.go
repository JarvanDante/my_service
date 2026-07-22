// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserBalanceLog is the golang structure of table user_balance_log for DAO operations like Where/Data.
type UserBalanceLog struct {
	g.Meta        `orm:"table:user_balance_log, do:true"`
	Id            any         //
	UserId        any         //
	Direction     any         //
	Scene         any         //
	Amount        any         //
	BalanceBefore any         //
	BalanceAfter  any         //
	RefId         any         //
	Remark        any         //
	CreatedAt     *gtime.Time //
}
