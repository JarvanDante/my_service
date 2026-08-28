-- 抖音二级导航种子：热门走榜单栏位；不种「关注/喜欢」（用户行为流）。
-- +goose Up
INSERT INTO douyin_category (site_id, name, kind, rank, status)
SELECT 1, v.name, v.kind, v.rank, 1
FROM (VALUES
    ('热门', 3, 100),
    ('发现', 0, 90)
) AS v(name, kind, rank)
WHERE NOT EXISTS (
    SELECT 1 FROM douyin_category x WHERE x.site_id = 1 AND x.name = v.name
);

-- +goose Down
DELETE FROM douyin_category WHERE site_id = 1 AND name IN ('热门', '发现');
