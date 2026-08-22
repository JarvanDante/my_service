// Package logic 用户业务实现。依赖 domain.Repository 接口。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/dao"
	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
	"github.com/JarvanDante/my_service/internal/shared/siteconf"
)

type sUser struct {
	repo domain.Repository
}

// New 注入仓储, 返回 service.IUser 实现。
func New(repo domain.Repository) service.IUser {
	return &sUser{repo: repo}
}

// Login 设备登录: 有则登录, 无则自动注册。
func (s *sUser) Login(ctx context.Context, in service.LoginInput) (*service.LoginDTO, error) {
	if in.DeviceId == "" {
		return nil, gerror.New("device_id 不能为空")
	}
	u, err := s.repo.FindByDeviceId(ctx, in.DeviceId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		// 自动注册。username 先用设备号占位避免撞唯一索引, 拿到自增 id 后再
		// 写成 dm 风格公开编号(encodeUserId); 设备可能曾被凭证找回换绑过。
		username := "device_" + in.DeviceId
		if exist, e := s.repo.FindByAccount(ctx, username); e != nil {
			return nil, e
		} else if exist != nil {
			username = username + "_" + grand.S(4)
		}
		now := gtime.Now()
		id, err := s.repo.Create(ctx, g.Map{
			"username":       username,
			"nickname":       randNickname(),
			"img":            randAvatar(ctx),
			"bg_img":         randBackground(ctx),
			"slat":           grand.S(8),
			"device_id":      in.DeviceId,
			"device_type":    in.DeviceType,
			"device_version": in.DeviceVersion,
			"register_at":    now,
			"register_date":  gconv.Int(now.Format("Ymd")),
			"register_ip":    in.Ip,
			"last_login_at":  now,
			"last_date":      gconv.Int(now.Format("Ymd")),
			"last_ip":        in.Ip,
			"login_num":      1,
		})
		if err != nil {
			return nil, err
		}
		if u, err = s.repo.FindById(ctx, id); err != nil {
			return nil, err
		}
	} else {
		if u.IsDisabled == 1 {
			msg := u.ErrorMsg
			if msg == "" {
				msg = "账号已被禁用"
			}
			return nil, gerror.New(msg)
		}
		_ = s.repo.UpdateLoginInfo(ctx, u.Id, in.Ip)
	}
	u = s.ensureEncodedUsername(ctx, u)
	token, err := kit.IssueToken(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{
		Token: token,
		User:  toUserInfo(u),
	}, nil
}

// Restore 凭证找回账号(公开): 校验 username==>md5(username_appid) 凭证,
// 通过后把账号换绑到当前设备并签发 token。重装后自动注册的空壳号会被让位(见 RebindDevice)。
func (s *sUser) Restore(ctx context.Context, in service.RestoreInput) (*service.LoginDTO, error) {
	cred := strings.TrimSpace(in.Credential)
	if cred == "" || in.DeviceId == "" {
		return nil, gerror.New("凭证与设备号不能为空")
	}
	parts := strings.SplitN(cred, "==>", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, gerror.New("凭证格式不正确")
	}
	username := parts[0]
	u, err := s.repo.FindByAccount(ctx, username)
	if err != nil {
		return nil, err
	}
	if u == nil || u.Username != username {
		return nil, gerror.New("凭证无效: 账号不存在")
	}
	// 与后台登录二维码同一套算法核验, 防止仅凭用户名冒领
	expect := buildAccountSlat(ctx, u.Username, u.SiteId, nil)
	if expect == "" || !strings.EqualFold(expect, username+"==>"+parts[1]) {
		return nil, gerror.New("凭证无效: 校验不通过")
	}
	if u.IsDisabled == 1 {
		msg := u.ErrorMsg
		if msg == "" {
			msg = "账号已被禁用"
		}
		return nil, gerror.New(msg)
	}
	if err = s.repo.RebindDevice(ctx, u.Id, in.DeviceId, in.DeviceType, in.DeviceVersion, in.Ip); err != nil {
		return nil, err
	}
	if u, err = s.repo.FindById(ctx, u.Id); err != nil {
		return nil, err
	}
	u = s.ensureEncodedUsername(ctx, u)
	token, err := kit.IssueToken(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{Token: token, User: toUserInfo(u)}, nil
}

func hashUserPassword(password, slat string) string {
	return gmd5.MustEncryptString(password + slat)
}

func (s *sUser) findLoginAccount(ctx context.Context, account string) (*entity.Users, error) {
	account = strings.TrimSpace(account)
	if account == "" {
		return nil, nil
	}
	u, err := s.repo.FindByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}
	if upper := strings.ToUpper(account); upper != account {
		u, err = s.repo.FindByAccount(ctx, upper)
		if err != nil {
			return nil, err
		}
		if u != nil {
			return u, nil
		}
	}
	if uid := kit.DecodePublicId(account); uid > 0 {
		return s.repo.FindById(ctx, uid)
	}
	if uid := kit.DecodeUserId(account); uid > 0 {
		return s.repo.FindById(ctx, uid)
	}
	if uid := kit.DecodeUserId(strings.ToUpper(account)); uid > 0 {
		return s.repo.FindById(ctx, uid)
	}
	return nil, nil
}

// AccountLogin 用户编号+密码登录(不是设备号)。通过后换绑到当前设备。
func (s *sUser) AccountLogin(ctx context.Context, in service.AccountLoginInput) (*service.LoginDTO, error) {
	account := strings.TrimSpace(in.Username)
	password := strings.TrimSpace(in.Password)
	if account == "" || password == "" || in.DeviceId == "" {
		return nil, gerror.New("用户名、密码与设备号不能为空")
	}
	u, err := s.findLoginAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户名或密码错误")
	}
	if u.Password == "" {
		return nil, gerror.New("该账号尚未设置密码, 请先在本机设置")
	}
	slat := u.Slat
	if slat == "" {
		return nil, gerror.New("用户名或密码错误")
	}
	if hashUserPassword(password, slat) != u.Password {
		return nil, gerror.New("用户名或密码错误")
	}
	if u.IsDisabled == 1 {
		msg := u.ErrorMsg
		if msg == "" {
			msg = "账号已被禁用"
		}
		return nil, gerror.New(msg)
	}
	if err = s.repo.RebindDevice(ctx, u.Id, in.DeviceId, in.DeviceType, in.DeviceVersion, in.Ip); err != nil {
		return nil, err
	}
	if u, err = s.repo.FindById(ctx, u.Id); err != nil {
		return nil, err
	}
	u = s.ensureEncodedUsername(ctx, u)
	token, err := kit.IssueToken(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{Token: token, User: toUserInfo(u)}, nil
}

// SetPassword 首次设密或修改密码。已设过则必须校验旧密码。
func (s *sUser) SetPassword(ctx context.Context, userId int64, oldPassword, password string) error {
	password = strings.TrimSpace(password)
	if len(password) < 6 || len(password) > 32 {
		return gerror.New("密码需6-32位")
	}
	u, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if u == nil {
		return gerror.New("用户不存在")
	}
	slat := u.Slat
	if slat == "" {
		slat = grand.S(8)
	}
	if u.Password != "" {
		if strings.TrimSpace(oldPassword) == "" {
			return gerror.New("请输入旧密码")
		}
		if hashUserPassword(oldPassword, u.Slat) != u.Password {
			return gerror.New("旧密码不正确")
		}
	}
	return s.repo.UpdateProfile(ctx, userId, g.Map{
		"password": hashUserPassword(password, slat),
		"slat":     slat,
	})
}

// Info 当前用户详情。
func (s *sUser) Info(ctx context.Context, userId int64) (*service.UserInfoDTO, error) {
	u, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户不存在")
	}
	u = s.ensureEncodedUsername(ctx, u)
	return toUserInfo(u), nil
}

