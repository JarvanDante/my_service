// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserCodeLogDao is the data access object for the table user_code_log.
type UserCodeLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserCodeLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserCodeLogColumns defines and stores column names for the table user_code_log.
type UserCodeLogColumns struct {
	Id        string //
	CodeId    string //
	Code      string //
	CodeKey   string //
	Name      string //
	Type      string //
	ObjectId  string //
	UserId    string //
	Username  string //
	AddNum    string //
	CreatedAt string //
}

// userCodeLogColumns holds the columns for the table user_code_log.
var userCodeLogColumns = UserCodeLogColumns{
	Id:        "id",
	CodeId:    "code_id",
	Code:      "code",
	CodeKey:   "code_key",
	Name:      "name",
	Type:      "type",
	ObjectId:  "object_id",
	UserId:    "user_id",
	Username:  "username",
	AddNum:    "add_num",
	CreatedAt: "created_at",
}

// NewUserCodeLogDao creates and returns a new DAO object for table data access.
func NewUserCodeLogDao(handlers ...gdb.ModelHandler) *UserCodeLogDao {
	return &UserCodeLogDao{
		group:    "default",
		table:    "user_code_log",
		columns:  userCodeLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserCodeLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserCodeLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserCodeLogDao) Columns() UserCodeLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserCodeLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserCodeLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserCodeLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
