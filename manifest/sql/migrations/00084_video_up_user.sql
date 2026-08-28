-- 抖音作品绑定前台 UP 主。上架时校验用户存在、is_up=1、未禁用。
-- +goose Up
ALTER TABLE video
    ADD COLUMN IF NOT EXISTS up_user_id bigint NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_video_up_user ON video (up_user_id) WHERE up_user_id > 0;
COMMENT ON COLUMN video.up_user_id IS '前台 UP 主 users.id；抖音上架必填';

-- +goose Down
DROP INDEX IF EXISTS idx_video_up_user;
ALTER TABLE video DROP COLUMN IF EXISTS up_user_id;