func (s *sUser) GetUser(ctx context.Context, id int64) (*service.UserDTO, error) {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户不存在")
	}
	return &service.UserDTO{ID: u.Id, Name: u.Username, Active: u.IsDisabled == 0}, nil
}

func (s *sUser) DisableUser(ctx context.Context, id int64) error {
	return s.repo.Disable(ctx, id, "后台禁用")
}

func toUserInfo(u *entity.Users) *service.UserInfoDTO {
	username := kit.EncodeUserId(u.Id)
	if username == "" {
		username = u.Username
	}
	return &service.UserInfoDTO{
		Id:           u.Id,
		Username:     username,
		Nickname:     u.Nickname,
		Phone:        u.Phone,
		Img:          u.Img,
		Signature:    u.Signature,
		Sex:          u.Sex,
		Level:        u.Level,
		Balance:      u.Balance,
		Credit:       u.Credit,
		GroupName:    u.GroupName,
		Fans:         u.Fans,
		Follow:       u.Follow,
		BgImg:        u.BgImg,
		ShareNum:     u.ShareNum,
		ChannelName:  u.ChannelName,
		GroupEndTime: u.GroupEndTime,
		HasPassword:  u.Password != "",
	}
}

// ensureEncodedUsername 把 username 回写成 dm 风格公开编号(encodeUserId)。
// 旧账号仍是 device_* 时, 登录/拉资料顺带迁移, 邀请码与前台编号一致。
func (s *sUser) ensureEncodedUsername(ctx context.Context, u *entity.Users) *entity.Users {
	if u == nil {
		return u
	}
	code := kit.EncodeUserId(u.Id)
	if code == "" || u.Username == code {
		return u
	}
	if err := s.repo.UpdateProfile(ctx, u.Id, g.Map{"username": code}); err != nil {
		g.Log().Warningf(ctx, "回写用户编号失败 id=%d: %v", u.Id, err)
		return u
	}
	u.Username = code
	return u
}

