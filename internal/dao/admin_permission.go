package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

func (r *adminRepo) ListPermissions(ctx context.Context) ([]*entity.AdminPermission, error) {
	var list []*entity.AdminPermission
	err := g.Model("admin_permission").Ctx(ctx).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *adminRepo) FindPermissionsByIds(ctx context.Context, ids []int64) ([]*entity.AdminPermission, error) {
	var list []*entity.AdminPermission
	if len(ids) == 0 {
		return list, nil
	}
	err := g.Model("admin_permission").Ctx(ctx).WhereIn("id", ids).Order("sort asc, id asc").Scan(&list)
	return list, err
}

func (r *adminRepo) FindPermissionById(ctx context.Context, id int64) (*entity.AdminPermission, error) {
	var p *entity.AdminPermission
	err := g.Model("admin_permission").Ctx(ctx).Where("id", id).Scan(&p)
	return p, err
}

func (r *adminRepo) CreatePermission(ctx context.Context, p *entity.AdminPermission) (int64, error) {
	res, err := g.Model("admin_permission").Ctx(ctx).Data(g.Map{
		"parent_id": p.ParentId, "name": p.Name, "route_url": p.RouteUrl,
		"component": p.Component, "method": p.Method, "icon": p.Icon,
		"is_menu": p.IsMenu, "hide_in_menu": p.HideInMenu, "affix_tab": p.AffixTab,
		"active_path": p.ActivePath, "sort": p.Sort, "status": p.Status,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *adminRepo) UpdatePermission(ctx context.Context, p *entity.AdminPermission) error {
	_, err := g.Model("admin_permission").Ctx(ctx).Where("id", p.Id).Data(g.Map{
		"parent_id": p.ParentId, "name": p.Name, "route_url": p.RouteUrl,
		"component": p.Component, "method": p.Method, "icon": p.Icon,
		"is_menu": p.IsMenu, "hide_in_menu": p.HideInMenu, "affix_tab": p.AffixTab,
		"active_path": p.ActivePath, "sort": p.Sort, "status": p.Status,
		"updated_at": gtime.Now(),
	}).Update()
	return err
}

func (r *adminRepo) DeletePermission(ctx context.Context, id int64) error {
	_, err := g.Model("admin_permission").Ctx(ctx).Where("id", id).Delete()
	return err
}

func (r *adminRepo) CountPermissionChildren(ctx context.Context, parentId int64) (int, error) {
	return g.Model("admin_permission").Ctx(ctx).Where("parent_id", parentId).Count()
}
