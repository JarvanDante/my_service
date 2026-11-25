package paas

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaDrama 媒资中心短剧整部；短剧只走开放接口，不回退后台接口以免绕过选用校验。
type MediaDrama struct {
	Id           string         `json:"id"`
	Title        string         `json:"title"`
	CoverUrl     string         `json:"cover_url"`
	Intro        string         `json:"intro"`
	ChapterCount int            `json:"chapter_count"`
	EpisodeCount int            `json:"episode_count"`
	UpdateStatus int            `json:"update_status"`
	Picked       bool           `json:"picked"`
	Episodes     []MediaEpisode `json:"episodes"`
}

type MediaEpisode struct {
	Id          string `json:"id"`
	Seq         int    `json:"seq"`
	Title       string `json:"title"`
	DurationSec int    `json:"duration_sec"`
	CoverUrl    string `json:"cover_url"`
}

const mediaKindDrama = 5

func ListDramas(ctx context.Context, page, size int, keyword string) ([]MediaDrama, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("size", fmt.Sprintf("%d", size))
	q.Set("kind", fmt.Sprintf("%d", mediaKindDrama))
	if keyword != "" {
		q.Set("keyword", keyword)
	}
	var data struct {
		List  []MediaDrama `json:"list"`
		Total int          `json:"total"`
	}
	if err := doJSON(ctx, "GET", "/open/assets?"+q.Encode(), &data); err != nil {
		return nil, 0, err
	}
	return data.List, data.Total, nil
}

// PickDrama 选用整部（幂等），返回最新就绪剧集，同步剧集也走这里。
func PickDrama(ctx context.Context, code string) (*MediaDrama, error) {
	var a MediaDrama
	if err := doOpenPost(ctx, "/open/assets/"+url.PathEscape(code)+"/pick", g.Map{}, &a); err != nil {
		return nil, err
	}
	if a.Id == "" {
		return nil, gerror.New("媒资不存在")
	}
	return &a, nil
}

// DramaPlayToken 按剧集短码签发，媒资中心会校验整部已被本站选用；结果按集+试看秒数短时缓存以节省站点签发额度。
func DramaPlayToken(ctx context.Context, episodeCode string, previewSec int) (string, error) {
	cacheKey := fmt.Sprintf("paas:dramaplay:%s:%s:%d", siteCode(ctx), episodeCode, previewSec)
	if v, err := g.Redis().Do(ctx, "GET", cacheKey); err == nil && !v.IsNil() && v.String() != "" {
		return v.String(), nil
	}
	var data struct {
		PlayUrl   string `json:"play_url"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := doOpenPost(ctx, "/open/assets/"+url.PathEscape(episodeCode)+"/play-token",
		g.Map{"preview_sec": previewSec}, &data); err != nil {
		return "", err
	}
	if data.PlayUrl == "" {
		return "", gerror.New("播放网关未配置")
	}
	ttl := g.Cfg().MustGet(ctx, "paas.play_token_cache_sec", 300).Int64()
	if left := data.ExpiresAt - time.Now().Unix() - 60; left < ttl {
		ttl = left
	}
	if ttl > 0 {
		_, _ = g.Redis().Do(ctx, "SETEX", cacheKey, ttl, data.PlayUrl)
	}
	return data.PlayUrl, nil
}

func doOpenPost(ctx context.Context, path string, body any, out any) error {
	base, err := mediaBase(ctx)
	if err != nil {
		return err
	}
	key := g.Cfg().MustGet(ctx, "paas.app_key").String()
	secret := g.Cfg().MustGet(ctx, "paas.app_secret").String()
	if key == "" || secret == "" {
		return gerror.New("未配置 paas.app_key / app_secret")
	}
	r, err := g.Client().SetHeader("X-App-Key", key).SetHeader("X-App-Secret", secret).
		ContentJson().Post(ctx, base+path, body)
	if err != nil {
		return gerror.Wrap(err, "请求媒资中心失败")
	}
	defer r.Close()
	return parseEnvelope(r.ReadAll(), out)
}
