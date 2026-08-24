// Package logic 签到业务。
package logic

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/checkin/service"
)

const (
	checkinSiteId = 1
	checkinDays   = 15
)

type sCheckin struct{}

func New() service.ICheckin { return &sCheckin{} }

func effective(r *entity.CheckinReward) (gold, points int64, vipDays int) {
	if r.IsMilestone == 1 {
		return r.MsGold, r.MsPoints, r.MsVipDays
	}
	return r.Gold, r.Points, r.VipDays
}

func (s *sCheckin) Info(ctx context.Context, userId int64) (*service.InfoDTO, error) {
	var rewards []*entity.CheckinReward
	if err := g.Model("checkin_reward").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("status", 1).Where("user_type", 0).
		Order("day_num asc").Scan(&rewards); err != nil {
		return nil, err
	}
	var records []*entity.UserCheckin
	if err := g.Model("user_checkin").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("user_id", userId).
		Order("checkin_date desc").Limit(60).Scan(&records); err != nil {
		return nil, err
	}
	checked, streak := streakOf(records)
	out := &service.InfoDTO{TodayChecked: checked, ContinuouslyDays: streak}
	for _, r := range rewards {
		gold, points, vip := effective(r)
		label := r.Label
		if label == "" {
			label = fmt.Sprintf("第%d天", r.DayNum)
		}
		out.Rewards = append(out.Rewards, service.RewardCfgDTO{
			DayNum: r.DayNum, Label: label, UserType: r.UserType,
			Gold: gold, Points: points, VipDays: vip, IsMilestone: r.IsMilestone,
		})
	}
	for _, r := range records {
		out.Records = append(out.Records, service.RecordDTO{
			Date: dateStr(r.CheckinDate), ContinuouslyDays: r.ContinuouslyDays,
			RewardGold: r.RewardGold, RewardPoints: r.RewardPoints, RewardVipDays: r.RewardVipDays,
		})
	}
	return out, nil
}

func (s *sCheckin) Click(ctx context.Context, userId int64) (*service.ClickDTO, error) {
	today := gtime.Now().Format("Y-m-d")
	yesterday := gtime.Now().AddDate(0, 0, -1).Format("Y-m-d")

	var latest []*entity.UserCheckin
	if err := g.Model("user_checkin").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("user_id", userId).
		Order("checkin_date desc").Limit(1).Scan(&latest); err != nil {
		return nil, err
	}
	if len(latest) > 0 && dateStr(latest[0].CheckinDate) == today {
		return &service.ClickDTO{
			Message: "今日已经签到，请勿重复签到", TodayChecked: true,
			ContinuouslyDays: latest[0].ContinuouslyDays,
		}, nil
	}
	streak := 1
	if len(latest) > 0 && dateStr(latest[0].CheckinDate) == yesterday {
		streak = latest[0].ContinuouslyDays + 1
	}
	cycle := checkinDays
	day := (streak-1)%cycle + 1

	var rewards []*entity.CheckinReward
	if err := g.Model("checkin_reward").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("status", 1).
		Where("day_num", day).Where("user_type", 0).Scan(&rewards); err != nil {
		return nil, err
	}
	var totalGold, totalPoints int64
	var totalVip int
	grants := make([]service.RewardGrant, 0, len(rewards))
	for _, r := range rewards {
		gold, points, vip := effective(r)
		totalGold += gold
		totalPoints += points
		totalVip += vip
		grants = append(grants, service.RewardGrant{Gold: gold, Points: points, VipDays: vip})
	}
	cfg, _ := s.AdminConfig(ctx)
	vipGroupId := int64(0)
	if cfg != nil {
		vipGroupId = cfg.VipGroupId
	}

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("user_checkin").Ctx(ctx).Data(g.Map{
			"site_id": checkinSiteId, "user_id": userId, "checkin_date": today,
			"continuously_days": streak, "reward_gold": totalGold,
			"reward_points": totalPoints, "reward_vip_days": totalVip,
		}).Insert(); err != nil {
			return err
		}
		data := g.Map{}
		if totalGold > 0 {
			data["balance"] = &gdb.Counter{Field: "balance", Value: float64(totalGold)}
		}
		if totalPoints > 0 {
			data["credit"] = &gdb.Counter{Field: "credit", Value: float64(totalPoints)}
		}
		if len(data) > 0 {
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).Data(data).Update(); err != nil {
				return err
			}
		}
		if totalVip > 0 {
			if err := grantVipDays(ctx, tx, userId, totalVip, vipGroupId); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &service.ClickDTO{
		Message: "签到成功", TodayChecked: true, ContinuouslyDays: streak, Rewards: grants,
	}, nil
}

func grantVipDays(ctx context.Context, tx gdb.TX, userId int64, vipDays int, vipGroupId int64) error {
	type row struct {
		GroupEndTime int64  `orm:"group_end_time"`
		GroupId      int64  `orm:"group_id"`
	}
	var u row
	if err := tx.Model("users").Ctx(ctx).Where("id", userId).Scan(&u); err != nil {
		return err
	}
	now := gtime.Timestamp()
	base := now
	if u.GroupEndTime > base {
		base = u.GroupEndTime
	}
	data := g.Map{"group_end_time": base + int64(vipDays)*86400}
	if vipGroupId > 0 {
		var grp struct {
			Name string `orm:"name"`
		}
		_ = tx.Model("user_group").Ctx(ctx).Where("id", vipGroupId).Scan(&grp)
		data["group_id"] = vipGroupId
		if grp.Name != "" {
			data["group_name"] = grp.Name
		}
	}
	_, err := tx.Model("users").Ctx(ctx).Where("id", userId).Data(data).Update()
	return err
}

