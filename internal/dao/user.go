package dao

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
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

// GetSignDays 取某月已签到日, 及该月记录是否存在。
func (r *userRepo) GetSignDays(ctx context.Context, userId int64, yearMonth int) ([]int, bool, error) {
	one, err := g.Model("user_sign").Ctx(ctx).
		Where("user_id", userId).Where("year_month", yearMonth).One()
	if err != nil {
		return nil, false, err
	}
	if one.IsEmpty() {
		return []int{}, false, nil
	}
	return one["days"].Ints(), true, nil
}

// SaveSign 事务内: upsert 签到记录 + 发放签到积分。
func (r *userRepo) SaveSign(ctx context.Context, userId int64, yearMonth int, days []int, exists bool, credit float64) error {
	b, _ := json.Marshal(days)
	daysJson := string(b)
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if exists {
			if _, err := tx.Model("user_sign").Ctx(ctx).
				Where("user_id", userId).Where("year_month", yearMonth).
				Data(g.Map{"days": daysJson}).Update(); err != nil {
				return err
			}
		} else {
			if _, err := tx.Model("user_sign").Ctx(ctx).Data(g.Map{
				"user_id":    userId,
				"year_month": yearMonth,
				"days":       daysJson,
				"exchanges":  "[]",
			}).Insert(); err != nil {
				return err
			}
		}
		if credit > 0 {
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
				Data(g.Map{"credit": &gdb.Counter{Field: "credit", Value: credit}}).Update(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepo) ListTasks(ctx context.Context) ([]*entity.UserTask, error) {
	var list []*entity.UserTask
	err := g.Model("user_task").Ctx(ctx).Where("status", 1).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *userRepo) FindTask(ctx context.Context, taskId int64) (*entity.UserTask, error) {
	var t *entity.UserTask
	err := g.Model("user_task").Ctx(ctx).Where("id", taskId).Scan(&t)
	return t, err
}

func (r *userRepo) TaskDoneToday(ctx context.Context, userId, taskId int64, logDate int) (int, error) {
	return g.Model("user_task_log").Ctx(ctx).
		Where("user_id", userId).Where("task_id", taskId).Where("log_date", logDate).Count()
}

// AddTaskLog 事务内: 写任务完成记录 + 发放积分。
func (r *userRepo) AddTaskLog(ctx context.Context, userId, taskId int64, typ string, logDate int, credit float64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("user_task_log").Ctx(ctx).Data(g.Map{
			"user_id":  userId,
			"task_id":  taskId,
			"type":     typ,
			"num":      1,
			"log_date": logDate,
		}).Insert(); err != nil {
			return err
		}
		if credit > 0 {
			if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
				Data(g.Map{"credit": &gdb.Counter{Field: "credit", Value: credit}}).Update(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepo) TaskLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserTaskLog, int, error) {
	m := g.Model("user_task_log").Ctx(ctx).Where("user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserTaskLog
	err = m.Clone().Page(page, size).OrderDesc("id").Scan(&list)
	return list, total, err
}

// ---- 充值 ----

func (r *userRepo) ListRechargePackages(ctx context.Context) ([]*entity.RechargePackage, error) {
	var list []*entity.RechargePackage
	err := g.Model("recharge_package").Ctx(ctx).Where("status", 1).OrderAsc("sort").Scan(&list)
	return list, err
}

func (r *userRepo) FindRechargePackage(ctx context.Context, id int64) (*entity.RechargePackage, error) {
	var p *entity.RechargePackage
	err := g.Model("recharge_package").Ctx(ctx).Where("id", id).Where("status", 1).Scan(&p)
	return p, err
}

// CreateRechargeOrder 创建待支付订单(到账逻辑在支付回调, 见 TODO)。
func (r *userRepo) CreateRechargeOrder(ctx context.Context, orderNo string, userId, packageId int64, amount, coin float64) error {
	_, err := g.Model("recharge_order").Ctx(ctx).Data(g.Map{
		"order_no":   orderNo,
		"user_id":    userId,
		"package_id": packageId,
		"amount":     amount,
		"coin":       coin,
		"status":     0,
	}).Insert()
	return err
}

// ---- VIP ----

func (r *userRepo) ListVipPackages(ctx context.Context) ([]*entity.VipPackage, error) {
	var list []*entity.VipPackage
	err := g.Model("vip_package").Ctx(ctx).Where("status", 1).OrderAsc("sort").Scan(&list)
	return list, err
}

func (r *userRepo) FindVipPackage(ctx context.Context, id int64) (*entity.VipPackage, error) {
	var p *entity.VipPackage
	err := g.Model("vip_package").Ctx(ctx).Where("id", id).Where("status", 1).Scan(&p)
	return p, err
}

// OpenVip 用金币开通/续费 VIP: 事务内扣金币(余额不足则失败) + 设置用户组 + 写记录 + 记账。
func (r *userRepo) OpenVip(ctx context.Context, userId int64, pkg *entity.VipPackage, startAt, endAt int64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("users").Ctx(ctx).
			Where("id", userId).Where("balance >= ?", pkg.Price).
			Data(g.Map{
				"balance":        &gdb.Counter{Field: "balance", Value: -pkg.Price},
				"group_id":       pkg.GroupId,
				"group_name":     pkg.Name,
				"group_end_time": endAt,
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("金币余额不足")
		}
		if _, err = tx.Model("vip_log").Ctx(ctx).Data(g.Map{
			"user_id": userId, "package_id": pkg.Id, "days": pkg.Days,
			"price": pkg.Price, "start_at": startAt, "end_at": endAt,
		}).Insert(); err != nil {
			return err
		}
		_, err = tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
			"user_id": userId, "direction": 2, "scene": "vip",
			"amount": pkg.Price, "remark": "开通VIP:" + pkg.Name,
		}).Insert()
		return err
	})
}

func (r *userRepo) VipLogs(ctx context.Context, userId int64, page, size int) ([]*entity.VipLog, int, error) {
	m := g.Model("vip_log").Ctx(ctx).Where("user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.VipLog
	err = m.Clone().Page(page, size).OrderDesc("id").Scan(&list)
	return list, total, err
}

// ---- 兑换(积分->金币) ----

func (r *userRepo) ExchangeCreditToCoin(ctx context.Context, userId int64, creditCost, coinGain float64) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("users").Ctx(ctx).
			Where("id", userId).Where("credit >= ?", creditCost).
			Data(g.Map{
				"credit":  &gdb.Counter{Field: "credit", Value: -creditCost},
				"balance": &gdb.Counter{Field: "balance", Value: coinGain},
			}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return gerror.New("积分不足")
		}
		_, err = tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
			"user_id": userId, "direction": 1, "scene": "exchange",
			"amount": coinGain, "remark": "积分兑换金币",
		}).Insert()
		return err
	})
}

