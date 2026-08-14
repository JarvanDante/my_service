-- 用户管理下二级菜单「用户组与成长」对齐公司后台, 改名为「会员等级」。
-- 公司 yhdm/dm 等站该页面对应 Mongo user_group(VIP 套餐/会员卡), 后台菜单名就是「会员等级」。
-- 只改展示名, 路由 /user/growth 和页面 component=growth/index 不变。
-- +goose Up
UPDATE admin_permission
SET name = '会员等级'
WHERE is_menu = 1
  AND name = '用户组与成长';

-- +goose Down
UPDATE admin_permission
SET name = '用户组与成长'
WHERE is_menu = 1
  AND name = '会员等级'
  AND (route_url LIKE '%growth%' OR COALESCE(component, '') LIKE '%growth%');
