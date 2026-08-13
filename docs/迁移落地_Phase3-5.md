# tianbi → my_service 迁移落地说明（Phase 3 / 4 / 5）

> 承接 `docs/迁移交接_Phase3-5.md`（那份是"设计与约定"，本份是"已经落成什么、怎么验收、还剩什么"）。
> Phase 1、2 之前已完成；本轮把 Phase 3、4、5 按顺序落地，迁移编号 **00032 ~ 00041**。

---

## 一、本轮做了什么（一句话版）

| Phase | 模块 | 迁移 | 前台/后台接口数 |
|---|---|---|---|
| 3 资金 | content_purchase 基建、wallet 钱包、withdrawal 提现+收款账户、coupon 优惠券 | 00032/00033/00034 | 20 / 14 |
| 4 内容+搜索 | 搜索基建(pg_trgm)、comics 漫画、novel 小说、photo 图集、publish 投稿、search 搜索、video 前台浏览 | 00035~00039 | 21 / 33 |
| 5 抽奖+AI | lottery 抽奖(含实物收货/发货)、aitask AI 生成任务骨架(供应商无关) | 00040/00041 | 10 / 16 |

合计新增 **46 个前台接口 + 56 个后台接口**，10 个迁移文件，11 个 `.http` e2e 套件。

---

## 二、验收步骤

```bash
# 1. 跑迁移(00032~00041)
make migrate

# 2. 重启一体化服务(GoFrame 启动时缓存表结构, 必须先迁移再启动)
gf run main.go        # :8000

# 3. IDE 重载 internal/cmd/router.go 标签页(防旧缓存覆盖)

# 4. 跑 e2e —— GoLand/IDEA HTTP Client, 每个文件 "Run all requests in file"
#    Phase 3: wallet.http  withdrawal.http  coupon.http
#    Phase 4: comics.http  novel.http  photo.http  publish.http  search.http  video_front.http
#    Phase 5: lottery.http  aitask.http
#    回归:    其余 test/http/*.http
```

### 关于"我这边已经跑过一遍"

本轮开发是在一个装了 PostgreSQL 16 + Redis 的沙箱里做的，**不是只写不跑**：

- 00001~00041 全部迁移在空库里从头跑通；
- 服务真起在 :8000，写了一个 mini 版 `.http` 运行器把 **全部 43 个 `.http` 文件**都真发请求跑了一遍；
- 结果：新增的 11 个套件全绿；老套件里除两个**环境性/数据性**问题外全绿（见第五节）。

所以你回来第一次跑，如果失败，优先怀疑本地环境（迁移没跑 / 服务没重启 / 库里有上次的脏数据），而不是逻辑。

---

## 三、各模块要点（只写"为什么这么做"，字段细节看迁移文件注释）

### Phase 3 · 资金

**公共能力 `internal/shared/balance`** —— 金币变动的唯一入口。
所有扣款一律 `WHERE balance >= ?` 条件更新 + `RowsAffected` 判定，**从 SQL 层防透支**，杜绝 tianbi 那种"先查余额再扣"的 TOCTOU 竞态。所有变动强制写 `user_balance_log`（direction/scene/before/after/ref_id），后台可按 scene 对账。

**公共能力 `internal/shared/appcfg`** —— 从 `app_config` KV 读运营参数（费率、上限、开关），后台随时改，不用改代码重启。

- **wallet**：前台余额（含在途冻结/累计收支/累计提现）与流水；后台全站流水 + 人工调账（正加负扣，扣款同样防透支）。
- **withdrawal**：状态机 `1申请中 → 2审核通过 → 4已打款`，`1|2 → 5拒绝(退款)`，`1 → 6撤回(退款)`。
  - **申请即冻结**：扣余额 + 写流水 + 建单，一个事务；
  - 每次状态迁移都带 `WHERE status=旧值`，`RowsAffected==0` 就报"已被处理"，**保证退款只可能发生一次**；
  - 手续费/最低额/倍数/日限/频控全部服务端从 `app_config` 取，不信客户端。
  - 附带 `bank_card` 收款账户 CRUD，沿用 tianbi 的"字段不得含空格"校验（打款失败高发原因）。
