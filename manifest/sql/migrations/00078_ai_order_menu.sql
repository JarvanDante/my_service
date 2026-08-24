-- 子后台 AI管理 下增加「订单管理」。
-- 列表/重试复用 /backend/ai/tasks*，删除走 DELETE /backend/ai/tasks/{id}。
-- +goose Up

INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '订单管理', '/ai/order', 'ai/order', '', 'lucide:receipt', 1, 0, 0, 99
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = 'AI管理'
  AND NOT EXISTS (
    SELECT 1 FROM admin_permission x
    WHERE x.is_menu = 1 AND x.name = '订单管理' AND x.route_url = '/ai/order'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('订单列表', '/backend/ai/tasks',           'GET',    1),
    ('重试投递', '/backend/ai/tasks/{id}/retry', 'POST',   2),
    ('删除订单', '/backend/ai/tasks/{id}',       'DELETE', 3)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '订单管理' AND m.route_url = '/ai/order'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

UPDATE admin_role r
SET permissions = trim(both ',' from r.permissions || ',' || ids.joined)
FROM (
    SELECT string_agg(p.id::text, ',') AS joined
    FROM admin_permission p
    WHERE (p.is_menu = 1 AND p.name = '订单管理' AND p.route_url = '/ai/order')
       OR p.parent_id IN (
           SELECT id FROM admin_permission
           WHERE is_menu = 1 AND name = '订单管理' AND route_url = '/ai/order'
       )
) ids
WHERE r.code <> 'superadmin'
  AND ids.joined IS NOT NULL
  AND EXISTS (
      SELECT 1 FROM admin_permission m
      WHERE m.is_menu = 1 AND m.name = 'AI管理'
        AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0 AND parent_id IN (
    SELECT id FROM admin_permission
    WHERE is_menu = 1 AND name = '订单管理' AND route_url = '/ai/order'
);
DELETE FROM admin_permission
WHERE is_menu = 1 AND name = '订单管理' AND route_url = '/ai/order';
