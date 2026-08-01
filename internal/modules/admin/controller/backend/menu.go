package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/admin/v1"
)

// Menus 返回当前站点管理员可见的菜单树(动态菜单 method B)。
// 目前返回子后台完整菜单;后续可按 adminId 的角色/Casbin 权限做过滤。
func (c *Controller) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	if _, err = adminId(ctx); err != nil {
		return nil, err
	}
	return &v1.MenusRes{List: buildBackendMenus()}, nil
}

// buildBackendMenus 子后台菜单树。component 为相对 apps/web-ele/src/views 的路径(不含 .vue)。
func buildBackendMenus() []v1.MenuNode {
	return []v1.MenuNode{
		{
			Name:     "Dashboard",
			Path:     "/dashboard",
			Redirect: "/analytics",
			Meta:     v1.MenuMeta{Title: "仪表盘", Icon: "lucide:layout-dashboard", Order: -1},
			Children: []v1.MenuNode{
				{
					Name:      "Analytics",
					Path:      "/analytics",
					Component: "dashboard/analytics/index",
					Meta:      v1.MenuMeta{Title: "分析页", Icon: "lucide:area-chart", AffixTab: true},
				},
			},
		},
		{
			Name:      "UserManage",
			Path:      "/user",
			Component: "user/index",
			Meta:      v1.MenuMeta{Title: "用户管理", Icon: "lucide:users", Order: 10},
		},
		{
			Name:      "FinanceManage",
			Path:      "/finance",
			Component: "finance/index",
			Meta:      v1.MenuMeta{Title: "财务管理", Icon: "lucide:wallet", Order: 20},
		},
		{
			Name:      "PromoManage",
			Path:      "/promo",
			Component: "promo/index",
			Meta:      v1.MenuMeta{Title: "兑换码", Icon: "lucide:ticket", Order: 30},
		},
		{
			Name:      "GrowthManage",
			Path:      "/growth",
			Component: "growth/index",
			Meta:      v1.MenuMeta{Title: "用户组与成长", Icon: "lucide:trophy", Order: 40},
		},
		{
			Name:      "OpsManage",
			Path:      "/ops",
			Component: "ops/index",
			Meta:      v1.MenuMeta{Title: "运营管理", Icon: "lucide:megaphone", Order: 50},
		},
		{
			Name:     "System",
			Path:     "/system",
			Redirect: "/system/role",
			Meta:     v1.MenuMeta{Title: "系统管理", Icon: "lucide:settings", Order: 90},
			Children: []v1.MenuNode{
				{
					Name:      "SystemRole",
					Path:      "role",
					Component: "system/role",
					Meta:      v1.MenuMeta{Title: "角色权限", Icon: "lucide:shield-check"},
				},
				{
					Name:      "SystemAdmin",
					Path:      "admin",
					Component: "system/admin",
					Meta:      v1.MenuMeta{Title: "管理员", Icon: "lucide:users"},
				},
			},
		},
	}
}
