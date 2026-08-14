-- 用户列表批量冻结/解冻的接口权限(挂在「用户列表」菜单节点 id=79 下)。
-- 前台 /front/v1/user/restore 是公开接口, 不需要权限行。
-- +goose Up
SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort)
SELECT 79, '批量冻结/解冻', '/backend/users/batch-disable', 'POST', 0, 7
WHERE NOT EXISTS (
    SELECT 1 FROM admin_permission
    WHERE is_menu = 0 AND parent_id = 79
      AND route_url = '/backend/users/batch-disable' AND method = 'POST'
);

-- +goose Down
DELETE FROM admin_permission WHERE is_menu = 0 AND parent_id = 79
   AND route_url = '/backend/users/batch-disable' AND method = 'POST';
