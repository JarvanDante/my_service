-- 子后台菜单重构 + Phase1~5 全量后台接口挂载。
--
-- 背景: 菜单是 accessMode=backend 动态下发的(GET /backend/auth/menus 读 admin_permission),
--       所以「加菜单」= 往这张表插数据, 前端不需要改 router。
--
-- 这次把原来平铺的 8 个顶级目录(用户/财务/兑换码/成长/运营/视频/系统…)按业务域收敛成 8 组:
--   数据概览 / 用户管理 / 内容管理 / 社区管理 / 资金管理 / 运营管理 / AI管理 / 系统设置
-- 归类原则: 「同一个人一天里会连着点的页面放一起」——
--   * 视频/漫画/小说/图集/投稿/标签 都是内容运营在维护 → 合成"内容管理";
--   * 帖子/反馈 是 UGC 与用户声音 → "社区管理";
--   * 兑换码/商品兑换/推广应用/排行热搜/抽奖/运营配置 都是活动运营 → "运营管理";
--   * 钱包/提现/优惠券/充值套餐 都涉及钱, 必须同一组便于对账 → "资金管理"。
-- 原顶级目录「兑换码」「用户组与成长」「视频管理」降级成二级页面, 空目录删除。
--
-- id 分配: 已有数据占用 1~84, 本次菜单节点用 100~，接口权限用 200+。
-- +goose Up

-- ---------- 1. 新增顶级目录 ----------
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort) VALUES
    (100, 0, '内容管理', '/content',   '', '', 'lucide:clapperboard', 1, 0, 0, 20),
    (101, 0, '社区管理', '/community', '', '', 'lucide:message-square', 1, 0, 0, 30),
    (102, 0, 'AI管理',   '/ai',        '', '', 'lucide:sparkles',      1, 0, 0, 60);

-- ---------- 2. 调整既有目录的名称/图标/排序 ----------
UPDATE admin_permission SET name='数据概览', icon='lucide:layout-dashboard', sort=-1  WHERE id=1;
UPDATE admin_permission SET name='用户管理', icon='lucide:users',            sort=10  WHERE id=3;
UPDATE admin_permission SET name='资金管理', icon='lucide:wallet',           sort=40  WHERE id=4;
UPDATE admin_permission SET name='运营管理', icon='lucide:megaphone',        sort=50  WHERE id=7;
UPDATE admin_permission SET name='系统设置', icon='lucide:settings',         sort=90  WHERE id=8;

-- ---------- 3. 既有二级页面重新归位 ----------
-- 成长中心: 顶级"用户组与成长"(6) → 用户管理下
UPDATE admin_permission SET parent_id=3, route_url='/user/growth', name='用户组与成长', icon='lucide:trophy', sort=2 WHERE id=82;
-- 视频列表: 顶级"视频管理"(67) → 内容管理下
UPDATE admin_permission SET parent_id=100, route_url='/content/video', name='视频管理', icon='lucide:video', sort=1 WHERE id=84;
-- 兑换码管理: 顶级"兑换码"(5) → 运营管理下(这是老 promo 模块: 通用码+分享统计)
UPDATE admin_permission SET parent_id=7, route_url='/ops/promo', name='推广兑换码', icon='lucide:ticket', sort=3 WHERE id=81;
-- 其余保持
UPDATE admin_permission SET name='用户列表', icon='lucide:list',   sort=1 WHERE id=79;
UPDATE admin_permission SET name='财务中心', icon='lucide:landmark', sort=1 WHERE id=80;
UPDATE admin_permission SET name='运营中心', icon='lucide:megaphone', sort=1 WHERE id=83;

-- 接口权限节点跟着它们的页面走, 不用动(parent 是页面 id, 没变)

-- ---------- 4. 删除腾空的顶级目录 ----------
DELETE FROM admin_permission WHERE id IN (5, 6, 67);

