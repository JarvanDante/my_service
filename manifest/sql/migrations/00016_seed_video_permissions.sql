-- 视频管理: 一级目录 + 二级列表 + 视频/媒体接口权限。
-- +goose Up
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
VALUES
 (67, 0,  '视频管理', '/video',      '',            '', 'lucide:video',         1, 0, 0, 55),
 (84, 67, '视频列表', '/video/list', 'video/index', '', 'lucide:clapperboard',  1, 0, 0, 1)
ON CONFLICT (id) DO UPDATE SET
  parent_id = EXCLUDED.parent_id,
  name = EXCLUDED.name,
  route_url = EXCLUDED.route_url,
  component = EXCLUDED.component,
  icon = EXCLUDED.icon,
  is_menu = EXCLUDED.is_menu,
  sort = EXCLUDED.sort;

INSERT INTO admin_permission (id, parent_id, name, route_url, method, is_menu, sort)
VALUES
 (68, 84, '视频列表',     '/backend/videos',                  'GET',    0, 1),
 (69, 84, '新增视频',     '/backend/videos',                  'POST',   0, 2),
 (70, 84, '编辑视频',     '/backend/videos/{id}',             'PUT',    0, 3),
 (71, 84, '删除视频',     '/backend/videos/{id}',             'DELETE', 0, 4),
 (72, 84, '视频上下架',   '/backend/videos/{id}/status',      'PUT',    0, 5),
 (73, 84, '上传媒体',     '/backend/media/upload',            'POST',   0, 6),
 (74, 84, '分片初始化',   '/backend/media/multipart/init',    'POST',   0, 7),
 (75, 84, '分片预签名',   '/backend/media/multipart/presign', 'POST',   0, 8),
 (76, 84, '分片已传列表', '/backend/media/multipart/parts',   'GET',    0, 9),
 (77, 84, '分片合并完成', '/backend/media/multipart/complete','POST',   0, 10),
 (78, 84, '分片取消',     '/backend/media/multipart/abort',   'POST',   0, 11)
ON CONFLICT (id) DO UPDATE SET
  parent_id = EXCLUDED.parent_id,
  name = EXCLUDED.name,
  route_url = EXCLUDED.route_url,
  method = EXCLUDED.method,
  is_menu = EXCLUDED.is_menu,
  sort = EXCLUDED.sort;

-- 兼容旧扁平结构残留
UPDATE admin_permission SET component = '' WHERE id = 67 AND component = 'video/index';
UPDATE admin_permission SET parent_id = 84 WHERE id BETWEEN 68 AND 78 AND parent_id = 67;

SELECT setval(pg_get_serial_sequence('admin_permission','id'), (SELECT MAX(id) FROM admin_permission));

-- +goose Down
DELETE FROM admin_permission WHERE id BETWEEN 67 AND 78 OR id = 84;
