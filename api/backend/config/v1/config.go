// Package v1 后台基础配置契约(KV 管理)。
package v1

import "github.com/gogf/gf/v2/frame/g"

type Item struct {
	Id        int64  `json:"id"`
	Grp       string `json:"grp"`
	Key       string `json:"key"`
	Value     string `json:"value"` // 原始 JSON 文本
	Remark    string `json:"remark"`
	Status    int    `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type ListReq struct {
	g.Meta  `path:"/config" method:"get" tags:"Backend/Config" summary:"配置列表"`
	Grp     string `json:"grp"`     // 空=全部
	Status  string `json:"status"`  // 空=全部  0=禁用  1=启用
	Keyword string `json:"keyword"` // key/备注模糊搜索
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}
type ListRes struct {
	List  []Item `json:"list"`
	Total int    `json:"total"`
}

type CreateReq struct {
	g.Meta `path:"/config" method:"post" tags:"Backend/Config" summary:"新增配置(value 支持 JSON: 数字/布尔/字符串/对象)"`
	Grp    string `json:"grp"`
	Key    string `json:"key" v:"required#key必填"`
	Value  string `json:"value" v:"required#value必填"` // 合法 JSON 原样存; 普通文本自动转 JSON 字符串
	Remark string `json:"remark"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type CreateRes struct {
	Id int64 `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/config/{id}" method:"put" tags:"Backend/Config" summary:"更新配置"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID必填"`
	Grp    string `json:"grp"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
	Status int    `json:"status" v:"in:0,1#状态非法"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/config/{id}" method:"delete" tags:"Backend/Config" summary:"删除配置"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID必填"`
}
type DeleteRes struct{}
