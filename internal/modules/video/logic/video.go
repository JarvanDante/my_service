package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/video/domain"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
	"github.com/JarvanDante/my_service/internal/shared/paas"
)

func decodeTags(raw string) []string {
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

func parseCategories(raw string) []string {
	raw = strings.ReplaceAll(strings.TrimSpace(raw), "，", ",")
	if raw == "" {
		return []string{}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 4)
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func joinCategories(list []string) string {
	return strings.Join(parseCategories(strings.Join(list, ",")), ",")
}

func resolveCategory(in *service.SaveInput) string {
	if in.Categories != nil {
		return joinCategories(in.Categories)
	}
	return joinCategories(parseCategories(in.Category))
}

type sVideo struct{ repo domain.Repository }

func New(repo domain.Repository) service.IVideo { return &sVideo{repo: repo} }

func (s *sVideo) List(ctx context.Context, in service.ListInput) (*service.ListDTO, error) {
	page, size := in.Page, in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := s.repo.List(ctx, domain.ListFilter{
		Keyword: strings.TrimSpace(in.Keyword), MediaCode: strings.TrimSpace(in.MediaCode),
		Kind: normalizeKind(in.Kind), Status: in.Status,
	}, page, size)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	out := make([]*service.VideoDTO, 0, len(list))
	for _, v := range list {
		out = append(out, toDTO(ctx, v))
	}
	s.overlayMediaURLs(ctx, out)
	return &service.ListDTO{List: out, Total: total, Page: page, Size: size}, nil
}

// FrontList 前台列表。复用后台那套 repo.List, 只是把 Status 钉死成"已发布" ——
// 前台与后台的差异只有可见范围与排序口径, 没必要为此在 video 模块里再开一条 g.Model 直连。
func (s *sVideo) FrontList(ctx context.Context, in service.FrontListInput) (*service.ListDTO, error) {
	page, size := in.Page, in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 { // 前台是公开接口, 限一下单页上限, 免得被拿来整表拉取
		size = 100
	}
	tags := append([]string{}, in.Tags...)
	if t := strings.TrimSpace(in.Tag); t != "" {
		tags = append(tags, t)
	}
	list, total, err := s.repo.List(ctx, domain.ListFilter{
		Keyword:  strings.TrimSpace(in.Keyword),
		Category: strings.TrimSpace(in.Category),
		Tags:     tags,
		Kind:     normalizeKind(in.Kind),
		Status:   entity.VideoStatusPublished,
		Sort:     in.Sort,
	}, page, size)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	out := make([]*service.VideoDTO, 0, len(list))
	for _, v := range list {
		out = append(out, toDTO(ctx, v))
	}
	s.overlayMediaURLs(ctx, out)
	return &service.ListDTO{List: out, Total: total, Page: page, Size: size}, nil
}

// FrontDetail 前台详情。草稿/下架与不存在返回同一个错误, 避免前台通过错误文案
// 探测出"这个 id 有内容, 只是暂时下架"。
func (s *sVideo) FrontDetail(ctx context.Context, id int64) (*service.VideoDTO, error) {
	if id <= 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "ID必填")
	}
	v, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	if v == nil || v.Id == 0 || v.Status != entity.VideoStatusPublished {
		return nil, gerror.NewCode(gcode.CodeNotFound, "视频不存在或已下架")
	}
	d := toDTO(ctx, v)
	s.resolvePlay(ctx, d)
	return d, nil
}

func (s *sVideo) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	in.Category = resolveCategory(&in)
	if err := validateSave(in); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, &entity.Video{
		Title: strings.TrimSpace(in.Title), Description: strings.TrimSpace(in.Description),
		CoverUrl: in.CoverUrl, CoverKey: in.CoverKey, CoverMediaId: in.CoverMediaId,
		SourceUrl: in.SourceUrl, SourceKey: in.SourceKey, SourceMediaId: in.SourceMediaId,
		MediaCode: strings.TrimSpace(in.MediaCode), Kind: normalizeKind(in.Kind), Category: in.Category, Tags: encodeJSON(in.Tags),
		Duration: in.Duration, Sort: in.Sort, Status: in.Status, CreatedBy: in.OperatorId,
	})
}

