-- 子后台管理员谷歌验证器(TOTP)
-- +goose Up
ALTER TABLE admin_user
    ADD COLUMN totp_secret varchar(64) NOT NULL DEFAULT '',
    ADD COLUMN totp_bound_at timestamptz;

COMMENT ON COLUMN admin_user.totp_secret IS '谷歌验证器 TOTP 密钥(已绑定才写入)';
COMMENT ON COLUMN admin_user.totp_bound_at IS '谷歌验证器绑定时间';

INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, '解绑验证器', '/backend/admins/{id}/totp', 'DELETE', 0, 5
FROM admin_permission m
WHERE m.is_menu = 1 AND m.route_url = '/system/admin'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.route_url = '/backend/admins/{id}/totp' AND x.method = 'DELETE'
  )
LIMIT 1;

-- +goose Down
DELETE FROM admin_permission WHERE route_url = '/backend/admins/{id}/totp' AND method = 'DELETE';
ALTER TABLE admin_user
    DROP COLUMN IF EXISTS totp_bound_at,
    DROP COLUMN IF EXISTS totp_secret;