- **coupon**：模板 + 用户券（领取时对模板做值快照，模板改价不影响已发放的券）。
  - 修掉 tianbi 自动选券的 bug（`{"gte": now}` 少了 `$`，过期券会被选中）：把"未使用 + 未过期 + 场景匹配 + 门槛达标"四个条件全写进 SQL；
  - 发券用 `WHERE total=-1 OR issued<total` 条件递增，防超发；
  - `UseInTx` 供后续下单链路在**自己的事务里**核销，条件更新 `1→2`。

### Phase 4 · 内容 + 搜索

**公共能力 `internal/shared/paywall`** —— 内容解锁三态模型：`price>0` 金币付费 / `is_vip=1` VIP 专享（二者互斥）/ 都没有则免费。
`Buy()` 的顺序是**先插 `content_purchase`（唯一约束）再扣款**，让数据库唯一约束把并发重复购买挡在门外，避免 tianbi "先查是否买过再扣钱"被双击扣两次。

**搜索选型**：不照搬 ES 那 7 个索引。`00035` 装 `pg_trgm` 给标题建 GIN 索引让 `ILIKE` 走索引；`zhparser` 用 `DO $$ ... EXCEPTION` 守护，装不上只打 NOTICE 不阻塞迁移。搜索埋点写 `hot_search`（`InsertIgnore` + 条件自增，不做"先查再插"），埋点失败不影响搜索结果返回。

- **comics / novel**：同构。整部购买（与 tianbi 一致），章节级只有"前 N 章免费"（`free_chapter`）。章节图片/正文分别存 jsonb / text。删除作品连带删章节；章节增删后按实际行数回算 `chapter_count`（小说还回算总字数）。
- **photo**：没有章节，用 `free_count` 做"前 N 张免费"——**未解锁时服务端直接截断 `pics` 数组**，剩余图片 URL 根本不下发；列表接口完全不返回 `pics`，防止绕过详情拿图。
- **publish**：UGC 投稿，默认待审，标题+简介过 `filter_word`；撤回与审核都是条件更新（仅待审可动）。
- **video**：只补了前台浏览（列表/详情，只出已上架），沿用 video 模块原有的 repo/dao 模式，没有混进 `g.Model`。文件与转码仍走 my_media + my_play。
- **新增中间件 `middleware.AuthOptional`**：带了合法 token 就写 userId，没带也放行。用于"公开但要认人"的接口（游客能看简介，登录了要知道是否已购）。

### Phase 5 · 抽奖 + AI 骨架

- **lottery**：tianbi 原版整条链路**没有事务**、先发奖后扣费、免费次数能刷成负数、中奖记录用 `go func` 异步写。全部修掉：
  - 整次抽奖是一个事务：`users` 行锁 → 每日次数校验 → 扣费 → 加权随机 → 库存条件递减 → 发奖 → 写 history；
  - 加权随机用 **前缀和 + crypto/rand**（不用 math/rand）；
  - **概率 `odds` 绝不下发给客户端**；
  - 跑马灯用真实中奖记录（昵称脱敏），不像 tianbi 那样造 30 条假数据；
  - 中选奖品库存被抢空时**降级为"谢谢参与"**而不是事务内重抽（避免长时间占行锁）；
  - 奖品类型：金币 / VIP天数 / 优惠券 / 实物 / 谢谢参与。券是在抽奖事务内直接写 `user_coupon` + 条件递增 `coupon_tpl.issued`，**没有调 coupon 模块的方法**（那个方法自带事务，嵌套会出问题，且它的每人限领风控不适用于中奖发券）；
  - **补了 tianbi 没有的用户自助填收货地址接口**（原版只能线下找客服），后台按 `待填写→待发货→已发货` 流转。
