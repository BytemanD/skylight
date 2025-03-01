package utility

import (
	"io"
	"os"

	"github.com/gogf/gf/v2/os/gfile"
)

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	return err
}
func MakesureDir(path string) {
	if gfile.Exists(path) {
		return
	}
	if err := gfile.Mkdir(path); err != nil {
		panic(err)
	}
}
