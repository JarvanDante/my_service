-- 已有帖子没有分类的，归到「最新」，H5 按分类能筛到。
-- +goose Up
UPDATE post SET category = '最新' WHERE category = '';

-- +goose Down
UPDATE post SET category = '' WHERE category = '最新';
