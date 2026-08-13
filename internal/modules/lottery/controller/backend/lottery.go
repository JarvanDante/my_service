// Package backend 后台抽奖控制器。
// 筛选参数一律 string 接收: 空串=全部, 有值才转成条件(int 零值区分不了"没传"和"传了0")。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/lottery/v1"
	"github.com/JarvanDante/my_service/internal/modules/lottery/service"
)

type Controller struct{ svc service.ILottery }

func New(svc service.ILottery) *Controller { return &Controller{svc: svc} }

// atoiOr 空串或非法值返回 def(用于"空=全部"的筛选语义)。
func atoiOr(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func atoi64Or(s string, def int64) int64 {
	if s == "" {
		return def
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return v
}

// ---------------------------------------------------------------- 活动

func (c *Controller) ActivityList(ctx context.Context, req *v1.ActivityListReq) (res *v1.ActivityListRes, err error) {
	list, err := c.svc.Activities(ctx, atoiOr(req.Status, -1))
	if err != nil {
		return nil, err
	}
	res = &v1.ActivityListRes{List: make([]v1.ActivityItem, 0, len(list))}
	for _, a := range list {
		res.List = append(res.List, v1.ActivityItem{
			Id: a.Id, Name: a.Name, LotteryType: a.LotteryType, PayType: a.PayType,
			CostGold: a.CostGold, DailyFree: a.DailyFree, DailyLimit: a.DailyLimit,
			Notice: a.Notice, Status: a.Status, CreatedAt: a.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) ActivityCreate(ctx context.Context, req *v1.ActivityCreateReq) (res *v1.ActivityCreateRes, err error) {
	id, err := c.svc.ActivityCreate(ctx, service.ActivityInput{
		Name: req.Name, LotteryType: req.LotteryType, PayType: req.PayType,
		CostGold: req.CostGold, DailyFree: req.DailyFree, DailyLimit: req.DailyLimit,
		Notice: req.Notice, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ActivityCreateRes{Id: id}, nil
}

func (c *Controller) ActivityUpdate(ctx context.Context, req *v1.ActivityUpdateReq) (res *v1.ActivityUpdateRes, err error) {
	if err = c.svc.ActivityUpdate(ctx, service.ActivityInput{
		Id: req.Id, Name: req.Name, PayType: req.PayType, CostGold: req.CostGold,
		DailyFree: req.DailyFree, DailyLimit: req.DailyLimit, Notice: req.Notice,
		Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.ActivityUpdateRes{}, nil
}

func (c *Controller) ActivityDelete(ctx context.Context, req *v1.ActivityDeleteReq) (res *v1.ActivityDeleteRes, err error) {
	if err = c.svc.ActivityDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ActivityDeleteRes{}, nil
}

// ---------------------------------------------------------------- 奖品

func (c *Controller) PrizeList(ctx context.Context, req *v1.PrizeListReq) (res *v1.PrizeListRes, err error) {
	list, err := c.svc.Prizes(ctx, atoi64Or(req.ActivityId, 0))
	if err != nil {
		return nil, err
	}
	res = &v1.PrizeListRes{List: make([]v1.PrizeItem, 0, len(list))}
	for _, p := range list {
		res.List = append(res.List, v1.PrizeItem{
			Id: p.Id, ActivityId: p.ActivityId, Name: p.Name, Cover: p.Cover, Desc: p.Desc,
			Type: p.Type, Amount: p.Amount, CouponTplId: p.CouponTplId, Odds: p.Odds,
			Stock: p.Stock, Awarded: p.Awarded, Rank: p.Rank, Status: p.Status,
			CreatedAt: p.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) PrizeCreate(ctx context.Context, req *v1.PrizeCreateReq) (res *v1.PrizeCreateRes, err error) {
	id, err := c.svc.PrizeCreate(ctx, service.PrizeInput{
		ActivityId: req.ActivityId, Name: req.Name, Cover: req.Cover, Desc: req.Desc,
		Type: req.Type, Amount: req.Amount, CouponTplId: req.CouponTplId,
		Odds: req.Odds, Stock: req.Stock, Rank: req.Rank, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PrizeCreateRes{Id: id}, nil
}

func (c *Controller) PrizeUpdate(ctx context.Context, req *v1.PrizeUpdateReq) (res *v1.PrizeUpdateRes, err error) {
	if err = c.svc.PrizeUpdate(ctx, service.PrizeInput{
		Id: req.Id, Name: req.Name, Cover: req.Cover, Desc: req.Desc, Type: req.Type,
		Amount: req.Amount, CouponTplId: req.CouponTplId, Odds: req.Odds,
		Stock: req.Stock, Rank: req.Rank, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.PrizeUpdateRes{}, nil
}

func (c *Controller) PrizeDelete(ctx context.Context, req *v1.PrizeDeleteReq) (res *v1.PrizeDeleteRes, err error) {
	if err = c.svc.PrizeDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.PrizeDeleteRes{}, nil
}

// ---------------------------------------------------------------- 记录 / 发货

func (c *Controller) HistoryList(ctx context.Context, req *v1.HistoryListReq) (res *v1.HistoryListRes, err error) {
	list, total, err := c.svc.Histories(ctx, service.HistoryFilter{
		UserId:      atoi64Or(req.UserId, 0),
		LotteryType: atoiOr(req.LotteryType, 0),
		PrizeType:   atoiOr(req.PrizeType, 0),
		Status:      atoiOr(req.Status, 0),
		Page:        req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.HistoryListRes{Total: total, List: make([]v1.HistoryItem, 0, len(list))}
	for _, h := range list {
		res.List = append(res.List, v1.HistoryItem{
			Id: h.Id, UserId: h.UserId, Nickname: h.Nickname, ActivityId: h.ActivityId,
			LotteryType: h.LotteryType, PayType: h.PayType, CostGold: h.CostGold,
			PrizeId: h.PrizeId, PrizeName: h.PrizeName, PrizeType: h.PrizeType,
			PrizeText: service.PrizeTypeText(h.PrizeType), PrizeAmount: h.PrizeAmount,
			Status: h.Status, Remark: h.Remark, CreatedAt: h.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) AddrList(ctx context.Context, req *v1.AddrListReq) (res *v1.AddrListRes, err error) {
	list, total, err := c.svc.Addrs(ctx, service.AddrFilter{
		DeliveryStatus: atoiOr(req.DeliveryStatus, -1),
		UserId:         atoi64Or(req.UserId, 0),
		Page:           req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AddrListRes{Total: total, List: make([]v1.AddrItem, 0, len(list))}
	for _, a := range list {
		res.List = append(res.List, v1.AddrItem{
			Id: a.Id, HistoryId: a.HistoryId, UserId: a.UserId, Nickname: a.Nickname,
			PrizeName: a.PrizeName, Receiver: a.Receiver, Phone: a.Phone, Address: a.Address,
			DeliveryStatus: a.DeliveryStatus, ExpressNo: a.ExpressNo, CreatedAt: a.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Ship(ctx context.Context, req *v1.ShipReq) (res *v1.ShipRes, err error) {
	if err = c.svc.Ship(ctx, req.Id, req.ExpressNo); err != nil {
		return nil, err
	}
	return &v1.ShipRes{}, nil
}
