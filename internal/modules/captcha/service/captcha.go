package service

import "context"

type ICaptcha interface {
	Issue(ctx context.Context, ip string) (id, image string, err error)
	Verify(ctx context.Context, id, code string) error
}
