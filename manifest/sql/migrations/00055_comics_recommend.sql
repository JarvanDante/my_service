-- 漫画推荐位: 仅 is_recommend=1 的作品进入 H5「推荐」栏，顺序看 rank。
-- +goose Up
ALTER TABLE comics ADD COLUMN IF NOT EXISTS is_recommend smallint NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_comics_recommend ON comics (site_id, status, is_recommend, rank DESC);
COMMENT ON COLUMN comics.is_recommend IS '1=进入H5推荐栏';

-- +goose Down
DROP INDEX IF EXISTS idx_comics_recommend;
ALTER TABLE comics DROP COLUMN IF EXISTS is_recommend;
