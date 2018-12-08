package ilua

import (
	"os"

	"github.com/iglev/ilua/log"
)

func getFileModtime(file string) (int64, error) {
	stat, err := os.Stat(file)
	if err != nil {
		log.Error("stat fail, file=%v err=%v", file, err)
		return 0, err
	}
	return stat.ModTime().Unix(), nil
}
