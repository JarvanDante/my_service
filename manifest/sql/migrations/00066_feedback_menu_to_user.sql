-- 「意见反馈」挪到用户管理下，改名为「用户反馈」，插在会员等级和站内消息之间。
-- 页面 component 仍是 community/feedback，不改文件。
-- +goose Up
UPDATE admin_permission
SET parent_id = (
        SELECT id FROM admin_permission
        WHERE is_menu = 1 AND parent_id = 0 AND name = '用户管理'
        LIMIT 1
    ),
    name = '用户反馈',
    route_url = '/user/feedback',
    sort = 3
WHERE is_menu = 1
  AND name = '意见反馈'
  AND route_url = '/community/feedback';

UPDATE admin_permission
SET sort = 4
WHERE is_menu = 1 AND name = '站内消息' AND route_url = '/user/message';

-- +goose Down
UPDATE admin_permission
SET parent_id = (
        SELECT id FROM admin_permission
        WHERE is_menu = 1 AND parent_id = 0 AND name = '社区管理'
        LIMIT 1
    ),
    name = '意见反馈',
    route_url = '/community/feedback',
    sort = 3
WHERE is_menu = 1
  AND name = '用户反馈'
  AND route_url = '/user/feedback';

UPDATE admin_permission
SET sort = 3
WHERE is_menu = 1 AND name = '站内消息' AND route_url = '/user/message';
