package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/video/domain"
	"github.com/JarvanDante/my_service/internal/modules/video/service"
)

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
		Keyword: strings.TrimSpace(in.Keyword), Status: in.Status,
	}, page, size)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	out := make([]*service.VideoDTO, 0, len(list))
	for _, v := range list {
		out = append(out, toDTO(v))
	}
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
	list, total, err := s.repo.List(ctx, domain.ListFilter{
		Keyword: strings.TrimSpace(in.Keyword),
		Status:  entity.VideoStatusPublished,
		Sort:    in.Sort,
	}, page, size)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeDbOperationError, err, "查询视频失败")
	}
	out := make([]*service.VideoDTO, 0, len(list))
	for _, v := range list {
		out = append(out, toDTO(v))
	}
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
	return toDTO(v), nil
}

func (s *sVideo) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if err := validateSave(in); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, &entity.Video{
		Title: strings.TrimSpace(in.Title), Description: strings.TrimSpace(in.Description),
		CoverUrl: in.CoverUrl, CoverKey: in.CoverKey, CoverMediaId: in.CoverMediaId,
		SourceUrl: in.SourceUrl, SourceKey: in.SourceKey, SourceMediaId: in.SourceMediaId,
		Duration: in.Duration, Sort: in.Sort, Status: in.Status, CreatedBy: in.OperatorId,
	})
}

func (s *sVideo) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "ID必填")
	}
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
	return s.repo.SetStatus(ctx, id, status)
}

func validateSave(in service.SaveInput) error {
	if strings.TrimSpace(in.Title) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "标题必填")
	}
	if strings.TrimSpace(in.SourceUrl) == "" && strings.TrimSpace(in.SourceKey) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "请先上传视频")
	}
	if in.Status < 0 || in.Status > 2 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "状态不合法")
	}
	return nil
}

func toDTO(v *entity.Video) *service.VideoDTO {
	if v == nil {
		return nil
	}
	d := &service.VideoDTO{
		Id: v.Id, Title: v.Title, Description: v.Description,
		CoverUrl: v.CoverUrl, CoverKey: v.CoverKey, CoverMediaId: v.CoverMediaId,
		SourceUrl: v.SourceUrl, SourceKey: v.SourceKey, SourceMediaId: v.SourceMediaId,
		Duration: v.Duration, Sort: v.Sort, Status: v.Status, CreatedBy: v.CreatedBy,
	}
	if v.CreatedAt != nil {
		d.CreatedAt = v.CreatedAt.String()
	}
	if v.UpdatedAt != nil {
		d.UpdatedAt = v.UpdatedAt.String()
	}
	return d
}
