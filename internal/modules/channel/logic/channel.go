package logic

import (
	"context"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/channel/service"
	"github.com/JarvanDante/my_service/internal/shared/appcfg"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

const channelSiteId = 1

type sChannel struct{}

func New() service.IChannel { return &sChannel{} }

func statusText(status int) string {
	if status == 1 {
		return "启用"
	}
	return "停用"
}

func buildH5Link(ctx context.Context, code string) string {
	base := strings.TrimSpace(appcfg.String(ctx, "share_url", ""))
	if strings.Contains(strings.ToLower(base), "example.com") {
		base = ""
	}
	base = strings.TrimRight(base, "/")
	if base == "" {
		return "?source=" + url.QueryEscape(code)
	}
	u, err := url.Parse(base)
	if err != nil {
		return base + "?source=" + url.QueryEscape(code)
	}
	q := u.Query()
	q.Set("source", code)
	q.Del("invite")
	u.RawQuery = q.Encode()
	return u.String()
}

func toDTO(ctx context.Context, r *entity.PromoChannel, userCount int) *service.ItemDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.ItemDTO{
		Id:             r.Id,
		Code:           r.Code,
		Name:           r.Name,
		Remark:         r.Remark,
		Status:         r.Status,
		StatusText:     statusText(r.Status),
		UserCount:      userCount,
		H5Link:         buildH5Link(ctx, r.Code),
		ClipboardValue: kit.SourceClipboard(r.Code),
		CreatedAt:      created,
	}
}

func userCountByCode(ctx context.Context, codes []string) (map[string]int, error) {
	out := make(map[string]int, len(codes))
	if len(codes) == 0 {
		return out, nil
	}
	type row struct {
		Code  string `json:"code"`
		Count int    `json:"count"`
	}
	var list []row
	if err := g.Model("users").Ctx(ctx).
		Fields("channel_name as code, count(*) as count").
		WhereIn("channel_name", codes).
		Group("channel_name").
		Scan(&list); err != nil {
		return nil, err
	}
	for _, r := range list {
		out[r.Code] = r.Count
	}
	return out, nil
}

func (s *sChannel) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("promo_channel").Ctx(ctx).Where("site_id", channelSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if kw := strings.TrimSpace(f.Keyword); kw != "" {
		m = m.Where("(code ILIKE ? OR name ILIKE ?)", "%"+kw+"%", "%"+kw+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.PromoChannel
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	codes := make([]string, 0, len(list))
	for _, r := range list {
		codes = append(codes, r.Code)
	}
	counts, err := userCountByCode(ctx, codes)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(ctx, r, counts[r.Code]))
	}
	return out, total, nil
}

func (s *sChannel) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	code := kit.ParseSource(in.Code)
	if code == "" {
		return 0, gerror.New("渠道码须为 2-32 位字母数字或 _-")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, gerror.New("渠道名必填")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	exist, err := g.Model("promo_channel").Ctx(ctx).
		Where("site_id", channelSiteId).
		Where("lower(code) = ?", strings.ToLower(code)).
		Count()
	if err != nil {
		return 0, err
	}
	if exist > 0 {
		return 0, gerror.New("渠道码已存在")
	}
	return g.Model("promo_channel").Ctx(ctx).Data(g.Map{
		"site_id": channelSiteId,
		"code":    code,
		"name":    name,
		"remark":  strings.TrimSpace(in.Remark),
		"status":  in.Status,
	}).InsertAndGetId()
}

func (s *sChannel) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return gerror.New("渠道名必填")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	n, err := g.Model("promo_channel").Ctx(ctx).
		Where("site_id", channelSiteId).Where("id", in.Id).
		Data(g.Map{
			"name":       name,
			"remark":     strings.TrimSpace(in.Remark),
			"status":     in.Status,
			"updated_at": gtime.Now(),
		}).UpdateAndGetAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return gerror.New("渠道不存在")
	}
	return nil
}

func (s *sChannel) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	var row *entity.PromoChannel
	if err := g.Model("promo_channel").Ctx(ctx).
		Where("site_id", channelSiteId).Where("id", id).Scan(&row); err != nil {
		return err
	}
	if row == nil {
		return gerror.New("渠道不存在")
	}
	used, err := g.Model("users").Ctx(ctx).Where("channel_name", row.Code).Count()
	if err != nil {
		return err
	}
	if used > 0 {
		return gerror.New("已有用户归属该渠道，请停用而不是删除")
	}
	_, err = g.Model("promo_channel").Ctx(ctx).
		Where("site_id", channelSiteId).Where("id", id).Delete()
	return err
}
