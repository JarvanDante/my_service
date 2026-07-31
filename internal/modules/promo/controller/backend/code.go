// Package backend 后台兑换码控制器(B3)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/promo/v1"
	"github.com/JarvanDante/my_service/internal/modules/promo/service"
)

type Controller struct{ promo service.IPromo }

func New(svc service.IPromo) *Controller { return &Controller{promo: svc} }

func (c *Controller) CodeList(ctx context.Context, req *v1.CodeListReq) (res *v1.CodeListRes, err error) {
	dto, err := c.promo.Codes(ctx, service.CodeListInput{
		Keyword: req.Keyword, CodeKey: req.CodeKey, Type: req.Type, Status: req.Status,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.CodeItem, 0, len(dto.List))
	for _, c0 := range dto.List {
		items = append(items, v1.CodeItem{
			Id: c0.Id, Name: c0.Name, Code: c0.Code, CodeKey: c0.CodeKey,
			Type: c0.Type, ObjectId: c0.ObjectId, AddNum: c0.AddNum,
			CanUseNum: c0.CanUseNum, UsedNum: c0.UsedNum, Status: c0.Status,
			ExpiredAt: c0.ExpiredAt, CreatedAt: c0.CreatedAt,
		})
	}
	return &v1.CodeListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}

func (c *Controller) CodeGen(ctx context.Context, req *v1.CodeGenReq) (res *v1.CodeGenRes, err error) {
	dto, err := c.promo.GenerateCodes(ctx, service.CodeGenInput{
		Name: req.Name, Type: req.Type, ObjectId: req.ObjectId, AddNum: req.AddNum,
		CanUseNum: req.CanUseNum, Count: req.Count, ExpiredAt: req.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CodeGenRes{CodeKey: dto.CodeKey, Count: dto.Count, Codes: dto.Codes}, nil
}

func (c *Controller) CodeVoid(ctx context.Context, req *v1.CodeVoidReq) (res *v1.CodeVoidRes, err error) {
	if err = c.promo.VoidCode(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CodeVoidRes{}, nil
}

func (c *Controller) CodeLogs(ctx context.Context, req *v1.CodeLogListReq) (res *v1.CodeLogListRes, err error) {
	dto, err := c.promo.CodeLogs(ctx, service.CodeLogListInput{
		CodeId: req.CodeId, UserId: req.UserId, Code: req.Code, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	items := make([]v1.CodeLogItem, 0, len(dto.List))
	for _, l := range dto.List {
		items = append(items, v1.CodeLogItem{
			Id: l.Id, CodeId: l.CodeId, Code: l.Code, Name: l.Name, Type: l.Type,
			UserId: l.UserId, Username: l.Username, AddNum: l.AddNum, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.CodeLogListRes{List: items, Total: dto.Total, Page: dto.Page, Size: dto.Size}, nil
}
