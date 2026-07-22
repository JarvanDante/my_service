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
