package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/JarvanDante/my_service/internal/cmd"
)

func main() {
	cmd.ManageAPI.Run(gctx.GetInitCtx())
}
