package cmd

// 注册数据库/缓存适配器。GoFrame v2 将驱动/适配器拆为独立 contrib 模块, 必须显式空导入。
import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2" // PostgreSQL 驱动
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"   // Redis 适配器(g.Redis 需要)
)
