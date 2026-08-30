package kit

import "strings"

// NamesCSV 拆中英文逗号并去掉空白。
func NamesCSV(raw string) []string {
	raw = strings.ReplaceAll(strings.TrimSpace(raw), "，", ",")
	if raw == "" {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, 4)
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

// MergeNames 合并多个名称切片并去重。
func MergeNames(parts ...[]string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, list := range parts {
		for _, p := range list {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}

// CategoryOverlapWhere 作品 category（逗号分隔名称）命中任一所选分类，会 trim 空格。
func CategoryOverlapWhere() string {
	return `EXISTS (
		SELECT 1 FROM unnest(string_to_array(replace(coalesce(category, ''), '，', ','), ',')) AS t
		WHERE btrim(t) <> '' AND btrim(t) = ANY(string_to_array(?, ','))
	)`
}
