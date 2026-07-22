// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFollowDao is the data access object for the table user_follow.
type UserFollowDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserFollowColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserFollowColumns defines and stores column names for the table user_follow.
type UserFollowColumns struct {
	Id        string //
	UserId    string //
	HomeId    string //
	CreatedAt string //
	UpdatedAt string //
}

// userFollowColumns holds the columns for the table user_follow.
var userFollowColumns = UserFollowColumns{
	Id:        "id",
	UserId:    "user_id",
	HomeId:    "home_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserFollowDao creates and returns a new DAO object for table data access.
func NewUserFollowDao(handlers ...gdb.ModelHandler) *UserFollowDao {
	return &UserFollowDao{
		group:    "default",
		table:    "user_follow",
		columns:  userFollowColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserFollowDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserFollowDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserFollowDao) Columns() UserFollowColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserFollowDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserFollowDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserFollowDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
