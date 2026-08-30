// Package logic 评论业务(移植自 tianbi comment)。
// 树形: parent_id 保留树结构, root_id 反范式化到顶层, 列表两查合并, 无需递归 CTE。
// UGC 过滤: 命中 filter_word 直接拒绝。
// 上墙规则(社区二期): VIP 评论/回复直接上墙; 普通用户待审。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/comment/service"
	msglogic "github.com/JarvanDante/my_service/internal/modules/message/logic"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const (
	cmtSiteId           = 1
	collectMediaComment = 6 // user_collect.media_type, 评论点赞
	collectLike         = 2
)

const (
	statusPending = 0
	statusLive    = 1
	statusReject  = 2
)

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

func decodePics(raw string) []string {
	out := []string{}
	if raw == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func encodePics(list []string) string {
	clean := make([]string, 0, len(list))
	for _, p := range list {
		p = strings.TrimSpace(p)
		if p != "" {
			clean = append(clean, p)
		}
	}
	if len(clean) > 3 {
		clean = clean[:3]
	}
	b, err := json.Marshal(clean)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func (s *sComment) rejectIfMuted(ctx context.Context, userId int64) error {
	var u *entity.Users
	if err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("comment_muted").Scan(&u); err != nil {
		return err
	}
	if u != nil && u.CommentMuted == 1 {
		return gerror.New("你已被禁言，暂时不能评论")
	}
	return nil
}

func (s *sComment) Add(ctx context.Context, in service.AddInput) (int64, int, error) {
	if err := s.rejectIfMuted(ctx, in.UserId); err != nil {
		return 0, 0, err
	}
	content := strings.TrimSpace(in.Content)
	pics := encodePics(in.Pics)
	if content == "" && pics == "[]" {
		return 0, 0, gerror.New("评论内容不能为空")
	}
	if hit, err := hitFilterWord(ctx, content); err != nil {
		return 0, 0, err
	} else if hit != "" {
		return 0, 0, gerror.New("内容包含违禁词, 请修改后重试")
	}
	rootId := int64(0)
	if in.ParentId > 0 {
		var parent *entity.Comment
		if err := g.Model("comment").Ctx(ctx).
			Where("site_id", cmtSiteId).Where("id", in.ParentId).Where("status", statusLive).
			Scan(&parent); err != nil {
			return 0, 0, err
		}
		if parent == nil {
			return 0, 0, gerror.New("被回复的评论不存在")
		}
		if parent.MediaType != in.MediaType || parent.ContentId != in.ContentId {
			return 0, 0, gerror.New("回复目标与内容不匹配")
		}
		rootId = parent.Id
		if parent.RootId > 0 { // 回复的回复, 锚到同一顶层
			rootId = parent.RootId
		}
	}
	isVip, err := paywall.IsVipActive(ctx, in.UserId)
	if err != nil {
		return 0, 0, err
	}
	status := statusPending
	if isVip {
		status = statusLive
	}
	var newId int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("comment").Ctx(ctx).Data(g.Map{
			"site_id": cmtSiteId, "user_id": in.UserId, "media_type": in.MediaType,
			"content_id": in.ContentId, "parent_id": in.ParentId, "root_id": rootId,
			"content": content, "pics": pics, "status": status,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newId = id
		if status == statusLive {
			return bumpLiveCounts(ctx, tx, in.MediaType, in.ContentId, rootId)
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	if status == statusLive {
		notifyCommentLive(ctx, in.UserId, in.MediaType, in.ContentId, newId, in.ParentId, rootId, content, pics)
	}
	return newId, status, nil
}

func bumpLiveCounts(ctx context.Context, tx gdb.TX, mediaType int, contentId, rootId int64) error {
	if rootId > 0 {
		if _, err := tx.Model("comment").Ctx(ctx).Where("id", rootId).
			Data(g.Map{"reply_count": &gdb.Counter{Field: "reply_count", Value: 1}}).
			Update(); err != nil {
			return err
		}
	}
	if mediaType == 2 {
		if _, err := tx.Model("post").Ctx(ctx).
			Where("site_id", cmtSiteId).Where("id", contentId).
			Data(g.Map{"comment_count": &gdb.Counter{Field: "comment_count", Value: 1}}).
			Update(); err != nil {
			return err
		}
	}
	return nil
}

func toDTO(r *entity.Comment) service.ItemDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return service.ItemDTO{
		Id: r.Id, UserId: r.UserId, ParentId: r.ParentId, RootId: r.RootId,
		Content: r.Content, Pics: decodePics(r.Pics), LikeCount: r.LikeCount, ReplyCount: r.ReplyCount,
		CreatedAt: created,
	}
}

func (s *sComment) List(ctx context.Context, mediaType int, contentId int64, page, size, sort int, viewerId int64) ([]service.ItemDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	base := g.Model("comment").Ctx(ctx).
		Where("site_id", cmtSiteId).Where("media_type", mediaType).
		Where("content_id", contentId).Where("status", statusLive)
	// 外层数量=主评+回复；分页仍按顶层
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	top := base.Clone().Where("root_id", 0)
	var tops []*entity.Comment
	q := top.Clone()
	if sort == 1 {
		q = q.OrderDesc("like_count").OrderDesc("id")
	} else {
		q = q.OrderDesc("id")
	}
	if err := q.Page(page, size).Scan(&tops); err != nil {
		return nil, 0, err
	}
	if len(tops) == 0 {
		return []service.ItemDTO{}, total, nil
	}
	// 本页顶层的全部已上墙回复
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
	fillFrontUsers(ctx, out, viewerId)
	return out, total, nil
}

func (s *sComment) Like(ctx context.Context, userId, commentId int64, flag bool) (int, bool, error) {
	if userId <= 0 || commentId <= 0 {
		return 0, false, gerror.New("参数非法")
	}
	var row *entity.Comment
	if err := g.Model("comment").Ctx(ctx).Where("site_id", cmtSiteId).Where("id", commentId).
		Where("status", statusLive).Scan(&row); err != nil {
		return 0, false, err
	}
	if row == nil {
		return 0, false, gerror.New("评论不存在")
	}
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if flag {
			res, err := tx.Model("user_collect").Ctx(ctx).Data(g.Map{
				"site_id": cmtSiteId, "user_id": userId, "op_type": collectLike,
				"media_type": collectMediaComment, "content_id": commentId,
			}).InsertIgnore()
			if err != nil {
				return err
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				_, err = tx.Model("comment").Ctx(ctx).Where("id", commentId).
					Data(g.Map{"like_count": &gdb.Counter{Field: "like_count", Value: 1}}).Update()
			}
			return err
		}
		res, err := tx.Model("user_collect").Ctx(ctx).
			Where("site_id", cmtSiteId).Where("user_id", userId).
			Where("op_type", collectLike).Where("media_type", collectMediaComment).
			Where("content_id", commentId).Delete()
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			_, err = tx.Model("comment").Ctx(ctx).Where("id", commentId).Where("like_count > ?", 0).
				Data(g.Map{"like_count": &gdb.Counter{Field: "like_count", Value: -1}}).Update()
		}
		return err
	})
	if err != nil {
		return 0, false, err
	}
	msglogic.NotifyCommentLike(ctx, userId, row.UserId, row.MediaType, row.ContentId, commentId, commentSnippet(row.Content, row.Pics), flag)
	var fresh *entity.Comment
	_ = g.Model("comment").Ctx(ctx).Where("id", commentId).Fields("like_count").Scan(&fresh)
	count := 0
	if fresh != nil {
		count = fresh.LikeCount
	}
	return count, flag, nil
}

func walkItems(list []service.ItemDTO, fn func(*service.ItemDTO)) {
	for i := range list {
		fn(&list[i])
		if len(list[i].Replies) > 0 {
			walkItems(list[i].Replies, fn)
		}
	}
}

func fillFrontUsers(ctx context.Context, list []service.ItemDTO, viewerId int64) {
	ids := make([]int64, 0)
	cids := make([]int64, 0)
	seenU, seenC := map[int64]struct{}{}, map[int64]struct{}{}
	byId := map[int64]*service.ItemDTO{}
	walkItems(list, func(d *service.ItemDTO) {
		byId[d.Id] = d
		if d.UserId > 0 {
			if _, ok := seenU[d.UserId]; !ok {
				seenU[d.UserId] = struct{}{}
				ids = append(ids, d.UserId)
			}
		}
		if _, ok := seenC[d.Id]; !ok {
			seenC[d.Id] = struct{}{}
			cids = append(cids, d.Id)
		}
	})
	type author struct{ nickname, img string }
	users := map[int64]author{}
	if len(ids) > 0 {
		rows, err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Fields("id,nickname,img").All()
		if err == nil {
			for _, row := range rows {
				users[row["id"].Int64()] = author{row["nickname"].String(), row["img"].String()}
			}
		}
		vipSet, _ := paywall.ActiveVipSet(ctx, ids)
		walkItems(list, func(d *service.ItemDTO) {
			if a, ok := users[d.UserId]; ok {
				d.Nickname, d.Img = a.nickname, a.img
			}
			d.IsVip = vipSet[d.UserId]
			if d.ParentId > 0 {
				if p := byId[d.ParentId]; p != nil {
					d.ReplyUserId = p.UserId
					if a, ok := users[p.UserId]; ok && a.nickname != "" {
						d.ReplyNickname = a.nickname
					}
				}
			}
		})
	}
	if viewerId <= 0 || len(cids) == 0 {
		return
	}
	likedRows, err := g.Model("user_collect").Ctx(ctx).
		Where("site_id", cmtSiteId).Where("user_id", viewerId).
		Where("op_type", collectLike).Where("media_type", collectMediaComment).
		WhereIn("content_id", cids).Fields("content_id").Array()
	if err != nil {
		return
	}
	liked := map[int64]struct{}{}
	for _, v := range likedRows {
		if id := v.Int64(); id > 0 {
			liked[id] = struct{}{}
		}
	}
	walkItems(list, func(d *service.ItemDTO) {
		_, d.Liked = liked[d.Id]
	})
}

func (s *sComment) AdminList(ctx context.Context, f service.AdminListFilter) ([]*service.AdminItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("comment").Ctx(ctx).Where("site_id", cmtSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	switch f.Kind {
	case "main":
		m = m.Where("parent_id", 0)
	case "reply":
		m = m.Where("parent_id > ?", 0)
	}
	if f.Keyword != "" {
		m = m.Where("content ILIKE ?", "%"+f.Keyword+"%")
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	filtered, err := applyAdminMediaFilter(ctx, m, f.MediaType)
	if err != nil {
		return nil, 0, err
	}
	m = filtered
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Comment
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.AdminItemDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.AdminItemDTO{
			Id: r.Id, UserId: r.UserId, MediaType: r.MediaType, ContentId: r.ContentId,
			ParentId: r.ParentId, RootId: r.RootId, Content: r.Content,
			LikeCount: r.LikeCount, ReplyCount: r.ReplyCount, Status: r.Status,
			BelongLabel: mediaBelongLabel(r.MediaType), CreatedAt: created,
		})
	}
	fillAdminUsers(ctx, out)
	fillBelongLabels(ctx, out)
	return out, total, nil
}

const (
	adminFilterCartoon = 8 // 审核筛选项：video.kind=动漫
	adminFilterDouyin  = 9 // 审核筛选项：video.kind=抖音
)

func applyAdminMediaFilter(ctx context.Context, m *gdb.Model, mediaType int) (*gdb.Model, error) {
	if mediaType <= 0 {
		return m, nil
	}
	kind := 0
	switch mediaType {
	case adminFilterCartoon:
		kind = entity.VideoKindCartoon
	case adminFilterDouyin:
		kind = entity.VideoKindDouyin
	case 1:
		kind = entity.VideoKindVideo
	default:
		return m.Where("media_type", mediaType), nil
	}
	ids, err := videoIDsByKind(ctx, kind)
	if err != nil {
		return m, err
	}
	m = m.Where("media_type", 1)
	if len(ids) == 0 {
		return m.Where("1 = 0"), nil
	}
	return m.WhereIn("content_id", ids), nil
}

func videoIDsByKind(ctx context.Context, kind int) ([]int64, error) {
	arr, err := g.Model("video").Ctx(ctx).Where("kind", kind).Fields("id").Array()
	if err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(arr))
	for _, v := range arr {
		if id := v.Int64(); id > 0 {
			out = append(out, id)
		}
	}
	return out, nil
}

