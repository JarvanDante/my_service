-- 搜索基建(替代 tianbi 的 7 个 ES 索引)。
--
-- 选型说明:
--   ES 那套「多索引 + IK 分词 + 回查 Mongo」在 PG 里没必要照搬。这里分两层:
--   1) pg_trgm(contrib 自带, 一定装得上): 给标题建 GIN trgm 索引, 让 `ILIKE '%词%'`
--      走索引而不是全表扫, 中文短词场景够用, 是默认路径;
--   2) zhparser(需要单独编译安装): 装上了就能建中文分词的 tsvector 列做真正的全文检索。
--      这里用 DO 块守护, 装不上只打 NOTICE, 不阻塞迁移 —— 迁移必须在任何环境都能跑过。
-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'pg_trgm 未能安装(%), 搜索将退化为无索引 ILIKE', SQLERRM;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS zhparser;
    -- 装上了才建中文全文检索配置
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'chinese_zh') THEN
        CREATE TEXT SEARCH CONFIGURATION chinese_zh (PARSER = zhparser);
        ALTER TEXT SEARCH CONFIGURATION chinese_zh ADD MAPPING FOR n,v,a,i,e,l,j WITH simple;
    END IF;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'zhparser 不可用(%), 中文分词全文检索跳过, 走 pg_trgm ILIKE 路径', SQLERRM;
END
$$;
-- +goose StatementEnd

-- 已有内容表的标题 trgm 索引(表不存在时跳过, 保证迁移在任何顺序下都不报错)
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_trgm') THEN
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'video') THEN
            CREATE INDEX IF NOT EXISTS idx_video_title_trgm ON video USING GIN (title gin_trgm_ops);
        END IF;
        IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'post') THEN
            CREATE INDEX IF NOT EXISTS idx_post_title_trgm ON post USING GIN (title gin_trgm_ops);
        END IF;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_post_title_trgm;
DROP INDEX IF EXISTS idx_video_title_trgm;
