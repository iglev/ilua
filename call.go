package ilua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iglev/ilua/luar"
	glua "github.com/yuin/gopher-lua"
)

func newArgs(L *glua.LState, preargs []glua.LValue, args ...interface{}) []glua.LValue {
	presize := len(preargs)
	argsize := len(args)
	size := argsize + presize
	res := make([]glua.LValue, 0, size)
	for i := 0; i < presize; i++ {
		res = append(res, preargs[i])
	}
	for i := 0; i < argsize; i++ {
		res = append(res, luar.New(L, args[i]))
	}
	return res
}

func call(L *glua.LState, funcname string, args ...interface{}) (glua.LValue, error) {
	sp := strings.Split(funcname, ".")
	size := len(sp)
	var funcval glua.LValue
	if size == 1 {
		if fv := L.GetGlobal(sp[0]); fv != glua.LNil {
			funcval = fv
		}
	} else if size == 2 {
		mod, ok := L.GetGlobal(sp[0]).(*glua.LTable)
		if !ok {
			return nil, fmt.Errorf("module name invalid, name=%v", funcname)
		}
		if fv := mod.RawGetString(sp[1]); fv != glua.LNil {
			funcval = fv
		}
	} else {
		return nil, fmt.Errorf("function name invalid, name=%v", funcname)
	}
	if funcval == nil {
		return nil, fmt.Errorf("not found func, name=%v", funcname)
	}
	largs := newArgs(L, nil, args...)
	err := L.CallByParam(glua.P{
		Fn:      funcval,
		NRet:    2,
		Protect: true,
	}, largs...)
	if err != nil {
		return nil, err
	}
	res := L.Get(-2)
	resErr := L.Get(-1)
	defer L.Pop(2)
	switch resErr.(type) {
	case glua.LString:
		return nil, errors.New(string(resErr.(glua.LString)))
	case *glua.LUserData:
		return nil, fmt.Errorf("%v", (resErr.(*glua.LUserData)).Value)
	}
	return res, nil
}