-- ---------- 5. 新增二级页面 ----------
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort) VALUES
    -- 用户管理
    (110, 3,   '站内消息',   '/user/message',       'user/message',        '', 'lucide:mail',           1,0,0,3),
    -- 内容管理
    (120, 100, '漫画管理',   '/content/comics',     'content/comics',      '', 'lucide:book-image',     1,0,0,2),
    (121, 100, '小说管理',   '/content/novel',      'content/novel',       '', 'lucide:book-open',      1,0,0,3),
    (122, 100, '图集管理',   '/content/photo',      'content/photo',       '', 'lucide:images',         1,0,0,4),
    (123, 100, '投稿审核',   '/content/publish',    'content/publish',     '', 'lucide:file-check',     1,0,0,5),
    (124, 100, '标签管理',   '/content/tag',        'content/tag',         '', 'lucide:tags',           1,0,0,6),
    -- 社区管理
    (130, 101, '帖子管理',   '/community/post',     'community/post',      '', 'lucide:notebook-pen',   1,0,0,1),
    (131, 101, '意见反馈',   '/community/feedback', 'community/feedback',  '', 'lucide:message-circle', 1,0,0,2),
    -- 资金管理
    (140, 4,   '金币钱包',   '/finance/wallet',     'finance/wallet',      '', 'lucide:coins',          1,0,0,2),
    (141, 4,   '提现审核',   '/finance/withdrawal', 'finance/withdrawal',  '', 'lucide:banknote',       1,0,0,3),
    (142, 4,   '优惠券',     '/finance/coupon',     'finance/coupon',      '', 'lucide:ticket-percent', 1,0,0,4),
    -- 运营管理
    (150, 7,   '运营配置',   '/ops/config',         'ops/config',          '', 'lucide:settings-2',     1,0,0,2),
    (151, 7,   '兑换码',     '/ops/redeemcode',     'ops/redeemcode',      '', 'lucide:gift',           1,0,0,4),
    (152, 7,   '商品兑换',   '/ops/redeem-goods',   'ops/redeem-goods',    '', 'lucide:shopping-bag',   1,0,0,5),
    (153, 7,   '推广应用',   '/ops/application',    'ops/application',     '', 'lucide:app-window',     1,0,0,6),
    (154, 7,   '排行热搜',   '/ops/rank',           'ops/rank',            '', 'lucide:trending-up',    1,0,0,7),
    (155, 7,   '抽奖活动',   '/ops/lottery',        'ops/lottery',         '', 'lucide:dices',          1,0,0,8),
    -- AI 管理
    (160, 102, 'AI模板',     '/ai/template',        'ai/template',         '', 'lucide:layout-template',1,0,0,1),
    (161, 102, 'AI任务',     '/ai/task',            'ai/task',             '', 'lucide:cpu',            1,0,0,2),
    -- 系统设置
    (170, 8,   '基础配置',   '/system/config',      'system/config',       '', 'lucide:sliders',        1,0,0,0);

-- 手工指定过 id 的插入不会推进 IDENTITY 序列, 必须先对齐再做自增插入,
-- 否则下面不带 id 的 INSERT 会从 1 开始生成、直接撞上已有主键。
SELECT setval(pg_get_serial_sequence('admin_permission', 'id'),
              (SELECT MAX(id) FROM admin_permission));

