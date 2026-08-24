package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/checkin/v1"
	"github.com/JarvanDante/my_service/internal/modules/checkin/service"
)

type Controller struct{ svc service.ICheckin }

func New(svc service.ICheckin) *Controller { return &Controller{svc: svc} }

func (c *Controller) GetConfig(ctx context.Context, _ *v1.GetConfigReq) (res *v1.GetConfigRes, err error) {
	d, err := c.svc.AdminConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetConfigRes{Config: v1.Config{
		MakeupPoints: d.MakeupPoints, MakeupLimit: d.MakeupLimit,
		MakeupDesc: d.MakeupDesc, VipGroupId: d.VipGroupId,
	}}, nil
}

func (c *Controller) SaveConfig(ctx context.Context, req *v1.SaveConfigReq) (res *v1.SaveConfigRes, err error) {
	if err = c.svc.SaveConfig(ctx, service.ConfigDTO{
		MakeupPoints: req.MakeupPoints, MakeupLimit: req.MakeupLimit,
		MakeupDesc: req.MakeupDesc, VipGroupId: req.VipGroupId,
	}); err != nil {
		return nil, err
	}
	return &v1.SaveConfigRes{}, nil
}

func (c *Controller) RewardList(ctx context.Context, _ *v1.RewardListReq) (res *v1.RewardListRes, err error) {
	list, err := c.svc.AdminRewards(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.RewardListRes{List: make([]v1.RewardRow, 0, len(list))}
	for _, r := range list {
		res.List = append(res.List, v1.RewardRow{
			DayNum: r.DayNum, Label: r.Label, Points: r.Points, Gold: r.Gold, VipDays: r.VipDays,
			IsMilestone: r.IsMilestone, MsPoints: r.MsPoints, MsGold: r.MsGold, MsVipDays: r.MsVipDays,
		})
	}
	return res, nil
}

func (c *Controller) SaveRewards(ctx context.Context, req *v1.SaveRewardsReq) (res *v1.SaveRewardsRes, err error) {
	rows := make([]service.RewardRowDTO, 0, len(req.List))
	for _, r := range req.List {
		rows = append(rows, service.RewardRowDTO{
			DayNum: r.DayNum, Label: r.Label, Points: r.Points, Gold: r.Gold, VipDays: r.VipDays,
			IsMilestone: r.IsMilestone, MsPoints: r.MsPoints, MsGold: r.MsGold, MsVipDays: r.MsVipDays,
		})
	}
	if err = c.svc.SaveRewards(ctx, rows); err != nil {
		return nil, err
	}
	return &v1.SaveRewardsRes{}, nil
}
