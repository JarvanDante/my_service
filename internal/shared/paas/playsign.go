package paas

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type playConf struct {
	base   string
	secret string
	ttl    int64
}

func playCfg(ctx context.Context) playConf {
	ttl := g.Cfg().MustGet(ctx, "play_gateway.token_ttl_sec", 14400).Int64()
	if ttl <= 0 {
		ttl = 14400
	}
	return playConf{
		base:   strings.TrimRight(g.Cfg().MustGet(ctx, "play_gateway.base_url", "").String(), "/"),
		secret: g.Cfg().MustGet(ctx, "play_gateway.secret", "").String(),
		ttl:    ttl,
	}
}

func siteCode(ctx context.Context) string {
	s := strings.TrimSpace(g.Cfg().MustGet(ctx, "site.code", "my").String())
	if s == "" {
		s = "my"
	}
	return strings.ToUpper(s)
}

func playSign(secret, code, site string, exp int64, d int, ip string, iat int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = fmt.Fprintf(mac, "%s|%s|%d|%d|%s|%d", code, site, exp, d, ip, iat)
	return hex.EncodeToString(mac.Sum(nil))
}

func playURL(ctx context.Context, code, site, file string, previewSec int) string {
	cfg := playCfg(ctx)
	if cfg.base == "" || cfg.secret == "" || code == "" || file == "" {
		return ""
	}
	if previewSec < 0 {
		previewSec = 0
	}
	now := time.Now().Unix()
	exp := now + cfg.ttl
	u := fmt.Sprintf("%s/hls/%s/%s?e=%d&s=%s&t=%d&sig=%s",
		cfg.base, url.PathEscape(code), file, exp, url.QueryEscape(site), now, playSign(cfg.secret, code, site, exp, previewSec, "", now))
	if previewSec > 0 {
		u += fmt.Sprintf("&d=%d", previewSec)
	}
	return u
}

// CoverURL 媒资封面走 my_play 签名地址, 下发 cover.bnc(AES 密文)。
// 后台预览把路径改成 cover.jpg 即可, 网关会解密直出; 签名不含文件名。
func CoverURL(ctx context.Context, code string) string {
	return playURL(ctx, code, siteCode(ctx), "cover.bnc", 0)
}

// PageURL 漫画页图走 my_play。objectKey 形如 comics/{code}/ch001/page_001.bnc。
// 统一签发 .bnc：旧库若仍是 jpg，网关会按同名候选加密后下发。
func PageURL(ctx context.Context, code, objectKey string) string {
	rel := comicRelPath(code, objectKey)
	if rel == "" {
		return ""
	}
	rel = forceBncExt(rel)
	return playURL(ctx, code, siteCode(ctx), rel, 0)
}

func comicRelPath(code, key string) string {
	key = strings.TrimLeft(strings.TrimSpace(key), "/")
	if i := strings.Index(key, "?"); i >= 0 {
		key = key[:i]
	}
	if code == "" || key == "" {
		return ""
	}
	prefix := "comics/" + code + "/"
	if !strings.HasPrefix(key, prefix) {
		if j := strings.Index(key, "/"+prefix); j >= 0 {
			key = key[j+1:]
		}
	}
	if !strings.HasPrefix(key, prefix) {
		return ""
	}
	rel := key[len(prefix):]
	if rel == "" || strings.Contains(rel, "..") {
		return ""
	}
	if strings.Count(rel, "/") > 1 {
		return ""
	}
	return rel
}

func forceBncExt(rel string) string {
	low := strings.ToLower(rel)
	if strings.HasSuffix(low, ".bnc") || strings.HasSuffix(low, ".ceb") {
		return rel
	}
	if i := strings.LastIndex(rel, "."); i >= 0 {
		return rel[:i] + ".bnc"
	}
	return rel + ".bnc"
}

// PlaylistURL 媒资播放地址走 my_play 签名清单。
func PlaylistURL(ctx context.Context, code, raw string, previewSec int) string {
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
	if u := playURL(ctx, code, siteCode(ctx), file, previewSec); u != "" {
		return u
	}
	return raw
}

func isMinioHls(u string) bool {
	if strings.Contains(u, "/my-storage/") {
		return false
	}
	return strings.Contains(u, "/my-media/media/hls/") ||
		strings.Contains(u, "/my-media/video/") ||
		strings.Contains(u, "/my-media/cartoon/") ||
		strings.Contains(u, "/my-media/comics/") ||
		strings.Contains(u, ":19000/my-media/") ||
		(strings.Contains(u, "host.docker.internal") && strings.Contains(u, "/my-media/"))
}

func isSiteStorageCover(cover string) bool {
	return strings.Contains(cover, "/my-storage/")
}

// ApplyGatewayURLs 有媒资短码时封面/播放一律改写为网关签名地址。
func ApplyGatewayURLs(ctx context.Context, cover, play, code string) (string, string) {
	return ApplyGatewayURLsPreview(ctx, cover, play, code, 0)
}

// ApplyGatewayURLsPreview 与 ApplyGatewayURLs 相同，但播放清单可带试看秒数(d=)。
func ApplyGatewayURLsPreview(ctx context.Context, cover, play, code string, previewSec int) (string, string) {
	if code == "" {
		return cover, play
	}
	if u := CoverURL(ctx, code); u != "" && !isSiteStorageCover(cover) && (cover == "" || isMinioHls(cover) ||
		strings.Contains(cover, "/hls/"+code+"/cover.jpg") ||
		strings.Contains(cover, "/hls/"+code+"/cover.bnc") ||
		strings.Contains(cover, "/comics/"+code+"/cover.")) {
		cover = u
	}
	if u := PlaylistURL(ctx, code, play, previewSec); u != "" && (play == "" || isMinioHls(play) || strings.Contains(play, "/hls/"+code+"/")) {
		play = u
	}
	return cover, play
}
