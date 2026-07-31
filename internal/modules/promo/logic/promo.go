// Package logic 推广/兑换码模块业务实现(B3)。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/promo/domain"
	"github.com/JarvanDante/my_service/internal/modules/promo/service"
)

type sPromo struct {
	repo domain.Repository
}

func New(repo domain.Repository) service.IPromo { return &sPromo{repo: repo} }

func (s *sPromo) Codes(ctx context.Context, in service.CodeListInput) (*service.PageDTO[service.CodeDTO], error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.ListCodes(ctx, domain.CodeFilter{
		Keyword: in.Keyword, CodeKey: in.CodeKey, Type: in.Type, Status: in.Status,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.CodeDTO, 0, len(list))
	for _, c := range list {
		out = append(out, &service.CodeDTO{
			Id: c.Id, Name: c.Name, Code: c.Code, CodeKey: c.CodeKey,
			Type: c.Type, ObjectId: c.ObjectId, AddNum: c.AddNum,
			CanUseNum: c.CanUseNum, UsedNum: c.UsedNum, Status: c.Status,
			ExpiredAt: c.ExpiredAt, CreatedAt: fmtTime(c.CreatedAt),
		})
	}
	return &service.PageDTO[service.CodeDTO]{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

// GenerateCodes 批量生成: 一个批次共用 code_key, 码为 16 位大写字母数字。
func (s *sPromo) GenerateCodes(ctx context.Context, in service.CodeGenInput) (*service.CodeGenDTO, error) {
	if in.Name == "" {
		return nil, gerror.New("名称必填")
	}
	if in.Type != "point" && in.Type != "group" {
		return nil, gerror.New("type 仅支持 point/group")
	}
	if in.Type == "group" && in.ObjectId <= 0 {
		return nil, gerror.New("type=group 时必须指定用户组 object_id")
	}
	if in.AddNum <= 0 {
		return nil, gerror.New("面额(add_num)必须大于0")
	}
	if in.Count <= 0 || in.Count > 1000 {
		return nil, gerror.New("生成数量须在 1~1000")
	}
	if in.CanUseNum <= 0 {
		in.CanUseNum = 1
	}
	if in.ExpiredAt > 0 && in.ExpiredAt < gtime.Timestamp() {
		return nil, gerror.New("过期时间不能早于当前时间")
	}
	codeKey := "B" + gtime.Now().Format("YmdHis") + grand.Digits(4)
	rows := make([]*entity.UserCode, 0, in.Count)
	codes := make([]string, 0, in.Count)
	seen := make(map[string]struct{}, in.Count)
	for len(rows) < in.Count {
		code := strings.ToUpper(grand.Letters(4) + grand.Digits(4) + grand.Letters(4) + grand.Digits(4))
		if _, dup := seen[code]; dup {
			continue
		}
		seen[code] = struct{}{}
		rows = append(rows, &entity.UserCode{
			Name: in.Name, Code: code, CodeKey: codeKey, Type: in.Type,
			ObjectId: in.ObjectId, AddNum: in.AddNum, CanUseNum: in.CanUseNum,
			ExpiredAt: in.ExpiredAt,
		})
		codes = append(codes, code)
	}
	if err := s.repo.BatchCreateCodes(ctx, rows); err != nil {
		return nil, err
	}
	return &service.CodeGenDTO{CodeKey: codeKey, Count: in.Count, Codes: codes}, nil
}

func (s *sPromo) VoidCode(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("兑换码ID无效")
	}
	c, err := s.repo.FindCodeById(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return gerror.New("兑换码不存在")
	}
	if c.Status == -1 {
		return gerror.New("兑换码已作废")
	}
	return s.repo.VoidCode(ctx, id)
}

func (s *sPromo) CodeLogs(ctx context.Context, in service.CodeLogListInput) (*service.PageDTO[service.CodeLogDTO], error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.ListCodeLogs(ctx, domain.CodeLogFilter{
		CodeId: in.CodeId, UserId: in.UserId, Code: in.Code,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.CodeLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.CodeLogDTO{
			Id: l.Id, CodeId: l.CodeId, Code: l.Code, Name: l.Name, Type: l.Type,
			UserId: l.UserId, Username: l.Username, AddNum: l.AddNum,
			CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return &service.PageDTO[service.CodeLogDTO]{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

// ==================== B6 分享 / 拉新 ====================

func (s *sPromo) ShareLogs(ctx context.Context, in service.ShareLogListInput) (*service.PageDTO[service.ShareLogAdminDTO], error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.ShareLogList(ctx, domain.ShareLogFilter{
		UserId: in.UserId, Type: in.Type, Channel: in.Channel,
		StartDate: in.StartDate, EndDate: in.EndDate,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	out := make([]*service.ShareLogAdminDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.ShareLogAdminDTO{
			Id: l.Id, UserId: l.UserId, Type: l.Type, TargetId: l.TargetId,
			Channel: l.Channel, CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return &service.PageDTO[service.ShareLogAdminDTO]{List: out, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sPromo) ShareStats(ctx context.Context, startDate, endDate string, top int) (*service.ShareStatsDTO, error) {
	totalShares, sharerCount, channels, err := s.repo.ShareStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	rank, err := s.repo.InviteRank(ctx, startDate, endDate, top)
	if err != nil {
		return nil, err
	}
	chs := make([]service.ChannelCountDTO, 0, len(channels))
	for _, c := range channels {
		chs = append(chs, service.ChannelCountDTO{Channel: c.Channel, Count: c.Count})
	}
	ranks := make([]service.InviteRankDTO, 0, len(rank))
	for _, r := range rank {
		ranks = append(ranks, service.InviteRankDTO{
			UserId: r.UserId, Username: r.Username, InviteCount: r.InviteCount,
		})
	}
	return &service.ShareStatsDTO{
		TotalShares: totalShares, SharerCount: sharerCount, Channels: chs, InviteRank: ranks,
	}, nil
}
