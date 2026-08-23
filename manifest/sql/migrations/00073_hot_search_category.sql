-- 热搜词按内容分类投放: 漫画/动漫/小说/短剧/视频 各自一套榜。
-- 空 category 表示通用, 前台某分类不足 10 条时用通用词补齐。
-- +goose Up
ALTER TABLE hot_search
    ADD COLUMN IF NOT EXISTS category varchar(32) NOT NULL DEFAULT '';
COMMENT ON COLUMN hot_search.category IS '投放分类: 空=通用 comic/cartoon/novel/short/video';

ALTER TABLE hot_search DROP CONSTRAINT IF EXISTS hot_search_site_id_keyword_key;
DROP INDEX IF EXISTS hot_search_site_id_keyword_key;
CREATE UNIQUE INDEX IF NOT EXISTS uk_hot_search_site_cat_kw ON hot_search (site_id, category, keyword);

DROP INDEX IF EXISTS idx_hot_search_top;
CREATE INDEX idx_hot_search_top ON hot_search (site_id, category, status, heat DESC, search_count DESC);

-- +goose Down
DROP INDEX IF EXISTS uk_hot_search_site_cat_kw;
DROP INDEX IF EXISTS idx_hot_search_top;
ALTER TABLE hot_search DROP COLUMN IF EXISTS category;
CREATE UNIQUE INDEX IF NOT EXISTS hot_search_site_id_keyword_key ON hot_search (site_id, keyword);
CREATE INDEX idx_hot_search_top ON hot_search (site_id, status, heat DESC, search_count DESC);
