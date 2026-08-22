-- 帖子增加排序权重，H5 前台按 rank 倒序。
-- +goose Up
ALTER TABLE post ADD COLUMN IF NOT EXISTS rank int NOT NULL DEFAULT 0;
COMMENT ON COLUMN post.rank IS '排序权重, 数值越大 H5 越靠前';
CREATE INDEX IF NOT EXISTS idx_post_feed_rank ON post (site_id, status, rank DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_post_feed_rank;
ALTER TABLE post DROP COLUMN IF EXISTS rank;
