-- 推广落地页地址，与分享页 share_url 分开。
-- 子后台渠道管理「复制链接」优先用这个，拼 ?source=渠道码。
-- +goose Up

INSERT INTO app_config (site_id, grp, key, value, remark, status)
SELECT 1, 'share', 'landing_url', '""', '推广落地页（渠道 source，不要带邀请参数）', 1
WHERE NOT EXISTS (
    SELECT 1 FROM app_config WHERE site_id = 1 AND key = 'landing_url'
);

-- +goose Down
DELETE FROM app_config WHERE site_id = 1 AND key = 'landing_url';