// Logout 退出登录: 撤销当前会话 token。
func (s *sUser) Logout(ctx context.Context, userId int64) error {
	return kit.RevokeByUserId(ctx, userId)
}

// Refresh 刷新 token: 重新签发(旧 token 自动失效)。
func (s *sUser) Refresh(ctx context.Context, userId int64) (string, error) {
	u, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", gerror.New("用户不存在")
	}
	if u.IsDisabled == 1 {
		return "", gerror.New("账号已被禁用")
	}
	return kit.IssueToken(ctx, userId)
}

// BindPhone 绑定手机号(校验唯一性)。
func (s *sUser) BindPhone(ctx context.Context, userId int64, phone, code string) error {
	// TODO: 校验短信验证码 code(接入短信服务后启用)
	_ = code
	other, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if other != nil && other.Id != userId {
		return gerror.New("该手机号已被其他账号绑定")
	}
	return s.repo.UpdatePhone(ctx, userId, phone)
}

// Home 他人主页: 公开信息 + 当前用户是否已关注。
func (s *sUser) Home(ctx context.Context, viewerId, homeId int64) (*service.HomeDTO, error) {
	u, err := s.repo.FindById(ctx, homeId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户不存在")
	}
	followed := false
	if viewerId > 0 && viewerId != homeId {
		followed, _ = s.repo.ExistsFollow(ctx, viewerId, homeId)
	}
	return &service.HomeDTO{User: toPublic(u), IsFollowed: followed}, nil
}

// UpdateProfile 改资料: 仅更新传了值的字段。
func (s *sUser) UpdateProfile(ctx context.Context, userId int64, in service.UpdateProfileInput) error {
	data := g.Map{}
	if in.Nickname != "" {
		data["nickname"] = in.Nickname
	}
	if in.Img != "" {
		if !isDefaultAvatar(ctx, in.Img) {
			return gerror.New("请从系统头像中选择")
		}
		data["img"] = in.Img
	}
	if in.BgImg != "" {
		data["bg_img"] = in.BgImg
	}
	if in.Signature != "" {
		data["signature"] = in.Signature
	}
	if in.Sex == 1 || in.Sex == 2 {
		data["sex"] = in.Sex
	}
	if len(data) == 0 {
		return gerror.New("没有要更新的内容")
	}
	return s.repo.UpdateProfile(ctx, userId, data)
}

// Images 用户图片: 目前返回头像+背景图。
// TODO: 接入相册/帖子媒体表后, 改为查真实图片列表。
func (s *sUser) Images(ctx context.Context, userId int64) ([]string, error) {
	u, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, gerror.New("用户不存在")
	}
	imgs := []string{}
	if u.Img != "" {
		imgs = append(imgs, u.Img)
	}
	if u.BgImg != "" {
		imgs = append(imgs, u.BgImg)
	}
	return imgs, nil
}

// FindByAccount 按账号(用户名/手机号)查找, 未找到返回 (nil, nil)。
func (s *sUser) FindByAccount(ctx context.Context, account string) (*service.PublicUserDTO, error) {
	u, err := s.repo.FindByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return toPublic(u), nil
}

func toPublic(u *entity.Users) *service.PublicUserDTO {
	return &service.PublicUserDTO{
		Id:        u.Id,
		Nickname:  u.Nickname,
		Img:       u.Img,
		BgImg:     u.BgImg,
		Signature: u.Signature,
		Sex:       u.Sex,
		Level:     u.Level,
		Fans:      u.Fans,
		Follow:    u.Follow,
		ShareNum:  u.ShareNum,
	}
}