func mediaBelongLabel(mediaType int) string {
	switch mediaType {
	case 1:
		return "视频"
	case 2:
		return "帖子"
	case 4:
		return "漫画"
	case 7:
		return "小说"
	default:
		return ""
	}
}

func fillBelongLabels(ctx context.Context, list []*service.AdminItemDTO) {
	ids := make([]int64, 0)
	seen := map[int64]struct{}{}
	for _, d := range list {
		if d == nil || d.MediaType != 1 || d.ContentId <= 0 {
			continue
		}
		if _, ok := seen[d.ContentId]; ok {
			continue
		}
		seen[d.ContentId] = struct{}{}
		ids = append(ids, d.ContentId)
	}
	if len(ids) == 0 {
		return
	}
	rows, err := g.Model("video").Ctx(ctx).WhereIn("id", ids).Fields("id,kind").All()
	if err != nil {
		return
	}
	kinds := map[int64]int{}
	for _, row := range rows {
		kinds[row["id"].Int64()] = row["kind"].Int()
	}
	for _, d := range list {
		if d == nil || d.MediaType != 1 {
			continue
		}
		switch kinds[d.ContentId] {
		case entity.VideoKindCartoon:
			d.BelongLabel = "动漫"
		case entity.VideoKindDouyin:
			d.BelongLabel = "抖音"
		default:
			d.BelongLabel = "视频"
		}
	}
}

