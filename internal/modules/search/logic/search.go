// Package logic 全站搜索业务(替代 tianbi 的多 ES 索引方案)。
//
// 关键取舍:
//   - 不建搜索表: 直接查各内容表的 title, 走 00035 建的 pg_trgm GIN 索引,
//     ILIKE '%词%' 能吃到索引, 中文短词场景够用, 也省掉"索引与库不同步"的一整类问题;
//   - 不经过各内容模块的 service/repo: 搜索是跨 5 张表的只读聚合, 走各模块接口既拿不到
//     统一字段, 又会退化成 5 次业务查询(带各自的付费判定/计数副作用)。这里只读标题与展示字段,
//     用同一套 SQL 形态直查, 是本模块自己的读模型;
//   - 热搜词复用 hot_search 表(00030), 不另起统计表。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/search/service"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const (
	schSiteId       = 1
	schAggEach      = 5  // type=0 时每类取几条(聚合页只做"有没有", 深翻由 type>0 分页承担)
	schSuggestLimit = 10 // 联想词上限
	schKeywordMax   = 64 // hot_search.keyword 是 varchar(64), 超长的词不入库免得写失败
)

type sSearch struct{}

func New() service.ISearch { return &sSearch{} }

// searchRow 各表查询结果的统一落点。
// 表里没有的列不出现在 Fields 里, 对应字段就保持零值 —— 比在 SQL 里拼空串常量占位
// (SELECT 空串 AS author 那种写法)干净, 也不用担心 PG 推断不出字面量的类型。
type searchRow struct {
	Id        int64       `orm:"id"`
	Title     string      `orm:"title"`
	Cover     string      `orm:"cover"`
	Author    string      `orm:"author"`
	Price     float64     `orm:"price"`
	IsVip     int         `orm:"is_vip"`
	ViewCount int64       `orm:"view_count"`
	CreatedAt *gtime.Time `orm:"created_at"`
}

// source 一个可搜索的内容源。fields 决定该表能填上统一结构的哪些字段。
type source struct {
	table     string
	mediaType int
	fields    string
}

// 上架状态: 五张表恰好都用 1 表示"可见"(video=已发布, post=审核通过, 其余=上架),
// 所以这里统一用一个常量, 不必按表分支。
const schOnline = entity.ContentStatusOnline

// 顺序即 type=0 聚合结果的分组顺序。
var schSources = []source{
	// video 表没有 author/price/is_vip/view_count 这几列(付费与计数尚未接入视频),
	// 缺的字段留零值, 前端按 media_type 自行决定是否展示。
	{table: "video", mediaType: paywall.MediaVideo, fields: "id, title, cover_url AS cover, created_at"},
	// post 是 UGC 帖子: 封面在 pics 数组里, 单条记录没有独立封面列, 故 cover 留空;
	// author 是 user_id 而非昵称, 需要联表才有意义, 搜索结果页用不上, 一并留空。
	{table: "post", mediaType: paywall.MediaPost, fields: "id, title, view_count, created_at"},
	{table: "comics", mediaType: paywall.MediaComics, fields: "id, title, cover, author, price, is_vip, view_count, created_at"},
	{table: "novel", mediaType: paywall.MediaNovel, fields: "id, title, cover, author, price, is_vip, view_count, created_at"},
	// photo_album 没有 author 列(图集按分类/标签组织, 不记作者)。
	{table: "photo_album", mediaType: paywall.MediaPhoto, fields: "id, title, cover, price, is_vip, view_count, created_at"},
}

// likeEscaper 转义 LIKE 通配符。用户搜 "100%" 时不该被当成通配符扩大命中面,
// PG 的 LIKE/ILIKE 默认转义字符就是反斜杠。
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func toDTO(r *searchRow, mediaType int) *service.ItemDTO {
	d := &service.ItemDTO{
		Id: r.Id, MediaType: mediaType, Title: r.Title, Cover: r.Cover,
		Author: r.Author, Price: r.Price, IsVip: r.IsVip, ViewCount: r.ViewCount,
	}
	if r.CreatedAt != nil {
		d.CreatedAt = r.CreatedAt.String()
	}
	return d
}

// queryOne 查单个内容源。needTotal=false 时跳过 count(聚合页每类只要一小撮, 但仍要
// 计入 total_hit, 所以目前两条路径都要 count; 参数留着是为了以后聚合页想省一半查询时好关)。
func (s *sSearch) queryOne(ctx context.Context, src source, kw string, page, size int, needTotal bool) ([]*service.ItemDTO, int, error) {
	base := g.Model(src.table).Ctx(ctx).
		Where("site_id", schSiteId).
		Where("status", schOnline).
		// 走 00035 / 00036~00038 建的 GIN trgm 索引; 扩展没装时退化为顺序扫描, 结果不变。
		Where("title ILIKE ?", "%"+likeEscaper.Replace(kw)+"%")
	total := 0
	if needTotal {
		n, err := base.Clone().Count()
		if err != nil {
			return nil, 0, err
		}
		total = n
	}
	var rows []*searchRow
	if err := base.Clone().Fields(src.fields).OrderDesc("id").Page(page, size).Scan(&rows); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ItemDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, toDTO(r, src.mediaType))
	}
	return out, total, nil
}

