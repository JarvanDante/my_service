// Package logic UGC投稿业务。
//
// 关键取舍:
//   - 投稿是 UGC 入口, 标题/简介先过 filter_word 敏感词(与 post 模块同一套判定), 命中直接拒绝,
//     脏内容不进库比进库后再删干净;
//   - 撤回与审核都走「条件更新」(WHERE status=0 ...), 用 RowsAffected==0 判失败:
//     状态流转的并发安全交给数据库, 不做"先查后改"的两步判断(那是 TOCTOU, 双击会重复处理);
//   - 撤回不删行而是置 status=3, 投稿痕迹要留(统计/风控/申诉都要用);
//   - resource/tags 都是只整体读写的数组, 用 jsonb 存。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/publish/service"
)

const pbSiteId = 1

type sPublish struct{}

func New() service.IPublish { return &sPublish{} }

// hitFilterWord 命中返回该敏感词, 未命中返回空串(与 post 模块同一实现口径:
// 词量小, 全量拉出来做 Contains 足够, 不引入额外的 AC 自动机依赖)。
func hitFilterWord(ctx context.Context, text string) (string, error) {
	var words []*entity.FilterWord
	if err := g.Model("filter_word").Ctx(ctx).
		Where("site_id", pbSiteId).Fields("word").Scan(&words); err != nil {
		return "", err
	}
	for _, w := range words {
		if w.Word != "" && strings.Contains(text, w.Word) {
			return w.Word, nil
		}
	}
	return "", nil
}

func decodeStrings(raw string) []string {
	out := []string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func encodeJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func toDTO(r *entity.UserPublish) *service.PublishDTO {
	return &service.PublishDTO{
		Id: r.Id, UserId: r.UserId, Type: r.Type, Title: r.Title, Intro: r.Intro,
		Cover: r.Cover, Resource: decodeStrings(r.Resource), Tags: decodeStrings(r.Tags),
		Status: r.Status, RejectReason: r.RejectReason, AuditBy: r.AuditBy,
		AuditAt: fmtTime(r.AuditAt), CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sPublish) Submit(ctx context.Context, in service.SubmitInput) (int64, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return 0, gerror.New("标题不能为空")
	}
	if in.Type < entity.PublishTypeVideo || in.Type > entity.PublishTypePhoto {
		return 0, gerror.New("投稿类型非法")
	}
	if hit, err := hitFilterWord(ctx, title+" "+in.Intro); err != nil {
		return 0, err
	} else if hit != "" {
		return 0, gerror.New("内容包含违禁词, 请修改后重试")
	}
	return g.Model("user_publish").Ctx(ctx).Data(g.Map{
		"site_id": pbSiteId, "user_id": in.UserId, "type": in.Type, "title": title,
		"intro": in.Intro, "cover": in.Cover, "resource": encodeJSON(in.Resource),
		"tags": encodeJSON(in.Tags), "status": entity.PublishStatusPending,
	}).InsertAndGetId()
}

func (s *sPublish) query(ctx context.Context, f service.ListFilter) ([]*service.PublishDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("user_publish").Ctx(ctx).Where("site_id", pbSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Type > 0 {
		m = m.Where("type", f.Type)
	}
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.UserPublish
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.PublishDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sPublish) My(ctx context.Context, userId int64, f service.ListFilter) ([]*service.PublishDTO, int, error) {
	if userId <= 0 {
		return nil, 0, gerror.New("未登录")
	}
	f.UserId = userId // 我的列表永远锁死当前用户, 不接受客户端指定
	return s.query(ctx, f)
}

func (s *sPublish) Cancel(ctx context.Context, userId, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	// 条件更新: 只有"我的 + 待审"才能撤回。这样重复点撤回 / 撤回别人的 / 撤回已审的
	// 都会落到 RowsAffected==0 这一条分支上, 不需要额外查询。
	res, err := g.Model("user_publish").Ctx(ctx).
		Where("site_id", pbSiteId).Where("id", id).
		Where("user_id", userId).Where("status", entity.PublishStatusPending).
		Data(g.Map{
			"status": entity.PublishStatusCanceled, "updated_at": gtime.Now(),
		}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("投稿不存在或已被处理, 无法撤回")
	}
	return nil
}

// ---------------- 后台 ----------------

func (s *sPublish) List(ctx context.Context, f service.ListFilter) ([]*service.PublishDTO, int, error) {
	return s.query(ctx, f)
}

func (s *sPublish) Audit(ctx context.Context, id, adminId int64, pass bool, rejectReason string) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{
		"audit_by": adminId, "audit_at": gtime.Now(), "updated_at": gtime.Now(),
	}
	if pass {
		data["status"] = entity.PublishStatusPass
		data["reject_reason"] = "" // 通过时清掉历史理由, 避免前端展示到脏数据
	} else {
		if strings.TrimSpace(rejectReason) == "" {
			return gerror.New("拒绝需填写原因")
		}
		data["status"] = entity.PublishStatusReject
		data["reject_reason"] = rejectReason
	}
	// 条件更新: 仅待审可审, 双人同时点审核只有一次生效。
	res, err := g.Model("user_publish").Ctx(ctx).
		Where("site_id", pbSiteId).Where("id", id).
		Where("status", entity.PublishStatusPending).
		Data(data).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("投稿不存在或已被处理")
	}
	return nil
}
