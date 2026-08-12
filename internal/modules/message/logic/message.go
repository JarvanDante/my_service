// Package logic 系统消息业务。
// 可见性: 前台可见 = (user_id=0 全员 OR user_id=我) AND status=1。
// 已读: 独立 read 表 + 唯一约束, 标记用 InsertIgnore 幂等。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/message/service"
)

const msgSiteId = 1 // 单站点样板

type sMessage struct{}

func New() service.IMessage { return &sMessage{} }

// visibleModel 前台可见消息(全员 + 发给我的, 已发布)。
func visibleModel(ctx context.Context, userId int64) *gdb.Model {
	return g.Model("sys_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("status", 1).
		Where("(user_id = 0 OR user_id = ?)", userId)
}

// MyList 我的消息列表(带已读标记)。
func (s *sMessage) MyList(ctx context.Context, userId int64, page, size int) ([]service.MsgDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := visibleModel(ctx, userId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.SysMessage
	if err := m.Clone().OrderDesc("created_at").OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	// 本页消息的已读集合(两查合并, 避免 JOIN 拼串)
	readSet := map[int64]bool{}
	if len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, r := range list {
			ids = append(ids, r.Id)
		}
		var reads []*entity.SysMessageRead
		if err := g.Model("sys_message_read").Ctx(ctx).
			Where("site_id", msgSiteId).Where("user_id", userId).
			WhereIn("message_id", ids).Scan(&reads); err != nil {
			return nil, 0, err
		}
		for _, r := range reads {
			readSet[r.MessageId] = true
		}
	}
	out := make([]service.MsgDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, service.MsgDTO{
			Id: r.Id, UserId: r.UserId, Type: r.Type, Title: r.Title, Content: r.Content,
			IsRead: readSet[r.Id], Status: r.Status, CreatedAt: created,
		})
	}
	return out, total, nil
}

// UnreadCount 未读数 = 可见消息数 - 已读数。
func (s *sMessage) UnreadCount(ctx context.Context, userId int64) (int, error) {
	cnt, err := visibleModel(ctx, userId).
		Where("id NOT IN (SELECT message_id FROM sys_message_read WHERE site_id = ? AND user_id = ?)",
			msgSiteId, userId).Count()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// MarkRead 标记已读(幂等)。
func (s *sMessage) MarkRead(ctx context.Context, userId, id int64, all bool) error {
	if !all && id <= 0 {
		return gerror.New("消息ID必填(或 all=true)")
	}
	m := visibleModel(ctx, userId)
	if !all {
		m = m.Where("id", id)
	}
	var list []*entity.SysMessage
	if err := m.Fields("id").Scan(&list); err != nil {
		return err
	}
	if len(list) == 0 {
		if !all {
			return gerror.New("消息不存在")
		}
		return nil
	}
	rows := make([]g.Map, 0, len(list))
	for _, r := range list {
		rows = append(rows, g.Map{
			"site_id": msgSiteId, "user_id": userId, "message_id": r.Id,
		})
	}
	_, err := g.Model("sys_message_read").Ctx(ctx).Data(rows).InsertIgnore()
	return err
}

func (s *sMessage) List(ctx context.Context, f service.ListFilter) ([]*service.MsgDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("sys_message").Ctx(ctx).Where("site_id", msgSiteId)
	if f.UserId >= 0 { // -1=全部
		m = m.Where("user_id", f.UserId)
	}
	if f.Status >= 0 { // -1=全部
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("title ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.SysMessage
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.MsgDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.MsgDTO{
			Id: r.Id, UserId: r.UserId, Type: r.Type, Title: r.Title, Content: r.Content,
			Status: r.Status, CreatedAt: created,
		})
	}
	return out, total, nil
}

func (s *sMessage) Create(ctx context.Context, in service.CreateInput) (int64, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return 0, gerror.New("标题不能为空")
	}
	if in.UserId < 0 {
		return 0, gerror.New("user_id 非法")
	}
	if in.Type <= 0 {
		in.Type = 1
	}
	id, err := g.Model("sys_message").Ctx(ctx).Data(g.Map{
		"site_id": msgSiteId, "user_id": in.UserId, "type": in.Type,
		"title": title, "content": in.Content, "status": 1,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sMessage) Update(ctx context.Context, in service.UpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"updated_at": gtime.Now()}
	if in.Type > 0 {
		data["type"] = in.Type
	}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Content != "" {
		data["content"] = in.Content
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("sys_message").Ctx(ctx).
		Where("site_id", msgSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sMessage) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("sys_message").Ctx(ctx).
			Where("site_id", msgSiteId).Where("id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("sys_message_read").Ctx(ctx).
			Where("site_id", msgSiteId).Where("message_id", id).Delete()
		return err
	})
}
