# 站点字段差异配置（Nacos）

多产品 SaaS 的**响应字段差异**走 Nacos（本地开发走 `manifest/config/config.yaml`）。  
功能/接口开关本期不做。

## 配置位置

| 环境 | 来源 |
|------|------|
| 部署 | Nacos `dataId=<SITE_CODE>.yaml`（`Watch=true` 热更新） |
| 本地 | `manifest/config/config.yaml` |
| 模板 | `docs/nacos-config-my.yaml` |

## 字段差异：`response.*_extra` + `ext`

```yaml
response:
  comics_list_extra:
    - topic_follow
    - update_date_label
  comics_detail_extra:
    - topic_follow
```

代码：

```go
import "github.com/JarvanDante/my_service/internal/shared/siteconf"

type Item struct {
    ID    int64                  `json:"id"`
    Title string                 `json:"title"`
    Ext   map[string]interface{} `json:"ext,omitempty"`
}

item.Ext = siteconf.PickExt(ctx, "comics_list", map[string]interface{}{
    "topic_follow":      row.TopicFollow,
    "update_date_label": row.UpdateDateLabel,
    "internal_only":     row.Secret, // 未在白名单则不会进 ext
})
```

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