// upsertConversation 会话 upsert(存在则更新, 否则插入)。
func upsertConversation(ctx context.Context, tx gdb.TX, userId, peerId int64, lastContent string, addUnread int) error {
	n, err := tx.Model("chat_conversation").Ctx(ctx).
		Where("user_id", userId).Where("peer_id", peerId).Count()
	if err != nil {
		return err
	}
	if n == 0 {
		_, err = tx.Model("chat_conversation").Ctx(ctx).Data(g.Map{
			"user_id": userId, "peer_id": peerId, "last_content": lastContent,
			"last_at": gtime.Now(), "unread": addUnread, "deleted": 0,
		}).Insert()
		return err
	}
	data := g.Map{"last_content": lastContent, "last_at": gtime.Now(), "deleted": 0}
	if addUnread > 0 {
		data["unread"] = &gdb.Counter{Field: "unread", Value: float64(addUnread)}
	}
	_, err = tx.Model("chat_conversation").Ctx(ctx).
		Where("user_id", userId).Where("peer_id", peerId).Data(data).Update()
	return err
}

// SendMessage 事务内: 写消息 + upsert 双方会话(对方未读+1)。
func (r *userRepo) SendMessage(ctx context.Context, fromId, toId int64, content string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("chat_message").Ctx(ctx).Data(g.Map{
			"from_id": fromId, "to_id": toId, "content": content,
		}).Insert(); err != nil {
			return err
		}
		if err := upsertConversation(ctx, tx, fromId, toId, content, 0); err != nil {
			return err
		}
		return upsertConversation(ctx, tx, toId, fromId, content, 1)
	})
}

func (r *userRepo) ListConversations(ctx context.Context, userId int64, page, size int) ([]*entity.ChatConversation, int, error) {
	m := g.Model("chat_conversation").Ctx(ctx).Where("user_id", userId).Where("deleted", 0)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ChatConversation
	err = m.Clone().Page(page, size).OrderDesc("last_at").Scan(&list)
	return list, total, err
}

