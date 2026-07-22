// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserCodeDao is the data access object for the table user_code.
type UserCodeDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserCodeColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserCodeColumns defines and stores column names for the table user_code.
type UserCodeColumns struct {
	Id        string //
	Name      string //
	Code      string //
	CodeKey   string //
	Type      string //
	ObjectId  string //
	AddNum    string //
	CanUseNum string //
	UsedNum   string //
	Status    string //
	ExpiredAt string //
	CreatedAt string //
	UpdatedAt string //
}

// userCodeColumns holds the columns for the table user_code.
var userCodeColumns = UserCodeColumns{
	Id:        "id",
	Name:      "name",
	Code:      "code",
	CodeKey:   "code_key",
	Type:      "type",
	ObjectId:  "object_id",
	AddNum:    "add_num",
	CanUseNum: "can_use_num",
	UsedNum:   "used_num",
	Status:    "status",
	ExpiredAt: "expired_at",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserCodeDao creates and returns a new DAO object for table data access.
func NewUserCodeDao(handlers ...gdb.ModelHandler) *UserCodeDao {
	return &UserCodeDao{
		group:    "default",
		table:    "user_code",
		columns:  userCodeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserCodeDao) Columns() UserCodeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserCodeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
