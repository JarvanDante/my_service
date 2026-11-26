// Package logic 收藏/点赞业务(移植自 tianbi collectser)。
// 幂等设计: 添加用 ON CONFLICT DO NOTHING(唯一约束), 取消用 DELETE, 重复操作不报错。
// 漫画/小说/图集的 like_count 记收藏数, 帖子的 like_count 记点赞数; 只在真正写入/删除后回写。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/collect/service"
	msglogic "github.com/JarvanDante/my_service/internal/modules/message/logic"
)

const (
	colSiteId   = 1   // 单站点样板
	colBatchMax = 100 // 单次批量上限
	opFav       = 1
	opLike      = 2
	mediaPost   = 2
	mediaComics = 3
	mediaNovel  = 4
	mediaPhoto  = 5
	mediaDrama  = 7 // 短剧按整部收藏/点赞, content_id 为 drama.id
)

type sCollect struct{}

func New() service.ICollect { return &sCollect{} }

// countTarget 哪种操作要回写哪张表的 like_count。
func countTarget(mediaType, opType int) (table string, ok bool) {
	switch {
	case mediaType == mediaPost && opType == opLike:
		return "post", true
	case mediaType == mediaComics && opType == opFav:
		return "comics", true
	case mediaType == mediaNovel && opType == opFav:
		return "novel", true
	case mediaType == mediaPhoto && opType == opFav:
		return "photo", true
	case mediaType == mediaDrama && opType == opLike:
		return "drama", true
	default:
		return "", false
	}
}

func bumpLikeCount(ctx context.Context, tx gdb.TX, mediaType, opType int, contentId, delta int64) error {
	table, ok := countTarget(mediaType, opType)
	if !ok || contentId <= 0 || delta == 0 {
		return nil
	}
	m := tx.Model(table).Ctx(ctx).Where("id", contentId)
	if delta < 0 {
		m = m.Where("like_count > ?", 0)
	}
	_, err := m.Data(g.Map{"like_count": &gdb.Counter{Field: "like_count", Value: float64(delta)}}).Update()
	return err
}

// checkDramaOnline 添加短剧收藏/点赞前确认每部都存在且已上架。
func checkDramaOnline(ctx context.Context, ids []int64) error {
	uniq := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		uniq[id] = struct{}{}
	}
	n, err := g.Model("drama").Ctx(ctx).Where("site_id", colSiteId).
		Where("status", entity.ContentStatusOnline).WhereIn("id", ids).Count()
	if err != nil {
		return err
	}
	if n != len(uniq) {
		return gerror.New("短剧不存在或已下架")
	}
	return nil
}

// Operate 添加/取消 收藏/点赞/踩。
func (s *sCollect) Operate(ctx context.Context, in service.OperateInput) error {
	if len(in.Ids) == 0 {
		return gerror.New("对象ID必填")
	}
	if len(in.Ids) > colBatchMax {
		return gerror.New("单次最多操作100条")
	}
	if !in.Flag {
		err := s.Delete(ctx, in.UserId, in.Ids, in.MediaType, in.Type)
		if err == nil && in.Type == opLike && in.MediaType == mediaPost {
			for _, id := range in.Ids {
				msglogic.NotifyPostLike(ctx, in.UserId, id, false)
			}
		}
		return err
	}
	ids := make([]int64, 0, len(in.Ids))
	for _, id := range in.Ids {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return gerror.New("对象ID非法")
	}
	if in.MediaType == mediaDrama {
		if err := checkDramaOnline(ctx, ids); err != nil {
			return err
		}
	}
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range ids {
			res, err := tx.Model("user_collect").Ctx(ctx).Data(g.Map{
				"site_id": colSiteId, "user_id": in.UserId, "op_type": in.Type,
				"media_type": in.MediaType, "content_id": id,
			}).InsertIgnore()
			if err != nil {
				return err
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				if err = bumpLikeCount(ctx, tx, in.MediaType, in.Type, id, 1); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err == nil && in.Type == opLike && in.MediaType == mediaPost {
		for _, id := range ids {
			msglogic.NotifyPostLike(ctx, in.UserId, id, true)
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
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, id := range ids {
			if id <= 0 {
				continue
			}
			res, err := tx.Model("user_collect").Ctx(ctx).
				Where("site_id", colSiteId).Where("user_id", userId).
				Where("op_type", opType).Where("media_type", mediaType).
				Where("content_id", id).Delete()
			if err != nil {
				return err
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				if err = bumpLikeCount(ctx, tx, mediaType, opType, id, -1); err != nil {
					return err
				}
			}
		}
		return nil
	})
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
