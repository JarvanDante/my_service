# 后台(运营管理)API 开发计划

> 前台(`/front`)已完成 P1~P7。本文规划**后台**接口。
> 门面:`/backend` = 运营后台。平台超管(`/manage`)已迁至 `my_manage_service`。
> 鉴权**独立于前台**:后台用独立的管理员登录 + Token + Casbin 权限,不复用前台用户 Token。
> 标注:✅ 用现成表 / 🆕 需新建表。

---

## 建议开发顺序(按优先级从上到下)

### B0 · 后台基础设施 【最先,必做】
先把「谁能登录后台、能干什么」立起来,后面所有接口都挂在这套鉴权下。

- 表 🆕：`admin_user`(管理员)、`admin_role`(角色)、`casbin_rule`(权限策略)、`admin_log`(操作日志)
- 接口:
  - `POST /backend/auth/login` 管理员登录(账号密码 → 后台 token)
  - `POST /backend/auth/logout` 退出
  - `GET  /backend/auth/info` 当前管理员信息 + 权限
  - `GET/POST/PUT/DELETE /backend/admins` 管理员管理
  - `GET/POST/PUT/DELETE /backend/roles` 角色 & 权限分配
- 中间件:`AdminAuth`(独立 token) + Casbin RBAC 校验 + 操作日志记录
- 说明:当前前后台共用一套用户 token,这一步把后台拆成独立管理员体系。

### B1 · 用户管理 【核心,高】
表 ✅ `users`、`user_balance_log`

- `GET  /backend/users` 用户列表(分页 + 搜索:用户名/手机/渠道/用户组/状态 + 时间范围)
- `GET  /backend/users/{id}` 用户详情 *(已示范)*
- `POST /backend/users/{id}/disable` 禁用/解禁 *(已示范)*
- `POST /backend/users/{id}/group` 调整用户组 / VIP 到期
- `POST /backend/users/{id}/balance` 调整金币/积分(写流水 + 记操作人)
- `GET  /backend/users/{id}/balance-logs` 用户余额流水

### B2 · 财务 / 订单 【高,营收相关】
表 ✅ `recharge_package`、`recharge_order`、`vip_package`、`vip_log`、`user_balance_log`

- 充值套餐 CRUD:`GET/POST/PUT/DELETE /backend/recharge-packages`
- VIP 套餐 CRUD:`/backend/vip-packages`
- `GET /backend/recharge-orders` 充值订单(状态/用户/时间筛选,对账)
- `GET /backend/balance-logs` 全站余额流水
- `POST /backend/pay/callback` 支付回调:标记订单已支付 + 到账
  —— **前台 P6 里 `recharge/do` 留的 TODO 在这落地**

### B3 · 兑换码管理 【中】
表 ✅ `user_code`、`user_code_log`

- `GET  /backend/codes` 兑换码列表
- `POST /backend/codes` 批量生成(类型/数量/面额/有效期)
- `POST /backend/codes/{id}/void` 作废
- `GET  /backend/code-logs` 兑换记录

### B4 · 用户组配置 【中】
表 🆕 `user_group`(用户组/等级定义:名称/折扣/权益)

- 用户组 CRUD:`/backend/user-groups`
- 说明:`users` 表已有 `group_id/group_name/group_rate` 等**快照字段**,这里补的是「组的定义表」,后台改组时写快照到用户。

### B5 · 成长体系配置 【中】
表 ✅ `user_task`、`user_task_log`、`user_sign`

- 任务 CRUD:`/backend/tasks`
- `GET /backend/task-logs` 任务完成记录
- `GET /backend/sign-stats` 签到统计

### B6 · 社交 / 推广运营 【中低】
表 ✅ `user_follow`、`user_share_log`

- `GET /backend/follows` 关注关系查询
- `GET /backend/share-stats` 推广/分享数据(拉新排行等)

### B7 · 消息 / 客服 【低】
表 ✅ `chat_conversation`、`chat_message`;🆕 `system_notice`(系统公告/推送)

- `GET  /backend/messages` 消息记录查询/监控
- `PUT  /backend/config/customer-url` 客服链接配置
- `POST /backend/push` 系统推送 / 群发公告

### B8 · 数据统计 / Dashboard 【最后,锦上添花】
- `GET /backend/stats/overview` 概览(新增/活跃/充值/DAU)
- `GET /backend/stats/...` 各维度趋势(充值趋势、留存、渠道分析)

---

## 需要新建的表汇总

| 阶段 | 新表 | 用途 |
|---|---|---|
| B0 | admin_user / admin_role / casbin_rule / admin_log | 管理员、角色、权限、操作日志 |
| B4 | user_group | 用户组/等级定义 |
| B7 | system_notice | 系统公告/推送 |

其余阶段全部复用前台已建的表。

---

## 落地建议

1. **先做 B0 → B1 → B2**:门立起来 + 用户管理 + 财务,后台就能支撑日常运营了。
2. 其余(B3~B8)按运营节奏排。
3. **鉴权拆分**:给后台单独的 `AdminAuth` 中间件(解析后台 token + Casbin),不要复用前台 `Auth`。平台超管见 `my_manage_service`。
4. **复用已有 service**:禁用用户、调组/调余额等,直接复用 `modules/user` 的 service;新的管理类(套餐、兑换码、管理员)建议新建 `modules/admin` 模块,别再全堆进 user。
5. 每个后台接口延续 module-first 分层:`api/backend/... → controller/backend → service → logic → dao`。

## 与前台的对应关系(速查)

| 前台阶段 | 对应后台 |
|---|---|
| P1 认证 / P2 资料 | B1 用户管理 |
| P4 兑换码 | B3 兑换码管理 |
| P5 签到/任务 | B5 成长配置 |
| P6 充值/VIP/兑换 | B2 财务/订单 |
| P3 社交 / P4 分享 | B6 社交/推广运营 |
| P7 私信 | B7 消息/客服 |
