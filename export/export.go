package export

import (
	glua "github.com/yuin/gopher-lua"
)

// OpenLib export lib to lua
func OpenLib(L *glua.LState, modName string, funcs map[string]glua.LGFunction) {
	L.Push(L.NewFunction(func(L *glua.LState) int {
		mod := L.RegisterModule(modName, funcs).(*glua.LTable)
		L.Push(mod)
		return 1
	}))
	L.Call(0, 0)
}
