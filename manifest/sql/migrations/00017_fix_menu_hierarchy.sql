-- 扁平业务菜单 → 一级目录 + 二级页面菜单(与系统管理一致)。
-- 注意: goose 按分号拆语句, 禁止使用 DO $$ ... $$ 块。
-- +goose Up

-- 一级改为纯目录(清空 component)
UPDATE admin_permission SET component = '', name = '概览' WHERE id = 1 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 3 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 4 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 5 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 6 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 7 AND is_menu = 1;
UPDATE admin_permission SET component = '' WHERE id = 67 AND is_menu = 1 AND component = 'video/index';

-- 二级页面菜单
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
VALUES
  (79, 3,  '用户列表',   '/user/list',    'user/index',    '', 'lucide:list',          1, 0, 0, 1),
  (80, 4,  '财务中心',   '/finance/list', 'finance/index', '', 'lucide:wallet',        1, 0, 0, 1),
  (81, 5,  '兑换码管理', '/promo/list',   'promo/index',   '', 'lucide:ticket',        1, 0, 0, 1),
  (82, 6,  '成长中心',   '/growth/list',  'growth/index',  '', 'lucide:trophy',        1, 0, 0, 1),
  (83, 7,  '运营中心',   '/ops/list',     'ops/index',     '', 'lucide:megaphone',     1, 0, 0, 1)
ON CONFLICT (id) DO UPDATE SET
  parent_id = EXCLUDED.parent_id,
  name = EXCLUDED.name,
  route_url = EXCLUDED.route_url,
  component = EXCLUDED.component,
  icon = EXCLUDED.icon,
  is_menu = 1,
  sort = EXCLUDED.sort;

-- 视频二级(00016 可能已写入 id=84)
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
VALUES (84, 67, '视频列表', '/video/list', 'video/index', '', 'lucide:clapperboard', 1, 0, 0, 1)
ON CONFLICT (id) DO UPDATE SET
  parent_id = 67,
  name = '视频列表',
  route_url = '/video/list',
  component = 'video/index',
  is_menu = 1;

-- 接口权限挂到二级菜单下
UPDATE admin_permission SET parent_id = 79 WHERE id BETWEEN 16 AND 21;
UPDATE admin_permission SET parent_id = 80 WHERE id BETWEEN 22 AND 31;
UPDATE admin_permission SET parent_id = 81 WHERE id BETWEEN 32 AND 35;
UPDATE admin_permission SET parent_id = 82 WHERE id BETWEEN 36 AND 45;
UPDATE admin_permission SET parent_id = 83 WHERE id BETWEEN 46 AND 54;
UPDATE admin_permission SET parent_id = 84 WHERE id BETWEEN 68 AND 78;

-- 角色已勾选一级时补上二级菜单 id
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',79')
  WHERE (',' || permissions || ',') LIKE '%,3,%' AND (',' || permissions || ',') NOT LIKE '%,79,%';
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',80')
  WHERE (',' || permissions || ',') LIKE '%,4,%' AND (',' || permissions || ',') NOT LIKE '%,80,%';
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',81')
  WHERE (',' || permissions || ',') LIKE '%,5,%' AND (',' || permissions || ',') NOT LIKE '%,81,%';
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',82')
  WHERE (',' || permissions || ',') LIKE '%,6,%' AND (',' || permissions || ',') NOT LIKE '%,82,%';
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',83')
  WHERE (',' || permissions || ',') LIKE '%,7,%' AND (',' || permissions || ',') NOT LIKE '%,83,%';
UPDATE admin_role SET permissions = trim(both ',' from permissions || ',84')
  WHERE (',' || permissions || ',') LIKE '%,67,%' AND (',' || permissions || ',') NOT LIKE '%,84,%';

SELECT setval(pg_get_serial_sequence('admin_permission','id'), (SELECT MAX(id) FROM admin_permission));

-- +goose Down
UPDATE admin_permission SET parent_id = 3 WHERE id BETWEEN 16 AND 21;
UPDATE admin_permission SET parent_id = 4 WHERE id BETWEEN 22 AND 31;
UPDATE admin_permission SET parent_id = 5 WHERE id BETWEEN 32 AND 35;
UPDATE admin_permission SET parent_id = 6 WHERE id BETWEEN 36 AND 45;
UPDATE admin_permission SET parent_id = 7 WHERE id BETWEEN 46 AND 54;
UPDATE admin_permission SET component = 'user/index', route_url = '/user' WHERE id = 3;
UPDATE admin_permission SET component = 'finance/index', route_url = '/finance' WHERE id = 4;
UPDATE admin_permission SET component = 'promo/index', route_url = '/promo' WHERE id = 5;
UPDATE admin_permission SET component = 'growth/index', route_url = '/growth' WHERE id = 6;
UPDATE admin_permission SET component = 'ops/index', route_url = '/ops' WHERE id = 7;
DELETE FROM admin_permission WHERE id IN (79, 80, 81, 82, 83);
