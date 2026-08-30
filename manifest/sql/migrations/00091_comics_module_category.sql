-- 漫画模块增加分类检索, 与标签并列。
-- +goose Up
ALTER TABLE comics_module
    ADD COLUMN IF NOT EXISTS category_ids jsonb NOT NULL DEFAULT '[]';
COMMENT ON COLUMN comics_module.category_ids IS '检索分类ID列表, 空则不限分类';
COMMENT ON TABLE comics_module IS '漫画首页运营模块(标题+样式+分类/标签检索)';

-- +goose Down
ALTER TABLE comics_module DROP COLUMN IF EXISTS category_ids;
