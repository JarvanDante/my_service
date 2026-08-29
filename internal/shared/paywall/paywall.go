// Package paywall 内容付费解锁的统一判定与购买(漫画/小说/图集/视频共用)。
//
// tianbi 那套 purchaseser 混了 VIP 卡分类、观影券、免费次数、预售卡、分成…太多历史包袱。
// 这里只保留三态模型, 其余留给后续按需扩展:
//
//	price > 0   → 金币付费(买过 content_purchase 即解锁)
//	is_vip = 1  → VIP 专享(有生效中的 vip_log 即解锁), 与 price 互斥
//	两者皆无     → 完全免费
//
// 章节级别的"前 N 章免费"由各模块自己按 free_count 判断, 不进这里。
package paywall

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/shared/balance"
)

const pwSiteId = 1

// media_type 全局约定(与 user_collect / content_purchase 同一套编码)
const (
	MediaVideo  = 1
	MediaPost   = 2
	MediaComics = 3
	MediaNovel  = 4
	MediaPhoto  = 5
)

// Access 解锁判定结果, 直接下发给前端决定 UI(去购买 / 去开会员 / 直接看)。
type Access struct {
	Playable bool    `json:"playable"` // 是否可看
	IsBuy    bool    `json:"is_buy"`   // 是否已购买
	NeedPay  bool    `json:"need_pay"` // 需要金币购买
	NeedVip  bool    `json:"need_vip"` // 需要开通VIP
	Price    float64 `json:"price"`    // 价格(金币)
	Enough   bool    `json:"enough"`   // 余额是否够买
	Reason   string  `json:"reason"`   // 不可看的原因文案
}

// IsVipActive 是否有生效中的 VIP。
// 后台改用户组写的是 users.group_id / group_end_time；前台开通还会写 vip_log。两处任一未过期即算 VIP。
func IsVipActive(ctx context.Context, userId int64) (bool, error) {
	if userId <= 0 {
		return false, nil
	}
	now := gtime.Now().Unix()
	n, err := g.Model("users").Ctx(ctx).Where("id", userId).
		Where("group_id > ?", 0).Where("group_end_time > ?", now).Count()
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	n, err = g.Model("vip_log").Ctx(ctx).
		Where("site_id", pwSiteId).Where("user_id", userId).
		Where("end_at > ?", now).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Purchased 是否已购买该内容。
func Purchased(ctx context.Context, userId int64, mediaType int, contentId int64) (bool, error) {
	if userId <= 0 {
		return false, nil
	}
	n, err := g.Model("content_purchase").Ctx(ctx).
		Where("site_id", pwSiteId).Where("user_id", userId).
		Where("media_type", mediaType).Where("content_id", contentId).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// PurchasedSet 批量判断已购(列表页用, 避免 N+1)。
func PurchasedSet(ctx context.Context, userId int64, mediaType int, ids []int64) (map[int64]bool, error) {
	out := make(map[int64]bool, len(ids))
	if userId <= 0 || len(ids) == 0 {
		return out, nil
	}
	all, err := g.Model("content_purchase").Ctx(ctx).
		Where("site_id", pwSiteId).Where("user_id", userId).
		Where("media_type", mediaType).WhereIn("content_id", ids).
		Fields("content_id").Array()
	if err != nil {
		return nil, err
	}
	for _, v := range all {
		out[v.Int64()] = true
	}
	return out, nil
}

// Check 解锁判定。isVip/price 来自内容本身。
func Check(ctx context.Context, userId int64, mediaType int, contentId int64,
	isVip bool, price float64) (*Access, error) {
	a := &Access{Price: price}
	// 完全免费
	if !isVip && price <= 0 {
		a.Playable = true
		return a, nil
	}
	if userId <= 0 {
		a.NeedPay = price > 0
		a.NeedVip = isVip
		a.Reason = "请先登录"
		return a, nil
	}
	// VIP 专享
	if isVip {
		ok, err := IsVipActive(ctx, userId)
		if err != nil {
			return nil, err
		}
		if ok {
			a.Playable = true
			return a, nil
		}
		a.NeedVip = true
		a.Reason = "该内容为会员专享, 请先开通会员"
		return a, nil
	}
	// 金币付费
	bought, err := Purchased(ctx, userId, mediaType, contentId)
	if err != nil {
		return nil, err
	}
	if bought {
		a.Playable, a.IsBuy = true, true
		return a, nil
	}
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return nil, err
	}
	a.NeedPay = true
	a.Enough = bal != nil && bal.Float64() >= price
	a.Reason = "该内容需要购买后观看"
	return a, nil
}

// ErrAlreadyBought 重复购买(调用方通常提示"已购买, 无需重复购买")。
var ErrAlreadyBought = gerror.New("已购买过该内容")

// Buy 购买内容: 一个事务里完成「插购买记录(唯一约束防重) → 条件扣款 → 写流水」。
//
// 顺序是先插记录再扣款: 唯一约束先把并发重复购买挡在门外, 避免 tianbi
// "先查是否买过再扣钱" 的 TOCTOU 竞态(双击/重试会扣两次钱)。
func Buy(ctx context.Context, userId int64, mediaType int, contentId int64,
	title string, price float64) error {
	if price <= 0 {
		return gerror.New("该内容无需购买")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Model("content_purchase").Ctx(ctx).Data(g.Map{
			"site_id": pwSiteId, "user_id": userId, "media_type": mediaType,
			"content_id": contentId, "title": title, "amount": price,
		}).InsertIgnore()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return ErrAlreadyBought
		}
		id, _ := res.LastInsertId()
		return balance.Deduct(ctx, tx, userId, price, balance.SceneContentBuy,
			gconv.String(id), "购买:"+title)
	})
}
