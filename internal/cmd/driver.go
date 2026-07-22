package cmd

// 注册数据库驱动。
// GoFrame v2 将各数据库驱动拆为独立 contrib 模块, 必须显式(空)导入才会注册,
// 否则运行时报 "cannot find database driver for specified database type pgsql"。
// 所有 app/* 二进制都会经 internal/cmd 加载此文件, 故驱动只需在此引一次。
import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	// 如需其他库再按需加, 例如:
	// _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)
