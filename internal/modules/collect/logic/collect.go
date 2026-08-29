// Package logic 收藏/点赞业务(移植自 tianbi collectser)。
// 幂等设计: 添加用 ON CONFLICT DO NOTHING(唯一约束), 取消用 DELETE, 重复操作不报错。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/collect/service"
	msglogic "github.com/JarvanDante/my_service/internal/modules/message/logic"
)

const (
	colSiteId   = 1   // 单站点样板
	colBatchMax = 100 // 单次批量上限
)

type sCollect struct{}

func New() service.ICollect { return &sCollect{} }

// Operate 添加/取消 收藏/点赞/踩。
func (s *sCollect) Operate(ctx context.Context, in service.OperateInput) error {
	if len(in.Ids) == 0 {
		return gerror.New("对象ID必填")
	}
	if len(in.Ids) > colBatchMax {
		return gerror.New("单次最多操作100条")
	}
	if in.Flag {
		// 添加: 唯一约束 + ON CONFLICT DO NOTHING 幂等
		rows := make([]g.Map, 0, len(in.Ids))
		for _, id := range in.Ids {
			if id <= 0 {
				continue
			}
			rows = append(rows, g.Map{
				"site_id": colSiteId, "user_id": in.UserId, "op_type": in.Type,
				"media_type": in.MediaType, "content_id": id,
			})
		}
		if len(rows) == 0 {
			return gerror.New("对象ID非法")
		}
		// InsertIgnore 在 PG 生成 ON CONFLICT DO NOTHING, 重复添加幂等
		_, err := g.Model("user_collect").Ctx(ctx).Data(rows).InsertIgnore()
		if err == nil && in.Type == 2 && in.MediaType == 2 {
			for _, id := range in.Ids {
				msglogic.NotifyPostLike(ctx, in.UserId, id, true)
			}
		}
		return err
	}
	err := s.Delete(ctx, in.UserId, in.Ids, in.MediaType, in.Type)
	if err == nil && in.Type == 2 && in.MediaType == 2 {
		for _, id := range in.Ids {
			msglogic.NotifyPostLike(ctx, in.UserId, id, false)
		}
	}
	return err
}

// Delete 批量取消。
func (s *sCollect) Delete(ctx context.Context, userId int64, ids []int64, mediaType, opType int) error {
	if len(ids) == 0 {
		return gerror.New("ids必填")
	}
	if len(ids) > colBatchMax {
		return gerror.New("单次最多操作100条")
	}
	_, err := g.Model("user_collect").Ctx(ctx).
		Where("site_id", colSiteId).Where("user_id", userId).
		Where("op_type", opType).Where("media_type", mediaType).
		WhereIn("content_id", ids).Delete()
	return err
}

// List 我的收藏/点赞列表。
func (s *sCollect) List(ctx context.Context, f service.ListFilter) ([]service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("user_collect").Ctx(ctx).
		Where("site_id", colSiteId).Where("user_id", f.UserId).Where("op_type", f.Type)
	if f.MediaType > 0 {
		m = m.Where("media_type", f.MediaType)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserCollect
	if err := m.Clone().OrderDesc("created_at").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]service.ItemDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, service.ItemDTO{
			ContentId: r.ContentId, MediaType: r.MediaType, CreatedAt: created,
		})
	}
	return out, total, nil
}
