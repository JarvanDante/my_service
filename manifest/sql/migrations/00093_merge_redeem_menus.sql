-- 运营管理下「推广兑换码」与「兑换码」合并为一个「兑换码」菜单。
-- 接口权限挂到新菜单下；两套表/API 仍保留，后台用页签切换。
-- +goose Up
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '兑换码', '/ops/codes', 'ops/codes', '', 'lucide:ticket', 1, 0, 0, 3
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '运营管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '兑换码' AND x.route_url = '/ops/codes'
  )
LIMIT 1;

UPDATE admin_permission SET parent_id = (
    SELECT id FROM admin_permission
    WHERE is_menu = 1 AND name = '兑换码' AND route_url = '/ops/codes'
    LIMIT 1
)
WHERE is_menu = 0
  AND parent_id IN (
      SELECT id FROM admin_permission
      WHERE is_menu = 1 AND (
          (name = '推广兑换码' AND route_url = '/ops/promo')
          OR (name = '兑换码' AND route_url = '/ops/redeemcode')
      )
  );

UPDATE admin_role r
SET permissions = trim(both ',' from r.permissions || ',' || m.id::text)
FROM (
    SELECT id FROM admin_permission
    WHERE is_menu = 1 AND name = '兑换码' AND route_url = '/ops/codes'
    LIMIT 1
) m
WHERE r.code <> 'superadmin'
  AND (',' || r.permissions || ',') NOT LIKE '%,' || m.id::text || ',%'
  AND EXISTS (
      SELECT 1 FROM admin_permission p
      WHERE p.is_menu = 1
        AND (
            (p.name = '推广兑换码' AND p.route_url = '/ops/promo')
            OR (p.name = '兑换码' AND p.route_url = '/ops/redeemcode')
        )
        AND (',' || r.permissions || ',') LIKE '%,' || p.id::text || ',%'
  );

UPDATE admin_role r
SET permissions = coalesce((
    SELECT string_agg(tok, ',')
    FROM unnest(string_to_array(r.permissions, ',')) AS tok
    WHERE tok <> ''
      AND tok ~ '^[0-9]+$'
      AND tok::bigint NOT IN (
          SELECT p.id
          FROM admin_permission p
          WHERE p.is_menu = 1 AND (
              (p.name = '推广兑换码' AND p.route_url = '/ops/promo')
              OR (p.name = '兑换码' AND p.route_url = '/ops/redeemcode')
          )
      )
), '');

DELETE FROM admin_permission
WHERE is_menu = 1 AND (
    (name = '推广兑换码' AND route_url = '/ops/promo')
    OR (name = '兑换码' AND route_url = '/ops/redeemcode')
);

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '推广兑换码', '/ops/promo', 'promo/index', '', 'lucide:ticket', 1, 0, 0, 3
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '运营管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '推广兑换码' AND x.route_url = '/ops/promo'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '兑换码', '/ops/redeemcode', 'ops/redeemcode', '', 'lucide:gift', 1, 0, 0, 4
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '运营管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '兑换码' AND x.route_url = '/ops/redeemcode'
  )
LIMIT 1;

UPDATE admin_permission SET parent_id = (
    SELECT id FROM admin_permission
    WHERE is_menu = 1 AND name = '推广兑换码' AND route_url = '/ops/promo'
    LIMIT 1
)
WHERE is_menu = 0 AND route_url LIKE '/backend/codes%';

UPDATE admin_permission SET parent_id = (
    SELECT id FROM admin_permission
    WHERE is_menu = 1 AND name = '兑换码' AND route_url = '/ops/redeemcode'
    LIMIT 1
)
WHERE is_menu = 0 AND route_url LIKE '/backend/redeemcode%';

DELETE FROM admin_permission
WHERE is_menu = 1 AND name = '兑换码' AND route_url = '/ops/codes';

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));
