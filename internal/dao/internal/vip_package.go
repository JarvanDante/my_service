// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VipPackageDao is the data access object for the table vip_package.
type VipPackageDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  VipPackageColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// VipPackageColumns defines and stores column names for the table vip_package.
type VipPackageColumns struct {
	Id        string //
	SiteId    string //
	Name      string //
	Days      string //
	Price     string //
	GroupId   string //
	Sort      string //
	Status    string //
	CreatedAt string //
}

// vipPackageColumns holds the columns for the table vip_package.
var vipPackageColumns = VipPackageColumns{
	Id:        "id",
	SiteId:    "site_id",
	Name:      "name",
	Days:      "days",
	Price:     "price",
	GroupId:   "group_id",
	Sort:      "sort",
	Status:    "status",
	CreatedAt: "created_at",
}

// NewVipPackageDao creates and returns a new DAO object for table data access.
func NewVipPackageDao(handlers ...gdb.ModelHandler) *VipPackageDao {
	return &VipPackageDao{
		group:    "default",
		table:    "vip_package",
		columns:  vipPackageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *VipPackageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *VipPackageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *VipPackageDao) Columns() VipPackageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *VipPackageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *VipPackageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *VipPackageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
