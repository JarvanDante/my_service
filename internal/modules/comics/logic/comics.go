// Package logic 漫画业务(移植自 tianbi comicser)。
//
// 关键取舍:
//   - 付费是「整部购买」(与 tianbi 一致), 章节级只有"前 N 章免费";
//   - 解锁判定全部走 shared/paywall, 与小说/图集/视频同一套语义;
//   - 购买链路的价格取自库里的 price, 绝不接受客户端传值;
//   - 章节图片存 jsonb, 读写时整体序列化, 不拆表。
package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/comics/service"
	"github.com/JarvanDante/my_service/internal/shared/kit"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const cmSiteId = 1

type sComics struct{}

func New() service.IComics { return &sComics{} }

func decodeTags(raw string) []string {
	out := []string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func encodeJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodePics(raw string) []service.PicDTO {
	out := []service.PicDTO{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func parseCategories(raw string) []string {
	raw = strings.ReplaceAll(strings.TrimSpace(raw), "，", ",")
	if raw == "" {
		return []string{}
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

func joinCategories(list []string) string {
	return strings.Join(parseCategories(strings.Join(list, ",")), ",")
}

func resolveCategory(in *service.SaveInput) string {
	if in.Categories != nil {
		return joinCategories(in.Categories)
	}
	return joinCategories(parseCategories(in.Category))
}

func toDTO(r *entity.Comics) *service.ComicsDTO {
	cates := parseCategories(r.Category)
	return &service.ComicsDTO{
		Id: r.Id, Title: r.Title, Author: r.Author, Cover: r.Cover, Intro: r.Intro,
		Category: r.Category, Categories: cates, Tags: decodeTags(r.Tags), IsVip: r.IsVip, Price: r.Price,
		FreeChapter: r.FreeChapter, ChapterCount: r.ChapterCount, ViewCount: r.ViewCount,
		BuyCount: r.BuyCount, LikeCount: r.LikeCount, UpdateStatus: r.UpdateStatus,
		Rank: r.Rank, IsRecommend: r.IsRecommend, Status: r.Status, PublishId: r.PublishId, MediaCode: r.MediaCode,
		CreatedAt: fmtTime(r.CreatedAt),
	}
}

func toChapterDTO(r *entity.ComicsChapter) *service.ChapterDTO {
	return &service.ChapterDTO{
		Id: r.Id, ComicsId: r.ComicsId, Seq: r.Seq, Title: r.Title,
		Pics: decodePics(r.Pics), PicCount: r.PicCount, Status: r.Status,
		CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sComics) query(ctx context.Context, f service.ListFilter) ([]*service.ComicsDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	base := g.Model("comics").Ctx(ctx).Where("site_id", cmSiteId)
	if f.Status >= 0 {
		base = base.Where("status", f.Status)
	}
	cates := kit.MergeNames(kit.NamesCSV(f.Category), f.Categories)
	if len(cates) > 0 {
		base = base.Where(kit.CategoryOverlapWhere(), strings.Join(cates, ","))
	}
	tagNames := kit.MergeNames(kit.NamesCSV(f.Tag), f.Tags)
	if len(tagNames) == 1 {
		base = base.Where("tags @> ?::jsonb", encodeJSON([]string{tagNames[0]}))
	} else if len(tagNames) > 1 {
		ors := make([]string, 0, len(tagNames))
		args := make([]any, 0, len(tagNames))
		for _, t := range tagNames {
			ors = append(ors, "tags @> ?::jsonb")
			args = append(args, encodeJSON([]string{t}))
		}
		base = base.Where("("+strings.Join(ors, " OR ")+")", args...)
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		base = base.Where("(title ILIKE ? OR author ILIKE ?)", kw, kw)
	}
	if f.OnlyRecommend {
		base = base.Where("is_recommend", 1)
	}
	switch f.PayType {
	case 1:
		base = base.Where("is_vip", 1)
	case 2:
		base = base.Where("is_vip", 0).Where("price > ?", 0)
	case 3:
		base = base.Where("is_vip", 0).Where("price <= ?", 0)
	}
	if f.Sort == 1 {
		base = base.Where("view_count > 0")
	}
	total, err := base.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	m := base.Clone()
	switch f.Sort {
	case 1:
		m = m.OrderDesc("view_count")
	case 2:
		m = m.OrderDesc("updated_at").OrderDesc("id")
	case 3:
		m = m.OrderDesc("like_count")
	default:
		m = m.OrderDesc("rank")
	}
	var list []*entity.Comics
	if err := m.OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ComicsDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	s.hydrateMedia(ctx, out)
	return out, total, nil
}

func (s *sComics) FrontList(ctx context.Context, userId int64, f service.ListFilter) ([]*service.ComicsDTO, int, error) {
	f.Status = entity.ContentStatusOnline
	list, total, err := s.query(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	if userId > 0 && len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, d := range list {
			ids = append(ids, d.Id)
		}
		bought, err := paywall.PurchasedSet(ctx, userId, paywall.MediaComics, ids)
		if err != nil {
			return nil, 0, err
		}
		for _, d := range list {
			d.IsBuy = bought[d.Id]
		}
	}
	return list, total, nil
}

// find 取上架作品(前台用)。
func (s *sComics) find(ctx context.Context, id int64, onlineOnly bool) (*entity.Comics, error) {
	m := g.Model("comics").Ctx(ctx).Where("site_id", cmSiteId).Where("id", id)
	if onlineOnly {
		m = m.Where("status", entity.ContentStatusOnline)
	}
	var r *entity.Comics
	if err := m.Scan(&r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, gerror.New("漫画不存在或已下架")
	}
	return r, nil
}

func (s *sComics) Detail(ctx context.Context, userId, id int64) (*service.DetailDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return nil, err
	}
	acc, err := paywall.Check(ctx, userId, paywall.MediaComics, r.Id, r.IsVip == 1, r.Price)
	if err != nil {
		return nil, err
	}
	// 观看数 +1(失败不影响主流程)
	_, _ = g.Model("comics").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"view_count": &gdb.Counter{Field: "view_count", Value: 1}}).Update()
	d := toDTO(r)
	s.hydrateMedia(ctx, []*service.ComicsDTO{d})
	d.IsBuy = acc.IsBuy
	d.ViewCount++
	return &service.DetailDTO{
		Comics: d, Playable: acc.Playable, NeedPay: acc.NeedPay,
		NeedVip: acc.NeedVip, Enough: acc.Enough, Reason: acc.Reason,
	}, nil
}

func (s *sComics) Chapters(ctx context.Context, userId, id int64, desc bool) (string, []*service.ChapterDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return "", nil, err
	}
	acc, err := paywall.Check(ctx, userId, paywall.MediaComics, r.Id, r.IsVip == 1, r.Price)
	if err != nil {
		return "", nil, err
	}
	m := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", id).Where("status", 1)
	if desc {
		m = m.OrderDesc("seq")
	} else {
		m = m.OrderAsc("seq")
	}
	var list []*entity.ComicsChapter
	if err := m.Scan(&list); err != nil {
		return "", nil, err
	}
	out := make([]*service.ChapterDTO, 0, len(list))
	for _, c := range list {
		d := toChapterDTO(c)
		d.Pics = nil // 目录不下发图片
		d.IsFree = c.Seq <= r.FreeChapter
		d.Playable = d.IsFree || acc.Playable
		out = append(out, d)
	}
	return r.Title, out, nil
}

func (s *sComics) Read(ctx context.Context, userId, chapterId int64) (*service.ReadDTO, error) {
	var c *entity.ComicsChapter
	if err := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", chapterId).Where("status", 1).
		Scan(&c); err != nil {
		return nil, err
	}
	if c == nil {
		return nil, gerror.New("章节不存在或已下架")
	}
	r, err := s.find(ctx, c.ComicsId, true)
	if err != nil {
		return nil, err
	}
	if c.Seq > r.FreeChapter { // 非免费章节才做解锁判定
		acc, err := paywall.Check(ctx, userId, paywall.MediaComics, r.Id, r.IsVip == 1, r.Price)
		if err != nil {
			return nil, err
		}
		if !acc.Playable {
			return nil, gerror.New(acc.Reason)
		}
	}
	prev, _ := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", c.ComicsId).Where("status", 1).
		Where("seq < ?", c.Seq).OrderDesc("seq").Fields("id").Value()
	next, _ := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", c.ComicsId).Where("status", 1).
		Where("seq > ?", c.Seq).OrderAsc("seq").Fields("id").Value()
	out := &service.ReadDTO{
		ChapterId: c.Id, ComicsId: c.ComicsId, Seq: c.Seq, Title: c.Title,
		Pics: s.refreshPics(ctx, r.MediaCode, c.Seq, decodePics(c.Pics)),
	}
	if prev != nil {
		out.PrevId = prev.Int64()
	}
	if next != nil {
		out.NextId = next.Int64()
	}
	return out, nil
}

func (s *sComics) Buy(ctx context.Context, userId, id int64) (float64, float64, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return 0, 0, err
	}
	if r.IsVip == 1 {
		return 0, 0, gerror.New("该作品为会员专享, 请开通会员")
	}
	if r.Price <= 0 {
		return 0, 0, gerror.New("该作品免费, 无需购买")
	}
	if err := paywall.Buy(ctx, userId, paywall.MediaComics, r.Id, r.Title, r.Price); err != nil {
		return 0, 0, err
	}
	_, _ = g.Model("comics").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"buy_count": &gdb.Counter{Field: "buy_count", Value: 1}}).Update()
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return r.Price, 0, nil
	}
	return r.Price, bal.Float64(), nil
}

