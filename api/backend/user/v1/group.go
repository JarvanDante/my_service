// Package v1 后台会员等级(用户组)接口契约(B4)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type UserGroupItem struct {
	Id                int64   `json:"id"`
	Name              string  `json:"name"`
	TitleHeat         string  `json:"title_heat"`
	TitleDescription  string  `json:"title_description"`
	TitlePicture      string  `json:"title_picture"`
	Img               string  `json:"img"`
	Level             int     `json:"level"`
	LevelText         string  `json:"level_text"`
	PromotionType     int     `json:"promotion_type"`
	PromotionTypeText string  `json:"promotion_type_text"`
	Price             float64 `json:"price"`
	OldPrice          float64 `json:"old_price"`
	Rate              int     `json:"rate"`
	DayNum            int     `json:"day_num"`
	GiftNum           int     `json:"gift_num"`
	DownloadNum       int     `json:"download_num"`
	DayTips           string  `json:"day_tips"`
	PriceTips         string  `json:"price_tips"`
	Rights            string  `json:"rights"`
	Remark            string  `json:"remark"`
	Sort              int     `json:"sort"`
	Status            int     `json:"status"`
	IsDisabledText    string  `json:"is_disabled_text"`
	UpdatedAt         string  `json:"updated_at"`
}

// 会员等级列表
type GroupListReq struct {
	g.Meta `path:"/user-groups" method:"get" tags:"Backend/UserGroup" summary:"会员等级列表"`
	Name   string `json:"name" in:"query"`
}
type GroupListRes struct {
	List []UserGroupItem `json:"list"`
}

// 创建会员等级
type GroupCreateReq struct {
	g.Meta           `path:"/user-groups" method:"post" tags:"Backend/UserGroup" summary:"创建会员等级"`
	Name             string  `json:"name"              v:"required#名称必填"`
	TitleHeat        string  `json:"title_heat"`
	TitleDescription string  `json:"title_description"`
	TitlePicture     string  `json:"title_picture"`
	Img              string  `json:"img"`
	Level            int     `json:"level"`
	PromotionType    int     `json:"promotion_type"`
	Price            float64 `json:"price"`
	OldPrice         float64 `json:"old_price"`
	Rate             int     `json:"rate"              v:"between:-2,100#折扣须在-2~100"`
	DayNum           int     `json:"day_num"           v:"min:1#可用天数必须大于0"`
	GiftNum          int     `json:"gift_num"`
	DownloadNum      int     `json:"download_num"`
	DayTips          string  `json:"day_tips"`
	PriceTips        string  `json:"price_tips"`
	Rights           string  `json:"rights"`
	Remark           string  `json:"remark"`
	Sort             int     `json:"sort"`
	Status           int     `json:"status"            v:"in:0,1#status 仅支持 0/1"`
}
type GroupCreateRes struct {
	Id int64 `json:"id"`
}

// 更新会员等级(同步组内用户快照)
type GroupUpdateReq struct {
	g.Meta           `path:"/user-groups/{id}" method:"put" tags:"Backend/UserGroup" summary:"更新会员等级"`
	Id               int64   `json:"id"                v:"required|min:1#组ID必填|组ID必须大于0"`
	Name             string  `json:"name"              v:"required#名称必填"`
	TitleHeat        string  `json:"title_heat"`
	TitleDescription string  `json:"title_description"`
	TitlePicture     string  `json:"title_picture"`
	Img              string  `json:"img"`
	Level            int     `json:"level"`
	PromotionType    int     `json:"promotion_type"`
	Price            float64 `json:"price"`
	OldPrice         float64 `json:"old_price"`
	Rate             int     `json:"rate"              v:"between:-2,100#折扣须在-2~100"`
	DayNum           int     `json:"day_num"           v:"min:1#可用天数必须大于0"`
	GiftNum          int     `json:"gift_num"`
	DownloadNum      int     `json:"download_num"`
	DayTips          string  `json:"day_tips"`
	PriceTips        string  `json:"price_tips"`
	Rights           string  `json:"rights"`
	Remark           string  `json:"remark"`
	Sort             int     `json:"sort"`
	Status           int     `json:"status"            v:"in:0,1#status 仅支持 0/1"`
}
type GroupUpdateRes struct{}

// 删除会员等级
type GroupDeleteReq struct {
	g.Meta `path:"/user-groups/{id}" method:"delete" tags:"Backend/UserGroup" summary:"删除会员等级"`
	Id     int64 `json:"id" v:"required|min:1#组ID必填|组ID必须大于0"`
}
type GroupDeleteRes struct{}
