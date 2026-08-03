// Package backend 后台角色控制器(RBAC 菜单驱动)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/admin/v1"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
)

// ListRoles 角色列表。
func (c *Controller) ListRoles(ctx context.Context, req *v1.ListRolesReq) (res *v1.ListRolesRes, err error) {
	dtos, err := c.admin.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]v1.RoleItem, 0, len(dtos))
	for _, d := range dtos {
		list = append(list, v1.RoleItem{Id: d.Id, Name: d.Name, Code: d.Code, Remark: d.Remark, Status: d.Status, Permissions: d.Permissions})
	}
	return &v1.ListRolesRes{List: list}, nil
}

// PermTree 权限树。
func (c *Controller) PermTree(ctx context.Context, req *v1.PermTreeReq) (res *v1.PermTreeRes, err error) {
	list, err := c.admin.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PermTreeRes{List: toPermNodes(list)}, nil
}

func toPermNodes(list []*service.PermissionDTO) []v1.PermNode {
	out := make([]v1.PermNode, 0, len(list))
	for _, p := range list {
		out = append(out, v1.PermNode{
			Id: p.Id, ParentId: p.ParentId, Name: p.Name, RouteUrl: p.RouteUrl,
			Component: p.Component, Method: p.Method, Icon: p.Icon, IsMenu: p.IsMenu,
			HideInMenu: p.HideInMenu, AffixTab: p.AffixTab, ActivePath: p.ActivePath,
			Sort: p.Sort, Status: p.Status, Children: toPermNodes(p.Children),
		})
	}
	return out
}

// PermCreate 新增权限节点。
func (c *Controller) PermCreate(ctx context.Context, req *v1.PermCreateReq) (res *v1.PermCreateRes, err error) {
	id, err := c.admin.CreatePermission(ctx, service.PermissionInput{
		ParentId: req.ParentId, Name: req.Name, RouteUrl: req.RouteUrl, Component: req.Component,
		Method: req.Method, Icon: req.Icon, IsMenu: req.IsMenu, HideInMenu: req.HideInMenu,
		AffixTab: req.AffixTab, ActivePath: req.ActivePath, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PermCreateRes{Id: id}, nil
}

// PermUpdate 更新权限节点。
func (c *Controller) PermUpdate(ctx context.Context, req *v1.PermUpdateReq) (res *v1.PermUpdateRes, err error) {
	if err = c.admin.UpdatePermission(ctx, service.PermissionInput{
		Id: req.Id, ParentId: req.ParentId, Name: req.Name, RouteUrl: req.RouteUrl, Component: req.Component,
		Method: req.Method, Icon: req.Icon, IsMenu: req.IsMenu, HideInMenu: req.HideInMenu,
		AffixTab: req.AffixTab, ActivePath: req.ActivePath, Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.PermUpdateRes{}, nil
}

// PermDelete 删除权限节点。
func (c *Controller) PermDelete(ctx context.Context, req *v1.PermDeleteReq) (res *v1.PermDeleteRes, err error) {
	if err = c.admin.DeletePermission(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.PermDeleteRes{}, nil
}
