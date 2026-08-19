-- 漫画分类改为多选: category 存逗号分隔(韩漫,日漫)，筛选用数组包含。
-- +goose Up
ALTER TABLE comics ALTER COLUMN category TYPE varchar(128);

-- +goose Down
ALTER TABLE comics ALTER COLUMN category TYPE varchar(32);
