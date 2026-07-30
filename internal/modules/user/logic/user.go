// Package logic 用户业务实现。依赖 domain.Repository 接口。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
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
		// 自动注册
		now := gtime.Now()
		id, err := s.repo.Create(ctx, g.Map{
			"username":       "device_" + in.DeviceId,
			"nickname":       "用户" + grand.Digits(6),
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
	token, err := kit.IssueToken(ctx, u.Id)
	if err != nil {
		return nil, err
	}
	return &service.LoginDTO{
		Token: token,
		User:  toUserInfo(u),
	}, nil
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
	return &service.UserInfoDTO{
		Id:        u.Id,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Phone:     u.Phone,
		Img:       u.Img,
		Signature: u.Signature,
		Sex:       u.Sex,
		Level:     u.Level,
		Balance:   u.Balance,
		Credit:    u.Credit,
		GroupName: u.GroupName,
		Fans:      u.Fans,
		Follow:    u.Follow,
	}
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
