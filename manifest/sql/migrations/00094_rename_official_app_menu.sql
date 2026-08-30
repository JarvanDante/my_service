-- 运营管理「福利应用」改名为「官方应用」，排到「官方社群」下面。
-- +goose Up
UPDATE admin_permission
SET name = '官方应用', sort = 10
WHERE is_menu = 1
  AND route_url = '/ops/application'
  AND name IN ('福利应用', '推广应用');

UPDATE admin_permission
SET sort = 11
WHERE is_menu = 1 AND name = '签到配置' AND route_url = '/ops/checkin';

-- +goose Down
UPDATE admin_permission
SET name = '福利应用', sort = 6
WHERE is_menu = 1 AND name = '官方应用' AND route_url = '/ops/application';

UPDATE admin_permission
SET sort = 10
WHERE is_menu = 1 AND name = '签到配置' AND route_url = '/ops/checkin';
