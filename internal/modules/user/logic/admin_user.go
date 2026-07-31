// Package logic — B1 后台用户管理实现。
package logic

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

func (s *sUser) AdminListUsers(ctx context.Context, in service.AdminUserListInput) (*service.AdminUserListDTO, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	list, total, err := s.repo.AdminListUsers(ctx, domain.AdminUserFilter{
		Keyword: in.Keyword, Channel: in.Channel, GroupId: in.GroupId,
		Status: in.Status, StartDate: in.StartDate, EndDate: in.EndDate,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*service.AdminUserItemDTO, 0, len(list))
	for _, u := range list {
		items = append(items, toAdminItem(u))
	}
	return &service.AdminUserListDTO{List: items, Total: total, Page: in.Page, Size: in.Size}, nil
}

func (s *sUser) AdminUserDetail(ctx context.Context, id int64) (*service.AdminUserDetailDTO, error) {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户不存在")
	}
	d := &service.AdminUserDetailDTO{
		AdminUserItemDTO: *toAdminItem(u),
		Sex:              u.Sex,
		Signature:        u.Signature,
		Img:              u.Img,
		Fans:             u.Fans,
		Follow:           u.Follow,
		ShareNum:         u.ShareNum,
		ParentId:         u.ParentId,
		ParentName:       u.ParentName,
		GroupRate:        u.GroupRate,
		GroupEndTime:     u.GroupEndTime,
		ErrorMsg:         u.ErrorMsg,
		RegisterIp:       u.RegisterIp,
		LastIp:           u.LastIp,
		LoginNum:         u.LoginNum,
	}
	return d, nil
}

func (s *sUser) AdminSetDisabled(ctx context.Context, id int64, disable bool, reason string) error {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return gerror.New("用户不存在")
	}
	flag := 0
	if disable {
		flag = 1
		if reason == "" {
			reason = "后台禁用"
		}
	}
	return s.repo.SetDisabled(ctx, id, flag, reason)
}

func (s *sUser) AdminSetGroup(ctx context.Context, in service.AdminSetGroupInput) error {
	if in.UserId <= 0 {
		return gerror.New("用户ID无效")
	}
	u, err := s.repo.FindById(ctx, in.UserId)
	if err != nil {
		return err
	}
	if u == nil {
		return gerror.New("用户不存在")
	}
	if in.GroupId < 0 || in.GroupRate < 0 {
		return gerror.New("参数不合法")
	}
	return s.repo.UpdateGroup(ctx, in.UserId, in.GroupId, in.GroupName, in.GroupRate, in.GroupEndTime)
}

func (s *sUser) AdminAdjustBalance(ctx context.Context, in service.AdminAdjustBalanceInput) error {
	if in.UserId <= 0 {
		return gerror.New("用户ID无效")
	}
	if in.Target != "balance" && in.Target != "credit" {
		return gerror.New("target 仅支持 balance/credit")
	}
	if in.Amount == 0 {
		return gerror.New("调整数额不能为 0")
	}
	refId := fmt.Sprintf("admin:%d", in.OperatorId)
	remark := in.Remark
	if remark == "" {
		remark = "后台调整"
	}
	return s.repo.AdminAdjustBalance(ctx, in.UserId, in.Target, in.Amount, refId, remark)
}

func (s *sUser) AdminBalanceLogs(ctx context.Context, userId int64, page, size int) ([]*service.BalanceLogDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := s.repo.BalanceLogs(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.BalanceLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.BalanceLogDTO{
			Id: l.Id, Direction: l.Direction, Scene: l.Scene, Amount: l.Amount,
			BalanceBefore: l.BalanceBefore, BalanceAfter: l.BalanceAfter,
			RefId: l.RefId, Remark: l.Remark, CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

func toAdminItem(u *entity.Users) *service.AdminUserItemDTO {
	return &service.AdminUserItemDTO{
		Id: u.Id, Username: u.Username, Nickname: u.Nickname, Phone: u.Phone,
		Channel: u.ChannelName, GroupId: u.GroupId, GroupName: u.GroupName,
		Level: u.Level, Balance: u.Balance, Credit: u.Credit, MoneyCount: u.MoneyCount,
		IsDisabled: u.IsDisabled,
		RegisterAt: fmtTime(u.RegisterAt), LastLoginAt: fmtTime(u.LastLoginAt),
	}
}
