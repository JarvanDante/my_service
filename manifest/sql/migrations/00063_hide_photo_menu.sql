-- 子后台侧栏去掉「图集管理」及「图集列表」, 页面和接口保留。
-- +goose Up
UPDATE admin_permission
SET hide_in_menu = 1
WHERE is_menu = 1
  AND (
      (name = '图集管理' AND route_url = '/photo')
      OR (name = '图集列表' AND route_url = '/photo/list')
  );

-- +goose Down
UPDATE admin_permission
SET hide_in_menu = 0
WHERE is_menu = 1
  AND (
      (name = '图集管理' AND route_url = '/photo')
      OR (name = '图集列表' AND route_url = '/photo/list')
  );