func (r *userRepo) Messages(ctx context.Context, meId, peerId int64, page, size int) ([]*entity.ChatMessage, int, error) {
	m := g.Model("chat_message").Ctx(ctx).
		Where("(from_id=? AND to_id=?) OR (from_id=? AND to_id=?)", meId, peerId, peerId, meId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ChatMessage
	err = m.Clone().Page(page, size).OrderDesc("id").Scan(&list)
	return list, total, err
}

func (r *userRepo) MarkRead(ctx context.Context, userId, peerId int64) error {
	_, err := g.Model("chat_conversation").Ctx(ctx).
		Where("user_id", userId).Where("peer_id", peerId).
		Data(g.Map{"unread": 0}).Update()
	return err
}

func (r *userRepo) DeleteConversation(ctx context.Context, userId, peerId int64) error {
	_, err := g.Model("chat_conversation").Ctx(ctx).
		Where("user_id", userId).Where("peer_id", peerId).
		Data(g.Map{"deleted": 1}).Update()
	return err
}

// ==================== 后台管理(B1) ====================

// AdminListUsers 后台用户列表(筛选+分页)。
func (r *userRepo) AdminListUsers(ctx context.Context, f domain.AdminUserFilter, page, size int) ([]*entity.Users, int, error) {
	m := g.Model("users").Ctx(ctx)
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(username ILIKE ? OR phone ILIKE ? OR nickname ILIKE ?)", kw, kw, kw)
	}
	if f.Channel != "" {
		m = m.Where("channel_name", f.Channel)
	}
	if f.GroupId > 0 {
		m = m.Where("group_id", f.GroupId)
	}
	switch f.Status {
	case 1:
		m = m.Where("is_disabled", 0)
	case 2:
		m = m.Where("is_disabled", 1)
	}
	if f.StartDate > 0 {
		m = m.Where("register_date >= ?", f.StartDate)
	}
	if f.EndDate > 0 {
		m = m.Where("register_date <= ?", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Users
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// SetDisabled 禁用(1)/解禁(0)。解禁时清空禁用原因。
func (r *userRepo) SetDisabled(ctx context.Context, id int64, disabled int, reason string) error {
	if disabled == 0 {
		reason = ""
	}
	_, err := g.Model("users").Ctx(ctx).Where("id", id).Data(g.Map{
		"is_disabled": disabled,
		"error_msg":   reason,
	}).Update()
	return err
}

// UpdateGroup 调整用户组快照字段(组定义表 B4 再建)。
func (r *userRepo) UpdateGroup(ctx context.Context, id, groupId int64, groupName string, groupRate int, groupEndTime int64) error {
	_, err := g.Model("users").Ctx(ctx).Where("id", id).Data(g.Map{
		"group_id":       groupId,
		"group_name":     groupName,
		"group_rate":     groupRate,
		"group_end_time": groupEndTime,
	}).Update()
	return err
}

// AdminAdjustBalance 后台调整金币(balance)/积分(credit): 事务+行锁, 记录前后值流水。
func (r *userRepo) AdminAdjustBalance(ctx context.Context, userId int64, target string, amount float64, refId, remark string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var u *entity.Users
		if err := tx.Model("users").Ctx(ctx).Where("id", userId).LockUpdate().Scan(&u); err != nil {
			return err
		}
		if u == nil {
			return gerror.New("用户不存在")
		}
		before := u.Balance
		if target == "credit" {
			before = u.Credit
		}
		after := before + amount
		if after < 0 {
			return gerror.New("调整后余额为负, 拒绝执行")
		}
		if _, err := tx.Model("users").Ctx(ctx).Where("id", userId).
			Data(g.Map{target: after}).Update(); err != nil {
			return err
		}
		direction, amt := 1, amount
		if amount < 0 {
			direction, amt = 2, -amount
		}
		scene := "admin_balance"
		if target == "credit" {
			scene = "admin_credit"
		}
		_, err := tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
			"user_id": userId, "direction": direction, "scene": scene,
			"amount": amt, "balance_before": before, "balance_after": after,
			"ref_id": refId, "remark": remark,
		}).Insert()
		return err
	})
}

