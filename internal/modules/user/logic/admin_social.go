// Package logic — B6 社交查询实现。
package logic

import (
	"context"

	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

func (s *sUser) AdminFollows(ctx context.Context, userId, homeId int64, page, size int) ([]*service.FollowAdminDTO, int, error) {
	page, size = normalizePage(page, size)
	list, total, err := s.repo.FollowList(ctx, domain.FollowFilter{UserId: userId, HomeId: homeId}, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.FollowAdminDTO, 0, len(list))
	for _, f := range list {
		out = append(out, &service.FollowAdminDTO{
			Id: f.Id, UserId: f.UserId, UserName: f.UserName,
			HomeId: f.HomeId, HomeName: f.HomeName, CreatedAt: f.CreatedAt,
		})
	}
	return out, total, nil
}
