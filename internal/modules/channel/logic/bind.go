package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/shared/kit"
)

// FindEnabled 按 source 查启用中的渠道，大小写不敏感，返回库里的规范码。
func FindEnabled(ctx context.Context, raw string) (string, error) {
	code := kit.ParseSource(raw)
	if code == "" {
		return "", nil
	}
	var row *entity.PromoChannel
	err := g.Model("promo_channel").Ctx(ctx).
		Where("site_id", channelSiteId).
		Where("status", 1).
		Where("lower(code) = ?", strings.ToLower(code)).
		Scan(&row)
	if err != nil || row == nil {
		return "", err
	}
	return row.Code, nil
}

// BindIfEmpty 仅在用户当前渠道为空且 source 对应启用渠道时写入，已有渠道不改。
func BindIfEmpty(ctx context.Context, userId int64, current, raw string) error {
	if userId <= 0 || strings.TrimSpace(current) != "" {
		return nil
	}
	code, err := FindEnabled(ctx, raw)
	if err != nil || code == "" {
		return err
	}
	_, err = g.Model("users").Ctx(ctx).
		Where("id", userId).
		Where("channel_name", "").
		Data(g.Map{
			"channel_name": code,
			"updated_at":   gtime.Now(),
		}).Update()
	return err
}
