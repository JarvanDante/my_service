// Package logic — B4 用户组定义管理实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

func (s *sUser) AdminGroups(ctx context.Context) ([]*service.UserGroupDTO, error) {
	list, err := s.repo.GroupList(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.UserGroupDTO, 0, len(list))
	for _, ug := range list {
		out = append(out, &service.UserGroupDTO{
			Id: ug.Id, Name: ug.Name, Rate: ug.Rate, Rights: ug.Rights,
			Remark: ug.Remark, Sort: ug.Sort, Status: ug.Status,
		})
	}
	return out, nil
}

func checkGroup(in service.UserGroupInput) error {
	if in.Name == "" {
		return gerror.New("组名必填")
	}
	if in.Rate < 0 || in.Rate > 100 {
		return gerror.New("折扣须在 0~100")
	}
	if in.Status < 0 || in.Status > 1 {
		return gerror.New("status 仅支持 0/1")
	}
	return nil
}

func toGroupEntity(in service.UserGroupInput) *entity.UserGroup {
	rights := in.Rights
	if rights == "" {
		rights = "{}"
	}
	return &entity.UserGroup{
		Id: in.Id, Name: in.Name, Rate: in.Rate, Rights: rights,
		Remark: in.Remark, Sort: in.Sort, Status: in.Status,
	}
}

func (s *sUser) AdminCreateGroup(ctx context.Context, in service.UserGroupInput) (int64, error) {
	if err := checkGroup(in); err != nil {
		return 0, err
	}
	return s.repo.GroupCreate(ctx, toGroupEntity(in))
}

// AdminUpdateGroup 更新组定义并同步组内用户快照(dao 内事务)。
func (s *sUser) AdminUpdateGroup(ctx context.Context, in service.UserGroupInput) error {
	if in.Id <= 0 {
		return gerror.New("组ID无效")
	}
	if err := checkGroup(in); err != nil {
		return err
	}
	ug, err := s.repo.GroupFind(ctx, in.Id)
	if err != nil {
		return err
	}
	if ug == nil {
		return gerror.New("用户组不存在")
	}
	return s.repo.GroupUpdate(ctx, toGroupEntity(in))
}

// AdminDeleteGroup 删除组定义(组内仍有用户则拒绝)。
func (s *sUser) AdminDeleteGroup(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("组ID无效")
	}
	ug, err := s.repo.GroupFind(ctx, id)
	if err != nil {
		return err
	}
	if ug == nil {
		return gerror.New("用户组不存在")
	}
	cnt, err := s.repo.GroupUserCount(ctx, id)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return gerror.New("该组下仍有用户, 不能删除")
	}
	return s.repo.GroupDelete(ctx, id)
}
