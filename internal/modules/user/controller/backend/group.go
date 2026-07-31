// Package backend 后台用户组配置控制器(B4)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

// GroupList 用户组列表。
func (c *Controller) GroupList(ctx context.Context, req *v1.GroupListReq) (res *v1.GroupListRes, err error) {
	list, err := c.user.AdminGroups(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.UserGroupItem, 0, len(list))
	for _, ug := range list {
		items = append(items, v1.UserGroupItem{
			Id: ug.Id, Name: ug.Name, Rate: ug.Rate, Rights: ug.Rights,
			Remark: ug.Remark, Sort: ug.Sort, Status: ug.Status,
		})
	}
	return &v1.GroupListRes{List: items}, nil
}

// GroupCreate 创建用户组。
func (c *Controller) GroupCreate(ctx context.Context, req *v1.GroupCreateReq) (res *v1.GroupCreateRes, err error) {
	id, err := c.user.AdminCreateGroup(ctx, service.UserGroupInput{
		Name: req.Name, Rate: req.Rate, Rights: req.Rights,
		Remark: req.Remark, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.GroupCreateRes{Id: id}, nil
}

// GroupUpdate 更新用户组(同步组内用户快照)。
func (c *Controller) GroupUpdate(ctx context.Context, req *v1.GroupUpdateReq) (res *v1.GroupUpdateRes, err error) {
	if err = c.user.AdminUpdateGroup(ctx, service.UserGroupInput{
		Id: req.Id, Name: req.Name, Rate: req.Rate, Rights: req.Rights,
		Remark: req.Remark, Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.GroupUpdateRes{}, nil
}

// GroupDelete 删除用户组。
func (c *Controller) GroupDelete(ctx context.Context, req *v1.GroupDeleteReq) (res *v1.GroupDeleteRes, err error) {
	if err = c.user.AdminDeleteGroup(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.GroupDeleteRes{}, nil
}
