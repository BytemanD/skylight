package main

import (
	_ "skylight/internal/packed"

	"skylight/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	cmd.Main.Run(gctx.New())
}
