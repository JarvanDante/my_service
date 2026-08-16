-- 视频对接 PaaS 媒资中心: 落 media_code, 支持按媒资 ID 查询与同步。
-- +goose Up
ALTER TABLE video ADD COLUMN IF NOT EXISTS media_code varchar(32) NOT NULL DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS uk_video_media_code ON video (media_code) WHERE media_code <> '';
COMMENT ON COLUMN video.media_code IS 'PaaS 媒资短码(my_media asset.code)';

SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT 84, '媒资同步', '/backend/videos/sync-media', 'POST', 0, 12
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission
    WHERE is_menu = 0 AND route_url = '/backend/videos/sync-media' AND method = 'POST'
);
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT 84, '媒资库列表', '/backend/media-assets', 'GET', 0, 13
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission
    WHERE is_menu = 0 AND route_url = '/backend/media-assets' AND method = 'GET'
);
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT 84, '选用媒资', '/backend/media-assets/{id}/pick', 'POST', 0, 14
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission
    WHERE is_menu = 0 AND route_url = '/backend/media-assets/{id}/pick' AND method = 'POST'
);

-- +goose Down
DELETE FROM admin_permission WHERE is_menu = 0 AND route_url IN (
    '/backend/videos/sync-media',
    '/backend/media-assets',
    '/backend/media-assets/{id}/pick'
);
DROP INDEX IF EXISTS uk_video_media_code;
ALTER TABLE video DROP COLUMN IF EXISTS media_code;
