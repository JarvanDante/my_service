// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type ChatMessage struct {
	Id        int64       `json:"id"        orm:"id"`
	FromId    int64       `json:"fromId"    orm:"from_id"`
	ToId      int64       `json:"toId"      orm:"to_id"`
	Content   string      `json:"content"   orm:"content"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
