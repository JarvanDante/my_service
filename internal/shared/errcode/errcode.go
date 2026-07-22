// Package errcode 统一业务错误码。
package errcode

import "github.com/gogf/gf/v2/errors/gcode"

var (
	NotFound  = gcode.New(40400, "resource not found", nil)
	Forbidden = gcode.New(40300, "forbidden", nil)
	Internal  = gcode.New(50000, "internal error", nil)
)
