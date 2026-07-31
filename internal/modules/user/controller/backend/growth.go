// Package backend 后台成长配置控制器(B5)。
package backend

import (
	"context"

	v1 "github.com/JarvanDante/my_service/api/backend/user/v1"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

// TaskList 任务列表(含下线)。
func (c *Controller) TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	list, err := c.user.AdminTasks(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.TaskItem, 0, len(list))
	for _, t := range list {
		items = append(items, v1.TaskItem{
			Id: t.Id, Name: t.Name, Type: t.Type, Description: t.Description,
			MaxNum: t.MaxNum, Reward: t.Reward, Status: t.Status, Sort: t.Sort,
		})
	}
	return &v1.TaskListRes{List: items}, nil
}

// TaskCreate 创建任务。
func (c *Controller) TaskCreate(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error) {
	id, err := c.user.AdminCreateTask(ctx, service.AdminTaskInput{
		Name: req.Name, Type: req.Type, Description: req.Description,
		MaxNum: req.MaxNum, Reward: req.Reward, Status: req.Status, Sort: req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateRes{Id: id}, nil
}

// TaskUpdate 更新任务。
func (c *Controller) TaskUpdate(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error) {
	if err = c.user.AdminUpdateTask(ctx, service.AdminTaskInput{
		Id: req.Id, Name: req.Name, Type: req.Type, Description: req.Description,
		MaxNum: req.MaxNum, Reward: req.Reward, Status: req.Status, Sort: req.Sort,
	}); err != nil {
		return nil, err
	}
	return &v1.TaskUpdateRes{}, nil
}

// TaskDelete 删除任务。
func (c *Controller) TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error) {
	if err = c.user.AdminDeleteTask(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.TaskDeleteRes{}, nil
}

// TaskLogList 任务完成记录。
func (c *Controller) TaskLogList(ctx context.Context, req *v1.TaskLogListAdminReq) (res *v1.TaskLogListAdminRes, err error) {
	list, total, err := c.user.AdminTaskLogs(ctx, service.AdminTaskLogInput{
		UserId: req.UserId, TaskId: req.TaskId, Type: req.Type,
		StartDate: req.StartDate, EndDate: req.EndDate, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	items := make([]v1.TaskLogAdminItem, 0, len(list))
	for _, l := range list {
		items = append(items, v1.TaskLogAdminItem{
			Id: l.Id, UserId: l.UserId, TaskId: l.TaskId, Type: l.Type,
			Num: l.Num, LogDate: l.LogDate, CreatedAt: l.CreatedAt,
		})
	}
	return &v1.TaskLogListAdminRes{List: items, Total: total, Page: page, Size: size}, nil
}

// SignStats 签到统计。
func (c *Controller) SignStats(ctx context.Context, req *v1.SignStatsReq) (res *v1.SignStatsRes, err error) {
	dto, err := c.user.AdminSignStats(ctx, req.YearMonth)
	if err != nil {
		return nil, err
	}
	days := make([]v1.SignDayCount, 0, len(dto.Days))
	for _, d := range dto.Days {
		days = append(days, v1.SignDayCount{Day: d.Day, Count: d.Count})
	}
	return &v1.SignStatsRes{
		YearMonth: dto.YearMonth, UserCount: dto.UserCount, SignCount: dto.SignCount, Days: days,
	}, nil
}
