package main

import (
	_ "realtime_quiz_system/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"realtime_quiz_system/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
