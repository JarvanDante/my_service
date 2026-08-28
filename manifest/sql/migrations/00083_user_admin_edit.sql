-- 子后台用户编辑：资料可改；评论禁言字段先落库预留，社区评论上线后再接逻辑。
-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS comment_muted smallint NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS comment_mute_until timestamptz,
    ADD COLUMN IF NOT EXISTS violate_count int NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS today_comment_count int NOT NULL DEFAULT 0;

COMMENT ON COLUMN users.comment_muted IS '社区评论禁言 0否 1是(预留)';
COMMENT ON COLUMN users.comment_mute_until IS '禁言截止时间，空=未禁言或永久(预留)';
COMMENT ON COLUMN users.violate_count IS '评论违规次数(预留)';
COMMENT ON COLUMN users.today_comment_count IS '当日已评条数(预留)';

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT 79, '编辑用户', '/backend/users/{id}', 'PUT', 0, 8
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission
    WHERE is_menu = 0 AND parent_id = 79
      AND route_url = '/backend/users/{id}' AND method = 'PUT'
);

UPDATE admin_role r
SET permissions = r.permissions || ',' || p.id::text
FROM admin_permission p
WHERE p.is_menu = 0 AND p.parent_id = 79
  AND p.route_url = '/backend/users/{id}' AND p.method = 'PUT'
  AND r.code <> 'superadmin'
  AND (',' || r.permissions || ',') NOT LIKE '%,' || p.id::text || ',%'
  AND EXISTS (
      SELECT 1 FROM admin_permission m
      WHERE m.id = 16 AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
  );

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0 AND parent_id = 79
  AND route_url = '/backend/users/{id}' AND method = 'PUT';

ALTER TABLE users
    DROP COLUMN IF EXISTS comment_muted,
    DROP COLUMN IF EXISTS comment_mute_until,
    DROP COLUMN IF EXISTS violate_count,
    DROP COLUMN IF EXISTS today_comment_count;
