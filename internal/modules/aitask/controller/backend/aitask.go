// Package backend 后台 AI 模板/任务控制器 + 供应商回调控制器。
package backend

import (
	"context"
	"strconv"

	v1 "github.com/JarvanDante/my_service/api/backend/aitask/v1"
	"github.com/JarvanDante/my_service/internal/modules/aitask/service"
)

type Controller struct{ svc service.IAiTask }

func New(svc service.IAiTask) *Controller { return &Controller{svc: svc} }

// atoiOr 后台筛选参数一律 string 接收: 空串=不筛, 这样 "0" 也能作为合法筛选值传进来。
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

func bizTypeText(bizType int, params map[string]any) string {
	switch bizType {
	case 1:
		if media, _ := params["media_type"].(string); media == "video" {
			return "视频换脸"
		}
		return "图片换脸"
	case 2:
		return "脱衣"
	case 3:
		return "文生图"
	case 4:
		return "图生视频"
	case 5:
		return "文生小说"
	case 6:
		return "AI对话"
	default:
		return strconv.Itoa(bizType)
	}
}

func statusText(status int) string {
	switch status {
	case 1:
		return "待处理"
	case 2:
		return "处理中"
	case 3:
		return "成功"
	case 4:
		return "异常"
	case 5:
		return "退款"
	case 6:
		return "已取消"
	default:
		return strconv.Itoa(status)
	}
}

func toTaskItem(d *service.TaskDTO) v1.TaskItem {
	sets := d.Sets
	if sets <= 0 {
		sets = 1
	}
	return v1.TaskItem{
		Id: d.Id, UserId: d.UserId, TaskNo: d.TaskNo, ClientToken: d.ClientToken,
		BizType: d.BizType, BizTypeText: bizTypeText(d.BizType, d.Params),
		TemplateId: d.TemplateId, Params: d.Params, InputUrl: d.InputUrl,
		CostGold: d.CostGold, Status: d.Status, StatusText: statusText(d.Status),
		Provider: d.Provider, ProviderTaskId: d.ProviderTaskId, Result: d.Result,
		ErrMsg: d.ErrMsg, RetryCount: d.RetryCount, SubmittedAt: d.SubmittedAt,
		FinishedAt: d.FinishedAt, CreatedAt: d.CreatedAt,
		Nickname: d.Nickname, Phone: d.Phone, Avatar: d.Avatar, GroupName: d.GroupName,
		ChannelName: d.ChannelName, DeviceType: d.DeviceType, Sets: sets,
	}
}

// ---------------- 模板 ----------------

func (c *Controller) TemplateList(ctx context.Context, req *v1.TemplateListReq) (res *v1.TemplateListRes, err error) {
	list, total, err := c.svc.TemplateList(ctx, service.TemplateFilter{
		BizType: atoiOr(req.BizType, 0), Status: atoiOr(req.Status, -1),
		Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.TemplateListRes{Total: total, List: make([]v1.TemplateItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.TemplateItem{
			Id: d.Id, Name: d.Name, BizType: d.BizType, Cover: d.Cover, Preview: d.Preview,
			Params: d.Params, CostGold: d.CostGold, Sort: d.Sort, Status: d.Status,
			UsageCount: d.UsageCount, CreatedAt: d.CreatedAt,
		})
	}
	return res, nil
}

func (c *Controller) TemplateCreate(ctx context.Context, req *v1.TemplateCreateReq) (res *v1.TemplateCreateRes, err error) {
	id, err := c.svc.TemplateCreate(ctx, service.TemplateInput{
		Name: req.Name, BizType: req.BizType, Cover: req.Cover, Preview: req.Preview,
		Params: req.Params, CostGold: req.CostGold, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TemplateCreateRes{Id: id}, nil
}

func (c *Controller) TemplateUpdate(ctx context.Context, req *v1.TemplateUpdateReq) (res *v1.TemplateUpdateRes, err error) {
	if err = c.svc.TemplateUpdate(ctx, service.TemplateInput{
		Id: req.Id, Name: req.Name, BizType: req.BizType, Cover: req.Cover,
		Preview: req.Preview, Params: req.Params, CostGold: req.CostGold,
		Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.TemplateUpdateRes{}, nil
}

func (c *Controller) TemplateDelete(ctx context.Context, req *v1.TemplateDeleteReq) (res *v1.TemplateDeleteRes, err error) {
	if err = c.svc.TemplateDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.TemplateDeleteRes{}, nil
}

// ---------------- 任务 ----------------

func (c *Controller) TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	list, total, stats, err := c.svc.TaskList(ctx, service.TaskFilter{
		UserId: atoi64Or(req.UserId, 0), BizType: atoiOr(req.BizType, 0),
		Status: atoiOr(req.Status, 0), TaskNo: req.TaskNo, Nickname: req.Nickname,
		ChannelName: req.ChannelName, DeviceType: req.DeviceType,
		StartTime: req.StartTime, EndTime: req.EndTime,
		RegisterStart: req.RegisterStart, RegisterEnd: req.RegisterEnd,
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.TaskListRes{Total: total, List: make([]v1.TaskItem, 0, len(list))}
	if stats != nil {
		res.Stats = v1.TaskStats{
			Total: stats.Total, Success: stats.Success, Refund: stats.Refund, Abnormal: stats.Abnormal,
			TotalGold: stats.TotalGold, SuccessGold: stats.SuccessGold,
			RefundGold: stats.RefundGold, AbnormalGold: stats.AbnormalGold,
		}
	}
	for _, d := range list {
		res.List = append(res.List, toTaskItem(d))
	}
	return res, nil
}

func (c *Controller) TaskRetry(ctx context.Context, req *v1.TaskRetryReq) (res *v1.TaskRetryRes, err error) {
	d, err := c.svc.TaskRetry(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.TaskRetryRes{Task: toTaskItem(d)}, nil
}

func (c *Controller) TaskRefund(ctx context.Context, req *v1.TaskRefundReq) (res *v1.TaskRefundRes, err error) {
	refund, err := c.svc.TaskRefund(ctx, req.Id, req.Remark)
	if err != nil {
		return nil, err
	}
	return &v1.TaskRefundRes{Refund: refund}, nil
}

func (c *Controller) TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error) {
	if err = c.svc.TaskDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.TaskDeleteRes{}, nil
}

// ---------------- 回调 ----------------

// Callback 供应商回调入口(公开路由, 见 router.RegisterCallback)。
// 这里只做参数搬运, 验签与幂等全在 logic —— 回调的正确性属于业务规则, 不该散落在控制器里。
func (c *Controller) Callback(ctx context.Context, req *v1.CallbackReq) (res *v1.CallbackRes, err error) {
	if err = c.svc.Callback(ctx, service.CallbackInput{
		Provider: req.Provider, ProviderTaskId: req.ProviderTaskId, Status: req.Status,
		Result: req.Result, ErrMsg: req.ErrMsg, Sign: req.Sign,
	}); err != nil {
		return nil, err
	}
	return &v1.CallbackRes{}, nil
}
