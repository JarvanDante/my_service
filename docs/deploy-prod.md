# 生产部署说明(一站一部署 + Nacos)

## 一、引导配置 = 环境变量(不进 Nacos、不进 git)

应用启动只靠这几个环境变量找到 Nacos,其余全部从 Nacos 拉。

| 变量 | 必填 | 说明 | 示例 |
|---|---|---|---|
| NACOS_ADDR | 是 | Nacos 地址(不设则用本地 config.yaml) | `nacos:8848` |
| NACOS_NAMESPACE | 是 | 命名空间ID(=环境 dev/test/prod) | `prod 的命名空间ID` |
| NACOS_GROUP | 否 | 分组,默认 DEFAULT_GROUP | `DEFAULT_GROUP` |
| SITE_CODE | 是 | 站点编码 → dataId=`<code>.yaml` | `jh` |
| NACOS_USER | 开鉴权必填 | Nacos 账号 | `nacos` |
| NACOS_PASS | 开鉴权必填 | Nacos 密码(走 Secret/.env,勿写死) | — |

**要点**:一个镜像 + 不同环境变量 = 不同环境 / 不同商户。切商户只换 `SITE_CODE`,切环境只换 `NACOS_NAMESPACE`(或 `NACOS_ADDR`)。

## 二、Nacos 生产要求(dev 我们简化了,生产要补)

1. **开鉴权**:`NACOS_AUTH_ENABLE=true`,并在控制台创建管理员账号;应用侧设 `NACOS_USER`/`NACOS_PASS`。
2. **MySQL 持久化 + 集群**:别用 dev 的 standalone 内嵌存储;生产用 MySQL 存储 + 至少 3 节点,避免配置丢失/单点。
3. **每个环境一套命名空间**:dev/test/prod 三个命名空间隔离;每个站点一条配置(dataId=站点编码.yaml),`database.link` 指各自库。

## 三、应用 Dockerfile(多阶段构建,放项目根)

```dockerfile
# ---- build ----
FROM golang:1.23-alpine AS build
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 构建前台二进制(其余 backendapi/manageapi/cron 各构建一个镜像或同镜像多入口)
RUN go build -o /app/frontapi ./app/frontapi

# ---- run ----
FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=build /app/frontapi .
# 本地兜底配置(Nacos 不可用时的回退,可选)
COPY manifest/config/config.yaml ./manifest/config/config.yaml
EXPOSE 8001
ENTRYPOINT ["/app/frontapi"]
```

> 四个二进制(front/backend/manage/cron)可以各打一个镜像,或同一镜像用不同 ENTRYPOINT/启动参数。

## 四、一站一部署 docker-compose 示例

```yaml
# 每个站点一份 service, 换 SITE_CODE + 端口即可
services:
  site-jh-front:
    image: my_service:latest
    container_name: site-jh-front
    restart: always
    environment:
      NACOS_ADDR: "nacos:8848"
      NACOS_NAMESPACE: "prod命名空间ID"
      NACOS_GROUP: "DEFAULT_GROUP"
      NACOS_USER: "nacos"
      NACOS_PASS: "${NACOS_PASS}"     # 从 .env / secret 注入
      SITE_CODE: "jh"                 # -> dataId=jh.yaml
    ports:
      - "18001:8001"
    networks: [prod-net]

  site-abc-front:
    image: my_service:latest
    container_name: site-abc-front
    restart: always
    environment:
      NACOS_ADDR: "nacos:8848"
      NACOS_NAMESPACE: "prod命名空间ID"
      NACOS_USER: "nacos"
      NACOS_PASS: "${NACOS_PASS}"
      SITE_CODE: "abc"                # -> dataId=abc.yaml
    ports:
      - "28001:8001"
    networks: [prod-net]

networks:
  prod-net:
```

`.env`(gitignore)放敏感值:
```
NACOS_PASS=你的nacos密码
```

## 五、上线步骤

1. 在生产 Nacos 的 `prod` 命名空间,为每个站点建配置 `<site>.yaml`(照 `docs/nacos-config-my.yaml` 模板,`database.link` 指该站点的生产库)。
2. 构建镜像:`docker build -t my_service:latest .`
3. 每个站点起一份容器,注入自己的 `SITE_CODE` + 端口。
4. 前面挂 Nginx,按域名反代到对应站点容器(域名→站点的映射在总后台维护)。
5. 改配置在 Nacos 改+发布,`Watch:true` 热更新,无需重启(改库连接等需重启的除外)。

## 六、本地开发(对照)

不设 `NACOS_ADDR` → 用本地 `manifest/config/config.yaml`,`site_id` 默认走库默认值,零依赖 Nacos。
