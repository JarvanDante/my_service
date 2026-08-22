// Package logic 意见反馈业务(移植自 tianbi feedbackser)。
package logic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/feedback/service"
	msglogic "github.com/JarvanDante/my_service/internal/modules/message/logic"
	msgsvc "github.com/JarvanDante/my_service/internal/modules/message/service"
)

const feedbackSiteId = 1 // 单站点样板

type sFeedback struct{}

func New() service.IFeedback { return &sFeedback{} }

// Add 提交反馈(1 分钟内限 1 条, 同 tianbi)。
func (s *sFeedback) Add(ctx context.Context, in service.AddInput) (int64, error) {
	if in.Content == "" {
		return 0, gerror.New("反馈内容不能为空")
	}
	cnt, err := g.Model("feedback").Ctx(ctx).
		Where("site_id", feedbackSiteId).Where("user_id", in.UserId).
		Where("created_at >= ?", gtime.Now().Add(-time.Minute)).Count()
	if err != nil {
		return 0, err
	}
	if cnt >= 1 {
		return 0, gerror.New("操作太频繁，请稍后再试")
	}
	if in.Type == 0 {
		in.Type = 1
	}
	picsJSON := "[]"
	if len(in.Pics) > 0 {
		if b, e := json.Marshal(in.Pics); e == nil {
			picsJSON = string(b)
		}
	}
	id, err := g.Model("feedback").Ctx(ctx).Data(g.Map{
		"site_id": feedbackSiteId, "user_id": in.UserId, "type": in.Type,
		"problem_type": in.ProblemType, "content": in.Content, "pics": picsJSON,
		"sys_info": in.SysInfo, "media_id": in.MediaId, "media_title": in.MediaTitle,
		"status": 1,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sFeedback) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("feedback").Ctx(ctx).Where("site_id", feedbackSiteId)
	if f.Status > 0 {
		m = m.Where("status", f.Status)
	}
	if f.Type > 0 {
		m = m.Where("type", f.Type)
	}
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Feedback
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		pics := []string{}
		if r.Pics != "" {
			_ = json.Unmarshal([]byte(r.Pics), &pics)
		}
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.ItemDTO{
			Id: r.Id, UserId: r.UserId, Type: r.Type, ProblemType: r.ProblemType,
			Content: r.Content, Pics: pics, SysInfo: r.SysInfo, MediaId: r.MediaId,
			MediaTitle: r.MediaTitle, Status: r.Status, Reply: r.Reply, CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sFeedback) Handle(ctx context.Context, id int64, reply string, status int) error {
	if status != 1 && status != 2 {
		status = 2
	}
	var row entity.Feedback
	if err := g.Model("feedback").Ctx(ctx).Where("id", id).Scan(&row); err != nil {
		return err
	}
	if row.Id == 0 {
		return gerror.New("反馈不存在")
	}
	_, err := g.Model("feedback").Ctx(ctx).Where("id", id).Data(g.Map{
		"reply": reply, "status": status, "updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return err
	}
	notifyFeedbackReply(ctx, row.UserId, row.Content, reply)
	return nil
}

func notifyFeedbackReply(ctx context.Context, userId int64, content, reply string) {
	reply = strings.TrimSpace(reply)
	if userId <= 0 || reply == "" {
		return
	}
	runes := []rune(strings.TrimSpace(content))
	snippet := string(runes)
	if len(runes) > 40 {
		snippet = string(runes[:40]) + "…"
	}
	body := "客服回复：" + reply
	if snippet != "" {
		body = "您的反馈：" + snippet + "\n" + body
	}
	if _, err := msglogic.New().Create(ctx, msgsvc.CreateInput{
		UserId: userId, Type: 1, Title: "反馈已处理", Content: body,
	}); err != nil {
		g.Log().Warningf(ctx, "feedback reply notify user=%d: %v", userId, err)
	}
}
