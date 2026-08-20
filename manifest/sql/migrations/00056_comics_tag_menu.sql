-- 漫画标签: 复用 tag 表 content_type=4, 在「漫画管理」下单独挂一份菜单方便运营维护。
-- +goose Up
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '漫画标签', '/comics/tag', 'content/comics-tag', '', 'lucide:tags', 1, 0, 0, 3
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '漫画管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '漫画标签' AND x.route_url = '/comics/tag'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('标签列表', '/backend/tag',     'GET',    1),
    ('新增标签', '/backend/tag',     'POST',   2),
    ('编辑标签', '/backend/tag/{id}', 'PUT',    3),
    ('删除标签', '/backend/tag/{id}', 'DELETE', 4)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '漫画标签' AND m.route_url = '/comics/tag'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0
  AND parent_id IN (
      SELECT id FROM admin_permission
      WHERE is_menu = 1 AND name = '漫画标签' AND route_url = '/comics/tag'
  );
DELETE FROM admin_permission
WHERE is_menu = 1 AND name = '漫画标签' AND route_url = '/comics/tag';
