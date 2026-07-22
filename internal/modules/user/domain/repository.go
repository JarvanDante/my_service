// Package domain 用户领域层。
// 说明: 为顺着 GoFrame 代码生成, 仓储直接以生成的 entity.Users 作为数据载体,
// 因此本包会 import model/entity(含 gtime)。若要严格「领域零框架」, 可再引入独立
// 领域实体 + 映射层, 这里按务实优先。
package domain

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/JarvanDante/my_service/internal/model/entity"
)

type Repository interface {
	// FindById 按主键查, 未找到返回 (nil, nil)。
	FindById(ctx context.Context, id int64) (*entity.Users, error)
	// FindByDeviceId 按设备号查(游客登录), 未找到返回 (nil, nil)。
	FindByDeviceId(ctx context.Context, deviceId string) (*entity.Users, error)
	// Create 新建用户, 返回自增 id。
	Create(ctx context.Context, data g.Map) (int64, error)
	// UpdateLoginInfo 登录成功后更新登录次数/时间/IP。
	UpdateLoginInfo(ctx context.Context, id int64, ip string) error
	// Disable 禁用用户。
	Disable(ctx context.Context, id int64, reason string) error
}
