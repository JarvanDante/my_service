package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/message/service"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const (
	channelComment = "comment"
	channelLike    = "like"

	subWorkComment = "work_comment"
	subReply       = "reply"
	subMention     = "mention"
	subCommentLike = "comment_like"
	subPostLike    = "post_like"

	targetComment = "comment"
	targetReply   = "reply"
	targetPost    = "post"

	jumpPageSize = 15 // 与帖子详情定位拉评一致(公司 latest 算法)
	statusLive   = 1
)

func (s *sMessage) UnreadAll(ctx context.Context, userId int64) (service.UnreadBreakdown, error) {
	out := service.UnreadBreakdown{}
	sys, err := s.UnreadCount(ctx, userId)
	if err != nil {
		return out, err
	}
	out.Sys = sys
	comment, err := g.Model("interact_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("user_id", userId).
		Where("channel", channelComment).Where("is_read", 0).Count()
	if err != nil {
		return out, err
	}
	like, err := g.Model("interact_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("user_id", userId).
		Where("channel", channelLike).Where("is_read", 0).Count()
	if err != nil {
		return out, err
	}
	out.Comment = comment
	out.Like = like
	return out, nil
}

func (s *sMessage) InteractList(ctx context.Context, userId int64, channel string, page, size int) ([]service.InteractDTO, int, error) {
	channel, err := normalizeChannel(channel)
	if err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("interact_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("user_id", userId).Where("channel", channel)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var rows []*entity.InteractMessage
	if err := m.Clone().OrderDesc("updated_at").OrderDesc("id").Page(page, size).Scan(&rows); err != nil {
		return nil, 0, err
	}
	out := make([]service.InteractDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toInteractDTO(r))
	}
	fillActors(ctx, out)
	fillJumps(ctx, out)
	return out, total, nil
}

func (s *sMessage) MarkInteractRead(ctx context.Context, userId, id int64, all bool, channel string) error {
	if !all && id <= 0 {
		return gerror.New("消息ID必填(或 all=true)")
	}
	m := g.Model("interact_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("user_id", userId).Where("is_read", 0)
	if !all {
		m = m.Where("id", id)
		n, err := m.Clone().Count()
		if err != nil {
			return err
		}
		if n == 0 {
			var exists int
			exists, err = g.Model("interact_message").Ctx(ctx).
				Where("site_id", msgSiteId).Where("user_id", userId).Where("id", id).Count()
			if err != nil {
				return err
			}
			if exists == 0 {
				return gerror.New("消息不存在")
			}
			return nil
		}
	} else if channel != "" {
		ch, err := normalizeChannel(channel)
		if err != nil {
			return err
		}
		m = m.Where("channel", ch)
	}
	_, err := m.Data(g.Map{"is_read": 1}).Update()
	return err
}

// NotifyWorkComment 作品新评论：通知作者（自己评自己不发）。
func NotifyWorkComment(ctx context.Context, actorId int64, mediaType int, contentId, commentId int64, snippet string) {
	if actorId <= 0 || contentId <= 0 || commentId <= 0 {
		return
	}
	authorId, title := resolveObject(ctx, mediaType, contentId)
	if authorId <= 0 || authorId == actorId {
		return
	}
	insertCommentMessage(ctx, g.Map{
		"user_id": authorId, "sub_type": subWorkComment, "actor_id": actorId,
		"media_type": mediaType, "content_id": contentId, "object_title": title,
		"target_type": targetComment, "comment_id": commentId, "root_comment_id": commentId,
		"snippet": snippet,
	})
}

// NotifyReply 回顶层评 → 通知父评作者；回回复 → 通知被回复人。不通知作品作者「别人互回」。
func NotifyReply(ctx context.Context, actorId, receiverId int64, mediaType int, contentId, commentId, rootId int64, mention bool, snippet string) {
	if actorId <= 0 || receiverId <= 0 || receiverId == actorId || commentId <= 0 {
		return
	}
	sub := subReply
	if mention {
		sub = subMention
	}
	_, title := resolveObject(ctx, mediaType, contentId)
	if rootId <= 0 {
		rootId = commentId
	}
	insertCommentMessage(ctx, g.Map{
		"user_id": receiverId, "sub_type": sub, "actor_id": actorId,
		"media_type": mediaType, "content_id": contentId, "object_title": title,
		"target_type": targetReply, "comment_id": commentId, "root_comment_id": rootId,
		"snippet": snippet,
	})
}

// NotifyCommentLike 评论点赞：按目标 + 北京时间自然日聚合。取消赞不新建。
func NotifyCommentLike(ctx context.Context, actorId, receiverId int64, mediaType int, contentId, commentId int64, snippet string, isLike bool) {
	if actorId <= 0 || receiverId <= 0 || receiverId == actorId || commentId <= 0 {
		return
	}
	_, title := resolveObject(ctx, mediaType, contentId)
	aggKey := fmt.Sprintf("like:%d:comment:%d:%s", receiverId, commentId, bizDate())
	upsertLike(ctx, upsertLikeIn{
		ReceiverId: receiverId, ActorId: actorId, SubType: subCommentLike,
		MediaType: mediaType, ContentId: contentId, Title: title,
		TargetType: targetComment, CommentId: commentId, RootId: commentId,
		Snippet: snippet, AggKey: aggKey, IsLike: isLike,
	})
}

// NotifyPostLike 帖子点赞：按帖 + 北京时间自然日聚合。
func NotifyPostLike(ctx context.Context, actorId, postId int64, isLike bool) {
	if actorId <= 0 || postId <= 0 {
		return
	}
	authorId, title := resolveObject(ctx, 2, postId)
	if authorId <= 0 || authorId == actorId {
		return
	}
	aggKey := fmt.Sprintf("like:%d:post:%d:%s", authorId, postId, bizDate())
	upsertLike(ctx, upsertLikeIn{
		ReceiverId: authorId, ActorId: actorId, SubType: subPostLike,
		MediaType: 2, ContentId: postId, Title: title,
		TargetType: targetPost, CommentId: 0, RootId: 0,
		Snippet: "", AggKey: aggKey, IsLike: isLike,
	})
}

func insertCommentMessage(ctx context.Context, data g.Map) {
	now := gtime.Now()
	data["site_id"] = msgSiteId
	data["channel"] = channelComment
	data["actor_ids"] = encodeIDs([]int64{toInt64(data["actor_id"])})
	data["like_count"] = 0
	data["is_read"] = 0
	data["agg_key"] = ""
	data["created_at"] = now
	data["updated_at"] = now
	if _, err := g.Model("interact_message").Ctx(ctx).Data(data).Insert(); err != nil {
		g.Log().Warningf(ctx, "interact comment insert: %v", err)
	}
}

type upsertLikeIn struct {
	ReceiverId int64
	ActorId    int64
	SubType    string
	MediaType  int
	ContentId  int64
	Title      string
	TargetType string
	CommentId  int64
	RootId     int64
	Snippet    string
	AggKey     string
	IsLike     bool
}

func upsertLike(ctx context.Context, in upsertLikeIn) {
	var row *entity.InteractMessage
	_ = g.Model("interact_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("agg_key", in.AggKey).Scan(&row)
	if !in.IsLike {
		revokeLikeActor(ctx, row, in.ActorId)
		return
	}
	now := gtime.Now()
	if row != nil {
		ids := decodeIDs(row.ActorIds)
		if containsID(ids, in.ActorId) {
			return
		}
		ids = append(ids, in.ActorId)
		if _, err := g.Model("interact_message").Ctx(ctx).Where("id", row.Id).Data(g.Map{
			"actor_id": in.ActorId, "actor_ids": encodeIDs(ids),
			"like_count": row.LikeCount + 1, "snippet": in.Snippet,
			"is_read": 0, "updated_at": now,
		}).Update(); err != nil {
			g.Log().Warningf(ctx, "interact like update: %v", err)
		}
		return
	}
	if _, err := g.Model("interact_message").Ctx(ctx).Data(g.Map{
		"site_id": msgSiteId, "user_id": in.ReceiverId, "channel": channelLike,
		"sub_type": in.SubType, "actor_id": in.ActorId, "actor_ids": encodeIDs([]int64{in.ActorId}),
		"like_count": 1, "media_type": in.MediaType, "content_id": in.ContentId,
		"object_title": in.Title, "target_type": in.TargetType,
		"comment_id": in.CommentId, "root_comment_id": in.RootId,
		"snippet": in.Snippet, "is_read": 0, "agg_key": in.AggKey,
		"created_at": now, "updated_at": now,
	}).Insert(); err != nil {
		g.Log().Warningf(ctx, "interact like insert: %v", err)
	}
}

func revokeLikeActor(ctx context.Context, row *entity.InteractMessage, actorId int64) {
	if row == nil {
		return
	}
	ids := decodeIDs(row.ActorIds)
	if !containsID(ids, actorId) {
		return
	}
	next := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id != actorId {
			next = append(next, id)
		}
	}
	last := int64(0)
	if n := len(next); n > 0 {
		last = next[n-1]
	}
	count := row.LikeCount - 1
	if count < 0 {
		count = 0
	}
	if _, err := g.Model("interact_message").Ctx(ctx).Where("id", row.Id).Data(g.Map{
		"actor_ids": encodeIDs(next), "actor_id": last, "like_count": count,
	}).Update(); err != nil {
		g.Log().Warningf(ctx, "interact like revoke: %v", err)
	}
}

func resolveObject(ctx context.Context, mediaType int, contentId int64) (int64, string) {
	if contentId <= 0 {
		return 0, ""
	}
	if mediaType == 2 {
		var p *entity.Post
		if err := g.Model("post").Ctx(ctx).Where("site_id", msgSiteId).Where("id", contentId).
			Fields("user_id,title").Scan(&p); err != nil || p == nil {
			return 0, ""
		}
		return p.UserId, p.Title
	}
	if mediaType == 1 {
		var v *entity.Video
		if err := g.Model("video").Ctx(ctx).Where("id", contentId).
			Fields("up_user_id,title").Scan(&v); err != nil || v == nil {
			return 0, ""
		}
		return v.UpUserId, v.Title
	}
	if mediaType == 4 {
		var c *entity.Comics
		if err := g.Model("comics").Ctx(ctx).Where("site_id", msgSiteId).Where("id", contentId).
			Fields("title").Scan(&c); err != nil || c == nil {
			return 0, ""
		}
		return 0, c.Title
	}
	if mediaType == 7 {
		var n *entity.Novel
		if err := g.Model("novel").Ctx(ctx).Where("site_id", msgSiteId).Where("id", contentId).
			Fields("title").Scan(&n); err != nil || n == nil {
			return 0, ""
		}
		return 0, n.Title
	}
	return 0, ""
}

func toInteractDTO(r *entity.InteractMessage) service.InteractDTO {
	created := ""
	if r.UpdatedAt != nil {
		created = r.UpdatedAt.String()
	} else if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	count := r.LikeCount
	if count <= 0 {
		count = 1
	}
	return service.InteractDTO{
		Id: r.Id, Channel: r.Channel, SubType: r.SubType, IsRead: r.IsRead == 1,
		CreatedAt: created, ActorId: r.ActorId, ActorCount: count,
		MediaType: r.MediaType, ContentId: r.ContentId, ObjectTitle: r.ObjectTitle,
		TargetType: r.TargetType, CommentId: r.CommentId, RootCommentId: r.RootCommentId,
		Snippet: r.Snippet, Page: 1, PageSize: jumpPageSize,
	}
}

func fillJumps(ctx context.Context, list []service.InteractDTO) {
	for i := range list {
		page, deleted := locateJump(ctx, list[i])
		list[i].Page = page
		list[i].PageSize = jumpPageSize
		list[i].Deleted = deleted
	}
}

// locateJump 按评论列表 latest(id DESC) 现算顶层评页码。楼中楼回复已随该页一次带回，不再算 reply_page。
func locateJump(ctx context.Context, d service.InteractDTO) (int, bool) {
	if d.ContentId <= 0 || d.MediaType <= 0 {
		return 1, d.CommentId > 0
	}
	rootId := d.RootCommentId
	if rootId <= 0 {
		rootId = d.CommentId
	}
	if rootId <= 0 {
		return 1, false
	}
	var root *entity.Comment
	if err := g.Model("comment").Ctx(ctx).Where("id", rootId).
		Fields("id,status,media_type,content_id").Scan(&root); err != nil || root == nil || root.Status != statusLive {
		return 1, true
	}
	newer, err := g.Model("comment").Ctx(ctx).
		Where("site_id", msgSiteId).Where("media_type", d.MediaType).Where("content_id", d.ContentId).
		Where("status", statusLive).Where("root_id", 0).Where("id > ?", root.Id).Count()
	if err != nil {
		return 1, true
	}
	page := newer/jumpPageSize + 1
	if d.CommentId > 0 && d.CommentId != rootId {
		var reply *entity.Comment
		if err := g.Model("comment").Ctx(ctx).Where("id", d.CommentId).Fields("id,status").Scan(&reply); err != nil || reply == nil || reply.Status != statusLive {
			return page, true
		}
	}
	return page, false
}

func fillActors(ctx context.Context, list []service.InteractDTO) {
	ids := make([]int64, 0, len(list))
	seen := map[int64]struct{}{}
	for i := range list {
		id := list[i].ActorId
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return
	}
	rows, err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Fields("id,nickname,img,sex").All()
	if err != nil {
		return
	}
	type actor struct {
		name, img string
		sex       int
	}
	m := map[int64]actor{}
	for _, row := range rows {
		m[row["id"].Int64()] = actor{row["nickname"].String(), row["img"].String(), row["sex"].Int()}
	}
	for i := range list {
		if a, ok := m[list[i].ActorId]; ok {
			list[i].ActorName = a.name
			list[i].ActorAvatar = a.img
			list[i].ActorSex = a.sex
		}
	}
	if set, err := paywall.ActiveVipSet(ctx, ids); err == nil {
		for i := range list {
			list[i].ActorIsVip = set[list[i].ActorId]
		}
	}
}

func normalizeChannel(channel string) (string, error) {
	channel = strings.ToLower(strings.TrimSpace(channel))
	if channel == "" {
		channel = channelComment
	}
	if channel != channelComment && channel != channelLike {
		return "", gerror.New("channel 仅支持 comment 或 like")
	}
	return channel, nil
}

func bizDate() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	return time.Now().In(loc).Format("2006-01-02")
}

func encodeIDs(ids []int64) string {
	if ids == nil {
		ids = []int64{}
	}
	b, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodeIDs(raw string) []int64 {
	out := []int64{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func containsID(ids []int64, id int64) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	default:
		return 0
	}
}
