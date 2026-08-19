-- 漫画对接 PaaS 媒资中心: 落 media_code; 分类由子站自定, 同步后待上架。
-- +goose Up
ALTER TABLE comics ADD COLUMN IF NOT EXISTS media_code varchar(32) NOT NULL DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS uk_comics_media_code ON comics (media_code) WHERE media_code <> '';
COMMENT ON COLUMN comics.media_code IS 'PaaS 媒资短码(my_media asset.code)；空=本站自建';

-- 挂到「漫画列表」菜单下(id 可能因环境不同)
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT m.id, v.name, v.route_url, v.method, 0, v.sort
FROM admin_permission m
CROSS JOIN (VALUES
    ('漫画媒资库',   '/backend/media-comics',           'GET',  12),
    ('选用漫画媒资', '/backend/media-comics/{id}/pick', 'POST', 13)
) AS v(name, route_url, method, sort)
WHERE m.is_menu = 1 AND m.name = '漫画列表'
  AND NOT EXISTS (
      SELECT 1 FROM admin_permission x
      WHERE x.is_menu = 0 AND x.route_url = v.route_url AND x.method = v.method
  );

-- +goose Down
DELETE FROM admin_permission WHERE is_menu = 0 AND route_url IN (
    '/backend/media-comics',
    '/backend/media-comics/{id}/pick'
);
DROP INDEX IF EXISTS uk_comics_media_code;
ALTER TABLE comics DROP COLUMN IF EXISTS media_code;
