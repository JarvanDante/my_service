// Package ratelimit 基于 Redis 的固定窗口限流(分钟窗 + 小时窗双限)。
// Redis 不可用时 fail-open(放行)并告警, 避免限流器故障拖垮前台可用性。
package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Allow 记一次并判断是否放行。bucket 区分场景(ip/user), id 为主体标识。
// 任一窗超限即拒。返回: 放行 / 建议重试秒 / 命中原因。
func Allow(ctx context.Context, bucket, id string, perMin, perHour int) (allowed bool, retryAfter int, reason string) {
	now := time.Now().Unix()
	if perMin > 0 {
		key := fmt.Sprintf("rl:%s:%s:m:%d", bucket, id, now/60)
		n, err := incr(ctx, key, 120)
		if err != nil {
			g.Log().Warningf(ctx, "ratelimit redis err(min) key=%s: %v; fail-open", key, err)
			return true, 0, ""
		}
		if int(n) > perMin {
			return false, 60, "请求过于频繁(每分钟上限)"
		}
	}
	if perHour > 0 {
		key := fmt.Sprintf("rl:%s:%s:h:%d", bucket, id, now/3600)
		n, err := incr(ctx, key, 3700)
		if err != nil {
			g.Log().Warningf(ctx, "ratelimit redis err(hour) key=%s: %v; fail-open", key, err)
			return true, 0, ""
		}
		if int(n) > perHour {
			return false, 3600, "请求过于频繁(每小时上限)"
		}
	}
	return true, 0, ""
}

// incr INCR + 首次设置 TTL。用 Do 以兼容不同 gf 版本的 gredis API。
func incr(ctx context.Context, key string, ttlSec int) (int64, error) {
	v, err := g.Redis().Do(ctx, "INCR", key)
	if err != nil {
		return 0, err
	}
	n := v.Int64()
	if n == 1 {
		_, _ = g.Redis().Do(ctx, "EXPIRE", key, ttlSec)
	}
	return n, nil
}
