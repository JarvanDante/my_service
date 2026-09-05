package logic

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var channelNoiseRe = regexp.MustCompile(`\b(official|404)\b`)

const maxChannelNameLen = 64

// ParseChannel 解析渠道名: 去掉 channel:// / agent:// 前缀, 过滤官方占位词。
func ParseChannel(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "sign=xxx", "")
	s = strings.Trim(s, "&")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "channel://")
	s = strings.TrimPrefix(s, "agent://")
	s = strings.TrimSpace(s)
	s = channelNoiseRe.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if utf8.RuneCountInString(s) > maxChannelNameLen {
		rs := []rune(s)
		s = string(rs[:maxChannelNameLen])
	}
	return s
}

func parseShareParent(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "share://")
	return strings.TrimSpace(s)
}
