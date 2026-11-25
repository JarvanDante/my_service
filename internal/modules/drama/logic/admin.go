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
	"github.com/JarvanDante/my_service/internal/modules/drama/service"
	"github.com/JarvanDante/my_service/internal/shared/paas"
	"github.com/JarvanDante/my_service/internal/shared/paywall"
)

func (s *sDrama) List(ctx context.Context, f service.ListFilter) ([]*service.DramaDTO, int, error) {
	return s.query(ctx, f)
}

func playableEpisodeCount(ctx context.Context, dramaId int64) (int, error) {
	return g.Model("drama_episode").Ctx(ctx).
		Where("site_id", dmSiteId).Where("drama_id", dramaId).Where("status", 1).Where("media_ready", 1).Count()
}

// checkOnline 上架前须有本站分类且至少一集可播。
func checkOnline(ctx context.Context, dramaId int64, category string) error {
	if strings.TrimSpace(category) == "" {
		return gerror.New("请先选择本站分类后再上架")
	}
	n, err := playableEpisodeCount(ctx, dramaId)
	if err != nil {
		return err
	}
	if n == 0 {
		return gerror.New("没有可播剧集, 请先同步剧集")
	}
	return nil
}

func (s *sDrama) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	if in.IsVip == 1 && in.Price > 0 {
		return gerror.New("VIP专享与金币定价互斥, 只能二选一")
	}
	if in.Price < 0 {
		return gerror.New("价格不能为负")
	}
	if in.FreeEpisode < 0 {
		in.FreeEpisode = 0
	}
	if in.IsRecommend != 1 {
		in.IsRecommend = 0
	}
	old, err := s.find(ctx, in.Id, false)
	if err != nil {
		return err
	}
	category := old.Category
	if in.Categories != nil || strings.TrimSpace(in.Category) != "" {
		category = resolveCategory(&in)
	}
	if in.Status == entity.ContentStatusOnline {
		if err := checkOnline(ctx, old.Id, category); err != nil {
			return err
		}
	}
	data := g.Map{
		"is_vip": in.IsVip, "price": in.Price, "free_episode": in.FreeEpisode, "rank": in.Rank,
		"is_recommend": in.IsRecommend, "status": in.Status, "category": category, "updated_at": gtime.Now(),
	}
	if t := strings.TrimSpace(in.Title); t != "" {
		data["title"] = t
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
	_, err = g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sDrama) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	// 已有单集购买时禁止删除，避免重新选用后剧集ID变化导致已购失效。
	bought, err := g.Model("content_purchase").Ctx(ctx).
		Where("site_id", dmSiteId).Where("media_type", paywall.MediaDrama).
		Where("content_id IN (SELECT id FROM drama_episode WHERE site_id = ? AND drama_id = ?)", dmSiteId, id).
		Count()
	if err != nil {
		return err
	}
	if bought > 0 {
		return gerror.New("该短剧已有购买记录, 请下架而不是删除")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model("drama_episode").Ctx(ctx).
			Where("site_id", dmSiteId).Where("drama_id", id).Delete(); err != nil {
			return err
		}
		_, err := tx.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("id", id).Delete()
		return err
	})
}

func (s *sDrama) Audit(ctx context.Context, id int64, status int) error {
	if status < 0 || status > 2 {
		return gerror.New("状态非法")
	}
	r, err := s.find(ctx, id, false)
	if err != nil {
		return err
	}
	if status == entity.ContentStatusOnline {
		if err := checkOnline(ctx, r.Id, r.Category); err != nil {
			return err
		}
	}
	_, err = g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("id", id).
		Data(g.Map{"status": status, "updated_at": gtime.Now()}).Update()
	return err
}

func (s *sDrama) EpisodeList(ctx context.Context, dramaId int64, page, size int) ([]*service.EpisodeDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}
	m := g.Model("drama_episode").Ctx(ctx).Where("site_id", dmSiteId).Where("drama_id", dramaId)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.DramaEpisode
	if err := m.Clone().OrderAsc("seq").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.EpisodeDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toEpisodeDTO(ctx, r))
	}
	return out, total, nil
}

