package export

import (
	log "github.com/iglev/ilog"
	"github.com/iglev/ilua/luar"
	glua "github.com/yuin/gopher-lua"
)

// OpenLib export lib to lua
func OpenLib(L *glua.LState, modName string, args map[string]interface{}) *glua.LTable {
	mod, ok := L.RegisterModule(modName, map[string]glua.LGFunction{}).(*glua.LTable)
	if !ok {
		log.Error("register %v module fail", modName)
		return nil
	}
	for k, v := range args {
		mod.RawSetString(k, luar.New(L, v))
	}
	return mod
}

// NewType new type for global
func NewType(L *glua.LState, typename string, instance interface{}) {
	L.SetGlobal(typename, luar.NewType(L, instance))
}

// NewModType new type for module
func NewModType(L *glua.LState, mod *glua.LTable, typename string, instance interface{}) {
	mod.RawSetString(typename, luar.NewType(L, instance))
}
