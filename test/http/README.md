# 真实接口(E2E)测试 · .http

用 JetBrains(GoLand/IDEA)或 VSCode 的 HTTP Client 直接对**运行中的服务**发请求。

## 运行

1. 启动服务(任选其一):

   ```bash
   gf run main.go            # 一体化, :8000(user.http 默认地址)
   # 或单独跑前台
   go run ./app/frontapi     # :8001, 需把 user.http 顶部 @baseUrl 改成 :8001
   ```

2. 打开 `user.http`:
   - 单个请求:点请求左侧 ▶。
   - 整条流程:右上角「Run all requests in file」,从登录到退出顺序执行,查看每步断言结果。

## 说明

- 顶部 `@baseUrl` 控制目标地址,换环境改这一行即可。
- 登录/刷新的响应脚本会把 token 存到全局变量 `{{token}}`,后续请求自动接力,无需手动复制。
- 每个请求带 `client.test(...)` 断言,结果在 IDE 的 Services/Run 面板可见。
- 这是**黑盒测试**:不启动 server、不 import 业务代码,只当外部客户端打真实服务,连的是 `config.yaml` 里配置的真实库,请对着开发/测试环境跑,勿对生产。

## 多环境(可选)

可在本目录放 `http-client.env.json` 定义多套环境:

```json
{
  "dev":  { "baseUrl": "http://127.0.0.1:8000" },
  "test": { "baseUrl": "http://192.168.x.x:8000" }
}
```

然后把 `user.http` 里的 `{{baseUrl}}` 改用环境变量,运行时右上角切换环境。