func (s *sVideo) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "ID必填")
	}
	in.Category = resolveCategory(&in)
	if err := validateSave(in); err != nil {
		return err
	}
	old, err := s.repo.Find(ctx, in.Id)
	if err != nil {
		return gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	if old == nil || old.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "视频不存在")
	}
	return s.repo.Update(ctx, &entity.Video{
		Id: in.Id, Title: strings.TrimSpace(in.Title), Description: strings.TrimSpace(in.Description),
		CoverUrl: in.CoverUrl, CoverKey: in.CoverKey, CoverMediaId: in.CoverMediaId,
		SourceUrl: in.SourceUrl, SourceKey: in.SourceKey, SourceMediaId: in.SourceMediaId,
		MediaCode: strings.TrimSpace(in.MediaCode), Kind: old.Kind, Category: in.Category, Tags: encodeJSON(in.Tags),
		Duration: in.Duration, Sort: in.Sort, Status: in.Status,
	})
}

func (s *sVideo) Delete(ctx context.Context, id int64) error {
	old, err := s.repo.Find(ctx, id)
	if err != nil {
		return gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	if old == nil || old.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "视频不存在")
	}
	return s.repo.Delete(ctx, id)
}

func (s *sVideo) SetStatus(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "状态不合法")
	}
	old, err := s.repo.Find(ctx, id)
	if err != nil {
		return gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	if old == nil || old.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "视频不存在")
	}
	if status == entity.VideoStatusPublished && strings.TrimSpace(old.Category) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "请先编辑并选择本站分类后再上架")
	}
	return s.repo.SetStatus(ctx, id, status)
}

func validateSave(in service.SaveInput) error {
	if strings.TrimSpace(in.Title) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "标题必填")
	}
	if strings.TrimSpace(in.SourceUrl) == "" && strings.TrimSpace(in.SourceKey) == "" && strings.TrimSpace(in.MediaCode) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "请先上传视频或选用媒资")
	}
	if in.Status < 0 || in.Status > 2 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "状态不合法")
	}
	if in.Status == entity.VideoStatusPublished && in.Category == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "上架前请选择本站分类")
	}
	return nil
}

func toDTO(ctx context.Context, v *entity.Video) *service.VideoDTO {
	if v == nil {
		return nil
	}
	d := &service.VideoDTO{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, CoverKey: v.CoverKey, CoverMediaId: v.CoverMediaId,
		SourceUrl: v.SourceUrl, SourceKey: v.SourceKey, SourceMediaId: v.SourceMediaId,
		MediaCode: v.MediaCode, Category: v.Category, Categories: parseCategories(v.Category), Tags: decodeTags(v.Tags),
		Duration: v.Duration, Sort: v.Sort, Status: v.Status, CreatedBy: v.CreatedBy,
	}
	if v.CreatedAt != nil {
		d.CreatedAt = v.CreatedAt.String()
	}
	if v.UpdatedAt != nil {
		d.UpdatedAt = v.UpdatedAt.String()
	}
	d.CoverUrl, d.SourceUrl = paas.ApplyGatewayURLs(ctx, d.CoverUrl, d.SourceUrl, d.MediaCode)
	return d
}

func (s *sVideo) overlayMediaURLs(ctx context.Context, list []*service.VideoDTO) {
	need := false
	for _, d := range list {
		if d != nil && d.MediaCode != "" {
			need = true
			break
		}
	}
	if !need {
		return
	}
	picks, _, err := paas.ListPicks(ctx, 1, 200)
	if err != nil {
		picks, _, err = paas.ListAssets(ctx, 1, 200, "", 0)
	}
	m := make(map[string]paas.MediaAsset, len(picks))
	for _, a := range picks {
		m[a.Id] = a
	}
	for _, d := range list {
		if d == nil || d.MediaCode == "" {
			continue
		}
		if a, ok := m[d.MediaCode]; ok {
			if a.CoverUrl != "" && !isSiteCustomCover(d.CoverUrl) {
				d.CoverUrl = a.CoverUrl
			}
			if a.PlayUrl != "" {
				d.SourceUrl = a.PlayUrl
			}
		}
		d.CoverUrl, d.SourceUrl = paas.ApplyGatewayURLs(ctx, d.CoverUrl, d.SourceUrl, d.MediaCode)
	}
}

func (s *sVideo) resolvePlay(ctx context.Context, d *service.VideoDTO) {
	if d == nil || d.MediaCode == "" {
		return
	}
	if a, err := paas.AssetDetail(ctx, d.MediaCode); err == nil && a != nil {
		if a.CoverUrl != "" && !isSiteCustomCover(d.CoverUrl) {
			d.CoverUrl = a.CoverUrl
		}
		if a.PlayUrl != "" {
			d.SourceUrl = a.PlayUrl
		}
	}
	if url, err := paas.PlayToken(ctx, d.MediaCode); err == nil && url != "" {
		d.SourceUrl = url
	}
	d.CoverUrl, d.SourceUrl = paas.ApplyGatewayURLs(ctx, d.CoverUrl, d.SourceUrl, d.MediaCode)
}

