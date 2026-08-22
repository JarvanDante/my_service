package kit

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// 与 H5 src/utils/idcrypt.ts 一致：前台 URL / 邀请码用的公开串，不是数据库主键。
const publicIdKey = "myh5v1k"

func xorPublicId(in []byte) []byte {
	key := []byte(publicIdKey)
	out := make([]byte, len(in))
	for i, b := range in {
		out[i] = b ^ key[i%len(key)]
	}
	return out
}

// EncodePublicId 数字 id → 前台加密串（帖子/视频/用户/邀请码）。
func EncodePublicId(id int64) string {
	if id <= 0 {
		return ""
	}
	raw := fmt.Sprintf("%d:%s", id, strconv.FormatInt(id*31+7, 36))
	return base64.RawURLEncoding.EncodeToString(xorPublicId([]byte(raw)))
}

// DecodePublicId 还原前台加密串；无法识别时返回 0。
func DecodePublicId(token string) int64 {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		padded := token
		if m := len(token) % 4; m != 0 {
			padded += strings.Repeat("=", 4-m)
		}
		raw, err = base64.URLEncoding.DecodeString(padded)
		if err != nil {
			return 0
		}
	}
	parts := strings.SplitN(string(xorPublicId(raw)), ":", 2)
	if len(parts) != 2 {
		return 0
	}
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || n <= 0 {
		return 0
	}
	if strconv.FormatInt(n*31+7, 36) != parts[1] {
		return 0
	}
	return n
}
