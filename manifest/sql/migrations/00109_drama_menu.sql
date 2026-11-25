-- 子后台“短剧管理”菜单与接口权限（列表 / 分类），幂等可重复执行。
-- +goose Up
INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT 0, '短剧管理', '/drama', '', '', 'lucide:clapperboard', 1, 0, 0, 25
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission x WHERE x.is_menu = 1 AND x.parent_id = 0 AND x.route_url = '/drama'
);

INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, v.name, v.route_url, v.component, '', v.icon, 1, 0, 0, v.sort
FROM admin_permission p
CROSS JOIN (VALUES
    ('短剧列表', '/drama/list',     'content/drama',          'lucide:clapperboard', 1),
    ('短剧分类', '/drama/category', 'content/drama-category', 'lucide:folders',      2)
) AS v(name, route_url, component, icon, sort)
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.route_url = '/drama'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x WHERE x.is_menu = 1 AND x.route_url = v.route_url
  );

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('短剧列表',     '/backend/dramas',                    'GET',    1),
    ('编辑短剧',     '/backend/dramas/{id}',               'PUT',    2),
    ('删除短剧',     '/backend/dramas/{id}',               'DELETE', 3),
    ('短剧上下架',   '/backend/dramas/{id}/audit',         'POST',   4),
    ('剧集列表',     '/backend/dramas/{id}/episodes',      'GET',    5),
    ('编辑剧集',     '/backend/drama-episodes/{id}',       'PUT',    6),
    ('同步剧集',     '/backend/dramas/{id}/sync-episodes', 'POST',   7),
    ('短剧媒资库',   '/backend/media-dramas',              'GET',    8),
    ('选用短剧媒资', '/backend/media-dramas/{id}/pick',    'POST',   9),
    ('短剧分类选项', '/backend/drama-categories',          'GET',    10)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.route_url = '/drama/list'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('分类列表', '/backend/drama-categories',      'GET',    1),
    ('新增分类', '/backend/drama-categories',      'POST',   2),
    ('编辑分类', '/backend/drama-categories/{id}', 'PUT',    3),
    ('删除分类', '/backend/drama-categories/{id}', 'DELETE', 4)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.route_url = '/drama/category'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

-- +goose Down
DELETE FROM admin_permission
WHERE parent_id IN (SELECT id FROM admin_permission WHERE is_menu = 1 AND route_url IN ('/drama/list', '/drama/category'));
DELETE FROM admin_permission WHERE is_menu = 1 AND route_url IN ('/drama/list', '/drama/category');
DELETE FROM admin_permission WHERE is_menu = 1 AND parent_id = 0 AND route_url = '/drama';