func (s *sSearch) Search(ctx context.Context, in service.SearchInput) (*service.ResultDTO, error) {
	kw := strings.TrimSpace(in.Keyword)
	if kw == "" {
		return nil, gerror.New("关键词不能为空")
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}
	// 埋点放 defer: 热搜计数是搜索的副产品, 写失败(库抖动/字段超长)只该丢一次统计,
	// 绝不能让用户搜不出东西。touchHot 内部吞掉所有错误, 只打日志。
	defer s.touchHot(ctx, kw)

	res := &service.ResultDTO{}
	if in.Type > 0 {
		for _, src := range schSources {
			if src.mediaType != in.Type {
				continue
			}
			list, total, err := s.queryOne(ctx, src, kw, in.Page, in.Size, true)
			if err != nil {
				return nil, err
			}
			res.List, res.Total, res.TotalHit = list, total, total
			return res, nil
		}
		return nil, gerror.New("资源类型非法")
	}

	// type=0: 每类各取 schAggEach 条 + 各自命中数, 汇总成 total_hit。
	// 5 类串行查(每类 count+list 两条 SQL), 单条都走索引, 相比并发省下 goroutine 与
	// 连接争用; 真扛不住时再改并发, 而不是现在就上复杂度。
	for _, src := range schSources {
		list, total, err := s.queryOne(ctx, src, kw, 1, schAggEach, true)
		if err != nil {
			return nil, err
		}
		res.TotalHit += total
		switch src.mediaType {
		case paywall.MediaVideo:
			res.Videos = list
		case paywall.MediaPost:
			res.Posts = list
		case paywall.MediaComics:
			res.Comics = list
		case paywall.MediaNovel:
			res.Novels = list
		case paywall.MediaPhoto:
			res.Photos = list
		}
	}
	return res, nil
}

// touchHot 搜索埋点: hot_search 里没有该词就插一条(heat=0, status=1), 有则计数 +1。
//
// 实现方式是「InsertIgnore + 无条件递增」两步, 而不是"先 Count 再决定插/更"。
// 00030 迁移给 hot_search 建了 UNIQUE(site_id, keyword), PG 驱动会把 InsertIgnore
// 编译成 ON CONFLICT DO NOTHING, 因此:
//   - 并发同词搜索不会撞唯一约束报错, 也不会插出重复行(先查再插的 TOCTOU 竞态);
//   - 第二步的 UPDATE ... search_count = search_count + 1 是数据库侧自增, 不会丢更新。
//
// 新词插入时 search_count 给 0, 紧接着被递增到 1 —— 语义上"创建这条记录的这次搜索
// 也算一次", 免得新词永远停在 0 显得没人搜过。
// 全程吞错: 调用方在 defer 里用它, 统计失败不能影响搜索结果。
func (s *sSearch) touchHot(ctx context.Context, keyword string) {
	if l := len([]rune(keyword)); l == 0 || l > schKeywordMax {
		return
	}
	if _, err := g.Model("hot_search").Ctx(ctx).Data(g.Map{
		"site_id": schSiteId, "keyword": keyword, "category": "", "heat": 0, "search_count": 0, "status": 1,
	}).InsertIgnore(); err != nil {
		g.Log().Warningf(ctx, "热搜埋点插入失败 keyword=%s: %v", keyword, err)
		return
	}
	if _, err := g.Model("hot_search").Ctx(ctx).
		Where("site_id", schSiteId).Where("keyword", keyword).
		Data(g.Map{
			"search_count": &gdb.Counter{Field: "search_count", Value: 1},
			"updated_at":   gtime.Now(),
		}).Update(); err != nil {
		g.Log().Warningf(ctx, "热搜埋点计数失败 keyword=%s: %v", keyword, err)
	}
}

// Suggest 前缀联想。只出 status=1 的词, 排序与热搜榜一致(人工权重优先, 其次真实搜索量),
// 这样运营置顶的词在联想里也能排前面。
func (s *sSearch) Suggest(ctx context.Context, keyword string) ([]string, error) {
	kw := strings.TrimSpace(keyword)
	out := []string{}
	if kw == "" {
		return out, nil
	}
	var list []*entity.HotSearch
	if err := g.Model("hot_search").Ctx(ctx).
		Where("site_id", schSiteId).Where("status", 1).
		Where("keyword ILIKE ?", likeEscaper.Replace(kw)+"%").
		OrderDesc("heat").OrderDesc("search_count").
		Limit(schSuggestLimit).Scan(&list); err != nil {
		return nil, err
	}
	for _, r := range list {
		out = append(out, r.Keyword)
	}
	return out, nil
}
