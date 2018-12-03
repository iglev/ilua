package ilua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iglev/ilua/luar"
	glua "github.com/yuin/gopher-lua"
)

func call(L *glua.LState, funcname string, args ...interface{}) (glua.LValue, error) {
	functb := L.NewTable()
	sp := strings.Split(funcname, ".")
	size := len(sp)
	if size == 1 {
		functb.RawSetString("f", glua.LString(sp[0]))
	} else if size == 2 {
		functb.RawSetString("m", glua.LString(sp[0]))
		functb.RawSetString("f", glua.LString(sp[1]))
	} else {
		return nil, fmt.Errorf("invalid funcname=%v", funcname)
	}
	largs := newArgs(L, []glua.LValue{functb}, args...)
	err := L.CallByParam(glua.P{
		Fn:      L.GetGlobal(LuaFuncCall),
		NRet:    1,
		Protect: true,
	}, largs...)
	if err != nil {
		return nil, err
	}
	restb, ok := L.Get(-1).(*glua.LTable)
	defer L.Pop(1)
	if !ok {
		return nil, fmt.Errorf("Call lua `func=LFGCall` must return {r=xxx, err=xxx}")
	}
	errStr, errOK := restb.RawGetString("err").(glua.LString)
	if errOK {
		return nil, errors.New(string(errStr))
	}
	return restb.RawGetString("r"), nil
}

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