// DoFollow 关注/取关切换, 返回操作后是否处于已关注状态。
func (s *sUser) DoFollow(ctx context.Context, userId, homeId int64) (bool, error) {
	if homeId <= 0 {
		return false, gerror.New("目标用户不合法")
	}
	if userId == homeId {
		return false, gerror.New("不能关注自己")
	}
	target, err := s.repo.FindById(ctx, homeId)
	if err != nil {
		return false, err
	}
	if target == nil {
		return false, gerror.New("用户不存在")
	}
	exists, err := s.repo.ExistsFollow(ctx, userId, homeId)
	if err != nil {
		return false, err
	}
	if exists {
		if err = s.repo.Unfollow(ctx, userId, homeId); err != nil {
			return false, err
		}
		return false, nil
	}
	if err = s.repo.Follow(ctx, userId, homeId); err != nil {
		return false, err
	}
	return true, nil
}

// Following 我的关注列表。
func (s *sUser) Following(ctx context.Context, userId int64, page, size int) ([]*service.PublicUserDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.FollowingList(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toPublicList(list), total, nil
}

// Fans 我的粉丝列表。
func (s *sUser) Fans(ctx context.Context, userId int64, page, size int) ([]*service.PublicUserDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.FansList(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	return toPublicList(list), total, nil
}

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 15
	}
	return page, size
}

func toPublicList(us []*entity.Users) []*service.PublicUserDTO {
	out := make([]*service.PublicUserDTO, 0, len(us))
	for _, u := range us {
		out = append(out, toPublic(u))
	}
	return out
}

// BindParent 绑定推荐人(by account/邀请码), 已绑定则不可改。
func (s *sUser) BindParent(ctx context.Context, userId int64, account string) error {
	if account == "" {
		return gerror.New("请填写推荐人/邀请码")
	}
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if me == nil {
		return gerror.New("用户不存在")
	}
	inviter, err := s.repo.FindByAccount(ctx, account)
	if err != nil {
		return err
	}
	if inviter == nil {
		if uid := kit.DecodeUserId(account); uid > 0 {
			inviter, err = s.repo.FindById(ctx, uid)
			if err != nil {
				return err
			}
		}
	}
	if inviter == nil {
		if uid := kit.DecodePublicId(account); uid > 0 {
			inviter, err = s.repo.FindById(ctx, uid)
			if err != nil {
				return err
			}
		}
	}
	if inviter == nil {
		return gerror.New("推荐人不存在")
	}
	if inviter.Id == userId {
		return gerror.New("不能绑定自己")
	}
	if err = s.rejectInviteCycle(ctx, userId, inviter); err != nil {
		return err
	}
	if me.ParentId != 0 {
		return gerror.New("已绑定推荐人, 不可修改")
	}
	parentName := kit.EncodeUserId(inviter.Id)
	if parentName == "" {
		parentName = inviter.Username
	}
	return s.repo.BindInviter(ctx, userId, inviter.Id, parentName)
}

// rejectInviteCycle 邀请只能单向：先绑定的生效，反向或沿邀请链回到自己则拒绝。
func (s *sUser) rejectInviteCycle(ctx context.Context, userId int64, inviter *entity.Users) error {
	if inviter.ParentId == userId {
		return gerror.New("不能互相邀请")
	}
	seen := map[int64]struct{}{inviter.Id: {}}
	cur := inviter.ParentId
	for i := 0; cur > 0 && i < 64; i++ {
		if cur == userId {
			return gerror.New("不能互相邀请")
		}
		if _, ok := seen[cur]; ok {
			break
		}
		seen[cur] = struct{}{}
		u, err := s.repo.FindById(ctx, cur)
		if err != nil {
			return err
		}
		if u == nil || u.ParentId <= 0 {
			break
		}
		cur = u.ParentId
	}
	return nil
}

