// Package service 用户模块对外能力接口。
package service

import "context"

type LoginInput struct {
	DeviceId      string
	DeviceType    string
	DeviceVersion string
	Ip            string
}

// UserInfoDTO 当前用户详情(含私密字段)。
type UserInfoDTO struct {
	Id        int64
	Username  string
	Nickname  string
	Phone     string
	Img       string
	Signature string
	Sex       int
	Level     int
	Balance   float64
	Credit    float64
	GroupName string
	Fans      int
	Follow    int
	// 以下供站点差异字段(ext)候选, 是否返回由 Nacos response.user_info_extra 决定
	BgImg        string
	ShareNum     int
	ChannelName  string
	GroupEndTime int64
}

// PublicUserDTO 对外公开信息(看他人时用, 不含手机/余额)。
type PublicUserDTO struct {
	Id        int64
	Nickname  string
	Img       string
	BgImg     string
	Signature string
	Sex       int
	Level     int
	Fans      int
	Follow    int
	ShareNum  int
}

// HomeDTO 他人主页。
type HomeDTO struct {
	User       *PublicUserDTO
	IsFollowed bool
}

// UpdateProfileInput 改资料入参(空值表示不改)。
type UpdateProfileInput struct {
	Nickname  string
	Img       string
	BgImg     string
	Signature string
	Sex       int // 1男 2女, 0 表示不改
}

type LoginDTO struct {
	Token string
	User  *UserInfoDTO
}

// ---- P4 推广/兑换码 ----

type RedeemDTO struct {
	Type   string // point / group
	AddNum int
	Name   string
}

type CodeLogDTO struct {
	Id        int64
	Code      string
	Name      string
	Type      string
	AddNum    int
	CreatedAt string
}

type ShareDTO struct {
	ShareCode string
	ShareUrl  string
	ShareNum  int
}

type ShareLogDTO struct {
	Id        int64
	Type      string
	TargetId  int64
	Channel   string
	CreatedAt string
}

// ---- P5 成长(签到/任务) ----

type SignDTO struct {
	Today  int
	Days   []int
	Count  int
	Reward float64
}

type TaskDTO struct {
	Id          int64
	Name        string
	Type        string
	Description string
	MaxNum      int
	DoneToday   int
}

type TaskDoneDTO struct {
	DoneToday int
	MaxNum    int
	Reward    float64
}

type TaskLogDTO struct {
	Id        int64
	TaskId    int64
	Type      string
	Num       int
	LogDate   int
	CreatedAt string
}

type UpDTO struct {
	Level             int
	Credit            float64
	Balance           float64
	SignDaysThisMonth int
	Fans              int
	Follow            int
	ShareNum          int
}

// ---- P6 资产(充值/VIP/兑换) ----

type RechargePackageDTO struct {
	Id     int64
	Name   string
	Amount float64
	Coin   float64
	Bonus  float64
}

type RechargeOrderDTO struct {
	OrderNo string
	Amount  float64
	Coin    float64
}

type VipPackageDTO struct {
	Id    int64
	Name  string
	Days  int
	Price float64
}

type VipLogDTO struct {
	Id        int64
	PackageId int64
	Days      int
	Price     float64
	StartAt   int64
	EndAt     int64
	CreatedAt string
}

type ExchangeInfoDTO struct {
	Rate    int
	Credit  float64
	Balance float64
}

// ---- P7 私信 ----

type ChatDTO struct {
	PeerId      int64
	Nickname    string
	Img         string
	LastContent string
	LastAt      string
	Unread      int
}

type MessageDTO struct {
	Id        int64
	FromId    int64
	ToId      int64
	Content   string
	Mine      bool
	CreatedAt string
}

// ---- B1 后台用户管理 ----

// AdminUserItemDTO 后台用户列表项。
type AdminUserItemDTO struct {
	Id             int64
	Username       string
	Nickname       string
	Phone          string
	Sex            int
	Tag            string
	Img            string
	AccountSlat    string
	Balance        float64
	GiftCount      float64
	Credit         float64
	MoneyCount     float64
	IsUp           int
	IsValid        int
	HasBuy         int
	Level          int
	GroupId        int64
	GroupName      string
	GroupRate      int
	GroupStartTime int64
	GroupEndTime   int64
	ParentId       int64
	ParentName     string
	Channel        string
	DeviceType     string
	DeviceExt      string
	DeviceVersion  string
	MovieFeeRate   int
	PostFeeRate    int
	ShareNum       int
	IsDisabled     int
	RegisterAt     string
	RegisterIp     string
	RegisterArea   string
	LastLoginAt    string
	LastIp         string
	LoginNum       int
}

