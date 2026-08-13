// Package front 前台抽奖控制器。
package front

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/lottery/v1"
	"github.com/JarvanDante/my_service/internal/modules/lottery/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.ILottery }

func New(svc service.ILottery) *Controller { return &Controller{svc: svc} }

// optionalUid 公开接口: 带 token 就取到用户(用于回填我的次数/余额), 不带返回 0。
func optionalUid(ctx context.Context) int64 {
	return ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
}

func uid(ctx context.Context) (int64, error) {
	id := optionalUid(ctx)
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

// Info 活动信息。这里是唯一一处「奖品出网关」的地方, 刻意逐字段拷贝而不是整体转换:
// service.PrizeDTO 里带着 odds(中奖权重), 逐字段拷贝能保证概率永远漏不出去 ——
// 概率是商业机密, 泄漏等于把奖池期望值和作弊思路一起送人。
func (c *Controller) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	d, err := c.svc.Info(ctx, optionalUid(ctx), req.LotteryType)
	if err != nil {
		return nil, err
	}
	res = &v1.InfoRes{
		ActivityId: d.Activity.Id, Name: d.Activity.Name, LotteryType: d.Activity.LotteryType,
		PayType: d.Activity.PayType, CostGold: d.Activity.CostGold,
		DailyFree: d.Activity.DailyFree, DailyLimit: d.Activity.DailyLimit,
		Notice: d.Activity.Notice, FreeLeft: d.FreeLeft, DrawLeft: d.DrawLeft,
		Balance: d.Balance, LoggedIn: d.LoggedIn,
		Prizes:  make([]v1.PrizeItem, 0, len(d.Prizes)),
		Marquee: make([]v1.MarqueeItem, 0, len(d.Marquee)),
	}
	for _, p := range d.Prizes {
		res.Prizes = append(res.Prizes, v1.PrizeItem{
			Id: p.Id, Name: p.Name, Cover: p.Cover, Desc: p.Desc,
			Type: p.Type, Amount: p.Amount, Rank: p.Rank, // 注意: 没有 Odds
		})
	}
	for _, m := range d.Marquee {
		res.Marquee = append(res.Marquee, v1.MarqueeItem{
			Nickname: m.Nickname, PrizeName: m.PrizeName, PrizeType: m.PrizeType,
			Amount: m.PrizeAmount, CreatedAt: m.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Draw(ctx context.Context, req *v1.DrawReq) (res *v1.DrawRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Draw(ctx, userId, req.LotteryType, req.PayType)
	if err != nil {
		return nil, err
	}
	return &v1.DrawRes{
		HistoryId: d.HistoryId, PrizeId: d.PrizeId, PrizeName: d.PrizeName,
		PrizeType: d.PrizeType, PrizeAmount: d.PrizeAmount, PrizeCover: d.PrizeCover,
		PrizeDesc: d.PrizeDesc, Status: d.Status, NeedAddr: d.NeedAddr,
		CostGold: d.CostGold, FreeLeft: d.FreeLeft, DrawLeft: d.DrawLeft, Balance: d.Balance,
	}, nil
}

func (c *Controller) My(ctx context.Context, req *v1.MyReq) (res *v1.MyRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.My(ctx, userId, req.LotteryType, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.MyRes{Total: total, List: make([]v1.MyItem, 0, len(list))}
	for _, h := range list {
		res.List = append(res.List, v1.MyItem{
			Id: h.Id, PrizeName: h.PrizeName, PrizeType: h.PrizeType,
			PrizeText: service.PrizeTypeText(h.PrizeType), PrizeAmount: h.PrizeAmount,
			PayType: h.PayType, CostGold: h.CostGold, Status: h.Status,
			Remark: h.Remark, CreatedAt: h.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) Address(ctx context.Context, req *v1.AddressReq) (res *v1.AddressRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	if err = c.svc.FillAddr(ctx, userId, req.HistoryId, req.Receiver, req.Phone, req.Address); err != nil {
		return nil, err
	}
	return &v1.AddressRes{}, nil
}

func (c *Controller) MyAddr(ctx context.Context, req *v1.MyAddrReq) (res *v1.MyAddrRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.MyAddrs(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	res = &v1.MyAddrRes{Total: total, List: make([]v1.AddrItem, 0, len(list))}
	for _, a := range list {
		res.List = append(res.List, v1.AddrItem{
			Id: a.Id, HistoryId: a.HistoryId, PrizeName: a.PrizeName, Receiver: a.Receiver,
			Phone: a.Phone, Address: a.Address, DeliveryStatus: a.DeliveryStatus,
			ExpressNo: a.ExpressNo, CreatedAt: a.CreatedAt,
		})
	}
	return res, nil
}
