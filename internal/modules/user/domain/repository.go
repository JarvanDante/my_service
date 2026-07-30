// Package domain 用户领域层。
package domain

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	FindById(ctx context.Context, id int64) (*entity.Users, error)
	FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Users, error)
	FindByAccount(ctx context.Context, account string) (*entity.Users, error)
	Create(ctx context.Context, data g.Map) (int64, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
	UpdatePhone(ctx context.Context, id int64, phone string) error
	UpdateProfile(ctx context.Context, id int64, data g.Map) error
	Disable(ctx context.Context, id int64, reason string) error
	// 社交
	ExistsFollow(ctx context.Context, userId, homeId int64) (bool, error)
	Follow(ctx context.Context, userId, homeId int64) error
	Unfollow(ctx context.Context, userId, homeId int64) error
	FollowingList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error)
	FansList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error)
	// 推广 / 兑换码
	BindInviter(ctx context.Context, userId, inviterId int64, inviterName string) error
	FindCodeByCode(ctx context.Context, code string) (*entity.UserCode, error)
	HasRedeemed(ctx context.Context, codeId, userId int64) (bool, error)
	RedeemCode(ctx context.Context, userId int64, username string, code *entity.UserCode) error
	CodeLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserCodeLog, int, error)
	AddShareLog(ctx context.Context, userId int64, typ string, targetId int64, channel string) error
	ShareLogList(ctx context.Context, userId int64, page, size int) ([]*entity.UserShareLog, int, error)
}
