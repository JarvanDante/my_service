package kit

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	tokenKeyPrefix = "auth:token:"        // token  -> userId
	uidKeyPrefix   = "auth:uid:"          // userId -> 当前有效 token(用于单点登录踢旧)
	tokenTTL       = int64(7 * 24 * 3600) // 7 天(秒)
)

func normalize(token string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(token), "Bearer "))
}

// IssueToken 签发新 token(单点登录: 先失效该用户的旧 token, 再签发新的)。
func IssueToken(ctx context.Context, userId int64) (string, error) {
	uidKey := uidKeyPrefix + gconv.String(userId)

	// 1) 踢掉旧 token
	if old, err := g.Redis().Do(ctx, "GET", uidKey); err == nil && !old.IsNil() && old.String() != "" {
		_, _ = g.Redis().Do(ctx, "DEL", tokenKeyPrefix+old.String())
	}

	// 2) 签发新 token, 同时写正向和反向映射
	token := guid.S()
	if _, err := g.Redis().Do(ctx, "SETEX", tokenKeyPrefix+token, tokenTTL, userId); err != nil {
		return "", err
	}
	if _, err := g.Redis().Do(ctx, "SETEX", uidKey, tokenTTL, token); err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 校验 token, 返回 userId, 并对正/反向 key 滑动续期。
func ParseToken(ctx context.Context, token string) (int64, error) {
	token = normalize(token)
	if token == "" {
		return 0, gerror.New("empty token")
	}
	tokenKey := tokenKeyPrefix + token
	v, err := g.Redis().Do(ctx, "GET", tokenKey)
	if err != nil {
		return 0, err
	}
	if v.IsNil() || v.String() == "" {
		return 0, gerror.New("token not found or expired")
	}
	uid := v.Int64()
	// 滑动续期
	_, _ = g.Redis().Do(ctx, "EXPIRE", tokenKey, tokenTTL)
	_, _ = g.Redis().Do(ctx, "EXPIRE", uidKeyPrefix+gconv.String(uid), tokenTTL)
	return uid, nil
}

// RevokeToken 退出登录: 删除正向和反向 key。
func RevokeToken(ctx context.Context, token string) error {
	token = normalize(token)
	if token == "" {
		return nil
	}
	tokenKey := tokenKeyPrefix + token
	if v, err := g.Redis().Do(ctx, "GET", tokenKey); err == nil && !v.IsNil() && v.String() != "" {
		_, _ = g.Redis().Do(ctx, "DEL", uidKeyPrefix+v.String())
	}
	_, err := g.Redis().Do(ctx, "DEL", tokenKey)
	return err
}

// RevokeByUserId 按用户ID撤销其当前 token(退出登录, 单点场景直接用)。
func RevokeByUserId(ctx context.Context, userId int64) error {
	uidKey := uidKeyPrefix + gconv.String(userId)
	if v, err := g.Redis().Do(ctx, "GET", uidKey); err == nil && !v.IsNil() && v.String() != "" {
		_, _ = g.Redis().Do(ctx, "DEL", tokenKeyPrefix+v.String())
	}
	_, err := g.Redis().Do(ctx, "DEL", uidKey)
	return err
}
