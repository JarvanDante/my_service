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
	"github.com/JarvanDante/my_service/internal/shared/siteconf"
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
	return &v1.LoginRes{Token: dto.Token, User: toApiUser(ctx, dto.User)}, nil
}

// Avatars 需登录: 默认头像列表。
func (c *Controller) Avatars(ctx context.Context, req *v1.AvatarsReq) (res *v1.AvatarsRes, err error) {
	return &v1.AvatarsRes{List: c.user.DefaultAvatars(ctx)}, nil
}

// Restore 公开: 凭证找回账号。
func (c *Controller) Restore(ctx context.Context, req *v1.RestoreReq) (res *v1.RestoreRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	dto, err := c.user.Restore(ctx, service.RestoreInput{
		Credential:    req.Credential,
		DeviceId:      req.DeviceId,
		DeviceType:    req.DeviceType,
		DeviceVersion: req.DeviceVersion,
		Ip:            r.GetClientIp(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.RestoreRes{Token: dto.Token, User: toApiUser(ctx, dto.User)}, nil
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
	return &v1.InfoRes{UserInfo: toApiUser(ctx, dto)}, nil
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

// DoSign 每日签到(需登录)。
func (c *Controller) DoSign(ctx context.Context, req *v1.SignReq) (res *v1.SignRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.DoDaySign(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.SignRes{Today: dto.Today, Days: dto.Days, Count: dto.Count, Reward: dto.Reward}, nil
}

// Tasks 任务列表(需登录)。
func (c *Controller) Tasks(ctx context.Context, req *v1.TasksReq) (res *v1.TasksRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, err := c.user.Tasks(ctx, id)
	if err != nil {
		return nil, err
	}
	items := make([]v1.Task, 0, len(list))
	for _, t := range list {
		items = append(items, v1.Task{
			Id: t.Id, Name: t.Name, Type: t.Type, Description: t.Description,
			MaxNum: t.MaxNum, DoneToday: t.DoneToday,
		})
	}
	return &v1.TasksRes{List: items}, nil
}

// DoTask 完成任务(需登录)。
func (c *Controller) DoTask(ctx context.Context, req *v1.TaskDoReq) (res *v1.TaskDoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.DoTask(ctx, id, req.TaskId)
	if err != nil {
		return nil, err
	}
	return &v1.TaskDoRes{DoneToday: dto.DoneToday, MaxNum: dto.MaxNum, Reward: dto.Reward}, nil
}

// TaskLogs 任务记录(需登录)。
func (c *Controller) TaskLogs(ctx context.Context, req *v1.TaskLogsReq) (res *v1.TaskLogsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.TaskLogs(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.TaskLog, 0, len(list))
	for _, l := range list {
		items = append(items, v1.TaskLog{
			Id: l.Id, TaskId: l.TaskId, Type: l.Type, Num: l.Num, LogDate: l.LogDate, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.TaskLogsRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Up 成长信息(需登录)。
func (c *Controller) Up(ctx context.Context, req *v1.UpReq) (res *v1.UpRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.Up(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.UpRes{
		Level: dto.Level, Credit: dto.Credit, Balance: dto.Balance,
		SignDaysThisMonth: dto.SignDaysThisMonth, Fans: dto.Fans, Follow: dto.Follow, ShareNum: dto.ShareNum,
	}, nil
}

// Recharge 充值套餐(需登录)。
func (c *Controller) Recharge(ctx context.Context, req *v1.RechargeReq) (res *v1.RechargeRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	list, err := c.user.RechargePackages(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.RechargePackage, 0, len(list))
	for _, p := range list {
		items = append(items, v1.RechargePackage{Id: p.Id, Name: p.Name, Amount: p.Amount, Coin: p.Coin, Bonus: p.Bonus})
	}
	return &v1.RechargeRes{List: items}, nil
}

// RechargeDo 发起充值(需登录)。
func (c *Controller) RechargeDo(ctx context.Context, req *v1.RechargeDoReq) (res *v1.RechargeDoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.DoRecharge(ctx, id, req.PackageId)
	if err != nil {
		return nil, err
	}
	return &v1.RechargeDoRes{OrderNo: dto.OrderNo, Amount: dto.Amount, Coin: dto.Coin}, nil
}

// RechargeMockPay 开发环境 mock 支付到账(需登录, 仅本人订单)。
func (c *Controller) RechargeMockPay(ctx context.Context, req *v1.RechargeMockPayReq) (res *v1.RechargeMockPayRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.MockPay(ctx, id, req.OrderNo); err != nil {
		return nil, err
	}
	return &v1.RechargeMockPayRes{}, nil
}

// Vip VIP套餐(需登录)。
func (c *Controller) Vip(ctx context.Context, req *v1.VipReq) (res *v1.VipRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	list, err := c.user.VipPackages(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.VipPackage, 0, len(list))
	for _, p := range list {
		items = append(items, v1.VipPackage{Id: p.Id, Name: p.Name, Days: p.Days, Price: p.Price})
	}
	return &v1.VipRes{List: items}, nil
}

// VipDo 开通/续费 VIP(需登录)。
func (c *Controller) VipDo(ctx context.Context, req *v1.VipDoReq) (res *v1.VipDoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.DoVip(ctx, id, req.PackageId); err != nil {
		return nil, err
	}
	return &v1.VipDoRes{}, nil
}

// VipLogs VIP记录(需登录)。
func (c *Controller) VipLogs(ctx context.Context, req *v1.VipLogsReq) (res *v1.VipLogsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.VipLogs(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.VipLog, 0, len(list))
	for _, l := range list {
		items = append(items, v1.VipLog{
			Id: l.Id, PackageId: l.PackageId, Days: l.Days, Price: l.Price,
			StartAt: l.StartAt, EndAt: l.EndAt, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.VipLogsRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Exchange 兑换信息(需登录)。
func (c *Controller) Exchange(ctx context.Context, req *v1.ExchangeReq) (res *v1.ExchangeRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	dto, err := c.user.ExchangeInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.ExchangeRes{Rate: dto.Rate, Credit: dto.Credit, Balance: dto.Balance}, nil
}

// ExchangeDo 积分兑换金币(需登录)。
func (c *Controller) ExchangeDo(ctx context.Context, req *v1.ExchangeDoReq) (res *v1.ExchangeDoRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.DoExchange(ctx, id, req.Coin); err != nil {
		return nil, err
	}
	return &v1.ExchangeDoRes{}, nil
}

// Chats 会话列表(需登录)。
func (c *Controller) Chats(ctx context.Context, req *v1.ChatsReq) (res *v1.ChatsRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.Chats(ctx, id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.Chat, 0, len(list))
	for _, d := range list {
		items = append(items, v1.Chat{
			PeerId: d.PeerId, Nickname: d.Nickname, Img: d.Img,
			LastContent: d.LastContent, LastAt: d.LastAt, Unread: d.Unread,
		})
	}
	return &v1.ChatsRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// ChatMessages 会话消息(需登录)。
func (c *Controller) ChatMessages(ctx context.Context, req *v1.ChatMessagesReq) (res *v1.ChatMessagesRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.user.Messages(ctx, id, req.PeerId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.Message, 0, len(list))
	for _, m := range list {
		items = append(items, v1.Message{
			Id: m.Id, FromId: m.FromId, ToId: m.ToId, Content: m.Content, Mine: m.Mine, CreatedAt: m.CreatedAt,
		})
	}
	return &v1.ChatMessagesRes{List: items, Total: total, Page: req.Page, Size: req.Size}, nil
}

// ChatSend 发消息(需登录)。
func (c *Controller) ChatSend(ctx context.Context, req *v1.ChatSendReq) (res *v1.ChatSendRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.SendMessage(ctx, id, req.ToId, req.Content); err != nil {
		return nil, err
	}
	return &v1.ChatSendRes{}, nil
}

// ChatDel 删除会话(需登录)。
func (c *Controller) ChatDel(ctx context.Context, req *v1.ChatDelReq) (res *v1.ChatDelRes, err error) {
	id, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.user.DelChat(ctx, id, req.PeerId); err != nil {
		return nil, err
	}
	return &v1.ChatDelRes{}, nil
}

// CustomerUrl 客服链接(需登录)。
func (c *Controller) CustomerUrl(ctx context.Context, req *v1.CustomerUrlReq) (res *v1.CustomerUrlRes, err error) {
	if _, err = uid(ctx); err != nil {
		return nil, err
	}
	url, err := c.user.CustomerUrl(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CustomerUrlRes{Url: url}, nil
}

func toApiUser(ctx context.Context, d *service.UserInfoDTO) v1.UserInfo {
	if d == nil {
		return v1.UserInfo{}
	}
	return v1.UserInfo{
		Id: d.Id, Username: d.Username, Nickname: d.Nickname, Phone: d.Phone,
		Img: d.Img, Signature: d.Signature, Sex: d.Sex, Level: d.Level,
		Balance: d.Balance, Credit: d.Credit, GroupName: d.GroupName,
		Fans: d.Fans, Follow: d.Follow,
		// 站点差异字段: 仅返回 Nacos response.user_info_extra 白名单内的 key
		Ext: siteconf.PickExt(ctx, "user_info", map[string]interface{}{
			"bg_img":         d.BgImg,
			"share_num":      d.ShareNum,
			"channel_name":   d.ChannelName,
			"group_end_time": d.GroupEndTime,
		}),
	}
}