// BalanceLogs 用户余额流水(倒序分页)。
func (r *userRepo) BalanceLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserBalanceLog, int, error) {
	m := g.Model("user_balance_log").Ctx(ctx).Where("user_id", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserBalanceLog
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// ==================== 用户组定义(B4) ====================

func (r *userRepo) GroupList(ctx context.Context) ([]*entity.UserGroup, error) {
	var list []*entity.UserGroup
	err := g.Model("user_group").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *userRepo) GroupFind(ctx context.Context, id int64) (*entity.UserGroup, error) {
	var ug *entity.UserGroup
	err := g.Model("user_group").Ctx(ctx).Where("id", id).Scan(&ug)
	return ug, err
}

func (r *userRepo) GroupCreate(ctx context.Context, ug *entity.UserGroup) (int64, error) {
	res, err := g.Model("user_group").Ctx(ctx).Data(g.Map{
		"name": ug.Name, "rate": ug.Rate, "rights": ug.Rights,
		"remark": ug.Remark, "sort": ug.Sort, "status": ug.Status,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GroupUpdate 更新组定义, 并同步该组所有用户的快照字段(name/rate)。
func (r *userRepo) GroupUpdate(ctx context.Context, ug *entity.UserGroup) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("user_group").Ctx(ctx).Where("id", ug.Id).Data(g.Map{
			"name": ug.Name, "rate": ug.Rate, "rights": ug.Rights,
			"remark": ug.Remark, "sort": ug.Sort, "status": ug.Status,
			"updated_at": gtime.Now(),
		}).Update(); err != nil {
			return err
		}
		_, err := tx.Model("users").Ctx(ctx).Where("group_id", ug.Id).Data(g.Map{
			"group_name": ug.Name, "group_rate": ug.Rate,
		}).Update()
		return err
	})
}

func (r *userRepo) GroupDelete(ctx context.Context, id int64) error {
	_, err := g.Model("user_group").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *userRepo) GroupUserCount(ctx context.Context, groupId int64) (int, error) {
	return g.Model("users").Ctx(ctx).Where("group_id", groupId).Count()
}

// ==================== 成长配置(B5) ====================

func (r *userRepo) TaskListAll(ctx context.Context) ([]*entity.UserTask, error) {
	var list []*entity.UserTask
	err := g.Model("user_task").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *userRepo) TaskCreate(ctx context.Context, t *entity.UserTask) (int64, error) {
	res, err := g.Model("user_task").Ctx(ctx).Data(g.Map{
		"name": t.Name, "type": t.Type, "description": t.Description,
		"max_num": t.MaxNum, "reward": t.Reward, "status": t.Status, "sort": t.Sort,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *userRepo) TaskUpdate(ctx context.Context, t *entity.UserTask) error {
	_, err := g.Model("user_task").Ctx(ctx).Where("id", t.Id).Data(g.Map{
		"name": t.Name, "type": t.Type, "description": t.Description,
		"max_num": t.MaxNum, "reward": t.Reward, "status": t.Status, "sort": t.Sort,
		"updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *userRepo) TaskDelete(ctx context.Context, id int64) error {
	_, err := g.Model("user_task").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *userRepo) TaskLogList(ctx context.Context, f domain.TaskLogFilter, page, size int) ([]*entity.UserTaskLog, int, error) {
	m := g.Model("user_task_log").Ctx(ctx)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.TaskId > 0 {
		m = m.Where("task_id", f.TaskId)
	}
	if f.Type != "" {
		m = m.Where("type", f.Type)
	}
	if f.StartDate > 0 {
		m = m.Where("log_date >= ?", f.StartDate)
	}
	if f.EndDate > 0 {
		m = m.Where("log_date <= ?", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserTaskLog
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// SignStats 某月签到统计: 签到用户数 / 总签到人次 / 按日分布(jsonb 展开聚合)。
func (r *userRepo) SignStats(ctx context.Context, yearMonth int) (int, int, []domain.SignDayCount, error) {
	// 签到用户数 + 总人次
	one, err := g.DB().GetOne(ctx,
		`SELECT count(*) AS user_count, coalesce(sum(jsonb_array_length(days)), 0) AS sign_count
		   FROM user_sign WHERE year_month = ?`, yearMonth)
	if err != nil {
		return 0, 0, nil, err
	}
	userCount := one["user_count"].Int()
	signCount := one["sign_count"].Int()
	// 按日分布
	all, err := g.DB().GetAll(ctx,
		`SELECT d::int AS day, count(*) AS cnt
		   FROM user_sign, jsonb_array_elements_text(days) AS d
		  WHERE year_month = ?
		  GROUP BY 1 ORDER BY 1`, yearMonth)
	if err != nil {
		return 0, 0, nil, err
	}
	days := make([]domain.SignDayCount, 0, len(all))
	for _, rec := range all {
		days = append(days, domain.SignDayCount{Day: rec["day"].Int(), Count: rec["cnt"].Int()})
	}
	return userCount, signCount, days, nil
}

// ==================== 社交查询(B6) ====================

// FollowList 关注关系(联表出双方用户名)。
func (r *userRepo) FollowList(ctx context.Context, f domain.FollowFilter, page, size int) ([]*domain.FollowItem, int, error) {
	m := g.Model("user_follow f").Ctx(ctx)
	if f.UserId > 0 {
		m = m.Where("f.user_id", f.UserId)
	}
	if f.HomeId > 0 {
		m = m.Where("f.home_id", f.HomeId)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*domain.FollowItem
	err = m.Clone().
		LeftJoin("users u1", "u1.id=f.user_id").
		LeftJoin("users u2", "u2.id=f.home_id").
		Fields("f.id, f.user_id, u1.username AS user_name, f.home_id, u2.username AS home_name, f.created_at").
		OrderDesc("f.id").Page(page, size).Scan(&list)
	return list, total, err
}
