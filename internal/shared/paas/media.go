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

func mediaBase(ctx context.Context) (string, error) {
	base := strings.TrimRight(g.Cfg().MustGet(ctx, "paas.media_base").String(), "/")
	if base == "" {
		return "", gerror.New("未配置 paas.media_base")
	}
	return base, nil
}

func parseEnvelope(raw []byte, out any) error {
	s := strings.TrimSpace(string(raw))
	if i := strings.Index(s, "{"); i >= 0 {
		s = s[i:]
	}
	var env envelope
	if err := json.Unmarshal([]byte(s), &env); err != nil {
		msg := s
		if len(msg) > 180 {
			msg = msg[:180]
		}
		return gerror.Newf("媒资中心错误: %s", msg)
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

func doReq(ctx context.Context, method, path string, headers map[string]string, out any) error {
	base, err := mediaBase(ctx)
	if err != nil {
		return err
	}
	c := g.Client()
	for k, v := range headers {
		c = c.SetHeader(k, v)
	}
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
	return parseEnvelope(raw, out)
}

func doJSON(ctx context.Context, method, path string, out any) error {
	key := g.Cfg().MustGet(ctx, "paas.app_key").String()
	secret := g.Cfg().MustGet(ctx, "paas.app_secret").String()
	if key == "" || secret == "" {
		return gerror.New("未配置 paas.app_key / app_secret")
	}
	return doReq(ctx, method, path, map[string]string{
		"X-App-Key": key, "X-App-Secret": secret,
	}, out)
}

func doAdminJSON(ctx context.Context, method, path string, out any) error {
	token := g.Cfg().MustGet(ctx, "paas.media_admin_token").String()
	if token == "" {
		return gerror.New("未配置 paas.media_admin_token")
	}
	return doReq(ctx, method, path, map[string]string{"X-Admin-Token": token}, out)
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
	if err := doJSON(ctx, "GET", "/open/assets?"+q.Encode(), &data); err == nil {
		return data.List, data.Total, nil
	}
	// 本站 paas_client 未同步时, 走媒资中心后台列表(只拉就绪)。
	q.Set("status", "2")
	if err := doAdminJSON(ctx, "GET", "/admin/assets?"+q.Encode(), &data); err != nil {
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
	if err := doJSON(ctx, "POST", "/open/assets/"+url.PathEscape(id)+"/pick", &a); err == nil && a.Id != "" {
		return &a, nil
	}
	if err := doAdminJSON(ctx, "GET", "/admin/assets/"+url.PathEscape(id), &a); err != nil {
		return nil, err
	}
	if a.Id == "" {
		return nil, gerror.New("媒资不存在")
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
