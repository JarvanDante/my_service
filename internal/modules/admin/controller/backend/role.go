// Package backend 后台角色 / 权限控制器。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/admin/v1"
)

// ListRoles 角色列表。
func (c *Controller) ListRoles(ctx context.Context, req *v1.ListRolesReq) (res *v1.ListRolesRes, err error) {
	dtos, err := c.admin.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]v1.RoleItem, 0, len(dtos))
	for _, d := range dtos {
		list = append(list, v1.RoleItem{Id: d.Id, Name: d.Name, Code: d.Code, Remark: d.Remark, Status: d.Status})
	}
	return &v1.ListRolesRes{List: list}, nil
}

// ListPerms 角色权限列表。
func (c *Controller) ListPerms(ctx context.Context, req *v1.ListPermsReq) (res *v1.ListPermsRes, err error) {
	dtos, err := c.admin.ListPerms(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	list := make([]v1.PermItem, 0, len(dtos))
	for _, d := range dtos {
		list = append(list, v1.PermItem{Path: d.Path, Method: d.Method})
	}
	return &v1.ListPermsRes{List: list}, nil
}

// AddPerm 新增角色权限。
func (c *Controller) AddPerm(ctx context.Context, req *v1.AddPermReq) (res *v1.AddPermRes, err error) {
	if err = c.admin.AddPerm(ctx, req.Code, req.Path, req.Method); err != nil {
		return nil, err
	}
	return &v1.AddPermRes{}, nil
}

// RemovePerm 删除角色权限。
func (c *Controller) RemovePerm(ctx context.Context, req *v1.DelPermReq) (res *v1.DelPermRes, err error) {
	if err = c.admin.RemovePerm(ctx, req.Code, req.Path, req.Method); err != nil {
		return nil, err
	}
	return &v1.DelPermRes{}, nil
}
