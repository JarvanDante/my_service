package kit

import (
	"regexp"
	"strings"
)

const sourceClipboardPrefix = "channel://"

var sourceCodeRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{1,31}$`)

var reservedSources = map[string]struct{}{
	"official": {},
	"404":      {},
	"invite":   {},
	"source":   {},
	"channel":  {},
	"share":    {},
	"agent":    {},
}

// ParseSource 清洗登录/落地页带来的 source。
// 只剥 channel://（APP 剪贴板），不认 share:// / agent://，避免和邀请混在一起。
func ParseSource(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "share://") || strings.Contains(s, "agent://") {
		return ""
	}
	if i := strings.Index(s, sourceClipboardPrefix); i >= 0 {
		s = strings.TrimSpace(s[i+len(sourceClipboardPrefix):])
	}
	if i := strings.IndexAny(s, "?&#/ \t\n"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if !sourceCodeRe.MatchString(s) {
		return ""
	}
	if _, reserved := reservedSources[strings.ToLower(s)]; reserved {
		return ""
	}
	return s
}

// SourceClipboard 落地页/APP 剪贴板口令，与分享邀请码无关。
func SourceClipboard(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	return sourceClipboardPrefix + code
}
