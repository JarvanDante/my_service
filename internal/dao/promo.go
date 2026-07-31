package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	promodomain "github.com/JarvanDante/my_service/internal/modules/promo/domain"
)

type promoRepo struct{}

// NewPromoRepo 返回 promo 领域仓储实现。
func NewPromoRepo() promodomain.Repository { return &promoRepo{} }

func (r *promoRepo) ListCodes(ctx context.Context, f promodomain.CodeFilter, page, size int) ([]*entity.UserCode, int, error) {
	m := g.Model("user_code").Ctx(ctx)
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(code ILIKE ? OR name ILIKE ?)", kw, kw)
	}
	if f.CodeKey != "" {
		m = m.Where("code_key", f.CodeKey)
	}
	if f.Type != "" {
		m = m.Where("type", f.Type)
	}
	switch f.Status { // 0全部 1可用 2已使用 3作废
	case 1:
		m = m.Where("status", 0)
	case 2:
		m = m.Where("status", 1)
	case 3:
		m = m.Where("status", -1)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCode
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *promoRepo) FindCodeById(ctx context.Context, id int64) (*entity.UserCode, error) {
	var c *entity.UserCode
	err := g.Model("user_code").Ctx(ctx).Where("id", id).Scan(&c)
	return c, err
}

func (r *promoRepo) BatchCreateCodes(ctx context.Context, rows []*entity.UserCode) error {
	data := make([]g.Map, 0, len(rows))
	for _, c := range rows {
		data = append(data, g.Map{
			"name": c.Name, "code": c.Code, "code_key": c.CodeKey,
			"type": c.Type, "object_id": c.ObjectId, "add_num": c.AddNum,
			"can_use_num": c.CanUseNum, "status": 0, "expired_at": c.ExpiredAt,
		})
	}
	_, err := g.Model("user_code").Ctx(ctx).Data(data).Insert()
	return err
}

func (r *promoRepo) VoidCode(ctx context.Context, id int64) error {
	_, err := g.Model("user_code").Ctx(ctx).Where("id", id).Data(g.Map{
		"status": -1, "updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *promoRepo) ListCodeLogs(ctx context.Context, f promodomain.CodeLogFilter, page, size int) ([]*entity.UserCodeLog, int, error) {
	m := g.Model("user_code_log").Ctx(ctx)
	if f.CodeId > 0 {
		m = m.Where("code_id", f.CodeId)
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Code != "" {
		m = m.Where("code", f.Code)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCodeLog
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// ==================== 分享 / 拉新(B6) ====================

func (r *promoRepo) ShareLogList(ctx context.Context, f promodomain.ShareLogFilter, page, size int) ([]*entity.UserShareLog, int, error) {
	m := g.Model("user_share_log").Ctx(ctx)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Type != "" {
		m = m.Where("type", f.Type)
	}
	if f.Channel != "" {
		m = m.Where("channel", f.Channel)
	}
	if f.StartDate != "" {
		m = m.Where("created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		m = m.Where("created_at < ?::date + interval '1 day'", f.EndDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserShareLog
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

// ShareStats 分享总次数 / 分享人数 / 渠道分布。
func (r *promoRepo) ShareStats(ctx context.Context, startDate, endDate string) (int, int, []promodomain.ChannelCount, error) {
	m := g.Model("user_share_log").Ctx(ctx)
	if startDate != "" {
		m = m.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		m = m.Where("created_at < ?::date + interval '1 day'", endDate)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, 0, nil, err
	}
	sharerCount, err := m.Clone().Fields("DISTINCT user_id").Count()
	if err != nil {
		return 0, 0, nil, err
	}
	all, err := m.Clone().Fields("channel, count(*) AS cnt").Group("channel").OrderDesc("cnt").All()
	if err != nil {
		return 0, 0, nil, err
	}
	channels := make([]promodomain.ChannelCount, 0, len(all))
	for _, rec := range all {
		channels = append(channels, promodomain.ChannelCount{
			Channel: rec["channel"].String(), Count: rec["cnt"].Int(),
		})
	}
	return total, sharerCount, channels, nil
}

// InviteRank 拉新排行: 按推荐人聚合被邀请注册的用户数(register_date 范围可选, YYYY-MM-DD)。
func (r *promoRepo) InviteRank(ctx context.Context, startDate, endDate string, top int) ([]*promodomain.InviteRankItem, error) {
	if top <= 0 || top > 100 {
		top = 10
	}
	m := g.Model("users").Ctx(ctx).Where("parent_id > 0")
	if startDate != "" {
		m = m.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		m = m.Where("created_at < ?::date + interval '1 day'", endDate)
	}
	all, err := m.Fields("parent_id, max(parent_name) AS parent_name, count(*) AS cnt").
		Group("parent_id").OrderDesc("cnt").Limit(top).All()
	if err != nil {
		return nil, err
	}
	out := make([]*promodomain.InviteRankItem, 0, len(all))
	for _, rec := range all {
		out = append(out, &promodomain.InviteRankItem{
			UserId:      rec["parent_id"].Int64(),
			Username:    rec["parent_name"].String(),
			InviteCount: rec["cnt"].Int(),
		})
	}
	return out, nil
}
