package export

import (
	"os"

	glua "github.com/yuin/gopher-lua"
)

var fileFuncs = map[string]glua.LGFunction{
	"ModTime": exportFileGetModTime,
}

// OpenFileLib export file lib
func OpenFileLib(L *glua.LState, modName string) {
	L.Push(L.NewFunction(exportFileLoader))
	L.Push(glua.LString(modName))
	L.Call(1, 0)
}

func exportFileLoader(L *glua.LState) int {
	modName, ok := L.Get(-1).(glua.LString)
	if !ok {
		return 0
	}
	mod := L.RegisterModule(string(modName), fileFuncs).(*glua.LTable)
	L.Push(mod)
	return 1
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
