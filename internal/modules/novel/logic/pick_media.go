package logic

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/novel/service"
	"github.com/JarvanDante/my_service/internal/shared/paas"
)

func (s *sNovel) ListMediaNovels(ctx context.Context, page, size int, keyword string, picked int) ([]service.MediaNovelDTO, int, error) {
	keyword = strings.TrimSpace(keyword)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if picked == 1 {
		return s.listLocalPickedNovels(ctx, page, size, keyword)
	}
	if picked == 2 {
		return s.listUnpickedMediaNovels(ctx, page, size, keyword)
	}
	return s.listAllMediaNovels(ctx, page, size, keyword)
}

func (s *sNovel) listAllMediaNovels(ctx context.Context, page, size int, keyword string) ([]service.MediaNovelDTO, int, error) {
	list, total, err := paas.ListAssets(ctx, page, size, keyword, 4)
	if err != nil {
		return nil, 0, err
	}
	out := make([]service.MediaNovelDTO, 0, len(list))
	for _, a := range list {
		out = append(out, s.mediaNovelOf(ctx, a))
	}
	return out, total, nil
}

func (s *sNovel) listLocalPickedNovels(ctx context.Context, page, size int, keyword string) ([]service.MediaNovelDTO, int, error) {
	m := g.Model("novel").Ctx(ctx).Where("site_id", nvSiteId).Where("media_code <> ?", "")
	if keyword != "" {
		m = m.WhereLike("title", "%"+keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var rows []*entity.Novel
	if err := m.OrderDesc("id").Page(page, size).Scan(&rows); err != nil {
		return nil, 0, err
	}
	out := make([]service.MediaNovelDTO, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		out = append(out, service.MediaNovelDTO{
			Id: r.MediaCode, Title: r.Title, CoverUrl: r.Cover, Intro: r.Intro,
			ChapterCount: r.ChapterCount, Picked: true, LocalId: r.Id,
		})
	}
	return out, total, nil
}

func (s *sNovel) listUnpickedMediaNovels(ctx context.Context, page, size int, keyword string) ([]service.MediaNovelDTO, int, error) {
	local, err := s.localMediaCodes(ctx)
	if err != nil {
		return nil, 0, err
	}
	need := page * size
	collected := make([]service.MediaNovelDTO, 0, need)
	srcTotal := 0
	for p := 1; p <= 30 && len(collected) < need; p++ {
		chunk, total, e := paas.ListAssets(ctx, p, size, keyword, 4)
		if e != nil {
			return nil, 0, e
		}
		srcTotal = total
		if len(chunk) == 0 {
			break
		}
		for _, a := range chunk {
			if local[a.Id] > 0 {
				continue
			}
			item := s.mediaNovelOf(ctx, a)
			item.LocalId = 0
			item.Picked = false
			collected = append(collected, item)
		}
		if p*size >= total {
			break
		}
	}
	start := (page - 1) * size
	if start > len(collected) {
		start = len(collected)
	}
	end := start + size
	if end > len(collected) {
		end = len(collected)
	}
	total := srcTotal - len(local)
	if total < len(collected) {
		total = len(collected)
	}
	if total < 0 {
		total = 0
	}
	return collected[start:end], total, nil
}

func (s *sNovel) mediaNovelOf(ctx context.Context, a paas.MediaAsset) service.MediaNovelDTO {
	item := service.MediaNovelDTO{
		Id: a.Id, Title: a.Title, CoverUrl: a.CoverUrl, Intro: a.Intro,
		ChapterCount: a.ChapterCount, Picked: a.Picked,
	}
	if local, _ := s.findByMediaCode(ctx, a.Id); local != nil {
		item.LocalId = local.Id
		item.Picked = true
	}
	return item
}

func (s *sNovel) localMediaCodes(ctx context.Context) (map[string]int64, error) {
	var rows []*entity.Novel
	err := g.Model("novel").Ctx(ctx).
		Fields("id", "media_code").
		Where("site_id", nvSiteId).Where("media_code <> ?", "").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		if r == nil || r.MediaCode == "" {
			continue
		}
		out[r.MediaCode] = r.Id
	}
	return out, nil
}

func (s *sNovel) PickMedia(ctx context.Context, code string) (int64, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, gerror.New("媒资ID必填")
	}
	a, err := paas.PickAsset(ctx, code)
	if err != nil {
		return 0, err
	}
	if a == nil || a.Id == "" {
		return 0, gerror.New("媒资无效")
	}
	if len(a.Chapters) == 0 {
		if d, e := paas.AssetDetail(ctx, code); e == nil && d != nil {
			if d.Intro != "" {
				a.Intro = d.Intro
			}
			if d.CoverUrl != "" {
				a.CoverUrl = d.CoverUrl
			}
			a.Chapters = d.Chapters
			a.ChapterCount = d.ChapterCount
			a.Kind = d.Kind
		}
	}
	if a.Kind != 4 {
		return 0, gerror.New("该媒资不是小说，请到对应列表选用")
	}
	return s.upsertFromAsset(ctx, a)
}

