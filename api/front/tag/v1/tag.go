// Package v1 前台标签契约(移植自 tianbi tag repo/list)。
package v1

import "github.com/gogf/gf/v2/frame/g"

// RepoTagItem 前台仅暴露 id + name(与 tianbi RepoTagInfo 一致)。
type RepoTagItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// RepoListReq 按内容类型拉取标签列表(公开, 无需登录)。
type RepoListReq struct {
	g.Meta `path:"/tag/repo/list" method:"post" tags:"Front/Tag" summary:"按内容类型拉取标签"`
	Type   int `json:"type" v:"required|in:1,2,3,4,5,6,7#内容类型必填|内容类型非法"` // 1影片 2抖音 3动漫 4漫画 5图集 6帖子 7小说
	Page   int `json:"page"`
	Size   int `json:"size"`
}
type RepoListRes struct {
	List []RepoTagItem `json:"list"`
}