// RedeemCode 使用兑换码。
func (s *sUser) RedeemCode(ctx context.Context, userId int64, code string) (*service.RedeemDTO, error) {
	if code == "" {
		return nil, gerror.New("请填写兑换码")
	}
	c, err := s.repo.FindCodeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, gerror.New("兑换码不存在")
	}
	if c.Status == -1 {
		return nil, gerror.New("兑换码已作废")
	}
	if c.ExpiredAt > 0 && c.ExpiredAt < gtime.Timestamp() {
		return nil, gerror.New("兑换码已过期")
	}
	if c.UsedNum >= c.CanUseNum {
		return nil, gerror.New("兑换码已用完")
	}
	used, err := s.repo.HasRedeemed(ctx, c.Id, userId)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, gerror.New("你已使用过该兑换码")
	}
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, gerror.New("用户不存在")
	}
	if err = s.repo.RedeemCode(ctx, userId, me.Username, c); err != nil {
		return nil, err
	}
	return &service.RedeemDTO{Type: c.Type, AddNum: c.AddNum, Name: c.Name}, nil
}

// CodeLogs 兑换记录。
func (s *sUser) CodeLogs(ctx context.Context, userId int64, page, size int) ([]*service.CodeLogDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.CodeLogs(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.CodeLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.CodeLogDTO{
			Id: l.Id, Code: l.Code, Name: l.Name, Type: l.Type, AddNum: l.AddNum,
			CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

// ShareInfo 分享信息(分享码=用户名)。
func (s *sUser) ShareInfo(ctx context.Context, userId int64) (*service.ShareDTO, error) {
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, gerror.New("用户不存在")
	}
	code := kit.EncodePublicId(me.Id)
	if code == "" {
		code = me.Username
	}
	return &service.ShareDTO{
		ShareCode: code,
		ShareUrl:  "",
		ShareNum:  me.ShareNum,
	}, nil
}

// ReportShare 上报一次分享。
func (s *sUser) ReportShare(ctx context.Context, userId int64, typ string, targetId int64, channel string) error {
	if typ == "" {
		typ = "app"
	}
	return s.repo.AddShareLog(ctx, userId, typ, targetId, channel)
}

// ShareLogs 分享记录。
func (s *sUser) ShareLogs(ctx context.Context, userId int64, page, size int) ([]*service.ShareLogDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.ShareLogList(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.ShareLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.ShareLogDTO{
			Id: l.Id, Type: l.Type, TargetId: l.TargetId, Channel: l.Channel,
			CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

func (s *sUser) Invitees(ctx context.Context, userId int64, page, size int) ([]*service.InviteeDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.InviteeList(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.InviteeDTO, 0, len(list))
	for _, u := range list {
		code := kit.EncodePublicId(u.Id)
		if code == "" {
			code = kit.EncodeUserId(u.Id)
		}
		name := strings.TrimSpace(u.Nickname)
		if name == "" {
			name = "用户" + code
		}
		at := u.CreatedAt
		if u.RegisterAt != nil {
			at = u.RegisterAt
		}
		out = append(out, &service.InviteeDTO{
			Nickname: name, InviteCode: code, CreatedAt: fmtTime(at),
		})
	}
	return out, total, nil
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

const (
	signCredit = 5.0  // 每日签到奖励积分
	taskCredit = 10.0 // 完成任务奖励积分
)

// DoDaySign 每日签到。
func (s *sUser) DoDaySign(ctx context.Context, userId int64) (*service.SignDTO, error) {
	now := gtime.Now()
	ym := gconv.Int(now.Format("Ym"))
	today := now.Day()
	days, exists, err := s.repo.GetSignDays(ctx, userId, ym)
	if err != nil {
		return nil, err
	}
	for _, d := range days {
		if d == today {
			return nil, gerror.New("今日已签到")
		}
	}
	days = append(days, today)
	if err = s.repo.SaveSign(ctx, userId, ym, days, exists, signCredit); err != nil {
		return nil, err
	}
	return &service.SignDTO{Today: today, Days: days, Count: len(days), Reward: signCredit}, nil
}

// Tasks 任务列表(含当日完成次数)。
func (s *sUser) Tasks(ctx context.Context, userId int64) ([]*service.TaskDTO, error) {
	list, err := s.repo.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	logDate := gconv.Int(gtime.Now().Format("Ymd"))
	out := make([]*service.TaskDTO, 0, len(list))
	for _, t := range list {
		done, _ := s.repo.TaskDoneToday(ctx, userId, t.Id, logDate)
		out = append(out, &service.TaskDTO{
			Id: t.Id, Name: t.Name, Type: t.Type, Description: t.Description,
			MaxNum: t.MaxNum, DoneToday: done,
		})
	}
	return out, nil
}

// DoTask 完成任务领奖(受单日上限约束)。
func (s *sUser) DoTask(ctx context.Context, userId, taskId int64) (*service.TaskDoneDTO, error) {
	t, err := s.repo.FindTask(ctx, taskId)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, gerror.New("任务不存在")
	}
	if t.Status != 1 {
		return nil, gerror.New("任务已下线")
	}
	logDate := gconv.Int(gtime.Now().Format("Ymd"))
	done, err := s.repo.TaskDoneToday(ctx, userId, taskId, logDate)
	if err != nil {
		return nil, err
	}
	if t.MaxNum > 0 && done >= t.MaxNum {
		return nil, gerror.New("今日该任务已完成")
	}
	reward := t.Reward
	if reward <= 0 {
		reward = taskCredit // 兜底: 未配置奖励时用默认值
	}
	if err = s.repo.AddTaskLog(ctx, userId, taskId, t.Type, logDate, reward); err != nil {
		return nil, err
	}
	return &service.TaskDoneDTO{DoneToday: done + 1, MaxNum: t.MaxNum, Reward: reward}, nil
}

// TaskLogs 任务记录。
func (s *sUser) TaskLogs(ctx context.Context, userId int64, page, size int) ([]*service.TaskLogDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.TaskLogs(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.TaskLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.TaskLogDTO{
			Id: l.Id, TaskId: l.TaskId, Type: l.Type, Num: l.Num, LogDate: l.LogDate,
			CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

// Up 成长信息。
func (s *sUser) Up(ctx context.Context, userId int64) (*service.UpDTO, error) {
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, gerror.New("用户不存在")
	}
	ym := gconv.Int(gtime.Now().Format("Ym"))
	days, _, _ := s.repo.GetSignDays(ctx, userId, ym)
	return &service.UpDTO{
		Level: me.Level, Credit: me.Credit, Balance: me.Balance,
		SignDaysThisMonth: len(days), Fans: me.Fans, Follow: me.Follow, ShareNum: me.ShareNum,
	}, nil
}

const creditPerCoin = 100.0 // 100 积分 = 1 金币

func (s *sUser) RechargePackages(ctx context.Context) ([]*service.RechargePackageDTO, error) {
	list, err := s.repo.ListRechargePackages(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.RechargePackageDTO, 0, len(list))
	for _, p := range list {
		out = append(out, &service.RechargePackageDTO{Id: p.Id, Name: p.Name, Amount: p.Amount, Coin: p.Coin, Bonus: p.Bonus})
	}
	return out, nil
}

// DoRecharge 发起充值: 创建待支付订单。到账在支付回调, 此处不加金币。
func (s *sUser) DoRecharge(ctx context.Context, userId, packageId int64) (*service.RechargeOrderDTO, error) {
	pkg, err := s.repo.FindRechargePackage(ctx, packageId)
	if err != nil {
		return nil, err
	}
	if pkg == nil {
		return nil, gerror.New("充值套餐不存在")
	}
	orderNo := "R" + gconv.String(gtime.Timestamp()) + grand.Digits(6)
	coin := pkg.Coin + pkg.Bonus
	if err = s.repo.CreateRechargeOrder(ctx, orderNo, userId, packageId, pkg.Amount, coin); err != nil {
		return nil, err
	}
	// TODO: 返回真实支付参数(接支付网关); 支付成功回调里再给用户加金币并置订单为已支付。
	return &service.RechargeOrderDTO{OrderNo: orderNo, Amount: pkg.Amount, Coin: coin}, nil
}

// MockPay 开发环境模拟支付到账: 仅当未配置 pay.callbackSecret 时可用, 且订单必须属于当前用户。
func (s *sUser) MockPay(ctx context.Context, userId int64, orderNo string) error {
	if orderNo == "" {
		return gerror.New("订单号必填")
	}
	secret := g.Cfg().MustGet(ctx, "pay.callbackSecret").String()
	if secret != "" {
		return gerror.New("正式环境请走支付回调")
	}
	order, err := s.repo.FindRechargeOrder(ctx, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return gerror.New("订单不存在")
	}
	if order.UserId != userId {
		return gerror.New("订单不属于当前用户")
	}
	return dao.NewFinanceRepo().MarkOrderPaid(ctx, orderNo)
}

func (s *sUser) VipPackages(ctx context.Context) ([]*service.VipPackageDTO, error) {
	list, err := s.repo.ListVipPackages(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.VipPackageDTO, 0, len(list))
	for _, p := range list {
		out = append(out, &service.VipPackageDTO{Id: p.Id, Name: p.Name, Days: p.Days, Price: p.Price})
	}
	return out, nil
}

// DoVip 用金币开通/续费 VIP。
func (s *sUser) DoVip(ctx context.Context, userId, packageId int64) error {
	pkg, err := s.repo.FindVipPackage(ctx, packageId)
	if err != nil {
		return err
	}
	if pkg == nil {
		return gerror.New("VIP套餐不存在")
	}
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}
	if me == nil {
		return gerror.New("用户不存在")
	}
	now := gtime.Timestamp()
	base := now
	if me.GroupEndTime > base {
		base = me.GroupEndTime // 未过期则续期
	}
	endAt := base + int64(pkg.Days)*86400
	return s.repo.OpenVip(ctx, userId, pkg, now, endAt)
}

func (s *sUser) VipLogs(ctx context.Context, userId int64, page, size int) ([]*service.VipLogDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.VipLogs(ctx, userId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.VipLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.VipLogDTO{
			Id: l.Id, PackageId: l.PackageId, Days: l.Days, Price: l.Price,
			StartAt: l.StartAt, EndAt: l.EndAt, CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

func (s *sUser) ExchangeInfo(ctx context.Context, userId int64) (*service.ExchangeInfoDTO, error) {
	me, err := s.repo.FindById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if me == nil {
		return nil, gerror.New("用户不存在")
	}
	return &service.ExchangeInfoDTO{Rate: int(creditPerCoin), Credit: me.Credit, Balance: me.Balance}, nil
}

// DoExchange 积分兑金币: coin 为想兑换的金币数, 花费 coin*rate 积分。
func (s *sUser) DoExchange(ctx context.Context, userId int64, coin int) error {
	if coin <= 0 {
		return gerror.New("兑换数量不合法")
	}
	cost := float64(coin) * creditPerCoin
	return s.repo.ExchangeCreditToCoin(ctx, userId, cost, float64(coin))
}

// SendMessage 发私信。
func (s *sUser) SendMessage(ctx context.Context, meId, toId int64, content string) error {
	if content == "" {
		return gerror.New("消息内容不能为空")
	}
	if toId == meId {
		return gerror.New("不能给自己发消息")
	}
	peer, err := s.repo.FindById(ctx, toId)
	if err != nil {
		return err
	}
	if peer == nil {
		return gerror.New("对方不存在")
	}
	return s.repo.SendMessage(ctx, meId, toId, content)
}

// Chats 会话列表(附对方昵称/头像)。
func (s *sUser) Chats(ctx context.Context, meId int64, page, size int) ([]*service.ChatDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.ListConversations(ctx, meId, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.ChatDTO, 0, len(list))
	for _, c := range list {
		d := &service.ChatDTO{
			PeerId: c.PeerId, LastContent: c.LastContent, Unread: c.Unread, LastAt: fmtTime(c.LastAt),
		}
		if peer, _ := s.repo.FindById(ctx, c.PeerId); peer != nil {
			d.Nickname = peer.Nickname
			d.Img = peer.Img
		}
		out = append(out, d)
	}
	return out, total, nil
}

// Messages 会话消息(读取后清零未读)。
func (s *sUser) Messages(ctx context.Context, meId, peerId int64, page, size int) ([]*service.MessageDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.Messages(ctx, meId, peerId, page, size)
	if err != nil {
		return nil, 0, err
	}
	_ = s.repo.MarkRead(ctx, meId, peerId)
	out := make([]*service.MessageDTO, 0, len(list))
	for _, m := range list {
		out = append(out, &service.MessageDTO{
			Id: m.Id, FromId: m.FromId, ToId: m.ToId, Content: m.Content,
			Mine: m.FromId == meId, CreatedAt: fmtTime(m.CreatedAt),
		})
	}
	return out, total, nil
}

// DelChat 删除会话(仅对自己)。
func (s *sUser) DelChat(ctx context.Context, meId, peerId int64) error {
	return s.repo.DeleteConversation(ctx, meId, peerId)
}

// CustomerUrl 客服链接(从配置读, 缺省占位)。
func (s *sUser) CustomerUrl(ctx context.Context) (string, error) {
	def := g.Cfg().MustGet(ctx, "customer.url", "https://example.com/kefu").String()
	return siteconf.Get(ctx, "customer_url", def), nil
}
