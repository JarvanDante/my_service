// Package siteconf 站点键值配置(site_config 表)读写。
// 定位为共享基础设施: 各模块均可读, 写入走后台配置接口。
package siteconf

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Get 读配置; 无记录时返回 def。
func Get(ctx context.Context, key, def string) string {
	v, err := g.Model("site_config").Ctx(ctx).Where("conf_key", key).Value("conf_value")
	if err != nil || v == nil || v.String() == "" {
		return def
	}
	return v.String()
}

// Set 写配置(upsert)。
func Set(ctx context.Context, key, value string) error {
	n, err := g.Model("site_config").Ctx(ctx).Where("conf_key", key).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		_, err = g.Model("site_config").Ctx(ctx).Where("conf_key", key).Data(g.Map{
			"conf_value": value, "updated_at": gtime.Now(),
		}).Update()
		return err
	}
	_, err = g.Model("site_config").Ctx(ctx).Data(g.Map{
		"conf_key": key, "conf_value": value,
	}).Insert()
	return err
}