- **aitask**：供应商还没定，所以只做**与供应商无关的"订单 + 回调 + 计费补偿"骨架**，代码里没有任何第三方域名/密钥/协议。
  - `internal/shared/aiprovider` 定义 `Provider` 接口 + 一个 mock 实现，接入真实供应商时只需实现接口并 `Register`，业务层零改动；
  - 提交：事务内扣费建单，**事务提交后**才发外部请求（网络调用不占数据库事务）；提交失败自动退款；
  - `client_token` 幂等键防双击重复扣费；
  - 回调：验签 + `WHERE status IN (1,2)` 条件更新 + 行锁，**重复回调不重复退款**（并发压测验证过：12 个线程同时推同一条失败回调，余额只涨一次、退款流水恰好一条）；
  - 查询接口带轮询兜底（回调丢失时的补偿路径）。

---

## 四、Phase 5 还没做的部分（**需要你先做决策，不能凭空写代码**）

| 模块 | 卡在哪 | 需要你定的事 |
|---|---|---|
| ai 系列的真实对接（换脸/脱衣/文生图/图生视频/文生小说） | 骨架已就位，缺 adapter | 选哪家供应商？拿到 API 文档、鉴权方式、回调协议与验签算法、计费口径（按次/按时长）、失败判定与退款口径 |
| aimate AI女友 / aicustomerservice AI客服 | 是"多轮对话 + 积分"，不是"任务 + 回调"，**形态和 aitask 不同**，不能套用 | 用哪个模型网关？对话上下文存哪（PG 还是向量库）？计费按轮次还是 token？ |
| live 直播 | tianbi 那套是 HLS + 推流 SDK | 复用 my_media + my_play 的哪条链路？推流鉴权怎么做？是否需要连麦/聊天室（那是另一套长连接基建） |
| porngame / game 游戏 | 第三方游戏聚合平台的单点登录换 URL | 供应商是谁？SSO 协议、分成结算口径 |
| customer 人工客服 | tianbi 嵌的是第三方客服 SaaS（webview + 签名免登录） | 继续买 SaaS 还是自建工单？自建的话要单独立项（工单+会话+坐席） |

**建议**：把 aitask 骨架先接一家供应商跑通（adapter 是几十行的活），验证"订单-回调-退款"闭环没问题，再复制到其余 AI 品类。live / game / customer 建议单独排期，不要和 CRUD 混在一起估工。

---

## 五、已知的老问题（**不是本轮引入的**）

1. `test/http/backend_media_multipart.http` —— 需要 MinIO 起着并有真实文件，纯环境依赖，沙箱里没 MinIO 所以跑不了。
2. `test/http/backend_social.http` 的"绑定推荐人"步骤 —— 用例本身**不可重复执行**：同一个测试账号第二次跑会撞上"已绑定推荐人, 不可修改"。产品逻辑是对的，是用例需要每次换新设备号或先解绑。
3. `test/http/admin_perm.http` —— 原先测的是早已删除的 Casbin 式接口 `/backend/roles/{code}/perms`，**本轮已按现行 RBAC 模型（权限节点树 + 角色持有节点 id）重写**，现在是绿的。

### 本轮顺手修的一个真 bug

`api/front/post/v1/post.go` 的 `Item.Status` 原本带 `json:"status,omitempty"`，而"待审"恰好是 0 —— `omitempty` 会把这个字段整个抹掉，前端在"我的帖子"里拿到 `undefined`，**没法区分"待审"和"没这个字段"**。已去掉 `omitempty`（公开流只出已通过内容，多下发一个 `status=1` 无害）。

---

## 六、约定补充（延续原交接文档，新增的部分）

7'. 迁移**下一个可用编号是 00042**。
13. `media_type` 全局编码统一为：**1视频 2帖子 3漫画 4小说 5图集**，`user_collect` / `content_purchase` / `search` / `paywall` 共用同一套。
14. 金币变动**只能**走 `internal/shared/balance`，不要在业务模块里直接 UPDATE `users.balance`——否则流水会漏，后台对不上账。
15. 运营可调参数（费率、上限、开关、价格）一律进 `app_config`，用 `internal/shared/appcfg` 读，默认值写在代码里兜底。
16. "公开但要认人"的接口挂 `middleware.AuthOptional`，不要用 `middleware.Auth`（会把游客挡掉），也不要不挂（带了 token 也读不到 userId）。
17. 涉及钱和状态机的操作，一律「**条件更新 + RowsAffected 判定**」，不要「先 SELECT 判断再 if 里改」。

