-- 漫画模块检索条件 JSON，位置可挂到 H5 分类 Tab。
-- +goose Up
ALTER TABLE comics_module
    ADD COLUMN IF NOT EXISTS filter jsonb NOT NULL DEFAULT '{}';
COMMENT ON COLUMN comics_module.filter IS '检索条件 JSON: tag_id/cat_id/order/is_end/pay_type/ids/keywords/recommend，查 comics 表属性';
COMMENT ON COLUMN comics_module.position IS '展示位置: comic_home=漫画首页, cat_{id}=H5 分类 Tab';

-- +goose Down
ALTER TABLE comics_module DROP COLUMN IF EXISTS filter;
