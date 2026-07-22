// Package entity 数据库实体, 与表一一对应。正式项目由 `gf gen dao` 生成。
package entity

type User struct {
	Id     int64  `orm:"id"     json:"id"`
	Name   string `orm:"name"   json:"name"`
	Email  string `orm:"email"  json:"email"`
	Status int    `orm:"status" json:"status"`
}
