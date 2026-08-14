// Package backend 后台会员等级配置控制器(B4)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
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

func groupInputFromCreate(req *v1.GroupCreateReq) service.UserGroupInput {
	return service.UserGroupInput{
		Name: req.Name, Rate: req.Rate, Rights: req.Rights, Remark: req.Remark,
		Sort: req.Sort, Status: req.Status, Img: req.Img, TitleHeat: req.TitleHeat,
		TitleDescription: req.TitleDescription, TitlePicture: req.TitlePicture,
		Level: req.Level, PromotionType: req.PromotionType, Price: req.Price,
		OldPrice: req.OldPrice, DayNum: req.DayNum, GiftNum: req.GiftNum,
		DownloadNum: req.DownloadNum, DayTips: req.DayTips, PriceTips: req.PriceTips,
	}
}

func groupInputFromUpdate(req *v1.GroupUpdateReq) service.UserGroupInput {
	in := groupInputFromCreate(&v1.GroupCreateReq{
		Name: req.Name, Rate: req.Rate, Rights: req.Rights, Remark: req.Remark,
		Sort: req.Sort, Status: req.Status, Img: req.Img, TitleHeat: req.TitleHeat,
		TitleDescription: req.TitleDescription, TitlePicture: req.TitlePicture,
		Level: req.Level, PromotionType: req.PromotionType, Price: req.Price,
		OldPrice: req.OldPrice, DayNum: req.DayNum, GiftNum: req.GiftNum,
		DownloadNum: req.DownloadNum, DayTips: req.DayTips, PriceTips: req.PriceTips,
	})
	in.Id = req.Id
	return in
}

func groupItem(ug *service.UserGroupDTO) v1.UserGroupItem {
	disabled := "禁用"
	if ug.Status == 1 {
		disabled = "正常"
	}
	return v1.UserGroupItem{
		Id: ug.Id, Name: ug.Name, TitleHeat: ug.TitleHeat,
		TitleDescription: ug.TitleDescription, TitlePicture: ug.TitlePicture,
		Img: ug.Img, Level: ug.Level, LevelText: userLevelText[ug.Level],
		PromotionType: ug.PromotionType, PromotionTypeText: promotionTypeText[ug.PromotionType],
		Price: ug.Price, OldPrice: ug.OldPrice, Rate: ug.Rate, DayNum: ug.DayNum,
		GiftNum: ug.GiftNum, DownloadNum: ug.DownloadNum, DayTips: ug.DayTips,
		PriceTips: ug.PriceTips, Rights: ug.Rights, Remark: ug.Remark,
		Sort: ug.Sort, Status: ug.Status, IsDisabledText: disabled, UpdatedAt: ug.UpdatedAt,
	}
}

// GroupList 会员等级列表。
func (c *Controller) GroupList(ctx context.Context, req *v1.GroupListReq) (res *v1.GroupListRes, err error) {
	list, err := c.user.AdminGroups(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	items := make([]v1.UserGroupItem, 0, len(list))
	for _, ug := range list {
		items = append(items, groupItem(ug))
	}
	return &v1.GroupListRes{List: items}, nil
}

// GroupCreate 创建会员等级。
func (c *Controller) GroupCreate(ctx context.Context, req *v1.GroupCreateReq) (res *v1.GroupCreateRes, err error) {
	id, err := c.user.AdminCreateGroup(ctx, groupInputFromCreate(req))
	if err != nil {
		return nil, err
	}
	return &v1.GroupCreateRes{Id: id}, nil
}

// GroupUpdate 更新会员等级(同步组内用户快照)。
func (c *Controller) GroupUpdate(ctx context.Context, req *v1.GroupUpdateReq) (res *v1.GroupUpdateRes, err error) {
	if err = c.user.AdminUpdateGroup(ctx, groupInputFromUpdate(req)); err != nil {
		return nil, err
	}
	return &v1.GroupUpdateRes{}, nil
}

// GroupDelete 删除会员等级。
func (c *Controller) GroupDelete(ctx context.Context, req *v1.GroupDeleteReq) (res *v1.GroupDeleteRes, err error) {
	if err = c.user.AdminDeleteGroup(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.GroupDeleteRes{}, nil
}
