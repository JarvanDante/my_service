// Package front 前台 AI 生成任务控制器。
package front

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/JarvanDante/my_service/api/front/aitask/v1"
	"github.com/JarvanDante/my_service/internal/modules/aitask/service"
	"github.com/JarvanDante/my_service/internal/shared/consts"
)

type Controller struct{ svc service.IAiTask }

func New(svc service.IAiTask) *Controller { return &Controller{svc: svc} }

func uid(ctx context.Context) (int64, error) {
	id := ghttp.RequestFromCtx(ctx).GetCtxVar(consts.CtxUserId).Int64()
	if id <= 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "未登录")
	}
	return id, nil
}

// atoiOr 把 string 型筛选参数转成 int; 空串/非法值返回 def(表示"不筛")。
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

// toItem 前台视图: 刻意不带 provider / provider_task_id。
func toItem(d *service.TaskDTO) v1.TaskItem {
	return v1.TaskItem{
		Id: d.Id, TaskNo: d.TaskNo, BizType: d.BizType, TemplateId: d.TemplateId,
		Params: d.Params, InputUrl: d.InputUrl, CostGold: d.CostGold, Status: d.Status,
		Result: d.Result, ErrMsg: d.ErrMsg, RetryCount: d.RetryCount,
		SubmittedAt: d.SubmittedAt, FinishedAt: d.FinishedAt, CreatedAt: d.CreatedAt,
	}
}

func (c *Controller) Templates(ctx context.Context, req *v1.TemplatesReq) (res *v1.TemplatesRes, err error) {
	list, err := c.svc.FrontTemplates(ctx, req.BizType)
	if err != nil {
		return nil, err
	}
	res = &v1.TemplatesRes{List: make([]v1.TemplateItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, v1.TemplateItem{
			Id: d.Id, Name: d.Name, BizType: d.BizType, Cover: d.Cover, Preview: d.Preview,
			Params: d.Params, CostGold: d.CostGold, Sort: d.Sort,
		})
	}
	return res, nil
}

func (c *Controller) Submit(ctx context.Context, req *v1.SubmitReq) (res *v1.SubmitRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	out, err := c.svc.Submit(ctx, userId, service.SubmitInput{
		BizType: req.BizType, TemplateId: req.TemplateId, Params: req.Params,
		InputUrl: req.InputUrl, ClientToken: req.ClientToken,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SubmitRes{Task: toItem(out.Task), Balance: out.Balance, Repeated: out.Repeated}, nil
}

func (c *Controller) Task(ctx context.Context, req *v1.TaskReq) (res *v1.TaskRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	d, err := c.svc.Task(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.TaskRes{Task: toItem(d)}, nil
}

func (c *Controller) Tasks(ctx context.Context, req *v1.TasksReq) (res *v1.TasksRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	list, total, err := c.svc.Tasks(ctx, userId, service.TaskFilter{
		BizType: atoiOr(req.BizType, 0), Status: atoiOr(req.Status, 0),
		Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.TasksRes{Total: total, List: make([]v1.TaskItem, 0, len(list))}
	for _, d := range list {
		res.List = append(res.List, toItem(d))
	}
	return res, nil
}

func (c *Controller) Cancel(ctx context.Context, req *v1.CancelReq) (res *v1.CancelRes, err error) {
	userId, err := uid(ctx)
	if err != nil {
		return nil, err
	}
	refund, bal, err := c.svc.Cancel(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.CancelRes{Refund: refund, Balance: bal}, nil
}
