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
