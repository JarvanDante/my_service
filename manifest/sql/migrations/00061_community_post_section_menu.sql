-- 社区管理下增加「帖子模块」：维护发帖可选板块(tag.content_type=6)。
-- +goose Up
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '帖子模块', '/community/section', 'community/section', '', 'lucide:layout-grid', 1, 0, 0, 3
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '社区管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '帖子模块' AND x.route_url = '/community/section'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('板块列表', '/backend/tag',      'GET',    1),
    ('新增板块', '/backend/tag',      'POST',   2),
    ('编辑板块', '/backend/tag/{id}', 'PUT',    3),
    ('删除板块', '/backend/tag/{id}', 'DELETE', 4)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '帖子模块' AND m.route_url = '/community/section'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

-- 已有「帖子管理」的角色一并勾上新菜单和接口, 超管本来就能看到全部菜单。
UPDATE admin_role r
SET permissions = trim(both ',' from r.permissions || ',' || ids.joined)
FROM (
    SELECT string_agg(p.id::text, ',') AS joined
    FROM admin_permission p
    WHERE (p.is_menu = 1 AND p.name = '帖子模块' AND p.route_url = '/community/section')
       OR p.parent_id IN (
           SELECT id FROM admin_permission
           WHERE is_menu = 1 AND name = '帖子模块' AND route_url = '/community/section'
       )
) ids
WHERE r.code <> 'superadmin'
  AND ids.joined IS NOT NULL
  AND (
      (',' || r.permissions || ',') LIKE '%,130,%'
      OR EXISTS (
          SELECT 1 FROM admin_permission m
          WHERE m.is_menu = 1 AND m.name = '帖子管理' AND m.route_url = '/community/post'
            AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
      )
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0
  AND parent_id IN (
      SELECT id FROM admin_permission
      WHERE is_menu = 1 AND name = '帖子模块' AND route_url = '/community/section'
  );
DELETE FROM admin_permission
WHERE is_menu = 1 AND name = '帖子模块' AND route_url = '/community/section';