func (s *sDrama) EpisodeUpdate(ctx context.Context, in service.EpisodeInput) error {
	if in.PayMode < entity.DramaPayFollow || in.PayMode > entity.DramaPayVip {
		return gerror.New("付费方式非法")
	}
	if in.Price < 0 {
		return gerror.New("价格不能为负")
	}
	if in.PayMode != entity.DramaPayCoin && in.Price > 0 {
		return gerror.New("只有金币集可以设置单集价")
	}
	if in.Status != 0 && in.Status != 1 {
		return gerror.New("状态非法")
	}
	res, err := g.Model("drama_episode").Ctx(ctx).Where("site_id", dmSiteId).Where("id", in.Id).
		Data(g.Map{"pay_mode": in.PayMode, "price": in.Price, "status": in.Status, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return gerror.New("剧集不存在")
	}
	return nil
}

func (s *sDrama) SyncEpisodes(ctx context.Context, id int64) (*service.SyncDTO, error) {
	r, err := s.find(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if r.MediaCode == "" {
		return nil, gerror.New("该短剧未关联媒资")
	}
	a, err := paas.PickDrama(ctx, r.MediaCode)
	if err != nil {
		return nil, err
	}
	return syncFrom(ctx, r.Id, a)
}

// syncFrom 按集号增量同步：新集插入，已有集只更新媒资字段，媒资侧没了的集置 media_ready=0，不删行。
func syncFrom(ctx context.Context, dramaId int64, a *paas.MediaDrama) (*service.SyncDTO, error) {
	out := &service.SyncDTO{}
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		seqs := make([]int, 0, len(a.Episodes))
		for _, ep := range a.Episodes {
			if ep.Seq <= 0 || ep.Id == "" {
				continue
			}
			title := strings.TrimSpace(ep.Title)
			if title == "" {
				title = "第" + strconv.Itoa(ep.Seq) + "集"
			}
			res, err := tx.Model("drama_episode").Ctx(ctx).Data(g.Map{
				"site_id": dmSiteId, "drama_id": dramaId, "seq": ep.Seq, "title": title,
				"media_code": ep.Id, "duration": ep.DurationSec, "media_ready": 1,
			}).InsertIgnore()
			if err != nil {
				return err
			}
			if n, _ := res.RowsAffected(); n > 0 {
				out.Created++
			} else {
				if _, err := tx.Model("drama_episode").Ctx(ctx).
					Where("site_id", dmSiteId).Where("drama_id", dramaId).Where("seq", ep.Seq).
					Data(g.Map{
						"title": title, "media_code": ep.Id, "duration": ep.DurationSec,
						"media_ready": 1, "updated_at": gtime.Now(),
					}).Update(); err != nil {
					return err
				}
				out.Updated++
			}
			seqs = append(seqs, ep.Seq)
		}
		miss := tx.Model("drama_episode").Ctx(ctx).
			Where("site_id", dmSiteId).Where("drama_id", dramaId).Where("media_ready", 1)
		if len(seqs) > 0 {
			miss = miss.WhereNotIn("seq", seqs)
		}
		res, err := miss.Data(g.Map{"media_ready": 0, "updated_at": gtime.Now()}).Update()
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n > 0 {
			out.Missing = int(n)
		}
		cnt, err := tx.Model("drama_episode").Ctx(ctx).
			Where("site_id", dmSiteId).Where("drama_id", dramaId).Where("media_ready", 1).Count()
		if err != nil {
			return err
		}
		out.Total = cnt
		_, err = tx.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("id", dramaId).
			Data(g.Map{"episode_count": cnt, "updated_at": gtime.Now()}).Update()
		return err
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *sDrama) ListMediaDramas(ctx context.Context, page, size int, keyword string) ([]service.MediaDramaDTO, int, error) {
	list, total, err := paas.ListDramas(ctx, page, size, strings.TrimSpace(keyword))
	if err != nil {
		return nil, 0, err
	}
	codes := make([]string, 0, len(list))
	for _, a := range list {
		codes = append(codes, a.Id)
	}
	local := map[string]int64{}
	if len(codes) > 0 {
		rows, err := g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).WhereIn("media_code", codes).
			Fields("id, media_code").All()
		if err != nil {
			return nil, 0, err
		}
		for _, row := range rows {
			local[row["media_code"].String()] = row["id"].Int64()
		}
	}
	out := make([]service.MediaDramaDTO, 0, len(list))
	for _, a := range list {
		out = append(out, service.MediaDramaDTO{
			Id: a.Id, Title: a.Title, CoverUrl: a.CoverUrl, Intro: a.Intro,
			ChapterCount: a.ChapterCount, EpisodeCount: a.EpisodeCount, UpdateStatus: a.UpdateStatus,
			Picked: a.Picked || local[a.Id] > 0, LocalId: local[a.Id],
		})
	}
	return out, total, nil
}

// PickMedia 选用媒资整部入库（按媒资短码幂等），并首次同步剧集。
func (s *sDrama) PickMedia(ctx context.Context, code string) (int64, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return 0, gerror.New("媒资ID必填")
	}
	a, err := paas.PickDrama(ctx, code)
	if err != nil {
		return 0, err
	}
	title := strings.TrimSpace(a.Title)
	if title == "" {
		title = a.Id
	}
	updateStatus := a.UpdateStatus
	if updateStatus != 1 && updateStatus != 2 {
		updateStatus = 1
	}
	if _, err := g.Model("drama").Ctx(ctx).Data(g.Map{
		"site_id": dmSiteId, "title": title, "intro": a.Intro, "tags": "[]",
		"update_status": updateStatus, "status": entity.ContentStatusPending, "media_code": a.Id,
	}).InsertIgnore(); err != nil {
		return 0, err
	}
	id, err := g.Model("drama").Ctx(ctx).Where("site_id", dmSiteId).Where("media_code", a.Id).Fields("id").Value()
	if err != nil {
		return 0, err
	}
	if id == nil || id.Int64() == 0 {
		return 0, gerror.New("短剧入库失败")
	}
	if _, err := syncFrom(ctx, id.Int64(), a); err != nil {
		return 0, err
	}
	return id.Int64(), nil
}
