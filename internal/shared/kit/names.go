package kit

import (
	"strconv"
	"strings"
)

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

// ParseI64CSV 拆中英文逗号的正整数 ID。
func ParseI64CSV(raw string) []int64 {
	raw = strings.ReplaceAll(strings.TrimSpace(raw), "，", ",")
	if raw == "" {
		return nil
	}
	seen := map[int64]struct{}{}
	out := make([]int64, 0, 8)
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// ReplaceNameCSV 把逗号分隔名称里的 old 整词换成 new，未命中则原样返回。
func ReplaceNameCSV(raw, old, new string) string {
	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)
	if old == "" || new == "" || old == new {
		return raw
	}
	next, changed := replaceNameList(NamesCSV(raw), old, new)
	if !changed {
		return raw
	}
	return strings.Join(next, ",")
}

// ReplaceNameList 把名称切片里的 old 整词换成 new；changed 表示是否有替换。
func ReplaceNameList(list []string, old, new string) (next []string, changed bool) {
	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)
	if old == "" || new == "" || old == new {
		return list, false
	}
	return replaceNameList(list, old, new)
}

func replaceNameList(list []string, old, new string) ([]string, bool) {
	changed := false
	tmp := make([]string, 0, len(list))
	for _, p := range list {
		if p == old {
			p = new
			changed = true
		}
		tmp = append(tmp, p)
	}
	if !changed {
		return list, false
	}
	return NamesCSV(strings.Join(tmp, ",")), true
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
