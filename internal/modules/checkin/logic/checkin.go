// Package logic 签到业务(移植自 tianbi checkinser)。
// 与 Mongo 版差异: 用 PG 表 + gf ORM; 连续天数直接存快照(不逐条回溯); 事务发金币到 users.balance。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/checkin/service"
)

const (
	checkinSiteId = 1  // 单站点样板; 多租户改为从上下文取 site_id
	checkinCycle  = 30 // 30 天周期(与 tianbi maxCheckinDay 一致)
)

type sCheckin struct{}

func New() service.ICheckin { return &sCheckin{} }

// Info 返回阶梯奖励配置 + 当前连续天数 + 用户签到记录。
func (s *sCheckin) Info(ctx context.Context, userId int64) (*service.InfoDTO, error) {
	var rewards []*entity.CheckinReward
	if err := g.Model("checkin_reward").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("status", 1).
		Order("day_num asc").Scan(&rewards); err != nil {
		return nil, err
	}
	var records []*entity.UserCheckin
	if err := g.Model("user_checkin").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("user_id", userId).
		Order("checkin_date desc").Limit(checkinCycle + 1).Scan(&records); err != nil {
		return nil, err
	}
	checked, streak := streakOf(records)
	out := &service.InfoDTO{TodayChecked: checked, ContinuouslyDays: streak}
	for _, r := range rewards {
		out.Rewards = append(out.Rewards, service.RewardCfgDTO{
			DayNum: r.DayNum, UserType: r.UserType, Gold: r.Gold, VipDays: r.VipDays,
		})
	}
	for _, r := range records {
		out.Records = append(out.Records, service.RecordDTO{
			Date: dateStr(r.CheckinDate), ContinuouslyDays: r.ContinuouslyDays, RewardGold: r.RewardGold,
		})
	}
	return out, nil
}

// Click 签到: 计算连续天数 -> 匹配阶梯奖励 -> 事务插记录+发金币。
func (s *sCheckin) Click(ctx context.Context, userId int64) (*service.ClickDTO, error) {
	today := gtime.Now().Format("Y-m-d")
	yesterday := gtime.Now().AddDate(0, 0, -1).Format("Y-m-d")

	var latest []*entity.UserCheckin
	if err := g.Model("user_checkin").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("user_id", userId).
		Order("checkin_date desc").Limit(1).Scan(&latest); err != nil {
		return nil, err
	}
	// 今日已签
	if len(latest) > 0 && dateStr(latest[0].CheckinDate) == today {
		return &service.ClickDTO{
			Message: "今日已经签到，请勿重复签到", TodayChecked: true,
			ContinuouslyDays: latest[0].ContinuouslyDays,
		}, nil
	}
	// 连续天数: 昨天签过则 +1, 否则重置为 1
	streak := 1
	if len(latest) > 0 && dateStr(latest[0].CheckinDate) == yesterday {
		streak = latest[0].ContinuouslyDays + 1
	}
	// 周期内第几天(1~30)
	day := (streak-1)%checkinCycle + 1

	// 匹配奖励(user_type=0 通用; VIP 档 user_type=1 需 VIP 系统, 样板略)
	var rewards []*entity.CheckinReward
	if err := g.Model("checkin_reward").Ctx(ctx).
		Where("site_id", checkinSiteId).Where("status", 1).
		Where("day_num", day).Where("user_type", 0).Scan(&rewards); err != nil {
		return nil, err
	}
	var totalGold int64
	grants := make([]service.RewardGrant, 0, len(rewards))
	for _, r := range rewards {
		totalGold += r.Gold
		grants = append(grants, service.RewardGrant{Gold: r.Gold, VipDays: r.VipDays})
	}

	// 事务: 插签到记录 + 发金币(余额)。unique(site,user,date) 兜底防并发重复。
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("user_checkin").Ctx(ctx).Data(g.Map{
			"site_id": checkinSiteId, "user_id": userId, "checkin_date": today,
			"continuously_days": streak, "reward_gold": totalGold,
		}).Insert(); err != nil {
			return err
		}
		if totalGold > 0 {
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
				Data(g.Map{"balance": &gdb.Counter{Field: "balance", Value: float64(totalGold)}}).Update(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &service.ClickDTO{
		Message: "success", TodayChecked: true, ContinuouslyDays: streak, Rewards: grants,
	}, nil
}

// streakOf 从倒序记录推断(今日是否已签, 当前连续天数)。
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
		return false, records[0].ContinuouslyDays // 昨天签过, 今日可续
	}
	return false, 0 // 断签
}

func dateStr(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("Y-m-d")
}
