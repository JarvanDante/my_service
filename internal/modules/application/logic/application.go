// Package logic 推广应用业务(移植自 tianbi StatConfigApp/advertise click)。
package logic

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/application/service"
)

const appSiteId = 1 // 单站点样板

type sApplication struct{}

func New() service.IApplication { return &sApplication{} }

func toDTO(r *entity.Application) *service.AppDTO {
	locIds := []int64{}
	if r.LocIds != "" {
		_ = json.Unmarshal([]byte(r.LocIds), &locIds)
	}
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.AppDTO{
		Id: r.Id, Name: r.Name, Tag: r.Tag, Intro: r.Intro, Avatar: r.Avatar,
		DownloadUrl: r.DownloadUrl, IosUrl: r.IosUrl, AndroidUrl: r.AndroidUrl,
		LocIds: locIds, Rank: r.Rank, DownTotal: r.DownTotal, Status: r.Status,
		CreatedAt: created,
	}
}

func locJSON(ids []int64) string {
	if len(ids) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(ids)
	return string(b)
}

// FrontList 前台上架应用, loc>0 时按投放位置筛(jsonb 包含判断)。
func (s *sApplication) FrontList(ctx context.Context, loc int64) ([]*service.AppDTO, error) {
	m := g.Model("application").Ctx(ctx).
		Where("site_id", appSiteId).Where("status", 1)
	if loc > 0 {
		m = m.Where("loc_ids @> ?", locJSON([]int64{loc}))
	}
	var list []*entity.Application
	if err := m.OrderDesc("rank").OrderDesc("id").Limit(100).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.AppDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

// Click 下载点击计数(仅统计, 不做防刷; 前台组已有 IP 限流)。
func (s *sApplication) Click(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("应用ID非法")
	}
	_, err := g.Model("application").Ctx(ctx).
		Where("site_id", appSiteId).Where("id", id).
		Data(g.Map{"down_total": &gdb.Counter{Field: "down_total", Value: 1}}).Update()
	return err
}

func (s *sApplication) List(ctx context.Context, f service.ListFilter) ([]*service.AppDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("application").Ctx(ctx).Where("site_id", appSiteId)
	if f.Status >= 0 { // -1=全部
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.Application
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.AppDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sApplication) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("应用名不能为空")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("application").Ctx(ctx).Data(g.Map{
		"site_id": appSiteId, "name": in.Name, "tag": in.Tag, "intro": in.Intro,
		"avatar": in.Avatar, "download_url": in.DownloadUrl,
		"ios_url": in.IosUrl, "android_url": in.AndroidUrl,
		"loc_ids": locJSON(in.LocIds), "rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sApplication) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{
		"tag": in.Tag, "rank": in.Rank, "updated_at": gtime.Now(),
		"loc_ids": locJSON(in.LocIds),
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
	if in.DownloadUrl != "" {
		data["download_url"] = in.DownloadUrl
	}
	if in.IosUrl != "" {
		data["ios_url"] = in.IosUrl
	}
	if in.AndroidUrl != "" {
		data["android_url"] = in.AndroidUrl
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("application").Ctx(ctx).
		Where("site_id", appSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sApplication) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("application").Ctx(ctx).
		Where("site_id", appSiteId).Where("id", id).Delete()
	return err
}