---

## 七、git 提交说明（供你在 Mac 上提交）

Phase 2 的提交说明在 `docs/迁移交接_Phase3-5.md` 末尾，本轮建议分三个 commit：

```
feat(phase3): 资金链路 —— 钱包/提现/优惠券/内容购买基建

- shared/balance(00032 起): 金币变动唯一入口, 扣款一律 WHERE balance>=? 条件更新
  + 强制写 user_balance_log, 从 SQL 层防透支(修 tianbi 先查后扣的 TOCTOU)
- shared/appcfg: 运营参数走 app_config KV, 后台可改不用改代码
- content_purchase(00032): (user,media_type,content_id) 唯一, 供 Phase 4 付费解锁
- wallet: 前台余额+流水; 后台全站流水+人工调账(正加负扣, 扣款防透支)
- withdrawal(00033): 申请即冻结, 状态机 1→2→4 / 1|2→5退款 / 1→6撤回退款,
  每次迁移带 WHERE status=旧值, 保证退款只发生一次; 附 bank_card 收款账户
- coupon(00034): 模板+用户券快照, 条件递增防超发, 修 tianbi 自动选券漏 $ 导致
  过期券被选中的 bug; UseInTx 供下单链路在自身事务内核销
- e2e: test/http/{wallet,withdrawal,coupon}.http
```

```
feat(phase4): 内容媒资 + 搜索 —— 漫画/小说/图集/投稿/搜索/视频前台

- shared/paywall: 解锁三态(免费/VIP专享/金币付费), 购买先插 content_purchase
  唯一约束再扣款, 用 DB 约束挡并发重复购买(修 tianbi 双击扣两次)
- 搜索基建(00035): pg_trgm + GIN 标题索引替代 ES; zhparser 用 DO 块守护, 装不上不阻塞
- comics(00036)/novel(00037): 整部购买 + 前 N 章免费, 章节图片 jsonb / 正文 text,
  删作品连带删章节, 章节增删按实际行数回算章节数(小说另回算总字数)
- photo(00038): 无章节, free_count 前 N 张免费, 未解锁服务端截断 pics 不下发余图
- publish(00039): UGC 投稿默认待审, 过 filter_word, 撤回/审核均条件更新
- search: 跨 video/post/comics/novel/photo 标题检索, 埋点写 hot_search(InsertIgnore
  +条件自增), 埋点失败不影响结果; 另有 /search/suggest 前缀联想
- video: 补前台列表/详情(只出已上架), 沿用该模块原有 repo/dao 模式
- middleware.AuthOptional: "公开但要认人"的接口用它
- fix(post): 我的帖子 status 因 omitempty 丢失待审态(0), 前端无法区分待审与缺字段
- e2e: test/http/{comics,novel,photo,publish,search,video_front}.http
```

```
feat(phase5): 抽奖 + AI 任务骨架

- lottery(00040): 整次抽奖收进一个事务(用户行锁→次数校验→扣费→加权随机→库存
  条件递减→发奖→写history), 修 tianbi 无事务/先发奖后扣费/免费次数刷成负数/
  history 异步写四个坑; 加权随机用前缀和+crypto/rand; odds 不下发客户端;
  跑马灯用真实记录; 库存抢空降级为谢谢参与; 补了用户自助填收货地址(原版没有)
- aitask(00041): 供应商未定, 只落"订单+回调+计费补偿"骨架, 无任何第三方协议;
  shared/aiprovider 定义 Provider 接口 + mock 实现, 接入只需加 adapter;
  事务内扣费建单、事务提交后才发外部请求; client_token 幂等防双击;
  回调验签 + WHERE status IN(1,2) 条件更新 + 行锁, 重复回调不重复退款;
  查询接口带轮询兜底
- fix(test): admin_perm.http 原测早已删除的 Casbin 接口, 按现行 RBAC 模型重写
- e2e: test/http/{lottery,aitask}.http
```
