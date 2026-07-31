// Package logic — B5 成长配置(任务/签到统计)管理实现。
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/user/domain"
	"github.com/JarvanDante/my_service/internal/modules/user/service"
)

func (s *sUser) AdminTasks(ctx context.Context) ([]*service.AdminTaskDTO, error) {
	list, err := s.repo.TaskListAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*service.AdminTaskDTO, 0, len(list))
	for _, t := range list {
		out = append(out, &service.AdminTaskDTO{
			Id: t.Id, Name: t.Name, Type: t.Type, Description: t.Description,
			MaxNum: t.MaxNum, Reward: t.Reward, Status: t.Status, Sort: t.Sort,
		})
	}
	return out, nil
}

func checkTask(in service.AdminTaskInput) error {
	if in.Name == "" || in.Type == "" {
		return gerror.New("任务名与类型必填")
	}
	if in.MaxNum <= 0 {
		return gerror.New("单日上限必须大于0")
	}
	if in.Reward <= 0 {
		return gerror.New("奖励必须大于0")
	}
	if in.Status < 0 || in.Status > 1 {
		return gerror.New("status 仅支持 0/1")
	}
	return nil
}

func toTaskEntity(in service.AdminTaskInput) *entity.UserTask {
	return &entity.UserTask{
		Id: in.Id, Name: in.Name, Type: in.Type, Description: in.Description,
		MaxNum: in.MaxNum, Reward: in.Reward, Status: in.Status, Sort: in.Sort,
	}
}

func (s *sUser) AdminCreateTask(ctx context.Context, in service.AdminTaskInput) (int64, error) {
	if err := checkTask(in); err != nil {
		return 0, err
	}
	return s.repo.TaskCreate(ctx, toTaskEntity(in))
}

func (s *sUser) AdminUpdateTask(ctx context.Context, in service.AdminTaskInput) error {
	if in.Id <= 0 {
		return gerror.New("任务ID无效")
	}
	if err := checkTask(in); err != nil {
		return err
	}
	t, err := s.repo.FindTask(ctx, in.Id)
	if err != nil {
		return err
	}
	if t == nil {
		return gerror.New("任务不存在")
	}
	return s.repo.TaskUpdate(ctx, toTaskEntity(in))
}

func (s *sUser) AdminDeleteTask(ctx context.Context, id int64) error {
	if id <= 0 {
		return gerror.New("任务ID无效")
	}
	t, err := s.repo.FindTask(ctx, id)
	if err != nil {
		return err
	}
	if t == nil {
		return gerror.New("任务不存在")
	}
	return s.repo.TaskDelete(ctx, id)
}

func (s *sUser) AdminTaskLogs(ctx context.Context, in service.AdminTaskLogInput) ([]*service.AdminTaskLogDTO, int, error) {
	page, size := normalizePage(in.Page, in.Size)
	list, total, err := s.repo.TaskLogList(ctx, domain.TaskLogFilter{
		UserId: in.UserId, TaskId: in.TaskId, Type: in.Type,
		StartDate: in.StartDate, EndDate: in.EndDate,
	}, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*service.AdminTaskLogDTO, 0, len(list))
	for _, l := range list {
		out = append(out, &service.AdminTaskLogDTO{
			Id: l.Id, UserId: l.UserId, TaskId: l.TaskId, Type: l.Type,
			Num: l.Num, LogDate: l.LogDate, CreatedAt: fmtTime(l.CreatedAt),
		})
	}
	return out, total, nil
}

// AdminSignStats 签到统计, yearMonth=0 时取当月。
func (s *sUser) AdminSignStats(ctx context.Context, yearMonth int) (*service.SignStatsDTO, error) {
	if yearMonth == 0 {
		yearMonth = gconv.Int(gtime.Now().Format("Ym"))
	}
	if yearMonth < 200001 || yearMonth > 300012 {
		return nil, gerror.New("year_month 格式应为 YYYYMM")
	}
	userCount, signCount, days, err := s.repo.SignStats(ctx, yearMonth)
	if err != nil {
		return nil, err
	}
	out := make([]service.SignDayCountDTO, 0, len(days))
	for _, d := range days {
		out = append(out, service.SignDayCountDTO{Day: d.Day, Count: d.Count})
	}
	return &service.SignStatsDTO{
		YearMonth: yearMonth, UserCount: userCount, SignCount: signCount, Days: out,
	}, nil
}
