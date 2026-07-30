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
	// 后台
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
	DisableUser(ctx context.Context, id int64) error
}
