// Package domain 用户领域层。
package domain

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

// AdminUserFilter 后台用户列表筛选条件。
type AdminUserFilter struct {
	Keyword   string // 用户名/手机/昵称 模糊
	Channel   string // 渠道
	GroupId   int64  // 用户组
	Status    int    // 0全部 1正常 2禁用
	StartDate int    // 注册日起 YYYYMMDD
	EndDate   int    // 注册日止 YYYYMMDD
}

// TaskLogFilter 任务记录筛选(B5)。
type TaskLogFilter struct {
	UserId    int64
	TaskId    int64
	Type      string
	StartDate int // log_date YYYYMMDD
	EndDate   int
}

// SignDayCount 某日签到人数(B5)。
type SignDayCount struct {
	Day   int
	Count int
}

// FollowFilter 关注关系筛选(B6)。
type FollowFilter struct {
	UserId int64 // 关注人
	HomeId int64 // 被关注人
}

// FollowItem 关注关系(联表用户名, B6)。
type FollowItem struct {
	Id        int64  `orm:"id"`
	UserId    int64  `orm:"user_id"`
	UserName  string `orm:"user_name"`
	HomeId    int64  `orm:"home_id"`
	HomeName  string `orm:"home_name"`
	CreatedAt string `orm:"created_at"`
}

type Repository interface {
	FindById(ctx context.Context, id int64) (*entity.Users, error)
	FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Users, error)
	FindByAccount(ctx context.Context, account string) (*entity.Users, error)
	Create(ctx context.Context, data g.Map) (int64, error)
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
	UpdatePhone(ctx context.Context, id int64, phone string) error
	UpdateProfile(ctx context.Context, id int64, data g.Map) error
	Disable(ctx context.Context, id int64, reason string) error
	// 后台管理(B1)
	AdminListUsers(ctx context.Context, f AdminUserFilter, page, size int) ([]*entity.Users, int, error)
	SetDisabled(ctx context.Context, id int64, disabled int, reason string) error
	UpdateGroup(ctx context.Context, id, groupId int64, groupName string, groupRate int, groupEndTime int64) error
	AdminAdjustBalance(ctx context.Context, userId int64, target string, amount float64, refId, remark string) error
	BalanceLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserBalanceLog, int, error)
	// 用户组定义(B4)
	GroupList(ctx context.Context) ([]*entity.UserGroup, error)
	GroupFind(ctx context.Context, id int64) (*entity.UserGroup, error)
	GroupCreate(ctx context.Context, g *entity.UserGroup) (int64, error)
	GroupUpdate(ctx context.Context, g *entity.UserGroup) error
	GroupDelete(ctx context.Context, id int64) error
	GroupUserCount(ctx context.Context, groupId int64) (int, error)
	// 成长配置(B5)
	TaskListAll(ctx context.Context) ([]*entity.UserTask, error)
	TaskCreate(ctx context.Context, t *entity.UserTask) (int64, error)
	TaskUpdate(ctx context.Context, t *entity.UserTask) error
	TaskDelete(ctx context.Context, id int64) error
	TaskLogList(ctx context.Context, f TaskLogFilter, page, size int) ([]*entity.UserTaskLog, int, error)
	SignStats(ctx context.Context, yearMonth int) (userCount int, signCount int, days []SignDayCount, err error)
	// 社交查询(B6)
	FollowList(ctx context.Context, f FollowFilter, page, size int) ([]*FollowItem, int, error)
	// 社交
	ExistsFollow(ctx context.Context, userId, homeId int64) (bool, error)
	Follow(ctx context.Context, userId, homeId int64) error
	Unfollow(ctx context.Context, userId, homeId int64) error
	FollowingList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error)
	FansList(ctx context.Context, userId int64, page, size int) ([]*entity.Users, int, error)
	// 推广 / 兑换码
	BindInviter(ctx context.Context, userId, inviterId int64, inviterName string) error
	FindCodeByCode(ctx context.Context, code string) (*entity.UserCode, error)
	HasRedeemed(ctx context.Context, codeId, userId int64) (bool, error)
	RedeemCode(ctx context.Context, userId int64, username string, code *entity.UserCode) error
	CodeLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserCodeLog, int, error)
	AddShareLog(ctx context.Context, userId int64, typ string, targetId int64, channel string) error
	ShareLogList(ctx context.Context, userId int64, page, size int) ([]*entity.UserShareLog, int, error)
	// 成长: 签到 / 任务
	GetSignDays(ctx context.Context, userId int64, yearMonth int) ([]int, bool, error)
	SaveSign(ctx context.Context, userId int64, yearMonth int, days []int, exists bool, credit float64) error
	ListTasks(ctx context.Context) ([]*entity.UserTask, error)
	FindTask(ctx context.Context, taskId int64) (*entity.UserTask, error)
	TaskDoneToday(ctx context.Context, userId, taskId int64, logDate int) (int, error)
	AddTaskLog(ctx context.Context, userId, taskId int64, typ string, logDate int, credit float64) error
	TaskLogs(ctx context.Context, userId int64, page, size int) ([]*entity.UserTaskLog, int, error)
	// 资产: 充值 / VIP / 兑换
	ListRechargePackages(ctx context.Context) ([]*entity.RechargePackage, error)
	FindRechargePackage(ctx context.Context, id int64) (*entity.RechargePackage, error)
	CreateRechargeOrder(ctx context.Context, orderNo string, userId, packageId int64, amount, coin float64) error
	ListVipPackages(ctx context.Context) ([]*entity.VipPackage, error)
	FindVipPackage(ctx context.Context, id int64) (*entity.VipPackage, error)
	OpenVip(ctx context.Context, userId int64, pkg *entity.VipPackage, startAt, endAt int64) error
	VipLogs(ctx context.Context, userId int64, page, size int) ([]*entity.VipLog, int, error)
	ExchangeCreditToCoin(ctx context.Context, userId int64, creditCost, coinGain float64) error
	// 私信
	SendMessage(ctx context.Context, fromId, toId int64, content string) error
	ListConversations(ctx context.Context, userId int64, page, size int) ([]*entity.ChatConversation, int, error)
	Messages(ctx context.Context, meId, peerId int64, page, size int) ([]*entity.ChatMessage, int, error)
	MarkRead(ctx context.Context, userId, peerId int64) error
	DeleteConversation(ctx context.Context, userId, peerId int64) error
}
