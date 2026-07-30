// Package rbac 基于 Casbin 的后台 RBAC。
// 使用 gdb 自实现适配器(不引 GORM), 复用 g.DB()/g.Model 直接读写 casbin_rule 表。
package rbac

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// RBAC 模型: sub=角色码, obj=请求路径(keyMatch2 支持 /a/:id 通配), act=方法; p.act="*" 放行所有方法。
const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

type casbinRow struct {
	Ptype string `orm:"ptype"`
	V0    string `orm:"v0"`
	V1    string `orm:"v1"`
	V2    string `orm:"v2"`
	V3    string `orm:"v3"`
	V4    string `orm:"v4"`
	V5    string `orm:"v5"`
}

type adapter struct{}

var cols = []string{"v0", "v1", "v2", "v3", "v4", "v5"}

func rowData(rule []string) g.Map {
	m := g.Map{}
	for i, v := range rule {
		if i >= len(cols) {
			break
		}
		m[cols[i]] = v
	}
	return m
}

// LoadPolicy 从 casbin_rule 全量读入内存。
func (a *adapter) LoadPolicy(m model.Model) error {
	var rows []casbinRow
	if err := g.Model("casbin_rule").Ctx(gctx.New()).Scan(&rows); err != nil {
		return err
	}
	for _, row := range rows {
		line := row.Ptype
		for _, v := range []string{row.V0, row.V1, row.V2, row.V3, row.V4, row.V5} {
			if v == "" {
				break
			}
			line += ", " + v
		}
		if err := persist.LoadPolicyLine(line, m); err != nil {
			return err
		}
	}
	return nil
}

// SavePolicy 全量覆盖写回(清空后重插)。
func (a *adapter) SavePolicy(m model.Model) error {
	ctx := gctx.New()
	if _, err := g.Model("casbin_rule").Ctx(ctx).Where("1=1").Delete(); err != nil {
		return err
	}
	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			d := rowData(rule)
			d["ptype"] = ptype
			if _, err := g.Model("casbin_rule").Ctx(ctx).Data(d).Insert(); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddPolicy 增量落库一条策略。
func (a *adapter) AddPolicy(sec, ptype string, rule []string) error {
	d := rowData(rule)
	d["ptype"] = ptype
	_, err := g.Model("casbin_rule").Ctx(gctx.New()).Data(d).Insert()
	return err
}

// RemovePolicy 删除一条策略。
func (a *adapter) RemovePolicy(sec, ptype string, rule []string) error {
	m := g.Model("casbin_rule").Ctx(gctx.New()).Where("ptype", ptype)
	for k, v := range rowData(rule) {
		m = m.Where(k, v)
	}
	_, err := m.Delete()
	return err
}

// RemoveFilteredPolicy 按字段条件批量删除。
func (a *adapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	m := g.Model("casbin_rule").Ctx(gctx.New()).Where("ptype", ptype)
	for i, v := range fieldValues {
		if v == "" {
			continue
		}
		idx := fieldIndex + i
		if idx < len(cols) {
			m = m.Where(cols[idx], v)
		}
	}
	_, err := m.Delete()
	return err
}

var (
	once  sync.Once
	en    *casbin.Enforcer
	enErr error
)

func get() (*casbin.Enforcer, error) {
	once.Do(func() {
		m, err := model.NewModelFromString(modelText)
		if err != nil {
			enErr = err
			return
		}
		en, enErr = casbin.NewEnforcer(m, &adapter{})
	})
	return en, enErr
}

// Enforce 判断角色 role 是否允许以 method 访问 path。
func Enforce(role, path, method string) (bool, error) {
	e, err := get()
	if err != nil {
		return false, err
	}
	return e.Enforce(role, path, method)
}

// AddPolicy 给角色加一条权限并落库。
func AddPolicy(role, path, method string) (bool, error) {
	e, err := get()
	if err != nil {
		return false, err
	}
	return e.AddPolicy(role, path, method)
}

// RemovePolicy 删除角色一条权限并落库。
func RemovePolicy(role, path, method string) (bool, error) {
	e, err := get()
	if err != nil {
		return false, err
	}
	return e.RemovePolicy(role, path, method)
}

// PermsForRole 返回角色的全部权限, 每项为 [path, method]。
func PermsForRole(role string) [][]string {
	e, err := get()
	if err != nil {
		return nil
	}
	rows, _ := e.GetFilteredPolicy(0, role)
	out := make([][]string, 0, len(rows))
	for _, r := range rows {
		if len(r) >= 3 {
			out = append(out, []string{r[1], r[2]})
		}
	}
	return out
}

// Reload 改动策略后重载(当前 Add/Remove 已即时生效, 预留手动重载)。
func Reload() error {
	e, err := get()
	if err != nil {
		return err
	}
	return e.LoadPolicy()
}
