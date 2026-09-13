package boot

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/shared/aesapi"
)

func LoadApiAES(ctx context.Context) {
	raw := strings.TrimSpace(g.Cfg().MustGet(ctx, "api_aes.key", aesapi.DefaultKey).String())
	aesapi.SetKey(raw)
	if raw != "" && aesapi.ActiveKey() != raw {
		g.Log().Warningf(ctx, "api_aes.key 必须是 16 字节, 已回退默认密钥")
	}
}
