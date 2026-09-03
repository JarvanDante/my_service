package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"

	"github.com/JarvanDante/my_service/internal/modules/captcha/service"
	"github.com/JarvanDante/my_service/internal/shared/ratelimit"
)

const (
	codeTTLSec = 180
	issueMin   = 30
	issueHour  = 200
)

type sCaptcha struct{}

func New() service.ICaptcha { return &sCaptcha{} }

func redisKey(id string) string { return "h5:captcha:" + id }

func (s *sCaptcha) Issue(ctx context.Context, ip string) (string, string, error) {
	if ip != "" {
		if ok, _, reason := ratelimit.Allow(ctx, "captcha", ip, issueMin, issueHour); !ok {
			return "", "", gerror.NewCode(gcode.CodeInvalidOperation, reason)
		}
	}
	code := pickCode()
	img, err := renderDataURI(code)
	if err != nil {
		return "", "", gerror.Wrap(err, "生成验证码失败")
	}
	id := uuid.NewString()
	if err = g.Redis().SetEX(ctx, redisKey(id), code, codeTTLSec); err != nil {
		g.Log().Warningf(ctx, "captcha redis set id=%s: %v", id, err)
		return "", "", gerror.NewCode(gcode.CodeInternalError, "验证码服务暂不可用")
	}
	return id, img, nil
}

func (s *sCaptcha) Verify(ctx context.Context, id, code string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "验证码已过期，请刷新")
	}
	key := redisKey(id)
	v, err := g.Redis().Get(ctx, key)
	if err != nil {
		g.Log().Warningf(ctx, "captcha redis get id=%s: %v", id, err)
		return gerror.NewCode(gcode.CodeInternalError, "验证码服务暂不可用")
	}
	_, _ = g.Redis().Del(ctx, key)
	if v.IsNil() || v.String() == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "验证码已过期，请刷新")
	}
	if !codesEqual(code, v.String()) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "验证码错误")
	}
	return nil
}
