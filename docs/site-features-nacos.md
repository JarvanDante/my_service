# 站点字段差异配置（Nacos）

多产品 SaaS 的**响应字段差异**走 Nacos（本地开发走 `manifest/config/config.yaml`）。  
功能/接口开关本期不做。

## 配置位置

| 环境 | 来源 |
|------|------|
| 部署 | Nacos `dataId=<SITE_CODE>.yaml`（`Watch=true` 热更新） |
| 本地 | `manifest/config/config.yaml` |
| 模板 | `docs/nacos-config-my.yaml` |

## 样板：`GET /front/user/info`

```yaml
response:
  user_info_extra:
    - bg_img
    - share_num
    - channel_name
    # - group_end_time
```

返回示例（公共字段不变，差异在 `ext`）：

```json
{
  "code": 0,
  "data": {
    "id": 2,
    "username": "device_64",
    "nickname": "用户864648",
    "balance": 0,
    "ext": {
      "bg_img": "",
      "share_num": 0,
      "channel_name": ""
    }
  }
}
```

实现位置：`controller/front/user.go` → `toApiUser` → `siteconf.PickExt(ctx, "user_info", ...)`。

改 Nacos / 本地 `response.user_info_extra` 即可增减 `ext` 字段，无需改接口路径。

| API | 作用 |
|-----|------|
| `ExtraFields(ctx, scene)` | 读白名单 |
| `HasExtra(ctx, scene, field)` | 是否声明某字段 |
| `PickExt(ctx, scene, candidates)` | 按白名单挑进 `ext` |

约定：公共字段固定；差异只进 `ext`；客户端忽略未知 key。

## 与 site_config 表的分工

| 配置 | 存放 |
|------|------|
| 扩展字段名单、基础设施 | **Nacos YAML** |
| 运营可改的文案/链接等 KV | `site_config` 表（现有 `siteconf.Get/Set`） |
