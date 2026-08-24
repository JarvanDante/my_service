-- AI 换脸本地工人: 排队/处理超时后扫表失败并退款。
-- +goose Up
INSERT INTO app_config (site_id, grp, key, value, remark, status) VALUES
    (1, 'ai', 'ai_task_timeout_sec', '600', 'AI任务排队/处理超时秒数, 超时失败并退款', 1)
ON CONFLICT (site_id, key) DO NOTHING;

-- +goose Down
DELETE FROM app_config WHERE site_id = 1 AND key = 'ai_task_timeout_sec';
