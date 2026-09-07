package logic

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const totpPendingTTL = int64(10 * 60)

type totpPending interface {
	Get(ctx context.Context, adminId int64) (string, error)
	Set(ctx context.Context, adminId int64, secret string) error
	Del(ctx context.Context, adminId int64) error
}

func totpPendingKey(adminId int64) string {
	return "admin:totp:pending:" + gconv.String(adminId)
}

type redisTotpPending struct{}

func (redisTotpPending) Get(ctx context.Context, adminId int64) (string, error) {
	v, err := g.Redis().Do(ctx, "GET", totpPendingKey(adminId))
	if err != nil {
		return "", err
	}
	if v.IsNil() || v.String() == "" {
		return "", nil
	}
	return v.String(), nil
}

func (redisTotpPending) Set(ctx context.Context, adminId int64, secret string) error {
	_, err := g.Redis().Do(ctx, "SETEX", totpPendingKey(adminId), totpPendingTTL, secret)
	return err
}

func (redisTotpPending) Del(ctx context.Context, adminId int64) error {
	_, err := g.Redis().Do(ctx, "DEL", totpPendingKey(adminId))
	return err
}

type memTotpPending struct {
	mu sync.Mutex
	m  map[int64]string
}

func newMemTotpPending() *memTotpPending {
	return &memTotpPending{m: make(map[int64]string)}
}

func (p *memTotpPending) Get(_ context.Context, adminId int64) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.m[adminId], nil
}

func (p *memTotpPending) Set(_ context.Context, adminId int64, secret string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.m[adminId] = secret
	return nil
}

func (p *memTotpPending) Del(_ context.Context, adminId int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.m, adminId)
	return nil
}
