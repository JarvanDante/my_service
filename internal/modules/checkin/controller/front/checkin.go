// Package front 前台签到控制器(薄适配层)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/checkin/v1"
	"github.com/JarvanDante/my_service/internal/modules/checkin/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.ICheckin }

func New(svc service.ICheckin) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

func (c *Controller) Click(ctx context.Context, req *v1.ClickReq) (res *v1.ClickRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Click(ctx, id)
	if err != nil {
		return nil, err
	}
	res = &v1.ClickRes{Message: d.Message, TodayChecked: d.TodayChecked, ContinuouslyDays: d.ContinuouslyDays}
	for _, x := range d.Rewards {
		res.Rewards = append(res.Rewards, v1.RewardItem{Gold: x.Gold, VipDays: x.VipDays})
	}
	return res, nil
}

func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Info(ctx, id)
	if err != nil {
		return nil, err
	}
	res = &v1.InfoRes{TodayChecked: d.TodayChecked, ContinuouslyDays: d.ContinuouslyDays}
	for _, x := range d.Rewards {
		res.Rewards = append(res.Rewards, v1.RewardCfg{DayNum: x.DayNum, UserType: x.UserType, Gold: x.Gold, VipDays: x.VipDays})
	}
	for _, x := range d.Records {
		res.Records = append(res.Records, v1.RecordItem{Date: x.Date, ContinuouslyDays: x.ContinuouslyDays, RewardGold: x.RewardGold})
	}
	return res, nil
}