func (s *sCheckin) AdminConfig(ctx context.Context) (*service.ConfigDTO, error) {
	var list []*entity.CheckinConfig
	if err := g.Model("checkin_config").Ctx(ctx).Where("site_id", checkinSiteId).Scan(&list); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return &service.ConfigDTO{MakeupPoints: 5, MakeupLimit: 3}, nil
	}
	c := list[0]
	return &service.ConfigDTO{
		MakeupPoints: c.MakeupPoints, MakeupLimit: c.MakeupLimit,
		MakeupDesc: c.MakeupDesc, VipGroupId: c.VipGroupId,
	}, nil
}

func (s *sCheckin) SaveConfig(ctx context.Context, in service.ConfigDTO) error {
	if in.MakeupPoints < 0 {
		in.MakeupPoints = 0
	}
	if in.MakeupLimit < 0 {
		in.MakeupLimit = 0
	}
	if len(in.MakeupDesc) > 21 {
		in.MakeupDesc = in.MakeupDesc[:21]
	}
	cnt, err := g.Model("checkin_config").Ctx(ctx).Where("site_id", checkinSiteId).Count()
	if err != nil {
		return err
	}
	data := g.Map{
		"makeup_points": in.MakeupPoints, "makeup_limit": in.MakeupLimit,
		"makeup_desc": in.MakeupDesc, "vip_group_id": in.VipGroupId, "updated_at": gtime.Now(),
	}
	if cnt == 0 {
		data["site_id"] = checkinSiteId
		_, err = g.Model("checkin_config").Ctx(ctx).Data(data).Insert()
		return err
	}
	_, err = g.Model("checkin_config").Ctx(ctx).Where("site_id", checkinSiteId).Data(data).Update()
	return err
}

func (s *sCheckin) AdminRewards(ctx context.Context) ([]service.RewardRowDTO, error) {
	var list []*entity.CheckinReward
	if err := g.Model("checkin_reward").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("user_type", 0).
		Order("day_num asc").Scan(&list); err != nil {
		return nil, err
	}
	byDay := map[int]*entity.CheckinReward{}
	for _, r := range list {
		byDay[r.DayNum] = r
	}
	out := make([]service.RewardRowDTO, 0, checkinDays)
	for d := 1; d <= checkinDays; d++ {
		row := service.RewardRowDTO{DayNum: d, Label: fmt.Sprintf("第%d天", d)}
		if r, ok := byDay[d]; ok {
			row.Label = r.Label
			if row.Label == "" {
				row.Label = fmt.Sprintf("第%d天", d)
			}
			row.Points, row.Gold, row.VipDays = r.Points, r.Gold, r.VipDays
			row.IsMilestone, row.MsPoints, row.MsGold, row.MsVipDays = r.IsMilestone, r.MsPoints, r.MsGold, r.MsVipDays
		}
		out = append(out, row)
	}
	return out, nil
}

func (s *sCheckin) SaveRewards(ctx context.Context, rows []service.RewardRowDTO) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, row := range rows {
			if row.DayNum < 1 || row.DayNum > checkinDays {
				continue
			}
			if row.Label == "" {
				row.Label = fmt.Sprintf("第%d天", row.DayNum)
			}
			if row.IsMilestone != 1 {
				row.IsMilestone = 0
			}
			cnt, err := tx.Model("checkin_reward").Ctx(ctx).
				Where("site_id", checkinSiteId).Where("user_type", 0).Where("day_num", row.DayNum).Count()
			if err != nil {
				return err
			}
			data := g.Map{
				"label": row.Label, "points": row.Points, "gold": row.Gold, "vip_days": row.VipDays,
				"is_milestone": row.IsMilestone, "ms_points": row.MsPoints, "ms_gold": row.MsGold,
				"ms_vip_days": row.MsVipDays, "status": 1, "updated_at": gtime.Now(),
			}
			if cnt == 0 {
				data["site_id"] = checkinSiteId
				data["day_num"] = row.DayNum
				data["user_type"] = 0
				if _, err := tx.Model("checkin_reward").Ctx(ctx).Data(data).Insert(); err != nil {
					return err
				}
				continue
			}
			if _, err := tx.Model("checkin_reward").Ctx(ctx).
				Where("site_id", checkinSiteId).Where("user_type", 0).Where("day_num", row.DayNum).
				Data(data).Update(); err != nil {
				return err
			}
		}
		return nil
	})
}

func streakOf(records []*entity.UserCheckin) (checked bool, streak int) {
	if len(records) == 0 {
		return false, 0
	}
	today := gtime.Now().Format("Y-m-d")
	yesterday := gtime.Now().AddDate(0, 0, -1).Format("Y-m-d")
	d := dateStr(records[0].CheckinDate)
	if d == today {
		return true, records[0].ContinuouslyDays
	}
	if d == yesterday {
		return false, records[0].ContinuouslyDays
	}
	return false, 0
}

func dateStr(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("Y-m-d")
}