func (s *sComics) MayLike(ctx context.Context, id int64, size int) ([]*service.ComicsDTO, error) {
	if size <= 0 {
		size = 6
	}
	m := g.Model("comics").Ctx(ctx).
		Where("site_id", cmSiteId).Where("status", entity.ContentStatusOnline)
	if id > 0 {
		if r, err := s.find(ctx, id, true); err == nil {
			m = m.Where("id != ?", id)
			if r.Category != "" {
				m = m.Where("string_to_array(replace(category, '，', ','), ',') && string_to_array(?, ',')", r.Category)
			}
		}
	}
	var list []*entity.Comics
	if err := m.OrderRandom().Limit(size).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.ComicsDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

// ---------------- 后台 ----------------

func (s *sComics) List(ctx context.Context, f service.ListFilter) ([]*service.ComicsDTO, int, error) {
	return s.query(ctx, f)
}

func normalize(in *service.SaveInput) error {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return gerror.New("标题不能为空")
	}
	in.Category = resolveCategory(in)
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	if in.Price < 0 {
		return gerror.New("价格不能为负")
	}
	if in.FreeChapter < 0 {
		in.FreeChapter = 0
	}
	if in.UpdateStatus != 1 && in.UpdateStatus != 2 {
		in.UpdateStatus = 1
	}
	if in.IsRecommend != 1 {
		in.IsRecommend = 0
	}
	if in.Status == entity.ContentStatusOnline && in.Category == "" {
		return gerror.New("请先选择本站分类后再上架")
	}
	return nil
}

