package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	admindomain "github.com/JarvanDante/my_service/internal/modules/admin/domain"
)

type adminRepo struct{}

func NewAdminRepo() admindomain.Repository { return &adminRepo{} }

func (r *adminRepo) FindByUsername(ctx context.Context, username string) (*entity.AdminUser, error) {
	var a *entity.AdminUser
	err := g.Model("admin_user").Ctx(ctx).Where("username", username).Scan(&a)
	return a, err
}

func (r *adminRepo) FindById(ctx context.Context, id int64) (*entity.AdminUser, error) {
	var a *entity.AdminUser
	err := g.Model("admin_user").Ctx(ctx).Where("id", id).Scan(&a)
	return a, err
}

func (r *adminRepo) UpdateLoginInfo(ctx context.Context, id int64, ip string) error {
	_, err := g.Model("admin_user").Ctx(ctx).Where("id", id).Data(g.Map{
		"last_login_at": gtime.Now(),
		"last_ip":       ip,
	}).Update()
	return err
}
