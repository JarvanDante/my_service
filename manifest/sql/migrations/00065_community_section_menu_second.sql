-- 社区管理子菜单：帖子模块排到意见反馈上面（第二位）。
-- +goose Up
UPDATE admin_permission SET sort = 2
WHERE is_menu = 1 AND name = '帖子模块' AND route_url = '/community/section';
UPDATE admin_permission SET sort = 3
WHERE is_menu = 1 AND name = '意见反馈' AND route_url = '/community/feedback';

-- +goose Down
UPDATE admin_permission SET sort = 2
WHERE is_menu = 1 AND name = '意见反馈' AND route_url = '/community/feedback';
UPDATE admin_permission SET sort = 3
WHERE is_menu = 1 AND name = '帖子模块' AND route_url = '/community/section';
