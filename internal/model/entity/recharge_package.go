// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type RechargePackage struct {
	Id        int64       `json:"id"        orm:"id"`
	Name      string      `json:"name"      orm:"name"`
	Amount    float64     `json:"amount"    orm:"amount"`
	Coin      float64     `json:"coin"      orm:"coin"`
	Bonus     float64     `json:"bonus"     orm:"bonus"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
