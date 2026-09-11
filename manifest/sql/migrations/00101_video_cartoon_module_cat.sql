-- 视频/动漫模块位置改挂分类 Tab，并加检索 JSON。
-- +goose Up
ALTER TABLE video_module
    ADD COLUMN IF NOT EXISTS filter jsonb NOT NULL DEFAULT '{}';
ALTER TABLE cartoon_module
    ADD COLUMN IF NOT EXISTS filter jsonb NOT NULL DEFAULT '{}';

UPDATE video_module AS m
SET position = 'cat_' || c.id
FROM (
    SELECT DISTINCT ON (site_id) site_id, id
    FROM video_category
    WHERE status = 1
    ORDER BY site_id, rank DESC, id DESC
) AS c
WHERE m.position = 'video_home'
  AND m.site_id = c.site_id;

UPDATE cartoon_module AS m
SET position = 'cat_' || c.id
FROM (
    SELECT DISTINCT ON (site_id) site_id, id
    FROM cartoon_category
    WHERE status = 1
    ORDER BY site_id, rank DESC, id DESC
) AS c
WHERE m.position = 'cartoon_home'
  AND m.site_id = c.site_id;

ALTER TABLE video_module ALTER COLUMN position SET DEFAULT '';
ALTER TABLE cartoon_module ALTER COLUMN position SET DEFAULT '';
COMMENT ON COLUMN video_module.position IS '展示位置: cat_{id}=H5 分类 Tab';
COMMENT ON COLUMN cartoon_module.position IS '展示位置: cat_{id}=H5 分类 Tab';
COMMENT ON COLUMN video_module.filter IS '检索条件 JSON: tag_id/cat_id/order/ids/keywords';
COMMENT ON COLUMN cartoon_module.filter IS '检索条件 JSON: tag_id/cat_id/order/ids/keywords';

-- +goose Down
ALTER TABLE video_module DROP COLUMN IF EXISTS filter;
ALTER TABLE cartoon_module DROP COLUMN IF EXISTS filter;
UPDATE video_module SET position = 'video_home' WHERE position LIKE 'cat_%';
UPDATE cartoon_module SET position = 'cartoon_home' WHERE position LIKE 'cat_%';
