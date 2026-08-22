-- 帖子分类不需要「榜单」(kind=3)，去掉种子「热帖」。
-- +goose Up
DELETE FROM post_category WHERE kind = 3;

-- +goose Down
INSERT INTO post_category (site_id, name, kind, rank, status)
SELECT 1, '热帖', 3, 80, 1
WHERE NOT EXISTS (
    SELECT 1 FROM post_category WHERE site_id = 1 AND name = '热帖'
);
