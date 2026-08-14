-- 给没有主页背景的存量用户随机补一张内置渐变背景(/static/bg/bg01~16.jpg)。
-- 新注册由代码分配(logic/avatar.go randBackground), 本迁移只兜历史数据。
-- +goose Up
UPDATE users
SET bg_img = '/static/bg/bg' || lpad((1 + floor(random() * 16))::int::text, 2, '0') || '.jpg'
WHERE (bg_img = '' OR bg_img IS NULL) AND deleted_at IS NULL;

-- +goose Down
-- 数据回填不可逆, 无需回滚。
SELECT 1;
