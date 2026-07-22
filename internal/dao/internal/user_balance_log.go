// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserBalanceLogDao is the data access object for the table user_balance_log.
type UserBalanceLogDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  UserBalanceLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// UserBalanceLogColumns defines and stores column names for the table user_balance_log.
type UserBalanceLogColumns struct {
	Id            string //
	UserId        string //
	Direction     string //
	Scene         string //
	Amount        string //
	BalanceBefore string //
	BalanceAfter  string //
	RefId         string //
	Remark        string //
	CreatedAt     string //
}

// userBalanceLogColumns holds the columns for the table user_balance_log.
var userBalanceLogColumns = UserBalanceLogColumns{
	Id:            "id",
	UserId:        "user_id",
	Direction:     "direction",
	Scene:         "scene",
	Amount:        "amount",
	BalanceBefore: "balance_before",
	BalanceAfter:  "balance_after",
	RefId:         "ref_id",
	Remark:        "remark",
	CreatedAt:     "created_at",
}

// NewUserBalanceLogDao creates and returns a new DAO object for table data access.
func NewUserBalanceLogDao(handlers ...gdb.ModelHandler) *UserBalanceLogDao {
	return &UserBalanceLogDao{
		group:    "default",
		table:    "user_balance_log",
		columns:  userBalanceLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserBalanceLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserBalanceLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserBalanceLogDao) Columns() UserBalanceLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserBalanceLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserBalanceLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserBalanceLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
