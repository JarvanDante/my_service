// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SiteDao is the data access object for the table site.
type SiteDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SiteColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SiteColumns defines and stores column names for the table site.
type SiteColumns struct {
	Id        string //
	Name      string //
	Domain    string //
	Appid     string //
	Status    string //
	Config    string //
	CreatedAt string //
	UpdatedAt string //
}

// siteColumns holds the columns for the table site.
var siteColumns = SiteColumns{
	Id:        "id",
	Name:      "name",
	Domain:    "domain",
	Appid:     "appid",
	Status:    "status",
	Config:    "config",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSiteDao creates and returns a new DAO object for table data access.
func NewSiteDao(handlers ...gdb.ModelHandler) *SiteDao {
	return &SiteDao{
		group:    "default",
		table:    "site",
		columns:  siteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SiteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SiteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SiteDao) Columns() SiteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SiteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SiteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SiteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
