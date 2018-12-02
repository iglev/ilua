package ilua

import (
	"github.com/iglev/ilua/luar"
	glua "github.com/yuin/gopher-lua"
)

func call(L *LState, funcname string, args ...interface{}) (interface{}, error) {
	// largs := newArgs(L.L(), args...)
	// TODO
	// ...
	return nil, nil
}

func newArgs(L *glua.LState, args ...interface{}) []glua.LValue {
	size := len(args)
	res := make([]glua.LValue, 0, size)
	for i := 0; i < size; i++ {
		res = append(res, luar.New(L, args[i]))
	}
	return res
}
