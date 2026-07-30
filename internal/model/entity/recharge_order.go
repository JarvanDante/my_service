// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
package entity

import "github.com/gogf/gf/v2/os/gtime"

type RechargeOrder struct {
	Id        int64       `json:"id"        orm:"id"`
	OrderNo   string      `json:"orderNo"   orm:"order_no"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	PackageId int64       `json:"packageId" orm:"package_id"`
	Amount    float64     `json:"amount"    orm:"amount"`
	Coin      float64     `json:"coin"      orm:"coin"`
	Status    int         `json:"status"    orm:"status"`
	PayAt     *gtime.Time `json:"payAt"     orm:"pay_at"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
