-- 小说模块示例楼层, 对齐漫画模块的标题/样式写法。
-- 新更/推荐/小说榜不按作品分类筛, 直接按最新/热度/点赞出书。
-- +goose Up
INSERT INTO novel_category (site_id, name, kind, rank, status)
SELECT 1, '新更', 1, 100, 1
WHERE NOT EXISTS (
    SELECT 1 FROM novel_category WHERE site_id = 1 AND name = '新更'
);

INSERT INTO novel_module (site_id, name, position, style, icon, category_ids, tag_ids, filter, size, rank, status)
SELECT 1, v.name, 'cat_' || c.id, v.style, v.icon, '[]'::jsonb, '[]'::jsonb, v.filter::jsonb, v.size, v.rank, 1
FROM novel_category c
JOIN (VALUES
    ('新更',     '今日上新：新章已经准时送达 👉',           7, 1, '{"order":"new"}',                 6, 110),
    ('新更',     '热门连载：一口气看到停不下来',             1, 3, '{"order":"hot","is_end":"n"}',    5, 100),
    ('新更',     '随机发现：今晚读哪一本',                   6, 2, '{"order":"rand"}',                9,  90),
    ('推荐',     '编辑精选：本周必看文',                     7, 2, '{"order":"like"}',                9, 100),
    ('推荐',     '人气推荐：大家都在追',                     4, 3, '{"order":"hot"}',                 6,  90),
    ('小说榜',   '人气小说榜【点击更多查看榜单】',           7, 3, '{"order":"hot"}',                 9, 100),
    ('小说榜',   '完结佳作：一口气看完更过瘾',               4, 1, '{"order":"new","is_end":"y"}',    6,  80),
    ('都市生活', '都市情色：办公室里的秘密',                 7, 1, '{"order":"new"}',                 9,  80),
    ('经典激情', '经典激情：老味道更上头',                   4, 3, '{"order":"hot"}',                 6,  80),
    ('情欲性爱', '情欲开胃：今晚不想睡',                     1, 3, '{"order":"new"}',                 5,  80),
    ('日本情色', '日式情色：细腻又上头',                     6, 2, '{"order":"new"}',                 9,  80),
    ('西洋情色', '西洋情色：开放尺度更大',                   7, 1, '{"order":"new"}',                 9,  80)
) AS v(cat_name, name, style, icon, filter, size, rank)
  ON c.site_id = 1 AND c.name = v.cat_name
WHERE NOT EXISTS (
    SELECT 1 FROM novel_module m WHERE m.site_id = 1 AND m.name = v.name
);

-- +goose Down
DELETE FROM novel_module
WHERE site_id = 1 AND name IN (
    '今日上新：新章已经准时送达 👉',
    '热门连载：一口气看到停不下来',
    '随机发现：今晚读哪一本',
    '编辑精选：本周必看文',
    '人气推荐：大家都在追',
    '人气小说榜【点击更多查看榜单】',
    '完结佳作：一口气看完更过瘾',
    '都市情色：办公室里的秘密',
    '经典激情：老味道更上头',
    '情欲开胃：今晚不想睡',
    '日式情色：细腻又上头',
    '西洋情色：开放尺度更大'
);
