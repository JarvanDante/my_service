-- 去掉子后台「运营配置」菜单（公告/跳转位/敏感词后台页）。
-- 运营中心系统公告、审核「禁用词库」保留；filter_word / announcement 表仍给前台与禁词校验用。
-- +goose Up
UPDATE admin_role r
SET permissions = coalesce((
    SELECT string_agg(tok, ',')
    FROM unnest(string_to_array(r.permissions, ',')) AS tok
    WHERE tok <> ''
      AND tok ~ '^[0-9]+$'
      AND tok::bigint NOT IN (
          SELECT p.id
          FROM admin_permission p
          WHERE (p.is_menu = 1 AND p.name = '运营配置' AND p.route_url = '/ops/config')
             OR p.parent_id IN (
                 SELECT id FROM admin_permission
                 WHERE is_menu = 1 AND name = '运营配置' AND route_url = '/ops/config'
             )
      )
), '');

DELETE FROM admin_permission
WHERE is_menu = 0
  AND parent_id IN (
      SELECT id FROM admin_permission
      WHERE is_menu = 1 AND name = '运营配置' AND route_url = '/ops/config'
  );

DELETE FROM admin_permission
WHERE is_menu = 1 AND name = '运营配置' AND route_url = '/ops/config';

-- +goose Down
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '运营配置', '/ops/config', 'ops/config', '', 'lucide:settings-2', 1, 0, 0, 2
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '运营管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '运营配置' AND x.route_url = '/ops/config'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('公告列表', '/backend/announcement',      'GET',    1),
    ('新增公告', '/backend/announcement',      'POST',   2),
    ('编辑公告', '/backend/announcement/{id}', 'PUT',    3),
    ('删除公告', '/backend/announcement/{id}', 'DELETE', 4),
    ('跳转位列表', '/backend/jumptab',         'GET',    5),
    ('新增跳转位', '/backend/jumptab',         'POST',   6),
    ('编辑跳转位', '/backend/jumptab/{id}',    'PUT',    7),
    ('删除跳转位', '/backend/jumptab/{id}',    'DELETE', 8),
    ('敏感词列表', '/backend/filterword',      'GET',    9),
    ('新增敏感词', '/backend/filterword',      'POST',   10),
    ('删除敏感词', '/backend/filterword/{id}', 'DELETE', 11)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '运营配置' AND m.route_url = '/ops/config'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));