// AdminUserListInput 后台用户列表入参。
type AdminUserListInput struct {
	Keyword     string
	UserId      int64
	Username    string
	Phone       string
	ParentId    int64
	Channel     string
	GroupId     int64
	IsUp        int // 0全部 1是 2否
	IsValid     int
	HasBuy      int
	Status      int // 0全部 1正常 2禁用
	DeviceType  string
	StartDate   int // YYYYMMDD
	EndDate     int
	MinLoginNum int
	MaxLoginNum int
	Page        int
	Size        int
}

// AdminUserListDTO 后台用户列表。
type AdminUserListDTO struct {
	List  []*AdminUserItemDTO
	Total int
	Page  int
	Size  int
}

// AdminUserDetailDTO 后台用户详情(含状态/组/资产/轨迹)。
type AdminUserDetailDTO struct {
	AdminUserItemDTO
	Signature string
	Fans      int
	Follow    int
	ErrorMsg  string
}

// AdminSetGroupInput 调整用户组入参。
type AdminSetGroupInput struct {
	UserId       int64
	GroupId      int64
	GroupName    string
	GroupRate    int
	GroupEndTime int64 // epoch 秒, 0 表示不限
}

// AdminAdjustBalanceInput 调整余额入参。
type AdminAdjustBalanceInput struct {
	UserId     int64
	Target     string  // balance / credit
	Amount     float64 // 正加负减, 不为 0
	Remark     string
	OperatorId int64 // 操作管理员
}

// BalanceLogDTO 余额流水。
type BalanceLogDTO struct {
	Id            int64
	Direction     int
	Scene         string
	Amount        float64
	BalanceBefore float64
	BalanceAfter  float64
	RefId         string
	Remark        string
	CreatedAt     string
}

// ---- B4 用户组定义 ----

type UserGroupDTO struct {
	Id     int64
	Name   string
	Rate   int
	Rights string
	Remark string
	Sort   int
	Status int
}

type UserGroupInput struct {
	Id     int64 // 更新时用
	Name   string
	Rate   int
	Rights string
	Remark string
	Sort   int
	Status int
}

// ---- B5 成长配置 ----

type AdminTaskDTO struct {
	Id          int64
	Name        string
	Type        string
	Description string
	MaxNum      int
	Reward      float64
	Status      int
	Sort        int
}

type AdminTaskInput struct {
	Id          int64 // 更新时用
	Name        string
	Type        string
	Description string
	MaxNum      int
	Reward      float64
	Status      int
	Sort        int
}

type AdminTaskLogDTO struct {
	Id        int64
	UserId    int64
	TaskId    int64
	Type      string
	Num       int
	LogDate   int
	CreatedAt string
}

type AdminTaskLogInput struct {
	UserId    int64
	TaskId    int64
	Type      string
	StartDate int
	EndDate   int
	Page      int
	Size      int
}

type SignDayCountDTO struct {
	Day   int
	Count int
}

type SignStatsDTO struct {
	YearMonth int
	UserCount int // 签到用户数
	SignCount int // 总签到人次
	Days      []SignDayCountDTO
}

// ---- B6 社交查询 ----

type FollowAdminDTO struct {
	Id        int64
	UserId    int64
	UserName  string
	HomeId    int64
	HomeName  string
	CreatedAt string
}

// ---- B7 消息监控 ----

type AdminMessageDTO struct {
	Id        int64
	FromId    int64
	ToId      int64
	Content   string
	CreatedAt string
}

type AdminMessageInput struct {
	FromId    int64
	ToId      int64
	UserId    int64
	Keyword   string
	StartDate string
	EndDate   string
	Page      int
	Size      int
}

type UserDTO struct {
	ID     int64
	Name   string
	Active bool
}

