-- 菜单调整: 按用户要求把「内容管理」拆开, 视频/漫画/小说/图集/帖子各占一个一级菜单
-- (参考公司后台的组织方式, 但去掉那边"标签/分类/评论散落在每个内容菜单里"的重复)。
--
-- 调整后的一级结构:
--   数据概览 / 用户管理 / 视频管理 / 漫画管理 / 小说管理 / 图集管理 / 社区管理 /
--   审核管理 / 资金管理 / 运营管理 / AI管理 / 系统设置
--
-- 说明:
--   * 投稿审核单独提为「审核管理」一级(对齐公司后台的"审核管理", 以后评论审核也挂这);
--   * 标签是跨内容共用的一张表(content_type 区分), 不像公司那样每个内容菜单里各放一个
--     "标签管理"—— 单独一份挂在运营管理下, 避免同一页面在侧栏出现四次;
--   * 只动 route_url(菜单路径), component 不变, 前端零改动;
--   * 想再微调归属, 直接在 系统设置→菜单权限 页面里改 parent_id 即可, 不用再写迁移。
-- +goose Up

-- ---------- 1. 新增一级目录(视频/漫画/小说/图集/审核) ----------
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort) VALUES
    (300, 0, '视频管理', '/video',  '', '', 'lucide:video',      1,0,0,20),
    (301, 0, '漫画管理', '/comics', '', '', 'lucide:book-image', 1,0,0,21),
    (302, 0, '小说管理', '/novel',  '', '', 'lucide:book-open',  1,0,0,22),
    (303, 0, '图集管理', '/photo',  '', '', 'lucide:images',     1,0,0,23),
    (304, 0, '审核管理', '/audit',  '', '', 'lucide:file-check', 1,0,0,35);

-- ---------- 2. 内容页面各回各家 ----------
UPDATE admin_permission SET parent_id=300, route_url='/video/list',  name='视频列表', sort=1 WHERE id=84;
UPDATE admin_permission SET parent_id=301, route_url='/comics/list', name='漫画列表', sort=1 WHERE id=120;
UPDATE admin_permission SET parent_id=302, route_url='/novel/list',  name='小说列表', sort=1 WHERE id=121;
UPDATE admin_permission SET parent_id=303, route_url='/photo/list',  name='图集列表', sort=1 WHERE id=122;
UPDATE admin_permission SET parent_id=304, route_url='/audit/publish', name='投稿审核', sort=1 WHERE id=123;
-- 标签管理 → 运营管理
UPDATE admin_permission SET parent_id=7, route_url='/ops/tag', name='标签管理', icon='lucide:tags', sort=9 WHERE id=124;

-- ---------- 3. 删除腾空的「内容管理」目录 ----------
DELETE FROM admin_permission WHERE id = 100;

-- 再次对齐 IDENTITY 序列(本迁移又手工指定了 id)
SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- +goose Down
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort)
VALUES (100, 0, '内容管理', '/content', '', '', 'lucide:clapperboard', 1, 0, 0, 20);
UPDATE admin_permission SET parent_id=100, route_url='/content/video',   name='视频管理', sort=1 WHERE id=84;
UPDATE admin_permission SET parent_id=100, route_url='/content/comics',  name='漫画管理', sort=2 WHERE id=120;
UPDATE admin_permission SET parent_id=100, route_url='/content/novel',   name='小说管理', sort=3 WHERE id=121;
UPDATE admin_permission SET parent_id=100, route_url='/content/photo',   name='图集管理', sort=4 WHERE id=122;
UPDATE admin_permission SET parent_id=100, route_url='/content/publish', name='投稿审核', sort=5 WHERE id=123;
UPDATE admin_permission SET parent_id=100, route_url='/content/tag',     name='标签管理', sort=6 WHERE id=124;
DELETE FROM admin_permission WHERE id IN (300,301,302,303,304);
