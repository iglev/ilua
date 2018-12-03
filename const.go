package ilua

///////////////////////////////////////////////////////////////////

// Lua State Options

const (
	// DefaultHotfix default hotfix time
	DefaultHotfix = 1 // 10 // 10 second
	// DefaultCallStackSize default call stack size
	DefaultCallStackSize = 256
	// DefaultRegistrySize default data stack size
	DefaultRegistrySize = 256 * 20
)

///////////////////////////////////////////////////////////////////

// LoadLibs

const (
	// LuaBaseMainFileName basemain.lua
	LuaBaseMainFileName = "G_BASE_MAIN_FILE"
	// LuaBaseMainModule basemain module name
	LuaBaseMainModule = "G_BASE_MAIN"
	// LuaMainFileName main.lua
	LuaMainFileName = "G_MAIN_FILE"
	// LuaMainModule main module name
	LuaMainModule = "G_MAIN"
	// LuaValLuaFiles obj luafiles
	LuaValLuaFiles = "luafiles"
	// LuaValDir obj dir
	LuaValDir = "dir"
	// LuaValFiles obj files
	LuaValFiles = "files"
	// LuaValDecode var decode
	LuaValDecode = "decode"
)

///////////////////////////////////////////////////////////////////

// lua for golang global function

const (
	// LuaFuncHotfix lua func "LFGHotFix"
	LuaFuncHotfix = "LFGHotFix"
	// LuaFuncCall lua func "LFGCall"
	LuaFuncCall = "LFGCall"
)

///////////////////////////////////////////////////////////////////
