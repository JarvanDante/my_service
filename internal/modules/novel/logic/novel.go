// Package logic 小说业务(与 comics 同构, 只换章节载体)。
//
// 关键取舍:
//   - 付费同样是「整部购买」, 章节级只有"前 N 章免费";
//   - 解锁判定全部走 shared/paywall, 与漫画/图集/视频同一套语义;
//   - 购买链路的价格取自库里的 price, 绝不接受客户端传值;
//   - 章节正文存 text 整读整写, 字数由正文实时算(人工填必然会和正文对不上);
//   - 作品的 word_count 是全书汇总值, 在章节增删改后按实际行重算, 不做增量累加。
package logic

import (
	"context"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/novel/service"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

const nvSiteId = 1

type sNovel struct{}

func New() service.INovel { return &sNovel{} }

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

// wordCount 按 rune 计数, 中文一字算一字(len() 会把一个汉字算成 3)。
func wordCount(s string) int {
	return utf8.RuneCountInString(s)
}

func fmtTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.String()
}

func toDTO(r *entity.Novel) *service.NovelDTO {
	return &service.NovelDTO{
		Id: r.Id, Title: r.Title, Author: r.Author, Cover: r.Cover, Intro: r.Intro,
		Category: r.Category, Tags: decodeTags(r.Tags), IsVip: r.IsVip, Price: r.Price,
		FreeChapter: r.FreeChapter, ChapterCount: r.ChapterCount, WordCount: r.WordCount,
		IsAudio: r.IsAudio, ViewCount: r.ViewCount, BuyCount: r.BuyCount, LikeCount: r.LikeCount,
		UpdateStatus: r.UpdateStatus, Rank: r.Rank, Status: r.Status, PublishId: r.PublishId,
		CreatedAt: fmtTime(r.CreatedAt),
	}
}

func toChapterDTO(r *entity.NovelChapter) *service.ChapterDTO {
	return &service.ChapterDTO{
		Id: r.Id, NovelId: r.NovelId, Seq: r.Seq, Title: r.Title,
		Content: r.Content, WordCount: r.WordCount, AudioUrl: r.AudioUrl,
		Status: r.Status, CreatedAt: fmtTime(r.CreatedAt),
	}
}

func (s *sNovel) query(ctx context.Context, f service.ListFilter) ([]*service.NovelDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	base := g.Model("novel").Ctx(ctx).Where("site_id", nvSiteId)
	if f.Status >= 0 {
		base = base.Where("status", f.Status)
	}
	if f.Category != "" {
		base = base.Where("category", f.Category)
	}
	if f.Tag != "" {
		base = base.Where("tags @> ?::jsonb", encodeJSON([]string{f.Tag}))
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		base = base.Where("(title ILIKE ? OR author ILIKE ?)", kw, kw)
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
		m = m.OrderDesc("id")
	case 3:
		m = m.OrderDesc("like_count")
	default:
		m = m.OrderDesc("rank")
	}
	var list []*entity.Novel
	if err := m.OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.NovelDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sNovel) FrontList(ctx context.Context, userId int64, f service.ListFilter) ([]*service.NovelDTO, int, error) {
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
		bought, err := paywall.PurchasedSet(ctx, userId, paywall.MediaNovel, ids)
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
func (s *sNovel) find(ctx context.Context, id int64, onlineOnly bool) (*entity.Novel, error) {
	m := g.Model("novel").Ctx(ctx).Where("site_id", nvSiteId).Where("id", id)
	if onlineOnly {
		m = m.Where("status", entity.ContentStatusOnline)
	}
	var r *entity.Novel
	if err := m.Scan(&r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, gerror.New("小说不存在或已下架")
	}
	return r, nil
}

func (s *sNovel) Detail(ctx context.Context, userId, id int64) (*service.DetailDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return nil, err
	}
	acc, err := paywall.Check(ctx, userId, paywall.MediaNovel, r.Id, r.IsVip == 1, r.Price)
	if err != nil {
		return nil, err
	}
	// 观看数 +1(失败不影响主流程)
	_, _ = g.Model("novel").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"view_count": &gdb.Counter{Field: "view_count", Value: 1}}).Update()
	d := toDTO(r)
	d.IsBuy = acc.IsBuy
	d.ViewCount++
	return &service.DetailDTO{
		Novel: d, Playable: acc.Playable, NeedPay: acc.NeedPay,
		NeedVip: acc.NeedVip, Enough: acc.Enough, Reason: acc.Reason,
	}, nil
}

