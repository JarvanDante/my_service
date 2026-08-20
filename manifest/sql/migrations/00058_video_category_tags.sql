-- 视频作品挂本站分类/标签。分类逗号分隔, 标签 jsonb 数组, 对齐漫画。
-- +goose Up
ALTER TABLE video
    ADD COLUMN IF NOT EXISTS category varchar(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS tags jsonb NOT NULL DEFAULT '[]'::jsonb;
COMMENT ON COLUMN video.category IS '本站分类, 逗号分隔(电影,短剧)';
COMMENT ON COLUMN video.tags IS '本站标签名数组';

-- +goose Down
ALTER TABLE video
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS tags;
