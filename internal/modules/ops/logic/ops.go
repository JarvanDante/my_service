// Package logic 运营配置业务(移植自 tianbi announcement/jumptab/filterword)。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/ops/service"
)

const opsSiteId = 1 // 单站点样板

type sOps struct{}

func New() service.IOps { return &sOps{} }

func ts(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func annDTO(r *entity.Announcement) *service.AnnDTO {
	return &service.AnnDTO{
		Id: r.Id, Title: r.Title, Content: r.Content, TextNode: r.TextNode,
		Cover: r.Cover, JumpUrl: r.JumpUrl, SysType: r.SysType,
		StartAt: ts(r.StartAt), EndAt: ts(r.EndAt), Status: r.Status, CreatedAt: ts(r.CreatedAt),
	}
}

func jtDTO(r *entity.Jumptab) *service.JtDTO {
	return &service.JtDTO{
		Id: r.Id, CnName: r.CnName, EnName: r.EnName, Avatar: r.Avatar,
		Link: r.Link, PicJumpLink: r.PicJumpLink, Location: r.Location,
		Rank: r.Rank, Status: r.Status, CreatedAt: ts(r.CreatedAt),
	}
}

// ---------- 前台 ----------

// LiveAnnouncements 有效期内启用公告(移植自 tianbi AnnouncementList)。
func (s *sOps) LiveAnnouncements(ctx context.Context, sysType string) ([]service.AnnDTO, error) {
	now := gtime.Now()
	m := g.Model("announcement").Ctx(ctx).
		Where("site_id", opsSiteId).Where("status", 1).
		Where("start_at <= ?", now).Where("end_at >= ?", now)
	if sysType != "" {
		m = m.Where("sys_type", sysType)
	}
	var list []*entity.Announcement
	if err := m.OrderDesc("id").Limit(20).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]service.AnnDTO, 0, len(list))
	for _, r := range list {
		out = append(out, *annDTO(r))
	}
	return out, nil
}

// FrontJumptabs 启用跳转位(rank desc)。
func (s *sOps) FrontJumptabs(ctx context.Context, location int) ([]service.JtDTO, error) {
	m := g.Model("jumptab").Ctx(ctx).
		Where("site_id", opsSiteId).Where("status", 1)
	if location > 0 {
		m = m.Where("location", location)
	}
	var list []*entity.Jumptab
	if err := m.OrderDesc("rank").OrderDesc("id").Limit(50).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]service.JtDTO, 0, len(list))
	for _, r := range list {
		out = append(out, *jtDTO(r))
	}
	return out, nil
}

// ---------- 公告管理 ----------

func (s *sOps) AnnList(ctx context.Context, f service.PageFilter) ([]*service.AnnDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("announcement").Ctx(ctx).Where("site_id", opsSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Announcement
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.AnnDTO, 0, len(list))
	for _, r := range list {
		out = append(out, annDTO(r))
	}
	return out, total, nil
}

