package ilua

import (
	"github.com/iglev/ilua/export"
	glua "github.com/yuin/gopher-lua"
)

// LState lua state interface
type LState struct {
	gl             *glua.LState
	opts           *Options
	lastHotfixTime int64
}

// LStateMod lua state module type
type LStateMod map[string]interface{}

// L get glua state
func (L *LState) L() *glua.LState {
	return L.gl
}

// Close close lua state
func (L *LState) Close() {
	L.gl.Close()
}

// RegMod register module
func (L *LState) RegMod(modName string, args LStateMod) {
	export.OpenLib(L.L(), modName, args)
}

// RegType register golang type for lua
func (L *LState) RegType(typename string, ins interface{}) {
	export.NewType(L.L(), typename, ins)
}

// DoProFiles do lua files
func (L *LState) DoProFiles(argsFile string) error {
	return doProFiles(L.L(), argsFile)
}

// CheckHotfix check hot fix
func (L *LState) CheckHotfix() error {
	return checkHotfix(L)
}

// Call golang call lua function
func (L *LState) Call(funcname string, args ...interface{}) (glua.LValue, error) {
	return call(L.L(), funcname, args...)
}

// UnmarshalLTB unmarshal lua table to golang struct
func (L *LState) UnmarshalLTB(script string, val interface{}) (interface{}, error) {
	return unmarshalLTB(L.L(), script, val)
}

// UnmarshalLTB unmarshal lua table to golang struct
func UnmarshalLTB(script string, val interface{}) (interface{}, error) {
	gluaOpts := glua.Options{
		CallStackSize: 64,
		RegistrySize:  64,
	}
	L := glua.NewState(gluaOpts)
	defer L.Close()
	return unmarshalLTB(L, script, val)
}

// NewState new lua state
func NewState(opts ...Option) *LState {
	do := &Options{
		HotfixTime:    DefaultHotfix,
		CallStackSize: DefaultCallStackSize,
		RegistrySize:  DefaultRegistrySize,
	}
	for _, option := range opts {
		option.f(do)
	}
	gluaOpts := glua.Options{
		CallStackSize: do.CallStackSize,
		RegistrySize:  do.RegistrySize,
	}
	L := &LState{
		gl:   glua.NewState(gluaOpts),
		opts: do,
	}
	L.openlibs()
	return L
}

func (L *LState) openlibs() {
	// export log
	export.OpenLogLib(L.L())
	// export file
	export.OpenFileLib(L.L())
}

////////////////////////////////////////////////////////////

// Options ilua options
type Options struct {
	HotfixTime    int64 // -1: no need hotfix, >=0: time for check
	CallStackSize int   // Call stack size
	RegistrySize  int   // Data stack size
}

// Option ilua option
type Option struct {
	f func(*Options)
}

// SetHotfixTime set hotfix time
func SetHotfixTime(dur int64) Option {
	return Option{func(do *Options) {
		do.HotfixTime = dur
	}}
}

// SetCallStackSize set call stack size
func SetCallStackSize(size int) Option {
	return Option{func(do *Options) {
		do.CallStackSize = size
	}}
}

// SetRegistrySize set data stack size
func SetRegistrySize(size int) Option {
	return Option{func(do *Options) {
		do.RegistrySize = size
	}}
}
