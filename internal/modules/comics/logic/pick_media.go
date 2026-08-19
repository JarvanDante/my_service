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
	"github.com/JarvanDante/my_service/internal/modules/comics/service"
	"github.com/JarvanDante/my_service/internal/shared/paas"
)

func (s *sComics) ListMediaComics(ctx context.Context, page, size int, keyword string) ([]service.MediaComicsDTO, int, error) {
	list, total, err := paas.ListAssets(ctx, page, size, strings.TrimSpace(keyword), 1)
	if err != nil {
		return nil, 0, err
	}
	out := make([]service.MediaComicsDTO, 0, len(list))
	for _, a := range list {
		item := service.MediaComicsDTO{
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

func (s *sComics) PickMedia(ctx context.Context, code string) (int64, error) {
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
	if a.Kind == 0 && len(a.Chapters) == 0 {
		return 0, gerror.New("该媒资不是漫画，请到视频列表选用")
	}
	return s.upsertFromAsset(ctx, a)
}

func (s *sComics) findByMediaCode(ctx context.Context, code string) (*entity.Comics, error) {
	if code == "" {
		return nil, nil
	}
	var r *entity.Comics
	err := g.Model("comics").Ctx(ctx).
		Where("site_id", cmSiteId).Where("media_code", code).Scan(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *sComics) upsertFromAsset(ctx context.Context, a *paas.MediaAsset) (int64, error) {
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
			if _, e := tx.Model("comics").Ctx(ctx).Where("id", id).Data(g.Map{
				"title":         title,
				"cover":         a.CoverUrl,
				"intro":         a.Intro,
				"chapter_count": len(a.Chapters),
				"updated_at":    gtime.Now(),
			}).Update(); e != nil {
				return e
			}
		} else {
			newID, e := tx.Model("comics").Ctx(ctx).Data(g.Map{
				"site_id":       cmSiteId,
				"title":         title,
				"author":        "",
				"cover":         a.CoverUrl,
				"intro":         a.Intro,
				"category":      "", // 分类由子站自定, 同步后待编辑才能上架
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
		return replaceMediaChapters(ctx, tx, id, a.Chapters)
	})
	return id, err
}

func replaceMediaChapters(ctx context.Context, tx gdb.TX, comicsId int64, chapters []paas.MediaChapter) error {
	if _, err := tx.Model("comics_chapter").Ctx(ctx).
		Where("site_id", cmSiteId).Where("comics_id", comicsId).Delete(); err != nil {
		return err
	}
	for _, ch := range chapters {
		pics := make([]service.PicDTO, 0, len(ch.Pages))
		for _, p := range ch.Pages {
			pics = append(pics, service.PicDTO{Url: p.Url, Key: p.Key})
		}
		title := strings.TrimSpace(ch.Title)
		if title == "" {
			title = "第" + strconv.Itoa(ch.Seq) + "话"
		}
		if _, err := tx.Model("comics_chapter").Ctx(ctx).Data(g.Map{
			"site_id":   cmSiteId,
			"comics_id": comicsId,
			"seq":       ch.Seq,
			"title":     title,
			"pics":      encodeJSON(pics),
			"pic_count": len(pics),
			"status":    1,
		}).Insert(); err != nil {
			return err
		}
	}
	_, err := tx.Model("comics").Ctx(ctx).Where("id", comicsId).Data(g.Map{
		"chapter_count": len(chapters), "updated_at": gtime.Now(),
	}).Update()
	return err
}

func (s *sComics) hydrateMedia(ctx context.Context, list []*service.ComicsDTO) {
	for _, d := range list {
		if d == nil || d.MediaCode == "" {
			continue
		}
		if u := paas.CoverURL(ctx, d.MediaCode); u != "" {
			d.Cover = u
			continue
		}
		a, err := paas.AssetDetail(ctx, d.MediaCode)
		if err != nil || a == nil {
			continue
		}
		if a.CoverUrl != "" {
			d.Cover = a.CoverUrl
		}
	}
}

func (s *sComics) refreshPics(ctx context.Context, mediaCode string, seq int, fallback []service.PicDTO) []service.PicDTO {
	if mediaCode == "" {
		return fallback
	}
	a, err := paas.AssetDetail(ctx, mediaCode)
	if err != nil || a == nil {
		return fallback
	}
	for _, ch := range a.Chapters {
		if ch.Seq != seq {
			continue
		}
		out := make([]service.PicDTO, 0, len(ch.Pages))
		for _, p := range ch.Pages {
			out = append(out, service.PicDTO{Url: p.Url, Key: p.Key})
		}
		if len(out) > 0 {
			return out
		}
	}
	return fallback
}
