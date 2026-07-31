// Package v1 后台财务接口契约(B2): 充值/VIP 套餐 CRUD。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 充值套餐 ----------

type RechargePackageItem struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Coin   float64 `json:"coin"`
	Bonus  float64 `json:"bonus"`
	Sort   int     `json:"sort"`
	Status int     `json:"status"`
}

type RechargePkgListReq struct {
	g.Meta `path:"/recharge-packages" method:"get" tags:"Backend/Finance" summary:"充值套餐列表"`
}
type RechargePkgListRes struct {
	List []RechargePackageItem `json:"list"`
}

type RechargePkgCreateReq struct {
	g.Meta `path:"/recharge-packages" method:"post" tags:"Backend/Finance" summary:"创建充值套餐"`
	Name   string  `json:"name"   v:"required#套餐名必填"`
	Amount float64 `json:"amount" v:"required|min:0.01#价格必填|价格必须大于0"`
	Coin   float64 `json:"coin"   v:"required|min:0.01#到账金币必填|到账金币必须大于0"`
	Bonus  float64 `json:"bonus"  v:"min:0#赠送不合法"`
	Sort   int     `json:"sort"`
	Status int     `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type RechargePkgCreateRes struct {
	Id int64 `json:"id"`
}

type RechargePkgUpdateReq struct {
	g.Meta `path:"/recharge-packages/{id}" method:"put" tags:"Backend/Finance" summary:"更新充值套餐"`
	Id     int64   `json:"id"     v:"required|min:1#套餐ID必填|套餐ID必须大于0"`
	Name   string  `json:"name"   v:"required#套餐名必填"`
	Amount float64 `json:"amount" v:"required|min:0.01#价格必填|价格必须大于0"`
	Coin   float64 `json:"coin"   v:"required|min:0.01#到账金币必填|到账金币必须大于0"`
	Bonus  float64 `json:"bonus"  v:"min:0#赠送不合法"`
	Sort   int     `json:"sort"`
	Status int     `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type RechargePkgUpdateRes struct{}

type RechargePkgDeleteReq struct {
	g.Meta `path:"/recharge-packages/{id}" method:"delete" tags:"Backend/Finance" summary:"删除充值套餐"`
	Id     int64 `json:"id" v:"required|min:1#套餐ID必填|套餐ID必须大于0"`
}
type RechargePkgDeleteRes struct{}

// ---------- VIP 套餐 ----------

type VipPackageItem struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Days    int     `json:"days"`
	Price   float64 `json:"price"`
	GroupId int64   `json:"group_id"`
	Sort    int     `json:"sort"`
	Status  int     `json:"status"`
}

type VipPkgListReq struct {
	g.Meta `path:"/vip-packages" method:"get" tags:"Backend/Finance" summary:"VIP套餐列表"`
}
type VipPkgListRes struct {
	List []VipPackageItem `json:"list"`
}

type VipPkgCreateReq struct {
	g.Meta  `path:"/vip-packages" method:"post" tags:"Backend/Finance" summary:"创建VIP套餐"`
	Name    string  `json:"name"  v:"required#套餐名必填"`
	Days    int     `json:"days"  v:"required|min:1#时长必填|时长必须大于0"`
	Price   float64 `json:"price" v:"required|min:0.01#价格必填|价格必须大于0"`
	GroupId int64   `json:"group_id" v:"min:0#用户组不合法"`
	Sort    int     `json:"sort"`
	Status  int     `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type VipPkgCreateRes struct {
	Id int64 `json:"id"`
}

type VipPkgUpdateReq struct {
	g.Meta  `path:"/vip-packages/{id}" method:"put" tags:"Backend/Finance" summary:"更新VIP套餐"`
	Id      int64   `json:"id"    v:"required|min:1#套餐ID必填|套餐ID必须大于0"`
	Name    string  `json:"name"  v:"required#套餐名必填"`
	Days    int     `json:"days"  v:"required|min:1#时长必填|时长必须大于0"`
	Price   float64 `json:"price" v:"required|min:0.01#价格必填|价格必须大于0"`
	GroupId int64   `json:"group_id" v:"min:0#用户组不合法"`
	Sort    int     `json:"sort"`
	Status  int     `json:"status" v:"in:0,1#status 仅支持 0/1"`
}
type VipPkgUpdateRes struct{}

type VipPkgDeleteReq struct {
	g.Meta `path:"/vip-packages/{id}" method:"delete" tags:"Backend/Finance" summary:"删除VIP套餐"`
	Id     int64 `json:"id" v:"required|min:1#套餐ID必填|套餐ID必须大于0"`
}
type VipPkgDeleteRes struct{}
