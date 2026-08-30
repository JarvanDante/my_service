-- 视频/动漫模块增加分类检索, 对齐漫画模块。
-- +goose Up
ALTER TABLE video_module
    ADD COLUMN IF NOT EXISTS category_ids jsonb NOT NULL DEFAULT '[]';
ALTER TABLE cartoon_module
    ADD COLUMN IF NOT EXISTS category_ids jsonb NOT NULL DEFAULT '[]';
COMMENT ON COLUMN video_module.category_ids IS '检索分类ID列表, 空则不限分类';
COMMENT ON COLUMN cartoon_module.category_ids IS '检索分类ID列表, 空则不限分类';

-- +goose Down
ALTER TABLE video_module DROP COLUMN IF EXISTS category_ids;
ALTER TABLE cartoon_module DROP COLUMN IF EXISTS category_ids;
