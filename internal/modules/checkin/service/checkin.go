// Package service 签到对外接口。
package service

import "context"

type RewardGrant struct {
	Gold    int64
	VipDays int
}
type ClickDTO struct {
	Message          string
	TodayChecked     bool
	ContinuouslyDays int
	Rewards          []RewardGrant
}
type RewardCfgDTO struct {
	DayNum   int
	UserType int
	Gold     int64
	VipDays  int
}
type RecordDTO struct {
	Date             string
	ContinuouslyDays int
	RewardGold       int64
}
type InfoDTO struct {
	Rewards          []RewardCfgDTO
	TodayChecked     bool
	ContinuouslyDays int
	Records          []RecordDTO
}

type ICheckin interface {
	Click(ctx context.Context, userId int64) (*ClickDTO, error)
	Info(ctx context.Context, userId int64) (*InfoDTO, error)
}
