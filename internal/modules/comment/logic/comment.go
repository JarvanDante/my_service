// Package logic 评论业务(移植自 tianbi comment)。
// 树形: parent_id 保留树结构, root_id 反范式化到顶层, 列表两查合并, 无需递归 CTE。
// UGC 过滤: 命中 filter_word 直接拒绝。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/comment/service"
)

const cmtSiteId = 1 // 单站点样板

type sComment struct{}

func New() service.IComment { return &sComment{} }

// hitFilterWord 命中敏感词返回该词(词表小, 全量载入内存匹配)。
func hitFilterWord(ctx context.Context, text string) (string, error) {
	var words []*entity.FilterWord
	if err := g.Model("filter_word").Ctx(ctx).
		Where("site_id", cmtSiteId).Fields("word").Scan(&words); err != nil {
		return "", err
	}
	for _, w := range words {
		if w.Word != "" && strings.Contains(text, w.Word) {
			return w.Word, nil
		}
	}
	return "", nil
}

func (s *sComment) Add(ctx context.Context, in service.AddInput) (int64, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return 0, gerror.New("评论内容不能为空")
	}
	if hit, err := hitFilterWord(ctx, content); err != nil {
		return 0, err
	} else if hit != "" {
		return 0, gerror.New("内容包含违禁词, 请修改后重试")
	}
	rootId := int64(0)
	if in.ParentId > 0 {
		var parent *entity.Comment
		if err := g.Model("comment").Ctx(ctx).
			Where("site_id", cmtSiteId).Where("id", in.ParentId).Where("status", 1).
			Scan(&parent); err != nil {
			return 0, err
		}
		if parent == nil {
			return 0, gerror.New("被回复的评论不存在")
		}
		if parent.MediaType != in.MediaType || parent.ContentId != in.ContentId {
			return 0, gerror.New("回复目标与内容不匹配")
		}
		rootId = parent.Id
		if parent.RootId > 0 { // 回复的回复, 锚到同一顶层
			rootId = parent.RootId
		}
	}
	var newId int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("comment").Ctx(ctx).Data(g.Map{
			"site_id": cmtSiteId, "user_id": in.UserId, "media_type": in.MediaType,
			"content_id": in.ContentId, "parent_id": in.ParentId, "root_id": rootId,
			"content": content, "status": 1,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newId = id
		if rootId > 0 { // 顶层评论回复数+1
			if _, err := tx.Model("comment").Ctx(ctx).Where("id", rootId).
				Data(g.Map{"reply_count": &gdb.Counter{Field: "reply_count", Value: 1}}).
				Update(); err != nil {
				return err
			}
		}
		if in.MediaType == 2 { // 帖子评论数+1
			if _, err := tx.Model("post").Ctx(ctx).
				Where("site_id", cmtSiteId).Where("id", in.ContentId).
				Data(g.Map{"comment_count": &gdb.Counter{Field: "comment_count", Value: 1}}).
				Update(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func toDTO(r *entity.Comment) service.ItemDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return service.ItemDTO{
		Id: r.Id, UserId: r.UserId, ParentId: r.ParentId, RootId: r.RootId,
		Content: r.Content, LikeCount: r.LikeCount, ReplyCount: r.ReplyCount,
		CreatedAt: created,
	}
}

func (s *sComment) List(ctx context.Context, mediaType int, contentId int64, page, size int) ([]service.ItemDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	base := g.Model("comment").Ctx(ctx).
		Where("site_id", cmtSiteId).Where("media_type", mediaType).
		Where("content_id", contentId).Where("status", 1)
	// 顶层分页
	top := base.Clone().Where("root_id", 0)
	total, err := top.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var tops []*entity.Comment
	if err := top.Clone().OrderDesc("id").Page(page, size).Scan(&tops); err != nil {
		return nil, 0, err
	}
	if len(tops) == 0 {
		return []service.ItemDTO{}, total, nil
	}
	// 本页顶层的全部回复(一查带回, parent_id 保留树形)
	rootIds := make([]int64, 0, len(tops))
	for _, t := range tops {
		rootIds = append(rootIds, t.Id)
	}
	var replies []*entity.Comment
	if err := base.Clone().WhereIn("root_id", rootIds).OrderAsc("id").Scan(&replies); err != nil {
		return nil, 0, err
	}
	replyMap := map[int64][]service.ItemDTO{}
	for _, r := range replies {
		replyMap[r.RootId] = append(replyMap[r.RootId], toDTO(r))
	}
	out := make([]service.ItemDTO, 0, len(tops))
	for _, t := range tops {
		d := toDTO(t)
		d.Replies = replyMap[t.Id]
		out = append(out, d)
	}
	return out, total, nil
}
