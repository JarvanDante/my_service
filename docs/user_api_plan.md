# User 模块 · 接口开发计划

> 参照 `mh51-php`(Hyperf + MongoDB)的 `UserController` 动作,迁移到 my_service
> (GoFrame 模块化单体 + 多门面)。前台走 `/front`(限流),后台走 `/backend`(鉴权),
> 总后台走 `/manage`(鉴权)。业务逻辑一份共享,门面只是接口面。

## 数据表(已生成)

`manifest/sql/migrations/` 下:

- `00001_create_users.sql` — `users` 主表
- `00002_create_user_satellites.sql` — `user_follow` / `user_code` / `user_code_log` /
  `user_sign` / `user_task` / `user_task_log` / `user_balance_log`

原项目是 MongoDB 单文档大宽表,这里按关系型拆成主表 + 卫星表:主表放账号/资料/资产快照/
状态,流水与关系类拆表(关注、兑换码、签到、任务、余额流水),便于对账与索引。

---

## 阶段计划(前台 /front)

### P0 · 基础设施(先做)
- [ ] users 及卫星表迁移落库(`goose up`)
- [ ] `gf gen dao` 生成 model/entity/dao
- [ ] JWT 鉴权中间件落地(替换 `shared/middleware.Auth` 的 TODO)
- [ ] 统一响应结构 + 错误码(`shared/errcode`)
- [ ] 设备/IP 解析工具

### P1 · 账号与认证 `[高]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| POST | /front/user/login | loginAction | 设备登录/游客自动注册, 返回 token |
| POST | /front/user/logout | — | 退出(可选) |
| POST | /front/user/token/refresh | — | 刷新 token(可选) |
| POST | /front/user/bind-phone | bindPhoneAction | 绑定手机 |

### P2 · 个人资料 `[高]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| GET  | /front/user/info | infoAction | 当前用户信息 |
| GET  | /front/user/home/{id} | homeAction | 他人主页 |
| POST | /front/user/update | updateInfoAction | 改昵称/头像/签名/性别等 |
| GET  | /front/user/images | imagesAction | 用户图片 |
| GET  | /front/user/find | findByAccountAction | 按账号查找用户 |

### P3 · 社交关系 `[中]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| POST | /front/user/follow | doFollowAction | 关注/取关(切换) |
| GET  | /front/user/follows | followAction | 我的关注列表 |
| GET  | /front/user/fans | fansAction | 我的粉丝列表 |

### P4 · 推广 / 兑换码 `[中]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| POST | /front/user/bind-parent | bindParentAction | 绑定推荐人 |
| POST | /front/user/bind-code | bindCodeAction | 绑定邀请码 |
| POST | /front/user/code/redeem | doCodeAction | 使用兑换码(加组/加币) |
| GET  | /front/user/code/logs | codeLogsAction | 兑换记录 |
| GET  | /front/user/share | shareInfoAction | 分享信息/海报 |
| GET  | /front/user/share/logs | shareLogsAction | 分享记录 |

### P5 · 成长体系(签到/任务) `[中]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| POST | /front/user/sign | doDaySignAction | 每日签到 |
| GET  | /front/user/tasks | taskAction | 任务列表 |
| POST | /front/user/task/do | doTaskAction | 完成任务领奖 |
| GET  | /front/user/task/logs | doTaskLogAction | 任务记录 |
| GET  | /front/user/up | upAction | 升级/成长信息 |

### P6 · 资产(充值/VIP/兑换) `[高]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| GET  | /front/user/recharge | rechargeAction | 充值套餐 |
| POST | /front/user/recharge/do | doRechargeAction | 发起充值(接支付) |
| GET  | /front/user/vip | vipAction | VIP 套餐 |
| POST | /front/user/vip/do | doVipAction | 开通/续费 VIP |
| GET  | /front/user/vip/logs | vipLogsAction | VIP 记录 |
| GET  | /front/user/exchange | exchangeInfoAction | 兑换信息 |
| POST | /front/user/exchange/do | doExchangeAction | 积分/金币兑换 |

### P7 · 消息 / 客服 `[低, 建议独立 message 模块]`
| 方法 | 路径 | 参照动作 | 说明 |
|---|---|---|---|
| GET  | /front/user/chats | chatsAction | 会话列表 |
| GET  | /front/user/chat/messages | chatMessagesAction | 会话消息 |
| POST | /front/user/chat/send | sendMessageAction | 发消息 |
| POST | /front/user/chat/del | delChatAction | 删会话 |
| GET  | /front/user/customer-url | getCustomerUrlAction | 客服链接 |

---

## 后台 /backend(运营管理, 参照 Backend/UserController)
| 方法 | 路径 | 说明 |
|---|---|---|
| GET  | /backend/users | 用户列表(分页/搜索/筛选) |
| GET  | /backend/users/{id} | 用户详情 |
| POST | /backend/users/{id}/disable | 禁用/解禁 |
| POST | /backend/users/{id}/group | 调整用户组/VIP |
| POST | /backend/users/{id}/balance | 调整余额(写 user_balance_log) |
| GET  | /backend/user-codes | 兑换码列表 |
| POST | /backend/user-codes | 生成兑换码 |
| GET  | /backend/user-tasks | 任务配置 |

## 总后台 /manage(平台超管)
| 方法 | 路径 | 说明 |
|---|---|---|
| GET  | /manage/users | 跨渠道用户总览 |
| POST | /manage/users/{id}/reset | 高危操作(重置/清算) |
| GET  | /manage/user-groups | 用户组/等级体系配置 |

---

## 落地顺序建议

1. **P0 → P1 → P2**:先把「登录 + 拿到自己的信息 + 改资料」跑通,这是一切的基础。
2. **P6 资产**:营收相关,优先级高,但依赖支付接入,可与 P3/P4 并行设计。
3. **P3 社交 / P4 推广 / P5 成长**:按运营节奏排。
4. **P7 消息**:字段和 user 关系较弱,建议单独开 `modules/message`,不要塞进 user 模块。
5. 后台/总后台接口:每做完一个前台能力,顺手补对应的后台管理接口(复用同一 service)。

## 开发约定(呼应 module-first 脚手架)

- 每个接口: `api/<门面>/user/v1` 定 req/res → `controller/<门面>` 适配 → 统一调
  `modules/user/service.IUser` → `logic` 实现 → `domain` 放规则 → `dao` 落库。
- 涉及金币/余额变动一律走 `user_balance_log` 记账, 并在事务内更新 `users.balance`。
- 关注/签到/任务这类高频写, 计数字段(fans/follow/login_num)用增量更新, 避免全表统计。