func fillAdminUsers(ctx context.Context, list []*service.AdminItemDTO) {
	ids := make([]int64, 0, len(list))
	seen := map[int64]struct{}{}
	for _, d := range list {
		if d == nil || d.UserId <= 0 {
			continue
		}
		if _, ok := seen[d.UserId]; ok {
			continue
		}
		seen[d.UserId] = struct{}{}
		ids = append(ids, d.UserId)
	}
	if len(ids) == 0 {
		return
	}
	rows, err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Fields("id,nickname,img").All()
	if err != nil {
		return
	}
	type author struct{ nickname, img string }
	m := map[int64]author{}
	for _, row := range rows {
		m[row["id"].Int64()] = author{row["nickname"].String(), row["img"].String()}
	}
	vipSet, _ := paywall.ActiveVipSet(ctx, ids)
	for _, d := range list {
		if d == nil {
			continue
		}
		if a, ok := m[d.UserId]; ok {
			d.Nickname, d.Img = a.nickname, a.img
		}
		d.IsVip = vipSet[d.UserId]
	}
}

func (s *sComment) Audit(ctx context.Context, id int64, pass bool) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	newStatus := statusLive
	if !pass {
		newStatus = statusReject
	}
	var row *entity.Comment
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := tx.Model("comment").Ctx(ctx).
			Where("site_id", cmtSiteId).Where("id", id).Where("status", statusPending).
			LockUpdate().Scan(&row); err != nil {
			return err
		}
		if row == nil {
			return gerror.New("评论不存在或已审核过")
		}
		if _, err := tx.Model("comment").Ctx(ctx).Where("id", id).
			Data(g.Map{"status": newStatus}).Update(); err != nil {
			return err
		}
		if newStatus == statusLive {
			return bumpLiveCounts(ctx, tx, row.MediaType, row.ContentId, row.RootId)
		}
		return nil
	})
	if err == nil && pass && row != nil {
		notifyCommentLive(ctx, row.UserId, row.MediaType, row.ContentId, row.Id, row.ParentId, row.RootId, row.Content, row.Pics)
	}
	return err
}

func notifyCommentLive(ctx context.Context, actorId int64, mediaType int, contentId, commentId, parentId, rootId int64, content, pics string) {
	snippet := commentSnippet(content, pics)
	if parentId <= 0 {
		msglogic.NotifyWorkComment(ctx, actorId, mediaType, contentId, commentId, snippet)
		return
	}
	var parent *entity.Comment
	if err := g.Model("comment").Ctx(ctx).Where("id", parentId).Fields("id,user_id,root_id").Scan(&parent); err != nil || parent == nil {
		return
	}
	mention := parent.RootId > 0
	root := rootId
	if root <= 0 {
		if parent.RootId > 0 {
			root = parent.RootId
		} else {
			root = parent.Id
		}
	}
	msglogic.NotifyReply(ctx, actorId, parent.UserId, mediaType, contentId, commentId, root, mention, snippet)
}

func commentSnippet(content, pics string) string {
	text := strings.TrimSpace(content)
	if text != "" {
		runes := []rune(text)
		if len(runes) > 80 {
			return string(runes[:80])
		}
		return text
	}
	if pics != "" && pics != "[]" {
		return "[图片]"
	}
	return ""
}
