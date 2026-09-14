-- 小说对齐漫画列表: 推荐位 + 分类多选(逗号分隔)。
-- +goose Up
ALTER TABLE novel ALTER COLUMN category TYPE varchar(128);
ALTER TABLE novel ADD COLUMN IF NOT EXISTS is_recommend smallint NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_novel_recommend ON novel (site_id, status, is_recommend, rank DESC);
COMMENT ON COLUMN novel.is_recommend IS '1=进入H5推荐栏';
COMMENT ON COLUMN novel.category IS '本站分类, 逗号分隔可多选';

-- +goose Down
DROP INDEX IF EXISTS idx_novel_recommend;
ALTER TABLE novel DROP COLUMN IF EXISTS is_recommend;
ALTER TABLE novel ALTER COLUMN category TYPE varchar(32);