func (s *sComics) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if err := normalize(&in); err != nil {
		return 0, err
	}
	return g.Model("comics").Ctx(ctx).Data(g.Map{
		"site_id": cmSiteId, "title": in.Title, "author": in.Author, "cover": in.Cover,
		"intro": in.Intro, "category": in.Category, "tags": encodeJSON(in.Tags),
		"is_vip": in.IsVip, "price": in.Price, "free_chapter": in.FreeChapter,
		"update_status": in.UpdateStatus, "rank": in.Rank, "is_recommend": in.IsRecommend, "status": in.Status,
		"publish_id": in.PublishId,
	}).InsertAndGetId()
}

func (s *sComics) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	if in.IsRecommend != 1 {
		in.IsRecommend = 0
	}
	in.Category = resolveCategory(&in)
	data := g.Map{
		"is_vip": in.IsVip, "price": in.Price, "free_chapter": in.FreeChapter,
		"rank": in.Rank, "is_recommend": in.IsRecommend, "status": in.Status, "category": in.Category, "updated_at": gtime.Now(),
	}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Author != "" {
		data["author"] = in.Author
	}
	if in.Cover != "" {
		data["cover"] = in.Cover
	}
	if in.Intro != "" {
		data["intro"] = in.Intro
	}
	if in.Tags != nil {
		data["tags"] = encodeJSON(in.Tags)
	}
	if in.UpdateStatus == 1 || in.UpdateStatus == 2 {
		data["update_status"] = in.UpdateStatus
	}
	if in.Status == entity.ContentStatusOnline {
		if in.Category == "" {
			old, err := s.find(ctx, in.Id, false)
			if err != nil {
				return err
			}
			in.Category = strings.TrimSpace(old.Category)
			data["category"] = in.Category
		}
		if in.Category == "" {
			return gerror.New("请先选择本站分类后再上架")
		}
	}
	_, err := g.Model("comics").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sComics) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("comics_chapter").Ctx(ctx).
			Where("site_id", cmSiteId).Where("comics_id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("comics").Ctx(ctx).
			Where("site_id", cmSiteId).Where("id", id).Delete()
		return err
	})
}

