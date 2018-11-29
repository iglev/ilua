package ilua

import (
	glua "github.com/yuin/gopher-lua"
)

// Options ilua options
type Options struct {
}

// NewState new lua state
func NewState(opts ...Options) *glua.LState {
	return glua.NewState()
}

// Close close lua state
func Close(L *glua.LState) {
	L.Close()
}
