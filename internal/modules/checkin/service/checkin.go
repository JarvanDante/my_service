package service

import "context"

type RewardGrant struct {
	Gold    int64
	Points  int64
	VipDays int
}
type ClickDTO struct {
	Message          string
	TodayChecked     bool
	ContinuouslyDays int
	Rewards          []RewardGrant
}
type RewardCfgDTO struct {
	DayNum      int
	Label       string
	UserType    int
	Gold        int64
	Points      int64
	VipDays     int
	IsMilestone int
}
type RecordDTO struct {
	Date             string
	ContinuouslyDays int
	RewardGold       int64
	RewardPoints     int64
	RewardVipDays    int
}
type InfoDTO struct {
	Rewards          []RewardCfgDTO
	TodayChecked     bool
	ContinuouslyDays int
	Records          []RecordDTO
}

type ConfigDTO struct {
	MakeupPoints int
	MakeupLimit  int
	MakeupDesc   string
	VipGroupId   int64
}

type RewardRowDTO struct {
	DayNum      int
	Label       string
	Points      int64
	Gold        int64
	VipDays     int
	IsMilestone int
	MsPoints    int64
	MsGold      int64
	MsVipDays   int
}

type ICheckin interface {
	Click(ctx context.Context, userId int64) (*ClickDTO, error)
	Info(ctx context.Context, userId int64) (*InfoDTO, error)
	AdminConfig(ctx context.Context) (*ConfigDTO, error)
	SaveConfig(ctx context.Context, in ConfigDTO) error
	AdminRewards(ctx context.Context) ([]RewardRowDTO, error)
	SaveRewards(ctx context.Context, rows []RewardRowDTO) error
}
