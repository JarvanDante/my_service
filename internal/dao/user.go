package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
)

type userRepo struct{}

// NewUserRepo 返回 user 领域仓储实现。
func NewUserRepo() domain.Repository { return &userRepo{} }

func (r *userRepo) FindById(ctx context.Context, id int64) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Scan(&u)
	return u, err
}

func (r *userRepo) FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error) {
	var u *entity.Users
	err := Users.Ctx(ctx).Where(Users.Columns().DeviceId, deviceId).Scan(&u)
	return u, err
}

func (r *userRepo) Create(ctx context.Context, data g.Map) (int64, error) {
	return Users.Ctx(ctx).Data(data).InsertAndGetId()
}

func (r *userRepo) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(g.Map{
		Users.Columns().LoginNum:    &gdb.Counter{Field: Users.Columns().LoginNum, Value: 1},
		Users.Columns().LastLoginAt: gtime.Now(),
		Users.Columns().LastIp:      ip,
	}).Update()
	return err
}

func (r *userRepo) Disable(ctx context.Context, id int64, reason string) error {
	_, err := Users.Ctx(ctx).Where(Users.Columns().Id, id).Data(g.Map{
		Users.Columns().IsDisabled: 1,
		Users.Columns().ErrorMsg:   reason,
	}).Update()
	return err
}