-- ---------- 6. 新页面对应的接口权限(is_menu=0), 供非超管角色勾选 ----------
INSERT INTO admin_permission (parent_id, name, route_url, method, is_menu, sort) VALUES
    -- 站内消息
    (110,'消息列表','/backend/message','GET',0,1),
    (110,'新增消息','/backend/message','POST',0,2),
    (110,'编辑消息','/backend/message/{id}','PUT',0,3),
    (110,'删除消息','/backend/message/{id}','DELETE',0,4),
    -- 漫画
    (120,'漫画列表','/backend/comics','GET',0,1),
    (120,'新增漫画','/backend/comics','POST',0,2),
    (120,'编辑漫画','/backend/comics/{id}','PUT',0,3),
    (120,'删除漫画','/backend/comics/{id}','DELETE',0,4),
    (120,'漫画上下架','/backend/comics/{id}/audit','POST',0,5),
    (120,'章节列表','/backend/comics/{id}/chapters','GET',0,6),
    (120,'新增章节','/backend/comics/{id}/chapters','POST',0,7),
    (120,'编辑章节','/backend/comics-chapters/{id}','PUT',0,8),
    (120,'删除章节','/backend/comics-chapters/{id}','DELETE',0,9),
    -- 小说
    (121,'小说列表','/backend/novels','GET',0,1),
    (121,'新增小说','/backend/novels','POST',0,2),
    (121,'编辑小说','/backend/novels/{id}','PUT',0,3),
    (121,'删除小说','/backend/novels/{id}','DELETE',0,4),
    (121,'小说上下架','/backend/novels/{id}/audit','POST',0,5),
    (121,'章节列表','/backend/novels/{id}/chapters','GET',0,6),
    (121,'新增章节','/backend/novels/{id}/chapters','POST',0,7),
    (121,'编辑章节','/backend/novel-chapters/{id}','PUT',0,8),
    (121,'删除章节','/backend/novel-chapters/{id}','DELETE',0,9),
    -- 图集
    (122,'图集列表','/backend/photos','GET',0,1),
    (122,'新增图集','/backend/photos','POST',0,2),
    (122,'编辑图集','/backend/photos/{id}','PUT',0,3),
    (122,'删除图集','/backend/photos/{id}','DELETE',0,4),
    (122,'图集上下架','/backend/photos/{id}/audit','POST',0,5),
    -- 投稿
    (123,'投稿列表','/backend/publishes','GET',0,1),
    (123,'投稿审核','/backend/publishes/{id}/audit','POST',0,2),
    -- 标签
    (124,'标签列表','/backend/tag','GET',0,1),
    (124,'新增标签','/backend/tag','POST',0,2),
    (124,'编辑标签','/backend/tag/{id}','PUT',0,3),
    (124,'删除标签','/backend/tag/{id}','DELETE',0,4),
    -- 帖子
    (130,'帖子列表','/backend/post','GET',0,1),
    (130,'帖子审核','/backend/post/{id}/audit','POST',0,2),
    (130,'删除帖子','/backend/post/{id}','DELETE',0,3),
    -- 反馈
    (131,'反馈列表','/backend/feedback','GET',0,1),
    (131,'处理反馈','/backend/feedback/{id}/handle','POST',0,2),
    -- 钱包
    (140,'金币流水','/backend/wallet/logs','GET',0,1),
    (140,'人工调账','/backend/wallet/adjust','POST',0,2),
    -- 提现
    (141,'提现列表','/backend/withdrawals','GET',0,1),
    (141,'提现审核','/backend/withdrawals/{id}/audit','POST',0,2),
    (141,'标记打款','/backend/withdrawals/{id}/mark-paid','POST',0,3),
    (141,'打款失败退款','/backend/withdrawals/{id}/refund','POST',0,4),
    -- 优惠券
    (142,'券模板列表','/backend/coupons','GET',0,1),
    (142,'新增券模板','/backend/coupons','POST',0,2),
    (142,'编辑券模板','/backend/coupons/{id}','PUT',0,3),
    (142,'删除券模板','/backend/coupons/{id}','DELETE',0,4),
    (142,'发放优惠券','/backend/coupons/{id}/grant','POST',0,5),
    (142,'用户券记录','/backend/coupons/users','GET',0,6),
    -- 运营配置
    (150,'公告列表','/backend/announcement','GET',0,1),
    (150,'新增公告','/backend/announcement','POST',0,2),
    (150,'编辑公告','/backend/announcement/{id}','PUT',0,3),
    (150,'删除公告','/backend/announcement/{id}','DELETE',0,4),
    (150,'跳转位列表','/backend/jumptab','GET',0,5),
    (150,'新增跳转位','/backend/jumptab','POST',0,6),
    (150,'编辑跳转位','/backend/jumptab/{id}','PUT',0,7),
    (150,'删除跳转位','/backend/jumptab/{id}','DELETE',0,8),
    (150,'敏感词列表','/backend/filterword','GET',0,9),
    (150,'新增敏感词','/backend/filterword','POST',0,10),
    (150,'删除敏感词','/backend/filterword/{id}','DELETE',0,11),
    -- 兑换码(新)
    (151,'兑换码列表','/backend/redeemcode','GET',0,1),
    (151,'新增兑换码','/backend/redeemcode','POST',0,2),
    (151,'编辑兑换码','/backend/redeemcode/{id}','PUT',0,3),
    (151,'删除兑换码','/backend/redeemcode/{id}','DELETE',0,4),
    (151,'核销记录','/backend/redeemcode/records','GET',0,5),
    -- 商品兑换
    (152,'商品列表','/backend/redeem-goods','GET',0,1),
    (152,'新增商品','/backend/redeem-goods','POST',0,2),
    (152,'编辑商品','/backend/redeem-goods/{id}','PUT',0,3),
    (152,'删除商品','/backend/redeem-goods/{id}','DELETE',0,4),
    (152,'兑换记录','/backend/redeem-goods/orders','GET',0,5),
    -- 推广应用
    (153,'应用列表','/backend/application','GET',0,1),
    (153,'新增应用','/backend/application','POST',0,2),
    (153,'编辑应用','/backend/application/{id}','PUT',0,3),
    (153,'删除应用','/backend/application/{id}','DELETE',0,4),
    -- 排行热搜
    (154,'热搜列表','/backend/hotsearch','GET',0,1),
    (154,'新增热搜','/backend/hotsearch','POST',0,2),
    (154,'编辑热搜','/backend/hotsearch/{id}','PUT',0,3),
    (154,'删除热搜','/backend/hotsearch/{id}','DELETE',0,4),
    (154,'刷新排行缓存','/backend/rank/refresh','POST',0,5),
    -- 抽奖
    (155,'活动列表','/backend/lottery/activities','GET',0,1),
    (155,'新增活动','/backend/lottery/activities','POST',0,2),
    (155,'编辑活动','/backend/lottery/activities/{id}','PUT',0,3),
    (155,'删除活动','/backend/lottery/activities/{id}','DELETE',0,4),
    (155,'奖品列表','/backend/lottery/prizes','GET',0,5),
    (155,'新增奖品','/backend/lottery/prizes','POST',0,6),
    (155,'编辑奖品','/backend/lottery/prizes/{id}','PUT',0,7),
    (155,'删除奖品','/backend/lottery/prizes/{id}','DELETE',0,8),
    (155,'中奖记录','/backend/lottery/histories','GET',0,9),
    (155,'收货单列表','/backend/lottery/addresses','GET',0,10),
    (155,'标记发货','/backend/lottery/addresses/{id}/ship','POST',0,11),
    -- AI 模板
    (160,'模板列表','/backend/ai/templates','GET',0,1),
    (160,'新增模板','/backend/ai/templates','POST',0,2),
    (160,'编辑模板','/backend/ai/templates/{id}','PUT',0,3),
    (160,'删除模板','/backend/ai/templates/{id}','DELETE',0,4),
    -- AI 任务
    (161,'任务列表','/backend/ai/tasks','GET',0,1),
    (161,'重新提交','/backend/ai/tasks/{id}/retry','POST',0,2),
    (161,'人工退款','/backend/ai/tasks/{id}/refund','POST',0,3),
    -- 基础配置
    (170,'配置列表','/backend/config','GET',0,1),
    (170,'新增配置','/backend/config','POST',0,2),
    (170,'编辑配置','/backend/config/{id}','PUT',0,3),
    (170,'删除配置','/backend/config/{id}','DELETE',0,4);

-- +goose Down
DELETE FROM admin_permission WHERE parent_id IN (110,120,121,122,123,124,130,131,140,141,142,150,151,152,153,154,155,160,161,170);
DELETE FROM admin_permission WHERE id BETWEEN 110 AND 170;
DELETE FROM admin_permission WHERE id IN (100,101,102);
INSERT INTO admin_permission (id, parent_id, name, route_url, component, method, icon, is_menu, hide_in_menu, affix_tab, sort) VALUES
    (5, 0, '兑换码',       '/promo',  '', '', 'lucide:ticket', 1,0,0,30),
    (6, 0, '用户组与成长', '/growth', '', '', 'lucide:trophy', 1,0,0,40),
    (67,0, '视频管理',     '/video',  '', '', 'lucide:video',  1,0,0,55);
UPDATE admin_permission SET parent_id=6,  route_url='/growth/list', name='成长中心'   WHERE id=82;
UPDATE admin_permission SET parent_id=67, route_url='/video/list',  name='视频列表'   WHERE id=84;
UPDATE admin_permission SET parent_id=5,  route_url='/promo/list',  name='兑换码管理' WHERE id=81;
