// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RechargeOrderDao is the data access object for the table recharge_order.
type RechargeOrderDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  RechargeOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// RechargeOrderColumns defines and stores column names for the table recharge_order.
type RechargeOrderColumns struct {
	Id        string //
	SiteId    string //
	OrderNo   string //
	UserId    string //
	PackageId string //
	Amount    string //
	Coin      string //
	Status    string //
	PayAt     string //
	CreatedAt string //
}

// rechargeOrderColumns holds the columns for the table recharge_order.
var rechargeOrderColumns = RechargeOrderColumns{
	Id:        "id",
	SiteId:    "site_id",
	OrderNo:   "order_no",
	UserId:    "user_id",
	PackageId: "package_id",
	Amount:    "amount",
	Coin:      "coin",
	Status:    "status",
	PayAt:     "pay_at",
	CreatedAt: "created_at",
}

// NewRechargeOrderDao creates and returns a new DAO object for table data access.
func NewRechargeOrderDao(handlers ...gdb.ModelHandler) *RechargeOrderDao {
	return &RechargeOrderDao{
		group:    "default",
		table:    "recharge_order",
		columns:  rechargeOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RechargeOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RechargeOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RechargeOrderDao) Columns() RechargeOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RechargeOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RechargeOrderDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *RechargeOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
