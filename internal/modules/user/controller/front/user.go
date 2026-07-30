// Package front 前台用户控制器(适配层, 薄)。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ user service.IUser }

func New(svc service.IUser) *Controller { return &Controller{user: svc} }

// uid 从 ctx 取当前用户(Auth 中间件写入)。
func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

// Login 公开。
func (c *Controller) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	dto, err := c.user.Login(ctx, service.LoginInput{
		DeviceId:      req.DeviceId,
		DeviceType:    req.DeviceType,
		DeviceVersion: req.DeviceVersion,
		Ip:            r.GetClientIp(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{Token: dto.Token, User: toApiUser(dto.User)}, nil
}

// Info 需登录。
func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.Info(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.InfoRes{UserInfo: toApiUser(dto)}, nil
}

// Logout 退出登录。
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.Logout(ctx, id); err != nil {
		return nil, err
	}
	return &v1.LogoutRes{}, nil
}

// Refresh 刷新 token。
func (c *Controller) Refresh(ctx context.Context, req *v1.RefreshReq) (res *v1.RefreshRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	token, err := c.user.Refresh(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.RefreshRes{Token: token}, nil
}

// BindPhone 绑定手机。
func (c *Controller) BindPhone(ctx context.Context, req *v1.BindPhoneReq) (res *v1.BindPhoneRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.BindPhone(ctx, id, req.Phone, req.Code); err != nil {
		return nil, err
	}
	return &v1.BindPhoneRes{}, nil
}

// Home 他人主页(需登录)。
func (c *Controller) Home(ctx context.Context, req *v1.HomeReq) (res *v1.HomeRes, err error) {
	viewer, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.Home(ctx, viewer, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.HomeRes{User: toPublicApi(dto.User), IsFollowed: dto.IsFollowed}, nil
}

// Update 修改资料(需登录)。
func (c *Controller) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.UpdateProfile(ctx, id, service.UpdateProfileInput{
		Nickname:  req.Nickname,
		Img:       req.Img,
		BgImg:     req.BgImg,
		Signature: req.Signature,
		Sex:       req.Sex,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Images 用户图片(需登录)。
func (c *Controller) Images(ctx context.Context, req *v1.ImagesReq) (res *v1.ImagesRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	imgs, err := c.user.Images(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.ImagesRes{Images: imgs}, nil
}

// Find 按账号查找(需登录)。
func (c *Controller) Find(ctx context.Context, req *v1.FindReq) (res *v1.FindRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	dto, err := c.user.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if dto == nil {
		return &v1.FindRes{Found: false}, nil
	}
	p := toPublicApi(dto)
	return &v1.FindRes{Found: true, User: &p}, nil
}

func toPublicApi(d *service.PublicUserDTO) v1.PublicUser {
	if d == nil {
		return v1.PublicUser{}
	}
	return v1.PublicUser{
		Id: d.Id, Nickname: d.Nickname, Img: d.Img, BgImg: d.BgImg,
		Signature: d.Signature, Sex: d.Sex, Level: d.Level,
		Fans: d.Fans, Follow: d.Follow, ShareNum: d.ShareNum,
	}
}

// DoFollow 关注/取关(需登录)。
func (c *Controller) DoFollow(ctx context.Context, req *v1.FollowReq) (res *v1.FollowRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	followed, err := c.user.DoFollow(ctx, id, req.HomeId)
	if err != nil {
		return nil, err
	}
	return &v1.FollowRes{Followed: followed}, nil
}

// Follows 我的关注列表(需登录)。
func (c *Controller) Follows(ctx context.Context, req *v1.FollowsReq) (res *v1.FollowsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.Following(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.FollowsRes{List: toPublicApiList(list), Total: total, Page: req.Page, Size: req.Size}, nil
}

// Fans 我的粉丝列表(需登录)。
func (c *Controller) Fans(ctx context.Context, req *v1.FansReq) (res *v1.FansRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.Fans(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.FansRes{List: toPublicApiList(list), Total: total, Page: req.Page, Size: req.Size}, nil
}

func toPublicApiList(ds []*service.PublicUserDTO) []v1.PublicUser {
	out := make([]v1.PublicUser, 0, len(ds))
	for _, d := range ds {
		out = append(out, toPublicApi(d))
	}
	return out
}

// BindParent 绑定推荐人(需登录)。
func (c *Controller) BindParent(ctx context.Context, req *v1.BindParentReq) (res *v1.BindParentRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.BindParent(ctx, id, req.Account); err != nil {
		return nil, err
	}
	return &v1.BindParentRes{}, nil
}

// BindCode 绑定邀请码(需登录, 复用绑定推荐人逻辑)。
func (c *Controller) BindCode(ctx context.Context, req *v1.BindCodeReq) (res *v1.BindCodeRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.BindParent(ctx, id, req.Code); err != nil {
		return nil, err
	}
	return &v1.BindCodeRes{}, nil
}

// Redeem 使用兑换码(需登录)。
func (c *Controller) Redeem(ctx context.Context, req *v1.RedeemReq) (res *v1.RedeemRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.RedeemCode(ctx, id, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.RedeemRes{Type: dto.Type, AddNum: dto.AddNum, Name: dto.Name}, nil
}

// CodeLogs 兑换记录(需登录)。
func (c *Controller) CodeLogs(ctx context.Context, req *v1.CodeLogsReq) (res *v1.CodeLogsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.CodeLogs(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.CodeLog, 0, len(list))
	for _, l := range list {
		items = append(items, v1.CodeLog{
			Id: l.Id, Code: l.Code, Name: l.Name, Type: l.Type, AddNum: l.AddNum, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.CodeLogsRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Share 分享信息(需登录)。
func (c *Controller) Share(ctx context.Context, req *v1.ShareReq) (res *v1.ShareRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.ShareInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.ShareRes{ShareCode: dto.ShareCode, ShareUrl: dto.ShareUrl, ShareNum: dto.ShareNum}, nil
}

// ShareLogs 分享记录(需登录)。
func (c *Controller) ShareLogs(ctx context.Context, req *v1.ShareLogsReq) (res *v1.ShareLogsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.ShareLogs(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.ShareLog, 0, len(list))
	for _, l := range list {
		items = append(items, v1.ShareLog{
			Id: l.Id, Type: l.Type, TargetId: l.TargetId, Channel: l.Channel, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.ShareLogsRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// ShareReport 上报分享(需登录)。
func (c *Controller) ShareReport(ctx context.Context, req *v1.ShareReportReq) (res *v1.ShareReportRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.ReportShare(ctx, id, req.Type, req.TargetId, req.Channel); err != nil {
		return nil, err
	}
	return &v1.ShareReportRes{}, nil
}

func toApiUser(d *service.UserInfoDTO) v1.UserInfo {
	if d == nil {
		return v1.UserInfo{}
	}
	return v1.UserInfo{
		Id: d.Id, Username: d.Username, Nickname: d.Nickname, Phone: d.Phone,
		Img: d.Img, Signature: d.Signature, Sex: d.Sex, Level: d.Level,
		Balance: d.Balance, Credit: d.Credit, GroupName: d.GroupName,
		Fans: d.Fans, Follow: d.Follow,
	}
}
