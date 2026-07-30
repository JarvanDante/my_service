// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type ChatConversation struct {
	Id          int64       `json:"id"          orm:"id"`
	UserId      int64       `json:"userId"      orm:"user_id"`
	PeerId      int64       `json:"peerId"      orm:"peer_id"`
	LastContent string      `json:"lastContent" orm:"last_content"`
	LastAt      *gtime.Time `json:"lastAt"      orm:"last_at"`
	Unread      int         `json:"unread"      orm:"unread"`
	Deleted     int         `json:"deleted"     orm:"deleted"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}
