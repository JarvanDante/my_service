// Package backend 后台用户控制器(B1)。
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ user service.IUser }

func New(svc service.IUser) *Controller { return &Controller{user: svc} }

// operatorId 从 ctx 取当前管理员ID(AdminAuth 写入)。
func operatorId(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxAdminId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "管理员未登录")
	}
	return id, nil
}

// List 用户列表。
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	dto, err := c.user.AdminListUsers(ctx, service.AdminUserListInput{
		Keyword: req.Keyword, Channel: req.Channel, GroupId: req.GroupId,
		Status: req.Status, StartDate: req.StartDate, EndDate: req.EndDate,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]v1.AdminUserItem, 0, len(dto.List))
	for _, u := range dto.List {
		list = append(list, toItem(u))
	}
	return &v1.ListRes{List: list, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

// Detail 用户详情。
func (c *Controller) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	d, err := c.user.AdminUserDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DetailRes{
		AdminUserItem: toItem(&d.AdminUserItemDTO),
		Sex:           d.Sex, Signature: d.Signature, Img: d.Img,
		Fans: d.Fans, Follow: d.Follow, ShareNum: d.ShareNum,
		ParentId: d.ParentId, ParentName: d.ParentName,
		GroupRate: d.GroupRate, GroupEndTime: d.GroupEndTime,
		ErrorMsg: d.ErrorMsg, RegisterIp: d.RegisterIp, LastIp: d.LastIp, LoginNum: d.LoginNum,
	}, nil
}

// Disable 禁用/解禁。
func (c *Controller) Disable(ctx context.Context, req *v1.DisableReq) (res *v1.DisableRes, err error) {
	if err = c.user.AdminSetDisabled(ctx, req.Id, req.Op == "disable", req.Reason); err != nil {
		return nil, err
	}
	return &v1.DisableRes{}, nil
}

// Group 调整用户组。
func (c *Controller) Group(ctx context.Context, req *v1.GroupReq) (res *v1.GroupRes, err error) {
	if err = c.user.AdminSetGroup(ctx, service.AdminSetGroupInput{
		UserId: req.Id, GroupId: req.GroupId, GroupName: req.GroupName,
		GroupRate: req.GroupRate, GroupEndTime: req.GroupEndTime,
	}); err != nil {
		return nil, err
	}
	return &v1.GroupRes{}, nil
}

// Balance 调整金币/积分。
func (c *Controller) Balance(ctx context.Context, req *v1.BalanceReq) (res *v1.BalanceRes, err error) {
	op, err := operatorId(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.AdminAdjustBalance(ctx, service.AdminAdjustBalanceInput{
		UserId: req.Id, Target: req.Target, Amount: req.Amount,
		Remark: req.Remark, OperatorId: op,
	}); err != nil {
		return nil, err
	}
	return &v1.BalanceRes{}, nil
}

// BalanceLogs 余额流水。
func (c *Controller) BalanceLogs(ctx context.Context, req *v1.BalanceLogsReq) (res *v1.BalanceLogsRes, err error) {
	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := c.user.AdminBalanceLogs(ctx, req.Id, page, size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.BalanceLogItem, 0, len(list))
	for _, l := range list {
		items = append(items, v1.BalanceLogItem{
			Id: l.Id, Direction: l.Direction, Scene: l.Scene, Amount: l.Amount,
			BalanceBefore: l.BalanceBefore, BalanceAfter: l.BalanceAfter,
			RefId: l.RefId, Remark: l.Remark, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.BalanceLogsRes{List: items, Total: total, Page: page, Size: size}, nil
}

func toItem(u *service.AdminUserItemDTO) v1.AdminUserItem {
	return v1.AdminUserItem{
		Id: u.Id, Username: u.Username, Nickname: u.Nickname, Phone: u.Phone,
		Channel: u.Channel, GroupId: u.GroupId, GroupName: u.GroupName,
		Level: u.Level, Balance: u.Balance, Credit: u.Credit, MoneyCount: u.MoneyCount,
		IsDisabled: u.IsDisabled, RegisterAt: u.RegisterAt, LastLoginAt: u.LastLoginAt,
	}
}
