// Package logic 基础配置业务(移植自 tianbi configser, 改为通用 KV)。
// value 存 jsonb: 后台录入时合法 JSON 原样存(数字/布尔/对象), 普通文本自动转 JSON 字符串;
// 前台下发时反序列化, 客户端拿到原始类型。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/config/service"
)

const cfgSiteId = 1 // 单站点样板

type sConfig struct{}

func New() service.IConfig { return &sConfig{} }

// normalizeValue 合法 JSON 原样返回, 否则转 JSON 字符串("hello" → "\"hello\"")。
func normalizeValue(v string) string {
	v = strings.TrimSpace(v)
	if json.Valid([]byte(v)) {
		return v
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// Info 前台全量配置(移植自 tianbi ConfigInfo, KV 化)。
func (s *sConfig) Info(ctx context.Context, grp string) (map[string]interface{}, error) {
	m := g.Model("app_config").Ctx(ctx).
		Where("site_id", cfgSiteId).Where("status", 1)
	if grp != "" {
		m = m.Where("grp", grp)
	}
	var list []*entity.AppConfig
	if err := m.Scan(&list); err != nil {
		return nil, err
	}
	out := make(map[string]interface{}, len(list))
	for _, r := range list {
		var v interface{}
		if err := json.Unmarshal([]byte(r.Value), &v); err != nil {
			v = r.Value // 兜底: 解析失败按原文本下发
		}
		out[r.Key] = v
	}
	return out, nil
}

func (s *sConfig) List(ctx context.Context, f service.ListFilter) ([]*service.ItemDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("app_config").Ctx(ctx).Where("site_id", cfgSiteId)
	if f.Grp != "" {
		m = m.Where("grp", f.Grp)
	}
	if f.Status >= 0 { // -1=全部
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		m = m.Where("(key ILIKE ? OR remark ILIKE ?)", kw, kw)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.AppConfig
	if err := m.Clone().OrderAsc("grp").OrderAsc("key").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(list))
	for _, r := range list {
		updated := ""
		if r.UpdatedAt != nil {
			updated = r.UpdatedAt.String()
		}
		out = append(out, &service.ItemDTO{
			Id: r.Id, Grp: r.Grp, Key: r.Key, Value: r.Value,
			Remark: r.Remark, Status: r.Status, UpdatedAt: updated,
		})
	}
	return out, total, nil
}

func (s *sConfig) Create(ctx context.Context, in service.CreateInput) (int64, error) {
	key := strings.TrimSpace(in.Key)
	if key == "" {
		return 0, gerror.New("key不能为空")
	}
	if in.Grp == "" {
		in.Grp = "base"
	}
	cnt, err := g.Model("app_config").Ctx(ctx).
		Where("site_id", cfgSiteId).Where("key", key).Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.New("该 key 已存在")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("app_config").Ctx(ctx).Data(g.Map{
		"site_id": cfgSiteId, "grp": in.Grp, "key": key,
		"value": normalizeValue(in.Value), "remark": in.Remark, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sConfig) Update(ctx context.Context, in service.UpdateInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"updated_at": gtime.Now()}
	if in.Grp != "" {
		data["grp"] = in.Grp
	}
	if in.Value != "" {
		data["value"] = normalizeValue(in.Value)
	}
	if in.Remark != "" {
		data["remark"] = in.Remark
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("app_config").Ctx(ctx).
		Where("site_id", cfgSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sConfig) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("app_config").Ctx(ctx).
		Where("site_id", cfgSiteId).Where("id", id).Delete()
	return err
}
