-- 评论配图(社区帖子评论可带图, 写入 my-storage 的 .bnc 地址)。
-- +goose Up
ALTER TABLE comment
    ADD COLUMN IF NOT EXISTS pics jsonb NOT NULL DEFAULT '[]';

COMMENT ON COLUMN comment.pics IS '评论配图 URL 列表';

-- +goose Down
ALTER TABLE comment DROP COLUMN IF EXISTS pics;
