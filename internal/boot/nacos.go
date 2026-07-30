// Package boot 启动引导。按需从 Nacos 加载配置。
// 未设置环境变量 NACOS_ADDR 时, 走本地 manifest/config/config.yaml, 不影响本地开发。
//
// 环境变量:
//
//	NACOS_ADDR       Nacos 地址, 如 127.0.0.1:8848 (不设则用本地文件)
//	NACOS_NAMESPACE  命名空间 ID (对应 dev/test/prod 命名空间的 ID; public 为空)
//	NACOS_GROUP      分组, 默认 DEFAULT_GROUP
//	SITE_CODE        站点编码, 作为 dataId 前缀, 默认 my
//	NACOS_DATAID     直接指定 dataId, 默认 <SITE_CODE>.yaml
package boot

import (
	"strconv"
	"strings"

	nacos "github.com/gogf/gf/contrib/config/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func init() {
	addr := genv.Get("NACOS_ADDR").String()
	if addr == "" {
		return // 本地开发: 用 manifest/config/config.yaml
	}
	ctx := gctx.GetInitCtx()
	host, port := parseAddr(addr)

	namespace := genv.Get("NACOS_NAMESPACE").String()
	group := genv.Get("NACOS_GROUP", "DEFAULT_GROUP").String()
	site := genv.Get("SITE_CODE", "my").String()
	dataId := genv.Get("NACOS_DATAID", site+".yaml").String()

	adapter, err := nacos.New(ctx, nacos.Config{
		ServerConfigs: []constant.ServerConfig{{IpAddr: host, Port: port}},
		ClientConfig: constant.ClientConfig{
			NamespaceId:         namespace,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./temp/nacos/log",
			CacheDir:            "./temp/nacos/cache",
			LogLevel:            "warn",
		},
		ConfigParam: vo.ConfigParam{DataId: dataId, Group: group},
		Watch:       true, // 配置变更热更新
	})
	if err != nil {
		g.Log().Fatalf(ctx, "nacos 配置加载失败: %+v", err)
	}
	g.Cfg().SetAdapter(adapter)
	g.Log().Infof(ctx, "配置来源=Nacos addr=%s ns=%s group=%s dataId=%s", addr, namespace, group, dataId)
}

func parseAddr(a string) (string, uint64) {
	a = strings.TrimSpace(a)
	host, portStr := a, "8848"
	if i := strings.LastIndex(a, ":"); i >= 0 {
		host, portStr = a[:i], a[i+1:]
	}
	p, _ := strconv.ParseUint(portStr, 10, 64)
	if p == 0 {
		p = 8848
	}
	return host, p
}