func (s *sNovel) Chapters(ctx context.Context, userId, id int64, desc bool) (string, []*service.ChapterDTO, error) {
	r, err := s.find(ctx, id, true)
	if err != nil {
		return "", nil, err
	}
	acc, err := paywall.Check(ctx, userId, paywall.MediaNovel, r.Id, r.IsVip == 1, r.Price)
	if err != nil {
		return "", nil, err
	}
	m := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", id).Where("status", 1)
	if desc {
		m = m.OrderDesc("seq")
	} else {
		m = m.OrderAsc("seq")
	}
	var list []*entity.NovelChapter
	if err := m.Scan(&list); err != nil {
		return "", nil, err
	}
	out := make([]*service.ChapterDTO, 0, len(list))
	for _, c := range list {
		d := toChapterDTO(c)
		d.Content = "" // 目录不下发正文, 否则整本书就白送了
		d.IsFree = c.Seq <= r.FreeChapter
		d.Playable = d.IsFree || acc.Playable
		out = append(out, d)
	}
	return r.Title, out, nil
}

func (s *sNovel) Read(ctx context.Context, userId, chapterId int64) (*service.ReadDTO, error) {
	var c *entity.NovelChapter
	if err := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("id", chapterId).Where("status", 1).
		Scan(&c); err != nil {
		return nil, err
	}
	if c == nil {
		return nil, gerror.New("章节不存在或已下架")
	}
	r, err := s.find(ctx, c.NovelId, true)
	if err != nil {
		return nil, err
	}
	if c.Seq > r.FreeChapter { // 非免费章节才做解锁判定
		acc, err := paywall.Check(ctx, userId, paywall.MediaNovel, r.Id, r.IsVip == 1, r.Price)
		if err != nil {
			return nil, err
		}
		if !acc.Playable {
			return nil, gerror.New(acc.Reason)
		}
	}
	prev, _ := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", c.NovelId).Where("status", 1).
		Where("seq < ?", c.Seq).OrderDesc("seq").Fields("id").Value()
	next, _ := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", c.NovelId).Where("status", 1).
		Where("seq > ?", c.Seq).OrderAsc("seq").Fields("id").Value()
	out := &service.ReadDTO{
		ChapterId: c.Id, NovelId: c.NovelId, Seq: c.Seq, Title: c.Title,
		Content: c.Content, WordCount: c.WordCount, AudioUrl: c.AudioUrl,
	}
	if prev != nil {
		out.PrevId = prev.Int64()
	}
	if next != nil {
		out.NextId = next.Int64()
	}
	return out, nil
}

func (s *sNovel) Buy(ctx context.Context, userId, id int64) (float64, float64, error) {
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
	if err := paywall.Buy(ctx, userId, paywall.MediaNovel, r.Id, r.Title, r.Price); err != nil {
		return 0, 0, err
	}
	_, _ = g.Model("novel").Ctx(ctx).Where("id", r.Id).
		Data(g.Map{"buy_count": &gdb.Counter{Field: "buy_count", Value: 1}}).Update()
	bal, err := g.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").Value()
	if err != nil {
		return r.Price, 0, nil
	}
	return r.Price, bal.Float64(), nil
}

func (s *sNovel) MayLike(ctx context.Context, id int64, size int) ([]*service.NovelDTO, error) {
	if size <= 0 {
		size = 6
	}
	m := g.Model("novel").Ctx(ctx).
		Where("site_id", nvSiteId).Where("status", entity.ContentStatusOnline)
	if id > 0 {
		if r, err := s.find(ctx, id, true); err == nil {
			m = m.Where("id != ?", id)
			if r.Category != "" {
				m = m.Where("category", r.Category)
			}
		}
	}
	var list []*entity.Novel
	if err := m.OrderRandom().Limit(size).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.NovelDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, nil
}

// ---------------- 后台 ----------------

func (s *sNovel) List(ctx context.Context, f service.ListFilter) ([]*service.NovelDTO, int, error) {
	return s.query(ctx, f)
}

func normalize(in *service.SaveInput) error {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return gerror.New("标题不能为空")
	}
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
	return nil
}

func (s *sNovel) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if err := normalize(&in); err != nil {
		return 0, err
	}
	return g.Model("novel").Ctx(ctx).Data(g.Map{
		"site_id": nvSiteId, "title": in.Title, "author": in.Author, "cover": in.Cover,
		"intro": in.Intro, "category": in.Category, "tags": encodeJSON(in.Tags),
		"is_vip": in.IsVip, "price": in.Price, "free_chapter": in.FreeChapter,
		"is_audio": in.IsAudio, "update_status": in.UpdateStatus, "rank": in.Rank,
		"status": in.Status, "publish_id": in.PublishId,
	}).InsertAndGetId()
}

