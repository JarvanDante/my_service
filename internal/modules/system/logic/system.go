// Package logic 系统模块业务实现(B7)。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/system/domain"
	"github.com/JarvanDante/my_service/internal/modules/system/service"
	"github.com/JarvanDante/my_service/internal/shared/siteconf"
)

type sSystem struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.ISystem { return &sSystem{repo: repo} }

// Push 发布系统公告/推送。
// TODO: type=push 时对接真实推送通道(APNs/FCM/极光), 当前仅落库。
func (s *sSystem) Push(ctx context.Context, in service.PushInput) (int64, error) {
	if in.Title == "" || in.Content == "" {
		return 0, gerror.New("标题与内容必填")
	}
	if in.Type != "notice" && in.Type != "push" {
		return 0, gerror.New("type 仅支持 notice/push")
	}
	return s.repo.NoticeCreate(ctx, &entity.SystemNotice{
		Title: in.Title, Content: in.Content, Type: in.Type, CreatedBy: in.OperatorId,
	})
}

func (s *sSystem) Notices(ctx context.Context, in service.NoticeListInput) (*service.NoticeListDTO, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.NoticeList(ctx, domain.NoticeFilter{Type: in.Type, Status: in.Status}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.NoticeDTO, 0, len(list))
	for _, n := range list {
		out = append(out, &service.NoticeDTO{
			Id: n.Id, Title: n.Title, Content: n.Content, Type: n.Type,
			Status: n.Status, CreatedBy: n.CreatedBy, CreatedAt: fmtTime(n.CreatedAt),
		})
	}
	return &service.NoticeListDTO{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sSystem) FrontNotices(ctx context.Context) ([]*service.NoticeDTO, error) {
	list, _, err := s.repo.NoticeList(ctx, domain.NoticeFilter{Type: "notice", Status: 1}, 1, 20)
	if err != nil {
		return nil, err
	}
	out := make([]*service.NoticeDTO, 0, len(list))
	for _, n := range list {
		out = append(out, &service.NoticeDTO{
			Id: n.Id, Title: n.Title, Content: n.Content, Type: n.Type,
			Status: n.Status, CreatedBy: n.CreatedBy, CreatedAt: fmtTime(n.CreatedAt),
		})
	}
	return out, nil
}

func (s *sSystem) SetNoticeStatus(ctx context.Context, id int64, status int) error {
	if id <= 0 {
		return gerror.New("公告ID无效")
	}
	if status != 0 && status != 1 {
		return gerror.New("status 仅支持 0/1")
	}
	n, err := s.repo.NoticeFind(ctx, id)
	if err != nil {
		return err
	}
	if n == nil {
		return gerror.New("公告不存在")
	}
	return s.repo.NoticeSetStatus(ctx, id, status)
}

func (s *sSystem) GetCustomerUrl(ctx context.Context) (string, error) {
	def := g.Cfg().MustGet(ctx, "customer.url", "https://example.com/kefu").String()
	return siteconf.Get(ctx, "customer_url", def), nil
}

func (s *sSystem) SetCustomerUrl(ctx context.Context, url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return gerror.New("链接必填")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return gerror.New("链接须以 http:// 或 https:// 开头")
	}
	return siteconf.Set(ctx, "customer_url", url)
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}
