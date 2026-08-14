// Package logic — B1 后台用户管理实现。
package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

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
		Keyword: in.Keyword, UserId: in.UserId, Username: in.Username, Phone: in.Phone,
		ParentId: in.ParentId, Channel: in.Channel, GroupId: in.GroupId,
		IsUp: in.IsUp, IsValid: in.IsValid, HasBuy: in.HasBuy, Status: in.Status,
		DeviceType: in.DeviceType, StartDate: in.StartDate, EndDate: in.EndDate,
		MinLoginNum: in.MinLoginNum, MaxLoginNum: in.MaxLoginNum,
	}, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	appidCache := map[int64]string{}
	items := make([]*service.AdminUserItemDTO, 0, len(list))
	for _, u := range list {
		items = append(items, toAdminItem(ctx, u, appidCache))
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
		AdminUserItemDTO: *toAdminItem(ctx, u, map[int64]string{}),
		Signature:        u.Signature,
		Fans:             u.Fans,
		Follow:           u.Follow,
		ErrorMsg:         u.ErrorMsg,
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

// AdminBatchSetDisabled 批量冻结/解冻, 返回实际变更人数。
func (s *sUser) AdminBatchSetDisabled(ctx context.Context, ids []int64, disable bool, reason string) (int, error) {
	if len(ids) == 0 {
		return 0, gerror.New("用户ID列表不能为空")
	}
	if len(ids) > 500 {
		return 0, gerror.New("单次最多操作500个用户")
	}
	flag := 0
	if disable {
		flag = 1
		if reason == "" {
			reason = "后台批量冻结"
		}
	}
	n, err := s.repo.SetDisabledBatch(ctx, ids, flag, reason)
	return int(n), err
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
	if in.GroupId == 0 {
		// 移出用户组: 清空快照
		return s.repo.UpdateGroup(ctx, in.UserId, 0, "", 0, 0)
	}
	// B4: 未显式给名称时, 从用户组定义表回填名称/折扣
	if in.GroupName == "" {
		ug, err := s.repo.GroupFind(ctx, in.GroupId)
		if err != nil {
			return err
		}
		if ug == nil {
			return gerror.New("用户组不存在, 请先在用户组配置中创建")
		}
		if ug.Status != 1 {
			return gerror.New("用户组已停用")
		}
		in.GroupName = ug.Name
		in.GroupRate = ug.Rate
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

func toAdminItem(ctx context.Context, u *entity.Users, appidCache map[int64]string) *service.AdminUserItemDTO {
	return &service.AdminUserItemDTO{
		Id: u.Id, Username: u.Username, Nickname: u.Nickname, Phone: u.Phone,
		Sex: u.Sex, Tag: u.Tag, Img: u.Img,
		AccountSlat: buildAccountSlat(ctx, u.Username, u.SiteId, appidCache),
		Balance:     u.Balance, GiftCount: u.GiftCount, Credit: u.Credit, MoneyCount: u.MoneyCount,
		IsUp: u.IsUp, IsValid: u.IsValid, HasBuy: u.HasBuy, Level: u.Level,
		GroupId: u.GroupId, GroupName: u.GroupName, GroupRate: u.GroupRate,
		GroupStartTime: u.GroupStartTime, GroupEndTime: u.GroupEndTime,
		ParentId: u.ParentId, ParentName: u.ParentName, Channel: u.ChannelName,
		DeviceType: u.DeviceType, DeviceExt: u.DeviceExt, DeviceVersion: u.DeviceVersion,
		MovieFeeRate: u.MovieFeeRate, PostFeeRate: u.PostFeeRate, ShareNum: u.ShareNum,
		IsDisabled: u.IsDisabled,
		RegisterAt: fmtTime(u.RegisterAt), RegisterIp: u.RegisterIp, RegisterArea: u.RegisterArea,
		LastLoginAt: fmtTime(u.LastLoginAt), LastIp: u.LastIp, LoginNum: u.LoginNum,
	}
}

// buildAccountSlat 与公司后台一致: username==>md5(username_appid)。
func buildAccountSlat(ctx context.Context, username string, siteId int64, cache map[int64]string) string {
	if username == "" {
		return ""
	}
	appid := siteAppid(ctx, siteId, cache)
	sum := md5.Sum([]byte(username + "_" + appid))
	return username + "==>" + hex.EncodeToString(sum[:])
}

var defaultAppidOnce sync.Once
var defaultAppid string

func siteAppid(ctx context.Context, siteId int64, cache map[int64]string) string {
	if siteId <= 0 {
		siteId = 1
	}
	if cache != nil {
		if v, ok := cache[siteId]; ok {
			return v
		}
	}
	v, err := g.Model("site").Ctx(ctx).Where("id", siteId).Value("appid")
	appid := ""
	if err == nil && v != nil {
		appid = v.String()
	}
	if appid == "" {
		defaultAppidOnce.Do(func() {
			dv, _ := g.Model("site").Ctx(ctx).OrderAsc("id").Value("appid")
			if dv != nil {
				defaultAppid = dv.String()
			}
		})
		appid = defaultAppid
	}
	if cache != nil {
		cache[siteId] = appid
	}
	return appid
}
