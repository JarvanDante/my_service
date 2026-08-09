# 总后台(/manage 控制面)API 开发计划

> **已迁出**:控制面实现见仓库 `my_manage_service`。本仓库(`my_service`)已移除站点侧 `/manage` 门面与 `manageapi` 二进制,下文仅作历史设计备忘。
>
> 定位:**控制面(Control Plane)**——管理所有商户/站点/域名/配置下发。
> 与站点后台(/backend, 数据面)的关系:一站一部署一库,站点只管自己的业务;总后台跨站点,独立部署、独立库。
>
> 两个已定决策:
> 1. 控制面数据存**独立库 `my_manage`**(GoFrame 双数据源 `g.DB("manage")`),与站点业务库彻底分离。
> 2. 开站走**半自动**:手动建库跑迁移,总后台登记连接信息 + 连通性校验 + 配置下发。
>
> 标注:🆕 需新建表(全部建在 my_manage 库,迁移目录独立 `manifest/sql/manage_migrations/`)。

---

## 架构速览

```
                    ┌────────────────────────────┐
                    │  总后台 manageapi (:8003)   │
                    │  库: my_manage(控制面)      │
                    │  商户/站点/域名/配置发布      │
                    └───────┬────────────┬───────┘
                     发布配置 │            │ 健康检查/汇总
                            ▼            ▼
                    ┌──────────┐   ┌──────────────────────┐
                    │  Nacos   │   │ 站点A部署(my库)        │
                    │ ns=env   │──▶│ 站点B部署(站B库)...    │
                    │ dataId=  │   │ front/backend 二进制   │
                    │ {code}.yaml   └──────────────────────┘
                    └──────────┘   站点侧零改造: boot/nacos.go
                                   已按 SITE_CODE 读配置
```

---

## 建议开发顺序

### M0 · 总后台基础设施 【最先,必做】
把「谁能登录总后台」立起来,并接好双数据源。

- 前置(手动):PG 建 `my_manage` 库;`config.yaml` + Nacos `my.yaml` 增加 `database.manage` 数据源
- 表 🆕:`manage_admin`、`manage_role`、`manage_casbin_rule`、`manage_log`(结构复用 B0 模式)
- 迁移:新目录 `manifest/sql/manage_migrations/`,goose 单独跑(连 my_manage)
- 接口:
  - `POST /manage/auth/login` / `POST /manage/auth/logout` / `GET /manage/auth/info`
- 中间件:`ManageAuth`(独立 token 前缀 `manage:token:`)+ `ManagePerm`(超管放行)+ `ManageOpLog`
- 模块:`internal/modules/manage`(domain/service/logic/controller/router),dao 全部走 `g.DB("manage")`
- 顺手修正:现 `/manage` 挂的是**前台用户 Auth** 的示范路由(user RegisterManage),撤掉,换独立管理员体系

### M1 · 商户管理 【高】
- 表 🆕:`merchant`(name / contact / phone / status / remark / created_at...)
- 接口:
  - `GET  /manage/merchants` 列表(名称模糊/状态/分页)
  - `POST /manage/merchants` 创建
  - `PUT  /manage/merchants/{id}` 更新
  - `POST /manage/merchants/{id}/disable` 停用/启用
- 保护:名下有站点的商户不可删/停用需确认(仅置状态,不物理删)

### M2 · 站点管理 + 域名分配 【高】
- 表 🆕:
  - `site`:merchant_id / **site_code(全局唯一)** / name / env(dev|test|prod) / status(0筹备 1上线 2停用) / db_host / db_port / db_name / db_user / db_pass(**AES 加密存**) / remark
  - `site_domain`:site_id / domain(唯一) / is_main / https / status
- 接口:
  - 站点 CRUD:`GET/POST /manage/sites`、`GET/PUT /manage/sites/{id}`、`POST /manage/sites/{id}/status`
  - 域名:`POST /manage/sites/{id}/domains` 绑定、`DELETE /manage/sites/{id}/domains/{domainId}` 解绑、`POST .../domains/{domainId}/main` 设主域名
  - `POST /manage/sites/{id}/db-check` **连通性校验**(用登记信息真连一次站点库,半自动开站的关键环节)
- 说明:`site_code` 就是 Nacos dataId 前缀 + 前台按 host 识别站点的标识(漫隐= my);域名分配即「把域名指到哪个商户站点」的登记源

