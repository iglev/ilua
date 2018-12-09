package ilua

import (
	"context"

	"github.com/iglev/ilua/export"
	glua "github.com/yuin/gopher-lua"
)

// LState lua state interface
type LState struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	gl             *glua.LState
	opts           *Options
	hfMgr          hotfixMgr
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
	L.cancelFunc()
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
	return doProFiles(L, argsFile)
}

// RegHotfix lua file register hotfix
func (L *LState) RegHotfix(file string) {
	if L.hfMgr != nil {
		L.hfMgr.reg(file)
	}
}

// Call golang call lua function
func (L *LState) Call(funcname string, args ...interface{}) (glua.LValue, error) {
	if L.hfMgr != nil {
		L.hfMgr.check(L)
	}
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
	return NewStateWithOpts(do)
}

// NewStateWithOpts new lua state with options
func NewStateWithOpts(do *Options) *LState {
	gluaOpts := glua.Options{
		CallStackSize: do.CallStackSize,
		RegistrySize:  do.RegistrySize,
	}
	ctx, cancel := context.WithCancel(context.Background())
	L := &LState{
		ctx:        ctx,
		cancelFunc: cancel,
		gl:         glua.NewState(gluaOpts),
		opts:       do,
	}
	L.openlibs()
	// hotfix
	if L.opts.NeedHotfix {
		L.hfMgr = newHotfixMgr(ctx, L.opts.NeedHotfixCoro)
	}
	return L
}

func (L *LState) openlibs() {
	// export log
	export.OpenLogLib(L.L())
}

////////////////////////////////////////////////////////////

// Options ilua options
type Options struct {
	HotfixTime     int64 // -1: no need hotfix, >=0: time for check
	CallStackSize  int   // Call stack size
	RegistrySize   int   // Data stack size
	GoCoroutine    bool  // need coroutine
	NeedHotfix     bool  // need hotfix
	NeedHotfixCoro bool  // need hotfix with coroutine
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

// SetHotfix set hotfix
func SetHotfix(needHotfixCoro bool) Option {
	return Option{func(do *Options) {
		do.NeedHotfix = true
		do.NeedHotfixCoro = needHotfixCoro
	}}
}
