-- 漫画模块不再使用 comic_home，旧数据挂到权重最高的分类。
-- +goose Up
UPDATE comics_module AS m
SET position = 'cat_' || c.id
FROM (
    SELECT DISTINCT ON (site_id) site_id, id
    FROM comics_category
    WHERE status = 1
    ORDER BY site_id, rank DESC, id DESC
) AS c
WHERE m.position = 'comic_home'
  AND m.site_id = c.site_id;

ALTER TABLE comics_module ALTER COLUMN position SET DEFAULT '';
COMMENT ON COLUMN comics_module.position IS '展示位置: cat_{id}=H5 分类 Tab';
COMMENT ON TABLE comics_module IS '漫画分类运营模块(标题+样式+检索条件)';

-- +goose Down
UPDATE comics_module SET position = 'comic_home' WHERE position LIKE 'cat_%';
