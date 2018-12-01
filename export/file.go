package export

import (
	"os"

	glua "github.com/yuin/gopher-lua"
)

const (
	// FileLibName file module name
	FileLibName = "MFile"
)

// OpenFileLib export file lib
func OpenFileLib(L *glua.LState) {
	OpenLib(L, FileLibName, map[string]glua.LGFunction{
		"ModTime": exportFileGetModTime,
	})
}

func exportFileGetModTime(L *glua.LState) int {
	path, ok := L.Get(-1).(glua.LString)
	if !ok {
		L.Push(glua.LNil)
		L.Push(glua.LString("param error"))
		return 2
	}
	stat, err := os.Stat(string(path))
	if err != nil {
		L.Push(glua.LNil)
		L.Push(glua.LString(err.Error()))
		return 2
	}
	L.Push(glua.LNumber(stat.ModTime().Unix()))
	L.Push(glua.LNil)
	return 2
}
