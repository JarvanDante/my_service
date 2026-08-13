// Package logic 商品兑换业务(移植自 tianbi redeem)。
// 防超卖: stock 条件递减(stock>0 或 -1 不限量); 防透支: balance 条件扣款;
// 均以影响行数判定, 事务内完成并写 user_balance_log 流水。
package logic

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/redeemgoods/service"
)

const rgSiteId = 1

type sRedeemGoods struct{}

func New() service.IRedeemGoods { return &sRedeemGoods{} }

func toDTO(r *entity.RedeemGoods) *service.GoodsDTO {
	created := ""
	if r.CreatedAt != nil {
		created = r.CreatedAt.String()
	}
	return &service.GoodsDTO{
		Id: r.Id, Name: r.Name, Cover: r.Cover, Intro: r.Intro, CostGold: r.CostGold,
		Stock: r.Stock, Exchanged: r.Exchanged, Rank: r.Rank, Status: r.Status,
		CreatedAt: created,
	}
}

func (s *sRedeemGoods) FrontList(ctx context.Context, page, size int) ([]*service.GoodsDTO, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.Model("redeem_goods").Ctx(ctx).Where("site_id", rgSiteId).Where("status", 1)
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RedeemGoods
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.GoodsDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

// Exchange 兑换。
func (s *sRedeemGoods) Exchange(ctx context.Context, userId, goodsId int64) (int64, error) {
	var goods *entity.RedeemGoods
	if err := g.Model("redeem_goods").Ctx(ctx).
		Where("site_id", rgSiteId).Where("id", goodsId).Where("status", 1).
		Scan(&goods); err != nil {
		return 0, err
	}
	if goods == nil {
		return 0, gerror.New("商品不存在或已下架")
	}
	if goods.Stock == 0 {
		return 0, gerror.New("商品已兑完")
	}
	var orderId int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 余额条件扣款(防透支)
		balRes, err := tx.Model("users").Ctx(ctx).
			Where("id", userId).Where("balance >= ?", goods.CostGold).
			Data(g.Map{"balance": &gdb.Counter{Field: "balance", Value: -goods.CostGold}}).Update()
		if err != nil {
			return err
		}
		if n, _ := balRes.RowsAffected(); n == 0 {
			return gerror.New("金币余额不足")
		}
		// 2. 库存条件递减(防超卖; -1 不限量只加已兑换数)
		stockCond := tx.Model("redeem_goods").Ctx(ctx).Where("id", goods.Id)
		data := g.Map{
			"exchanged":  &gdb.Counter{Field: "exchanged", Value: 1},
			"updated_at": gtime.Now(),
		}
		if goods.Stock > 0 {
			stockCond = stockCond.Where("stock > 0")
			data["stock"] = &gdb.Counter{Field: "stock", Value: -1}
		}
		stRes, err := stockCond.Data(data).Update()
		if err != nil {
			return err
		}
		if n, _ := stRes.RowsAffected(); n == 0 {
			return gerror.New("商品已兑完")
		}
		// 3. 兑换记录
		oid, err := tx.Model("redeem_goods_order").Ctx(ctx).Data(g.Map{
			"site_id": rgSiteId, "user_id": userId, "goods_id": goods.Id,
			"goods_name": goods.Name, "cost_gold": goods.CostGold,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		orderId = oid
		// 4. 余额流水
		one, err := tx.Model("users").Ctx(ctx).Where("id", userId).Fields("balance").One()
		if err != nil {
			return err
		}
		after := one["balance"].Float64()
		_, err = tx.Model("user_balance_log").Ctx(ctx).Data(g.Map{
			"site_id": rgSiteId, "user_id": userId, "direction": 2, "scene": "redeem_goods",
			"amount": goods.CostGold, "balance_before": after + goods.CostGold,
			"balance_after": after, "ref_id": gconv.String(oid),
			"remark": "兑换商品:" + goods.Name,
		}).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}
	return orderId, nil
}

func (s *sRedeemGoods) History(ctx context.Context, userId int64, page, size int) ([]*service.OrderDTO, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var list []*entity.RedeemGoodsOrder
	if err := g.Model("redeem_goods_order").Ctx(ctx).
		Where("site_id", rgSiteId).Where("user_id", userId).
		OrderDesc("id").Page(page, size).Scan(&list); err != nil {
		return nil, err
	}
	out := make([]*service.OrderDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.OrderDTO{
			Id: r.Id, UserId: r.UserId, GoodsId: r.GoodsId, GoodsName: r.GoodsName,
			CostGold: r.CostGold, CreatedAt: created,
		})
	}
	return out, nil
}

func (s *sRedeemGoods) List(ctx context.Context, f service.ListFilter) ([]*service.GoodsDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("redeem_goods").Ctx(ctx).Where("site_id", rgSiteId)
	if f.Status >= 0 {
		m = m.Where("status", f.Status)
	}
	if f.Keyword != "" {
		m = m.Where("name ILIKE ?", "%"+f.Keyword+"%")
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RedeemGoods
	if err := m.Clone().OrderDesc("rank").OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.GoodsDTO, 0, len(list))
	for _, r := range list {
		out = append(out, toDTO(r))
	}
	return out, total, nil
}

func (s *sRedeemGoods) Create(ctx context.Context, in service.SaveInput) (int64, error) {
	if strings.TrimSpace(in.Name) == "" {
		return 0, gerror.New("商品名不能为空")
	}
	if in.CostGold <= 0 {
		return 0, gerror.New("金币价需大于0")
	}
	if in.Status != 0 && in.Status != 1 {
		in.Status = 1
	}
	id, err := g.Model("redeem_goods").Ctx(ctx).Data(g.Map{
		"site_id": rgSiteId, "name": in.Name, "cover": in.Cover, "intro": in.Intro,
		"cost_gold": in.CostGold, "stock": in.Stock, "rank": in.Rank, "status": in.Status,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *sRedeemGoods) Update(ctx context.Context, in service.SaveInput) error {
	if in.Id <= 0 {
		return gerror.New("ID非法")
	}
	data := g.Map{"rank": in.Rank, "stock": in.Stock, "updated_at": gtime.Now()}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Cover != "" {
		data["cover"] = in.Cover
	}
	if in.Intro != "" {
		data["intro"] = in.Intro
	}
	if in.CostGold > 0 {
		data["cost_gold"] = in.CostGold
	}
	if in.Status == 0 || in.Status == 1 {
		data["status"] = in.Status
	}
	_, err := g.Model("redeem_goods").Ctx(ctx).
		Where("site_id", rgSiteId).Where("id", in.Id).Data(data).Update()
	return err
}

func (s *sRedeemGoods) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("ID非法")
	}
	_, err := g.Model("redeem_goods").Ctx(ctx).
		Where("site_id", rgSiteId).Where("id", id).Delete()
	return err
}

func (s *sRedeemGoods) Orders(ctx context.Context, f service.ListFilter) ([]*service.OrderDTO, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Size <= 0 {
		f.Size = 20
	}
	m := g.Model("redeem_goods_order").Ctx(ctx).Where("site_id", rgSiteId)
	if f.UserId > 0 {
		m = m.Where("user_id", f.UserId)
	}
	if f.GoodsId > 0 {
		m = m.Where("goods_id", f.GoodsId)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return nil, 0, err
	}
	var list []*entity.RedeemGoodsOrder
	if err := m.Clone().OrderDesc("id").Page(f.Page, f.Size).Scan(&list); err != nil {
		return nil, 0, err
	}
	out := make([]*service.OrderDTO, 0, len(list))
	for _, r := range list {
		created := ""
		if r.CreatedAt != nil {
			created = r.CreatedAt.String()
		}
		out = append(out, &service.OrderDTO{
			Id: r.Id, UserId: r.UserId, GoodsId: r.GoodsId, GoodsName: r.GoodsName,
			CostGold: r.CostGold, CreatedAt: created,
		})
	}
	return out, total, nil
}
