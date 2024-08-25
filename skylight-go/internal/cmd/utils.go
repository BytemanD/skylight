package cmd

import (
	"github.com/gogf/gf/v2/os/gfile"
)

func MakesureDir(path string) {
	if gfile.Exists(path) {
		return
	}
	if err := gfile.Mkdir(path); err != nil {
		panic(err)
	}
}
