package main

import (
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
	cmd.Version = Version
	cmd.GoVersion = GoVersion
	cmd.BuildDate = BuildDate
	cmd.BuildPlatform = BuildPlatform

	cmd.Main.Run(gctx.New())
}
