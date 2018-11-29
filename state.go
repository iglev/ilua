package ilua

import (
	glua "github.com/yuin/gopher-lua"
)

// Options ilua options
type Options struct {
}

// NewState new lua state
func NewState(opts ...Options) *glua.LState {
	L := glua.NewState()
	for _, pair := range []struct {
		n string
		f glua.LGFunction
	}{
		{glua.LoadLibName, glua.OpenPackage},
		{glua.BaseLibName, glua.OpenBase},
		{glua.TabLibName, glua.OpenTable},
	} {
		if err := L.CallByParam(glua.P{
			Fn:      L.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, glua.LString(pair.n)); err != nil {
			logerror("open package fail, pack=%v err=%v", pair.n, err)
			return nil
		}
	}
	return L
}

// Close close lua state
func Close(L *glua.LState) {
	L.Close()
}
