// Package paas 调用平台能力(媒资中心 open API)。
package paas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type MediaAsset struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	CoverUrl    string `json:"cover_url"`
	PlayUrl     string `json:"play_url"`
	PlayKey     string `json:"play_key"`
	DurationSec int    `json:"duration_sec"`
	Picked      bool   `json:"picked"`
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func mediaCfg(ctx context.Context) (base, key, secret string, err error) {
	base = strings.TrimRight(g.Cfg().MustGet(ctx, "paas.media_base").String(), "/")
	key = g.Cfg().MustGet(ctx, "paas.app_key").String()
	secret = g.Cfg().MustGet(ctx, "paas.app_secret").String()
	if base == "" || key == "" || secret == "" {
		return "", "", "", gerror.New("未配置 paas.media_base / app_key / app_secret")
	}
	return base, key, secret, nil
}

func doJSON(ctx context.Context, method, path string, out any) error {
	base, key, secret, err := mediaCfg(ctx)
	if err != nil {
		return err
	}
	c := g.Client().SetHeader("X-App-Key", key).SetHeader("X-App-Secret", secret)
	var raw []byte
	switch method {
	case "GET":
		r, e := c.Get(ctx, base+path)
		if e != nil {
			return gerror.Wrap(e, "请求媒资中心失败")
		}
		defer r.Close()
		raw = r.ReadAll()
	case "POST":
		r, e := c.ContentJson().Post(ctx, base+path, "{}")
		if e != nil {
			return gerror.Wrap(e, "请求媒资中心失败")
		}
		defer r.Close()
		raw = r.ReadAll()
	default:
		return gerror.New("不支持的方法")
	}
	var env envelope
	if err = json.Unmarshal(raw, &env); err != nil {
		return gerror.Wrapf(err, "媒资中心响应无法解析: %s", string(raw))
	}
	if env.Code != 0 {
		msg := env.Message
		if msg == "" {
			msg = fmt.Sprintf("媒资中心错误 %d", env.Code)
		}
		return gerror.New(msg)
	}
	if out == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	return json.Unmarshal(env.Data, out)
}

func ListAssets(ctx context.Context, page, size int, keyword string) ([]MediaAsset, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("size", fmt.Sprintf("%d", size))
	if keyword != "" {
		q.Set("keyword", keyword)
	}
	var data struct {
		List  []MediaAsset `json:"list"`
		Total int          `json:"total"`
	}
	if err := doJSON(ctx, "GET", "/open/assets?"+q.Encode(), &data); err != nil {
		return nil, 0, err
	}
	return data.List, data.Total, nil
}

func ListPicks(ctx context.Context, page, size int) ([]MediaAsset, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("size", fmt.Sprintf("%d", size))
	var data struct {
		List  []MediaAsset `json:"list"`
		Total int          `json:"total"`
	}
	if err := doJSON(ctx, "GET", "/open/picks?"+q.Encode(), &data); err != nil {
		return nil, 0, err
	}
	return data.List, data.Total, nil
}

func AssetDetail(ctx context.Context, id string) (*MediaAsset, error) {
	var a MediaAsset
	if err := doJSON(ctx, "GET", "/open/assets/"+url.PathEscape(id), &a); err != nil {
		return nil, err
	}
	if a.Id == "" {
		return nil, nil
	}
	return &a, nil
}

func PickAsset(ctx context.Context, id string) (*MediaAsset, error) {
	var a MediaAsset
	if err := doJSON(ctx, "POST", "/open/assets/"+url.PathEscape(id)+"/pick", &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func PlayToken(ctx context.Context, id string) (string, error) {
	var data struct {
		PlayUrl string `json:"play_url"`
	}
	if err := doJSON(ctx, "POST", "/open/assets/"+url.PathEscape(id)+"/play-token", &data); err != nil {
		return "", err
	}
	return data.PlayUrl, nil
}
