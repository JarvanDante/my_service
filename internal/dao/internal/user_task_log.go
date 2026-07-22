// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTaskLogDao is the data access object for the table user_task_log.
type UserTaskLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserTaskLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserTaskLogColumns defines and stores column names for the table user_task_log.
type UserTaskLogColumns struct {
	Id        string //
	UserId    string //
	TaskId    string //
	Type      string //
	Num       string //
	LogDate   string //
	CreatedAt string //
}

// userTaskLogColumns holds the columns for the table user_task_log.
var userTaskLogColumns = UserTaskLogColumns{
	Id:        "id",
	UserId:    "user_id",
	TaskId:    "task_id",
	Type:      "type",
	Num:       "num",
	LogDate:   "log_date",
	CreatedAt: "created_at",
}

// NewUserTaskLogDao creates and returns a new DAO object for table data access.
func NewUserTaskLogDao(handlers ...gdb.ModelHandler) *UserTaskLogDao {
	return &UserTaskLogDao{
		group:    "default",
		table:    "user_task_log",
		columns:  userTaskLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserTaskLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserTaskLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserTaskLogDao) Columns() UserTaskLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserTaskLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserTaskLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserTaskLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
