# my_service · 漫隐 API

GoFrame 模块化单体(Module-First)。每个 `internal/modules/<模块>/` 是一个自包含业务边界,
将来可整目录拆为微服务。

## 目录约定

```
api/<模块>/v1/        接口契约(req/res + 校验)
internal/
  cmd/                启动与依赖装配(唯一注入 dao、注册路由处)
  modules/<模块>/
    controller/       接入适配, ghttp 只在这
    service/          对外接口, 别的模块只认此契约
    logic/            业务实现
    domain/           领域实体 + 仓储接口, 纯 Go 零框架
    router.go         模块自注册
  dao/                统一仓储实现, g.Model/gdb 只在这
  model/entity|do/    数据模型(gf gen dao 生成)
  shared/             全局中间件/错误码/常量/工具
  event/  mq/         模块间事件 / Kafka
manifest/config/      运行配置(PG/Redis/Kafka)
hack/config.yaml      gf gen 多模块生成配置
```

## 框架隔离铁律

- `ghttp` 只出现在 controller,`g.Model`/`gdb` 只出现在 dao,`domain/` 零框架 import。
- 模块间只经 service 接口或 event/mq 通信,禁止跨模块 import logic/dao、禁止读对方表。
- 跨模块依赖在 cmd 集中注入。

## 快速开始

```bash
go mod tidy
gf run main.go          # 或 go run main.go
# 生成代码:
gf gen dao              # 由 PG 反向生成 dao/model
gf gen service          # 由各模块 logic 生成 service 接口
gf gen ctrl             # 由各模块 api 生成 controller 骨架
```

## 新增模块

1. 复制 `internal/modules/user` 改名
2. `hack/config.yaml` 的 service/ctrl 段各加一条映射
3. `internal/cmd/cmd.go` 加一行 `xxxmod.Register(group, dao.NewXxxRepo())`
