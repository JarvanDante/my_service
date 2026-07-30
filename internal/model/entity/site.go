// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Site is the golang structure for table site.
type Site struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Domain    string      `json:"domain"    orm:"domain"     description:""` //
	Appid     string      `json:"appid"     orm:"appid"      description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Config    string      `json:"config"    orm:"config"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