type IUser interface {
	// 认证
	Login(ctx context.Context, in LoginInput) (*LoginDTO, error)
	Logout(ctx context.Context, userId int64) error
	Refresh(ctx context.Context, userId int64) (string, error)
	// 资料
	Info(ctx context.Context, userId int64) (*UserInfoDTO, error)
	Home(ctx context.Context, viewerId, homeId int64) (*HomeDTO, error)
	UpdateProfile(ctx context.Context, userId int64, in UpdateProfileInput) error
	Images(ctx context.Context, userId int64) ([]string, error)
	FindByAccount(ctx context.Context, account string) (*PublicUserDTO, error)
	BindPhone(ctx context.Context, userId int64, phone, code string) error
	// 社交
	DoFollow(ctx context.Context, userId, homeId int64) (bool, error)
	Following(ctx context.Context, userId int64, page, size int) ([]*PublicUserDTO, int, error)
	Fans(ctx context.Context, userId int64, page, size int) ([]*PublicUserDTO, int, error)
	// 推广 / 兑换码
	BindParent(ctx context.Context, userId int64, account string) error
	RedeemCode(ctx context.Context, userId int64, code string) (*RedeemDTO, error)
	CodeLogs(ctx context.Context, userId int64, page, size int) ([]*CodeLogDTO, int, error)
	ShareInfo(ctx context.Context, userId int64) (*ShareDTO, error)
	ShareLogs(ctx context.Context, userId int64, page, size int) ([]*ShareLogDTO, int, error)
	ReportShare(ctx context.Context, userId int64, typ string, targetId int64, channel string) error
	// 成长(签到/任务)
	DoDaySign(ctx context.Context, userId int64) (*SignDTO, error)
	Tasks(ctx context.Context, userId int64) ([]*TaskDTO, error)
	DoTask(ctx context.Context, userId, taskId int64) (*TaskDoneDTO, error)
	TaskLogs(ctx context.Context, userId int64, page, size int) ([]*TaskLogDTO, int, error)
	Up(ctx context.Context, userId int64) (*UpDTO, error)
	// 资产(充值/VIP/兑换)
	RechargePackages(ctx context.Context) ([]*RechargePackageDTO, error)
	DoRecharge(ctx context.Context, userId, packageId int64) (*RechargeOrderDTO, error)
	VipPackages(ctx context.Context) ([]*VipPackageDTO, error)
	DoVip(ctx context.Context, userId, packageId int64) error
	VipLogs(ctx context.Context, userId int64, page, size int) ([]*VipLogDTO, int, error)
	ExchangeInfo(ctx context.Context, userId int64) (*ExchangeInfoDTO, error)
	DoExchange(ctx context.Context, userId int64, coin int) error
	// 私信
	SendMessage(ctx context.Context, meId, toId int64, content string) error
	Chats(ctx context.Context, meId int64, page, size int) ([]*ChatDTO, int, error)
	Messages(ctx context.Context, meId, peerId int64, page, size int) ([]*MessageDTO, int, error)
	DelChat(ctx context.Context, meId, peerId int64) error
	CustomerUrl(ctx context.Context) (string, error)
	// 后台(B1 用户管理)
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
	DisableUser(ctx context.Context, id int64) error
	AdminListUsers(ctx context.Context, in AdminUserListInput) (*AdminUserListDTO, error)
	AdminUserDetail(ctx context.Context, id int64) (*AdminUserDetailDTO, error)
	AdminSetDisabled(ctx context.Context, id int64, disable bool, reason string) error
	AdminSetGroup(ctx context.Context, in AdminSetGroupInput) error
	AdminAdjustBalance(ctx context.Context, in AdminAdjustBalanceInput) error
	AdminBalanceLogs(ctx context.Context, userId int64, page, size int) ([]*BalanceLogDTO, int, error)
	// 后台(B4 用户组定义)
	AdminGroups(ctx context.Context) ([]*UserGroupDTO, error)
	AdminCreateGroup(ctx context.Context, in UserGroupInput) (int64, error)
	AdminUpdateGroup(ctx context.Context, in UserGroupInput) error
	AdminDeleteGroup(ctx context.Context, id int64) error
	// 后台(B5 成长配置)
	AdminTasks(ctx context.Context) ([]*AdminTaskDTO, error)
	AdminCreateTask(ctx context.Context, in AdminTaskInput) (int64, error)
	AdminUpdateTask(ctx context.Context, in AdminTaskInput) error
	AdminDeleteTask(ctx context.Context, id int64) error
	AdminTaskLogs(ctx context.Context, in AdminTaskLogInput) ([]*AdminTaskLogDTO, int, error)
	AdminSignStats(ctx context.Context, yearMonth int) (*SignStatsDTO, error)
	// 后台(B6 社交查询)
	AdminFollows(ctx context.Context, userId, homeId int64, page, size int) ([]*FollowAdminDTO, int, error)
	// 后台(B7 消息监控)
	AdminMessages(ctx context.Context, in AdminMessageInput) ([]*AdminMessageDTO, int, error)
}
