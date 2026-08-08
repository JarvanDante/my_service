// Package siteconf 中与 Nacos/本地 YAML 相关的响应字段扩展。
//
// 配置约定(写在 <SITE_CODE>.yaml / manifest/config/config.yaml):
//
//	response:
//	  comics_list_extra: [topic_follow, update_date_label]
//
// Watch=true 时 Nacos 变更会热更新到 g.Cfg(), 本包每次直接读配置, 无需本地缓存。
package siteconf

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// ExtraFields 读取某接口的扩展字段名单(Nacos response.<scene>_extra)。
// 例: ExtraFields(ctx, "comics_list") → response.comics_list_extra
func ExtraFields(ctx context.Context, scene string) []string {
	scene = strings.TrimSpace(scene)
	if scene == "" {
		return nil
	}
	key := "response." + scene + "_extra"
	raw := g.Cfg().MustGet(ctx, key).Strings()
	out := make([]string, 0, len(raw))
	seen := map[string]struct{}{}
	for _, f := range raw {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, ok := seen[f]; ok {
			continue
		}
		seen[f] = struct{}{}
		out = append(out, f)
	}
	return out
}

// HasExtra 判断某 scene 是否声明了指定扩展字段。
func HasExtra(ctx context.Context, scene, field string) bool {
	field = strings.TrimSpace(field)
	for _, f := range ExtraFields(ctx, scene) {
		if f == field {
			return true
		}
	}
	return false
}

// PickExt 按 Nacos 白名单从候选 map 中挑出扩展字段; 无命中返回 nil(便于 omitempty)。
func PickExt(ctx context.Context, scene string, candidates map[string]interface{}) map[string]interface{} {
	fields := ExtraFields(ctx, scene)
	if len(fields) == 0 || len(candidates) == 0 {
		return nil
	}
	ext := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		if v, ok := candidates[f]; ok {
			ext[f] = v
		}
	}
	if len(ext) == 0 {
		return nil
	}
	return ext
}
