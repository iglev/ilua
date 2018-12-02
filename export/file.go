package export

import (
	"os"

	glua "github.com/yuin/gopher-lua"
)

const (
	// FileLibName file module name
	FileLibName = "mfile"
)

// OpenFileLib export file lib
func OpenFileLib(L *glua.LState) {
	OpenLib(L, FileLibName, map[string]interface{}{
		"ModTime": exportFileGetModTime,
	})
}

func exportFileGetModTime(path string) (glua.LNumber, glua.LValue) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, glua.LString(err.Error())
	}
	return glua.LNumber(stat.ModTime().Unix()), glua.LNil
}
