-- 分析页扩展维度的接口权限(挂在「分析页」菜单节点 id=2 下, 供非超管角色勾选)。
-- +goose Up
SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort) VALUES
    (2, '时段分布',     '/backend/stats/hour-dist',      'GET', 0, 5),
    (2, '设备分布',     '/backend/stats/devices',        'GET', 0, 6),
    (2, '内容库概览',   '/backend/stats/content',        'GET', 0, 7),
    (2, '金币流水构成', '/backend/stats/balance-scenes', 'GET', 0, 8);

-- +goose Down
DELETE FROM admin_permission WHERE is_menu = 0 AND parent_id = 2
   AND route_url IN ('/backend/stats/hour-dist', '/backend/stats/devices',
                     '/backend/stats/content', '/backend/stats/balance-scenes');
