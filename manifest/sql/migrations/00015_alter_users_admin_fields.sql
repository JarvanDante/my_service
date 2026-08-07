-- 对齐公司后台用户列表: UP/有效用户/分成比例。
-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_up          smallint NOT NULL DEFAULT 0,  -- 是否 UP 主
    ADD COLUMN IF NOT EXISTS is_valid       smallint NOT NULL DEFAULT 0,  -- 是否有效用户
    ADD COLUMN IF NOT EXISTS movie_fee_rate int      NOT NULL DEFAULT 0,  -- 视频分成(%)
    ADD COLUMN IF NOT EXISTS post_fee_rate  int      NOT NULL DEFAULT 0;  -- 帖子分成(%)

CREATE INDEX IF NOT EXISTS idx_users_is_up ON users (is_up);
CREATE INDEX IF NOT EXISTS idx_users_has_buy ON users (has_buy);

-- +goose Down
DROP INDEX IF EXISTS idx_users_has_buy;
DROP INDEX IF EXISTS idx_users_is_up;
ALTER TABLE users
    DROP COLUMN IF EXISTS post_fee_rate,
    DROP COLUMN IF EXISTS movie_fee_rate,
    DROP COLUMN IF EXISTS is_valid,
    DROP COLUMN IF EXISTS is_up;
