-- 收藏/点赞只写了 user_collect，内容表 like_count 一直是 0。按已有记录回填。
-- 漫画/小说/图集: like_count = 收藏数(op_type=1)；帖子: like_count = 点赞数(op_type=2)。
-- +goose Up
UPDATE comics c
SET like_count = COALESCE((
    SELECT COUNT(*) FROM user_collect uc
    WHERE uc.site_id = c.site_id AND uc.media_type = 3 AND uc.op_type = 1 AND uc.content_id = c.id
), 0);

-- +goose StatementBegin
DO $$
BEGIN
    IF to_regclass('public.novel') IS NOT NULL THEN
        UPDATE novel n
        SET like_count = COALESCE((
            SELECT COUNT(*) FROM user_collect uc
            WHERE uc.site_id = n.site_id AND uc.media_type = 4 AND uc.op_type = 1 AND uc.content_id = n.id
        ), 0);
    END IF;
    IF to_regclass('public.photo') IS NOT NULL THEN
        UPDATE photo p
        SET like_count = COALESCE((
            SELECT COUNT(*) FROM user_collect uc
            WHERE uc.site_id = p.site_id AND uc.media_type = 5 AND uc.op_type = 1 AND uc.content_id = p.id
        ), 0);
    END IF;
    IF to_regclass('public.post') IS NOT NULL THEN
        UPDATE post p
        SET like_count = COALESCE((
            SELECT COUNT(*) FROM user_collect uc
            WHERE uc.site_id = p.site_id AND uc.media_type = 2 AND uc.op_type = 2 AND uc.content_id = p.id
        ), 0);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
SELECT 1;
