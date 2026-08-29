-- 抖音用户上传：来源 + 待审/拒绝状态，审核仍在抖音列表。
-- status 新增 3待审核 4已拒绝；0草稿/2下架在抖音列表都归「已下架」。
-- 视频/动漫管理仍只用 0/1/2，不受影响。
-- +goose Up
ALTER TABLE video
    ADD COLUMN IF NOT EXISTS submit_source smallint NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reject_reason varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS audit_by bigint NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS audit_at timestamptz;

COMMENT ON COLUMN video.status IS '0草稿 1上架 2下架 3待审核 4已拒绝';
COMMENT ON COLUMN video.submit_source IS '0后台录入/媒资 1用户上传';
COMMENT ON COLUMN video.reject_reason IS '抖音用户上传拒绝原因';

CREATE INDEX IF NOT EXISTS idx_video_kind_source_status
    ON video (kind, submit_source, status);

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, '抖音审核', '/backend/videos/{id}/audit', 'POST', 0, 6
FROM admin_permission m
WHERE m.is_menu = 1 AND m.name = '抖音列表' AND m.route_url = '/douyin/list'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.parent_id = m.id AND x.route_url = '/backend/videos/{id}/audit' AND x.method = 'POST' AND x.is_menu = 0
  );

UPDATE admin_role r
SET permissions = r.permissions || ',' || p.id::text
FROM admin_permission p
WHERE p.is_menu = 0
  AND p.route_url = '/backend/videos/{id}/audit'
  AND p.method = 'POST'
  AND r.code <> 'superadmin'
  AND (',' || r.permissions || ',') NOT LIKE '%,' || p.id::text || ',%'
  AND EXISTS (
      SELECT 1 FROM admin_permission m
      WHERE m.is_menu = 1 AND m.name = '抖音列表' AND m.route_url = '/douyin/list'
        AND (',' || r.permissions || ',') LIKE '%,' || m.id::text || ',%'
  );

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission
WHERE is_menu = 0 AND route_url = '/backend/videos/{id}/audit' AND method = 'POST';

DROP INDEX IF EXISTS idx_video_kind_source_status;

ALTER TABLE video
    DROP COLUMN IF EXISTS audit_at,
    DROP COLUMN IF EXISTS audit_by,
    DROP COLUMN IF EXISTS reject_reason,
    DROP COLUMN IF EXISTS submit_source;

COMMENT ON COLUMN video.status IS '0草稿 1上架 2下架';