func (s *sNovel) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	// word_count 不在这里改: 它是章节汇总出来的, 人工改会和正文对不上
	data := g.Map{
		"is_vip": in.IsVip, "price": in.Price, "free_chapter": in.FreeChapter,
		"is_audio": in.IsAudio, "rank": in.Rank, "status": in.Status, "updated_at": gtime.Now(),
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
	if in.Category != "" {
		data["category"] = in.Category
	}
	if in.Tags != nil {
		data["tags"] = encodeJSON(in.Tags)
	}
	if in.UpdateStatus == 1 || in.UpdateStatus == 2 {
		data["update_status"] = in.UpdateStatus
	}
	_, err := g.Model("novel").Ctx(ctx).
		Where("site_id", nvSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sNovel) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("novel_chapter").Ctx(ctx).
			Where("site_id", nvSiteId).Where("novel_id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("novel").Ctx(ctx).
			Where("site_id", nvSiteId).Where("id", id).Delete()
		return err
	})
}

func (s *sNovel) Audit(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.New("状态非法")
	}
	res, err := g.Model("novel").Ctx(ctx).
		Where("site_id", nvSiteId).Where("id", id).
		Data(g.Map{"status": status, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("小说不存在")
	}
	return nil
}

func (s *sNovel) ChapterList(ctx context.Context, novelId int64, page, size int) ([]*service.ChapterDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}
	m := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", novelId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.NovelChapter
	if err := m.Clone().OrderAsc("seq").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.ChapterDTO, 0, len(list))
	for _, c := range list {
		out = append(out, toChapterDTO(c))
	}
	return out, total, nil
}

// syncChapterCount 章节增删改后同步作品的章节数与总字数(以实际行为准, 不做增量累加)。
func syncChapterCount(ctx context.Context, tx gdb.TX, novelId int64) error {
	var stat struct {
		Cnt   int   `json:"cnt"`
		Words int64 `json:"words"`
	}
	if err := tx.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", novelId).
		Fields("COUNT(1) AS cnt, COALESCE(SUM(word_count),0) AS words").Scan(&stat); err != nil {
		return err
	}
	_, err := tx.Model("novel").Ctx(ctx).Where("id", novelId).Data(g.Map{
		"chapter_count": stat.Cnt, "word_count": stat.Words, "updated_at": gtime.Now(),
	}).Update()
	return err
}

func (s *sNovel) ChapterCreate(ctx context.Context, in service.ChapterInput) (int64, error) {
	if in.NovelId <= 0 || in.Seq <= 0 {
		return 0, gerror.New("小说ID与章节序号必填")
	}
	if _, err := s.find(ctx, in.NovelId, false); err != nil {
		return 0, err
	}
	dup, err := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", in.NovelId).Where("seq", in.Seq).Count()
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
		id, err := tx.Model("novel_chapter").Ctx(ctx).Data(g.Map{
			"site_id": nvSiteId, "novel_id": in.NovelId, "seq": in.Seq,
			"title": in.Title, "content": in.Content, "word_count": wordCount(in.Content),
			"audio_url": in.AudioUrl, "status": in.Status,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newId = id
		return syncChapterCount(ctx, tx, in.NovelId)
	})
	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *sNovel) ChapterUpdate(ctx context.Context, in service.ChapterInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	// 先取原行: 改正文会改变全书总字数, 需要拿 novel_id 回算(漫画没这一步, 图片数不进作品表)
	var c *entity.NovelChapter
	if err := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("id", in.Id).Scan(&c); err != nil {
		return err
	}
	if c == nil {
		return gerror.New("章节不存在")
	}
	data := g.Map{"status": in.Status, "updated_at": gtime.Now()}
	if in.Seq > 0 {
		data["seq"] = in.Seq
	}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.AudioUrl != "" {
		data["audio_url"] = in.AudioUrl
	}
	contentChanged := in.Content != ""
	if contentChanged {
		data["content"] = in.Content
		data["word_count"] = wordCount(in.Content)
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("novel_chapter").Ctx(ctx).
			Where("site_id", nvSiteId).Where("id", in.Id).Data(data).Update(); err != nil {
			return err
		}
		if !contentChanged {
			return nil
		}
		return syncChapterCount(ctx, tx, c.NovelId)
	})
}

func (s *sNovel) ChapterDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	var c *entity.NovelChapter
	if err := g.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("id", id).Scan(&c); err != nil {
		return err
	}
	if c == nil {
		return gerror.New("章节不存在")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("novel_chapter").Ctx(ctx).Where("id", id).Delete(); err != nil {
			return err
		}
		return syncChapterCount(ctx, tx, c.NovelId)
	})
}
