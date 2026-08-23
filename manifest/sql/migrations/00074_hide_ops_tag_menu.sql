-- 运营管理下的「标签管理」与各模块自己的标签页重复, 从侧栏拿掉。
-- 页面 content/tag 和 /backend/tag 接口保留, 漫画/动漫/视频标签菜单不动。
-- +goose Up
UPDATE admin_permission
SET hide_in_menu = 1
WHERE is_menu = 1
  AND name = '标签管理'
  AND route_url IN ('/ops/tag', '/content/tag');

-- +goose Down
UPDATE admin_permission
SET hide_in_menu = 0
WHERE is_menu = 1
  AND name = '标签管理'
  AND route_url IN ('/ops/tag', '/content/tag');
