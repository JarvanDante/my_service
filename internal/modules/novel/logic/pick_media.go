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

func (s *sNovel) ListMediaNovels(ctx context.Context, page, size int, keyword string) ([]service.MediaNovelDTO, int, error) {
	list, total, err := paas.ListAssets(ctx, page, size, strings.TrimSpace(keyword), 4)
	if err != nil {
		return nil, 0, err
	}
	out := make([]service.MediaNovelDTO, 0, len(list))
	for _, a := range list {
		item := service.MediaNovelDTO{
			Id: a.Id, Title: a.Title, CoverUrl: a.CoverUrl, Intro: a.Intro,
			ChapterCount: a.ChapterCount, Picked: a.Picked,
		}
		if local, _ := s.findByMediaCode(ctx, a.Id); local != nil {
			item.LocalId = local.Id
			item.Picked = true
		}
		out = append(out, item)
	}
	return out, total, nil
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
