// Package v1 后台成长配置接口契约(B5): 任务 CRUD / 任务记录 / 签到统计。
package v1

import "github.com/gogf/gf/v2/frame/g"

// ---------- 任务 CRUD ----------

type TaskItem struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	MaxNum      int     `json:"max_num"`
	Reward      float64 `json:"reward"`
	Status      int     `json:"status"` // 1上架 0下线
	Sort        int     `json:"sort"`
}

type TaskListReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"Backend/Growth" summary:"任务列表(后台, 含下线)"`
}
type TaskListRes struct {
	List []TaskItem `json:"list"`
}

type TaskCreateReq struct {
	g.Meta      `path:"/tasks" method:"post" tags:"Backend/Growth" summary:"创建任务"`
	Name        string  `json:"name"        v:"required#任务名必填"`
	Type        string  `json:"type"        v:"required#类型必填"`
	Description string  `json:"description"`
	MaxNum      int     `json:"max_num"     v:"required|min:1#单日上限必填|单日上限必须大于0"`
	Reward      float64 `json:"reward"      v:"required|min:0.01#奖励必填|奖励必须大于0"`
	Status      int     `json:"status"      v:"in:0,1#status 仅支持 0/1"`
	Sort        int     `json:"sort"`
}
type TaskCreateRes struct {
	Id int64 `json:"id"`
}

type TaskUpdateReq struct {
	g.Meta      `path:"/tasks/{id}" method:"put" tags:"Backend/Growth" summary:"更新任务"`
	Id          int64   `json:"id"          v:"required|min:1#任务ID必填|任务ID必须大于0"`
	Name        string  `json:"name"        v:"required#任务名必填"`
	Type        string  `json:"type"        v:"required#类型必填"`
	Description string  `json:"description"`
	MaxNum      int     `json:"max_num"     v:"required|min:1#单日上限必填|单日上限必须大于0"`
	Reward      float64 `json:"reward"      v:"required|min:0.01#奖励必填|奖励必须大于0"`
	Status      int     `json:"status"      v:"in:0,1#status 仅支持 0/1"`
	Sort        int     `json:"sort"`
}
type TaskUpdateRes struct{}

type TaskDeleteReq struct {
	g.Meta `path:"/tasks/{id}" method:"delete" tags:"Backend/Growth" summary:"删除任务"`
	Id     int64 `json:"id" v:"required|min:1#任务ID必填|任务ID必须大于0"`
}
type TaskDeleteRes struct{}

// ---------- 任务完成记录 ----------

type TaskLogAdminItem struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	TaskId    int64  `json:"task_id"`
	Type      string `json:"type"`
	Num       int    `json:"num"`
	LogDate   int    `json:"log_date"`
	CreatedAt string `json:"created_at"`
}

type TaskLogListAdminReq struct {
	g.Meta    `path:"/task-logs" method:"get" tags:"Backend/Growth" summary:"任务完成记录"`
	UserId    int64  `json:"user_id"    v:"min:0#user_id 不合法"`
	TaskId    int64  `json:"task_id"    v:"min:0#task_id 不合法"`
	Type      string `json:"type"`
	StartDate int    `json:"start_date"` // log_date YYYYMMDD
	EndDate   int    `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}
type TaskLogListAdminRes struct {
	List  []TaskLogAdminItem `json:"list"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// ---------- 签到统计 ----------

type SignDayCount struct {
	Day   int `json:"day"`
	Count int `json:"count"`
}

type SignStatsReq struct {
	g.Meta    `path:"/sign-stats" method:"get" tags:"Backend/Growth" summary:"签到统计"`
	YearMonth int `json:"year_month" v:"min:0#year_month 不合法"` // YYYYMM, 0=当月
}
type SignStatsRes struct {
	YearMonth int            `json:"year_month"`
	UserCount int            `json:"user_count"` // 签到用户数
	SignCount int            `json:"sign_count"` // 总签到人次
	Days      []SignDayCount `json:"days"`
}
