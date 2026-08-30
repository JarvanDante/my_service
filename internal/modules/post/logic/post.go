// Package logic 帖子业务(移植自 tianbi post/community 核心面)。
// UGC 过滤: 标题+内容过 filter_word; 发布默认待审; 审核用条件更新保证幂等。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/post/service"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
	"github.com/JarvanDante/my_service/internal/shared/storage"
)

const postSiteId = 1 // 单站点样板
const defaultPostCategory = "最新"

type sPost struct{}

func New() service.IPost { return &sPost{} }

func hitFilterWord(ctx context.Context, text string) (string, error) {
	var words []*entity.FilterWord
	if err := g.Model("filter_word").Ctx(ctx).
		Where("site_id", postSiteId).Fields("word").Scan(&words); err != nil {
		return "", err
	}
	for _, w := range words {
		if w.Word != "" && strings.Contains(text, w.Word) {
			return w.Word, nil
		}
	}
	return "", nil
}

func toDTO(r *entity.Post) *service.PostDTO {
	pics := []string{}
	if r.Pics != "" {
		_ = json.Unmarshal([]byte(r.Pics), &pics)
	}
	topics := []string{}
	if r.Topics != "" {
		_ = json.Unmarshal([]byte(r.Topics), &topics)
	}
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.PostDTO{
		Id: r.Id, UserId: r.UserId, Title: r.Title, Content: r.Content, Pics: pics,
		Topics: topics, Category: r.Category, VideoUrl: r.VideoUrl, MediaId: r.MediaId, ViewCount: r.ViewCount,
		Rank: r.Rank,
		LikeCount: r.LikeCount, CommentCount: r.CommentCount, Status: r.Status,
		RejectReason: r.RejectReason, CreatedAt: created,
	}
}

func signPostDTO(ctx context.Context, d *service.PostDTO) {
	if d == nil {
		return
	}
	d.VideoUrl = storage.SignPlayURL(ctx, d.VideoUrl)
}

func fillAuthors(ctx context.Context, list []*service.PostDTO) {
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
	rows, err := g.Model("users").Ctx(ctx).WhereIn("id", ids).Fields("id,nickname,img,sex").All()
	if err != nil {
		return
	}
	type author struct {
		nickname string
		img      string
		sex      int
	}
	m := map[int64]author{}
	for _, row := range rows {
		m[row["id"].Int64()] = author{row["nickname"].String(), row["img"].String(), row["sex"].Int()}
	}
	for _, d := range list {
		if d == nil {
			continue
		}
		if a, ok := m[d.UserId]; ok {
			d.Nickname, d.Img, d.Sex = a.nickname, a.img, a.sex
		}
	}
	if set, err := paywall.ActiveVipSet(ctx, ids); err == nil {
		for _, d := range list {
			if d != nil {
				d.IsVip = set[d.UserId]
			}
		}
	}
}

