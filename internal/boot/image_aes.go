package boot

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/shared/aesbnc"
)

func LoadImageAES(ctx context.Context) {
	raw := strings.TrimSpace(g.Cfg().MustGet(ctx, "image_aes.key", aesbnc.DefaultKey).String())
	aesbnc.SetKey(raw)
	if raw != "" && aesbnc.ActiveKey() != raw {
		g.Log().Warningf(ctx, "image_aes.key 必须是 16 字节, 已回退默认密钥")
	}
}
