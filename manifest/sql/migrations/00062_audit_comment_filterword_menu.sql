-- 审核管理二级菜单清空「投稿审核」, 改为「评论审核」+「禁用词库」。
-- 投稿审核页面仍保留(hide_in_menu=1), 接口权限不动, 只从侧栏拿掉。
-- +goose Up

UPDATE admin_permission
SET hide_in_menu = 1, sort = 99
WHERE is_menu = 1
  AND name = '投稿审核'
  AND route_url = '/audit/publish';

INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '评论审核', '/audit/comment', 'audit/comment', '', 'lucide:message-square', 1, 0, 0, 1
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '审核管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '评论审核' AND x.route_url = '/audit/comment'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
SELECT p.id, '禁用词库', '/audit/filterword', 'audit/filterword', '', 'lucide:ban', 1, 0, 0, 2
FROM admin_permission p
WHERE p.is_menu = 1 AND p.parent_id = 0 AND p.name = '审核管理'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 1 AND x.name = '禁用词库' AND x.route_url = '/audit/filterword'
  )
LIMIT 1;

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('评论列表', '/backend/comment',            'GET',  1),
    ('评论审核', '/backend/comment/{id}/audit', 'POST', 2)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '评论审核' AND m.route_url = '/audit/comment'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('禁词列表', '/backend/filterword',        'GET',    1),
    ('新增禁词', '/backend/filterword',        'POST',   2),
    ('删除禁词', '/backend/filterword/{id}',   'DELETE', 3)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '禁用词库' AND m.route_url = '/audit/filterword'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = v.route_url AND x.method = v.method AND x.is_menu = 0
  );

COMMENT ON COLUMN comment.status IS '0待审核 1已上墙 2已拒绝';

-- 已有「审核管理 / 投稿审核」的角色一并勾上新菜单和接口。
UPDATE admin_role r
SET permissions = trim(both ',' from r.permissions || ',' || ids.joined)
FROM (
    SELECT string_agg(p.id::text, ',') AS joined
    FROM admin_permission p
    WHERE (p.is_menu = 1 AND p.name IN ('评论审核', '禁用词库')
           AND p.route_url IN ('/audit/comment', '/audit/filterword'))
       OR p.parent_id IN (
           SELECT id FROM admin_permission
           WHERE is_menu = 1 AND name IN ('评论审核', '禁用词库')
             AND route_url IN ('/audit/comment', '/audit/filterword')
       )
) ids
WHERE r.code <> 'superadmin'
  AND ids.joined IS NOT NULL
  AND (
      (',' || r.permissions || ',') LIKE '%,304,%'
      OR (',' || r.permissions || ',') LIKE '%,123,%'
      OR EXISTS (
          SELECT 1 FROM admin_permission m
          WHERE m.is_menu = 1 AND m.name = '审核管理'
            AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
      )
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
UPDATE admin_permission
SET hide_in_menu = 0
WHERE is_menu = 1
  AND name = '投稿审核'
  AND route_url = '/audit/publish';

DELETE FROM admin_permission
WHERE is_menu = 0
  AND parent_id IN (
      SELECT id FROM admin_permission
      WHERE is_menu = 1 AND (
          (name = '评论审核' AND route_url = '/audit/comment')
          OR (name = '禁用词库' AND route_url = '/audit/filterword')
      )
  );
DELETE FROM admin_permission
WHERE is_menu = 1 AND (
    (name = '评论审核' AND route_url = '/audit/comment')
    OR (name = '禁用词库' AND route_url = '/audit/filterword')
);
COMMENT ON COLUMN comment.status IS '1正常 0隐藏';
