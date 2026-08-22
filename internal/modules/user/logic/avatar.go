// Package logic — 随机默认头像(注册分配)。
// 与公司做法同构: 一张头像 URL 列表 + rand 随机取一张(tianbi 是 picture.yaml 的
// portrait 数组)。差异在于列表放 app_config(key=default_avatars, jsonb 字符串数组),
// 后台改配置即可换头像池、不用发版; 未配置时用内置的 /static/avatar 兜底包。
package logic

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/util/grand"

	"github.com/JarvanDante/my_service/internal/shared/appcfg"
)

// builtinAvatars 内置兜底头像(resource/public/avatar, 由 /static 静态目录提供)。
var builtinAvatars = func() []string {
	out := make([]string, 0, 48)
	for i := 1; i <= 48; i++ {
		out = append(out, fmt.Sprintf("/static/avatar/av%02d.png", i))
	}
	return out
}()

// avatarPool 当前生效的默认头像池: 配置优先, 内置兜底。
func avatarPool(ctx context.Context) []string {
	return appcfg.StringSlice(ctx, "default_avatars", builtinAvatars)
}

// builtinBackgrounds 内置兜底主页背景(resource/public/bg, 16 张渐变图)。
var builtinBackgrounds = func() []string {
	out := make([]string, 0, 16)
	for i := 1; i <= 16; i++ {
		out = append(out, fmt.Sprintf("/static/bg/bg%02d.jpg", i))
	}
	return out
}()

// randBackground 注册用: 随机取一张默认主页背景(配置 default_backgrounds 优先)。
// 对齐公司 getUserbackGround() 的做法。
func randBackground(ctx context.Context) string {
	pool := appcfg.StringSlice(ctx, "default_backgrounds", builtinBackgrounds)
	if len(pool) == 0 {
		return ""
	}
	return pool[grand.Intn(len(pool))]
}

// DefaultAvatars 默认头像列表(前台改资料时挑选, 对齐公司 /user/avatar 接口)。
func (s *sUser) DefaultAvatars(ctx context.Context) []string {
	return avatarPool(ctx)
}

func isDefaultAvatar(ctx context.Context, img string) bool {
	if img == "" {
		return false
	}
	for _, item := range avatarPool(ctx) {
		if item == img {
			return true
		}
	}
	return false
}

// randAvatar 注册用: 随机取一张默认头像。
func randAvatar(ctx context.Context) string {
	pool := avatarPool(ctx)
	if len(pool) == 0 {
		return ""
	}
	return pool[grand.Intn(len(pool))]
}
