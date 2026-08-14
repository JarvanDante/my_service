// Package appcfg 读取 app_config KV 配置的轻量工具(供资金/抽奖等模块取运营参数)。
// value 列是 jsonb: 数字存 `10`, 字符串存 `"abc"`, 布尔存 `true`。
// 取不到 / 解析失败一律返回调用方给的默认值, 不让配置缺失阻断业务。
package appcfg

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
)

const cfgSiteId = 1 // 单站点样板, 与各业务模块的 xxSiteId 保持一致

// raw 取原始 jsonb 文本; 未配置或已禁用返回空串。
func raw(ctx context.Context, key string) string {
	v, err := g.Model("app_config").Ctx(ctx).
		Where("site_id", cfgSiteId).Where("key", key).Where("status", 1).
		Fields("value").Value()
	if err != nil || v == nil {
		return ""
	}
	return v.String()
}

// Float 取浮点配置(如费率、单价)。
func Float(ctx context.Context, key string, def float64) float64 {
	s := raw(ctx, key)
	if s == "" {
		return def
	}
	var f float64
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return def
	}
	return f
}

// Int 取整数配置(如每日次数上限)。
func Int(ctx context.Context, key string, def int) int {
	s := raw(ctx, key)
	if s == "" {
		return def
	}
	var f float64
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return def
	}
	return int(f)
}

// String 取字符串配置。
func String(ctx context.Context, key, def string) string {
	s := raw(ctx, key)
	if s == "" {
		return def
	}
	var v string
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return def
	}
	return v
}

// StringSlice 取字符串数组配置(如默认头像列表)。
func StringSlice(ctx context.Context, key string, def []string) []string {
	s := raw(ctx, key)
	if s == "" {
		return def
	}
	var v []string
	if err := json.Unmarshal([]byte(s), &v); err != nil || len(v) == 0 {
		return def
	}
	return v
}

// Bool 取布尔配置(如提现开关)。
func Bool(ctx context.Context, key string, def bool) bool {
	s := raw(ctx, key)
	if s == "" {
		return def
	}
	var v bool
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return def
	}
	return v
}
