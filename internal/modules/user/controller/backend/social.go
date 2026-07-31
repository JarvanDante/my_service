// Package backend 后台社交查询控制器(B6)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
)

// FollowList 关注关系查询。
func (c *Controller) FollowList(ctx context.Context, req *v1.FollowListReq) (res *v1.FollowListRes, err error) {
	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	list, total, err := c.user.AdminFollows(ctx, req.UserId, req.HomeId, page, size)
	if err != nil {
		return nil, err
	}
	items := make([]v1.FollowItem, 0, len(list))
	for _, f := range list {
		items = append(items, v1.FollowItem{
			Id: f.Id, UserId: f.UserId, UserName: f.UserName,
			HomeId: f.HomeId, HomeName: f.HomeName, CreatedAt: f.CreatedAt,
		})
	}
	return &v1.FollowListRes{List: items, Total: total, Page: page, Size: size}, nil
}
