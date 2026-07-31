package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	sysdomain "github.com/JarvanDante/my_service/internal/modules/system/domain"
)

type systemRepo struct{}

// NewSystemRepo 返回 system 领域仓储实现。
func NewSystemRepo() sysdomain.Repository { return &systemRepo{} }

func (r *systemRepo) NoticeCreate(ctx context.Context, n *entity.SystemNotice) (int64, error) {
	res, err := g.Model("system_notice").Ctx(ctx).Data(g.Map{
		"title": n.Title, "content": n.Content, "type": n.Type,
		"status": 1, "created_by": n.CreatedBy,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *systemRepo) NoticeList(ctx context.Context, f sysdomain.NoticeFilter, page, size int) ([]*entity.SystemNotice, int, error) {
	m := g.Model("system_notice").Ctx(ctx)
	if f.Type != "" {
		m = m.Where("type", f.Type)
	}
	switch f.Status {
	case 1:
		m = m.Where("status", 1)
	case 2:
		m = m.Where("status", 0)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.SystemNotice
	err = m.Clone().OrderDesc("id").Page(page, size).Scan(&list)
	return list, total, err
}

func (r *systemRepo) NoticeFind(ctx context.Context, id int64) (*entity.SystemNotice, error) {
	var n *entity.SystemNotice
	err := g.Model("system_notice").Ctx(ctx).Where("id", id).Scan(&n)
	return n, err
}

func (r *systemRepo) NoticeSetStatus(ctx context.Context, id int64, status int) error {
	_, err := g.Model("system_notice").Ctx(ctx).Where("id", id).Data(g.Map{
		"status": status, "updated_at": gtime.Now(),
	}).Update()
	return err
}
