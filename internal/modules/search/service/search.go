// Package service 全站搜索对外接口。
package service

import "context"

// ItemDTO 跨表统一结果项(字段含义见 api/front/search/v1.Item)。
type ItemDTO struct {
	Id        int64
	MediaType int
	Title     string
	Cover     string
	Author    string
	Price     float64
	IsVip     int
	ViewCount int64
	CreatedAt string
}

// ResultDTO 两种形态共用: type=0 填分组 + TotalHit, type>0 填 List + Total。
type ResultDTO struct {
	Videos   []*ItemDTO
	Posts    []*ItemDTO
	Comics   []*ItemDTO
	Novels   []*ItemDTO
	Photos   []*ItemDTO
	TotalHit int

	List  []*ItemDTO
	Total int
}

type SearchInput struct {
	Keyword string
	Type    int // 0全部 1视频 2帖子 3漫画 4小说 5图集
	Page    int
	Size    int
}

type ISearch interface {
	// Search 跨表标题模糊搜(只出已上架内容), 顺带把关键词计入热搜。
	Search(ctx context.Context, in SearchInput) (*ResultDTO, error)
	// Suggest 按前缀从热搜词联想, 最多 10 个。
	Suggest(ctx context.Context, keyword string) ([]string, error)
}
