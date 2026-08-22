package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/group/service"
)

const groupSiteId = 1

type sGroup struct{}

func New() service.IGroup { return &sGroup{} }

func toDTO(r *entity.OfficialGroup) *service.ItemDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.ItemDTO{
		Id: r.Id, Name: r.Name, Intro: r.Intro, Avatar: r.Avatar,
		Link: r.Link, Platform: r.Platform, Rank: r.Rank, Status: r.Status,
		CreatedAt: created,
	}
}

func (s *sGroup) FrontList(ctx context.Context) ([]*service.ItemDTO, error) {
	var list []*entity.OfficialGroup
	if err := g.Model("official_group").Ctx(ctx).
		Where("site_id", groupSiteId).Where("status", 1).
		OrderDesc("rank").OrderDesc("id").Limit(100).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

func (s *sGroup) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("official_group").Ctx(ctx).Where("site_id", groupSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.OfficialGroup
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sGroup) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("社群名不能为空")
	}
	if in.Link == "" {
		return 0, gerror.New("跳转链接不能为空")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	return g.Model("official_group").Ctx(ctx).Data(g.Map{
		"site_id": groupSiteId, "name": in.Name, "intro": in.Intro,
		"avatar": in.Avatar, "link": in.Link, "platform": in.Platform,
		"rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
}

func (s *sGroup) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{
		"rank": in.Rank, "platform": in.Platform, "updated_at": gtime.Now(),
	}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Intro != "" {
		data["intro"] = in.Intro
	}
	if in.Avatar != "" {
		data["avatar"] = in.Avatar
	}
	if in.Link != "" {
		data["link"] = in.Link
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("official_group").Ctx(ctx).
		Where("site_id", groupSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sGroup) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("official_group").Ctx(ctx).
		Where("site_id", groupSiteId).Where("id", id).Delete()
	return err
}
