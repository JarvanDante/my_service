// Package worklabel 分类/标签改名时回写作品上拷贝的名称字符串。
package worklabel

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/shared/kit"
)

const siteID = 1

// RewriteCategoryCSV 把作品 category（逗号分隔名称）里的旧名整词换成新名。
func RewriteCategoryCSV(ctx context.Context, table, oldName, newName string, extra map[string]any) error {
	oldName = strings.TrimSpace(oldName)
	newName = strings.TrimSpace(newName)
	if table == "" || oldName == "" || newName == "" || oldName == newName {
		return nil
	}
	m := g.Model(table).Ctx(ctx).Where("site_id", siteID)
	for k, v := range extra {
		m = m.Where(k, v)
	}
	var rows []struct {
		Id       int64  `orm:"id"`
		Category string `orm:"category"`
	}
	if err := m.Fields("id", "category").Where(kit.CategoryOverlapWhere(), oldName).Scan(&rows); err != nil {
		return err
	}
	now := gtime.Now()
	for _, r := range rows {
		next := kit.ReplaceNameCSV(r.Category, oldName, newName)
		if next == r.Category {
			continue
		}
		if _, err := g.Model(table).Ctx(ctx).Where("id", r.Id).Data(g.Map{
			"category": next, "updated_at": now,
		}).Update(); err != nil {
			return err
		}
	}
	return nil
}

// RewriteJSONNames 把 jsonb 名称数组（tags / topics）里的旧名整词换成新名。
func RewriteJSONNames(ctx context.Context, table, field, oldName, newName string, extra map[string]any) error {
	oldName = strings.TrimSpace(oldName)
	newName = strings.TrimSpace(newName)
	if table == "" || field == "" || oldName == "" || newName == "" || oldName == newName {
		return nil
	}
	needle, err := json.Marshal([]string{oldName})
	if err != nil {
		return err
	}
	m := g.Model(table).Ctx(ctx).Where("site_id", siteID).
		Where(field+" @> ?::jsonb", string(needle))
	for k, v := range extra {
		m = m.Where(k, v)
	}
	all, err := m.Fields("id", field).All()
	if err != nil {
		return err
	}
	now := gtime.Now()
	for _, rec := range all {
		raw := rec[field].String()
		list := []string{}
		if raw != "" {
			_ = json.Unmarshal([]byte(raw), &list)
		}
		next, changed := kit.ReplaceNameList(list, oldName, newName)
		if !changed {
			continue
		}
		b, err := json.Marshal(next)
		if err != nil {
			return err
		}
		if _, err := g.Model(table).Ctx(ctx).Where("id", rec["id"].Int64()).Data(g.Map{
			field: string(b), "updated_at": now,
		}).Update(); err != nil {
			return err
		}
	}
	return nil
}
