-- 帖子增加分类字段；后台可改分类和浏览量。
-- +goose Up
ALTER TABLE post ADD COLUMN IF NOT EXISTS category varchar(32) NOT NULL DEFAULT '';
COMMENT ON COLUMN post.category IS '帖子分类名(对应 post_category.name)';
CREATE INDEX IF NOT EXISTS idx_post_category ON post (site_id, category);

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, '编辑帖子', '/backend/post/{id}', 'PUT', 0, 4
FROM admin_permission m
WHERE m.is_menu = 1 AND m.name = '帖子列表' AND m.route_url = '/community/post'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = '/backend/post/{id}' AND x.method = 'PUT' AND x.is_menu = 0
  );

UPDATE admin_role r
SET permissions = trim(both ',' from r.permissions || ',' || ids.joined)
FROM (
    SELECT string_agg(p.id::text, ',') AS joined
    FROM admin_permission p
    WHERE p.is_menu = 0 AND p.route_url = '/backend/post/{id}' AND p.method = 'PUT'
) ids
WHERE r.code <> 'superadmin'
  AND ids.joined IS NOT NULL
  AND EXISTS (
      SELECT 1 FROM admin_permission m
      WHERE m.is_menu = 1 AND m.name = '帖子列表' AND m.route_url = '/community/post'
        AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0 AND route_url = '/backend/post/{id}' AND method = 'PUT' AND name = '编辑帖子';
DROP INDEX IF EXISTS idx_post_category;
ALTER TABLE post DROP COLUMN IF EXISTS category;
