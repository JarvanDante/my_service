// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersDao is the data access object for the table users.
type UsersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UsersColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UsersColumns defines and stores column names for the table users.
type UsersColumns struct {
	Id             string //
	Username       string //
	Nickname       string //
	Phone          string //
	CountryCode    string //
	Password       string //
	Slat           string //
	DeviceId       string //
	DeviceType     string //
	DeviceVersion  string //
	DeviceExt      string //
	Img            string //
	BgImg          string //
	Signature      string //
	Sex            string //
	Level          string //
	Tag            string //
	Balance        string //
	Credit         string //
	GiftCount      string //
	MoneyCount     string //
	SendCount      string //
	HasBuy         string //
	Rights         string //
	WithdrawInfo   string //
	GroupId        string //
	GroupName      string //
	GroupRate      string //
	GroupStartTime string //
	GroupEndTime   string //
	ParentId       string //
	ParentName     string //
	ChannelName    string //
	ShareNum       string //
	Fans           string //
	Follow         string //
	Country        string //
	Province       string //
	City           string //
	Location       string //
	RegisterArea   string //
	IsChina        string //
	IsDisabled     string //
	ErrorMsg       string //
	IsSystem       string //
	RegisterAt     string //
	RegisterDate   string //
	RegisterIp     string //
	LoginNum       string //
	LastLoginAt    string //
	LastDate       string //
	LastIp         string //
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// usersColumns holds the columns for the table users.
var usersColumns = UsersColumns{
	Id:             "id",
	Username:       "username",
	Nickname:       "nickname",
	Phone:          "phone",
	CountryCode:    "country_code",
	Password:       "password",
	Slat:           "slat",
	DeviceId:       "device_id",
	DeviceType:     "device_type",
	DeviceVersion:  "device_version",
	DeviceExt:      "device_ext",
	Img:            "img",
	BgImg:          "bg_img",
	Signature:      "signature",
	Sex:            "sex",
	Level:          "level",
	Tag:            "tag",
	Balance:        "balance",
	Credit:         "credit",
	GiftCount:      "gift_count",
	MoneyCount:     "money_count",
	SendCount:      "send_count",
	HasBuy:         "has_buy",
	Rights:         "rights",
	WithdrawInfo:   "withdraw_info",
	GroupId:        "group_id",
	GroupName:      "group_name",
	GroupRate:      "group_rate",
	GroupStartTime: "group_start_time",
	GroupEndTime:   "group_end_time",
	ParentId:       "parent_id",
	ParentName:     "parent_name",
	ChannelName:    "channel_name",
	ShareNum:       "share_num",
	Fans:           "fans",
	Follow:         "follow",
	Country:        "country",
	Province:       "province",
	City:           "city",
	Location:       "location",
	RegisterArea:   "register_area",
	IsChina:        "is_china",
	IsDisabled:     "is_disabled",
	ErrorMsg:       "error_msg",
	IsSystem:       "is_system",
	RegisterAt:     "register_at",
	RegisterDate:   "register_date",
	RegisterIp:     "register_ip",
	LoginNum:       "login_num",
	LastLoginAt:    "last_login_at",
	LastDate:       "last_date",
	LastIp:         "last_ip",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewUsersDao creates and returns a new DAO object for table data access.
func NewUsersDao(handlers ...gdb.ModelHandler) *UsersDao {
	return &UsersDao{
		group:    "default",
		table:    "users",
		columns:  usersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UsersDao) Columns() UsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UsersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
