-- 给没有头像的存量用户随机补一张内置默认头像(/static/avatar/av01~48.png)。
-- 新注册用户由代码分配(logic/avatar.go), 本迁移只兜历史数据。幂等: 跑完不再有空头像。
-- +goose Up
UPDATE users
SET img = '/static/avatar/av' || lpad((1 + floor(random() * 48))::int::text, 2, '0') || '.png'
WHERE (img = '' OR img IS NULL) AND deleted_at IS NULL;

-- +goose Down
-- 数据回填不可逆(无法区分哪些是回填的), 也无需回滚。
SELECT 1;
