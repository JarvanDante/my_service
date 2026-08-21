-- 一级菜单「漫画管理」排到「动漫管理」上面。
-- +goose Up
UPDATE admin_permission SET sort = 21 WHERE is_menu = 1 AND parent_id = 0 AND name = '漫画管理' AND route_url = '/comics';
UPDATE admin_permission SET sort = 22 WHERE is_menu = 1 AND parent_id = 0 AND name = '动漫管理' AND route_url = '/cartoon';
UPDATE admin_permission SET sort = 23 WHERE is_menu = 1 AND parent_id = 0 AND name = '小说管理' AND route_url = '/novel';

-- +goose Down
UPDATE admin_permission SET sort = 21 WHERE is_menu = 1 AND parent_id = 0 AND name = '动漫管理' AND route_url = '/cartoon';
UPDATE admin_permission SET sort = 22 WHERE is_menu = 1 AND parent_id = 0 AND name = '漫画管理' AND route_url = '/comics';
UPDATE admin_permission SET sort = 22 WHERE is_menu = 1 AND parent_id = 0 AND name = '小说管理' AND route_url = '/novel';
