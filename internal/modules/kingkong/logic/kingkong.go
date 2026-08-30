package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/kingkong/service"
)

const siteId = 1

var positionNames = map[string]string{
	entity.KingkongPosComics:  "漫画",
	entity.KingkongPosCartoon: "动漫",
	entity.KingkongPosMovie:   "视频",
	entity.KingkongPosNovel:   "小说",
	entity.KingkongPosShort:   "短剧",
}

var openModeNames = map[string]string{
	entity.KingkongModeBlock:  "模块",
	entity.KingkongModeList:   "列表",
	entity.KingkongModeDouyin: "抖音流",
}

var linkPresets = map[string]string{
	"activityLand": "活动专区",
	"selected":     "精选",
	"day":          "每日",
	"checkin":      "签到",
	"invite":       "邀请",
	"collect":      "收藏",
	"submission":   "投稿",
	"aiExperience": "AI体验馆",
	"huangyou":     "黄油",
	"douyin":       "抖音",
	"vipUpgrade":   "补差价升级",
	"memberCenter": "会员中心",
}

type sKingkong struct{}

func New() service.IKingkong { return &sKingkong{} }

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func normalizePosition(p string) (string, error) {
	p = strings.TrimSpace(p)
	if _, ok := positionNames[p]; !ok {
		return "", gerror.New("展示位置非法")
	}
	return p, nil
}

func normalizeOpenMode(m string) string {
	m = strings.TrimSpace(m)
	if _, ok := openModeNames[m]; ok {
		return m
	}
	return entity.KingkongModeBlock
}

func linkLabel(link, appLink string) string {
	link = strings.TrimSpace(link)
	appLink = strings.TrimSpace(appLink)
	if link == "" && appLink == "" {
		return "不跳转"
	}
	parts := make([]string, 0, 2)
	if link != "" {
		if name, ok := linkPresets[link]; ok {
			parts = append(parts, link+" | "+name)
		} else {
			parts = append(parts, link)
		}
	}
	if appLink != "" {
		parts = append(parts, "app: "+appLink)
	}
	return strings.Join(parts, "；")
}

func toDTO(r *entity.KingkongItem) *service.ItemDTO {
	statusText := "禁用"
	if r.Status == 1 {
		statusText = "正常"
	}
	return &service.ItemDTO{
		Id:           r.Id,
		Name:         r.Name,
		IconUrl:      r.IconUrl,
		OpenMode:     r.OpenMode,
		OpenModeName: openModeNames[r.OpenMode],
		Link:         r.Link,
		AppLink:      r.AppLink,
		LinkLabel:    linkLabel(r.Link, r.AppLink),
		Position:     r.Position,
		PositionName: positionNames[r.Position],
		Sort:         r.Sort,
		Status:       r.Status,
		StatusText:   statusText,
		CreatedAt:    fmtTime(r.CreatedAt),
		UpdatedAt:    fmtTime(r.UpdatedAt),
	}
}

func (s *sKingkong) FrontList(ctx context.Context, position string) ([]*service.ItemDTO, error) {
	pos, err := normalizePosition(position)
	if err != nil {
		return nil, err
	}
	var list []*entity.KingkongItem
	if err := g.Model("kingkong_item").Ctx(ctx).
		Where("site_id", siteId).
		Where("position", pos).
		Where("status", 1).
		OrderDesc("sort").OrderDesc("id").
		Limit(20).
		Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

func (s *sKingkong) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("kingkong_item").Ctx(ctx).Where("site_id", siteId)
	if name := strings.TrimSpace(f.Name); name != "" {
		m = m.Where("name ILIKE ?", "%"+name+"%")
	}
	if pos := strings.TrimSpace(f.Position); pos != "" {
		m = m.Where("position", pos)
	}
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.KingkongItem
	if err := m.Clone().OrderDesc("sort").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sKingkong) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, gerror.New("名称不能为空")
	}
	icon := strings.TrimSpace(in.IconUrl)
	if icon == "" {
		return 0, gerror.New("请上传图标")
	}
	pos, err := normalizePosition(in.Position)
	if err != nil {
		return 0, err
	}
	status := in.Status
	if status != 0 && status != 1 {
		status = 1
	}
	id, err := g.Model("kingkong_item").Ctx(ctx).Data(g.Map{
		"site_id":   siteId,
		"name":      name,
		"icon_url":  icon,
		"open_mode": normalizeOpenMode(in.OpenMode),
		"link":      strings.TrimSpace(in.Link),
		"app_link":  strings.TrimSpace(in.AppLink),
		"position":  pos,
		"sort":      in.Sort,
		"status":    status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sKingkong) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return gerror.New("名称不能为空")
	}
	icon := strings.TrimSpace(in.IconUrl)
	if icon == "" {
		return gerror.New("请上传图标")
	}
	pos, err := normalizePosition(in.Position)
	if err != nil {
		return err
	}
	status := in.Status
	if status != 0 && status != 1 {
		status = 1
	}
	_, err = g.Model("kingkong_item").Ctx(ctx).
		Where("site_id", siteId).Where("id", in.Id).
		Data(g.Map{
			"name":       name,
			"icon_url":   icon,
			"open_mode":  normalizeOpenMode(in.OpenMode),
			"link":       strings.TrimSpace(in.Link),
			"app_link":   strings.TrimSpace(in.AppLink),
			"position":   pos,
			"sort":       in.Sort,
			"status":     status,
			"updated_at": gtime.Now(),
		}).Update()
	return err
}

func (s *sKingkong) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("kingkong_item").Ctx(ctx).
		Where("site_id", siteId).Where("id", id).Delete()
	return err
}
