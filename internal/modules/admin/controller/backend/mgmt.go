// Package backend 后台角色/管理员账号管理控制器。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/admin/v1"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
)

// ---------- 角色管理 ----------

func (c *Controller) CreateRole(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	id, err := c.admin.CreateRole(ctx, service.RoleCreateInput{Name: req.Name, Code: req.Code, Remark: req.Remark, Permissions: req.Permissions})
	if err != nil {
		return nil, err
	}
	return &v1.RoleCreateRes{Id: id}, nil
}

func (c *Controller) UpdateRole(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	if err = c.admin.UpdateRole(ctx, service.RoleUpdateInput{Id: req.Id, Name: req.Name, Remark: req.Remark, Status: req.Status, Permissions: req.Permissions}); err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}

func (c *Controller) DeleteRole(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	if err = c.admin.DeleteRole(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}

// ---------- 管理员账号管理 ----------

func (c *Controller) ListAdmins(ctx context.Context, req *v1.AdminListReq) (res *v1.AdminListRes, err error) {
	dto, err := c.admin.ListAdmins(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	list := make([]v1.AdminItem, 0, len(dto.List))
	for _, a := range dto.List {
		list = append(list, v1.AdminItem{
			Id: a.Id, Username: a.Username, Nickname: a.Nickname,
			RoleId: a.RoleId, RoleName: a.RoleName, Status: a.Status, LastLoginAt: a.LastLoginAt,
		})
	}
	return &v1.AdminListRes{List: list, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) CreateAdmin(ctx context.Context, req *v1.AdminCreateReq) (res *v1.AdminCreateRes, err error) {
	id, err := c.admin.CreateAdmin(ctx, service.AdminCreateInput{
		Username: req.Username, Password: req.Password, Nickname: req.Nickname, RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminCreateRes{Id: id}, nil
}

func (c *Controller) UpdateAdmin(ctx context.Context, req *v1.AdminUpdateReq) (res *v1.AdminUpdateRes, err error) {
	if err = c.admin.UpdateAdmin(ctx, service.AdminUpdateInput{
		Id: req.Id, Nickname: req.Nickname, RoleId: req.RoleId, Status: req.Status, Password: req.Password,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminUpdateRes{}, nil
}

func (c *Controller) DeleteAdmin(ctx context.Context, req *v1.AdminDeleteReq) (res *v1.AdminDeleteRes, err error) {
	operator, err := adminId(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.admin.DeleteAdmin(ctx, req.Id, operator); err != nil {
		return nil, err
	}
	return &v1.AdminDeleteRes{}, nil
}
