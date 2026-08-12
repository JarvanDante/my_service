// Package v1 前台基础配置契约(移植自 tianbi ping/config, 改为 KV 全量下发)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// InfoReq 拉取全量配置(公开, 客户端启动第一个协议)。
type InfoReq struct {
	g.Meta `path:"/config/info" method:"get" tags:"Front/Config" summary:"全量基础配置"`
	Grp    string `json:"grp"` // 可选: 只取某分组
}
type InfoRes struct {
	Configs map[string]interface{} `json:"configs"` // key → 原始类型值
}

// CheckReq 健康检查(移植自 tianbi ping/check, 客户端测域名可用性)。
type CheckReq struct {
	g.Meta `path:"/config/check" method:"get" tags:"Front/Config" summary:"健康检查"`
}
type CheckRes struct{}