func (s *sComics) Audit(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.New("状态非法")
	}
	if status == entity.ContentStatusOnline {
		r, err := s.find(ctx, id, false)
		if err != nil {
			return err
		}
		if strings.TrimSpace(r.Category) == "" {
			return gerror.New("请先编辑并选择本站分类后再上架")
		}
	}
	res, err := g.Model("comics").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", id).
		Data(g.Map{"status": status, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("漫画不存在")
	}
	return nil
}

func (s *sComics) ChapterList(ctx context.Context, comicsId int64, page, size int) ([]*service.ChapterDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}
	m := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", comicsId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.ComicsChapter
	if err := m.Clone().OrderAsc("seq").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ChapterDTO, 0, len(list))
	mediaCode := ""
	if r, e := s.find(ctx, comicsId, false); e == nil && r != nil {
		mediaCode = r.MediaCode
	}
	for _, c := range list {
		d := toChapterDTO(c)
		d.Pics = s.refreshPics(ctx, mediaCode, c.Seq, d.Pics)
		out = append(out, d)
	}
	return out, total, nil
}

// syncChapterCount 章节增删后同步作品的章节数(以实际行数为准, 不做增量累加)。
func syncChapterCount(ctx context.Context, tx gdb.TX, comicsId int64) error {
	n, err := tx.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", comicsId).Count()
	if err != nil {
		return err
	}
	_, err = tx.Model("comics").Ctx(ctx).Where("id", comicsId).
		Data(g.Map{"chapter_count": n, "updated_at": gtime.Now()}).Update()
	return err
}

func (s *sComics) ChapterCreate(ctx context.Context, in service.ChapterInput) (int64, error) {
	if in.ComicsId <= 0 || in.Seq <= 0 {
		return 0, gerror.New("漫画ID与章节序号必填")
	}
	if _, err := s.find(ctx, in.ComicsId, false); err != nil {
		return 0, err
	}
	dup, err := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", in.ComicsId).Where("seq", in.Seq).Count()
	if err != nil {
		return 0, err
	}
	if dup > 0 {
		return 0, gerror.New("该章节序号已存在")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	var newId int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("comics_chapter").Ctx(ctx).Data(g.Map{
			"site_id": cmSiteId, "comics_id": in.ComicsId, "seq": in.Seq,
			"title": in.Title, "pics": encodeJSON(in.Pics), "pic_count": len(in.Pics),
			"status": in.Status,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newId = id
		return syncChapterCount(ctx, tx, in.ComicsId)
	})
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *sComics) ChapterUpdate(ctx context.Context, in service.ChapterInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"status": in.Status, "updated_at": gtime.Now()}
	if in.Seq > 0 {
		data["seq"] = in.Seq
	}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Pics != nil {
		data["pics"] = encodeJSON(in.Pics)
		data["pic_count"] = len(in.Pics)
	}
	res, err := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", in.Id).Data(data).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("章节不存在")
	}
	return nil
}

func (s *sComics) ChapterDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	var c *entity.ComicsChapter
	if err := g.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("id", id).Scan(&c); err != nil {
		return err
	}
	if c == nil {
		return gerror.New("章节不存在")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("comics_chapter").Ctx(ctx).Where("id", id).Delete(); err != nil {
			return err
		}
		return syncChapterCount(ctx, tx, c.ComicsId)
	})
}
