// Package logic — B4 会员等级(用户组)管理实现。
package logic

import (
	"context"
	"unicode/utf8"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

var userLevelText = map[int]string{
	1: "普通",
	2: "普通+暗网",
}

var promotionTypeText = map[int]string{
	0: "正常价格",
	1: "新人专享",
}

func (s *sUser) AdminGroups(ctx context.Context, name string) ([]*service.UserGroupDTO, error) {
	list, err := s.repo.GroupList(ctx, name)
	if err != nil {
		return nil, err
	}
	out := make([]*service.UserGroupDTO, 0, len(list))
	for _, ug := range list {
		out = append(out, toGroupDTO(ug))
	}
	return out, nil
}

func checkGroup(in service.UserGroupInput) error {
	if in.Name == "" {
		return gerror.New("名称必填")
	}
	if utf8.RuneCountInString(in.TitleHeat) > 5 {
		return gerror.New("会员头衔最多输入5个字符")
	}
	if utf8.RuneCountInString(in.TitleDescription) > 21 {
		return gerror.New("头部会员卡描述最多输入21个字符")
	}
	if in.Rate < -2 || in.Rate > 100 {
		return gerror.New("折扣须在 -2~100")
	}
	if in.DayNum < 1 {
		return gerror.New("可用天数必须大于0")
	}
	if in.Level == 0 {
		in.Level = 1
	}
	if _, ok := userLevelText[in.Level]; !ok {
		return gerror.New("等级配置错误")
	}
	if _, ok := promotionTypeText[in.PromotionType]; !ok {
		return gerror.New("促销类型配置错误")
	}
	if in.Status < 0 || in.Status > 1 {
		return gerror.New("status 仅支持 0/1")
	}
	return nil
}

func toGroupEntity(in service.UserGroupInput) *entity.UserGroup {
	rights := in.Rights
	if rights == "" {
		rights = "[]"
	}
	level := in.Level
	if level == 0 {
		level = 1
	}
	return &entity.UserGroup{
		Id:               in.Id,
		Name:             in.Name,
		Rate:             in.Rate,
		Rights:           rights,
		Remark:           in.Remark,
		Sort:             in.Sort,
		Status:           in.Status,
		Img:              in.Img,
		TitleHeat:        in.TitleHeat,
		TitleDescription: in.TitleDescription,
		TitlePicture:     in.TitlePicture,
		Level:            level,
		PromotionType:    in.PromotionType,
		Price:            in.Price,
		OldPrice:         in.OldPrice,
		DayNum:           in.DayNum,
		GiftNum:          in.GiftNum,
		DownloadNum:      in.DownloadNum,
		DayTips:          in.DayTips,
		PriceTips:        in.PriceTips,
	}
}

func toGroupDTO(ug *entity.UserGroup) *service.UserGroupDTO {
	updated := ""
	if ug.UpdatedAt != nil {
		updated = ug.UpdatedAt.Format("2006-01-02")
	}
	return &service.UserGroupDTO{
		Id:               ug.Id,
		Name:             ug.Name,
		Rate:             ug.Rate,
		Rights:           ug.Rights,
		Remark:           ug.Remark,
		Sort:             ug.Sort,
		Status:           ug.Status,
		Img:              ug.Img,
		TitleHeat:        ug.TitleHeat,
		TitleDescription: ug.TitleDescription,
		TitlePicture:     ug.TitlePicture,
		Level:            ug.Level,
		PromotionType:    ug.PromotionType,
		Price:            ug.Price,
		OldPrice:         ug.OldPrice,
		DayNum:           ug.DayNum,
		GiftNum:          ug.GiftNum,
		DownloadNum:      ug.DownloadNum,
		DayTips:          ug.DayTips,
		PriceTips:        ug.PriceTips,
		UpdatedAt:        updated,
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
