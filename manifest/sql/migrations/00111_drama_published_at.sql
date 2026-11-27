-- 短剧增加可编辑的上架时间，供分类列表按上架时间倒序。
-- +goose Up
ALTER TABLE drama ADD COLUMN IF NOT EXISTS published_at timestamptz;
COMMENT ON COLUMN drama.published_at IS '上架时间(可后台编辑；待上架为空)';
UPDATE drama SET published_at = created_at WHERE status IN (1, 2) AND published_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_drama_published ON drama (site_id, status, published_at DESC NULLS LAST, id DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_drama_published;
ALTER TABLE drama DROP COLUMN IF EXISTS published_at;