### M3 · 配置管理与 Nacos 下发 【核心】
- 表 🆕:`site_config_publish`(site_id / env / content(yaml) / version / operator / published_at)
- 接口:
  - `GET  /manage/sites/{id}/config` 读取该站当前 Nacos 配置
  - `POST /manage/sites/{id}/config/publish` 发布(namespace=env 的 namespaceId,dataId=`{site_code}.yaml`,group 同站点侧)——发布前自动把 M2 登记的 DB 连接信息渲染进 `database` 段
  - `GET  /manage/sites/{id}/config/history` 发布历史
  - `POST /manage/sites/{id}/config/rollback` 回滚(取历史版本重新发布)
- 实现:总后台引 `nacos-sdk-go` ConfigClient.PublishConfig(依赖已在,boot 用过)
- 站点侧**零改造**:现有 `internal/boot/nacos.go` 按 SITE_CODE 拉配置且支持热更新,发布即生效
- 安全:发布记录里 DB 密码脱敏显示

### M4 · 开站流程串联(半自动) 【中】
- `POST /manage/sites/{id}/provision-check` 开站清单校验,逐项返回:
  - 商户状态有效 ✓/✗
  - 至少绑定一个域名且有主域名 ✓/✗
  - DB 可连通(复用 db-check)✓/✗
  - 站点库迁移版本 = 本仓库最新 goose 版本(连站点库查 `goose_db_version` 对比)✓/✗,不一致提示"需手动跑迁移"
  - Nacos 配置已发布 ✓/✗
- 全部通过才允许 `status → 1 上线`(状态机校验)
- 产出一份《开站 SOP》文档(手动步骤:建库→跑迁移→总后台登记→校验→发布配置→上线)

### M5 · 跨站监控 / 汇总 【中低】
- `GET /manage/sites/{id}/health` 站点健康:DB ping + (可选)HTTP 探活站点 `/health` 端点(站点侧需加一个公开 health 路由,顺手做)
- `GET /manage/overview` 跨站汇总:遍历已上线站点的登记库,汇总用户数/今日新增/累计充值(近期只有漫隐一个站,实现简单;站多了再改并发+缓存)
- `GET /manage/logs` 总后台操作审计查询(manage_log)

### M6 · 总后台权限细化 【可选,人多再做】
- manage_role + Casbin 细分(如:商务只管商户、运营只读)
- 注:`shared/rbac` 当前是单例且写死 `g.DB()` 默认库,需小改造支持传入库分组(manage 用 `manage_casbin_rule`);M0~M5 期间超管一人即可,不着急

---

## 新表汇总(全部在 my_manage 库)

| 阶段 | 新表 | 用途 |
|---|---|---|
| M0 | manage_admin / manage_role / manage_casbin_rule / manage_log | 总后台管理员/角色/权限/审计 |
| M1 | merchant | 商户 |
| M2 | site / site_domain | 站点 / 域名分配 |
| M3 | site_config_publish | 配置发布记录(含回滚依据) |

---

## 落地顺序与里程碑

1. **主线 M0 → M1 → M2 → M3**:做完即形成完整闭环——建商户 → 建站点绑域名登记库 → 校验 → 发配置到 Nacos → 站点部署读配置跑起来。
2. M4(开站清单)/M5(监控汇总)是增强,按需排。
3. M6 权限细化等总后台多人使用时再做。

## 依赖与前置清单(开工前手动做一次)

- [ ] PG 建库:`CREATE DATABASE my_manage;`
- [ ] `manifest/config/config.yaml` 与 Nacos `my.yaml` 增加 `database.manage` 数据源
- [ ] 确认 Nacos 三个环境 namespace 的 namespaceId(dev 已有,test/prod 发布功能要用)
- [ ] 选定 AES 密钥来源(env 变量 MANAGE_SECRET,用于站点 DB 密码加密)

## 与站点后台(B 系列)的边界

| 事项 | 站点后台 /backend | 总后台 /manage |
|---|---|---|
| 用户/订单/兑换码/公告 | ✅ 管本站 | ❌ 不直接管 |
| 商户/站点/域名 | ❌ | ✅ |
| 配置(客服链接等运营开关) | ✅ site_config 表 | ❌(基础设施配置走 Nacos 下发) |
| 数据库连接/Nacos 配置 | 只读(启动拉取) | ✅ 唯一写入方 |
| 管理员体系 | admin_user(站点库) | manage_admin(my_manage 库),互不相通 |
