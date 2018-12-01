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

// L get glua state
func (L *LState) L() *glua.LState {
	return L.gl
}

// CheckHotfix check hot fix
func (L *LState) CheckHotfix() error {
	return checkHotfix(L)
}

// Close close lua state
func (L *LState) Close() {
	L.gl.Close()
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
