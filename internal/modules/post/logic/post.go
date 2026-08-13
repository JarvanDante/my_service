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
)

const postSiteId = 1 // 单站点样板

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
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.PostDTO{
		Id: r.Id, UserId: r.UserId, Title: r.Title, Content: r.Content, Pics: pics,
		MediaId: r.MediaId, ViewCount: r.ViewCount, LikeCount: r.LikeCount,
		CommentCount: r.CommentCount, Status: r.Status, RejectReason: r.RejectReason,
		CreatedAt: created,
	}
}

func (s *sPost) Create(ctx context.Context, in service.CreateInput) (int64, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return 0, gerror.New("标题不能为空")
	}
	if hit, err := hitFilterWord(ctx, title+" "+in.Content); err != nil {
		return 0, err
	} else if hit != "" {
		return 0, gerror.New("内容包含违禁词, 请修改后重试")
	}
	picsJSON := "[]"
	if len(in.Pics) > 0 {
		if b, e := json.Marshal(in.Pics); e == nil {
			picsJSON = string(b)
		}
	}
	id, err := g.Model("post").Ctx(ctx).Data(g.Map{
		"site_id": postSiteId, "user_id": in.UserId, "title": title,
		"content": in.Content, "pics": picsJSON, "media_id": in.MediaId, "status": 0,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sPost) FrontList(ctx context.Context, f service.ListFilter) ([]*service.PostDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("post").Ctx(ctx).Where("site_id", postSiteId).Where("status", 1)
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	q := m.Clone()
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
		out = append(out, d)
	}
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
		out = append(out, toDTO(r))
	}
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
		out = append(out, toDTO(r))
	}
	return out, total, nil
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
