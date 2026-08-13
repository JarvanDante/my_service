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

## 写用例时的坑(踩过的)

- **`client.global` 只能存字符串/数字, 存不了函数或对象。**
  IntelliJ HTTP Client 的全局变量是要跨请求持久化的, 存函数进去, 后面的脚本
  `client.global.get(...)` 拿到的不是函数, 一调用整个脚本块就抛异常; 该块里后续的
  `client.global.set` 全部不执行, 于是后面引用这些变量的请求会因为 `{{变量}}` 解析不出来而
  **"Failed to start"**, 一路雪崩。需要跨块复用的工具函数(比如自己实现的 md5), 只能在每个
  用到它的脚本块里各写一份。
- 断言里跟 `client.global.get()` 做 `===` 比较时, 注意两边类型要一致。
- 用例要能**重复执行**: 固定设备号的账号会累积数据, 涉及"每日上限""每人限领""已绑定"这类
  一次性状态的步骤, 要么先把上限调高跑完再恢复, 要么每轮用随机后缀避开。