func (s *sOps) AnnCreate(ctx context.Context, in service.AnnSaveInput) (int64, error) {
	if strings.TrimSpace(in.Title) == "" {
		return 0, gerror.New("标题不能为空")
	}
	end := gtime.NewFromStr(in.EndAt)
	if end == nil || end.IsZero() {
		return 0, gerror.New("结束时间格式非法, 如: 2027-12-31 23:59:59")
	}
	start := gtime.Now()
	if in.StartAt != "" {
		if st := gtime.NewFromStr(in.StartAt); st != nil && !st.IsZero() {
			start = st
		}
	}
	if in.SysType == "" {
		in.SysType = "app"
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("announcement").Ctx(ctx).Data(g.Map{
		"site_id": opsSiteId, "title": in.Title, "content": in.Content,
		"text_node": in.TextNode, "cover": in.Cover, "jump_url": in.JumpUrl,
		"sys_type": in.SysType, "start_at": start, "end_at": end, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sOps) AnnUpdate(ctx context.Context, in service.AnnSaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"updated_at": gtime.Now()}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Content != "" {
		data["content"] = in.Content
	}
	if in.TextNode != "" {
		data["text_node"] = in.TextNode
	}
	if in.Cover != "" {
		data["cover"] = in.Cover
	}
	if in.JumpUrl != "" {
		data["jump_url"] = in.JumpUrl
	}
	if in.SysType != "" {
		data["sys_type"] = in.SysType
	}
	if in.StartAt != "" {
		if st := gtime.NewFromStr(in.StartAt); st != nil && !st.IsZero() {
			data["start_at"] = st
		}
	}
	if in.EndAt != "" {
		end := gtime.NewFromStr(in.EndAt)
		if end == nil || end.IsZero() {
			return gerror.New("结束时间格式非法")
		}
		data["end_at"] = end
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("announcement").Ctx(ctx).
		Where("site_id", opsSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sOps) AnnDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("announcement").Ctx(ctx).
		Where("site_id", opsSiteId).Where("id", id).Delete()
	return err
}

// ---------- 跳转位管理 ----------

func (s *sOps) JtList(ctx context.Context, f service.PageFilter) ([]*service.JtDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("jumptab").Ctx(ctx).Where("site_id", opsSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Location > 0 {
		m = m.Where("location", f.Location)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Jumptab
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.JtDTO, 0, len(list))
	for _, r := range list {
		out = append(out, jtDTO(r))
	}
	return out, total, nil
}

func (s *sOps) JtCreate(ctx context.Context, in service.JtSaveInput) (int64, error) {
	if strings.TrimSpace(in.CnName) == "" {
		return 0, gerror.New("名称不能为空")
	}
	if in.Location <= 0 {
		return 0, gerror.New("位置非法")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("jumptab").Ctx(ctx).Data(g.Map{
		"site_id": opsSiteId, "cn_name": in.CnName, "en_name": in.EnName,
		"avatar": in.Avatar, "link": in.Link, "pic_jump_link": in.PicJumpLink,
		"location": in.Location, "rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sOps) JtUpdate(ctx context.Context, in service.JtSaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"rank": in.Rank, "updated_at": gtime.Now()}
	if in.CnName != "" {
		data["cn_name"] = in.CnName
	}
	if in.EnName != "" {
		data["en_name"] = in.EnName
	}
	if in.Avatar != "" {
		data["avatar"] = in.Avatar
	}
	if in.Link != "" {
		data["link"] = in.Link
	}
	if in.PicJumpLink != "" {
		data["pic_jump_link"] = in.PicJumpLink
	}
	if in.Location > 0 {
		data["location"] = in.Location
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("jumptab").Ctx(ctx).
		Where("site_id", opsSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sOps) JtDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("jumptab").Ctx(ctx).
		Where("site_id", opsSiteId).Where("id", id).Delete()
	return err
}

// ---------- 敏感词管理 ----------

func (s *sOps) FwList(ctx context.Context, f service.PageFilter) ([]*service.FwDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 50
	}
	m := g.Model("filter_word").Ctx(ctx).Where("site_id", opsSiteId)
	if f.Keyword != "" {
		m = m.Where("word ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.FilterWord
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.FwDTO, 0, len(list))
	for _, r := range list {
		out = append(out, &service.FwDTO{Id: r.Id, Word: r.Word, CreatedAt: ts(r.CreatedAt)})
	}
	return out, total, nil
}

// FwAdd 批量添加(重复词 InsertIgnore 跳过), 返回实际新增数。
func (s *sOps) FwAdd(ctx context.Context, words []string) (int, error) {
	rows := make([]g.Map, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		rows = append(rows, g.Map{"site_id": opsSiteId, "word": w})
	}
	if len(rows) == 0 {
		return 0, gerror.New("有效词为空")
	}
	res, err := g.Model("filter_word").Ctx(ctx).Data(rows).InsertIgnore()
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return int(affected), nil
}

func (s *sOps) FwDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("filter_word").Ctx(ctx).
		Where("site_id", opsSiteId).Where("id", id).Delete()
	return err
}
