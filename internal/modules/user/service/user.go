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
	Id          int64
	Username    string
	Nickname    string
	Phone       string
	Channel     string
	GroupId     int64
	GroupName   string
	Level       int
	Balance     float64
	Credit      float64
	MoneyCount  float64
	IsDisabled  int
	RegisterAt  string
	LastLoginAt string
}

// AdminUserListInput 后台用户列表入参。
type AdminUserListInput struct {
	Keyword   string
	Channel   string
	GroupId   int64
	Status    int // 0全部 1正常 2禁用
	StartDate int // YYYYMMDD
	EndDate   int
	Page      int
	Size      int
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
	Sex          int
	Signature    string
	Img          string
	Fans         int
	Follow       int
	ShareNum     int
	ParentId     int64
	ParentName   string
	GroupRate    int
	GroupEndTime int64
	ErrorMsg     string
	RegisterIp   string
	LastIp       string
	LoginNum     int
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
}