func isSiteCustomCover(cover string) bool {
	cover = strings.TrimSpace(cover)
	if cover == "" {
		return false
	}
	if strings.Contains(cover, "/my-storage/") {
		return true
	}
	low := strings.ToLower(cover)
	if (strings.Contains(low, ".bnc") || strings.Contains(low, ".ceb")) && !strings.Contains(low, "/hls/") {
		return true
	}
	return false
}

func normalizeKind(kind int) int {
	if kind == entity.VideoKindCartoon {
		return entity.VideoKindCartoon
	}
	return entity.VideoKindVideo
}

func (s *sVideo) ListMediaAssets(ctx context.Context, page, size int, keyword string, kind int) ([]service.MediaAssetDTO, int, error) {
	kind = normalizeKind(kind)
	list, total, err := paas.ListAssets(ctx, page, size, strings.TrimSpace(keyword), kind)
	if err != nil {
		return nil, 0, err
	}
	out := make([]service.MediaAssetDTO, 0, len(list))
	for _, a := range list {
		cover, play := paas.ApplyGatewayURLs(ctx, a.CoverUrl, a.PlayUrl, a.Id)
		item := service.MediaAssetDTO{
			Id: a.Id, Title: a.Title, CoverUrl: cover, PlayUrl: play,
			DurationSec: a.DurationSec, Picked: a.Picked,
		}
		if local, _ := s.repo.FindByMediaCode(ctx, a.Id, kind); local != nil {
			item.LocalId = local.Id
		}
		out = append(out, item)
	}
	return out, total, nil
}

func (s *sVideo) PickMedia(ctx context.Context, code string, operatorId int64, kind int) (int64, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "媒资ID必填")
	}
	a, err := paas.PickAsset(ctx, code)
	if err != nil {
		return 0, err
	}
	return s.upsertFromAsset(ctx, a, operatorId, normalizeKind(kind))
}

func (s *sVideo) SyncMedia(ctx context.Context, operatorId int64, kind int) (*service.SyncMediaDTO, error) {
	kind = normalizeKind(kind)
	created, updated, total := 0, 0, 0
	for page := 1; ; page++ {
		list, cnt, err := paas.ListAssets(ctx, page, 50, "", kind)
		if err != nil {
			return nil, err
		}
		if page == 1 {
			total = cnt
		}
		if len(list) == 0 {
			break
		}
		for _, a := range list {
			picked, err := paas.PickAsset(ctx, a.Id)
			if err != nil || picked == nil {
				picked = &a
			}
			existed, _ := s.repo.FindByMediaCode(ctx, a.Id, kind)
			if _, err = s.upsertFromAsset(ctx, picked, operatorId, kind); err != nil {
				return nil, err
			}
			if existed != nil && existed.Id > 0 {
				updated++
			} else {
				created++
			}
		}
		if page*50 >= cnt {
			break
		}
	}
	return &service.SyncMediaDTO{Created: created, Updated: updated, Total: total}, nil
}

func (s *sVideo) upsertFromAsset(ctx context.Context, a *paas.MediaAsset, operatorId int64, kind int) (int64, error) {
	if a == nil || a.Id == "" {
		return 0, gerror.New("媒资无效")
	}
	kind = normalizeKind(kind)
	old, err := s.repo.FindByMediaCode(ctx, a.Id, kind)
	if err != nil {
		return 0, err
	}
	title := strings.TrimSpace(a.Title)
	if title == "" {
		title = a.Id
	}
	if old != nil && old.Id > 0 {
		old.Title = title
		old.CoverUrl = a.CoverUrl
		old.SourceUrl = a.PlayUrl
		old.SourceKey = a.PlayKey
		old.MediaCode = a.Id
		if a.DurationSec > 0 {
			old.Duration = a.DurationSec
		}
		return old.Id, s.repo.Update(ctx, old)
	}
	return s.repo.Create(ctx, &entity.Video{
		Title: title, CoverUrl: a.CoverUrl, SourceUrl: a.PlayUrl, SourceKey: a.PlayKey,
		MediaCode: a.Id, Kind: kind, Category: "", Tags: "[]", Duration: a.DurationSec,
		Status: entity.VideoStatusDraft, CreatedBy: operatorId,
	})
}
