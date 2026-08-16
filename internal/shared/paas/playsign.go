package paas

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type playConf struct {
	base   string
	secret string
	ttl    int64
}

var (
	pc    playConf
	ponce sync.Once
)

func loadPlay() {
	ctx := context.Background()
	pc.base = strings.TrimRight(g.Cfg().MustGet(ctx, "play_gateway.base_url", "").String(), "/")
	pc.secret = g.Cfg().MustGet(ctx, "play_gateway.secret", "").String()
	pc.ttl = g.Cfg().MustGet(ctx, "play_gateway.token_ttl_sec", 14400).Int64()
	if pc.ttl <= 0 {
		pc.ttl = 14400
	}
}

func siteCode(ctx context.Context) string {
	s := g.Cfg().MustGet(ctx, "site.code", "my").String()
	if s == "" {
		return "my"
	}
	return s
}

func playSign(code, site string, exp int64, d int, ip string, iat int64) string {
	mac := hmac.New(sha256.New, []byte(pc.secret))
	_, _ = fmt.Fprintf(mac, "%s|%s|%d|%d|%s|%d", code, site, exp, d, ip, iat)
	return hex.EncodeToString(mac.Sum(nil))
}

func playURL(code, site, file string) string {
	ponce.Do(loadPlay)
	if pc.base == "" || pc.secret == "" || code == "" || file == "" {
		return ""
	}
	now := time.Now().Unix()
	exp := now + pc.ttl
	return fmt.Sprintf("%s/hls/%s/%s?e=%d&s=%s&t=%d&sig=%s",
		pc.base, url.PathEscape(code), file, exp, url.QueryEscape(site), now, playSign(code, site, exp, 0, "", now))
}

// CoverURL 媒资封面走 my_play 签名地址, 由网关代理直出 JPEG, 避免 MinIO 私有桶直链 403。
func CoverURL(ctx context.Context, code string) string {
	return playURL(code, siteCode(ctx), "cover.jpg")
}

// PlaylistURL 媒资播放地址走 my_play 签名清单。
func PlaylistURL(ctx context.Context, code, raw string) string {
	file := "master.m3u8"
	if raw != "" {
		if i := strings.LastIndex(raw, "/"); i >= 0 {
			name := raw[i+1:]
			if j := strings.IndexAny(name, "?#"); j >= 0 {
				name = name[:j]
			}
			if strings.HasSuffix(name, ".m3u8") {
				file = name
			}
		}
	}
	if u := playURL(code, siteCode(ctx), file); u != "" {
		return u
	}
	return raw
}

func isMinioHls(u string) bool {
	return strings.Contains(u, ":19000/") || strings.Contains(u, "/my-media/media/hls/")
}

// ApplyGatewayURLs 有媒资短码时封面/播放一律改写为网关签名地址。
func ApplyGatewayURLs(ctx context.Context, cover, play, code string) (string, string) {
	if code == "" {
		return cover, play
	}
	if u := CoverURL(ctx, code); u != "" && (cover == "" || isMinioHls(cover) || strings.Contains(cover, "/hls/"+code+"/cover.jpg")) {
		cover = u
	}
	if u := PlaylistURL(ctx, code, play); u != "" && (play == "" || isMinioHls(play) || strings.Contains(play, "/hls/"+code+"/")) {
		play = u
	}
	return cover, play
}
