package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
)

type userRepo struct{}

// NewUserRepo 返回 user 领域仓储实现。
func NewUserRepo() domain.Repository { return &userRepo{} }

func (r *userRepo) FindById(ctx context.Context, id int64) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Scan(&u)
	return u, err
}

func (r *userRepo) FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).Where(Users.Columns().DeviceId, deviceId).Scan(&u)
	return u, err
}

func (r *userRepo) Create(ctx context.Context, data g.Map) (int64, error) {
	return Users.Ctx(ctx).Data(data).InsertAndGetId()
}

func (r *userRepo) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(g.Map{
		Users.Columns().LoginNum:    &gdb.Counter{Field: Users.Columns().LoginNum, Value: 1},
		Users.Columns().LastLoginAt: gtime.Now(),
		Users.Columns().LastIp:      ip,
	}).Update()
	return err
}

func (r *userRepo) Disable(ctx context.Context, id int64, reason string) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(g.Map{
		Users.Columns().IsDisabled: 1,
		Users.Columns().ErrorMsg:   reason,
	}).Update()
	return err
}

func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).Where(Users.Columns().Phone, phone).Scan(&u)
	return u, err
}

func (r *userRepo) UpdatePhone(ctx context.Context, id int64, phone string) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(g.Map{
		Users.Columns().Phone: phone,
	}).Update()
	return err
}

func (r *userRepo) FindByAccount(ctx context.Context, account string) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).
		Where(Users.Columns().Username, account).
		WhereOr(Users.Columns().Phone, account).
		Scan(&u)
	return u, err
}

func (r *userRepo) UpdateProfile(ctx context.Context, id int64, data g.Map) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(data).Update()
	return err
}

func (r *userRepo) ExistsFollow(ctx context.Context, userId, homeId int64) (bool, error) {
	n, err := g.Model("user_follow").Ctx(ctx).
		Where("user_id", userId).
		Where("home_id", homeId).
		Count()
	return n > 0, err
}

// Follow 关注: 事务内写关注关系 + 双方计数。
func (r *userRepo) Follow(ctx context.Context, userId, homeId int64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("user_follow").Ctx(ctx).Data(g.Map{
			"user_id": userId,
			"home_id": homeId,
		}).Insert(); err != nil {
			return err
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
			Data(g.Map{"follow": &gdb.Counter{Field: "follow", Value: 1}}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", homeId).
			Data(g.Map{"fans": &gdb.Counter{Field: "fans", Value: 1}}).Update(); err != nil {
			return err
		}
		return nil
	})
}

// Unfollow 取关: 事务内删关注关系 + 双方计数。
func (r *userRepo) Unfollow(ctx context.Context, userId, homeId int64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("user_follow").Ctx(ctx).
			Where("user_id", userId).Where("home_id", homeId).Delete()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return nil // 本来就没关注, 不改计数
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
			Data(g.Map{"follow": &gdb.Counter{Field: "follow", Value: -1}}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", homeId).
			Data(g.Map{"fans": &gdb.Counter{Field: "fans", Value: -1}}).Update(); err != nil {
			return err
		}
		return nil
	})
}

// FollowingList 我关注的人。
func (r *userRepo) FollowingList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error) {
	m := g.Model("user_follow f").Ctx(ctx).
		LeftJoin("users u", "u.id=f.home_id").
		Where("f.user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Users
	err = m.Clone().Fields("u.*").Page(page, size).OrderDesc("f.id").Scan(&list)
	return list, total, err
}

// FansList 关注我的人。
func (r *userRepo) FansList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error) {
	m := g.Model("user_follow f").Ctx(ctx).
		LeftJoin("users u", "u.id=f.user_id").
		Where("f.home_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Users
	err = m.Clone().Fields("u.*").Page(page, size).OrderDesc("f.id").Scan(&list)
	return list, total, err
}

// BindInviter 绑定推荐人: 事务内设置 parent + 推荐人 share_num+1。
func (r *userRepo) BindInviter(ctx context.Context, userId, inviterId int64, inviterName string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).Data(g.Map{
			"parent_id":   inviterId,
			"parent_name": inviterName,
		}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", inviterId).Data(g.Map{
			"share_num": &gdb.Counter{Field: "share_num", Value: 1},
		}).Update(); err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepo) FindCodeByCode(ctx context.Context, code string) (*entity.UserCode, error) {
	var c *entity.UserCode
	err := g.Model("user_code").Ctx(ctx).Where("code", code).Scan(&c)
	return c, err
}

func (r *userRepo) HasRedeemed(ctx context.Context, codeId, userId int64) (bool, error) {
	n, err := g.Model("user_code_log").Ctx(ctx).
		Where("code_id", codeId).Where("user_id", userId).Count()
	return n > 0, err
}

// RedeemCode 使用兑换码: 事务内发放(金币/用户组) + 改用量 + 写记录。
func (r *userRepo) RedeemCode(ctx context.Context, userId int64, username string, c *entity.UserCode) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		switch c.Type {
		case "point": // 加金币(balance)
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
				Data(g.Map{"balance": &gdb.Counter{Field: "balance", Value: float64(c.AddNum)}}).Update(); err != nil {
				return err
			}
		case "group": // 加/续用户组, add_num 为天数
			now := gtime.Timestamp()
			base := now
			v, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("group_end_time").Value()
			if err != nil {
				return err
			}
			if cur := v.Int64(); cur > base {
				base = cur
			}
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).Data(g.Map{
				"group_id":       c.ObjectId,
				"group_end_time": base + int64(c.AddNum)*86400,
			}).Update(); err != nil {
				return err
			}
		}
		// 用量 +1, 用完置为已使用
		status := c.Status
		if c.UsedNum+1 >= c.CanUseNum {
			status = 1
		}
		if _, err := tx.Model("user_code").Ctx(ctx).Where("id", c.Id).Data(g.Map{
			"used_num": &gdb.Counter{Field: "used_num", Value: 1},
			"status":   status,
		}).Update(); err != nil {
			return err
		}
		// 记录
		if _, err := tx.Model("user_code_log").Ctx(ctx).Data(g.Map{
			"code_id":   c.Id,
			"code":      c.Code,
			"code_key":  c.CodeKey,
			"name":      c.Name,
			"type":      c.Type,
			"object_id": c.ObjectId,
			"user_id":   userId,
			"username":  username,
			"add_num":   c.AddNum,
		}).Insert(); err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepo) CodeLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserCodeLog, int, error) {
	m := g.Model("user_code_log").Ctx(ctx).Where("user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCodeLog
	err = m.Clone().Page(page, size).OrderDesc("id").Scan(&list)
	return list, total, err
}

func (r *userRepo) AddShareLog(ctx context.Context, userId int64, typ string, targetId int64, channel string) error {
	_, err := g.Model("user_share_log").Ctx(ctx).Data(g.Map{
		"user_id":   userId,
		"type":      typ,
		"target_id": targetId,
		"channel":   channel,
	}).Insert()
	return err
}

func (r *userRepo) ShareLogList(ctx context.Context, userId int64, page, size int) ([]*entity.UserShareLog, int, error) {
	m := g.Model("user_share_log").Ctx(ctx).Where("user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserShareLog
	err = m.Clone().Page(page, size).OrderDesc("id").Scan(&list)
	return list, total, err
}