func (s *sPost) Create(ctx context.Context, in service.CreateInput) (int64, error) {
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" {
		return 0, gerror.New("标题不能为空")
	}
	if content == "" {
		return 0, gerror.New("请输入内容")
	}
	topics := trimTopics(in.Topics)
	if len(topics) == 0 {
		return 0, gerror.New("请选择话题")
	}
	if len(in.Pics) == 0 {
		return 0, gerror.New("请上传图片")
	}
	if len(in.Pics) > 9 {
		return 0, gerror.New("最多上传9张图片")
	}
	if hit, err := hitFilterWord(ctx, title+" "+content); err != nil {
		return 0, err
	} else if hit != "" {
		return 0, gerror.New("内容包含违禁词, 请修改后重试")
	}
	picsJSON, _ := json.Marshal(in.Pics)
	topicsJSON, _ := json.Marshal(topics)
	id, err := g.Model("post").Ctx(ctx).Data(g.Map{
		"site_id": postSiteId, "user_id": in.UserId, "title": title,
		"content": content, "pics": string(picsJSON), "topics": string(topicsJSON),
		"video_url": strings.TrimSpace(in.VideoUrl), "media_id": in.MediaId,
		"category": defaultPostCategory, "status": 0,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func trimTopics(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, raw := range in {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
		if len(out) >= 10 {
			break
		}
	}
	return out
}

func (s *sPost) FrontList(ctx context.Context, f service.ListFilter) ([]*service.PostDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("post").Ctx(ctx).Where("site_id", postSiteId).Where("status", 1)
	if f.FollowOnly {
		if f.ViewerId <= 0 {
			return []*service.PostDTO{}, 0, nil
		}
		ids, err := g.Model("user_follow").Ctx(ctx).
			Where("user_id", f.ViewerId).Fields("home_id").Limit(500).Array()
		if err != nil {
			return nil, 0, err
		}
		followed := make([]int64, 0, len(ids))
		for _, v := range ids {
			if id := v.Int64(); id > 0 {
				followed = append(followed, id)
			}
		}
		if len(followed) == 0 {
			return []*service.PostDTO{}, 0, nil
		}
		m = m.WhereIn("user_id", followed)
	} else if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	if cat := strings.TrimSpace(f.Category); cat != "" {
		m = m.Where("category", cat)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	q := m.Clone().OrderDesc("rank")
	if f.Sort == "hot" {
		q = q.OrderDesc("like_count").OrderDesc("id")
	} else {
		q = q.OrderDesc("id")
	}
	var list []*entity.Post
	if err := q.Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.PostDTO, 0, len(list))
	for _, r := range list {
		d := toDTO(r)
		d.Status, d.RejectReason = 0, "" // 前台流不暴露审核态
		signPostDTO(ctx, d)
		out = append(out, d)
	}
	fillAuthors(ctx, out)
	return out, total, nil
}

func (s *sPost) Detail(ctx context.Context, id, viewerId int64) (*service.PostDTO, error) {
	var r *entity.Post
	if err := g.Model("post").Ctx(ctx).
		Where("site_id", postSiteId).Where("id", id).Scan(&r); err != nil {
		return nil, err
	}
	if r == nil || r.Status == 3 {
		return nil, gerror.New("帖子不存在")
	}
	if r.Status != 1 && r.UserId != viewerId {
		return nil, gerror.New("帖子不存在或未通过审核")
	}
	_, _ = g.Model("post").Ctx(ctx).Where("id", id).
		Data(g.Map{"view_count": &gdb.Counter{Field: "view_count", Value: 1}}).Update()
	d := toDTO(r)
	d.ViewCount++
	signPostDTO(ctx, d)
	fillAuthors(ctx, []*service.PostDTO{d})
	return d, nil
}

func (s *sPost) My(ctx context.Context, userId int64, page, size int) ([]*service.PostDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("post").Ctx(ctx).
		Where("site_id", postSiteId).Where("user_id", userId).WhereNot("status", 3)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Post
	if err := m.Clone().OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.PostDTO, 0, len(list))
	for _, r := range list {
		d := toDTO(r)
		signPostDTO(ctx, d)
		out = append(out, d)
	}
	fillAuthors(ctx, out)
	return out, total, nil
}

func (s *sPost) DeleteOwn(ctx context.Context, userId, id int64) error {
	res, err := g.Model("post").Ctx(ctx).
		Where("site_id", postSiteId).Where("id", id).Where("user_id", userId).
		WhereNot("status", 3).
		Data(g.Map{"status": 3, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("帖子不存在或无权删除")
	}
	return nil
}

func (s *sPost) List(ctx context.Context, f service.ListFilter) ([]*service.PostDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("post").Ctx(ctx).Where("site_id", postSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Post
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.PostDTO, 0, len(list))
	for _, r := range list {
		d := toDTO(r)
		signPostDTO(ctx, d)
		out = append(out, d)
	}
	return out, total, nil
}

func (s *sPost) Update(ctx context.Context, in service.UpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.ViewCount < 0 {
		return gerror.New("浏览量不能为负")
	}
	if in.Rank < 0 {
		return gerror.New("权重不能为负")
	}
	cat := strings.TrimSpace(in.Category)
	if cat != "" {
		cnt, err := g.Model("post_category").Ctx(ctx).
			Where("site_id", postSiteId).Where("name", cat).Where("status", 1).Count()
		if err != nil {
			return err
		}
		if cnt == 0 {
			return gerror.New("分类不存在或已禁用")
		}
	}
	res, err := g.Model("post").Ctx(ctx).
		Where("site_id", postSiteId).Where("id", in.Id).
		Data(g.Map{
			"category": cat, "view_count": in.ViewCount, "rank": in.Rank, "updated_at": gtime.Now(),
		}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("帖子不存在")
	}
	return nil
}

func (s *sPost) Audit(ctx context.Context, id int64, pass bool, reason string) error {
	newStatus, data := 1, g.Map{"updated_at": gtime.Now()}
	if !pass {
		if strings.TrimSpace(reason) == "" {
			return gerror.New("拒绝需填写原因")
		}
		newStatus = 2
		data["reject_reason"] = reason
	}
	data["status"] = newStatus
	res, err := g.Model("post").Ctx(ctx).
		Where("site_id", postSiteId).Where("id", id).Where("status", 0). // 仅待审可审
		Data(data).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("帖子不存在或已审核过")
	}
	return nil
}

func (s *sPost) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("post").Ctx(ctx).
			Where("site_id", postSiteId).Where("id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("comment").Ctx(ctx).
			Where("site_id", postSiteId).Where("media_type", 2).Where("content_id", id).Delete()
		return err
	})
}
