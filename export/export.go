package export

import (
	"github.com/iglev/ilua/log"
	"github.com/iglev/ilua/luar"
	glua "github.com/yuin/gopher-lua"
)

// OpenLib export lib to lua
func OpenLib(L *glua.LState, modName string, args map[string]interface{}) {
	mod, ok := L.RegisterModule(modName, map[string]glua.LGFunction{}).(*glua.LTable)
	if !ok {
		log.Error("register %v module fail", modName)
		return
	}
	for k, v := range args {
		mod.RawSetString(k, luar.New(L, v))
	}
}

// NewType new type for lua
func NewType(L *glua.LState, typename string, instance interface{}) {
	L.SetGlobal(typename, luar.NewType(L, instance))
}
