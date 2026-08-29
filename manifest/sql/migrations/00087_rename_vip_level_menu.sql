-- 用户管理下二级菜单「会员等级」改名为「VIP等级」。
-- 只改展示名, 路由 /user/growth 和页面 component=growth/index 不变。
-- H5 会员中心售卖项也改读 user_group, 不再读独立 vip_package。
-- +goose Up
UPDATE admin_permission
SET name = 'VIP等级'
WHERE is_menu = 1
  AND name = '会员等级'
  AND (route_url LIKE '%growth%' OR COALESCE(component, '') LIKE '%growth%');

-- +goose Down
UPDATE admin_permission
SET name = '会员等级'
WHERE is_menu = 1
  AND name = 'VIP等级'
  AND (route_url LIKE '%growth%' OR COALESCE(component, '') LIKE '%growth%');
