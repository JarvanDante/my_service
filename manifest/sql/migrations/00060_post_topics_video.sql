-- 发帖：话题标签 + 可选视频。
-- +goose Up
ALTER TABLE post
    ADD COLUMN IF NOT EXISTS topics jsonb NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS video_url varchar(512) NOT NULL DEFAULT '';

COMMENT ON COLUMN post.topics IS '帖子话题/板块名数组';
COMMENT ON COLUMN post.video_url IS '可选视频地址，空=无';

INSERT INTO tag (site_id, content_type, name, rank, status) VALUES
    (1, 6, '日常分享', 100, 1),
    (1, 6, '最新公告', 95, 1),
    (1, 6, '活动公告', 90, 1),
    (1, 6, 'AI内容', 85, 1),
    (1, 6, '动漫情报站', 80, 1),
    (1, 6, '同人漫画', 75, 1),
    (1, 6, '新番预告', 70, 1),
    (1, 6, '自拍日常', 65, 1),
    (1, 6, '游戏补给站', 60, 1),
    (1, 6, '小说创作', 55, 1),
    (1, 6, '广场', 50, 1)
ON CONFLICT (site_id, content_type, name) DO NOTHING;

-- +goose Down
ALTER TABLE post
    DROP COLUMN IF EXISTS topics,
    DROP COLUMN IF EXISTS video_url;
