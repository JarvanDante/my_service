package paas

import (
	"context"
	"encoding/json"
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
	Playlist    string `json:"playlist"`
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

// DramaPlayToken 按剧集短码向媒资中心签发(不缓存)，仅在本站未配置播放网关时兜底使用。
func DramaPlayToken(ctx context.Context, episodeCode string, previewSec int) (string, error) {
	var data struct {
		PlayUrl string `json:"play_url"`
	}
	if err := doOpenPost(ctx, "/open/assets/"+url.PathEscape(episodeCode)+"/play-token",
		g.Map{"preview_sec": previewSec}, &data); err != nil {
		return "", err
	}
	if data.PlayUrl == "" {
		return "", gerror.New("播放网关未配置")
	}
	return data.PlayUrl, nil
}

// dramaAuth 整部短剧在本站的播放授权缓存，只存选用标记和各集清单文件名，不存签名地址。
type dramaAuth struct {
	Picked    bool              `json:"picked"`
	Lists     map[string]string `json:"lists"`
	FetchedAt int64             `json:"fetched_at"`
}

func dramaAuthKey(ctx context.Context, seriesCode string) string {
	return fmt.Sprintf("paas:dramaauth:%s:%s", siteCode(ctx), seriesCode)
}

func loadDramaAuth(ctx context.Context, seriesCode string) (*dramaAuth, error) {
	var a struct {
		Picked   bool           `json:"picked"`
		Episodes []MediaEpisode `json:"episodes"`
	}
	if err := doJSON(ctx, "GET", "/open/assets/"+url.PathEscape(seriesCode), &a); err != nil {
		return nil, err
	}
	out := &dramaAuth{Picked: a.Picked, Lists: map[string]string{}, FetchedAt: time.Now().Unix()}
	for _, e := range a.Episodes {
		name := e.Playlist
		if name == "" {
			name = "master.m3u8"
		}
		out.Lists[e.Id] = name
	}
	if ttl := g.Cfg().MustGet(ctx, "paas.drama_auth_cache_sec", 300).Int64(); ttl > 0 {
		if b, err := json.Marshal(out); err == nil {
			_, _ = g.Redis().Do(ctx, "SETEX", dramaAuthKey(ctx, seriesCode), ttl, string(b))
		}
	}
	return out, nil
}

// DramaPlaylist 校验整部已被本站选用且该集在媒资中心可播，返回清单文件名；同一部剧按整部缓存，集缺失时最多每 30 秒回源一次。
func DramaPlaylist(ctx context.Context, seriesCode, episodeCode string) (string, error) {
	var auth *dramaAuth
	if v, err := g.Redis().Do(ctx, "GET", dramaAuthKey(ctx, seriesCode)); err == nil && !v.IsNil() && v.String() != "" {
		var a dramaAuth
		if json.Unmarshal([]byte(v.String()), &a) == nil {
			auth = &a
		}
	}
	if auth == nil || (auth.Lists[episodeCode] == "" && time.Now().Unix()-auth.FetchedAt >= 30) {
		a, err := loadDramaAuth(ctx, seriesCode)
		if err != nil {
			return "", err
		}
		auth = a
	}
	if !auth.Picked {
		return "", gerror.New("请先选用该媒资")
	}
	name := auth.Lists[episodeCode]
	if name == "" {
		return "", gerror.New("剧集在媒资中心不可播")
	}
	return name, nil
}

// InvalidateDramaAuth 选用或同步剧集后清掉整部授权缓存。
func InvalidateDramaAuth(ctx context.Context, seriesCode string) {
	_, _ = g.Redis().Do(ctx, "DEL", dramaAuthKey(ctx, seriesCode))
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
