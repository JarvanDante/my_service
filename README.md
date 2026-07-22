# my_service · 漫隐 API

GoFrame 模块化单体(Module-First)+ 多二进制入口。业务代码一份共享,
前台/后台/总后台/定时任务各编译为独立二进制、独立进程,互不影响资源。

## 二进制与端口

| 二进制 | 入口 | 路由前缀 | 默认端口 | 说明 |
|---|---|---|---|---|
| frontapi   | app/frontapi   | /front   | 8001 | 面向 C 端, 限流 |
| backendapi | app/backendapi | /backend | 8002 | 运营/管理, 鉴权 |
| manageapi   | app/manageapi   | /manage   | 8003 | 平台超管(总后台), 鉴权 |
| cron       | app/cron       | 无        | 无    | 定时任务, 与 API 隔离 |
| (dev)      | main.go        | 全部      | 8000 | 本地一体化, 单进程调试 |

## 构建 / 运行

```bash
make all         # 构建四个二进制到 bin/
make front       # 单独构建前台
make dev         # 本地一体化运行(gf run main.go)
go run ./app/frontapi   # 或单独跑某个入口
```

## 目录约定

```
app/<bin>/main.go       各二进制入口(薄, 只调 internal/cmd)
internal/cmd/           各入口启动装配: frontapi/backendapi/manageapi/cron + cmd.go(dev)
api/{front,backend,manage}/<模块>/v1/   ★接口契约按门面拆★
internal/modules/<模块>/
  controller/{front,backend,manage}/    ★控制器按门面拆★
  service|logic|domain/                ★业务共享, 一份★
  router.go             RegisterFront/Backend/Admin
internal/dao|model/     统一数据访问与模型(共享)
internal/shared/        中间件(CORS/Auth/RateLimit)/错误码/工具
internal/event|mq/      模块间事件 / Kafka
manifest/config/        运行配置(各门面端口 / PG / Redis / Kafka)
hack/config.yaml        gf gen 多模块生成配置
```

## 隔离要点

- 四个进程独立: 后台批量/导出等重操作不占用前台 API 的 CPU/协程, 可分机部署、独立扩缩容。
- 后台/总后台可在 config 指向**只读副本库**, 避免拖慢前台主库。
- 各门面独立中间件: 前台限流, 后台鉴权。

## 框架隔离铁律

- `ghttp` 只在 controller,`g.Model`/`gdb` 只在 dao,`domain/` 零框架 import。
- 模块间只经 service 接口或 event/mq, 禁止跨模块 import logic/dao、禁止读对方表。
- 跨模块依赖在 internal/cmd 各入口集中注入。

## 新增门面接口(以已有 user 模块为例)

1. `api/<门面>/user/v1/` 加 req/res
2. `internal/modules/user/controller/<门面>/` 加控制器方法
3. `router.go` 对应 RegisterXxx 里 Bind
4. 对应 `internal/cmd/<门面>api.go` 已自动挂载, 无需改动

## 快速开始

```bash
go mod tidy && make all
```
