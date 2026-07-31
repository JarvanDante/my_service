-- B5: 任务配置化 —— 奖励/状态/排序入表(默认值与原硬编码行为一致)。
-- +goose Up
ALTER TABLE user_task
    ADD COLUMN reward numeric(14,2) NOT NULL DEFAULT 10,  -- 完成奖励积分(原硬编码 10)
    ADD COLUMN status smallint      NOT NULL DEFAULT 1,   -- 1上架 0下线
    ADD COLUMN sort   int           NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE user_task DROP COLUMN IF EXISTS sort;
ALTER TABLE user_task DROP COLUMN IF EXISTS status;
ALTER TABLE user_task DROP COLUMN IF EXISTS reward;
