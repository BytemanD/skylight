package main

import (
	"skylight/internal/consts"
	_ "skylight/internal/packed"

	"skylight/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

var (
	Version       string
	GoVersion     string
	BuildDate     string
	BuildPlatform string
)

func main() {
	consts.Version = Version
	consts.GoVersion = GoVersion
	consts.BuildDate = BuildDate
	consts.BuildPlatform = BuildPlatform
	if consts.Version == "" {
		consts.Version = "dev"
	}

	cmd.Main.Run(gctx.New())
}