func (s *sNovel) findByMediaCode(ctx context.Context, code string) (*entity.Novel, error) {
	if code == "" {
		return nil, nil
	}
	var r *entity.Novel
	err := g.Model("novel").Ctx(ctx).
		Where("site_id", nvSiteId).Where("media_code", code).Scan(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *sNovel) upsertFromAsset(ctx context.Context, a *paas.MediaAsset) (int64, error) {
	title := strings.TrimSpace(a.Title)
	if title == "" {
		title = a.Id
	}
	old, err := s.findByMediaCode(ctx, a.Id)
	if err != nil {
		return 0, err
	}

	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if old != nil && old.Id > 0 {
			id = old.Id
			data := g.Map{
				"title":         title,
				"intro":         a.Intro,
				"chapter_count": len(a.Chapters),
				"updated_at":    gtime.Now(),
			}
			if a.CoverUrl != "" && !strings.Contains(old.Cover, "/my-storage/") {
				data["cover"] = a.CoverUrl
			}
			if _, e := tx.Model("novel").Ctx(ctx).Where("id", id).Data(data).Update(); e != nil {
				return e
			}
		} else {
			newID, e := tx.Model("novel").Ctx(ctx).Data(g.Map{
				"site_id":       nvSiteId,
				"title":         title,
				"author":        "",
				"cover":         a.CoverUrl,
				"intro":         a.Intro,
				"category":      "",
				"tags":          "[]",
				"chapter_count": len(a.Chapters),
				"status":        entity.ContentStatusPending,
				"media_code":    a.Id,
			}).InsertAndGetId()
			if e != nil {
				return e
			}
			id = newID
		}
		return replaceMediaNovelChapters(ctx, tx, id, a.Chapters)
	})
	return id, err
}

func replaceMediaNovelChapters(ctx context.Context, tx gdb.TX, novelId int64, chapters []paas.MediaChapter) error {
	if _, err := tx.Model("novel_chapter").Ctx(ctx).
		Where("site_id", nvSiteId).Where("novel_id", novelId).Delete(); err != nil {
		return err
	}
	totalWords := 0
	for _, ch := range chapters {
		title := strings.TrimSpace(ch.Title)
		if title == "" {
			title = "第" + strconv.Itoa(ch.Seq) + "章"
		}
		words := ch.WordCount
		if words <= 0 {
			words = wordCount(ch.Content)
		}
		totalWords += words
		if _, err := tx.Model("novel_chapter").Ctx(ctx).Data(g.Map{
			"site_id":    nvSiteId,
			"novel_id":   novelId,
			"seq":        ch.Seq,
			"title":      title,
			"content":    ch.Content,
			"word_count": words,
			"status":     1,
		}).Insert(); err != nil {
			return err
		}
	}
	_, err := tx.Model("novel").Ctx(ctx).Where("id", novelId).Data(g.Map{
		"chapter_count": len(chapters),
		"word_count":    totalWords,
		"updated_at":    gtime.Now(),
	}).Update()
	return err
}
