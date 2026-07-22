// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserSignDao is the data access object for the table user_sign.
type UserSignDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserSignColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserSignColumns defines and stores column names for the table user_sign.
type UserSignColumns struct {
	Id        string //
	UserId    string //
	YearMonth string //
	Days      string //
	Exchanges string //
	CreatedAt string //
	UpdatedAt string //
}

// userSignColumns holds the columns for the table user_sign.
var userSignColumns = UserSignColumns{
	Id:        "id",
	UserId:    "user_id",
	YearMonth: "year_month",
	Days:      "days",
	Exchanges: "exchanges",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserSignDao creates and returns a new DAO object for table data access.
func NewUserSignDao(handlers ...gdb.ModelHandler) *UserSignDao {
	return &UserSignDao{
		group:    "default",
		table:    "user_sign",
		columns:  userSignColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserSignDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserSignDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserSignDao) Columns() UserSignColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserSignDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserSignDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserSignDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
