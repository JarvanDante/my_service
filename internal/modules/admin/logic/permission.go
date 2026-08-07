// Package logic — RBAC 菜单+接口权限 实现(子后台)。
package logic

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
)

// parseIds 把 "1,2,3" 解析为 []int64。
func parseIds(s string) []int64 {
	out := make([]int64, 0)
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			out = append(out, v)
		}
	}
	return out
}

// ---------- 权限树 CRUD ----------

func (s *sAdmin) ListPermissions(ctx context.Context) ([]*service.PermissionDTO, error) {
	rows, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	nodeMap := make(map[int64]*service.PermissionDTO, len(rows))
	for _, p := range rows {
		nodeMap[p.Id] = &service.PermissionDTO{
			Id: p.Id, ParentId: p.ParentId, Name: p.Name, RouteUrl: p.RouteUrl,
			Component: p.Component, Method: p.Method, Icon: p.Icon, IsMenu: p.IsMenu,
			HideInMenu: p.HideInMenu, AffixTab: p.AffixTab, ActivePath: p.ActivePath,
			Sort: p.Sort, Status: p.Status, Children: []*service.PermissionDTO{},
		}
	}
	var roots []*service.PermissionDTO
	for _, p := range rows {
		node := nodeMap[p.Id]
		if p.ParentId == 0 {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[p.ParentId]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}
	return roots, nil
}

func (s *sAdmin) CreatePermission(ctx context.Context, in service.PermissionInput) (int64, error) {
	if in.Name == "" {
		return 0, gerror.New("名称必填")
	}
	return s.repo.CreatePermission(ctx, toPermEntity(in))
}

func (s *sAdmin) UpdatePermission(ctx context.Context, in service.PermissionInput) error {
	if in.Id <= 0 {
		return gerror.New("id 必填")
	}
	if in.Name == "" {
		return gerror.New("名称必填")
	}
	return s.repo.UpdatePermission(ctx, toPermEntity(in))
}

func (s *sAdmin) DeletePermission(ctx context.Context, id int64) error {
	cnt, err := s.repo.CountPermissionChildren(ctx, id)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return gerror.New("请先删除子节点")
	}
	return s.repo.DeletePermission(ctx, id)
}

func toPermEntity(in service.PermissionInput) *entity.AdminPermission {
	return &entity.AdminPermission{
		Id: in.Id, ParentId: in.ParentId, Name: in.Name, RouteUrl: in.RouteUrl,
		Component: in.Component, Method: in.Method, Icon: in.Icon, IsMenu: in.IsMenu,
		HideInMenu: in.HideInMenu, AffixTab: in.AffixTab, ActivePath: in.ActivePath,
		Sort: in.Sort, Status: in.Status,
	}
}

// ---------- 当前管理员菜单(/auth/menus) ----------

func (s *sAdmin) MenusForAdmin(ctx context.Context, adminId int64) ([]*service.MenuNodeDTO, error) {
	admin, err := s.repo.FindById(ctx, adminId)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return []*service.MenuNodeDTO{}, nil
	}
	role, err := s.repo.FindRoleById(ctx, admin.RoleId)
	if err != nil {
		return nil, err
	}

	all, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var menus []*entity.AdminPermission
	if role != nil && role.Code == superAdminCode {
		for _, p := range all {
			if p.IsMenu == 1 && p.Status == 1 {
				menus = append(menus, p)
			}
		}
	} else if role != nil {
		// 勾选接口权限时也要带上祖先菜单, 否则侧栏为空/缺父级
		menus = menusForRole(all, parseIds(role.Permissions))
	}
	return buildMenuTree(menus), nil
}

// menusForRole 根据角色勾选的权限 id, 展开祖先后筛出启用菜单节点。
func menusForRole(all []*entity.AdminPermission, selected []int64) []*entity.AdminPermission {
	if len(selected) == 0 {
		return nil
	}
	byID := make(map[int64]*entity.AdminPermission, len(all))
	for _, p := range all {
		byID[p.Id] = p
	}
	keep := make(map[int64]struct{}, len(selected)*2)
	for _, id := range selected {
		for cur := id; cur > 0; {
			if _, ok := keep[cur]; ok {
				break
			}
			p, ok := byID[cur]
			if !ok {
				break
			}
			keep[cur] = struct{}{}
			cur = p.ParentId
		}
	}
	out := make([]*entity.AdminPermission, 0, len(keep))
	for id := range keep {
		p := byID[id]
		if p != nil && p.IsMenu == 1 && p.Status == 1 {
			out = append(out, p)
		}
	}
	return out
}

func buildMenuTree(menus []*entity.AdminPermission) []*service.MenuNodeDTO {
	nodeMap := make(map[int64]*service.MenuNodeDTO, len(menus))
	for _, p := range menus {
		nodeMap[p.Id] = &service.MenuNodeDTO{
			Name:       routeName(p.RouteUrl),
			Path:       p.RouteUrl,
			Component:  p.Component,
			Title:      p.Name,
			Icon:       p.Icon,
			Order:      p.Sort,
			AffixTab:   p.AffixTab == 1,
			HideInMenu: p.HideInMenu == 1,
			ActivePath: p.ActivePath,
			Children:   []*service.MenuNodeDTO{},
		}
	}
	var roots []*service.MenuNodeDTO
	for _, p := range menus {
		node := nodeMap[p.Id]
		if p.ParentId == 0 {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[p.ParentId]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}
	// 纯目录(自身无 component)才重定向到首个子菜单
	for _, n := range nodeMap {
		if len(n.Children) > 0 && n.Component == "" {
			n.Redirect = n.Children[0].Path
		}
	}
	return roots
}

// routeName 由 route_url 生成唯一路由名: "/system/role" -> "system-role"。
func routeName(routeUrl string) string {
	s := strings.TrimPrefix(routeUrl, "/")
	s = strings.ReplaceAll(s, "/", "-")
	if s == "" {
		s = "root"
	}
	return s
}
